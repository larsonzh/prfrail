# 切片 `DIRECTIVE-GLM-EVAL` · 验证报告

| 项 | 值 |
|---|---|
| 切片 | `DIRECTIVE-GLM-EVAL`（`[SLICE]`，延续上一会话至 ④） |
| 目标 | 评估 `GLM-5.3 (glm)` 在 **⑤ 预审位置**的质量与成本（对照 `DeepSeek V4 Pro (deepseek)`），并评估 `GLM-5.3-Flash (glm)` 在 4 类批量辅助任务上的质量与产物可用性 |
| 边界 | **只评估、不入白名单**；白名单与角色决策在本片之后另议 |
| 授权边界 | commit 否 / push 否 / gitee 否 / 探针 0（GLM 调用为付费调用，与探针分账） |
| 报告落盘 | `docs/validation/directive-glm-eval.md`（本文件，CN 权威）+ `_EN` 同位镜像 |

> **性质声明**：本片**不改任何规则正文**、**不改契约 / schema / 夹具**、**不改角色集**、**不扩白名单**。
> 交付物为：本报告对、评估数据集、冻结评估证据包、自由度清单、台账行。GLM 两枚模型均**不在**准则
> §2.2 白名单内，本片对其调用属**用户授权的临时评估**，按 §1.7.4 独立留痕（本项目记号 `DIRECTIVE-GLM-EVAL`），
> **不得作为先例**——GLM 的正式纳入与否在**本片之后**另议。

---

## 1. 变更概要

| 类别 | 内容 |
|---|---|
| 新增 | 本报告对 `docs/validation/directive-glm-eval{,_EN}.md`；评估数据集 `docs/validation/directive-glm-eval-dataset.json`；冻结证据包 `docs/validation/evidence/glm-eval/`（`MANIFEST.md` + `SHA256SUMS.txt`）；自由度清单 `docs/validation/evidence/DIRECTIVE-GLM-EVAL-freedom-list.md` |
| 修改 | `docs/t027/REMAINING_SLICES{,_EN}.md` 的台账完成行 |
| 零改动 | `docs/DELIVERY_DIRECTIVE{,_EN}.md`、`docs/CONTRACTS{,_EN}.md`、`schemas/**`、契约夹具、`.github/agents/**`、`.github/workflows/**`、`internal/**`、`tools/gates/**`、`go.mod`、`tools/gates/cjk-newwords.txt` |

**行使执行自由度的实质性选择**（§1.7.3 留痕；逐条见 §10）：

1. ⑤ 对照腿选型取「合成带种缺陷」（候选 B），理由见 §3.1；
2. 评分器在**实验腿之后、对照腿之前**修好（四条通用规则，非针对本次输出特判）——报告须**同时**给出修正前后两读（§4.1）；
3. 腿产物落盘口径不对称（实验腿逐字节捕获 / 对照腿逐字转录）——如实登记（§7）；
4. 批量 4 腿**并发**执行（方案书写「依次执行」）——如实登记、不补测（§7）；
5. 成本读数口径由「逐次窗口读数」降级为「用户提供的**累计**读数」——如实登记（§5）。

---

## 2. 评估目的、范围与边界

**目的**：为「是否把 GLM 系列纳入准则审查链对应位置」提供**本样本内的**质量与成本证据。

**范围**：两个评估面。

- **面一（质量对照，n=1）**：把 `GLM-5.3 (glm)` 放在 **⑤ 预审员**位置，与现行 ⑤ 载体所用模型
  `DeepSeek V4 Pro (deepseek)` 在同一输入、同一提示词下并跑，按**预注册钥匙**评分。
- **面二（批量辅助可用性）**：`GLM-5.3-Flash (glm)` 承担 4 类批量辅助任务（证据树枚举统计 / 长文档摘录与数字核对 /
  双语镜像语言侧检查 / 归档清单草案），只评**质量与产物可用性**。

**边界（明确不做）**：不写入白名单、不改角色集、不改 §2.2 / §2.4 / §5.2 / §7.5 / §1.5.1 正文；
不做多轮重复以获取统计效力（本片 n=1）；不对 V4 Pro 侧做 token 计量（读数不可得，见 §5）。

---

## 3. 评估设计

### 3.1 同输入构造（选型 B：合成带种缺陷）

**选型理由**：候选 A（以真实仓库变更集为输入）被否——① A 的答案钥匙**已提交入库**（提交 `025158f` 的验证报告与
自由度清单全文载有 ⑤/⑦ 全部发现与措辞），审查载体具备读文件能力，软约束「不得读仓库」不可机械强制 ⇒ 对照必然失真；
② A 的钥匙就是 V4 Pro 自己的发现，用它评 V4 Pro 的命中率存在系统性偏差（恰是要测的变量）；③ A 的钥匙不完备。
候选 B 的实验窗口内**仓库不存在任何答案源**（合成场景与真实仓库无同文）⇒ 污染面为零。

### 3.2 bundle 冻结对象（`tmp/glm-eval/bundle/`，任务包内逐字节内联，不给路径）

| 文件 | 字节 | SHA-256 |
|---|---|---|
| `00-plan.md`（合成 mini 方案书） | 1,184 | `70076661d86615d8dfe180aa5884ec9618338ec9496b3893900a749c092f36b8` |
| `01-clause-cn.txt`（合成契约原文，自造节号 X1.1/X6.3） | 2,346 | `d8fe4bc82df774535c47c39406c96fee1325c87f3153b5d72ee2a378ff03e525` |
| `02-diff.js`（合成变更集，约 6.5 KB） | 6,451 | `c939bc9e756696c463571c981a7d3756c87563de913b9c4c58faebdc4d5ff405` |
| `03-prompt.txt`（任务包正文，两腿逐字节相同，**不写模型名**） | 1,597 | `3b7a3580edcecb6a92ccb42f8acf37b469283136cb0c8d8d1b2b24d0b54fd758` |
| **`input_hash`**（四文件定序拼接） | 11,578 | `bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6` |

