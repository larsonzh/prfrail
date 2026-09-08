# T017 Verifier Baseline Independent Review Record

[中文](t017-baseline-review.md)

Date: 2026-09-08. Verdict: `PASS FOR BASELINE REVIEW`. No reproducible `BLOCKER` or `HIGH` finding was identified. The reviewed implementation may become a verifier-baseline candidate commit, but this verdict is not T017 completion, release authorization, publisher authentication, or proof of successful GitHub CI.

## Review Source and Boundary

1. The reviewed object was the uncommitted T017 worktree based on HEAD `fc99f553c3066b24230a093d14bc17af0c6a398e`, including the workflow, release verifier, workflow contract/mutation tests, fixtures, and related documentation. It is not immutable until committed.
2. The user initiated a strict read-only GPT-5.3 Codex review requiring adversarial analysis, findings-first reporting, no file edits/network/commit/push/release, and an explicit distinction between local checks and real GitHub CI.
3. This page archives the conversational verdict and local command evidence. It is neither a signed receipt nor independently authenticated model identity evidence. Before landing the baseline, the maintainer must still inspect the actual diff and exact file list.
4. Status remains `IMPLEMENTED_AWAITING_CI`. T017 can close only after an earlier verifier-baseline commit evaluates a different candidate commit and the first native GitHub Windows/Linux matrix succeeds with archived run evidence.

See also the matching [DeepSeek V4 Pro independent review record](t017-deepseek-baseline-review_EN.md). The two records corroborate each other but do not constitute a signature or attestation.

## Findings

1. `MEDIUM`: full-script SHA-256 values are content-drift gates, not PowerShell semantic proofs. A malicious same-change update to both workflow scripts and expected test digests could make local tests green again. The control is independent baseline review and a prohibition on deriving expected digests from the workflow under test. A protected external read-only digest baseline may be considered later.
2. `LOW`: production offline static checks cover direct `net` imports and enumerated downloader/networked-Git literals, but the command vocabulary is not a complete OS network sandbox. The current executor advertises no network-isolation capability, so `deny/loopback/allowlist` fail closed before execution.
3. `LOW`: the local toolchain was Go 1.27.0 while CI pins Go 1.22.12. Local success cannot replace target Windows/Linux runner results.

## Local Validation Evidence

The reviewer ran the following checks from `d:\LZProjects\prfrail` in read-only mode or with only ordinary Go build-cache effects. Every command exited 0:

1. `CGO_ENABLED=0 go build ./...`
2. `CGO_ENABLED=0 go vet ./...`
3. `CGO_ENABLED=0 go test -count=1 ./...`
4. Focused workflow contract and mutation tests: passed, with 540 mutation subtests.
5. `git diff --check`
6. Editor diagnostics for critical files: no errors.
7. Changed-file encoding/EOL checks: passed; Markdown uses UTF-8 BOM + LF, while Go/YAML/JSON use UTF-8 without BOM + LF.
8. Production Go static spot-check: zero direct `net` imports; `os/exec` appeared only in approved execution boundaries, with no common downloader or networked-Git command found.
9. `tmp/` contained only `.gitkeep`; `sessbridge` was unchanged.

## Baseline Landing Gates

1. Before committing, inspect the complete worktree diff relative to `fc99f553c3066b24230a093d14bc17af0c6a398e` and confirm no post-review code drift beyond this record and its links.
2. Stage an exact file list and create a dedicated verifier-baseline commit; do not mix unreviewed changes, temporary files, or the candidate commit into the baseline.
3. After the baseline commit reaches `origin/main`, create a content-different candidate commit with a different SHA.
4. Manually dispatch the workflow from `main`, recording verifier SHA, candidate SHA, run URL, platform results, and final artifact/SHA256SUMS digests.
5. Any platform, input-pin, or evidence failure keeps the gate blocked. This PASS is not release authorization.

## Actions Not Performed

This review did not run `git commit`, `git push`, tagging, signing, attestation, GitHub Release, workflow dispatch, or real-runner attack exercises.