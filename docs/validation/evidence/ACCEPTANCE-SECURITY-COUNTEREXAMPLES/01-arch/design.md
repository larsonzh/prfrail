# ① 架构设计原文（逐字归档）

载体：runSubagent(agentName="deep-reasoner", model="DeepSeek V4 Pro (deepseek)")。只读设计，未修改任何文件。

DESIGN: DONE

# P52 安全反例切片 · ① 架构设计（只读方案，零生产代码改动）

**前置事实核查（与 GLM 输入不符处如实指出）**：GAP-70 所称两条缺口已被覆盖——`TestReadinessAcceptanceRejectsInvalidSealBinding`（readiness_test.go L280，3 子例）与 `TestReadinessProductionRejectsMalformedAT23`（L328，5 个**非空**非法 AT23 子例）。本片只补残余个案（§3-F）。GAP-51 真正活缺口是**中间组件符号链接**（→ V-1）。

## 1 威胁模型（逐 GAP）
- **GAP-51（Unix dirfd）**：攻击者=对证据目录子树有写权限者（模型/工具输出即敌意输入）。控制面：目录内任意文件与符号链接。预期 fail-closed：绑定读永不解析到绑定对象之外（bounddir_clean.go L22 注释契约）。
- **GAP-52（Windows reparse）**：同权限模型，攻击面=reparse point（junction）。预期：`confirmOpenedArtifact`（bounddir_windows.go L118）以句柄最终名拒绝界外解析。
- **GAP-55（原语层）**：攻击者=误用句柄的调用方（关闭后读、增长越界）。预期：三哨兵（boundread.go L31/33/35）精确报出。
- **GAP-60（K 绑定）**：攻击者=能改写磁盘 seal 文件者。预期：acceptance 回放路径拒绝他方 K（engine.go L352-353）。
- **GAP-64（ProbePermit）**：攻击者=双进程同时探测（配额翻倍）。预期：磁盘 O_EXCL marker 为线性化点，第二进程得 `ErrProbePermitConsumed`（probe_permit.go L235-238）。
- **GAP-70（readiness）**：攻击者=伪造 seal/AT23 输入。预期：四哨兵（readiness.go L53/56/59/62）拒绝。

## 2 攻击载荷（逐步 + 未阻断后果）
**GAP-51-B 中间符号链接逃逸（红→上报 V-1）**：① mkdir D，写 D/original；外目录 O 写 secret；② `OpenBoundDir(D)`；③ mkdir D/sub；`ln -s O D/link`；④ `OpenBoundDirFile(dir,"link/secret")`。现码路径：`cleanBoundRelativePath` 只查拼写（L20-33）→ `openat(fd,"link/secret",O_NOFOLLOW)`（bounddir_unix.go L66）——内核仅对**末组件**拒绝 symlink，中间组件照跟→fd 指向 O/secret；`NewBoundFile` 只拒目录/非正则/超界（boundread.go L103-111）→ 返回界外句柄。**后果**：任意可读常规文件（含他人证据、宿主敏感文件）字节进入 ReadAll/Hash 被当作本绑定证据，绑定隔离承诺整体失效。
**GAP-51-A 祖先替换、GAP-52-B 祖先改名竞态**：既有一致已绿（evidence_bind_dir_{unix,windows}_test.go），不重复。
**GAP-52-A junction 逃逸（应绿）**：① 绑 D；② 外目录 O 写 secret；③ `cmd /c mklink /J D\j O`；④ `OpenBoundDirFile(dir,"j\secret")`。CreateFile 跟随 reparse，`GetFinalPathNameByHandle` 返回 O\secret，前缀检查拒绝（L128）并关句柄。**后果（若未阻断）**：Windows 半部唯一防护失效，界外字节成为证据。
**GAP-60 他方 K 回放**：见 T60-1。
**GAP-64 真跨进程双写**：见 T64-6。

