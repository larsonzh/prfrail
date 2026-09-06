# ProofRail 迭代开发计划

[English](DEV_PLAN_EN.md)

日期：2026-09-06；全部实现任务尚未完成。输入：[需求](PRODUCT_REQUIREMENTS.md)、[架构](ARCHITECTURE.md)、[契约](CONTRACTS.md)、[安全](SECURITY.md)、[ADR](ADR_REGISTER.md)。规划遵循修订后的 RFC §14/19；新增 S1 最小范围需重新确认预算，不将 S2/S3 增强移入首版。

## 1. 工作方式与约束

Go 1.22 模块兼容基线、纯 Go、CLI/TUI、内容寻址文件存储，具体依赖由 S0 冻结。不存在框架迁移，本计划不套用 Java/.NET 迁移规则。遵循只读源目录、single writer、first-fail-stop、独立批准、无默认 Git 写操作；不存在架构例外。

建议每个开发会话只处理一个任务；大于一次可验证修改的任务继续按子测试拆分，不一次生成全部模块。采用红测→最小实现→原用例→包测试→全门禁→更新双语文档/证据。先运行离线确定性测试，再申请模型调用。日历估算以首次两个任务实测速度更新，不承诺未经测量的完成日期。

WIP=1；每个任务执行前设定本轮人工认可的费用/时长上限，默认不发收费调用；两次无新证据的同类失败暂停升级。最强模型仅用于协议/安全歧义评审，小模型只执行已冻结契约的小切片。所有 Git commit/push 仍需当轮授权，默认只 origin。

## 2. 迭代与任务

每行包含责任路径（尚未存在的路径是计划创建位置）、依赖、验收。仅 T001-T004 是 S0；T004 是隔离技术探针，不是 S1 生产实现。T005 之前必须 §17 全绿且人工批准。

### P0.1 文档与授权基线（REQ-001/014/017/023/026/028）

- [x] T001 [Plan:P0.1] 评审 `docs/PRODUCT_REQUIREMENTS*.md`、`docs/SECURITY*.md`、`docs/ADR_REGISTER*.md`；记录所有者/安全评审者、预算、支持/签名和影子范围；依赖：本包；验收：§17.1/6/8 签署，依赖许可决策可查；禁止：代用户批准。完成：2026-09-06，所有者批准 S0 规格冻结及 ADR §4 决策单；独立安全结论仍由 T003/§17.3 单独关闭。

T001 状态：`COMPLETE`。本次仅批准 S0 规格冻结并允许进入 T002；不授权新增付费调用、具体 whois 数据访问、Git 操作、发布或 S1 生产编码。签名主体、公钥、精确依赖及独立安全评审仍按后续门禁落实。

### P0.2 机器契约冻结（REQ-001/004/005/007/008/010/012/013/017/018/020/021/022/024/025/028）

- [ ] T002 [Plan:P0.2] 先修订 `docs/RFC-proofrail-unattended-ai-engineering-product.md`，同步 `docs/CONTRACTS*.md`/ADR，建立 `schemas/*.schema.json`；冻结 ADR-002/004 的字段、版本、canonical 向量、状态表、错误与队列；依赖：T001；验收：AT-01 结构校验及独立复核；禁止：以草案字段实现生产代码。
- [ ] T003 [Plan:P0.2] 建立 `testdata/contracts/{valid,invalid}/`、`tools/contracts/` 独立 validator/语义 fixture oracle，纳入三任务、四 kind、双/三语言、handoff、doc/waiver、各 receipt/manifest；依赖：T002；验收：AT-01 正反例及 RFC §17.2/4/9-12；禁止：让未来候选自己生成 expected 结果。

T002 状态：`IN_PROGRESS`。已冻结共同版本/ID/路径/JCS/UTF-8 规则，并建立 chain、hook、target、
workspace 四份 Draft 2020-12 结构 Schema；尚待 canonical 测试向量、状态机、错误码、队列、其余持久
对象及 ADR-004。T003 独立 validator/语义 checker 与正反 golden 尚未执行，因此不得标记 T002、AT-01
或 RFC §17.2 完成。

T001 同时评审 ADR-008 与 PC-01–PC-06 新增工作量；T002 冻结 RFC §19 的授权/预览/导出/副作用/成本/备份记录；T003 增加对应独立正反黄金样例，覆盖 §17.13。此处只冻结契约，不提前执行 S1 验收。

### P0.3 技术探针（REQ-006/007/009/015/017/021/023/024/027）

- [ ] T004 [Plan:P0.3] 在根 `tmp/` 隔离验证 Go/TUI、Windows 原子替换/进程树停止、Linux 路径、IPC/文件队列；结论入 `docs/validation/s0-spikes.md`，删除实验文件；依赖：T002，S1 前还需 T003/T001 全部门禁；验收：可重复命令、实际平台、失败能力明确；禁止：未授权模型付费、原型冒充生产实现。

### P1.1 证据和快照底座（REQ-004/023/024/025/028）

