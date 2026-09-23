# DR-1 切片验证报告：可用性探针的 CLI 额度下限

> 切片 `DR-1`（业务切片；含一次准则台账补登）| 分级 `[SLICE]` | 规模 S | 日期 2026-09-24。
> 状态：**已停止——未完成**（⑩ 停点；commit / push 须同轮授权）。

## 1. 变更概要（含执行自由度清单）

把 `prfrail ai check` 传给候选 CLI 的 `--max-ai-credits` 由"产品请求预算"改为 `max(30, maximumRequests)`，使产品可用性探针首次能在已 pin 候选上产出 `available` 记录。

- 实现：`internal/adapters/ai_probe_copilot.go` 新增常量 `copilotCLICreditFloor = 30`（`:23`），参数点由 `:111` 位移至 `:117`；`internal/console/ai.go` **零改动**。
- 测试：argv 断言由 `"1"` 改为 `"30"`；新增 2 用例（3 子例 + CLI 拒绝 fail-closed）。
- 文档与台账：v1.14 台账补登（OB-46–50）、ADR-014、DR-1 切口定义、链配置、证据包与自洽清单。
- 自由度清单：`docs/validation/evidence/DR-1-freedom-list.md`（F-1…F-6 与用户裁定类 D-1…D-6）。

## 2. 执行流水线

| 环 | 角色 / 模型 | 结果 | 额度 |
|---|---|---|---|
| ① 架构 | **无独立方案书**（方案基线 = 用户 R1–R5 裁定） | 见 §3 第 1 条 | — |
| ②a 协议先行 | 未触发（R1=A 冻结记录字段与 schema） | — | — |
| ②b 代码 | 主控 | 常量 + 参数点 | — |
| ③ 测试 | 主控 | 2 用例 + 变异检验 | — |
| ④ 门禁 | 主控 | 见 §9 | — |
| ⑤ 预审 / 复审 | `deepseek-v4-pro`（Max） | `PASS WITH FIXES` → 整改 → `PASS` | 2/2 |
| ⑥ 独立扫描 | `mai-code-1.1-flash`（盲扫） | `INDEPENDENT SCAN: FINDINGS` | 1/3 |
| ⑦ 终审 | `deepseek-v4-pro`（**降级**，ADR-014） | `FINAL REVIEW: PASS` | 1/1 |
| ⑨ 原生验证 | 主控 | 真探针 1 次 ⇒ `available` | 1/5 |

## 3. 异常与兜底（逐条）

1. **无独立 ① 方案书**（偏离 §1.2）：用户 R1–R5 裁定即方案基线，已在 `tmp/dr1/rulings.md` 与自由度清单 D-1 留痕。
2. **执行顺序偏离**（F-6）：先 ⑨ 真探针、后写本报告，避免报告出现 G2 所禁的占位形态；不改任何规则语义。
3. ⑤ Medium①（变异证据未随输入包）：补齐 `tmp/dr1/mutation-check.txt` 后闭合。
4. ⑤ Low②（§12.4 引用悬空）：按用户裁定改走 `docs/ADR_REGISTER{,_EN}.md` 的 ADR-014，台账指针同步修正。
5. ⑥ 层未宣告通过：按 `docs/DELIVERY_DIRECTIVE.md` §7.2 的替代路径转入 ⑦，由 ⑦ 结论闭环（用户裁定方案 b）。
6. ⑦ 两项 Low：§7.2 分支未逐字对应（留痕）、`DEV_PLAN{,_EN}` 与 `:111 → :117` 引用漂移（见 §13）。
7. **⑧b 回写脚本缺陷（自查发现）**：首版脚本对同一台账文件读取两次、后一次写入覆盖前一次 ⇒ 缺陷行改写与收尾行丢失；`git diff --numstat` 显示台账仅 `+17/-0` 时暴露，已用“断言旧文在且新文不在”的补丁脚本重放，并核实为 `+20/-1`；未影响代码、判据与结论。

## 4. 环境与档位

- 平台：Windows 本机、Go 1.22 工具链；pin = `copilot.exe` 1.0.83（sha256 `d3f3bb7b…0f671ee2`，探针前复核逐字节一致）。
- 档位：主控标准 High；⑤ 与 ⑦ = `deepseek-v4-pro`（Max，不占 Copilot 预算）；⑥ = `mai-code-1.1-flash`。
- **⑦ 降级**：本片 ⑦ 由非白名单模型承担，**独立性降低**；按 OB-37 元规则以 ADR-014 独立留痕（四要素 + 不得作为先例 + 三项标注），交 ⑧c 复核。
- **试点顺延第 3 次**：本片为业务切片验证、非准则治理 ⇒ 不作为试点；四要素：裁决人 = 用户；日期 = 2026-09-24；理由 = 业务切片验证；依据 = 用户同轮裁定原文；**不得作为先例**。

## 5. 审查结论摘要

