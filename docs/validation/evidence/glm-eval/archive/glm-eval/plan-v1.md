# 切片 `DIRECTIVE-GLM-EVAL` · ① 架构师评估方案书

DESIGN: DONE

> 前提声明（读入回执后定稿时核对）：本方案不执行任何付费调用、不改任何文件；合成场景的答案钥匙在本方案 §1 预注册，② 实现者构建 bundle 时逐字节嵌入，报告落盘先于调用。调用串 = `GLM-5.3 (glm)` / `GLM-5.3-Flash (glm)`（权威限定名，已由探测定）；对照腿 = `DeepSeek V4 Pro (deepseek)`；载体外壳统一 `deep-reasoner`（只读），实验变量仅 `model` 参数。GLM 不在 §2.2 白名单 ⇒ 本片调用属**用户授权的临时评估**，按 §1.7.4 五要素登记（OB/ADR，含「不得作为先例」），**不写入任何规则/白名单/角色集**。

## 1. ⑤ 对照腿的同输入构造（选型：候选 B，合成带种缺陷）

**选择候选 B，理由（A 被否）**：① A 的答案钥匙**已提交入库**（`025158f` 的验证报告 + 自由度清单全文载有 ⑤/⑦ 全部发现、编号与措辞），审查载体有 `read` 工具，软约束「不得读仓库」不可机械强制 ⇒ 对照必然失真；② A 的钥匙 = **V4 Pro 自己的发现**，用它评 V4 Pro 召回有系统性偏倚（恰是我们要测的变量）；③ A 的钥匙不完备（非全量真值）。B 实验窗口内**仓库不存在任何答案源**（合成场景与真实仓库无同文），污染面为零。

**Bundle（`tmp/glm-eval/bundle/`，全部 UTF-8 无 BOM + LF，任务包内**逐字节内联**，不给路径）**：
`00-plan.md`（合成 mini 方案书）、`01-clause-cn.txt`（合成契约原文 40–60 行，自造节号 X1.1/X6.3，措辞不与真实准则重合）、`02-diff.js`（约 150–200 行合成 `gx.js` + selftest 变更 hunk）、`03-prompt.txt`（任务包正文，两腿逐字节相同，**不写模型名**，只含 B.5 的输入/产出/禁止/结论行契约 + 「不得读取任何工作区或外部文件，仅依据本消息文本；报告以最终消息返回」）。

**预注册种子钥匙（写入报告「答案钥匙」节，先于两腿调用落盘；bundle 内不得出现 S 编号）**：

