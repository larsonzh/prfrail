# FIX-GAP-77-SYMLINK-ESCAPE (P53) fix validation report

Slice: `FIX-GAP-77-SYMLINK-ESCAPE` (low-level security fix `[SLICE]`, tier M). Started 2026-11-03, directive v1.60, anchor `HEAD` = `dc785ed`. Controller model = `DeepSeek V4 Pro`.

## 1. Slice overview

- Defect: `internal/adapters/bounddir_unix.go` constrained the path with a single `openat(O_NOFOLLOW)`, which covers only the final component; a middle symlink could carry resolution outside the bound directory. P52 already demonstrated a read of bytes outside the boundary on Linux.
- Fix: replaced by a single constrained `openat2(2)` resolution, mask `RESOLVE_NO_XDEV|RESOLVE_NO_MAGICLINKS|RESOLVE_NO_SYMLINKS|RESOLVE_BENEATH` = `0x0F`, `SYS_openat2` = 437, a private 24-byte `open_how` layout, with no new dependency.
- Degradation: off linux/amd64, or when the kernel returns `ENOSYS`, the call fails closed with `PathError{Err: ENOSYS}`; no component-by-component walk fallback remains, and `allowUnsafeFallback` stays disabled in this slice.
- Lifecycle: a private `mu sync.Mutex` closes the fd hand-out and close paths; public signatures are frozen and the Windows and darwin implementations are untouched.
- Change set: production side is `bounddir_unix.go` only (measured `git diff --numstat` = `+158/-12`); test side reworks the standing regression plus one new case file each for payload, redbaseline, failclosed and storm.

## 2. Goal and rulings

The user issued the rulings on 2026-11-04; the verbatim text is archived in `01-arch/rulings.md`.

| Item | Ruling | Landing |
|---|---|---|
| (a) | The threat boundary includes external rename, so `openat2` is a hard requirement; degraded hosts fail closed | The main path runs only on linux/amd64; `ENOSYS` is refused outright |
| (b) | Mount import is part of the attack surface, so `NO_XDEV` is kept | Accepts the cost that a legitimate in-bound submount is refused |
| (c) | No depth cap is introduced | The depth-cap idea is dropped |
| (d) | Document sync belongs to ⑧ of this slice | This report and its `_EN` mirror |
| (e) | `openat2` is signed off as beyond the two-option list | Recorded in `freedom-list.md` |

## 3. Dispatch table

| Stage | agentName | model | Conclusion |
|---|---|---|---|
| ① r1 | `deep-reasoner` | DeepSeek V4 Pro | A/B hybrid design |
| ① red team | `independent-reviewer` | GPT-6.1 Sol | `FINDINGS` (2 High + 4 Medium) |
| ① r2 | `deep-reasoner` | DeepSeek V4 Pro | `openat2` finalised |
| ③ | `prfrail-tester` | DeepSeek V4.1 Flash | Real red plus the F1 guard; F2 recorded as SKIP |
| ②a | `prfrail-implementer` | DeepSeek V4.1 Flash | Production implementation landed |
| ③b | `prfrail-tester` | DeepSeek V4.1 Flash | Standing always-red case removed |
| ③c | `prfrail-tester` | DeepSeek V4.1 Flash | Warnings cleared, ⑤ findings disposed |
| ⑤ r1 | `deep-reasoner` | DeepSeek V4 Pro | 2 Medium and 3 Low |
| ⑤ r2 | `deep-reasoner` | DeepSeek V4 Pro | `PRE-REVIEW: PASS` |
| ⑥ | `independent-reviewer` | MAI-Code-1.1-Flash | Independent scan, verbatim in 5.1 |
| ⑦ | `independent-reviewer` | GPT-5.3-Codex | Final review, no findings |
| Batch assist ×2 | `deep-reasoner` | GLM-5.3-Flash | `openat` usage inventory and Linux compatibility facts |

## 4. Red and green evidence

Readings come from the Ubuntu VM (kernel 6.8, go1.22.12, run through `bash -lc`); the Windows host only covers compilation and static checks.

