# GLM-1 — ci.yml 结构与 Race 腿历史事实检索（批量辅助 1/2 · 逐字归档）

> 载体：`deep-reasoner` + `GLM-5.3-Flash (glm)`。任务：只读检索。以下为代理汇报原文，未改写。

---

检索完成，全程只读。以下为汇报。

## A. ci.yml Race 腿现状

**Race 步骤原文**（`.github/workflows/ci.yml`）：
```yaml
      - name: Race
        if: runner.os == 'Linux'
        run: CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/...
```

**workflow 顶层 env 块原文**：
```yaml
env:
  GO_VERSION: "1.22.12"
  CGO_ENABLED: "0"
  GOTOOLCHAIN: "local"
```

- **matrix / runs-on**：`test` job `strategy.matrix.os: [ubuntu-latest, windows-latest]`，`runs-on: ${{ matrix.os }}` ⇒ 双腿（ubuntu + windows）均含。
- **if 条件**：Race 步骤 `if: runner.os == 'Linux'` ⇒ 仅 Ubuntu 腿执行，Windows 腿跳过。
- **步骤依赖**：GitHub Actions 步骤严格顺序执行、失败即中止 ⇒ 顺序为 Checkout → Setup Go → Build → Vet → Test → **Race** → Contract fixtures（也是 Linux-only）→ Gates。Race 无显式 `needs`，纯顺序链：Race 失败 ⇒ 后续 Contract fixtures/Gates 不执行。
- **其它 CGO 配置**：`test` job 内无其它步骤显式写 CGO/cgo；全局 `env.CGO_ENABLED: "0"` 使两腿所有步骤（含 Windows 腿的 Build/Vet/Test）都以 CGO_ENABLED=0 运行，仅 Race 步骤行内覆盖为 `CGO_ENABLED=1`。

## B. Race 腿历史事实

- **[docs/validation/p45-e2e-zero-model.md:166]** 「…CI `41315652405` 首跑双腿全绿；ubuntu 腿 `ok chain 0.564s` / `ok console 2.856s` / `ok adapters 12.291s` / `ok gates 0.159s` + **Race 腿 `ok adapters 89.514s` / `ok chain 2.521s`**，windows 腿 `ok chain 5.011s` / `ok console 20.796s` / `ok adapters 34.297s` / `ok gates 0.786s`…」（提交 `4273d6e`）——唯一 Race 腿逐包耗时记录。
- **[docs/validation/p45-e2e-zero-model.md:168]** 「`-race` 腿扩包至 `internal/console`：**本机缺 cgo**，并发腿未过竞态检测。」——p45 候选登记，P49 的直接依据。
- **[docs/t027/REMAINING_SLICES.md:82]** 「| P49 RACE-LEG-EXPAND | 🔄 进行中 | P48 | M | — | docs/validation/race-leg-expand{,_EN}.md（进行中） |」
- **[docs/t027/REMAINING_SLICES.md:1494]** 「目标 | 把 CI 的 Race 腿（`CGO_ENABLED=1 go test -race`）从 `internal/adapters`、`internal/chain` 扩包到 `internal/console`、`internal/gates`、`internal/guard`、`internal/tickets`（共六个包）」
- **[docs/t027/REMAINING_SLICES.md:1498]** 「现状缺口 | Race 腿只覆盖 `internal/adapters` 与 `internal/chain`；`console`/`gates`/`guard`/`tickets` 四包并发腿未过竞态检测（p45 候选登记）」
- **[docs/t027/REMAINING_SLICES.md:1500]** 验收证据：③ 先行断言证明当前 CI 未覆盖四包 Race；`ci.yml` 修改后本地跑通扩包后 race 测试；YAML 语法检查通过；build/vet/test 全绿；TOTAL_FAIL=0 且 SELFTEST 697/697。[:1501] 边界：仅改 Race 腿目标包，不删步骤、不改生产代码。
- **[docs/t027/REMAINING_SLICES.md:352-354]** A7 run `35118531930` 全绿：Ubuntu Race 步骤逐字同现 ci.yml 命令；进套件 5 项 A7 测试已在 Linux `-race` 下执行；边界订正「仅 `a7native` 不进 CI」。
- **[docs/t027/REMAINING_SLICES.md:1523]** 一次 ubuntu 腿 Test 步红（attempt 2 复绿）登记：「同 run 的 Race 腿覆盖 `internal/adapters` 且绿」（DR13-FIX 已修复）。
- **[docs/DEV_PLAN.md:81]** T014 补充验证：2026-09-08 Ubuntu VM 远端全量 build/vet/test 与 race（**gcc 全仓、clang 关键包**）全部通过——历史上全仓 `-race` 曾在 VM 验证可行。
- **[docs/DEV_PLAN.md:120]** A1：「Linux CI 新增 `-race` 门禁并同步工作流冻结契约」（run `34809403487`）——Race 门禁起源。
- **[docs/DEV_PLAN.md:136]** 代理包 15 项测试进两腿默认套件但「**不在 `-race` 子集**」——印证 race 子集仅含 adapters/chain。
- **docs/ARCHITECTURE.md**：`race|Race|CGO_ENABLED|cgo` 0 命中。

