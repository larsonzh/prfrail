# P74 SANDBOX-INTERFACE — 设计文档 r1

- 切片：P74 `SANDBOX-INTERFACE`
- 环节：① 设计（deep-reasoner）
- 仓库：d:\LZProjects\prfrail（Go 1.22+，module github.com/larsonzh/prfrail）
- 日期：2026-11-15
- 产出性质：本文件为唯一产出；不修改任何仓库源码或文档
- 结论：`DESIGN: DONE`（停升级边界 7 条全不触）

## 0. 引用事实（已核实，本环节不再调研）

- `internal/guard/`：无 interface、无 Registry 命名类型；`ProcessSpec` 含 `OnSpawn` 回调（spawn 前可拒绝）；`ProcessIdentity{PID,StartToken}`；`ManagedProcess` 不导出字段。
- `internal/chain/agent_runner_launcher.go:14`：`AgentRunnerLauncher` 接口与 `AgentRunnerLaunchRequest`。
- `internal/adapters/agent_runner_pinned_cli.go:218` 调 `guard.StartManaged`；`:300` 调 `guard.StopProcessIdentity`。
- `internal/gates/probe_permit.go:210` `Require`：cache 快拒绝 → `persistClaim()` 线性化点，持久化失败即 `ErrProbePermitConsumed`；`:253` `Complete` CAS 转 spent；配额不退。
- `internal/console/runtime_agent.go`：Step 0 就绪门（配额不烧）；Step 1 预检块；既有注释明言「a configuration refusal must not burn the one-shot quota」；`requireProbePermit` 在全部预检通过之后；Step 2 marker 为「Step 2: the run-root marker comes first」。
- `internal/console/runtime_supervised.go`：仅 `freshStart` 为真才 `requireProbePermit` + defer。
- `docs/CONTRACTS.md`：§7 = `## 7. AgentRunner 与辅助消息通道边界`（:210–:322）；§7 内唯一 `###` 为 `### 7.1 supervised-black-box 显式模式`（:310）；§8 起于 :322。

## 1. 接口语义定位：检查层，不是执行层

- Sandbox 检查层只回答一个问题：「当前环境是否满足 enforcement 要求」。不执行任何 enforcement。
- S1 落点：环境确定不满足（Denied）或无法判定（Unknown）⇒ fail-closed 拒绝派发；检查能力不可用（err）⇒ 同样 fail-closed。检查通过 ⇒ 不改变任何既有行为。
- 明确不做：不引入 Create/Exec/Destroy 进程全生命周期；不与 guard 现有类型重叠（`ProcessSpec`/`ProcessIdentity`/`ManagedProcess` 零改动、零新字段；不新增 Registry 类命名类型）。

| 层 | 职责 | S1 范围 |
|---|---|---|
| 检查层（本切片） | 判定环境是否满足 enforcement 要求 | 仅 fail-closed 拒绝 |
| 执行层（不在本切片） | 真实 enforcement / 生命周期 / 工具网络控制 | 不落地 |

## 2. 最终 Go 接口签名与 ③0 骨架（完整代码草案）

### 2.1 接口方法集最终签名

```go
Check(ctx context.Context, spec CheckSpec) (Verdict, error)
```

- Verdict：Allowed / Denied(reason) / Unknown(reason)；**Denied 与 Unknown 均为 fail-closed**。
- 错误分类：`ErrSandboxUnavailable`（检查能力不可用）；门函数统一产出 `ErrSandboxRejected`（非 Allowed 的 fail-closed 拒绝）。

### 2.2 internal/guard/sandbox.go（③0 骨架，完整代码草案）

