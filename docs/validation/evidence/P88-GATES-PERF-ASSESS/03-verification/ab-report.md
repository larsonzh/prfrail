# P88 GATES-PERF-ASSESS / (3) tester report - gates output equivalence + timing

Encoding of this file: **pure ASCII, LF only** (no BOM, no non-ASCII byte; ASCII makes the
BOM question moot). Verified with `[IO.File]::ReadAllBytes` in `_verify-encoding` step below.

Role: (3) tester for slice `P88-GATES-PERF-ASSESS`.
Claim under test: the O-1 (batch blob prefetch) + O-3c (read memo) optimisation of
`tools/gates/lib/ctx.js` changes **no** gate-criterion output, and makes the cold corpus path
materially faster.

Nothing was staged, stashed or committed; no product file, `.md`, `internal/**` or
`.github/**` file was written by this role. All artefacts live under
`tmp/P88-GATES-PERF-ASSESS/03-ab/` (gitignored).

---

## 0. Preconditions (verified BEFORE any measurement)

| check | result |
|---|---|
| HEAD | `fed0f5c4daac03c78b7ba87ad0770a0813730783` |
| node | `v24.17.0` |
| `gate.js` old vs new | sha256 identical (`0d43b5cc...e2e927`) - the CLI is untouched |
| old snapshot integrity (`03-ab/integrity.txt`) | 100/100 files byte-identical to `git cat-file blob HEAD:tools/gates/...`, 0 mismatch, 0 extra, 0 missing |
| working tree vs HEAD | exactly **1** differing product file: `lib/ctx.js` (`5525eb46a8f9` -> `62ecdd40e6f7`); every other file identical |
| `tools/gates/lib/ctx.js` sha256 (start = end) | `62ecdd40e6f74a0787a2efbfb41ffcdbb5c633b871151426f08265c581352261` |
| G4b cache dir | `tmp/gates-cache/g4b`, 100 JSON entries (LRU full), 0 non-JSON files |
| cache key, old vs new (`03-ab/precheck.txt`) | `keysEqual=true` for all four scopes: tree `64f622431f22`, index `b2b542360e86`, range `a337ec35c303`, ci `59c285d4cdad` |
| `git status --porcelain` | unchanged by this role (see section 7) |

The key precheck matters: it proves **before** any run that the cache key is not affected
by the optimisation (it cannot be - `sourceDigest()` hashes only `g4b.js` and the whitelist).
It also tells us the cache state each cell will see, which is what makes the pairing sound.

---

## 1. Execution matrix

Redirection form (to the letter): `cmd /c "node <script> <args> > <file> 2>&1"`, realised by
generating one `.cmd` file per cell (`03-ab/_cmd/<cell>.cmd`) so PowerShell never re-encodes
the stream and quoting cannot drift. Captures therefore contain stdout+stderr merged; section 6
shows stderr was empty in every cell.

