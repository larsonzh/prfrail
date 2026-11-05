# ⑦ 独立终审原文（r2，整改后复审）— `G4B-CORPUS-DECOUPLE`

性质：`independent-reviewer` + `GPT-5.3-Codex (copilot)` 的**最终消息逐字转录**。整改后复审第 1 轮（收窄输入：只核 r1 的三项发现闭合）。

---

读数来源：读取主控落盘的 [docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md](docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md)、[docs/validation/evidence/G4B-CORPUS-DECOUPLE/SHA256SUMS.txt](docs/validation/evidence/G4B-CORPUS-DECOUPLE/SHA256SUMS.txt)、[docs/validation/evidence/G4B-CORPUS-DECOUPLE](docs/validation/evidence/G4B-CORPUS-DECOUPLE)、[tools/gates/lib/ctx.js](tools/gates/lib/ctx.js)、[tools/gates/selftest.js](tools/gates/selftest.js)、[docs/validation/g4b-corpus-decouple.md](docs/validation/g4b-corpus-decouple.md)、[tmp/G4B-CORPUS-DECOUPLE/review-r5/input-manifest.txt](tmp/G4B-CORPUS-DECOUPLE/review-r5/input-manifest.txt)、[tmp/G4B-CORPUS-DECOUPLE/timings.txt](tmp/G4B-CORPUS-DECOUPLE/timings.txt)、[tmp/G4B-CORPUS-DECOUPLE/after-selftest-r2.txt](tmp/G4B-CORPUS-DECOUPLE/after-selftest-r2.txt)。

三方一致核验结果：MANIFEST 条目数 7（[MANIFEST.md:15-21](docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md#L15-L21)），SHA 条目数 8（[SHA256SUMS.txt](docs/validation/evidence/G4B-CORPUS-DECOUPLE/SHA256SUMS.txt)），包内载荷数 8（不含两件元数据）；差异为 MANIFEST 缺少 [docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/05-pre-review-r2.md](docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/05-pre-review-r2.md)（SHA 见 [SHA256SUMS.txt:4](docs/validation/evidence/G4B-CORPUS-DECOUPLE/SHA256SUMS.txt#L4)）。

HIGH | [docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md:7](docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md#L7) | “条目集三方相等（…）” | 在归档表补入 04-reviews/05-pre-review-r2.md（位于归档清单区）并同步重算 SHA，再复核三方集合为 8=8=8
LOW | [tmp/G4B-CORPUS-DECOUPLE/review-r5/input-manifest.txt:11](tmp/G4B-CORPUS-DECOUPLE/review-r5/input-manifest.txt#L11) | “tmp/G4B-CORPUS-DECOUPLE/timings.txt|bytes=2050|lines=31” | 重写 input-manifest 全部 lines 字段，按实际末行号去除统一 +1 偏差（例如 [timings.txt:30](tmp/G4B-CORPUS-DECOUPLE/timings.txt#L30)、[ctx.js:772](tools/gates/lib/ctx.js#L772)、[selftest.js:6729](tools/gates/selftest.js#L6729)）

H-1 → 未闭合。依据：三方计数 7/8/8，且缺失项为 04-reviews/05-pre-review-r2.md。
M-1 → 闭合。依据：T118 明确要求 index 走 ls-tree HEAD（[selftest.js:6684](tools/gates/selftest.js#L6684)）且禁止任何 ls-files 调用（[selftest.js:6685](tools/gates/selftest.js#L6685)）；同时假桩把 ls-files 设为错误集合（[selftest.js:6671](tools/gates/selftest.js#L6671)），因此“加 ls-files 特判”必红；实现端 index/tree 枚举统一到 ls-tree（[ctx.js:547](tools/gates/lib/ctx.js#L547)、[ctx.js:548](tools/gates/lib/ctx.js#L548)）。
L-1 → 闭合。依据：§2 首句已改为“本围栏是步骤 0 的快照，⑧ 阶段按最终变更集重排并逐件对齐”（[docs/validation/g4b-corpus-decouple.md:21](docs/validation/g4b-corpus-decouple.md#L21)）。

独立复算补充：SELFTEST 读数一致，timings 为 703/703（[timings.txt:30](tmp/G4B-CORPUS-DECOUPLE/timings.txt#L30)），after-selftest-r2 为 703/703（[after-selftest-r2.txt:708](tmp/G4B-CORPUS-DECOUPLE/after-selftest-r2.txt#L708)）。

RE-REVIEW: FINDINGS