```go
package guard

import (
	"context"
	"errors"
)

// ErrSandboxUnavailable 表示沙箱检查能力不可用（如平台无实现/实现缺失）。
// 调用方必须按 fail-closed 处理。
var ErrSandboxUnavailable = errors.New("guard: sandbox check unavailable")

// VerdictStatus 是沙箱检查结论的枚举。
type VerdictStatus uint8

const (
	// VerdictUnknown 表示无法判定环境是否满足 enforcement 要求（fail-closed）。
	VerdictUnknown VerdictStatus = iota
	// VerdictAllowed 表示环境满足要求，允许继续派发。
	VerdictAllowed
	// VerdictDenied 表示环境确定不满足要求（fail-closed）。
	VerdictDenied
)

func (s VerdictStatus) String() string {
	switch s {
	case VerdictAllowed:
		return "allowed"
	case VerdictDenied:
		return "denied"
	default:
		return "unknown"
	}
}

// Verdict 是 Sandbox.Check 的裁决。Denied 与 Unknown 都必须被调用方视为
// fail-closed 拒绝；只有 Allowed 允许继续。
type Verdict struct {
	Status VerdictStatus
	Reason string // 人类可读理由，随拒绝错误进入审计与诊断信息。
}

// IsAllowed 便捷判定：仅 Allowed 为 true。
func (v Verdict) IsAllowed() bool { return v.Status == VerdictAllowed }

// CheckSpec 是一次沙箱检查的输入。S1 仅携带审计标签，不含任何进程
// 生命周期参数（检查层不创建、不执行、不销毁进程）。
type CheckSpec struct {
	// Reason 是本次检查的触发原因标签（如 "agent-runner-dispatch"）。
	Reason string
}

// Sandbox 是沙箱检查层接口：只回答「当前环境是否满足 enforcement 要求」。
//
// 契约：
//   - 不提供 Create/Exec/Destroy/Terminate 等进程生命周期方法；
//     与 ProcessSpec/ProcessIdentity/ManagedProcess 无类型重叠。
//   - err != nil 或 Verdict 非 Allowed 均视为 fail-closed 拒绝。
type Sandbox interface {
	Check(ctx context.Context, spec CheckSpec) (Verdict, error)
}

// NewSandbox 返回当前平台的默认实现（见 sandbox_unsupported.go /
// sandbox_windows.go 中的 newPlatformSandbox）。
func NewSandbox() Sandbox {
	return newPlatformSandbox()
}
```

### 2.3 internal/guard/sandbox_unsupported.go（③0，`//go:build !windows`）

```go
//go:build !windows

package guard

import (
	"context"
)

// unsupportedSandbox 是 !windows 平台的检查实现：检查能力不存在，
// 恒返回 Unknown 与 ErrSandboxUnavailable（fail-closed）。
type unsupportedSandbox struct{}

func (unsupportedSandbox) Check(_ context.Context, _ CheckSpec) (Verdict, error) {
	return Verdict{
		Status: VerdictUnknown,
		Reason: "sandbox check not supported on this platform",
	}, ErrSandboxUnavailable
}

func newPlatformSandbox() Sandbox { return unsupportedSandbox{} }
```

### 2.4 internal/guard/sandbox_windows.go（③0，`//go:build windows`，形态 α）

```go
//go:build windows

package guard

import (
	"context"
)

// alphaSandbox 是 Windows 平台的 S1 实现（形态 α）：S1 未执行任何检查，
// 恒返回 Unknown——如实报告「不知道环境是否符合要求」，不谎报「确定不符合」。
// Unknown 同样被调用方 fail-closed，安全强度与 Denied 相同、诚实性更高。
type alphaSandbox struct{}

func (alphaSandbox) Check(_ context.Context, _ CheckSpec) (Verdict, error) {
	return Verdict{
		Status: VerdictUnknown,
		Reason: "S1 sandbox enforcement checks not implemented; environment compliance unknown",
	}, nil
}

func newPlatformSandbox() Sandbox { return alphaSandbox{} }
```

### 2.5 ②b 追加：门函数（追加到 sandbox.go）

