# ProofRail 验证报告 — P42-T027B-EVIDENCE-BIND

| 字段 | 内容 |
|---|---|
| 切片 | `P42-T027B-EVIDENCE-BIND` |
| 类型 | 实现类 `[SLICE]` |
| 准则版本 | v1.60 |
| 基线 HEAD | `b355081` |
| 日期 | 2026-10-31 |

## 目标与边界

- **目标**：闭合 C3 与 N2 两个缺口——同一份证据的限长判定、内容、截断标志、摘要与身份必须出自同一个已打开的文件对象，既不得「先按路径 stat、再按路径读」，也不得「打开后再按路径重读」；`RunEvidenceDir` 的目录绑定不得在返回字符串处终止，须延续到实际 I/O。
- **收口形态**：`internal/evidence` 只增不改地新增绑定读原语 `BoundFile`（`OpenBoundFile` / `OpenBoundFileTruncating` / `NewBoundFile` / `ReadAll` / `Hash` / `Identity` / `Close`），身份取自句柄 `Stat()`，`ReadAll` 与 `Hash` 只 `Seek` 后读同一句柄；`internal/adapters` 新增目录绑定原语（`BoundDir` / `OpenBoundDir` / `OpenBoundDirFile`），Windows 以 `CreateFile` 加 `FILE_FLAG_BACKUP_SEMANTICS` 取目录句柄、用 `GetFinalPathNameByHandle` 派生真名相对打开、再以文件句柄最终路径复核，Unix 以 `openat` 加 `O_NOFOLLOW` 与 `O_CLOEXEC` 相对 dirfd 解析。
- **范围**：`internal/adapters/**` 与 `internal/evidence/**`。
- **边界**：不碰契约 / schema / fixtures / `ci.yml`；`internal/chain` / `gates` / `snapshot` / `console` 零改动；公开签名（`RunEvidenceDir` 的字符串返回、三个 seam 函数变量）逐字不变；不推 gitee；`commit` / `push` 停 ⑩ 等授权。

## 审查链结论

- **① 架构**：`deep-reasoner` + DeepSeek V4 Pro；《C3+N2 句柄与目录绑定方案》，`DESIGN: DONE`，含修订注记 D1 / D4 / F1，并声明上文与注记冲突时以注记为准。
- **③ 红测**：`RED-CONFIRMED`（不变式 I1–I5 先钉住同句柄读旧内容、替换后拒绝、祖先目录替换、身份自洽、等价回归，随后进入实现）。
- **② 实现两轮**：主体落成（绑定读原语 + 目录绑定原语 + 读取点替换 + 冻结摘要兜底）→ F 修复轮（Windows 逐读复核 `confirmOpenedArtifact` 落地、注释与实际强度对齐）。
- **④ 门禁**：`go build` / `go vet` / `go test ./...` 全绿；三作用域与 SELFTEST 读数见 §门禁读数。
- **⑤ r1**：`PRE-REVIEW: PASS WITH FIXES`（1 Medium + 3 Low + 2 Info）。
- **⑤ r2**：`PRE-REVIEW: PASS`（F1–F4 逐条核实闭合，并追加 N1–N6 登记要求）。
- **⑥**：首轮输出格式不合规（无 A–E 段）已按契约重发，第二轮 `INDEPENDENT SCAN: PASS`，本片 ⑥ 额度 2/3。
- **⑦**：`FINAL REVIEW: PASS WITH FIXES`（3 Low + 4 Info，均非阻断）。

## 六段扫描转录

来源 `docs/validation/evidence/P42-T027B-EVIDENCE-BIND/06-scan/review.md`；围栏内逐字转录、零改写，范围自 `## Section A` 起、至结论行止。

