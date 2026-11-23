# P85 `SECURITY-PROOF-SUBSYSTEM` slice report

[中文](security-proof-subsystem.md)

**Conclusion: `CLOSED`** (T2 / T3 / T8a all CLOSED; F1–F8 all closed). This report = the docs/validation landing of the evidence package `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/report.md` (2026-11-24 v70 backfill, user ruling). Differences from the evidence package total only six: the title, the [English] link line, the candidate line (point-in-time updated to three commits), Appendix B (sixstage fence), Appendix C (artifacts fence), and three G6-4 stale-literal micro-adjustments (each replacing the adjudicated literal with "total/changed files/source files" inline); the rest of the body matches the evidence package report.md verbatim.

- Slice: P85 (implementation-class `[SLICE]`, tier L, protocol-first); goal = close F3/F6/F8 carried over from P83 §8 (T2 frozen-policy artifact / T3 immutable install source / T8a BYOK proof carrier).
- Conclusion (code-review accounting): **T2 / T3 / T8a all CLOSED**; ⑨ evidence debt F17/F22 was settled on 2026-11-23 on dual-leg CI evidence (see §12; the no-evidence-no-close accounting is thereby closed).
- Commits: main commit `cc7d90f` + R=8 test-side fix `6480bb1` + backfill `18b2613`, pushed to origin (origin only, gitee forbidden).
- Date: 2026-11-23. Driver = DeepSeek V4 Pro Max (same-turn user authorization, fixed after kickoff).

## §1 Exception registry

| # | Item | Disposition |
|---|---|---|
| E-1 | R count reached 7 (§9.2 R>3 pause) | Every round stops for manual review + same-turn user authorization to continue; OB-142 registered (must not serve as precedent) |
| E-2 | After ⑦ quota 1+3 was exhausted, four extra reads were added (Carbon/Neon/Sol/Codex) | OB-141 (Codex fourth read), Sol fifth read recorded under the §7.8.3 user-authorization exception; the three-read trial is a governance-slice trial authorization; none may serve as precedent |
| E-3 | Targeted mini-read after R=7 (Carbon High, only reads the R=7 diff + tests) | OB-143 (§1.7.4 temporary authorization, must not serve as precedent) |
| E-4 | L4 comment-wording change happened after the review loop closed | Per P26 ruling D-2: attached diff excerpt (`04-reviews/09-l4-diff-excerpt.md`) + substantive advisor verification + user final review (registered at the ⑩ receipt) |
| E-5 | The Linux directory-slot fixture's initial minimal ACL was equivalent to mode, kernel did not persist, first VM run of that case failed | Fixture changed to extended ACL (named user u:4242 + mask) and both platforms went green; nature = test-fixture correction, not a production-code rollback |
| E-6 | Windows publication `CreateTemp`/`Link` path-based write residual | User final ruling: register the **S1 independent boundary** (see §8), no F04 exemption attached, no code path changed |
| E-7 | gopls unused cross-GOOS-variant diagnostics of 12 symbols (rescanned after fix) | Three-class disposal completed, three views cleared to zero (§10/§12) |
| E-8 | Main-commit CI windows leg 2 environmental failures (volume-root contention + 8.3 short-name spelling; local/VM both green) | **R=8** test-side fix (zero production change): volume-root probe contention retry+skip, parent-directory naming downgraded to object-identity comparison; commit `6480bb1` (user authorization 2026-11-23); re-run CI 45208001981 dual legs green (§12) |

## §2 Change summary

