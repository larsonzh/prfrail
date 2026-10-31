# ProofRail Validation Report — P42-T027B-EVIDENCE-BIND

| Field | Content |
|---|---|
| Slice | `P42-T027B-EVIDENCE-BIND` |
| Type | implementation class `[SLICE]` |
| Directive version | v1.60 |
| Baseline HEAD | `b355081` |
| Date | 2026-10-31 |

## Goals and boundaries

- **Goal**: close the two gaps C3 and N2 - one artifact's bound decision, its bytes, truncation flag, digest and identity must all come from the same already-open file object, so neither "stat by path, then read by path" nor "open, then read by path again" is allowed; the directory binding of `RunEvidenceDir` must not end when the string is returned and must carry through to the actual I/O.
- **Closure shape**: `internal/evidence` adds the bound-read primitive `BoundFile` add-only (`OpenBoundFile` / `OpenBoundFileTruncating` / `NewBoundFile` / `ReadAll` / `Hash` / `Identity` / `Close`), taking identity from the handle `Stat()`, while `ReadAll` and `Hash` only seek and then read that one handle; `internal/adapters` adds the directory-binding primitive (`BoundDir` / `OpenBoundDir` / `OpenBoundDirFile`): on Windows `CreateFile` with `FILE_FLAG_BACKUP_SEMANTICS` opens the directory handle, `GetFinalPathNameByHandle` derives the real name for a relative open, and the file handle's final path re-checks it; on Unix `openat` with `O_NOFOLLOW` and `O_CLOEXEC` resolves relative to the dirfd.
- **Scope**: `internal/adapters/**` and `internal/evidence/**`.
- **Boundary**: no touch to contracts / schema / fixtures / `ci.yml`; `internal/chain` / `gates` / `snapshot` / `console` unchanged; public signatures (the string return of `RunEvidenceDir` and the three seam function variables) verbatim unchanged; no gitee push; `commit` / `push` stop at 10 for authorisation.

## Review chain conclusions

- **Round 1 architecture**: `deep-reasoner` + DeepSeek V4 Pro; the C3+N2 handle and directory binding design book, `DESIGN: DONE`, including revision notes D1 / D4 / F1 and the clause that the notes win over the body on conflict.
- **Red tests**: `RED-CONFIRMED` (invariants I1–I5 first pin reading the old content through one handle, refusal under replacement, ancestor-directory replacement, self-consistent identity and equivalent regression, then the implementation starts).
- **2 implementation, two passes**: the main body lands (bound-read primitive + directory-binding primitive + read-site replacement + frozen-digest backstop) -> the F fix pass (the Windows per-read re-check `confirmOpenedArtifact` lands, comments aligned with the real strength).
- **Gates**: `go build` / `go vet` / `go test ./...` green; three-scope and SELFTEST readings in section Gate readings.
- **Pre-review r1**: `PRE-REVIEW: PASS WITH FIXES` (1 Medium + 3 Low + 2 Info).
- **Pre-review r2**: `PRE-REVIEW: PASS` (F1–F4 verified closed item by item, plus the N1–N6 registration requirements).
- **Scan**: the first round was format-non-compliant (no A–E sections) and was reissued per contract; the second round is `INDEPENDENT SCAN: PASS`, this slice's scan budget 2/3.
- **Final review**: `FINAL REVIEW: PASS WITH FIXES` (3 Low + 4 Info, none blocking).

## Six-stage scan transcript

Source `docs/validation/evidence/P42-T027B-EVIDENCE-BIND/06-scan/review.md`; the fence transcribes it verbatim with zero rewriting, from `## Section A` through the conclusion line.

