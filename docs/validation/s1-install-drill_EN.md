# S1 Windows Portable-Package Installation Drill

[中文](s1-install-drill.md)

Date: 2026-09-09; updated 2026-09-10. Verdict: `PASS AS CANDIDATE DRILL`. This report proves that archived trusted Windows amd64 artifacts can complete offline verification, first run, upgrade, rollback, and evidence-preserving uninstall under an explicit temporary directory. It does not authorize a GitHub Release. The product owner subsequently selected an S1 model of manually extracting a portable ZIP into a user-selected independent directory, invoking it by explicit path, and not modifying PATH.

## Inputs

1. Old package: `prfrail-Windows-1e7af676e8e84028c7ecfef1ccde91213738c20c` from GitHub Actions run `34277671704`.
2. New package: `prfrail-Windows-603e8dc03e9fb74ce9c53243d38e704fd8469c6a` from Node 24 Actions run `34325758548`.
3. Both contain `prfrail.exe`, `SHA256SUMS`, `sbom.spdx.json`, and `licenses.json`. A `prfrail-release verify` binary built from the current reviewed source verified each package against its full candidate SHA with exit code 0.
4. The drill root was the ignored `tmp/t018-install-drill/` only. It did not write user PATH, the registry, system directories, the source tree, or Git.

## Results

| Step | Result | Evidence |
|---|---|---|
| First install | PASS | Copied the old package into isolated `installed/current`; `version --json` returned `0.1.0+1e7af676e8e84028c7ecfef1ccde91213738c20c` |
| Upgrade | PASS | Copied the new package to staging before retaining current as rollback; the new binary returned `0.1.0+603e8dc03e9fb74ce9c53243d38e704fd8469c6a` |
| Rollback | PASS | Removed the new current directory and restored the old directory; the old binary again returned its original full version |
| Uninstall | PASS | Removed the installed binary directory; independent `data/retained-evidence.txt` remained present with content `retain-me` |
| Installation-guide replay | PASS | Zipped the four-file artifact from run `34325758548`, then followed `INSTALLATION_EN.md` from a clean directory; per-file SHA-256, both JSON documents, `version`, and Go-example `validate/preview` passed, with process and user PATH unchanged |
| `gh` download path | PASS | `gh run download 34325758548 --repo larsonzh/prfrail --name prfrail-Windows-603e8dc03e9fb74ce9c53243d38e704fd8469c6a` succeeded and automatically extracted exactly the four expected files; the temporary directory was removed |

The old `prfrail.exe` SHA-256 is `5b30f2b3f0de23ce662d6c75a27395cc59d5087678016faf3daebbda7ffc2db1`; the new value is `31335030152a417757990fd36e64f882d891a14112e0a2b865e1fc7a4280bfc8`. Checksums prove agreement with a manifest, not publisher identity.

## Final Candidate Replay (Before Nested-Directory Change)

On 2026-09-10, the Windows artifact `prfrail-Windows-bcc235496e0bece825047e95f3165a36585a59c0` from trusted workflow [run `34358467083`](https://github.com/larsonzh/prfrail/actions/runs/34358467083) completed a second lifecycle and five-step quick-start replay. The unexpired artifact has ID `10106867086`, size `2,398,240` bytes, and GitHub archive digest `sha256:10d902164194a8b1c6d1121adba5363bb69f829ac3516f7b891c6f5f41d08e06`.

| Check | Result | Evidence |
|---|---|---|
| Upgrade | PASS | Upgraded `0.1.0+603e8dc03e9fb74ce9c53243d38e704fd8469c6a` to `0.1.0+bcc235496e0bece825047e95f3165a36585a59c0`; old/new EXE SHA-256 values were `31335030152a417757990fd36e64f882d891a14112e0a2b865e1fc7a4280bfc8` and `c77a5780abbc381a7a696c24d3f428286fa1ebe9d77f21f06bcb6ac2ce41d692` |
| Rollback | PASS | After restoration, `version --json` again returned the old full version |
| Installation side effects | PASS | Process PATH, user PATH, and ProofRail process IDs were unchanged before and after the drill |
| Five-step quick start | PASS | `init`, `validate`, `preview`, noop-only `run`, and `report` all passed; the run had `chainState: COMPLETED` and `events: 11`, while the report had `taskState: PASSED=1` and `stepState: NOOP_RECORDED=1` |
| Cleanup | PASS | The temporary drill directory was removed and the Git worktree remained clean |

This run's verified artifact still used the old layout with four files directly at the archive root. The single top-level `prfrail/` directory introduced on 2026-09-10 was not exercised by this run. These results are final-candidate content and lifecycle evidence, not final-carrier acceptance of the new ZIP layout.

## Boundaries

1. This is a candidate portable-install drill using existing CI artifacts, not a formal release package or final user tutorial.
2. S1 defines no default install directory, does not modify PATH, and provides no implicit version switch. The download entry point and final filename remain subject to release approval.
3. The same flow must be rerun against the new final ZIP containing the top-level `prfrail/` directory, with version, commit, run, filename, support matrix, and release notes bound to that carrier.
4. Temporary artifacts and drill directories are removed after recording evidence and are not versioned.