# ⑦ 终审（原样归档）

> 载体：runSubagent(agentName="independent-reviewer", model="GPT-5.3-Codex (copilot)")。偏移 120591。

## A 终审结论
FINAL REVIEW: PASS

## B 缺陷
| 编号 | 严重度 | file:line | 缺陷 | 修法 |
|---|---|---|---|---|
无（未发现 Critical/High/Medium 真实缺陷）。

## C 核验结论
1. D1 是否真正关闭原竞态  
结论：是，且与 D2 组合后对本次“双赢家”路径形成结构性封堵。  
关键证据：同进程所有权门在 [internal/adapters/agent_runner_pinned_cli.go#L390](internal/adapters/agent_runner_pinned_cli.go#L390)；Recover 在 claim 前强制执行于 [internal/adapters/agent_runner_pinned_cli.go#L148](internal/adapters/agent_runner_pinned_cli.go#L148) 与 [internal/adapters/agent_runner_pinned_cli.go#L151](internal/adapters/agent_runner_pinned_cli.go#L151)。onSpawn 注册与 running 推进在 [internal/adapters/agent_runner_pinned_cli.go#L209](internal/adapters/agent_runner_pinned_cli.go#L209)、[internal/adapters/agent_runner_pinned_cli.go#L266](internal/adapters/agent_runner_pinned_cli.go#L266)，marker 关闭在 [internal/adapters/agent_runner_pinned_cli.go#L272](internal/adapters/agent_runner_pinned_cli.go#L272)。  
关于 Stop 失败分支：无条件 release 在 [internal/adapters/agent_runner_pinned_cli.go#L292](internal/adapters/agent_runner_pinned_cli.go#L292)、[internal/adapters/agent_runner_pinned_cli.go#L301](internal/adapters/agent_runner_pinned_cli.go#L301)，会导致“D1 依赖的内存在册信息提前消失”的残余风险；但该路径不再重现本次双赢家根因（后续 claim 仍受 live mirror 拒绝链约束于 [internal/adapters/agent_runner_pinned_cli.go#L605](internal/adapters/agent_runner_pinned_cli.go#L605)-[internal/adapters/agent_runner_pinned_cli.go#L616](internal/adapters/agent_runner_pinned_cli.go#L616) 与 slot 规则 [internal/adapters/agent_runner_pinned_cli.go#L712](internal/adapters/agent_runner_pinned_cli.go#L712) 起）。

2. D1 是否误伤孤儿收养  
结论：未见误伤。D1 依赖进程内 registry（lookup 定义于 [internal/adapters/agent_runner_process_registry.go#L145](internal/adapters/agent_runner_process_registry.go#L145)），跨进程崩溃恢复不在册，仍走原收养判定。  
覆盖面：  
- 通用 I5： [internal/adapters/spawn_recovery_test.go#L138](internal/adapters/spawn_recovery_test.go#L138)、[internal/adapters/spawn_recovery_test.go#L171](internal/adapters/spawn_recovery_test.go#L171)  
- Unix live-host refusal / orphan adopt： [internal/adapters/spawn_recovery_live_host_unix_test.go#L118](internal/adapters/spawn_recovery_live_host_unix_test.go#L118)、[internal/adapters/spawn_recovery_live_host_unix_test.go#L160](internal/adapters/spawn_recovery_live_host_unix_test.go#L160)  
- Windows crash/live-host： [internal/adapters/spawn_recovery_live_host_windows_test.go#L75](internal/adapters/spawn_recovery_live_host_windows_test.go#L75)、[internal/adapters/spawn_recovery_live_host_windows_test.go#L124](internal/adapters/spawn_recovery_live_host_windows_test.go#L124)  
说明：本机 Windows 无法本地执行 unix build-tag 用例，需 ⑨ Ubuntu 腿举证。

3. D2 分类是否 fail-closed 且不掩盖真错  
结论：fail-closed 成立；“类型层面区分竞争/IO”仍是既有契约折中。  
证据：锁后 stat，NotExist 直接拒绝不回收于 [internal/adapters/agent_runner_pinned_cli.go#L367](internal/adapters/agent_runner_pinned_cli.go#L367)-[internal/adapters/agent_runner_pinned_cli.go#L373](internal/adapters/agent_runner_pinned_cli.go#L373)；其他 IO 错误包装为 conflict 且保留原始错误文本于 [internal/adapters/agent_runner_pinned_cli.go#L375](internal/adapters/agent_runner_pinned_cli.go#L375)。冲突哨兵本身定义即“单跑根排他拒绝”总类于 [internal/adapters/agent_runner_pinned_cli.go#L27](internal/adapters/agent_runner_pinned_cli.go#L27)。

4. E1 是否可证伪且无假绿  
结论：是。  
反例测试： [internal/adapters/spawn_recovery_test.go#L221](internal/adapters/spawn_recovery_test.go#L221)。  
其构造是“真实活进程 + 真实 slot claim + 真实 registry 在册 + 无锁 marker 写入”：  
- 活进程来源 [internal/adapters/spawn_recovery_test.go#L89](internal/adapters/spawn_recovery_test.go#L89) 与 [internal/adapters/spawn_recovery_test.go#L228](internal/adapters/spawn_recovery_test.go#L228)  
- 真实 claim 于 [internal/adapters/spawn_recovery_test.go#L230](internal/adapters/spawn_recovery_test.go#L230)  
- registry 在册于 [internal/adapters/spawn_recovery_test.go#L246](internal/adapters/spawn_recovery_test.go#L246)  
- 四条断言锚点：拒绝不报错 [internal/adapters/spawn_recovery_test.go#L266](internal/adapters/spawn_recovery_test.go#L266)、不 reclaim [internal/adapters/spawn_recovery_test.go#L270](internal/adapters/spawn_recovery_test.go#L270)、不杀活进程 [internal/adapters/spawn_recovery_test.go#L278](internal/adapters/spawn_recovery_test.go#L278)、不清 marker 与不夺 slot [internal/adapters/spawn_recovery_test.go#L281](internal/adapters/spawn_recovery_test.go#L281)、[internal/adapters/spawn_recovery_test.go#L283](internal/adapters/spawn_recovery_test.go#L283)、[internal/adapters/spawn_recovery_test.go#L287](internal/adapters/spawn_recovery_test.go#L287)。  
旧实现下可被证伪（unix 先在“活进程”断言红，windows 在“marker/slot”断言红）。

5. 并发回归是否结构性成立  
结论：成立（同进程并发 Start 语义下）。  
并发测试锚点： [internal/adapters/agent_runner_pinned_cli_test.go#L857](internal/adapters/agent_runner_pinned_cli_test.go#L857)、[internal/adapters/agent_runner_pinned_cli_test.go#L901](internal/adapters/agent_runner_pinned_cli_test.go#L901)、[internal/adapters/agent_runner_pinned_cli_test.go#L904](internal/adapters/agent_runner_pinned_cli_test.go#L904)。  
实测：  
- ConcurrentStartsHaveOneWinner，count=20 通过  
- NeverReclaimsALaunchThisProcessStillOwns，count=20 通过  
- RecoverStarting 相关组，count=1 通过  
残余窗口：D3 未做（release 仍是 unlock 后 remove，见 [internal/adapters/agent_runner_starting_marker.go#L92](internal/adapters/agent_runner_starting_marker.go#L92)-[internal/adapters/agent_runner_starting_marker.go#L99](internal/adapters/agent_runner_starting_marker.go#L99)），以及跨 launcher/跨进程同 runRoot 场景，需继续登记。

6. 登记是否诚实（是否以登记代替修复）  
结论：当前证据中未见“掩盖性结论”；但 ⑧ 正式报告尚未落地，若缺登记则不能算完成闭环。  
已看到的诚实披露载体： [docs/validation/evidence/DR13-FIX/01-diagnosis/review.md](docs/validation/evidence/DR13-FIX/01-diagnosis/review.md)（已明确 D3 未做、W1/W2 与跨进程残余）。  
仍待 ⑧/⑨明确落地：B1（ubuntu 腿证据移交）、B2（D2 无专属确定性测试）、D3 defer、跨进程 W1 残留、⑥ §D 四项。

7. 是否构成放行  
结论：可进入 ⑧/⑨，前提是满足下节放行条件并完成证据登记。

## D 放行条件
1. ⑨ 必须提供 Ubuntu 腿证据，至少覆盖并发单赢家与 unix recovery 路径，并在日志中提取目标包级 ok 行。  
2. 在 ⑧ 报告显式登记 B1：本机无法执行 unix 用例，ubuntu 腿为最终判据。  
3. 在 ⑧ 报告显式登记 B2：D2 目前无专属“孤儿 inode”确定性测试，仅概率覆盖。  
4. 在 ⑧ 报告显式登记 D3 defer：为何本片不做“remove 再 unlock”，以及后续片计划。  
5. 在 ⑧ 报告显式登记跨进程 W1 残留与适用边界（同 runRoot、跨 launcher）。  
6. 在 ⑧ 报告显式登记并跟踪 ⑥ §D 四项（Stop 失败分支 releaseLaunch、W2 flock 语义未实证、unix StopProcessIdentity 本机未复现、跨 launcher/跨进程镜像失效窗口）。

## E 可证伪性
1. 去掉 D1 门（[internal/adapters/agent_runner_pinned_cli.go#L390](internal/adapters/agent_runner_pinned_cli.go#L390)-[internal/adapters/agent_runner_pinned_cli.go#L396](internal/adapters/agent_runner_pinned_cli.go#L396)）后，E1 会在至少一个断言点失败，且并发单赢家回归风险重现。  
2. 去掉 D2 stat 复核（[internal/adapters/agent_runner_pinned_cli.go#L367](internal/adapters/agent_runner_pinned_cli.go#L367)-[internal/adapters/agent_runner_pinned_cli.go#L375](internal/adapters/agent_runner_pinned_cli.go#L375)）后，锁住孤儿 inode 的路径可回归。  
3. 保留现实现状并复跑 count 压测（本次已做 count=20）仍稳定绿，说明“修复依赖调度运气”的证伪目前不成立。  
4. 若 ⑨ Ubuntu 腿出现 winners>1 或 unix orphan 用例异常，即可直接推翻本次 PASS 并回退到 BLOCKED。

FINAL REVIEW: PASS