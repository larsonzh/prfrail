# 切片验证报告 · `DIRECTIVE-TMP-LIFECYCLE`（v1.23）

> 类型：准则治理类（`tmp/` 生命周期执行：归档 + 引用加注记 + 台账登记）。分级 `[SLICE]`。规模 M。探针 0。
> **独立性标注**：本片 ⑦ 由 `deepseek-v4-pro` 承担（ADR-023，用户直接裁定），与 ⑤ 预审同模型族 ⇒ **独立性降低**，如实标注并交 ⑧c 复核。

## 1. 变更概要

- 按 §1.5.1 的 ② 类口径把 `tmp/` 的存量归入 15 个归档单元，落 `docs/validation/evidence/directive-history/`；每单元含 `MANIFEST.md`（归档条目 / 排除清单 / 未入库文件清单 / 重建说明 / 对应引用点 / 后缀与编码映射）与 `SHA256SUMS.txt`，载荷根 = `archive/`。
- 引用加注记：准则内三处既有注记之后追加 v1.23 注记（指向 `.txt` 实际路径）；9 对验证报告各追加 1 行归档注记；历史措辞一律不改（§11.1）。
- 台账与版本包：v1.23 版本包（头版本号 / §0.3 新行 / 附录 C.0m）、§12.4 登记 OB-76 / OB-77 / OB-78 / OB-79 / OB-80、`ADR-023`、完成行（只记报告路径）。
- 同批顺带：`DIRECTIVE-MODELID-MAP-mutation.md` 的注记措辞改「**重建归档**」（用户明令）；`ledger-backfill{,_EN}` 的过程文件计数补注为 21。
- 自由度清单：`docs/validation/evidence/DIRECTIVE-TMP-LIFECYCLE-freedom-list.md`（F-1 至 F-13）。

## 2. ① 判定与口径

- **归档口径** = 证据包：被受版本控制文档引用其路径者 ∪ 切片级文档（切片定义 / 方案 / 计划 / 评审包 / 门禁读数 / 冻结补丁）；运行时残留与一次性脚本只入排除清单。
- `dce` 与 `assembly-min` **按 ② 定义归档**（§1.5.1 的 ③ 例示与定义不一致，已登记 OB-78）。
- **二进制与忽略类型不入库**：只在单元 MANIFEST 记路径 + 字节 + sha256 + 重建说明；被引用的二进制路径在报告注记中指回 MANIFEST，不留死引用。
- **后缀与编码**（方案 3）：归档副本的 `.md` 一律改存 `.txt`（避开只扫 `.md` 的文档判据），并按归档后类型归一编码——`.ps1` 带 BOM、其余无 BOM、一律 LF；每单元逐条记录「原路径 / 归档路径 / 原始字节哈希 / 归档字节哈希 / 变更种类」。
- **边界拒绝**：`dr-fix-enforcement-proxy{,_EN}` 不加注记（该报告无 `artifacts` 围栏块，改动即触 G7-a）；证据树内的既有冻结工件不改动。

## 3. 交付面（归档单元）

| 单元 | 归档条目 | 排除条目 | 未入库条目 |
|---|---|---|---|
| `assembly-min` | 12 / 659,581 B | 72 / 0 B | 0 / 0 B |
| `b4` | 50 / 396,972 B | 300 / 0 B | 3 / 0 B |
| `dce` | 11 / 215,448 B | 6 / 0 B | 0 / 0 B |
| `dr1` | 14 / 137,253 B | 17 / 0 B | 0 / 0 B |
| `drfix` | 19 / 335,863 B | 19 / 0 B | 0 / 0 B |
| `gate-js` | 1 / 8,931 B | 0 / 0 B | 0 / 0 B |
| `gates-ext` | 31 / 671,077 B | 769 / 0 B | 2 / 0 B |
| `lbf` | 12 / 273,553 B | 9 / 0 B | 0 / 0 B |
| `master-subagent-architecture` | 1 / 26,124 B | 0 / 0 B | 0 / 0 B |
| `pre-directive-package` | 1 / 15,483 B | 0 / 0 B | 0 / 0 B |
| `reviews` | 6 / 100,926 B | 0 / 0 B | 0 / 0 B |
| `s1-reach` | 8 / 177,101 B | 6 / 0 B | 0 / 0 B |
| `tmp-lifecycle-workdir` | 15 / 744,747 B | 41 / 0 B | 0 / 0 B |
| `v112` | 22 / 450,619 B | 18 / 0 B | 0 / 0 B |
| `v113` | 76 / 2,602,143 B | 70 / 0 B | 0 / 0 B |
| **合计** | **279 / 6,815,821 B** | **1327 / 0 B** | **5 / 0 B** |

