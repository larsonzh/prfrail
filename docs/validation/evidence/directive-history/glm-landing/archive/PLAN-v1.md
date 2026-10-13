# 切片 `DIRECTIVE-GLM-LANDING` · ① 架构方案书（只读分析，v1）

**角色与载体**：架构师 ①；`deep-reasoner` / `DeepSeek V4 Pro (deepseek)`（Max）。本片产物即本文档正文，未写任何文件。
**取证基线**：HEAD = `e9b0875`（步骤 0 后、未提交）；`docs/DELIVERY_DIRECTIVE.md` 头版 v1.26（行 1–9 区），§0.3 最新行 = v1.26（CN/EN 均行 72）。

**两处引述更正（以原文为准，已逐字复核）**：
1. 起跑消息曾引述「§5.2 的『⑦ 三读在准则治理片不降级』」——**§5.2 现文无此字样**。实存口径为：§5.2 矩阵的 ⑦ 终审列 = **必**（行 417–421 区）、§5.2 说明段「**⑦ 独立终审不可省**」（行 424）、§5.2 ⑦ 三模型分档表（行 426–433 区）、§7.5 的 ⑦ 行（行 610）。本方案书 §5 据此派单，不再沿用被误引的表述。
2. **OB-83 处置列现文 = `已登记`**（CN 行 911 尾列：`**已登记**（落点 = 本条；调用分账与本片证据见 …）`），不是「临时授权」。本方案书 §2-C6 按现文给出变更前逐字。

---

## §1 目标与范围

**落地（八项目标逐条对应）**：
- ① §2.2 白名单增 `GLM-5.3 (glm)` 作 **⑤ 预审首选**；
- ② §2.2 白名单增 `GLM-5.3-Flash (glm)` 作 **批量辅助首选**；
- ③ §2.2 现表 `DeepSeek V4 Pro (deepseek)` 的 ⑤ 位改注为 **备用**；
- ④ §5.2 增 ⑤ 相关说明段（现文**无** ⑤ 说明段——见 §2-决策 4）；
- ⑤ §6.3 的 G5-b 白名单闭集扩入 `glm-5.3` / `glm-5.3-flash`（准则文本 + `g5b.js` 实现 + 夹具 + selftest 断言同步）；
- ⑥ §12.4 的 OB-83 处置列由 `已登记` 翻为**已落地（常设）**，并说明临时授权终止、升级触发关闭（形态裁决见 §2-决策 1）；
- ⑦ 批量辅助**限 4 类**（① 枚举/统计；② 长文档摘录与数字核对；③ 双语镜像语言侧检查；④ MANIFEST 清单草案），**明确排除「排序与 Top-N 择优」**（依据 = 上一片 BATCH-a 反例，`docs/validation/directive-glm-eval.md` §4.2，行 151–163 区）；
- ⑧ **① / ②③ 的模型位不动**（CN 行 222 架构师行 = `DeepSeek V4 Pro (deepseek)` 保持原样；②③⑧ 行不动）；**成本比较不在本片**（`directive-glm-eval.md` §5 的相对判据不成立结论与 OB-84 数据点均不进入本片改动）。

**明确不落地**：不触 §2.4 黑名单**正文**（行 251–262 区零 diff）、§5.2 矩阵本体（行 417–421 区零 diff）、§7.5 数值（行 603–615 区零 diff）、§1.5.1、§1.7、§7.7、角色集与载体（`.github/agents/**` 零 diff）、`internal/**`、`docs/CONTRACTS{,_EN}.md`、`schemas/**`、契约夹具、`.github/workflows/**`、`go.mod`；不新增判据（G5-b 只扩集合元素）；不追加 ⑦ 盲审；不写 `tmp/` 之外的新目录结构。

---

## §2 变更清单（逐文件逐节 + 逐字草稿）

### 决策 1 — (a)「首选 / 备用」在 §2.2 表中的表达（**推荐方案 + 反方案**）

- **现表结构（取证 CN 行 220–231）**：7 列 = `角色 | 调用点 | agentName（载体） | 限定名（调用串，唯一权威） | modelId（内部标识） | 与主控同模型？ | 工具能力`；v1.22 已对调限定名 / modelId 两列语义（CN 行 232 注记段，`DIRECTIVE-MODELID-MAP`）。⑦ 已有**三行先例**：`独立终审员` / `独立终审员（深度档）` / `独立终审员（平行备选）`（行 228–230），即「行标签承载档位、默认行无括注、分档判据在 §5.2」。
- **推荐方案**：**不加新列**。仿 ⑦ 三行先例——默认行「预审员」= 首选（改限定名 / modelId 为 GLM 两值），新增行「预审员（备用）」= 原 V4 Pro 三元组；「首选 / 备用」的选用口径写进 §5.2 新注记（单一权威，与 ⑦ 分档判据同构）。理由：① 不改表结构（7 列不变，G1-a 列数判据零风险）；② 限定名单元格保持**纯净调用串**（若在单元格内写「（首选）」会污染调用串——任务包照抄单元格将写出非法串，违反 R2.2 行 232–236）；③ 复用既有 ⑦ 行标签体例，零新概念。
- **反方案**：新增第 8 列「⑤ 首选/备用」——需改全部 9 行 + 表头，diff 面大，且「在限定名列旁加并列」有把注解混入权威列的漂移风险。**不推荐**。

### 决策 2 — (b) 两个模型的落位（⑤ 位 vs 独立批量槽位）

- `GLM-5.3 (glm)` → **⑤ 位**：修改现有「预审员」行（行 223）的限定名与 modelId 两列。
- `GLM-5.3-Flash (glm)` → **独立「批量辅助」行**（新行，调用点列 = `—`）。理由（硬约束）：§6.3 的 G5-b 文本自带「（与 §2.2 的 `modelId`（内部标识）列**一致**）」（CN 行 506）——若 `glm-5.3-flash` 只出现在注记段而非表行，该一致性声明即失效；给它一行是维持「G5-b 闭集 ≡ §2.2 modelId 列」的最小代价。批量辅助 ≠ ⑤（上片实测：⑤ 实验腿 = `glm-5.3`、批量 4 腿 = `glm-5.3-flash`，`directive-glm-eval.md` §4.1/§4.2）。

### 决策 3 — (c) 批量辅助的落点（独立行 + 注记段，**推荐**）