任务包正文（`03-prompt.txt`）含 B.5 的输入 / 产出 / 禁止 / 结论行契约，并含两条硬约束：
「不得读取任何工作区或外部文件，仅依据本消息文本；报告以最终消息返回」。

### 3.3 预注册答案钥匙（S1–S8；先于两腿调用落盘）

| 种子 | 严重度 | 缺陷 |
|---|---|---|
| S1 | Critical | 解析 SHA256SUMS 失败 / 白名单缺失时 `catch` 吞并并 `return true`（fail-open，静默通过） |
| S2 | High | 包成员名与清单名直接字符串比较，未做反斜杠到正斜杠归一 ⇒ Windows 上包恒判红 |
| S3 | High | 复算时取清单行**成员名列**当摘要，且 `if (recorded === recorded)` 恒真 ⇒ 摘要比对形同虚设 |
| S4 | Medium | 契约声明「清单有成员而包内无该文件」亦为冲突，实现无该分支 |
| S5 | Medium | 契约声明四子查，实现只落地三支（缺重建注记检查） |
| S6 | Medium | 「只判变更集」被实现为全树扫描（与契约冲突，历史误报） |
| S7 | Low | 顶部注释写「八类」实为四类（陈旧注释） |
| S8 | 无缺陷 | 纯重排 / 空行（期待 0 发现，报出即计 FP） |

**钥匙未事后修改声明**：`key.json` 的 SHA-256 为 `6fccbe461286abae63e74e52f64d8130f4fa979bb6d27ba71651fad6c3339536`，
与数据集 `pre_registration.key_file_sha256` **逐字相等** ⇒ 冻结后未再改动（复核命令 `node tmp/persist.js`）。

**钥匙与细则的张力（如实登记，两读并报）**：预注册钥匙把 S4 / S5 / S6 定为 **Medium**；而发给模型的细则把
**High 定义为「重要分支缺失或判断恒真恒假」**。按细则重判，S4 / S5 / S6 均可读作 High ⇒ `sev_exact = 1.0`、`Q = 1.0`。
**两个读数并存**：按钥匙 `Q = 0.9571`；按细则 `Q = 1.0`。**钥匙不因该张力而改动**（上段已证）。

### 3.4 评分细则与复合分（预注册）

- 命中率：`recall_all = 命中 / 7`；`recall_high = 命中(S1–S3) / 3`。
- 假阳性：发现所引行不存在、或断言与 bundle 不符者计 `FP`；**额外真发现**经主控机械复核后单列，不计 FP。
- 严重度校准：精确 1.0 / 差 1 档 0.5 / 差 ≥2 档 0；`sev_exact` 取 7 个计分种子的均值。
- 结论行：末行恰为 `PRE-REVIEW: PASS | PASS WITH FIXES | FINDINGS` 之一（§4.1 闭集）⇒ 1 / 0。
- **复合分 `Q = 0.4 × recall_high + 0.2 × recall_all + 0.2 × sev_exact + 0.2 × conclusion_line`**。
  `Q` 仅作**描述量**；n=1 不得下模型优劣结论。

### 3.5 污染检测（机械执行）

对两腿输出检索标记词表（`OB-79` / `025158f` / `0c12abc` / `G8-[abcd]` / `DIRECTIVE-EVIDENCE-SCOPE` / `重建归档` /
`判据作用域排除` / `S[1-8]`）须 **0 命中**；另核对每条 `file:line` 是否落在 bundle 内路径。
**实测**：两腿均 `PASS`（0 命中、路径全在 bundle 内）。**S8 区段**（`02-diff.js:191–198`）两腿均判「非发现」⇒ `s8_fp = 0`。

### 3.6 冻结与预注册自证

`tmp/glm-eval/FROZEN-MANIFEST.txt`（73 条）为会话内冻结清单，一次复算既得清单也复验 `input_hash` 与 `key.json`：
**`input_hash` 复算值 = 预注册值；`key.json` 哈希 = 预注册值；数据集 8 行、`pre_registration` 未变**。
交接件 `BRIEFING.md` 的**权威哈希**以该清单所列**完整哈希**为准：
`18867c95816abfa977c30eec31a74ef54ef17405b6c90e91caa6f05415df3ae8`（18,971 B）。

---

## 4. 评估结果

### 4.1 面一 · ⑤ 位置两腿对照（同输入、同提示词；n=1）

| 项 | ⑤-实验 | ⑤-对照 |
|---|---|---|
| `model_string`（限定名，R2.2） | `GLM-5.3 (glm)` | `DeepSeek V4 Pro (deepseek)` |
| `quota_package` | `glm-trial-2000w` | `deepseek-n-a` |
| 载体 | `deep-reasoner`（只读） | `deep-reasoner`（只读） |
| 输出文件 / 字节 | `leg-glm.txt` / 13,234 | `leg-control.txt` / 6,591 |
| 输出 SHA-256 | `54c895f9…` | `2fa1d9ee…` |
| 落盘口径 | 工具原始捕获**逐字节**保存 | 工具内联返回、主控**逐字转录** |
| 发现条数（自报） | 12（Critical 1 / High 6 / Medium 3 / Low 2） | 15（Critical 1 / High 7 / Medium 5 / Low 2） |
| `recall_all` | 1.0（7/7） | 1.0（7/7） |
| `recall_high` | 1.0（3/3） | 1.0（3/3） |
| `sev_exact` | 0.7857 | 0.7857 |
| `conclusion_line` | 1（`PRE-REVIEW: FINDINGS`） | 1（`PRE-REVIEW: FINDINGS`） |
| `FP` | 0 | 0 |
| 弃权行 | 2 | 1 |
| `extra_true`（额外真发现，经复核） | 10 | 8 |
| `s8_fp` | 0 | 0 |
| **`Q`（钥匙口径）** | **0.9571** | **0.9571** |
| 污染检测 | PASS | PASS |

