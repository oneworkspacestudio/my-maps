package directions

import (
	"archive/zip"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

const limit = uint64(2 << 30)

var segmentPattern = regexp.MustCompile(`^[EW](0|[1-9][0-9]{0,2})_[NS](0|[1-9][0-9]?)\.rd5$`)
var idPattern = regexp.MustCompile(`^[a-f0-9]{32}$`)

type Choice struct {
	Coverage string  `json:"coverage"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Radius   float64 `json:"radius"`
}
type File struct {
	Name   string `json:"name"`
	Bytes  uint64 `json:"bytes"`
	SHA256 string `json:"sha256,omitempty"`
	ETag   string `json:"etag,omitempty"`
	URL    string `json:"-"`
}
type Info struct {
	ID        string     `json:"id"`
	Choice    Choice     `json:"choice"`
	Bounds    [4]float64 `json:"bounds"`
	Files     []File     `json:"files"`
	Bytes     uint64     `json:"bytes"`
	Installed string     `json:"installed"`
	Engine    string     `json:"engine"`
}
type Plan struct {
	ID            string     `json:"id"`
	Choice        Choice     `json:"choice"`
	Bounds        [4]float64 `json:"bounds"`
	Files         []File     `json:"files"`
	DownloadBytes uint64     `json:"downloadBytes"`
	RequiredBytes uint64     `json:"requiredBytes"`
	Expires       time.Time  `json:"expires"`
}
type Progress struct {
	Phase   string `json:"phase"`
	Message string `json:"message"`
	Done    uint64 `json:"done"`
	Total   uint64 `json:"total"`
	Busy    bool   `json:"busy"`
}
type Status struct {
	Installed *Info    `json:"installed"`
	Progress  Progress `json:"progress"`
	Ready     bool     `json:"ready"`
	Engine    string   `json:"engine"`
}

type Manager struct {
	root      string
	diskFree  func(string) (uint64, error)
	client    *http.Client
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.Mutex
	plan      *Plan
	info      *Info
	progress  Progress
	jobCancel context.CancelFunc
	wg        sync.WaitGroup
	routeGate chan struct{}
	engineMu  sync.Mutex
	engine    *engineProcess
}

func New(root string, diskFree func(string) (uint64, error)) (*Manager, error) {
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, err
	}
	if st, err := os.Lstat(root); err != nil || st.Mode()&os.ModeSymlink != 0 {
		return nil, errors.New("directions directory must be a real local folder")
	}
	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{root: root, diskFree: diskFree, ctx: ctx, cancel: cancel, routeGate: make(chan struct{}, 1), progress: Progress{Phase: "choose", Message: "Download local directions to drive, walk or cycle."}}
	m.client = &http.Client{Transport: &http.Transport{Proxy: nil, ResponseHeaderTimeout: 30 * time.Second, TLSHandshakeTimeout: 15 * time.Second}, CheckRedirect: func(r *http.Request, via []*http.Request) error {
		if len(via) > 4 || r.URL.Scheme != "https" || r.URL.User != nil {
			return errors.New("invalid download redirect")
		}
		host := r.URL.Hostname()
		first := via[0].URL.Hostname()
		if host != first && !(first == "github.com" && host == "release-assets.githubusercontent.com") {
			return errors.New("download left its approved source")
		}
		return nil
	}}
	if raw, err := os.ReadFile(filepath.Join(root, "active.json")); err == nil {
		var info Info
		if json.Unmarshal(raw, &info) == nil && validInfo(&info) && m.filesPresent(&info) {
			m.info = &info
			m.progress = Progress{Phase: "ready", Message: "Local directions are installed."}
		} else {
			m.progress = Progress{Phase: "error", Message: "Directions data needs a fresh download. Saved places are unaffected."}
		}
	} else if !os.IsNotExist(err) {
		cancel()
		return nil, err
	}
	entries, _ := os.ReadDir(root)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".partial") && idPattern.MatchString(strings.TrimSuffix(e.Name(), ".partial")) {
			_ = removeStage(filepath.Join(root, e.Name()))
		}
	}
	return m, nil
}

func validInfo(i *Info) bool {
	if !idPattern.MatchString(i.ID) || len(i.Files) < 3 || len(i.Files) > 26 || i.Bytes > limit+(512<<20) {
		return false
	}
	for _, f := range i.Files {
		if !segmentPattern.MatchString(f.Name) && f.Name != "brouter.zip" && f.Name != "runtime.zip" {
			return false
		}
		if len(f.SHA256) != 64 {
			return false
		}
	}
	return true
}
func (m *Manager) filesPresent(i *Info) bool {
	base := filepath.Join(m.root, i.ID)
	for _, f := range i.Files {
		if segmentPattern.MatchString(f.Name) {
			s, e := os.Lstat(filepath.Join(base, "segments", f.Name))
			if e != nil || !s.Mode().IsRegular() || uint64(s.Size()) != f.Bytes {
				return false
			}
		}
	}
	for _, name := range []string{"runtime/bin/java.exe", "brouter.jar", "profiles/lookups.dat"} {
		s, e := os.Lstat(filepath.Join(base, filepath.FromSlash(name)))
		if e != nil || !s.Mode().IsRegular() {
			return false
		}
	}
	return true
}
func (m *Manager) Status() Status {
	m.mu.Lock()
	defer m.mu.Unlock()
	return Status{Installed: m.info, Progress: m.progress, Ready: m.info != nil, Engine: Version}
}
func (m *Manager) Close() {
	m.cancel()
	m.mu.Lock()
	if m.jobCancel != nil {
		m.jobCancel()
	}
	m.mu.Unlock()
	m.stopEngine()
	m.wg.Wait()
	m.client.CloseIdleConnections()
}
func randomID() string {
	var b [16]byte
	if _, e := rand.Read(b[:]); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b[:])
}

func choiceBounds(c Choice) ([4]float64, error) {
	if c.Coverage == "uk" {
		return [4]float64{-8.7, 49.8, 2.1, 60.9}, nil
	}
	if c.Coverage != "area" || math.IsNaN(c.Lat) || math.IsNaN(c.Lon) || math.IsNaN(c.Radius) || math.IsInf(c.Lat, 0) || math.IsInf(c.Lon, 0) || math.IsInf(c.Radius, 0) || math.Abs(c.Lat) > 80 || math.Abs(c.Lon) > 180 || c.Radius < 1 || c.Radius > 500 {
		return [4]float64{}, errors.New("choose the United Kingdom or an area with a 1–500 km radius")
	}
	dy := c.Radius / 111.32
	dx := dy / math.Cos(c.Lat*math.Pi/180)
	if c.Lon-dx < -180 || c.Lon+dx >= 180 {
		return [4]float64{}, errors.New("choose an area away from the date line")
	}
	return [4]float64{c.Lon - dx, math.Max(-85, c.Lat-dy), c.Lon + dx, math.Min(85, c.Lat+dy)}, nil
}
func segmentName(lon, lat int) string {
	ew, ns := "E", "N"
	if lon < 0 {
		ew = "W"
		lon = -lon
	}
	if lat < 0 {
		ns = "S"
		lat = -lat
	}
	return fmt.Sprintf("%s%d_%s%d.rd5", ew, lon, ns, lat)
}
func segmentsFor(bounds [4]float64) ([]string, error) {
	var names []string
	for lon := int(math.Floor(bounds[0]/5)) * 5; float64(lon) <= bounds[2]; lon += 5 {
		for lat := int(math.Floor(bounds[1]/5)) * 5; float64(lat) <= bounds[3]; lat += 5 {
			names = append(names, segmentName(lon, lat))
			if len(names) > 24 {
				return nil, errors.New("choose a smaller directions area")
			}
		}
	}
	return names, nil
}
func (m *Manager) Plan(ctx context.Context, c Choice) (*Plan, error) {
	if runtime.GOOS != "windows" || runtime.GOARCH != "amd64" {
		return nil, errors.New("local directions downloads currently require Windows x64")
	}
	bounds, err := choiceBounds(c)
	if err != nil {
		return nil, err
	}
	names, err := segmentsFor(bounds)
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
	busy := m.progress.Busy
	m.mu.Unlock()
	if busy {
		return nil, errors.New("a directions download is already running")
	}
	plan := &Plan{ID: randomID(), Choice: c, Bounds: bounds, Expires: time.Now().Add(20 * time.Minute)}
	for _, part := range components {
		plan.Files = append(plan.Files, File{Name: part.Name, URL: part.URL, SHA256: part.SHA256, Bytes: part.Bytes})
		plan.DownloadBytes += part.Bytes
	}
	found := 0
	for _, name := range names {
		headCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		req, _ := http.NewRequestWithContext(headCtx, "HEAD", segmentOrigin+name, nil)
		res, e := m.client.Do(req)
		cancel()
		if e != nil {
			return nil, errors.New("directions data source is unavailable; retry the download check")
		}
		res.Body.Close()
		if res.StatusCode == 404 {
			continue
		}
		if res.StatusCode != 200 || res.ContentLength <= 0 || uint64(res.ContentLength) > limit {
			return nil, errors.New("directions source returned an invalid download size")
		}
		etag := res.Header.Get("ETag")
		if etag == "" || strings.HasPrefix(etag, "W/") {
			return nil, errors.New("directions source did not provide a stable file version")
		}
		plan.Files = append(plan.Files, File{Name: name, URL: segmentOrigin + name, Bytes: uint64(res.ContentLength), ETag: etag})
		plan.DownloadBytes += uint64(res.ContentLength)
		found++
		if plan.DownloadBytes > limit {
			return nil, errors.New("this directions area exceeds 2 GB; choose a smaller area")
		}
	}
	if found == 0 {
		return nil, errors.New("no routing data is available for this area")
	}
	plan.RequiredBytes = plan.DownloadBytes*3 + (512 << 20)
	if free, e := m.diskFree(m.root); e != nil || free < plan.RequiredBytes {
		return nil, errors.New("not enough free space for directions and a safe replacement; choose a smaller area")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.progress.Busy {
		return nil, errors.New("a download started while this plan was prepared")
	}
	m.plan = plan
	return plan, nil
}
func (m *Manager) Install(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.progress.Busy {
		return errors.New("a directions download is already running")
	}
	plan := m.plan
	if plan == nil || plan.ID != id || time.Now().After(plan.Expires) {
		return errors.New("check the download size again to create a fresh plan")
	}
	if free, e := m.diskFree(m.root); e != nil || free < plan.RequiredBytes {
		return errors.New("free space changed; create a fresh directions plan")
	}
	ctx, cancel := context.WithTimeout(m.ctx, 45*time.Minute)
	m.jobCancel = cancel
	m.plan = nil
	m.progress = Progress{Phase: "downloading", Message: "Downloading your selected local directions...", Total: plan.DownloadBytes, Busy: true}
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer cancel()
		err := m.runInstall(ctx, plan)
		m.mu.Lock()
		defer m.mu.Unlock()
		m.jobCancel = nil
		if err != nil {
			phase, message := "error", "Directions download failed. Retry with a fresh plan; previous directions were kept."
			if ctx.Err() != nil {
				phase, message = "cancelled", "Download cancelled. Previous directions were kept."
			}
			m.progress = Progress{Phase: phase, Message: message}
		} else {
			m.progress = Progress{Phase: "ready", Message: "Local driving, walking and cycling directions are ready.", Done: plan.DownloadBytes, Total: plan.DownloadBytes}
		}
	}()
	return nil
}
func (m *Manager) Cancel() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.jobCancel != nil {
		m.jobCancel()
	}
}
func (m *Manager) progressPhase(phase, message string) {
	m.mu.Lock()
	m.progress.Phase = phase
	m.progress.Message = message
	m.mu.Unlock()
}

func (m *Manager) download(ctx context.Context, file File, path string) (File, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", file.URL, nil)
	if err != nil {
		return file, err
	}
	if file.ETag != "" {
		req.Header.Set("If-Match", file.ETag)
	}
	res, err := m.client.Do(req)
	if err != nil {
		return file, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 || (res.ContentLength >= 0 && uint64(res.ContentLength) != file.Bytes) || (file.ETag != "" && res.Header.Get("ETag") != file.ETag) {
		return file, errors.New("download changed; create a fresh plan")
	}
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return file, err
	}
	h := sha256.New()
	reader := io.LimitReader(res.Body, int64(file.Bytes)+1)
	buffer := make([]byte, 64<<10)
	var total uint64
	for {
		n, e := reader.Read(buffer)
		if n > 0 {
			if _, err = out.Write(buffer[:n]); err != nil {
				break
			}
			h.Write(buffer[:n])
			total += uint64(n)
			m.mu.Lock()
			m.progress.Done += uint64(n)
			m.mu.Unlock()
		}
		if e == io.EOF {
			break
		}
		if e != nil {
			err = e
			break
		}
	}
	if err == nil {
		err = out.Sync()
	}
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return file, err
	}
	digest := hex.EncodeToString(h.Sum(nil))
	if total != file.Bytes || (file.SHA256 != "" && digest != file.SHA256) {
		return file, errors.New("directions checksum or length mismatch")
	}
	file.SHA256 = digest
	file.URL = ""
	return file, nil
}

func extractRuntime(archive, dest string) error {
	z, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer z.Close()
	var total uint64
	if len(z.File) > 4096 {
		return errors.New("runtime archive has too many files")
	}
	for _, entry := range z.File {
		parts := strings.SplitN(entry.Name, "/", 2)
		if len(parts) != 2 || parts[1] == "" {
			continue
		}
		name := parts[1]
		if strings.Contains(name, "\\") || strings.Contains(name, ":") || strings.HasPrefix(name, "/") || entry.Mode()&os.ModeSymlink != 0 {
			return errors.New("unsafe runtime archive path")
		}
		clean := filepath.Clean(filepath.FromSlash(name))
		if clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
			return errors.New("runtime archive escapes its folder")
		}
		path := filepath.Join(dest, clean)
		if entry.FileInfo().IsDir() {
			if err = os.MkdirAll(path, 0700); err != nil {
				return err
			}
			continue
		}
		total += entry.UncompressedSize64
		if total > 512<<20 || entry.UncompressedSize64 > 256<<20 {
			return errors.New("runtime archive expands beyond its limit")
		}
		if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}
		if err = extractFile(entry, path); err != nil {
			return err
		}
	}
	return nil
}
func extractFile(entry *zip.File, path string) error {
	in, e := entry.Open()
	if e != nil {
		return e
	}
	defer in.Close()
	out, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		return e
	}
	n, err := io.Copy(out, io.LimitReader(in, int64(entry.UncompressedSize64)+1))
	closeErr := out.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	if uint64(n) != entry.UncompressedSize64 {
		return errors.New("archive file length mismatch")
	}
	return nil
}
func extractRouter(archive, dest string) error {
	z, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer z.Close()
	count := 0
	wanted := map[string]string{"brouter-1.7.10-all.jar": "brouter.jar", "profiles2/lookups.dat": "profiles/lookups.dat"}
	for _, p := range []string{"car-fast", "car-eco", "trekking", "fastbike", "gravel", "mtb", "shortest", "hiking-mountain"} {
		wanted["profiles2/"+p+".brf"] = "profiles/" + p + ".brf"
	}
	for _, entry := range z.File {
		rel := strings.TrimPrefix(entry.Name, "brouter-1.7.10/")
		target, ok := wanted[rel]
		if !ok {
			continue
		}
		if entry.UncompressedSize64 > 16<<20 || entry.Mode()&os.ModeSymlink != 0 {
			return errors.New("invalid routing component")
		}
		path := filepath.Join(dest, filepath.FromSlash(target))
		if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}
		if err = extractFile(entry, path); err != nil {
			return err
		}
		count++
	}
	if count != len(wanted) {
		return errors.New("routing component is incomplete")
	}
	return nil
}
func removeStage(path string) error {
	if !strings.HasSuffix(path, ".partial") || !idPattern.MatchString(strings.TrimSuffix(filepath.Base(path), ".partial")) {
		return errors.New("refusing unknown staging path")
	}
	if err := filepath.WalkDir(path, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.Type()&os.ModeSymlink != 0 {
			return errors.New("refusing reparse path")
		}
		return nil
	}); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.RemoveAll(path)
}
func (m *Manager) runInstall(ctx context.Context, plan *Plan) error {
	stage := filepath.Join(m.root, plan.ID+".partial")
	if err := os.Mkdir(stage, 0700); err != nil {
		return err
	}
	defer removeStage(stage)
	if err := os.Mkdir(filepath.Join(stage, "segments"), 0700); err != nil {
		return err
	}
	info := &Info{ID: plan.ID, Choice: plan.Choice, Bounds: plan.Bounds, Installed: time.Now().UTC().Format(time.RFC3339), Engine: Version}
	for _, file := range plan.Files {
		path := filepath.Join(stage, file.Name)
		if segmentPattern.MatchString(file.Name) {
			path = filepath.Join(stage, "segments", file.Name)
		}
		saved, err := m.download(ctx, file, path)
		if err != nil {
			return err
		}
		info.Files = append(info.Files, saved)
	}
	m.progressPhase("verifying", "Checking the local routing engine and data integrity...")
	if err := extractRuntime(filepath.Join(stage, "runtime.zip"), filepath.Join(stage, "runtime")); err != nil {
		return err
	}
	if err := extractRouter(filepath.Join(stage, "brouter.zip"), stage); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(stage, "MapsRoutingBridge.class"), bridge, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(stage, "LICENSE-BRouter.txt"), routerLicense, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(stage, "NOTICES.txt"), []byte("BRouter 1.7.10: MIT; see LICENSE-BRouter.txt. OpenStreetMap-derived routing data: ODbL, © OpenStreetMap contributors, https://www.openstreetmap.org/copyright\nTemurin JRE 21.0.12+8: GPL2 with Classpath exception and bundled components; complete notices in runtime/legal/. Corresponding source: https://github.com/adoptium/jdk21u/tree/jdk-21.0.12%2B8\n"), 0600); err != nil {
		return err
	}
	os.Remove(filepath.Join(stage, "runtime.zip"))
	os.Remove(filepath.Join(stage, "brouter.zip"))
	if err := filepath.WalkDir(stage, func(path string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.Type()&os.ModeSymlink != 0 {
			return errors.New("unexpected link in directions data")
		}
		if !d.IsDir() {
			st, e := d.Info()
			if e != nil {
				return e
			}
			info.Bytes += uint64(st.Size())
		}
		return nil
	}); err != nil {
		return err
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	select {
	case m.routeGate <- struct{}{}:
		defer func() { <-m.routeGate }()
	case <-ctx.Done():
		return ctx.Err()
	}
	m.engineMu.Lock()
	defer m.engineMu.Unlock()
	// Validate the staged graph before switching the active pointer. Existing
	// directions stay available if download, extraction or startup fails.
	candidate, err := m.launch(ctx, stage)
	if err != nil {
		return err
	}
	candidate.stop()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	final := filepath.Join(m.root, plan.ID)
	if err = os.Rename(stage, final); err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = m.removePack(plan.ID)
		}
	}()
	if err = os.WriteFile(filepath.Join(final, "info.json"), mustJSON(info), 0600); err != nil {
		return err
	}
	tmp := filepath.Join(m.root, plan.ID+".json.partial")
	if err = os.WriteFile(tmp, mustJSON(info), 0600); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if ctx.Err() != nil {
		os.Remove(tmp)
		return ctx.Err()
	}
	if err = os.Rename(tmp, filepath.Join(m.root, "active.json")); err != nil {
		os.Remove(tmp)
		return err
	}
	if m.engine != nil {
		m.engine.stop()
		m.engine = nil
	}
	previous := m.info
	m.info = info
	committed = true
	if previous != nil && previous.ID != info.ID {
		_ = m.removePack(previous.ID)
	}
	return nil
}

// Only installer-created, direct-child IDs can be removed. Refuse links at any depth.
func (m *Manager) removePack(id string) error {
	if !idPattern.MatchString(id) {
		return errors.New("invalid directions pack ID")
	}
	path := filepath.Join(m.root, id)
	if filepath.Dir(path) != filepath.Clean(m.root) {
		return errors.New("directions pack escaped its root")
	}
	if err := filepath.WalkDir(path, func(_ string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.Type()&os.ModeSymlink != 0 {
			return errors.New("refusing linked directions pack")
		}
		return nil
	}); err != nil {
		return err
	}
	return os.RemoveAll(path)
}
func mustJSON(v any) []byte {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return append(b, '\n')
}
