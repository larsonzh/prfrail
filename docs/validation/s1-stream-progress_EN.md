# ProofRail S1-STREAM-PROGRESS Report (P70, layer-2 governance side, implementation class [SLICE], tier S)

Date: 2026-11-13. Status: documentation assembly finalised (steps 1-7 closed). Slice: `S1-STREAM-PROGRESS` (P70, layer-2 governance side, implementation class `[SLICE]`, tier S).

## 1. Executive Summary

- Slice: `S1-STREAM-PROGRESS` (P70, layer-2 governance side, implementation class `[SLICE]`, tier S), primary model `DeepSeek V4.1 Flash (deepseek)`.
- Delivery: the lightweight `run --progress` stream switch, off by default; new file `internal/console/progress.go`.
- Conclusion: form A landed; default behaviour unchanged; `--json` semantics unchanged; gates all green.
- Metering: step 5 once (`PRE-REVIEW: PASS`), step 6 once (`INDEPENDENT SCAN: PASS`), step 7 final review 1 + re-review 1 (`FINAL REVIEW: PASS WITH FIXES` followed by `RE-REVIEW: FINDINGS`); blind review 0; GLM batch assistance 2 times.
- Counts: **R = 1, D = 2** (R = 1 counts the CI first-run red and its fix repush under §9.2; temporary-file D accumulated twice; both were capture files the controller misplaced at the `tmp/` root inside the post-wrap-up capture slot, each removed item by item with zero residue, registered under the §5.1 non-self-healing tier for a user ruling).

### 1.1 Exceptions and Deviations

a) During the step-3 test remediation the test carrier misused `git checkout -- internal/console/progress_test.go` while that path sat in the index as an empty intent-to-add entry, so the working copy was emptied; it was restored from VS Code Local History (790 entries) and confirmed identical to its pre-probe state by three proofs: identical line numbering, a green baseline after the restore, and the pre-change `git diff --stat` showing 790 insertions.
b) After the Low remediation of the step-7 re-review, **no third step-7 round was opened**: the tier-S quota of final review 1 + re-review 1 was exhausted, so the item was closed per §7.8.3 "quota exhausted", with the remediation self-verified by mechanical readings (see §4).
c) **Temporary-file D = 1**: in the capture window between the cleared wrap-up and the step-10 receipt the controller wrote one gate capture to the `tmp/` root (`tmp/g71.txt`, which is not inside `tmp/<slice-id>/`), violating the §5.1 location closed set; it was self-reported on the spot, deleted item by item with `Remove-Item -LiteralPath`, and the zero-residue readings were submitted: `Test-Path` moved from `True` to `False`, the recursive file count under `tmp` went from 102 to 101 (1 once the `gates-cache/g4b` exception is excluded), the top-level names returned to `gates-cache,.gitkeep`, and the untracked path count was 0. Per the §5.1 self-healing tier a controller-owned violation only takes effect after **same-round confirmation by an independent review layer**, and this slice had already spent its step-7 quota (1 + 1) with no review layer left in the round, so it is registered under the **non-self-healing tier** and escalated for a user ruling; the violation touched no deliverable bytes (the capture was a plain-ASCII gate digest and never entered the index). **Second occurrence (same class)**: while re-reading the gate results the controller again wrote two captures to the `tmp/` root (`tmp/gi.txt`, `tmp/gt.txt`); both were self-reported on the spot and deleted item by item, with zero-residue readings: `Test-Path` was `False` for both paths, the recursive file count under `tmp` went from 103 to 101 (1 once the `gates-cache/g4b` exception is excluded), the top-level names returned to `gates-cache,.gitkeep`, and the untracked path count was 0. The two occurrences total D = 2, close to the §5.1 ceiling that stops the slice at an accumulated D of 3; the tail of this slice switched to inline script printing only and wrote no further in-repo files.

## 2. Switch Design

- Data source: every engine state event flows through `Engine.transition` → `stateWriter.append` → `FileEventStore.Append`, and the last one **calls `file.Sync()` per event** ⇒ events land in order.
- Form A: `progressStore` decorates `*chain.FileEventStore`; `Append` forwards to the durable store first and prints only on success; `Load` is a pure forward.
- **SealSink fidelity**: the decorator must implement and forward `SealTaskEvidence` / `LoadTaskSeal` (three `Options.Events` type assertions live in `internal/chain/engine.go` and `seal_verify.go`); without the forwarding the seal replay is silently disabled or the acceptance closure fails closed.
- Output split: `--progress` writes **stderr** only; stdout keeps its original `writeRunSummary` / `writeJSONResponse` paths ⇒ default and `--json` output stay byte-identical.
- Line format (fixed): `progress: <occurredAt> <entity> <fromState>-><toState> <reasonCode>`; `entity` is `chain` / `task <taskId>` / `step <taskId>/<stepId>`; empty fields are filled with `-`.
- Three injection points: `ExecuteNoopRun`, `executeAgentRunnerRun`, and the read-only replay point `supervisedOpenWindowAt`, which is **not decorated**.

