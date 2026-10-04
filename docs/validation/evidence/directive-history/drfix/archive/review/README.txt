DR-FIX pilot slice — ⑦ final-review input bundle (2026-09-22)
============================================================

PRIMARY TARGET (this slice):
  DRFIX-IMPL.diff                        = raw diff of the DR-FIX change set (2 files: main.go, main_test.go; +33/-15).
  Full files to read: tools/agent-probe/enforcement-proxy/{main.go,main_test.go}
  Slice definition:   tmp/drfix/SLICE-DEFINITION.md

FOLDED-IN OPEN-ITEM INPUTS (OB-9: "the next ⑦ input must contain these; do not run a separate 4th round"):
  OB9-T3T4-criteria-52c28af.diff         = T3/T4 CRITERION changes (§6.3 G5-a hardening + header-note position anchoring).
  OB9-T4-note-move-e0b08ec.diff           = the T4 role-file bytes (3 `prfrail-*` note blocks moved above the H1).
                                           NOTE: these bytes live in commit e0b08ec (v1.5), NOT in 52c28af (docs-only).
  OB9-v1.7-c473fd2.diff                   = v1.7 (blacklisted default `model` removal + G5-b + identifier disambiguation).
  OB9-v1.8-30cf799.diff                   = v1.8 (Appendix B.7 startup card).
  OB9-v1.9-6f186fd..a366fe3.diff          = v1.9 in full (3 commits: body fb12ff9 + ledger close-out b738e6b + report a366fe3).
  `_EN` mirror:                            docs/DELIVERY_DIRECTIVE_EN.md (strict line-for-line mirror of the CN authority).

Already-reviewed context (NOT a fresh input; listed to avoid re-doing work):
  docs/validation/sw5-work-tiers{,_EN}.md = the v1.9 slice report (⑥ PASS + ⑦ final review + re-review all closed).

Environment facts (immutable, for your judgement):
  - CI race step: CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/...
    => does NOT include tools/agent-probe/enforcement-proxy.
  - Local Windows host has no gcc => that package cannot be `-race`-tested locally either.