# DR13-FIX-D3 验证报告

- 定位：v37 层 1 第 1 项；P45 前置 `DR13-FIX` 的 D3 遗留片。
- 准则：v1.60；起跑锚点 `e742224`；日期 2026-11-02。
- 实现类：`[SLICE]`；规模档 S；台账行 `P48`。

## 1. 目标与范围

- **目标**：`release()` 的释放序改为「先 `os.Remove`，再 unlock，后 Close」，配合既有 D1 同进程所有权门与 D2 Stat 复核，闭合毫秒级并发窗口 W1。
- **落点**：`internal/adapters/agent_runner_starting_marker.go:108` 函数起始；`:112` 先删；`:113` 再解锁；`:114` 后关闭。
- **回归边界**：Windows remove-while-open 全量回归；双平台编译与全包测试。
- **不改边界**：D1 / D2 实现、`winners==1` 断言、CONTRACTS / Schema / fixtures、`ci.yml`；不推 gitee。

## 2. 流水线执行

| 步骤 | agentName | model | 产出 |
|---|---|---|---|
| 0 | 主控 | `DeepSeek V4 Pro (deepseek)` | 台账登记、证据包骨架、自由度清单 |
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `01-arch/design.md` |
| GLM 1/2 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `09-glm/glm-1-case-scan.md` |
| GLM 2/2 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `09-glm/glm-2-impact-scan.md` |
| ③ | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | 红基线、新增用例集 |
| ② | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 生产改动、`04-impl/implementer-report.md` |
| ④ | `prfrail-implementer` / `prfrail-tester` | 同上 | `go build` / `go vet` / `go test ./...` 全绿 |
| ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `PRE-REVIEW: PASS` |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS` |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `FINAL REVIEW: PASS` |
| ⑧ | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | 本报告对与文档同步 |

## 3. 红绿反转证据

- 红测首跑（未修复代码）：`TestStartingMarkerReleaseRemovesBeforeUnlock` 判 FAIL，判据为 unlock 位于 `os.Remove` 之前；日志 `docs/validation/evidence/DR13-FIX-D3/03-tests/logs/go-test-red.txt`（原始过程件经 ⑧c 收档）。
- 压力红证（补充，概率口径）：`TestStartingMarkerReleaseNeverExposesReleasedContent` 首跑命中违例 4/35，日志同上。
- 主控独立复核：`docs/validation/evidence/DR13-FIX-D3/03-tests/logs/maincontroller-red-verify.txt`。
- 修复后聚焦转绿：`docs/validation/evidence/DR13-FIX-D3/03-tests/logs/mc-focused-green.txt`，6/6 全绿。
- 绿基线：remove-while-open 直测、幂等直测、Windows fail-closed 占名反测、ShareFact 探针。

## 4. 审查链结果

- ⑤：`PRE-REVIEW: PASS`（2 Low + 3 Info）。Low-1 顺序门假绿理论形态，交 ⑦ 确认且不改代码；Low-2 unix 运行时腿未本机验证，由双平台 CI 兜底；Info 三项无修法。
- ⑥：`INDEPENDENT SCAN: PASS`，发现 0；独立实测 6/6 全绿。
- ⑦：`FINAL REVIEW: PASS`，发现 0；登记非阻断观察项——顺序门按最早调用位置判定，若出现不可达分支前置调用的反例且行为测试未拦住，应升级为缺陷并补可达路径级断言。

```sixstage
## Section A

The verdict line of the transcript is merged into the closing line of this fence.

## Section B

(no findings)

## Section C

From the code and the targeted tests, the D3 reorder holds and introduces no new concurrency window or caller-semantics drift. The release sequence is remove first, then unlock, then close; when removal fails the lock and handle are released unconditionally, so the failure direction is fail-closed. The boundary audit found only the release function changed; the return semantics are unchanged.

## Section D

The judgment is overturned by any evidence that unlock runs before the remove, or that the lock stays held after a failed removal; the structural gate and the fail-closed counter-test fail immediately in those cases.

## Section E

The unix runtime leg was not verified on this host; a dual-platform CI failure in a targeted path overturns the pass.