## 3. Data Source and Form Choice

- `events/state-events.jsonl`: readable while running (per-event sync, non-exclusive append, close after each write).
- AgentRunner's `events.jsonl`: written by the external pinned CLI without sync and read by prfrail only at the final hop ⇒ unusable as a live progress source.
- Rejected B (polling): needs a goroutine lifetime, ctx cancellation, poll cost and a last-line ordering race.
- Rejected C (AgentRunner events only): unreachable while running and blind to task/step start and end.

## 4. Counterexamples and Mutation Checks

Test list (R1-R9, one line each):

- R1 with the switch off the stdout stays byte-identical and stderr stays empty; R2 `--json` and `--json --progress` produce byte-equal stdout; R3 a broken writer changes neither the Append result nor the event bytes; R4 the decorator leaves the event bytes untouched; R5 the line-format snapshot; R6a/R6b the SealSink capability and the write-path forwarding; R7 the usage text and flag acceptance; R8 the AgentRunner route keeps its projection, event shape and line count; R9 the Supervised route keeps stdout and the `--json` report, with progress on stderr only.

Mutations (M1-M8, one line each; every one turned red in a real run and was restored verbatim with an unchanged SHA256 afterwards):

- M1 the nil guard; M2 writing to stdout; M3 returning the print error; M4 mutating event bytes; M5 the placeholder; M6 not forwarding `SealTaskEvidence`; M7 disabling JSON normalisation; M8 disabling the injection.

## 5. Review Chain

- Step 5 `PRE-REVIEW: PASS` (3 Low + 3 Info); F1/F5 were remediated, F2 is covered by the later R8/R9, and F3/F4/F6 stay as Info without remediation.
- Step 6 `INDEPENDENT SCAN: PASS` (0 findings).
- Step 7 final review `FINAL REVIEW: PASS WITH FIXES` (1 Low: the two routes lacked runtime coverage) → R8/R9 added → re-review `RE-REVIEW: FINDINGS` (the Low stayed open: no explicit JSON invariance) → JSON invariance assertions added → self-verified by mechanical readings.
- **Blind review ruling for step 7**: this slice touches no ownership / stop / identity semantics ⇒ the mandatory draw does not apply; the proportional draw is **not** taken either - the P67 blind review counts as a drawn slice inside the §7.7 window (the §12.4 clearing note "not counted as the §7.7 mandatory +1" means only that it does not consume a mandatory slot, not that it sits outside the §7.7 window; the P69 fact pack already applied the same reading: "window P67=1 / P68=2 / P69=3, fewer than 4 drawn, no draw") ⇒ within the window P68 / P69 / P70 are 3 undrawn slices, fewer than 4.
- **Whether P68 counts in the window**: it does - it is registered in the ledger's "later batch slices" section as `✅ 已完成`, not `⏸ 待评估` (only the latter marks an informal slice) ⇒ from P71 on the window is P68 / P69 / P70 / P71 as four undrawn slices, which §7.7 requires to be drawn.
- The step-7 re-review reported 1 Low (the AgentRunner and Supervised routes lacked an explicit JSON-invariance assertion); the controller added R8/R9 assertions and self-verified them with a green `go test ./internal/console/ -run Progress`; the step-7 quota of 1 + 1 was exhausted, so no third confirming round was opened => handled under §7.8.3 "quota exhausted", to do: fold into the next step-7 input (P71 must cite this line at kickoff).

## 6. Gate Readings

Raw step-4 batch readings (Windows host, 2026-11-13; the step-9 rerun after this documentation assembly gives the finalised readings):