```sixstage
## Section A — 总体结论
本次只读复核未发现 P42-T027B-EVIDENCE-BIND 的 Medium+ 缺陷，当前实现与设计一致，并已按 C3+N2 句柄绑定/目录绑定语义落地。

## Section B — 发现清单
- 无

## Section C — 逐条证据与推理
在 [internal/evidence/boundread.go](internal/evidence/boundread.go#L23-L152) 中，`OpenBoundFile` 只保留单个已打开的 `*os.File`，随后 `ReadAll`、`Hash`、`Identity` 都从同一句柄读取，不再按路径做第二次解析；这与“打开后再按照路径重读”的 TOCTOU 缺陷相反。 [internal/adapters/bounddir_windows.go](internal/adapters/bounddir_windows.go#L39-L196) 和 [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L18-L99) 分别通过 Windows 的 `GetFinalPathNameByHandle` 复核与 Unix 的 `openat(..., O_NOFOLLOW)` 维持目录绑定，确保祖先目录被替换时读取仍落到原对象。 [internal/adapters/agent_runner_evidence.go](internal/adapters/agent_runner_evidence.go#L100-L201) 中的 `FrozenFacts` 也要求缺口统一返回 `ErrAgentRunnerEvidenceGap`，避免把“未完成”伪装成空摘要；其既定边界在 [internal/adapters/agent_runner_postflight.go](internal/adapters/agent_runner_postflight.go#L554-L613) 和 [internal/adapters/agent_runner_run.go](internal/adapters/agent_runner_run.go#L317-L346) 中被明确标注为本 slice 的 residual boundary，而非未修复漏洞。 [docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md](docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md#L31-L101) 的 D1/D4/F1 修订注记也保持一致：Q2 绑定对象收敛到基目录、释放外部不在本 slice 迁移、Windows 逐读复核已实现。

## Section D — 可证伪性审计
一个直接反例是：在 [internal/evidence/boundread.go](internal/evidence/boundread.go#L83-L152) 中删去 `Seek(0, io.SeekStart)` 和 `io.LimitReader(bound.file, readLimit(bound.bound))` 的同句柄读逻辑，改回 `os.ReadFile(path)`，或者删除 [internal/evidence/boundread_test.go](internal/evidence/boundread_test.go#L19-L92) 里 `TestBoundFileReadsOriginalAfterPathReplacement` 的断言；该变更会让“路径替换后仍读到旧文件”的安全判定失效，结论会从 PASS 翻成 FAIL。另一个反例是去掉 [internal/adapters/evidence_bind_dir_windows_test.go](internal/adapters/evidence_bind_dir_windows_test.go#L31-L88) / [internal/adapters/evidence_bind_dir_unix_test.go](internal/adapters/evidence_bind_dir_unix_test.go#L31-L147) 的祖先替换测试，或删除 [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L43-L66) 的 `O_NOFOLLOW`/`openat` 绑定；那么目录绑定就无法证明“原对象仍被读取”，PASS 也会被推翻。

## Section E — 我无法验证的部分与什么会推翻 PASS
1. Unix 上的 `openat(..., O_NOFOLLOW)` 与 symlink 反制场景需要在 CI/Unix 环境中真实执行；本次只读复核无法替代该端到端验证，因此 [internal/adapters/evidence_bind_dir_unix_test.go](internal/adapters/evidence_bind_dir_unix_test.go#L98-L147) 的 Unix 分支仍是“需 CI 保留”的部分。
2. 本 slice 明确保留的 residual boundary 仍处于后续读路径 sweep 范围，例如 [internal/adapters/agent_runner_run.go](internal/adapters/agent_runner_run.go#L316-L346) 中的 `readAgentRunnerChildPIDs` 以及 [internal/adapters/agent_runner_replay_store.go](internal/adapters/agent_runner_replay_store.go#L220-L270) 的部分读入口；如果后续实现把这些位置也改回 `os.ReadFile`/`os.Stat` 后再读，或如果以上红测在 Windows/Unix 上任一失败，PASS 将被推翻。
3. 只要同句柄绑定、目录绑定和缺口门槛仍保持不变，且 Windows 的全量/聚焦检查结果维持为 0 FAIL，则本结论继续成立。

INDEPENDENT SCAN: PASS
```
- Transcript note: the first scan round was format-non-compliant (no A–E sections) and was reissued per contract, so the second compliant output is transcribed here; the archive header states truthfully that the first round was not adopted, budget 2/3.

## F-item closure