- **T2 (F3) frozen-policy artifact**: freeze before assembly + re-read by run/task/step/attempt + re-check before publication; `persistAgentRunnerFrozenPolicyProven` production assembly injects no proof (the same-name no-proof wrapper was deleted); the four-slot digest `:2` fixes `[策略目录|父目录|制品|锚]`, empty slot `∅`, refusal path keeps four slots (R=7).
- **T3 (F6) immutable install source**: `openInstallSourceLeaf` changed to `OpenBoundDirVerified` (same-handle classification/digest, errors carry directoryID); carrier publication re-binds and cross-checks the root chain before and after; Windows `verifyInstallSourceImmutability` keeps the leaf, probes `FileIDAt == leaf.FileID()` before and after (R=7); shebang 256-byte window + ELF PT_INTERP unchanged.
- **T8a (F8) BYOK proof carrier**: production assembly injects no Sandbox, both platforms default `VerdictUnknown` ⇒ zero-start refusal before credential parsing, classification `byok_confinement_unproven`; S1 keeps BYOK disabled and claims no proven OS write restriction.
- **Contracts**: `docs/CONTRACTS{,_EN}.md` :214 (byok_confinement_unproven closed set), :216 (T3 install source + Windows root chain), :220 (BYOK zero-start gate + fd3), :231 (ManagedWriteProofShape create/reuse), :233 (domain :2 four-slot).
- **Directive ledger**: `docs/DELIVERY_DIRECTIVE{,_EN}.md` OB-140/141/142/143 (four elements, must not serve as precedent).
- **Tests**: R4-phase M13–M17, R6-phase slot_order/sandwich/windows_forensics/windows_publish_recheck, R7-phase `p85_r7_role_suffix` / `p85_r7_refusal_slot` / `p85_r7_open_bounddir_verified` + directory-slot fixtures (extended-ACL discrimination + closed-end verification leg) totaling 20+ cases.
- **IDE warnings cleared to zero**: 9 dead-code deletions (including the consequential `openAgentRunnerPolicyFile`); 3 symbol relocations (`managedAgentPolicyWriteProofDomain`→`linux||windows` file; `boundRootName`/`errBoundRootArtifact`→`!darwin` file (freebsd compilation preserved); unsupported deletes `installSourceInterpreter`; `p85R4SwapPolicyDirectory`→linux test file); zero `//nolint`.

## §3 Execution pipeline (steps / roles / status / rework)

① design Rev1–Rev4 (V4 Pro Max) → ②a protocol-first (CONTRACTS five bilingual spots) → ②b implementation (Flash High) → ③ tests →
④ gates → ⑤ pre-review r1/r2 → ⑥ scan r3/r4 → ⑦ blind review (22 items) → ⑦ final review (Sol Max) → re-review r1/r2/r3 →
R=4/R=5/R=6 fixes → four-read same input (Carbon/Neon/Sol/Codex) → **R=7 merged fix batch** (OB-142) →
targeted mini-read (Carbon High, OB-143) → ⑧a first draft → ⑨ native verification → ⑧b final draft → ⑧c closeout → ⑩ stop point.
**Rework count R = 7** (after r3, three batches R=4/5/6 + the R=7 merged batch; user authorization every round, OB-140/142 on record).

## §4 Environment and tiers

- Local Windows 11 amd64 (non-elevated shell); Ubuntu VM 10.0.0.199 (larson non-root, go1.22.12, `GOPROXY=goproxy.cn`).
- Driver/① = DeepSeek V4 Pro Max; ②③ = DeepSeek V4.1 Flash High; ⑤ = DeepSeek V4 Pro;
  ⑥ = MAI Code 2 Flash (+1 round OB-140); ⑦ = GPT-6.2 Sol Max (1+3) + four reads (§5);
  targeted mini-read = Gemini 4 Carbon High (OB-143). Probes = 0.

## §5 Review conclusions summary (⑤ ⑥ ⑦)

- ⑤ pre-review r1 FINDINGS → r2 PASS; ⑥ scan r3/r4; ⑦ blind review FINDINGS (22 items) → final review FINDINGS (3H+2M) →
  re-review r1/r2/r3 FINDINGS → R=4/R=5/R=6 → **four-read same-input matrix** → R=7 → targeted mini-read.
- **Four-read matrix (same input = R6 batch; recorded by "verdict divergence", n=1 single-input observation, no model-capability ranking)**:

