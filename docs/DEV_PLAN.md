# ProofRail 迭代开发计划

[English](DEV_PLAN_EN.md)

日期：2026-09-08；S1 实现尚未全部完成。输入：[需求](PRODUCT_REQUIREMENTS.md)、[架构](ARCHITECTURE.md)、[契约](CONTRACTS.md)、[安全](SECURITY.md)、[ADR](ADR_REGISTER.md)。规划遵循修订后的 RFC §14/19；新增 S1 最小范围需重新确认预算，不将 S2/S3 增强移入首版。

## 1. 工作方式与约束

Go 1.22 模块兼容基线、纯 Go、CLI/TUI、内容寻址文件存储，具体依赖由 S0 冻结。不存在框架迁移，本计划不套用 Java/.NET 迁移规则。遵循只读源目录、single writer、first-fail-stop、独立批准、无默认 Git 写操作；不存在架构例外。

建议每个开发会话只处理一个任务；大于一次可验证修改的任务继续按子测试拆分，不一次生成全部模块。采用红测→最小实现→原用例→包测试→全门禁→更新双语文档/证据。先运行离线确定性测试，再申请模型调用。日历估算以首次两个任务实测速度更新，不承诺未经测量的完成日期。

WIP=1；每个任务执行前设定本轮人工认可的费用/时长上限，默认不发收费调用；两次无新证据的同类失败暂停升级。最强模型仅用于协议/安全歧义评审，小模型只执行已冻结契约的小切片。所有 Git commit/push 仍需当轮授权，默认只 origin。

## 2. 迭代与任务

每行包含责任路径（尚未存在的路径是计划创建位置）、依赖、验收。仅 T001-T004 是 S0；T004 是隔离技术探针，不是 S1 生产实现。T005 之前必须 §17 全绿且人工批准。

### P0.1 文档与授权基线（REQ-001/014/017/023/026/028）

- [x] T001 [Plan:P0.1] 评审 `docs/PRODUCT_REQUIREMENTS*.md`、`docs/SECURITY*.md`、`docs/ADR_REGISTER*.md`；记录所有者/安全评审者、预算、支持/签名和影子范围；依赖：本包；验收：§17.1/6/8 签署，依赖许可决策可查；禁止：代用户批准。完成：2026-09-06，所有者批准 S0 规格冻结及 ADR §4 决策单；独立安全结论仍由 T003/§17.3 单独关闭。

T001 状态：`COMPLETE`。本次仅批准 S0 规格冻结并允许进入 T002；不授权新增付费调用、具体 whois 数据访问、Git 操作、发布或 S1 生产编码。精确依赖及独立安全评审仍按后续门禁落实；2026-09-08 批准的简化发布信任策略不要求 S1 签名主体或公钥。

### P0.2 机器契约冻结（REQ-001/004/005/007/008/010/012/013/017/018/020/021/022/024/025/028）

- [x] T002 [Plan:P0.2] 先修订 `docs/RFC-proofrail-unattended-ai-engineering-product.md`，同步 `docs/CONTRACTS*.md`/ADR，建立 `schemas/*.schema.json`；冻结 ADR-002/004 的字段、版本、canonical 向量、状态表、错误与队列；依赖：T001；验收：结构与本地一致性检查完成，独立 fixtures/validator 由 T003 承接；禁止：以草案字段实现生产代码。
- [x] T003 [Plan:P0.2] 建立 `testdata/contracts/{valid,invalid}/`、`tools/contracts/` 独立 validator/语义 fixture oracle，纳入三任务、四 kind、双/三语言、handoff、doc/waiver、各 receipt/manifest；依赖：T002；验收：AT-01 正反例及 RFC §17.2/4/9-12；禁止：让未来候选自己生成 expected 结果。完成：2026-09-07，Ajv 2020 独立编译 29 份 Schema；当前 82 个手写预期 fixture 全部通过，`canonicalize` 独立复算两条 JCS/摘要向量。

