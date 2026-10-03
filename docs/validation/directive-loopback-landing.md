# 切片验证报告 · `DIRECTIVE-LOOPBACK-LANDING`（轻量回流机制落地）

> 依据 `docs/DELIVERY_DIRECTIVE.md` v1.20（本片起跑版本；本片落 v1.21 版本包）的 §5.1 串行流水线与 §6.2 证据落盘义务。**CN 为权威版**，`_EN` 为同位镜像；围栏块见文末「产物清单」。规模 S、分级 `[SLICE]`、探针 0。**起跑消息经用户同轮修正**：D-3 模板落点改为 `docs/` 根（项目级工具，不专属 T027），并新增 ① 判定第 4 项（新模板是否同批登记进文档清单）。

## 1. 变更概要（含执行自由度清单）

**目标**：落地四组轻量机制——D-1 主控主动建议登记原则、D-2 `DOCUMENTATION_PLAN` 跨层指针、D-3 结项报告模板 + T027 结束锚点、C 台账完成行格式。**不触规则本体**（§2.2 / §2.4 / §5.2 / §6.3 / §1.7 的规则语义零改动；⑦ 复核）。

**本片行使的执行自由度（§1.7.1，实质选择）**：

- **F-1 版本包**：本片改 §11.2 正文 ⇒ 按 §0.2 升 **v1.21**（头版本号 + §0.3 变更日志行 + 附录 C.0k 台账块）。属台账元数据类副产物。
- **F-2 D-1 与 C 的落点（① 判定 1 与 2）**：两者同落 **§11.2**、**各自独立成段**（行首加粗主题标签），**不新增 `### 11.2.x`**（避免编号体系扰动与标题行号序列变化）。
- **F-3 模板落点（① 判定 3 已按用户同轮修正作废）**：落 **`docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md`**（`docs/` 根，与 `DELIVERY_DIRECTIVE` / `DOCUMENTATION_PLAN` / `CONTRACTS` 同级）。理由：本模板是**项目级治理工具**，不专属 T027；**T027 结束只是首个应用锚点，不是模板归属依据**。① 首轮所判的 `docs/t027/` 随之作废并如实留存；`docs/validation/` 仍因属 G7 扫描域而被排斥。
- **F-4 OB-49 处置列按用户裁定改为「替换旧片名」而非「补句保留旧名」**：① 原案是在保留 `DIRECTIVE-LEDGER-ARCHIVE` 的前提下补「归档三源」；本片按用户切片定义的原文（「把处置列里的旧片名改为与评估方向一致的表述」）**删除旧片名**，改指「`tmp/` 生命周期评估」片并并列三类处置。该差异在此声明。
- **F-5 本片台账完成行首次适用新格式**：只记报告路径，**不含 commit hash 与 CI run 号**。首例即暴露 C 的残留自指问题（commit hash 与 CI 号同属「提交前不可知」）——已如实登记并附建议（见 §7 末）。
- **F-6 §11.2 由单段扩为三段**（回写口径 / 登记建议与留痕 / 台账完成行格式）：未新增编号小节，与①②的判定一致。
- **F-7 ⑦ 降级**：本片 ⑦ 用 `deepseek-v4-pro`（用户切片定义明示），按 OB-37 元规则以 **ADR-022** 留痕。
- **F-8 D-2 同批登记新模板（① 判定 4）与其连带项**：按 ① 判定在 `DOCUMENTATION_PLAN{,_EN}` §2 文档清单新增模板行；连带把该文件正文中**现行生效**的「12 对文档」改为「13 对」（§3 的「文档基线 12 对」属已执行阶段的记录，未改）。该项属 ① 判定未列出的连带补齐，为免「地图已改、计数仍旧」的最小化处理，在此声明并交 ⑩ 裁定。

## 2. 启动前清单核对（步骤 0；逐条）

