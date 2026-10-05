# 切片 `DIRECTIVE-GLM-EVAL` · 会话间简报（handoff）

> **用途**：上一会话封板 → 新会话起跑审查链的**唯一交接件**。
> **生成**：2026-10-05（会话最后动作）｜**切片目标**：评估 `GLM-5.3 (glm)` 在 **⑤ 预审位置**的质量与成本（对照 `DeepSeek V4 Pro (deepseek)`），并评估 `GLM-5.3-Flash (glm)` 在 4 类批量辅助任务上的质量与产物可用性 —— **只评估、不入白名单**；白名单/角色决策在本片之后另议。
> **本会话授权边界**：commit 否 / push 否 / gitee 否 / 探针 0（新会话授权由用户另行声明，本文件不预设）。
> **编码口径**：本文件与 `plan-v1.md` 是 `tmp/` 会话内临时物，按逐字节口径保存（UTF-8 **无 BOM** + LF）；⑧a 入库报告仍遵 `.md` = UTF-8 **with BOM** + LF。

## 0. 状态一览

| 步骤 | 状态 | 产物 / 证据 |
|---|---|---|
| 0 前置校验 | ✅ | 角色/额度/边界确认（见 §5） |
| 0a 限定名探测 | ✅ | `PROBE-NOTES.txt`；权威名 = `GLM-5.3 (glm)` |
| 1 ① 方案书 | ✅ 落盘 | `tmp/glm-eval/plan-v1.md`（此前仅存在于会话内） |
| 2 ⑤ 两腿（实验 + 对照） | ✅ 已跑已评分 | `tmp/glm-eval-out/leg-{glm,control}.txt` |
| 3 批量 4 腿（Flash） | ✅ 已跑已裁决 | 见 §4 |
| 4 门禁 | ✅ 全绿 | `TOTAL_FAIL=0`、`SELFTEST: PASS 328/328` |
| 5 ⑤ 报告预审 | ❌ **未开始** | 新会话第一件事 |
| 6 ⑥ / 7 ⑦×3 / 8a ⑧a / 9 ⑨ / 8b ⑧b / 8c ⑧c / 10 ⑩ | ❌ 未开始 | 见 §9 |

**结论量级（n=1，不得外推）**：两腿复合分 `Q=0.9571` 并列；**成本维度未闭合**（见 §6）。

## 1. 新会话起跑清单（按序）

1. 重读 `docs/DELIVERY_DIRECTIVE.md` §5.1 / §7.3 / 附录 B.7.3（⑤⑥⑦ 审查链定义与输入包规则）；
2. 读本文件 → `tmp/glm-eval/plan-v1.md`（⑤ 审查对象）→ `PLAN-PROVENANCE.txt` → `FROZEN-MANIFEST.txt`；
3. 读 `docs/validation/directive-glm-eval-dataset.json`（8 行）；
4. 复核 §8 的三条硬前置；
5. 进 ⑤ 报告预审 → ⑥ → ⑦ 三读（1 盲审 + 1 终审 + 1 复审）→ ⑧a → ⑨ → ⑧b → ⑧c → ⑩。

- 注意 A：`tmp/glm-eval/` **已回到工作区内**（全部调用已结束；钥匙在库内不再影响已完成的调用）。
- 注意 B：⑦ 预留 3 次**未动用**；审查链模型 = ⑥ `MAI-Code-1.1-Flash (copilot)`、⑦ `GPT-5.3-Codex (copilot)`。

## 2. 冻结物与哈希自证（SHA-256 of raw bytes）