## C. 六包测试规模（`*_test.go`，Recurse）

| 包 | 文件数 | 总 KB | 并发关键词命中文件数 |
|---|---|---|---|
| adapters | 55 | 662.7 | 28 |
| chain | 14 | 195.7 | 0 |
| console | 13 | 191.9 | 3 |
| gates | 7 | 56.2 | 2 |
| guard | 8 | 41.9 | 4 |
| tickets | 5 | 76.2 | 4 |

关键词：`race`/`Goroutine`/`Concurrent`/`Parallel`/`sync.`/`go func`（大小写不敏感；chain 零命中已经 grep 交叉复核确认，非统计错误）。

命中文件清单：
- **adapters (28)**：agent_runner_admission、agent_runner_b3b_candidates_windows、agent_runner_b3b_variant_degraded、agent_runner_capability、agent_runner_completion、agent_runner_enforcement、agent_runner_event、agent_runner_launch_receipt、agent_runner_native_a6、agent_runner_native_a7_fencing、agent_runner_native_a7、agent_runner_native_b3b_windows、agent_runner_native_race、agent_runner_pinned_cli、agent_runner_postflight、agent_runner_replay_dispatcher、agent_runner_replay_root、agent_runner_replay_root_windows、agent_runner_replay_store_terminals、agent_runner_replay_store、agent_runner_run、agent_runner_starting_marker_order、agent_runner_terminal_publisher、agent_runner_timeout（以上均 `_test.go`）+ ai_availability_test.go、spawn_recovery_test.go、spawn_recovery_live_host_unix_test.go、spawn_recovery_live_host_windows_test.go
- **chain (0)**：无
- **console (3)**：runtime_agent_concurrent_e2e_test.go、runtime_agent_e2e_test.go、stopper_test.go
- **gates (2)**：executor_test.go、probe_permit_test.go
- **guard (4)**：lease_test.go、process_reaper_test.go、process_spawn_ownership_test.go、process_test.go
- **tickets (4)**：cost_test.go、durable_cost_ledger_test.go、durable_cost_lock_order_test.go、durable_cost_lock_reentry_test.go

## 结论一句话

当前 Race 腿（Linux-only，`CGO_ENABLED=1 go test -race -count=1`）仅覆盖 `internal/adapters` 与 `internal/chain`；缺 `internal/console`、`internal/gates`、`internal/guard`、`internal/tickets` 四包（其中 chain 测试面 14 文件零显式并发关键词，console/gates/guard/tickets 分别有 3/2/4/4 个含并发语义的测试文件等待纳入竞态检测）。
