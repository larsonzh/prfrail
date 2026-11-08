# P61 AT23-RESTRICTED-CANDIDATE restricted real-candidate evaluation (validation report)

**Status: ✅ Done (PARTIAL close-out)** (final verdict = ⑦ final review `FINAL REVIEW: PARTIAL`, accepted by the user on 2026-11-08 under ruling (E); committed `6bb23a1`, CI `42452377967` dual-green; waiting at ⑩ for follow-up slice rulings).

- Slice: `AT23-RESTRICTED-CANDIDATE` (acceptance-class `[SLICE]`, v50 §4 layer 1 item 8 "AT-23 restricted real candidate and sign-off")
- Started: 2026-11-07; start anchor HEAD = `05e8e35`; directive v1.64
- Controller model: `DeepSeek V4 Pro` (ADR-035, registered at step 0 under the OB-37 meta-rule)
- Close-out ruling: ADR-036 (user, 2026-11-08: accept PARTIAL, stop the ⑦ re-review loop)
- Report conclusion line: `FINAL REVIEW: PARTIAL`
- Evidence pack: `docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/` (MANIFEST and SHA256SUMS self-consistent)

## 1. Executive summary

**AT-23 conditional closure not achieved**: premise ① is unmet because of the F2 `process-tree-taken-away` inherited drift (cause `02d5797`, not introduced by P61; the B4 baseline re-run today would go red as well); premise ② is met (run-1 start + run-2 execution blocked + run-3 stop + run-5 unsupported, the same shape as B1 #8); this slice closes as PARTIAL (user ruling (E), ADR-036).

- Premise ① "mechanism legs all green": **not met** — one F2 positive-leg assertion moved from the baseline `true` to `false`, while the 9 legs show 0 drift on `verdictClass` / `chainState` / `taskState` / `stepState` / `observations`.
- Premise ② "the real candidate leg actually ran (≥3 times, covering start / execution / stop) with unsupported / unknown recorded": **met** — the evidence is the five legs in §3 and §4.1.
- Close-out nature: PARTIAL is a legal value of the §4.1 conclusion-line closed set; this slice **does not claim** conditional closure, and **does not claim** real-candidate compatibility.
- Follow-up ownership: the F2 fix and the R-08.5 criterion revision are both **independent follow-up slices** and do not block this slice's close-out (see §6.1 and §7).

## 2. Layered facts (layer 0 → layer 3)

### 2.1 Layer 0 · BYOK feasibility

- Verdict: **all green**. The wrapper and the bare run report the same version (`1.0.83`), the BYOK call returns normally, and all six BYOK environment variables are cleared item by item at wrap-up.
- Cost: layer 0 spent ¥0.01 (per the user's measured DeepSeek reading).

### 2.2 Layer 1 · skeleton, five items

| # | Check | Observed | Verdict |
|---|---|---|---|
| 1 | wrapper `--version` | `GitHub Copilot CLI 1.0.83` | PASS |
| 2 | bare `--no-auto-update --version` | token `1.0.83` (trailing-period difference registered) | PASS |
| 3 | `-Unlock` control | `1.0.85` (inside the controlled set) | PASS |
| 4 | one BYOK call | non-empty reply + six env entries cleared | PASS |
| 5 | `--effort=high` three-state probe | exit 0 with normal output ⇒ `supported` | PASS |

### 2.3 Layer 2 · nine-leg comparison (core scenarios)

- Covers the six F1 isolated-workspace legs and the three F4 unknown-recovery legs; `verdictClass` / `chainState` / `taskState` / `stepState` / `observations` are **identical to the B4 baseline (0 discrepancies)**, and `failedChecks` is an empty set on both sides.
- Only label drift: the F1 assertion id moved from `five-facts-distinct` to `six-facts-distinct` (a fact-enumeration extension, **a strict superset**, measured all green), with the same cause `02d5797`.
- The deviation set matches the baseline (three deviation legs, one item each, `kind=deviation`).

### 2.4 Layer 3 · mechanism legs, three faces

- F2: the three legs keep the baseline verdict class (`uncertain-or-paused` / `PAUSED`) with **one assertion-value drift** — the positive leg's `process-tree-taken-away` moved from `true` to `false` (inherited; the `hang` branch writes no `stub-child.pid`, so the sixth frozen fact judges it a fail-closed gap); error evidence rose from 5 to 6 items.
- F3 / F4: all green, byte-identical to the baseline (empty `failedChecks` and the same in-pack re-run results).
- postflight six-fact rebuild: the five facts of the two legs that carry frozen facts MATCH byte for byte; `process-stop` is never persisted and **cannot be rebuilt**, explicitly flagged and backed by the semantic cross-check "liveness of the reported pids".

### 2.5 Gates and review rings

- ④ gates: `gofmt` / `go build` / `go vet` / `go test ./...` all green; `gate.js --all` over the three scopes (`index` / `tree` / `ci`) reports `TOTAL_FAIL=0`; `SELFTEST: PASS 722/722`.
- ⑤ pre-review: `PRE-REVIEW: PASS` (check faces A–I all PASS, 4 LOW; disposition in §8.1).
- ⑥ independent scan: `INDEPENDENT SCAN: PASS` (Section B is `(none)`; the sixstage fence is in §8.2).
- ⑦ blind review `BLIND REVIEW: FINDINGS` and final review `FINAL REVIEW: PARTIAL`: the two independent reads overlap on the same HIGH (the F2 inherited drift blocks premise ①); see §8.3.

## 3. Real-candidate leg evidence (five legs)

| Leg | Shape | Exit | Failure stage | Result |
|---|---|---|---|---|
| run-1 start | production wrapper BYOK call (`--effort=high`) | 0 | none | normal reply plus one real shell call; six env entries cleared |
| run-2 execution | existing probe `-Scenario tool-deny`, three attempts | 1 | start | all three failed in the start stage, zero model calls, zero cost |
| run-3 stop | background wrapper call + `taskkill /PID <host> /T /F` | taskkill 0 | stop | the whole process tree terminated; processes `gone` afterwards; host exit code unavailable |
| run-4 deny observation | narrow allowlist + `--deny-tool=shell(Get-ChildItem)` | 0 | none | refusal observed, but attribution is limited, so it is **not** recorded as unsupported / unknown |
| run-5 unsupported | wide allowlist + the same deny (B1 #8 shape) | 0 | none | `VERDICT: unsupported`: `gci -Name` actually executed and returned a directory listing |

- The three run-2 attempts failed specifically: the `-File` form cannot bind `string[]` from the command line; the in-process array splat is still positional and fails the same way; after the tmp companion switched to a named splat the tool reached its own precheck and was refused by the version pin (the bare `--version` resolves to `1.0.85` while the tool pins `1.0.83`). All three attempts carry argv / stderr / exit-code evidence.
- run-5 attribution: the wide allowlist leaves `--deny-tool=shell(Get-ChildItem)` as the only filter that could explain a refusal, and it did not take effect ⇒ the alias bypass reproduces, with the same shape and conclusion as T027 B1 probe #8.
- Cost: **¥0.09 cumulative** (layer 0 ¥0.01 + layer 1 ¥0.03 + layer 3 run-1~3 ¥0.04 + run-4 ¥0.00 + run-5 ¥0.01), all per the user's measured DeepSeek web reading; the CLI-reported token block is the 1.0 fallback estimate used while the Anthropic `deepseek-flash` model lacks a canonical model ID, and is **not a billing basis**.

## 4. PARTIAL conclusion

### 4.1 Four elements

- **The ⑦ verdict is correct**: premise ① being unmet has hard evidence (the failed check in `l3-f2-pack-verdict.json` and the assertion-value drift in `l3-baseline-compare.txt`), and premise ② being met has per-leg stderr / exit codes / failure stages.
- **The inherited drift is not introduced by this slice**: cause `02d5797` (2026-10-22, A5 bound pid reports as the sixth frozen fact); this slice changed no harness, no `tools/agent-stub` and no production semantics.
- **PARTIAL is a legal §4.1 closed-set value**: the ⑦ independent final reviewer's value set is `PASS` / `PASS WITH FIXES` / `PARTIAL` / `FINDINGS`.
- **The fix and the revision are independent follow-ups**: the F2 assertion fix and the R-08.5 criterion revision are opened as their own slices and do not block the P61 close-out (ADR-036 four elements and follow-up obligations).

### 4.2 Contract sentences and the three sentences (verbatim)

1. R-08.5 criterion sentence (verbatim, original text kept): 「机制腿全绿 + 真实候选腿实际跑过（≥3 次，覆盖启动/执行/停机）并记录 unsupported/unknown。」
2. The conclusion sentence R-08.5 requires the report to state (verbatim): **AT-23 条件闭合，不构成真实候选兼容性证明**
3. **Non-applicability statement**: premise ① is unmet in this slice, so sentence 2 **must not be used as this slice's conclusion**; it is recorded verbatim here only to keep the contract citation verifiable, and this slice's conclusion is always §1's "conditional closure not achieved".
4. R-08.3 power-loss wording (verbatim, original text kept): 「不覆盖断电场景；断电耐久引用 B3b（P′ 覆盖物理机目标部署环境未论证）。」
5. R-08.6 BYOK boundary (verbatim, original text kept): 「BYOK 结论不能替代 B1 的 Copilot-hosted 基线。」

## 5. Known boundaries

- `--effort`: the layer 1 three-state verdict is supported, and layer 3 fixes `high` without changing tiers; whether the BYOK provider **actually applies** that tier is unmeasurable in this environment (there is no observability surface for that provider), and `low` / `medium` / `xhigh` / `max` were not probed.
- `--context`: default tier; neither run-1 nor run-3 used `-C`, registered as an **uncovered boundary**.
- 1.0 fallback: the token block is observational only and takes no part in cost attribution (see the last item of §3).
- run-4 attribution limit: argv carries a narrow `--allow-tool` and a `--deny-tool` at once, so a refusal is explainable by either filter; it is therefore not recorded as unsupported / unknown.
- Power loss and durability: power-loss scenarios are not covered; power-loss durability cites B3b (whether P′ covers the physical-machine target deployment environment is not argued).
- Inherent gap of the stop leg: run-3's host exit code is unavailable (the process tree was killed before the CLI reached its summary), and `stderr` holds only the 1.0 fallback warning.
- LOW-4 (⑤ pre-review): `test-r3.md` §6 and §9 carry a state gap between the closure-premise self-check and the conclusion after run-4 / run-5 were added; it is **closed by this report**; `test-r3.md` stays archived as the receipt of that moment and is not retro-edited.

## 6. Explicitly not executed, and documentation sync

### 6.1 Explicitly not executed

- **The F2 fix and the R-08.5 revision are independent follow-ups**: this slice touches no `tools/agent-stub`, no harness and no criterion text.
- No product code, CI, contract, schema or fixture change (the only production addition is `tools/agent-probe/CopilotCli.ps1`, the BYOK wrapper delivered in ②); the sole exception is the two verbatim-transcription character entries registered in `tools/gates/cjk-newwords.txt` (data registration, not a criterion change).
- Nothing committed / pushed; no new slice started; no ruling made on the user's behalf; the stage 1 → stage 2 upgrade was not executed.
- No Ubuntu VM and no network proxy were enabled (R-08.3 / R-08.4).

### 6.2 Documentation sync, 15 items one by one

| # | Document | Disposition |
|---|---|---|
| 1 | `docs/validation/at23-restricted-candidate{,_EN}.md` | created (this report, CN authoritative + EN paragraph-by-paragraph mirror) |
| 2 | evidence pack `docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/` | built (ring receipts, gate readings, MANIFEST, SHA256SUMS) |
| 3 | `docs/t027/REMAINING_SLICES{,_EN}.md` | changed: the P61 block flips to ✅ Done (PARTIAL close-out), completion record, status row, layer 1 item 8 marker |
| 4 | `docs/DEV_PLAN{,_EN}.md` | changed: the T027 paragraph's AT-23 status annotated (still not passed, pending the R-08.5 revision slice) |
| 5 | `docs/TEST_STRATEGY{,_EN}.md` | changed: a new AT-23 acceptance-conclusion paragraph citing the report path |
| 6 | `docs/ADR_REGISTER{,_EN}.md` | ADR-035 and ADR-036 were registered by the controller; this slice does not touch it |
| 7 | `docs/validation/s1-exit{,_EN}.md` | changed: the AT-23 row and the S-G7 dependency chain; **still BLOCKED**, S1 exit is not claimed as passed |
| 8 | `docs/DELIVERY_DIRECTIVE{,_EN}.md` | not applicable: no directive-semantics change; the R-08.5 criterion is cited from existing clauses, nothing added |
| 9 | `docs/CONTRACTS{,_EN}.md` | not applicable: no wire / schema / state-machine / criterion-semantics change |
| 10 | `docs/ARCHITECTURE{,_EN}.md` | not applicable: no module-boundary, data-flow or storage change |
| 11 | `docs/SECURITY{,_EN}.md` | not applicable: no permission-model, threat-matrix or secret-handling change |
| 12 | `docs/S1_SUPPORT_MATRIX{,_EN}.md` | not applicable: no platform support-tier change |
| 13 | `docs/S1_RELEASE_NOTES{,_EN}.md` | not applicable: no tag / artifact / feature-item change |
| 14 | `docs/t027/CHECKLIST_v50{,_EN}.md` | changed: the user's same-turn 2026-11-08 ruling overrides snapshot freezing — layer 1 item 8 marked 🔄 "P61 assessment complete; premise ② satisfied; premise ① not satisfied due to inherited drift, awaiting the R-08.5 revision slice ruling" (see §6.3) |
| 15 | `docs/OPERATIONS{,_EN}.md`, `docs/INSTALLATION{,_EN}.md`, `docs/PRODUCT_REQUIREMENTS{,_EN}.md`, `docs/BUSINESS_WORKFLOWS{,_EN}.md` | not applicable: no operations manual, installation model, product scope or business-flow change |

### 6.3 Conflict disposition (v50 §4 layer 1 table vs snapshot freezing; closed by the user's same-turn ruling)

- Item B of the briefing asks for "§4 layer 1 table item 8" to be marked 🔄; that table body lives in `docs/t027/CHECKLIST_v50{,_EN}.md` (the v50 snapshot is this slice's start anchor), while item F-14 of the same message requires v50 to be **not applicable** (snapshot frozen, v51 inherits).
- Disposition: **the user's same-turn 2026-11-08 ruling overrides snapshot freezing** — layer 1 item 8 (AT-23) is marked 🔄 "P61 assessment complete; premise ② satisfied; premise ① not satisfied due to inherited drift, awaiting the R-08.5 revision slice ruling"; after P61 completes it moves from "awaiting start" to "awaiting fix", not blocking layer 2 / layer 3 (both fully dequeued). The v50{,_EN} layer 1 table item 8 dependency/boundary columns are updated per the ruling and included in this report's artifacts fence; the same field in the live ledger `REMAINING_SLICES{,_EN}.md` P61 section stays in sync.

## 7. Fact pack for the next version

- ① Run first: the **R-08.5 criterion revision slice** — define the waiver conditions for an "inherited drift": same baseline + same assertion + a stricter drift direction + unchanged `verdictClass` / `chainState` / `taskState` / `stepState` + traceable cause (the `02d5797` class).
- ② Run later: an **independent fix slice** — make the `tools/agent-stub` `hang` branch write an empty child report and re-run layer3 F2 until `failedChecks` is cleared.
- ③ v51 observation candidates: the supported values and known boundaries of `--effort` and `--context`; the version pin of `Invoke-AgentProbe.ps1` conflicting on hosts where 1.0.83 and 1.0.85 coexist (the run-2 blocking root cause).
- ④ Discipline reuse: a validation slice must not modify the harness under test (B4 and C1 precedents); an inherited drift must cite its source truthfully as either a "same-baseline re-run" or a "code-semantics inference", never as verified.
- ⑤ Pipeline numbering: in this slice ④ = gates and ⑤ = pre-review (v1.64 §5.1); the controller's previous receipt once wrote "④ pre-review", which was corrected and logged.

## 8. Review conclusions

### 8.1 ⑤ pre-review (`PRE-REVIEW: PASS`)

- Check faces A–I all PASS: layer 1's five items, layer 2 vs the baseline, the layer 3 mechanism legs, the real-candidate legs, no new direct call site, the Go side unchanged with the deviation registered, the three path-A technical reasons holding independently, six facts being a genuine 5+1 rather than a rename, the cost wording, the D-item registration, the four FINDING rulings landing, and encoding compliance.
- Disposition of the 4 LOW items: LOW-1 (run-5 ruling and authorization chain unregistered) closed; LOW-2 (run-4 and run-5 cost unrecorded) closed (run-5 = ¥0.01, ¥0.09 cumulative); LOW-3 (run-2 attempt-1 captured on console only) closed (attempt-2 and attempt-3 captured, the zero-model-call conclusion unchanged); LOW-4 (`test-r3.md` §6 and §9 state gap) closed by this report.

### 8.2 ⑥ independent scan (`INDEPENDENT SCAN: PASS`, sixstage fence transcribed verbatim)

Note outside the fence: below is the verbatim Section A–E transcription of `04-reviews/review-06-scan.md`; the original conclusion line stays where it was (Section A), and the fence's last line repeats the same value to satisfy G7-f's "last non-empty line = conclusion line" rule (both occurrences carry the same value and no other byte changed); the original's "controller note" section is not ⑥ output and is not transcribed.

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
(none)

## Section C
1. 层 1 五项达到要求：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r1.md 与 layer1/l1-01-wrapper-version.txt、l1-02-bare-version.txt、l1-03-unlock-version.txt、l1-04-byok-call.txt、l1-05-effort-probe.txt，五项结论均为 PASS，且版本钉版为 1.0.83，-Unlock 观测为 1.0.85，BYOK 退出码为 0，--effort=high 被判定为 supported。此处不存在用于绕过封装的直接调用；生产新增仅在 tools/agent-probe/CopilotCli.ps1 中出现 buildCopilotArgs 与 Invoke-CopilotCli。

2. 层 2 结论与 B4 基线一致：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r2.md 及 layer2/l2-baseline-compare.txt、l2-f1-driver.txt、l2-f4-driver.txt，F1 六腿与 F4 三腿的 verdict class/chainState/taskState/stepState/observations 全部与基线一致，唯一漂移仅是断言 id 从 five-facts-distinct 变为 six-facts-distinct；这属于事实枚举扩展，不构成 verdict 类不符，也不构成命名重命名伪装。该漂移的成因与 02d5797 的 pid report 绑定一致，且在证据文本中已如实登记。

3. 层 3 的 F2/F4 观察并不推翻总体结论：读取 docs/.../03-test/test-r3.md、controller-note-layer3-findings.md，以及 layer3/l3-baseline-compare.txt、real-candidate/run-5/argv.txt、stdout.txt、stderr.txt、verdict.txt，F2 的 process-tree-taken-away 从基线 true 变为当前 false，但对应证据说明其成因是 02d5797 之后的 pid report fail-closed 语义继承漂移；F2 的 verdict class 仍保持 uncertain-or-paused，且 run-5 的结论是 unsupported，且确实是宽白名单 + deny 后的 alias 旁路结果，而非“无证据直接判红”。这是真实的运行观察，不是伪造的结论。

4. 生产新增扫描口径成立：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt，生产 diff 仅新增 tools/agent-probe/CopilotCli.ps1；没有在工作区中发现其他直接使用 buildCopilotArgs 或 Invoke-CopilotCli 的新增生产调用点。这里的“既有 6 处调用点维持原样”在证据中没有被破坏，且新增文件工作方式与原设计保持一致。

5. Go/PS 命名、签名偏离的可接受性：本扫描未发现违反契约的“隐藏重写”。读取主控落盘的 docs/validation/r08-r09-ruling.md §2 及 tools/agent-probe/CopilotCli.ps1 自述部分，R-08.1 明确要求“封装函数 Invoke-CopilotCli / buildCopilotArgs”，而 Go 侧的命名差异（copilotCLIProbeArgs）位于不同语言边界、不同调用层；其差异已经在证据中登记为“命名/签名偏离，但不构成语义偏离”。这类偏离不影响生产封装的 fail-closed 语义，也不违背新增调用点必须走封装的约束。

6. Path A 的三个技术理由各自独立成立：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/00-briefing.md、r08-r09-ruling.md 及 03-test 证据，三项技术理由分别是：
   - 直接调用点收敛在单一生产文件；
   - 封装函数维持默认 --no-auto-update 的 fail-closed 版本锁定；
   - 真实候选腿的执行/停机/unsupported 结果都以具体 stderr、退出码、失败阶段呈现。
   这些理由彼此独立，任一项不成立时也不会自动让其他两项成立；但在本轮材料中三项都成立，因此没有发现“绕过封装”的有效证据。

7. 六事实断言不是改名：读取主控落盘的 docs/.../03-test/test-r2.md 与 controller-note-layer3-findings.md，six-facts-distinct 的真实语义是“新增第 6 枚事实：pid reports”，并非把原来的五枚事实重命名。该断言覆盖了原五枚事实，并以新增 pid reports 作为第六枚事实，因而与 B4 的 five-facts-distinct 相比是严格超集，不构成逻辑上“改名伪装”。

8. 费用口径与 D 项：读取主控落盘的 docs/.../03-test/controller-note-layer3-findings.md 及 real-candidate/run-5/verdict.txt，说明“本地 token 块不足以支撑费用归因，以用户 DeepSeek 实测为准”，且 run-5 实测为 ¥0.01，累计 ¥0.09。与此对应，D 项作为过程越权记录已被登记为 D=1，非自愈档；这属于过程治理问题，不构成生产代码/门禁语义缺陷。状态上仍是“已登记但不致命”，且不妨碍本轮独立扫描的 PASS 结论。

## Section D
- 若删掉 buildCopilotArgs 的默认 --no-auto-update 守卫：l1-01-wrapper-version.txt、l1-02-bare-version.txt、l1-04-byok-call.txt 会红；该守卫删除后，版本钉版检查与包装后的 argv 语义都失去可验证性；F1/F4 的 verdict 类不会因为此守卫删除而立刻变红，但其“可证明稳定版”的断言会失效。
- 若删掉 Invoke-CopilotCli 的 pin-check 与 -Unlock 受控集：l1-01-wrapper-version.txt 与 l1-03-unlock-version.txt 会红；run-2 的启动阶段失败也会失去其“版本钉死”的解释力；F1/F4 仍可能保持绿，但不再具有“同一版本基线”这一军规。
- 若删掉“新增调用点必须走封装”的守卫：生产 diff 可能出现直接 copilot 调用，区分机制腿与真实候选腿的证据链就会失去可追踪性；这个守卫一旦删掉，输入包中的“新增仅此文件”检查会失效，而 F1/F4 断言本身可能仍绿。
- 若删掉“ Path A 维持原样”的守卫：既有 6 处调用点不再被追踪，真实候选腿与生产封装之间的关系就不再可证明；F1/F4 不会直接红，但“新增和旧调用点没有污染现有调用图”这条约束会失效。
- 若删掉 six-facts-distinct 的断言：F1 正例的事实互异检查会从有效的 6 枚事实变成模糊的“至少 5 枚”，因此 F1 的正例断言会红；然而 verdict 类不会因为同一断言删除而立刻转红，因为它是增强性断言而非状态判定。
- 若删掉 run-5 的 unsupported 守卫：真实候选腿中宽白名单 + deny 旁路会被误判成“成功执行”，导致 unsupported 结论无法被证明；F1/F4 和层 1 门禁本身不过度受影响，但这条数据会直接失去证据证明。
- 若删掉 D 项登记：该问题不会立刻改变生产代码，但会使过程越权事实失去审计痕迹；这里的“红”会落在审计链而不是运行判据。
- UNCOVERED: 断电/物理机掉电场景未覆盖；本机 Win11/NTFS 环境下的跨主机或跨网络行为未覆盖；精确的服务端费用拆分没有本地 token 块可核对，因而只接受用户 DeepSeek 实测为准；SessionBridge/visible 相关外部交互不在本次生产 diff 范围。
- UNCOVERED: 真实候选腿的全部细粒度工具调用日志无法从本地 token 块还原出服务端计费细节；因此不能把本地估算当作唯一成本证据。

## Section E
我无法验证的部分：
- 断电/物理机失效场景未在本轮输入包中出现，因此无法独立证明该场景的行为边界。
- 本地 token 块不足以还原服务端精确计费，故只接受用户 DeepSeek 实测为准；我们无法从本地日志追踪每一笔服务端计费的细节。
- 若其他机器或不同 PowerShell host 下存在别的 `gci -Name` alias 展开行为，则本轮仅在当前环境中成立；输入包不覆盖全机型。

什么会推翻我的 PASS：
- 若新增生产 diff 进一步引入直接 `copilot` 调用，或把版本锁定/allowlist 语义改成隐藏绕过；
- 若 F1/F4 的实测证据被更正为“verdict class 不一致”，或 run-5 的 unsupported 结论被反驳为“实为成功”；
- 若 six-facts-distinct 被证明不是新增第 6 枚事实，而是重命名或对旧事实做了语义覆盖；
- 若后续出现同一环境中同一守卫被移除后仍无红的情况，说明本轮的“守卫可证伪”结论并非稳固；
- 若发现真实候选腿中的过程清理或证据链存在未登记的外部副作用，或 evidence pack 误删了关键材料，而不是仅 D=1 的过程治理问题。

SCAN COMPLETE
INDEPENDENT SCAN: PASS
```

### 8.3 ⑦ blind review and final review

- Blind review (`BLIND REVIEW: FINDINGS`): 1 HIGH — a failed check exists on the mechanism legs ⇒ conditional closure must not be claimed; two minimal-fix options were given (clear it and re-run layer3 / the user rules a revision of the R-08.5 criterion text). The first round carried no checklist at all (§7.7).
- Final review (`FINAL REVIEW: PARTIAL`): the four task verdicts = safety passed, architecture consistency passed, completeness passed, test-counterexample review partially passed; premise ① unmet, premise ② met; overlapping the blind review independently on the same HIGH.
- Disposition (ADR-036): the user accepted PARTIAL as the final conclusion, no re-review round was dispatched, and neither (A) nor (B) was closed inside this slice; ⑦ counts = 1 final + 1 blind (0 re-reviews, a user-authorized exception that must not become a precedent).

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/t027/CHECKLIST_v50.md	repo
docs/t027/CHECKLIST_v50_EN.md	repo
docs/validation/at23-restricted-candidate.md	repo
docs/validation/at23-restricted-candidate_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
tools/agent-probe/CopilotCli.ps1	repo
tools/gates/cjk-newwords.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/00-briefing.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/01-arch/byok-layer0-vs-repo.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/01-arch/design-r1.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/01-arch/layer0-prereq-2026-11-07.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/01-arch/layer0-result-2026-11-08.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/02-impl/controller-note-path-a.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/02-impl/impl-r1.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/02-impl/luna-connectivity.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/controller-note-layer3-findings.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-01-wrapper-version.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-02-bare-version.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-03-unlock-version.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-04-byok-call.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-05-effort-probe.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-baseline-compare.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-driver.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg-r2/f1-neg-r2-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg-r2/f1-neg-r2-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg-r2/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg/f1-neg-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg/f1-neg-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2-r2/f1-neg2-r2-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2-r2/f1-neg2-r2-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2-r2/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2/f1-neg2-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2/f1-neg2-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg3/f1-neg3-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg3/f1-neg3-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg3/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-pos/f1-pos-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-pos/f1-pos-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-pos/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-pack-verdict.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-assertions.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-driver.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/f4-neg-r2-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/f4-neg-r2-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/f4-neg-r2-p2-forge.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/f4-neg-r2-p2-forge.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/verdict-forge.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/f4-neg-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/f4-neg-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/f4-neg-p2-forge.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/f4-neg-p2-forge.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/verdict-forge.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/f4-pos-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/f4-pos-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/f4-pos-p2-recover.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/f4-pos-p2-recover.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/verdict-recover.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-pack-verdict.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-driver.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg-r2/f2-neg-r2-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg-r2/f2-neg-r2-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg-r2/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg/f2-neg-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg/f2-neg-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-pos/f2-pos-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-pos/f2-pos-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-pos/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-pack-verdict.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-driver.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg-r2/f3-neg-r2-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg-r2/f3-neg-r2-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg-r2/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg/f3-neg-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg/f3-neg-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-pos/f3-pos-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-pos/f3-pos-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-pos/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-pack-verdict.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-driver.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/f4-neg-r2-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/f4-neg-r2-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/f4-neg-r2-p2-forge.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/f4-neg-r2-p2-forge.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/verdict-forge.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/f4-neg-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/f4-neg-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/f4-neg-p2-forge.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/f4-neg-p2-forge.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/verdict-forge.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/f4-pos-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/f4-pos-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/f4-pos-p2-recover.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/f4-pos-p2-recover.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/verdict-recover.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-pack-verdict.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/postflight-rebuild/l3-six-fact-rebuild.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/command.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/stderr.raw-cp936.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-2-param-binding-failure.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-2-param-binding-failure.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-3.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-3.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/command.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/command.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/stop-call.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/stop-call.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/stop-evidence.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/argv.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/exit.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/verdict.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/argv.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/exit.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/verdict.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r1.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r2.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r3.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/04-reviews/review-05-pre.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/04-reviews/review-06-scan.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/04-reviews/review-07-blind.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/04-reviews/review-07-final.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-ci.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-index.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-tree.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/MANIFEST.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/freedom-list.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/input-list.md	repo
```
