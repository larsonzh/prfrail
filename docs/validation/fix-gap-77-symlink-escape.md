# FIX-GAP-77-SYMLINK-ESCAPE（P53）修复验证报告

切片：`FIX-GAP-77-SYMLINK-ESCAPE`（底层安全修复类 `[SLICE]`，M 档）。起于 2026-11-03，准则 v1.60，锚点 `HEAD` = `dc785ed`。主控模型 = `DeepSeek V4 Pro`。

## 1. 切片概览

- 缺陷：`internal/adapters/bounddir_unix.go` 以单次 `openat(O_NOFOLLOW)` 约束路径，只覆盖末组件；中间符号链接可把解析带出绑定目录。P52 已在 Linux 上实证读出界外字节。
- 修复：改为 `openat2(2)` 单次受限解析，受限掩码 `RESOLVE_NO_XDEV|RESOLVE_NO_MAGICLINKS|RESOLVE_NO_SYMLINKS|RESOLVE_BENEATH` = `0x0F`，`SYS_openat2` = 437，24 字节 `open_how` 私有常量，零新依赖。
- 降级语义：非 linux/amd64 或内核返回 `ENOSYS` 时 fail-closed，返回 `PathError{Err: ENOSYS}`，不保留逐级 walk 回退；`allowUnsafeFallback` 本片不启用。
- 生命周期：私有 `mu sync.Mutex` 闭合 fd 的取用与关闭；公共签名冻结，Windows 与 darwin 实现零改动。
- 改动集：生产侧只有 `bounddir_unix.go`（`git diff --numstat` 实测 `+158/-12`）；测试侧为常驻回归改造，另新增 payload、redbaseline、failclosed、storm 各一支用例文件。

## 2. 目标与裁定

裁定由用户在 2026-11-04 下达，原文归档于 `01-arch/rulings.md`。

| 条目 | 裁定 | 落地 |
|---|---|---|
| (a) | 威胁边界包含外部 rename，`openat2` 为硬要求；降级环境 fail-closed | 主路径只在 linux/amd64 启用，`ENOSYS` 直接拒绝 |
| (b) | 挂载导入属攻击面，保留 `NO_XDEV` | 接受「界内合法子挂载被拒」的代价 |
| (c) | 不引入深度上限 | 深度上限设想删除 |
| (d) | 文档同步纳入本片 ⑧ | 本报告与 `_EN` 镜像 |
| (e) | 签认 `openat2` 超出二选一清单 | 记入 `freedom-list.md` |

## 3. 派单表

| 环节 | agentName | model | 结论 |
|---|---|---|---|
| ① r1 | `deep-reasoner` | DeepSeek V4 Pro | A/B 混合方案 |
| ① 红队 | `independent-reviewer` | GPT-6.1 Sol | `FINDINGS`（2 High + 4 Medium） |
| ① r2 | `deep-reasoner` | DeepSeek V4 Pro | `openat2` 定稿 |
| ③ | `prfrail-tester` | DeepSeek V4.1 Flash | 真实红与 F1 护栏；F2 登记 SKIP |
| ②a | `prfrail-implementer` | DeepSeek V4.1 Flash | 生产实现落地 |
| ③b | `prfrail-tester` | DeepSeek V4.1 Flash | 消除常驻恒红 |
| ③c | `prfrail-tester` | DeepSeek V4.1 Flash | 告警消除与 ⑤ 发现处置 |
| ⑤ r1 | `deep-reasoner` | DeepSeek V4 Pro | 2 Medium 与 3 Low |
| ⑤ r2 | `deep-reasoner` | DeepSeek V4 Pro | `PRE-REVIEW: PASS` |
| ⑥ | `independent-reviewer` | MAI-Code-1.1-Flash | 独立扫描，原文见 5.1 节 |
| ⑦ | `independent-reviewer` | GPT-5.3-Codex | 终审无发现 |
| 批辅 ×2 | `deep-reasoner` | GLM-5.3-Flash | `openat` 用法枚举与 Linux 兼容性事实 |

## 4. 红绿证据

读数取自 Ubuntu VM（内核 6.8，go1.22.12，经 `bash -lc` 执行）；本机 Windows 只承担编译与静态检查。

