# P83 input list

Slice: SECURITY-FINDINGS-REMEDIATION (P83), started 2026-11-20 by user designation.

## Anchors verified at step 0

| # | Anchor | Expected | Actual | Result |
|---|---|---|---|---|
| 1 | HEAD | 7a63962 | 7a63962 | OK |
| 2 | origin/main | 7a63962 | 7a63962 | OK |
| 3 | Directive version | v1.71 | v1.71 (docs/DELIVERY_DIRECTIVE.md:3) | OK |
| 4 | SELFTEST | 740/740 | PASS 740/740 | OK |
| 5 | Three scopes | TOTAL_FAIL=0 for all | index / tree / ci all 0 | OK |
| 6 | v64 snapshot lines | 224 | 224 / 224 (CN/EN) | OK |
| 7 | v65 snapshot lines | 220 | 220 / 220 (CN/EN) | OK |
| 8 | v66 snapshot lines | verified by v67 | 216 / 216 (CN/EN), recorded only | RECORDED |

Extra: working tree porcelain empty; tmp top level = .gitkeep + gates-cache.

## Inputs to read

1. docs/validation/s2-pre-independent-security-review{,_EN}.md - section 9 (F1-F8 table), section 15 (ruling record).
2. docs/t027/CHECKLIST_v66{,_EN}.md - section 4 queue item for P83.
3. docs/CONTRACTS{,_EN}.md - cited sentences (:147 / :212 / :214 / :216 / :391, plus F4 row line refs).
4. docs/SECURITY.md - :23 / :60.
5. docs/DELIVERY_DIRECTIVE.md - sections 1.4 / 5.2 / 5.3 / 7.5 / 7.7 / 8.2 / 9.2.
6. Code state of the seven files:
   - internal/gates/runner.go (:208, :73)
   - internal/console/runtime_agent.go (:1245, :731)
   - internal/console/runtime_supervised.go (:1642)
   - internal/chain/review.go (:116)
   - internal/console/ai.go (:94)
   - internal/guard/process_windows.go (:457)
   - internal/adapters/ai_probe_copilot_executor.go (:53)

## Findings under remediation (verbatim source = P79 report section 9)

Priority: F2 / F3 (face 2 High) -> F6 (face 3 High) -> F1 / F4 / F5 / F7 / F8.