```sixstage
## Section A — 总体结论
本次只读复核未发现 P42-T027B-EVIDENCE-BIND 的 Medium+ 缺陷，当前实现与设计一致，并已按 C3+N2 句柄绑定/目录绑定语义落地。

## Section B — 发现清单
- 无

## Section C — 逐条证据与推理
在 [internal/evidence/boundread.go](internal/evidence/boundread.go#L23-L152) 中，`OpenBoundFile` 只保留单个已打开的 `*os.File`，随后 `ReadAll`、`Hash`、`Identity` 都从同一句柄读取，不再按路径做第二次解析；这与“打开后再按照路径重读”的 TOCTOU 缺陷相反。 [internal/adapters/bounddir_windows.go](internal/adapters/bounddir_windows.go#L39-L196) 和 [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L18-L99) 分别通过 Windows 的 `GetFinalPathNameByHandle` 复核与 Unix 的 `openat(..., O_NOFOLLOW)` 维持目录绑定，确保祖先目录被替换时读取仍落到原对象。 [internal/adapters/agent_runner_evidence.go](internal/adapters/agent_runner_evidence.go#L100-L201) 中的 `FrozenFacts` 也要求缺口统一返回 `ErrAgentRunnerEvidenceGap`，避免把“未完成”伪装成空摘要；其既定边界在 [internal/adapters/agent_runner_postflight.go](internal/adapters/agent_runner_postflight.go#L554-L613) 和 [internal/adapters/agent_runner_run.go](internal/adapters/agent_runner_run.go#L317-L346) 中被明确标注为本 slice 的 residual boundary，而非未修复漏洞。 [docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md](docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md#L31-L101) 的 D1/D4/F1 修订注记也保持一致：Q2 绑定对象收敛到基目录、释放外部不在本 slice 迁移、Windows 逐读复核已实现。

## Section D — 可证伪性审计
一个直接反例是：在 [internal/evidence/boundread.go](internal/evidence/boundread.go#L83-L152) 中删去 `Seek(0, io.SeekStart)` 和 `io.LimitReader(bound.file, readLimit(bound.bound))` 的同句柄读逻辑，改回 `os.ReadFile(path)`，或者删除 [internal/evidence/boundread_test.go](internal/evidence/boundread_test.go#L19-L92) 里 `TestBoundFileReadsOriginalAfterPathReplacement` 的断言；该变更会让“路径替换后仍读到旧文件”的安全判定失效，结论会从 PASS 翻成 FAIL。另一个反例是去掉 [internal/adapters/evidence_bind_dir_windows_test.go](internal/adapters/evidence_bind_dir_windows_test.go#L31-L88) / [internal/adapters/evidence_bind_dir_unix_test.go](internal/adapters/evidence_bind_dir_unix_test.go#L31-L147) 的祖先替换测试，或删除 [internal/adapters/bounddir_unix.go](internal/adapters/bounddir_unix.go#L43-L66) 的 `O_NOFOLLOW`/`openat` 绑定；那么目录绑定就无法证明“原对象仍被读取”，PASS 也会被推翻。

## Section E — 我无法验证的部分与什么会推翻 PASS
1. Unix 上的 `openat(..., O_NOFOLLOW)` 与 symlink 反制场景需要在 CI/Unix 环境中真实执行；本次只读复核无法替代该端到端验证，因此 [internal/adapters/evidence_bind_dir_unix_test.go](internal/adapters/evidence_bind_dir_unix_test.go#L98-L147) 的 Unix 分支仍是“需 CI 保留”的部分。
2. 本 slice 明确保留的 residual boundary 仍处于后续读路径 sweep 范围，例如 [internal/adapters/agent_runner_run.go](internal/adapters/agent_runner_run.go#L316-L346) 中的 `readAgentRunnerChildPIDs` 以及 [internal/adapters/agent_runner_replay_store.go](internal/adapters/agent_runner_replay_store.go#L220-L270) 的部分读入口；如果后续实现把这些位置也改回 `os.ReadFile`/`os.Stat` 后再读，或如果以上红测在 Windows/Unix 上任一失败，PASS 将被推翻。
3. 只要同句柄绑定、目录绑定和缺口门槛仍保持不变，且 Windows 的全量/聚焦检查结果维持为 0 FAIL，则本结论继续成立。

INDEPENDENT SCAN: PASS
```
- 转录口径：⑥ 首轮输出格式不合规（无 A–E 段）已按契约重发，此处转录第二轮合规输出；归档文件头部如实标注首轮未采纳，额度 2/3。