- 红（改生产代码之前跑出）：`TestBoundDirMiddleSymlinkEscapesBound` 失败，文案逐字为 `GAP-77 regression: middle symlink escape succeeded and read outside the bound directory`，退出码 1；越界句柄逐字节读出界外 23 字节 `SECRET-OUTSIDE-BOUNDARY`。
- 绿（修复后同一用例）：通过；拒绝 errno 落在闭集 `{ELOOP, ENOTDIR, EXDEV}` 内，实测为 `ELOOP`。
- 正向对照：界内多层路径 `sub/nested/inside` 读取成功，且与 canary `IN-SIDE-BOUNDARY` 全串逐字节相等。
- 变异对照：`PRFRAIL_P53_RED_BASELINE=1 go test -run TestBoundDirRedBaselineMiddleSymlinkEscape` 在固定树上失败，文案 `red baseline not reproduced`；默认套件下该用例按环境变量门控 SKIP。
- 并发与风暴：`TestBoundDirConcurrentCloseAndRead` 以 4 个读者循环读取、主协程执行 16 次 `Close`，结果只允许「界内精确字节」或 `adapters: nil bound directory`；`TestBoundDirRenameStormNeverYieldsOutsideBytes` 跑 200 轮，界内成功读 6336 次、拒绝 1495 次，无界外字节；`-race` 对 `adapters` 与 `chain` 无 DATA RACE。
- fail-closed 契约：`TestBoundDirFailClosedWithoutOpenat2` 锁存 `openat2Disabled` 后断言 `*os.PathError`、`Op` 为 `openat`、`errors.Is(..., syscall.ENOSYS)`、无 panic，二次读取同形拒绝；非 linux/amd64 分支由常量平台门直接拒绝。

## 5. 审查链

- ① 架构：r1 出 A/B 混合方案，红队（一次性临时授权）报 2 High 与 4 Medium，r2 改为 `openat2` 单次受限解析定稿；结论 `DESIGN: DONE`。
- ③ 测试先行：改生产代码之前先跑出真实红，红证据归档于 `03-tests/red-baseline.md`。
- ②a 实施：生产单文件落地，公共签名与平台分支不变。
- ③b 与 ③c：先消除常驻恒红，再按 ⑤ 的发现处置测试口径；记录在 `03-tests/remediation-report.md`。
- ⑤ 预审：r1 报 2 Medium 与 3 Low（`PRE-REVIEW: FINDINGS`），③c 处置后 r2 复核为 `PRE-REVIEW: PASS`。
- ⑥ 独立扫描：结论见 5.1 节，原文归档于 `06-scan/review.md`。
- ⑦ 终审：无发现，`FINAL REVIEW: PASS`（`07-final/review.md`）。

### 5.1 ⑥ 输出原文

以下为 ⑥ 独立扫描输出的逐字原文（保留原始语言），两侧逐字节相同；围栏末行为结论行，围栏外不作评价。