## 3 测试清单与断言矩阵（新文件 7 个）
A. `internal/evidence/boundread_counterexample_test.go`（无 tag，package evidence）
- **T55-1** TestBoundFileReadAfterCloseRefused：open→ReadAll→Close→ReadAll 与 Hash 均 `errors.Is(err, ErrBoundFileClosed)`（L183-188）。
- **T55-2** TestBoundFileHashRefusesGrowthPastBound：open(bound=4)→第二写者增至 5 字节→ReadAll：truncated=true 且 len==4；Hash→`errors.Is(err, ErrBoundArtifactExceedsBound)`（L157）。Windows 若共享模式阻第二写者则 skip（预留）。
- **T55-3** TestBoundFileNewRefusesNilAndNegativeBound：NewBoundFile(nil,0,false)→`ErrBoundFileNotRegular`（L90）；NewBoundFile(f,-1,false)→`ErrBoundArtifactExceedsBound`（L94）。
- **T55-4** TestBoundFileOpenRefusesDirectoryWithSentinel：OpenBoundFile(目录)→`ErrBoundFileNotRegular`（补 L127 既有测试缺哨兵断言）。
- **T55-5** TestBoundFileOpenRefusesOverBoundWithSentinel：bound=8、9 字节→`ErrBoundArtifactExceedsBound`（补 L89）。
B. `internal/adapters/bounddir_identity_test.go`（`//go:build !darwin`）
- **T55-6** TestBoundDirVerifyRefusesSwappedDirectory：OpenBoundDir(P)→rename P→P.bak→mkdir P→`boundDirVerify(dir,P)`：err!=nil 且文案含 "is not"（unix L107-110 / windows L168+sameObjectAs L240-242）。
- **T55-7** TestCleanBoundRelativePathRefusesEscapeShapes：`""`/`"/etc/passwd"`/`".."`/`"a/../../b"` 四拒绝（L24/27/31）；控制 `"sub/./x"`→`"sub/x"`。
C. `internal/adapters/bounddir_escape_unix_test.go`（`//go:build !windows && !darwin`）
- **T51-1** TestBoundDirMiddleSymlinkEscapesBound（§2 GAP-51-B 载荷）：断言方向=「攻击不得得手」。本片期望**红**；② 在 Linux 环境采集红输出（含读到的字节与路径）归档 V-1 证据，文件以现状记录注释保留，不修生产码。
D. `internal/adapters/bounddir_escape_windows_test.go`（`//go:build windows`）
- **T52-1** TestBoundDirRefusesJunctionEscape（§2 GAP-52-A）：err!=nil 且含 "outside the bound directory"（L128）；mklink 失败即 skip（先例 agent_runner_replay_root_windows_test.go L38-41）。
- **T52-2** TestBoundDirJunctionInsideBoundAdmitted：j→D\sub（界内）放行且字节正确（防过度收紧）。
E. `internal/gates/probe_permit_counterexample_test.go`（无 tag）
- **T64-1** Require(nil)→`errors.Is(err, ErrProbePermitMissing)`（L214-216）；**T64-2** (*ProbePermit)(nil).Require(ctx)→同哨兵（L211-213）。
- **T64-3** TestProbePermitCompleteUnknownState：白盒 `p.state.Store(3)`→Complete→`ErrProbePermitConsumed`（L266）。
- **T64-4** TestProbePermitLoadNonNotExistReadFailure：marker 路径为目录→err!=nil 且 NOT errors.Is(三个 permit 哨兵)（契约 L49-53：读失败无哨兵）。
- **T64-5** LoadProbePermit("")→`ErrProbePermitMissing`（L142）。
- **T64-6** TestProbePermitCrossProcessSingleClaim（真跨进程，先例 executor_test.go L83/L109 子进程模式）：父 NewProbePermit(path)；子=exec.Command(os.Executable(),"-test.run=TestProbePermitCrossProcessHelper") env 带 path：子 Require 成功→写 claimed 信号→等 release 信号→Complete→退出。父见信号后 Require→`ErrProbePermitConsumed`（O_EXCL 输给子，L235-238）；父发 release 等子退出；父**新**实例再 Require→`ErrProbePermitConsumed`（进程死亡配额不复活）。
F. `internal/gates/readiness_counterexample_test.go`（无 tag）
- **T70-1** TestReadinessProductionRejectsShaPrefixedMalformedAT23："sha256:"+63hex / 65hex / 63hex+"g" / "SHA256:"+64hex → `ErrReadinessEvidenceIncomplete`（L209）。
- **T70-2** TestReadinessProductionValidatesSealBeforeAT23：SealBinding{Attempt:0,...}+合法 AT23→`ErrReadinessEvidenceIncomplete`（L202），readiness==nil。
- **T70-3** TestReadinessAcceptanceRejectsPaddedSealHash：SealHash=64hex+" "→`ErrReadinessEvidenceIncomplete`（validate 不 trim，L229-230 注释）。
- **T70-4** TestReadinessAT23TrimsSurroundingWhitespace（行为钉）：AT23="  "+64hex+"  "→放行且 Release().AT23Hash==裸 64hex（TrimSpace L205），佐证 V-3。
G. `internal/chain/state_seal_counterexample_test.go`（无 tag，复用 state_seal_test.go 的 fileStore/sealTestOptions）
- **T60-1** TestAcceptanceReplayRefusesForeignContextKey：store.SealTaskEvidence 写入与引擎同构但 ContextKey=他方 64hex、ReviewDecision="approve"、PromotionDecision="completed" 的 seal；直接调 `engine.acceptFromExistingSeal(ctx, taskEntityAt("run-one","task-one",1), "task-one", 1, parent, workspace)`（parent.Hash==record.ParentHash、workspace 取 testOptions 返回实例，仅 K 差）→ ok==false 且 TaskStates[taskKey] != "PASSED"（REPAIR_PENDING，reason "evidence-seal-mismatch"，engine.go L352-353）。
- **T60-2** TestSealValidateRefusesNonHashContextKey：ContextKey="not-a-hash"→`errors.Is(err, ErrSealIncomplete)`（seal.go L81-82）。
- **T60-3** TestSealValidateAdmitsEmptyContextKey：ContextKey=""→validate()==nil，钉现状佐证 V-2。

