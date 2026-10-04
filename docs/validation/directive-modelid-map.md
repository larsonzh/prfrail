# 切片验证报告 · `DIRECTIVE-MODELID-MAP`（`modelId` 口径反转 + 完成行重建 + `tmp/` 生命周期机制）

> 依据 `docs/DELIVERY_DIRECTIVE.md` v1.22 的 §5.1 串行流水线与 §6.2 证据落盘义务。**CN 为权威版**，`_EN` 为同位镜像；围栏块见文末「产物清单」。规模 M+、分级 `[SLICE]`、探针 0、⑦ = `gpt-5.3-codex`（**不降级**）。

## 1. 变更概要（含执行自由度清单）

**目标**：三件同族规则修正同批落地——① **R2.2 修正**（调用串 = 限定名 `Model Name (vendor)`，`modelId` 降为内部标识）；② **OB-72**（台账完成行结构重建）；③ **OB-49**（`tmp/` 生命周期机制）。改动面 = `DELIVERY_DIRECTIVE{,_EN}` 两份文件（§0.3 / §1.5 / §2.2 / §2.3 / §2.4 / §5.0 / §5.2 / §6.2 / §6.3 / §11.2 / §12.4 / 附录 A.1 / 附录 B / 附录 C.0l）＋本报告对＋自由度清单。

**本片行使的执行自由度**（§1.7.3；与本片自由度清单 `docs/validation/evidence/DIRECTIVE-MODELID-MAP-freedom-list.md` 内容一致）：

- **F-1 版本包**：按 §0.2 升 **v1.22**（头版本号 + §0.3 追加一条 + 附录 C.0l 台账块）。台账元数据类副产物。
- **F-2 §2.2 表形态**：**原地对调两列语义**（列数、列序不变；九行数据仅两列换值、取值不变）。
- **F-3 维护机制落点**：限定名维护机制写进 **R2.2 本体**（而非 §2.2 注记或新增小节）。
- **F-4 OB-72 取案 A′**：＝案 A 的结构骨架 + 案 B 的补全条款（完成行只写报告路径；hash 与 CI 号只进报告，由本片主控在 CI 双绿后以独立 docs-only 提交补全，且该提交仍须 §10 同轮授权）。用户给的是案 A / 案 B 二选一，① 判定案 A 仍有自指（报告与完成行同处一个提交）⇒ 取第三案，属**升级边界**并如实声明。
- **F-5 OB-49 落点**：**新增 `### 1.5.1`**（使 G6-5 有可解析节号目标），而非扩写 §1.5 三条或塞入 §11.2。
- **F-6 引用处理**：三处准则内 `tmp/` 引用**只加注记、不改路径**，实际归档与路径改写留待归档执行片。
- **F-7 边界拒绝**：**拒绝** ⑧ 对 5 处**基线遗留**的 CN / `_EN` 反引号级漂移的「顺手修」（逐处回退），依 §1.7.2 判为超出切片范围的顺手修，改为登记候选 OB 交用户裁定。
- **F-8 历史行处置**：§0.3 的 v1.12 条目**不改原文**，只追加「口径更正指针」（依 §11.1）。
- **F-9 编排工件**：报告与清单的结构与措辞。

## 2. ① 判定与落点（方案书 `DESIGN: DONE`）

| # | 判定项 | 结论 |
|---|---|---|
| 1 | **三件合并是否可控** | **可控，不拆分**（改动面集中在同一权威文件对；三件同属规则语义且落在同一审查盲区；⑦ 一次盲审 + 复审可整体覆盖；若拆三片会把 ⑦ 额度从 1+1 增至 3+3，且片间互相锁序） |
| 2 | R2.2 表形态 | 两列**语义对调**（限定名列 = 调用串唯一权威；`modelId` 列 = 内部标识），列数列序不变 |
| 3 | R2.2 本体内核 | **保留**「必须显式传 `model`」，外层形态改为限定名格式，并新增**限定名维护机制** |
| 4 | 候选判据 G8 | **不实现**（§6.3 判据不动），只把其口径同步为「附录 B 的 `model=` 集合 ⊆ §2.2 限定名列」 |
| 5 | OB-72 案 | 案 A′（见 F-4） |
| 6 | OB-49 落点 | 新增 §1.5.1；三处引用只加注记；本片不执行归档 |