| cell | impl | command (relative to repo root) | exit | elapsed ms | sha256 | bytes | G4b cache line | same as peer |
|---|---|---|---|---|---|---|---|---|
| A1 | old | `node tmp\P88-GATES-PERF-ASSESS\orig-gates\gate.js --all --scope=tree` | 1 | 75000 | `11562a7f5c0f4f6411f9ed8ca719639783bf6a47c847c8a189924761ea767772` | 4645 | `hit 64f622431f22 corpus=1466` | A2 YES |
| A2 | new | `node tools\gates\gate.js --all --scope=tree` | 1 | 79628 | `11562a7f...ea767772` | 4645 | `hit 64f622431f22 corpus=1466` | A1 YES |
| B1 | old | `... --all --skip-base --scope=index` | 1 | 2430 | `3fb9ba1c3bdf6a7b16259060b5e316cf654515fecd56648ac9da4e8683b15d8e` | 4147 | `hit b2b542360e86 corpus=1466` | B2 YES |
| B2 | new | `... --all --skip-base --scope=index` | 1 | 1750 | `3fb9ba1c...b15d8e` | 4147 | `hit b2b542360e86 corpus=1466` | B1 YES |
| B3 | old | `... --all --skip-base --scope=range:HEAD~3..HEAD` | 1 | 59272 | `b6fb364ee4d4a88809d8a6f6707c9c4d30438e1fdb92439db49d1ce65baf5f33` | 7061 | `hit a337ec35c303 corpus=1466` | B4 YES |
| B4 | new | `... --all --skip-base --scope=range:HEAD~3..HEAD` | 1 | 33083 | `b6fb364e...5f33` | 7061 | `hit a337ec35c303 corpus=1466` | B3 YES |
| C1 | old | `... --all --skip-base --scope=ci` (cache cleared => COLD) | 0 | 170277 | `a95474da6707381178034c9866a1802005268076127ec47d53450ea0fdb72f2a` | 4905 | `miss 59c285d4cdad corpus=1466` | C2 YES |
| C1w | new | same as C1, cache **warm** (entry written by C1) | 0 | 10137 | `b5d8174cbed09b3e99aa82b36d4f60568fcc901d8d5428d6bab5437cd2bcf181` | 4904 | `hit 59c285d4cdad corpus=1466` | differs from C1 by hit/miss only |
| C2 | new | same as C1 (cache cleared again => COLD) | 0 | 12929 | `a95474da...b72f2a` | 4905 | `miss 59c285d4cdad corpus=1466` | C1 YES |
| D1 | old | `node ... --selftest` | 0 | 209992 | `36af71e52db1dfca0efa844cffe75fa53fa1f1ea8cdae25424c9f3f053bd9d2f` | 72726 | n/a (`SELFTEST: PASS 744/744`) | D2 YES |
| D2 | new | `node tools\gates\gate.js --selftest` | 0 | 128881 | `36af71e5...bd2f` | 72726 | n/a (`SELFTEST: PASS 744/744`) | D1 YES |

Machined source of the table: `03-ab/matrix.tsv`; per-file hashes: `03-ab/sha256.txt`.

Byte comparison was done twice per pair, independently: sha256 equality plus `fc /b`
(exit 0 = no differing byte). No whitespace/EOL-insensitive mode was used anywhere.

### Pre-existing FAIL that both sides share (not caused by the optimisation)

A1/A2/B1..B4 exit 1 because of `G8-a evidence-tree path shape` on this slice's own untracked
evidence directory (`docs/validation/evidence/P88-GATES-PERF-ASSESS/*` is not yet a registered
evidence pack). The old (== HEAD) implementation reports the identical record, so this is a
property of the in-flight slice, not a regression. C1/C2 (`--scope=ci`) show `TOTAL_FAIL=0`.

---

## 2. Pairwise verdict (byte-exact)

| assertion | evidence | verdict |
|---|---|---|
| A1 == A2 | sha `11562a7f...`, `fc /b` exit 0 | **IDENTICAL** |
| B1 == B2 | sha `3fb9ba1c...`, `fc /b` exit 0 | **IDENTICAL** |
| B3 == B4 | sha `b6fb364e...`, `fc /b` exit 0 | **IDENTICAL** |
| C1 == C2 | sha `a95474da...`, `fc /b` exit 0 | **IDENTICAL** (incl. the `INFO G4b-cache ... miss 59c285d4cdad corpus=1466` line, i.e. the cache key did not move) |
| D1 == D2 | sha `36af71e5...`, `fc /b` exit 0, both `SELFTEST: PASS 744/744`, both exit 0 | **IDENTICAL** |

C1 vs C1w differ in exactly one line (`miss` vs `hit`, same key12, same `corpus=1466`); that is
the *expected* difference and it doubles as the cross-implementation cache-compatibility proof:
the new reader consumes an entry the old writer produced, byte-identically.

No inconsistency was found, so the "optimisation must be rolled back" branch does not apply.

---

## 3. Timing (old vs new, same HEAD, same machine, one run per cell)

