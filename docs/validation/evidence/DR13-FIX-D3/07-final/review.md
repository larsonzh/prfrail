FINAL REVIEW: PASS

| SEVERITY | file:line | 描述 | 最小修法 |
|---|---|---|---|
| - | - | 未发现需整改项（本轮无 Medium+/Low 新增缺陷） | - |

1. fail-closed 组合安全  
结论：通过。  
依据：释放顺序已变为先删路径、再解锁、再关闭句柄，且关闭前后都不做重试，语义清晰且无死锁/句柄泄漏新路径。[internal/adapters/agent_runner_starting_marker.go](internal/adapters/agent_runner_starting_marker.go#L108) [internal/adapters/agent_runner_starting_marker.go](internal/adapters/agent_runner_starting_marker.go#L112) [internal/adapters/agent_runner_starting_marker.go](internal/adapters/agent_runner_starting_marker.go#L113) [internal/adapters/agent_runner_starting_marker.go](internal/adapters/agent_runner_starting_marker.go#L114)  
删除失败后会进入既有收敛链：同进程所有权门拒收（D1），路径复核（D2）避免孤儿 inode 误判，后续扫描再收敛。[internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L367) [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L391)

2. 竞态残差  
结论：通过（残差与既有口径一致，未引入新路径）。  
remove 成功后到 Close 前的 Windows delete-pending 窗口：扫描侧会先做路径复核，路径已消失即返回 false,nil，不会误回收。[internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L367) [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L373)  
扫描者持柄导致 remove 失败：已有 Windows 反测验证“返回 false 且锁必须释放”。[internal/adapters/agent_runner_starting_marker_order_windows_test.go](internal/adapters/agent_runner_starting_marker_order_windows_test.go#L55) [internal/adapters/agent_runner_starting_marker_order_windows_test.go](internal/adapters/agent_runner_starting_marker_order_windows_test.go#L77) [internal/adapters/agent_runner_starting_marker_order_windows_test.go](internal/adapters/agent_runner_starting_marker_order_windows_test.go#L96)  
“扫描者误回收他人 marker”未见新通道：清理分支在重开+加锁后做内容匹配再删路径，仍是“只删自己看过的记录”。[internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L530) [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L554) [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L567)

3. 测试可证伪性  
结论：通过（含一条非阻断观察项）。  
顺序门为确定性 AST 门，能卡住常规“unlock 早于 remove”的回归，且排除了 defer/go 伪顺序。[internal/adapters/agent_runner_starting_marker_order_test.go](internal/adapters/agent_runner_starting_marker_order_test.go#L154) [internal/adapters/agent_runner_starting_marker_order_test.go](internal/adapters/agent_runner_starting_marker_order_test.go#L188) [internal/adapters/agent_runner_starting_marker_order_test.go](internal/adapters/agent_runner_starting_marker_order_test.go#L216)  
Windows fail-closed 反测确实构造了无 FILE_SHARE_DELETE 占名（CreateFile share 掩码仅读写），并验证删失败后锁可再获。[internal/adapters/agent_runner_starting_marker_order_windows_test.go](internal/adapters/agent_runner_starting_marker_order_windows_test.go#L45) [internal/adapters/agent_runner_starting_marker_order_windows_test.go](internal/adapters/agent_runner_starting_marker_order_windows_test.go#L86)  
压力回归是概率门（文内已声明），本轮聚焦输出显示 281 次释放、0 次命中 W1 观测、8 次 residual delete，口径自洽。[tmp/DR13-FIX-D3/07-final/inputs/test-focused.txt](tmp/DR13-FIX-D3/07-final/inputs/test-focused.txt#L8)  
ShareFactProbe 仅记录不做断言，定位为探针而非门禁，当前可接受。[internal/adapters/agent_runner_starting_marker_order_windows_test.go](internal/adapters/agent_runner_starting_marker_order_windows_test.go#L118)

4. 边界合规  
结论：通过。  
输入 diff 显示仅 1 个生产文件改动，另有 2 个台账文件更新；未触碰 D1/D2 实现文件、winners==1 断言文件、CONTRACTS 文档。[tmp/DR13-FIX-D3/07-final/inputs/slice.diff](tmp/DR13-FIX-D3/07-final/inputs/slice.diff#L1) [tmp/DR13-FIX-D3/07-final/inputs/slice.diff](tmp/DR13-FIX-D3/07-final/inputs/slice.diff#L41) [tmp/DR13-FIX-D3/07-final/inputs/slice.diff](tmp/DR13-FIX-D3/07-final/inputs/slice.diff#L81)  
返回值语义未漂移，仍是“本次是否成功删除路径”。[tmp/DR13-FIX-D3/07-final/inputs/head-marker.go.txt](tmp/DR13-FIX-D3/07-final/inputs/head-marker.go.txt#L99) [tmp/DR13-FIX-D3/07-final/inputs/work-internal_adapters_agent_runner_starting_marker.go.txt](tmp/DR13-FIX-D3/07-final/inputs/work-internal_adapters_agent_runner_starting_marker.go.txt#L116)

5. 回归完整性  
结论：通过（条件：双平台 CI 由后续环节继续兜底）。  
本输入包内证据显示：门禁 tree 总失败 0、base go test 通过、6 个聚焦用例全绿。[tmp/DR13-FIX-D3/07-final/inputs/gate-tree.txt](tmp/DR13-FIX-D3/07-final/inputs/gate-tree.txt#L48) [tmp/DR13-FIX-D3/07-final/inputs/gate-tree.txt](tmp/DR13-FIX-D3/07-final/inputs/gate-tree.txt#L49) [tmp/DR13-FIX-D3/07-final/inputs/test-focused.txt](tmp/DR13-FIX-D3/07-final/inputs/test-focused.txt#L16)  
本机未执行 unix 专腿这一点属已知事实；但本次改动函数本身无平台分支，平台差异点（reclaim/open sharing）仍由现有 windows/unix 实现约束，未见新增平台回归面。[internal/adapters/agent_runner_starting_marker_unix.go](internal/adapters/agent_runner_starting_marker_unix.go#L60) [internal/adapters/agent_runner_starting_marker_windows.go](internal/adapters/agent_runner_starting_marker_windows.go#L87)

6. 台账与证据链  
结论：通过。  
P48 条目为进行中、R/D=0、临时件路径声明与输入包一致。[docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L81) [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L1466) [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L1483) [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L1484)  
G9-e 两条是 INFO（未暂存提示），且 TOTAL_FAIL 仍为 0；不构成本片放行阻断。[tmp/DR13-FIX-D3/07-final/inputs/gate-tree.txt](tmp/DR13-FIX-D3/07-final/inputs/gate-tree.txt#L42) [tmp/DR13-FIX-D3/07-final/inputs/gate-tree.txt](tmp/DR13-FIX-D3/07-final/inputs/gate-tree.txt#L43) [tmp/DR13-FIX-D3/07-final/inputs/gate-tree.txt](tmp/DR13-FIX-D3/07-final/inputs/gate-tree.txt#L49)

可证伪性审计与推翻条件  
1. 若后续证据表明本次 diff 还包含第 2 个生产文件改动，则本结论立即推翻（边界失真）。  
2. 若双平台任一 CI 腿出现与本改动相关的 targeted 失败（特别是 release 顺序/RecoverStarting 路径），本结论推翻并回到 FINDINGS。  
3. 若能构造“同进程仍持有该 launch 句柄却被 RecoverStarting 回收”的复现实验，本结论推翻（D1 失效）。  
4. 若能构造“路径已切换到他人记录却被清理分支删除”的复现实验，本结论推翻（clear 只删本记录的证明失效）。  
5. 观察项（非阻断）：顺序门当前按最早调用位置判定，理论上存在不可达分支前置调用的假绿形态；若出现该类反例且行为测试未拦住，应把该观察项升级为缺陷并补 CFG/可达路径级断言。[internal/adapters/agent_runner_starting_marker_order_test.go](internal/adapters/agent_runner_starting_marker_order_test.go#L192)
