# CI 扩包方案 — RACE-LEG-EXPAND（① 架构 · 逐字归档）

> 载体：① 架构 `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`。以下为代理输出原文，未改写。

---

# CI 扩包方案 — RACE-LEG-EXPAND（① 架构）

## 一、cgo 两腿影响评估

1. **Ubuntu 腿：`CGO_ENABLED=1` 行内覆盖仍然正确且必要。** `ci.yml:46` 顶层 `CGO_ENABLED: "0"`，Race 步骤 `ci.yml:84` 用 `run:` 行内 `CGO_ENABLED=1 go test ...` 前缀覆盖（bash 环境前缀，仅作用于该命令）。扩包只追加 `run:` 行尾部包清单，前缀不动，语义不变。`-race` 依赖 cgo 构建的 runtime/race 库，`ubuntu-latest` 镜像自带 gcc，历史三次 Ubuntu Race 全绿（t028×2、haiku-bump）实证可行，扩包不改变该前提。不得把 `CGO_ENABLED=1` 上移到 job 级——那会波及 Windows 腿的 Build/Vet/Test 步骤，破坏既有设计。
2. **Windows 腿：完全不受影响。** Race 步骤带 `if: runner.os == 'Linux'`（`ci.yml:83`），Windows 腿直接跳过、不解析 `run:` 内容；本片仅修改 `run:` 行字符串，不触碰 `if:`、不触碰 Windows 腿任何步骤。
3. **无波及 Windows 腿的风险点。** 扩包不改 env、不改 matrix、不改其它步骤。

## 二、超时风险评估

**历史证据**（`docs/validation/evidence/directive-history/t028/archive/ci-log.txt`，Ubuntu 腿）：
- 非 race（Test 步骤）：console 4.065s（行 925）、gates 0.588s（927）、guard 0.618s（928）、tickets 0.008s（933）、adapters 13.001s（922）、chain 0.253s（924）。
- race（Race 步骤）：adapters 83.702s（943）、chain 1.703s（944）。另两份日志同量级（adapters 84.3s/84.1s，console 2.1–3.5s，guard 0.6–1.4s）。
- race 放大系数实测 6.0–7.2×。取保守 8× 估计：console ≈33s、gates ≈5s、guard ≈12s、tickets ≈0.1s（子进程 re-exec race 二进制启动略慢，给 5s 裕度）。**四新包合计 ≈50s；Race 步骤总时长 ≈84+2+50 ≈136s；2 倍安全裕度 ≈272s ≈4.5 分钟。**

**长时测试核查**：`internal/console|gates|guard|tickets` 六处 `time.Sleep(time.Hour)` 与 `time.Minute`（guard/process_test.go:224、gates/executor_test.go:95–107、console/runtime_agent_concurrent_e2e_test.go:43、guard/process_linux_test.go:31、guard/process_reclaim_test.go:22）全部位于 **re-exec 子进程 helper**（子进程自身存活期，被父测试 SIGKILL/终止，父测试不睡满）；`testing.Short()` 六包零命中。

**包级 10m 风险结论：不构成风险。** `go test` 默认 `-timeout 10m` 按**包独立计时**；最大包 adapters 历史 84–90s，裕度 >6×；四新包每包估计 ≤33s。DR-12 的 600s 阻塞是 Windows 专属（P40 已修复），Race 腿仅在 Ubuntu，且 guard 非 race 实测仅 0.6–1.4s。

**job 级 timeout-minutes：无。** 复读全文确认 prfrail `ci.yml` 中 `timeout-minutes` **零命中**（唯一的 `timeout-minutes: 15` 在 sessbridge 仓库 ci.yml:97，与本仓库无关）；GitHub 默认 job 上限 360 分钟。test job（`ci.yml:50–94`）无显式超时，无风险。

## 三、ci.yml 精确修改

- 修改前（逐字，`ci.yml:84`）：
```
        run: CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/...
```
- 修改后（逐字）：
```
        run: CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/... ./internal/console/... ./internal/gates/... ./internal/guard/... ./internal/tickets/...
```
仅此一行改动，无其它改动。**保持 `-count=1`**：`go test` 结果缓存对 `-race` 同样生效，`-count=1` 保证每次真实执行，与既有步骤一致。**不引入 `./internal/...` 通配**：会把 applier/evidence/repair/snapshot/taskdef/release 等未批准包拖入，越界，禁止。包序按裁定六包顺序。

## 四、YAML 语法与 if 条件

1. `run:` 值是块标量，含 `=`、`-`、`/` 与空格，但**无冒号+空格组合、无 `#`**（路径中 `./` 与 `CGO_ENABLED=1` 均非 YAML 特殊结构），**无需引号**。
2. `if: runner.os == 'Linux'`（`ci.yml:83`）**保持不变**。
3. 步骤顺序（Build→Vet→Test→Race→Contract fixtures→Gates，`ci.yml:73–94`）不变，不删除任何既有步骤。

## 五、否决 / 降级条件检查

1. 不删任何 CI 步骤 ✓；不改生产代码 ✓；仅改 Race `run:` 行 ✓。
2. 无活体一致性校验引用该行：`tools/gates` 不含 `ci.yml` 引用（零命中）；其余 22 处匹配均为历史证据文档/归档日志（描述当时状态，不失效）；`P44...review.md:82` 遗留项正是本片动因。
3. 超时/环境：§二已排除。
4. **残余风险（非否决）**：四新包首次进 `-race` 可能暴露**未发现的真实数据竞争**导致 CI 失败——这正是本片验证目标；但若修复需改生产代码，则越出「不改生产代码」边界，须上报升级，不得在本片擅自修。guard 子进程测试在 `-race` 二进制 re-exec 下存在轻微 flake 可能（5s WithTimeout 裕度大，风险低）；若首跑出现，由 ② 报告并按裁定降级处置。
5. **结论：无客观否决因素，批准六包全量扩包。**

## 六、结论

批准：授权 ② 按 §三「修改后」文本实施——仅将 `ci.yml:84` 的 `run:` 行扩为六包清单，保持 `CGO_ENABLED=1` 前缀、`-race -count=1`、`if: runner.os == 'Linux'` 不变，不删除任何既有步骤，不改生产代码。
