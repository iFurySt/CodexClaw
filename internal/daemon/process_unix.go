//go:build unix

package daemon

import (
	"errors"
	"syscall"
)

func detachedSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}

func signalTerminate(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}

func isProcessRunning(pid int) bool {
	err := syscall.Kill(pid, syscall.Signal(0))
	return err == nil || errors.Is(err, syscall.EPERM)
}