- Red (produced before any production change): `TestBoundDirMiddleSymlinkEscapesBound` failed with the verbatim text `GAP-77 regression: middle symlink escape succeeded and read outside the bound directory` and exit code 1; the escaped handle read 23 bytes of `SECRET-OUTSIDE-BOUNDARY` outside the boundary.
- Green (same case after the fix): passed; the refusal errno lands in the closed set `{ELOOP, ENOTDIR, EXDEV}`, measured as `ELOOP`.
- Positive control: the in-bound multi-level path `sub/nested/inside` reads successfully and equals the `IN-SIDE-BOUNDARY` canary byte for byte over the whole string.
- Mutation control: `PRFRAIL_P53_RED_BASELINE=1 go test -run TestBoundDirRedBaselineMiddleSymlinkEscape` fails on the fixed tree with `red baseline not reproduced`; in the default suite the case SKIPs behind the environment variable gate.
- Concurrency and storm: `TestBoundDirConcurrentCloseAndRead` runs 4 readers against 16 `Close` calls from the main goroutine and only allows in-bound exact bytes or `adapters: nil bound directory`; `TestBoundDirRenameStormNeverYieldsOutsideBytes` ran 200 rounds with 6336 successful in-bound reads and 1495 refusals, no outside bytes; `-race` reports no DATA RACE for `adapters` and `chain`.
- Fail-closed contract: after latching `openat2Disabled`, `TestBoundDirFailClosedWithoutOpenat2` asserts `*os.PathError`, `Op` equal to `openat`, `errors.Is(..., syscall.ENOSYS)` and no panic, and a second read is refused the same way; on non linux/amd64 the constant platform gate refuses directly.

## 5. Review chain

- ① Architecture: r1 produced the A/B hybrid, the red team (one-off temporary authorisation) reported 2 High and 4 Medium, and r2 finalised the single constrained `openat2` resolution; verdict `DESIGN: DONE`.
- ③ Tests first: the real red was produced before any production change; the red evidence is archived in `03-tests/red-baseline.md`.
- ②a Implementation: the single production file landed, with public signatures and platform branches unchanged.
- ③b and ③c: the standing always-red case was removed first, then the ⑤ findings were disposed in the test layer; recorded in `03-tests/remediation-report.md`.
- ⑤ Pre-review: r1 reported 2 Medium and 3 Low (`PRE-REVIEW: FINDINGS`); after the ③c dispositions r2 re-checked and returned `PRE-REVIEW: PASS`.
- ⑥ Independent scan: the verdict is in section 5.1 and the verbatim text is archived in `06-scan/review.md`.
- ⑦ Final review: no findings, `FINAL REVIEW: PASS` (`07-final/review.md`).

### 5.1 ⑥ raw output

Below is the verbatim output of the ⑥ independent scan (original language kept); both sides are byte-identical, the last fence line is the conclusion line, and nothing is added outside the fence.

