# ⑦ 追加第四读·Codex 基准读原文（P85 · GPT-5.3-Codex Extra High · 法定模型同输入基线）

已完成盲读核查：三组修复的实现链、R6 变异红绿对照、双平台腿与门禁日志都已复核，且未触碰你禁读的四份他模结论文件。下一步我将给出闭集结论（PASS/FINDINGS）和三组逐项判定、三项终裁与 ⑨ 口径声明。

FINAL REVIEW: PASS

1. Linux ACL 取证对证：CLOSED。
证据：Linux 槽位观测已形成“开名对证→按名 ACL 读→闭名对证”的身份三明治，[internal/adapters/agent_runner_policy_write_guard_unix.go#L152](internal/adapters/agent_runner_policy_write_guard_unix.go#L152)、[internal/adapters/agent_runner_policy_write_guard_unix.go#L155](internal/adapters/agent_runner_policy_write_guard_unix.go#L155)、[internal/adapters/agent_runner_policy_write_guard_unix.go#L170](internal/adapters/agent_runner_policy_write_guard_unix.go#L170)；S1 窗口残余已如实登记于实现注释 [internal/adapters/agent_runner_policy_write_guard_unix.go#L100](internal/adapters/agent_runner_policy_write_guard_unix.go#L100)；终态 ReassertIdentities 在创建/复用/加载三条主路径均执行 [internal/console/agent_runner_policy.go#L486](internal/console/agent_runner_policy.go#L486)、[internal/console/agent_runner_policy.go#L589](internal/console/agent_runner_policy.go#L589)、[internal/console/agent_runner_policy.go#L936](internal/console/agent_runner_policy.go#L936)，且实现覆盖目录/父目录/制品/锚 [internal/console/agent_runner_policy.go#L798](internal/console/agent_runner_policy.go#L798)。

2. 摘要槽位对齐契约：CLOSED。
证据：契约句已固定为四槽并含空标记占位 [docs/CONTRACTS.md#L233](docs/CONTRACTS.md#L233)、[docs/CONTRACTS_EN.md#L227](docs/CONTRACTS_EN.md#L227)；实现域版本为 :2 [internal/adapters/agent_runner_policy_write_guard.go#L61](internal/adapters/agent_runner_policy_write_guard.go#L61)，文件槽按角色折叠而非按路径排序 [internal/adapters/agent_runner_policy_write_guard.go#L223](internal/adapters/agent_runner_policy_write_guard.go#L223)，并以 .policy/.anchor 后缀识别角色 [internal/adapters/agent_runner_policy_write_guard.go#L161](internal/adapters/agent_runner_policy_write_guard.go#L161)、[internal/adapters/agent_runner_policy_write_guard.go#L193](internal/adapters/agent_runner_policy_write_guard.go#L193)；锚-only 场景强制制品槽为 ∅ 的白盒用例在位 [internal/adapters/p85_r6_slot_order_test.go#L101](internal/adapters/p85_r6_slot_order_test.go#L101)；相应变异可打红且恢复后通过 [tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L15](tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L15)、[tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L27](tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L27)。

3. T3 Windows 取证与发布复核走根链：CLOSED。
证据：可执行文件取证入口已改为 OpenBoundDirVerified 同链绑定 [internal/adapters/ai_install_proof.go#L414](internal/adapters/ai_install_proof.go#L414)、[internal/adapters/ai_install_proof.go#L426](internal/adapters/ai_install_proof.go#L426)；其底层为逐组件根链分类并拒绝 reparse [internal/evidence/bounddir_windows.go#L102](internal/evidence/bounddir_windows.go#L102)、[internal/evidence/bounddir_windows.go#L135](internal/evidence/bounddir_windows.go#L135)、[internal/evidence/bounddir_windows.go#L179](internal/evidence/bounddir_windows.go#L179)；发布复核在路径写前后都执行根链重绑定校验 [internal/console/ai_install_proof_publish_other.go#L50](internal/console/ai_install_proof_publish_other.go#L50)、[internal/console/ai_install_proof_publish_other.go#L78](internal/console/ai_install_proof_publish_other.go#L78)、[internal/console/ai_install_proof_publish_other.go#L98](internal/console/ai_install_proof_publish_other.go#L98)、[internal/console/ai_install_proof_publish_other.go#L102](internal/console/ai_install_proof_publish_other.go#L102)；回退为 FileIDAt 的变异可打红且恢复后通过 [tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L52](tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L52)、[tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L61](tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L61)。

新发现：无（0 项）。

三项终裁复核（基于本次三组修复后的状态）
1. T2：CLOSED。
2. T3：CLOSED。
3. T8a：CLOSED（代码复审口径维持；本次三组修复未引入回退迹象，R6 终跑通过 [tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L63](tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L63)、[tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L69](tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt#L69)）。

⑨口径声明
1. F17/F22 口径维持：不单独阻断代码复审。
2. 但证据债仍是证据债，缺证不得关片；⑨ 归档阶段必须按既定口径补齐并保持可追溯。
3. 机械门禁与自检在输入中为全绿（TOTAL_FAIL=0 与 SELFTEST 740/740）[tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/gate-index.txt#L46](tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/gate-index.txt#L46)、[tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/gate-tree.txt#L46](tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/gate-tree.txt#L46)、[tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/gate-ci.txt#L46](tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/gate-ci.txt#L46)、[tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/selftest.txt#L1](tmp/P85-SECURITY-PROOF-SUBSYSTEM/05-input/selftest.txt#L1)。
