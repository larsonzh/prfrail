# 07-final-review-r1 — ⑦ 独立终审第 1 轮（原文归档）

归档说明：本文件为该轮审查方最终消息的**逐字原文**。审查方 = `independent-reviewer`；调用模型 = `GPT-5.3-Codex (copilot)`；时点 2026-11-05。块末结论行为 `FINAL REVIEW: PASS WITH FIXES`。

```text
终审报告（P55 EXIT-SEED-BUILD-EVIDENCE，⑦ 独立终审）

**发现清单（按严重度）**
1. Low - 文档状态断言偏强（“已就地修正”与正文读数未完全一致）  
位置：[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md#L5)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report.md#L60)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report-r2.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report-r2.md#L98)  
说明：MANIFEST 写明 R2-1..R2-6“已就地修正”，但 03-test-report 仍保留“无 EXIT 字段=6”的旧计数；R2 报告已更正为 7。  
修法：二选一即可闭合。  
- 方案 A：在 MANIFEST 改成“R2-3 以勘误记录闭合（历史报告保留原值）”。  
- 方案 B：在 03-test-report 增补勘误注记，显式链接到 R2 更正。

**独立核验动作与读数**
1. 实读你指定文件：  
[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md#L1)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/anomalies.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/anomalies.md#L1)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report.md#L1)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report-r2.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report-r2.md#L1)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/index.json](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/index.json#L78)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-21-bindings.json](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-21-bindings.json#L1)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md#L1)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/SHA256SUMS.txt](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/SHA256SUMS.txt#L1)。
2. SHA256 抽查（满足“两目录各≥2 + 包根≥1”）：共 5 条，5/5 命中。  
样本含包根 1 条（claim-scope）、Windows 2 条、Linux 2 条。
3. 五枚夹具哈希独立复算并三方对照：  
- 实盘夹具 5 条  
- 对照 Linux bindings 5/5 一致（[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-21-bindings.json](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-21-bindings.json#L15)）  
- 对照 Windows bindings 5/5 一致（[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/bindings.json](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/bindings.json#L1)）  
- 对照 CI 输入表 5/5 一致（[docs/validation/s1-selfhost.md](docs/validation/s1-selfhost.md#L93)）。
4. 代码层核对：  
- overlap 在 seed 哈希前： [internal/release/verify.go](internal/release/verify.go#L94) 先于 [internal/release/verify.go](internal/release/verify.go#L110)。  
- VerifySelfHost 为只读路径（读取/哈希/解析）：[internal/release/verify.go](internal/release/verify.go#L82)。  
- probe flag 集与 02a-05 索引忠实命令一致：  
  [cmd/prfrail-release/main.go](cmd/prfrail-release/main.go#L41)、[cmd/prfrail-release/main.go](cmd/prfrail-release/main.go#L42)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/index.json](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/index.json#L80)。  
- 同时确认捕获原文的缩写命令已由勘误绑定：  
  [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-05-probe-actual.txt](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-05-probe-actual.txt#L3)、[docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/anomalies.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/anomalies.md#L22)。
5. 越界核对：暂存 53 项均为文档与证据包，未见 .go、未见 .github/workflows 路径变更。

**“证过了头”判定**
1. 技术断言层面：未发现超出 S-1..S-7 的新增事实断言。  
边界收敛（N-1..N-7）在 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md#L21) 明确，且与代码与捕获读数一致。
2. 元叙述层面：存在 1 处偏强状态断言（见 Low 发现），不影响 AT-13 核心证据成立，但建议修正文案精度。

**对 ③/⑤/⑥ 结论采信判定**
1. ③ 采信：支持。  
理由：AT-13 核心链条（分离、oracle 比较、失败不损 seed 的当前口径）与实读/复算一致。参考 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md#L9)。
2. ⑤ 采信：部分支持。  
理由：F-05-1/F-05-3 的勘误闭合证据充分（E-1/E-2），见 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/anomalies.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/anomalies.md#L16)；其余 ⑤条目在当前包内非完整原文链，仅能按现存证据部分采信。
3. ⑥ 采信：部分支持。  
理由：在其“扫描无发现”范围内可成立；但本轮独立终审新增 1 条文档一致性 Low（不推翻 AT-13 技术结论）。

**⑧阶段六项文档同步清单核对**
1. TEST_STRATEGY：适用，需同步 AT-13 状态/证据口径。  
文件：[docs/TEST_STRATEGY.md](docs/TEST_STRATEGY.md#L43)
2. OPERATIONS：可 N/A（若无运行流程、值守、发布运维动作变化）；需写明“仅证据补链，无运维流程变更”。  
文件：[docs/OPERATIONS.md](docs/OPERATIONS.md)
3. DEV_PLAN：适用，需同步本片完成记录与 R/D 收口。  
文件：[docs/DEV_PLAN.md](docs/DEV_PLAN.md)
4. S1_SUPPORT_MATRIX：可 N/A（若无平台/版本支持矩阵变化）；需写明“无支持范围增减”。  
文件：[docs/S1_SUPPORT_MATRIX.md](docs/S1_SUPPORT_MATRIX.md)
5. S1_RELEASE_NOTES：可 N/A（若无发布件/产品行为变更，仅证据补链）；需写明“无发行物与功能变更”。  
文件：[docs/S1_RELEASE_NOTES.md](docs/S1_RELEASE_NOTES.md)
6. REMAINING_SLICES：已同步（进行中条目存在）。  
文件：[docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L1612)、[docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L1608)

**载体备注回应（同一只读载体，不同模型）**
1. 已追加独立复核动作，不直接采信⑥结论。  
2. 追加动作包括：独立终端复算（样本 SHA、夹具三方一致性）、源码锚点核对（verify/main 与捕获命令对齐）、越界清单核对（53 项路径级审计）。  
3. 结论上保持“同载体下的二次独立性”：方法、取样与判定均由本轮重新执行。

**无法验证部分与推翻条件**
1. 无法验证部分：  
- ⑤与⑥的完整原始审查报告正文未在当前证据包中完整归档（当前可验证的是其在包内的落点与勘误）。  
- 本轮按只读边界未重跑 clean-room 构建与 VM 现场流程，只做证据与代码一致性复核。  
2. 推翻条件（任一成立即需重开⑦）：  
- 任一后续全量 SHA256 复算出现不一致。  
- ⑧文档同步出现超出 N-1..N-7 的新断言（尤其是“Windows 核验者独立性已证明”这类被 N-1 明确禁止的表述）。  
- 暂存集出现 .go 或 .github/workflows 路径改动而未单独授权。

FINAL REVIEW: PASS WITH FIXES
```
