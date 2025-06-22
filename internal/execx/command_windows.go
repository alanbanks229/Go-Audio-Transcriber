//go:build windows
// +build windows

package execx

import (
	"os/exec"
	"syscall"
)

// Command returns *exec.Cmd with console-window suppression baked in.
func Command(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)

	// Works on Go 1.21+.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	return cmd
}