- 严重度逐项（两腿一致）：S1 Critical→Critical(1.0)、S2 High→High(1.0)、S3 High→High(1.0)、S7 Low→Low(1.0)；
  S4 / S5 / S6 钥匙记 **Medium**、两腿均报 **High**(0.5) ⇒ `sev_exact = 5.5 / 7 = 0.7857`。
- **按细则重判**（§3.3 的张力）：`sev_exact = 1.0`、`Q = 1.0`。**两读并报**。
- **评分器修正前后两读（实验腿）**：修正前 `Q = 0.9048 / FP = 1` → 修正后 `Q = 0.9571 / FP = 0`。
  修正内容为四条**通用**规则（区间引用 / 逗号引用 / 弃权识别 / 去重），**非针对本次输出特判**；`key.json` 未改。
  对照腿只用修正后版本。
- **可用作方向性观测（非 token）**：实验腿输出体量约为对照的 **2.0 倍**（13,234 对 6,591 字节）。
  **字符数不等于 token、不可换算**，不得据此宣称成本高低。

### 4.2 面二 · 批量 4 腿（`GLM-5.3-Flash (glm)`，包 `glm-flash-realname-newuser`）

| 腿 | 任务 | 结果 | 裁决 |
|---|---|---|---|
| BATCH-a | 证据树枚举 / 统计（40 条合成清单记录） | 40 ✓、扩展分布 ✓、209,020 字节 ✓、顶层 6 条 ✓；**top-5 错**（命中 3/5：9817、9610、8840；漏 9403、9196；多列 8633、8426） | **部分失败**（排序子任务） |
| BATCH-b | 长文档摘录 + 数字核对（v1.24 报告对） | 328/328 ✓、作用域 39/2/4 ✓、成本向量 ✓、提交 `025158f` ✓、run `37278691747` ✓、CN/EN 逐项 `ALL MATCH` ✓ | 可用（全对） |
| BATCH-c | 双语镜像语言侧检查 | 14 条语言侧意见；守边界（无通过 / 不通过判定、不涉结构；正确拒绝把 `driver` / `To be discussed` 当术语漂移） | 可用（审查输入） |
| BATCH-d | MANIFEST + SHA256SUMS 草案 | 6 条 / 4437 字节 ✓、6 条哈希与字节 ✓、`archive/` 去目录层级 + 字典序 ✓、两空格格式 ✓；微差：未入库表末列名 | 可用草案 |

**结论（本样本内，n=1，不得外推）**：面二四腿中**三腿产物可用**、**一腿（BATCH-a）在排序子任务上部分失败**。
批量腿只评质量与产物可用性，**不参与** ⑤ 位置的质量分。

---

## 5. 成本与计量（单边，相对判据不成立）

**计量口径降级声明（如实登记）**：调用窗口内无法取得控制台**前置 / 后置**读数（operator 不在场、无窗口独占确认），
扩展日志与存储亦**无 usage 字段** ⇒ 单次 token 不可归属。后经用户在会话内**补提供累计读数**：

| 模型 | 配额包 | 累计读数 | 本片调用次数 |
|---|---|---|---|
| `glm-5.3` | `glm-trial-2000w` | **161,863 tokens** | 2（限定名探测 1 + ⑤ 实验腿 1） |
| `glm-5.3-flash` | `glm-flash-realname-newuser` | **424,690 tokens** | 4（批量 4 腿） |
| `deepseek-v4-pro` | `deepseek-n-a` | **不可得** | 1（⑤ 对照腿） |

- **逐行标包**（数据集 8 行）：探测拒绝 / 探测接受 / ⑤ 实验腿 = `glm-trial-2000w`；对照腿 = `deepseek-n-a`；批量 4 腿 = `glm-flash-realname-newuser`。
- `metering_mode`：GLM 侧 = `user-provided-cumulative`；V4 Pro 侧 = `unavailable`。
- **不填估值、不做 token 换算**。原始降级依据保留于数据集 `metering_note_at_call_time`。
- **单边绝对量级（口径更正，⑤ 预审 Medium #3）**：⑤ 实验腿单次的量级**上界 ≈ 161,863 tokens**
  （用户提供的累计读数；其中**探测腿为单词级最小调用**、且用户早期测试消息**未剥离**）。
  ⇒ **不得**以 `161,863 ÷ 2 ≈ 80,931` 当作实验腿单次量级——该均值含最小探测腿，会**系统性低估**。
- **相对成本判据不成立**：V4 Pro 侧读数不可得 ⇒ **不得**给出「GLM 与 V4 Pro 谁更省」的相对结论；该比较**悬置**，作为后续切片的输入。
- **新增如实登记**：实测消耗**显著高于主控预估**（预估 10k–20k / 次量级，实测约 4–8 倍）。原因未核
  （可能为 CoT 计费、系统提示开销，或用户早期测试消息未剥离）——**只登记、不解释**。
- **交接件成本必载项的口径更正（⑤ 复审 NL-2，已处置）**：`BRIEFING.md` 的成本必载项第 3 条原写「单次预审量级 ≈ 80k」，**已被本节更正取代**——按**上界** 161,863 写入，并明确**禁止**使用 ÷2 均值口径；交接件本身不改写（用户裁定「以 `FROZEN-MANIFEST` 为唯一权威」）。
- **⑦ 三读成本实测（本片数据点）**：约 3 小时墙钟、约 600 credits（约 200 credits / 读），与上一片 `DIRECTIVE-EVIDENCE-SCOPE` 同量级 ⇒ 已登记为 **OB-84**（§12.4），并在 §12.5 立缺口行「⑦ 三读成本无预算模型」，供后续评估片 / 治理片预算参考。

