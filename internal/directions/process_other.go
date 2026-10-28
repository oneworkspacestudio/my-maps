//go:build !windows

package directions

import (
	"errors"
	"io"
	"os/exec"
)

func hideWindow(*exec.Cmd) {}
func ownProcess(*exec.Cmd) (io.Closer, error) {
	return nil, errors.New("downloaded local directions currently require Windows")
}
