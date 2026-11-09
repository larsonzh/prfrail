# design-r1 — F2-FIX-STUB-HANG ① 架构（hang 分支恒定写修复设计）

## A. 修复落点（实测行号 + 逐行事实）

全部行号经磁盘实测（`tools/agent-stub/main.go`）：

1. **`case "hang"` 分支 = L112–116**，当前代码逐字：

```go
	case "hang":
		// No artifacts beyond its own pid: the run is stopped from the outside.
		run.writeOwnPID()
		time.Sleep(time.Hour)
		return 0
```

- L114 `run.writeOwnPID()`：已写 `stub.pid`（L254–258：十进制 pid + `\n`，写失败 `os.Exit(11)`）。
- L115 `time.Sleep(time.Hour)` 后 L116 `return 0`：**该分支早退，永远到不了 L135 起的公共路径** ⇒ 整个进程生命周期内 `stub-child.pid` 从不落盘（连 0 字节都没有）——f2/pos 消费方按 A6 判「证据缺口」fail-closed。

2. **公共路径写序 = L135–142**（run/spawn/spawn-exit/spawn-wait/no-usage/ask/mutate-workspace/torn-events 等模式），逐字：

```go
	run.writeOwnPID()
	// A run that cannot spawn writes the empty report as a fact. The spawn modes
	// deliberately write nothing here: an empty report sitting next to a child
	// that is still coming up would be a false "no descendant", and a run stopped
	// mid-spawn must stay a missing report (a gap), never a fabricated fact.
	if !run.spawnsChild() {
		run.writeEmptyChildPIDReport()
	}
```

- 为什么 run 满足恒定写而 hang 不满足：run 经 switch 落入公共路径，`spawnsChild()`（L164–166，仅 spawn/spawn-exit/spawn-wait 为 true）= false ⇒ 必写 0 字节 `stub-child.pid`；hang 在 L116 提前 `return 0`，完全绕过该段。
- 对照语义（L136–139 注释自证）：无子进程写空报告是**正向事实**；spawn 族故意不提前写（防「子进程未起先立空报告」的假事实），真实子 pid 由 `spawnChild()` 启动后覆盖写入（L182 为 Start 失败时补写空报告的 fail-closed 错误路径）。

3. **`writeEmptyChildPIDReport()` = L260–267**：`os.WriteFile(..., childPIDName, nil, 0o600)`——`nil` 数据 ⇒ **0 字节空文件**；失败 `os.Exit(11)`。`writeOwnPID()` = L254–258。

4. **CONTRACTS.md A6（L256，§7 内）恒定写口径原文要点**：

> pid 报告写入口径（恒定写）：在执行到写入点时，无子进程写**空文件（0 字节）**，有子进程写子 pid 列表……若在写入前遭遇外部中断或写入失败，报告可缺失且必须按证据缺口 fail-closed（`uncertain`）处理。消费方三态语义：文件不存在即证据缺口……文件存在且为空（0 字节）即确实无子进程；文件存在且非空即按内容解析子 pid 列表。**不得以「文件不存在」与「文件为空」互推**——前者是缺陷，后者是事实。

5. **漂移事实（`docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt`）**：`CHECK-PASS-DIFF f2-pos/process-tree-taken-away: now=false baseline=true`，now detail = `missing required pid report(s) stub-child.pid so the managed processes cannot be proven gone`（出处 = `tools/b4-e2e/postflight.go` `alivePIDs` L222–259：`stub.pid`/`stub-child.pid` 任一缺失即回 `ErrInvalidAgentRunnerEvidence` 包裹错误，绝不读作「无存活进程」）。f2 三腿 verdict 类/四态与基线零漂移（uncertain-or-paused/PAUSED），failedChecks 仅该 1 项。

## B. 修复方案

最小改法：在 `case "hang"` 分支内、`run.writeOwnPID()`（L114）之后、`time.Sleep(time.Hour)`（L115）之前插入 1 次调用 + 3 行注释（hang 从不 spawn 子进程 ⇒ 空报告是其**真实事实**）。精确 diff：

```diff
 	case "hang":
 		// No artifacts beyond its own pid: the run is stopped from the outside.
 		run.writeOwnPID()
+		// A hang never spawns a child, so the empty report is the positive fact
+		// the constant-write contract requires: consumers can tell "no child"
+		// apart from "report never written" (the latter must stay fail-closed).
+		run.writeEmptyChildPIDReport()
 		time.Sleep(time.Hour)
 		return 0
```

性质说明：

