# T013 Adapters 验证报告

[English](t013-adapters_EN.md)

日期：2026-09-07。结论：`T013 COMPLETE`。本任务实现 file-queue 与 sessbridge 两类适配器入口，统一 request/claim/result 信封与 dispatch/takeover receipt 的摘要规则，并完成 AT-08 离线关键场景验证。

## 实现范围

1. `internal/adapters/filequeue.go`：
   - 实现 request/claim/result 三类 envelope 的校验、摘要、编解码。
   - 实现 dispatch/takeover receipt 的校验、摘要、编解码。
   - 实现文件队列目录（requests/inflight/results/archive）与 no-replace 原子写。
   - 实现 claim 领取、generation 递增、fencing token 约束、result 提交与幂等冲突判定。
2. `internal/adapters/sessbridge.go`：
   - 定义消费侧 silent client 端口。
   - 实现 sessbridge dispatch 结果到 adapter dispatch receipt 的映射。
   - 覆盖 accepted/rejected/uncertain 结果与 requestId 绑定校验、传输错误证据补全。
3. `internal/adapters/{filequeue_test,sessbridge_test}.go`：
   - 覆盖 file-queue 与 sessbridge 关键回归测试。

## AT-08 覆盖

1. 信封一致性：request/claim/result 统一使用 `adapter-envelope.schema.json` 对应字段与 domain hash（request/claim/result 分域）。
2. 重复应用防护：同 `requestId` 同内容重投幂等返回；同 `requestId` 异内容冲突阻断。
3. 领取与 fencing：
   - 首次 claim 需要请求进入 inflight。
   - generation 必须递增。
   - `(claimId,generation)` 不匹配时拒绝 result。
4. takeover 约束：更换 claimId 的后续 claim 必须提供 takeoverReceiptHash。
5. sessbridge 边界：
   - `transport=sessbridge` 强制非空 transportRequestId。
   - busy/timeout/poll_timeout/lm_api_unavailable/未知状态统一记为 uncertain。
   - response `requestId` 与本地 requestId 不一致时 fail-close 为 uncertain。
6. 无 GUI 回退：本实现仅消费 silent client 端口，不引入 visible/GUI 自动回退路径。

## 测试用例

- `internal/adapters/filequeue_test.go`
  - `TestFileQueueDispatchClaimResultFlow`
  - `TestFileQueueDispatchIsIdempotentForSamePayload`
  - `TestFileQueueDispatchRejectsDifferentPayloadForSameRequestID`
  - `TestFileQueueClaimTakeoverRequiresReceipt`
  - `TestFileQueueSubmitResultRejectsFencingMismatch`
  - `TestFileQueueSubmitResultRejectsWrongRequestHash`
  - `TestWriteCanonicalLineNoReplaceRejectsExistingFile`
  - `TestFileQueueClaimGeneration1IsImmutable`
   - `TestFileQueueClaimGeneration1RejectsTakeoverReceiptHash`
   - `TestFileQueueClaimRenewRejectsTakeoverReceiptHash`
- `internal/adapters/sessbridge_test.go`
  - `TestSessbridgeDispatchAccepted`
  - `TestSessbridgeDispatchRejected`
  - `TestSessbridgeDispatchBusyReturnsUncertain`
  - `TestSessbridgeDispatchRequestIDMismatchReturnsUncertain`
  - `TestSessbridgeDispatchClientErrorGeneratesTransportID`
  - `TestSessbridgeDispatchRequiresClient`

## 门禁

1. `go test -count=1 ./internal/adapters`：通过。
2. `go build ./...`：通过。
3. `go vet ./...`：通过。
4. `go test -count=1 ./...`：通过。

## 审查修复

2026-09-07 审查发现并修复以下问题：

1. no-replace 发布不成立：`os.Rename` 在 Windows 与 Unix 上都会替换已存在目标文件（违反“已存在的正式文件不得覆盖”），此前 `fs.ErrExist` 回退分支不可达。已改为硬链接 no-replace 发布（目标存在时 `os.Link` 失败），临时文件清理为 best-effort；request→inflight 移动改为 link+remove 且失败回滚。
2. claim 续租语义遗漏：按照契约，初始 claim 与同 claimId 续租时 `takeoverReceiptHash` 必须为 `null`，此前实现允许非空值。现已在 `Claim` 路径强制拒绝此类输入，仅在 claimId 变化时允许并要求 `takeoverReceiptHash`。
3. 新增回归测试：`TestWriteCanonicalLineNoReplaceRejectsExistingFile`、`TestFileQueueClaimGeneration1IsImmutable`、`TestFileQueueClaimGeneration1RejectsTakeoverReceiptHash`、`TestFileQueueClaimRenewRejectsTakeoverReceiptHash`。

## 边界

1. 本切片未修改 `sessbridge` 仓库代码。
2. 本轮为离线契约与单元测试验证；授权宿主 smoke 按计划在具备授权环境时执行，不改变本切片实现结论。