| Read | Conclusion | Group1/Group2/Group3 | T2/T3/T8a |
|---|---|---|---|
| Gemini 4 Carbon | PASS WITH FIXES | C/C/C | C/C/C |
| Gemini 4 Neon | FINDINGS (4M+2L) | N/C/C | N/N/C |
| GPT-6.2 Sol | FINDINGS (1M+3L) | N/N/C | N/N/C |
| GPT-5.3-Codex (blind) | PASS | C/C/C | C/C/C |

  Verdict-divergence rows (P87 data points): Group 1 Carbon/Codex CLOSED vs Neon/Sol NOT-CLOSED; Group 2 Carbon/Neon/Codex
  CLOSED vs Sol NOT-CLOSED; T2/T3 Carbon/Codex CLOSED vs Neon/Sol NOT-CLOSED. Core observation = under identical input,
  PASS and FINDINGS coexist (verdict-philosophy divergence: fail-closed fallback vs contract completeness); the verdict-standard consensus and dispatch policy
  are ruled by the P87 governance slice. P87 raw data archived at `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/P87-four-read-data/`.
- **R=7 merged batch**: Neon 4M+2L + Sol 3L + user final ruling on the S1 boundary (§8); Codex new findings = 0.
- **Targeted mini-read (⑧c input)**: `TARGETED REVIEW: PASS WITH FIXES` — 9 items all CLOSED, 0 Medium+, 4 Low non-blocking
  (L1 p85_t3_windows_f04 no longer enters the ancestor-refusal branch after passing nil leaf; production is covered by VerifyAIInstallSource;
  L2 nominal entry uppercase suffix role=-1 (non-production entry only); L3 aclErr drops probeErr without Join (diagnostic only); L4 comment wording (fixed, §1 E-4)).
  L1/L2/L3 accepted-unfixed by user ruling, recorded faithfully in this report.

## §6 ⑧a self-check (four items, present/absent + evidence paths)

| # | Item | Conclusion | Evidence path |
|---|---|---|---|
| 1 | Mechanical audit of the cost table | Present | §14 cost table; credits user-registered accounting (Sol 67.9/14m50s, Codex 129.9/26m58s, four-read total ≈958) |
| 2 | "Disposed" file:line | Present | `docs/DELIVERY_DIRECTIVE.md` L1184–L1187 (OB-140/141/142/143); CONTRACTS :214/:216/:220/:231/:233 |
| 3 | Old-wording grep | Present | Command and results = `05-input/old-wording-grep.txt` (only the new-wording "not an F04 exemption" sentence and three legitimate F04 registrations in `bounddir_windows.go` — spawn final name / artifact-open final name / directory flush, none on publication paths) |
| 4 | ⑥/⑦ every-round Info/Low all in the report | Present | See §5: Carbon 3L + Neon 2L + Sol 3L + targeted-mini-read 4L all listed; ⑥ Info see R4–R6 records |

## §7 A-NN numbering integrity

- The 22 blind-review items F01–F22 are disposed of one by one with records (all landed before ⑨); R-rounds MUT-1..MUT-5; OB-140/141/142/143; no new DRs; numbering continuous with no gaps.

## §8 Independent registration of the Windows publication-path S1 boundary (user final ruling 2026-11-23)

- **Boundary fact**: on Windows, `os.CreateTemp` + `os.Link` cannot be executed relative to a directory handle; publication writes resolve via user-mode paths;
  the before/after root-chain re-bind cross-check (`OpenBoundDirVerified` twice) runs before write and after write, but if chain insertion/replacement happens between the two cross-checks,
  the link has already happened and the post-hoc refusal cannot be withdrawn. The no-replace hard link provides atomic visibility, **not power-loss durability**.
- **Registration accounting**: this entry is the **S1 user-mode path publication boundary** (user final ruling 2026-11-23, independently registered), **not** an F04 exemption;
  5 code comments have been re-registered under this accounting (`agent_runner_policy_publish_other.go:14-20`, `ai_install_proof_publish_other.go:13-19`,
  `ai_install_proof.go:126/258` etc.), the three existing F04 mentions in `bounddir_windows.go` (spawn final name / artifact-open final name /
  directory flush) are **other already-registered residuals**, not part of this boundary.
