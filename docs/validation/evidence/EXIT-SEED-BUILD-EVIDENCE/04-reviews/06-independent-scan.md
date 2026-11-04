# 06-independent-scan — ⑥ 独立扫描（原文归档）

归档说明：本文件为该轮审查方最终消息的**逐字原文**。审查方 = `independent-reviewer`；调用模型 = `MAI-Code-1.1-Flash (copilot)`；时点 2026-11-05。块末结论行为 `INDEPENDENT SCAN: PASS`。

```text
## Section A

时间点：2026-11-05 本轮只读审计期间。

扫描范围：按 `git status --porcelain -uall` 的实际状态，53 项已暂存变更均落在两份台账文档与一组证据包；我实际读取了：

- [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md)
- [docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md)
- [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md)
- [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/SHA256SUMS.txt](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/SHA256SUMS.txt)
- [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md)
- [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/anomalies.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/anomalies.md)
- [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report.md)
- [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report-r2.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report-r2.md)
- [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/index.json](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/index.json)
- [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/index.json](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/index.json)
- [cmd/prfrail-release/main.go](cmd/prfrail-release/main.go#L13-L82)
- [internal/release/verify.go](internal/release/verify.go#L82-L129)
- [internal/release/probe.go](internal/release/probe.go#L24-L77)
- [internal/release/verify_test.go](internal/release/verify_test.go#L15-L57)

结论：在这次只读审计中，变更集未出现 `.go`、`.github/workflows/**` 或其他运行时代码越界；修订内容被限制在文档与证据包范围内。 

## Section B

(no findings)

## Section C

- 变更集证据：
  - `git status --porcelain -uall` 返回 53 条暂存项，且 `git diff --cached --name-only` 仅命中两份文档和证据包文件；没有 `.go` / workflow / 其它越界文件。
  - 这与用户说明一致：当前 HEAD 为 `baa144d`，且未 commit / 未 push。

- 载荷与哈希三方相等，机械读数为：
  - `disk_count=49`
  - `sum_count=49`
  - `manifest_count=49`
  - `disk_set_equals_sum=True`
  - `manifest_equals_disk=True`
  - `hash_ok=49`
  - `hash_bad=0`

- 对照证据：
  - [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md) 明确写出“条目集三方相等（MANIFEST 归档条目集合 ≡ SHA256SUMS 条目集合 ≡ 包内载荷文件集合）”。
  - [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/SHA256SUMS.txt](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/SHA256SUMS.txt) 的 49 行条目与磁盘载荷一致，且我用文件级 SHA256 实测全部匹配。
  - [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md) 的 N-1 至 N-7 都被明确标注为“不 assert”，与代码层限制一致，不存在把“未证明结论”写成“已证明事实”。

- 代码级 fail-close 证据：
  - [internal/release/verify.go](internal/release/verify.go#L82-L129) 里，`VerifySelfHost` 在 `pathsOverlap` 检查后才做 `HashTree` 和 `seedHash != plan.ExpectedSeedHash` 比较，且 `existingRealPath` 只做 `Abs` + `EvalSymlinks`，没有写入逻辑。
  - [internal/release/verify.go](internal/release/verify.go#L214-L235) 的 `equalOracle` 只比较结构化 `schemaVersion` 和逐项 case 对象，不把「oracle 一致」伪装成「无恶意证明」。
  - [internal/release/verify_test.go](internal/release/verify_test.go#L15-L57) 明确测试了 overlap、oracle mismatch、seed mutation 的 fail-closed 行为；这正是“可证伪的真实钉住条目”。

## Section D

- 双侧 fail-close：已被真实钉住。
  - 钉住条件：`pathsOverlap`、`!pathWithin(candidateRoot, actualOracle)`、`seedHash != plan.ExpectedSeedHash`、`!equalOracle(...)` 都直接返回 `ErrInvalidSelfHost`。
  - 证据： [internal/release/verify.go](internal/release/verify.go#L82-L129) 与 [internal/release/verify_test.go](internal/release/verify_test.go#L15-L57)。
  - 这不是哨兵常量诱导的虚假逻辑；它是运行时比较和返回错误的真实判定。

- 失败不损 seed：已被真实钉住。
  - 钉住条件：在失败分支中，`seedHash` 先拿到、再和 `ExpectedSeedHash` 对比，且 `02a-11` 与 `02b-19` 作为“失败后 seed 不变”的对照证据被收录。
  - 证据： [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md) 的 S-6；对应捕获件 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-11-seed-intact.txt](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-11-seed-intact.txt) 与 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-19-seed-intact.txt](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-19-seed-intact.txt)。
  - 这也不是“种子值被常量写死”的哨兵用法；它是回读当前树哈希并比较实际值。

- oracle 一致：已被真实钉住，但范围严格受限。
  - 钉住条件：`equalOracle` 比较 `schemaVersion` 和 `Cases` 串行元素；它证明“结构化行为相同”，不是“无恶意证明”。
  - 证据： [internal/release/verify.go](internal/release/verify.go#L214-L235) 与 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md) 的 N-3、N-7。
  - 这条结论被文档清楚限定为“行为漂移比较”，不扩展为安全性结论。

- 工具链同形：审计上成立但范围有限，不是全局说法。
  - 真正钉住的对象是 Linux 侧构建期：`GOPROXY=off` 下完成，且 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md) 明确写出 N-6。
  - 这意味着同形结论只在 Linux 侧被证明，不是 Windows 与 CI 的无条件同形声明。

- 载荷哈希三方相等：已机械钉住。
  - 证据：`disk_count=49`、`sum_count=49`、`manifest_count=49`、`hash_ok=49`、`hash_bad=0`。
  - “常量恒等”在这里不是哨兵；它是对真实磁盘内容摘要的计算结果，且对应的是文件内容，不是标志位。

## Section E

- 无法验证部分：
  - 远端 VM 的现场环境在证据包完成后已清理，当前只读会话无法再次复现那台机器的活状态；因此只能验证“归档证据包 + 代码逻辑 + 哈希稳定性”，不能证明远端环境的实时再现性。
  - [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md) 已显式拒绝将 “oracle 一致” 扩展成 “无恶意”。
  - 这意味着本片只证明“行为一致 + fail-close + seed 未变”，不证明更高层的安全/恶意性结论。

- 推翻条件：
  1. 若 `git status` 里出现任何 `.go`、workflow 或其它非文档/证据包文件；
  2. 若 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/SHA256SUMS.txt](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/SHA256SUMS.txt) 与实际磁盘哈希不再一致；
  3. 若 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md) 与磁盘载荷集不一致；
  4. 若 [internal/release/verify.go](internal/release/verify.go#L82-L129) 的判定顺序或写入语义改变；
  5. 若失败路径中的 seed 对照再出现被改写的现象（例如 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-11-seed-intact.txt](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-11-seed-intact.txt) 或 [docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-19-seed-intact.txt](docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-19-seed-intact.txt) 的“相等性”被否定）。

INDEPENDENT SCAN: PASS
```