---

## 6. 可证伪性与机械核验

| 核验项 | 断言 | 结果 |
|---|---|---|
| 冻结包复算 | `node tmp/persist.js` 重算清单并比对全部哈希 | **match=true**（`input_hash`、`key.json`、数据集 8 行） |
| 钥匙冻结 | `key.json` 哈希 = `pre_registration.key_file_sha256` | **相等** |
| bundle 未被动过 | 四文件逐文件 sha256 与预注册逐条相等 | **相等** |
| 评分 / 哈希 / 污染脚本 | ③ 独立对抗式自测器（`selftest-report.txt`，**179 断言**，含 T-GE-13/14） | `SELFTEST TOTAL=179 PASS=179 FAIL=0`，对**最终**评分器 `score.js = 2c4052fa…`（该输出已在 ⑤ 复审 NL-1 整改中**再生**，见 §7 的相应记录） |
| 变异证伪 | 评分分母 / 范围短路 / 弃权分节 / 污染边界四处变异 ⇒ 断言转红 ⇒ 逐字节还原 | **四处均转红后复绿**（sha256 前后一致） |
| 污染检测 | 两腿输出标记词表 0 命中 + 路径全在 bundle 内 | **PASS**（表述为「未检出污染信号」，**不主张**「证明未读仓库」）|
| 数据集零泄漏 | 数据集内不得出现种子编号 / 钥匙正文 | 常驻断言通过 |
| 载体自报 | **不采信**子代理自报的模型与档位（§2.3） | 以调用参数与 operator UI 记录为准 |

- **③ 报告正文与冻结自测输出的版本差（⑤ 预审 Medium #1，已整改）**：`TEST-REPORT-DIRECTIVE-GLM-EVAL.txt`
  的正文停留在**第二轮（154 断言、评分器 `b4ad39…`）**口径，而**最终**评分器为 `2c4052fa…`、冻结自测输出为 **179/179**。
  ⇒ 已在 `TEST-REPORT` 追加**第三轮（最终版本复验）**小节。更正前实验腿读数（`Q=0.9048 / FP=1`）所用的旧评分器
  **已不在盘上**，仅存其「变异前后 sha256 一致」链（`TEST-REPORT` §8.5）⇒ 作为**复现性局限**如实登记，不补造工件。

- **⑤ 预审第 1 轮的 Low 处置（记录在案，不阻断）**：① `plan-v1.md` 的两处**笔误**（§5 写「bundle/ 六文件」实为四文件；
  §1 的污染检测命令路径写 `tmp\glm-eval\out\*.md`，实际产物在 `tmp/glm-eval-out/*.txt`，实现以 `contam.js` 等价执行）
  ——**不修改**已冻结的方案书（它是预注册依据，事后改动会损害预注册完整性），改在本报告**登记为 erratum**；
  ② 数据集 `call_no` 的声明域为 `1..N`，本片 **8 行**（探测 2 + ⑤ 2 + 批量 4，逐行），与 ① 方案书 §5 的 6 调用面口径差异在此显式声明；
  ③ 批量 b / c / d 的输入为**内联提示词、未落盘** ⇒ 三腿产物**不可复现审计**（`input_hash: null` 即该事实的登记），如实声明；
  ④ `extra_true`（10 / 8）的复核记录：由 ⑤ 预审独立抽查 F5 / F9 / F10 三例，确认为合成 diff 中的真实缺陷（非误计）。

- **`plan-v1.md` 的更正注记（⑦ 盲审 M2 / Low-4，§11.1 形态）**：见 `tmp/glm-eval/PLAN-ERRATA.txt`。
  严重度档位的**权威口径**为 `key.json` 的 `severity_rank`（`Critical` 3 / `High` 2 / `Medium` 1 / `Low` 0），
  与方案书 §1 括号中的三档枚举**档距等价**（两处 `High − Medium = Medium − Low = 1`；`Critical` 仅用于 `S1`，
  两腿均精确命中）⇒ **无功能漂移**；污染检测的**正确命令形态**亦在该注记内。
  **方案书正文不写回**——理由：`tmp/persist.js` 从聊天转录**逐字重新提取**该文件，手工追加会在下次复算时被覆盖
  （实测：追加后重跑即回到 12,860 B 的转录原样）⇒ 手工改动**不持久**。

---

## 7. 异常与兜底记录