- 这是**让实现符合 CONTRACTS A6（§7）恒定写**：hang 是「无子进程」运行，恒定写口径要求执行到写入点时写 0 字节空文件。修复前 hang 违反该口径（漏写）；修复后与 run 等模式同口径。
- 不是绕过断言：`alivePIDs` 的缺失判据与 postflight Leg 7 双报告对账（`postflight.go` L113–126、L128–148）都**按契约执行**。修复后 `stub-child.pid` 存在且 0 字节：`alivePIDs` 两报告 present、子报告解析出 0 个 pid，仅对 `stub.pid` 内 pid 做存活检查——watchdog 已 kill 该 pid ⇒ gone=true 是**契约成立后的真实测量**，非豁免。
- 消费方三态语义变化（hang 专属）：修复前 hang 运行 = 「缺失」⇒ 证据缺口 ⇒ `uncertain` fail-closed；修复后 = 「存在且 0 字节」⇒ 确实无子进程（事实）。语义从「无法证明」升级为「可证明」，正是恒定写的目的；不改变其它模式三态。
- 边界：不改 harness（`tools/b4-e2e/**` 零改动）、不改 `crash`/`fail-fast`/`child`/spawn 族模式（各有契约理由：crash/fail-fast 建模「写之前死」；child 是子代镜像；spawn 族写真实 pid）、不改任何生产代码（仅测试工具 `tools/agent-stub`）。

## C. 回归测试方案（main_test.go）

新测试 `TestHangWritesEmptyChildPIDReport`（放在 `TestSpawnWritesChildPIDList` 之后，复用 `buildStub`/`readReport`/`childPIDReportName` 等既有常量）：

关键点：既有 `runStub`（L124）是同步形态（`cmd.CombinedOutput()` 阻塞到进程退出），而 hang 睡眠 1h 不退出 ⇒ 新测试**不能**用 `runStub`，须后台启动 + 轮询 + kill：

```go
func TestHangWritesEmptyChildPIDReport(t *testing.T) {
	dir := t.TempDir()
	cmd := exec.Command(buildStub(t), "-mode", "hang")
	cmd.Env = append(os.Environ(), evidenceDirVarName+"="+dir)
	if err := cmd.Start(); err != nil { t.Fatalf(...) }
	defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }() // 兜底清理，防泄漏 1h 进程

	childPath := filepath.Join(dir, childPIDReportName)
	deadline := time.Now().Add(5 * time.Second)   // 报告必须在写入点立即出现
	for {
		info, err := os.Stat(childPath)
		if err == nil {
			if info.IsDir() { t.Fatalf("%s is a directory, want a 0-byte regular file", childPath) }
			if info.Size() != 0 { t.Fatalf("%s has %d bytes, want exactly 0", childPath, info.Size()) }
			break
		}
		if !time.Now().Before(deadline) {
			t.Fatalf("child pid report %s must exist after a hang run started (constant write)", childPath)
		}
		time.Sleep(20 * time.Millisecond)
	}

	own := readReport(t, filepath.Join(dir, ownPIDReportName))
	if _, err := strconv.Atoi(strings.TrimSpace(string(own))); err != nil {
		t.Fatalf("own pid report holds %q, want a decimal pid: %v", own, err)
	}
}
```

断言三要素：`stub-child.pid` 存在、0 字节、非目录；`stub.pid` 存在且十进制非空。两报告同时存在 = 恒定写对 hang 成立（与 CONTRACTS 三态「空即无子进程」一致）。

变异检验（机械可证伪）：删除插入行（恢复 L112–116 原样）重跑本测试——`stub-child.pid` 永不出现 ⇒ 5s 轮询截止触发 `t.Fatalf("... must exist ...")` ⇒ **变红**。变异操作 = 仅删 `run.writeEmptyChildPIDReport()` 一行（连同其注释）；预期失败 = 上述轮询超时断言。

既有测试不回归：4 个既有测试（L159 run 空报告、L194 spawn 写子 pid、L263 不可写 fail-closed、L299 spawn 不可写杀子）只走 run/spawn 路径与同步 `runStub`；本改动仅插入 hang 分支，与这些路径无交集；`go test ./tools/agent-stub` 应 5/5 通过。

## D. 重跑方案（layer3 F2）

- 命令（Windows 本机，仓库根执行，PS 5.1 兼容）：
  `powershell -ExecutionPolicy Bypass -File tools\b4-e2e\run-faces.ps1 -OnlyFace f2`
  `-OnlyFace f2` 精确命中三腿组 f2-pos/f2-neg/f2-neg-r2（`$groups` 过滤 `$group.face -ne $OnlyFace` 即 skip；三腿各单 phase `run`）。