| 对象 | 字节 | SHA-256 |
|---|---|---|
| `tmp/glm-eval/plan-v1.md` | 12,860 | `8be638d27cbf467f0b337de7e4362bad6b090b43e876377d1ad7ec7429181d04` |
| `tmp/glm-eval/PLAN-PROVENANCE.txt` | 1,500 | `66f1a3f5e81fe4619c16cf340f4add5e41a3eff2a5af0d7a09fa593874e18dfa` |
| `tmp/glm-eval/bundle/00-plan.md` | 1,184 | `70076661d86615d8dfe180aa5884ec9618338ec9496b3893900a749c092f36b8` |
| `tmp/glm-eval/bundle/01-clause-cn.txt` | 2,346 | `d8fe4bc82df774535c47c39406c96fee1325c87f3153b5d72ee2a378ff03e525` |
| `tmp/glm-eval/bundle/02-diff.js` | 6,451 | `c939bc9e756696c463571c981a7d3756c87563de913b9c4c58faebdc4d5ff405` |
| `tmp/glm-eval/bundle/03-prompt.txt` | 1,597 | `3b7a3580edcecb6a92ccb42f8acf37b469283136cb0c8d8d1b2b24d0b54fd758` |
| **`input_hash`**（四文件定序拼接） | 11,578 | `bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6` |
| `tmp/PROMPT-EXACT.txt`（= 上述拼接字节） | 11,578 | 同上（同一对象） |
| `tmp/glm-eval/key.json`（预注册钥匙） | 3,843 | `6fccbe461286abae63e74e52f64d8130f4fa979bb6d27ba71651fad6c3339536` |
| `tmp/glm-eval-out/leg-glm.txt` | 13,234 | `54c895f9df821231d9049ac7819077b44335db83383b266d86f3c28fee137b5b` |
| `tmp/glm-eval-out/leg-control.txt` | 6,591 | `2fa1d9eeedae134c758174fff900311983a52b0eda8470889f4ff5d67f463a7f` |
| `tmp/glm-eval/hash.js` | 1,378 | `ff608449685e5150e5338a8bcad5e36e0e8910f172ed77182b8e1bf3c4c073b3` |
| `tmp/glm-eval/score.js` | 12,339 | `2c4052fa8e10273b09465b668a09b0c1e4402065576ec48ca4357341d27ee051` |
| `tmp/glm-eval/contam.js` | 2,550 | `dfcb76e737b6ba5dab0453977fa6c1307bb2391e7890e5d1bd9aca3525099455` |
| `tmp/glm-eval/selftest/selftest-report.txt` | 25,932 | `d8f59ef338d00fd6b98da70e2356d3355c4d9bc12f98ecd5a5c8fcd356f3d405` |
| `tmp/glm-eval/TEST-REPORT-DIRECTIVE-GLM-EVAL.txt` | 25,019 | `12002624ceb2491dcde6db7b94658ef7afe34c9273cef41fddcec5f4326ddc3f` |
| `docs/validation/directive-glm-eval-dataset.json` | 9,917 | `3b61c53e75cb337148def19c0ef36f77bf3283c1a5164ebbd82b9077202b05d0` |

- 全量清单（冻结物 + 交接文档 + 承重脚本 + 批量钥匙）：`tmp/glm-eval/FROZEN-MANIFEST.txt`；一次复算既得清单也复验 `input_hash` / `key.json`：`node tmp/persist.js`。
- 交接文档自身哈希：`PROBE-NOTES.txt`、`BRIEFING.md`（本文件，避免自指故不内嵌）→ 见 `FROZEN-MANIFEST.txt` 末段。
- 已核验：`input_hash` 复算值 = 预注册值；`key.json` 哈希 = `pre_registration.key_file_sha256`（即**冻结后未再改动**）。

## 3. ⑤ 位置两腿读数（同输入、同提示词；n=1）

| 项 | ⑤-实验 | ⑤-对照 |
|---|---|---|
| `model_string` | `GLM-5.3 (glm)` | `DeepSeek V4 Pro (deepseek)` |
| `quota_package` | `glm-trial-2000w` | `deepseek-n-a` |
| 输出文件 / 字节 | `leg-glm.txt` / 13,234 | `leg-control.txt` / 6,591 |
| 落盘口径 | 工具原始捕获**逐字节**保存 | 工具内联返回、主控**逐字转录** |
| `recall_all` | 1.0（7/7） | 1.0（7/7） |
| `recall_high` | 1.0（3/3） | 1.0（3/3） |
| `sev_exact` | 0.7857 | 0.7857 |
| `conclusion_line` | 1（`PRE-REVIEW: FINDINGS`） | 1（`PRE-REVIEW: FINDINGS`） |
| `FP` | 0 | 0 |
| 弃权行 | 2 | 1 |
| `extra_true` | 10 | 8 |
| `s8_fp` | 0 | 0 |
| **`Q`** | **0.9571** | **0.9571** |
| 污染检测 | PASS | PASS |

