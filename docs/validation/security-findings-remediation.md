# P83 安全发现修复验证报告（SECURITY-FINDINGS-REMEDIATION）

生成日期：2026-11-21。切片：P83。分级：`[SLICE]`，实现类（安全修复），规模 M-L。

**结论：`PARTIALLY CLOSED`**（部分闭合）——F2 / F5 / F7 与 T9 收窄已闭合；F1 / F4 部分实现、残余顺延 P84；F3 / F6 / F8 顺延 P84。**不得读作「闭环」或「安全闭合」**。

---

## 1. 切片标识与授权

| 项 | 内容 |
|---|---|
| 标识 | `SECURITY-FINDINGS-REMEDIATION`（承接 P79 报告 §9 的 F1–F8） |
| 授权 | 用户 2026-11-20 点名起跑；同日两次追加裁定（分流批准、追加复审 1 轮）；2026-11-21 终裁 |
| 依据 | P79 报告 `docs/validation/s2-pre-independent-security-review{,_EN}.md` §9；v66 快照 §四 队列项 |
| 边界 | 只修 F1–F8；不改 S1 Exit 八门判据、不改威胁模型基线、不升准则版本、不做 gitee |

## 2. 步骤 0 锚点实测

| # | 锚点 | 期望 | 实得 | 结论 |
|---|---|---|---|---|
| 1 | HEAD | `7a63962` | `7a63962` | 相符 |
| 2 | origin/main | `7a63962` | `7a63962` | 相符 |
| 3 | 准则版本 | v1.71 | v1.71 | 相符 |
| 4 | SELFTEST | 740/740 | 740/740 | 相符 |
| 5 | 三作用域 | 全 0 | index / tree / ci 全 0 | 相符 |
| 6 | 快照 v64 | 见 v66 §一 | 见 v66 §一（按 LF 计数） | 相符 |
| 7 | 快照 v65 | 见 v66 §一 | 见 v66 §一（按 LF 计数） | 相符 |
| 8 | 快照 v66 | 由 v67 复核 | 见 v66 §一（记录） | 记录 |

## 3. 目标与范围

修复 P79 登记的八条安全发现，优先级 F2 / F3 → F6 → F1 / F4 / F5 / F7 / F8；每 F 一条可证伪闭合判据与至少一条反例回归测试加变异检验。

## 4. 修复实现与闭合状态

| 编号 | 面 | 级别 | 主要改动位置 | 闭合状态 |
|---|---|---|---|---|
| F1 | 面 1 | Medium | `internal/gates/runner.go`、`internal/evidence/bounddir_*.go` | **部分实现**：执行前绑定根目录、无产物不打开根、非阻塞打开加类型校验、FIFO 拒绝；T4 残余（大小写折叠、中间 reparse、根 `""` 语义）顺延 P84 |
| F2 | 面 2 | High | `internal/console/runtime_agent.go`、`internal/adapters/agent_runner_replay_dispatcher.go`、`internal/adapters/agent_runner_admission.go`、`internal/console/trusted_ledger.go` | **已闭合**：账本路径由持久 root 派生、生产不回退内存副本、读取后取时钟、复核为启动前最后一次外部观察、失败保留 R 且零启动 |
| F3 | 面 2 | High | `internal/chain/review.go`、`internal/console/runtime_agent.go`、`internal/console/runtime_supervised.go` | **顺延 P84**：只认最新 grant、风险确认与策略批准分离、发布前重读均落地；**T2 STOP**——AgentRunner 无可持久复验的冻结策略制品，`grant.PolicyHash` 无独立对照物 |
| F4 | 面 2 | Medium | `internal/console/runtime_agent.go`、`internal/console/runtime_supervised.go`、`internal/guard/lease.go` | **部分实现**：进程级持有者校验、锁内接管加代次推进、逐代租约留档、证明按 run 归属校验；T5#4 残余（批量回写缺前驱验证、跨进程中断恢复）顺延 P84 |
| F5 | 面 2 | Medium | `internal/gates/runner.go`、`internal/gates/adapter.go` | **已闭合**：分类只取冻结 hook、准入无条件执行、定义摘要与恢复保证/外部系统冲突即拒、适配器不再伪造占位摘要 |
| F6 | 面 3 | High | `internal/console/ai.go`、`internal/adapters/ai_probe_copilot_executor.go`、`internal/adapters/ai_preflight.go` | **顺延 P84**：审批材料含 profile 摘要与完整绑定、审批源由链配置目录派生、凭据解析前校验、批准摘要与文件身份随启动传递；**T3 残余**——Windows 名称到启动映像的绑定窗口、Linux 同 inode 不等于字节不可变 |
| F7 | 面 3 | Medium | `internal/guard/process_windows.go`、`internal/console/secret_input_windows.go` | **已闭合**：终止符余量、逐条 scratch 与弃用数组清零、终端输入增长副本清零、局部失败块清零，均有观测缝与反例 |
| F8 | 面 3 | Medium | `internal/adapters/ai_probe_copilot_executor.go`、`internal/adapters/ai_probe_copilot.go` | **顺延 P84**：关闭 `--log-dir` 落盘、拒绝含日志目的地的参数、拒绝证据只用闭合分类；**T8a 未实施**——无法证明子进程禁写未脱敏日志时拒发 BYOK |
| T9 | 范围 | Medium | `internal/chain/review.go`、`internal/console/ai.go` | **已闭合**：证据根绑定要求限于 `policy` 模式；无 `baseUrl` 的 host-auth profile 不再被接收端门禁误拒 |

