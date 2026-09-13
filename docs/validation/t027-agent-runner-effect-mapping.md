# T027 AgentRunner effect mapping 前置切片验证报告

[English](t027-agent-runner-effect-mapping_EN.md)

日期：2026-09-13。结论：`PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`。

## 实现范围

1. 冻结 v1 canonical mapping body：`{"local-process":"local-discardable","workspace-write":"local-discardable"}`；域为 `proofrail:agent-runner-effect-mapping:1\n`，SHA-256 为 `sha256:2c9609ee374bfc79bd0ed997aea6cfe4865c82ed0216d1a4a288139e0b0c92d4`。
2. `agent-runner-request` 强制绑定 `effectMappingVersion="1"` 与上述 `effectMappingHash`。映射是实现内固定契约，不新增 runtime mapping record。
3. Schema、negative goldens、canonical vector 与 semantic validator 覆盖缺失、未知或篡改绑定、未知 operation、grant class 覆盖不足，以及 grant class 超集准入。
4. Go request 构造/严格解码验证同一固定绑定；composite admission 将每个 `allowedEffects` operation 映射为授权 class，并在 request replay identity 写入和 dispatch 之前验证 grant 覆盖。

## Fail-closed 结果

- `local-process` 与 `workspace-write` 仅映射到 `local-discardable`。
- 缺失/未知/篡改 mapping version 或 hash、未知 operation 均在 request record 验证阶段拒绝。
- 有效 authorization grant 未覆盖映射 class 时在 authorization admission 阶段拒绝；拒绝后使用同一 replay index 的健康请求仍可首次准入，证明 request identity 未被消费。
- grant `effectClasses` 可包含 `read-only`、`external-write` 等额外 class，只要包含全部所需 `local-discardable` 即可通过。

## 验证结果

1. `node tools/contracts/contracts.test.js`：3/3 通过；36 份 Schema、123 个 fixture、4 条 canonical vector。
2. `go test -count=1 -json ./internal/adapters`：124 通过、0 失败、0 跳过。
3. `go build ./...`：通过。
4. `go vet ./...`：通过。
5. `go test -count=1 -json ./...`：13 个含测试包共 393 通过、0 失败、2 跳过（`TestApplyRejectsSymlinkParent`、`TestCaptureSymlinkValid`）；另 3 个包无测试。
6. 对修改的 Go 文件执行 `gofmt`：通过。

## 边界与剩余阻断

1. 本切片没有实现或验证真实工具/网络 enforcement，固定候选仍不可 dispatch。
2. 持久 dispatch/completion ledger 与 terminal receipt 收敛尚未接入 AgentRunner port。
3. resume 的 prior completion/session/workspace 连续性验证尚未实现。
4. 未调用真实模型 CLI，未读取 SecretStore，未访问网络，未执行 commit/push/publish。
5. T027 保持 `BLOCKED / NOT IMPLEMENTED`；本报告不声称 AT-23 通过。