## F 项闭合

⑦ 结论闭集为 `FINAL REVIEW: PASS WITH FIXES`（3 Low + 4 Info）；下表 FR 行按 ⑦ 报告编号逐条对应到同一批登记项与证据锚点（Low = FR-01..FR-03；Info = FR-04..FR-07）。

| # | 项 | 严重度 | 位置 | 处置 |
|---|---|---|---|---|
| ⑤ r1 F1 | 设计反例乙规定的逐读复核未实现 | Medium | `internal/adapters/bounddir_windows.go:78-100` | 已闭合：`confirmOpenedArtifact` 以文件句柄最终路径对目录句柄 `currentName()` 比较，不符即 `CloseHandle` 并拒绝；回归 I3 PASS |
| ⑤ r1 F2 | `bindEvidenceDir` 注释过度声称 | Low | `internal/adapters/agent_runner_replay_store.go` 绑定段 | 已闭合：注释与实现逐条吻合 |
| ⑤ r1 F3 | 不变式 4 按字面未兑现（`hashLogs` / `fileHash` / `readPIDs` 仍按路径解析） | Low | `internal/adapters/agent_runner_postflight.go:552` / `:637` | 已登记至已知限制：三处边界注释如实标注 P42 boundary；迁移归后续读路径 sweep 片（见 §后续片候选登记） |
| ⑤ r1 F4 | D3 注释高估验证强度 | Low | `internal/adapters/agent_runner_replay_root.go:381-399` | 已闭合：改写为「保留一次 `Lstat` 类型门、残余窗口极小」 |
| ⑤ r1 F5 | unix 红测本机未跑、`CloseBoundDirs` 未接 run 收尾 | Info | `internal/adapters/evidence_bind_dir_unix_test.go` 与 store 生命周期 | 已登记至已知限制：unix 分支由 CI 腿兜底；生产接线归后续片 |
| ⑤ r1 F6 | 双工具链差异 | Info | `go.mod`（`go 1.22`）与 ⑤ r1 复现条件 | 已闭合：`go1.22.12` 下双绿复现，差异消除 |
| ⑤ r2 N1 | `CloseBoundDirs` 无生产调用方，注释仍称 per-request 句柄 | Low | `internal/adapters/agent_runner_replay_store.go:220-224` | 已登记至已知限制：句柄存活期以 store / 进程为上限；生产接线与措辞统一归后续片 |
| ⑤ r2 N2 | Windows 复核为前缀比较，reparse 解析边界未显式登记 | Low | `internal/adapters/bounddir_windows.go:127-162` | 已登记至已知限制：对象身份级校验归后续片 |
| ⑤ r2 N3 | Unix `runs/<requestId>` 绝对目标符号链接向量 | Low | `internal/adapters/bounddir_unix.go:43-66` | 已登记至已知限制：两段式 dirfd 归后续片 |
| ⑤ r2 N4 | `agentRunnerFileHash` 边界注释措辞 | Info | `internal/adapters/agent_runner_postflight.go:643` | 已登记至已知限制：措辞层，无语义影响 |
| ⑤ r2 N5 | `tmp` 仅存档全量一腿 | Info | `tmp/P42-T027B-EVIDENCE-BIND/full1.txt` | 已登记至已知限制：第二腿由 CI 双腿承担，`tmp` 于 ⑧c 前清除 |
| ⑤ r2 N6 | `OpenBoundFileTruncating` 注释措辞 | Info | `internal/evidence/boundread.go:71-76` | 已登记至已知限制：措辞层，无语义影响 |
| ⑦ FR-01 | `CloseBoundDirs` 无生产调用方且注释口径与 D1 冲突 | Low | `internal/adapters/agent_runner_replay_store.go:220-224` | 已登记至已知限制，并列为后续片候选（生产接线 + 注释统一） |
| ⑦ FR-02 | Windows 句柄级复核为路径前缀比较，reparse 边界未证明 | Low | `internal/adapters/bounddir_windows.go:127-162` | 已登记至已知限制，并列为后续片候选（对象身份级校验） |
| ⑦ FR-03 | Unix 中间组件与绝对目标链接向量不在覆盖内 | Low | `internal/adapters/bounddir_unix.go:43-66` | 已登记至已知限制，并列为后续片候选（两段式 dirfd） |
| ⑦ FR-04 | 冻结摘要措辞泛化（frozen manifest 与 `Facts` 摘要字段） | Info | `internal/adapters/agent_runner_postflight.go:643` | 已登记至已知限制：措辞层，比对行为不变 |
| ⑦ FR-05 | 全量测试第二腿未入 `tmp` 存档 | Info | `tmp/P42-T027B-EVIDENCE-BIND/full1.txt` | 已登记至已知限制：以 CI 双腿为权威读数 |
| ⑦ FR-06 | 截断类开函数注释措辞 | Info | `internal/evidence/boundread.go:71-76` | 已登记至已知限制：措辞层，契约不变 |
| ⑦ FR-07 | darwin 交叉编译失效、旧验证文档条目失效 | Info | `internal/adapters/bounddir_unix.go:61` 与 `docs/validation/t027-replay-convergence.md:38` | 已登记至已知限制：本片仅登记「口径替代」，旧文档修订归后续文档片 |