- `gofmt -l .` empty output; `go build ./...` exit 0; `go vet ./...` exit 0; `go test ./... -count=1` ok packages 15 / FAIL 0.
- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 737/737`.
- `node tools/gates/gate.js --all --scope=tree` ⇒ `changed=19`, `TOTAL_FAIL=0`; `--scope=ci` ⇒ `changed=2`, `TOTAL_FAIL=0`.
- `GOOS=linux GOARCH=amd64 go build ./...` and `go vet ./...` both empty with exit 0 (host-side evidence of zero Unix impact; the CI ubuntu leg backs it up).
- Focused tests `go test ./internal/console/... -count=1 -run Progress -v`: all nine `--- PASS` lines are green.
- The ⑨ closing readings (rerun after the ⑧ documentation assembly, with an explicit stage before reading the index): `--scope=index` => `changed=33`, `TOTAL_FAIL=0`; `--scope=tree` => `changed=33`, `TOTAL_FAIL=0`; `--scope=ci` => `changed=2`, `TOTAL_FAIL=0`; `SELFTEST: PASS 737/737`; `go test ./... -count=1` exit 0; `GOOS=linux` `go build` and `go vet` both exit 0.
- Post-commit CI chain: the main commit `f77ab49` first run turned the **ubuntu `Test` leg red** (run `43400311736`, exit 1; the windows leg stayed green) => attempt 2 from `gh run rerun --failed` on the same run kept **ubuntu red and windows green with the same test, the same line and the same assertion string** (not a flake; registered as DR-14) => the fix commit `6a405bd` in run `43403891205` made **both legs succeed** (ubuntu 3m17s including Race / Contract fixtures / Gates, windows 4m28s; the package-level `ok github.com/larsonzh/prfrail/internal/console` lines read `3.477s` / `15.783s` and `24.133s`, with `--- FAIL` line counts of 0).

## 7. Known Limits

- Process exit codes of gate commands: the `StateEvent` field set carries no exitCode and `reason.code` is a semantic code ⇒ unreachable; supplying it needs an event-model change (out of scope).
- AI tool-call summary lines: the external child CLI batch-writes them without sync and prfrail reads them only at the final hop ⇒ unreachable.
- Start and end of in-port command gates (baseline / acceptance / review / scan): no matching event ⇒ only hook/gate shaped steps are covered.
- The AgentRunner **CLI success leg** is not asserted at the CLI layer: on Windows `replayStorePublicationDurability()` is always `Unproven` ⇒ the production dispatcher fails closed; R8 therefore covers the same production entry the `run` command actually calls.
- The remediation of the step-7 re-review Low did not receive a third independent round (quota exhausted, §7.8.3) and is folded into the next step-7 input.

## 8. Explicitly Not Executed

- No change to `internal/chain/**`, `internal/adapters/**`, `internal/gates/**`, `internal/tickets/**`, `internal/guard/**`, `schemas/**`, `testdata/contracts/**`, `.github/workflows/**`, `go.mod`.
- No change to `docs/DELIVERY_DIRECTIVE{,_EN}.md` (the directive body is untouched) and no change to `docs/CONTRACTS{,_EN}.md`.
- No change to the v50-v56 historical snapshots; no commit / no push; nothing pushed to gitee.
- The self-healing confirmer for D=2 is missing: the step-7 quota of 1 + 1 was already exhausted, so no same-round independent review layer could confirm it; it is registered under the §5.1 non-self-healing tier and escalated, and the user accepted it on 2026-11-13 with one condition - if a third occurrence of the same class (D=3) happens, the §5.1 accumulated ceiling forces the slice to stop.
- Recurrence prevention: inside the post-wrap-up capture window every gate capture must use inline printing (`node -e` / `execFileSync`) or land in `tmp/<slice-id>/`; never at the `tmp/` root. Both D=2 occurrences in this slice were instances of that violation.
- P68 counts in the §7.7 window (see §5) => **P71 must draw a blind review at kickoff**; P71 has to cite this line.
- DR-14 (CI first-run red, fixed in this slice): carrier = the `progressSummaryJSON` guard starting at `internal/console/progress_test.go:172`; the ledger registration is the DR-14 entry in the `docs/t027/REMAINING_SLICES{,_EN}.md` known-defects table; evidence = run `43400311736` (head `f77ab49`) with the ubuntu `Test` step red, attempt 2 still red and the windows leg green; the fix commit `6a405bd` (run `43403891205`) is green on both legs.
- **P71 landing pointer (2026-11-13)**: the "P71 must draw a blind review at kickoff" line recorded in §8 of this report has been executed - the P71 `S1-SCOPE-ANNOTATE` slice drew one blind review before its final review, with the verdict in `docs/validation/s1-scope-annotate{,_EN}.md` §5. The original line above stays untouched (§11.1).

Artifact list of this slice (repo-relative path + disposition, TAB separated; the fence is byte-identical in this file and its `_EN` mirror).

```artifacts
docs/OPERATIONS.md	repo
docs/OPERATIONS_EN.md	repo
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/s1-stream-progress.md	repo
docs/validation/s1-stream-progress_EN.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/00-briefing.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/01-arch/design-r1.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/02-batch/batch-01-state-events-enumeration.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/02-batch/batch-02-agentrunner-events-enumeration.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/04-reviews/scan-r1.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/04-reviews/final-r1.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/04-reviews/final-r1b.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/06-gates/gate-readings.txt	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/06-gates/mutation-evidence.txt	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/MANIFEST.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/SHA256SUMS.txt	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/freedom-list.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/input-list.md	repo
internal/console/cli.go	repo
internal/console/progress.go	repo
internal/console/progress_test.go	repo
internal/console/runtime.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_supervised.go	repo
internal/console/runtime_zero_model_e2e_test.go	repo
internal/console/stopper_test.go	repo
internal/console/supervised_cmd.go	repo
```

The following is the verbatim Section A-E transcript of the step-6 independent scan r1 (transcribed from `docs/validation/evidence/P70-S1-STREAM-PROGRESS/04-reviews/scan-r1.md`).

```sixstage
## Section A

INDEPENDENT SCAN: PASS

## Section B

(no findings)

## Section C

- 透明性判断：在 [internal/console/progress.go](internal/console/progress.go) 里，装饰器先转发到内层事件存储，只有在 `inner.Append(...)` 成功后才写一行 progress；它没有重编码、没有改写 event，且不吞掉真实事件错误。真正会返回的错误仍然来自 `inner.Append`，这保留了原始失败语义。没有新 goroutine、没有后台 poll、没有额外状态对象，也没有泄漏路径；它是纯观察器，不改变事件日志的持久化字节。
- SealSink 保护：装饰器显式声明 `var _ chain.SealSink`，并且两个方法都转发到内层存储；因此 [internal/chain/engine.go](internal/chain/engine.go) 中的三处 `Events.(SealSink)` 断言仍然成立，且 [internal/chain/seal_verify.go](internal/chain/seal_verify.go) 的恢复校验也仍可读到 seal。若装饰器未转发，R6 相关断言会立刻红；当前代码没有这个缺口。
- `--progress` / `--json`：在 [internal/console/cli.go](internal/console/cli.go) 中，`progressWriter` 只在 `--progress` 真值时绑定到 `stderr`，默认关闭时保持 `nil`，因此无 progress 输出，stdout 字节流保持原样；`--json` 仍走 stdout，和原有输出路径一致，不存在"把 progress 丢进 stdout"的回归。
- 改动面界定：输入包中的原始 diff 只涉及 docs 台账与 [internal/console](internal/console) 目录；它没有越过到 [internal/chain](internal/chain) 等生产代码。也没有修改运行时协议、事件模型、CONTRACTS 语义或外部 adapters 层。
- 测试审计：新增的 [internal/console/progress_test.go](internal/console/progress_test.go) 设计了 R1–R7，对默认关闭、stderr-only、writer 失败、事件字节不变、占位格式、SealSink 转发、usage 接受等关键面都做了断言；它们不是恒真语句，而且用"删守卫/改写输出/改写 SealSink"来制造显式红灯，未发现过度掩码或假绿。

## Section D

- 删除 `if progress != nil` 这个守卫：R1 会红，且会出现 nil writer panic；默认关闭路径会被破坏。
- 把 progress 改写到 stdout：R2 会红；`--json` 与非 JSON 的 stdout 对比直接不一致。
- 去掉 "打印失败忽略" 的处理：R3 会红；broken writer 直接暴露错误。
- 在 `Append` 前 mutate 事件：R4 会红；事件文件字节发生变化。
- 把 placeholder `"-"` 改为空：R5 会红；固定快照断言失败。
- 删除 `SealTaskEvidence` 转发：R6b 会红；seal 文件不会落盘且 hash 不一致。
- 删除 `--progress` 用法接受：R7 会红；usage 断言直接失败。
- 反之，当前代码在这些守卫下都能保持对齐；没有看到"删掉某守卫后仍旧绿"的情况。

## Section E

- 我无法验证的部分：本次结论基于静态读取、设计基线和新增测试的断言意图，不是基于真实进程执行输出的现场复现；因此对特定 live 环境下的并发/终端交互仍保留运行时观察限制。
- 什么会推翻我的 PASS：如果未来提交把 progress 改回 stdout、移除 SealSink 转发、修改 event 记录前后字节，或者在 [internal/console](internal/console) 之外引入生产代码修改；若测试被改成掩码式"只看 stdout 不看磁盘"，或有真实并发下的 race/泄漏出现，也会推翻本结论。按现有读取材料，没有发现这些信号。

INDEPENDENT SCAN: PASS
```
