# ⑦ 独立终审原文（r1）— `G4B-CORPUS-DECOUPLE`

性质：`independent-reviewer` + `GPT-5.3-Codex (copilot)` 的**最终消息逐字转录**（首轮，未附任何 ⑤ / ⑥ 清单）。时点 2026-11-05。

---

我已完成只读证据交叉：实现、自测段落、报告对、证据包清单和实际目录都核对完了。接下来给你正式终审结论，先列严重度发现，再给四段审查回答和最终判定。

读数来源：读取主控落盘的 [tmp/G4B-CORPUS-DECOUPLE/review-r3/diff.patch](tmp/G4B-CORPUS-DECOUPLE/review-r3/diff.patch)、[tools/gates/lib/ctx.js](tools/gates/lib/ctx.js)、[tools/gates/lib/criteria/g4b.js](tools/gates/lib/criteria/g4b.js)、[tools/gates/selftest.js](tools/gates/selftest.js)、[docs/validation/g4b-corpus-decouple.md](docs/validation/g4b-corpus-decouple.md)、[docs/validation/g4b-corpus-decouple_EN.md](docs/validation/g4b-corpus-decouple_EN.md)、[docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md](docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md)、[docs/validation/evidence/G4B-CORPUS-DECOUPLE/SHA256SUMS.txt](docs/validation/evidence/G4B-CORPUS-DECOUPLE/SHA256SUMS.txt)、[docs/validation/evidence/G4B-CORPUS-DECOUPLE/freedom-list.md](docs/validation/evidence/G4B-CORPUS-DECOUPLE/freedom-list.md)、[docs/validation/evidence/G4B-CORPUS-DECOUPLE/input-list.md](docs/validation/evidence/G4B-CORPUS-DECOUPLE/input-list.md)、[tmp/G4B-CORPUS-DECOUPLE/timings.txt](tmp/G4B-CORPUS-DECOUPLE/timings.txt)、[tmp/G4B-CORPUS-DECOUPLE/ab-probe-adr.txt](tmp/G4B-CORPUS-DECOUPLE/ab-probe-adr.txt)、[tmp/G4B-CORPUS-DECOUPLE/probe-staged-delete-g4b.txt](tmp/G4B-CORPUS-DECOUPLE/probe-staged-delete-g4b.txt)。

High | [docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md:7](docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md#L7) | 被违反句子：“条目集三方相等（MANIFEST 归档条目集合 ≡ SHA256SUMS 条目集合 ≡ 包内载荷文件集合）”；事实：读取主控落盘的 [docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/05-pre-review-r1.md](docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/05-pre-review-r1.md) 与 [docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/06-independent-scan.md](docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/06-independent-scan.md) 已存在，但 [MANIFEST.md:15](docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md#L15)-[17](docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md#L17) 与 [SHA256SUMS.txt:1](docs/validation/evidence/G4B-CORPUS-DECOUPLE/SHA256SUMS.txt#L1)-[3](docs/validation/evidence/G4B-CORPUS-DECOUPLE/SHA256SUMS.txt#L3) 未登记 | 最小修复：将现有 04-reviews 载荷补入 MANIFEST 与 SHA256SUMS，或移出证据包目录并在报告显式声明不入包。

Medium | [tools/gates/selftest.js:6554](tools/gates/selftest.js#L6554) | 被违反事实：T114-T117 仅覆盖 tree/range（见 [tools/gates/selftest.js:6554](tools/gates/selftest.js#L6554) 与 [tools/gates/selftest.js:6642](tools/gates/selftest.js#L6642)）；若把 [tools/gates/lib/ctx.js:547](tools/gates/lib/ctx.js#L547) 改为 index 分支回退到 ls-files 或非 HEAD，现有 T108-T117 可不变红 | 最小修复：新增 index 作用域枚举断言（含 staged delete 冷/热一致性）和 ci 作用域 baseRef 枚举断言。

Low | [docs/validation/g4b-corpus-decouple.md:21](docs/validation/g4b-corpus-decouple.md#L21) | 被违反句子：“以下为本片完整变更集条目”；与当前目录事实不一致（读取主控落盘目录仅见 7 个文件） | 最小修复：改为“计划清单/预期清单”，或补齐后再使用“完整变更集”。

**1. 安全性**
缓存命中路径未见掩盖违规：实现在 [tools/gates/lib/criteria/g4b.js:200](tools/gates/lib/criteria/g4b.js#L200) 先做白名单错误 fail-closed，再在 [tools/gates/lib/criteria/g4b.js:138](tools/gates/lib/criteria/g4b.js#L138) 限定“仅缓存语料字符集”，且异常走 fail-open 重算。
枚举变更后，tree/index 走 [tools/gates/lib/ctx.js:548](tools/gates/lib/ctx.js#L548) 的对象库枚举，语料不会因暂存删除而收缩。
你问的后者我明确回答：语料“大于真值而洗白违规”的路径未见；当前风险方向是 git 异常时返回空列表导致更严（假阳性），不是洗白。

**2. 架构一致性**
与判据文本一致：tree/index 参考 HEAD、range/ci 参考区间基点（见 [tools/gates/lib/ctx.js:547](tools/gates/lib/ctx.js#L547)）。
影响面一致：listAllMarkdown 生产消费点仅 G4b（见 [tools/gates/lib/criteria/g4b.js:236](tools/gates/lib/criteria/g4b.js#L236)）。
kind 段保留自洽：和 [docs/validation/evidence/G4B-CORPUS-DECOUPLE/freedom-list.md:14](docs/validation/evidence/G4B-CORPUS-DECOUPLE/freedom-list.md#L14)、[docs/validation/evidence/G4B-CORPUS-DECOUPLE/freedom-list.md:43](docs/validation/evidence/G4B-CORPUS-DECOUPLE/freedom-list.md#L43) 的 FL-4 一致，属于命名空间隔离，不改变判据语义。

**3. 完整性**
调用点、注释、版本号主链已覆盖：版本升到 v2（见 [tools/gates/lib/criteria/g4b.js:43](tools/gates/lib/criteria/g4b.js#L43)、[tools/gates/lib/criteria/g4b.js:154](tools/gates/lib/criteria/g4b.js#L154)、[tools/gates/lib/criteria/g4b.js:243](tools/gates/lib/criteria/g4b.js#L243)），台账已登记（见 [docs/t027/REMAINING_SLICES.md:89](docs/t027/REMAINING_SLICES.md#L89)）。
但证据包三方一致性不成立（见 High），因此完整性结论为未通过。

**4. 测试反例**
T108-T113 对缓存键三要素和 fail-open 基本敏感；T114/T115 对 tree 对象库枚举与 staged delete 冷热一致敏感；T116 对缓存版本写入敏感；T117 对 range 的 -z 枚举敏感。
具体红例：把 [tools/gates/lib/ctx.js:548](tools/gates/lib/ctx.js#L548) 改回 ls-files，T114/T115 会红。把 [tools/gates/lib/criteria/g4b.js:43](tools/gates/lib/criteria/g4b.js#L43) 改回 v1，T116 会红。
具体不红缺口：只破坏 index 分支（例如改 [tools/gates/lib/ctx.js:547](tools/gates/lib/ctx.js#L547) 令 index 不走 HEAD 对象库），现有 T114-T117 不一定变红。

FINAL REVIEW: FINDINGS