## 5. 契约影响判定

八项修复与 T9 收窄均判 **DOES-NOT-TOUCH**：被违反句（`CONTRACTS` :147 / :150 / :212 / :214 / :216 / :391，`SECURITY` :23 / :60）已包含所要求语义。本片**未改** `docs/CONTRACTS{,_EN}.md`、schemas 或夹具。

## 6. 分域升级与用户裁定

| 事项 | 裁定 | 留痕 |
|---|---|---|
| ⑦ 额度（1 终审 + 3 复审） | 用户 2026-11-20 裁定 | 超出 §7.5「M-L 视同 M」字面，系用户直接裁定、不得作为先例 |
| 盲审 | 必抽 1 次（F3 / F6 触身份语义） | 另计 +1，不占终审与复审额度，终审前执行 |
| T8 半句 | 批准「拒绝证据只用稳定分类」；「无法证明禁写即拒发 BYOK」不批准、顺延 | 与 T8a 分开登记 |
| 追加复审 1 轮 | 用户 2026-11-20 裁定（复审 3/3 用尽后） | 超出 §7.5，系用户直接裁定、不得作为先例 |
| T4 / T5#4 | 用户 2026-11-21 裁定：不修、顺延 P84 | F1 / F4 不得记为闭合 |

## 7. 审查链

| 环节 | 载体 | 轮次 | 结论 |
|---|---|---|---|
| ⑤ 预审 | `DeepSeek V4 Pro (deepseek)` | 2 / 2（用尽） | r1 FINDINGS（1 Medium）→ r2 PASS |
| ⑥ 独立扫描 | `MAI Code 2 Flash (copilot)` | 2 / 3 | r1 FINDINGS（4 阻断）→ r2 PASS |
| ⑦ 盲审 | `GPT-6.2 Sol (copilot)` Max | 1 | `BLIND REVIEW: FINDINGS` |
| ⑦ 终审 | `GPT-6.2 Sol (copilot)` Max | 1 / 1 | `FINAL REVIEW: FINDINGS`（9 项） |
| ⑦ 复审 | `GPT-6.2 Sol (copilot)` Max | 3 / 3（用尽） | r1 / r2 / r3 均 FINDINGS |
| ⑦ 追加复审 | `GPT-6.2 Sol (copilot)` Max | 1（用户裁定） | `ADDITIONAL RE-REVIEW: FINDINGS`（范围内 2 项 Medium） |

逐字原文：终审见 `tmp/P83-SECURITY-FINDINGS-REMEDIATION/07-final-review.md`；三轮复审见 `08-reviews-r1-r3.md`；追加轮见 `09-additional-review.md`（均为临时件，按 §11.2 不入库）。

## 8. 顺延 P84 候选清单

