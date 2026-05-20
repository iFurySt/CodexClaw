# 功能发布记录

## 2026-05

| 日期 | 功能域 | 用户价值 | 变更摘要 |
| --- | --- | --- | --- |
| 2026-05-20 | Daemon 生命周期 | 可以像 nginx 一样让 CodexClaw 长期驻留后台，并用命令查看、停止或重启进程。 | 新增 `codexclaw daemon start/stop/restart/status`，后台日志写入 `state_dir/daemon.log`，`daemon.lock` 复用为 pid 状态文件。 |

## 2026-04

| 日期 | 功能域 | 用户价值 | 变更摘要 |
| --- | --- | --- | --- |
| 2026-04-08 | 模板仓库 | 提供了一套可直接用于新项目启动的 Agent-first 基础模板。 | 补齐了 AGENTS 入口、execution plan、history、release note、CI/CD 和供应链安全骨架。 |
