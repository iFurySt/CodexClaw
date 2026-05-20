# CodexClaw

English version: [`CodexClaw`](https://github.com/iFurySt/CodexClaw)

## 简介

一个面向 Agent 协作开发的基础仓库模板，可以用来启动任何你想做的产品或服务。

## 快速开始

当前仓库已经包含一个最小 CodexClaw daemon：

```toml
# ~/.codexclaw/config.toml
workspace = "/Users/bytedance/projects/github/aifi"
interval = "30m"
timeout = "20m"
prompt = "Read the repo instructions and continue any documented work that should be done now."
```

```sh
go run ./cmd/codexclaw config instructions
go run ./cmd/codexclaw daemon once --dry-run
go run ./cmd/codexclaw daemon start
go run ./cmd/codexclaw daemon status
go run ./cmd/codexclaw daemon stop
```

`daemon start` 会把 daemon 长期驻留到后台，并按 `~/.codexclaw/config.toml` 定时执行 `codex exec --cd <workspace> <prompt>`。如果需要交给 launchd/systemd 这类进程管理器托管，可以直接使用前台命令 `codexclaw daemon run`。设计说明见 [`docs/design-docs/cli-daemon.md`](docs/design-docs/cli-daemon.md)。

可以在这个仓库右上角直接使用 GitHub 的模板流程：

1. 选择 **Use this template**。
2. 选择 [**Create a new repository**](https://github.com/new?template_name=CodexClaw&template_owner=iFurySt)。

也可以在新仓库或已有仓库里用 [`harness-cli`](https://github.com/iFurySt/harness-cli) 初始化。先通过 npm 安装：

```sh
npm install -g @ifuryst/harness-cli
```

然后运行：

```sh
harness-cli init --language zh
```

`harness-cli` 需要 Node.js 18+，并且本机 `PATH` 中需要有 Go。

## 许可证

[MIT](LICENSE)

## 备注

这套方法主要来自我们自己的持续实践和整理，同时也吸收了 OpenAI 在 [harness engineering 文章](https://openai.com/index/harness-engineering/) 中的一部分思路，最后汇总成了这个模板。