```sixstage
## Section A

范围与输入：
- 生产实现： [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L1-L253)
- 关键回归测试： [internal/adapters/bounddir_escape_unix_test.go](internal/adapters/bounddir_escape_unix_test.go#L1-L197), [internal/adapters/bounddir_failclosed_unix_test.go](internal/adapters/bounddir_failclosed_unix_test.go#L1-L226), [internal/adapters/bounddir_symlink_storm_test.go](internal/adapters/bounddir_symlink_storm_test.go#L1-L220)
- 共享载荷： [internal/adapters/bounddir_escape_payload_unix_test.go](internal/adapters/bounddir_escape_payload_unix_test.go#L1-L97), [internal/adapters/bounddir_escape_redbaseline_test.go](internal/adapters/bounddir_escape_redbaseline_test.go#L1-L95)
- 权威证据： [docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/design-r2.md](docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/design-r2.md), [docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/rulings.md](docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/rulings.md), [docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/green-suite-report.md](docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/green-suite-report.md)

本轮只读，不做任何文件修改。已检查对象范围内的 git 状态与 diff；未发现 Windows/darwin 代码、go.mod、其他生产文件被触碰，且目标文件集中在本次审查集。

## Section B

(no findings)

本次独立扫描未发现需要修正的真实缺陷：openat2 ABI、fail-closed 门、fd 生命周期锁、回归测试闭集和可证伪性均与设计和实现相一致；当前状态不构成越界改动。

## Section C

证据摘要：
- ABI 与 resolve mask： [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L18-L61) 明确写入 SYS_openat2 437、open_how 24B、RESOLVE_NO_XDEV/NO_MAGICLINKS/NO_SYMLINKS/BENEATH。与 linux/openat2.h 语义一致。
- Fail-closed 语义： [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L95-L176) 仅在 linux/amd64 上通往 openat2；其余平台和 ENOSYS 直接返回，且 openat2Disabled 只在 ENOSYS 时置位，不扩散到其他错误。
- fd 生命周期与 dir.mu： [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L73-L110) 与 [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L181-L253) 显示 BoundDir 的 fd/Close/verify 都在同一私有 mutex 保护下，未出现遗漏的读写点，且 Close 幂等。
- 回归测试闭集与正向控制： [internal/adapters/bounddir_escape_unix_test.go](internal/adapters/bounddir_escape_unix_test.go#L34-L197) 对 middle symlink 采用 {ELOOP, ENOTDIR, EXDEV} 的拒绝闭集，并用 in-bound canary 做正向验证；[internal/adapters/bounddir_failclosed_unix_test.go](internal/adapters/bounddir_failclosed_unix_test.go#L24-L226) 额外锁定 ENOSYS 契约和二次拒绝稳定性。
- 风暴与 Rename/F2 约束： [internal/adapters/bounddir_symlink_storm_test.go](internal/adapters/bounddir_symlink_storm_test.go#L78-L220) 证明 rename storm 不会泄出 outside bytes，且 F2/环境门控逻辑是有条件地跳过，不是静默放行。
- 真实验证：已执行 go test -count=1 ./internal/adapters，结果为：
  ok      github.com/larsonzh/prfrail/internal/adapters   33.473s

## Section D

可证伪性审计：
- 若去掉 openat2 runtime gate 或把 openat2Disabled 语义去掉，fail-closed 契约测试会立即失效；这已被 [internal/adapters/bounddir_failclosed_unix_test.go](internal/adapters/bounddir_failclosed_unix_test.go#L24-L226) 钉住。
- 若去掉 dir.mu 或关闭的 fd 失效，concurrent close/read 测试会在 [internal/adapters/bounddir_failclosed_unix_test.go](internal/adapters/bounddir_failclosed_unix_test.go#L100-L226) 中暴露 EBADF、泄露字节或 panic。
- 若误把拒绝闭集放宽到任意错误， [internal/adapters/bounddir_escape_unix_test.go](internal/adapters/bounddir_escape_unix_test.go#L86-L126) 立即判红；这正是该类测试的可证伪性。
- git 证据与范围核对也没有发现越界修改；本轮不存在 go.mod、Windows/darwin 或其他生产文件的覆盖改动。

## Section E

无法验证部分：
- F2 bind-mount 的真实 namespace 环境需要特权/namespace 支持，因环境限制会 skip；这不是代码缺陷，而是宿主能力限制。
- 这次扫描只审查了目标改动集，不扩展到其他 repo 组件；如果后续出现新的宿主环境差异，需重新做平台分支验证。

推翻条件：
- 只要出现下述任一事实，即可推翻当前结论：
  - 目标实现未真正使用单次 openat2 受限解析；
  - 非 linux/amd64 或 ENOSYS 分支被绕过并泄漏到正常路径；
  - dir.mu 未覆盖所有 fd 生命周期/close 读写点；
  - git status/diff 显示出了 go.mod、Windows/darwin 或其他生产文件的无授权改动。

INDEPENDENT SCAN: PASS
```

## 6. Gate readings

- Host (Windows): `gofmt -l` prints nothing; `go build ./...` and `go vet ./...` exit 0; `go test ./...` reports ok for all fifteen packages.
- VM (go1.22.12): all fifteen packages ok; `-race` reports no DATA RACE for `adapters` and `chain`.
- Gates: the ④ and ⑨ phase readings (code state before this report pair landed) are `TOTAL_FAIL=0`, with `SELFTEST 697/697` matching the baseline (the production change does not touch gate implementations).
- Re-running the tree scope after this report pair landed reads `TOTAL_FAIL=2`: `G7-b` is red only because the fence rows are not yet in the index (it turns green once ⑧b stages them), and `G8-b` is red only because `MANIFEST.md` does not yet cover `06-scan/review.md` (the controller recomputes MANIFEST and SHA256SUMS after this report lands).
- Criteria on this report pair itself: the `G3(a)` six shapes (lines, headings, bold, pipes, blank lines, backticks) agree on both sides; `G4a` reports UTF-8 with BOM and LF endings; `G4b` reports zero unknown ideographs; both `G7-e` and `G7-f` pass.

## 7. Deviations and open items

