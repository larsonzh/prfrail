GATES-EXT — ⑥ INDEPENDENT SCAN (MAI) INPUT BUNDLE (sanitized: NO reviewer lists attached)
======================================================================
Repo: d:\LZProjects\prfrail (branch main). Slice: GATES-EXT, the delivery directive v1.9's first `[SLICE]`.
HEAD when this bundle was produced: d9536e0

## What the slice delivers (as written into the directive and the tooling)
1. `docs/DELIVERY_DIRECTIVE.md` §6.3 gains two criteria classes, mirrored in `docs/DELIVERY_DIRECTIVE_EN.md`
   (the pair is a strict positional mirror; 1037 lines each, G3(a.2) shape identical, 0 per-line mismatches):
   - **G6 台账元数据一致性** — G6-1 status tense, G6-2 cost-clause consistency, G6-3 run-id **format**, G6-4 expiring literals.
   - **G7 审查包完整性** — G7-a artifacts block present, G7-b path/disposition resolvable, G7-c finding disposition cell, G7-d stopped-state marker.
   Both judge **added lines only** (same change-line scope discipline as G1-a/G3/G4; historical rows are never blamed).
2. `tools/gates/**` — a new offline gate runner (single entry `gate.js`, `lib/criteria/*.js` per criterion,
   `selftest.js`, `testdata/**`, `g6-whitelist.txt`; the pre-existing `cjk-newwords.txt` table is extended
   in place). It migrates the master's throwaway `tmp/gate.js` (kept in the tree, unmodified, for comparison)
   and fixes its known gaps: no hardcoded repo root, real exit codes (0/1/2), untracked and staged-only files are
   seen, deletions/renames parsed, `git check-ignore` instead of a hardcoded skip, every git call wrapped,
   BOM/CRLF normalised on read. Base **G1** (three-way status consistency) and base **G5** (change-set contract)
   are deliberately **not** implemented and are printed as `SKIP` (never PASS).
   G6-1 accepts the **union** of the completion shapes the directive enumerates
   (`（实得）` / `（**实得**）` / `**实得**` / `已闭合` / `✅ 已完成` / `✅ **已完成**`), each still
   gated by the same-label 40-character window; `selftest.js` asserts every shape fires and that a
   different label inside the window stays clean.
   Scope guards: a `range` endpoint that does not resolve to a commit, a failing `git diff`, a failing
   `git status` (tree and index), an inconclusive `--is-shallow-repository` probe, and a SHALLOW
   `--scope=ci` clone (where the EMPTY_TREE fallback would classify the whole repository as added
   lines) are all usage/environment errors (exit 2). A command failure must never be read as "no
   changes": that collapses into the "empty change set" vacuous-pass path. `createContext` takes an
   optional `spawnGit` runner so each guard is asserted independently by the selftest.
3. `.github/workflows/ci.yml` gains **one** step (no `continue-on-error`, no new job, no new Action):
   `- name: Gates` / `run: node tools/gates/gate.js --selftest && node tools/gates/gate.js --all --scope=ci`,
   both matrix legs. Its `Checkout source` step additionally gains `persist-credentials: false` and
   `fetch-depth: 0` — a full history is what lets `--scope=ci` resolve `HEAD^1` instead of degrading to the
   whole tree. `internal/release/workflow_test.go` (the frozen CI contract) gains the matching step-table entry,
   the `test/Gates` script digest entry and the `test/Checkout source` parameter entry: **three** edits to a
   deliberately frozen file, each one separately authorised by the user as a slice-scoped exception.
4. §12.4 gains observation items **OB-15**, **OB-16**, **OB-17**, **OB-18** and **OB-19**; §12.5 gained the
   infrastructure-gap row.

## Raw evidence you may rely on (all pasted outputs, produced with the working tree as it stands)
- `tmp/gates-ext/review/GATE-all-tree.txt` — `--all --scope=tree` (expect TOTAL_FAIL=0).
- `tmp/gates-ext/review/GATE-selftest.txt` — the in-process selftest (one assertion per line).
- `tmp/gates-ext/review/GATE-range-historical.txt` — the two **real historical** retro-tests (A1 ⇒ G6-2 FAIL on
  `range:b738e6b^..b738e6b`; C1 ⇒ G7-a FAIL on `range:1b01115^..1b01115`) and the fixed side (green).
- `tmp/gates-ext/review/IMPL-tracked.diff` — diff of the tracked modifications.
- `tmp/gates-ext/review/IMPL-newfiles.diff` — unified diffs of every new production file under `tools/gates/**`.
- `tmp/gates-ext/review/NEWFILES.txt` — inventory (path/bytes/sha256) of all new paths.
- `tmp/gates-ext/review/MUTATION-2f-guard-removal.txt` — a guard-removal mutation on the real file with
  before/after SHA-256 (byte-identical restore) and the selftest going red while the guard is removed.
- `tmp/gates-ext/review/MUTATION-f1-donere.txt` — the same two-way mutation for the G6-1
  completion-marker set: the real file is narrowed to the older bold-only regex, the selftest goes red on
  exactly the two shape assertions, the file is restored byte-identically (SHA-256 shown before/after).
- `tmp/gates-ext/review/MUTATION-f1-range-vacuous.txt` — two-way mutation plus a **process-level exit-code
  matrix** for the scope guards: with the guards removed, `--scope=range:deadbeef…` exits **0** while
  reporting "empty change set" (the fail-open), and the four negative assertions go red; with the guards in
  place the same command exits **2**. `ctx.js` is restored byte-identically.
- `tmp/gates-ext/review/MUTATION-guards-independent.txt` — **per-guard** mutation evidence: each of the
  four guards is reverted on its own and the selftest must turn red on that guard's own assertions
  (the endpoint guard → 4 T17 reds; the diff guard → 1 T18 red; the status guard → 2 T18 reds; the
  conclusive-probe requirement → 1 T18 red), with every restore byte-identical.
- `tmp/gates-ext/ARCH-DESIGN.md` — the design proposal; `tmp/gates-ext/SLICE-DEFINITION.md` — the slice definition
  including the three blocking acceptance steps.
- The working tree itself: read `tools/gates/**` and both directive files directly.

## Notes that prevent misreading the evidence
- The older `tmp/gates-ext/step1/**` files predate a later fix to the added-line counting; the *post-fix* parity
  evidence lives in `tmp/gates-ext/step3/F4-count-parity.txt` (tree vs range counts identical).
- The `test/Gates` workflow digest is recomputed independently in `tmp/gates-ext/step3/` (see its report file);
  `go test ./internal/release/` passing also pins the step body.
- `--scope=ci` was only ever exercised locally on a **full** clone; its behaviour under a shallow clone is the
  defect the `fetch-depth: 0` edit above addresses, and the fix itself is not yet observed on the real runner.
- You have **no execution capability**: do not claim to have run anything. Base verdicts on the files and the
  pasted outputs, and say explicitly when something cannot be verified from the bundle.

## Output
Follow the format your role requires. Report every finding with severity, the exact `file:line`, and a
falsifiable failure scenario. Keep identifiers, paths, commands and criterion names verbatim.