## 4. 执行流水线与环境

① 判定（主控 + V4 Pro）⇒ ②③ 归档与引用加注记（主控脚本，确定性、可复算）⇒ ④ 门禁 ⇒ ⑤ 预审（V4 Pro）⇒ ⑥ 独立扫描（MAI）⇒ ⑦ 终审（V4 Pro，ADR-023）⇒ 整改 ⇒ ⑧ 本报告对 ⇒ ⑨ 原生验证 ⇒ ⑧c 收尾清洁 ⇒ ⑩ 停点。

环境：Windows 11 / go1.27.0 / Node 24；仓库 `github.com/larsonzh/prfrail`，分支 `main`。

## 5. 审查结论

- **⑤ 预审**（V4 Pro）：`PASS WITH FIXES`。1 High = 报告对尚未落盘（本批同批落地即闭合）；1 Medium = OB-79 用户已批但无台账行（已登记）；2 Low = `tmp/tl` 漂移与 F-6 待追认。
- **⑥ 独立扫描**（MAI）：`PASS`。读数：哈希 `279 / 0`、覆盖 `279 / 0`、被引归档路径 `8 / 0`、忽略类型入库 `0`。
- **⑦ 终审**（V4 Pro，ADR-023）：`PASS WITH FIXES`。3 Low 已整改：三处「OB-76/77/78」计数补入 OB-79、单元 MANIFEST 的源路径与重建说明措辞、报告注记加「此前切片已清除路径不在本次范围」括注。
- ⑦ 的反例审查 6 条中 2 条成立（`dml` 既有死链、OB 台账漏登），均已在整改中闭合；其余 4 条实测不成立。

## 6. 可证伪性与门禁结果