**评估项结论（用户要求「若结论是仍需维护机制，须在报告里明确写出」）**：**限定名并不比 `modelId` 更稳**——它只是**唯一被 `runSubagent` 接受的形态**（三次 `modelId` 形式实测全被拒 + 官方文档明文 `handouts.model` 用 `Model Name (vendor)`）；`modelId` 退为最稳的内部标识（白名单闭集 / G5-b / ⑦ 分档表 / `LEGACY_FLASH_MODEL_IDS` 全不受版本漂移影响）。**限定名带版本号（如 `DeepSeek V4.1 Flash`）⇒ 存在版本漂移，仍需维护机制**，机制落 R2.2 本体：调用被工具拒绝时**停下并按 §1.7.2 升级用户裁决**，不得自行换串猜测；§2.2 限定名列与附录 B 模板串的更新须经用户授权、随规则修订进行。**既有盲区（如实）**：工具若静默接受旧串却路由到别的模型，该机制不触发——属映射式规则的固有盲区，本片不解决、只登记。

## 3. 交付面与逐条落点

| 项 | 落点 |
|---|---|
| R2.2 本体与维护机制 | `DELIVERY_DIRECTIVE{,_EN}` 的 R2.2 行 + §2.2 表下限定名注记 |
| §2.2 两列对调 | 同表表头 + 九行数据（值不变） |
| 引用点同步 | §2.3 尾句 / §2.4 两处（黑名单本体不动）/ §5.0 步 2 与步 4 / §6.3 的 G5-b 括注 / §5.2 的 ⑦ 分档表列名 |
| 附录 B 模板 | 六处角色行 + 七处 `model=` 调用串 + B.0 两条清单行 |
| OB-72 | §11.2「台账完成行格式」段重建；OB-72 翻 `已处置` |
| OB-49 | 新增 §1.5.1；附录 A.1 首行与 §0.3 的 v1.10 / v1.12 条目加归档注记；OB-49 翻 `待议` |
| 台账翻转 | OB-67 `待议 → 已处置`；OB-70 保持 `已处置` 并追加落地句；OB-29 追加候选 G8 口径同步句 |
| 版本包 | 头版本 v1.22 + §0.3 条目 + 附录 C.0l 台账块 |
| 自由度清单 | `docs/validation/evidence/DIRECTIVE-MODELID-MAP-freedom-list.md` |

## 4. 执行流水线与环境

① 架构（V4 Pro，拆分判定 + 三件设计）⇒ ② 文档实现（CN 权威稿，逐组 assert 命中数）⇒ ③ 文档一致性核对（`modelId` 残留分类：除历史行与「内部标识」用法外无残留）⇒ ④ 门禁 ⇒ ⑤ 预审（V4 Pro，**两轮**：首轮 `PASS WITH FIXES` ⇒ 整改 ⇒ 复审 `PASS`）⇒ ⑥ 独立扫描（MAI，**两轮**：首轮 `FINDINGS` ⇒ 整改 ⇒ 复审 `PASS`）⇒ ⑦ 终审（Codex，**两轮**：首轮 `PASS WITH FIXES` ⇒ 整改 ⇒ 复审 `PASS`）⇒ ⑧a/⑧b 报告与本 `_EN` 镜像 ⇒ ⑨ 原生验证 ⇒ ⑧c 收尾清洁 ⇒ ⑩ 停点。

