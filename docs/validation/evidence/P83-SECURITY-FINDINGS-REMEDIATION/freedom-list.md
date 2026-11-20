# P83 freedom list (execution degrees of freedom)

Scope: what this slice MAY decide freely versus what the brief and the directive fix.

## Fixed (not free)

- The eight findings F1–F8 and their file:line anchors, verbatim from the P79 report section 9. No expansion, no shrinking.
- Priority order: F2 / F3, then F6, then F1 / F4 / F5 / F7 / F8.
- Step chain and ordering, including the blind review before the final review.
- Model and tier assignments; no mid-slice tier change.
- (7) budget: 1 final review plus 3 re-reviews (user ruling, not a precedent); blind review 1, counted separately; one additional re-review authorized by the user on 2026-11-20 (not a precedent).
- Boundaries: no S1 Exit gate change, no threat-model baseline change, no directive version bump, no CI or `DEV_PLAN` change, no gitee.
- Contract-first rule: a fix touching contract semantics changes `docs/CONTRACTS{,_EN}.md` (plus schemas or fixtures) before code. In this slice all eight were judged DOES-NOT-TOUCH, so no contract text changed.

## Free (decided by the implementer, recorded here)

- The internal shape of the new helpers: function names, signatures, package placement, unexported seams.
- Error sentinels for the refusal paths, provided each refusal stays observable and testable.
- Test file placement and naming.
- Whether a check is expressed as a re-read plus verification helper or as an invariant at the call site.
- Comment wording (English), provided each comment names the finding id it serves.

## Not free (stop and escalate instead)

- Any new schema record type, new CLI command, new third-party dependency, or contract wording change.
- Changing exported behavior beyond what a finding requires.
- Touching files outside the seven target files plus their tests, the necessary wiring, the contract documents, and the report and ledger artifacts.
- Registering a new OB or ADR without a user ruling.
- Any real BYOK or network probe (probe budget 0; section 7.6 pre-authorization required).

## Escalations actually taken

| Trigger | Outcome |
|---|---|
| T2 STOP: no frozen policy artifact for AgentRunner runs | reported, not implemented; user ruled deferred to P84 |
| T8a: proving the child cannot write unredacted logs needs a new subsystem | reported; user ruled deferred to P84 |
| T3 Windows/Linux launch binding | partially implemented; user ruled residuals deferred to P84 |
| T4 / T5#4 residuals after the additional re-review | user ruled deferred to P84; F1 / F4 may not be recorded as closed |