- **推荐**：§2.2 **独立行**（决策 2）+ 表下**新增注记段**承载「限 4 类 + 排除排序/Top-N」的规则正文。理由：行解决 `modelId` 一致性，注记段解决任务类限定；批量辅助**不进入 §5.2 矩阵、不占 §7.5 额度行**（非流水线角色）。
- **反方案**：只在 §2.2 写注记段、不加表行——必须同时改写 §6.3 G5-b 括注「与 §2.2 modelId 列一致」为「与 §2.2 表行及注记段一致」，多动一处判据文本且使闭集与列的对齐失去机械可读性。**不推荐**。

### 决策 4 — §5.2 的 ⑤ 说明段（现文无此段，取证）

- **取证**：§5.2（行 413–434）现文只有矩阵、说明（文档型切片）段（行 424）、⑦ 分档表（行 426–433）；**不存在** ⑤ 说明段。目标 ④ 的「同步」= **新增**一段（模型名**不写进 §5.2**，只作指针，保持 `glm` 串不出现在 §5.2——见 §4 白名单位置闭合）。此为该目标的唯一可行落法。

### 逐文件草稿（CN / EN 逐字；节后括注 = 现文行号锚点）

**C1 · §2.2 表行**（CN 行 222–223 / EN 行 222–223）：

变更前：
```
| 架构师 | ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `deepseek-v4-pro` | 否 | 只读 |
| 预审员 | ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `deepseek-v4-pro` | 否 | 只读 |
```
变更后（CN）：
```
| 架构师 | ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `deepseek-v4-pro` | 否 | 只读 |
| 预审员 | ⑤ | `deep-reasoner` | `GLM-5.3 (glm)` | `glm-5.3` | 否 | 只读 |
| 预审员（备用） | ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `deepseek-v4-pro` | 否 | 只读 |
| 批量辅助（非流水线角色） | — | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `glm-5.3-flash` | 否 | 只读 |
```
变更后（EN）：
```
| Architect | ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `deepseek-v4-pro` | No | Read-only |
| Pre-reviewer | ⑤ | `deep-reasoner` | `GLM-5.3 (glm)` | `glm-5.3` | No | Read-only |
| Pre-reviewer (backup) | ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `deepseek-v4-pro` | No | Read-only |
| Batch assistance (non-pipeline role) | — | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `glm-5.3-flash` | No | Read-only |
```
（插入点：现「预审员」行之后、「实现者」行之前；每行 8 个 `|`，与两侧镜像逐行同构。）

**C2 · §2.2 表下注记段**（插在「⑦ 三档分档」段之后、`> **已知限制（过渡态）**` 之前）：

CN：
```
**⑤ 首选 / 备用（2026-10-06 新增）**：⑤ 预审**首选** = `GLM-5.3 (glm)`（本表『预审员』行）、**备用** = `DeepSeek V4 Pro (deepseek)`（本表『预审员（备用）』行）；启用备用行（首选缺位 / 不可用）时，主控在任务包中记录实际限定名（R2.2）。

**批量辅助（2026-10-06 新增；非流水线角色，不进 §5.2 矩阵）**：`GLM-5.3-Flash (glm)` 为批量辅助**首选**，**仅限**以下四类任务：① 枚举 / 统计；② 长文档摘录与数字核对；③ 双语镜像语言侧检查；④ MANIFEST 清单草案。**明确排除**「排序与 Top-N 择优」——依据 = 切片 `DIRECTIVE-GLM-EVAL` 的 BATCH-a 反例（见 `docs/validation/directive-glm-eval{,_EN}.md`）。批量辅助的调用仍须显式传 `model`（R2.1 / R2.2），且深度 ≤ 1（R2.3）。
```
EN：
```
**⑤ preferred / backup (added 2026-10-06)**: the ⑤ pre-review **preferred** = `GLM-5.3 (glm)` (the "Pre-reviewer" row of this table) and the **backup** = `DeepSeek V4 Pro (deepseek)` (the "Pre-reviewer (backup)" row of this table); when the backup row is activated (the preferred is missing / unavailable), the master records the actual qualified name in the task package (R2.2).

**Batch assistance (added 2026-10-06; a non-pipeline role, not in the §5.2 matrix)**: `GLM-5.3-Flash (glm)` is the batch-assistance **preferred** model, **limited to** the following four task classes: ① enumeration / statistics; ② long-document excerpting and number checking; ③ language-side checks of bilingual mirrors; ④ MANIFEST list drafts. **Explicitly excluded**: "sorting and Top-N selection" — basis = the BATCH-a counter-example of slice `DIRECTIVE-GLM-EVAL` (see `docs/validation/directive-glm-eval{,_EN}.md`). A batch-assistance call must still pass `model` explicitly (R2.1 / R2.2), with depth ≤ 1 (R2.3).
```

**C3 · §5.2 新增 ⑤ 说明段**（插在「本表的成本与档位边界」段（行 433）之后）：

CN：
```
**⑤ 预审模型说明（2026-10-06 新增）**：⑤ 预审**首选**与**备用**的模型取值见 §2.2 的『预审员』与『预审员（备用）』两行；本矩阵的 ⑤ 列与 §7.5 的预审员额度**不变**（本片不动矩阵与数值）。批量辅助**非流水线角色、不进本矩阵**，其模型与四类限定见 §2.2 的批量辅助注记。
```
EN：
```
**⑤ pre-review model note (added 2026-10-06)**: the ⑤ pre-review **preferred** and **backup** model values are in the "Pre-reviewer" and "Pre-reviewer (backup)" rows of §2.2; the ⑤ column of this matrix and the pre-reviewer cap of §7.5 are **unchanged** (this slice does not touch the matrix or the figures). Batch assistance is **a non-pipeline role and does not enter this matrix**; its model and the four-class limit are in the batch-assistance note of §2.2.
```

**C4 · §6.3 的 G5-b 闭集枚举**（CN 行 506 / EN 行 506，同行内改）：