| # | 项 | 状态 |
|---|---|---|
| 1 | D-1 原则现状盘点（准则无此原则） | **已处理**（落 §11.2 新段） |
| 2 | D-2 跨层指针现状盘点（跨层引用近零） | **已处理**（落 `DOCUMENTATION_PLAN{,_EN}` §1 表后） |
| 3 | D-3 模板与锚点现状盘点（无模板、无锚点行） | **已处理**（新建项目级模板对，落 `docs/` 根 + `REMAINING_SLICES{,_EN}` 用法节第 4 条锚点） |
| 4 | C 台账完成行格式现状盘点 | **已处理**（落 §11.2 新段；本片完成行首次适用） |
| 5 | OB-49 指针修正（旧片名 → 评估方向） | **已处理**（一行改动，未单独成项） |
| 6 | 台账 CI 回填（累计两处） | **已处理**：`DIRECTIVE-DESIGN-LOOPBACK` 行回填 `a815970` 与 run `37107634723`（run 号出处 = `gh run view` 读数 `conclusion=success` 与 `headSha` 前缀 `a815970` 同台账一致，读数已归档为 `docs/validation/evidence/ci-run-37107634723.txt`（`repo` 处置）；本片未重跑 CI）；本片新行按新格式只记报告路径 |
| 7 | OB-71 状态翻转 `待议 → 已处置` | **已处理**（D-3 落地即其处置） |
| 8 | 上一片遗留 `DIRECTIVE-CLEANUP-EVAL` 行 | **不适用**（上一片 D-5 已回填；本片只读复核） |
| 9 | `DIRECTIVE-MODELID-MAP`、OB-68、OB-53⑤ | **不适用**（不在本片范围） |
| 10 | `tmp/` 根下残留 | **不适用**（留待 OB-49 立项；本片仅改其指针） |

## 3. ① 判定与落点（方案书 `DESIGN: DONE`；含用户修正后的补充判定轮）

| # | 判定项 | 结论 |
|---|---|---|
| 1 | D-1 落 §1.7 还是 §11.2 | **§11.2**（登记建议属「回写 / 登记」域；§1.7 是自由度与升级边界，塞入会破坏其语义边界） |
| 2 | C 与 D-1 同段还是独立成段 | **同节、各自独立成段**（两主题语义异族；不新增编号小节以免扰动标题行号序列） |
| 3 | D-3 模板路径 | **原判 `docs/t027/` 已作废**（① 首轮结论如实留存）；按用户同轮修正改判 **`docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md`**（`docs/` 根，项目级、通用形态） |
| 4 | CN / `_EN` 双语对称方案 | 每处 CN 改动在 `_EN` 同逻辑位置同构落地；加粗逐句对齐、表格不改列结构、空行节奏一致、ID 计数两侧相等（实测：五形状全等）。模板迁至 `docs/` 根后**不再受 G7 扫描域覆盖**，围栏义务转由本报告的 `repo` 处置行承担 |
| 5 | 新模板是否同批登记进 `DOCUMENTATION_PLAN` 文档清单（修正轮新增） | **登记**（文档地图须列全项目级文档，避免「文档在盘而地图未列」的空窗；D-2 本批已动该文件，同批省一次往返）⇒ 落其 §2 清单一行，CN / `_EN` 各 +1 行 |

## 4. 交付面与逐条落点

| 项 | 落点 |
|---|---|
| D-1 原则 + 首例 | `docs/DELIVERY_DIRECTIVE{,_EN}.md` §11.2 新段「登记建议与留痕」 |
| C 台账完成行格式 | 同节新段「台账完成行格式」 |
| D-2 跨层指针 | `docs/DOCUMENTATION_PLAN{,_EN}.md` §1 权威分域表后 |
| D-3 模板 | 新建 `docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md`（项目级、通用形态；回看 = 主控；前瞻 = 主控清单 + 用户裁决） |
| D-2 同批登记 | `docs/DOCUMENTATION_PLAN{,_EN}.md` §2 文档清单新增结项报告模板行（① 判定 5） |
| D-3 锚点 | `docs/t027/REMAINING_SLICES{,_EN}.md` 用法节第 4 条 |
| OB-71 / OB-49 | `docs/DELIVERY_DIRECTIVE{,_EN}.md` §12.4 的两行处置列 |
| 版本包 | 头版本 v1.21 + §0.3 行 + 附录 C.0k 台账块 |
| ADR | `docs/ADR_REGISTER{,_EN}.md` 的 ADR-022 |
| 自由度清单 | `docs/validation/evidence/DIRECTIVE-LOOPBACK-LANDING-freedom-list.md` |

## 5. 执行流水线与环境