| 种子 | 严重度 | 缺陷 | 可定位性 |
|---|---|---|---|
| S1 | Critical | `gx.js` 解析 SHA256SUMS 失败/白名单缺失时 `catch { return true }`（fail-open，静默通过） | 唯一 catch 块 |
| S2 | High | 包成员 vs 清单比较未做 `\`→`/` 归一 ⇒ Windows 包恒判红 | 路径比较行 |
| S3 | High | 复算哈希时读 SHA256SUMS 行第 1 列（文件名列）当哈希 ⇒ 比较恒真 | split 取列行 |
| S4 | Medium | 契约声明「未入库清单与包内文件冲突判红」，实现无该分支 | 分支缺失 |
| S5 | Medium | 契约声明四子查，实现只落地三支（缺重建注记检查） | 枚举缺失 |
| S6 | Medium | 「只判变更集」被实现为全树扫描（与契约冲突，历史误报） | 作用域行 |
| S7 | Low | 顶部注释写「八类」实为四类（陈旧注释） | 注释行 |
| S8 | — | 纯重排/空行（无缺陷） | 期待 0 发现 |

**评分细则（预注册，五维）**：召回 `recall_all=命中/7`、`recall_high=命中(S1–S3)/3`；假阳性 = 引用行不存在或断言与 bundle 不符的发现（**额外真发现**经主控机械复核后单列，不计 FP）；severity 校准（High=3/Medium=2/Low=1，精确=1.0，差 1 档=0.5，差 ≥2=0）；结论行形态 = 末行恰为 `PRE-REVIEW: PASS|PASS WITH FIXES|FINDINGS` 之一（§4.1 闭集）⇒ 1/0；越界自检 = 输出含「改规则/白名单/角色集」建议、引用本工作区路径、或声称执行命令 ⇒ 逐条记录（不进 Q）。**复合分 Q = 0.4×recall_high + 0.2×recall_all + 0.2×sev_exact + 0.2×conclusion_line**；Q 仅作描述量，n=1 不得下模型优劣结论。

**污染检测（对两腿输出机械执行）**：`Select-String -Path tmp\glm-eval\out\*.md -Pattern 'OB-79|025158f|0c12abc|G8-[abcd]|DIRECTIVE-EVIDENCE-SCOPE|重建归档|判据作用域排除|S[1-8]'` 须 0 命中；另核对每条 `file:line` 是否落在 bundle 内路径。命中即记污染指标并降级该腿结论。

## 2. token 计量口径

**GLM 控制台（累计读数）逐步规程**（每次调用，串行执行）：
1. 声明窗口独占：本片按 R5.1 串行，调用前后 5 分钟内无其它 GLM 活动（用户侧并行使用须如实记录）；
2. 调用前：用量页（按模型×包筛）导出截图/CSV → `tmp/glm-eval/metering/{seq}-pre.*`，记录 `pre_total` + 时戳；
3. 发起调用（记录 `startedAt`）；
4. 调用后立即读 `post_total`，**隔 10 分钟复核**一次，两次一致才采信；
5. 归属：`measured_tokens = post_total − pre_total`，落在 GLM-5.3 的**尝鲜包**或 Flash 的**实名+新用户包**（数据集逐行标包）；控制台若提供按请求明细导出则取明细为第二来源；
6. 读数不更新或仅聚合 ⇒ 降级 `metering_mode=interval-sum`、`measured_tokens=NULL`，行内如实写窗口与两次读数，**不得估摸**；
7. 两腿输入逐字节相同 ⇒ 输入 token 相等按构造成立，只比输出 token。

**V4 Pro 对照腿：声明「token 不可得」**——理由：① 经第三方模型通道（非 Copilot 计费），会话调试日志无该通道 token 字段；② §7.5 的 R7.1 真值源（会话调用回执/日志）只记调用次数、不记 token。替代口径：两腿均记**输出字符数**（CJK+ASCII 计码点）与耗时，明确标注「字符数 ≠ token、不可换算」；若用户事后提供 DeepSeek 平台侧账单明细可作第二来源补记并标来源。

## 3. ④/⑤ 的顺序与 ⑥ 判定

**本片流水线序**：①（本方案）→ ②（构建 bundle + 数据集脚手架 + 评分/哈希脚本，落 `tmp/glm-eval/`）→ ③（上述脚本的机械测试）→ ④（`gate.js --all --scope=tree` 全绿；bundle 冻结 + 哈希落盘 + **钥匙预注册进报告草稿**）→ ⑤ = 两条实验腿（GLM 先、V4 Pro 后；同输入、输出分别落盘不互锚）→ ⑥ → ⑦ 三读 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩。

**⑥ 逐条核对（§7.2.1）⇒ 不可跳过**：① 不成立/存疑——数据集 `.json` 是机器结构文件，不在「文档与配置文案类」字面范畴；② 不成立——本片**新增 `.json` 数据集**（结构化、承载评分与哈希）；③ 不成立——本片产出**⑤ 位置的评估方法论与证据数据集**，直接服务审查链模型选型判断、新增证据形态（评估数据集 + 冻结评估包），触「审查体系/证据模型」边缘。另 §7.2.1 规定跳过必须 ⑦/用户确认、主控不得自行认定——本片不寻求跳过。

**⑥ 输入包**：本方案冻结版 + 数据集 schema + 评分细则与预注册钥匙 + 本片变更 diff（报告/台账/冻结包）+ 哈希复算脚本输出 + 污染检测结果；实验腿原始输出只给路径与哈希（抽查读原文件）。

## 4. 失效处置（无重试额度）

**闭集枚举（`failure_mode` 列取值）**：`OK / EMPTY（空输出）/ TRUNCATED（截断）/ NO_CONCLUSION_LINE / FORMAT_DEVIATION（非 Section 形态或结论行非法）/ REFUSAL（拒答）/ PRODUCT_UNUSABLE（批量腿：产物不可用）`。

**判定规则**：末行未命中 §4.1 闭集 ⇒ `NO_CONCLUSION_LINE` 或 `FORMAT_DEVIATION`；`TRUNCATED` 指输出未以结论行结尾且明显中段截断。**一律记为结果**——不重试、不换提示词「救活」、不改写提示词重新计价；token 仍按窗口读数记入该行（能归属则记，不能则 NULL）。

**影响口径**：失效腿的 Q 与结论行均为空，`quality_verdict=该失效形态`；GLM 腿失效 ⇒ 该位置记为「GLM 本样本失效 n=1」，对照腿照常评分；**双腿均失效 ⇒ 无对照结论**，报告只输出失效事实；批量腿失效 ⇒ 该类任务「产物不可用」，不得二次调用补测。

## 5. 数据集 schema（`docs/validation/directive-glm-eval-dataset.json`，无 BOM + LF，逐次调用一行）

| 列 | 取值域 |
|---|---|
| `call_no` | 1–6（⑤-实验/⑤-对照/a/b/c/d） |
| `role_position` | `⑤-EXPERIMENT / ⑤-CONTROL / BATCH-a / BATCH-b / BATCH-c / BATCH-d` |
| `model_string` | `GLM-5.3 (glm) / DeepSeek V4 Pro (deepseek) / GLM-5.3-Flash (glm)` |
| `quota_package` | `glm-trial-2000w / glm-flash-realname-newuser / deepseek-n-a` |
| `input_summary` | bundle 结构描述 + prompt 首 120 字符 |
| `input_hash` / `output_hash` | SHA-256（64 hex） |
| `conclusion_line` | 原样末行（或失效形态） |
| `measured_tokens` / `metering_mode` | 整数或 `null`；`exact-delta / interval-sum / unavailable` |
| `duration_s` | 调用墙钟耗时 |
| `quality_verdict` | `OK/EMPTY/…（§4 闭集）`；⑤ 腿另记 `Q` 与五维分 |
| `notes` | 污染标记、窗口读数、第二来源、越界记录 |

**哈希与固化**：SHA-256；输入对象 = `bundle/` 六文件按 `00…03` 固定序拼接字节 + 拼接 `03-prompt.txt` 原文，命令：`$b = (Get-Content tmp\glm-eval\bundle\* -Raw -Encoding UTF8) -join "`n"; [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData([Text.Encoding]::UTF8.GetBytes($b)))`（或逐文件 `Get-FileHash -Algorithm SHA256` + 目录 `SHA256SUMS.txt`）；输出对象 = 返回文本原文（无 BOM）哈希 + 归一化副本哈希各一。bundle 与两腿输出随后归档入 `docs/validation/evidence/glm-eval/` 冻结包（`MANIFEST.md` + `SHA256SUMS.txt`，过 G8-a/G8-b）。

## 6. 越界自检与报告结构

**报告 `docs/validation/directive-glm-eval{,_EN}.md` 按 §11.3**：变更概要 → 执行流水线 → 异常与兜底 → 环境与档位 → 审查结论摘要（⑤=实验腿对照表、⑥、⑦）→ 审查输入包摘要（含钥匙预注册时戳与哈希）→ 可证伪性（哈希复算 + 评分脚本测试 + 污染 grep）→ 门禁结果 → ⑧c 记录 → 成本与计量（两包分账）→ 明确未执行事项 → 下一步建议（含「GLM 纳入与否另议」）。

**必须证明未改清单 + 机械核验命令**（收尾时逐条执行并留痕）：
```
git status --porcelain -- docs/DELIVERY_DIRECTIVE.md docs/DELIVERY_DIRECTIVE_EN.md docs/CONTRACTS.md docs/CONTRACTS_EN.md .github/agents .github/workflows internal schemas go.mod tools/gates/testdata
# 全部无输出；tools/gates/lib/**、gate.js、selftest.js、evidence-forms.txt 亦须无输出
grep -l "^model:" .github/agents/prfrail-*.agent.md          # 无输出（G5-a）
node tools/gates/gate.js --selftest                            # 必须仍为 SELFTEST: PASS 328/328
node tools/gates/gate.js --all --scope=tree                    # TOTAL_FAIL=0
node tools/gates/gate.js --check=G8-b --scope=tree             # glm-eval 包 self-consistent
git grep -n "GLM" docs/DELIVERY_DIRECTIVE.md docs/CONTRACTS.md schemas .github/agents  # 无输出
```
即：§2.2/§2.4/§5.2/§6.3/§7.5 与 §1.5.1 正文零改动；`.github/agents/**`、`internal/**`、`schemas/**`、契约夹具、`workflows`、`go.mod` 零改动；`cjk-newwords.txt` 不扩（报告用词复用既有语料，新字以既有词改写——freedom-list 异常记录 1 先例）。仓库落盘仅限：报告对、数据集 JSON、`docs/validation/evidence/glm-eval/` 冻结包、`REMAINING_SLICES{,_EN}` 台账行/完成行、本片自由度清单（根级 `*-freedom-list.md` 形态，G8-a 认可）。

## 7. 风险与缓解

| 风险 | 缓解 |
|---|---|
| 污染/泄题（钥匙或真实仓库被读） | 选 B（仓库无答案源）；钥匙仅存于 tmp 与聊天转录；任务包禁读仓库；输出 grep 标记词表 + bundle 内路径核对 |
| n=1 统计效力 | 结论只用「本样本中…」句式；禁写「GLM 优于/劣于 V4 Pro」；Q 标为描述量 |
| 控制台累计口径 vs 实测差（滞后/聚合/外部流量） | 前后双读 + 10 分钟复核 + 窗口独占声明；不更新即降级 `interval-sum` 并如实标注 |
| V4 Pro token 不可得 ⇒ 对比不对称 | 提前声明「不可得 + 理由」；以输出字符数/耗时为公共代理，输入按构造相等 |
| 语言侧检查 vs G3 职责边界 | (c) 腿任务包明示「只做语言侧意见（漏译/错译/增译/术语），**不得输出通过/不通过判定**」；G3 结论以 `gate.js` 为准，Flash 输出仅作质量信号，防判据外判 |
| 失效形态误读为模型能力 | 失效即记结果、不重试，报告区分「失效形态」与「质量分」，避免幸存者偏差 |
| 评分脚本自身错误 | ③ 对评分/哈希/污染脚本做正反夹具测试；评分结果由 ⑥/⑦ 复核钥匙一致性 |
| 载体/模型串不确定性（§2.3） | 以 operator UI 设置与 `runSubagent` 参数记录为准，不采信自报；GLM 调用串仅用已探测成功的 `GLM-5.3 (glm)` 形态 |