- ⑤ 首轮 `PASS WITH FIXES`（0 Critical / 0 High / 1 Medium / 5 Low）；复审 `PASS`（六条逐条闭合，无新增偏离项）。
- ⑥ `INDEPENDENT SCAN: FINDINGS`（High×1 + Medium×1，主题 = CLI 有效额度 ≠ 产品预算且记录不含费用维度）。
- ⑦ `FINAL REVIEW: PASS`（四项全通过；2 Low + 2 Note，均属后续阶段义务）。
- 额度：⑤ 2/2、⑥ 1/3、⑦ 1/1、探针 1/5（余额 4 归还）。

## 6. §7.2 处置记录（⑥ 层 FINDINGS 的替代路径）

- ⑥ 的两条发现指向同一差异：CLI 侧有效额度 30 ≠ 产品请求预算 1，且记录 `requestsUsed` 只是请求计数、不含费用维度；其“最小修复”（把有效额度写入记录、或要求该情形返回 `unavailable`）与 **R1=A** 直接冲突。
- 处置（用户裁定方案 b）：**维持 R1=A，不修**；⑥ 层**不宣告通过**；按 §7.2 的替代路径转入 ⑦ 并由 `FINAL REVIEW: PASS` 闭环。
- 正向信号：⑤（带裁定）与 ⑥（盲扫）**独立同源**地指出同一条差异 ⇒ 该差异真实、明显，非我方虚构。

## 7. 审查输入包摘要（含盲扫隔离）

- ⑤ 与 ⑦ 输入：切口定义、R1–R5 裁定整理、原始 diff（`tmp/dr1/review-input-r2.diff`，9 文件）、门禁与聚焦测试原始输出、变异证据、链配置自证、契约原句；**不含**① 对话历史与 ②b 实现推理。
- ⑥ 首轮：**只给** diff、契约原句与只读约束，**不附任何他人清单**。
- 三个审查层均被硬约束"不得创建 / 修改 / 删除文件；不得声称执行 `node` / `go` / `gate` 命令，只可引用主控落盘路径"。

## 8. 可证伪性（机械检查 + 变异检验）

- 变异：断言式脚本要求"恰好 1 处"命中（否则抛错）⇒ 变异后 **3 处 FAIL**（`ai_probe_copilot_test.go:80`、`:318`、`:342`）⇒ 逐字节还原 ⇒ sha256 前后一致 `B805EA42…A4EB7` ⇒ 复绿。
- 未覆盖不变量（已登记、不并入本片）：`invalid FailureCode → unknown` 与 `maximumRequests < 1 → error` 两条既有分支无直接用例。

## 9. 门禁结果（⑩ 前终值）

- `gofmt -l` 空、`go build ./...` 干净、`go vet ./...` exit 0。
- 聚焦 `go test -count=1 -run CopilotCLIProbe ./internal/adapters/` ⇒ ok；全量 `go test -count=1 ./...` ⇒ 14 包全 ok、FAIL 0。
- `node tools/gates/gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0`；G4a 编码判据全绿（文件数由该判据自身输出）；G6-5 scanned 2 / refs 4 / 未解析 0；G6-6 scanned 10；G3(a) 四对（指令 `+7/-1`、ADR `+1/-0`、DEV_PLAN `+1/-1`、台账 `+20/-1`）；G1-b 头部与日志同为 v1.14。
- **镜像口径（如实）**：`DELIVERY_DIRECTIVE{,_EN}` 为**严格同位镜像**（1309 / 1309，逐行 0 失配）；`ADR_REGISTER{,_EN}` 与 `REMAINING_SLICES{,_EN}` **不做整文件位置镜像**（HEAD 即有 37 / 226 处位置性差异，非本片引入），对其适用 G3(a) 的“改动行数对称”与新增块形状一致（ADR-014 两行为 14 bold / 5 bars 全等；DR-1 定义块逐行形状全等）。
- **未执行项**：`gate.js --selftest` 本片未跑（未改动任何判据 ⇒ 无触发条件），将在提交前与 CI 执行。

## 10. ⑧c 收尾清洁记录

- 全仓逐文件编码：`tracked=2923 text=2834 binary=89 violations=3`，3 处均为**冻结证据包**（OB-27 临时豁免：`b2-2026-09-17` 的 CRLF 与 `b3b-2026-09-21` 的两处 BOM）；本片新增文件 0 违规。
- 证据包自洽性：`packs=5 entries=2331 missing=0 mismatch=0 not-self-checkable=1`；本片包由不可自核变为 `entries=4 missing=0 mismatch=0 -> PASS`。
- `tmp/` **未清理**：按用户既有指令（交给 `DIRECTIVE-LEDGER-ARCHIVE` 与 OB-49），本片草稿与原始捕获保留并以 `local-only` 声明。

## 11. 已知差异与候选 OB 预留

