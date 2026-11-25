# ③ 测试报告 ③c 原文（P85 · T8a 反例与变异）

新增文件（唯一改动）：`internal/adapters/p85_t8a_tester_test.go`（611 行，UTF-8 无 BOM + LF）。未改实现、未改测试基础设施、未 commit。

### A. 覆盖矩阵
- P1 `TestP85T8AMissingProofRefusesBeforeResolveOrSpawn`（2 子例，反例·证明缺失）
- P2 `TestP85T8AUnreadableSandboxVerdictIsNotAdmission`（反例·证明不可读 + 正向对照；guard 三判 + checker 两判）
- P3 `TestP85T8AConfinementProofRefusesEachBrokenElement`（9 子例，反例·pre-spawn 五元逐一破坏）
- P4 `TestP85T8AProcessSpecCarriesNoExtraInheritedHandle`（反射钉住无 ExtraFiles/Files/SysProcAttr）
- P5 `TestP85T8AExecuteRefusesBrokenBYOKSpawnWithoutSpawning`（端到端零启动）
- P6 `TestP85T8AExecuteRefusesUnprovableVerdictWithoutSpawning`（S1 判定端到端；主机证明放行则 Skip）
- P7 `TestP85T8ASameClassRefusalsShareEvidenceAcrossValues`（T8b 兼容：不同 model 同类拒绝证据逐字节相等）
- P8 `TestP85T8AHostAuthenticationAndOtherChannelsBypassTheGate`（正向对照·边界）

### B. 验证输出
`go test -count=1 ./internal/adapters/...` → `ok ... 34.417s` 全绿；`GOOS=linux go vet ./internal/adapters/` 无输出。

### C. 变异检验（3 条）
- 撤 checker 门禁（`AuthMode == "secret" && false`）⇒ 红（`Status:available RequestsUsed:1` 凭据已解析、探针已跑；P2 unknown 子例与 P8 同红）
- 撤 executor pre-spawn 断言块 ⇒ 红（`outcome="exited"` 子进程真被创建；P6 同红）
- 撤元素 2 的 override 键循环 ⇒ 红（错误变成 sandbox 类 ⇒ 证明「点名元素而非平台判定」判据有效）
恢复校验：`ai_preflight.go` = `7237694…DDB169`、`ai_probe_copilot_executor.go` = `60305D4…0AD69F260` 与基线逐一相等；备份已删除。

### D. 不可测/边界
1. adapters 无沙箱注入 seam ⇒ 用 `guard.Sandbox` 接口替身证明 gate 语义（P2），P6 用真实平台判定兜底带 Skip 分支
2. `ExtraFiles 非空` 不能经 `ProcessSpec` 构造 ⇒ 结构化反射断言（P4）
3. 元素 2 端到端非逐元素变异敏感（被 S1 沙箱拒绝遮蔽）⇒ 逐元素灵敏度由 P3 白盒用例承担
4. Linux 原生腿未在本机执行；不可证明时 Skip 而非假绿，需 Ubuntu VM（10.0.0.199）复跑