**额度实得**：探针 0；① 一轮；⑤ 两轮（额度 1+1）；⑥ 两次（上限 3）；⑦ 两次（额度 1+1）。
环境：Windows 本机；Node v24.17.0；Go 工具链 `go1.27.0 windows/amd64`；仓库 `github.com/larsonzh/prfrail`；分支 `main`，起点 `4a72c8f`。

### 4.1 异常与兜底记录（§11.3 的报告要素）

- **角色越界与回退（F-7）**：⑧ 镜像轮除本片镜像外，另对 **5 处基线遗留**的反引号级 CN / `_EN` 漂移做了纯格式对齐（其中一处改动了英文措辞）⇒ 逐处**回退**，使改动集只含本片镜像；差异登记为候选 OB（见 §6）。
- **`create_file` 落盘编码异常（工具侧）**：自由度清单落盘为无 BOM + CRLF ⇒ ⑦ 首轮判红（G4a）⇒ 归一为 BOM + LF 后复绿；此后新增文本一律先归一。
- **终端命令链中断（工具侧）**：长命令链被终端误判为中断，且链前半可能已执行 ⇒ 受影响步骤按「每步一次调用」重跑，并在中断后先核对工作区与索引是否一致。
- **① 判定与用户设想的偏差（如实）**：用户给出的翻转项含「OB-70 `待议 → 已处置`」，但 OB-70 在 v1.19 已是 `已处置` ⇒ 实际动作改为 **OB-67 翻 `已处置`**、OB-70 追加落地句，并同步更正 C.0l 附注措辞（⑤ 首轮 Low 提出、已改）。
- **无其他兜底**：未触发 §9.1 的其余场景（无预算超限、无依赖阻塞、无原生反复失败）。

## 5. 审查结论

- ⑤ 预审（V4 Pro，两轮）：首轮 `PASS WITH FIXES`（2 Medium：§11.2 补全提交未衔接 §10 授权、§1.5.1 缺 §1.5 显式豁免；6 Low：无 CI 切片不补 hash 与 run 号、OB-67 处置列措辞、C.0l 附注 OB-70 不一致、归档需先过 ⑧c 编码门禁、清理时机需先全量盘点、§6.2 路径规则指向 §1.5.1）⇒ 逐条整改 ⇒ 复审 **`PASS`**（全部 closed，无新问题）。
- ⑥ 独立扫描（MAI，两轮）：首轮 `INDEPENDENT SCAN: FINDINGS`（1 High：§0.3 的 v1.12 历史条目仍称「附录 B 七处调用串改 modelId」；1 Medium：5 处基线反引号级镜像漂移）⇒ High 以**追加口径更正指针**闭合（不改历史原文）、Medium 按**边界拒绝**并登记候选 OB ⇒ 复审 **`INDEPENDENT SCAN: PASS`**（明确接受两条处置，无新增 High/Medium）。
- ⑦ 终审（Codex，两轮）：首轮 `FINAL REVIEW: PASS WITH FIXES`（1 Medium：自由度清单缺 BOM + CRLF ⇒ G4a 判红；其余四项与 F-1 至 F-9 逐条判无越界）⇒ 归一编码 ⇒ 复审 **`RE-REVIEW: PASS`**（无 finding，双作用域复现 `TOTAL_FAIL=0`）。**⑦ 不降级 ⇒ 本片无 ADR**（① 曾建议登记 ADR-023，按用户回报格式「若不降级，无 ADR」不登记）。

### 5.1 审查输入包摘要（含盲审隔离证明）

变更集 diff = `tmp/dml/change-set.diff`（配 `tmp/dml/numstat.txt`）；契约原句 = `DELIVERY_DIRECTIVE.md` 的 §0.2 / §1.5 / §1.7.2 / §1.7.3 / §2.2 / §2.3 / §2.4 / §5.0 / §5.1 / §5.2 / §6.2 / §6.3 / §7.2 / §7.5 / §10 / §11.1 / §11.2 / §12.4；自由度清单 = §3 末行的固定路径。**盲审隔离证明**：本片触「身份」语义（模型白名单 / 黑名单即身份准入）⇒ ⑦ **首轮以盲态执行**——输入仅含切片定义、原始 diff（路径）、契约原句、自由度清单与只读约束，**未附** ⑤ 预审报告与 ⑥ 扫描报告；复审轮按附录 B.6 允许附上「独立扫描完成后读取」的摘要。输入消毒 = 材料全为仓库内文本，无凭据 / token / SSH 路径。

