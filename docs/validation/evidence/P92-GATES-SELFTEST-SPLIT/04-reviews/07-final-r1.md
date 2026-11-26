# P92 ⑦ 终审 r1（原文，逐字）

- 载体：`independent-reviewer` + `GPT-5.3-Codex (copilot)`；时点：2026-11-26；轮次：1。

---

1. FINAL REVIEW: FINDINGS

2. 发现表：SEVERITY | file:line | 问题 | 最小修复  
MEDIUM | [docs/validation/gates-selftest-split.md:85](docs/validation/gates-selftest-split.md#L85) | sixstage 围栏第 85/86/89 行不是对 [docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r3.md:21](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r3.md#L21)、[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r3.md:22](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r3.md#L22)、[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r3.md:25](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r3.md#L25) 的逐字转录；违反 [docs/DELIVERY_DIRECTIVE.md:1030](docs/DELIVERY_DIRECTIVE.md#L1030) 的“块内逐字转录、注记放围栏外”契约 | 将 sixstage 围栏 85/86/89 三行回填为 r3 原文三行，围栏外保留归一说明  
MEDIUM | [docs/validation/gates-selftest-split_EN.md:85](docs/validation/gates-selftest-split_EN.md#L85) | 英文镜像报告同位行存在同样的围栏内非逐字改写（85/86/89） | 同步把 EN 围栏对应三行恢复为 r3 原文

3. VERDICT: has Medium+

4. 逐项裁定 Q1–Q6  
1. Q1 判据结论不变：当前快照下语义等价成立。  
依据是缓存入口仅把固定实参调用收敛到一次：[tools/gates/selftest.js:81](tools/gates/selftest.js#L81)、[tools/gates/selftest.js:82](tools/gates/selftest.js#L82)；源函数仍是 startDir 求根逻辑：[tools/gates/lib/ctx.js:151](tools/gates/lib/ctx.js#L151)-[tools/gates/lib/ctx.js:163](tools/gates/lib/ctx.js#L163)。  
744 断言未见受影响，A/B 都是 PASS 744/744：[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-A.txt:749](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-A.txt#L749)、[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-B1.txt:749](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-B1.txt#L749)、[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-B2.txt:749](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-B2.txt#L749)。

2. Q2 收益论证闸成立，无“收益不足却做”或“默认全做”。  
五项逐项有收益/代价/判定/依据：[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md:19](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md#L19)-[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md:23](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md#L23)。  
关键实测支撑存在：(c) 子进程热点与收益、(b) 夹具读收益不足：[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md:21](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md#L21)、[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md:41](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md#L41)、[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md:47](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md#L47)-[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md:51](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md#L51)。

3. Q3 版本包处理可接受，未越界。  
起跑边界明确给了触判据实现文件时的 MINOR 触发口径：[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md:75](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md#L75)、[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md:154](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md#L154)-[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md:155](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md#L155)。  
① 的 MINOR_bump:no 与主控裁注冲突已显式登记：[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md:8](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md#L8)、[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md:13](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md#L13)、[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/freedom-list.md:10](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/freedom-list.md#L10)。

4. Q4 构成 §11.3 转录契约违规。  
围栏内存在 3 处非逐字改写（85/86/89 对 21/22/25），属于围栏内措辞改写而非纯转录。  
即使围栏外有归一注记：[docs/validation/gates-selftest-split.md:104](docs/validation/gates-selftest-split.md#L104)，也不能替代 [docs/DELIVERY_DIRECTIVE.md:1030](docs/DELIVERY_DIRECTIVE.md#L1030) 的逐字要求。

5. Q5 边界总体满足。  
实现 diff 仅一文件：[tmp/P92-GATES-SELFTEST-SPLIT/review/impl.diff:1](tmp/P92-GATES-SELFTEST-SPLIT/review/impl.diff#L1)；且方案书实现范围同样限定在 selftest：[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md:25](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md#L25)。  
门禁记录也显示禁改面未触及，例如 G5-b 仅 checked 0 file(s)：[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-index-final.txt:18](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-index-final.txt#L18)。

6. Q6 残留风险已被注释与报告材料覆盖，但属于持续约束。  
代码注释已明确“无 chdir 前提”和“未来若依赖 cwd 变化必须直调 discoverRepoRoot”：[tools/gates/selftest.js:76](tools/gates/selftest.js#L76)、[tools/gates/selftest.js:78](tools/gates/selftest.js#L78)。  
文档侧也登记了同类风险：[docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md:51](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md#L51)，报告正文保留了当前等价前提：[docs/validation/gates-selftest-split.md:43](docs/validation/gates-selftest-split.md#L43)。

5. 可证伪性（≤4）  
1. 若把 sixstage 围栏与 r3 做逐行全等比对后不存在任何差异，则本次 Medium 发现失效。  
2. 若你提供的准则文本中存在比 [docs/DELIVERY_DIRECTIVE.md:1030](docs/DELIVERY_DIRECTIVE.md#L1030) 更新且明确允许“围栏内机械归一”的条款，本次 Q4 判定可被推翻。  
3. 若重新核验 A/B 发现任一不再是 PASS 744/744（见 [docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-B1.txt:749](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-B1.txt#L749)），Q1 结论需重判。  
4. 若出现新的实现 diff 触及 selftest 之外代码文件，本次 Q5 边界结论需重判。
