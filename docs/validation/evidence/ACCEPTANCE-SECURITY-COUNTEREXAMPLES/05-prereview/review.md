# ⑤ 预审结论原文（逐字归档）

载体：runSubagent(agentName="deep-reasoner", model="DeepSeek V4 Pro (deepseek)")。只读审查，未修改任何文件。

核验完毕。全部事实已收集：7 个测试文件独立重跑全绿、gofmt/vet 干净、编码门禁合规、零生产代码改动、击穿证据与 skip 文案逐字核对。以下为 ⑤ 预审结论。

---

PRE-REVIEW: PASS

## A. 测试造假审计 — 通过
- 包名全部与被测包同名（`evidence`/`adapters`/`gates`/`chain`），获得包内访问（`boundDirVerify`、`cleanBoundRelativePath`、`permit.state`、`acceptFromExistingSeal` 路径、`taskKey` 等）。无 `_test` 外部包伪装。
- Build tag 与生产文件逐一对齐并实测：`bounddir_identity_test.go` `!darwin`（因 `cleanBoundRelativePath` 仅 `bounddir_clean.go` 定义、darwin 无该符号；darwin 桩 `boundDirVerify` 返回 `ErrBoundDirUnsupported`，exclude 合理）；`bounddir_escape_unix_test.go` `!windows && !darwin` ≡ `bounddir_unix.go`；`bounddir_escape_windows_test.go` `windows` ≡ `bounddir_windows.go`。`GOOS=linux go vet ./internal/adapters/` 独立复跑通过。
- 行号引用经独立核实真实：`engine.go:230` 调用 `acceptFromExistingSeal`、`engine.go:352-359` K 比较与 `evidence-seal-mismatch` 报修（grep 实证）；`boundread.go` 三哨兵与拒绝分支（nil/负 bound/目录/超界/关闭后）逐条与测试断言路径吻合；`probe_permit.go` Require nil/nil-ctx/unknown-state/persistClaim、`readiness.go` seal-before-AT23 顺序与 TrimSpace 不对称均与断言一致。
- 未发现恒真断言或从未执行的错误分支：独立 `go test -count=1` 四包全绿，verbose 逐条 PASS（junction 两条真实执行未 skip、跨进程 0.17s、T60-1 两个子测试 0.03s）。

## B. 哨兵纪律 — 通过
- 应绿反例全部 `errors.Is` 绑定特定哨兵：T55（三哨兵）、T64（`ErrProbePermitMissing`/`ErrProbePermitConsumed`）、T70（`ErrReadinessEvidenceIncomplete`）、T60（`ErrSealIncomplete`/`ErrPaused`）。
- 例外合理且有注释依据：**T64-4** 三连否定（契约原文"读失败无哨兵"写入测试注释，并补 `strings.Contains(err,"load probe permit")` 锚定）；**T55-6** 文案断言（`boundDirVerify` 无哨兵，红基线②注明"属预期设计"）；**T51-1** skip 后断言落在攻击结果（P53 才引入 openat 形态错误，注释写明）；另 **T52-1** 文案断言（`confirmOpenedArtifact` 无哨兵，文件头注释明示）、**T55-7**（`cleanBoundRelativePath` 无哨兵，红基线②注明）同属此列，均达标。

## C. 攻击载荷真实性 — 通过
- **T52-1/2**：真实 `cmd /c mklink /J` 建 junction、真实界外 `O/secret` 与界内 `D/sub` 文件、经 `CreateFile` 跟随 reparse 后 `GetFinalPathNameByHandle` 前缀检查拒绝/放行；mklink 失败才 skip（先例注明）。本机实测 0.18s/0.08s 双 PASS，非 mock。
- **T64-6**：真实双进程（`exec.Command(os.Args[0], "-test.run=...Helper")` + 环境变量传 marker），子进程走真实 `LoadProbePermit`→`Require`（O_CREATE|O_EXCL 落盘）→`os.Exit(0)`；父进程 `LoadProbePermit` 后 `Require` 被拒 `ErrProbePermitConsumed`。线性化点真实。
- **T60-1**：走生产入口 `Engine.Run` → `runTask` → `acceptFromExistingSeal`（engine.go:230），封条由真实 `FileEventStore.SealTaskEvidence` 持久化；拒绝发生在任何伪端口之前（除 `Workspaces.Materialize` 供 `TaskContextKey` 的 workspaceRoot，对照子测试证明派生一致）。等效性依据**成立**。