| path (what it measures) | old ms | new ms | speed-up (old/new) | saved ms |
|---|---|---|---|---|
| A `--all --scope=tree` (warm G4b + 4 base gates) | 75000 | 79628 | 0.94x | -4628 (see note) |
| B `--all --skip-base --scope=index` (warm) | 2430 | 1750 | 1.39x | 680 |
| B `--all --skip-base --scope=range:HEAD~3..HEAD` (warm) | 59272 | 33083 | 1.79x | 26189 |
| C `--all --skip-base --scope=ci` (**cold corpus**) | 170277 | 12929 | **13.17x** | 157348 |
| D `--selftest` | 209992 | 128881 | 1.63x | 81111 |
| (extra) C1w new, warm ci | - | 10137 | - | - |

Note on A: the tree run is dominated by the four base gates (`gofmt`, `go build ./...`,
`go vet ./...`, `go test ./...`) and the G4b corpus is a cache **hit**, so this row carries no
optimisation upside; the -6% is single-sample noise on a build/test-dominated workload. Both
sides produce byte-identical output, so the difference cannot be a semantic one.

CI critical path (`--selftest && --all --scope=ci`, cold runner, warm cache):

| | old | new |
|---|---|---|
| selftest | 209992 ms | 128881 ms |
| cold ci gate | 170277 ms | 12929 ms |
| **total** | **380269 ms (380.3 s)** | **141810 ms (141.8 s)** |

=> **2.68x** on the CI gate step, **238.5 s saved (62.7%)** per run, of which 157.3 s comes
from the cold corpus path (O-1) and 81.1 s from the selftest path (O-3c memo, repeated reads in
the criteria the suite exercises).

Timing discipline: one sample per cell, no repetition (per instruction). Variance between the
two sides of A/B is therefore not separable from noise; the **correctness** verdicts do not
depend on timing at all.

---

## 4. Coverage matrix

| # | test point (from the dispatch) | cell(s) | type | verdict |
|---|---|---|---|---|
| T1 | Full `--all` output identical (base four-piece + criteria lines) | A1/A2 | positive (warm cache) | PASS |
| T2 | Criteria-only output identical, two different scopes | B1/B2 (index), B3/B4 (range) | positive (warm cache, two key namespaces) | PASS |
| T3 | Cold path (`miss` branch) output identical, including the cache-key line | C1/C2 | boundary (empty cache, 1146-file corpus rescan) | PASS |
| T4 | New implementation reads an entry the old one wrote | C1 -> C1w | positive (cache compatibility) | PASS (only hit/miss differs) |
| T5 | Selftest parity | D1/D2 | positive (744 suite) | PASS |
| T6 | Cheap-flag / usage surface unaffected | implicit in every cell (exit codes, headers); `gate.js` sha identical, `$Args`-bug cells in `_invalidated-run1/` printed the same usage text from both impls | negative-ish (sanity) | PASS |

Falsifiability: each PASS is refuted by a single differing byte. Demonstrated mechanically in
`03-ab/comparator-sanity.txt`: flipping **one bit** of a copy of `A2.out.txt` makes `fc /b`
against `A1.out.txt` exit 1, so the comparator does discriminate. If O-1/O-3c had changed a
single character of any criterion's text, or any `corpus=` count, the corresponding pair would
go red.

Explicitly **not** covered (cannot be covered in this session, stated rather than implied):

1. Non-Windows / non-node-24 leg (CI runner node version is not pinned; `spawnSync` encoding
   behaviour differs across node majors - the architect notes already flag this).
2. Simultaneous concurrent gate processes sharing one cache key (not asserted by this slice;
   the change touches only reads, but this is a design claim, not a measured one).
3. Corpora above the `maxBuffer` (2e8) or above `ARG_MAX`-style limits - the fail-open paths
   are exercised only by the implementer's harness (`verify-o1-o3c.js`) and by selftest
   fixtures, not by a real oversized repository.
4. Repeated timing samples / statistical significance (dispatch forbids repetition).

---

## 5. Deviations from the instruction (all documented, none affecting the verdict)