- `Q = 0.4×recall_high + 0.2×recall_all + 0.2×sev_exact + 0.2×conclusion_line`（预注册五维）。
- 严重度逐项（两腿一致）：S1 Critical→Critical(1.0)、S2 High→High(1.0)、S3 High→High(1.0)、S7 Low→Low(1.0)；S4/S5/S6 钥匙记 **Medium**、两腿均报 **High**(0.5) ⇒ `sev_exact = 5.5/7 = 0.7857`。
- **钥匙/细则张力**：按预注册钥匙 `Q=0.9571`；按发给模型的细则（High =「重要分支缺失或判断恒真恒假」）重判 S4/S5/S6 ⇒ `sev_exact=1.0`、`Q=1.0`。**两读并报；钥匙未事后修改**（`key.json` 哈希 = 预注册值）。
- 可用作方向性观测（**非 token**）：实验腿输出体量 ≈ 对照 2.0×（13,234 vs 6,591 字节）；字符数 ≠ token、不可换算，不得据此宣称成本高低。

## 4. 批量 4 腿（`GLM-5.3-Flash (glm)`，包 `glm-flash-realname-newuser`）

| 腿 | 任务 | 结果 | 裁决 |
|---|---|---|---|
| BATCH-a | 证据树枚举/统计（40 行合成清单） | 40 ✓、扩展分布 ✓、209,020 字节 ✓、顶层 6 条 ✓；**top-5 错**（命中 3/5：9817、9610、8840；漏 9403、9196；多列 8633、8426） | **部分失败**（排序子任务） |
| BATCH-b | 长文档摘录 + 数字核对（v1.24 报告对） | 328/328 ✓、作用域 39/2/4 ✓、成本向量 ✓、commit `025158f` ✓、run `37278691747` ✓、CN/EN 逐项 `ALL MATCH` ✓ | 可用（全对） |
| BATCH-c | 双语镜像语言侧检查 | 14 条语言侧意见；守边界（无通过/不通过判定、不涉结构；正确拒绝把 `driver` / `To be discussed` 当术语漂移） | 可用（审查输入） |
| BATCH-d | MANIFEST + SHA256SUMS 草案 | 6 条 / 4437 字节 ✓、6 条哈希与字节 ✓、`archive/` 扁平化 + 字典序 ✓、两空格格式 ✓；微差：未入库表末列名 | 可用草案 |

## 5. 额度与预算消耗（两类包分开记账）

- 限定名探测：**2 次尝试**（尝试 1 被拒、未触发调用无消耗；尝试 2 接受）—— 独立计数，不占 ⑤ 上限。
- ⑤ 实验 `GLM-5.3 (glm)`：**1/1**；⑤ 对照 `DeepSeek V4 Pro (deepseek)`：**1/1**；批量 `GLM-5.3-Flash (glm)`：**4/4**。
- ⑦ 预留 3 次（1 盲 + 1 终 + 1 复审）：**未动用**。
- 探针：**0**（未运行 `.vscode/scripts/agent-probe`）。

## 6. 成本维度：**未闭合 → 悬置**（报告须照此声明）

- 本会话按裁决尝试闭合：以三问形式向 operator 索取读数 ——（A）控制台是否可读（B）200 万通用包累计消耗量级（C）1500 万 Flash 包累计消耗量级。**三次均返回自动应答（operator 不在场）⇒ 未取得任何读数**。
- ⇒ ⑧a 报告须写：「**成本判据未闭合，纳入 ⑤ 的成本问题悬置**」；**不得填估值**；`measured_tokens = null`、`metering_mode = unavailable`（⑤ 两腿）/ `interval-sum`（批量腿，逐次耗时与 token 均不可分）。
- 原始降级依据：调用窗口内无控制台前置读数；扩展日志/存储**无 usage 字段** ⇒ 单次 token 不可归属。

## 7. 六项偏离与定性

