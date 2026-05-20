# 供应链安全

这份文档定义模板默认采用的供应链安全做法。

## 当前控制项

- 通过 `scripts/check-action-pinning.sh` 要求 workflow 中的 GitHub Actions 固定到不可变 commit SHA。
- Homebrew 发布 workflow 只使用已 pin 的 checkout/setup-go action，其余发布逻辑在仓库脚本和 shell 中完成。
- 不在仓库中保存 Homebrew tap 写入凭据；tap 推送权限来自 GitHub secret 或后续 trusted publishing 配置。
- 所有 GitHub Actions 都固定到不可变的 commit SHA，而不是漂移的版本标签。

## 当前对应关系

- `scripts/check-action-pinning.sh`：如果 workflow 里出现浮动 tag 而不是 SHA，直接让 CI 失败。
- `.github/workflows/homebrew.yml`：用 tag 源码包生成 formula 并推送到 tap。

## 限制和前提

- 当前没有单独的 dependency review、OSV、SBOM 或 provenance workflow；第一版先避免为了模板感堆过重的 CI/CD。
- Homebrew tap 写入权限需要单独配置，推荐优先用短期或受限 token；如果平台支持 trusted publishing，再切到 OIDC/trust 模式。
- OpenSSF Scorecard 默认不启用，因为当前仓库还没有真实分支保护、release 历史和 SAST 姿态可以评分；等仓库规则配置完成后再按需加回。

## 项目落地后建议继续做的事

- 锁定并提交项目真实依赖的 lockfile。
- 让构建过程尽量可重复、可验证。
- 如果发布链路开始分发二进制制品，再补 SBOM、provenance 和对应校验。