The final-review closed set is `FINAL REVIEW: PASS WITH FIXES` (3 Low + 4 Info); the FR rows below map the 7 report numbering one-to-one onto the same batch of registered items and evidence anchors (Low = FR-01..FR-03; Info = FR-04..FR-07).

| # | Item | Severity | Location | Disposition |
|---|---|---|---|---|
| 5 r1 F1 | the per-read re-check required by design counter-example B was not implemented | Medium | `internal/adapters/bounddir_windows.go:78-100` | Closed: `confirmOpenedArtifact` compares the file handle's final path against the directory handle's `currentName()`, and closes plus refuses on a mismatch; regression I3 PASS |
| 5 r1 F2 | `bindEvidenceDir` comments over-claim | Low | `internal/adapters/agent_runner_replay_store.go` binding block | Closed: comments and implementation match item by item |
| 5 r1 F3 | invariant 4 is not met literally (`hashLogs` / `fileHash` / `readPIDs` still resolve by path) | Low | `internal/adapters/agent_runner_postflight.go:552` / `:637` | Registered to known limitations: the three boundary comments state the P42 boundary truthfully; migration belongs to the later read-path sweep slice (see section Later-slice candidate registrations) |
| 5 r1 F4 | the D3 comment over-rates the verification strength | Low | `internal/adapters/agent_runner_replay_root.go:381-399` | Closed: rewritten as "one `Lstat` type gate retained, residual window minimal" |
| 5 r1 F5 | the unix red tests did not run natively, and `CloseBoundDirs` is not wired into run close-out | Info | `internal/adapters/evidence_bind_dir_unix_test.go` and the store lifetime | Registered to known limitations: the unix branch is covered by the CI leg; production wiring belongs to a later slice |
| 5 r1 F6 | dual-toolchain difference | Info | `go.mod` (`go 1.22`) and the pre-review r1 reproduction condition | Closed: reproduced green on both under `go1.22.12`, the difference is cleared |
| 5 r2 N1 | `CloseBoundDirs` has no production caller and the comment still says per-request handle | Low | `internal/adapters/agent_runner_replay_store.go:220-224` | Registered to known limitations: the handle lives with the store / process; production wiring and wording unification belong to a later slice |
| 5 r2 N2 | the Windows re-check is a prefix comparison and the reparse resolution boundary is not registered explicitly | Low | `internal/adapters/bounddir_windows.go:127-162` | Registered to known limitations: object-identity-grade verification belongs to a later slice |
| 5 r2 N3 | Unix `runs/<requestId>` absolute-target symlink vector | Low | `internal/adapters/bounddir_unix.go:43-66` | Registered to known limitations: two-stage dirfd belongs to a later slice |
| 5 r2 N4 | `agentRunnerFileHash` boundary comment wording | Info | `internal/adapters/agent_runner_postflight.go:643` | Registered to known limitations: wording only, no semantic impact |
| 5 r2 N5 | `tmp` archives only one leg of the full run | Info | `tmp/P42-T027B-EVIDENCE-BIND/full1.txt` | Registered to known limitations: the second leg is the CI legs, `tmp` is cleared before 8c |
| 5 r2 N6 | `OpenBoundFileTruncating` comment wording | Info | `internal/evidence/boundread.go:71-76` | Registered to known limitations: wording only, no semantic impact |
| 7 FR-01 | `CloseBoundDirs` has no production caller and its comment conflicts with D1 | Low | `internal/adapters/agent_runner_replay_store.go:220-224` | Registered to known limitations and listed as a later-slice candidate (production wiring plus comment unification) |
| 7 FR-02 | the Windows handle-grade re-check is a path prefix comparison and the reparse boundary is unproven | Low | `internal/adapters/bounddir_windows.go:127-162` | Registered to known limitations and listed as a later-slice candidate (object-identity-grade verification) |
| 7 FR-03 | Unix intermediate-component and absolute-target link vectors stay out of coverage | Low | `internal/adapters/bounddir_unix.go:43-66` | Registered to known limitations and listed as a later-slice candidate (two-stage dirfd) |
| 7 FR-04 | frozen-digest wording is generic (frozen manifest versus the `Facts` digest field) | Info | `internal/adapters/agent_runner_postflight.go:643` | Registered to known limitations: wording only, comparison behaviour unchanged |
| 7 FR-05 | the second leg of the full suite is not archived under `tmp` | Info | `tmp/P42-T027B-EVIDENCE-BIND/full1.txt` | Registered to known limitations: the CI legs are the authoritative reading |
| 7 FR-06 | truncating open-function comment wording | Info | `internal/evidence/boundread.go:71-76` | Registered to known limitations: wording only, contract unchanged |
| 7 FR-07 | darwin cross-compilation no longer works and the old validation entry is void | Info | `internal/adapters/bounddir_unix.go:61` and `docs/validation/t027-replay-convergence.md:38` | Registered to known limitations: this slice only registers the "criterion replacement"; the old document revision belongs to a later documentation slice |

