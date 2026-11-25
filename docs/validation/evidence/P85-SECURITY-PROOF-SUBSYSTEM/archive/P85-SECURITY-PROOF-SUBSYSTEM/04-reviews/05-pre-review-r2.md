# ⑤ 预审 r2 原文（P85 · DeepSeek V4 Pro · 原位重跑）

**PRE-REVIEW: PASS**

F1 | CLOSED | `p85_t3_linux_tester_test.go:166` 磁盘实况 = `EXEMPT-REGISTERED P85/F1` 四要素齐全（reason=root 属主夹具需 chown 且 VM 无免密 sudo；carried-by 两测试；destination=报告 §8 + ⑦ 输入包），文件头注释同步声明为登记而非静默 gap。
F2 | CLOSED | `CONTRACTS.md:216`「读取方 = 检查命令（每次启动前重读）与探针执行器（每次启动前实时重验路径→文件身份→摘要）」；`CONTRACTS_EN.md:210` 同位镜像（re-read / re-verifies path, file identity, and digest live before every start），行号精确、语义一致。
F3 | CLOSED | `ai_install_proof_unix.go` 头读 4 字节；非 `#!` 须 `\x7fELF` 魔数，否则 fail-closed；<4 字节短文件同样拒绝（不再返回 nil,nil）。新反例 `TestP85T3LinuxNonELFWithoutShebangIsNotProven`（:111，plain-text + short-file 双候选、双断言 err+观察层）与正例 `TestP85T3LinuxELFObjectIsThePositiveControl`（test binary 自身作 ELF）俱在。调用点全量 3 处（生产 1 + 测试 2），合法路径（ELF 或合法 shebang）均保留，无误伤。
F4 | CLOSED | `ai_test.go:539` `EXEMPT-REGISTERED P85/F4` 四要素齐全；承接测试 `TestP83T3cExecutableReplacedAfterApprovalRefused` 实证存在（`p83_fixround5_tester_test.go:40`）。
F5 | CLOSED | `ai_probe_copilot_executor.go:467` 偏离注释在 `assertCopilotCLIProbeBYOKConfinement` 前：明确「有意不落地域常量、复用 T8b 闭集机制、已申报非遗漏」。
F6 | CLOSED(流程承诺) | 磁盘未见三份报告原文/变异红留痕归档，与主控「⑦ 前归档」声明一致；判可接受但 ⑦ 前须核到。
F7 | N/A | ⑤ r1 已认可无需处置，维持。

新发现（1 条）：
- **Low** | `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/diff.patch` 滞后于磁盘：diff 内 `installSourceInterpreter` 仍为 2 字节头 + 非 shebang 即 `(nil,nil)`（pre-F3），`p85_t3_linux_tester_test.go` 仍为 110 行裸 skip 版（无 F3 反例、无四要素）。主控已声明重建——**⑦ 前必须重建并逐文件核对 diff 与磁盘一致**，否则终审输入包失真（F3 修复在 diff 中不可见）；若未重建即交 ⑦，此项升级 Medium。

结论：r2 收窄范围内七条处置全部真实落地、无新增缺陷；唯一遗留为 diff.patch 重建（已知、有承诺）。
