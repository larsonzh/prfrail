# P84 · 最终收窄复审请求与 OB-37 登记（2026-11-21）

## ⑦ 超额度登记（OB-37 元规则，独立留痕）

| 要素 | 内容 |
|---|---|
| 裁决人 | 用户（2026-11-21 授权消息「P84 最终收窄复审 · 授权消息」直接裁定） |
| 日期 | 2026-11-21 |
| 理由 | 前两轮收窄复审后仍余两项 Minor 需独立核验；本片 ⑦ 既定额度（盲审 1 + 终审 1 + 复审 3 + 窄域例外 1）已全部用尽 |
| 依据 | 用户同轮裁定原文（仅限 Minor 1 与 Minor 2 两项，其余一律不重开） |

**性质标注**：超出 §2.2 字面标准 / 本片既定 ⑦ 额度，系用户直接裁定；**不得作为先例**
（其它切片仍按 §7.5 与 §2.2 白名单额度执行）。
**后续义务**：P84 报告如实标注「本片 ⑦ 追加轮次超出既有额度、系用户直接裁定、独立性结论
按本次载体记载」；本轮结论不推广。

## 复审基线（冻结）
- diff：`docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1r2-final-diff-r3.txt`
  （SHA256 = `669fe92141cf8c851c483af716701677b68699cae6a2fba74d8d29c1fa6b7fdf`）。
- HEAD = origin = `d597f1e`；仅两文件暂存（实现 +454/−2、测试 +509/−0）。
- 只读复审：不改代码、不重跑全量测试。

## 范围（仅两项）
1. Minor 1：`seedIsFixedVolumeRoot`（仅凭已打开种子句柄取证：GetVolumeInformationByHandleW
   答 NTFS + GetFileInformationByHandle 低 32 位文件索引 == 5；driveTypeOf 不参与判定）。
2. Minor 2：生产 seam `directoryEntryShortNameFn` + 判别用例
   `TestP84W12AliasBranchRefusesMissingEntryEvidence`（M7 判别）与旧 M3 标注 HISTORICAL 更正。

## 口径
- 即使总评 0 Critical / 0 Major，仍逐条核 Minor 是否达 Medium+（§4.3：Medium+ 即阻断）；
  出现任何 Medium+ 必须停下报用户，不得自行加轮。
- 复审报告落盘 `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/13-final-narrow-review.md`，附每条 FINDING 的严重度、
  判定依据、结论（放行/阻断）。

---

# 追加登记：方案 A 窄复审（2026-11-22）

## ⑦ 超额度登记（OB-37 元规则，追加条目，独立留痕）

| 要素 | 内容 |
|---|---|
| 裁决人 | 用户（2026-11-21 裁定「P84 · 补正与窄复审授权」直接指定载体 independent-reviewer + GPT-6.2 Sol Max） |
| 日期 | 2026-11-22 |
| 理由 | 上轮终审阻断后按用户裁定实施方案 A（48 位记录号），需独立核验；本片 ⑦ 额度此前已用尽（见上条登记），本轮再次超出 |
| 依据 | 用户裁定原文（仅方案 A 范围，不重开 Minor 2 seam、不重开上轮 Major、不重开其它任何项） |

**性质标注**：再次超出 §2.2 字面标准 / 本片既定 ⑦ 额度，系用户直接裁定；**不得作为先例**
（其它切片仍按 §7.5 与 §2.2 白名单额度执行）。
**后续义务**：P84 报告如实标注「本片 ⑦ 再次追加轮次超出既有额度、系用户直接裁定、独立性结论
按本次载体记载」；本轮结论不推广。

## 复审基线（冻结）
- diff：`docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1r2-final-diff-r5.txt`（51,291 字节）
  SHA256 = `88e866260f7cd03c82e41da2fe82f05aff24f93ac6fbc77d7c54b0f91c98558c`
  （与 r4 逐字节一致，同内容重新冻结；候选 `c6429341…` 不作基线）。
- HEAD = origin = `d597f1e`；仅两文件暂存（实现 +471/−2、测试 +543/−0）。
- 只读复审：不改代码、不重跑全量测试。

## 范围（仅方案 A）
- `ntfsRootRecordOf` 纯函数、`seedIsFixedVolumeRoot` 调用点、`ntfsRootDirectoryFileIndex`
  常量注释 48 位语义、`TestP84W12NtfsRootRecordOf` 判别用例、
  `TestP84W12SeedIsFixedVolumeRoot` 两行恢复、两处 Info、M9 变异证据。
- **不重开**：Minor 2（seam）、上轮 Major、其它任何项。

## 口径（不变）
- 任何 Medium+ 即阻断 ⇒ 停下报用户，不得自行加轮、不得自行修复。
- 报告落盘 `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/15-final-review-A.md`，附每条 FINDING 严重度、判定依据、结论。

## 载体追认（2026-11-22 用户裁定，不新增模型）
- 裁定署名「GPT-6.2 Sol Max」在可用模型列表中不存在；本轮窄复审实际执行载体为
  `GPT-6.2 Sol (copilot)`（与前几轮复审同一载体），结论已落盘 `15-final-review-A.md`
  （RELEASE，0 Critical / 0 Major / 0 Medium，2 Info）。
- 最终载体如实追认：**GPT-6.2 Sol**；不新增模型、不新增额度（超额登记见上）。