## 偏差与登记

| 编号 | 内容 | 来源 | 判定 |
|---|---|---|---|
| D1 | 绑定对象由 per-request 证据目录收敛为 `runs` 基目录 | ② 实现期实测 + ① 修订注记 | 安全：Windows 上子目录存在打开句柄时，重命名其祖先返回 `ERROR_ACCESS_DENIED`（即便以 `FILE_SHARE_DELETE` 打开），per-request 句柄会把祖先替换攻击变成拒绝服务；绑定基目录则祖先可被改名而句柄继续指向原对象树。读取解析定为 `<requestId>/<rel>` |
| D2 | `Close` / `NewBoundFile` / `OpenBoundFileTruncating` 三分 | ② 实现期 | 安全：拒绝类与截断类契约必须可区分，记录 / PID 类打开即拒，日志 / 事件类保留 truncated 降级 |
| D3 | `loadAgentRunnerReplayOwnership` 保留一次 `Lstat` 类型门 | ② 实现期 | 边界可接受：保住 symlink 叶拒绝语义；残余窗口极小，注释已按实际强度改写 |
| D4 | `BuildReleaseInventory` 不迁移 | ① 修订注记 | 设计清单项收窄：adapters 到 evidence 的依赖方向不允许反向复用；release 树由操作员选定，不属候选可写面，归后续读路径 sweep 片 |
| D5 | 上界 64 MiB 与 1 MiB | ② 实现期 | 安全：`agentRunnerWorkspaceEntryBound` 为 64 MiB，`agentRunnerMetadataBound` / `replayStoreRecordBound` / `fileQueueEnvelopeBound` 为 1 MiB，fail-closed 方向 |
| D6 | `CloseBoundDirs` 无生产调用方 | ⑤ r2 N1 | 已登记：句柄存活期以 store / 进程为上限；注释仍称 per-request，与 D1 的基目录口径矛盾，随生产接线一并修正 |
| D7 | darwin 交叉编译失效 | ⑤ r1 实测 | 口径替代：`GOOS=darwin` 失败于 `bounddir_unix.go:61` 的 `undefined: syscall.Openat`，`GOOS=linux` 通过；支持矩阵仅 windows / linux，`docs/validation/t027-replay-convergence.md:38` 旧「Darwin/amd64 交叉编译通过」条目失效，旧文档修订归后续文档片 |
| D8 | `allowTruncation` 由调用方按类判定 | ② 实现期 | 安全：契约由调用点显式选择，原语不猜类别 |

- **残余面登记**：`runs/<requestId>` 同层替换属兄弟节点，不在基目录句柄覆盖内，由唯一生产调用方以冻结摘要比对 fail-closed，替换退化为失败而非可信证据；Windows「派生名到 `CreateFile`」之间的再替换窗口由 `confirmOpenedArtifact` 的句柄级复核封住；Unix 中间组件与绝对目标链接向量归后续两段式 dirfd 片。

## 模型与计量