T002/T003 状态：`COMPLETE`。已冻结共同版本/ID/路径/JCS/UTF-8 规则，并建立 chain、hook、target、
workspace、state-event、error、adapter-envelope、adapter-receipt、snapshot-manifest、evidence-manifest、
review-receipt、promotion-receipt、handoff-receipt、hook-result、run-manifest、signature-receipt、error-set、ticket-ledger、repair-transaction、harness、toolchain、change-set、verification-report、plan-preview、export-record、authorization-record、effect-record、cost-ledger、lifecycle-record 二十九份 Draft 2020-12 结构 Schema；首批 canonical 向量、单记录状态转换、37 个错误 code/主错误排序与 CLI 退出码、
文件队列 framing/claim/result/dispatch/takeover receipt、快照边界、评审前证据根及 approve/reject/waive
评审结论和原子 promotion reader contract、manual-handoff step 判别分支与归还 receipt、hook 执行事实/门禁判断/策略处置分层、run 有效配置来源与运行绑定、harness/toolchain 注册表、detached Ed25519 signature receipt、失败指纹/ticket ledger/三段修复预算、四阶段 repair transaction、确定性配置合并、change-set、跨事件/跨记录/信任链 checker 报告及 PC-01–PC-06 持久记录已冻结。T003 的独立结构 validator、静态语义 checker、向量 fixture 与正反 golden 已执行；AT-01 的 S0 契约部分及 RFC §17.2/4/9-13 已有机器证据。真实文件、进程、能力与平台事实仍属 T004/S1 runtime gate，整体 S0 仍为 `NOT_READY`。

T001 同时评审 ADR-008 与 PC-01–PC-06 新增工作量；T002 冻结 RFC §19 的授权/预览/导出/副作用/成本/备份记录；T003 已增加对应独立正反黄金样例，覆盖 §17.13。此处只冻结契约，不提前执行 S1 验收。

### P0.3 技术探针（REQ-006/007/009/015/017/021/023/024/027）

- [x] T004 [Plan:P0.3] 在根 `tmp/` 隔离验证 Go/TUI、Windows 原子替换/进程树停止、Linux 路径、IPC/文件队列；结论入 `docs/validation/s0-spikes.md`，删除实验文件；依赖：T002，S1 前还需 T003/T001 全部门禁；验收：可重复命令、实际平台、失败能力明确；禁止：未授权模型付费、原型冒充生产实现。完成：2026-09-07，Windows 11/NTFS 与 Ubuntu 24.04/ext4 原生探针、Bubble Tea v1.3.4 隔离原型及 SessionBridge 40 项无模型测试通过；限制与 fail-close 规则见报告。

T004 状态：`COMPLETE`。技术可行性已验证，但整体 S0 仍为 `NOT_READY`；依赖供应链、bootstrap、CI 发布证据与 S1 生产故障注入仍按 T017/后续任务关闭，本结果不授权 T005。

### P1.1 证据和快照底座（REQ-004/023/024/025/028）

- [x] T005 [Plan:P1.1] 实现 `internal/evidence/{canonical,event,receipt,verify}.go` 与同名 `_test.go`，覆盖 canonical/事件链/对象引用和只读核验；依赖：S0 全通过；验收：AT-04 证据部分；禁止：自由文本 PASS 或未知版本写入。完成：2026-09-07，冻结向量、严格 JSON/JCS、事件投影链、evidence manifest 及内容存储重读核验通过；报告见 `docs/validation/t005-evidence.md`。

T005 状态：`COMPLETE`。结果只覆盖证据层的 AT-04 子集；journal、写入事务、回滚及各崩溃点仍由 T008 验收，不得由本状态推导为通过。

- [x] T006 [Plan:P1.1] 实现 `internal/snapshot/{capture,restore,store,types,reparse_*}.go` 与测试，路径验证/排除/配额/保留及 GC dry-run；依赖：T005；验收：AT-02；禁止：读 Git 历史恢复、自动清理有引用对象。完成：2026-09-07，未提交文件树捕获、CAS 对象存储、保留名/大小写/长路径/软硬链接检测、捕获中变动校验、配额限制、无 Git 恢复及 Windows/Ubuntu 双平台验证通过；报告见 `docs/validation/t006-snapshot.md`。

T006 状态：`COMPLETE`。AT-02 完整验收通过；无 Git 恢复与不可变快照已闭环；工作区写事务与崩溃回滚点由 T008 承接。

### DG-01 独立产品文档治理门禁

