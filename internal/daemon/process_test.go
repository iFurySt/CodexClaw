package daemon

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestReadProcessStatusDetectsRunningPID(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, lockFileName), []byte(" "+strconv.Itoa(os.Getpid())+" \n"), 0o600); err != nil {
		t.Fatal(err)
	}

	status, err := ReadProcessStatus(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !status.Running || status.PID != os.Getpid() {
		t.Fatalf("status = %#v", status)
	}
}

func TestRemoveStaleLockRemovesInvalidPID(t *testing.T) {
	dir := t.TempDir()
	lockPath := filepath.Join(dir, lockFileName)
	if err := os.WriteFile(lockPath, []byte("not-a-pid\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	status, err := ReadProcessStatus(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !status.Stale {
		t.Fatal("expected stale status")
	}

	if err := RemoveStaleLock(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(lockPath); !os.IsNotExist(err) {
		t.Fatalf("expected stale lock to be removed, stat err=%v", err)
	}
}
