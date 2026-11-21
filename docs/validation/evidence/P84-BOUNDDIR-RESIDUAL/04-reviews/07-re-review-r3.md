# P84 ⑦ 复审 r3 逐字转录 · 2026-11-21

- 载体：`independent-reviewer`（只读）；模型 `GPT-6.2 Sol (copilot)`，档位 Max（用户 2026-11-21 排片裁定）。
- 额度：⑦ = 1 终审 + 3 复审（本轮 = 复审 3/3，额度用尽）；盲审另计 +1 已在前。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。

---

只读核对完成。N11 和 F4/N5 的入口验收已补上，安装文档也与哨兵一致；T4 仍有一个可静态确认的规范卷根误拒，我会把它列为阻断项。

RE-REVIEW: FINDINGS

P84 ⑦ 原位复审 r3（3/3）：**T5#4 闭合；T4 未闭合。** R2 的路径折叠和发布入口验收缺口均已补上。本轮新增 1 项 Medium：规范盘符根通过规范形检查后，仍会被等值分支误判为界外。

**逐项复核**
| 项 | 严重度 | 状态 | 独立核对 |
|---|---|---|---|
| F1 | High | CLOSED（原反例） | [走链:676](internal/evidence/bounddir_windows.go#L676) 从绑定句柄重走后缀，不让内核跟随中间重定向。[W5b:318](internal/evidence/p84_bounddir_residual_test.go#L318) 区分界内两跳和第二跳越界。 |
| F2 | High | CLOSED | [历史校验:2135](internal/console/runtime_agent.go#L2135) 先验真实前驱，异进程再走锁内接管。[C8:709](internal/console/p84_t5h4_residual_test.go#L709) 核对接管后的代次、持有者和摘要。 |
| F3 | Medium | CLOSED | [资源锁释放:285](internal/guard/lease.go#L285) 不删除锁文件。[C9:781](internal/console/p84_t5h4_residual_test.go#L781) 区分进程退出后锁可重获与已提交记录导致的 `ErrLeaseHeld`。 |
| F4 | Medium | CLOSED | [Publish:1982](internal/console/runtime_supervised.go#L1982) 与 [claim/adopt:1173](internal/console/runtime_supervised.go#L1173) 均使用锁内写入；C18 补上实际发布入口验收。 |
| F5 | Medium | CLOSED | [根枚举:169](internal/evidence/bounddir_windows.go#L169) 直接走绑定句柄，不按名重开；[W6:387](internal/evidence/p84_bounddir_residual_test.go#L387) 保留替身对照。 |
| F6 | Low | CLOSED | [UTF-16 长度守卫:991](internal/evidence/bounddir_windows.go#L991) 为终止 NUL 留两字节。[W7:453](internal/evidence/p84_bounddir_residual_test.go#L453) 能区分删掉预留的回退。 |
| N1 | High | CLOSED | [Windows:81](internal/guard/lease_lock_windows.go#L81) 与 [Unix:50](internal/guard/lease_lock_unix.go#L50) 对稳定文件对象取排他锁；[C11:904](internal/console/p84_t5h4_residual_test.go#L904) 能区分活锁和未取 OS 锁。 |
| N2 | Medium | CLOSED | [打开锁文件:60](internal/guard/lease_lock_windows.go#L60) 不等于持锁；取得锁后由 OS 回收，不再有占锁后等待写 PID 标记的状态。 |
| N3 | Medium | CLOSED | [目标选择:802](internal/evidence/bounddir_windows.go#L802) 只用 `SubstituteName`；[C12:616](internal/evidence/p84_bounddir_residual_test.go#L616) 的双向对照能抓字段互换。名称解码另见 N14。 |
| N4 | Medium | CLOSED | [恒等也批:2001](internal/console/runtime_supervised.go#L2001)，批内跳过已当前项。[C13:968](internal/console/p84_t5h4_residual_test.go#L968) 能抓恒等早退。不承诺整批回滚。 |
| N5 | Medium | CLOSED | [单写:2364](internal/console/runtime_agent.go#L2364)、[批写:1230](internal/console/runtime_supervised.go#L1230) 都在锁内读当前记录、比 token/hash 并落盘，与代次推进共用 [`withLock`:262](internal/guard/lease.go#L262)。[C14:1009](internal/console/p84_t5h4_residual_test.go#L1009) 钉批写，[C16:1142](internal/console/p84_t5h4_residual_test.go#L1142) 钉单写，C18 钉入口。 |
| N6 | Medium | CLOSED | [续租入口:2096](internal/console/runtime_agent.go#L2096)、[接管回调:2166](internal/console/runtime_agent.go#L2166) 与[发布终验:2401](internal/console/runtime_agent.go#L2401) 均查资源和 consumer；[C15:1086](internal/console/p84_t5h4_residual_test.go#L1086) 用真实前驱验证异 consumer 仍被拒。 |
| N7 | Medium | CLOSED（原反例） | [前缀归一化:868](internal/evidence/bounddir_windows.go#L868) 不裁空白；[W8:677](internal/evidence/p84_bounddir_residual_test.go#L677) 检查字节保留、近名兄弟拒绝和界内对照。 |
| N8 | Medium | CLOSED | [哨兵:323](internal/guard/lease.go#L323) 与 [C17:1366](internal/console/p84_t5h4_residual_test.go#L1366) 对应：旧目录不等待、不删除、不写记录，移除后同一 Claim 可成功。 |
| N9 | Low | CLOSED（原范围） | 归档已有 [Rev 3b:125](docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/01-arch/01-arch-design.md#L125) 和 [Rev 4:134](docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/01-arch/01-arch-design.md#L134)；Rev 5 归档缺口并入 N12。 |
| N10 | Medium | CLOSED（原折叠反例） | [规范形守卫:736](internal/evidence/bounddir_windows.go#L736) 拒绝会被 `Clean` 改写的目标，[比较:952](internal/evidence/bounddir_windows.go#L952) 不再拿改写后的目标判归属。普通正斜杠、尾分隔符等即使能被 Windows 解析，仍是 Rev 5 有意拒绝；规范卷根误拒另见 N13。 |
| N11 | Medium | CLOSED | [C18:1239](internal/console/p84_t5h4_residual_test.go#L1239) 复用发布夹具，只推进记录一代；同进程真实 N−1 proof [追平:2142](internal/console/runtime_agent.go#L2142) 后，在[阈值判断:2180](internal/console/runtime_agent.go#L2180) 直接返回，无 Renew/Reclaim。`refreshed != proof` 因此进入实际 `Publish` 的本任务写分支。 |
| N12 | Low | OPEN（非阻断） | [MANIFEST:15](docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/MANIFEST.md#L15) 仍只列 Rev 1/2；[归档方案:141](docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/01-arch/01-arch-design.md#L141) 止于 Rev 4，[Rev 5:143](tmp/P84-BOUNDDIR-RESIDUAL/01-arch-design.md#L143) 只在临时副本。收尾时须同步正文和清单，再重算摘要；G8-b 自洽不证明设计版本齐全。 |
| 安装联动 | Medium（N8 验收） | CLOSED | [中文:202](docs/INSTALLATION.md#L202) 与 [英文:202](docs/INSTALLATION_EN.md#L202) 同位、同义，与[哨兵及诊断:329](internal/guard/lease.go#L329) 一致：识别 `legacy lock directory`，先证明旧 run 停机，再人工移除遗留目录。 |

**独立发现**
| 项 | 严重度 | 发现 | 最小修复 |
|---|---|---|---|
| N13 | Medium（阻断） | [bounddir_windows.go:953](internal/evidence/bounddir_windows.go#L953)：当基准和目标都是规范的 `D:\` 时，`Clean` 保留根分隔符，但 `TrimRight` 把基准变成 `D:`。等值不成立，严格后代分支又因目标刚好等于前缀而拒绝。这是确定的界内误拒，不是 Rev 5 有意拒绝的非规范拼写。[W9 正例:832](internal/evidence/p84_bounddir_residual_test.go#L832) 没覆盖盘符根。 | 等值比较保留规范根名的必要分隔符，仅在构造严格后代前缀时处理分隔符。加盘符根等值与严格后代正例，不放宽现有拒绝条件。 |
| N14 | Low（原生可达性未证） | [bounddir_windows.go:844](internal/evidence/bounddir_windows.go#L844)：按声明长度提取的名称交给 `UTF16ToString`，若声明范围内含 NUL，会先截断，再进入规范形检查。解析器因此没有独立保证按完整名称判定。未证实受支持的文件系统会返回这种缓冲，也未发现界外读取。 | 解码前对 `SubstituteName` 声明范围内的 NUL 明确拒绝；用解析器测试区分拒绝与截断，不依赖原生文件系统生成该输入。 |

**可证伪性审计**
- [W9:794](internal/evidence/p84_bounddir_residual_test.go#L794) 的五个负例能区分 M1 恢复 `Clean` 后再比较；M2 只删生产侧理由检查时，后缀守卫仍拒绝，变红的是 `cannot be interpreted losslessly` 断言。真实 reparse 数据用例直调生产后缀入口，不等于公开 API 全链路；[P52:84](internal/adapters/bounddir_escape_windows_test.go#L84) 和 W8 补公开入口对照。
- [C18:1315](internal/console/p84_t5h4_residual_test.go#L1315) 持锁时要求 `DeadlineExceeded`、无发布决策、proof 字节不变；释放后要求 `completed`、一条租约证据及当前 proof。M3 即使仍被后续批写锁挡住，也会先改变字节，因此不能靠超时断言蒙混过关。C18 不单独证明所有并发交错；锁内等值与落盘由源码核对。
- [C9:781](internal/console/p84_t5h4_residual_test.go#L781) 单测不能排除从未取锁，须与 C11 合看；[W4:40](internal/gates/p83_fixround5_tester_test.go#L40) 的成功断言也看不到一次失败后被忽略的根打开。

**无法验证**
已读 [R10 树门禁:45](tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-postR10.txt#L45) 的 `TOTAL_FAIL=0`、[安装镜像门禁:10](tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-postR10.txt#L10) 的五形状通过和 [R10 全量日志:3](tmp/P84-BOUNDDIR-RESIDUAL/test-all-postR10.txt#L3) 的 15 个包级 `ok`。未独立核验逐例 PASS/SKIP、M1–M3 实跑及恢复、`SELFTEST 740/740`、最新完整 diff、摘要复算、原生竞态或旧 run 实际升级。末轮额度已用尽，按 [§7.8.3](docs/DELIVERY_DIRECTIVE.md#L824) 将未闭合项交主控登记和用户裁决，不自行加轮。

未运行测试 / 未修改文件
