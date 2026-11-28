# P98 `DIRECTIVE-REVIEW-TIER-ALIGN` ① 架构方案书（停点 1）

- 载体：deep-reasoner；模型：DeepSeek V4 Pro (deepseek)（Max 档，用户 2026-11-28 裁定，不得作先例）。
- 本环产出：仅本方案书；不编辑任何其它文件。停点 1 四项候选逐项给候选结论 + 理由，交用户裁定。
- 行号引用均以 2026-11-28 磁盘实况为准（CN 侧用 `#L` 形态；EN 侧与 CN 同位）。
- 编码注记：本件由写入工具产出，可能为 UTF-8 无 BOM + LF；主控在 ② 前按仓库约定归一为 UTF-8 with BOM + LF。
- 结论行：`DESIGN: DONE`（回执见文末）。

---

## ① 目标与边界

**目标**（对应起跑消息「修订要求 a–d」）：

1. **O-2（High，Hydrogen #1）**：§5.2 深度档行档位列拆明「Sol 主读取 Max（须同轮授权）；Neon 交叉复核取 Extra High」，与 §8.2 逐字一致；§8.2 取档授权句同步对齐，消解「Neon 究竟 Max 还是 Extra High、是否需授权」与「Luna 唯一档 Max / Sol 深度档 Max 的授权口径」两处歧义。
2. **O-3（Medium，Hydrogen #2）**：§5.2 矩阵补字面可判性措辞（五子项：轻量档术语定义 / 大上下文文档档逻辑连接词与和轻量档的优先级 / 常规档「其余 N 档」不漏计 / Sol 备选与 Luna 平行备选的预算触发区分门槛 / 「最高风险轮」与一般高风险轮次的界线）。
3. **O-4（Medium，Hydrogen #3）**：§7.5 Neon 两档单价齐备 + 主表 300 / 1500 标注档位归属。
4. **O-5 #4（Low）**：§7.5 备选额度口径首句措辞消除误读。
5. **O-5 #5（Low）**：P97 `numstat.txt` #L16 自记行数陈旧——**冻结件不修，只登记「已知、不追改」**。
6. 双语镜像：三节修订文本 CN/EN 同构（G3(a.2) 六形状：行数 / 标题行号 / 加粗标记数 / 竖线数 / 空行数 / 反引号数），新增 CJK 字全在 HEAD 全库 `.md` 语料 ∪ `tools/gates/cjk-newwords.txt` 内（G4b），不得新增白名单字。