- 运行前置：默认构建三件到 `tmp/b4/bin`（`prfrail.exe`/`agent-stub.exe`/`b4-e2e.exe`；已有产物可 `-SkipBuild` 复用）；自动准备 `tmp/b4`（input/runs/evidence/seed）；每组 run root 先 `Remove-Item` 重建；结尾脚本自检 `[Parser]::ParseFile` 须 0 错误。
- 输出路径：聚合 `tmp/b4/pack-verdict.json`（`failedChecks`/`deviations`/`cases`）+ `tmp/b4/SHA256SUMS.txt` + `tmp/b4/driver.log`；逐腿 `tmp/b4/evidence/f2-pos/verdict-run.json`、`f2-neg/`、`f2-neg-r2/`；每条命令/退出码/stdout/stderr 记入 `invocations`。
- 红基线（修复前，预期与 `l3-baseline-compare.txt` 一致）：`failedChecks = ["f2-pos/process-tree-taken-away"]`；f2-pos checks=6、f2-neg/f2-neg-r2 checks=4；三腿 verdict 类=uncertain-or-paused、chainState=PAUSED（红基线只记录，不构成判定依据）。
- 绿结果（修复后同一命令）：`failedChecks = []`；`process-tree-taken-away = true`（`recordedPIDsGone`→`alivePIDs`：两报告 present、0 存活）；三腿 verdict 类/chainState/taskState/stepState/observations 与 B4 基线逐字段一致（对比 `docs/validation/evidence/directive-history/b4/archive/pack-verdict.json`，generatedAt 2026-09-22T16:57:51.437Z，failedChecks 为空）。
- 建议同时重跑 `-OnlyFace f1` 作 six-facts 对照腿（`six-facts-distinct` 仅存在于 f1 面）——见 E④。

## E. ⑦ 判定问题清单（不预填答案）

⑦（独立盲审）须按起跑消息 3.3 判定四项，各自证据路径：

1. **F2 三腿 verdict 类/chainState/taskState/stepState/observations 与 B4 基线逐字段一致？**
   证据：新 `tmp/b4/pack-verdict.json` cases[face=f2] vs 基线 `docs/validation/evidence/directive-history/b4/archive/pack-verdict.json`；历史漂移记录 `docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt`。
2. **`process-tree-taken-away` 回到 true？**
   证据：`tmp/b4/evidence/f2-pos/verdict-run.json` 内该 check 的 pass/detail（期望 detail `reported pids gone: true`）；旧红记录同上 txt。
3. **`failedChecks` 空集？**
   证据：新 `tmp/b4/pack-verdict.json` 顶层 `failedChecks`（期望 `[]`）；对照 B4 基线同字段。
4. **`six-facts-distinct` 等其它断言是否受影响？**
   事实侧注（供 ⑦ 核对，非结论）：`six-facts-distinct` 仅存在于 f1 面 checks（`main.go` L690–697）；f2 面 checks 集 = outcome-not-completed / no-frozen-facts /（pos 另加）timeout-fired、error-evidence-present / process-tree-taken-away / verdict-class-recorded——本片不改 harness，f2 checks 集不变。P61 §2.4 记录的「F2 错误证据 5→6 枚」随缺口消失预期回 5，属 observations 恢复而非新漂移，是否计入问题① 由 ⑦ 判定。建议附 `-OnlyFace f1` 重跑结果作对照证据。
5. **最终判定：前提①满足 / 不满足**（⑦ 唯一主体，两读独立收敛；主控与①不预填、不下结论）。

## F. 风险与边界

- **hang 进程须被外部 kill**：已确认存在——harness 侧 `waitAndMap`（`tools/b4-e2e/runmap.go` L233–288）按 spec 超时 watch（f2/pos：timeout=5s，`main.go` L70）后走 `mirrorNeedsDurableStop` → durable stop（grace=2s），`stopOutcome`/`StopProof` 记入 verdict（`main.go` L557–560）；stub 侧 hang 本身永不自杀，依赖该停机链。
- **改动范围**：仅测试工具 `tools/agent-stub/main.go` 一行插入（+3 行注释）与 `main_test.go` 一个新增测试；不触 wire/schema/判据/生产代码；F2 检查表达式在 `e8cbff6` 与 HEAD 逐字节相同（本片不改 harness，保持该事实）。
- **与 P63 观察项 1（R-08.5 要素②过严）的关系**：本片从根修复 F2 漂移，修复后「继承漂移豁免是否成立」问题自然消失（不再存在待豁免的漂移），不重判五要素；观察项 1 仍留作 v53 治理片候选，本片不评估、不触碰判据。
- 残留风险：hang 分支新增写入使「写失败」路径多一个 `os.Exit(11)` 可能（与 run 模式同构，属既有 fail-closed 语义）；`crash`/`fail-fast`/`child`/spawn 族模式行为零变化。

DESIGN: DONE