| # | 偏离 | 定性 |
|---|---|---|
| 1 | 计量不可得（无前置读数 + 无 usage 字段） | 见 §6：成本维度**未闭合 → 悬置**；不填估值 |
| 2 | 批量 4 腿**并发**执行（方案书写「依次执行」） | 入报告、**不补测**；逐次耗时不可分；不影响结论有效性（批量腿只评质量与产物可用性） |
| 3 | 腿产物落盘口径不对称（GLM 逐字节捕获 / V4 Pro 逐字转录） | 入报告；哈希口径不同**如实登记** |
| 4 | 评分器在实验腿之后、对照腿之前修好（区间引用 / 逗号引用 / 弃权识别 / 去重四条**通用**规则，非针对本次输出特判；`key.json` 未改） | **合理**；报告须**同时列出修正前后两读**：实验腿 `Q=0.9048 / FP=1` → `Q=0.9571 / FP=0`；对照腿只用修正后版本 |
| 5 | 钥匙/细则张力（S4/S5/S6：钥匙 Medium vs 细则 High） | 两读并报 + 声明钥匙未事后修改（见 §3） |
| 6 | ① 方案书此前**未落盘** | **已修正**：本会话落盘 `plan-v1.md` + `PLAN-PROVENANCE.txt`（见 §2） |

## 8. ⑤ 前硬前置三项（本会话收尾时状态）

1. ① 方案书落盘 → ✅ `tmp/glm-eval/plan-v1.md`（12,860 B，`8be638d2…`，7 节 / 8 种子行，无 BOM + LF）
2. 本片简报落盘 → ✅ `tmp/glm-eval/BRIEFING.md`（本文件）
3. 数据集完整性 → ✅ `docs/validation/directive-glm-eval-dataset.json` 共 **8 行**（探测 2 + ⑤ 2 + 批量 4），`pre_registration` 未变

## 9. 未完成事项（新会话按序）

1. ⑤ 报告预审（`DeepSeek V4 Pro (deepseek)` 审**本片评估报告**，非审 GLM 输出）；
2. ⑥ 越界/方法论审（`MAI-Code-1.1-Flash (copilot)`）；
3. ⑦ ×3（`GPT-5.3-Codex (copilot)`：盲审不附任何清单 → 终审 → 复审）；
4. ⑧a 报告对 `docs/validation/directive-glm-eval{,_EN}.md` + 自由度清单 `docs/validation/evidence/DIRECTIVE-GLM-EVAL-freedom-list.md` + 台账行 + 证据包 `docs/validation/evidence/glm-eval/`（`MANIFEST.md` + `SHA256SUMS.txt`，须过 G8-a/G8-b）；
5. ⑨ 原生核验 → ⑧b 回写 → ⑧c 清理（`tmp/` 仅留 `.gitkeep`；平台脚本先归档）→ ⑩ 回执（含 §11.2 回写检查行）。

## 10. 报告必载清单（用户裁决口径，⑧a 逐条核对）

1. 限定名探测记录（尝试 1 拒绝 + 可用清单作命名证据；注明「会话内转录，非原始载荷」）→ `PROBE-NOTES.txt`
2. 额度分账**逐行标包**（`glm-trial-2000w` / `glm-flash-realname-newuser` / `deepseek-n-a`）
3. 计量降级声明 + 成本判据未闭合 ⇒ 悬置（§6）
4. 钥匙/细则张力**两个 Q** + 钥匙未改声明（§3）
5. **n=1 边界**：只用「本样本中…」句式；禁写「GLM 优于/劣于 V4 Pro」
6. 批量腿并发偏离；7. 腿产物落盘口径不对称；8. 评分器修正前后两读；9. 越界自检结论（§12）
10. 「**只评估、不入白名单**」；白名单/角色决策推迟到本片之后

## 11. 环境与坑（实测）

