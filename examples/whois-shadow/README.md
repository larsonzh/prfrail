# whois Read-Only Shadow Example

This `AT-15` example compares a frozen, non-secret whois test fixture with an offline ProofRail shadow projection. It never invokes whois, performs network requests, writes the whois repository, or changes its production flow.

The baseline was exported on 2026-09-09 from `D:\LZProjects\whois\testdata\cidr_matrix_cases_draft.tsv` under a one-session read-only authorization. Source SHA-256: `8a6abecb0178a970677b84b3e5403ce8830f281e1a197d792ae54089b3e9efb0` (1,981 bytes, 10 lines, zero secret-marker matches).

Run from the ProofRail repository root:

```powershell
go run ./examples/whois-shadow --baseline ./examples/whois-shadow/baseline.json --shadow ./examples/whois-shadow/shadow.json
go test ./examples/whois-shadow
```

The positive fixture must preserve all nine case IDs, inputs, results, and failure classifications. Tests also prove that input, result, and failure-classification drift each produce a non-zero comparison failure.

This is deterministic fixture parity evidence only. It is not a live whois execution, a production cutover, or evidence that ProofRail generated the historical whois result.