```sixstage
## Section A

范围与输入：
- 生产实现： [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L1-L253)
- 关键回归测试： [internal/adapters/bounddir_escape_unix_test.go](internal/adapters/bounddir_escape_unix_test.go#L1-L197), [internal/adapters/bounddir_failclosed_unix_test.go](internal/adapters/bounddir_failclosed_unix_test.go#L1-L226), [internal/adapters/bounddir_symlink_storm_test.go](internal/adapters/bounddir_symlink_storm_test.go#L1-L220)
- 共享载荷： [internal/adapters/bounddir_escape_payload_unix_test.go](internal/adapters/bounddir_escape_payload_unix_test.go#L1-L97), [internal/adapters/bounddir_escape_redbaseline_test.go](internal/adapters/bounddir_escape_redbaseline_test.go#L1-L95)
- 权威证据： [docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/design-r2.md](docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/design-r2.md), [docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/rulings.md](docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/rulings.md), [docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/green-suite-report.md](docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/green-suite-report.md)

本轮只读，不做任何文件修改。已检查对象范围内的 git 状态与 diff；未发现 Windows/darwin 代码、go.mod、其他生产文件被触碰，且目标文件集中在本次审查集。

## Section B

(no findings)

本次独立扫描未发现需要修正的真实缺陷：openat2 ABI、fail-closed 门、fd 生命周期锁、回归测试闭集和可证伪性均与设计和实现相一致；当前状态不构成越界改动。

## Section C

证据摘要：
- ABI 与 resolve mask： [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L18-L61) 明确写入 SYS_openat2 437、open_how 24B、RESOLVE_NO_XDEV/NO_MAGICLINKS/NO_SYMLINKS/BENEATH。与 linux/openat2.h 语义一致。
- Fail-closed 语义： [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L95-L176) 仅在 linux/amd64 上通往 openat2；其余平台和 ENOSYS 直接返回，且 openat2Disabled 只在 ENOSYS 时置位，不扩散到其他错误。
- fd 生命周期与 dir.mu： [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L73-L110) 与 [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L181-L253) 显示 BoundDir 的 fd/Close/verify 都在同一私有 mutex 保护下，未出现遗漏的读写点，且 Close 幂等。
- 回归测试闭集与正向控制： [internal/adapters/bounddir_escape_unix_test.go](internal/adapters/bounddir_escape_unix_test.go#L34-L197) 对 middle symlink 采用 {ELOOP, ENOTDIR, EXDEV} 的拒绝闭集，并用 in-bound canary 做正向验证；[internal/adapters/bounddir_failclosed_unix_test.go](internal/adapters/bounddir_failclosed_unix_test.go#L24-L226) 额外锁定 ENOSYS 契约和二次拒绝稳定性。
- 风暴与 Rename/F2 约束： [internal/adapters/bounddir_symlink_storm_test.go](internal/adapters/bounddir_symlink_storm_test.go#L78-L220) 证明 rename storm 不会泄出 outside bytes，且 F2/环境门控逻辑是有条件地跳过，不是静默放行。
- 真实验证：已执行 go test -count=1 ./internal/adapters，结果为：
  ok      github.com/larsonzh/prfrail/internal/adapters   33.473s

## Section D

可证伪性审计：
- 若去掉 openat2 runtime gate 或把 openat2Disabled 语义去掉，fail-closed 契约测试会立即失效；这已被 [internal/adapters/bounddir_failclosed_unix_test.go](internal/adapters/bounddir_failclosed_unix_test.go#L24-L226) 钉住。
- 若去掉 dir.mu 或关闭的 fd 失效，concurrent close/read 测试会在 [internal/adapters/bounddir_failclosed_unix_test.go](internal/adapters/bounddir_failclosed_unix_test.go#L100-L226) 中暴露 EBADF、泄露字节或 panic。
- 若误把拒绝闭集放宽到任意错误， [internal/adapters/bounddir_escape_unix_test.go](internal/adapters/bounddir_escape_unix_test.go#L86-L126) 立即判红；这正是该类测试的可证伪性。
- git 证据与范围核对也没有发现越界修改；本轮不存在 go.mod、Windows/darwin 或其他生产文件的覆盖改动。

## Section E

无法验证部分：
- F2 bind-mount 的真实 namespace 环境需要特权/namespace 支持，因环境限制会 skip；这不是代码缺陷，而是宿主能力限制。
- 这次扫描只审查了目标改动集，不扩展到其他 repo 组件；如果后续出现新的宿主环境差异，需重新做平台分支验证。

推翻条件：
- 只要出现下述任一事实，即可推翻当前结论：
  - 目标实现未真正使用单次 openat2 受限解析；
  - 非 linux/amd64 或 ENOSYS 分支被绕过并泄漏到正常路径；
  - dir.mu 未覆盖所有 fd 生命周期/close 读写点；
  - git status/diff 显示出了 go.mod、Windows/darwin 或其他生产文件的无授权改动。

INDEPENDENT SCAN: PASS
```

## 6. 门禁读数

- 本机（Windows）：`gofmt -l` 无输出；`go build ./...` 与 `go vet ./...` 退出码 0；`go test ./...` 全仓十五个包 ok。
- VM（go1.22.12）：全仓十五个包 ok；`-race` 对 `adapters` 与 `chain` 无 DATA RACE。
- 门禁：④ 与 ⑨ 阶段（本报告对落盘之前的代码态）读数为 `TOTAL_FAIL=0`，`SELFTEST 697/697` 与基线一致（生产改动未触及 gates 判据实现）。
- 本报告对落盘后重跑 tree 作用域，读数为 `TOTAL_FAIL=2`：`G7-b` 仅因围栏条目尚未入索引而红（⑧b 暂存后转绿），`G8-b` 仅因 `MANIFEST.md` 尚未收入 `06-scan/review.md` 而红（主控在本报告落盘后重算 MANIFEST 与 SHA256SUMS）。
- 本报告对自身的判据：`G3(a)` 六形状（行数、标题、加粗、竖线、空行、反引号）两侧相等；`G4a` 编码为 UTF-8 带 BOM、行尾 LF；`G4b` 未知汉字为 0；`G7-e` 与 `G7-f` 均 PASS。

## 7. 偏离与未闭环