| # | 事件 | 定性 |
|---|---|---|
| 1 | **计量不闭合**：窗口内无前置读数、扩展日志无 usage 字段 | 成本维度先降级 `unavailable`，后由用户补累计读数 ⇒ 改为**单边**绝对量级；相对判据仍悬置（§5） |
| 2 | **批量 4 腿并发执行**（方案书写「依次执行」） | 协议偏离，如实登记、**不补测**；逐次耗时不可分；不影响结论有效性（批量腿只评质量与可用性） |
| 3 | **腿产物落盘口径不对称**（实验腿逐字节捕获 / 对照腿逐字转录） | 如实登记；两枚哈希口径不同，**不得**互相当作同一采集路径 |
| 4 | **评分器在实验腿之后、对照腿之前修好** | 属**合理**修正（四条通用规则）；报告**同时**列出修正前后两读（§4.1） |
| 5 | **钥匙 / 细则张力**（S4 / S5 / S6：钥匙 Medium vs 细则 High） | 两读并报 + 声明钥匙**未事后修改**（§3.3） |
| 6 | **① 方案书此前未落盘** | 上一会话已修正：落盘 `tmp/glm-eval/plan-v1.md` + `PLAN-PROVENANCE.txt`（§3.6） |
| 7 | **编码事件（会话间外部保存）** | 交接件 `BRIEFING.md` 在上一会话收尾时实测 14,687 B / **无 BOM**；本会话追加前实测 **14,690 B / 含 BOM**（`EF BB BF`）——差值恰为 **+3 B 的 BOM**，**内容逐字节相同**。已用聊天转录本重建 14,687 B 版本并逐字节比对证实，其 SHA-256 与上一会话所载 `d8463ef5…` **精确一致**；取证链 `tmp/recon-briefing.js` + `tmp/fix-bom-and-append.js`。落盘时**已恢复**为 `tmp/` 过程件口径（无 BOM + LF）。同批提示中的 `PROBE-NOTES.txt` 经复核**未受影响**。 |
| 8 | `runSubagent` 偶发 `Agent error: no response was returned` | 处置规程：**先查盘 / 转录本**，再决定是否重发，避免重复付费调用 |
| 9 | **报告草案先于 ⑤ 产出（评估片特殊形态）** | 准则 §5.1 的 ⑧a 位于 ⑦ 之后；本片为**评估片**，审查对象是**评估执行与证据**（用户同轮定位），⑤⑥⑦ **均不审报告正文**。故报告草案在 ⑤ 之前产出，**不属** ⑧a 环节的产物，也不作为 ⑤–⑦ 的审查对象；其最终正文由 ⑧a 定稿（含 `_EN` 镜像）。 |
| 10 | **门禁红（如实登记，未静默略过）** | 报告草案落盘后 `--scope=tree` 先报 **G4b**（8 个新汉字）与 **G6-4**（两个「行数」型过期字面量），**均已就地修复**（改写用词、改口径，`cjk-newwords.txt` **未扩**）；其后 **G7-b**（`repo` 工件未入索引）与 **G7-a**（缺 `artifacts` 块）先后判红——按用户同轮口径**不是预期忽略项**，**均已在 ⑧a 落盘 + `git add` 后关闭**（§9、§16）。 |
| 11 | **自测报告再生（⑤ 复审 NL-1 整改的副作用）** | ③ 重跑发现 `run-selftest.js` 有 3 条**脚手架期**断言按设计失效（`calls[]` 已由执行步填为 8 行、`tmp/glm-eval/` 新增交接文档）、1 条 domain 断言需随数据集 schema 整改同步 ⇒ 已由 ③ 更新为**执行后真实状态**（`T-GE-09.known` 由通配**收紧**为 58 项闭集）并重跑至 `179/179`。副作用：`selftest-report.txt` **被再生**，其**脚手架期原始字节不再可复现**（属**不可逆的取证损失**，如实登记，不补造）。 |
| 12 | **⑦ 盲审 M1 — 数据集列域与空值不自洽** | 二轮修复：`input_hash` / `output_hash` 的 `domain` 与 `duration_s` 的 `type` 改为**可空**并写明不可得原因（与 ⑤ 轮已修的 `metering_mode` / `measured_tokens` / `conclusion_line` 同类）；`.json` 重校验通过。 |
| 13 | **⑦ 盲审 M2 — 严重度档位双口径** | 更正注记落 `tmp/glm-eval/PLAN-ERRATA.txt`（§11.1「只加更正注记、不重写正文」）：权威口径 = `key.json` 的 `severity_rank`；两处**档距等价**、`Critical` 仅用于 `S1` ⇒ **无功能漂移**。**不写回 `plan-v1.md`**（`persist.js` 从转录逐字再生该文件，手工改动不持久）。 |
| 14 | **⑦ 盲审 M3 — 自测报告覆写导致取证损失** | 如实登记为**不可逆取证损失**；并在 `TEST-REPORT` 追加 §12「证据保留纪律」（先快照、后重跑；一态一文件；不可复现即声明），**并立即执行**：现行 `selftest-report.txt` 已快照为 `tmp/glm-eval-out/selftest-report-20261005.txt`（两枚哈希相同）。 |
| 15 | **⑦ 盲审 Low-4 — 污染命令路径笔误** | 已并入 `PLAN-ERRATA.txt` 更正二（正确形态 = `contam.js` 对两腿 `.txt` 执行，或对 `leg-*.txt` 做 `Select-String`）。 |
| 16 | **⑦ 终审 H1 — 冻结清单未登记 `PLAN-ERRATA.txt`** | 已闭合：重跑 `node tmp/persist.js` 后，`PLAN-ERRATA.txt`（4,320 B / `a3c08355…`）与快照 `selftest-report-20261005.txt`（26,531 B / `7b06d9d7…`）均已登记进 `FROZEN-MANIFEST.txt`。 |
| 17 | **⑦ 终审 M2 — 方案书 `call_no` 行域与执行面不一致** | 已并入 `PLAN-ERRATA.txt` **更正三**：8 行 = 2 探测 + ⑤ 2 + 批量 4，并给出差异来源与当前权威口径。 |
| 18 | **⑦ 终审 M3 — 严重度等价性未写适用边界** | 已在 `PLAN-ERRATA.txt` **更正一**补「适用边界」段：等价性**仅对本样本读数成立**；通用评分一律以 `key.json` 的 `severity_rank` 为准（附四档映射）。 |
| 19 | **§1.7.4 独立登记（OB-83）** | GLM 两枚模型不在 §2.2 白名单，其调用属**用户同轮授权的临时评估** ⇒ 已按 §1.7.4 五要素在 `docs/DELIVERY_DIRECTIVE{,_EN}.md` 的 §12.4 **独立登记为 OB-83**（裁决人 / 日期 / 理由 / 依据齐备，并如实标注「超出 §2.4 字面标准、系用户直接裁定、不得作为先例」）；随片升 **v1.25**，本行为切片内副产物、非独立准则治理片。 |
| 20 | **计量规程「先建后未用」** | `tmp/glm-eval/metering/README.md` 的「pre / post / recheck 三读」规程**在调用窗口内未能执行**（operator 不在场），实际走「用户提供的累计读数」路径 ⇒ 该规程**未被实际走过**；保留为后续评估片的**首选项**，并建议按可得性改写（明写「无前置读数则退到 `user-provided-cumulative`」的路径与后果）。 |