- PowerShell 长命令链会被 `^C` 打断 ⇒ 每条命令尽量 ≤ ~250 字符、少用 `ForEach-Object` 管道。
- `create_file` 会写出 **CRLF** ⇒ `.md` / `.txt` 落盘后需归一为 LF（本会话用 node 一次性转换）。
- 终端 cwd 可能漂移到 `%TEMP%\glm-eval` ⇒ 先 `Set-Location d:\LZProjects\prfrail`。
- PS 控制台渲染 CJK 可能乱码 ⇒ 用「写文件 + 搜索工具读」而非看控制台。
- `tmp/**` 被 gitignore ⇒ 搜索需开 `includeIgnoredFiles`。
- `runSubagent` 可能返回 `Agent error: no response was returned` ⇒ **先查盘/转录本**再决定是否重发（避免重复付费调用）。
- 门禁命令：`node tools/gates/gate.js --all --scope=tree` ｜ `--selftest` ｜ `--check=G8-b --scope=tree`。
- `git check-ignore` 单次 ~417 ms ⇒ 门禁已批量化（勿改回逐文件）。

## 12. 越界自检（本会话实测结果）

- `git status --short` → 仅 `?? docs/validation/directive-glm-eval-dataset.json`（仓库侧唯一新增；报告/证据包尚未落盘）。
- `gate.js --all --scope=tree` → `changed=1`、`TOTAL_FAIL=0`（含 `PASS base go test ./... :: ok=14 FAIL=0`）。
- `gate.js --selftest` → `SELFTEST: PASS 328/328`。
- `git grep -n "GLM" -- docs/DELIVERY_DIRECTIVE.md docs/CONTRACTS.md schemas .github/agents` → **无输出**。
- `git grep -l '^model:' -- '.github/agents/prfrail-*'` → **无输出**（G5-a）。
- `git status --porcelain -- .github/agents internal schemas go.mod .github/workflows tools` → **无输出**。
- §2.2 / §2.4 / §5.2 / §6.3 / §7.5 与 §1.5.1 正文**零改动**；`.github/agents/**`、`internal/**`、`schemas/**`、契约夹具、`workflows`、`go.mod` 零改动；`cjk-newwords.txt` 未扩。
- 探针 **0**（未运行 `.vscode/scripts/agent-probe`）。
- 提示（预置免责）：`deep-reasoner.agent.md` / `independent-reviewer.agent.md` / `quick-verifier.agent.md` **含** `^model:` 行 —— 系**预先存在**、本片未触碰（`git status` 无输出），且不属 G5-a 管辖（G5-a 仅约束 `prfrail-*`）。

## 13. scratch 脚本清点（⑧c 清理前须先归档承重项）

- **承重（勿删，建议随证据包归档）**：`tmp/persist.js`（方案书逐字提取 + 冻结包复算 + 清单再生）、`tmp/freeze.js`（④ 冻结/预注册/隔离）、`tmp/save-leg.js`（腿产物逐字节保存 + 评分 + 污染）、`tmp/norm-score.js`、`tmp/batch-a-gen.js`、`tmp/batch-d-gen.js`、`tmp/fill-dataset.js`。
- **平台脚本**：`tmp/glm-eval/{hash.js,score.js,contam.js,selftest/**,metering/README.md}`。
- **一次性（本会话已删）**：`tmp/_scan*.js` 与 `tmp/_scan*/`（转录本检索中间物，内容冗余；方案书已落盘、命名证据已入 `PROBE-NOTES.txt`）。

## 追加（2026-10-05 补）·成本读数

- GLM 侧读数到位（见 `GLM-USAGE-READING.txt`）：`glm-5.3` 累计 161,863 tokens；`glm-5.3-flash` 累计 424,690 tokens。
- V4 Pro 对照腿成本读数：**不可得** ⇒ 成本对照为**单边**。
- 报告口径变更：
  - ⑤ 两腿 `metering_mode`：GLM 侧 = `user-provided-cumulative`；V4 Pro 侧保持 `unavailable`。
  - 成本判据：可给出 GLM 侧**绝对量级**；不可给出 GLM vs V4 Pro **相对结论**。
  - 不填估值、不做 token 换算。
