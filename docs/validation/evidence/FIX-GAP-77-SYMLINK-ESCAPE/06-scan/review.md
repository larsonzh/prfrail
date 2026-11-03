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
