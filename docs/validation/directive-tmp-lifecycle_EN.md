# Slice validation report - `DIRECTIVE-TMP-LIFECYCLE` (v1.23)

> Type: directive-governance (executing the `tmp/` lifecycle: archiving + reference annotation + ledger registration). Tier `[SLICE]`. Size M. Probes 0.
> **Independence note**: this slice ⑦ is carried by `deepseek-v4-pro` (ADR-023, a direct user ruling) and shares the model family with the ⑤ pre-review ⇒ **reduced independence**, stated truthfully and handed to ⑧c.

## 1. Change summary

- The `tmp/` stock is classified into 15 archive units under `docs/validation/evidence/directive-history/` per the class ② caliber of §1.5.1; each unit carries a `MANIFEST.md` (archived entries / exclusion list / untracked-file list / rebuild notes / citing references / suffix-and-encoding mapping) and a `SHA256SUMS.txt`, with the payload root at `archive/`.
- Reference annotation: a v1.23 archiving note is appended after each of the three existing notes in the directive (pointing at the actual `.txt` paths); 9 report pairs each gain one archiving note line; no historical wording is rewritten (§11.1).
- Ledger and version bundle: the v1.23 bundle (header version / new §0.3 row / appendix C.0m), §12.4 registers OB-76 / OB-77 / OB-78 / OB-79 / OB-80, `ADR-023`, and the completion row (recording the report path only).
- Same batch: the note wording of `DIRECTIVE-MODELID-MAP-mutation.md` becomes "**rebuilt archive**" (user instruction); the process-file count of `ledger-backfill{,_EN}` is annotated to 21.
- Freedom list: `docs/validation/evidence/DIRECTIVE-TMP-LIFECYCLE-freedom-list.md` (F-1 to F-13).

## 2. ① judgements and calibers

- **Archiving caliber** = evidence pack: files whose paths version-controlled documents cite, plus slice-level documents (slice definitions / designs / plans / review packs / gate readings / frozen patches); runtime residue and one-off scripts go only into the exclusion list.
- `dce` and `assembly-min` are **archived per the class ② definition** (the class ③ illustrations in §1.5.1 contradict the definition; registered as OB-78).
- **Binaries and ignored types are not committed**: the unit MANIFEST records path + bytes + sha256 + rebuild notes, and the report notes point back at the MANIFEST so no dangling reference is left.
- **Suffix and encoding** (option 3): archived copies of `.md` are stored as `.txt` (avoiding the document criteria that scan `.md` only) and encoded per archived type - `.ps1` keeps a BOM, everything else carries none, all LF; each unit records "original path / archived path / original-byte hash / archived-byte hash / change kind" per item.
- **Boundary declines**: `dr-fix-enforcement-proxy{,_EN}` receives no note (that report has no `artifacts` fence, so any change would trip G7-a); existing frozen artifacts inside the evidence tree are untouched.

## 3. Delivery surface (archive units)

| Unit | Archived | Excluded | Untracked |

| `assembly-min` | 12 / 659,581 B | 72 / 0 B | 0 / 0 B |
| `b4` | 50 / 396,972 B | 300 / 0 B | 3 / 0 B |
| `dce` | 11 / 215,448 B | 6 / 0 B | 0 / 0 B |
| `dr1` | 14 / 137,253 B | 17 / 0 B | 0 / 0 B |
| `drfix` | 19 / 335,863 B | 19 / 0 B | 0 / 0 B |
| `gate-js` | 1 / 8,931 B | 0 / 0 B | 0 / 0 B |
| `gates-ext` | 31 / 671,077 B | 769 / 0 B | 2 / 0 B |
| `lbf` | 12 / 273,553 B | 9 / 0 B | 0 / 0 B |
| `master-subagent-architecture` | 1 / 26,124 B | 0 / 0 B | 0 / 0 B |
| `pre-directive-package` | 1 / 15,483 B | 0 / 0 B | 0 / 0 B |
| `reviews` | 6 / 100,926 B | 0 / 0 B | 0 / 0 B |
| `s1-reach` | 8 / 177,101 B | 6 / 0 B | 0 / 0 B |
| `tmp-lifecycle-workdir` | 15 / 744,747 B | 41 / 0 B | 0 / 0 B |
| `v112` | 22 / 450,619 B | 18 / 0 B | 0 / 0 B |
| `v113` | 76 / 2,602,143 B | 70 / 0 B | 0 / 0 B |
| **Total** | **279 / 6,815,821 B** | **1327 / 0 B** | **5 / 0 B** |