1. **First A/B attempt invalidated by a harness bug.** The driver's `Invoke-Cell` parameter was
   named `$Args`, a PowerShell automatic variable, so every flag was silently dropped and all six
   cells ran as bare `gate.js` (exit 2, 122 bytes of usage text). Quarantined in
   `03-ab/_invalidated-run1/`; `matrix.tsv` was recreated; cache stayed 100 -> 100 in each aborted
   cell. Renamed to `$ArgLine` (comment in `run.ps1` explains why) and re-ran once.
2. **Extra cell C1w** (new implementation, warm ci cache). Not in the matrix; costs ~10 s and
   proves the new reader accepts the old writer's record (item T4). No extra old-side run.
3. **A second cache clear between C1 and C2.** A cold run **writes** the entry, so leaving the
   cache alone would have made C2 print `hit` and the pair would compare `miss` against `hit`.
   Both sides of a cold pair must start from an empty entry set; the backup is restored afterwards
   (100 entries, `diff_lines=0`, 0 non-JSON files - `03-ab/cache-restore.txt`,
   `03-ab/final-state.txt`).
4. **Redirection realised via generated `.cmd` files.** The executed line is literally
   `node <script> <args> > <file> 2>&1` under `cmd /c`; wrapping it in a `.cmd` file only removes
   PowerShell quoting ambiguity. Consequence: captures merge stderr, so stderr emptiness is
   evidenced indirectly in section 6.

---

## 6. stderr / contamination check

`03-ab/comparator-sanity.txt` records the first line of every capture: gate runs start with
`GATE REPORT scope=...` and the selftest captures start with `GATES SELFTEST (in-process; ...)`.
The only occurrences of the literal `gate.js:` in any capture are inside the **test name** of
`T57` in `D1/D2` line 542 (a selftest case that asserts a stderr diagnosis exists), identical on
both sides. No stray diagnostic line appears in any capture, on either side.

---

## 7. Final state (unchanged by this role)

```
$ git status --porcelain
 M docs/t027/REMAINING_SLICES.md
 M docs/t027/REMAINING_SLICES_EN.md
 M tools/gates/lib/ctx.js
?? docs/validation/evidence/P88-GATES-PERF-ASSESS/
```

Same four entries as at the start of this role; `tools/gates/lib/ctx.js` sha256 is unchanged
(`62ecdd40...52261`). Cache: 100 JSON entries, identical name+sha256 set to the backup.
`tmp/` root contains `.gitkeep` and one **pre-existing** foreign file `tmp/find-quotes.js`
(mtime 2026-11-22 00:10, i.e. three days before this run) - this role did not create it and did
not delete it (not ours to remove); flagged here because the repo convention keeps `tmp/` clean.

## 8. Evidence index (`tmp/P88-GATES-PERF-ASSESS/03-ab/`)

| file | content |
|---|---|
| `env.txt` | date, HEAD, node, sha256 of both `gate.js`, cache dir + count |
| `integrity.txt` | old snapshot vs HEAD blob-by-blob verification (100/100) |
| `precheck.txt` | cache key12 per scope for both impls + cached flag |
| `matrix.tsv` | machine table of all 11 cells |
| `A1..B4,C1,C1w,C2,D1,D2.out.txt` | raw captures (stdout+stderr) |
| `sha256.txt` | sha256 + byte size of every capture |
| `cache-restore.txt` | restored cache manifest (name + sha256 per entry) |
| `comparator-sanity.txt` | first lines, `gate.js:` hit count, 1-bit-flip detection proof |
| `final-state.txt` | porcelain, diffstat, product hash, cache/manifest diff, tmp hygiene |
| `run.ps1`, `precheck.js`, `integrity.js` | harness (all under the gitignored slice dir) |
| `_cmd/*.cmd` | the exact `cmd /c "node ... > file 2>&1"` line executed per cell |
| `_invalidated-run1/` | quarantined outputs of the harness-bug attempt |

TEST: DONE
