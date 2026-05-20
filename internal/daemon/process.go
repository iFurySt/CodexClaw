package daemon

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ProcessStatus struct {
	StateDir string
	LockPath string
	LogPath  string
	PID      int
	Running  bool
	Stale    bool
}

type StartResult struct {
	PID     int
	LogPath string
}

type StopResult struct {
	PID          int
	WasRunning   bool
	RemovedStale bool
}

func ReadProcessStatus(stateDir string) (ProcessStatus, error) {
	if stateDir == "" {
		stateDir = DefaultStateDir()
	}

	status := ProcessStatus{
		StateDir: stateDir,
		LockPath: LockPath(stateDir),
		LogPath:  LogPath(stateDir),
	}

	data, err := os.ReadFile(status.LockPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return status, nil
		}
		return status, fmt.Errorf("read daemon lock: %w", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || pid <= 0 {
		status.Stale = true
		return status, nil
	}

	status.PID = pid
	status.Running = IsProcessRunning(pid)
	status.Stale = !status.Running
	return status, nil
}

func RemoveStaleLock(stateDir string) error {
	status, err := ReadProcessStatus(stateDir)
	if err != nil {
		return err
	}
	if !status.Stale {
		return nil
	}
	if err := os.Remove(status.LockPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove stale daemon lock: %w", err)
	}
	return nil
}

func StartDetached(configPath string, stateDir string) (StartResult, error) {
	if stateDir == "" {
		stateDir = DefaultStateDir()
	}
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		return StartResult{}, fmt.Errorf("create state dir: %w", err)
	}
	if err := RemoveStaleLock(stateDir); err != nil {
		return StartResult{}, err
	}

	status, err := ReadProcessStatus(stateDir)
	if err != nil {
		return StartResult{}, err
	}
	if status.Running {
		return StartResult{}, fmt.Errorf("daemon already running pid=%d", status.PID)
	}

	executable, err := os.Executable()
	if err != nil {
		return StartResult{}, fmt.Errorf("find executable: %w", err)
	}

	logFile, err := os.OpenFile(LogPath(stateDir), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return StartResult{}, fmt.Errorf("open daemon log: %w", err)
	}
	defer logFile.Close()

	cmd := exec.Command(executable, "daemon", "run", "--config", configPath)
	cmd.Stdin = nil
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = detachedSysProcAttr()

	if err := cmd.Start(); err != nil {
		return StartResult{}, fmt.Errorf("start daemon: %w", err)
	}

	pid := cmd.Process.Pid
	if err := cmd.Process.Release(); err != nil {
		return StartResult{}, fmt.Errorf("release daemon process: %w", err)
	}

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		status, err := ReadProcessStatus(stateDir)
		if err != nil {
			return StartResult{}, err
		}
		if status.Running && status.PID == pid {
			return StartResult{PID: pid, LogPath: LogPath(stateDir)}, nil
		}
		if !IsProcessRunning(pid) {
			return StartResult{}, fmt.Errorf("daemon exited during startup; see log at %s", LogPath(stateDir))
		}
		time.Sleep(100 * time.Millisecond)
	}

	return StartResult{}, fmt.Errorf("daemon did not create lock within startup timeout; see log at %s", LogPath(stateDir))
}

func StopProcess(stateDir string, timeout time.Duration) (StopResult, error) {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	status, err := ReadProcessStatus(stateDir)
	if err != nil {
		return StopResult{}, err
	}
	result := StopResult{PID: status.PID}

	if status.Stale {
		if err := os.Remove(status.LockPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return result, fmt.Errorf("remove stale daemon lock: %w", err)
		}
		result.RemovedStale = true
		return result, nil
	}
	if !status.Running {
		return result, nil
	}

	result.WasRunning = true
	if err := signalTerminate(status.PID); err != nil {
		return result, fmt.Errorf("stop daemon pid=%d: %w", status.PID, err)
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		nextStatus, err := ReadProcessStatus(stateDir)
		if err != nil {
			return result, err
		}
		if !nextStatus.Running || nextStatus.PID != status.PID {
			if nextStatus.Stale {
				if err := os.Remove(nextStatus.LockPath); err != nil && !errors.Is(err, os.ErrNotExist) {
					return result, fmt.Errorf("remove stale daemon lock: %w", err)
				}
			}
			return result, nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return result, fmt.Errorf("daemon pid=%d did not stop within %s", status.PID, timeout)
}

func IsProcessRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	return isProcessRunning(pid)
}