INDEPENDENT SCAN: PASS
```
差异披露：⑥ 原文 Section A 的结论行与本围栏块末行合并，差异恰为一处。

## 5. 文档同步更新

- `docs/ARCHITECTURE.md` 与 `docs/ARCHITECTURE_EN.md` 的 §4.2 各增补锁序条目：先物理删除路径再释放锁，使扫描者无法观察到「锁空闲 + 路径仍在」；删除失败时锁无条件释放，残留由所有权门与下一次扫描收敛。
- `docs/DEV_PLAN.md` 与 `docs/DEV_PLAN_EN.md`：无 DR13 / starting marker / P48 相关里程碑描述，未改。

## 6. 已知限制

- 顺序门为源码顺序门，不是运行时证明；理论假绿形态见 ⑦ 观察项。
- 压力回归为概率口径（R3）：3s 窗口、无 `-race`，零违例不构成严格证明，结构保证来自顺序门与 D1 / D2 组合论证。
- unix 运行时腿本机未执行（Windows 宿主）；`RefusesLiveHostUnix` 与孤儿收养路径以双平台 CI 为最终判据。
- 事实登记 R-N1：本工具链 `os.OpenFile` 句柄阻止 `os.Remove`（errno 32），与 `agent_runner_starting_marker_windows.go:97-101` 注释一致；高并发下 `release` 返回 false 属设计内 R5 收敛径。

## 7. 门禁读数（⑨）

- 三作用域：index / tree / ci 均 `TOTAL_FAIL=0`（changed=27；G6-1 / G6-4 对新增行 scanned 282）。
- `SELFTEST: PASS 697/697`。
- base 三腿：gofmt 空、`go build ./...`、`go vet ./...`、`go test ./...` 全绿（ok=15、FAIL=0）。

## 8. 提交与 CI

- 主提交：`82e3ac8`（`feat: harden starting marker release order against W1 window (P48-DR13-FIX-D3)`），已推 `origin/main`（未推 gitee）。
- CI：run `41379998477`（head `82e3ac8`）**首跑双腿全绿**（attempt 1，未触发复跑）。
- 两腿 `Test` 均实测 `ok github.com/larsonzh/prfrail/internal/adapters`；ubuntu 腿 `Race` 步骤另见 `internal/adapters` 与 `internal/chain` 两行 ok。
- 两腿 `Gates` 实测：`SELFTEST: PASS 697/697`、`GATE REPORT scope=ci changed=27`、`TOTAL_FAIL=0`。
- 三个 dispatch 专用作业（candidate build / candidate probe / bootstrap release）按常态 skipped。

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/t027/CHECKLIST_v38.md	repo
docs/t027/CHECKLIST_v38_EN.md	repo
docs/validation/dr13-fix-d3.md	repo
docs/validation/dr13-fix-d3_EN.md	repo
docs/validation/evidence/DR13-FIX-D3/00-briefing.md	repo
docs/validation/evidence/DR13-FIX-D3/01-arch/design.md	repo
docs/validation/evidence/DR13-FIX-D3/03-tests/red-baseline.md	repo
docs/validation/evidence/DR13-FIX-D3/03-tests/tester-report.md	repo
docs/validation/evidence/DR13-FIX-D3/04-impl/implementer-report.md	repo
docs/validation/evidence/DR13-FIX-D3/05-prereview/review.md	repo
docs/validation/evidence/DR13-FIX-D3/06-scan/review.md	repo
docs/validation/evidence/DR13-FIX-D3/07-final/review.md	repo
docs/validation/evidence/DR13-FIX-D3/09-glm/glm-1-case-scan.md	repo
docs/validation/evidence/DR13-FIX-D3/09-glm/glm-2-impact-scan.md	repo
docs/validation/evidence/DR13-FIX-D3/MANIFEST.md	repo
docs/validation/evidence/DR13-FIX-D3/SHA256SUMS.txt	repo
docs/validation/evidence/DR13-FIX-D3/freedom-list.md	repo
docs/validation/evidence/DR13-FIX-D3/input-list.md	repo
internal/adapters/agent_runner_starting_marker.go	repo
internal/adapters/agent_runner_starting_marker_order_test.go	repo
internal/adapters/agent_runner_starting_marker_order_windows_test.go	repo
docs/validation/evidence/DR13-FIX-D3/03-tests/logs/go-test-red.txt	repo
docs/validation/evidence/DR13-FIX-D3/03-tests/logs/maincontroller-red-verify.txt	repo
docs/validation/evidence/DR13-FIX-D3/03-tests/logs/mc-focused-green.txt	repo
docs/validation/evidence/DR13-FIX-D3/08-gates/gate-readings.txt	repo
tmp/DR13-FIX-D3/	local-only
```

本片停在 ⑩ 等授权。