```go
// ErrSandboxRejected 表示沙箱检查以 fail-closed 方式拒绝本次派发。
var ErrSandboxRejected = errors.New("guard: sandbox check rejected dispatch")

// CheckSandbox 是派发路径的统一检查门：Check 非 Allowed 或出错即返回
// fail-closed 错误。s 为 nil 时回退 NewSandbox()。
func CheckSandbox(ctx context.Context, s Sandbox, spec CheckSpec) error {
	if s == nil {
		s = NewSandbox()
	}
	v, err := s.Check(ctx, spec)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSandboxUnavailable, err)
	}
	if !v.IsAllowed() {
		return fmt.Errorf("%w: %s", ErrSandboxRejected, v.Reason)
	}
	return nil
}
```

（②b 时 `fmt` 导入一并加入；③0 骨架不含 `CheckSandbox`，仍可编译。）

编译保证：③0 三个文件落地即 `go build ./...` 通过，无任何新依赖。

## 3. CONTRACTS §7 新增子段草案

英文标题：`### 7.2 Sandbox Check-Layer Interface`

中文正文草案（追加在 §7 末尾、7.1 之后、§8 之前）：

```markdown
### 7.2 Sandbox 检查层接口

检查层（`guard.Sandbox`）位于派发路径的探针许可（probe permit）之前，只回答
「当前环境是否满足 enforcement 要求」一个问题，不执行 enforcement 本身。S1
阶段仅落实：环境不满足或无法判定时，fail-closed 拒绝派发；检查通过时不改变
任何既有行为。

- 接口：`Check(ctx, CheckSpec) (Verdict, error)`。Verdict 区分 Allowed /
  Denied(reason) / Unknown(reason)；**Denied 与 Unknown 均为 fail-closed**，
  调用方必须拒绝派发；`err != nil`（ErrSandboxUnavailable）同样 fail-closed。
- 检查点：位于探针许可 Require 之前——Sandbox 拒绝属配置/环境类拒绝，
  不得烧掉不可退还的一次性探针配额（配置拒绝不得烧配额）。
- 平台无关性：接口不绑定任何平台；无实现的平台恒返回 Unknown +
  ErrSandboxUnavailable。S1 Windows 实现恒返回 Unknown（S1 未执行任何检查，
  不得谎报「确定不符合」）。
- 边界：本接口**不改 chain/adapters 判定本体，不改 A6 spawn 语义**；检查仅
  新增一条派发前短路路径，通过后的流程与原逻辑完全一致。
- 触发模型（S1）：仅当派发显式声明要求 Sandbox enforcement 时才调用检查；
  S1 无任何声明来源，默认不触发，运行时行为与现状一致。
```

## 4. 检查点时序与调用点

### 4.1 结论：方案 (i) —— permit 之前

论证：`gates.ProbePermit.Require`（probe_permit.go:210）先查 cache 快拒绝、再
`persistClaim()` 作为线性化点，配额一次性不可退（Complete 仅 CAS 转 spent）。
Sandbox 拒绝属「配置/环境类拒绝」——runtime_agent.go 既有注释明言
「a configuration refusal must not burn the one-shot quota」。若放在 permit 之后，
一次环境拒绝就永久烧掉一条探针配额，违反该原则。故选 (i)：检查插在
`requireProbePermit` 之前、Step 2 marker 之前。

### 4.2 runtime_agent.go 插入位置与调用代码草案

- 位置：`ExecuteAgentRunnerRun` 的 Step 1 预检块内、最后一项预检
  `agentRunnerValidateCLI(options.AgentCLI)` 之后、「The probe permit quota is
  claimed only after every precondition passed」注释与 `requireProbePermit(ctx,
  options.ProbePermit)` 之前。天然位于 Step 2 marker（「Step 2: the run-root
  marker comes first」）之前 ⇒ 拒绝直接 return，不进入 Step 2、不建 run-root
  marker、不 spawn。
- `AgentRunnerAssemblyOptions` 新增字段 `SandboxEnforced bool`（S1 无任何来源置真
  ⇒ 恒 false ⇒ 默认路径零行为变化）。
