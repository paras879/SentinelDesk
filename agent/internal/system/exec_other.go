//go:build !windows

package system

import "os/exec"

func newExecCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