## 4 红测规格（可证伪性）
| 测试 | 移除防御（生产码一处） | 红形状 |
|---|---|---|
| T55-1 | checkReadable 删 closed 判断 | 非哨兵错误/panic |
| T55-2 | Hash 删 L157 | 返回摘要而非哨兵 |
| T55-3/4/5 | NewBoundFile 删 L90/L94/L103/L111 | 返回 nil 或放行 |
| T55-6 | boundDirVerify 删身份比较 | 换目录后返回 nil |
| T55-7 | clean 删任一检查 | 对应 case 放行 |
| T51-1 | （现码已红，本片不修） | — |
| T52-1 | 删 confirmOpenedArtifact 调用 | 界外字节可读 |
| T60-1 | 删 L353 K 比较 | 他方 K 回放→PASSED |
| T60-2 | validate 删 L81-82 | 非哈希 K 通过 |
| T64-1/2/3 | 删 nil/default 分支 | panic 或 nil |
| T64-6 | 线性化点改内存态 | 双进程双赢 |
| T70-1/2/3 | 删 L209/L202/validate 检查 | 非法输入放行 |
本片判定：除 T51-1 外全绿；T51-1 红=真实漏洞证据，归档不修。

## 5 疑似漏洞上报清单（本片不修，报 GAP）
- **V-1（HIGH，GAP-51）**：bounddir_unix.go L66 openat 的 O_NOFOLLOW 仅护末组件，中间符号链接逃逸绑定（§2 载荷稳定复现）。
- **V-2（LOW，GAP-60）**：seal.go validate（L81-82）允许空 ContextKey；两消费端（engine.go L352-353、seal_verify.go L44-48）各自 fail-closed，但 validate 契约应收紧为「K 必非空且为哈希」。router/handoff（router.go L79-129、handoff.go）不消费 seal，属设计分工而非漏洞。
- **V-3（LOW，GAP-70）**：归一化不对称——AT23Hash TrimSpace（readiness.go L205）而 SealHash 拒绝空白（validate 注释 L229-230），建议统一拒绝空白（需改生产码，另开切片）。
- **V-4（观察级，GAP-52）**：isUnderDirectory（L136）拼写前缀而非对象身份，设计注释自证同 API 同形式；未能构造确定性绕过（8.3/大小写/`..` 均被拒），仅登记残余风险。

## 6 ③ 阶段拆分（纪律 V，串行）
- **P52-S1 原语层反例**：文件 A+B（T55-1..7）。验证：`go test ./internal/evidence/`、`go test ./internal/adapters/`；`go vet ./...`。
- **P52-S2 平台分流攻击**：文件 C+D（T51-1、T52-1/2）。T52 本机 Windows 跑；T51-1 需 Linux（远程 VM/CI，先例 prfrail-remote-validation），采集红输出并归档 V-1；`GOOS=linux go build ./...` 交叉编译自检。
- **P52-S3 门禁反例**：文件 E+F（T64-1..6、T70-1..4）。`go test ./internal/gates/`。
- **P52-S4 K 绑定反例**：文件 G（T60-1..3）。`go test ./internal/chain/`；收尾汇总 V-1..4 证据清单交 ②/⑧。
红线复述：四步均不动 `internal/**` 生产代码；T51-1 红输出按「疑似漏洞上报」停等授权。

DESIGN: END
