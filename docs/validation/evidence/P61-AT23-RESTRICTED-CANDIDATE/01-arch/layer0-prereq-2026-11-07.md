# 层 0 前置探针记录（2026-11-07，步骤 0 收尾）

状态：**层 0 BLOCKED（前置缺失，停下升级用户裁决）**。

## 连通性预检（R-08.4）

- 直连 `api.deepseek.com:443`：`TcpTestSucceeded=True`，远端地址 119.188.220.215（无代理、无 VM）。
- 判定：网络可达，不触发 R-08.4 的代理/端点失败路径。

## BYOK 密钥可用性检查（只报存在性，不读值）

| 通道 | 存在性 |
|---|---|
| 进程环境变量 `COPILOT_PROVIDER_API_KEY` / `DEEPSEEK_API_KEY` | 不存在 |
| Windows 凭据管理器 `ProofRail` 前缀条目（`cmdkey /list`） | 不存在 |
| 其余 5 项 BYOK 环境变量（TYPE / BASE_URL / MODEL / 两个 token 上限） | 均未设置 |

## Copilot CLI 现状

- 钉版二进制在位：`%APPDATA%\npm\node_modules\@github\copilot\node_modules\@github\copilot-win32-x64\copilot.exe`，144,796,448 字节（与 T026 钉版一致），`--version` 输出 1.0.85（版本字符串非稳定身份，见 B1 的 DR-2）。
- 封装函数 `Invoke-CopilotCli` / `buildCopilotArgs` 尚未落地（步骤 0 审计结论：3 个脚本 6 处直接调用点 + 1 处 Go 产品调用点，不足「散落 5 处以上」阈值 ⇒ 不触发步骤 0 前置收敛，由 ② 新建封装）。

## 结论

- R-08.6 步骤 0 最小验证（`copilot --no-auto-update -p "say hello"`）因缺少 `COPILOT_PROVIDER_API_KEY` 无法执行。
- 按 R-08.4 / R-08.6 失败路径与本片红线 3：主控不得自行启用代理、切端点或自行配置密钥 ⇒ **停下升级用户裁决**。