- 调用代码草案：

```go
	// Sandbox 检查层（S1：显式声明触发，默认不参与）。放在 permit 之前：
	// 配置/环境类拒绝不得烧一次性探针配额。
	if options.SandboxEnforced {
		if err := preflightSandbox(ctx, "agent-runner-dispatch"); err != nil {
			return RunSummary{}, fmt.Errorf("%w: sandbox preflight failed (fail-closed): %w", ErrInvalidConfig, err)
		}
	}
```

### 4.3 runtime_supervised.go 插入位置与调用代码草案

- 位置：`ExecuteSupervisedRun` 中 `freshStart` 为真的分支内、`requireProbePermit(ctx,
  options.ProbePermit)` 之前（`defer releaseProbePermit()` 及之后逻辑不变）。
  理由：仅 freshStart 才消费不可退配额且真正进入启动流程；再入路径无派发
  动作，不做检查（避免扩大行为面）。
- `SupervisedAssemblyOptions` 新增同名字段 `SandboxEnforced bool`。
- 调用代码草案：

```go
	if freshStart {
		// Sandbox 检查层（S1：显式声明触发，默认不参与）。放在 permit 之前。
		if options.SandboxEnforced {
			if err := preflightSandbox(ctx, "supervised-runtime-dispatch"); err != nil {
				return RunSummary{}, fmt.Errorf("%w: sandbox preflight failed (fail-closed): %w", ErrInvalidConfig, err)
			}
		}
		releaseProbePermit, permitErr := requireProbePermit(ctx, options.ProbePermit)
		if permitErr != nil {
			return RunSummary{}, permitErr
		}
		defer releaseProbePermit()
	}
```

### 4.4 console 侧可测装配（②b）：internal/console/sandbox_preflight.go

```go
package console

import (
	"context"

	"github.com/larsonzh/prfrail/internal/guard"
)

// sandboxForDispatch 是派发路径的检查实现注入点（白盒测试可替换）。
var sandboxForDispatch guard.Sandbox = guard.NewSandbox()

// preflightSandbox 执行派发前沙箱检查；非 Allowed 或出错即 fail-closed。
func preflightSandbox(ctx context.Context, reason string) error {
	return guard.CheckSandbox(ctx, sandboxForDispatch, guard.CheckSpec{Reason: reason})
}
```

fail-closed 处理：Denied/Unknown ⇒ `CheckSandbox` 返回 `ErrSandboxRejected`；
err ⇒ 返回 `ErrSandboxUnavailable`；两个调用点统一
`return RunSummary{}, fmt.Errorf(...)`，不进入 Step 2 / 不进入 permit 之后的
启动流程。

## 5. 平台实现策略

- `sandbox_unsupported.go`（`//go:build !windows`）：恒
  `(Verdict{Unknown, ...}, ErrSandboxUnavailable)`。
- `sandbox_windows.go`：选 **α（恒 Unknown，nil error）**。
  理由：S1 事实上未执行任何检查，无法断言「确定不符合」；β（Denied）会把
  「没查」谎报成「查了且确认不符合」，在审计上伪造确定性。α 的 Unknown 如实
  陈述「未检查、无法判定」，且同样 fail-closed——安全强度与 β 相同、诚实性
  更高。α 还与 !windows 路径同为 Unknown，跨平台拒绝语义单一。

## 6. 红测设计（TDD 红测先行）

