# P84 ⑤ 预审 r1 逐字转录 · 2026-11-21

- 载体：`deep-reasoner`（只读）；模型 `DeepSeek V4 Pro (deepseek)`。
- 额度：⑤ = 1 终审 + 1 复审（本轮 = 第 1 轮）。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。

---

# P84 `BOUNDDIR-RESIDUAL` ⑤ 预审报告（deep-reasoner / DeepSeek V4 Pro，只读）

## 结论

**PRE-REVIEW: PASS WITH FIXES**（仅 5 项 Low、均不阻断；无 Medium 及以上）

---

## 逐项核验（§4.2.1 最低动作）

### A. T4 走链（`bounddir_windows.go`）— 通过
- **DesiredAccess / class**：`accessListOrRead|accessReadAttributes|accessSynchronize`（:473-475、:629-633）与注释、Rev 1 逐字一致；`fileIdBothDirectoryInfo = 10` 正确（:486-492，GetFileInformationByHandleEx 枚举 FileIdBothDirectoryInfo 确为 10，非 NtQueryInformationFile 的 3），② 实证口径可信。
- **NTSTATUS→Win32**：`ntStatusError`（:665-676）经 `RtlNtStatusToDosError`，映射零值回退 ERROR_INVALID_PARAMETER(87)——未映射状态拒绝而非放行，fail-closed ✓；非零警告态亦拒绝（fail-closed）✓。
- **中间 reparse**：末级一律拒绝（:539-546，消息含 "reparse point"，兼容 `p83_fixround8_windows_test.go:156` 断言）；中间跟随 + `followIntermediateReparse`（:580-599）三重证明（目录类型 + `confirmOpenedArtifact` 的非 reparse/卷号/绑定根仍名/Ordinal 前缀）。P52 契约双向兼容：界外拒绝含 "outside the bound directory"（`bounddir_escape_windows_test.go:74` ✓）、界内放行 ✓；Rev 1 口径实现到位。
- **句柄纪律**：`openDescendantRelative`（:525-560）每步关父、失败路径零泄漏（reparse 句柄即关、follow 失败关 target、`dir.handle` 永不误关）；`BoundDirNames` 双路径均关（:197-201）；返回句柄唯一移交给调用方 ✓。
- **Ordinal 两处**（:683-693、:722-730）：`isUnderDirectory` 改字节相等——真后代同 API 同形态前缀必然字节相等，不误拒；case 孪生被拒 ✓。`boundDirVerify` 非字节相等落 `sameObjectAs` 对象比较，同对象（含大小写不同拼写）仍放行、孪生拒绝——方向正确，无残留放行窗口。
- **根 rel 语义**：`boundRootName`（`bounddir_clean.go:46-48`）三实现一致（win :95-98/:150-153、unix :140-143/:265-268 使用；darwin 不解析名、零改动）✓；`OpenBoundDirFile("")/(".")` 同走 `errBoundRootArtifact` 同一错误 ✓（无既有测试断言旧消息，grep 确认）。

### B. T5#4 前驱链（`runtime_supervised.go` / `runtime_agent.go`）— 通过
- **职责边界互补**：`verifyPublishLeaseHistory`（`runtime_agent.go:2206-2272`，多代链走，仅批量改写用）与 `verifyPublishLeasePredecessor`（:2166-2186，严格 N-1，仅同进程追平用）无重叠漏洞；History 覆盖 adoption（Reclaim 后 G-1 命中）与多代 skew ✓。
- **报错非跳过**：失败项错误含 proof 路径 + 检查名（`runtime_supervised.go:1234-1236`），`writeSupervisedJSON` 仅在验证通过后执行——失败项字节不改写 ✓；已通过项为验证过的真实历史，重跑幂等收敛 ✓。
- **三调用点**（:1148、:1164、:1940）均建 store + 传 ctx ✓；claim 分支传零值 previous 仅用于 `refreshed == previous` 短路，语义正确 ✓。
- **恢复不变量**：三条路径（同进程 `renewPublishLease:2104-2111` 追平 / 跨进程 Reclaim+批量 / 批量中途崩溃幂等重跑）覆盖 ✓；`:2104`（非本进程拒绝追平）、`:2120`（Reclaim 停机证明）、`:2202`（holder 复核）、`:1961-1963`（资源/consumer 归属）经现状全文核对**原样未放宽** ✓。
- 链走读 `RetainedLease` 带资源匹配 + hash 自证（`lease.go:366-380`），伪造 hash/token/异 run 不可达链 ✓。