## Deviations and registrations

| Id | Content | Source | Verdict |
|---|---|---|---|
| D1 | the bound object converges from the per-request evidence directory to the `runs` base directory | implementation measurements + revision note | Safe: on Windows renaming the ancestor of a directory that has an open handle returns `ERROR_ACCESS_DENIED` (even with `FILE_SHARE_DELETE`), so a per-request handle would turn the ancestor-replacement attack into a denial of service; binding the base directory lets the ancestor be renamed while the handle keeps addressing the original object tree. The read path resolves `<requestId>/<rel>` |
| D2 | the `Close` / `NewBoundFile` / `OpenBoundFileTruncating` split | implementation pass | Safe: the refusal class and the truncation class must be distinguishable, so record / pid class is refused at open while log / event class keeps the truncated degradation |
| D3 | `loadAgentRunnerReplayOwnership` keeps one `Lstat` type gate | implementation pass | Boundary acceptable: it preserves the leaf-symlink refusal; the residual window is minimal and the comment is rewritten to the real strength |
| D4 | `BuildReleaseInventory` is not migrated | revision note | Design item narrowed: the adapters-to-evidence dependency direction forbids reuse in the reverse direction; the release tree is operator-chosen and not a candidate-writable surface, so it belongs to the later read-path sweep slice |
| D5 | bounds of 64 MiB and 1 MiB | implementation pass | Safe: `agentRunnerWorkspaceEntryBound` is 64 MiB and `agentRunnerMetadataBound` / `replayStoreRecordBound` / `fileQueueEnvelopeBound` are 1 MiB, in the fail-closed direction |
| D6 | `CloseBoundDirs` has no production caller | pre-review r2 N1 | Registered: the handle lives with the store / process; the comment still says per-request, contradicting the D1 base-directory criterion, and is fixed together with the production wiring |
| D7 | darwin directory binding unavailable | pre-review r1 measurement | Criterion replacement: the standard library has no `syscall.Openat` on darwin, so `bounddir_unix.go` excludes darwin and `bounddir_darwin.go` supplies a fail-closed stub (`ErrBoundDirUnsupported`): darwin compiles but refuses directory binding at run time; the support matrix stays windows / linux, the old "Darwin/amd64 cross-compilation passes" entry at `docs/validation/t027-replay-convergence.md:38` is void, and revising the old document belongs to a later documentation slice |
| D8 | `allowTruncation` is decided by the caller per class | implementation pass | Safe: the contract is chosen explicitly at the call site and the primitive never guesses a class |

- **Residual surface registered**: a same-level replacement of `runs/<requestId>` is a sibling node outside the base-directory handle coverage, and the single production caller compares the frozen digest so the replacement fails closed instead of becoming trusted evidence; the Windows replacement window between the derived name and `CreateFile` is closed by the handle-grade re-check of `confirmOpenedArtifact`; Unix intermediate-component and absolute-target link vectors belong to the later two-stage dirfd slice.

## Models and metering

| Stage | Carrier | Model |
|---|---|---|
| 1 architecture | `deep-reasoner` | DeepSeek V4 Pro |
| 2 implementation | `prfrail-implementer` | DeepSeek V4.1 Flash |
| 3 tests | `prfrail-tester` | DeepSeek V4.1 Flash |
| 5 pre-review | `deep-reasoner` | DeepSeek V4 Pro |
| 6 scan | `independent-reviewer` | MAI-Code-1.1-Flash |
| 7 final review | `independent-reviewer` | GPT-5.3-Codex |
| 8 documentation | `prfrail-documenter` | DeepSeek V4.1 Flash |
| batch assist | `deep-reasoner` | GLM-5.3-Flash (2 runs: read-path census / counterexample check) |

