# 层 0 BYOK 可行性验证结果（2026-11-08）

判定：**层 0 全绿**（R-08.6 步骤 0 最小验证通过，可按闸门进层 1）。

## 前置

- 凭据：用户按裁定 (a) 写入 `windows-credential:ProofRail/deepseek`；主控经产品路径复核 `prfrail.exe secret status --json` ⇒ `"exists": true`。
- 注入：`tmp/P61-AT23-RESTRICTED-CANDIDATE/Invoke-ByokLayer0.ps1`（最终 sha256 `DA9C358F…0006DE4`，先经结构体布局修正、后加清理自检，完整哈希演变见 `01-arch/byok-layer0-vs-repo.txt`）。
- 凭据只读探针：`ok=True lastError=0 blobSize=35 type=1 persist=2`（只报元数据，未读值）。

## 调用与结果

- argv：`--no-auto-update -p say hello`（与起跑消息层 0 第 2 项逐字一致）。
- 版本（同二进制 `--version`）：`GitHub Copilot CLI 1.0.83.`
- 退出码：**0**。
- 输出：模型正常回复（含 "Hello! 👋 I'm your GitHub Copilot CLI assistant, ready to help with the prfrail (ProofRail) repository."）；无工具调用。
- CLI 计量行：`Tokens ↑ 18.1k • ↓ 37`（1.0 fallback 估算，见下）。
- 警告（逐字）：`! No token count multiplier is known for Anthropic model "deepseek-flash". Token estimates will use the 1.0 fallback; configure a recognized canonical model ID to enable model-specific token counting.`
- **真实计费（BYOK 铁证）**：本次调用实际计费 **¥0.01 CNY**（用户于 2026-11-08 从 DeepSeek 网页用量页核实）。host 认证不扣 DeepSeek 费用，故这是 BYOK 通道真实调用的直接证据。本笔 = R-08.2 阶段 1「逐次声明」的**第一笔真实付费**。
- **环境清理机械读数**：`BYOK env cleanup: verified (all six cleared)`（`tmp/P61-AT23-RESTRICTED-CANDIDATE/cleanup-check.txt`，调用后 finally 内逐变量断言全空）。
- **已知边界（层 3 预算控制以此为准）**：BYOK 下 CLI 的 token 计数为 1.0 fallback 估算，不是真实计费依据；真实计费以 DeepSeek 侧为准。

## 记录与观察

- 证据件：`tmp/P61-AT23-RESTRICTED-CANDIDATE/layer0/argv.txt`、`layer0/version.txt`、`tmp/P61-AT23-RESTRICTED-CANDIDATE/layer0-run2.txt`（原始输出）。
- 观察项候选（v51 §五）：① 与 T026 `:121` 同源的 `unknown_token_count_multiplier` 警告在 BYOK 下再现，且措辞指向「canonical model ID」——CLI 层 model-id/name 区分的新数据点；② 凭据管理器 CredReadW 结构体偏移坑（`long` vs 双 `uint` FILETIME）。
- 未读任何密钥值；密钥未出现在 argv/输出/文件/对话中。
- 时点注记（§11.1 历史不改写）：`01-arch/layer0-prereq-2026-11-07.md` 中「密钥不存在」的旧事实已被 2026-11-08 的 `exists=true` 事实取代；原文不删，本行作时点注记。
- 零计费校验注记：`copilot --version` 经 `-File` 数组参数传入时参数绑定不可靠，且缺 `-p` 时进入交互态挂起（已终止进程，无计费）；清理契约读数取自首轮校验件，不受影响。