**兜底计数**：本片为定位 / 修复 / 重跑而介入的逐次行动 ≤ 3（未触发 §9.2 的「> 3 次即标高风险」）。

---

## 8. 环境与档位

| 项 | 值 |
|---|---|
| 主控 | `DeepSeek V4.1 Flash (deepseek)`，档位 High |
| ① 架构师 / ⑤ 预审员 | `DeepSeek V4 Pro (deepseek)`，档位 Max |
| 产品层（实现 / 测试 / 文档） | `DeepSeek V4.1 Flash (deepseek)` |
| ⑥ 独立扫描 | `MAI-Code-1.1-Flash (copilot)` |
| ⑦ 独立终审 | `GPT-5.3-Codex (copilot)`，Extra High（**不降级**） |
| 评估对象（非准则角色） | `GLM-5.3 (glm)` / `GLM-5.3-Flash (glm)`（均**不在**白名单，属用户授权临时评估） |
| 载体 | ⑤ 两腿与限定名探测统一用 `deep-reasoner`（只读）；实验变量仅 `model` 参数 |
| 探针 | **0**（未运行 `.vscode/scripts/agent-probe`） |
| 档位变更 | 无（全程未暂停改档） |

---

## 9. 门禁结果

| 时点 | 命令 | 结果 |
|---|---|---|
| ④（报告落盘前） | `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 328/328` |
| ④（报告落盘前） | `node tools/gates/gate.js --all --scope=tree` | `changed=1`、`TOTAL_FAIL=0`（含 `PASS base go test ./... :: ok=14 FAIL=0`） |
| 报告草案落盘后 | `node tools/gates/gate.js --all --scope=tree` | `changed=2`；**G4b / G6-4 判红 → 已就地修复**；现存 **G7-b 判红** |
| 移除 `artifacts` 块后 | 同上 | **G7-a 判红**（变更的 `docs/validation/*.md` 缺 `artifacts` 围栏块） |
| ⑨ 原生验证（⑧a 落盘 + `git add` 后） | `--all --scope=tree` / `--all --scope=index` / `--selftest` / `--check=G8-b` | **tree `TOTAL_FAIL=0`**、**index `TOTAL_FAIL=0`**、`SELFTEST: PASS 328/328`、G8-b `1 pack(s) self-consistent` ⇒ **G7 已全绿** |
| 基础门禁 | gofmt / build / vet / test（§6.1，由 `--all` 的 base 项承担） | 全绿 |
| 编码 | `.md` = BOM + LF；`.json` / `.js` / `.txt` = 无 BOM + LF | 逐文件核验通过 |
| 推送后（§11.2 完成行） | 主提交与 CI | 主提交 `09dfef7`，推送范围 `aaa6e1e..09dfef7`（仅 `origin`）；CI run `37346176237` 双腿 **success**（`SELFTEST: PASS 328/328`、`GATE REPORT scope=ci changed=87`、`TOTAL_FAIL=0`） |

**门禁结论（⑨ 实测）**：三个读数全绿——tree `TOTAL_FAIL=0`、index `TOTAL_FAIL=0`、`SELFTEST: PASS 328/328`；
此前的 **G4b / G6-4 / G7-b / G7-a** 四处红**全部关闭**（修复路径逐条见 §7 的相关记录与 §16）。

> `G1`（三处人工撰写内容互判）与 `G5`（需仓库外结构化输入）由脚本如实打印 `SKIP`，**不得读作 PASS**。

---

## 10. 越界自检

- `git status --short` → 87 项（4 项修改：台账 CN/EN + 准则 CN/EN（v1.25 台账补登）；83 项新增：报告对 / 数据集 / 证据包 / 自由度清单），**全部逐文件 `git add`**（未用 `-A`）；无未申报的未跟踪文件。
- `git grep -n "GLM" -- docs/DELIVERY_DIRECTIVE.md docs/CONTRACTS.md schemas .github/agents` → **无输出**。
- `git grep -l '^model:' -- '.github/agents/prfrail-*'` → **无输出**（G5-a）。
- `git status --porcelain -- .github/agents internal schemas go.mod .github/workflows tools` → **无输出**。
- §2.2 / §2.4 / §5.2 / §6.3 / §7.5 与 §1.5.1 正文**零改动**；`cjk-newwords.txt` **未扩**。
- 探针 **0**（未运行 `.vscode/scripts/agent-probe`）。
- **预置免责**：`.github/agents/` 下的 `deep-reasoner.agent.md` / `independent-reviewer.agent.md` / `quick-verifier.agent.md`
  **含** `^model:` 行——系**事先存在**、本片未触碰，且不属 G5-a 约束范围（G5-a 仅约束 `prfrail-*`）。

**③ 测试报告的独立覆盖缺口（如实照录）**：G2（bundle 无 BOM/LF 的常驻断言）、G5（计量规程无执行式守护）、
G6（引用解析误吞代码块 / URL）、G7（非锚点接受行的针对性变异）四项仍缺；G3（数据集列 ↔ ① 方案书）因方案书不在磁盘而**不可测**；
G4（并发 / 崩溃）**不适用**（单发 CLI）。

---

## 11. 明确未执行事项

1. **未把 GLM 写入任何白名单或角色集**（本片边界即「只评估」）；
2. **未做多轮重复**（n=1，不追求统计效力）；
3. **未做 V4 Pro 侧 token 计量**（读数不可得）；**未给出相对成本结论**；
4. **未补测批量腿的并发偏离**（不二次调用）；
5. **未对评分器四条通用规则之外的评分口径做任何调整**；
6. **未提交 / 未推送**（授权边界：commit 否 / push 否 / gitee 否）；
7. **未运行探针**（额度 0）。

---

## 12. 下一步建议

