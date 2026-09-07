# ProofRail contract oracle

This directory is an independent S0 validation tool. It does not import future ProofRail production code and does not generate expected fixture outcomes.

## Run

```powershell
Push-Location tools\contracts
npm ci --ignore-scripts
npm test
Pop-Location
```

The suite compiles all 29 Draft 2020-12 schemas, checks every positive and negative catalog, verifies full per-schema coverage, and recomputes the frozen JCS vectors. Dependencies are pinned in `package-lock.json`: Ajv validates schemas and `canonicalize` supplies an RFC 8785 implementation.

Ajv keeps strict schema and format validation enabled. `strictTypes` is disabled because it is an Ajv authoring heuristic, not a Draft 2020-12 requirement, and rejects valid constraints whose type is supplied by a composed schema branch.

## Fixture layers

Each catalog entry contains exactly `fixtureId`, `schema`, `expected`, `rejectionLayer`, `reason`, and `instance`.

- `schema`: the independent validator accepts or rejects one JSON instance.
- `semantic`: `instance.documents` contains individually schema-valid records, then the oracle checks static IDs, references, ownership, dependency/generation cycles, explicit scope intersection, documentation rules, handoff write targets, and fixture-supplied review context.
- Runtime gates are intentionally outside this tool. Filesystem resolution, processes, capabilities, secrets, budgets, current identity/time, and actual publication remain S1 runtime evidence. Review fixtures supply frozen identity/time values only to test deterministic semantic decisions.

All repository JSON consumed here must be UTF-8 without BOM and LF-only. A malformed catalog or unknown schema is an infrastructure error; an unexpected accept/reject result is a test failure.
