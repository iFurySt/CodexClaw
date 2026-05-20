# CI/CD 说明

CodexClaw 的 CI/CD 现在只保留两条路径：基础 CI 和 Homebrew formula 发布。

## 默认包含的内容

- `ci.yml`：仓库级检查，覆盖 docs、repo hygiene、Go 测试、Markdown 和 shell 脚本校验。
- `homebrew.yml`：在 `v*` tag 或手动触发时生成 `codexclaw.rb`，并自动提交到 `iFurySt/homebrew-CodexClaw` tap 仓库。

## 设计原则

这套流水线的目标是保持少而清晰：CI 验证仓库，Homebrew workflow 负责发布 formula 到 tap。

Homebrew formula 直接引用 GitHub tag 源码包，并在安装时用 Go 构建 `./cmd/codexclaw`。这让第一版不需要维护额外二进制制品、SBOM 或 provenance 流水线。

所有 GitHub Actions 都已经 pin 到 commit SHA。后续升级 action 时，也要继续保持这个约束。

## Homebrew 发布

发布使用单独的 tap 仓库 `iFurySt/homebrew-CodexClaw`。Homebrew 官方建议 GitHub tap 仓库使用 `homebrew-` 前缀，这样用户可以用短格式 tap 名。

主仓库需要配置 `HOMEBREW_TAP_TOKEN` secret，用来写入 `iFurySt/homebrew-CodexClaw`。这个 token 至少需要对 tap 仓库有 contents write 权限。维护者可以用当前已登录的 `gh` 凭据配置：

```sh
gh auth token | gh secret set HOMEBREW_TAP_TOKEN --repo iFurySt/CodexClaw
```

`GITHUB_TOKEN` 只能写当前仓库，不能直接推送到独立 tap 仓库，所以这里显式使用 tap token。后续如果平台提供适合 GitHub-to-GitHub tap 写入的 trusted publishing，再切换过去。

本地或 tap workflow 可以用下面的命令预览 formula：

```sh
./scripts/generate-homebrew-formula.sh \
  --repo iFurySt/CodexClaw \
  --tag v0.2.0 \
  --sha256 "$(printf test | sha256sum | awk '{print $1}')"
```