- **Primary model facts (user ruling of 2026-10-31, registered truthfully)**: from P37 to P41 the primary model had been set to `DeepSeek V4 Pro` by accident, which made triage efficient but raised the API cost; from P42 the primary model is back to `DeepSeek V4.1 Flash`.
- **Gate perceptibility**: the directive sets no hard gate on the primary model because it is an operator environment setting the terminal gates cannot see; 1 and 5 are still carried by V4 Pro under the whitelist, which is a separate matter from the primary-model choice.
- **Long-term policy**: Flash by default (High tier); for slices with low-level concurrency / state-machine / handle complexity, a single design document over 10,000 characters, deep root-cause analysis of large CI logs, or two consecutive behaviour drifts on Flash, the primary model should propose switching to Pro or switch after authorisation; the Max tier is never the default.

## Gate readings

- tree scope: `TOTAL_FAIL=0` (re-run after precise staging, worktree and index identical).
- index scope: `TOTAL_FAIL=0` (changed 34).
- Self-test: `SELFTEST: PASS 697/697`.
- G6-1 status-tense scanned lines: 418; no pending/completion conflict.
- ci scope is covered by the two CI legs; see the submission and CI section.

## Commit and CI

- Main commit: `cab3e08` `feat: bind evidence reads to one handle and the evidence directory (P42-T027B-EVIDENCE-BIND)`; pushed to origin as `b355081..cab3e08`.
- CI run `41194972405` went double green on the first attempt: `Go ubuntu-latest` success, `Go windows-latest` success.
- ubuntu leg hard evidence: `ok ... internal/adapters 12.230s`, `ok ... internal/evidence 0.015s`, Race leg `ok ... internal/adapters 89.088s`, `SELFTEST: PASS 697/697`, `GATE REPORT scope=ci changed=34`, `TOTAL_FAIL=0`.
- windows leg hard evidence: `ok ... internal/adapters 31.296s`, `ok ... internal/evidence 0.181s`, `SELFTEST: PASS 697/697`, `GATE REPORT scope=ci changed=34`, `TOTAL_FAIL=0`.
- The back-fill commit behind this report is a documentation increment; its CI run is observe-only and does not count as a slice round.

## Later-slice candidate registrations

- Read-path sweep: the N2 binding migration of `hashLogs` / `readPIDs` / `fileHash` / `scanWorkspace`.
- Unix two-stage dirfd: open each level relative to its parent with `openat` to cover the intermediate-component symlink vector.
- Windows reparse object-identity-grade verification: replace the current prefix comparison.
- `CloseBoundDirs` production wiring: call it on the run close-out path and unify the comment wording.
- darwin old-document supersede: revise the cross-compilation entry in `docs/validation/t027-replay-convergence.md`.
- Primitive-level test gap: read-after-close, hash over-bound refusal, adversarial identity same-source case (from batch assist).

## Known limitations

