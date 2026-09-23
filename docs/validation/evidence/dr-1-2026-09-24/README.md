# DR-1 evidence bundle (2026-09-24)

Decisive artifacts for slice DR-1 (product availability probe: CLI credit floor for
`prfrail ai check`). The raw scratch captures used during the slice live under the
git-ignored `tmp/dr1/` directory and are declared `local-only` in the slice report.

## Provenance

| Artifact | What it shows |
|---|---|
| `proofrail.chain.json` | the chain configuration used as probe input: profile `copilot-cli-host` (host auth, model `mai-code-1.1-flash`) bound to channel `agent-runner-cli` |
| `availability-1083-available.json` | the record written by `ai check`: `status=available`, `requestsUsed=1`, no `reason`, `profileConfigHash` identical to the B1 main baseline |
| `probe-01-ai-check.json` | product stdout of the probe run (`command`/`ok`/`exitCode` wrapper around the same record) |
| `probe-01-ai-verify.json` | product stdout of `ai verify` against that record (freshness and request-budget checks pass, zero model calls) |

## Hashes (sha256)

| File | sha256 |
|---|---|
| `proofrail.chain.json` | `99d0dcec54e99b73c65324398a2ef99e51a525c2243cbdc0eef3fc895d8c6d6a` |
| `availability-1083-available.json` | `2dc3edbd9393b4e428ec4a9d3d98ec816b3b1d1c6aae374833fb334d64364418` |
| `probe-01-ai-check.json` | `93c661906d098e1ca669ef76e23c6ff958bf0ea569df4aa926bff45af7cfde20` |
| `probe-01-ai-verify.json` | `8c4fe63cf59d3d9a7b523153ad9f55dbd5d8873ace30f9bc86337f3ab6a6eabd` |

## Reproduction

```
go run ./cmd/prfrail validate --chain docs/validation/evidence/dr-1-2026-09-24/proofrail.chain.json --json
go run ./cmd/prfrail config explain --chain docs/validation/evidence/dr-1-2026-09-24/proofrail.chain.json --json
go run ./cmd/prfrail ai check --chain docs/validation/evidence/dr-1-2026-09-24/proofrail.chain.json \
  --channel agent-runner-cli --copilot "%APPDATA%\npm\node_modules\@github\copilot\node_modules\@github\copilot-win32-x64\copilot.exe" \
  --workspace tmp/dr1/ws --max-requests 1 --out <new-path>.json --json
```

The probe spends one paid request per successful run and is never triggered automatically;
the chain config carries no `secretRef` (host authentication).