## D. 零生产代码改动 — 通过
- 实测 `git diff --cached --name-only -- internal/`：恰 7 个文件、全部 `_test.go`，非 `_test.go` 改动为**空**。全片 22 个暂存文件均为文档/证据包/测试，无 `go.mod` 等其它变更。

## E. GAP 覆盖对照 — 通过
GAP-51→T51-1（skip 记录击穿 GAP-77）、GAP-52→T52-1/2、GAP-55→T55-1..7、GAP-60→T60-1..3、GAP-64→T64-1..6、GAP-70→T70-1..4。六条全部有对应反例，无遗漏；GAP-78/79 现状钉（T60-3/T70-4）与 GAP-79/80 注释登记齐全。

## F. T51-1 skip 合规性 — 通过
- skip 文案逐字为 `t.Skip("Known vulnerability GAP-77: tracked for P53")`，位于被测调用**之前**（函数首行，比要求更保守）。
- P53 启用路径自洽：删 skip 后 `OpenBoundDirFile` 修复为逐组件 openat 拒绝 → `err != nil` → 断言通过；攻击结果断言方向（`err == nil` ⇒ Fatal）正确。击穿证据 `05-attack/t51-1-linux-escape.md` 含逐字探针源码与 Linux 实跑输出（23 字节界外字节）。

## G. 测试稳健性 — 通过
- 全部 fixture 用 `t.TempDir()`；junction 以 `defer os.Remove(link)` 摘除；句柄显式/延迟 Close；T55-2 第二写者失败才 skip（本机实跑未 skip，超界增长路径真实触达）；helper 无 env 即 skip，正常 `go test ./...` 不会误触发。
- 小注：T64-6 子进程 `CombinedOutput` 无显式超时，依赖 go test 全局超时；实测 0.17s，风险可忽略（Info）。

---

**发现汇总**
- Critical: 0
- High: 0
- Medium: 0
- Low: 2
  1. **L1（input-list.md 文档级）**：GAP-70 代码面归属错误——`internal/chain/seal.go` 被描述为「GoalProduction / GoalAcceptance 判定入口」、`seal_verify.go` 被描述为「AT23Hash / SealBinding 校验」；两者实测无这些符号（seal_verify.go 是恢复期封条核验），真实主面 `internal/gates/readiness.go` 反未列入输入清单。不影响测试正确性，建议 ⑧ 文档轮勘正。
  2. **L2（设计-实现偏差未登记）**：① 设计 T64-6 规格含 claimed/release 信号握手与「父新实例再 Require」两段；② 实现简化为顺序式（子退出→父 Load+Require），S3 偏离清单未登记。覆盖目标（跨进程持久声明、进程死亡不复活配额）仍达成，属登记缺口。
- Info: 3
  1. 无哨兵断言例外实际含 4 处（T52-1/T55-6/T55-7/T64-4），均有文件注释或红基线依据；T51-1 为 skip+攻击结果断言。
  2. T64-6 子进程无显式超时（依赖 go test 全局超时），实测 0.17s。
  3. T55-4 本机走 Windows `fs.ErrPermission` 形态，unix 形态由 Linux CI 腿覆盖；双形态均为 `errors.Is` 绑定。

**独立验证记录**：`gofmt -l` 7 文件无输出；7 文件 BOM=false/CR=false；`go vet` 四包 VET-OK；`GOOS=linux go vet ./internal/adapters/` 通过；`go test -count=1` 四包独立复跑全绿；`zz_probe*` 残留 0；`git status` 无异常条目。

**结论行：PRE-REVIEW: PASS**（Critical 0 / High 0 / Medium 0 / Low 2 / Info 3；无 Medium+，⑤ 原位重跑条件不触发，可进入 ⑥）

PRE-REVIEW: END
