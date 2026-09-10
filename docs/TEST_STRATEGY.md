# ProofRail 测试验证方案

[English](TEST_STRATEGY_EN.md)

日期：2026-09-07；测试与验收基线，非测试通过报告。本文是 AT、测试分层、平台和故障注入验收规则的规范性权威。输入：[需求](PRODUCT_REQUIREMENTS.md)、[契约](CONTRACTS.md)、[安全](SECURITY.md)；[项目建议书](RFC-proofrail-unattended-ai-engineering-product.md) §14、§17 保留阶段验收来源。

## 1. 基线与分层

应用类型：本地 CLI/TUI + 文件/进程编排。T005–T008 已有核心包 `_test.go`，但尚无可用 CLI 的产品 E2E；不能把包测试替代产品验收。规范入口为 `go test ./...`，新 Go 测试必须由该入口发现；构建/静态门禁另跑 `go build ./...`、`go vet ./...`。本轮安装的 Go 为 1.27.0；模块声明 1.22，兼容性仍须用声明的最低工具链另验。

| 层级 | 工具/数据 | 证明内容 | 执行时机 |
|---|---|---|---|
| 文档 | 离线链接、BOM/LF、中英 ID、REQ→任务→AT 检查 | 交付可导航且不遗漏范围 | 每次文档变更 |
| Schema | 独立 JSON Schema validator + 语义 fixture oracle | 正反样例在正确层接受/拒绝 | S0 冻结前及每次协议修改 |
| 单元 | Go testing、table-driven、fake clock/agent/runner | 状态转换、预算、checker、哈希 | 每个小任务先跑相关包 |
| 组件集成 | t.TempDir、真实文件、Go 子进程 helper | journal、路径、停机、恢复、租约 | 涉及文件/进程的任务 |
| CLI E2E | 已构建候选 + 临时 workspace/store + 固定 oracle | 三任务链及 US1-US6 | 迭代与 S1 exit |
| AgentRunner 契约 | 确定性 fake + 固定版本真实 CLI Agent 能力探针 | 启动/事件/工具循环/停止/恢复/用量与独立验收 | 离线契约后；真实收费调用需另授权 |
| 可见宿主集成 | SessionBridge v0.1.1 + 受控测试实例 | 辅助 silent 与 visible 黑箱投递/归还；不推断过程完成 | 离线契约后；收费调用需另授权 |
| 原生平台 | Windows 11 amd64；Linux amd64 核心 CI | 文件系统/信号/路径真实语义 | 每个候选版本 |
| 自托管 | seed N、candidate N+1、外部 oracle | 两代信任与 clean-room 重放 | S1 exit/发布 |

S0/S1 不需要数据库、Docker 服务或浏览器；不安装这些测试设施。infra 轴使用真实本地文件系统和子进程，不用内存 mock 替代原子写证明。browser 轴不适用，TUI 用终端探针；S2 新增 Web 后再引入浏览器 E2E。Windows race 可能需要 C 工具链，能力不符时在支持的 CI 执行并明确本地跳过，不与 CGO_ENABLED=0 发布混淆。

## 2. 验收矩阵

AT 编号代表测试组，不是已存在的测试函数。每组必须拆正/反例；输出用例数、PASS/FAIL/SKIP、版本、证据路径。

