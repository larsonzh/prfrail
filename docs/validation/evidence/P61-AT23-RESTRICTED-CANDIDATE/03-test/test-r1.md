# ③ 测试回执（层 1 骨架 5 项）— P61 AT23-RESTRICTED-CANDIDATE

- 日期：2026-11-08
- 角色：③ 测试/探针（`prfrail-tester`，主控指定 `DeepSeek V4.1 Flash`）
- 范围：`01-arch/design-r1.md` §C 层 1 五项；层 2 未启动（按本单要求）
- 付费：项 4（第 3 笔）、项 5（第 4 笔）各 1 次，无重试、无追加

## 逐项结果

| # | 命令摘要 | 退出码 | 判定 | 证据 |
|---|---|---|---|---|
| 1 | dot-source `tools/agent-probe/CopilotCli.ps1` → `Invoke-CopilotCli -Subcommand version`（argv `--no-auto-update version`） | 0 | **PASS**：stdout 含 `GitHub Copilot CLI 1.0.83` | `03-test/layer1/l1-01-wrapper-version.txt` |
| 2 | 裸跑 `& <钉版exe> --no-auto-update --version` | 0 | **PASS**：版本 token `1.0.83` 与项 1 逐字一致（全行仅差句末句点，已登记） | `03-test/layer1/l1-02-bare-version.txt` |
| 3 | `Invoke-CopilotCli -Subcommand version -Unlock`（argv `version`，无 `--no-auto-update`） | 0 | **PASS**：观测 `1.0.85` ∈ {1.0.83, 1.0.85}，与预期一致 | `03-test/layer1/l1-03-unlock-version.txt` |
| 4 | `Invoke-ByokLayer0.ps1 -InvokeWrapper -WrapperPrompt 'say hello'`（封装 argv `--no-auto-update -p say hello`） | 0 | **PASS**：回复非空 + `BYOK env cleanup: verified (all six cleared)` | `03-test/layer1/l1-04-byok-call.txt` |
| 5 | 同上 + `-WrapperEffort 'high'`（封装 argv `--no-auto-update -p say hello --effort=high`） | 0 | **PASS（三态判定 = 支持）**：exit 0 且回复非空 ⇒ `effort=high supported` | `03-test/layer1/l1-05-effort-probe.txt` |

## 观察项（不阻断）

1. **版本串标点随调用形态变化**：子命令形态输出 `GitHub Copilot CLI 1.0.83`（无句点）；标志形态 `--version` 输出 `GitHub Copilot CLI 1.0.83.`（有句点）。版本 token 一致，判据取 token。
2. **更新提示**：CLI 提示 `Update available: 1.0.107`；`--no-auto-update` 下仍固定 `1.0.83`，符合 R-08.1 预期。
3. **canonical model ID 边界**：BYOK 两次调用均告警 `No token count multiplier is known for Anthropic model "deepseek-flash" ... configure a recognized canonical model ID`（与层 0 同一观察）。
4. **项 5 的「支持」级差**：本层只证明 CLI 接受 `--effort=high` 且调用正常完成；BYOK provider 是否真正应用该档位在本层**不可测**（无该 provider 的生效可观测面）。按纪律未追加 low/medium/xhigh/max（超出批准额度），层 3 固定 `high`。
5. **控制台捕获乱码**：项 4/5 输出中的 emoji/中文在 PowerShell 5.1 捕获下呈乱码，属代码页捕获现象，非调用失败；证据文件已标注。

## 层 1 闸门结论

- **全绿**：5/5 PASS，无红项 ⇒ 层 1 闸门通过，可进层 2（本单不启动）。

## 未做事项

- 未跑层 2 / 层 3。
- 未 `git commit` / 未 `git push`。
- 付费仅 2 笔（项 4、项 5），未超审批额度，未重试、未追加 `--effort` 其他档位。
- 未打印任何密钥值（凭据由 tmp 配套件从 Windows 凭据管理器读取，命令行与全部证据文件不含密钥）。
- 未修改任何生产/工具脚本；未注册 `MANIFEST.md` / `SHA256SUMS.txt`（主控负责）。

TEST: DONE (layer1 green)
