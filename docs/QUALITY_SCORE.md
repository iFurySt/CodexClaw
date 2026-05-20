# 质量评分

用这份文档按产品区域和架构层次记录当前质量水位，方便持续知道最薄弱的地方在哪。

## 建议的评分标准

- `A`：覆盖完整、行为稳定、文档清楚、运行风险低。
- `B`：整体可接受，但还有明确短板。
- `C`：能用，但需要针对性补强。
- `D`：脆弱、缺少规范，或很多行为尚未定义。

## 初始模板

| 区域 | 评分 | 原因 | 下一步 |
| --- | --- | --- | --- |
| 产品面 | B | 已有最小 daemon 用户路径：从 `~/.codexclaw/config.toml` 定时执行 Codex，并支持后台 `start/stop/restart/status`。 | 补真实使用反馈，再决定是否需要 service 安装。 |
| 架构文档 | B | 已替换为 CodexClaw daemon 的当前结构和边界。 | App Server runner 进入实现前先补协议设计。 |
| 测试 | C | Go 单元测试覆盖命令构造、config.toml 加载和 daemon 生命周期命令暴露面。 | 补更完整的 CLI smoke 和后续 service 安装测试。 |
| 可观测性 | C | daemon 有基础文本日志、timeout 和本地锁。 | 增加 JSON 日志与最近一次 tick 状态。 |
| 安全 | C | 默认复用 Codex CLI 认证，并使用 workspace-write sandbox。 | App Server/远程控制前补认证与监听边界验证。 |