① 架构（V4 Pro，**两轮**：首轮判定项 1 至 3 + 用户修正后的补充判定项 4 与 5）⇒ ② 文档实现 ⇒ ③ 清单核对与一致性 ⇒ ④ 门禁 ⇒ ⑤ **跳过**（S 档可选）⇒ ⑥ **跳过**（纯 `.md`、无代码语义，按 §7.2 的跳过条件，由 ⑦ 确认）⇒ ⑦ 终审（`deepseek-v4-pro`，**ADR-022 降级**；**独立性降低如实标注**）⇒ ⑧ 报告与 `_EN` 镜像 ⇒ ⑨ 原生验证 ⇒ ⑧c 收尾清洁 ⇒ ⑩ 停点。

**额度实得**：探针 0；付费调用 = ① 方案书**两轮**（首轮 + 用户修正后的补充判定；**超 §7.5「① ≤ 1」**，如实登记）、⑦ 终审一轮与复审一轮（如触发）、⑧ 镜像一轮。
环境：Windows 本机；Node v24.17.0；Go 工具链 `go1.27.0 windows/amd64`；仓库 `github.com/larsonzh/prfrail`；分支 `main`，起点 `a815970`。

### 5.1 异常与兜底记录（§11.3 的报告要素）

- **新建 `.md` 的落盘编码异常（工具侧）**：`create_file` 落盘两个模板文件时带 CRLF，暂存期由仓库规范化转为 LF/`BOM` 归一（该现象与上一片 `create_file` 的同类行为一致）⇒ 兜底 = 每次新建后**显式复核 BOM / LF / 行数**，并以脚本统一归一；本片两个模板文件在验证前即已归一（BOM + LF）。
- **核对顺序异常**：首次形状核验读到的模板对「空行数 22 / 1」不对称，原因是核验脚本在未归一 CRLF 的前提下计数 ⇒ 兜底 = 核验前先归一换行，再逐行对位；归一后模板对**两侧元素数相同（按 `split(/\n/)` 口径各 44 条）、空行各 22 条、标题行号序列全等**。
- **起跑消息中途修正 ⇒ ① 配额超限（如实登记）**：本片 ① 累计调用 **2 次**（首轮判定项 1 至 3 + 用户同轮修正后的补充判定）。① 架构师本片累计调用 2 次，超出 §7.5 单片 ≤1 上限，**系用户同轮修正指令驱动**，按「超出上限 ⇒ 升级到用户裁决」口径执行，**仅限本片、不得作为先例**；该升级已由用户同轮修正本身完成，故不另立 OB / ADR。两轮调用均已计入本节与「额度实得」。
- **G6-4 判红与纠错（未用白名单）**：模板迁根后首跑门禁，G6-4 判红两条（本报告与自由度清单内的行数字面量）；**该红态读数未单独留档**（已被后续复跑覆盖）⇒ 改写为「元素数 / 空行数」措辞后复绿，复绿后读数见 `tmp/ll/pre.txt`（新增行口径 `ADDED-LINE HITS=0`）与 `tmp/ll/lit.txt`（全文口径，命中均在历史行），**未登记白名单豁免**。
- **G2 占位符残留在整改新增行上判红**：⑦ 首轮整改新增的 §8 条目含以「待」字起头的同轮授权措辞（占位式），触 G2 的占位符词表 ⇒ 改写为「须用户同轮授权」后复绿，**未登记白名单**；该红态是本片第二次「整改引入新判红」的实例，一并记入自由度清单 §3。
- **连带项自查**：模板迁根后逐处核对路径引用（准则 §0.3 的 v1.21 变更日志条目、§12.4 的 OB-71 处置列、附录 C.0k 交付物行、台账用法节与完成记录行、本报告围栏），并按 ① 判定补齐 `DOCUMENTATION_PLAN` 的清单行与现行计数。

## 6. 审查结论

- ⑤：**跳过**（S 档可选）。
- ⑥：**跳过成立**（逐条核对，由 ⑦ 确认）：① 变更集全为 `.md` 文案类；② 无 `.go` / `.json` / `.yml` / schema / fixtures / 脚本逻辑；③ 不触状态机 / 门禁判据 / 证据模型 / 角色授权语义。
- ⑦：**结论见 §8**（本片 ⑦ 由 `deepseek-v4-pro` 承担 ⇒ 独立性降低，已在 ADR-022 与本节如实标注）。

### 6.1 审查输入包摘要（含盲审隔离证明）

