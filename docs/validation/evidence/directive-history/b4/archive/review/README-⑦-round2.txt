REVIEW INPUT BUNDLE - ⑦ final review, ROUND 2 (B4 slice + CONTRACTS §7 wording revision)
=========================================================================================

You are the independent final reviewer (⑦, GPT-5.3-Codex) for this slice. This is
ROUND 2. You did not see any previous review here: no ③/⑤/⑥ finding list is attached,
by design (directive §7.1). Judge the material below on its own merits.

Read only the material in this file. Everything you need is quoted verbatim; line
numbers are given for repo files you may open yourself under the repository root
D:\LZProjects\prfrail.

The four review axes you must answer, as in round 1: safety / architectural
consistency / completeness / test-counterexample review. Additionally, round 1 left
findings F1..F5 and recommendations R-1..R-4; section 5 lists them together with the
exact artifact that claims to address each. Verify or refute each claim; do not take
the master's word for it.

-----------------------------------------------------------------------------------------
SECTION 1 - WHAT THE USER RULED (2026-09-23), VERBATIM IN EFFECT
-----------------------------------------------------------------------------------------
The round-1 ruling "R-1 = the contract wording needs revision" triggered a
pre-declared hard stop in the delivery directive ("if ⑦ rules the contract needs
changing, stop and escalate: red line = changing protocol semantics"). The master
stopped and escalated. The user then ruled:

  (a) CONTRACTS §7 revision: AUTHORIZED, wording only, no code change. Unify the
      L205/L227 and L229 ordering statements to the MEASURED order
      (in-step gate -> postflight -> freeze/review/promotion). Do NOT touch
      engine.go - the measured behaviour is correct, changing code would turn
      correct into wrong and widen scope (red line). Protocol first (CN + _EN),
      plus gates, plus a ⑦ re-review.
  (b) B-1/B-3 wording: AUTHORIZED. B-1's claim scope binds to "unchanged OUTSIDE
      the capture exclusion set"; B-3 is annotated "applies only to this slice's
      platform refusal + substitute path, must not be extrapolated to production
      normal state". Both are additive descriptions, no semantic change.
  (c) Report annotation fixes: AUTHORIZED IMMEDIATELY, all four items from ⑦ round 1 -
      mark skipped legs as "non-equivalent substitute", mark the gate-port row
      "policy port uncovered", move the R-4 five items into the matrix in place,
      and add the sentence "this slice claims no equivalent coverage". ⑦'s reading
      prevails (it is the only final gate).
  (d) ⑦ quota: B4 and the CONTRACTS revision are counted separately. B4 = 1 (used)
      + 1 (re-review) = 2/2, inside quota.
  (e) OB-22 registered with corrected wording.
  (f) One commit, protocol first:
      docs(contracts+t027): unify the §7 gate/freeze ordering and annotate B4's
      non-equivalent fallbacks
  (g) The ⑤/⑦ disagreement on sufficiency is ALSO registered, as OB-23.

-----------------------------------------------------------------------------------------
SECTION 2 - THE DEFECT THAT TRIGGERED THE CONTRACT REVISION
-----------------------------------------------------------------------------------------
Claim (master, independently re-verified by the master before escalating):
`docs/CONTRACTS.md` contradicted itself about the in-step gate / postflight / freeze
order. Measured behaviour and the implementation follow L229, not L205/L227.

VERBATIM, `docs/CONTRACTS.md` BEFORE this round (line numbers are the pre-round file):
  L205: "代理报告 completed、退出 0 或最终文本只表示外部执行结束。ProofRail 必须先
        证明进程树停止，重新扫描 workspace，校验范围/秘密/副作用，冻结 candidate，
        再独立运行声明的 build/test/verify gates；通过后仍进入独立 review/promotion。..."
  L227: "...completed 只允许 step `TERMINAL_PENDING→PASSED`，任务仍依次经过 freeze、
        gates、review 与 completed promotion 才可 PASSED，chain 仅在全部任务接受后
        COMPLETED；..."
  L229: "...只有 postflight 通过才写 `REVIEW_PENDING`，随后依次经过既有 freeze、
        gates、review 与 completed promotion，task 才 PASSED，chain 仅在全部任务
        接受后 COMPLETED；..."

Implementation evidence that the measured order is "in-step gate -> postflight ->
freeze -> review -> promotion":
  `internal/chain/engine.go`
    - line 186: `for _, step := range task.Steps { ... engine.runStep(...) }`
      (the gate hook step is a member of the task step list and runs here)
    - line 197: `postflightInputs, postflightReason, postflightErr := engine.runTaskPostflight(...)`
      followed by `engine.transition(ctx, entity, "REVIEW_PENDING", ...)`
    - line 364: `candidate, err := engine.options.Acceptance.Accept(...)` (the
      candidate freeze) then `engine.options.Reviewer.Review(...)`
  `tmp/b4/runs/f5-pos/events/state-events.jsonl`
    - seq 13: step `b4-verify-step` RUNNING -> PASSED, reason `step-passed`
    - seq 14: task STEPS_RUNNING -> REVIEW_PENDING, reason `postflight-passed`
    - seq 15: task REVIEW_PENDING -> PASSED, reason `task-accepted`
    (the three sequences are adjacent, i.e. the gate hook step precedes
    REVIEW_PENDING by exactly one event)

-----------------------------------------------------------------------------------------
SECTION 3 - THE REVISION AS APPLIED (verbatim AFTER-text)
-----------------------------------------------------------------------------------------
`docs/CONTRACTS.md` line 5 (header date line, revision trail added):
  "日期：2026-09-07；S1 架构契约基线（2026-09-23 修订 §7 的门禁/postflight/freeze
   次序措辞，统一为实测次序；依据见准则 §12.4 的 OB-22）。本文是 wire、Schema、
   状态转换、事件、receipt 与跨记录不变式的规范性权威；[项目建议书](...) 保留
   来源、范围和设计理由。..."

`docs/CONTRACTS.md` L205 AFTER (the changed clause only):
  "ProofRail 必须先独立运行声明的 build/test/verify gates（门禁挂钩步是任务步骤表的
   一员，先于 postflight 与候选冻结执行），再证明进程树停止、重新扫描 workspace、
   校验范围/秘密/副作用（postflight），随后冻结 candidate；通过后仍进入独立
   review/promotion。"

`docs/CONTRACTS.md` L227 AFTER (the changed clause only):
  "completed 只允许 step `TERMINAL_PENDING→PASSED`，任务仍依次经过 gates（门禁挂钩步）、
   postflight、freeze、review 与 completed promotion 才可 PASSED，chain 仅在全部任务
   接受后 COMPLETED；"

`docs/CONTRACTS.md` L229 AFTER (the changed clause only):
  "只有 postflight 通过才写 `REVIEW_PENDING`，随后依次经过既有 freeze、review 与
   completed promotion（门禁挂钩步已在步骤阶段判过），task 才 PASSED，chain 仅在全部
   任务接受后 COMPLETED；"

`docs/CONTRACTS_EN.md` mirrors (same three statements, same sentence boundaries):
  L199 AFTER: "Agent-reported completion, exit zero or final text means only that
    external execution ended. ProofRail must first run the declared build/test/verify
    gates (gate hook steps are members of the task step list and run before postflight
    and before the candidate freeze), then prove process-tree termination, rescan the
    workspace and validate scope/secrets/effects (postflight), and only then freeze a
    candidate; review/promotion still follows. ..."
  L221 AFTER: "...completed only permits step `TERMINAL_PENDING→PASSED`, while the task
    still passes gates (gate hook steps), postflight, freeze, review, and a completed
    promotion before PASSED, and the chain reaches COMPLETED only after every task is
    accepted; ..."
  L223 AFTER: "Only a passed postflight writes `REVIEW_PENDING`, after which the existing
    freeze, review and completed-promotion sequence runs (the gate hook steps have
    already been judged in the step phase) and the task becomes PASSED, the chain
    COMPLETED only once all tasks are accepted; ..."

Diff shape (git numstat, staged): docs/CONTRACTS.md +4/-4, docs/CONTRACTS_EN.md +4/-4
(i.e. three ordering sentences + one header line on each side; no other line touched).

No Go file was modified in this round. `git diff --cached --numstat` contains no
`internal/` path.

-----------------------------------------------------------------------------------------
SECTION 4 - THE REPORT ANNOTATIONS AS APPLIED (⒞, verbatim)
-----------------------------------------------------------------------------------------
File: `docs/validation/t027-at23-e2e.md` (205 lines) and its strict `_EN` mirror
(205 lines).

§7 matrix, the six "skipped" rows AFTER (CN):
  | 全片 | 生产派发路径（`AgentRunnerReplayDispatcher` 发布 R） | 跳过 | 平台常量：Windows 复现库发布耐久 `unproven` | **非等价替代**：记录 fail-closed 原文后，改用同一生产 launcher 直启；见附录 A |
  | 全片 | 生产终局发布（`AgentRunnerTerminalPublisher`） | 跳过 | 同上，请求记录缺失 | **非等价替代**：用公开记录构造器加 `ToChainAgentRunnerTerminal` 构造终局链；链侧路由未改 |
  | 全片 | 生产门禁端口对（`gates.ChainPort` 加 `GuardExecutor`） | 跳过 | 能力声明不匹配：`memory limit` 缺口，任何 hook 都不可执行 | **非等价替代**（**策略端口未覆盖**）：harness 用同一 `guard.RunManaged` 边界执行真实门禁命令并留证 |
  | 全片 | 真实 AI 候选启动与计费调用 | 跳过 | B1/B2 冻结：候选不可合法启动、计费证据被 §8.3 与 §8.4 阻断 | **非等价替代**：透明声明的固定确定性 CLI 作为受管负载；**不宣称候选兼容** |
  | 全片 | Linux 原生验证 | 跳过 | 本片平台范围为 Windows；原生 Linux 属 C1 | **非等价替代**：不采集、不推断（同位给出，见 §9.6） |
  | 全片 | SessionBridge `visible` 人工交互 | 跳过 | 属 T028 或 AT-24 | **非等价替代**：本片零覆盖（同位给出，见 §9.6） |

New paragraph immediately AFTER the matrix (CN):
  "**本片不主张等价覆盖**：上表「跳过」行的替代证据**不是**对生产语义的等价覆盖，
   只是本片可诚实执行部分的取证；「充分性」判定口径的分歧见 OB-23（准则 §12.4）。"
The _EN mirror reads: "**This slice claims no equivalent coverage**: the alternative
evidence in the rows marked skipped above is **not** equivalent coverage of the
production semantics, only the evidence this slice can honestly collect; the
disagreement about how sufficiency is judged is registered as OB-23 (directive §12.4)."

§1 conclusion row AFTER (CN):
  "| 跳过面的替代证据是否充分 | **不充分**——⑦ 以「等价覆盖生产语义」口径判五项均不充分
   （⑤ 以「诚实 fail-closed」口径判充分，口径分歧登记为 OB-23）；**本片不主张等价覆盖** |"

§9.7 accounting conclusion AFTER (CN), the paragraph that previously said "this slice
only records and does not change it; the ruling is referred to ⑦; if ⑦ picks contract
revision the master stops and escalates" is replaced by:
  "⑦ 终审已裁定 **R-1：契约措辞需修订**，依据是契约内部两处次序陈述不一致（L205/L227
   与 L229），而实现与 L229 一致。按用户 2026-09-23 裁决 ⒜，本片**只改契约措辞、不动
   代码**：`CONTRACTS` §7 三处次序陈述已统一为实测次序「步内 gate → postflight →
   freeze / review / promotion」（CN 与 `_EN` 同步），登记见准则 §12.4 的 **OB-22**。"

§10 artifacts block AFTER (added four rows, all `repo`):
  docs/CONTRACTS.md, docs/CONTRACTS_EN.md, docs/DELIVERY_DIRECTIVE.md,
  docs/DELIVERY_DIRECTIVE_EN.md  (plus the pre-existing rows)
Total staged artifact rows resolved by G7-b: 30.

Directive §12.4 registrations added this round (CN and _EN):
  OB-22 | 2026-09-23 | CONTRACTS §7 internal ordering inconsistency (found by the B4 ⑦
    final review): L205/L227 say "freeze the candidate, then run the gates" while L229
    says "postflight before REVIEW_PENDING"; the measured behaviour and the
    implementation both agree with L229 (evidence: internal/chain/engine.go:186/197/364
    and tmp/b4/runs/f5-pos/events/state-events.jsonl:13/14/15) => the implementation is
    correct and the contract text must be unified. | Disposed (2026-09-23, user ruling
    ⒜): the CONTRACTS §7 wording is unified to the measured order "in-step gate ->
    postflight -> freeze / review / promotion", wording only, no code change, CN and
    _EN kept in sync |
  OB-23 | 2026-09-23 | B4 ⑤/⑦ disagreement on "sufficiency of the alternative evidence
    for skipped faces": within one slice ⑤ judged it sufficient under the honest
    fail-closed reading while ⑦ judged all five insufficient under the
    equivalent-coverage-of-production-semantics reading => the criterion lacks one
    agreed reading; both readings are reasonable but not equivalent. Suggested: adopt
    ⑦'s equivalence reading and state explicitly that a "non-equivalent substitute" is a
    legal but must-be-labelled status. | To be discussed: evaluated at the next §12
    revision |

-----------------------------------------------------------------------------------------
SECTION 5 - ROUND-1 FINDINGS AND WHERE EACH IS CLAIMED TO BE ADDRESSED
-----------------------------------------------------------------------------------------
  F1 (Critical, contract internal order conflict) -> CONTRACTS §7 three statements
     unified (section 3); registered as OB-22; report §9.7 rewritten (section 4).
  F2 (High, the two production-publication legs use a NON-EQUIVALENT substitute path)
     -> both rows now carry the explicit "non-equivalent substitute" label; the matrix
     is followed by the sentence "this slice claims no equivalent coverage"
     (section 4). No behaviour change is claimed.
  F3 (Medium, the gate-port leg is not equivalently covered) -> that row now carries
     "non-equivalent substitute" AND "policy port uncovered".
  F4 (Medium, B-1 needs contract-layer wording) -> B-1's claim scope is now bound in
     the report to "unchanged outside the capture exclusion set" (deviation 3, §8,
     unchanged from round 1) and the boundary is explicitly declared as the claim's
     limit. NOTE: the master did NOT add a new normative sentence to CONTRACTS for
     this; the user ruling ⒝ authorizes it as an "additive description", and the
     master's reading is that the report-level claim boundary IS that description.
     Judge whether that satisfies F4 or whether a CONTRACTS sentence is still owed.
  F5 (Low, R-4's Linux/SessionBridge items were not given in place in the matrix)
     -> two new matrix rows add "Native Linux validation" and "SessionBridge visible
     human interaction" with status/reason/alternative evidence in place.
  R-2 (B-1 boundary needs contract wording) -> same as F4; report-level declaration
     plus OB-23 registration. Judge sufficiency.
  R-3 (B-3 empty `launches/` must not be extrapolated) -> report deviation 4 remains
     three-part with the corrected criterion (dispatch-event count + identity mirror);
     OB-23 records the reading disagreement. The master did NOT add a CONTRACTS
     sentence. Judge sufficiency.
  R-4 (all five skipped faces judged insufficient under equivalence) -> the user ruling
     (g) makes ⑦'s equivalence reading the operative one; the matrix now labels all six
     skipped rows and the slice explicitly claims no equivalent coverage.
  Everything previously disclosed about AT-23 is unchanged: AT-23 is NOT passed, the
  real-candidate leg is measurably unreachable today (OB-21), and the report's §1
  "AT-23 未通过" statement is not dressed up.

-----------------------------------------------------------------------------------------
SECTION 6 - GATE / BUILD EVIDENCE FOR THIS ROUND
-----------------------------------------------------------------------------------------
Command: node tools/gates/gate.js --all --scope=tree   -> TOTAL_FAIL=0
Command: node tools/gates/gate.js --all --scope=index  -> TOTAL_FAIL=0
Command: node tools/gates/gate.js --all --scope=ci     -> TOTAL_FAIL=0
Command: node tools/gates/gate.js --selftest           -> SELFTEST: PASS
Selected details from the tree-scope report (changed=14):
  PASS G3(a) :: docs/CONTRACTS.md +4/-4 vs _EN +4/-4 OK || docs/DELIVERY_DIRECTIVE.md
       +4/-0 vs _EN +4/-0 OK || docs/validation/t027-at23-e2e.md +204/-0 vs _EN
       +204/-0 OK
  PASS G3(b), PASS G4b (newCjk=589 unknown=[]), PASS G6-1..G6-4 (scanned 424
       added lines; no pending/completion conflict, no cost clause, no malformed run
       id, no expiring literal)
  PASS G7-a (2 reports OK), G7-b (30 artifact rows OK), G7-c, G7-d
  PASS base gofmt -l . (empty), go build ./... (exit 0), go vet ./... (exit 0),
       go test ./... (ok=14 FAIL=0)
Independent structural mirror check (master-side script, positional):
  docs/DELIVERY_DIRECTIVE.md <-> _EN: CN=1068 EN=1068, shape all true, per-line
    mismatches = 0
  docs/validation/t027-at23-e2e.md <-> _EN: CN=205 EN=205, shape all true, per-line
    mismatches = 0
  docs/CONTRACTS.md <-> _EN: CN=312 EN=304 - THIS PAIR IS NOT A POSITIONAL MIRROR AND
    NEVER WAS (pre-existing, 8-line difference). G3-a only requires equal numstat for
    this pair (+4/-4 both sides), which holds. Disclosed as a coverage note.

-----------------------------------------------------------------------------------------
SECTION 7 - WHAT THE MASTER COULD NOT VERIFY, AND WHAT WOULD OVERTURN IT
-----------------------------------------------------------------------------------------
  - The master did not re-run the B4 driver this round: the report's §5/§6 numbers and
    the tmp/b4/** pack are unchanged from round 1 and remain local-only.
  - The master cannot compute sha256 of its own report text against a repo revision
    that does not exist yet (nothing is committed).
  - Linux native validation is out of scope (C1) and is claimed nowhere.
  - If any statement in sections 3 or 4 misdescribes the measured order, or if the
    unified CONTRACTS wording still contradicts `internal/chain/engine.go`, that
    falsifies the central claim of this round.
  - If the CONTRACTS revision is judged to need a Schema/fixture/ADR change, say so:
    the master deliberately changed wording only and touched no Go file.

-----------------------------------------------------------------------------------------
SECTION 8 - REQUIRED OUTPUT (same contract as round 1, directive §7.2)
-----------------------------------------------------------------------------------------
  A. one conclusion line: `INDEPENDENT SCAN: PASS` or `INDEPENDENT SCAN: FINDINGS`
  B. findings, one per line: `SEVERITY | file:line | verbatim quote | minimal fix`
  C. evidence you used
  D. falsifiability audit (mandatory): for each claim, the condition that would refute it
  E. what you could not verify, and what would overturn your conclusion
