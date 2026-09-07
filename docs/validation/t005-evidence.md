# T005 证据层验证报告

[English](t005-evidence_EN.md)

日期：2026-09-07。结论：`T005 COMPLETE`。本报告仅覆盖 AT-04 的证据层；journal、写入事务、回滚和崩溃点属于 T008。

## 实现范围

- `canonical.go`：严格 JSON 读取、重复键/非法 UTF-8/未配对 surrogate 拒绝、JCS UTF-16 键排序、ECMAScript 数字表示及域分离 SHA-256。
- `event.go`：`state-event` 构造/解码、版本和字段验证、摘要连续性、sequence 连续性、按实体投影的状态转换验证。
- `receipt.go`：`evidence-manifest` 构造/解码、不可变根摘要、item 唯一性与固定排序、kind/media type/redaction 约束。
- `verify.go`：从调用方提供的只读内容存储重新读取对象，验证引用闭包、内容摘要、canonical JSON、事件链和 attempt/run 绑定；只返回 `passed|failed|incomplete` 结构化检查，错误只使用冻结目录。

## 失败关闭证据

测试覆盖冻结 canonical/state-event 向量，以及重复键、尾随值、非法 UTF-8、surrogate、负零、错误版本、未知字段、manifest 篡改/乱序/重复、缺失对象、摘要不符、非 canonical JSON、事件 gap/reorder/corruption、跨实体状态不连续、manifest 顺序与事件 sequence 不同、run/attempt 不匹配和 nil reader。缺失能力返回 `incomplete`；结构、摘要或链冲突返回 `failed`。

## 验证命令

```powershell
go test ./internal/evidence
go vet ./internal/evidence
go test -cover ./internal/evidence
```

包测试和 vet 通过，语句覆盖率 79.5%。Windows race 构建未执行：当前主机未安装 `gcc`，`go test -race` 因 cgo 编译器缺失而停止；普通测试随后通过。最终仓库级门禁结果以任务结束时执行记录为准。

## 边界

本实现不生成完整 `verification-report` 持久记录，也不声称完成 RFC §16.9.26 的 identity/signature/review/promotion/snapshot/queue/lifecycle 全检查集；这些检查随所属后续任务接入。T005 不写状态、对象或授权，不把哈希完整性等同于身份可信。
