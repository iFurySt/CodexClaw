//go:build !unix

package daemon

import (
	"errors"
	"syscall"
)

func detachedSysProcAttr() *syscall.SysProcAttr {
	return nil
}

func signalTerminate(pid int) error {
	return errors.New("daemon process signals are not supported on this platform")
}

func isProcessRunning(pid int) bool {
	return false
}