## 6. 可证伪性与门禁结果

| 检查 | 结果 |
|---|---|
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0`（changed=7；报告对与台账行入暂存后的最终读数） |
| `gate.js --all --scope=index` | `TOTAL_FAIL=0`（changed=7） |
| `gate.js --selftest` | `SELFTEST: PASS 264/264`（本片不新增夹具） |
| 基础门禁（§6.1） | `gofmt` 空、`go build` 与 `go vet` 与 `go test` 均 exit 0（ok=14、FAIL=0） |
| 镜像同位（含五形状） | 读数见 `tmp/dml/pairs.txt`（两侧元素数相同、形状仅 5 处基线遗留反引号漂移、非本片改动行） |
| 变异检验 | 读数见 `tmp/dml/mutation.md`：① 改 OB-49 首状态词 ⇒ G6-6 判红；② 改 §0.3 内的节号引用 ⇒ G6-5 判红；③ 仅改 CN 增行 ⇒ G3(a.1) 与 G3(a.2) 判红；三例均以 sha256 复核后逐字节还原并复绿 |
| 新增行预检 | `tmp/dml/pre.txt`：`ADDED-LINE HITS=0`、`UNKNOWN-CJK=`（无未登记新字） |

**可证伪性边界（如实）**：本片三件中，① 的**调用串口径**与 ③ 的**三类处置口径**属**无新机械判据**的规范（用户约束「§6.3 判据不动」）⇒ 其核验 = ⑤ / ⑥ / ⑦ 人工复核 + G6-5 / G6-6 旁证；**变异检验不覆盖它们**，不得读作「已由判据守住」。② 的完成行格式同理（无判据）。

**已登记 OB（用户 2026-10-04 批准，首状态词 `待议`；本片已在 §12.4 登记）**：
- **OB-73**：CN / `_EN` 的反引号级镜像漂移（**5 处基线遗留**，其中 2 处含 EN 独有句子）——G3(a.2) 的五形状不含反引号计数，故门禁不覆盖；修正需动与本片无关的历史行 ⇒ 本片按边界拒绝（⑥ Medium 的处置）。
- **OB-74**：Haiku 作为 §2.2 之外唯一豁免，其**限定名未定义** ⇒ 启用时无调用串可写（既有缺口，本片未引入、未声称解决）。
- **OB-75（本片评估项的直接产物，须明确写出）**：R2.2 的维护机制**只在调用被拒绝时触发**；若工具**接受旧串却把请求路由到其它模型**，「调用串 = 限定名」在静默路由面前**没有可观测判据**——这是**已被识别的盲区**（不是缺陷，属映射式规则的固有面），本片只把它写进 R2.2 与台账、**未解决**，留待模型可观测性机制（operator 侧实际生效模型的记录与核对）处理。

## 7. 成本与计量

① 一轮、⑤ 两轮、⑥ 两次、⑦ 两次（Codex 常规档，未降级）、⑧ 两次（CN 权威稿由主控落地、`_EN` 镜像由文档工程师产出）；探针 0；无 Haiku 调用；**未发生额度例外，故本片无 ADR**。

## 8. 明确未执行事项

- 未实现候选判据 **G8**；未改 §6.3 任何判据本体、未动 `g5b.js` 与夹具。
- 未改 §2.4 黑名单本体；未改 §5.2 矩阵；未改 §2.2 表的**取值**（含 `GPT-6-Luna (copilot)` 的连字符形态与工具实测空格形态的分歧，只登记）。
- **未实际移动、归档或删除任何 `tmp/` 文件**（归档执行另立「tmp/ 生命周期执行」片）。
- 未修 5 处基线反引号级镜像漂移（OB-73，已登记）。
- 未登记 ADR（⑦ 未降级）；未 commit、未 push（待用户同轮授权）。

## 9. 下一步建议（含各角色档位建议）

- **下一片（建议）**：项目名 `tmp/` 生命周期执行片——按 §1.5.1 的三类口径执行存量 ↔ 引用全量盘点、归档与引用改写；规模 **M+**、⑦ 用 `gpt-5.3-codex`（触规则正文引用）。
- **⑦ 档位建议**：本片及归档执行片均属规则语义类 ⇒ 常规档 `gpt-5.3-codex` 足够；不触安全 / 机器状态 / 协议语义 ⇒ 无需深度档。
- **⑧ 角色建议**：镜像轮须先读 CN 权威稿的逐组改动，**禁止**对齐与本片无关的基线漂移（本片实测教训）。

## 10. ⑩ 停点报告（含回写核对行）

- 执行至 ⑨；**停在提交授权之前**：本片不 commit、不 push、不推 gitee。
- **回写核对**（§11.2 的 OB-63 口径）：验证报告 = **已回写**（本文件）；`DEV_PLAN` = **不适用**（准则治理片，按切片定义与既有先例不写 `DEV_PLAN` 段落）；台账 = **已回写**（`REMAINING_SLICES{,_EN}` 的完成记录行，按 §11.2 新格式只记报告路径）；`DELIVERY_DIRECTIVE{,_EN}` 的 §0.3 与 §12.4 已就地更新；ADR = **不适用**（无降级、无额度例外）。
- **⑦ 结论**：首轮 `PASS WITH FIXES`（1 Medium：自由度清单编码）⇒ 归一编码 ⇒ 复审 `RE-REVIEW: PASS`；⑦ 为 Codex 常规档、**未降级**，独立性不受损。
- **门禁读数**：`--scope=tree` 与 `--scope=index` 均 `TOTAL_FAIL=0`；`SELFTEST: PASS 264/264`（见 §6）。

## 11. 产物清单

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-modelid-map.md	repo
docs/validation/directive-modelid-map_EN.md	repo
docs/validation/evidence/DIRECTIVE-MODELID-MAP-freedom-list.md	repo
tmp/dml/bing.txt	local-only
tmp/dml/change-set.diff	local-only
tmp/dml/check_bt.js	local-only
tmp/dml/chuanshu.txt	local-only
tmp/dml/ci-run-37073985144.txt	local-only
tmp/dml/ci-run-37092310920.txt	local-only
tmp/dml/cn-shape.txt	local-only
tmp/dml/consistency-10.txt	local-only
tmp/dml/consistency-final.txt	local-only
tmp/dml/consistency-out4.txt	local-only
tmp/dml/consistency.js	local-only
tmp/dml/contract-excerpt.md	local-only
tmp/dml/design-1-v4pro.md	local-only
tmp/dml/directive_cn.diff	local-only
tmp/dml/directive_en.diff	local-only
tmp/dml/dump-cn.js	local-only
tmp/dml/dump.js	local-only
tmp/dml/dump2.js	local-only
tmp/dml/dump3.js	local-only
tmp/dml/dump6.js	local-only
tmp/dml/dump7.js	local-only
tmp/dml/en-stat.js	local-only
tmp/dml/en3.txt	local-only
tmp/dml/en58.txt	local-only
tmp/dml/en6.txt	local-only
tmp/dml/enrep.txt	local-only
tmp/dml/fence-cmp.js	local-only
tmp/dml/fence.js	local-only
tmp/dml/find.js	local-only
tmp/dml/find.txt	local-only
tmp/dml/fix-fixtures-lf.js	local-only
tmp/dml/fix-h5.js	local-only
tmp/dml/fix-ledger.js	local-only
tmp/dml/fix-prereview.js	local-only
tmp/dml/fix-r2-1.js	local-only
tmp/dml/fix-r7-1.js	local-only
tmp/dml/g4b.txt	local-only
tmp/dml/gate-index-final.txt	local-only
tmp/dml/gate-index.txt	local-only
tmp/dml/gate-index10.txt	local-only
tmp/dml/gate-tree.txt	local-only
tmp/dml/gate-tree10.txt	local-only
tmp/dml/gate-tree7.txt	local-only
tmp/dml/gate-tree9.txt	local-only
tmp/dml/gofmt.txt	local-only
tmp/dml/grep_misc.txt	local-only
tmp/dml/grep_modelid.txt	local-only
tmp/dml/grep_modelstr.txt	local-only
tmp/dml/grep_showname.txt	local-only
tmp/dml/gsum.js	local-only
tmp/dml/hunks.js	local-only
tmp/dml/hunks.txt	local-only
tmp/dml/land-8b.js	local-only
tmp/dml/land-adr019-conclusion.js	local-only
tmp/dml/land-adr019.js	local-only
tmp/dml/land-r2-2.js	local-only
tmp/dml/land-rules.js	local-only
tmp/dml/land1.js	local-only
tmp/dml/land10.js	local-only
tmp/dml/land11.js	local-only
tmp/dml/land12.js	local-only
tmp/dml/land2.js	local-only
tmp/dml/land3.js	local-only
tmp/dml/land4.js	local-only
tmp/dml/land5.js	local-only
tmp/dml/land6.js	local-only
tmp/dml/land7.js	local-only
tmp/dml/land7b.js	local-only
tmp/dml/land8.js	local-only
tmp/dml/land9.js	local-only
tmp/dml/leftovers.js	local-only
tmp/dml/leftovers.txt	local-only
tmp/dml/lib.js	local-only
tmp/dml/luna-probe	local-only
tmp/dml/m1.orig	local-only
tmp/dml/m2.orig	local-only
tmp/dml/m3.orig	local-only
tmp/dml/mk-fixtures-f5.js	local-only
tmp/dml/mk-fixtures-m4.js	local-only
tmp/dml/mk-freedom-list.js	local-only
tmp/dml/mk-report.js	local-only
tmp/dml/mkpack.js	local-only
tmp/dml/mut.js	local-only
tmp/dml/mutation.md	local-only
tmp/dml/norm-en.js	local-only
tmp/dml/norm.js	local-only
tmp/dml/numstat.txt	local-only
tmp/dml/ob-land.js	local-only
tmp/dml/ob-land2.js	local-only
tmp/dml/pairs.js	local-only
tmp/dml/pairs.txt	local-only
tmp/dml/pairs2.js	local-only
tmp/dml/pre.txt	local-only
tmp/dml/preflight.js	local-only
tmp/dml/r2-cn.diff	local-only
tmp/dml/r2-cn2.diff	local-only
tmp/dml/r2-consistency.txt	local-only
tmp/dml/r2-en.diff	local-only
tmp/dml/r2-gate-tree.txt	local-only
tmp/dml/report-stat.txt	local-only
tmp/dml/revert.js	local-only
tmp/dml/review-r1.md	local-only
tmp/dml/scan.js	local-only
tmp/dml/scan.txt	local-only
tmp/dml/scan2.js	local-only
tmp/dml/seg.js	local-only
tmp/dml/seg.txt	local-only
tmp/dml/selftest-final.txt	local-only
tmp/dml/selftest.txt	local-only
tmp/dml/selftest10.txt	local-only
tmp/dml/selftest7.txt	local-only
tmp/dml/src.txt	local-only
tmp/dml/src2.txt	local-only
tmp/dml/src3.txt	local-only
tmp/dml/src4.txt	local-only
tmp/dml/src5.txt	local-only
tmp/dml/strict-pairs.js	local-only
tmp/dml/test.txt	local-only
```
