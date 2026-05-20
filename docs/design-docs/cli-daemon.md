# CLI Daemon 设计

## 当前状态

CodexClaw 现在有一个最小 Go/Cobra 二进制入口：`codexclaw daemon`。

第一版只解决一件事：读取本机全局配置 `~/.codexclaw/config.toml`，按固定节奏在配置的 repo 里执行 `codex exec` 和配置的 prompt。

这等价于一个更可维护的 shell loop：配置文件负责描述“在哪个目录、多久一次、让 Codex 做什么”，daemon 只负责准时执行和避免重复启动。

## 设计判断

第一版应该优先走 Codex CLI，而不是直接接 Codex App Server。

原因：

- `codex exec` 已经是 Codex 官方提供的非交互自动化入口，支持 `--cd`、sandbox、model、profile 等必要控制。
- Codex App Server 的 JSON-RPC 协议能力更完整，但连接、初始化、thread 生命周期、事件流和兼容性处理会让第一版明显变重。
- CodexClaw 不需要闲置确认协议、任务账本、通道投递和多 agent 调度；定时执行一条配置好的 Codex 命令就够了。

因此 MVP 的边界是：CodexClaw 负责“什么时候、在哪个目录、带什么 prompt 执行 Codex”，Codex 负责“具体怎么完成 prompt 里的任务”。

## MVP 架构

```text
codexclaw
  daemon run
    read config.toml
    acquire lock
    run tick immediately
    wait interval
    run next tick

  daemon once
    run one tick and exit

tick
  build codex exec command
  apply per-run timeout
  capture stdout/stderr
```

模块边界：

- `cmd/codexclaw`：二进制入口，只负责执行 root command。
- `internal/cli`：Cobra 命令、flags、信号处理和用户输出。
- `internal/daemon`：配置读取、Codex 命令构造、单次执行、循环调度和本地锁。

## 配置文件

默认配置路径是 `~/.codexclaw/config.toml`。

```toml
workspace = "/Users/bytedance/projects/github/aifi"
interval = "30m"
timeout = "20m"
prompt = """
Read the repo instructions and continue any documented work that should be done now.
"""

codex_bin = "codex"
sandbox = "workspace-write"
# model = "gpt-5.4-codex"
# profile = "default"
# codex_args = ["--json"]
```

第一版只支持一个任务。要换 repo、间隔或任务，就直接改这份文件。

## 命令面

```sh
go run ./cmd/codexclaw --help
go run ./cmd/codexclaw config instructions
go run ./cmd/codexclaw daemon once --dry-run
go run ./cmd/codexclaw daemon run
```

常用参数：

- `--config`：配置文件路径，默认 `~/.codexclaw/config.toml`。
- `--dry-run`：只打印将执行的命令。

`config instructions` 是给 AI/Agent 使用的只读辅助命令。它打印 `config.toml` 的字段、示例和编辑规则，不读取或修改本机配置。

## 第一版不做

- 不内置任务队列。任务发现交给 Codex 读取仓库文件、计划和外部上下文。
- 不做远程控制面。先保证本机 daemon 稳。
- 不直接实现 App Server client。等 CLI 方案跑通后，再抽 `Runner` 接口接入。
- 不做多任务配置。先保持一个 config 对应一个 repo/prompt/interval。
- 不做系统服务安装器。后续可以补 launchd/systemd 单独命令。
- 不做闲置确认语义。Codex 输出什么就透传什么。

## 后续扩展点

1. `Runner` 抽象：保留 CLI runner，新增 App Server runner。
2. `install-service`：生成 launchd/systemd user service。
3. `status`：读取 lock/state，展示最近一次 tick 结果。
4. JSON 日志：给长期运行和 CI smoke 更稳定的机器可读输出。
5. 多任务配置：真实需要时再从单任务 config 演进。

这些能力都应该在真实使用中出现明确痛点后再加。
