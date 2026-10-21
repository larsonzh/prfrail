# C1 evidence pack (Linux AgentRunner native validation)

Slice: C1. Captured on an Ubuntu 24.04.5 LTS VM (go1.27.1 linux/amd64) at HEAD 6898f91.
Source tree of the run: a fresh clone at /home/larson/prfrail-c1; output directory born 2026-10-21 17:48:37 +0800.

## 归档条目清单

| 路径 | 记录内容 |
|---|---|
| `summary.tsv` | every harness invocation: face, case, phase, exit code, verdict path |
| `console.log` | the five summary lines printed by the driver |
| `toolchain-sha256.txt` | sha256 and version probe of the three built binaries |
| `baseline-diff.txt` | C1 terminal verdicts compared with the frozen B4 baseline |
| `driver-build.log` | per-face driver log: exact command line, cwd, exit code |
| `driver-cli.log` | per-face driver log: exact command line, cwd, exit code |
| `driver-f1.log` | per-face driver log: exact command line, cwd, exit code |
| `driver-f2.log` | per-face driver log: exact command line, cwd, exit code |
| `driver-f3.log` | per-face driver log: exact command line, cwd, exit code |
| `driver-f4.log` | per-face driver log: exact command line, cwd, exit code |
| `verdicts/f1-neg-run.json` | harness-written verdict, never edited here |
| `verdicts/f1-neg-r2-run.json` | harness-written verdict, never edited here |
| `verdicts/f1-neg2-run.json` | harness-written verdict, never edited here |
| `verdicts/f1-neg2-r2-run.json` | harness-written verdict, never edited here |
| `verdicts/f1-neg3-run.json` | harness-written verdict, never edited here |
| `verdicts/f1-pos-run.json` | harness-written verdict, never edited here |
| `verdicts/f2-neg-run.json` | harness-written verdict, never edited here |
| `verdicts/f2-neg-r2-run.json` | harness-written verdict, never edited here |
| `verdicts/f2-pos-run.json` | harness-written verdict, never edited here |
| `verdicts/f3-neg-run.json` | harness-written verdict, never edited here |
| `verdicts/f3-neg-r2-run.json` | harness-written verdict, never edited here |
| `verdicts/f3-pos-run.json` | harness-written verdict, never edited here |
| `verdicts/f4-neg-forge.json` | harness-written verdict, never edited here |
| `verdicts/f4-neg-park.json` | harness-written verdict, never edited here |
| `verdicts/f4-neg-r2-forge.json` | harness-written verdict, never edited here |
| `verdicts/f4-neg-r2-park.json` | harness-written verdict, never edited here |
| `verdicts/f4-pos-park.json` | harness-written verdict, never edited here |
| `verdicts/f4-pos-recover.json` | harness-written verdict, never edited here |

## 排除项

Run roots, the seed workspace, the emitted inputs file and the built binaries stay out of the pack: they are
regenerable from the script plus the frozen chain definitions, and the run-root trees are large.

## 复现

Run tools/validation/c1-linux/run-c1-faces.sh against a fresh output directory (see the known limits in the report
before reusing an existing directory), then re-run the comparison that produced baseline-diff.txt.