| 环节 | 载体 | 模型 |
|---|---|---|
| ① 架构 | `deep-reasoner` | DeepSeek V4 Pro |
| ② 实现 | `prfrail-implementer` | DeepSeek V4.1 Flash |
| ③ 测试 | `prfrail-tester` | DeepSeek V4.1 Flash |
| ⑤ 预审 | `deep-reasoner` | DeepSeek V4 Pro |
| ⑥ 扫描 | `independent-reviewer` | MAI-Code-1.1-Flash |
| ⑦ 终审 | `independent-reviewer` | GPT-5.3-Codex |
| ⑧ 文档 | `prfrail-documenter` | DeepSeek V4.1 Flash |
| 批量辅助 | `deep-reasoner` | GLM-5.3-Flash（2 次：读路径普查 / 反例核查） |

- **主控模型事实（用户 2026-10-31 裁定，如实登记）**：P37 至 P41 期间主控被意外设为 `DeepSeek V4 Pro`，排查效率高但 API 费用偏高；自 P42 起主控切回 `DeepSeek V4.1 Flash`。
- **门禁可感知性**：准则未对主控模型设硬门禁，属操作员环境设定，终端门禁不可感知；①⑤ 按白名单仍由 V4 Pro 承担，与主控选型是两件事。
- **长期策略**：默认 Flash（High 档）；遇底层并发 / 状态机 / 句柄类复杂片、单份超 10,000 字符的设计文档、大量 CI 日志的深度根因分析、或 Flash 连续 2 次行为漂移，主控应提示切换 Pro 或获授权后切换；Max 档不作默认。

## 门禁读数

- tree 作用域：`TOTAL_FAIL=0`（精确暂存后复跑，worktree 与 index 同源）。
- index 作用域：`TOTAL_FAIL=0`（changed 34）。
- 自检：`SELFTEST: PASS 697/697`。
- G6-1 状态时态扫描行数：418；无待办与完成冲突。
- ci 作用域由 CI 双腿承担，读数见 §提交与 CI。

## 提交与 CI

- 主体提交：`cab3e08` `feat: bind evidence reads to one handle and the evidence directory (P42-T027B-EVIDENCE-BIND)`；推送 origin `b355081..cab3e08`。
- CI 运行 `41194972405` 首轮即双绿：`Go ubuntu-latest` success、`Go windows-latest` success。
- ubuntu 腿硬证据：`ok … internal/adapters 12.230s`、`ok … internal/evidence 0.015s`、Race 腿 `ok … internal/adapters 89.088s`、`SELFTEST: PASS 697/697`、`GATE REPORT scope=ci changed=34`、`TOTAL_FAIL=0`。
- windows 腿硬证据：`ok … internal/adapters 31.296s`、`ok … internal/evidence 0.181s`、`SELFTEST: PASS 697/697`、`GATE REPORT scope=ci changed=34`、`TOTAL_FAIL=0`。
- 本报告回填提交属文档增量，其 CI 属观察运行，不计切片轮次。

## 后续片候选登记

- 读路径 sweep：`hashLogs` / `readPIDs` / `fileHash` / `scanWorkspace` 的 N2 绑定迁移。
- Unix 两段式 dirfd：以 `openat` 逐级相对打开，覆盖中间组件符号链接向量。
- Windows reparse 对象身份级校验：替代当前的前缀比较。
- `CloseBoundDirs` 生产接线：run 收尾调用并统一注释口径。
- darwin 旧文档 supersede：`docs/validation/t027-replay-convergence.md` 的交叉编译条目修订。
- 原语层测试补齐：`BoundFile` 关闭后拒读、`Hash` 超界拒摘要、身份同源对抗用例（批量辅助登记）。

## 已知限制

