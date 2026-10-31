# ⑥ 独立扫描（原样归档）

> 载体：runSubagent(agentName="independent-reviewer", model="MAI-Code-1.1-Flash (copilot)")。transcript 偏移 117734。

## A 结论

INDEPENDENT SCAN: PASS

## B 缺陷

| 编号 | 严重度 | file:line | 缺陷描述 | 修法 |
|---|---|---|---|---|
|  |  |  |  |  |

## C 核验结论

1. 成立 — C7 三级不可混用，且包外无法实现 `Readiness`。`Readiness` 是 sealed interface，只有同包的私有 tag 方法才允许满足；`GoalProbe`/`GoalAcceptance`/`GoalProduction` 只能通过 `EvaluateT027BReadiness` 返回，见 [internal/gates/readiness.go](internal/gates/readiness.go#L69-L212)。`GoalAcceptance` 会在 `Seal == nil` 时先拒绝探针证据，`GoalProduction` 会在 `Seal == nil` 时直接返回 `ErrReadinessAcceptanceMissing`，而不看 `ProbeEvidence`，这正是 R-1 的 fail-closed 行为，见 [internal/gates/readiness.go](internal/gates/readiness.go#L173-L212)。`ReadinessInput` 以哈希串和 `Seal` 绑定表达 L2/L3，未形成跨层类型别名；真正的隔离由 sealed interface + 三个方法集差异保持，见 [internal/gates/readiness.go](internal/gates/readiness.go#L81-L138) 与 [internal/gates/readiness_test.go](internal/gates/readiness_test.go#L20-L48)。

2. 成立 — N4 是严格一次性消费，且 `Require`/`Complete` 不存在“重置额度”路径。`Require` 仅可从 `Unused -> Claimed` 走 CAS；任何失败、超时、重启、O_EXCL 冲突都落到 `ErrProbePermitInUse`/`ErrProbePermitConsumed`，并把状态推进到 `Spent`，见 [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L129-L203)。`Complete` 仅能从 `Claimed -> Spent` 或 `Spent -> Spent`，且拒绝未要求过的 `Unused`，没有“Spent 复活”路径，见 [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L216-L243)。`LoadProbePermit` 对“缺失/存在/损坏”做 fail-closed：缺失返回 `Unused`；存在/损坏返回 `Spent`，不允许发放第二次配额，见 [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L129-L171)。空 `path`、`ctx` 取消、并发 CAS 失败、`O_EXCL` 冲突四类边界在实现中均有自洽语义，见 [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L184-L285)。

3. 成立 — `EvaluateT027BReadiness` 是纯函数。它只读取 `goal` 和 `input`，并在 `GoalAcceptance`/`GoalProduction` 调用 `SealBinding.validate` 与正则校验，没有任何文件 I/O、没有写状态、没有消费 permit，见 [internal/gates/readiness.go](internal/gates/readiness.go#L157-L212)。测试也做了目录快照对比，断言评估前后无副作用，见 [internal/gates/readiness_test.go](internal/gates/readiness_test.go#L171-L214)。这也验证了 R-2 对“纯函数/哈希串签名”的实现方向是成立的。

4. 成立 — 哨兵完整性与 `errors.Is` 可达。`ErrProbePermitMissing`、`ErrProbePermitInUse`、`ErrProbePermitConsumed` 在 [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L55-L67) 定义；`ErrReadinessProbeNotAcceptance`、`ErrReadinessEvidenceIncomplete`、`ErrReadinessAcceptanceMissing`、`ErrReadinessProductionAT23Required` 在 [internal/gates/readiness.go](internal/gates/readiness.go#L19-L49) 定义。`GoalAcceptance` 和 `GoalProduction` 的返回都用 `%w` 包装，未出现“把 permit 错误包成 readiness 错误”或反向错误，见 [internal/gates/readiness.go](internal/gates/readiness.go#L173-L212)。

5. 成立 — 测试强度是有效的，且没有“镜像实现”式的无意义断言。编译期 pin 在 [internal/gates/readiness_test.go](internal/gates/readiness_test.go#L20-L48) 明确锁定了 sealed interface 与三类访问器；`TestReadinessLevelsCannotBeMixed` 仍有一个运行期断言，但仅用于正向 L1 分支的校验，见 [internal/gates/readiness_test.go](internal/gates/readiness_test.go#L108-L140)。它不构成跨级混用通道；真正的“不可混用”由编译期 pin 保证。该测试文件中也没有把实现逻辑复制进断言。

6. 成立 — 边界和污染控制成立。`git status --short --untracked-files=all` 只枚举了本 slice 允许文件：四个 gates 实现/测试文件、两个台账文件以及证据目录中的载荷。未出现 any sibling package drift；没有出现 [internal](internal) 之外的越界改动。允许集合包括 [internal/gates/readiness.go](internal/gates/readiness.go)、[internal/gates/probe_permit.go](internal/gates/probe_permit.go)、[internal/gates/readiness_test.go](internal/gates/readiness_test.go)、[internal/gates/probe_permit_test.go](internal/gates/probe_permit_test.go)、[docs/validation/evidence/P44-T027B-GATE-READINESS/00-briefing.md](docs/validation/evidence/P44-T027B-GATE-READINESS/00-briefing.md)、[docs/validation/evidence/P44-T027B-GATE-READINESS/freedom-list.md](docs/validation/evidence/P44-T027B-GATE-READINESS/freedom-list.md)、[docs/validation/evidence/P44-T027B-GATE-READINESS/MANIFEST.md](docs/validation/evidence/P44-T027B-GATE-READINESS/MANIFEST.md)、[docs/validation/evidence/P44-T027B-GATE-READINESS/SHA256SUMS.txt](docs/validation/evidence/P44-T027B-GATE-READINESS/SHA256SUMS.txt)、[docs/validation/evidence/P44-T027B-GATE-READINESS/01-arch-design/review.md](docs/validation/evidence/P44-T027B-GATE-READINESS/01-arch-design/review.md)、[docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md)、[docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md)。R-3 的“console 接线后置”结论也成立：本 slice 未宣称端到端闭合，且边界守住了“只读实现 + 后续装配切片”，见 [docs/validation/evidence/P44-T027B-GATE-READINESS/01-arch-design/review.md](docs/validation/evidence/P44-T027B-GATE-READINESS/01-arch-design/review.md#L131-L196)。

7. 成立 — 编码硬规则满足：四个新 Go 文件均为 UTF-8 无 BOM + LF。已直接检查文件字节，结论为 `BOM=False`、`LFonly=True`、`CRLF=False`，且 `go test ./internal/gates -count=20 -cpu 1,2,4` 返回 `ok ... 25.351s`。这覆盖了本 slice 的关键并发/语义验证；R-4 仍是环境限制，未在本机执行 `-race`，但不是实现缺陷。

## D 未覆盖面

- `SealBinding.validate` 的分支缺少单元覆盖：`Attempt < 1`、`TaskKey` 为空、`SealHash` 非法值的三条边界都未被单独断言；这些是 L2 质量门禁的重要失败路径。
- `LoadProbePermit` 对 `os.ReadFile` 的非 `NotExist` 失败路径未覆盖，例如权限错误或目录路径；这会决定“读取失败时是否继续返回 permit”。
- `Require` 的 `ctx == nil` 分支和 `Complete` 的“未知状态”分支尚未被直接测试；这属于强制 fail-closed 的边界条件。
- `GoalProduction` 对非空但非法格式的 `AT23Hash`（如 `abc` 或 63/65 位）未覆盖；其正则拒绝与 `ErrReadinessEvidenceIncomplete` 构成 L3 门禁。
- `GoalAcceptance` 对非法 `SealBinding`（带值但 hash 非法）未单独覆盖；目前仅验证了“没有 seal”与“没有 at23”，未专门锁定 `validate` 的格式篡改门。

## E 可证伪性

本 PASS 结论可被以下反例推翻：

- 任何跨包实现 `Readiness` 的类型能通过编译，或 `GoalAcceptance`/`GoalProduction` 在缺 L2 时返回 nil error，说明 sealed interface 或 fail-closed 判定失效。
- `Require` 在失败、超时、重启后仍能重新返回有效 claim；或 `Complete` 能把 `Spent` 状态复活；或 `LoadProbePermit` 对损坏标记返回 `Unused`，则 N4 一次性约束被否定。
- `EvaluateT027BReadiness` 产生任何文件写入、调用 `os.*`，或 `SealBinding.validate`/hash regex 允许伪造值，说明纯函数性不成立。
- 任何包外代码能绕过 `GoalProduction` 的 `Seal == nil` 语义并拿到生产 readiness，也会推翻 R-1 与 C7 语义。

§A INDEPENDENT SCAN: PASS