- [ ] T005 [Plan:P1.1] 实现 `internal/evidence/{event,receipt,verify}.go` 与同名 `_test.go`，覆盖 canonical/事件链/对象引用和只读核验；依赖：S0 全通过；验收：AT-04；禁止：自由文本 PASS 或未知版本写入。
- [ ] T006 [Plan:P1.1] 实现 `internal/snapshot/{capture,restore,store}.go` 与测试，路径验证/排除/配额/保留及 GC dry-run；依赖：T005；验收：AT-02；禁止：读 Git 历史恢复、自动清理有引用对象。

### P1.2 写入与顺序闭环（REQ-002/003/004/012/018/023/024/026）

- [ ] T007 [Plan:P1.2] 实现 `internal/guard/{process,lease}.go`、平台文件及测试；证明进程树停机、接管和租约；依赖：T005、T004；验收：AT-03/06；禁止：杀用户进程、只凭 PID 或过期接管。
- [ ] T008 [Plan:P1.2] 实现 `internal/taskdef/checker.go`、`internal/applier/{apply,journal,recover}.go` 与测试；内存顺序预验→事务→回滚；依赖：T005/T006/T007；验收：AT-04 每崩溃点；禁止：边验边写或跳过断言。
- [ ] T009 [Plan:P1.2] 实现 `internal/chain/{engine,state,recover}.go` 与测试；三任务、四 kind、两变更模式、事件投影/暂停取消，先 fake agent/runner；依赖：T008；验收：AT-03；禁止：候选作为下一父快照。
- [ ] T010 [Plan:P1.2] 实现 `internal/gates/{runner,policy,result}.go` 与测试并接 chain；依赖：T007/T009；验收：AT-06；禁止：不具备能力时弱化网络/资源策略。

### P1.3 独立接受与有限修复（REQ-005/014/022/023/024/025/026）

- [ ] T011 [Plan:P1.3] 实现 `internal/chain/{review,publish}.go` 与测试，绑定候选/review/receipt，支持拒绝和有效 waiver；依赖：T010；验收：AT-05；禁止：agent 自批、技术通过即 PASSED。
- [ ] T012 [Plan:P1.3] 实现 `internal/tickets/{ledger,budget,fingerprint}.go`、`internal/repair/transaction.go` 与测试；依赖：T011；验收：AT-07；禁止：重置预算、改写正式定义绕过候选。

### P1.4 代理和人工交接（REQ-006/007/008/016/021/023/024/025/026）

- [ ] T013 [Plan:P1.4] 实现 `internal/adapters/{filequeue,sessbridge}.go` 与测试，核心端口由消费方定义；依赖：T012；验收：AT-08，先离线后授权实机；禁止：修改 SessionBridge、GUI/auto 回退、重复应用。
- [ ] T014 [Plan:P1.4] 实现 `internal/chain/handoff.go`、`internal/console/handoff.go` 与测试；依赖：T013；验收：AT-09/10；禁止：秘密经模型、到期自动通过、自动恢复并发写。

### P1.5 工程包与可用 CLI/TUI（REQ-009/010/011/013/015/016/018/020/022/027）

- [ ] T015 [Plan:P1.5] 建立 `harnesses/{generic,c,go}/`、`internal/taskdef/documentation.go` 与测试，模板 A、impact rule、生成新鲜度、C+Go task；依赖：T014；验收：AT-11；禁止：whois 专名进入核心或交付任意脚本 B。
- [ ] T016 [Plan:P1.5] 完成 `cmd/prfrail/main.go`、`internal/console/`、配置解析/解释与测试，同步 `README.md`/`docs/OPERATIONS*.md`；依赖：T015；验收：AT-12；禁止：视图直接改状态、占位命令报告成功。

### P1.6 可信发布与影子验收（REQ-001/002/006/011/014/015/019/024/025/028）

- [ ] T017 [Plan:P1.6] 建立 `.github/workflows/ci.yml`、`testdata/selfhost/`、`docs/validation/s1-selfhost.md`；固定工具链、两代隔离、外部 oracle、原生平台/签名检查；依赖：T024 和独立 bootstrap 批准；验收：AT-13/14；禁止：候选覆盖 seed、未授权发布。
- [ ] T018 [Plan:P1.6] 建立 `examples/{whois-shadow,go-minimal}/` 与 `docs/validation/s1-exit.md`；只读固定输入比较结果/失败分类，汇总全部 AT 和剩余风险；依赖：T017 及 whois 输入授权；验收：AT-15、RFC §14 S1 全部 exit；禁止：切换 whois 生产流程或“部分通过”冒充 S1 完成。

### P1.7 产品闭环与生命周期（在 P1.6 发布验收前执行）

编号保持稳定，实际依赖顺序为 T016 → T019–T024 → T017 → T018。以下全部待 S0 授权，责任路径是拟建位置；每个任务可再按一个失败反例拆小，不一次生成完整功能。