变更集 diff = `tmp/ll/change-set.diff`（配 `tmp/ll/numstat.txt`）；契约原句 = `docs/DELIVERY_DIRECTIVE.md` 的 §0.2 / §6.3 / §7.2 / §7.5 / §11.2 / §11.3 / §12.4；① 两轮判定的汇总以本报告 §3 为准（首轮判定 3 已由用户修正作废）；自由度清单 = `docs/validation/evidence/DIRECTIVE-LOOPBACK-LANDING-freedom-list.md`；**盲审隔离证明**：⑤ 与 ⑥ 均跳过 ⇒ 无中间清单可附，⑦ 首轮只收切片定义、diff、准则原句与只读约束（§7.3）；输入消毒 = 材料全为仓库内文本，无凭据 / token / SSH 路径。

## 7. 可证伪性与门禁结果

| 检查 | 结果 |
|---|---|
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0`（changed=13） |
| `gate.js --all --scope=index` | `TOTAL_FAIL=0`（changed=13） |
| `gate.js --selftest` | `SELFTEST: PASS 264/264`（本片不新增夹具） |
| 基础门禁（§6.1） | `gofmt` 空、`go build` 与 `go vet` 与 `go test` 均 exit 0（ok=14、FAIL=0） |
| 镜像同位（含五形状） | 读数见 `tmp/ll/mirror.txt` |
| 变异检验 | 读数见 `tmp/ll/mutation.md`（`docs/DELIVERY_DIRECTIVE.md:879` 判红、引用改号判红、模板对不对称判红；逐条 sha256 复核后逐字节还原并复绿；**文末另有「最终复跑」收尾段**） |
| G6-4 纠错 | 首跑判红两条（行数字面量，红态未单独留档）⇒ 改写措辞后复绿，**未用白名单**（复绿后读数见 `tmp/ll/pre.txt` 与 `tmp/ll/lit.txt`） |
| 清单核对（§2） | 十项逐条标注，其中 1 至 7 已处理、8 至 10 不适用 |

**C 与 D-1 的可证伪性边界（如实）**：两者属**无新机械判据**的格式与行为规范（本片硬约束不新增判据）⇒ 其核验 = ⑦ 人工复核 + ⑩ 回写核对行 + G2 / G6-4 旁证；**变异检验不覆盖它们**，不得读作「已由判据守住」。

**C 的残留自指问题（本片首例观察 + 建议）**：C 的立意是消除「CI run 号提交前不可知」造成的跨片回填，但 **commit hash 同属提交前不可知**（完成行写在提交之内）⇒ 本片完成行只能记报告路径，hash 仍需下一片回填。**建议**（按本片刚落地的 D-1 原则提出，交用户裁定）：完成行改记「**写入时的 HEAD（BASE）hash** + 报告路径」，或改为「**提交后由同轮的台账追加提交**」，二者均可让完成行在写入时即可确定。

**候选 OB（本片未登记，按约束只落地 OB-71 / OB-49 两处处置列）**：① 完成行的 BASE-hash 口径（见上）；② 「§11.2 扩为三段后是否拆节」的体例观察（① 已提示主题密度上升）。

### 7.1 ⑧c 收尾清洁记录（§5.1 固定环节，主控执行）

`gate.js --all --scope=tree` 全绿；编码门禁（G4a）按 `.md` = BOM+LF 对全部改动文件逐文件核验，无违规；诊断扫描（§3.4）无新增告警；工作区卫生 = 全部改动逐文件暂存；**无 git 可见的未暂存项与未跟踪项**（`tmp/ll/` 属 gitignore 目录）（读数见 `tmp/ll/eightc.txt`）；`tmp/ll/` 为一次性工作目录（gitignore），其本地证据在用户审阅与提交后按 §1.5 清除。

## 8. ⑩ 停点报告（含回写核对行）

- 执行至 ⑨；**停在提交授权之前**：本片不 commit、不 push、不推 gitee，等待用户同轮授权。
- **档位与返工**：规模 S、分级 `[SLICE]`；返工 = ① 追加一轮（用户修正驱动，已登记）+ G6-4 纠错一轮 + ⑦ 首轮整改一轮；无超时、无阻塞。
- **回写核对**（§11.2 的 OB-63 口径）：验证报告 = **已回写**（本文件）；`DEV_PLAN` = **不适用**（准则治理片，按切片定义与既有先例不写 `DEV_PLAN` 段落）；台账 = **已回写**（`REMAINING_SLICES{,_EN}` 的 T027 锚点行与完成记录行、`DELIVERY_DIRECTIVE{,_EN}` §12.4 的 OB-71 与 OB-49 处置列）；ADR = **已回写**（ADR-022）。
- **门禁读数**：`--scope=tree` 与 `--scope=index` 均 `TOTAL_FAIL=0`；`SELFTEST: PASS 264/264`（见 §7）。
- **明确未执行事项**：OB-49 的清理本体（本片只改指针）；`tmp/` 根下残留；`DOCUMENTATION_PLAN` §3 历史计数的追溯改写；本片的 commit 与 push（须用户同轮授权）。
- **下一步建议**：① 请用户裁定 C 的 BASE-hash 口径（⑦ 建议登记为 OB-72）；② 下一片按 OB-49 立项评估 `tmp/` 生命周期并落三类处置；③ ⑦ 复审后的整改由 ⑧b 记入本节；⑤ 与 ⑥ 不适用（已跳过），⑦ 与 ⑧ 沿用本片档位。
- **⑦ 结论**：首轮 `PASS WITH FIXES`（1 Medium 与 6 Low；Medium = §0.3 ④ 与 §12.4 的 OB-49 口径不一致，已改正）⇒ 整改 ⇒ 复审 `PASS WITH FIXES`（无 Medium 及以上；3 条 Low 已由 ⑧b 就地修：EN 的 `->`、证据文件 NOTE 与口径、CI run 引证）；⑦ 为 `deepseek-v4-pro` 降级（ADR-022），**独立性降低如实标注**，最终裁决归用户与 ⑧c。
### 8.1 ⑦ 两轮 finding 与整改记录（§5.3）

| # | ⑦ finding | 处置 |
|---|---|---|
| F-1 | §0.3 v1.21 变更日志的 ④ 项与 §12.4 的 OB-49 落盘口径不一致（Medium） | 已改为「改指 + 三类处置」，与附录 C.0k 附注同口径 |
| F-2 | OB-49 处置列的 EN 重复句与 `->`、CN 括注残缺（Low） | 双侧已改写；CN 括注补全为「本项列为下一片评估对象」 |
| F-3 | 变异证据以中间态红尾收档（Low） | 已补「最终复跑」收尾段与 NOTE 说明；三例 VERIFY 均 `TOTAL_FAIL=0` |
| F-4 | numstat 与门禁读数相对索引略陈旧（Low） | 已按整改后的索引重新生成全部读数 |
| F-5 | G6-4 纠错条目引用的读数文件为复绿后状态（Low） | 已注明「红态未单独留档」并改引 `tmp/ll/lit.txt` |
| F-6 | §8 缺档位与返工、未执行事项、下一步建议（Low） | 已补三条（见本节上文） |
| F-7 | §7.1「无未跟踪项」与并列读数冲突（Low） | 已限定为「无 git 可见的未暂存项与未跟踪项」并注明 `tmp/ll/` 属 gitignore |
| F-8 | 复审：EN 的 OB-49 处置列仍留 `->` 两处（Low） | 已改为 `⇒`；与 CN 的同类符号一致 |
| F-9 | 复审：`tmp/ll/lit.txt` 为全文口径且无 NOTE，易误读为门禁红态（Low） | 已加 NOTE 头，并改为同时引 `tmp/ll/pre.txt`（新增行口径） |
| F-10 | 复审：CI run 引证为自指（Low） | 已补 `gh run view` 读数出处并归档为 `docs/validation/evidence/ci-run-37107634723.txt`（`repo`） |
### 8.2 ⑩ 停点后裁定回执（事实记录）

本节为事实记录，**不改写 §8 的 ⑩ 停点当时表述**（依 §11.1）。

| # | 用户裁定（2026-10-03，⑩ 停点后同轮） | 处置 |
|---|---|---|
| 1 | ① 配额超限（2 次）**接受**，不重新起跑、不停下升级；留痕口径「仅限本片、不得作为先例、不另立 OB / ADR」亦接受 | 已接受；报告 §5.1 与「额度实得」的口径不变 |
| 2 | C 的 BASE-hash 自指**登记为 OB-72**（首状态词 `待议`） | 已登记于 `DELIVERY_DIRECTIVE{,_EN}` §12.4，并计入 v1.21 版本包 |
| 3 | OB-49 **不现在立项**，保留为下一片候选（建议排在 `DIRECTIVE-MODELID-MAP` 之后） | 指针保留（本片已完成）；立项时机不在本片 |
| 4 | 授权 commit 与 push origin（逐文件暂存、不推 gitee）；`tmp/ll/` 提交后清理，其中 CI run 读数归档为 `docs/validation/evidence/ci-run-37107634723.txt` | 提交与推送见台账完成行；归档件按 `repo` 处置入本报告 §9 围栏 |

**本片即为 OB-72 的活例**：本节的 commit hash 与 CI run 号在写入时同样不可知，只能由下一片或本节修订补全。

## 9. 产物清单

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/DOCUMENTATION_PLAN.md	repo
docs/DOCUMENTATION_PLAN_EN.md	repo
docs/SLICE-CLOSEOUT-TEMPLATE.md	repo
docs/SLICE-CLOSEOUT-TEMPLATE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-loopback-landing.md	repo
docs/validation/directive-loopback-landing_EN.md	repo
docs/validation/evidence/DIRECTIVE-LOOPBACK-LANDING-freedom-list.md	repo
docs/validation/evidence/ci-run-37107634723.txt	repo
tmp/ll/apply.js	local-only
tmp/ll/change-set.diff	local-only
tmp/ll/chk.js	local-only
tmp/ll/ci-run.txt	local-only
tmp/ll/cn-dump.txt	local-only
tmp/ll/dump.js	local-only
tmp/ll/dump2.js	local-only
tmp/ll/dump3.js	local-only
tmp/ll/dump4.js	local-only
tmp/ll/dump5.js	local-only
tmp/ll/eightc.js	local-only
tmp/ll/eightc.txt	local-only
tmp/ll/en-dump.txt	local-only
tmp/ll/en-full.txt	local-only
tmp/ll/en-new.md	local-only
tmp/ll/finalize.js	local-only
tmp/ll/fixmut.js	local-only
tmp/ll/g3a.txt	local-only
tmp/ll/gate-index.txt	local-only
tmp/ll/gate-tree.txt	local-only
tmp/ll/gsum.js	local-only
tmp/ll/gsum2.js	local-only
tmp/ll/land1.js	local-only
tmp/ll/land10.js	local-only
tmp/ll/land11.js	local-only
tmp/ll/land12.js	local-only
tmp/ll/land13.js	local-only
tmp/ll/land14.js	local-only
tmp/ll/land2.js	local-only
tmp/ll/land3.js	local-only
tmp/ll/land3.txt	local-only
tmp/ll/land4.js	local-only
tmp/ll/land4.txt	local-only
tmp/ll/land5.js	local-only
tmp/ll/land5.txt	local-only
tmp/ll/land6.js	local-only
tmp/ll/land7.js	local-only
tmp/ll/land7.txt	local-only
tmp/ll/land8.js	local-only
tmp/ll/land8.txt	local-only
tmp/ll/land9.js	local-only
tmp/ll/lit.js	local-only
tmp/ll/lit.txt	local-only
tmp/ll/llmirror.js	local-only
tmp/ll/m1.orig	local-only
tmp/ll/m1.txt	local-only
tmp/ll/m2.orig	local-only
tmp/ll/m3.orig	local-only
tmp/ll/mir.js	local-only
tmp/ll/mir2.js	local-only
tmp/ll/mirror.txt	local-only
tmp/ll/mut.js	local-only
tmp/ll/mutation.md	local-only
tmp/ll/numstat.txt	local-only
tmp/ll/pairs.js	local-only
tmp/ll/pre.txt	local-only
tmp/ll/pre2.txt	local-only
tmp/ll/preflight.js	local-only
tmp/ll/recon.js	local-only
tmp/ll/recon2.js	local-only
tmp/ll/recon2.txt	local-only
tmp/ll/rev2-en-directive.txt	local-only
tmp/ll/selftest.txt	local-only
tmp/ll/tail-dump.txt	local-only
tmp/ll/tail2.txt	local-only
tmp/ll/tpl-cn.txt	local-only
tmp/ll/tpl-en.txt	local-only
tmp/ll/tpl.js	local-only
tmp/ll/tpl.txt	local-only
```
