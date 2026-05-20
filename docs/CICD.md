# CI/CD 说明

CodexClaw 的 CI/CD 现在只保留两条路径：基础 CI 和 Homebrew formula 发布。

## 默认包含的内容

- `ci.yml`：仓库级检查，覆盖 docs、repo hygiene、Go 测试、Markdown 和 shell 脚本校验。
- `homebrew.yml`：在 `v*` tag 或手动触发时生成 `codexclaw.rb`，并推送到 Homebrew tap 仓库。

## 设计原则

这套流水线的目标是保持少而清晰：CI 只验证仓库，发布只更新 Homebrew tap。

Homebrew formula 直接引用 GitHub tag 源码包，并在安装时用 Go 构建 `./cmd/codexclaw`。这让第一版不需要维护额外二进制制品、SBOM 或 provenance 流水线。

所有 GitHub Actions 都已经 pin 到 commit SHA。后续升级 action 时，也要继续保持这个约束。

## Homebrew 发布

发布前需要准备一个 tap 仓库，例如 `iFurySt/homebrew-tap`。Homebrew 官方建议 GitHub tap 仓库使用 `homebrew-` 前缀，这样用户可以用短格式 `brew install owner/tap/formula`。

`homebrew.yml` 需要以下配置之一：

- 手动触发时填写 `tap_repository`。
- 或在仓库变量中设置 `HOMEBREW_TAP_REPOSITORY`。

推送 tap 需要写权限。当前 workflow 读取 `HOMEBREW_TAP_TOKEN` secret；如果后续 Homebrew/tap 支持 trusted publishing，可用 OBU 在 GitHub 界面补对应 trust，再收敛掉长期 token。

本地可以用下面的命令预览 formula：

```sh
./scripts/generate-homebrew-formula.sh \
  --repo iFurySt/CodexClaw \
  --tag v0.1.0 \
  --sha256 "$(printf test | sha256sum | awk '{print $1}')"
```