- [x] DG-01 明确 ProofRail 独立产品边界，新增 `BUSINESS_WORKFLOWS*.md` 与 `INSTALLATION_PLAN*.md`，同步 RFC、需求、导航、README 和操作手册；完成：2026-09-07。安装载体、目录、签名主体和升级方式保持待决，不生成虚假安装步骤。T 编号保持稳定，DG-01 是 T009 的文档前置门禁，不占用实施任务编号。

### DG-02 分域权威迁移门禁

- [x] DG-02 将原 RFC 定位为项目建议书与历史设计来源，在 `DOCUMENTATION_PLAN*.md` 建立分域权威矩阵，并把 T009 直接依赖的 chain/task/step 转换、事件投影、顺序调度、暂停取消和崩溃恢复迁入 `CONTRACTS*.md` §2.1；完成：2026-09-07。旧章节保留用于来源追踪，冲突必须阻断并修订，不由实现自行选择。

### P1.2 写入与顺序闭环（REQ-002/003/004/012/018/023/024/026）

- [x] T007 [Plan:P1.2] 实现 `internal/guard/{process,lease}.go`、平台文件及测试；证明进程树停机、接管和租约；依赖：T005、T004；验收：AT-03/06；禁止：杀用户进程、只凭 PID 或过期接管。验证：[中文](validation/t007-process-lease.md) / [English](validation/t007-process-lease_EN.md)。
- [x] T008 [Plan:P1.2] 实现 `internal/taskdef/checker.go`、`internal/applier/{apply,journal,recover}.go` 与测试；内存顺序预验→事务→回滚；依赖：T005/T006/T007；验收：AT-04 每崩溃点；禁止：边验边写或跳过断言。验证：[中文](validation/t008-applier.md) / [English](validation/t008-applier_EN.md)。
- [x] T009 [Plan:P1.2] 实现 `internal/chain/{engine,state,recover}.go` 与测试；三任务、四 kind、两变更模式、事件投影/暂停取消，先 fake agent/runner；依赖：T008、DG-01、DG-02；验收：AT-03；禁止：候选作为下一父快照。完成：2026-09-07。14 个切片及 Windows 11/NTFS、Ubuntu 24.04/ext4 原生验证通过。验证：[中文](validation/t009-chain.md) / [English](validation/t009-chain_EN.md)。
- [x] T010 [Plan:P1.2] 实现 `internal/gates/{runner,policy,result}.go` 与测试并接 chain；依赖：T007/T009；验收：AT-06；禁止：不具备能力时弱化网络/资源策略。完成：2026-09-07。12 个切片及 Windows 11/NTFS、Ubuntu 24.04/ext4 原生验证通过；未实现的网络/资源 enforcement 保持 preflight 阻断。验证：[中文](validation/t010-gates.md) / [English](validation/t010-gates_EN.md)。

### P1.3 独立接受与有限修复（REQ-005/014/022/023/024/025/026）

- [x] T011 [Plan:P1.3] 实现 `internal/chain/{review,publish}.go` 与测试，绑定候选/review/receipt，支持拒绝和有效 waiver；依赖：T010；验收：AT-05；禁止：agent 自批、技术通过即 PASSED。完成：2026-09-07。覆盖拒评、自批、过期/错候选 waiver 拒绝，以及有效独立批准后单次发布。验证：[中文](validation/t011-chain.md) / [English](validation/t011-chain_EN.md)。
- [x] T012 [Plan:P1.3] 实现 `internal/tickets/{ledger,budget,fingerprint}.go`、`internal/repair/transaction.go` 与测试；依赖：T011；验收：AT-07；禁止：重置预算、改写正式定义绕过候选。完成：2026-09-07。覆盖相同指纹耗尽、attempt/墙钟/费用耗尽、假修复、陈旧候选与 Promote 中断。验证：[中文](validation/t012-tickets-repair.md) / [English](validation/t012-tickets-repair_EN.md)。

### P1.4 代理和人工交接（REQ-006/007/008/016/021/023/024/025/026）