| # | 文件（路径） | 测试函数 | 断言要点 | 红/绿时点 |
|---|---|---|---|---|
| i | internal/guard/sandbox_interface_test.go | TestSandboxInterfaceContract | `NewSandbox()` 非 nil；三个 VerdictStatus 枚举互异；`CheckSandbox`：fakeAllowed→nil、fakeDenied→`ErrSandboxRejected`、fakeUnknown→`ErrSandboxRejected`、fakeErr→`ErrSandboxUnavailable`（均 errors.Is 判定）；`Verdict{Status:VerdictUnknown}.IsAllowed()==false` | 红：③0 前（编译失败）；绿：②b（门函数落地后全绿） |
| ii | internal/guard/sandbox_unsupported_test.go（`//go:build !windows`） | TestUnsupportedPlatformReturnsUnknown | `NewSandbox().Check` ⇒ `Status==VerdictUnknown`、`errors.Is(err, ErrSandboxUnavailable)`、`!IsAllowed()`、Reason 非空 | 红：③0 前；绿：③0（在 Linux CI leg 运行） |
| iii | internal/console/runtime_agent_sandbox_test.go | TestAgentRunnerDispatchSandboxFailClosed | 白盒替换 `sandboxForDispatch` 为记录型 fake（Denied/Unknown/err 三例），调 `preflightSandbox(ctx, "agent-runner-dispatch")`：错误分别 errors.Is `ErrSandboxRejected` / `ErrSandboxRejected` / `ErrSandboxUnavailable`；fake 观测到 `spec.Reason=="agent-runner-dispatch"`；全程不产生进程 | 红：②b 前（`preflightSandbox` 不存在，编译失败）；绿：②b |
| iv | internal/console/runtime_supervised_sandbox_test.go | TestSupervisedRuntimeSandboxFailClosed | 同上，reason=="supervised-runtime-dispatch"；连续两次调用互不残留状态 | 红：②b 前；绿：②b |

(iii)/(iv) 可测性说明（最小侵入、不改 chain/adapters 判定本体）：不调用
`ExecuteAgentRunnerRun`/`ExecuteSupervisedRun` 整条流程（避免真起进程与真实
permit 消耗），直接测试派发路径唯一调用点所依赖的装配函数
`preflightSandbox`，注入点即 console 包内可替换变量 `sandboxForDispatch`
（白盒测试直接赋值 fake）。fake 无任何 spawn 路径 ⇒ 天然不起进程。
「拒绝 ⇒ 不进入 Step 2 / 不 spawn」是插入点的结构保证（4.2/4.3 锚定在
`requireProbePermit` 与 Step 2 marker 之前），不另起进程验证。

## 7. 停升级边界 7 条逐条裁决

① 需改 CONTRACTS §7 之外契约段？**不触**。新增子段 7.2 位于 §7 内部（7.1 之后、§8 之前），不修改任何其他 § 的既有条款。
② 需改 chain/adapters 判定本体？**不触**。插入点全部在 internal/console；`AgentRunnerLauncher` 接口与 `agent_runner_pinned_cli.go`（:218 StartManaged / :300 StopProcessIdentity）零改动。
③ 需改 schema/fixtures/CI/候选版本/R-08/R-09？**不触**。仅新增 Go 源文件与测试，无新依赖；不改任何 JSON schema、fixtures、CI 配置；S1 默认零行为变化，不触及 R-08/R-09 与候选版本。
④ 需引入 P′？**不触**。无新增外部产品、组件或协议。
⑤ 需扩域（真实 enforcement/进程生命周期/工具网络控制）？**不触**。接口不含生命周期方法，S1 不执行任何真实检查。
⑥ 需改 S1 Exit 八门判定本身？**不触**。本切片默认不触发，S1 运行时行为与现状完全一致，Exit 八门无法感知。
⑦ Sandbox 检查改变现有拒绝语义？**不触**。S1 默认路径零行为变化；显式声明触发时仅新增 permit 之前的一条短路路径，permit 之后与 spawn 语义不变。

全部不触 ⇒ 设计继续，不写 STOP_ESCALATE。

## 8. ⑦ 盲审不抽的正式确认

设计者正式确认：`Sandbox.Check()` 属检查层，接口方法集仅含 `Check`，不含任何
进程生命周期方法；不修改 `ManagedProcess.Terminate`、`verifyGone`、
`StopProcessIdentity` 等停机路径及其调用方，所有权与身份语义亦不变 ⇒
本切片不触停机语义，⑦ 盲审不抽。