## 4. Pipeline and environment

① judgement (driver + V4 Pro) ⇒ ②③ archiving and reference annotation (driver scripts, deterministic and recomputable) ⇒ ④ gates ⇒ ⑤ pre-review (V4 Pro) ⇒ ⑥ independent scan (MAI) ⇒ ⑦ final review (V4 Pro, ADR-023) ⇒ remediation ⇒ ⑧ this report pair ⇒ ⑨ native verification ⇒ ⑧c close-out cleanup ⇒ ⑩ stopping point.

Environment: Windows 11 / go1.27.0 / Node 24; repository `github.com/larsonzh/prfrail`, branch `main`.

## 5. Review conclusions

- **⑤ pre-review** (V4 Pro): `PASS WITH FIXES`. 1 High = the report pair had not landed yet (closed by landing it in this batch); 1 Medium = OB-79 was approved but had no ledger row (registered); 2 Low = `tmp/tl` drift and the pending F-6 acknowledgement.
- **⑥ independent scan** (MAI): `PASS`. Readings: hashes `279 / 0`, coverage `279 / 0`, cited archive paths `8 / 0`, ignored types committed `0`.
- **⑦ final review** (V4 Pro, ADR-023): `PASS WITH FIXES`. Three Low items remediated: the OB-79 count added to three places, the source-path and rebuild-note wording in unit MANIFESTs, and a bracketed clause in the report notes covering paths cleared by earlier slices.
- Of the six counter-examples in ⑦, two held (the pre-existing `dml` dead links and the missing OB ledger row) and both were closed by the remediation; the other four were measured not to hold.

## 6. Falsifiability and gate results