- [x] T013 [Plan:P1.4] 实现 `internal/adapters/{filequeue,sessbridge}.go` 与测试，核心端口由消费方定义；依赖：T012；验收：AT-08，先离线后授权实机；禁止：修改 SessionBridge、GUI/auto 回退、重复应用。完成：2026-09-07。覆盖 request/claim/result 信封、file-queue 领取与 fencing、幂等冲突防护，以及 sessbridge silent 派发的 accepted/rejected/uncertain 映射。验证：[中文](validation/t013-adapters.md) / [English](validation/t013-adapters_EN.md)。
- [x] T014 [Plan:P1.4] 实现 `internal/chain/handoff.go`、`internal/console/handoff.go` 与测试；依赖：T013；验收：AT-09/10；禁止：秘密经模型、到期自动通过、自动恢复并发写。完成：2026-09-07。覆盖 handoffPolicy 约束、complete/abort/request-agent 收据语义、WAITING_FOR_OPERATOR 状态转换，以及 structured/secret-direct 输入策略与秘密摘要化。补充验证：2026-09-08 在 Ubuntu VM 通过临时代理链路完成远端全量 build/vet/test 与 race（gcc 全仓、clang 关键包）并全部通过。验证：[中文](validation/t014-handoff.md) / [English](validation/t014-handoff_EN.md)。

### P1.5 工程包与可用 CLI/TUI（REQ-009/010/011/013/015/016/018/020/022/027）

- [x] T015 [Plan:P1.5] 建立 `harnesses/{generic,c,go}/`、`internal/taskdef/documentation.go` 与测试，模板 A、impact rule、生成新鲜度、C+Go task；依赖：T014；验收：AT-11；禁止：whois 专名进入核心或交付任意脚本 B。完成：2026-09-08。覆盖文档协同影响规则、模板 A 锁定与生成物新鲜度检查、自批/手改生成物阻断，以及 C+Go 语言作用域任一失败触发同 task 整体恢复。验证：[中文](validation/t015-taskdef-documentation.md) / [English](validation/t015-taskdef-documentation_EN.md)。
- [x] T016 [Plan:P1.5] 完成 `cmd/prfrail/main.go`、`internal/console/`、配置解析/解释与测试，同步 `README.md`/`docs/OPERATIONS*.md`；依赖：T015；验收：AT-12；禁止：视图直接改状态、占位命令报告成功。完成：2026-09-08。覆盖无 IDE `init/validate/config explain/run/report` 命令、稳定 `--json` 输出、配置搜索顺序（`--chain`→`proofrail.chain.json`→`proofrail.json`）、noop-only 本地运行闭环与执行型 step fail-close 非零退出。验证：[中文](validation/t016-cli-console.md) / [English](validation/t016-cli-console_EN.md)。

### P1.6 可信发布与影子验收（REQ-001/002/006/011/014/015/019/024/025/028）