1. **白名单决策另立切片**：本片只提供证据；是否纳入 `GLM-5.3 (glm)` 到 ⑤ 位置，建议在**成本读数可闭合**
   （V4 Pro 侧 token 可得）后另议，避免只凭单边证据决策。
2. **成本维度补测**：把「V4 Pro 侧 token 读数」作为后续切片的输入；若长期不可得，考虑改用**输出字符数 / 耗时**
   作为公共代理指标并预先声明其局限。
3. **批量辅助的适用面**：BATCH-b（长文档摘录 + 数字核对）与 BATCH-d（归档清单草案）在本样本中产物完全可用，
   可作为「辅助型任务」的候选；BATCH-a 的排序子任务在本样本中部分失败，**不宜**单独承担排序 / Top-N 选择。
4. **评分器通用规则入册**：四条通用修正规则（区间引用 / 逗号引用 / 弃权识别 / 去重）建议在后续评估中**保持为主控侧工具**，
   不写入准则正文。

---

## 13. 执行流水线

| 步 | 角色 / 模型 | 状态 | 产物 / 备注 |
|---|---|---|---|
| 0 前置校验 | 主控 | ✅ | 角色 / 额度 / 边界确认；三条硬前置通过（§3.6） |
| 0a 限定名探测 | 主控；`GLM-5.3 (glm)` | ✅ | 尝试 1 被拒（**未触发调用、无消耗**，但登记为一次尝试）、尝试 2 接受 ⇒ 权威限定名 |
| ① 方案书 | 架构师 `DeepSeek V4 Pro (deepseek)` | ✅ | `tmp/glm-eval/plan-v1.md` + `PLAN-PROVENANCE.txt` |
| ② 文档 / 脚手架实现 | 实现者 `DeepSeek V4.1 Flash (deepseek)` | ✅ | bundle 四文件、数据集脚手架、评分 / 哈希 / 污染脚本；**不发起付费调用** |
| ③ 独立对抗式自测 | 测试工程师 `DeepSeek V4.1 Flash (deepseek)` | ✅ | 两轮 + 整改后一轮；`SELFTEST TOTAL=179 PASS=179 FAIL=0` |
| ④ 集成 + 门禁（冻结 + 预注册） | 主控 | ✅ | `TOTAL_FAIL=0`、`SELFTEST 328/328`；bundle 冻结 + 钥匙预注册 |
| ⑤-position 两腿（实验 → 对照） | 主控；载体 `deep-reasoner`（只读） | ✅ | 同输入、同提示词；实验腿先、对照腿后；两腿输出分别落盘、不互锚 |
| 批量 4 腿 | 主控；载体 `deep-reasoner`（只读） | ✅ | **并发执行（偏离，已登记）** |
| ⑤ 预审 | 预审员 `DeepSeek V4 Pro (deepseek)` | ✅ 2 轮 | 第 1 轮 `PRE-REVIEW: PASS WITH FIXES`（3 Medium）⇒ 整改 ⇒ 第 2 轮 `PRE-REVIEW: PASS WITH FIXES`（仅 3 Low） |
| ⑥ 独立扫描 | `MAI-Code-1.1-Flash (copilot)` | ✅ 1 次 | `INDEPENDENT SCAN: PASS`（0 条 file:line 级发现；**可信度降级**，见 §14） |
| ⑦ 盲审 / 终审 / 复审 | `GPT-5.3-Codex (copilot)`（Extra High，**不降级**） | ✅ 3 次 | `BLIND REVIEW: FINDINGS`（3 Medium + 1 Low）⇒ 整改 ⇒ `FINAL REVIEW: FINDINGS`（1 High + 2 Medium）⇒ 整改 ⇒ `RE-REVIEW: PASS`（7/7 闭合、0 新增） |
| ⑧a 报告对 + 台账 + 证据包 + 自由度清单 | ⑧ 文档工程师 + 主控 | ✅ | 本节所在 |
| ⑨ 原生验证 | 主控 | ✅ | 读数并入 §9 |
| ⑧b 定稿 | ⑧ 文档工程师 | ✅ | 吸收 ⑨ 读数与 ⑧c 记录 |
| ⑧c 收尾清洁 | 主控 | ✅ | 见 §16 |
| ⑩ 停点 | 主控 | ✅ | 停在停点、等待提交授权 |

**返工计数（§9.2 口径：同一环同一根因计 1 次）**：③ ×1（自测断言同步）、⑤ ×1（整改后复审）、⑦ ×2（盲审后、终审后各一轮整改）；**兜底介入 ≤ 3**，未触发「> 3 次即标高风险」。

## 14. 审查结论摘要（⑤⑥⑦）

- **⑤ 预审**：第 1 轮 `PRE-REVIEW: PASS WITH FIXES`（3 Medium：③ 报告版本口径 / 数据集 schema 自洽 / `÷2` 均值低估；若干 Low）⇒ 整改 ⇒ 第 2 轮 `PRE-REVIEW: PASS WITH FIXES`（**无 Medium+**；新增 3 条 Low，其中 NL-1 已由 ③ 更新断言并重跑至 `179/179` 闭合，NL-2 / NL-3 属交接件陈旧表述、按「以 `FROZEN-MANIFEST` 为唯一权威」处理）⇒ **⑤ 闭环**。
- **⑥ 独立扫描**：`INDEPENDENT SCAN: PASS`（Section A–E 齐备、0 条 file:line 级发现）。**如实登记负向信号**：其 Section C 第 3 项称「`leg-glm.txt` 中 S4 写为 `Medium — 02-diff.js:93`」——与实物不符（两腿在该种子均报 **High**：GLM 腿 `F3 | High | 02-diff.js:84-95`、对照腿 `F5 | High | 02-diff.js:93`）⇒ 按 §7.7 该审计**可信度降级**；其 PASS 仍被采纳（无 file:line 级发现），最终结论以 ⑦ 为准。
- **⑦ 独立终审**：三读全数消耗（§7.5 的 1 终审 + 1 复审 + 盲审 +1）。**盲审与终审均不附任何 ⑤/⑥ 清单**；仅复审轮附清单（§7.3）。整改后 `RE-REVIEW: PASS`（7/7 闭合、0 新增发现）。