**边界**：只改 `docs/DELIVERY_DIRECTIVE{,_EN}.md` 的 §5.2 / §7.5 / §8.2 相关行 + 版本包（头部 / §0.3 新行 / 附录 C 新块 C.0bf）+ 台账 `docs/t027/REMAINING_SLICES{,_EN}.md` P97 / P98 段 + 报告对与证据包。**不改** `tools/gates/**`（gate.js / selftest.js / lib/criteria/**）、`internal/**`、`ci.yml`、CONTRACTS / schema / fixtures / workflows、§2.2 / §2.4 / §6.3 判据本体与 G5-b 闭集、`cjk-newwords.txt` 与 `g6-whitelist.txt`；**不新增 ADR**（判定见 ⑥）；P87 报告对与证据包、CHECKLIST_v81 冻结不追改；不做 gitee。

---

## ② 现文问题定位（逐处行号 + 现文摘录 + Hydrogen #N 对应）

### O-2（Hydrogen #1 High）—— §5.2 深度档档位列 vs §8.2 直接冲突

| 位置 | 现文摘录（2026-11-28 磁盘实况） |
|---|---|
| CN §5.2 #L582（EN 同位 #L582） | `| 深度档 | \`gpt-6.2-sol\` 主读；交叉复核 \`gemini-4-neon\` | Max | 硬门 / 高风险轮次：A7 / B2 / B3 / B4 级切片、安全 / 机器状态 / 协议语义、发布前终审；最高风险轮必配 Neon 交叉复核；主读报 Medium+ 或存在结构疑点 ⇒ 追加 Neon 深挖 |` |
| CN §8.2 #L941（EN 同位 #L941） | `| \`gemini-4-neon\` | **Extra High** |` |
| CN §8.2 #L943（EN 同位 #L944） | `**取档授权**：\`gpt-6.2-sol\` 取 \`Max\` 须在起跑消息中显式声明并附同轮授权；未声明时默认 \`Extra High\`（与 §7.5「不得中途升档」一致）。三枚 Gemini 成员按 §5.2 选型矩阵的档位取档、无需追加取档授权；任何 \`Max\` 取档仍须同轮授权。` |

冲突三角：① 深度档行把 Sol 主读与 Neon 交叉复核合并一行、思考档位单写 `Max`，未拆明两者分别取档；② §8.2 对 `gemini-4-neon` 规定 `Extra High` ⇒ Neon 究竟取 Max 还是 Extra High 及是否需授权直接冲突；③ `gpt-6-luna` 唯一档为 `Max`（§5.2 #L585 / §8.2 #L940），被「任何 `Max` 取档仍须同轮授权」覆盖而 §5.2 平行备选档未注授权。

### O-3（Hydrogen #2 Medium）—— 矩阵字面可判性五子项

| 子项 | 位置 | 现文摘录 / 缺口 |
|---|---|---|
| ① 轻量档术语无定义 | CN §5.2 #L584 | `| 轻量档 | \`gemini-4-hydrogen\` | Extra High | 文档型 / 治理片的定稿复核与定向小读（低成本复核） |` ——「定稿复核」「定向小读」全文无定义 |
| ② 大上下文文档档逻辑连接词缺失 + 与轻量档无优先级 | CN §5.2 #L583 | `审查输入包超 400K token 且为纯文档 / S 档` —— `/` 未明示且或关系；>400K 纯文档 S 档定稿复核时本档与轻量档同时命中、无优先级 |
| ③ 常规档「其余三档」漏计平行备选档 | CN §5.2 #L581 | `默认档；其余三档判据均不满足时使用` —— 矩阵共五档，除常规档外有四档（深度 / 大上下文文档 / 轻量 / 平行备选），「三档」漏计第 5 行 |
| ④ Sol 备选与 Luna 平行备选预算触发无区分门槛 | CN §5.2 #L581 / #L585 | #L581「预算紧张轮取 Sol 备选」与 #L585「`gpt-5.3-codex` 不可用或 §7.5 预算受限时的平行替换」同以预算受限触发、无客观区分 |
| ⑤ 「最高风险轮」与一般高风险轮次无界线 | CN §5.2 #L582 | 「最高风险轮必配 Neon 交叉复核」中「最高风险轮」未定义（现文仅适用域排除句 #L591 定义了「深水轮」） |

### O-4（Hydrogen #3 Medium）—— §7.5 Neon 两档单价与档位归属缺

| 位置 | 现文摘录 |
|---|---|
| CN §7.5 #L852（EN 同位 #L852） | `| \`gemini-4-neon\` | 300 | 1500 |` —— 未标注属 200K 档还是 1M 档 |
| CN §7.5 #L854（EN 同位 #L855） | `…\`gemini-4-carbon\` 与 \`gemini-4-hydrogen\` 仅 1M 档，上表单价即实价；\`gemini-4-neon\` 分 200K / 1M 两档；\`gpt-5.3-codex\` 实采 400K。…` —— 分两档但只给一份单价（对比 Sol 明确给了 1M 档 400 / 1500）；另一档缺失将阻断 §7.5 兜底句下的事前预算估算 |

### O-5 #4（Low）—— §7.5 备选额度口径首句误读张力

| 位置 | 现文摘录 |
|---|---|
| CN §7.5 #L856（EN 同位 #L857） | `**备选额度口径（P87 落定）**：备选档（常规档的 Sol 备选）⑦ 额度**与深度档同口径**——均按 §7.5 的切片档行取值（L 档 / 治理片 / 高风险轮次 = 1 终审 + 3 复审）；轻量档（\`gemini-4-hydrogen\`）同常规档。模型档只决定模型与单价，不改变 §7.5 按切片档分层的次数上限。` —— 首句「与深度档同口径」与末句「不改变按切片档分层」并存，易误读为 Sol 备选恒享 1+3 |

### O-5 #5（Low）—— 冻结件不修

- 位置：`docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/…/numstat.txt` #L16 对自身仅记 `27 0`（实际写满 35 行后未回采）。
- 处置：**只登记「已知、不追改」**（冻结证据包件；登记位置 = P98 报告 + 台账 P97 段观察项处置行）。

---

## ③ 修订文本三件套（CN/EN 逐字拟文 + 锚点）

> 槽位说明：本节拟文中以 `〈…〉` 标出的槽位属停点 1 待裁定内容，**不得以槽位形态进入准则正文**（G2）；② 落定时按用户裁定填实后再行修订。
> G 判据自查：拟文无「待…回填 / 尚未提交 / 待同轮授权 / pending the push」形态（G2）；无「数字（2 位+）＋行」「数字＋文件」「共＋数字」形态（G6-4）；节号引用 §5.2 / §7.5 / §8.2 / §1.7.2 均存在（G6-5）；新增 CJK 字全部来自现文语料（G4b，机械预检结果见 §九）。

### 三件套之一：§5.2（CN 六处 + EN 同位六处）

**【替换 1｜CN §5.2 #L582｜EN §5.2 #L582】深度档行**

- 替换前文（CN）：`| 深度档 | \`gpt-6.2-sol\` 主读；交叉复核 \`gemini-4-neon\` | Max | 硬门 / 高风险轮次：A7 / B2 / B3 / B4 级切片、安全 / 机器状态 / 协议语义、发布前终审；最高风险轮必配 Neon 交叉复核；主读报 Medium+ 或存在结构疑点 ⇒ 追加 Neon 深挖 |`
- 拟文（CN）：`| 深度档 | \`gpt-6.2-sol\` 主读；交叉复核 \`gemini-4-neon\` | Sol 主读取 Max（须同轮授权）；Neon 交叉复核取 Extra High | 硬门 / 高风险轮次：A7 / B2 / B3 / B4 级切片、安全 / 机器状态 / 协议语义、发布前终审；深水轮（定义见下方适用域排除句）必配 Neon 交叉复核；其余高风险轮次在主读报 Medium+ 或存在结构疑点 ⇒ 追加 Neon 深挖 |`
- 替换前文（EN）：`| Deep tier | \`gpt-6.2-sol\` main read; cross-check \`gemini-4-neon\` | Max | hard-gate / high-risk rounds: A7 / B2 / B3 / B4-class slices, security / machine state / protocol semantics, pre-release final review; the highest-risk round always pairs a Neon cross-check; when the main read reports Medium+ or a structural doubt remains, add a Neon deep-dive |`
- 拟文（EN）：`| Deep tier | \`gpt-6.2-sol\` main read; cross-check \`gemini-4-neon\` | Sol main read takes Max (same-round authorization needed); Neon cross-check takes Extra High | hard-gate / high-risk rounds: A7 / B2 / B3 / B4-class slices, security / machine state / protocol semantics, pre-release final review; a deep-water round (defined in the applicability-exclusion sentence below) always pairs a Neon cross-check; other high-risk rounds add a Neon deep-dive when the main read reports Medium+ or a structural doubt remains |`
- 理由：档位列与 §8.2 修订后取档授权句**逐字一致**（`Sol 主读取 Max（须同轮授权）` / `Neon 交叉复核取 Extra High`）；「最高风险轮」→「深水轮（定义见适用域排除句）」消除未定义术语，并划出与一般高风险轮次的界线（深水轮必配、其余高风险轮触发式追加）——覆盖 O-2 与 O-3 ⑤。

**【替换 2｜CN §5.2 #L581｜EN §5.2 #L581】常规档行**

- 替换前文（CN）：`| 常规档 | \`gpt-5.3-codex\`；性价比备选 \`gpt-6.2-sol\` | Extra High（Sol 取 Max 须同轮授权） | 默认档；其余三档判据均不满足时使用；预算紧张轮取 Sol 备选 |`
- 拟文（CN）：`| 常规档 | \`gpt-5.3-codex\`；性价比备选 \`gpt-6.2-sol\` | Extra High（Sol 取 Max 须同轮授权） | 默认档；其余四档判据均不满足时使用；预算紧张轮取 Sol 备选（Codex 可用时的低成本替换） |`
- 替换前文（EN）：`| Regular tier | \`gpt-5.3-codex\`; cost-effective alternative \`gpt-6.2-sol\` | Extra High (Sol taking Max needs a same-round authorization) | default tier; used when none of the other three tiers' criteria is satisfied; budget-tight rounds take the Sol alternative |`
- 拟文（EN）：`| Regular tier | \`gpt-5.3-codex\`; cost-effective alternative \`gpt-6.2-sol\` | Extra High (Sol taking Max needs a same-round authorization) | default tier; used when none of the other four tiers' criteria is satisfied; budget-tight rounds take the Sol alternative (the lower-cost substitute while Codex is available) |`
- 理由：覆盖 O-3 ③（「其余四档」不漏计平行备选档）与 O-3 ④ 的 Sol 侧门槛（「Codex 可用时」= 与 Luna 的区分）。

**【替换 3｜CN §5.2 #L583｜EN §5.2 #L583】大上下文文档档行**

- 替换前文（CN）：`| 大上下文文档档 | \`gemini-4-carbon\` | High（固定） | 审查输入包超 400K token 且为纯文档 / S 档 |`
- 拟文（CN）：`| 大上下文文档档 | \`gemini-4-carbon\` | High（固定） | 审查输入包超 400K token 且（为纯文档或为 S 档） |`
- 替换前文（EN）：`| Large-context document tier | \`gemini-4-carbon\` | High (fixed) | the review input pack exceeds 400K tokens and the slice is pure-document / S-tier |`
- 拟文（EN）：`| Large-context document tier | \`gemini-4-carbon\` | High (fixed) | the review input pack exceeds 400K tokens and (the slice is pure-document or S-tier) |`
- 理由：覆盖 O-3 ② 的连接词缺失（「且（为纯文档或为 S 档）」= 与 + 或两级明示）；与轻量档的优先级由下方新增「档位取档次序」句承载。

**【替换 4｜CN §5.2 #L584｜EN §5.2 #L584】轻量档行**

- 替换前文（CN）：`| 轻量档 | \`gemini-4-hydrogen\` | Extra High | 文档型 / 治理片的定稿复核与定向小读（低成本复核） |`
- 拟文（CN）：`| 轻量档 | \`gemini-4-hydrogen\` | Extra High | 文档型 / 治理片的定稿复核（对 ⑧a / ⑧b 稿或报告终稿的最终一致性复核，不重审全文）与定向小读（输入收窄至单一发现闭合路径或单节文本）；低成本复核 |`
- 替换前文（EN）：`| Light tier | \`gemini-4-hydrogen\` | Extra High | finalization reviews and targeted small reads of document-type / governance slices (low-cost review) |`
- 拟文（EN）：`| Light tier | \`gemini-4-hydrogen\` | Extra High | finalization reviews (a final consistency review of the (8a)/(8b) drafts or the report final draft, not a full re-read) and targeted small reads (input narrowed to a single finding-closure path or a single section) of document-type / governance slices; low-cost review |`
- 理由：覆盖 O-3 ① —— 为「定稿复核」「定向小读」锚定可判定义（复核对象、范围、输入形态），不加新术语。

**【替换 5｜CN §5.2 #L585｜EN §5.2 #L585】平行备选档行**

- 替换前文（CN）：`| 平行备选档 | \`gpt-6-luna\` | Max | \`gpt-5.3-codex\` 不可用或 §7.5 预算受限时的平行替换 |`
- 拟文（CN）：`| 平行备选档 | \`gpt-6-luna\` | Max（取档须同轮授权） | \`gpt-5.3-codex\` 不可用或常规档两成员均不可承担时的平行替换 |`
- 替换前文（EN）：`| Parallel-alternative tier | \`gpt-6-luna\` | Max | the parallel substitute when \`gpt-5.3-codex\` is unavailable or the §7.5 budget is constrained |`
- 拟文（EN）：`| Parallel-alternative tier | \`gpt-6-luna\` | Max (taking it needs a same-round authorization) | the parallel substitute when \`gpt-5.3-codex\` is unavailable or both regular-tier members cannot carry the round |`
- 理由：覆盖 O-2 ③（Luna 唯一档 Max 注授权）与 O-3 ④ 的 Luna 侧门槛——触发词由「预算受限」改为「常规档两成员均不可承担」，与 Sol 备选（Codex 可用时的降本替换）形成两级区分，预算受限不再单独触发本档（由次序句明示）。

**【插入 1｜CN §5.2 矩阵表后、判级门槛共识句（#L587）之前｜EN 同位】新增「档位取档次序」段**

- 拟文（CN）：`**档位取档次序（P98 落定）**：同时命中多档时取优先档——深度档 > 大上下文文档档 > 轻量档 > 常规档（默认档）；平行备选档在常规档两成员均不可承担时启用（判据见该档行），预算受限不单独触发本档。`
- 拟文（EN）：`**Tier-taking order (landed by P98)**: when several tiers match, take the priority tier - deep tier > large-context document tier > light tier > regular tier (default); the parallel-alternative tier is enabled only when both regular-tier members cannot carry the round (see that tier's criterion), and a budget constraint alone does not trigger it.`
- 理由：一条次序句同时解决 O-3 ② 的优先级（大上下文 > 轻量）与 O-3 ④ 的触发互斥（平行备选不因预算单独触发）；常规档「其余四档均不满足」与新次序自洽（次序最低）。

**【替换 6｜CN §5.2 #L591｜EN §5.2 #L591】适用域排除句末句（「触发式」→「必配」对齐）**

- 替换前文（CN）：`（P85 四读实证两者同输入均漏已验证的真实 Medium；深水轮 = \`gpt-6.2-sol\` Max 主读 + \`gemini-4-neon\` 触发式交叉复核）。`
- 拟文（CN）：`（P85 四读实证两者同输入均漏已验证的真实 Medium；深水轮 = \`gpt-6.2-sol\` Max 主读 + \`gemini-4-neon\` 必配交叉复核）。`
- 替换前文（EN）：`(the P85 four-read evidence shows both missed a verified real Medium on the same input; a deep-water round = \`gpt-6.2-sol\` Max main read plus a triggered \`gemini-4-neon\` cross-check).`
- 拟文（EN）：`(the P85 four-read evidence shows both missed a verified real Medium on the same input; a deep-water round = \`gpt-6.2-sol\` Max main read plus a mandatory \`gemini-4-neon\` cross-check).`
- 理由：与替换 1 的「深水轮必配 Neon 交叉复核」逐字一致（消除现文「必配」vs「触发式」的残留张力）。

### 三件套之二：§8.2（CN 一处 + EN 同位一处）

**【替换 7｜CN §8.2 #L943｜EN §8.2 #L944】取档授权句**

- 替换前文（CN）：`**取档授权**：\`gpt-6.2-sol\` 取 \`Max\` 须在起跑消息中显式声明并附同轮授权；未声明时默认 \`Extra High\`（与 §7.5「不得中途升档」一致）。三枚 Gemini 成员按 §5.2 选型矩阵的档位取档、无需追加取档授权；任何 \`Max\` 取档仍须同轮授权。`
- 拟文（CN）：`**取档授权**：\`gpt-6.2-sol\` 取 \`Max\`（含深度档 Sol 主读取 Max）须在起跑消息中显式声明并附同轮授权；未声明时默认 \`Extra High\`（与 §7.5「不得中途升档」一致）。\`gpt-6-luna\` 唯一档为 \`Max\`、取档须同轮授权（与 Sol 同口径）。三枚 Gemini 成员按 §5.2 选型矩阵取档——\`gemini-4-neon\` 交叉复核取 \`Extra High\`、\`gemini-4-hydrogen\` 取 \`Extra High\`、\`gemini-4-carbon\` 取 \`High\`（固定）——均为非 \`Max\` 档、无需追加取档授权。`
- 替换前文（EN）：`**Tier authorization**: taking \`Max\` with \`gpt-6.2-sol\` requires an explicit declaration in the kickoff message plus a same-turn authorization; without that declaration the default is \`Extra High\` (consistent with "no mid-slice tier upgrade" in §7.5). The three Gemini members take their tier per the §5.2 selection matrix with no extra tier authorization; any \`Max\` tier still needs a same-round authorization.`
- 拟文（EN）：`**Tier authorization**: taking \`Max\` with \`gpt-6.2-sol\` (including the deep-tier Sol main read taking Max) requires an explicit declaration in the kickoff message plus a same-turn authorization; without that declaration the default is \`Extra High\` (consistent with "no mid-slice tier upgrade" in §7.5). \`gpt-6-luna\` has \`Max\` as its only tier and taking it needs a same-round authorization (same caliber as Sol). The three Gemini members take their tiers per the §5.2 selection matrix - \`gemini-4-neon\` cross-check takes \`Extra High\`, \`gemini-4-hydrogen\` takes \`Extra High\`, \`gemini-4-carbon\` takes \`High\` (fixed) - all non-\`Max\` tiers, with no extra tier authorization.`
- 理由：两处歧义一并消解——Neon = `Extra High`、无需追加授权（与 §5.2 新档位列逐字一致）；Luna 唯一档 `Max` 取档须同轮授权（与 Sol 深度档 Max 同口径）；原末句「任何 Max 取档仍须同轮授权」具体化进三成员表述，不再产生覆盖歧义。映射表 #L941 Neon 行维持 `Extra High` 不动。
- 修订不扩大解释到 §2.2 模型名单：本句只述档位与授权口径，不涉白名单闭集（G5-b 不动）。

### 三件套之三：§7.5（CN 两处 + EN 同位两处）

**【替换 8｜CN §7.5 #L854｜EN §7.5 #L855】上下文档位计价注记（Neon 两档单价 + 档位归属）**

> 归属与另一档单价均无实测真值源 ⇒ **不得估算**。两条候选口径见 §九（推荐 B；A 为兜底）。以下按口径分别给拟文（槽位在 ② 前填实）。

- 替换前文（CN）：`**上下文档位计价注记（P87 落定）**：\`gpt-6.2-sol\` 的 1M 档为 400 / 1500（输入为上表两倍、输出为上表 1.5 倍）；\`gemini-4-carbon\` 与 \`gemini-4-hydrogen\` 仅 1M 档，上表单价即实价；\`gemini-4-neon\` 分 200K / 1M 两档；\`gpt-5.3-codex\` 实采 400K。Cache 读写单价（Credits / 1M tokens）分别 = Neon 15 / 0、Carbon 10 / 0、Hydrogen 7.5 / 0。`
- 拟文（CN，口径 B）：`**上下文档位计价注记（P87 落定；P98 补 Neon 档位归属）**：\`gpt-6.2-sol\` 的 1M 档为 400 / 1500（输入为上表两倍、输出为上表 1.5 倍）；\`gemini-4-carbon\` 与 \`gemini-4-hydrogen\` 仅 1M 档，上表单价即实价；\`gemini-4-neon\` 分 200K / 1M 两档——上表 300 / 1500 为〈归属档〉、另一档单价为〈Input / Output，2026-11-28 用户 operator UI 实读〉；\`gpt-5.3-codex\` 实采 400K。Cache 读写单价（Credits / 1M tokens）分别 = Neon 15 / 0、Carbon 10 / 0、Hydrogen 7.5 / 0。`
- 拟文（CN，口径 A-甲：上表标 200K 档、1M 档兜底）：`**上下文档位计价注记（P87 落定；P98 补 Neon 档位归属）**：\`gpt-6.2-sol\` 的 1M 档为 400 / 1500（输入为上表两倍、输出为上表 1.5 倍）；\`gemini-4-carbon\` 与 \`gemini-4-hydrogen\` 仅 1M 档，上表单价即实价；\`gemini-4-neon\` 分 200K / 1M 两档——上表 300 / 1500 为 200K 档；1M 档单价计价页不可得 ⇒ 按本节兜底句停下升级用户裁决（不得估算或填写）；\`gpt-5.3-codex\` 实采 400K。Cache 读写单价（Credits / 1M tokens）分别 = Neon 15 / 0、Carbon 10 / 0、Hydrogen 7.5 / 0。`
- 拟文（CN，口径 A-乙：上表标 1M 档、200K 档兜底）：同 A-甲，仅「为 200K 档；1M 档单价」与「为 1M 档；200K 档单价」两处对调。
- 替换前文（EN）：`**Context-window pricing note (landed by P87)**: the 1M tier of \`gpt-6.2-sol\` is 400 / 1500 (input 2x and output 1.5x the table above); \`gemini-4-carbon\` and \`gemini-4-hydrogen\` have a 1M tier only, so the prices above are the real prices; \`gemini-4-neon\` has two tiers, 200K / 1M; \`gpt-5.3-codex\` is measured at 400K. Cache read / write prices (Credits / 1M tokens) are Neon 15 / 0, Carbon 10 / 0, Hydrogen 7.5 / 0 respectively.`
- 拟文（EN，口径 B）：`**Context-window pricing note (landed by P87; Neon band attribution added by P98)**: the 1M tier of \`gpt-6.2-sol\` is 400 / 1500 (input 2x and output 1.5x the table above); \`gemini-4-carbon\` and \`gemini-4-hydrogen\` have a 1M tier only, so the prices above are the real prices; \`gemini-4-neon\` has two tiers, 200K / 1M - the table's 300 / 1500 belongs to the <band>, and the other band costs <Input / Output, read from the user's operator UI on 2026-11-28>; \`gpt-5.3-codex\` is measured at 400K. Cache read / write prices (Credits / 1M tokens) are Neon 15 / 0, Carbon 10 / 0, Hydrogen 7.5 / 0 respectively.`
- 拟文（EN，口径 A-甲 / A-乙）：同 B 的槽位替换为 `- the table's 300 / 1500 belongs to the 200K tier; the 1M-tier price is unavailable on the pricing page ⇒ stop and escalate per this section's fallback sentence (no estimation or filling in);`（A-乙 对调 200K / 1M 两词）。
- 主表 #L852 Neon 行不动（归属在注记内标注，避免加列破坏六形状；该行已是「一表一行」形态）。
- 理由：覆盖 O-4——两档单价齐备（或明示兜底）+ 档位归属明确；Sol 类比（1M = 2x / 1.5x）与 Gemini 族先例（Carbon / Hydrogen「仅 1M 档、上表即实价」）方向相反，**归属不可推定**（详见 §九）。

**【替换 9｜CN §7.5 #L856｜EN §7.5 #L857】备选额度口径（首句精简）**

- 替换前文（CN）：`**备选额度口径（P87 落定）**：备选档（常规档的 Sol 备选）⑦ 额度**与深度档同口径**——均按 §7.5 的切片档行取值（L 档 / 治理片 / 高风险轮次 = 1 终审 + 3 复审）；轻量档（\`gemini-4-hydrogen\`）同常规档。模型档只决定模型与单价，不改变 §7.5 按切片档分层的次数上限。`
- 拟文（CN）：`**备选额度口径（P87 落定；P98 精简）**：所有 ⑦ 选型矩阵模型档的额度**均按 §7.5 的切片档行取值**——模型档只决定模型与单价，不改变 §7.5 按切片档分层的次数上限。`
- 替换前文（EN）：`**Alternative-tier budget (landed by P87)**: the alternative tier (the Sol alternative of the regular tier) takes the **same budget caliber as the deep tier** - both read their value from the §7.5 slice-tier row (L tier / governance / high risk = 1 final review + 3 re-reviews); the light tier (\`gemini-4-hydrogen\`) matches the regular tier. The model tier decides only the model and the price and does not change the slice-tier call caps of §7.5.`
- 拟文（EN）：`**Alternative-tier budget (landed by P87; streamlined by P98)**: the quota of every §5.2 selection-matrix model tier **is read from the §7.5 slice-tier row** - the model tier decides only the model and the price and does not change the slice-tier call caps of §7.5.`
- 理由：采纳 Hydrogen #4 建议原样——删除「与深度档同口径」与「轻量档同常规档」两句（它们已被末句收敛、且首句「1 终审 + 3 复审」括注是误读源头），单句口径与切片档表零张力。P97 自由度清单未登记该口径偏离的补记见 §九（冻结件不修）。

---

## ④ 升版判定

**候选结论：v1.83 = MINOR**（唯一候选；若用户选择加大修订范围另议）。

理由：

1. **MINOR 判例**：v1.33（「新增判据与说明句，未改角色集与流水线结构」= MINOR）与 v1.41（「7 项均为规则条款 / 说明句 / 台账翻转，未改角色集与流水线结构」= MINOR）——本片修订 §5.2 选档判据措辞、§7.5 计价注记、§8.2 授权句，属**规则条款与说明句级修订**；不新增判据（§6.3 与 `tools/gates/**` 零改动）、不改角色集（§2.2 零改动）、不改流水线结构（§5.1 零改动）⇒ 不满足 MAJOR 判据。
2. **PATCH 不适格**：非台账补登 / 措辞纠错单行（涉及三节规则条款与一处术语定义增补）。
3. **版本包**：头部 `版本：v1.82` → `v1.83`（G1-b 与 §0.3 最新行一致）；§0.3 表尾新行（拟文见下）；附录 C 新块 C.0bf（编号承接 C.0be / P96）；§12.4 **不新增 OB**（本片不构成临时授权 / 突破字面事件，见 ⑥）。

§0.3 新行拟文（CN）：

`| v1.83 | 2026-11-28 | **切片 P98 \`DIRECTIVE-REVIEW-TIER-ALIGN\` 落地（治理片 [SLICE]，M 档；用户 2026-11-28 点名起跑）**：① §5.2 深度档档位列拆明「Sol 主读取 Max（须同轮授权）；Neon 交叉复核取 Extra High」并与 §8.2 逐字对齐，常规档计口改「其余四档」、轻量档术语定义、大上下文文档档逻辑连接词与优先次序、平行备选档授权与触发门槛、适用域排除句「触发式」改「必配」，新增「档位取档次序」句；② §8.2 取档授权句对齐（Neon = Extra High 无需授权；Luna 唯一档 Max 取档须同轮授权）；③ §7.5 补 Neon 两档单价与档位归属、备选额度口径首句精简；④ 观察项 O-2..O-4 / O-5 台账处置翻「已处置（P98）」；⑤ 附录 C 增 C.0bf。**不改 §2.2 / §2.4 / G5-b 闭集与白名单 / §6.3 判据本体；不动 \`tools/gates/**\` / \`internal/**\` / \`ci.yml\` / CONTRACTS / schema / fixtures / workflows；不新增 ADR。** | \`P98 DIRECTIVE-REVIEW-TIER-ALIGN\`（本片） |`

（EN 镜像按现文 §0.3 行风格翻译，与 CN 同位。）

附录 C.0bf 拟文（指针块，同 C.0be 风格；正文以准则三节为准）：

`### C.0bf P98 DIRECTIVE-REVIEW-TIER-ALIGN（v1.83）` + 块体表格（切片 / 触发 / 目标 / 交付物 / 范围 / 审查 / 版本 / 附注），其中范围行写「只改 §5.2 / §7.5 / §8.2 相关行 + 版本包 + 台账；不改 §2.2 / §2.4 / G5-b / §6.3 / \`tools/gates/**\` / \`internal/**\` / \`ci.yml\`；不新增 ADR」，附注行写「本块只做指针；条款正文以准则 §5.2 / §7.5 / §8.2 为准，本块不复述」。

---

## ⑤ OB 处置口径（停点 1-3）

**候选结论：O-2..O-4 / O-5 统一翻「已处置（P98）」**；O-5 #5 只登记「已知、不追改」（不再单列为待处置）。翻动位置与拟文：

1. **台账 `REMAINING_SLICES.md` P97 段「观察项处置方向」行（#L2606）**：
   - 替换前文：`观察项处置方向（用户 2026-11-28 裁定）：O-2..O-4 尽快处理——并入下一治理片顺带修订，或另立独立小修订片（修订对象 = §5.2 深度档档位列与取档授权句对齐 §8.2、§7.5 的 Neon 单价档位归属、矩阵可判性措辞）；O-1 仅记录不追改。`
   - 拟文：`观察项处置方向（用户 2026-11-28 裁定）：O-2..O-4 **已处置（P98）**——由 P98 \`DIRECTIVE-REVIEW-TIER-ALIGN\` 修订落地（v1.83：§5.2 深度档档位列与取档授权句对齐 §8.2、§7.5 的 Neon 单价档位归属、矩阵可判性措辞）；O-5 **已处置（P98）**（#4 备选额度口径精简；#5 numstat 陈旧只登记「已知、不追改」）；O-1 仅记录不追改。`
2. **`REMAINING_SLICES_EN.md` P97 段同行（#L2593）**：
   - 拟文：`Observation disposition direction (user ruling, 2026-11-28): O-2..O-4 **handled (P98)** - landed by the P98 \`DIRECTIVE-REVIEW-TIER-ALIGN\` revision (v1.83: aligning the section 5.2 deep-tier tier column and tier-authorization sentence with section 8.2, the Neon price band attribution in section 7.5, and the matrix decidability wording); O-5 **handled (P98)** (#4 the alternative-quota wording streamlined; #5 the stale numstat is registered only as "known, not retro-fixed"); O-1 is recorded only and not retro-fixed.`
3. **P98 段状态行（闭环时翻转，拟文先行）**：CN #L2615 与 EN #L2602 的状态行在 ⑩ 前翻 `✅ 闭环`（完成记录载明 O-2..O-4 / O-5 已处置）；停点 1 阶段不动（P98 仍 🔄 执行中）。
4. **CHECKLIST_v81.md**：**冻结不改**（O-2..O-4 / O-5 处置翻转由 v82 快照吸收；本片不触 v81 快照）。
5. **`docs/validation/hydrogen-control-pilot{,_EN}.md`**：冻结不追改（该报告 §4 / §9 的观察项原文保持）。
6. O-5 #5 的「已知、不追改」登记句同时进 P98 报告（异常与偏离章）与台账 P97 段处置行（如上拟文），不进 §12.4（O 系列属快照 / 台账观察项轨道，非 OB 表）。

---

## ⑥ ADR 判定（停点 1-4）

**候选结论：不新增 ADR**（唯一候选）。

理由：ADR 登记表承载的是**决策记录与临时授权 / 字面突破独立留痕**（ADR-024 / ADR-037 / ADR-044–048 各条性质列）。本片：① 无临时授权事件（模型调用全在 §2.2 白名单与 §7.5 额度内）；② 无 §2.2 / §2.4 字面突破；③ 修订内容 = 消除 P87 落地物（ADR-048 决策②产物）内部的 High 级冲突 + 补全可判性，观察项 O-2..O-4 / O-5 的处置轨道已由快照（v81 §五）与台账承载；④ 对照 ADR 各条性质列，「档位对齐修订」不属「含规则语义变更、系用户直接裁定」类（本片系用户点名的常规治理片，规则语义为**冲突消解与补全**而非新增规则）。⇒ 无 ADR 事件。

**ADR-048 是否加修订注记 —— 两候选**：

- **候选甲（加注记；本方案推荐）**：在 ADR-048 行尾追加一行注记（双语，`ADR_REGISTER{,_EN}.md`）——CN：`**P98 修订注记（2026-11-28）**：本 ADR 决策②所落 §5.2 矩阵的深度档档位列与 §8.2 取档授权句已由 P98 \`DIRECTIVE-REVIEW-TIER-ALIGN\`（v1.83）对齐修订（Neon 交叉复核 = Extra High 无需授权；Luna 唯一档 Max 取档须同轮授权；决策①③与 §2.2 三行均未变）。`
  - 理由：① ADR-048 决策②的落地物（矩阵）已被本片修正，不加注记会使后续读者以 P87 产物形态误读该 ADR；② 先例 = ADR-024 加修订注记（v1.35 因 §1.7.4 ⑤ 修订其有效范围措辞）；③ 成本 = 单行双语。
- **候选乙（不加注记）**：理由：ADR-048 记录的三个决策项（扩列三枚 / 四档矩阵 / §5.4）**均未改变**；本片修正的是落地物内部冲突（P87 笔误类），不改变任何决策内容；ADR 台账按 §11.1 不追溯改写，加注记引入「决策 vs 落地物」双源维护。
- 推荐：**候选甲**——High 级冲突修正需可追溯，且「决策未变、产物修正」恰是修订注记的适用形态（ADR-024 先例）；风险仅限双语镜像义务（G3）。

---

## ⑦ ⑥ 跳过复核（起跑消息已裁定，本环复核确认）

**结论：⑥ 跳过成立。**

依据：

1. **文件面（§7.2.1 ①②）**：本片变更集仅 `.md`（DELIVERY_DIRECTIVE{,_EN}.md、REMAINING_SLICES{,_EN}.md、报告对、证据包新件）；不触 `tools/gates/**`（gate.js / selftest.js / lib/criteria/**）、`internal/**`、`ci.yml`、CONTRACTS / schema / fixtures / workflows；无任何代码 / 配置语义变更 ⇒ ①② 成立。
2. **§7.2.1 ③（不涉状态机 / 门禁判据 / 证据模型 / 角色授权契约语义）**：逐面论证——状态机（无代码改动 ✓）；门禁判据（§6.3 与 `tools/gates/**` 零改动，G 判据判定逻辑不变 ✓）；证据模型（不动 G8 / 冻结包 / 证据契约 ✓）；角色 / 授权契约（§2.2 角色-载体-模型绑定零改动、G5-b 白名单闭集零改动、**不新增任何授权义务**——Luna 取 Max 须同轮授权已由 §8.2 现文末句「任何 Max 取档仍须同轮授权」覆盖，本片仅消歧，授权要件集合不变；「Sol 主读取 Max 须同轮授权」现文常规档行已有同串）⇒ ③ 成立。
3. **先例**：P97（⑥ 跳过，经 ⑦ 确认）与 P72 同类；本片不触 P87 曾触的 §6.3 G5-b / `tools/gates` 实现面，比 P87 窄。
4. **义务**：§7.2.1「跳过必须由 ⑦ 或用户确认」——用户已于起跑消息裁定，本环复核一致；⑦ 终审轮对「③ 的边界论证」追加复核（写入 P98 报告）。

## ⑧ selftest 影响面

- 判据实现与 §6.3 判据文本**零改动** ⇒ `SELFTEST` 806/806 不变（不新增 / 不删减断言）。
- selftest 夹具文本（`tools/gates/selftest.js` 内嵌判据文本样例）：本片不改判据文本，无需同步夹具 ⇒ **预计不动**。
- 验收基线：三作用域 `TOTAL_FAIL=0` + CI 双腿 success（②/⑨ 环节执行，本环不跑）。

## ⑨ 风险与开放项

1. **Neon 另一档单价（O-4 核心开放项，停点 1 交用户裁定）**——两条候选口径：
   - **候选 A（标注归属 + 兜底句）**：主表 300 / 1500 标注归属档，另一档单价位按 §7.5 兜底句写「计价页不可得 ⇒ 停下升级用户裁决（不得估算或填写）」。归属子候选：A-甲 = 上表为 200K 档（依 Sol 先例：主表写最常用档、1M 档另注）；A-乙 = 上表为 1M 档（依 Gemini 族先例：Carbon / Hydrogen「仅 1M 档、上表即实价」）。两先例方向相反，**归属不可推定**。
   - **候选 B（停点 1 索取 operator UI 实时计价，推荐）**：请用户在停点 1 提供 operator UI 实时计价页的两项读数（主表 300 / 1500 属 200K 还是 1M 档 + 另一档 Input / Output 单价），② 落定时填入替换 8 槽位并双语同步。
   - **推荐 B**：验收证据②要求「两档单价齐备且档位归属明确」，A 口径下另一档仍是「明示不可得」而非齐备，O-4 只算部分处置；B 一次到位、零估算。若用户侧 UI 读数不可得，退回 A（子候选随归属标注）。
2. **O-5 #4 自由度清单补记**：P97 证据包 `freedom-list.md` 未登记「备选额度与起跑消息『与深度档同 1+3』口径的偏离」——P97 冻结件不修，本片在其报告与台账处置行如实登记该补记（随 ⑤ 拟文）。
3. **G 判据执行项（② 落定时门禁机械验证，本环不自证）**：G4b 新增 CJK 差集为空（拟文字全部来自现文语料，已逐字比对；② 前跑 `--check=G4b` 复核）；G6-4 无「数字＋行 / 文件 / 共」形态；G2 无占位（槽位在 ② 前填实）；G3(a.2) 六形状与 (b) 关键字段双语对齐（② 前跑 `shape-six` 与 numstat 对称核对）。
4. **§7.2.1 ③ 边界论证的 ⑦ 复核**：本片对「角色 / 授权契约」的消歧性论证（授权要件集合不变）最终由 ⑦ 终审确认，如 ⑦ 认为 ③ 不成立则补跑 ⑥（按 §7.2 失败处置执行，不自动豁免）。
5. **术语锚定风险**：轻量档新定义锚定「⑧a / ⑧b 稿或报告终稿」，若用户认为治理片「⑧ 稿」口径需要更宽（含 §0.3 行 / 台账行），在停点 1 裁定中修订该括注即可（本定义不改判据、只增可判性）。

---

## 停点 1 四项候选汇总

1. **修订文本三件套**：采纳本方案 §三 拟文（替换 1–9 + 插入 1；Neon 口径待项 1 裁定；槽位在 ② 前填实）。
2. **升版判定**：v1.83 = MINOR（头部 + §0.3 新行 + 附录 C.0bf；§12.4 不新增 OB）。
3. **OB 处置口径**：O-2..O-4 / O-5 翻「已处置（P98）」；翻动位置 = REMAINING_SLICES{,_EN}.md P97 段观察项处置方向行（拟文见 §五）+ P98 段状态（闭环时）；CHECKLIST_v81 与 P97 报告对冻结不改；O-5 #5 只登记「已知、不追改」。
4. **ADR 判定**：不新增 ADR；ADR-048 加修订注记——推荐候选甲（加），候选乙（不加）备选。

---

## 回执

**DESIGN: DONE**

1. **三件套要点**：§5.2 深度档档位列拆明「Sol 主读取 Max（须同轮授权）；Neon 交叉复核取 Extra High」+ 五子项可判性补全（轻量档定义 / 大上下文连接词与次序 / 其余四档 / Sol-Luna 两级触发门槛 / 深水轮界线）+ 次序句；§8.2 取档授权句三成员对齐（Neon = Extra High 免授权、Luna Max 须授权）；§7.5 Neon 两档归属 + 备选额度口径单句化（替换 1–9 + 插入 1，CN/EN 同构）。
2. **升版**：v1.83 = MINOR（头部 / §0.3 新行 / 附录 C.0bf；§12.4 不新增 OB）。
3. **OB**：O-2..O-4 / O-5 翻「已处置（P98）」；翻动 = REMAINING_SLICES 双语 P97 处置行 + P98 状态；CHECKLIST_v81 与 P97 报告对冻结；O-5 #5 只登记不追改。
4. **ADR**：不新增；ADR-048 加修订注记（推荐甲；乙备选）。
5. **开放项**：Neon 另一档单价推荐口径 B（停点 1 索取 operator UI 实时计价后填实；A-甲 / A-乙 兜底）；§7.2.1 ③ 边界论证交 ⑦ 终审复核；⑥ 跳过成立、SELFTEST 806/806 不变（判据零改动）。
