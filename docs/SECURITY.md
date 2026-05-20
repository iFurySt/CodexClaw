# 安全默认约束

这份文档用于把安全默认值讲清楚，避免实现逐步演进时越走越散。

## CodexClaw daemon

- 默认通过 `codex exec --sandbox workspace-write` 执行配置里的任务。
- 不在 CodexClaw 中保存 OpenAI/Codex 凭据，认证沿用本机 Codex CLI 的配置。
- `~/.codexclaw/config.toml` 是本机信任边界内的输入；不要把 daemon 直接暴露成远程执行 API。
- `workspace` 和 `prompt` 会直接影响 Codex 执行范围和行为，配置文件应按本机用户私有文件处理。
- 后续如果接入 Codex App Server，只允许本机 Unix socket 或明确认证的 loopback/WebSocket，不默认开放公网监听。

建议维护的内容：

- 认证与授权约束。
- 密钥和环境变量管理方式。
- 依赖治理与供应链安全要求。
- 数据分级、脱敏与保留策略。
- 对外 API、Webhook、文件上传和沙箱执行的规则。

仓库级的依赖、SBOM 和 provenance 默认能力，统一写在 `docs/SUPPLY_CHAIN_SECURITY.md`。