CN 变更前片段：
```
闭集 = `deepseek-v4-pro`、`deepseek-flash`、`mai-code-1.1-flash`、`gpt-5.3-codex`、`gpt-6.1-sol`、`gpt-6-luna`（与 §2.2 的 `modelId`（内部标识）列一致；
```
CN 变更后片段：
```
闭集 = `deepseek-v4-pro`、`deepseek-flash`、`mai-code-1.1-flash`、`gpt-5.3-codex`、`gpt-6.1-sol`、`gpt-6-luna`、`glm-5.3`、`glm-5.3-flash`（与 §2.2 的 `modelId`（内部标识）列一致；
```
EN 变更前片段：
```
the closed set = `deepseek-v4-pro`, `deepseek-flash`, `mai-code-1.1-flash`, `gpt-5.3-codex`, `gpt-6.1-sol`, `gpt-6-luna` (matching the §2.2 `modelId` (internal identifier) column;
```
EN 变更后片段：
```
the closed set = `deepseek-v4-pro`, `deepseek-flash`, `mai-code-1.1-flash`, `gpt-5.3-codex`, `gpt-6.1-sol`, `gpt-6-luna`, `glm-5.3`, `glm-5.3-flash` (matching the §2.2 `modelId` (internal identifier) column;
```

**C5 · §0.3 的 v1.26 行**（CN/EN 行 72；**扩展同一行**，原「不新增判据…」尾句系步骤 0 表述、随落地失效，整行重写）：