### C. 测试可证伪性 — 通过
8 用例反例均命中所声明守卫：W1（拼写直测 + fsutil case-sensitive 端到端，skip 条件明确）、W2 三 leg（artifact/enum/excursion——excursion leg 是中间证明的分离形状）、W3 根拼写、C1（异 claim/未保留 hash + 字节不变断言 + 控制组）、C2（两代 skew + promotion 验证）、C3（子进程自 exec，PID 不可检查即 skip）。7 项变异映射逐一写在测试注释且与守卫一一对应；登记项锚点核实存在（W4→`gates/p83_fixround5_tester_test.go:37-63`；C4→`console/p83_fixround8_tester_test.go` T5#4 腿；C5→`console/p83_fixround5_tester_test.go:419/:432`）。

### D. 边界与回退 — 通过
git status 实核：8 个已跟踪文件（6 Go 生产/测试 + 2 台账文档）+ 2 个新测试 + 证据目录，无 CONTRACTS/schema/go.mod/CI 改动 ✓。P83 守卫未回退：`confirmOpenedArtifact`/`confirmEnumeratedDirectory` 保留、末级与枚举目标 reparse 拒绝保留、`lease_open_*.go` 未动、`boundResolveMask` 未动、`retainGeneration` 未动 ✓。Unix 仅分类器 + 注释（唯一行为变化 = `BoundDirNames(dir,".")` 从拒绝变枚举根——正是 T4 项 4 修复目标）；darwin 零改动 ✓。

### E. 证据件 — 通过
实测：4 载荷 SHA256 全部复算 OK；`.md` 均 BOM+LF、`SHA256SUMS.txt` 无 BOM+LF（合规）；MANIFEST 条目集与磁盘三方相等 ✓。

## Low 发现（不阻断，建议整改）

1. **L-1** `bounddir_windows.go:222-240`：`fileIDBothDirInfo` 在 GOARCH=386 对齐不匹配（Go 386 int64 对齐 4 → FileId offset 94 ≠ 内核 96）。交付范围仅 `windows/amd64`（CONTRACTS.md:224），无实际影响；建议加 `//go:build windows && amd64` 或运行时 `Offsetof` 断言 fail-closed。
2. **L-2** `bounddir_windows.go:315-334`：`readDirNames` 未校验 `offset+nameOffset+FileNameLength ≤ size`，畸形内核响应可越界读。建议加界检查（fail-closed），成本一行。
3. **L-3** `bounddir_windows.go:113-115`：`os.PathError.Path` 由绝对 `target` 改为相对 `clean`，Op/Err 契约不变、无断言依赖（grep 确认）；建议报告登记该可观察变化。
4. **L-4** `runtime_agent.go:2180-2182`：链走上限 64 在 TTL=24h、renewBefore=12h 下需 run 连续存活 ≥32 天才可达；方向 fail-closed，但触及即死结。建议文档登记边界或后续放宽。
5. **L-5** C3 跨进程测试含 skip 路径（host 边界，方案书已声明降级登记）；7 项变异「实测变红并 restore」仅凭 ③ 声称，建议 ⑧c 将变异证据件随包留存。

## 已核实 / 未核实

- **已核实**：方案书（含 Rev 1）、全部 diff 与现状代码全文（3 个 bounddir 文件、supervised/agent 关键段、lease.go、darwin/unix 现状）、P52 与 P83 契约测试断言兼容、三调用点与四处接管条件原样、测试反例与变异注释映射、git 变更面、门禁读数文件、证据件编码与 SHA256SUMS 自洽。
- **未核实**：未重跑任何测试（只读约束；④ 门禁/TEST-ALL 为磁盘读数）；变异实测未重跑；C3 的 CI 稳定性与 W1 本机 leg 是否实际 skip 无法从全绿摘要区分；L-1 对齐仅静态分析。
