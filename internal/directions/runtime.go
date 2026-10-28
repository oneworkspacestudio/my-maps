package directions

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type engineProcess struct {
	cmd        *exec.Cmd
	url, token string
	done       chan struct{}
	job        io.Closer
}

func (e *engineProcess) stop() {
	if e == nil {
		return
	}
	if e.job != nil {
		e.job.Close()
	}
	_ = e.cmd.Process.Kill()
	select {
	case <-e.done:
	case <-time.After(5 * time.Second):
	}
}
func (m *Manager) stopEngine() {
	m.engineMu.Lock()
	defer m.engineMu.Unlock()
	if m.engine != nil {
		m.engine.stop()
		m.engine = nil
	}
}
func (m *Manager) launch(ctx context.Context, root string) (*engineProcess, error) {
	for _, name := range []string{"runtime/bin/java.exe", "brouter.jar", "profiles/lookups.dat"} {
		st, err := os.Lstat(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil || !st.Mode().IsRegular() {
			return nil, errors.New("local directions runtime needs a fresh download")
		}
	}
	path := filepath.Join(root, "MapsRoutingBridge.class")
	if st, err := os.Lstat(path); err == nil && !st.Mode().IsRegular() {
		return nil, errors.New("invalid local directions bridge")
	}
	if err := os.WriteFile(path, bridge, 0600); err != nil {
		return nil, err
	}
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(m.ctx, filepath.Join(root, "runtime", "bin", "java.exe"), "-Xmx384m", "-Xms32m", "--add-modules", "jdk.httpserver", "-cp", root+string(os.PathListSeparator)+filepath.Join(root, "brouter.jar"), "MapsRoutingBridge", filepath.Join(root, "segments"), filepath.Join(root, "profiles"))
	cmd.Dir = root
	cmd.Stderr = io.Discard
	hideWindow(cmd)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, errors.New("Windows could not start the downloaded directions runtime")
	}
	e := &engineProcess{cmd: cmd, token: hex.EncodeToString(secret[:]), done: make(chan struct{})}
	e.job, err = ownProcess(cmd)
	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return nil, errors.New("directions process could not be safely supervised")
	}
	go func() { _ = cmd.Wait(); close(e.done) }()
	if _, err = io.WriteString(stdin, e.token+"\n"); err != nil {
		e.stop()
		return nil, err
	}
	stdin.Close()
	ready := make(chan string, 1)
	go func() {
		line, _ := bufio.NewReader(io.LimitReader(stdout, 128)).ReadString('\n')
		ready <- strings.TrimSpace(line)
	}()
	deadline := time.NewTimer(2 * time.Minute)
	defer deadline.Stop()
	select {
	case line := <-ready:
		if !strings.HasPrefix(line, "READY ") {
			e.stop()
			return nil, errors.New("directions data integrity or runtime startup check failed")
		}
		port, parseErr := strconv.Atoi(strings.TrimPrefix(line, "READY "))
		if parseErr != nil || port < 1 || port > 65535 {
			e.stop()
			return nil, errors.New("directions runtime returned an invalid endpoint")
		}
		e.url = fmt.Sprintf("http://127.0.0.1:%d", port)
		req, _ := http.NewRequestWithContext(ctx, "GET", e.url+"/healthz", nil)
		req.Header.Set("Authorization", "Bearer "+e.token)
		client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{Proxy: nil}, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("local directions cannot redirect") }}
		response, callErr := client.Do(req)
		if callErr != nil {
			e.stop()
			return nil, errors.New("directions runtime health check failed")
		}
		response.Body.Close()
		client.CloseIdleConnections()
		if response.StatusCode != 200 {
			e.stop()
			return nil, errors.New("directions runtime is not ready")
		}
		return e, nil
	case <-ctx.Done():
		e.stop()
		return nil, ctx.Err()
	case <-m.ctx.Done():
		e.stop()
		return nil, m.ctx.Err()
	case <-e.done:
		e.stop()
		return nil, errors.New("directions data or runtime could not be opened")
	case <-deadline.C:
		e.stop()
		return nil, errors.New("checking directions data took too long; retry or install a smaller area")
	}
}

func (m *Manager) engineFor(ctx context.Context) (*engineProcess, error) {
	m.engineMu.Lock()
	defer m.engineMu.Unlock()
	if m.engine != nil {
		select {
		case <-m.engine.done:
			m.engine.job.Close()
			m.engine = nil
		default:
			return m.engine, nil
		}
	}
	m.mu.Lock()
	info := m.info
	m.mu.Unlock()
	if info == nil {
		return nil, ErrNoData
	}
	e, err := m.launch(ctx, filepath.Join(m.root, info.ID))
	if err != nil {
		return nil, err
	}
	m.engine = e
	return e, nil
}