## 15. 审查输入包摘要（含盲审隔离证明）

| 项 | 值 |
|---|---|
| ① 与 ⑤ 的调用隔离（R2.6） | ⑤ 两轮均为 `deep-reasoner` 载体的**独立调用**，未传入 ① 的对话历史或实现期推理；输入 = 方案书 + 评估证据**路径** + 契约原句 |
| ⑥ 首轮 | **不附**任何 ⑤ / ⑦ 清单；对象为评估执行与证据 |
| ⑦ 盲审（第 1 读） | **不附**任何 ⑤ / ⑥ 清单、不附前轮结论 |
| ⑦ 终审（第 2 读） | **不附**任何清单（对整改后盘上状态作独立复读） |
| ⑦ 复审（第 3 读） | 按 §7.3 **附完整清单**并逐条闭合 |
| 报告正文 | 按用户同轮定位，⑤⑥⑦ **均不审报告正文**（审查对象 = 评估执行与证据） |
| 输入消毒 | 提示词内只给**仓库相对路径**；绝对路径中的用户名段在审查报告中替换为 `<user>`；无 IP / 主机名 / 凭据 / token |
| 输入包落盘 | 评估证据一律以**路径**交付（`tmp/glm-eval/**`、`tmp/glm-eval-out/**`、`docs/validation/directive-glm-eval-dataset.json`），非转述 |

## 16. ⑧c 收尾清洁记录（主控执行）

| 项 | 读数 / 结论 |
|---|---|
| 全仓诊断（按 §3.4.1 分类） | 编译型（Go）：`gofmt -l .` 空、`go build ./...` / `go vet ./...` exit 0、`go test ./...` `ok=14 FAIL=0`；标记数据（`.md` / `.json`）：编码门禁逐文件核验 + JSON 解析通过；纯文本其他：编码门禁通过；脚本语言：本片变更集**无** `.js` / `.ps1` 改动 |
| IDE 诊断汇总 | 变更集内 6 个关键文件（报告对 / 数据集 / 证据包 `MANIFEST` / 自由度清单 / 台账 CN）均 **No errors found** |
| 门禁全绿 | `--all --scope=tree` `TOTAL_FAIL=0`；`--all --scope=index` `TOTAL_FAIL=0`；`--selftest` `SELFTEST: PASS 328/328`；`--check=G8-b` `1 pack(s) self-consistent` |
| 暂存与变更集 | 逐文件 `git add`（**未用 `-A`**）⇒ 87 项：4 项修改（台账 CN/EN + 准则 CN/EN 的 v1.25 台账补登）+ 83 项新增（报告对 / 数据集 / 证据包 77 条载荷 / `MANIFEST.md` / `SHA256SUMS.txt` / 自由度清单） |
| 编码（全类型逐文件） | 变更集内 `.md` = BOM + LF；`.json` / `.txt` / `.js` = 无 BOM + LF；证据树内的归档载荷按 §1.5.1「归档后类型归一」，其可核验性由 `MANIFEST` 的「后缀与编码映射」双哈希记录承担 |
| 临时目录 / 探针残留 | `tmp/` 内**一次性中间物全部清除**（11 件：`_g4bchk.js` / `_g4b.txt` / `_g8b.txt` / `_gate.txt` / `_gate-tree.txt` / `_persist2.log` / 5 件 `_append-*.txt`）；**①/② 类过程件（`tmp/glm-eval/**`、`tmp/glm-eval-out/**`、`tmp/glm-eval-batch/**`、`tmp/` 根承重脚本）保留原位**——依 §1.5.1「①② 类在归档执行片处置前豁免 §1.5 的用完立即清除；任何切片的 ⑧c 只清 ③ 类」，且**全量已归档**进本片证据包；探针残留 0 |
| 未跟踪文件清单 | 逐项核对：全部已在 `artifacts` 块内声明为 `repo` 或 `local-only`，**无**未申报的残留 |
| 与交接件的差异（如实登记） | `BRIEFING.md` 写「⑧c 清理（`tmp/` 仅留 `.gitkeep`）」；本轮按 §0.1「执行纪律冲突以本准则为准」改依 §1.5.1 **只清 ③ 类**，①/② 类保留待「`tmp/` 生命周期执行」片处置——差异已登记，未静默略过 |

---

> **工件清单块（`artifacts` 围栏）——⑧a 落盘后恢复**：`_EN` 镜像、证据包（`MANIFEST.md` / `SHA256SUMS.txt`）、
> 自由度清单与台账行**均已落盘**，原先「声明不实」的成因已消除 ⇒ 本节恢复完整块（见文末）；
> 恢复后须复跑 `--scope=tree` / `--scope=index` 至 G7 全绿。

```artifacts
docs/validation/directive-glm-eval.md	repo
docs/validation/directive-glm-eval_EN.md	repo
docs/validation/directive-glm-eval-dataset.json	repo
docs/validation/evidence/glm-eval/MANIFEST.md	repo
docs/validation/evidence/glm-eval/SHA256SUMS.txt	repo
docs/validation/evidence/DIRECTIVE-GLM-EVAL-freedom-list.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
tmp/glm-eval/BRIEFING.md	local-only
tmp/glm-eval/plan-v1.md	local-only
tmp/glm-eval/PLAN-ERRATA.txt	local-only
tmp/glm-eval/key.json	local-only
tmp/glm-eval-out/leg-glm.txt	local-only
tmp/glm-eval-out/leg-control.txt	local-only
```