| 项 | 读数 |
|---|---|
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0`（changed=337） |
| `gate.js --all --scope=index` | `TOTAL_FAIL=0`（changed=337） |
| `--selftest` | `SELFTEST: PASS 264/264` |
| 归档自洽 | 15 单元 / 载荷 279 条；`SHA256SUMS` 逐条复算失配 0、未列 0、缺失 0 |
| 归档保真 | `.md` → `.txt` 副本归一（去 BOM + LF）后与 `tmp/` 原件逐字节全等，字符内容 0 差异 |
| G4b | `newCjk=504 unknown=[]` |
| G7-b | `316 artifact row(s) OK` |
| 三表完备 | `b4` 353 = 50 + 300 + 3；`gates-ext` 802 = 31 + 769 + 2 |
| 变异检验 | 不适用（文档与归档片，无代码语义）；反证手段 = ⑥ 的六项机械核查与 ⑦ 的六条反例审查 |

**已知偏离（如实登记）**：归档副本的 BOM 形态按归档后类型归一（`.ps1` 带 BOM、其余无 BOM），与用户字面「BOM 形态保留」不一致 ⇒ 自由度清单 F-5 与 F-6 **已由用户追认**（2026-10-04），并登记 **OB-80**（归档载荷的编码形态口径缺失）。

## 7. 成本与计量

① 一轮、⑤ 一轮、⑥ 一轮、⑦ 一轮、⑧ 一轮（CN 权威稿与 `_EN` 镜像由主控一次生成）；探针 0；无 Haiku 调用；**未发生额度例外**（⑦ 降级以 ADR-023 留痕）。

## 8. 明确未执行事项

- 未改任何判据本体（未为证据树加判据作用域排除）⇒ 登记 OB-79，留待专门片。
- 未清任何 `tmp/` 文件（③ 类清理在 ⑧c 执行，见附录）。
- 未提交、未推送、未推 gitee。

## 9. 下一步建议（含档位建议）

- **下一片（建议）**：判据作用域与证据树（OB-79）——把证据树从文档判据的作用域中排除，改由冻结包自洽性核验承担；属判据本体变更 ⇒ 规模 **M+**、⑦ 用 `gpt-5.3-codex`，**不降级**。
- 若同批处理 OB-76（重建归档第四类处置）、OB-78（§1.5.1 例示修正）与 OB-80（归档载荷编码形态），可与上项合并为一片。

## 10. ⑩ 停点报告（含回写核对行）

- 执行至 ⑨ 与 ⑧c；**停在提交授权之前**：本片不 commit、不 push、不推 gitee。
- **回写核对**（§11.2 的 OB-63 口径）：验证报告 = **已回写**（本文件）；`DEV_PLAN` = **不适用**（准则治理片）；台账 = **已回写**（`REMAINING_SLICES{,_EN}` 完成行，只记报告路径）；准则 = **已回写**（§0.3 / §12.4 / 附录 C.0m / 三处注记）；ADR = **已回写**（ADR-023）。
- **⑨ 读数**：tree 与 index 均 `TOTAL_FAIL=0`；`SELFTEST: PASS 264/264`；归档自洽复算失配 0；⑧c 清理后非归档面复跑全绿。
- **提交与 CI**：commit `4284e94`；CI run `37209843851` 两腿全绿（`Go ubuntu-latest` / `Go windows-latest`），三条 `workflow_dispatch` 腿按事件条件跳过；两腿读数一致：`SELFTEST: PASS 264/264`、`GATE REPORT scope=ci changed=340`、`TOTAL_FAIL=0`。

## 11. 产物清单

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-tmp-lifecycle.md	repo
docs/validation/directive-tmp-lifecycle_EN.md	repo
docs/validation/evidence/DIRECTIVE-TMP-LIFECYCLE-freedom-list.md	repo
docs/validation/evidence/DIRECTIVE-MODELID-MAP-mutation.md	repo
docs/validation/ledger-backfill.md	repo
docs/validation/ledger-backfill_EN.md	repo
docs/validation/directive-cleanup-eval.md	repo
docs/validation/directive-cleanup-eval_EN.md	repo
docs/validation/directive-review-saturation.md	repo
docs/validation/directive-review-saturation_EN.md	repo
docs/validation/directive-v1.12.md	repo
docs/validation/directive-v1.12_EN.md	repo
docs/validation/dr-1-availability-probe-floor.md	repo
docs/validation/dr-1-availability-probe-floor_EN.md	repo
docs/validation/gates-ext.md	repo
docs/validation/gates-ext_EN.md	repo
docs/validation/s1-reach-reachability.md	repo
docs/validation/s1-reach-reachability_EN.md	repo
docs/validation/t027-assembly-min.md	repo
docs/validation/t027-assembly-min_EN.md	repo
docs/validation/t027-at23-e2e.md	repo
docs/validation/t027-at23-e2e_EN.md	repo
docs/validation/evidence/directive-history/tmp-lifecycle-workdir/CLEARED-INVENTORY.txt	repo
docs/validation/evidence/directive-history/assembly-min/MANIFEST.md	repo
docs/validation/evidence/directive-history/b4/MANIFEST.md	repo
docs/validation/evidence/directive-history/dce/MANIFEST.md	repo
docs/validation/evidence/directive-history/dr1/MANIFEST.md	repo
docs/validation/evidence/directive-history/drfix/MANIFEST.md	repo
docs/validation/evidence/directive-history/gate-js/MANIFEST.md	repo
docs/validation/evidence/directive-history/gates-ext/MANIFEST.md	repo
docs/validation/evidence/directive-history/lbf/MANIFEST.md	repo
docs/validation/evidence/directive-history/master-subagent-architecture/MANIFEST.md	repo
docs/validation/evidence/directive-history/pre-directive-package/MANIFEST.md	repo
docs/validation/evidence/directive-history/reviews/MANIFEST.md	repo
docs/validation/evidence/directive-history/s1-reach/MANIFEST.md	repo
docs/validation/evidence/directive-history/tmp-lifecycle-workdir/MANIFEST.md	repo
docs/validation/evidence/directive-history/v112/MANIFEST.md	repo
docs/validation/evidence/directive-history/v113/MANIFEST.md	repo
```

说明：归档单元以**单元 MANIFEST 行**声明（载荷逐条见各单元 `SHA256SUMS.txt`，合计 279 条）；`local-only` 类不另行声明（③ 类清理已在 ⑧c 执行，清单见附录）。

## 12. 附录 · ③ 类清理记录

- **清理时点清单工件**：`docs/validation/evidence/directive-history/tmp-lifecycle-workdir/CLEARED-INVENTORY.txt`——逐条记录路径、字节、sha256 与未归档理由（清理时点全量，含 MANIFEST 冻结后新增的过程文件）。
- 清理范围 = 本片排除清单（1327 条）＋ 本片工作目录读数之外的过程文件；清理后 `tmp/` 顶层只剩 `.gitkeep`。
- 清理后复跑 `--scope=index` 与 `--scope=tree`（读数见 §10）；未被引用的过程产物已按 §1.5 的「用完立即清除」归零。
- **脚注（清单口径）**：收尾期（③ 类清理时点之后）产生的临时过程文件由主控在每步执行后即时清除，**不计入** `CLEARED-INVENTORY.txt` 的清理时点集合。
