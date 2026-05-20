# 架构总览

CodexClaw 当前从 Agent 协作模板演进为一个轻量 Codex CLI daemon。

核心目标很简单：按固定间隔在指定仓库目录里执行一条 `codex exec`。要跑哪个 repo、多久跑一次、传什么 prompt，都由本机全局配置文件 `~/.codexclaw/config.toml` 决定。

## 当前仓库结构

- `cmd/codexclaw/`：Go 二进制入口。
- `internal/cli/`：Cobra 命令层，负责 flags、信号处理和输出。
- `internal/daemon/`：daemon 核心，负责读取配置、单次 Codex 调用、循环调度和本地锁。
- `scripts/`：仓库级自动化脚本，供人和 Agent 直接调用。
- `docs/`：仓库知识库，也是本地规则和上下文的正式来源。

## 运行拓扑

```text
operator
  -> codexclaw daemon start
      -> detached codexclaw daemon run
          -> read ~/.codexclaw/config.toml
          -> codex exec --cd <workspace> <prompt>

launchd / systemd
  -> codexclaw daemon run
      -> read ~/.codexclaw/config.toml
      -> codex exec --cd <workspace> <prompt>
```

第一版只对接 Codex CLI。Codex App Server 是后续 runner 适配方向，不进入当前最小路径。

## 边界约束

- CodexClaw 只负责调度和进程生命周期，不实现任务队列和 agent 决策逻辑。
- Codex 执行权限由 `codex exec` 的 sandbox、profile、model 等参数控制。
- daemon 默认串行执行 tick；一次 Codex run 未结束前不会并发启动下一次。
- 本地 lock/pid 防止同一个 state dir 下重复启动多个 daemon，并支持 `start`、`stop`、`restart`、`status` 生命周期命令。
- 第一版只支持一个配置任务；要换 repo 或 prompt，就直接编辑 `~/.codexclaw/config.toml`。
- 行为变化要同步更新 `docs/design-docs/cli-daemon.md`。

## 后续补齐

- launchd/systemd user service 安装命令。
- App Server runner 的协议封装和兼容性验证。
- 最近一次 tick 状态、JSON 日志和更完整的健康检查。