- [ ] T019 [Plan:P1.7] PC-01：实现 `internal/taskdef/preview.go`、`internal/console/preview.go` 与测试、无 AI 示例；依赖：T016；验收：AT-16；禁止：预览启动命令/网络/模型/自动安装。
- [ ] T020 [Plan:P1.7] PC-02：实现 `internal/snapshot/export.go`、`internal/evidence/delivery.go` 与测试；依赖：T019、T006/T011；验收：AT-17；禁止：导出候选冒充接受、覆盖目标、源树/Git 写入。
- [ ] T021 [Plan:P1.7] PC-03：实现 `internal/chain/authorization.go`、`internal/console/approvals.go` 与测试，接 guard 停机；依赖：T020、T007/T011；验收：AT-18；禁止：自授权、硬门禁 waiver、过期自动批准。
- [ ] T022 [Plan:P1.7] PC-04：实现 `internal/gates/effects.go`、`internal/evidence/diagnostics.go` 和 chain 恢复计划测试；依赖：T021、T010；验收：AT-19；禁止：外部写 runner、未知副作用盲重试、诊断删锁/日志。
- [ ] T023 [Plan:P1.7] PC-05：扩展 `internal/tickets/budget.go`、`internal/adapters/` 用量和 `internal/console/` 成本报告；依赖：T022、T012/T013；验收：AT-20；禁止：重启释放未知费用、跨 run 绕上限、默认遥测。
- [ ] T024 [Plan:P1.7] PC-06：实现 `internal/snapshot/backup.go`、`internal/evidence/disposition.go`、生命周期控制/测试及发布清单工具；依赖：T023；验收：AT-21、SBOM/支持说明；禁止：未知格式迁移、默认删证据/共享工具、擅自签名发布。

## 3. 需求映射

Implementation evidence 是计划证据，当前不代表实现。

| REQ ID | Plan Items | 任务 | 验收/实现证据 |
|---|---|---|---|
| REQ-001 | P0.1,P0.2,P1.6 | T001,T002,T003,T018 | AT-01/15；通用 Schema/非 C 样例 |
| REQ-002 | P1.2,P1.6 | T009,T018 | AT-03/15；chain/报告 |
| REQ-003 | P1.2 | T009 | AT-03；有序调度 |
| REQ-004 | P0.2,P1.1,P1.2 | T002,T006,T008 | AT-02/04；snapshot/applier |
| REQ-005 | P0.2,P1.3 | T002,T011 | AT-05；review/publish |
| REQ-006 | P0.3,P1.4,P1.6 | T004,T013,T017 | AT-08/14；adapter/CI |
| REQ-007 | P0.2,P0.3,P1.4 | T002,T004,T013 | AT-08；双通道 |
| REQ-008 | P0.2,P1.4 | T002,T013 | AT-08；context 恢复 |
| REQ-009 | P0.3,P1.5 | T004,T016 | AT-12；TUI |
| REQ-010 | P0.2,P1.5 | T002,T016 | AT-01/12；配置解释 |
| REQ-011 | P1.5,P1.6 | T015,T018 | AT-11/15；三 harness |
| REQ-012 | P0.2,P1.2 | T002,T007,T010 | AT-06；gate runner |
| REQ-013 | P0.2,P1.5 | T002,T015 | AT-11；模板生成 |
| REQ-014 | P0.1,P1.3,P1.6 | T001,T012,T018 | AT-07/15；阶段报告 |
| REQ-015 | P0.3,P1.5,P1.6 | T004,T016,T017 | AT-12/13/14；构建/TUI |
| REQ-016 | P1.4,P1.5 | T013,T016 | AT-08/12；无 IDE CLI |
| REQ-017 | P0.1,P0.2,P0.3 | T001,T002,T003,T004 | AT-01；readiness |
| REQ-018 | P0.2,P1.2,P1.5 | T003,T009,T016 | AT-01/03；四 kind |
| REQ-019 | P1.6 | T017 | AT-13/14；两代自托管 |
| REQ-020 | P0.2,P1.5 | T003,T015 | AT-01/11；组件/双语言 |
| REQ-021 | P0.2,P0.3,P1.4 | T003,T004,T014 | AT-09/10；handoff |
| REQ-022 | P0.2,P1.3,P1.5 | T003,T011,T015 | AT-01/05/11；文档协同 |
| REQ-023 | P0.1,P0.3,P1.1,P1.2,P1.3,P1.4 | T001,T004,T005,T007,T010,T011,T012,T014 | AT-05/06/07/09；安全 |
| REQ-024 | P0.2,P0.3,P1.1,P1.2,P1.3,P1.4,P1.6 | T002,T004,T005,T006,T008,T009,T011,T014,T017 | AT-02/03/04/09/13；恢复 |
| REQ-025 | P0.2,P1.1,P1.3,P1.4,P1.6 | T002,T005,T011,T013,T014,T017 | AT-04/06/10/14；证据 |
| REQ-026 | P0.1,P1.2,P1.3,P1.4 | T001,T009,T010,T012,T013 | AT-03/06/07；预算 |
| REQ-027 | P0.3,P1.5 | T004,T016 | AT-12；终端/JSON |
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

RFC §19 的后续候选：S2 图形上手/终端本地化、脱敏通知/远程审批身份设计、外部 PR 集成评估、外部写补偿、成本趋势和显式存储迁移；S3 商业后台/云管理须先验证需求。这些不在 T019–T024 的授权范围，后续再拆任务。当前共 10 个 plan items、24 个未完成任务、21 组计划验收；不是完成率。