- **Substance of section E in the scan**: the Unix `openat` / `O_NOFOLLOW` symlink countermeasures must really run in a CI/Unix environment and this read-only review cannot replace that; `readAgentRunnerChildPIDs` and part of the replay-store read entry points remain this slice's residual boundary; if any bound-handle binding, directory binding or gap gate falls back to reading by path, or a red test fails on either platform, the conclusion turns over.
- **N1 · status = accepted (later-slice wiring)**: `CloseBoundDirs` has no production caller and the comment still says per-request handle; anchor `internal/adapters/agent_runner_replay_store.go:220-224`.
- **N2 · status = accepted (later slice)**: the Windows re-check compares the file handle's final path against the directory handle's current real name by prefix, and neither the reparse resolution boundary nor the prefix-only limit is proven at object identity grade; anchor `internal/adapters/bounddir_windows.go:127-162`.
- **N3 · status = accepted (later slice)**: `openat` keeps an ancestor from being redirected, while the absolute-target symlink and intermediate-component vectors of `runs/<requestId>` stay out of coverage; anchor `internal/adapters/bounddir_unix.go:43-66`.
- **N4 · status = accepted (wording layer)**: the `agentRunnerFileHash` boundary comment uses frozen manifest generically for the frozen digest, while the real comparison anchor is the manifest digest field in `Facts`; anchor `internal/adapters/agent_runner_postflight.go:643`.
- **N5 · status = accepted**: `tmp/P42-T027B-EVIDENCE-BIND/` archives one leg of the full run only, the second leg is the CI legs, and `tmp` is cleared before 8c.
- **N6 · status = accepted (wording layer)**: the comment wording of the truncating open function, with the contract unchanged; anchor `internal/evidence/boundread.go:71-76`.
- **FR-01 · status = accepted + later slice**: same origin as N1, and the production wiring plus comment unification are handled in one slice.
- **FR-02 · status = accepted + later slice**: same origin as N2, object-identity-grade verification is listed as a later slice.
- **FR-03 · status = accepted + later slice**: same origin as N3, two-stage dirfd is listed as a later slice.
- **FR-04 · status = accepted (wording layer)**: same origin as N4, comparison behaviour unchanged.
- **FR-05 · status = accepted**: same origin as N5, the CI legs are the authoritative reading.
- **FR-06 · status = accepted (wording layer)**: same origin as N6.
- **Batch-assist finding · status = accepted + later slice**: three primitive-level test gaps (read-after-close / hash over-bound refusal / adversarial identity same-source); the main chain avoids them via single-read-then-close and the `ReadAll` path, new callers must self-verify.
- **FR-07 · status = accepted (criterion replacement)**: the support matrix is windows / linux only; revising the void entry at `docs/validation/t027-replay-convergence.md:38` belongs to a later documentation slice.
- **D7 · status = accepted (criterion replacement)**: the standard library has no `syscall.Openat` on darwin, so `internal/adapters/bounddir_unix.go` excludes darwin and `internal/adapters/bounddir_darwin.go` refuses directory binding with `ErrBoundDirUnsupported` instead of falling back to a path-based open; darwin now compiles and fails closed at run time, the supported platform surface stays windows / linux, and revising the stale document entry belongs to a later documentation slice.

```artifacts
docs/t027/REMAINING_SLICES_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/05-prereview/review.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/06-scan/review.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/07-final/review.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/freedom-list.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/MANIFEST.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/SHA256SUMS.txt	repo
docs/validation/p42-evidence-bind_EN.md	repo
docs/validation/p42-evidence-bind.md	repo
internal/adapters/agent_runner_evidence.go	repo
internal/adapters/agent_runner_pinned_cli.go	repo
internal/adapters/agent_runner_postflight.go	repo
internal/adapters/agent_runner_replay_root.go	repo
internal/adapters/agent_runner_replay_store_launches.go	repo
internal/adapters/agent_runner_replay_store_terminals.go	repo
internal/adapters/agent_runner_replay_store.go	repo
internal/adapters/agent_runner_run.go	repo
internal/adapters/ai_probe_copilot_executor.go	repo
internal/adapters/bounddir_darwin.go	repo
internal/adapters/bounddir_unix.go	repo
internal/adapters/bounddir_windows.go	repo
internal/adapters/bounddir.go	repo
internal/adapters/evidence_bind_dir_unix_test.go	repo
internal/adapters/evidence_bind_dir_windows_test.go	repo
internal/adapters/evidence_bind_semantics_test.go	repo
internal/adapters/evidence_seam_signatures_test.go	repo
internal/adapters/filequeue.go	repo
internal/adapters/spawn_recovery_test.go	repo
internal/evidence/boundopen_unix.go	repo
internal/evidence/boundopen_windows.go	repo
internal/evidence/boundread_test.go	repo
internal/evidence/boundread.go	repo
internal/evidence/disposition.go	repo
internal/guard/process_windows.go	repo
tools/gates/cjk-newwords.txt	repo
```