- **已知差异（R1=A，不修）**：契约 `docs/CONTRACTS.md:191` 的“在显式费用/请求预算内”与实际 CLI 允许 30 额度之间存在费用维度的差异，记录无该维度；已在代码注释（`:19-22`）、切口定义边界行与本报告三处标注。
- **候选 OB-52**（下一片登记）：该差异（合并 ⑥ 的 High 与 Medium）**并合** ⑥ 的 Low④（CLI 下限变更 ⇒ 收敛为 `unknown` 且仅存哈希、无机械发现性）——同属一条差异的两面。
- **候选 OB-51**（保留）：BYOK 与 secretStorage 五点只读确认，留待 DR-1 收尾后处理；若日后 `DIRECTIVE-MODEL-LISTS` 立项需编号，另起 **OB-53**。

## 12. 成本与计量

- 付费模型调用合计 **4 次**（⑤ 2 + ⑥ 1 + ⑦ 1）；其中 ⑥ 走 Copilot 预算。
- 探针 **1/5**（目标 1，余额 4 归还）；成功探针按最坏情形记 **1 个 Copilot premium**；**"BYOK 归属未证"如实标注**（本片未走 BYOK）。
- ① 无调用（无独立方案书）；⑧ 与 ⑨ 由主控执行，0 次付费。

## 13. 未执行事项与下一步建议

- 未执行：commit 与 push（须同轮授权）；`gate.js --selftest`（提交前与 CI）；⑥ 的逐字条目未给 ⑦（刻意隔离）；BYOK 与 secretStorage 五点确认（收尾后）。
- 建议：`DEV_PLAN{,_EN}` 的 B1 段回写与 `:111 → :117` 引用修正已在本片完成；下一片登记 OB-52 与 OB-51（若 BYOK 项仍待办）。

## 14. 产物清单

```artifacts
internal/adapters/ai_probe_copilot.go	repo
internal/adapters/ai_probe_copilot_test.go	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/validation/dr-1-availability-probe-floor.md	repo
docs/validation/dr-1-availability-probe-floor_EN.md	repo
docs/validation/evidence/DR-1-freedom-list.md	repo
docs/validation/evidence/dr-1-2026-09-24/proofrail.chain.json	repo
docs/validation/evidence/dr-1-2026-09-24/availability-1083-available.json	repo
docs/validation/evidence/dr-1-2026-09-24/probe-01-ai-check.json	repo
docs/validation/evidence/dr-1-2026-09-24/probe-01-ai-verify.json	repo
docs/validation/evidence/dr-1-2026-09-24/README.md	repo
docs/validation/evidence/dr-1-2026-09-24/SHA256SUMS.txt	repo
tmp/dr1/**	local-only
```

**`tmp/dr1/**` 的处置说明**：包含本片审查输入包、变异检验转录、链配置自证、探针原始捕获与草稿；按用户既有指令**不清理**，并留给 `DIRECTIVE-LEDGER-ARCHIVE`（双主题：台账归档与 `tmp/` 生命周期）处置。

## 15. 附录 A：⑨ 原生验证记录与对照

| 项 | B1 修复前基线 | 本片修复后 |
|---|---|---|
| 记录文件 | `docs/validation/evidence/b1-20260917/availability-1083-unknown.json` | `docs/validation/evidence/dr-1-2026-09-24/availability-1083-available.json` |
| `status` | `unknown` | `available` |
| `reason` | `unknown` | 缺失（`available` 禁止 `reason`） |
| `requestsUsed` | 0 | 1 |
| 退出码 | 1 | 0 |
| `profileConfigHash` | `sha256:d529c357…e436` | 同值（逐字一致） |
| `profileId` / `channel` | `copilot-cli-host` / `agent-runner-cli` | 同值 |

- 探针附加证据：`ai verify` 对同一记录 `ok=true`、`exitCode=0`（freshness 与预算检查通过，零模型调用）。
- 记录文件 745 字节、SHA-256 `2DC3EDBD…4418`；pin 二进制 sha256 复核 `D3F3BB7B…1EE2`。

## 16. 附录 B：⑤ / ⑥ / ⑦ 结论行与 ⑥ 的两条发现原文

- ⑤ 首轮：`PRE-REVIEW: PASS WITH FIXES`；复审：`PRE-REVIEW: PASS`。
- ⑥：`INDEPENDENT SCAN: FINDINGS`。
- ⑦：`FINAL REVIEW: PASS`。
- ⑥ 的发现原文（逐字，未改写）：

`High | internal/adapters/ai_probe_copilot.go:19-23 | 把候选 CLI 的平台最小 credits（30）强制覆盖原始 maximumRequests，但没有把 effectiveProbeCredits 与 declared budget 一并写进 record；available 可在 maximumRequests=1 时出现，违反"显式费用/请求预算内"的契约。 | 把 effectiveProbeCredits 与 declared budget 一并写进 record`

`Medium | internal/adapters/ai_probe_copilot_test.go:292-342 | 测试只断言 RequestsUsed == 1 和 Status == "available"，没有模拟真实 credit floor 与 declared budget 的冲突，也没有验证 effectiveProbeCredits <= budget | 增加一条"maximumRequests=1 且真实 effective credits=30 时必须返回 unavailable"的反例`
