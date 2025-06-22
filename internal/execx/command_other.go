//go:build !windows
// +build !windows

package execx

import "os/exec"

// Non-Windows: just hand back exec.Command unmodified.
func Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
