# ③ 测试报告 ③b 原文（P85 · T3 反例与变异 + 夹具修复）

### 1. 本机验证输出
`ok github.com/larsonzh/prfrail/internal/adapters 43.587s`；`ok .../internal/console 22.352s`；`GOOS=linux go vet ./internal/adapters/` 无输出。

### 2. 夹具修复（保持原断言语义）
注入点判定：console 侧无注入点（`ai.go:397-409` 硬接线）⇒ 正向夹具指向本机事实成立的系统二进制（helper `p85T3ImmutableProbeExecutable(t)`；`VerifyAIInstallSource` 失败即 `t.Skip`，含 euid=0 的 Linux 腿）。
- `console/ai_test.go`：2 例正向夹具改用 helper；swap 子例改 `t.Skip`
- `console/p83_fixround4_tester_test.go` / `p83_fixround6_tester_test.go`：同改用 helper
- `adapters/p83_fixround4_f8_tester_test.go`：`p83Fix4CanaryCommand` 改绝对路径并先自证，不可证则 skip
- `adapters/ai_preflight_test.go`：`fakeAIChannelProbe` 补 `VerifyBYOKConfinement() error { return nil }`
- 多修 1 例：`TestAIAvailabilityCheckBlocksMissingSecretWithoutProbe`（T8a 门禁使缺能力 fake 被判 `byok_confinement_unproven` 抢在 secret_missing 前；补能力后恢复原语义）

### 3. 新增测试
- `console/p85_t3_tester_test.go`：`TestP85T3RefusesUnprovableInstallSourceWithoutPublishingCarrier`（反例①）、`TestP85T3RefusesMalformedCarrierInsteadOfRecomputing`（反例②）、`TestP85T3RefusesCarrierThatNoLongerDescribesTheSource`（反例③）、`TestP85T3PublishesTheCarrierAndConverges`（正向对照）
- `adapters/p85_t3_tester_test.go`：`TestP85T3WritableInstallDirectoryIsNotProven`、`TestP85T3ExecuteRefusesUnprovableInstallSourceWithoutSpawning`（反例④ pre-spawn）、`TestP85T3CarrierDecodeRefusesAWeakenedOrMixedClaim`、`TestP85T3MatchesDetectsAChangedObservation`
- `adapters/p85_t3_linux_tester_test.go`（`//go:build !windows`）：`...RefusesWritableInstallDirectory`、`...RefusesADirectoryThisUserOwns`、`...ShebangInterpreterInWritableDirectoryRefused`、`TestP85T3LinuxSameInodeRewriteRequiresAPrivilegedFixture`（本版显式 skip，⑤ r1 后改豁免登记）
所有拒绝断言含「不发布 / 不推进 / 证据仍完整」三元。

### 4. 变异检验
基线：`ai_install_proof_windows.go = 6A80612F…92FE1`；`ai_probe_copilot_executor.go = 60305D44…9F260`。
- M1 撤 Windows ACL 判定 ⇒ 红（三测试同红）；恢复 sha256 复位 ✅
- M4 撤 pre-spawn 实时重验 ⇒ 红（`outcome="exited"` 子进程真被启动）；恢复 ✅
- M2 撤 Linux 模式检查 / M3 撤解释器验证：本机不可执行（Windows），执行证据待 Ubuntu VM（10.0.0.199）
全程未用 git checkout。

### 5. 跳过/不可测登记
1. `TestAICheckRefusesReceiverMismatch…/executor-refuses-a-swapped-executable-before-the-spawn` → `t.Skip`（console 无注入点）；等价覆盖 = `TestP83T3cExecutableReplacedAfterApprovalRefused`（adapters）
2. Linux 反向 + 解释器用例 Windows skip；euid=0 正向 skip
3. 同 inode 改写标「不可测」（本版；⑤ r1 后改豁免登记）
4. Windows CI 腿风险：`windows-latest` 提升权限 ⇒ `System32` 对 Administrators 含 (M) ⇒ 正向夹具会 skip；Ubuntu 腿（非 root）与本地非提升 Windows 覆盖完整