- **⑥ Section E 要旨**：Unix 侧 `openat` 与 `O_NOFOLLOW` 的 symlink 反制须在 CI/Unix 环境实跑，本机只读复核不能替代；`readAgentRunnerChildPIDs` 与回放存储的部分读入口仍属本片 residual boundary；任一同句柄绑定、目录绑定或缺口门槛退回按路径重读，或红测在任一平台失败，结论即翻转。
- **N1 · 状态 = 接受（后续片接线）**：`CloseBoundDirs` 无生产调用方，注释仍称 per-request 句柄；锚点 `internal/adapters/agent_runner_replay_store.go:220-224`。
- **N2 · 状态 = 接受（后续片）**：Windows 复核为句柄最终路径与目录句柄当前真名的前缀比较，reparse 解析边界与「仅前缀」限制未做对象身份级证明；锚点 `internal/adapters/bounddir_windows.go:127-162`。
- **N3 · 状态 = 接受（后续片）**：`openat` 保证祖先不被重定向，`runs/<requestId>` 的绝对目标符号链接与中间组件向量不在覆盖内；锚点 `internal/adapters/bounddir_unix.go:43-66`。
- **N4 · 状态 = 接受（措辞层）**：`agentRunnerFileHash` 边界注释以 frozen manifest 泛指冻结摘要，实际比对锚点为 `Facts` 的清单摘要字段；锚点 `internal/adapters/agent_runner_postflight.go:643`。
- **N5 · 状态 = 接受**：`tmp/P42-T027B-EVIDENCE-BIND/` 只存档全量一腿，第二腿由 CI 双腿承担，`tmp` 于 ⑧c 前清除。
- **N6 · 状态 = 接受（措辞层）**：截断类开函数的注释措辞，契约不变；锚点 `internal/evidence/boundread.go:71-76`。
- **FR-01 · 状态 = 接受 + 后续片**：与 N1 同源，生产接线与注释统一同片处理。
- **FR-02 · 状态 = 接受 + 后续片**：与 N2 同源，对象身份级校验列为后续片。
- **FR-03 · 状态 = 接受 + 后续片**：与 N3 同源，两段式 dirfd 列为后续片。
- **FR-04 · 状态 = 接受（措辞层）**：与 N4 同源，比对行为不变。
- **FR-05 · 状态 = 接受**：与 N5 同源，以 CI 双腿为权威读数。
- **FR-06 · 状态 = 接受（措辞层）**：与 N6 同源。
- **批量辅助发现 · 状态 = 接受 + 后续片**：原语层三处测试缺口（关闭后拒读 / 超界拒摘要 / 身份同源对抗），主链经单次读即关闭与 `ReadAll` 路径规避，新调用方需自证。
- **FR-07 · 状态 = 接受（口径替代）**：支持矩阵仅 windows / linux；`docs/validation/t027-replay-convergence.md:38` 旧条目失效的修订归后续文档片。

```artifacts
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/p42-evidence-bind.md	repo
docs/validation/p42-evidence-bind_EN.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/MANIFEST.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/SHA256SUMS.txt	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/05-prereview/review.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/06-scan/review.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/07-final/review.md	repo
docs/validation/evidence/P42-T027B-EVIDENCE-BIND/freedom-list.md	repo
internal/adapters/agent_runner_evidence.go	repo
internal/adapters/agent_runner_pinned_cli.go	repo
internal/adapters/agent_runner_postflight.go	repo
internal/adapters/agent_runner_replay_root.go	repo
internal/adapters/agent_runner_replay_store.go	repo
internal/adapters/agent_runner_replay_store_launches.go	repo
internal/adapters/agent_runner_replay_store_terminals.go	repo
internal/adapters/agent_runner_run.go	repo
internal/adapters/ai_probe_copilot_executor.go	repo
internal/adapters/filequeue.go	repo
internal/adapters/bounddir.go	repo
internal/adapters/bounddir_unix.go	repo
internal/adapters/bounddir_windows.go	repo
internal/adapters/evidence_bind_dir_unix_test.go	repo
internal/adapters/evidence_bind_dir_windows_test.go	repo
internal/adapters/evidence_bind_semantics_test.go	repo
internal/adapters/evidence_seam_signatures_test.go	repo
internal/evidence/boundread.go	repo
internal/evidence/boundopen_unix.go	repo
internal/evidence/boundopen_windows.go	repo
internal/evidence/boundread_test.go	repo
internal/evidence/disposition.go	repo
tools/gates/cjk-newwords.txt	repo
```