| ID | 需求 | 场景与断言 | 所有者/计划 |
|---|---|---|---|
| AT-01 | REQ-001/010/017/018/020/022/028 | Schema 合法配置；未知字段/版本、重复 ID、空 steps、非法 noop、重叠 target、依赖/生成环拒绝；C+Go/三语言/交接/文档样例 | taskdef；T002/T003 |
| AT-02 | REQ-004/024 | 含未提交文件树捕获；源树不变；保留名、大小写、symlink/reparse/hardlink、长路径、捕获中变动、配额不足 | snapshot；T006 |
| AT-03 | REQ-002/003/018/024/026 | 三任务混合 steps、noop 零进程；T2 失败不启 T3；暂停/取消/重启、失控子树与磁盘不足 | chain/guard；T007/T009 |
| AT-04 | REQ-004/024/025/028 | 事件缺失/重排、对象损坏、断尾、每个 journal/receipt/投影写点崩溃；完整接受才可重放 | evidence/applier；T005/T008 |
| AT-05 | REQ-005/022/023 | 拒评、自批、过期/错候选 waiver 拒绝；有效独立批准后一次发布；代码/文档原子接受 | chain；T011 |
| AT-06 | REQ-012/023/025/026 | argv 注入、环境泄密、越界 cwd、网络限制不可用、超时/无限输出/秘密扫描失败/缺产物均阻断 | gates/guard；T007/T010 |
| AT-07 | REQ-014/023/026 | 相同指纹耗尽、attempt/墙钟/费用耗尽、假修复、陈旧候选、Promote 中断；不重置预算 | tickets/repair；T012 |
| AT-08 | REQ-006/007/008/016 | IPC/文件队列同一业务结果；busy、错 requestId、旧结果、timeout、宿主重启、缓存丢失、历史截断；禁止 GUI 回退/重复应用 | adapters；T013 |
| AT-09 | REQ-021/023/024 | 人工成功归还、越界、冲突租约、离席/断线、abort/request-agent、崩溃；每次归还重跑 gates | chain/guard；T014 |
| AT-10 | REQ-021/025 | 未声明 prompt 停止；secret-direct 仅终端；无安全输入能力暂停；secret canary 不入证据/模型 | gates/console；T014 |
| AT-11 | REQ-011/013/020/022 | generic/C/Go、锁模板生成 A；自批/手改生成物/文档缺失阻断；C+Go 同 task 整体恢复 | taskdef/gates；T015 |
| AT-12 | REQ-009/010/015/016/027 | 无 IDE init/validate/run/report/config explain；键盘/无色/窄终端/--json；非实现命令不得冒充成功 | console/cmd；T016 |
| AT-13 | REQ-015/019/024 | seed 与 candidate 分离；外部 oracle 比较；clean-room 成功及失败不损 seed | harness/release；T017 |
| AT-14 | REQ-006/015/019/025/028 | Windows 原生/Linux 核心 CI；最低 Go/固定发布 Go；旧格式只读/未知版本拒写；固定 commit、SHA256SUMS、SBOM/许可证清单核验；校验和不冒充身份认证 | release；T017 |
| AT-15 | REQ-001/002/011/014 | whois 冻结夹具只读影子，同输入/结果/失败分类；小 Go 非 C 夹具闭环，不触碰正式仓库 | harness；T018 |
| AT-16 | REQ-009/010/016/027 | PC-01：无 AI 夹具预览；恶意 hook、版本探测、凭据/网络调用计数全为零；unknown 显示，源/store 不变，窄屏/JSON 同事实 | taskdef/console；T019 |
| AT-17 | REQ-002/004/025 | PC-02：已接受包离线核验；候选/秘密/缺对象/篡改/已有目标/别名重叠拒绝；写中断无完成回执，源树不变 | snapshot/evidence；T020 |
| AT-18 | REQ-005/021/023 | PC-03：撤销/过期/错摘要阻断下一动作和接受；在途进程停止失败保持暂停；硬门禁不可 waive，待审批重启可见 | chain/guard/console；T021 |
| AT-19 | REQ-012/023/024 | PC-04：外部写与无法执行限制拒绝，未知副作用不重投；诊断只读脱敏，恢复计划不删锁/改日志 | gates/chain/evidence；T022 |
| AT-20 | REQ-006/026 | PC-05：预留前后/投递后/结算点崩溃，超时保留余额占用，重复结算去重，两 run 竞争共享额度不超配；订阅模式不报精确账单 | tickets/adapters；T023 |
| AT-21 | REQ-019/025/028 | PC-06：闭包缺对象/活动写者/未知格式阻断；新 store 恢复验证，原件不变；卸载保留证据/共享工具，删除确认与保留冲突，离线撤销状态 unknown | snapshot/evidence/release；T024 |
| AT-22 | REQ-006/007/008/009/016/021/023/024/025/027 | 仅结构化请求可进入 WAITING_FOR_OPERATOR；TUI 展示允许响应并拒绝错主体/attempt/context hash、陈旧/重复响应；超时/断线/持久化失败保持暂停，重启重建待办；同 conversationId/新 requestId 恢复；人工写入走 handoff，secret-direct 不入模型/证据；没有 `@sbr-review` 仍可闭环 | adapters/chain/console/evidence；T025 |
| AT-23 | REQ-006/007/008/009/016/023/024/025/026 | fake AgentRunner 契约与真实固定 CLI 能力探针；隔离 workspace 工具迭代，记录版本/配置/session/process/events/logs/usage；超时、日志缺失、停机或恢复不明暂停；Agent 退出 0 后仍重扫并独立 gates/review | adapters/guard/chain/evidence；T026/T027 |
| AT-24 | REQ-007/009/021/023/024/025/026/027 | visible 黑箱模式启动前显示并持久化绑定风险确认；只在隔离 workspace，投递不等于完成；显式归还后全量 manifest/diff、范围/秘密/文件/副作用扫描及独立 gates/review；未确认、断线、未知进程/费用/外部副作用阻断，报告标记 reduced/unknown 且禁止自动降级 | console/adapters/chain/guard/evidence；T028 |