- Open item 1 (F2): the bind-mount counterexample SKIPs on the VM because `unshare(CLONE_NEWUSER)` returns `EINVAL`, `unshare(CLONE_NEWNS)` returns `EPERM` and passwordless sudo is unavailable; recorded as an environment limit, not as a pass.
- Open item 2 (kernel floor): a real host below kernel 5.6 is unavailable, so the `openat2Disabled` latch contract test stands in, mirroring the production constant gate one to one.
- Open item 3 (mutation check): the `PRFRAIL_P53_RED_BASELINE=1` control command is left for ⑨ to re-execute; this report records only the command and its reading rules.
- D1 (structural conflict): section 4.1 of the ②a report demanded three passing tests, which conflicts with the ③ deliverables; ③b resolved it by restructuring.
- D2 (historical archive): the `-tags p53_red_baseline` records from ③b stay unchanged inside `03-tests/green-suite-report.md` and `MANIFEST.md`; the code comments and this report carry the replacement command.
- D3 (lock-order wording): described as measured - `readBoundEvidence` and `CloseBoundDirs` release `store.mu` before touching `dir.mu`, while `bindEvidenceDir` takes `dir.mu` then `store.mu`, all sequential with no nested holding.
- D4 (transcription correction): the "27 bytes" in the task pack misreads "27 byte-exact reads" from the ③c report; the in-bound check compares the whole string byte by byte and depends on no byte length.
- Boundary: the production change figures follow this run of `git diff --numstat`; earlier snapshot readings differ.
- Known limit: the Windows and darwin implementations are untouched in this slice; their counterexample coverage stays with the existing cases and is not repeated here.

## 8. artifacts fence

The rows below are the complete change set of this slice (evidence payloads, this report pair, the ledger, the specification documents and the code/test edits; path + disposition, TAB separated, dictionary order).

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/CHECKLIST_v43.md	repo
docs/t027/CHECKLIST_v43_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/00-briefing.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/design-r2.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/design.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/red-team-sol-review.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/rulings.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/green-suite-report.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/red-baseline.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/remediation-report.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/05-prereview/review-r1.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/05-prereview/review-r2.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/06-scan/review.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/07-final/review.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/08-gates/gate-readings.txt	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/09-glm/glm-1-openat-inventory.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/09-glm/glm-2-linux-facts.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/MANIFEST.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/SHA256SUMS.txt	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/freedom-list.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/input-list.md	repo
docs/validation/fix-gap-77-symlink-escape.md	repo
docs/validation/fix-gap-77-symlink-escape_EN.md	repo
internal/adapters/bounddir_escape_payload_unix_test.go	repo
internal/adapters/bounddir_escape_redbaseline_test.go	repo
internal/adapters/bounddir_escape_unix_test.go	repo
internal/adapters/bounddir_failclosed_unix_test.go	repo
internal/adapters/bounddir_symlink_storm_test.go	repo
internal/adapters/bounddir_unix.go	repo
```

## 9. Commit and CI

- Main commit: `7802819` (`fix: close the middle-symlink escape on the bound dir via openat2 (P53)`, 37 files, +2873/−62), pushed to origin (never gitee).
- CI run `41707160563` (head `7802819`, attempt 1): green on the first attempt on both legs (`Go ubuntu-latest` and `Go windows-latest` both success; the 3 dispatch-only jobs skipped as usual).
- Gates measured on both legs: `SELFTEST: PASS 697/697`, `GATE REPORT scope=ci changed=37`, `TOTAL_FAIL=0` (ubuntu leg 19:54 / windows leg 19:56).
- Package-level hard evidence: windows leg `ok internal/adapters 35.644s`; ubuntu leg Test `ok internal/adapters 13.415s` plus Race leg `ok internal/adapters 91.205s`, `ok internal/chain 2.145s`.

## 10. Conclusion

- GAP-77 is closed in this slice: the single `openat` is replaced by a single constrained `openat2` resolution, non linux/amd64 and `ENOSYS` fail closed; the real red was produced before the code change and turned green after it, with the mutation control command retained.
- Review chain: the ① double pass, the ③ red-to-green reversal, the ②a implementation, the ③b and ③c rework, the two ⑤ rounds, the single ⑥ round and the ⑦ final review are all complete, with no ⑦ finding.
- The three open items (F2 environment limit, kernel-floor substitute criterion, mutation check left to ⑨) are recorded in section 7 and do not weaken this slice conclusion.
- This slice is committed as `7802819` and pushed to origin (never gitee), with CI run `41707160563` green on both legs (see section 9).
