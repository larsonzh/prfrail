# ⑤ 预审 r1 原文（deep-reasoner + DeepSeek V4 Pro）

## ⑤ 预审结论(P71 S1-SCOPE-ANNOTATE)

**四项审查 + 三项额外核实,逐条:**

INFO | docs/PRODUCT_REQUIREMENTS.md:69/71/81/83 | REQ-007 阶段列「S1/S2(黑箱候选子项推迟 S2)」、REQ-009「S1 不要求,推迟至 S2(基础终端保留)」、REQ-019 仍「S1」未动、REQ-021 仍「S1/S2」未动 | 与 P68 b1(`remaining-rulings.md:92/:105`、`r07-refined.md:42/43/63`)逐条一致:REQ-019 未被越界改动,REQ-007 未被整条标「不要求」;无需修复。

INFO | docs/PRODUCT_REQUIREMENTS{,_EN}.md(REQ 表) | 「可验收要求」列(第 4 列)在 diff 中零改动;REQ 条目 28 条(REQ-001–022 功能 + REQ-023–028 非功能)编号连续;本文件无 MVA 文本、diff 未引入,`s1-exit` 新增段仅陈述「不批正式定义」 | 不变量全部成立;无需修复。

INFO | docs/validation/s1-exit{,_EN}.md:52–60 | 门表 9 行 diff 为零:BLOCKED 2(核心任务链/产品所有者批准)、PASS 4、PARTIAL 3,与 ③ 读数一致;AT 表 PASS 22 + 门表 PASS 4 = 26 与 ③ 一致 | 门表 zero-diff 成立;无需修复。

INFO | docs/OPERATIONS.md:37 | 三条边界「visible 不是 S1 交付承诺」「不构成权威证据」「依赖 VS Code+Copilot Chat+桥扩展+风险确认+人工监督;真实 CLI 端到端仍依赖 T027-B」 | 与 ADR-012(`docs/ADR_REGISTER.md:22`「visible 只证明面板投递」「拒绝把 transcript/完成声明/退出 0 当权威事实」「真实 CLI 端到端依赖 T027-B」)逐项一致;无需修复。

INFO | docs/S1_SUPPORT_MATRIX.md(visible 行) | visible 行仍为「产品侧已实现(降级保证)、真实 CLI 端到端依赖 T027-B、尚不可运行发布」,diff 未动 | 仍为「降级保证」而非「正式支持」;无需修复。

LOW | docs/S1_SUPPORT_MATRIX.md:33(注记) / :9(Windows 行) | 注记自述「平台清单、支持层级、正式支持范围与 visible 降级保证本身均不变」,而 Windows 行约束列新增「完整一体化 TUI 移出 S1、由 S2 承接」——「正式支持范围不变」与「S1 功能承诺面缩小」存在措辞张力 | 注记改为「平台清单与支持层级不变;TUI 移出系 P68 已裁定的功能阶段归属同步」(状态列「正式支持候选」未变,故非裁定冲突)。

LOW | docs/validation/s1-exit.md:62(P71 适用性说明) | 「完整 TUI 移出不构成门级判定变更」是断言性声明,缺论证:AT-12/16/22 断言全部基于基础终端/`interactions tui`、不依赖完整一体化 TUI,故既有 PASS 不被削弱 | 在说明段补一句「AT-12/16/22 断言不依赖完整 TUI,既有 PASS 证据不受影响」。

LOW | docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/01-arch/design-r2.md(§7 X-1) | 判据 `grep -c '^\| REQ-'` 在 BRE 下 `\|` 为 alternation,`^` 命中每行 ⇒ 字面执行返回文件总行数而非 28,判据写法错误 | 改为 `grep -c '^| REQ-'`。

LOW | design-r2 §7 X-2 | 「正式支持」全文件 grep 会额外命中 P71 注记行的「正式支持范围」字样,「仅命中支持候选行」字面不满足 | 限定为表格行内断言:「Windows 行为 '正式支持候选' 的唯一点,visible 行不得含 '正式支持'」。

LOW | design-r2 §7 X-3/X-4 | 行号锚点(`:52`/`:60`/`:81`)当前正确,但文件前部任一行增删即漂移,机械判定脆弱 | 改为内容断言(门表行名+状态词、REQ-019 行第三列)而非行号。

INFO | docs/validation/remaining-rulings.md:107 / p57-closeout-assessment.md:93 / s1-stream-progress.md:83 | 三条加法式指针均不改原行、显式写「原行不改(§11.1)」;§11.1 字面即「历史版本不重写,只加指针或更正注记」;MATRIX 已有 P55/P57 同类注记先例 | 与「历史不改写」不冲突,加法式可接受;无需修复。

INFO | docs/PRODUCT_REQUIREMENTS.md:82(REQ-020) | P68 b1 ③「S2 收窄(Q3 语言冻结)」对应的「S2 承诺收窄登记(REQ-020)」未见于本片 diff、亦无落地指针(阶段列仍 S0-S2 无冻结注记);本片切片定义未列 REQ-020,属范围外 | 建议报告「已知边界」登记 REQ-020 收窄登记的承接状态,交用户裁定。

**额外三项核实结论:**
- (a) **不成立(不报 Medium+)**——追加属 P68 r07-refined §3 明示授权的「S1 范围基线(支持层级/平台行)若受影响则同步」,Windows 行状态列「正式支持候选」未变、头部「状态:CANDIDATE」无冲突;仅注记措辞歧义(见 LOW 1)。
- (b) **不成立**——门表 9 行 diff 为零、状态词与判据未变;适用性说明是描述性声明,未引入豁免或新判据(见 LOW 2 补论证)。
- (c) **不成立**——加法式指针即 §11.1 字面允许的「只加指针」,且原行未改(见 INFO 8)。

**可证伪性说明:** 删除「REQ-019 | R19 | S1」的阶段值、把 visible 行改为含「正式支持」、改动三条指针上方任一原行、或把门表任一行状态词从 BLOCKED/PASS/PARTIAL 改值,均会使对应 grep/diff 变红;X-1 修复后 `grep -c '^| REQ-'` 必须等于 28,X-5 以「第三列 diff 为空」可直接验证。

PRE-REVIEW: PASS WITH FIXES
