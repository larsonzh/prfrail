# P56 slice validation report — `G4B-CORPUS-DECOUPLE` (G4b corpus enumeration decoupled from the index)

| Field | Value |
|---|---|
| Slice / queue slot | `G4B-CORPUS-DECOUPLE` (layer-3 alternative / suspended item, activated by the user's same-turn ruling on 2026-11-05) |
| Size / class | M / `[SLICE]` (criterion optimisation) |
| Anchor | `HEAD` = `f021af1` (the v45 snapshot commit) = origin/main |
| Directive | `docs/DELIVERY_DIRECTIVE.md` v1.60 → v1.61 once this slice lands |
| Controller model | `DeepSeek V4.1 Flash (deepseek)` |
| Evidence package | `docs/validation/evidence/G4B-CORPUS-DECOUPLE/` |
| Status | Reached the 10 stop point and waiting for the same-turn commit authorisation (not committed, not pushed) |

## 1. Goal and hard boundary

Goal: close OB-120. The tree / index branch of `ctx.listAllMarkdown()` used `git ls-files` (an **index** enumeration) while every file body was read with `git show <ref>:<p>` (the **object store**); after a staged deletion of a `.md` file the cold-run corpus shrank below the HEAD truth, so G4b over-reported, and a cache hit and a cold run could return different verdicts for the same input. This slice switches that branch to object-store enumeration so the corpus no longer depends on the staging state.

Hard boundary: no other criterion body touched; no `internal/**`, CONTRACTS / schema / fixtures or workflows change; `DEV_PLAN` not updated; never push gitee; stop at the 10 stop point after 9 / 8c are green.

## 2. artifacts fence

The entries below are this slice's change-set rows (path + disposition, TAB separated, sorted by path). The 8b stage has re-sorted this fence against the final change set and reconciled every row.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/00-briefing.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/03-test-report.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/06-independent-scan.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/07-final-review-r2.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/07-final-review-r3.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/SHA256SUMS.txt	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/freedom-list.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/input-list.md	repo
docs/validation/g4b-corpus-decouple.md	repo
docs/validation/g4b-corpus-decouple_EN.md	repo
tools/gates/cjk-newwords.txt	repo
tools/gates/lib/criteria/g4b.js	repo
tools/gates/lib/ctx.js	repo
tools/gates/selftest.js	repo
```

## 3. Change summary

1. **Enumeration decoupled** (`tools/gates/lib/ctx.js`): both branches of `listAllMarkdown()` now run `git ls-tree -r -z --name-only <ref>`; tree / index use `ref = 'HEAD'`, range / ci use `ref = baseRef`; `-z` plus `\0` splitting, with the case-insensitive `.md` filter kept; a failed git call still returns an empty list (an empty corpus turns G4b red, i.e. the direction is fail-closed).
2. **Cache version bumped** (`tools/gates/lib/criteria/g4b.js`): `CACHE_VERSION` and the stored `v` field both move from `v1` to `v2`; the `kind` segment is kept but re-documented as a namespace separator; the stale `cacheUsable()` doc block now says "both scopes read the same corpus revision ⇒ staged adds, edits and deletions cannot change the corpus, so a hit and a cold run are identical on the corpus".
3. **Assertions** (`tools/gates/selftest.js`): added `T114` (tree enumeration source), `T115` (cold/hot equivalence under a staged `.md` deletion), `T116` (`v` field ↔ version prefix), `T117` (the range `-z` shape) and `T118` (index must use the HEAD object store, with a fake `ls-files` that deliberately answers with the WRONG set); three hardcoded `'v1'` expectations were rebound to `G4B.internals.CACHE_VERSION` (same strength, no drift).
4. **Criterion text**: §6.3 already said "the whole-repository `.md` corpus at HEAD", so this slice is the **implementation catching up with the text**; the criterion body is unchanged. Only §0.3 gained the v1.61 row, §12.4 flipped OB-120 to `Resolved`, and appendix C gained `C.0ar`.
5. **Execution freedom exercised** (see `freedom-list.md`): FL-1 unified enumeration shape, FL-2 `-z`, FL-3 cache bump, FL-4 keep the `kind` segment, FL-5 empty list on git failure, FL-6 mirrored report pair; all fall under "wrong ⇒ change code", none touched an escalation boundary.

## 4. Execution pipeline

| Stage | Executor | Result |
|---|---|---|
| Step 0 | controller | directories, `00-briefing.md`, ledger, `input-list.md`, `freedom-list.md`, `MANIFEST` + `SHA256SUMS`, fence |
| 1 architecture | **skipped** | start message [0, 3] explicitly disables the architecture double play |
| 2 implementation | `prfrail-implementer` ×2 | 2a `ctx.js` enumeration (+10/−5); 2b `g4b.js` cache bump (+19/−17) |
| 3 tests | `prfrail-tester` ×3 | T114–T116 → T117 (5 Low) → T118 (7 r1 Medium); `SELFTEST` 697 → 703 |
| 4 gates | controller | `node --check` on all three files; the tree scope kept only two expected mid-slice failures (G7-b unstaged, G7-f not transcribed) |
| 5 pre-review | `deep-reasoner` + `DeepSeek V4 Pro` ×2 | r1 `PASS WITH FIXES` (3 Low) → r2 `PASS WITH FIXES` (no Medium+, 1 Low) |
| 6 scan | `independent-reviewer` + `MAI-Code-1.1-Flash` | `INDEPENDENT SCAN: PASS` (Section B = `(no findings)`) |
| 7 final review | `independent-reviewer` + `GPT-5.3-Codex` ×3 | r1 `FINDINGS` (1 High + 1 Medium + 1 Low) → r2 `RE-REVIEW: FINDINGS` (1 High + 1 Low) → r3 `RE-REVIEW: PASS` |
| 8a/b docs | `prfrail-documenter` ×2 + controller assembly | directive CN/EN (same shape, same line positions), ledger CN/EN, this report pair |
| 9 native validation | controller | three scopes + `SELFTEST: PASS 703/703` |
| 8c cleanup | controller | `tmp/G4B-CORPUS-DECOUPLE/` removed item by item; the readback shows zero residue |

## 5. Anomalies and fallbacks

| Id | Class | Fact | Disposition |
|---|---|---|---|
| R-1 | rework (R count) | 5 r1 Low#3: the range / ci shape change had no dedicated assertion | the controller judged "record only" insufficient for regression capture ⇒ 3 added `T117` and `SELFTEST` was re-run |
| R-2 | rework (R count) | 7 r1 High (`MANIFEST` / `SHA256SUMS` did not register `04-reviews/**`) + Medium (no index assertion) + Low (fence wording) | three-way consistency repaired, 3 added `T118`, wording fixed; 7 r2 judged M-1 and L-1 closed but H-1 open because `05-pre-review-r2.md` was still missing ⇒ after registering it 7 r3 returned `PASS` |
| D-1 | temp-artefact discipline (D count) | the 8a-1 documentation subagent created `tmp/g4bcd/` (not under `tmp/G4B-CORPUS-DECOUPLE/`) and deleted it itself; the deletion preceded the controller's reading, so the "reading (2) baseline" of the 5.1 self-heal tier was unavailable | the controller re-verified mechanically afterwards: `Test-Path tmp/g4bcd` = `False` and `git status --porcelain -uall` shows no such path ⇒ **zero residue confirmed**; registered as **D = 1** and declared in the 10 receipt so the user can rule on the disposition at the stop point |

**R / D summary**: R = 2 (which triggers the §5.2 high-risk row ⇒ the 5 re-review is mandatory and was executed as r2), D = 1.

## 6. Environment and thinking tier

| Item | Value |
|---|---|
| Host / terminal | Windows, one serial terminal (no parallel subagents) |
| Node | v24.17.0 |
| Controller tier | High (never changed, no pause) |
| 5 / 7 tiers | set in the operator UI; the controller records the qualified names only and never trusts self-reports (§2.3) |

## 7. Review conclusions (5 / 6 / 7)

- **5 r1**: `PRE-REVIEW: PASS WITH FIXES`; A-1..A-5 all hold; 3 Low (capture mojibake, missing after-change wall clock, missing range/ci assertion), no Medium+.
- **5 r2**: `PRE-REVIEW: PASS WITH FIXES`; L-1 / L-2 / L-3 / H-1 / M-1 / L-4 each judged "closed"; one new Low (`timings.txt` still narrated 702), fixed in place.
- **6**: `INDEPENDENT SCAN: PASS`, Section B = `(no findings)`; Section D named one **unasserted** pre-existing invariant (a broken whitelist must fail closed before the cache), no assertion added (the shipped code already throws before the cache), recorded as an observation.
- **7 r1 → r3**: `FINAL REVIEW: FINDINGS` → `RE-REVIEW: FINDINGS` → `RE-REVIEW: PASS`; r3 counted three-way = 8 / 8 / 8 and line-number spot checks = 3 / 3.

Verbatim transcription of the 6 scan (full text in `04-reviews/06-independent-scan.md`; the scan itself ran in Chinese, so the transcript is kept in its original language):
```sixstage
## Section A
INDEPENDENT SCAN: PASS
根据只读核对的实现和断言，没有发现违反这两条硬边界的情况：G4b 语料枚举确实被改成对象库读取，并且缓存键/版本逻辑仍然绑定要件，而且自测对关键回归有较强变异敏感性。

## Section B
(no findings)

## Section C
1) 对象库枚举真实生效，且不再依赖索引。读取主控落盘的 tools/gates/lib/ctx.js:536-551，listAllMarkdown 的实现是 git(['ls-tree', '-r', '-z', '--name-only', ref])，随后用 '\0' 分割，并用 /\.md$/i 过滤；这正是“对象库 + NUL 分隔 + JS 过滤”的组合，且注释明确指出不能用 ls-files 或 *.md pathspec，原因是会让 G4b 语料收缩并制造假阳性。也就是说，问题不只是报告里说了，而是实现上确实走了 tree/object store ，而不是 index。读取主控落盘的 tools/gates/lib/ctx.js:540-551 还特别说明了 -z 的必要性：非 ASCII 路径会被 quotePath 包裹，newline 分割会把路径弄坏，故 NUL 分隔是 load-bearing。

2) 缓存键仍然绑定解析 sha / whitelist 文本摘要 / 判据源码摘要 / 格式版本，并且 fail-closed 在缓存前。读取主控落盘的 tools/gates/lib/criteria/g4b.js:43-89 和 99-170：CACHE_VERSION 是 v2，keyString() 返回 [CACHE_VERSION, rs, kind, source, sha256(text), sourceDigest()]，其中 rs 来自 resolvedSha(ctx)，text 来自 whitelist.text，sourceDigest() 哈希了 WHITELIST_REL、cjkChars.toString() 与 runBody()；同时 run() 在缓存前就先读取 whitelist，并在 whitelist.error 时抛 UsageError，避免缓存掩盖坏白名单。validValue() 也要求 v === 2 且 corpus 是字符串数组，和缓存写入时的 { v: 2, corpus: [...] } 保持一致。缓存命中只读 corpus，不写 verdict，且任何 get 失败/读坏/写坏都 fail-open 回全量计算，这也和硬边界一致。

3) 新增断言具备变异敏感性，且覆盖了关键回归面。读取主控落盘的 tools/gates/selftest.js:6579-6663，T114/T117 明确锁定 object-store + -z + .md 过滤；T115 锁定“staged markdown deletion 后 cold/warm 相等且 verdict 不变”；T109/T110/T111 测试解析 sha、whitelist 摘要、sourceDigest 变化时必须 miss，而只有同一 key/同一判据源码/同一白名单才可 hit；T112 锁定 fail-open 与 key 形状，并确认 g4bCacheDisabled / no-sha 时都退回全量计算；T116 锁定缓存存储的 v 值必须等于 CACHE_VERSION。检索 tools/gates 中 listAllMarkdown 的唯一消费者也只有 G4b 自身和 selftest，没有发现其它判据被这次变动连带地依赖到这个枚举函数；因此“其它判据被波及”这点在当前仓库代码中未发现。

## Section D
- 改回 ls-files / 改回 newline split：T114、T117、T115 会红；T116 不会红，因为它只看 v 值和 corpus 结构，不管文件枚举实现。
- 删掉 -z 与 '\0' 分割：T114、T117 会红；T115 也会红；T116 不会红。
- 删除 resolved sha 段：T109 会红；T110/T111 不会红；T116 不会红。
- 删除 whitelist 摘要段：T110 会红；T109/T111 不会红。
- 删除 sourceDigest 段：T111 会红；T109/T110 不会红。
- 删除 validValue 里的 v===2 或软编码 v 绑定：T112 fail-open 2 会红；T116 会红；T108/T109/T110/T111 一般不会因这条单独变红。
- 删除 cacheUsable() 对 no-sha 的保护：T112 的“git 无法解析 sha 时静默全量计算、无缓存记录”会红；其余 T109-T111 不会特定红。
- 删除 whitelist.error 的 fail-closed：我没有看到当下断言覆盖“坏白名单错误应在缓存之前失败”的案例；这是目前没有覆盖的不变量。另一个没有覆盖的不变量是“另一个 consumer 不会再调用 listAllMarkdown”，但代码搜索显示当前 tools/gates 内没有其它消费者，故这条未覆盖风险已被消除。

## Section E
我无法验证的部分：我只读取了任务明确列出的主控落盘文件和实现文件；没有读取别的审查方清单，也没有执行 git/selftest；因此我不能对“其它审查团队是否还看到了别的命题”作独立确认。另一个无法验证点是：我未读取 diff.patch 或其它审查包内容，所以我不声明其余结论。

会推翻我的 PASS 的条件：如果在实际仓库里还有一条不在当前读到的实现中、但在未落盘分支/未读文件里调用 listAllMarkdown；如果后续新增了第二个 listAllMarkdown consumer；如果实际对象库语义不同于这里的 HEAD/object-store 表述；如果白名单 loader 在某个代码路径中返回 error 后仍被缓存绕过；或者如果该实现和当前断言的对象库语义存在未覆盖的环境差异。
INDEPENDENT SCAN: PASS
```

## 8. 8a self-check

| Item | Result |
|---|---|
| Mechanical cost-table audit | yes: the dispatch counts in §13 reconcile item by item with the `runSubagent` receipts (counts and model strings are checkable) |
| "Resolved" `file:line` | yes: R-1 / R-2 / D-1 in §5 all carry a landing point; each of the 7 findings has a landing point (`selftest.js` T118, the `MANIFEST.md` archive table, the §2 wording) |
| Old-shape grep | yes: `grep -n "ls-files" tools/gates/lib/ctx.js` leaves only `resolveRepoArtifact` (G7-b's independent call, not the corpus enumeration); `CACHE_VERSION` is `v2` across the repo |
| Every 6 / 7 Info / Low captured | yes: the three 5 r1 Lows, the 5 r2 Low, the 7 r1 High/Medium/Low and the 7 r2 High/Low are all registered in §5 and here; the reference set is the four transcripts in `04-reviews/` |

## 9. Review input package summary

- Before every round `make-review-pack.js` writes `<round>/diff.patch` + `<round>/input-manifest.txt` (per-file bytes / last line number, exit 1 on a missing input); six packages `review-r1`..`review-r6` were produced, all with `missing=0`.
- Input sanitisation: no IP, user name, host name, SSH path, credential or token.
- Blind-review isolation: 7's first round (r1) received **no** 5 / 6 list, only the slice definition, diff paths, the verbatim criterion sentences and the read-only constraint; later rounds attach only its own prior findings.
- 5 and 7 only ever received read-only material (no `execute` ability) and created or modified no file.

## 10. Falsifiability

| Mutation | Expected red | Evidence |
|---|---|---|
| Revert tree / index to `git ls-files` | `T114` / `T115` / `T118` | the fake answers `ls-files` with the wrong set or 127; `ab-probe-adr.txt` records the old-shape unknown set |
| Revert `CACHE_VERSION` to `'v1'` | `T116` (`v` mismatch), `T112` (key head) | all three bindings go through `G4B.internals.CACHE_VERSION` |
| Revert `validValue`'s `v === 2` to `1` | `T115` (the warm run becomes a miss ⇒ the `hit` check reds) | the assertion compares `hit` / `miss` directly |
| Drop `-z` or revert `split('\0')` to `split('\n')` | `T114` / `T117` / `T118` | the NUL and newline shapes are mutually exclusive under the fakes |
| Add an `ls-files` special case for index | `T118` | the fake returns only `.txt` for `ls-files`, so the list mismatches immediately |

**Independent corroboration (A12)**: with `docs/ADR_REGISTER.md` staged-deleted, the old enumeration lost one corpus character and reported one unknown, while the new enumeration had an empty difference set and the gate returned `PASS` (`ab-probe-adr.txt` / `probe-staged-delete-g4b.txt`); the index was restored with `git reset`.

## 11. Gate results

| Gate | Reading |
|---|---|
| `--selftest` | `SELFTEST: PASS 703/703` (697/697 before the change) |
| `--scope=tree` / `--scope=index` | `TOTAL_FAIL=0` once staged (final readings in §12) |
| `--scope=ci` | `TOTAL_FAIL=0` |
| `go build` / `go vet` / `go test` | no Go code changed; still executed in full per §6.1 and recorded at 9 |
| Focused gate | `node --check` on `tools/gates/lib/ctx.js` / `g4b.js` / `selftest.js` all pass |

**Timing** (same host, serial): `--selftest` 278653 ms → 214574 ms; the G4b single-criterion tree cold pass 129151 ms → 104598 ms and the warm pass 1554 ms → 1697 ms. The cold-run drop comes from `ls-tree` avoiding one index resolution; the warm delta is noise, and neither `--selftest` nor the three scopes regressed materially.

## 12. 8c wrap-up cleanup record

8c was executed by the controller: repo-wide diagnostics (the §3.4 table), IDE diagnostics summary, `gate.js --all --scope=tree` green (once staged), per-file encoding verification and a `tmp/` inventory check. **Order**: records persisted first → item-by-item `Remove-Item -LiteralPath` → mechanical readback.

**Final readings**: `--all --scope=index` / `--all --scope=tree` / `--all --scope=ci` all report `TOTAL_FAIL=0`, with the base gates likewise green (`gofmt` empty output, `go build` / `go vet` exit 0, `go test` `ok=15 FAIL=0`); `--selftest` = `SELFTEST: PASS 703/703`.

## 13. Cost and metering

| Item | Count |
|---|---|
| 2 implementation | 2 (`prfrail-implementer`, Flash) |
| 3 tests | 3 (`prfrail-tester`, Flash) |
| 5 pre-review | 2 (`deep-reasoner` + `DeepSeek V4 Pro`; r1 plus the r2 remediation re-review) |
| 6 scan | 1 (`independent-reviewer` + `MAI-Code-1.1-Flash`) |
| 7 final review | 3 (`independent-reviewer` + `GPT-5.3-Codex`; r1 final plus r2 / r3 re-reviews, exhausting the 1+2 budget) |
| 8 docs | 2 (`prfrail-documenter`, Flash) |
| Batch auxiliary | 2 (`deep-reasoner` + `GLM-5.3-Flash`: call-site and impact enumeration / new-assertion cross-check, both ending with `BATCH: DONE`) |
| Probes / paid external calls | 0 |

The metering source of truth is the session receipts; each `agentName` and `model` is read verbatim into the dispatch table of the 10 receipt.

## 14. Explicitly not executed

1. No `internal/**` or Go production code changed, so no opposite-`GOOS` `go vet` was run (the precondition of discipline H does not hold).
2. No dedicated `ci`-scope assertion: its `baseRef` shares the same implementation line with `range`, which `T117` covers; the gap is registered in `03-test-report.md` §4.
3. No assertion was added for "a broken whitelist must fail closed before the cache" (a pre-existing invariant named by 6 Section D): a pre-existing coverage gap, registered here as an observation rather than an out-of-scope change.

## 15. Precedent-ban check

Clauses marked "must not be used as a precedent" cited by this slice: **none**. OB-120 was disposed of through the normal clause (§12.4 disposition flip); no temporary authorisation and no OB-37 meta-rule channel was used.

## 16. Next steps

1. If the user accepts the D-1 disposition, this slice can be committed under the standing discipline (message: `fix: decouple the G4b corpus enumeration from the index (P56)`), pushed to origin only.
2. Optional follow-ups: a constructed assertion for the `ci` scope; an assertion for "broken whitelist fails closed before the cache" (both registered here).
3. The next queue item remains layer-3 item 8 `AT23-RESTRICTED-CANDIDATE` (⛔ blocked on the R-08 / R-09 rulings).