| Item | Reading |
|---|---|
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0` (changed=337) |
| `gate.js --all --scope=index` | `TOTAL_FAIL=0` (changed=337) |
| `--selftest` | `SELFTEST: PASS 264/264` |
| Archive self-consistency | 15 units / 279 payload items; per-item `SHA256SUMS` recomputation found 0 mismatches, 0 unlisted, 0 missing |
| Archive fidelity | After normalization (BOM removed, LF) the `.md` → `.txt` copies are byte-identical to the `tmp/` originals, with zero character differences |
| G4b | `newCjk=504 unknown=[]` |
| G7-b | `316 artifact row(s) OK` |
| Three-table completeness | `b4` 353 = 50 + 300 + 3; `gates-ext` 802 = 31 + 769 + 2 |
| Mutation checks | Not applicable (a documentation and archiving slice with no code semantics); the falsification means are ⑥ six mechanical checks and ⑦ six counter-examples |

**Known deviation (registered truthfully)**: the archived copies' BOM form is normalized by archived type (`.ps1` with BOM, everything else without), which differs from the user's literal "keep the BOM form" ⇒ freedom-list F-5 and F-6 were **acknowledged by the user** (2026-10-04) and **OB-80** was registered (the archive-payload encoding form has no rule).

## 7. Cost and accounting

① one round, ⑤ one round, ⑥ one round, ⑦ one round, ⑧ one round (the CN authoritative draft and its `_EN` mirror produced by the driver in one pass); probes 0; no Haiku calls; **no quota exception** (the ⑦ downgrade is traced in ADR-023).

## 8. Explicitly not executed

- No criterion body was changed (no criterion-scope exclusion was added for the evidence tree) ⇒ OB-79 is registered for a dedicated slice.
- No `tmp/` file was deleted (the class ③ cleanup runs at ⑧c, see the appendix).
- Not committed, not pushed, not pushed to gitee.

## 9. Next-step recommendations (including tier advice)

- **Next slice (recommended)**: criterion scope and the evidence tree (OB-79) - excluding the evidence tree from the document criteria and letting frozen-pack self-consistency carry the check; a criterion-body change ⇒ size **M+**, ⑦ uses `gpt-5.3-codex` **without downgrade**.
- Handling OB-76 (the fourth disposition class for rebuilt archives), OB-78 (§1.5.1 illustration fix) and OB-80 (archive-payload encoding form) in the same batch could merge with the item above.

## 10. ⑩ stopping-point report (including the write-back check line)

- Executed through ⑨ and ⑧c; **stopped before commit authorization**: this slice does not commit, does not push, does not push to gitee.
- **Write-back check** (the OB-63 formulation of §11.2): validation report = **written back** (this file); `DEV_PLAN` = **not applicable** (a directive-governance slice); ledger = **written back** (the `REMAINING_SLICES{,_EN}` completion row, recording the report path only); directive = **written back** (§0.3 / §12.4 / appendix C.0m / the three notes); ADR = **written back** (ADR-023).
- **⑨ readings**: both tree and index report `TOTAL_FAIL=0`; `SELFTEST: PASS 264/264`; archive self-consistency recomputation shows 0 mismatches; after the ⑧c cleanup a re-run of the non-archive surface stays green.
- **Commit and CI**: commit `4284e94`; CI run `37209843851` is green on both legs (`Go ubuntu-latest` / `Go windows-latest`), with the three `workflow_dispatch` legs skipped by their event condition; both legs agree: `SELFTEST: PASS 264/264`, `GATE REPORT scope=ci changed=340`, `TOTAL_FAIL=0`.

## 11. Artifacts list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-tmp-lifecycle.md	repo
docs/validation/directive-tmp-lifecycle_EN.md	repo
docs/validation/evidence/DIRECTIVE-TMP-LIFECYCLE-freedom-list.md	repo
docs/validation/evidence/DIRECTIVE-MODELID-MAP-mutation.md	repo
docs/validation/ledger-backfill.md	repo
docs/validation/ledger-backfill_EN.md	repo
docs/validation/directive-cleanup-eval.md	repo
docs/validation/directive-cleanup-eval_EN.md	repo
docs/validation/directive-review-saturation.md	repo
docs/validation/directive-review-saturation_EN.md	repo
docs/validation/directive-v1.12.md	repo
docs/validation/directive-v1.12_EN.md	repo
docs/validation/dr-1-availability-probe-floor.md	repo
docs/validation/dr-1-availability-probe-floor_EN.md	repo
docs/validation/gates-ext.md	repo
docs/validation/gates-ext_EN.md	repo
docs/validation/s1-reach-reachability.md	repo
docs/validation/s1-reach-reachability_EN.md	repo
docs/validation/t027-assembly-min.md	repo
docs/validation/t027-assembly-min_EN.md	repo
docs/validation/t027-at23-e2e.md	repo
docs/validation/t027-at23-e2e_EN.md	repo
docs/validation/evidence/directive-history/tmp-lifecycle-workdir/CLEARED-INVENTORY.txt	repo
docs/validation/evidence/directive-history/assembly-min/MANIFEST.md	repo
docs/validation/evidence/directive-history/b4/MANIFEST.md	repo
docs/validation/evidence/directive-history/dce/MANIFEST.md	repo
docs/validation/evidence/directive-history/dr1/MANIFEST.md	repo
docs/validation/evidence/directive-history/drfix/MANIFEST.md	repo
docs/validation/evidence/directive-history/gate-js/MANIFEST.md	repo
docs/validation/evidence/directive-history/gates-ext/MANIFEST.md	repo
docs/validation/evidence/directive-history/lbf/MANIFEST.md	repo
docs/validation/evidence/directive-history/master-subagent-architecture/MANIFEST.md	repo
docs/validation/evidence/directive-history/pre-directive-package/MANIFEST.md	repo
docs/validation/evidence/directive-history/reviews/MANIFEST.md	repo
docs/validation/evidence/directive-history/s1-reach/MANIFEST.md	repo
docs/validation/evidence/directive-history/tmp-lifecycle-workdir/MANIFEST.md	repo
docs/validation/evidence/directive-history/v112/MANIFEST.md	repo
docs/validation/evidence/directive-history/v113/MANIFEST.md	repo
```

Note: archive units are declared by **unit MANIFEST rows** (per-item payload details live in each unit `SHA256SUMS.txt`, 279 items in total); no `local-only` rows are declared separately (the class ③ cleanup ran at ⑧c, see the appendix).

## 12. Appendix · class ③ cleanup record

- **Cleanup-time inventory artifact**: `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/CLEARED-INVENTORY.txt` - per-item path, bytes, sha256 and the not-archived reason (the full set at cleanup time, including process files added after the MANIFEST freeze).
- Cleanup scope = this slice's exclusion lists (1327 items) plus the process files outside the workdir readings; after cleanup the `tmp/` top level holds only `.gitkeep`.
- After cleanup both `--scope=index` and `--scope=tree` are re-run (readings in §10); the uncited process artifacts are zeroed per the §1.5 "removed immediately after use" rule.
- **Footnote (inventory caliber)**: temporary process files created after the class ③ cleanup point are removed by the driver immediately after each step and are **not counted** in the cleanup-time set of `CLEARED-INVENTORY.txt`.
