# P84 ⑤ 预审 r2（原位复核）逐字转录 · 2026-11-21

- 载体：`deep-reasoner`（只读）；模型 `DeepSeek V4 Pro (deepseek)`。
- 额度：⑤ = 1 + 1（本轮 = 第 2 轮，用尽）。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。

---

复核完毕。以下为第 2 轮原位复核结论。

---

```
PRE-REVIEW: PASS
```

**L-1（386 对齐 / FileName 偏移）— CLOSED**
- `internal/evidence/bounddir_windows.go:209-226`：struct 按内核 x64 布局注释逐字段标偏移（ShortName 止于 94、FileId 96..103、FileName 104）；我独立重算 x64 对齐（94→8 对齐到 96）与之一致，上轮 94 确系漏算 FileId。
- `:246-254` `assertFileIDBothDirInfoLayout`：`Offsetof(FileName)!=104` 或 `Sizeof<104` 即返回 ABI 错误，fail-closed；`:271` 为 `readDirNames` 入口首行调用，解任何条目前必过。不误拒合法条目：本机 amd64 `go test` 实测枚举正常。
- 附带核验：支持矩阵（`docs/S1_SUPPORT_MATRIX.md:9-10`）仅 Windows amd64 / Linux amd64，**无 386**；且 386 下偏移亦为 104（94→4 对齐到 96），断言不会误拒。

**L-2（解析越界）— CLOSED**
- `bounddir_windows.go:315-324`：解引用前 `offset+104>len(buf)` 拒绝；name 切片前 `offset+104+FileNameLength>len(buf)` 拒绝；`units==0` 拒绝；`NextEntryOffset==0` 收敛。缓冲 64KB 起步、1MB 上限（`:514-515`），内核成功应答时条目必然完整，正常目录不误拒；offset 递进为 uint32 有限值，无死循环、无新越界。

**L-3（Path 改相对 clean）— CLOSED（无代码动作已核实）**
- `bounddir_windows.go:114-116`：`&os.PathError{Op: "open", Path: clean, ...}` 仍为相对路径、未回退，与「登记在案、不改代码」口径一致。⑧a 报告文本本身非代码，代码面不可核。

**L-4（链上限 64 注释）— CLOSED**
- `internal/console/runtime_agent.go:2179-2189`：注释如实登记 TTL=24h/renewBefore=12h ⇒ 一天两续 ⇒ 连续存活 >32 天方可触及 64 代；触及为 fail-closed depth 错误；放宽需另片。口径与主控陈述一致。

**L-5（变异证据件留存）— CLOSED（登记面）**
- 无代码动作；`docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/` 目录已建。⑧c 收档动作不可从代码核实。

**回归确认（第 4 点）**
- 走链与中间 reparse 界内证明完整保留：`bounddir_windows.go:438-540`（openDescendantRelative）、`:545-600`（followIntermediateReparse）、`:600-700`（openRelativeComponent）语义与 P83 round 8 / P84 round 3 注释一致；批量前驱链 `runtime_agent.go:2155-2177`（verifyPublishLeasePredecessor 四道校验）完整。
- `go vet` 通过；`go test ./internal/evidence/ ./internal/console/` 全部 ok（evidence 2.463s / console 13.351s），新增 W1-W3 测试含 `BoundDirNames(dir,"")`/`(dir,".")` 实走 readDirNames。

**残余观察（Low，不阻塞）**：386 下 `offset+104+uintptr(FileNameLength)` 存在 uintptr 32 位溢出绕过第二道界校验的理论可能，但 386 不在支持矩阵内，仅登记不设 Fix。

**已核实**：L-1/L-2/L-3 代码面/L-4 注释/回归面（windows 半侧 + runtime_agent.go 注释处）。
**未核实**：L-3/L-5 的 ⑧a/⑧c 登记文本本身（非代码）；unix 半侧本轮未重审（属早前轮次已审通过范围）。
