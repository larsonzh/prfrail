# 07-final-r1 — P67 GATES-DEBT-CLEANUP（⑦ 终审第 1 轮，逐字归档）

> 归档纪律（§13.1.13 / §11.3）：⑤⑥⑦ 产出逐字原文归档；主控不改写。以下为 ⑦ 终审 r1 回执原文。

数据来源（只读）：读取主控落盘的 [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/00-briefing.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/00-briefing.md)、[tmp/P67-GATES-DEBT-CLEANUP/review-input/final-r1.diff](tmp/P67-GATES-DEBT-CLEANUP/review-input/final-r1.diff)、[docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md)、[docs/validation/gates-debt-cleanup.md](docs/validation/gates-debt-cleanup.md) 及对应镜像与证据件。

**发现（按严重度）**
1. Medium | file:line = [docs/validation/gates-debt-cleanup.md#L10](docs/validation/gates-debt-cleanup.md#L10), [docs/validation/gates-debt-cleanup.md#L23](docs/validation/gates-debt-cleanup.md#L23), [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/00-briefing.md#L262](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/00-briefing.md#L262), [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/00-briefing.md#L264](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/00-briefing.md#L264), [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/00-briefing.md#L265](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/00-briefing.md#L265)
逐字引用：
“交付物：§1 执行摘要、§2 四案对比与最终落点、§3 #1 判据本体改动、§4 #2 门禁坑正式化…§5 #3 §12.4 清账、§6 #4 ⑥ 提示词加固、§7 已知边界、§8 明确未执行、”
“## 1. 执行摘要”
“## 7. 已知边界”
“sixstage 围栏于 ⑧a 装配时逐字转录（§11.3），本行届时替换为围栏外注记。”
理由：按切片定义，报告交付面要求包含 §1~§8 与 sixstage 转录；当前报告仍是进行中骨架，仅见 §1 与 §7，关键章节与 sixstage 尚未落盘，审计闭环不完整。
最小修复：在进入 ⑩ 前补齐 [docs/validation/gates-debt-cleanup.md](docs/validation/gates-debt-cleanup.md) 与 [docs/validation/gates-debt-cleanup_EN.md](docs/validation/gates-debt-cleanup_EN.md) 的 §2~§6、§8，并把 [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/06-scan-r1.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/06-scan-r1.md) 的 A–E 转录为 sixstage 围栏；随后重算并回填 [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/MANIFEST.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/MANIFEST.md) 与 [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/SHA256SUMS.txt](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/SHA256SUMS.txt)。

2. Low | file:line = [tools/gates/selftest.js#L7182](tools/gates/selftest.js#L7182), [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md#L25](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md#L25)
逐字引用：
“const cap = G4B.internals.CACHE_CAP;”
“**判定**：LRU（案 B）在边界内可落地。**建议 N=100**（覆盖现状 96 条、留 4 条余量；”
理由：新增自测在 T127 使用实现内常量作为期望值，能验证“有上限”，但不能防止“上限值漂移”导致的假通过（例如常量被改成非 100，测试仍绿）。
最小修复：在 selftest 增加一条显式断言，固定检查 CACHE_CAP 等于 100；若未来要改值，则同步更新切片定义/准则文案并修改该断言。

**已核查无问题（本轮）**
- 安全性：未发现会改变 G4b verdict 语义的实质风险。命中仍经有效值约束后重算 verdict，见 [tools/gates/lib/criteria/g4b.js#L163](tools/gates/lib/criteria/g4b.js#L163), [tools/gates/lib/criteria/g4b.js#L307](tools/gates/lib/criteria/g4b.js#L307)；touch 与淘汰仅作用于缓存元数据，见 [tools/gates/lib/criteria/g4b.js#L243](tools/gates/lib/criteria/g4b.js#L243), [tools/gates/lib/criteria/g4b.js#L265](tools/gates/lib/criteria/g4b.js#L265)。
- 架构一致性：缓存仍是优化层而非第二真值源；异常路径维持 fail-open 取向。
- 完整性子项：A–G 留痕（A~F + consistency）、MANIFEST/SHA、自身中英镜像在当前输入包中可对齐。

FINAL REVIEW: FINDINGS
