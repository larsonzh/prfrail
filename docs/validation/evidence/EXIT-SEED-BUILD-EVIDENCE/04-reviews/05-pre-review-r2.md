# 05-pre-review-r2 — ⑤ 预审第 2 轮（整改后复审，原文归档）

归档说明：本文件为该轮审查方最终消息的**逐字原文**。审查方 = `deep-reasoner`；调用模型 = `DeepSeek V4 Pro (deepseek)`；时点 2026-11-05。

```text
# P55 ⑤ 预审 · 第 2 轮复审报告（只读）

## 一、四项整改逐条核验

### F-05-1（Medium）— **闭合** ✅
- `anomalies.md`「记录保真度勘误」节存在，E-1 行三列齐备（位置 / 不一致事实 / 真实情况）。
- `02-win-captures/index.json` 的 `02a-05`：`command` 已是忠实命令（含 `--noop-chain testdata\selfhost\noop-chain.json --executable-chain testdata\selfhost\executable-chain.json --work-root … --output … --expected-version 0.1.0-dev`），且 `"commandErrata": "E-1 (see anomalies.md)"` 存在。
- 保真纪律成立：捕获件 `02a-05-probe-actual.txt:3` 的 `COMMAND:` 仍为缩写（缺两链），磁盘 SHA256 = `28d24043…` 与 `SHA256SUMS.txt` 命中（未改写）。
- 忠实命令与实现签名一致：`cmd/prfrail-release/main.go:38-49` 的 flag 集（`--binary/--noop-chain/--executable-chain/--work-root/--output/--expected-version`）逐参数吻合；`internal/release/probe.go:47` = `existingRealPath(plan.NoopChainPath)`，`verify.go:214` 对空串报 `empty path` ⇒ 缩写命令照抄执行确会 exit 1，E-1 的论证链成立；且与 `replay.md:17-20` 的 Windows 命令、`02b-13` Linux 捕获逐参数同形。
- 微瑕（见新发现 N-2）：E-1 引文行号 `:47/49` 中 `:47` 精确，`:49` 实为 noop 报错 return，executable 的 `existingRealPath` 在第 51 行。

### F-05-2（Low）— **闭合** ✅
`input-list.md` §6.3「缺口编号表」G-1 至 G-8 完整存在，逐条给出闭合捕获位置；§6.4 保留原四条方向为历史记录。哈希已重算落盘（`98685946…` 命中）。

### F-05-3（Low）— **部分闭合** ⚠️
- E-2 存在且「真实情况」列有据：`02b-00-env.txt:16-19` 实测含 `GOPROXY_preload=https://goproxy.cn,direct`、`GOPROXY_build=off`、`GOSUMDB=off` 与双 note 行；`02b-05` 的 `TITLE:` 确有 `(GOPROXY=direct)` 残留。
- 但 E-2 声称 `02b-05/06/07` 三件残留，实测 `02b-06`（`…verifier worktree`）、`02b-07`（`…candidate worktree`）的 TITLE **无**该残留 —— 勘误范围过度包含（见 N-1）。

### F-05-4（Low）— **部分闭合（按预期）** ⏸
计数统一未在本轮做：`03-test-report-r2.md` 正文读数 48 保持哈希锁定的历史时点；但 `03-test-report-r2.md` 已提前补录入 `MANIFEST.md`/`SHA256SUMS.txt`（现口径 49），且主控已声明留待 ⑧b/⑨ 统一。⑧b/⑨ 引用 R2 读数时须注明现口径 49。

## 二、证据链完整性核验 — 通过 ✅
- 三方集合相等实测：`SHA256SUMS` 49 行 = `MANIFEST` 表 49 条 = 磁盘载荷 49 件；`manifest_only=0`、`sums_vs_disk_diff=0`。
- 哈希抽查 **8 条全命中**（超出要求的 4 条）：win 侧 `02a-05.txt`、`index.json`、`anomalies.md`、`input-list.md`；linux 侧 `02b-13.txt`、`02b-21-bindings.json`；根 `03-test-report.md`、`03-test-report-r2.md`。特别核验：本轮改写的 `anomalies.md`/`input-list.md`/`win index.json` 哈希均已同步重算。
- 编码合规：3 个 `.md` BOM=True+LF，2 个 `index.json` 与 `SHA256SUMS.txt` 无 BOM。
- `commandErrata` 全量 39 条中仅 `02a-05` 一条非 null，与 E-1 范围精确对应。
- 无新矛盾：index `command`（忠实）与捕获 `COMMAND:`（缩写）的不一致已被 E-1 + `commandErrata` 明确解释；MANIFEST「三方相等」宣称实测成立。

## 三、新发现（均可 ⑧b 前一行字修复）
| 编号 | 严重度 | 位置 | 修法 |
|---|---|---|---|
| N-1 | Low | `anomalies.md` E-2 行 | 「`02b-05/06/07` 的 TITLE」改为「`02b-05` 的 TITLE」（可加注 06/07 无残留）。**建议本轮修**：`anomalies.md` 已规定「⑧b 验证报告须引用本表 E-1/E-2」，若 ⑧b 照抄会将「06/07 有残留」这一错误事实带入验证报告 |
| N-2 | Low | `anomalies.md` E-1 行 | 引文「probe.go:47/49」改为「probe.go:47/51」或「:47（noop）/ :51（executable）」；`:49` 实为 noop 报错 return 行。语义正确，可不修 |

## 四、继承 Low 处置判定
| 编号 | 判定 | 理由 |
|---|---|---|
| F-08 | 接受为观察 | N-1/S-1 已收窄断言；`work/` 已删，本片边界内不可补 |
| F-10 | 接受为观察 | 叙述型残留不阻断；建议 ⑧b 引用 `02a-12`/`02b-20` 同文件的目录清单作旁证 |
| F-11 | 接受为观察 | 7 件无 `EXIT:` 行属捕获格式差异；计数更正已由 R2-3 落盘 |
| F-12 | 接受为观察 | VM 已清理不可补；S-7/A-5 已把结论限定在「构建期 GOPROXY=off」 |
| R2-3 | 已闭合 | 计数 6→7 更正已登记于 `03-test-report-r2.md` |

顺带确认：R2-4（A-1/A-2 复核状态措辞）主控已顺带修复——`anomalies.md` 两行现均为「**部分可核验**」，与 A-6/A-7 口径对齐。

## 五、无法验证事项与推翻条件
1. 主控的三个脚本读数（`payloads=49 meta=2`、win 16/linux 23、`changed=53 TOTAL_FAIL=0`）我未重跑其脚本本身（位于 `tmp/`，按授权未读），但以磁盘实况独立复算了关键子集（集合相等、8 条哈希、BOM、commandErrata 全量），与主控声明一致；若 ⑧b 重跑 gate 出现 FAIL，即推翻本报告第二部分。
2. 构建 / 探测未重跑（授权边界外）；`02a-05` 忠实命令的**端到端可执行性**依赖 `replay.md` 正确性，而 `replay.md` 已与 main.go 签名逐参数核对一致。
3. A-6 / A-7 / F-12 等远端不可复核项维持原状，若 VM 恢复或留存新证据方可重开。

---

**PRE-REVIEW: PASS WITH FIXES**（需修项：N-1 必修 + N-2 建议修，均 Low 一行字；其余 Low 接受为观察，F-05-4 计数统一按声明留待 ⑧b/⑨）
```