## 附：切片新增文件清单

- ③0：internal/guard/sandbox.go、sandbox_unsupported.go、sandbox_windows.go
- ②b：internal/console/sandbox_preflight.go；runtime_agent.go / runtime_supervised.go 各一处调用点 + 两个 options 结构各新增 `SandboxEnforced` 字段；4 个红测文件；docs/CONTRACTS.md §7.2 落盘（步骤另行排期）

DESIGN: DONE

## 9. r2 补正：可证伪性与死代码处置

> 本节取代 r1 §4.2/§4.3 的「触发」表述与 §6 红测 (iii)/(iv)，并补一条 §2.2/附录增量；
> 与 r1 已有正文冲突处以本节为准。其余（§1/§2/§3/§4.1/§4.4/§5/§7/§8）不变。

### 9.1 (A) 触发面事实核查与方案选择

只读核查事实（file:line）：

- `cmd/prfrail/main.go:13-14`：cmd 层无任何 flag 解析，仅委托 console：
  ```go
	cli := console.NewCLI(version)
	exitCode := cli.Execute(context.Background(), os.Args[1:], os.Stdout, os.Stderr)
  ```
- `internal/console/cli.go:276` `set, parseErr := newFlagSet("run")`；run flag 全集 :277-298（chain/run-dir/run-id/json/progress/request-record(s)/capability-enforcement-availability-record/authorization-ledger/cost-ledger/agent-cli/agent-version/probe-permit），**无任何 sandbox 类开关**。
- `internal/console/cli.go:410` `return cli.executeSupervisedRun(...)` / `:413` `return cli.executeAgentRunnerRun(...)` 三路分流；`:579` `options := AgentRunnerAssemblyOptions{...}` 构造点。
- `internal/console/runtime_agent.go:74-135` `AgentRunnerAssemblyOptions` 字段全集（RunID/RunRoot/Definition/RequestRecord(s)/Capability/Enforcement/Availability(+Policy)/AuthorizationLedger/CostLedger/AgentCLI/WorkspaceRoot/ProbePermit/Clock/IDs/Progress/testStub）——**无可复用触发字段**。
- `internal/console/runtime_supervised.go:217-265` `SupervisedAssemblyOptions` 同上，无 sandbox 类字段。
- `internal/console/supervised_cmd.go:34` `executeSupervisedRun`；`:39` `options.Progress = progress`、`:42` `options.ProbePermit = probePermit` 注入点；`:54-64` `supervisedAssemblyOptions` 构造点。
- 全仓库 `os.Getenv`/`LookupEnv` 仅用于测试 helper（`PROOFRAIL_*` 前缀）与 adapters 探针；产品级环境开关尚不存在。

**选择：方案 B**（新增环境变量开关）。

论证（对照判据逐条）：
1. 最小侵入：不新增 flag ⇒ `newFlagSet("run")`、runUsage、解析结构与 schema 全不动；仅新增一个读环境变量的小函数与两条路由内部签名各加一参。
2. 可达：用户可 `PRFRAIL_SANDBOX_CHECK=enforce prfrail run ...` 显式开启 ⇒ 检查调用点存在真实触发源。
3. 默认零行为变化：未设置 ≡ `off` ⇒ false ⇒ 与 r1 之前控制流一致。
4. 死代码禁令（P73 附条件 b）：`if options.SandboxEnforced` 存在明确可达入口（CLI 环境变量读取→注入），且 CONTRACTS §7.2 记载开关语义 ⇒ 非死代码。

方案 A 否决：核查证明两个 options 结构与 run flag 集均无可用触发字段；复用 `--progress` 等现成 flag 会污染语义。方案 C 不需要（B 已满足全部判据）。

### 9.2 触发面落地形态（取代 §4.2/§4.3「S1 无任何来源置真」表述）