- [x] T017 [Plan:P1.6] 建立 `.github/workflows/ci.yml`、`testdata/selfhost/`、`docs/validation/s1-selfhost.md`；固定工具链与 Actions 完整 SHA、两代隔离、较早可信 verifier commit、外部 oracle、原生平台、固定 candidate commit、SHA256SUMS、SBOM/许可证检查；依赖：T024 和独立 bootstrap 批准；验收：AT-13/14；禁止：候选构建/覆盖 verifier 或 seed、把校验和冒充身份认证、未授权发布。签名/attestation 属后续增强，不阻断 S1。完成：2026-09-09。baseline `12fa05d` 的 push CI run `34253473081` 在 Windows/Linux 均失败后，CI 修复版经 DeepSeek V4 Pro 与 GPT-5.3 Codex 双主审 `PASS FOR BASELINE REVIEW`，并以 `18c9340c60deb30049554491f8905df51d0a3180` 建立新的 verifier baseline；push CI run [`34268076912`](https://github.com/larsonzh/prfrail/actions/runs/34268076912) 的 Ubuntu/Windows job 均成功。不同的 docs-only candidate `1e7af676e8e84028c7ecfef1ccde91213738c20c` 随后由维护者从 `main` 手工 dispatch；[run `34277671704`](https://github.com/larsonzh/prfrail/actions/runs/34277671704) 的 test、candidate-build、candidate-probe、selfhost Windows/Linux 8 个 job 全部成功，6 个 artifact 的名称、大小和 GitHub archive digest 已归档。验证：[中文](validation/s1-selfhost.md) / [English](validation/s1-selfhost_EN.md)；初始 Codex 主审：[中文](validation/t017-baseline-review.md) / [English](validation/t017-baseline-review_EN.md)；初始 DeepSeek 主审：[中文](validation/t017-deepseek-baseline-review.md) / [English](validation/t017-deepseek-baseline-review_EN.md)；CI 修复版 DeepSeek 复审：[中文](validation/t017-ci-remediation-deepseek-review.md) / [English](validation/t017-ci-remediation-deepseek-review_EN.md)；CI 修复版 Codex 复审：[中文](validation/t017-ci-remediation-codex-review.md) / [English](validation/t017-ci-remediation-codex-review_EN.md)。
- [ ] T018 [Plan:P1.6] 建立 `examples/{whois-shadow,go-minimal}/` 与 `docs/validation/s1-exit.md`；只读固定输入比较结果/失败分类，汇总全部 AT 和剩余风险；依赖：T017、T025 及 whois 输入授权；验收：AT-15/22、RFC §14 S1 全部 exit；禁止：切换 whois 生产流程或“部分通过”冒充 S1 完成。

T018 状态：`IN_PROGRESS / BLOCKED`。2026-09-09 已在单轮只读授权内导出脱敏 whois 固定夹具并通过 9/9 结果对齐及输入/结果/失败分类漂移负例；generic/C/Go 三个最小 CLI 闭环与 Windows 候选便携包的首次运行、升级、回滚、保留数据卸载演练均通过。产品所有者已批准手工解压便携 ZIP、不修改 PATH 的安装模型，双语教程、支持矩阵和发行说明草案已建立。AT-01–AT-22 与剩余门禁见 [S1 exit 报告](validation/s1-exit.md)；T025/AT-22、最终 ZIP 绑定/复验、EOL 通知方式和产品所有者 S1 exit 批准仍缺失，因此不勾选 T018。

### P1.7 产品闭环与生命周期（在 P1.6 发布验收前执行）

编号保持稳定，实际依赖顺序为 T016 → T019–T024 → T017 → T025 → T018。以下在 S0 授权后执行；T019–T024 已完成。每个任务可再按一个失败反例拆小，不一次生成完整功能。

- [x] T019 [Plan:P1.7] PC-01：实现 `internal/taskdef/preview.go`、`internal/console/preview.go` 与测试、无 AI 示例；依赖：T016；验收：AT-16；禁止：预览启动命令/网络/模型/自动安装。完成：2026-09-08。覆盖离线只读 preview record 生成、unknown 可见、零调用计数与 JSON/文本同事实输出。验证：[中文](validation/t019-preview.md) / [English](validation/t019-preview_EN.md)。
- [x] T020 [Plan:P1.7] PC-02：实现 `internal/snapshot/export.go`、`internal/evidence/delivery.go` 与测试；依赖：T019、T006/T011；验收：AT-17；禁止：导出候选冒充接受、覆盖目标、源树/Git 写入。完成：2026-09-08。覆盖 accepted hash 绑定、目标必须 absent、源/run/store 重叠阻断、secret 路径阻断、缺对象/篡改阻断、写中断无 completed record。验证：[中文](validation/t020-export.md) / [English](validation/t020-export_EN.md)。
- [x] T021 [Plan:P1.7] PC-03：实现 `internal/chain/authorization.go`、`internal/console/approvals.go` 与测试，接 guard 停机；依赖：T020、T007/T011；验收：AT-18；禁止：自授权、硬门禁 waiver、过期自动批准。完成：2026-09-08。覆盖授权/撤销绑定摘要、撤销/过期/错摘要阻断下一动作与接受、在途停止失败保持暂停并可重启可见。验证：[中文](validation/t021-approvals.md) / [English](validation/t021-approvals_EN.md)。
- [x] T022 [Plan:P1.7] PC-04：实现 `internal/gates/effects.go`、`internal/evidence/diagnostics.go` 和 chain 恢复计划测试；依赖：T021、T010；验收：AT-19；禁止：外部写 runner、未知副作用盲重试、诊断删锁/日志。完成：2026-09-08。覆盖副作用分类与准入阻断、unknown/外部写 fail-close、只读脱敏诊断与恢复动作。验证：[中文](validation/t022-effects-diagnostics.md) / [English](validation/t022-effects-diagnostics_EN.md)。
- [x] T023 [Plan:P1.7] PC-05：实现 `internal/tickets/cost.go`、`internal/adapters/cost_usage.go`、`internal/console/cost.go` 与测试；依赖：T022、T012/T013；验收：AT-20；禁止：重启释放未知费用、跨 run 绕上限、默认遥测。完成：2026-09-08。覆盖调用前持久 reservation、settlement 去重、unknown hold 重启保留、跨 run 共享上限阻断与本地无默认遥测成本报告。验证：[中文](validation/t023-cost-ledger.md) / [English](validation/t023-cost-ledger_EN.md)。
- [x] T024 [Plan:P1.7] PC-06：实现 `internal/snapshot/backup.go`、`internal/evidence/disposition.go`、生命周期控制/测试及发布清单工具；依赖：T023；验收：AT-21、SBOM/支持说明；禁止：未知格式迁移、默认删证据/共享工具、擅自签名发布。完成：2026-09-08。覆盖完整对象/事件/引用闭包备份、新 store 复原核验、源 store 不变、活动写者/未知版本/缺闭包阻断、默认保留证据与共享工具、显式删除授权/保留冲突决定，以及离线撤销状态 unknown。验证：[中文](validation/t024-lifecycle.md) / [English](validation/t024-lifecycle_EN.md)。

### P1.8 AI/操作员交互闭环（REQ-006/007/008/009/016/021/023/024/025/027）

- [ ] T025 [Plan:P1.8] 冻结 `operator-interaction` Schema、正反 golden、哈希与重放契约；实现内嵌 SessionBridge v1 文件 IPC、稳定非空 `conversationId`、TUI 待处理交互收件箱/响应和重启恢复；依赖：T014、T016、T021、T024；验收：AT-22；禁止：自由文本提问触发状态转换、把聊天历史当权威状态、依赖 `@sbr-review`、绕过 authorization/handoff、秘密进入模型或证据。`@sbr-review` 可保留为 SessionBridge 独立诊断能力。状态：`PLANNED / NOT IMPLEMENTED`。

## 3. 需求映射

除明确标记完成的任务外，Implementation evidence 是计划证据，不代表实现。

| REQ ID | Plan Items | 任务 | 验收/实现证据 |
|---|---|---|---|
| REQ-001 | P0.1,P0.2,P1.6 | T001,T002,T003,T018 | AT-01/15；通用 Schema/非 C 样例 |
| REQ-002 | P1.2,P1.6 | T009,T018 | AT-03/15；chain/报告 |
| REQ-003 | P1.2 | T009 | AT-03；有序调度 |
| REQ-004 | P0.2,P1.1,P1.2 | T002,T006,T008 | AT-02/04；snapshot/applier |
| REQ-005 | P0.2,P1.3 | T002,T011 | AT-05；review/publish |
| REQ-006 | P0.3,P1.4,P1.6 | T004,T013,T017 | AT-08/14；adapter/CI |
| REQ-007 | P0.2,P0.3,P1.4,P1.8 | T002,T004,T013,T025 | AT-08/22；双通道/内嵌 IPC |
| REQ-008 | P0.2,P1.4,P1.8 | T002,T013,T025 | AT-08/22；context/交互恢复 |
| REQ-009 | P0.3,P1.5,P1.8 | T004,T016,T025 | AT-12/22；CLI/TUI 交互 |
| REQ-010 | P0.2,P1.5 | T002,T016 | AT-01/12；配置解释 |
| REQ-011 | P1.5,P1.6 | T015,T018 | AT-11/15；三 harness |
| REQ-012 | P0.2,P1.2 | T002,T007,T010 | AT-06；gate runner |
| REQ-013 | P0.2,P1.5 | T002,T015 | AT-11；模板生成 |
| REQ-014 | P0.1,P1.3,P1.6 | T001,T012,T018 | AT-07/15；阶段报告 |
| REQ-015 | P0.3,P1.5,P1.6 | T004,T016,T017 | AT-12/13/14；构建/TUI |
| REQ-016 | P1.4,P1.5,P1.8 | T013,T016,T025 | AT-08/12/22；无 IDE CLI/宿主交互 |
| REQ-017 | P0.1,P0.2,P0.3 | T001,T002,T003,T004 | AT-01；readiness |
| REQ-018 | P0.2,P1.2,P1.5 | T003,T009,T016 | AT-01/03；四 kind |
| REQ-019 | P1.6 | T017 | AT-13/14；两代自托管 |
| REQ-020 | P0.2,P1.5 | T003,T015 | AT-01/11；组件/双语言 |
| REQ-021 | P0.2,P0.3,P1.4,P1.8 | T003,T004,T014,T025 | AT-09/10/22；handoff/交互分流 |
| REQ-022 | P0.2,P1.3,P1.5 | T003,T011,T015 | AT-01/05/11；文档协同 |
| REQ-023 | P0.1,P0.3,P1.1,P1.2,P1.3,P1.4,P1.8 | T001,T004,T005,T007,T010,T011,T012,T014,T025 | AT-05/06/07/09/22；安全 |
| REQ-024 | P0.2,P0.3,P1.1,P1.2,P1.3,P1.4,P1.6,P1.8 | T002,T004,T005,T006,T008,T009,T011,T014,T017,T025 | AT-02/03/04/09/13/22；恢复 |
| REQ-025 | P0.2,P1.1,P1.3,P1.4,P1.6,P1.8 | T002,T005,T011,T013,T014,T017,T025 | AT-04/06/10/14/22；交互证据 |
| REQ-026 | P0.1,P1.2,P1.3,P1.4 | T001,T009,T010,T012,T013 | AT-03/06/07；预算 |
| REQ-027 | P0.3,P1.5,P1.8 | T004,T016,T025 | AT-12/22；终端/JSON/交互 |
| REQ-028 | P0.1,P0.2,P1.1,P1.6 | T001,T002,T003,T005,T017 | AT-01/04/14；兼容 |

## 4. 测试策略与完成定义

§3 原映射与下表取并集；checkpoint 的产品切片补充记录同样加入 P0.1/P0.2/P1.7，不另编号需求。

| PC | REQ | Plan Items | 任务 | 验收 |
|---|---|---|---|---|
| PC-01 | REQ-009/010/016/027 | P0.1,P0.2,P1.7 | T001,T002,T003,T019 | AT-16 |
| PC-02 | REQ-002/004/025 | P0.1,P0.2,P1.7 | T001,T002,T003,T020 | AT-17 |
| PC-03 | REQ-005/021/023 | P0.1,P0.2,P1.7 | T001,T002,T003,T021 | AT-18 |
| PC-04 | REQ-012/023/024 | P0.1,P0.2,P1.7 | T001,T002,T003,T022 | AT-19 |
| PC-05 | REQ-006/026 | P0.1,P0.2,P1.7 | T001,T002,T003,T023 | AT-20 |
| PC-06 | REQ-019/025/028 | P0.1,P0.2,P1.7 | T001,T002,T003,T024 | AT-21 |

完整方案见 [TEST_STRATEGY.md](TEST_STRATEGY.md)：Go testing、真实本地文件/子进程、独立 Schema validator、CLI oracle、授权宿主 smoke；无数据库/浏览器依赖。模块测试必须纳入 go test ./...；每任务报告测试数/失败/跳过。mock 通过不代替原生 Windows 进程树/文件系统或真实 adapter 能力。

完成任务须有：实现与测试、对应需求/契约、双语影响文档、命令/退出码/证据、范围无越界。完成迭代还须全 build/vet/test 与相关 AT；完成 S1 必须全强制测试和人审签署。计划覆盖率不是实现覆盖率。

## 5. 后续阶段和资源决策

S2 backlog 在 S1 exit 后细化，不预生成空实现：Linux 正式、多语言组件依赖编排/三语言实跑、supervised/Web/可选 VS Code 视图、模型规则切换和智能文档影响提示。S3 为生成脚本 B、多编辑器、inception。需求 ID 保持不变，通过 RFC/ADR 扩展任务和 AT；不得从 S1 悄悄搬走强制能力。

每次迭代报告实际人工时间、模型调用/费用估算、故障数、已验收需求和下一步；若超预算，暂停并由用户选择继续或经 RFC 缩范围，不自动加模型/服务。

RFC §19 的后续候选：S2 图形上手/终端本地化、脱敏通知/远程审批身份设计、外部 PR 集成评估、外部写补偿、成本趋势和显式存储迁移；S3 商业后台/云管理须先验证需求。这些不在 T019–T025 的授权范围，后续再拆任务。当前共 11 个 plan items、2 个未完成任务、22 组计划验收；不是完成率。