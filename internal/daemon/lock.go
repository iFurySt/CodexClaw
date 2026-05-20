package daemon

import (
	"fmt"
	"os"
	"path/filepath"
)

const lockFileName = "daemon.lock"
const logFileName = "daemon.log"

type Lock struct {
	path string
	file *os.File
}

func DefaultStateDir() string {
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".codexclaw")
	}
	return ".codexclaw"
}

func LockPath(stateDir string) string {
	if stateDir == "" {
		stateDir = DefaultStateDir()
	}
	return filepath.Join(stateDir, lockFileName)
}

func LogPath(stateDir string) string {
	if stateDir == "" {
		stateDir = DefaultStateDir()
	}
	return filepath.Join(stateDir, logFileName)
}

func AcquireLock(stateDir string) (*Lock, error) {
	if stateDir == "" {
		stateDir = DefaultStateDir()
	}
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		return nil, fmt.Errorf("create state dir: %w", err)
	}

	path := LockPath(stateDir)
	if err := RemoveStaleLock(stateDir); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("daemon lock exists at %s; another codexclaw daemon may be running", path)
		}
		return nil, fmt.Errorf("create daemon lock: %w", err)
	}
	if _, err := fmt.Fprintf(file, "%d\n", os.Getpid()); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return nil, fmt.Errorf("write daemon lock: %w", err)
	}
	return &Lock{path: path, file: file}, nil
}

func (l *Lock) Release() error {
	if l == nil {
		return nil
	}
	var firstErr error
	if l.file != nil {
		firstErr = l.file.Close()
	}
	if err := os.Remove(l.path); err != nil && !os.IsNotExist(err) && firstErr == nil {
		firstErr = err
	}
	return firstErr
}