- 开关：环境变量 `PRFRAIL_SANDBOX_CHECK` ∈ {未设置, `off`, `enforce`}；未设置 ≡ `off`（S1 默认）；**其它任何值 ⇒ CLI usage error**（fail-closed、零写入，发生在任何文件读取之前）。
- 读取点：`internal/console/cli.go` `executeRun` 内、`set.Parse` 之后、`cwd` 获取之前（约 :307），经新函数 `sandboxEnforcementFromEnv() (bool, error)`（落 `internal/console/sandbox_preflight.go`，与 §4.4 同文件）**读取一次**；结果经新增参数传入 `executeAgentRunnerRun` / `executeSupervisedRun`，分别注入 :579 options 字面量（`SandboxEnforced: sandboxEnforced`）与 supervised_cmd.go :39/:42 旁（`options.SandboxEnforced = sandboxEnforced`）。
- r1 §4.2/§4.3 已声明的 `SandboxEnforced bool` 字段保留为**程序化缝隙**：装配入口自身不读环境变量（避免测试受宿主环境污染），仅 CLI 层读一次注入；测试直接置字段即可，无需 `t.Setenv`。
- S1 语义注明：置 `enforce` 时，S1 平台实现（α 恒 Unknown / !windows 恒 Unknown+Err）⇒ `CheckSandbox` fail-closed ⇒ **所有派发被拒绝**。这是诚实语义：真实检查未实现前，显式要求 enforcement 的派发一律拒绝，优于静默放行。
- §3 CONTRACTS 7.2 草案「触发模型」句以本节为准落盘：开关名、取值、默认 off、enforce 在 S1 的效果。

### 9.3 (B) 调用点级红测（取代 §6 (iii)/(iv)；§6 (i)/(ii) 不变）

| # | 文件 | 测试函数（新） | 断言要点 |
|---|---|---|---|
| iii′ | internal/console/runtime_agent_sandbox_test.go | TestAgentRunnerRunSandboxRefusesBeforePermit | 经 `ExecuteAgentRunnerRun` **全入口**调用：复用 e2e 装配夹具补齐 Step 0/Step 1 合法输入 + `SandboxEnforced: true` + `testStub` 记录型 Dispatcher；白盒替换 `sandboxForDispatch`（§4.4）为记录型 fake；Denied / Unknown / err 三例 |
| iv′ | internal/console/runtime_supervised_sandbox_test.go | TestSupervisedRunSandboxRefusesBeforePermit | 经 `ExecuteSupervisedRun` **全入口**调用：复用 supervisedHarness（fresh run-root + StubVisibleClient + 假 GateRunner）+ `SandboxEnforced: true`；fake 同上三例 |

iii′ 断言（每例）：
1. `err != nil` 且 `errors.Is(err, guard.ErrSandboxRejected)`（Denied/Unknown）或 `errors.Is(err, guard.ErrSandboxUnavailable)`（err 例）；`summary == RunSummary{}`。
2. fake 恰被调用 1 次且 `spec.Reason == "agent-runner-dispatch"`。
3. **未消耗探针配额**：拒绝返回后再次 `requireProbePermit(ctx, fixture.options.ProbePermit)` 返回 nil（配额可再认领 ⇒ 检查点确在 permit 之前；若已烧则返回配额已消耗错误）。
4. **未产生 spawn / 零副作用**（最省事可测手段）：`os.Stat(runEventLogPath(runRoot))` 为 `ErrNotExist`（Step 2 marker 未建 ⇒ 流程在 Step 2 之前即止，launcher 构造（Step 5）与任何 spawn 结构不可达）；记录型 Dispatcher 0 次调用（下游端口未被触碰）。testStub 同时是变异场景的安全护栏：即便调用点被删，流程也只走 stub 端口，绝不真实 spawn。