- 新增必载项：实测消耗**显著高于主控预估**（约 4–8 倍），原因未核，如实登记。
- **编码事件披露（会话间外部保存）**：本文件在上一会话收尾时实测为 14,687 B / **无 BOM**；本会话追加前实测为 **14,690 B / 含 BOM**（`EF BB BF`）—— 差值恰为 **+3 B 的 BOM**，**内容逐字节相同**（已用聊天转录本重建 14,687 B 版本并逐字节比对证实，其 SHA-256 与上一会话所载 `d8463ef5…` 精确一致）。本块落盘时**已恢复为本文头部声明的「无 BOM + LF」口径**。同批变更提示中的 `PROBE-NOTES.txt` 经复核**未受影响**（2,023 B / 无 BOM / 哈希与上一会话记录一致）。如需改回仓库 `.md` 惯例（UTF-8 **with BOM**），转换会改变本节所载哈希，属预期变更而非偏差。

### 报告成本节必载项（⑧a 逐条写入）

1. GLM 侧绝对量级已获得（161,863 / 424,690，含少量用户早期测试消息）；
2. V4 Pro 侧读数**不可得** ⇒ **相对成本判据不成立**；
3. 对「是否值得把 GLM-5.3 纳入 ⑤」的成本问题给出**单边证据**：GLM-5.3 单次预审量级 ≈ 80k tokens，属可接受范围；
4. 相对成本比较**保留为悬置**，作为后续切片的输入。

- **主哈希变更（追加导致，非重写）**：原哈希 `d8463ef5…` → 新哈希（自排除口径）`8283094c40a7ab501a6c334f790ded0f77677357ca32f36d11ca063fd3be6ee0`（口径：本文件**末二行** —— 本行与 `append_sha256` 行 —— 不计入哈希计算，自指不可内嵌；`FROZEN-MANIFEST.txt` 所列为本文件**含**末行的完整哈希，二者口径不同、均自证；复算 `node tmp/persist.js`）
append_sha256=ba00db74cd52b0a8

## 追加（2026-10-05 补 2）·用户裁定（哈希权威 / 编码口径 / BOM 事件去向）

- **哈希权威（裁定三）**：本文件的验收值以 `FROZEN-MANIFEST.txt` 所列**完整哈希**为**唯一权威**；本文件内的 `append_sha256`（自排除口径，剔除末二行）**只作补充校验、非权威**。复核命令：`node tmp/persist.js`（重算清单并比对全部哈希），无需自排除口径。
- **编码口径（裁定二）**：`tmp/` 下**过程件**保持**无 BOM + LF**（本文件、`plan-v1.md`、`PROBE-NOTES.txt`、`GLM-USAGE-READING.txt` 等）。依据：§6.3 的 G4a 作用域 = 本次**待提交变更集**，而 `tmp/` 被 gitignore、**不在集内** ⇒ G4a 不适用；无 BOM 不影响 Node / 编辑器读取。若未来这些交接件**归档**进 `docs/validation/evidence/**`，届时在**归档执行片**内按 §1.5.1「归档后类型归一」（`.md` → with BOM + LF）统一转换并产生新哈希 —— **属预期流程、非偏差**。
- **BOM 事件去向（裁定一）**：外部保存定性确认（`base_sha = d8463ef5…` 精确一致 ⇒ 内容逐字节未变；已复原为无 BOM + LF；取证链 = `tmp/recon-briefing.js` + `tmp/fix-bom-and-append.js`）。⇒ **⑧a 报告须写入「异常与兜底记录」节**（本片新增必载项）。
- 其余确认（数据集 5 行累计口径、成本节必载项 4 条、冻结包逐字节未变、门禁全绿、授权边界）见上一节，用户已逐项确认。
- **主哈希变更（追加导致，非重写）**：本文件权威哈希 `580f0d9cb7d2f9a4ece41e894250084b8f374084dd84075f4306eca9c0231eea`（16,973 B）→ `6f0a081a7d795abc7b52ca2b48f582e99f72f95f3dd9392a3be912d2f399889d`（同链登记：`d8463ef5…` → `580f0d9c…` → 本值；口径：本文件**末二行** —— 本行与 `append_sha256` 行 —— 不计入哈希计算，自指不可内嵌；权威值始终以 `FROZEN-MANIFEST.txt` 为准）
append_sha256=b27db9b13c8841e7
