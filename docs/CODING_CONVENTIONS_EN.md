# ProofRail Coding and Engineering Conventions

[简体中文](CODING_CONVENTIONS.md)

Encoding/Git rules finalized 2026-09-04; pre-prototype workflow added 2026-09-06. Applies to all maintained repository files.

## 1. Encoding and Line Endings

| Type | Encoding | EOL |
|---|---|---|
| .md/.ps1/.json | UTF-8 with BOM | LF |
| Strict-parser JSON exceptions | UTF-8 without BOM | LF |
| Other code/text (.go/.py/.js/.txt/.gitignore/LICENSE etc.) | UTF-8 without BOM | LF |

The explicit exception .github/hooks/context-mode.json must have no BOM and retain "$schema": "../context-mode-schema.json". Check both after context-mode upgrades. JSON directly consumed by strict Node/Go json.Decoder parsers also requires no BOM. Do not apply arbitrary encoding changes to unrelated files.

Check BOM/LF before commits. Mechanical normalization must not change semantics. Local .vscode/scripts/ is ignored; shared engineering tools belong in tracked tools/.

## 2. Go Engineering

Module github.com/larsonzh/prfrail, Go 1.22+. Use gofmt and Go naming conventions; exported documentation begins with the symbol name. cmd/prfrail is entry/composition; internal contains adapters, applier, chain, console, evidence, gates, guard, repair, snapshot, taskdef and tickets. Required gates: go build ./..., go vet ./..., go test ./.... Follow RFC/existing path conventions instead of scattered magic paths.

## 3. Documentation and Protocol

The authoritative design is [RFC](RFC-proofrail-unattended-ai-engineering-product.md); revise it before protocol implementation. ProofRail consumes SessionBridge silent plus file queue only. visible is SessionBridge product functionality, not ProofRail's formal protocol.

## 4. Git Discipline

Push only origin by default. Gitee is a mirror; no Gitee push without explicit same-turn permission. No git commit/push without explicit same-turn authorization. Stage exact files with git add <file>, not all changes. Commit format: type: summary (chore/feat/fix/docs/test/build).

## 5. Temporary Files

Use root tmp/ for one-off build/test/debug artifacts; only tmp/.gitkeep is tracked. Remove temporary artifacts when finished. .vscode/ and tmp/ are local/ignored; shared scripts belong in tools/.

## 6. Implementation and Review

- Freeze behavior from RFC/contracts, then table-driven negative tests and minimal implementation. Do not implement unapproved defaults, permissions or wire fields.
- Consumers own interfaces; CLI composes, core never imports concrete adapters/console. Abstract real test/substitution boundaries only, not empty module scaffolds.
- Propagate cancellation/deadlines with context.Context. Processes/goroutines have explicit stop/join paths. Inject clock/process/filesystem faults; never guess completion with sleep.
- Preserve error causes/object identity; wrap and use errors.Is/As, not string matching for state control. Never swallow write/close/flush/rollback errors.
- No process/file side effects in package init, and no os.Exit/log.Fatal in libraries. Separate executable/args without implicit shell.
- Use structured parsers and reject unknown fields. Hash frozen canonical data/original content bytes, not map iteration or string concatenation.
- Validate traversal, aliases, symlinks/reparse points and check-use races. Follow journal/atomic contracts; rename alone is not a group transaction.
- Use t.TempDir/non-secret fixtures and terminate helpers only. Never change real source trees, released seeds or user processes; never restore files while relevant processes run.
- Use standard library/reviewed pure-Go dependencies, pinning version/license/minimum Go/source. No incidental toolchain upgrades or runtime services.

## 7. Low-Cost Delivery

Start at [documentation navigation](DOCUMENTATION_PLAN_EN.md), select T/REQ/AT and read only relevant contracts/tests. Form one falsifiable hypothesis, make one small edit, validate narrowly immediately, then build/vet/test. Use gofmt; never remove assertions, weaken gates, broaden paths or reset budgets for green results.

Report scope, actual commands/exits/counts, gaps and next action. Escalate two same-class failures without new evidence. No additional paid model calls without authorization. Missing tests are not acceptance; cross-compilation is not native execution; model prose is not evidence.

## 8. Bilingual Contributions

Update specifications before tests/examples/usage. Chinese retains original filenames, English uses _EN with reciprocal links and identical REQ/AT/T/ADR identifiers in one review. The current Chinese RFC remains authoritative; translation cannot alter state machines or fields.

Review submissions include requirement/task IDs, minimal behavior differences, test evidence, compatibility and documentation impact. Security/protocol changes need independent review; the owner approves stages/budgets. Document creation or plan coverage does not imply S0 readiness.