| 编号 | 顺延内容 | 来源 |
|---|---|---|
| P84-a | 策略制品（AgentRunner 可持久复验的冻结策略，绑定 run / task / step） | T2 STOP（F3） |
| P84-b | 不可变安装来源（Windows 稳定路径与 ACL、Linux 不可写安装目录） | T3 残余（F6） |
| P84-c | BYOK 证明载体（证明子进程禁写未脱敏日志；不能证明即拒发） | T8a（F8） |
| P84-d | bounddir 大小写折叠与中间 reparse 残余、根相对路径语义 | T4 残余（F1） |
| P84-e | proof 恢复持久化（批量回写前驱验证、代次到 proof 的可恢复事务） | T5#4 残余（F4） |

五项均为「F1–F8 修复义务的顺延部分」；P84 起跑时由用户点名并裁定优先级。

## 9. 文档同步

| 对象 | 处置 |
|---|---|
| 本报告对 | 新增 `docs/validation/security-findings-remediation{,_EN}.md`（CN / EN 严格同位） |
| 台账 | `docs/t027/REMAINING_SLICES{,_EN}.md` P83 段与完成记录（逐 F 状态、R / D、额度披露、顺延清单） |
| SECURITY 注记 | `docs/SECURITY{,_EN}.md` 第 5 行如实写修复状态（F2/F5/F7 修复；F1/F4 部分实现；F3/F6/F8 顺延） |
| ADR-047 | 保持 `COMPENSATED` 不动 |
| 未改 | `CHECKLIST_v66{,_EN}`、`s1-exit`、`p57` 门表、准则本体、CI、`DEV_PLAN` |

## 10. 门禁结果

| 检查 | 结果 |
|---|---|
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 740/740` |
| `--all --skip-base --scope=index` | `TOTAL_FAIL=0` |
| `--all --skip-base --scope=tree` | `TOTAL_FAIL=0` |
| `--all --skip-base --scope=ci` | `TOTAL_FAIL=0` |
| `gofmt -l internal cmd` | 无输出 |
| 编码 | 新增 / 改动 `.go` 为 UTF-8 无 BOM + LF；`.md` 为 UTF-8 with BOM + LF |

## 11. ⑨ 验证结果与本机不可执行项

| 项 | 结果 |
|---|---|
| `go build ./...` | 通过 |
| `go vet ./...` | 通过 |
| `go test ./... -count=1` | 15 个有测试的包全部通过，无 FAIL |
| `GOOS=linux go build ./...` / `go vet ./...` | 通过 |
| `GOOS=darwin go build ./...` | 通过 |
| `-race` | **未执行**（本机无 gcc） |

**本机不可执行项（验收证据改由 CI Ubuntu 腿承担）**：

| 用例 | 原因 |
|---|---|
| `TestP83Fix8T4cBoundDirNamesReadsFromDescriptor` | Linux 专属（build tag） |
| `TestP83Fix8EVerifiedDescriptorInheritedOnlyByItsChild` | Linux 专属（build tag） |
| `TestP83Fix7T3dDescriptorExecNeverResolvesTheReplacedPath` | Linux 专属（build tag） |
| `TestP83Fix8T4bReparseArtifactIsRefused` 的 file-symlink 子项 | 本机无 `SeCreateSymbolicLinkPrivilege`（junction 子项已执行通过） |
| `TestP83Fix7T5c1RetainedSymlinkIsRefused` | 同上（Windows 侧链接拒绝分支以静态与 junction 变体覆盖） |
| `TestP83Fix8T4c` / `E` 的变异复核 | 同上 |

## 12. 环境与档位

| 角色 | 载体与档位 | 设定时刻 | 是否暂停 |
|---|---|---|---|
| 主控 | `DeepSeek V4.1 Flash (deepseek)` High | 起跑时设定 | 否 |
| ① 架构 | `deep-reasoner` + `DeepSeek V4 Pro (deepseek)` Max | 方案书与判定书两轮 | 否 |
| ② 实现 | `prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)` | 全程 | 否 |
| ③ 测试 | `prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)` | 全程 | 否 |
| ⑤ 预审 | `DeepSeek V4 Pro (deepseek)` | 预审环 | 否 |
| ⑥ 扫描 | `MAI Code 2 Flash (copilot)` | 扫描环 | 否 |
| ⑦ 盲审 / 终审 / 复审 | `GPT-6.2 Sol (copilot)` Max | 2026-11-20 取档 | 否 |

## 13. 审查与计量

| 项 | 值 | 说明 |
|---|---|---|
| R | 2 | ⑤ 环整改重跑 1 次 + ⑥ 环整改重跑 1 次；⑦ 复审环按 §9.2 不计 R |
| D | 0 | 本片无临时件纪律违规 |
| 探针 | 0 | 无外部真实调用；F8 真实 BYOK 未调用 |
| ⑦ 额度 | 盲审 1 + 终审 1 + 复审 3 + 追加 1 | 追加轮系用户直接裁定、不得作为先例，报告与台账均须如实披露 |
| ⑤ / ⑥ 口径 | **PASS 绑定整改前快照** | ⑤ 2/2、⑥ 2/3 的 PASS 均早于后续多轮整改；最终快照仅由 ⑦ 复审覆盖，且仍为 FINDINGS |
| 异常 E-1 | `docs/t027/CHECKLIST_v66.md` 出现未暂存工作区改动 | 该文件在「不动」清单内；来源未能归因，已按用户裁定 `git checkout --` 回滚，工作区与暂存区恢复一致（文件集与暂存区相同），不计 D |

## 14. 明确未执行事项

未做 gitee 动作；未升准则版本；未新增 OB / ADR；未改 CI 与 `DEV_PLAN`；未跑真实 BYOK 探针；未跑 `-race`；未在 Linux 实机执行 Linux 专属用例（交由 CI Ubuntu 腿）。

## 15. ⑧a 自查

交付物 5 项齐备（代码与测试、本报告对、台账段与完成记录、SECURITY 注记、证据包）；每 F 有可证伪判据与反例；无测试被弱化或删除；除内部 API 外无契约面改动；结论行与逐 F 状态一致，未出现「闭环」「安全闭合」措辞。

## 16. 六段自查（sixstage）

```sixstage
## Section A 总体结论
PARTIALLY CLOSED