iv′ 断言：1 同（reason=="supervised-runtime-dispatch"）；3 同（`requireProbePermit` 再认领成功）；4 `os.Stat(SupervisedInteractionsLedgerPath(runRoot))` 与 `os.Stat(runEventLogPath(runRoot))` 均 `ErrNotExist`（无 ledger、无 marker）；附加：连续两次调用（首次拒绝后仍 freshStart）结果一致、fake 恰 2 次、无残留状态。

fake 替换需保存/恢复 `sandboxForDispatch`，测试不得 `t.Parallel`（白盒全局变量替换）。

**可证伪性自检句（iii′）**：
> 若删除 `runtime_agent.go` 中的检查调用点（`if options.SandboxEnforced { preflightSandbox(...) }` 整块），`TestAgentRunnerRunSandboxRefusesBeforePermit` **必变红**：流程不再于检查处短路，而是继续执行 `requireProbePermit`（烧掉一次性配额）并创建 Step 2 run-root marker，随后经 testStub 走完驱动（返回非 Sandbox 错误或 nil）；于是断言 1（err 非 ErrSandboxRejected/Unavailable）、断言 2（fake 实为 0 次）、断言 3（二次 requireProbePermit 报配额已烧）、断言 4（marker 已存在、Dispatcher 调用 >0）同时失败 ⇒ 红。

**可证伪性自检句（iv′）**：
> 若删除 `runtime_supervised.go` freshStart 块内的检查调用点，`TestSupervisedRunSandboxRefusesBeforePermit` **必变红**：流程越过检查继续 permit → eventLog marker → engine（假 GateRunner 驱动），返回非 Sandbox 错误或 nil，且 ledger/marker 已被写 ⇒ 断言 1/3/4 全失败 ⇒ 红。

变异检验硬交付物（⑥ Section D）：删点 → `go test ./internal/console -run '^TestAgentRunnerRunSandboxRefusesBeforePermit$'` 与 supervised 同名用例 → 红输出原文落 `tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt` / `mutation-supervised.txt` → 还原 → 绿复核。

§6 原 (iii)/(iv)（直测 `preflightSandbox`）降级为可选辅助单测，不作为变异检验载体。

### 9.4 (C) 默认零回归

等价性：`SandboxEnforced` 零值（未设置 / env 未设或 `off` 时 CLI 不注入 / 测试零值构造）时，`ExecuteAgentRunnerRun` 的 `if options.SandboxEnforced` 分支不进入——函数仅多一次布尔加载+比较，无分配、无系统调用、无协程；`ExecuteSupervisedRun` freshStart 分支同理。CLI 层每 `run` 多一次 `os.LookupEnv`，未设 ⇒ false ⇒ 语义与 r1 之前逐句等价；`executeRun` 既有解析/互斥/拒绝顺序（:314-380 判定块）零改动。

④ 阶段证明：`go build ./...`、`go vet ./...`、`go test ./...` 全绿；既有测试文件零改动（全部新增式）；既有装配/多任务/监督测试（`runtime_agent_e2e_test.go` 12 装配用例、`runtime_multitask_red_test.go`、`runtime_supervised_test.go` harness）零值字段即旧路径，天然充当等价性证人；另跑 `go test ./internal/console/... ./internal/guard/... -count=1`；CI 环境无该环境变量 ⇒ 默认 off。

### 9.5 对 r1 其余部分的增量（不改 r1 正文）

- §2.2 骨架：无需微调，③0 按原文落地；§2.5 `CheckSandbox` 不变。
- §附 清单增量（本节补）：②b 新增 `sandboxEnforcementFromEnv()`（sandbox_preflight.go 内）；cli.go `executeRun` 一处读取 + :410/:413 两条路由签名各加 `sandboxEnforced bool` 参数；:579 options 字面量与 supervised_cmd.go :39/:42 旁各一处字段注入；红测 (iii′)/(iv′) 两文件；CONTRACTS §7.2 落盘时含开关语义；`03-tests/mutation-*.txt` 两件（④/⑥ 产出）。

DESIGN-R2: DONE