- **Follow-up**: if S2 introduces handle-relative create/link or hands writes to a managed executor, this boundary is re-evaluated.

## §9 Review input package summary (including blind-review isolation proof)

- ⑦ inputs for every round: frozen diff (diff.patch 776KB) + bilingual contracts + dual-platform readouts; **the first blind-review round had no review checklist attached** (isolation holds);
  the four-read same input = R6 batch (three r3 items + three fix groups + the contract-slot sentence + dual-platform readouts); Codex is a blind read (forbidden to read the first three legs' reports, declared).
- R=7 targeted-mini-read input = only `r7-fix.diff` (22 changed files, 78KB) + R7 kickoff message + corresponding tests, no full re-review.

## §10 Falsifiability (mechanical checks + mutation tests)

- **Mutation matrix (all red → byte-for-byte restore → green)**: MUT-1 slot order reverts to path order (R6); MUT-2 Windows sandwich (R6);
  MUT-3 carrier-publication re-check reverts FileIDAt (R6); **MUT-4A** fd→name-based ACL (red = `TestP85R7LinuxDirectorySlotReadsTheACLFromItsDescriptor`,
  VM non-root); **MUT-4B** deletes the closed-name comparison (ruling-specified; red = `TestP85R7LinuxDirectorySlotClosingComparisonRefusesAMovedName`,
  VM non-root); **MUT-5** `OpenBoundDirVerified`→`OpenBoundDir` follows the chain (red = `TestP85R7OpenInstallSourceLeafRefusesALinkedChain`,
  Windows non-elevated). All restores are SHA-256 byte-for-byte identical (`05-input/mutation-r7-executed.txt`, `mutation-r6.txt`).
- **Counterexamples**: missing/malformed/shape-mismatched carrier refused; cross-role same-name claim refused; Role/suffix contradiction and case mismatch refused; junction/symlink chain binding refused.

## §11 Gate results

- `node tools/gates/gate.js --all --skip-base --scope=index|tree|ci`: **TOTAL_FAIL=0 × 3**; `--selftest` **PASS 740/740**.
- Three-variant cross-build `GOOS=windows/linux/darwin GOARCH=amd64 go build ./...` all green; three-view `go vet` all green.
- Three-view `gopls check -severity=hint` (411 source files under internal/**/*.go) filtered `unused` = **0 lines × 3**.
- Staging hygiene: porcelain all `M/A` precisely staged, MM/AM/?? = 0; nothing committed (awaiting ⑩ authorization).

## §12 ⑨ native verification results (both ends)

- **Windows local (non-elevated)**: `go build ./...` exit=0; `go test ./...` 15 packages ok / FAIL=0, exit=0;
  `go vet ./...` exit=0. Per leg (R7 new cases): `OpenInstallSourceLeafRefusesALinkedChain` PASS (junction leg),
  `UnobservableSlotKeepsItsPositionInTheDigest` PASS, `RoleSuffixContradictionIsRefused` 8/8 PASS; no SKIP.
  Raw logs = `05-input/local-full-log.txt`, `local-r7-legs.txt`.
- **Ubuntu VM (larson non-root, go1.22.12)**: `go build ./...` exit=0; `go test ./...` 15 packages ok / FAIL=0, exit=0;
  `go vet ./...` exit=0. Per leg: directory-slot three legs (`ReadsTheACLFromItsDescriptor` / `RefusesAReplacedName` /
  `ClosingComparisonRefusesAMovedName`) PASS; `OpenInstallSourceLeafRefusesALinkedChain` PASS (symlink leg);
  `UnobservableSlotKeepsItsPositionInTheDigest` PASS; `RoleSuffixContradictionIsRefused` 8/8 PASS.
  Raw logs = `05-input/vm-full-log.txt`, `vm-r7-legs.txt`.
- **F17 privileged fixture (four elements registered faithfully)**: adjudicator = user; date = 2026-11-23; scope = whole slice; reason = the F17 exemption covers only same-inode
  privileged legs — this slice did not run privileged fixtures requiring chown/SeTakeOwnership; the relevant positive controls were replaced by non-elevated/non-root native legs and recorded leg by leg;
  basis = ⑦ final-review accounting + user ruling. **Not executed is annotated as such; not counted as closed.**
- **CI dual legs (settled)**: main commit `cc7d90f` first CI 45201921016: ubuntu ✓, windows leg 2 environmental failures (§1 E-8) →
  R=8 test-side fix (`6480bb1`) → re-run CI **45208001981 dual legs green**: windows 4m55s / ubuntu 3m25s;
  Gates step **TOTAL_FAIL=0 × 2**; `adapters`/`console`/`guard` package-level `ok` lines present on both legs (plus ubuntu Race three packages ok);
  Test step `--- FAIL` = 0. Raw evidence = `05-input/ci-45208001981-evidence.txt` (archived in this package).
  **F17/F22 settlement**: dual-leg Gates + package-level ok hard evidence complete ⇒ both evidence debts settled per ⑨ accounting; the no-evidence-no-close condition closed.

## §13 ⑧c closeout cleanliness record

- Repo-wide `gofmt -l` zero output; three-view gopls unused cleared; encoding gate: all changed files BOM/EOL fully checked
  (md=BOM+LF, go/json/txt=noBOM+LF); `tmp/` root only `.gitkeep` + existing slice directories; untracked files 0 (all precisely staged).
- Evidence-package self-consistency: the P85 evidence directory MANIFEST and SHA256SUMS are in sync (§9 payload; the frozen package is excluded from per-file encoding rewrites).

## §14 Cost and metering (user-registered accounting)

- ⑦ statutory quota: final review 1 + re-review 3 (F-01 user ruling) + blind review 1 (F-03) — 5 Sol Max calls in total.
- Additions: four reads (Carbon / Neon / Sol 5th read / Codex) + targeted mini-read (Carbon High, OB-143).
- Credits: Sol baseline read 67.9 / 14m50s; Codex 129.9 / 26m58s; four-read total ≈ 958 (the remaining reads and the two ② Flash job credits are user-registered).
- Probes = 0; ⑥ = MAI Code 2 Flash (quota 3 + OB-140 addition 1).

## §15 Explicitly unexecuted items

| # | Item | Reason and accounting |
|---|---|---|
| 1 | F17 same-inode privileged fixture | Non-elevated/non-root cannot construct it; four-element registration (§12), not counted as closed |
| 2 | P85 dual-leg CI (Gates + package-level ok) | Requires post-commit trigger; backfilled after ⑩ authorization, no evidence no close |
| 3 | L1/L2/L3 (targeted-mini-read Low) | User ruling accepts-unfixed, recorded faithfully (§5) |
| 4 | Windows publication-path handle-relative create/link | User final ruling: S1 independent boundary registration (§8), no code fix |

## §16 No-precedent verification

- OB-140 (⑥ +1) / OB-141 (Codex fourth read) / OB-142 (R=7 batch) / OB-143 (targeted mini-read) are all "must not serve as precedent";
  F-01 ⑦ quota, F-04 driver model, four-read trial, D-2 additional ruling are all item-by-item user authorizations, not precedents;
  this report's four-read matrix is recorded only by verdict divergence (n=1), with no capability-ranking wording.

## §17 Next-step recommendations

1. ⑩ authorize commit (`git add` precise staging ready; commit message `<type>: <summary>` in English) → push origin (origin only, gitee forbidden).
2. Trigger CI dual legs after commit → archive the Gates step and the adapters/console/guard package-level `ok` lines → settle the F17/F22 evidence debts.
3. Backfill the P85 section and the SECURITY note in the ledger `REMAINING_SLICES{,_EN}.md` per existing process once CI evidence lands; the P87 governance-slice scheduling reuses this slice's four-read data.

## Appendix A Execution-freedom list (written back at ⑧b finalization, consistent with `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/freedom-list.md`)

| No. | Freedom | Final ruling |
|---|---|---|
| F-01 | ⑦ quota 1+3 (user ruling) | in-bounds (OB recorded) |
| F-02 | ⑦ = GPT-6.2 Sol Max deep tier | in-bounds |
| F-03 | Blind review must sample 1 | in-bounds |
| F-04 | Driver/① unified V4 Pro Max | in-bounds (user authorization) |
| F-05 | ②③ Flash High no mid-course tier change | in-bounds |
| F-06 | ⑤ V4 Pro 1+1 | in-bounds |
| F-07 | ⑥ ≤3 + OB-140 | in-bounds (addition registered) |
| F-08 | GLM-5.3-Flash ≥2 serial | in-bounds |
| F-09 | Protocol-first | in-bounds |
| F-10 | ⑨ includes Ubuntu VM + dual-platform CI | in-bounds (CI pending post-commit) |
| F-11 | Probes 0 | in-bounds |
| F-12 | No gitee push; English commit messages; ⑩ stop point | in-bounds |
| F-13 (added) | Four reads + targeted mini-read (model-list-expansion-type temporary authorization) | in-bounds (OB-141/143 + trial authorization, must not serve as precedent) |

## Appendix B Six-stage self-check (sixstage)

```sixstage
## Section A Overall conclusion
⑥ two rounds (r3 first scan + r4 re-scan; quota 2/4, one extra round via OB-140): r3 = `FINDINGS` (1 High + 4 Medium, the output opened with a bare `FINDINGS` token and was resent under §7.2.2 as a format deviation); r4 = `FINDINGS` (11 NOT-CLOSED items + platform evidence gaps, Section B in bullet form); per the §7.2.2 failure handling, the flow entered the ⑦ independent final review and closed out there.

## Section B Findings list

| SEVERITY | Location | Disposition |
|---|---|---|
| High | ⑥ r3 P85-F01 (no second artifact re-read after the Publish entry re-read) | re-scanned by r4 after fixes; closed with the R=4–R=7 batch |
| Medium | ⑥ r3 P85-F02/F03/F04/F05 (4 items) | re-scanned by r4 after fixes; closed item by item by ⑦ |
| — | ⑥ r4 11 NOT-CLOSED items (F01/F02/F04/F05/F08/F10/F16/F17/F18/F20/F22) | closed item by item via the ⑦ final/re-reviews + the four reads + the targeted mini-read (§5) |
| — | ⑥ r4 11 CLOSED items (F03/F06/F07/F09/F11/F12/F13/F14/F15/F19/F21) | verified CLOSED item by item |

## Section C Evidence and reasoning per item
Evidence chain: workspace HEAD `e35e50f`, gates index/tree both `TOTAL_FAIL=0`, SELFTEST `740/740`, contract tests 4/4 PASS, `go build ./...` and `go vet ./...` successful with no output. The P85 targeted T2, T8a and Windows F04/F05 tests all PASS; the full Go suite showed one transient P84 C7 Windows lease `Access is denied` on the first run and PASS on the standalone rerun.

## Section D Falsifiability audit
Falsification points include: the same UID changing the policy file and the anchor at once, a directory replaced with a junction/symlink, a pre-launch replacement with a same-digest different inode, a reparse switch before Windows process creation, named-ACL write grants, re-reading the non-persisted directory after a Windows crash, running the shared C:\ carrier test on Linux, running the same-inode rewrite under a privileged fixture, and not substituting gate-ci when the P85 CI is missing.

## Section E What I cannot verify
Cannot verify: the Ubuntu VM copy is an old/dirty baseline `dc785ed` without the latest P85 sources, so it is not comparable native Linux evidence; the package-level `ok` lines of the P85 dual-platform CI were not obtained; the independent mutations were not run from source on disk — the above conclusions rest on the static diff, the current Windows measurements and the existing tests.

INDEPENDENT SCAN: FINDINGS
```

## Appendix C Artifact fence (artifacts)

```artifacts
docs/validation/security-proof-subsystem.md	repo
docs/validation/security-proof-subsystem_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/report.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/MANIFEST.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/SHA256SUMS.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/freedom-list.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/input-list.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/00-briefing.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/P87-four-read-data/07-three-read-comparison.md	repo
```
