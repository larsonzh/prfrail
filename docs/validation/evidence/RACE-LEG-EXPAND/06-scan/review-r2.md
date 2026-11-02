# P49 RACE-LEG-EXPAND ⑥ 独立扫描 r2 结论（A–E 格式合规重发 · 逐字归档）

> 载体：⑥ 独立扫描 `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`。r1 格式偏离（§7.2.2 A–E 契约），本 r2 为合规重发（预算 2/3）。以下为代理输出原文，未改写。

---

## Section A
I reviewed the only-read diff for the CI Race expansion and the matching workflow freeze contract: one run line in .github/workflows/ci.yml, one digest entry in internal/release/workflow_test.go, and the ledger/evidence package for P49. The scope remains infrastructure-only: no production Go code changes, no CI step deletions, and no package coverage outside the six-package Linux race target.

## Section B
(no findings)

## Section C
git diff shows exactly the intended changes: ci.yml changed only the Race run target from two packages to six packages; workflow_test.go changed only the digest for test/Race from ff4826f2103a2457cd306662f2542e54d809546deb1f804964b63f2148185dd5 to 825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9; the ledger files and evidence directory were updated as expected. The SHA-256 recomputation over the exact run string "CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/... ./internal/console/... ./internal/gates/... ./internal/guard/... ./internal/tickets/..." produced 825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9, matching the frozen contract. For runtime evidence, the VM log in docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/vm-race-four-pkgs-green.txt shows rc=0 and these four package results: console ok 20.410s, gates ok 3.384s, guard ok 3.043s, tickets ok 7.420s. The evidence bundle also includes 04-impl/logs/04-release-workflow-tests-verbose.txt and SHA256SUMS.txt, and the file hashes are internally consistent with MANIFEST.md; the package body is UTF-8 BOM+LF as required.

## Section D
If the Linux guard were removed from the Race step, the workflow contract and red/green TDD assertion would fail because the review enforces runner.os == 'Linux' and the package list must be exactly the six targets. If the old digest were restored in workflow_test.go, the workflow contract test would fail immediately with a SHA mismatch at test/Race. If the package list were truncated back to adapters and chain, the red-baseline assertion and the VM coverage log would no longer match the accepted six-package scope. No uncovered invariant remains: the Linux-only condition, the six-package target set, the frozen digest, and the evidence-package integrity are all directly checked.

## Section E
I could not independently re-run GitHub Actions in the remote repo; only the archived VM and bundle evidence are available.
What would overturn this PASS: a current HEAD run where the workflow contract test fails on the digest or Linux guard; a mismatched SHA256SUMS.txt or MANIFEST.md for the evidence bundle; or any real race failure in the four newly added packages on the same Ubuntu toolchain.

INDEPENDENT SCAN: PASS