- 未闭环 1（F2）：bind mount 反例在 VM 因 `unshare(CLONE_NEWUSER)` 返回 `EINVAL`、`unshare(CLONE_NEWNS)` 返回 `EPERM`、且无免密 sudo 而 SKIP；已登记为环境限制，不计为通过。
- 未闭环 2（内核下限）：真实低于 5.6 的内核主机不可得，改以 `openat2Disabled` 锁存的契约测试替代，判据与生产常量门一一对应。
- 未闭环 3（变异检验）：`PRFRAIL_P53_RED_BASELINE=1` 对照命令留 ⑨ 复核执行；本报告只登记命令与判读口径。
- D1（构造性冲突）：②a 报告的 4.1 节要求三测全 PASS，与 ③ 的交付物定义相冲突，由 ③b 重构处置。
- D2（历史归档）：③b 期的 `-tags p53_red_baseline` 记载保留在 `03-tests/green-suite-report.md` 与 `MANIFEST.md` 内不改写；替代命令以代码注释与本报告为准。
- D3（锁序口径）：按实况描述——`readBoundEvidence` 与 `CloseBoundDirs` 先释 `store.mu` 再触 `dir.mu`，`bindEvidenceDir` 先取 `dir.mu` 后取 `store.mu`，均为顺序执行、无嵌套持有。
- D4（转录更正）：任务包中「27 字节」系对 ③c 报告「27 byte-exact reads」的误读；界内判定用全串逐字节比较，不依赖字节长度。
- 边界：生产改动数值以本次 `git diff --numstat` 实测为准，早前轮次的快照读数与之不同。
- 已知限制：Windows 与 darwin 实现本片零改动，两平台的反例覆盖沿用既有用例，本报告不重复。

## 8. artifacts 围栏

以下为本片完整变更集条目（证据载荷、报告对、台账、规范文档与代码/测试改动；路径 + 归属，TAB 分隔，按路径字典序）。

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/CHECKLIST_v43.md	repo
docs/t027/CHECKLIST_v43_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/00-briefing.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/design-r2.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/design.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/red-team-sol-review.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/rulings.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/green-suite-report.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/red-baseline.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/remediation-report.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/05-prereview/review-r1.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/05-prereview/review-r2.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/06-scan/review.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/07-final/review.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/08-gates/gate-readings.txt	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/09-glm/glm-1-openat-inventory.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/09-glm/glm-2-linux-facts.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/MANIFEST.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/SHA256SUMS.txt	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/freedom-list.md	repo
docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/input-list.md	repo
docs/validation/fix-gap-77-symlink-escape.md	repo
docs/validation/fix-gap-77-symlink-escape_EN.md	repo
internal/adapters/bounddir_escape_payload_unix_test.go	repo
internal/adapters/bounddir_escape_redbaseline_test.go	repo
internal/adapters/bounddir_escape_unix_test.go	repo
internal/adapters/bounddir_failclosed_unix_test.go	repo
internal/adapters/bounddir_symlink_storm_test.go	repo
internal/adapters/bounddir_unix.go	repo
```

## 9. 提交与 CI

- 主提交：`7802819`（`fix: close the middle-symlink escape on the bound dir via openat2 (P53)`，37 件文件，+2873/−62），已推 origin（未推 gitee）。
- CI run `41707160563`（head `7802819`，attempt 1）：**双腿全绿**（`Go ubuntu-latest` 与 `Go windows-latest` 均 success；3 个 dispatch 专用作业常态 skipped）。
- 双腿 Gates 实测：`SELFTEST: PASS 697/697`、`GATE REPORT scope=ci changed=37`、`TOTAL_FAIL=0`（ubuntu 腿 19:54 / windows 腿 19:56）。
- 包级硬证据：windows 腿 `ok internal/adapters 35.644s`；ubuntu 腿 Test `ok internal/adapters 13.415s` + Race 腿 `ok internal/adapters 91.205s`、`ok internal/chain 2.145s`。

## 10. 结论

- GAP-77 在本片闭合：生产侧由 `openat2` 单次受限解析替代单次 `openat`，非 linux/amd64 与 `ENOSYS` 走 fail-closed；真实红在改代码之前跑出，修复后转为绿，并保留变异对照命令。
- 审查链：① 双打、③ 红绿反转、②a 实施、③b 与 ③c 返工、⑤ 两轮、⑥ 一次、⑦ 终审均已完成，⑦ 无发现。
- 未闭环三项（F2 环境受限、内核下限替代判据、变异检验留 ⑨）已在 7 节登记，不构成对本片结论的削弱。
- 本片已提交 `7802819` 并推 origin（未推 gitee），CI run `41707160563` 双腿全绿（见 9 节）。