CN 变更前：
```
| v1.26 | 2026-10-06 | **台账补登（切片 `DIRECTIVE-GLM-LANDING` 的步骤 0；用户 2026-10-06 裁定）**：§12.4 登记 **OB-85**（试点继续顺延第 **5** 次，按 §1.7.4 / OB-37 元规则独立留痕，含四要素与如实标注）。**不新增判据、不改任何规则语义、不改 §6.3 判据作用域**。 | `DIRECTIVE-GLM-LANDING`（本片） |
```
CN 变更后：
```
| v1.26 | 2026-10-06 | **切片 `DIRECTIVE-GLM-LANDING` 落地（步骤 0 的用户裁定见下）**：步骤 0 登记 **OB-85**（试点继续顺延第 **5** 次，按 §1.7.4 / OB-37 元规则独立留痕，含四要素与如实标注）；本片落地——① §2.2 增 `GLM-5.3 (glm)`（⑤ 预审**首选**）与 `GLM-5.3-Flash (glm)`（**批量辅助首选**，限四类任务、排除排序与 Top-N 择优）两行，原预审员行改注 **备用**（① 与 ②③ 位不动）；② §5.2 增 ⑤ 预审模型说明段（**矩阵本体不动**）；③ §6.3 的 G5-b 白名单闭集扩入 `glm-5.3` / `glm-5.3-flash`，同步 `tools/gates` 实现与自检夹具；④ §12.4 的 OB-83 处置列翻 `已处置`（**已落地（常设）**，临时授权终止、升级触发关闭）；⑤ 附录 B.5 的 ⑤ 模板改首选限定名；⑥ 附录 C 增 C.0o。**不触 §2.4 黑名单正文、§5.2 矩阵、§7.5 数值、角色集与载体、`internal/**`、CONTRACTS / schemas / 契约夹具 / workflows / `go.mod`。** | `DIRECTIVE-GLM-LANDING`（本片） |
```
EN 变更前：
```
| v1.26 | 2026-10-06 | **Ledger registration (step 0 of slice `DIRECTIVE-GLM-LANDING`; user ruling 2026-10-06)**: §12.4 registers **OB-85** (the pilot deferral continues for the **5th** time, independently recorded per §1.7.4 and the OB-37 meta-rule, with the four elements and the truthful annotation). **No new criterion, no change to any rule semantics, no change to the §6.3 criterion scope**. | `DIRECTIVE-GLM-LANDING` (this slice) |
```
EN 变更后：
```
| v1.26 | 2026-10-06 | **Slice `DIRECTIVE-GLM-LANDING` landing (the step-0 user ruling is below)**: step 0 registers **OB-85** (the pilot deferral continues for the **5th** time, recorded independently per §1.7.4 / the OB-37 meta-rule, with the four elements and the truthful annotation); this slice lands - ① §2.2 gains `GLM-5.3 (glm)` (⑤ pre-review **preferred**) and `GLM-5.3-Flash (glm)` (**batch-assistance preferred**, limited to four task classes, sorting and Top-N selection excluded), and the former pre-reviewer row is annotated **backup** (① and the ②③ slots are untouched); ② §5.2 gains a ⑤ pre-review model note (**the matrix body is untouched**); ③ the §6.3 G5-b whitelist closed set gains `glm-5.3` / `glm-5.3-flash`, with the `tools/gates` implementation and self-test fixtures in sync; ④ the §12.4 OB-83 disposition flips to `Resolved` (**Landed (standing)**, the temporary authorisation is terminated and the upgrade trigger is closed); ⑤ appendix B.5's ⑤ template moves to the preferred qualified name; ⑥ appendix C gains C.0o. **No change to the §2.4 blacklist body, the §5.2 matrix, the §7.5 figures, the role set and carriers, `internal/**`, CONTRACTS / schemas / contract fixtures / workflows / `go.mod`.** | `DIRECTIVE-GLM-LANDING` (this slice) |
```
（G6-2 自护：两行均无 `⑥×n / ⑦×n` 成本向量——成本只进报告，F-9 先例；G6-4：无行数/文件数/合计字面量。）

**C6 · §12.4 的 OB-83 处置列**（CN/EN 行 911，同行内改；现象列**不动**）：

CN 变更前：
```
**已登记**（落点 = 本条；调用分账与本片证据见 `docs/validation/directive-glm-eval{,_EN}.md`）
```
CN 变更后：
```
**已处置**（**已落地（常设）**：`DIRECTIVE-GLM-LANDING`（v1.26）已将 `GLM-5.3 (glm)` / `GLM-5.3-Flash (glm)` 落入 §2.2 白名单 ⇒ **临时授权终止、升级触发关闭**——此后对两模型的调用按 §2.2 / §2.4 常规口径执行，不再走 §1.7.4；调用分账与本片证据见 `docs/validation/directive-glm-eval{,_EN}.md`）
```
EN 变更前：
```
**Registered** (landing point = this entry; the call breakdown and this slice's evidence are in `docs/validation/directive-glm-eval{,_EN}.md`)
```
EN 变更后：
```
**Resolved** (**Landed (standing)**: `DIRECTIVE-GLM-LANDING` (v1.26) has whitelisted `GLM-5.3 (glm)` / `GLM-5.3-Flash (glm)` in §2.2 ⇒ **the temporary authorisation is terminated and the upgrade trigger is closed** — from now on, calls to either model follow the regular §2.2 / §2.4 reading and no longer go through §1.7.4; the call breakdown and this slice's evidence are in `docs/validation/directive-glm-eval{,_EN}.md`)
```
**形态说明（必须写进报告）**：目标 ⑥ 的字面「已落地（常设）」首词 `已落地` **不在 G6-6 状态闭集** `观察中 / 待办 / 待议 / 已处置 / 已裁决 / 已登记`（§6.3 的 G6-6，CN 行 537–538 区；EN 闭集 `Observing / To do / To be discussed / Resolved / Ruling / Registered`）⇒ 处置列**首状态词取 `已处置`**、括注承载「已落地（常设）」。若用户坚持字面首词 `已落地`，须扩 G6-6 闭集——属判据语义变更、超本片八项目标 ⇒ 按 §1.7.2 升级用户裁决（见 §7-R1）。

**C7 · 附录 B.5**（CN/EN 行 1057–1058）：

CN 变更前：
```
角色：预审员。模型：DeepSeek V4 Pro (deepseek)。只读，不得修改任何文件。
载体：runSubagent(agentName='deep-reasoner', model='DeepSeek V4 Pro (deepseek)')
```
CN 变更后：
```
角色：预审员。模型：GLM-5.3 (glm)（⑤ 首选；备用 = DeepSeek V4 Pro (deepseek)，见 §2.2）。只读，不得修改任何文件。
载体：runSubagent(agentName='deep-reasoner', model='GLM-5.3 (glm)')
```
EN 变更前：
```
Role: pre-reviewer. Model: DeepSeek V4 Pro (deepseek). Read-only; must not modify any file.
Vehicle: runSubagent(agentName='deep-reasoner', model='DeepSeek V4 Pro (deepseek)')
```
EN 变更后：
```
Role: pre-reviewer. Model: GLM-5.3 (glm) (⑤ preferred; backup = DeepSeek V4 Pro (deepseek), see §2.2). Read-only; must not modify any file.
Vehicle: runSubagent(agentName='deep-reasoner', model='GLM-5.3 (glm)')
```
（载体不变 ⇒ R2.6 的「同载体两次独立调用」不受影响；B.1 架构师模板不动，① 位不动。）

**C8 · §5.1 流水线图**（CN/EN 行 399，同行内改，`glm` 串不进 §5.1）：
CN：`→ ⑤预审(V4 Pro)` ⇒ `→ ⑤预审(首选)`；EN：`→ ⑤pre-review (V4 Pro)` ⇒ `→ ⑤pre-review (preferred)`。

**C9 · §7.1 三层审查表**（CN/EN 行 559，同行内改）：
CN：`| ⑤ | 预审员（V4 Pro） | …` ⇒ `| ⑤ | 预审员（首选） | …`；EN：`| ⑤ | Pre-reviewer (V4 Pro) | …` ⇒ `| ⑤ | Pre-reviewer (preferred) | …`（其余单元格逐字不动）。

**C10 · §7.2.3 三级回退**（CN/EN 行 591，同行内改）：
CN：`① 退回 ⑤ 请 V4 Pro 补扫一轮（预审加一轮）` ⇒ `① 退回 ⑤ 请预审员补扫一轮（预审加一轮）`；EN：`① fall back to ⑤ and ask V4 Pro for one extra scan round` ⇒ `① fall back to ⑤ and ask the pre-reviewer for one extra scan round`。

**C11 · 附录 C.0o（新增块；体例模板 = C.0n，CN 行 1403–1419 / EN 行 1403–1419）**：

CN：
```
### C.0o GLM 白名单落地切片 `DIRECTIVE-GLM-LANDING`

**为何必需**：OB-83 现象列末句「同类情形第三次须立项为切片正式修订本准则」即本片路径——GLM 两模型的临时评估授权在 `DIRECTIVE-GLM-EVAL`（OB-83）已耗尽先例空间；本片按评估结论把 `GLM-5.3 (glm)`（⑤ 首选）与 `GLM-5.3-Flash (glm)`（批量辅助首选，限四类、排除排序与 Top-N 择优）落入 §2.2 白名单并同步 G5-b 闭集；OB-83 处置列随本片翻 `已处置`（已落地（常设））。

| 字段 | 内容 |
|---|---|
| 标识 | `DIRECTIVE-GLM-LANDING`（v1.26） |
| 状态 | ✅ **已完成**（2026-10-06；变更与审查证据见 `docs/validation/directive-glm-landing{,_EN}.md`） |
| 目标 | ① §2.2 增 GLM 两行（⑤ 首选 / 批量辅助首选）并把原预审员行改注备用；② §5.2 增 ⑤ 预审模型说明段；③ G5-b 闭集扩 `glm-5.3` / `glm-5.3-flash` 并同步实现与夹具；④ OB-83 处置列翻 `已处置`（已落地（常设））；⑤ 附录 B.5 同步；⑥ 本台账块 |
| 触发 | 用户 2026-10-06 启动声明（`[SLICE]`，准则治理 / 判据本体变更片、探针 0；试点顺延第 5 次 = OB-85） |
| 规模 | M（纯 `.md` 与 `tools/gates` 的 JS；含判据语义） |
| 交付物 | `docs/DELIVERY_DIRECTIVE{,_EN}.md` 的 v1.26（§0.3 / §2.2 / §5.2 / §6.3 / §7.1 / §7.2.3 / §12.4 / 附录 B.5 / 附录 C.0o）；`tools/gates/lib/criteria/g5b.js` 与夹具与 `selftest.js`；本报告对；自由度清单 |
| 验收证据 | `gate.js --all --scope=tree` 与 `--scope=index` 全绿；`gate.js --selftest` 全绿；G5-b 正反夹具（两枚 GLM `modelId` PASS、非白名单 FAIL）；CN / `_EN` 同位由 G3(a) 与 G3(b) 校验 |
| 边界 | 不触 §2.4 黑名单正文、§5.2 矩阵、§7.5 数值、§1.5.1、§1.7、§7.7、角色集与载体、`.github/agents/**`、`internal/**`、`docs/CONTRACTS{,_EN}.md`、`schemas/**`、契约夹具、`.github/workflows/**`、`go.mod`；成本比较不在本片 |
| 分级 | `[SLICE]` |
| 探针额度 | 无 |
| 自由度清单路径 | `docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md` |
| 附注 | ⑦ = `GPT-5.3-Codex (copilot)` 1 终审 + 1 复审、不追加盲审（本片不触 §7.7 必抽三语义；按比例抽查判定在启动前记录理由，见 OB-82）；⑤ 首次落地复审可用 `GLM-5.3 (glm)`；成本向量只进报告（G6-2 口径，F-9 先例） |
```
EN（同 18 行，逐行同位）：
```
### C.0o GLM whitelist landing slice `DIRECTIVE-GLM-LANDING`

**Why required**: the OB-83 phenomenon column's last sentence - "the third instance of the same kind must be chartered as a slice that formally revises this directive" - is this slice's path: the temporary evaluation authorisation for the two GLM models already exhausted its precedent space in `DIRECTIVE-GLM-EVAL` (OB-83); this slice lands `GLM-5.3 (glm)` (⑤ preferred) and `GLM-5.3-Flash (glm)` (batch-assistance preferred, limited to four classes, sorting and Top-N selection excluded) into the §2.2 whitelist per the evaluation result and synchronises the G5-b closed set; the OB-83 disposition flips to `Resolved` (Landed (standing)) with this slice.

| Field | Content |
|---|---|
| ID | `DIRECTIVE-GLM-LANDING` (v1.26) |
| Status | ✅ **Completed** (2026-10-06; changes and review evidence in `docs/validation/directive-glm-landing{,_EN}.md`) |
| Goal | (1) §2.2 gains two GLM rows (⑤ preferred / batch-assistance preferred) and the former pre-reviewer row is annotated as backup; (2) §5.2 gains a ⑤ pre-review model note; (3) the G5-b closed set gains `glm-5.3` / `glm-5.3-flash` with the implementation and fixtures in sync; (4) the OB-83 disposition flips to `Resolved` (Landed (standing)); (5) appendix B.5 is synchronised; (6) this ledger block |
| Trigger | the user 2026-10-06 start declaration (`[SLICE]`, a directive-governance / criterion-body slice, probes 0; 5th pilot deferral = OB-85) |
| Size | M (pure `.md` plus the JS under `tools/gates`; carries criterion semantics) |
| Deliverables | v1.26 of `docs/DELIVERY_DIRECTIVE{,_EN}.md` (§0.3 / §2.2 / §5.2 / §6.3 / §7.1 / §7.2.3 / §12.4 / appendix B.5 / appendix C.0o); `tools/gates/lib/criteria/g5b.js` with fixtures and `selftest.js`; this report pair; the freedom list |
| Acceptance evidence | `gate.js --all --scope=tree` and `--scope=index` green; `gate.js --selftest` green; G5-b positive / negative fixtures (both GLM `modelId` values PASS, non-whitelist FAIL); CN / `_EN` positional parity checked by G3(a) and G3(b) |
| Boundary | no change to the §2.4 blacklist body, the §5.2 matrix, the §7.5 figures, §1.5.1, §1.7, §7.7, the role set and carriers, `.github/agents/**`, `internal/**`, `docs/CONTRACTS{,_EN}.md`, `schemas/**`, contract fixtures, `.github/workflows/**`, `go.mod`; the cost comparison is not in this slice |
| Tier | `[SLICE]` |
| Probe budget | none |
| Freedom-list path | `docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md` |
| Note | ⑦ = `GPT-5.3-Codex (copilot)` 1 final review + 1 re-review, no blind review added (this slice does not touch the three mandatory-blind semantics of §7.7; the proportional-sampling judgement is recorded with its reason before start, see OB-82); the ⑤ landing re-review may use `GLM-5.3 (glm)`; cost vectors go only into the report (the G6-2 reading, the F-9 precedent) |
```

**C12 · §12.4 新增观察项（候选号 OB-86；②a 落盘时按台账顺序占用，主控在 ④ 定号）**：

CN：
```
| OB-86 | 2026-10-06 | **§2.2 的 ⑤ 首选 / 备用标注与批量辅助四类限定无机械判据覆盖**（本片新增的文档语义）：G1–G8 均不解析 §2.2 表行与注记段的模型指派；候选覆盖与 OB-29 同族（附录 B / §2.2 的 `model=` 取值集合 ⊆ 限定名列的候选判据）。本片不新增判据，由 ⑤⑥⑦ 审查链承担。 | 观察中 |
```
EN：
```
| OB-86 | 2026-10-06 | **The §2.2 ⑤ preferred / backup annotation and the batch-assistance four-class limit have no mechanical criterion coverage** (document semantics newly added by this slice): none of G1–G8 parses the model assignment of the §2.2 table rows or note paragraphs; the candidate coverage is the same family as OB-29 (the candidate criterion "the `model=` value set of appendix B / §2.2 ⊆ the qualified-name column"). This slice adds no criterion; the ⑤⑥⑦ review chain carries it. | Observing |
```

**(d) §2.4 括注是否需要变：不需要。** 取证：CN 行 253 括注与行 261「白名单闭集」均为**指针式**表述（「闭集按 §2.2 的 `modelId`（内部标识）列定义」），不枚举成员 ⇒ 新增白名单成员自动生效，§2.4 全节**零 diff**（黑名单正文行 251–262 区零改动，由 §4 检查 1 机械钉住）。

---

## §3 G5-b 判据与夹具设计

**闭集扩法**：只为集合加元素，**不新增子判据**。`tools/gates/lib/criteria/g5b.js` 的 `CLOSED_SET`（现 6 元，行 18–24）尾部追加两行：

```js
  'gpt-6-luna',
  'glm-5.3',
  'glm-5.3-flash',
];
```
判定逻辑（frontmatter 限定解析、全键收集、大小写敏感逐字相等、trim、引号包裹不算、空值让位 G5-a②、键缺失 PASS）与新增成员**正交**，不改动。

**新增夹具（4 个；落 `tools/gates/testdata/fixtures/`，`.txt` 后缀约定不变）**：
| 夹具 | frontmatter `model` 值 | 期望 |
|---|---|---|
| `g5b-closed-set-legal-glm.txt` | `glm-5.3` | PASS |
| `g5b-closed-set-legal-glm-flash.txt` | `glm-5.3-flash` | PASS |
| `g5b-illegal-glm-case.txt` | `GLM-5.3` | FAIL（大小写敏感） |
| `g5b-illegal-glm-qualified.txt` | `GLM-5.3 (glm)` | FAIL（限定名 ≠ modelId，R2.2 行 232–236 边界） |

夹具正文沿用现例形（`g5b-closed-set-legal-luna.txt`：`---` frontmatter + 一行 "Probe body. This file is test data, not a real agent role file."），**不含** G6/G2 会命中的内容 ⇒ T13 全语料锚（`EXPECTED_CORPUS_FAILS`，selftest.js 行 150–164）**不变**——G5-b 不在 `ALL_SUBCHECKS`（G6/G7/G8）内，新 FAIL 夹具不进入该锚。

**selftest.js 改动**：FIXTURES 注册表（行 75–87）追加 4 条（virtual 路径形如 `.github/agents/fixture-g5b-closed-set-legal-glm.agent.md`）；新增变异常量 `G5B_GLM_MEMBER = "  'glm-5.3',\n"`、`G5B_GLM_FLASH_MEMBER = "  'glm-5.3-flash',\n"`（替换串唯一：`glm-5.3` 行不包含 `glm-5.3-flash` 行的引号逗号尾缀，无串扰）；新增 **7 条 check**（编号 T61–T67，续现有 T60）：

| 断言 | 内容 | 删掉守卫后变红的断言 |
|---|---|---|
| T61 | 两枚新成员夹具均 PASS 且 findings=0 | 守卫 = `CLOSED_SET` 两新成员；删除 ⇒ T61 变红 |
| T62 | 变异删除 `'glm-5.3-flash',` 成员 ⇒ `-glm-flash` 夹具 FAIL | 加回成员 ⇒ T62 变红（期望 FAIL 变 PASS） |
| T63 | 变异还原：`g5b.js` 字节一致 + 缓存模块仍 PASS | 还原失败/字节漂移 ⇒ T63 变红 |
| T64 | 变异删除 `'glm-5.3',` 成员 ⇒ `-glm` 夹具 FAIL | 加回成员 ⇒ T64 变红 |
| T65 | 变异还原：字节一致 + 缓存模块仍 PASS | 还原失败 ⇒ T65 变红 |
| T66 | `g5b-illegal-glm-case` FAIL 且 detail 含 `model=GLM-5.3` | 判据改成大小写不敏感 ⇒ T66 变红 |
| T67 | `g5b-illegal-glm-qualified` FAIL | 有人把限定名串加进闭集 ⇒ T67 变红（钉 R2.2 边界） |

**变更后预期断言总数**：本任务书给定基线 `SELFTEST PASS 328/328` ⇒ **预期 `SELFTEST PASS 335/335`**（+7）。④ 落地时以 `--selftest` 实测输出为准：若实测基线 ≠ 328，按「实测基线 + 7」修正报告读数并如实注明差异，**不得**改判据来凑数。

**必须回答的问题**：§2.2 表内「首选 / 备用」标注（及批量辅助四类限定）**未被任何门禁判据覆盖**——G5-b 只读 `.github/agents/*.agent.md` 的 frontmatter（g5b.js 行 35–55）；G1-a 只判表格结构；G3 只判双语对称；G6-5/G6-6 只判节号与状态词 ⇒ 属**「文档语义、无机械判据」**。**评估结论**：不新增判据（本片边界），建议**登记观察项 OB-86（观察中）**并明示由 ⑤⑥⑦ 审查链承担；候选覆盖与 OB-29 同族（行 857）。若用户要求本片新增判据覆盖它，需论证必要性并扩片（本方案书不推荐，理由：成本高于收益，且 OB-29 已为同族候选留位）。

---

## §4 边界自检清单（④ 可机械执行）

1. **§2.4 黑名单正文零改动**（内容锚定提取，免疫行号漂移）：
   `git show HEAD:docs/DELIVERY_DIRECTIVE.md | awk '/^### 2.4 /{f=1} f{print} /^### 2.5 /{f=0}'` 与工作树同命令输出 `diff` ⇒ 空；`_EN` 同法。
2. **§5.2 矩阵本体零改动**：锚定 `| 切片类型 |` 表头行至 `**说明（文档型切片` 前，HEAD vs 工作树 `diff` ⇒ 空（CN+EN）。
3. **§7.5 数值零改动**：锚定 `### 7.5` 至 `### 7.6`，HEAD vs 工作树 `diff` ⇒ 空（CN+EN）。
4. **禁触文件零改动**：`git diff --stat -- .github/agents/ internal/ docs/CONTRACTS.md docs/CONTRACTS_EN.md schemas/ .github/workflows/ go.mod` ⇒ 空输出。
5. **`glm` 出现位置闭合**：`git diff -U0 -- docs/DELIVERY_DIRECTIVE.md docs/DELIVERY_DIRECTIVE_EN.md | grep '^+' | grep -in glm` 的每条命中必须落在六区 {§0.3 v1.26 行、§2.2 表行/注记段、§6.3 G5-b 行、§12.4 OB-83 处置列与 OB-86 行、附录 B.5、附录 C.0o}；反向断言 = 落盘后 `grep -n 'glm'` 的行号全集 ∈ 六区行号区间。
6. **`g5b.js` 改动面**：`git diff -- tools/gates/lib/criteria/g5b.js` 只允许 `CLOSED_SET` 追加两行（+2 行、0 删除）。
7. **selftest 全绿**：`node tools/gates/gate.js --selftest` ⇒ `SELFTEST PASS 335/335`（实测为准，见 §3）。
8. **三作用域门禁**：`node tools/gates/gate.js --all --scope=tree` 与 `--scope=index` ⇒ 退出码 0、`TOTAL_FAIL=0`；G1-b（头版 v1.26 = 日志末行 v1.26）、G1-a（新表行无 `||` 拼接、列数 8）、G3-a.1/a.2/b（行数对称 + 五形状）、G4a/G4b（BOM+LF；差集预检 `unknown=[]`）、G6-5（v1.26/C.0o 的 § 引用可解析）、G6-6（OB-83/OB-86 首状态词 ∈ 闭集）由 `--all` 承担；`--scope=ci` 由 CI 的 Gates 步在推送后执行（§6.4，CI 配置零改动）。
9. **G4b 新汉字预检**：本方案书草稿用字均为语料内字（择优/枚举/摘录/镜像/辅助/草案等已在 HEAD 全库 `.md` 语料，取证见下）；②a 落盘前先跑 G4b，差集非空 ⇒ **只改措辞、禁止扩白名单**（`tools/gates/cjk-newwords.txt` 零改动）。
10. **夹具后缀约定**：新增夹具全部 `.txt`（T13 的 FIXTURE_DISK_MUST_NOT_BE_MD 断言自动守）。

---

## §5 任务分解（②a → ②b → ③ → ④ → ⑤ → ⑥ → ⑦ → ⑧）

**顺序依据**：§1.4 协议先行 ⇒ ②a（准则文本 = 本片「契约」）先于 ②b（判据实现）先于 ③（测试）。

**②a 实现者**（`DeepSeek V4.1 Flash (deepseek)` / `prfrail-implementer`）：
- 输入：本方案书 §2 的 C1–C12 草稿 + 行号锚点 + 现文行 72/223/399/424/506/559/591/608/610/857/897/911/1057/1403 摘录。
- 期望输出：`docs/DELIVERY_DIRECTIVE.md`（**先改**）与 `docs/DELIVERY_DIRECTIVE_EN.md`（**后改**）diff——C1/C2/C3/C4/C5/C6/C7/C8/C9/C10/C11/C12；编码 BOM+LF。
- 验收：逐条草稿逐字比对；CN/EN 增删行数相等（G3-a.1）；G1-a/G6-5/G6-6 预跑绿；**不改** §2.4 正文 / §5.2 矩阵 / §7.5 数值。
- 禁止：自行扩范围（如改 §7.5 数值、改 G6-6 闭集、改黑名单正文）；`model=` 串只写限定名。

**②b 实现者**（`DeepSeek V4.1 Flash (deepseek)` / `prfrail-implementer`）：
- 输入：`tools/gates/lib/criteria/g5b.js` 现文 + §3 的闭集扩法。
- 期望输出：`CLOSED_SET` 追加 `'glm-5.3'`、`'glm-5.3-flash'` 两行（+2/−0），文件其余字节不动。
- 验收：`git diff --numstat` = 2/0/0；`node --check` 通过。

**③ 测试工程师**（`DeepSeek V4.1 Flash (deepseek)` / `prfrail-tester`）：
- 输入：§3 夹具清单与 T61–T67 断言定义、`selftest.js` 的 FIXTURES 注册表（行 75–87）与 G5B 变异辅助（行 666–730）。
- 期望输出：4 个新夹具 + FIXTURES 4 条注册 + T61–T67 断言 + 覆盖矩阵（编号/测试名/类型/对应测试点）。
- 验收：`--selftest` 全绿；**可证伪性逐条**（§3 表的「删守卫 → 变红」列全部实测成立，变异后字节还原、缓存判定成立）；`EXPECTED_CORPUS_FAILS` 无变动。

**④ 主控集成 + 门禁**（主控本体；不派子代理）：
- 落盘三重校验（前后 diff 行数一致 / `.md` 编码门禁 / 落盘后聚焦门禁）+ §4 全 10 项 + `--scope=tree`/`--scope=index`/`--selftest` + 原始输出落 `tmp/` 备审（§6.2 路径交付）。
- 集成后产出自由度清单 `docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md`（⑥ 后、⑦ 前，§1.7.3 行 202）。

**⑤ 预审员**（`DeepSeek V4 Pro (deepseek)`；**首次落地复审可用 `GLM-5.3 (glm)`**——注意：本片 ⑤ 首轮仍传 V4 Pro，因「⑤ 首选 = GLM」规则在本片 ②a 才落盘、对本片自身的派单不追溯适用，主控按任务包显式传 model 即为合规）：
- 输入：方案书 + ②a/②b/③ 的 diff **路径** + §2.2/§6.3 G5-b/§12.4 OB-83 契约原句（R2.6 调用隔离，不给 ① 推理）。
- 验收：`PRE-REVIEW: PASS|PASS WITH FIXES|FINDINGS` + 偏离项表 + 反例三类；额度 = 1 + 整改后复审 1（§7.5 行 608）。

**⑥ 独立扫描**（`MAI-Code-1.1-Flash (copilot)` / `independent-reviewer`）：Section A–E；有效调用 ≤ 3（§7.2.2）。**不得跳过**——本片改门禁判据（G5-b 闭集）⇒ 命中 §7.2.1 反面样例（行 573–577）同款，③ 不成立。

**⑦ 独立终审**（`GPT-5.3-Codex (copilot)` / `independent-reviewer`，**不降级**；额度 = 1 终审 + 1 复审，§7.5 行 610）：
- **不追加盲审**：本片不触 §7.7 必抽三语义（所有权 / 停机 / 身份）；「按比例抽查」（每 4 片抽 1）的判定与理由须在**启动前记录**（OB-82 行 914–916 的口径，锚点待议）——主控在回执中如实写「本片判定 = 不抽 + 理由（上片 `DIRECTIVE-EVIDENCE-SCOPE` 刚被抽中、本片不触必抽语义）」。
- 输入：切片定义 + 最终 diff 路径 + 契约原句 + **自由度清单**（§1.7.3 行 202）+ 首轮不附清单（§7.3）。

**⑧ 文档工程师**（`DeepSeek V4.1 Flash (deepseek)` / `prfrail-documenter`）：报告对 `docs/validation/directive-glm-landing{,_EN}.md`（§11.3 结构 + `artifacts` 围栏块，G7）+ 台账完成行（只写报告路径，§11.2 行 772）+ `DEV_PLAN` 回写标注「不适用 + 理由」（OB-12 行 866 口径：准则治理片回写附录 C，不回写 `DEV_PLAN`）+ C.0o 状态行。

---

## §6 自由度预先识别（体例同 `DIRECTIVE-EVIDENCE-SCOPE-freedom-list.md`，行 9–14）

| 项 | 内容 | 依据 | 判定 |
|---|---|---|---|
| F-1 | 首选/备用用**行标签**表达、不加新列（§2-决策 1） | §1.7.1 实现表达选择；⑦ 三行先例（行 228–230） | 授权内 |
| F-2 | §5.1（行 399）/§7.1（行 559）/§7.2.3（行 591）三处「V4 Pro」标签同步为「首选/预审员」 | 一致性义务（OB-69 行 897：「同步本片所改判据的过期描述属一致性义务，非扩范围」） | 授权内（三处均只改标签、不动正文语义；`glm` 串不进入这三节） |
| F-3 | OB-83 处置列首状态词取 `已处置`、括注承载「已落地（常设）」 | G6-6 状态闭集（§6.3 行 537–538）机械强制 | 授权内；字面「已落地」首词需扩 G6-6 闭集 ⇒ 见 §7-R1 升级路径 |
| F-4 | 批量辅助落 §2.2 独立行（调用点 `—`）+ 表下注记段 | 维持「G5-b 闭集 ≡ §2.2 modelId 列」（行 506） | 授权内 |
| F-5 | §0.3 v1.26 行整行重写（原「不新增判据…」尾句随落地失效） | 同一 v1.26 行内扩展（任务书 §2 要求）；历史行不追溯（§11.1） | 授权内 |
| F-6 | 成本向量（⑥×n/⑦×n）只进报告、不进 §0.3/C.0o | G6-2 双源冲突；F-9 先例（`DIRECTIVE-EVIDENCE-SCOPE-freedom-list.md` 行 16） | 实现选择 |
| F-7 | OB-86 观察项登记（候补编号，主控 ④ 时按台账顺序定号） | §12.3 观察项登记；G6-6 首状态词 `观察中` | 授权内（本片不新增判据，仅登记） |
| F-8 | 夹具命名与 T61–T67 编号 | §1.7.1 单文件实现选择 | 实现选择 |
| F-9 | 批量辅助额度/成本口径不在本片定义（注记段只写任务类限定） | 成本比较已顺延（目标 ⑧）；§7.5 不新增行 | 授权内（本片零预算条款改动） |
| F-10 | C.0n 等历史块、C.1 对照表（行 1426 的「§2.1 V4 Pro」）不追溯改写 | §11.1 历史不重写；C.1 不在 G6-5 目标区域（§6.3 行 523–526） | 授权内 |

---

## §7 风险与未决项

- **R1（中）**：OB-83 字面「已落地（常设）」与 G6-6 闭集冲突。处置已定（F-3）。**升级路径**：若用户裁定必须字面首词 `已落地` ⇒ 停止 + §1.7.2 升级，扩 G6-6 闭集另立项（判据语义变更，不在本片八项目标内）。
- **R2（低）**：G3(a.2) 五形状（加粗标记/竖线计数）CN/EN 漂移。处置：草稿两侧加粗对数已逐条对齐（§2 各草稿）；②a/⑧ 每步跑 `--check=G3-a.2`。
- **R3（低）**：G4b 差集非空（本方案书全部用字经语料核验为既有字；若 ② 措辞变化引入新字）⇒ 改措辞，**不扩** `cjk-newwords.txt`。
- **R4（低）**：任务书基线 328 与实测基线如有出入 ⇒ 以实测为准、报告如实修正（§3 已声明）。
- **R5（低）**：§7.7 按比例抽查锚点未定（OB-82 `待议`）⇒ 本片判定「不抽」并启动前记录理由；若 ⑦ 或用户要求补盲审 ⇒ 升级用户裁决（额度 +1，§7.5 行 610）。
- **R6（高，边界越界）**：若 ② 执行中发现必须触及边界外内容（如 §7.5 数值、§5.2 矩阵、黑名单正文、G6-6 闭集）⇒ **停下 + 升级用户裁决**（§1.7.2 行 184–193），不得自行扩范围。
- **R7（未决）**：批量辅助的**调用额度与成本上限**本片未定义——需要**用户裁定**或后续成本片评估（依据：§7.5 无批量辅助行；本片目标 ⑧ 明示成本比较顺延；`directive-glm-eval.md` §5 的相对判据不成立 + OB-84 数据点）。在裁定前，批量辅助调用按 §10/§7.6 常规授权与探针记账执行，不引用任何推定额度。
- **R8（低）**：⑤ 首选规则对现片无追溯力（F-2 相关）：本片 ⑤ 派单仍显式传 V4 Pro，规则自本片收尾后对后续切片生效。

---

## 执行摘要（主控据此生成 ② 的 ③…⑦ 任务包）

1. 版本行 v1.26 **扩展**（CN+EN 行 72）：登记 OB-85 保留、追加六项落地；②a 先改 CN 再改 EN。
2. §2.2：改「预审员」行为 GLM-5.3（首选）；新增「预审员（备用）」= V4 Pro 与「批量辅助（非流水线角色）」= GLM-5.3-Flash 两行；表下新增 ⑤ 首选/备用 + 批量辅助限四类/排除排序与 Top-N 的注记段。
3. §5.2 新增 ⑤ 预审模型说明段（只指针、不含模型名）；§6.3 G5-b 闭集扩 `glm-5.3`/`glm-5.3-flash`；§2.4 **零 diff**。
4. §12.4：OB-83 处置列 `已登记` → `已处置（已落地（常设）；临时授权终止、升级触发关闭）`；新增 OB-86 观察项（无机械判据覆盖，观察中）。
5. 附录 B.5 改 `GLM-5.3 (glm)`；§5.1/§7.1/§7.2.3 三处「V4 Pro」标签同步（一致性义务，OB-69）。
6. ②b：`g5b.js` 闭集 +2 行（+2/−0）；③：4 新夹具 + T61–T67 七断言。
7. selftest 预期 328 → **335/335**（实测为准，偏差如实记录）。
8. ④：§4 十项边界自检 + tree/index/selftest 三作用域；CI scope 推送后观察。
9. ⑤ = V4 Pro（首次落地复审可用 GLM-5.3）；⑥ = MAI；⑦ = Codex 1+1、**不追加盲审**、不降级。
10. 成本向量只进报告（G6-2/F-9）；批量辅助额度未定 ⇒ 用户裁定，本片不立预算条款。
11. 自由度清单 F-1…F-10 预登记，⑦ 输入包必含；R1/R6 触发即停 + §1.7.2 升级。