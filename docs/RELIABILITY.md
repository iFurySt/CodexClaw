# 稳定性与可运维性

这里用来定义项目的运行质量底线。

## CodexClaw daemon

- `daemon start` 会把 daemon detach 到后台长期驻留；`daemon run` 保持为前台长进程，适合 launchd/systemd 托管。
- daemon 从 `~/.codexclaw/config.toml` 读取 `workspace`、`interval`、`timeout` 和 `prompt`。
- 启动后会先立即执行一次 tick，然后按配置的 `interval` 串行执行后续 tick。
- 单次 Codex 调用受配置的 `timeout` 控制，默认 20 分钟。
- 同一个配置目录下通过 `daemon.lock` 防止重复启动；lock 文件保存 pid，供 `stop`、`restart` 和 `status` 使用。
- `daemon start` 的 stdout/stderr 默认写入 `state_dir/daemon.log`。
- 当前只提供文本日志；后续如果需要更强长期托管，应补 JSON 日志、最近一次 tick 状态和 launchd/systemd 健康检查。

建议维护的内容包括：

- 启动、健康检查和基本可用性要求。
- 日志、指标、链路的采集和访问约定。
- timeout、retry、backoff 的默认策略。
- 本地和 CI 的关键路径验证方式。
- 常见故障、排查路径和恢复步骤。

CI/CD 流程结构和 release 自动化的默认方案，统一写在 `docs/CICD.md`。