AT-16–AT-21 是 RFC §19 的 S1 最小验收，AT-22 是 AI/操作员交互，AT-23 是 AgentRunner 完整执行保证，AT-24 是受监督黑箱候选的降级保证。均先用无付费调用的确定性夹具；真实宿主/模型调用另授权。S0 T002/T003 只冻结记录和正反样例，不能将解析通过报成运行验证；T018 汇总全部 24 组结果。US7–US11 覆盖初次试用、交付、撤销、停用恢复和黑箱候选，收益指标用本地记录验证，不预设节省比例。

S2/S3 验收仍由 RFC §14 控制：三语言产品真实集成、supervised 断线、Web、更多 adapter 和生成脚本 B；不能用 S0 配置解析替代真实运行。

## 3. 数据、故障注入与证据

固定非秘密小夹具，用稳定 ID、内容与预期摘要；运行资源使用随机 run ID 和 t.TempDir 隔离。永不借当前仓库作为破坏性测试目标。whois 只读脱敏导出需用户授权；没有输入就标 SKIP/BLOCKED，不自造“等价通过”。

注入点包括：临时对象写前后、rename 前后、journal durable 前后、部分目标写入、review 后、publication receipt 后、事件后投影前、GC 引用更新和人工归还。每次重启断言源树不变、已接受引用不丢、不重复应用、无法证明则暂停。用 helper 子进程故意退出，不杀 IDE 或用户进程。

每次报告必须有工具链/OS/依赖版本、输入摘要、命令和 cwd、退出码、用例计数、失败证据、跳过原因、实际运行/仅设计区分。t.TempDir 自动清理；一次性人工产物放根 tmp 并用完删除。正式验收摘要保存到计划中的 `docs/validation/`，不可把待审原始秘密日志入库。

## 4. 门禁与排错顺序

每任务先相关包测试，再 build/vet/test 全门禁。协议改动额外跑独立 validator 和兼容 goldens；进程/路径改动跑原生平台；adapter 改动跑离线契约后实机。证据不全、任何强制安全失败、关键测试 SKIP 均不算阶段通过。

排错固定顺序：重现输入→首个失败层→错误/证据哈希→最小包/用例→一次小修复→原用例复验→全门禁。连续两次同类失败无新证据先停止调用模型，升级给维护者；不得扩大写范围、重置预算或关闭门禁取绿。

预算策略：先无 AI 离线 fixtures，再单次授权的实机 smoke；不以多个模型重复分析同一通过结果。性能方案沿用需求文档固定数据集，当前无性能结论。