## Section B 发现清单

| SEVERITY | 位置 | 残余 |
|---|---|---|
| High | F3 | T2 STOP：AgentRunner 无冻结策略制品 |
| High | F6 | T3：Windows 名称到启动映像的绑定窗口、Linux 同 inode 字节可变 |
| Medium | F8 | T8a：BYOK 拒发未实施 |
| Medium | F1 | T4：bounddir 大小写折叠与中间 reparse 残余 |
| Medium | F4 | T5#4：批量回写前驱验证、跨进程中断恢复 |

## Section C 逐条证据与推理
- F2 / F5 / F7 与 T9 收窄的守卫、反例测试见 §4 表格。
- ⑦ 六轮结论逐字归档于证据包 04-reviews/ 与 01-arch/ 载荷。
- ⑤ / ⑥ 的 PASS 绑定整改前快照，最终快照仅由 ⑦ 覆盖。

## Section D 可证伪性审计
- 删除任一 F 的主守卫：对应反例测试立即变红。
- 把本片结论改写为「闭环」或「安全闭合」：§4 与 §15 的一致性复核立即红。
- 移除 T2 STOP 或 T8a 的顺延登记：§8 与台账完成记录的互证立即红。
- UNCOVERED：Linux 专属用例与符号链接权限用例未在本机执行；真实 BYOK 未调用。


## Section E 我无法验证的部分
Windows 名称到启动映像的绑定窗口、Linux 同 inode 字节可变性、子进程自身落盘行为均无法在本机静态证明，已登记顺延 P84。

