# GATES-EXEMPT criteria-exemption and whitelist-threshold combined governance (P94)

Date: 2026-11-26. Status: **self-check complete**, ⑥ / ⑦ conclusions closed. Slice id: `GATES-EXEMPT` (ledger P94).

## 1. Summary of changes

- (a) Criteria-exemption generalisation: **do** the isomorphic extension - carry the P93 `sixstage` fence-body scope exemption from G6-4 to G2 and G4b; **do not** add a `gate-exempt` comment-marker mechanism (opening a mechanical exemption hole over the controller's own prose would hollow out the criterion, and the fail-closed cost of nested and empty-range shapes is disproportionate).
- (b) The generic quote-exemption registry handed over by P93: **do not** - the verbatim-transcription duty applies only inside the `sixstage` fence, while every other quotation is the controller's own summary and may be reworded.
- (c) The whitelist registration threshold: **do** - clauses only, no admission; `g6-whitelist.txt` keeps zero entries and the three T23 assertions stay byte-identical.
- Version package: MINOR bump to **v1.77** (header line / new §0.3 line / appendix C `C.0bd`), with the CN and `_EN` six shapes equal.

## 2. Execution pipeline

| Step | Role and model | Status | Result |
|---|---|---|---|
| 0 | Controller (DeepSeek V4.1 Flash High) | Done | anchor check; ledger first `⬜ 未立项` then `🔄 执行中` |
| ① | deep-reasoner + DeepSeek V4 Pro | Done | `DESIGN: DONE` (three-item combined assessment + version package) |
| ② | prfrail-implementer + Flash | Two rounds | `IMPLEMENT: DONE` (second round added the G2 / G4b exemption channel and assertions) |
| ③ | prfrail-tester + Flash | Done | `TEST: DONE` (existing-input regression byte-identical + mutation check + three scopes) |
| ④ | Controller | Done | `TOTAL_FAIL=0` in both the tree and the index scope |
| ⑤ | DeepSeek V4 Pro (GLM-5.3 took over after the channel fault) | 1 round | `PRE-REVIEW: PASS WITH FIXES` (1 Low fixed + 3 Info) |
| ⑥ | MAI Code 2 Flash | 1 round | `INDEPENDENT SCAN: PASS` (0 findings) |
| ⑦ | GPT-5.3-Codex standard tier | 3 rounds | r1 `FINAL REVIEW: PASS WITH FIXES`; r2 `RE-REVIEW: FINDINGS`; r3 `RE-REVIEW: PASS` |
| ⑧a | Controller (this report) | Done | report pair + transcript fence + freedom list |
| ⑨ | Controller | Done | three scopes + `gofmt` / `build` / `vet` / `test` all green |
| ⑧b | Controller | Done | ledger flipped to `✅ 闭环` + completion record + main commit |
| ⑧c | Controller | Done | temp files cleared, `tmp` top level back to three entries |
| ⑩ | Controller | Parked at the authorisation point | final receipt |

## 3. Anomalies and fallbacks

- A-1: the ⑤ channel fault - DeepSeek returned `UND_ERR_SOCKET` twice, so the §2.2 standby tier carried the review with GLM-5.3; the reduced independence is recorded as-is and **must not serve as a precedent**.
- A-2: the execution freedom list was written to disk after ⑦; this slice records that fact as-is and does not retro-edit the review record.
- A-3: the 1 Low from ⑦ r1 - the `fence-exempt N` counting caliber of G2; fixed.
- A-4: the 1 Low from ⑦ r1 - T155 was too weak; it was rebuilt as five structured threshold assertions plus a new T156, moving the assertion count from 769 to 770 and the id range to T144–T156, with §0.3 and appendix `C.0bd` kept in step.
- A-5: the two Lows from ⑦ r2 - the fallback path still over-counted without a line-number view, and a comment id was stale; both are CLOSED.
- A-6: `POLICY: DIFF` from batch helper 2 - the whitelist comment header lacked the division-of-labour sentence naming `cjk-newwords.txt`; added, then `PASS 770/770` and `TOTAL_FAIL=0` were re-run.

## 4. Environment and tiers

- Host: local Windows; the Ubuntu VM and the proxy are both disabled; probes 0; no gitee.
- Controller: DeepSeek V4.1 Flash (High tier); ①: DeepSeek V4 Pro.
- ② ③: Flash; ⑤: DeepSeek V4 Pro (GLM-5.3 took over on the §2.2 standby tier after the channel fault).
- ⑥: MAI Code 2 Flash; ⑦: GPT-5.3-Codex standard tier (1 final + 2 re-reviews).
- Batch helper: GLM-5.3-Flash used twice, carried by `deep-reasoner`.

## 5. Review conclusions

- ⑤ pre-review (1 round): `PRE-REVIEW: PASS WITH FIXES`, 1 Low (fixed) and 3 Info; nothing at Medium or above.
- ⑥ independent scan (1 round): `INDEPENDENT SCAN: PASS`, 0 findings; the original text is transcribed verbatim in the fence below.
- ⑦ final review and re-reviews (three rounds, quota exhausted): two Lows in r1, two Lows in r2, both CLOSED in r3 with 0 findings.
- Disclosure outside the fence: the difference from the ⑥ original is exactly one item - a conclusion line identical to Section A was appended at the block end.

```sixstage
## Section A
INDEPENDENT SCAN: PASS
## Section B
| SEVERITY | file:line | violated sentence or fact | minimal fix |
| - | - | - | none |
## Section C
- numstat：变更集恰 13 文件；仅有准则对、台账、证据包、`g2.js`、`g4b.js`、`selftest.js`、`g6-whitelist.txt`。
- 受保护边界核对：`g6.js`、`g7.js`、`ci.yml`、`internal/**`、CONTRACTS、schema、testdata、`cjk-newwords.txt` 均无 staged diff。
- 判定本体：g2谓词未改；g4b差集、白名单与`CACHE_VERSION=v2`未改；豁免均在候选集合剔除后执行。
- `T144-T155` 证据已逐条 PASS；最终 `SELFTEST: PASS 769/769`。
- gate index/tree 均 `TOTAL_FAIL=0`；G4b `newCjk=419 unknown=[]`。
- 三处政策与双语镜像：gate G3(a)/G3(a.2)、G1-b 均 PASS；CN/EN 新增条款形状一致。
- `node --check` 对 g2、g4b、selftest 均无语法错误。
## Section D
- 可证伪性审计：删除 g2/g4b 的豁免剔除、恢复判断谓词、改动受保护文件、使 G4b 出现未知字或让任一 `T144-T155`/门禁读数失败，结论立即转 `FINDINGS`。
- 特别是 T154 会在关闭 g2/g4b 剔除过滤器后把正例转红，证据承重。
## Section E
- 未执行外部 CI/VM 实跑；结论基于完整 staged diff、源码真值、本地 gate/selftest 与双语镜像读数。任何受保护文件变动或 CI 包级证据与本地读数冲突都会推翻本结论。
INDEPENDENT SCAN: PASS
```

## 6. ⑧a self-check

- Cost-table mechanical audit: yes - the call counts in §12 come from the per-stage receipts, and the batch helper ran twice.
- "already disposed" `file:line` scan: none - this slice flips no OB ledger row.
- Old-caliber grep: yes - the tracked occurrences of `fence-exempt` and `sixstageExemptSet` were checked and land in the criteria code, `selftest.js`, the directive pair, the P93 / P94 evidence packages and the P93 report pair.
- Whether every ⑥ / ⑦ Info and Low landed in the report: yes - see A-1 to A-6 in §3.

## 7. A-NN numbering integrity

This report uses A-1 to A-6, each with an entity definition in §3; no dangling reference.

## 8. Review input package

- Inputs = the directive pair diff, `g2.js`, `g4b.js`, `selftest.js`, `g6-whitelist.txt`, the P94 ledger row and the evidence package; sanitised per §7.3.
- The ⑤ input carries the ① design artifact but no other reviewer's conclusion; the ⑥ input carries no ⑤ conclusion; the ⑦ first read carries neither ⑤ nor ⑥ conclusions.
- Blind review: not drawn for this slice (no identity admission / model whitelist / G5-b closed set / ownership / stop condition involved).

## 9. Falsifiability (mechanical checks + mutation)

- T154 in-memory mutation: with the G2 / G4b exclusion filter disabled, the two positive cases T144 and T149 turn red; after restoration the bytes match the baseline sha256 exactly.
- T156 counting caliber: on the same-text collision the count follows line numbers, not occurrence counts, and both the positive and the negative case reproduce.
- All 34 fail-closed probes hand-built by ③ are green; the three scopes read 0.
- The set difference between the baseline-tree and working-tree `selftest` output lines is empty (new assertion lines aside), and deleted lines are 0.

## 10. Gate results

- `--selftest`: `PASS 770/770` (757 baseline + 13 new assertions, T156 included).
- `--all --scope=tree` and `--all --skip-base --scope=index`: both `TOTAL_FAIL=0` (`changed=13`).
- Base gates: `gofmt -l .` empty, and `go build ./...`, `go vet ./...`, `go test ./...` all green.
- `G4b`: `newCjk=419 unknown=[]`.
- The directive pair's six shapes are equal (2624 / 2624 on both sides).

## 11. ⑧c closing-cleanup record

- `tmp/P94-GATES-EXEMPT/` was cleared item by item, with no leftover process script or capture file.
- The `tmp` top level is back to three entries: `.gitkeep` / `gates-cache` (only `g4b` beneath it) / `find-quotes.js`.

## 12. Cost and metering

- Calls: ① 1, ② 2, ③ 1, ⑤ 1 (GLM standby), ⑥ 1, ⑦ 3 (1 final + 2 re-reviews, quota exhausted), batch helper 2, probes 0.
- **R = 6, D = 0**; rework sources: the ② second round adding the channel, the ⑤ 1 Low, the two ⑦ r1 Lows, the two ⑦ r2 Lows, and `POLICY: DIFF` from batch helper 2.

## 13. Explicitly not executed

- No gitee push; the Ubuntu VM and the proxy were not enabled.
- No change to `ci.yml` / `internal/**` / CONTRACTS / schema / fixtures / workflows; no new ADR.
- No whitelist entry registered; the G5-b closed set and the model whitelist are unchanged.

## 14. Precedent-ban check

- One "must not serve as a precedent" clause is cited - the ⑤ model standby tier and the after-the-fact landing of the execution freedom list; neither was used to expand authority, and both are recorded as-is in A-1 and A-2.

## 15. Next steps

- The v78 snapshot should absorb this slice (P94 leaves the queue).
- The next item, the S2 pre-re-evaluation (S-G1 / S-G8), waits for user scheduling.

## 16. Commit and CI

- Main commit: `49f1609` (16 paths; message `feat: extend the sixstage scope exemption to G2 and G4b and add the whitelist threshold clause (P94)`).
- Push: `origin/main` (SSH, gitee untouched), range `9217915..49f1609`.
- CI: run `45910223000` (head `49f1609`), both legs success; both legs' `Gates` measured `SELFTEST: PASS 770/770`, `GATE REPORT scope=ci changed=16`, `TOTAL_FAIL=0`; package-level ok ubuntu 21 / windows 15, FAIL 0; the other three candidate jobs are skipped as usual.
- This section is added by the docs-only commit after both legs went green (§11.2); that backfill commit's own CI run is an observation record only, and no further commit is made for it.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/gates-exempt.md	repo
docs/validation/gates-exempt_EN.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/00-briefing.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/01-design.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/input-list.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/freedom-list.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/MANIFEST.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/SHA256SUMS.txt	repo
tools/gates/lib/criteria/g2.js	repo
tools/gates/lib/criteria/g4b.js	repo
tools/gates/selftest.js	repo
tools/gates/g6-whitelist.txt	repo
```
