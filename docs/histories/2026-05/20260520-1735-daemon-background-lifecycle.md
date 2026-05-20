## [2026-05-20 17:35] | Task: daemon background lifecycle

### Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex CLI`

### User Query

> 希望 `codexclaw daemon` 像 nginx 一样长期驻留后台，并提供停止、重启进程的子命令。

### Changes Overview

**Scope:** CLI daemon lifecycle and process state.

**Key Actions:**

- **Lifecycle commands**: 新增 `daemon start`、`daemon stop`、`daemon restart`、`daemon status`。
- **Background process**: `start` 会 detached 启动前台 `daemon run`，并将 stdout/stderr 写入 `state_dir/daemon.log`。
- **PID state**: 复用 `daemon.lock` 保存 pid，支持状态查询、停止进程和清理 stale lock。
- **Docs and tests**: 更新 README、架构/设计/可靠性文档，并补充生命周期相关单元测试。
- **Release prep**: 将 CLI 版本提升到 `0.2.0`，并补充用户可见发布记录。

### Design Intent (Why)

保留 `daemon run` 作为 launchd/systemd 友好的前台长进程入口，同时提供面向手动使用的 nginx 风格后台生命周期命令。这样不会破坏已有托管模型，也能满足本地长期驻留和显式停止/重启的使用习惯。

### Files Modified

- `internal/cli/daemon.go`
- `internal/daemon/lock.go`
- `internal/daemon/process.go`
- `internal/daemon/process_unix.go`
- `internal/daemon/process_unsupported.go`
- `internal/cli/daemon_test.go`
- `internal/daemon/process_test.go`
- `README.md`
- `internal/cli/root.go`
- `docs/ARCHITECTURE.md`
- `docs/CICD.md`
- `docs/design-docs/cli-daemon.md`
- `docs/RELIABILITY.md`
- `docs/QUALITY_SCORE.md`
- `docs/releases/feature-release-notes.md`