INDEPENDENT SCAN: FINDINGS
```

## 17. 工件围栏

```artifacts
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/00-briefing.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/01-arch/remediation-plan.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/01-arch/ruling-triage.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/07-additional-review.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/07-final-review.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/07-re-reviews-r1-r3.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/MANIFEST.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/SHA256SUMS.txt	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/freedom-list.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/input-list.md	repo
docs/validation/security-findings-remediation.md	repo
docs/validation/security-findings-remediation_EN.md	repo
internal/adapters/agent_runner_admission.go	repo
internal/adapters/agent_runner_replay_dispatcher.go	repo
internal/adapters/ai_preflight.go	repo
internal/adapters/ai_probe_copilot.go	repo
internal/adapters/ai_probe_copilot_executor.go	repo
internal/adapters/ai_probe_copilot_test.go	repo
internal/adapters/bounddir.go	repo
internal/adapters/bounddir_alias.go	repo
internal/adapters/p83_fixround3_tester_test.go	repo
internal/adapters/p83_fixround4_f8_tester_test.go	repo
internal/adapters/p83_fixround5_tester_test.go	repo
internal/adapters/p83_fixround6_tester_test.go	repo
internal/adapters/p83_fixround7_tester_test.go	repo
internal/adapters/p83_fixround8_tester_test.go	repo
internal/chain/engine_test.go	repo
internal/chain/p83_t9a_tester_test.go	repo
internal/chain/review.go	repo
internal/chain/review_test.go	repo
internal/chain/state_seal_counterexample_test.go	repo
internal/console/ai.go	repo
internal/console/ai_test.go	repo
internal/console/approvals.go	repo
internal/console/cli.go	repo
internal/console/cost.go	repo
internal/console/p83_f2_tester_test.go	repo
internal/console/p83_f4_fixround_test.go	repo
internal/console/p83_f6_tester_test.go	repo
internal/console/p83_fixround3_tester_test.go	repo
internal/console/p83_fixround4_tester_test.go	repo
internal/console/p83_fixround5_tester_test.go	repo
internal/console/p83_fixround5_windows_test.go	repo
internal/console/p83_fixround6_tester_test.go	repo
internal/console/p83_fixround7_tester_test.go	repo
internal/console/p83_fixround8_tester_test.go	repo
internal/console/p83_security_findings_test.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_agent_e2e_test.go	repo
internal/console/runtime_multitask_agent_test.go	repo
internal/console/runtime_multitask_driver_test.go	repo
internal/console/runtime_noop.go	repo
internal/console/runtime_supervised.go	repo
internal/console/secret_input_windows.go	repo
internal/console/trusted_ledger.go	repo
internal/evidence/bounddir.go	repo
internal/evidence/bounddir_clean.go	repo
internal/evidence/bounddir_darwin.go	repo
internal/evidence/bounddir_escape_payload_unix_test.go	repo
internal/evidence/bounddir_escape_redbaseline_test.go	repo
internal/evidence/bounddir_escape_unix_test.go	repo
internal/evidence/bounddir_failclosed_unix_test.go	repo
internal/evidence/bounddir_fifo_linux_test.go	repo
internal/evidence/bounddir_identity_test.go	repo
internal/evidence/bounddir_symlink_storm_test.go	repo
internal/evidence/bounddir_unix.go	repo
internal/evidence/bounddir_windows.go	repo
internal/evidence/boundread.go	repo
internal/evidence/fileid_unix.go	repo
internal/evidence/fileid_windows.go	repo
internal/evidence/p83_fixround8_linux_test.go	repo
internal/evidence/p83_fixround8_windows_test.go	repo
internal/gates/adapter.go	repo
internal/gates/p83_f1_f5_tester_test.go	repo
internal/gates/p83_f5_fixround_test.go	repo
internal/gates/p83_fixround5_tester_test.go	repo
internal/gates/p83_fixround6_tester_test.go	repo
internal/gates/p83_fixround7_tester_test.go	repo
internal/gates/policy_test.go	repo
internal/gates/probe_permit_counterexample_test.go	repo
internal/gates/runner.go	repo
internal/gates/runner_artifact_bound_test.go	repo
internal/gates/runner_effect_policy_test.go	repo
internal/gates/runner_test.go	repo
internal/guard/executable_binding.go	repo
internal/guard/executable_binding_error.go	repo
internal/guard/lease.go	repo
internal/guard/lease_open_unix.go	repo
internal/guard/lease_open_windows.go	repo
internal/guard/lease_test.go	repo
internal/guard/p83_fixround5_tester_test.go	repo
internal/guard/p83_fixround5_windows_test.go	repo
internal/guard/p83_fixround6_tester_test.go	repo
internal/guard/p83_fixround6_windows_test.go	repo
internal/guard/p83_fixround7_linux_test.go	repo
internal/guard/p83_fixround7_tester_test.go	repo
internal/guard/p83_fixround7_windows_test.go	repo
internal/guard/p83_fixround8_linux_test.go	repo
internal/guard/p83_fixround8_windows_test.go	repo
internal/guard/process.go	repo
internal/guard/process_env_block_windows_test.go	repo
internal/guard/process_linux.go	repo
internal/guard/process_unsupported.go	repo
internal/guard/process_windows.go	repo
```
