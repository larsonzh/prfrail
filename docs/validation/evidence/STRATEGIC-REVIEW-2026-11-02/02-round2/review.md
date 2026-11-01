# 轮 2 · 缺口台账（原样归档，含两片）

> 载体：runSubagent(agentName="deep-reasoner", model="DeepSeek V4 Pro (deepseek)")。输出 8501 字符，主动分两片（批次 A GAP-01–33 / 批次 B GAP-34–76），片界标注原样保留。

本片为分片 1/2，总计 2 片（批次 A：§12.4 非终态 24 + §六 待裁 9，GAP-01–GAP-33）

| 编号 | 来源 | 缺口一句话 | 影响面 | 前置条件 | 规模 |
|---|---|---|---|---|---|
| GAP-01 | §12.4 OB-1 L1017（观察中） | ⑥ 在纯文档切片上的有效性无结论性评估（n=2 归因未定） | §7.2 纯文档豁免条款无法定案 | 后续纯 .md 切片 ⑤ 预审时观察，积累 2–3 片样本 | S |
| GAP-02 | §12.4 OB-7 L1023（观察中） | ⑥/⑦ 独立重叠仅单点正向信号，无系统性证据 | 「⑥ 独立价值」论证数据点不足 | 更多首轮不附清单的对照切片 | S |
| GAP-03 | §12.4 OB-10 L1026（观察中） | §1.6 两项新条款无机械判据（停止态报告头标注 / [PILOT] 指标槽） | 同类时态不一致仅靠人工目视（本片 3 次复发实证） | 下次 §6.3 修订开启（拟提升必做） | M |
| GAP-04 | §12.4 OB-33 L1057（观察中） | 跨位置一致性审查三层架构未建成 | 同一事实多处写、漏同步无机械覆盖 | 专门片评估三层架构 | M |
| GAP-05 | §12.4 OB-87 L1111（观察中） | ⑤ 位成本口径无准则依据（仅 n=2 手动数据点） | 治理/评估片 ⑤ 预算无法事前估算 | 后续片在 §12.5 缺口行给出成本口径 | S |
| GAP-06 | §12.4 OB-99 L1123（观察中） | 时序类临时件写入（仓库外/变量中转）在 CI 输入上零信息 | G9 抓不到，依赖人工纪律（OB-98 已证规则文本无效） | 再犯时评估 AST 污点追踪或 ⑧c 记录态收口 | M |
| GAP-07 | §12.4 OB-100 L1124（观察中） | 创建位置违规但未提交的临时件不被 G9 检测 | G9 只守提交形态，其余靠 §5.1 三段处置架 | ⑧c 后 ?? 项再次 >0 时评估 G9-d | M |
| GAP-08 | §12.4 OB-104 L1128（观察中） | G3(a.2) 第七/八形状未实现判据（漂移率 26%） | 双语镜像形状漂移无机械覆盖 | 下一轮 GATES-* 片由 ① 定形态 | M |
| GAP-09 | §12.4 OB-105 L1129（观察中） | gate.js `--list` 与 `--check` 判据 id 集不一致（base 不可单查） | 调用方无法据 `--list` 枚举可单独选取的判据 | 下个 GATES-* 片修入口一致性 | S |
| GAP-10 | §12.4 OB-107 L1131（观察中） | ③ 派单顺序无机械检查点（本片 ③ 晚于 ⑤⑥⑦ 补派） | ⑥/⑦ 输入包不含 ③ 测试结果 | ⑩ 回执加顺序自查；再犯评估机械判据 | S |
| GAP-11 | §12.4 OB-111 L1135（观察中） | 白名单门槛/替代词表流程化留待后续 GATES-*（维持现状） | 每次触发手工找替代词、可能引入新未知字 | 下次触发时 ① 评估或 GATES-* 同批 | S |
| GAP-12 | §12.4 OB-113 L1137（观察中） | 行状态滞后：G4b 缓存已由 G4b-CACHE 落地，判据异步化收益未验 | 台账时态与事实可能不一致；`--all` 全量仍 1–2 分钟 | 核对 G4b-CACHE 后收口本行 | S |
| GAP-13 | §12.4 OB-19 L1035（待议） | S2 设计输入四条实测验证机制未定议 | S2 设计启动前缺决策输入 | S2 设计启动前 + Sol 战略回看问题集 | M |
| GAP-14 | §12.4 OB-21 L1039（待议） | AT-23 真实 AI 候选腿不可达（CLI 未装配 AgentRunner、无可启动候选） | S1 exit 判读被阻塞 | 下次 S1 exit 评估 + 候选/装配进展 | M |
| GAP-15 | §12.4 OB-52 L1076（待议） | CONTRACTS:191 额度要求与 CLI 实际额度不一致（记录不含费用维度） | 可用性结论与产品请求预算脱节 | 另立切片修 CONTRACTS 语义 | M |
| GAP-16 | §12.4 OB-53 L1077（待议） | 成本上限重估 + 链配置 `model:"auto"` 处置未落 | auto 可间接引入黑名单模型 | 用户排片 | M |
| GAP-17 | §12.4 OB-55 L1079（待议） | enforcement 真实边界接入产品的 B 段归属未定 | T027 仅做拒绝路径，真实边界无落点 | B 段立项时 ① 评估（独立片/提前 Sandbox/S2） | M |
| GAP-18 | §12.4 OB-57 L1081（待议） | dispatcher spawn 后身份记录失败即返回且不回收进程 | 孤儿进程风险仅 console 层 best-effort | 另立切片修 dispatcher 根因 | M |
| GAP-19 | §12.4 OB-58 L1082（待议） | 行未收口：实质（pid 摘要入冻结事实）已由 v1.42 落地 | 台账时态滞后于事实 | 批量清账片收口 | S |
| GAP-20 | §12.4 OB-59 L1083（待议） | 真实 CLI 生产侧恒定写子 pid 报告契约未落（stub 写侧已修） | 无子进程时生产正向链判 uncertain | 另立切片补 CLI 契约 | S |
| GAP-21 | §12.4 OB-61 L1085（待议） | 不变量测试靠文本计数，可被等价改写绕过 | 源码级回归护栏可失效 | 评估 AST 级结构断言收益后立项 | M |
| GAP-22 | §12.4 OB-75 L1099（待议） | 限定名静默路由盲区无 operator 侧防护 | 调用被路由到其它模型时规则层无感知 | 下一片评估或与模型可观测性同批 | S |
| GAP-23 | §12.4 OB-120 L1144（待议） | G4b 语料枚举与索引耦合（tree/index 读索引） | 暂存删除后冷跑假阳性、缓存/冷跑判定不一致 | 实际误判或需 hit≡cold 保证时启动 | S |
| GAP-24 | §12.4 OB-69 L1093（待办） | §2.5 时态整改并入 GATES-PATCH ⑦ 输入的承诺无收口 | 整改独立复审覆盖存疑 | 核对 GATES-PATCH ⑦ 输入后收口 | S |
| GAP-25 | CHECKLIST_v35 §六 #0 L197 | P46 战略回看是否立项待裁 | 队列 0（本次回看本身）被阻塞 | 用户裁定 | S |
| GAP-26 | CHECKLIST_v35 §六 #1 L198 | `CONTRACTS-S7-DURABILITY` 立项节奏待裁 | C6 的 Windows 落地被阻塞 | 用户排片 | M |
| GAP-27 | CHECKLIST_v35 §六 #2 L199 | P45 后续片立项归属待裁（console 多任务装配 / Race 扩包 / `registry.len()` 断言） | 后续装配与验证片被阻塞 | 用户排片 | S |
| GAP-28 | CHECKLIST_v35 §六 #3 L200 | DR-13「已修复+观察」口径是否延续待裁 | CI 稳定度口径 | 用户裁定 | S |
| GAP-29 | CHECKLIST_v35 §六 #4 L201 | DR13-FIX D3（先 remove 再 unlock）是否升格独立加固片待裁 | 跨进程并发加固被阻塞 | 用户裁定（技术缺口见 GAP-71） | M |
| GAP-30 | CHECKLIST_v35 §六 #5 L202 | OB-55 归属与 T027-B 立项节奏待裁 | 后续实施被阻塞 | 用户排片（同 GAP-17） | M |
| GAP-31 | CHECKLIST_v35 §六 #6 L203 | MVA、S1 移出项与 S2 承诺收窄批准待裁 | 后续切片被阻塞 | 用户裁定 | M |
| GAP-32 | CHECKLIST_v35 §六 #7 L204 | 候选与 Windows 前提 P 及资源/探针授权待裁 | G2 门被阻塞 | 用户裁定 | S |
| GAP-33 | CHECKLIST_v35 §六 #8 L205 | AT-24 `unknown` 与「未知副作用保持暂停」接受边界待裁 | AT-24 收口被阻塞 | 用户裁定 | S |

ROUND2: 分片 1/2

本片为分片 2/2，总计 2 片（批次 B：P39–P45 候选 36 + ⑥ 未覆盖面 1 + DR13-FIX 遗留 6，GAP-34–GAP-76）

> 排除注记：P44 review:32（SealBinding 三边界已 FR-02 补测）、P45 限制③/review:45（证据包 06/07/08 由 ⑨ 闭环）、P45 ⑦MINOR 与 D9（已闭合/边界外裁）不列为缺口。

| 编号 | 来源 | 缺口一句话 | 影响面 | 前置条件 | 规模 |
|---|---|---|---|---|---|
| GAP-34 | p39 §3E(:38)/§4(:44) | Windows 目录屏障维持 `unproven`，是否升级耐久等级未决 | C6 Windows 落地口径无定论 | `CONTRACTS-S7-DURABILITY` 立项（同 GAP-26） | S |
| GAP-35 | p39 §4(:45) 观察项一 | 持锁进程崩溃后的锁释放语义无测试腿 | 崩溃释放语义无证据 | 补测试腿 | S |
| GAP-36 | p39 §4(:45) 观察项二 | rename 成功后报错的可见分支无测试腿 | 错误分支无覆盖证据 | 补测试腿 | S |
| GAP-37 | p39 §4(:45) 观察项三 | `newDurableAck` 的 `confirmed` 字段命名与语义不符 | 契约可读性/语义漂移 | 重命名或补语义注记 | S |
| GAP-38 | p39 §4(:45) 观察项四 | 矩阵「确认前取消」腿未落地 | 该状态转换无覆盖 | 补矩阵腿 | S |
| GAP-39 | p40-lock-reaper.md:148 | 静态锁序断言未纳入 CI 判据（F-07a） | 锁序回归无机械守护 | go/analysis 或等价脚本断言 | S |
| GAP-40 | p40-lock-reaper.md:149 | `lockPath` 未先 canonical 化再比较（F-11） | 路径形态差异可绕过锁判 | 判据内补 canonical 化 | S |
| GAP-41 | p40-lock-reaper.md:150 | Reclaim 未接生产（uncertain 路径无调用者） | 锁回收能力无生产调用 | `agent_runner_pinned_cli.go` uncertain 路径接线 | M |
| GAP-42 | p40:151 / p44 R-4(:156/165) / p45:167 | CI Race 腿未覆盖 `internal/gates`、`guard`、`tickets`、`console` | 这些包并发结论无 -race 证据 | ci.yml 扩展；console 需 cgo 环境 | S |
| GAP-43 | p41-spawn-ownership.md:131 | kill-on-close 仅半证（I1：宿主崩溃注入+零残留 PID） | 回收断言缺宿主崩溃腿 | 崩溃注入测试 | M |
| GAP-44 | p41:132 | 四缝宿主崩溃注入未做（I2） | Start/attach/OnSpawn/resume 各缝无覆盖 | 逐缝注入测试 | M |
| GAP-45 | p41:133 | `(false,nil)` 拒绝态未三态化（F-5/⑦F-01） | 失败语义不可区分 | 哨兵错误或状态枚举 | S |
| GAP-46 | p41:134 | rename-away 对象按 path 删除有覆盖窗口（⑦F-02） | 窗口期误删风险 | 绑定对象删除 | S |
| GAP-47 | p41:135 | mirror 删除路径缺显式断言（⑦F-03） | 成功路径断言缺位 | 补断言 | S |
| GAP-48 | p41:136 | Linux 崩溃回收延迟未收敛到内核即杀（pdeathsig） | Linux 下回收依赖重启间隔 | 独立切片落地 pdeathsig | M |
| GAP-49 | p41:137 | starting 相位防双进程兜底强度未复核（⑤Low-1） | 双进程防护强度存疑 | 槽位新鲜度兜底评估 | S |
| GAP-50 | p42-evidence-bind.md:130 | 读路径 N2 绑定未迁移（hashLogs/readPIDs/fileHash/scanWorkspace） | 大目录读路径性能风险 | 逐项迁移 | M |
| GAP-51 | p42:131 | Unix 两段式 dirfd 未落地 | 中间组件符号链接向量未覆盖 | openat 逐级打开实现 | M |
| GAP-52 | p42:132 | Windows reparse 仍前缀比较，非对象身份级校验 | reparse 身份伪造窗口 | 身份级校验替代前缀比较 | M |
| GAP-53 | p42:133 | CloseBoundDirs 未接生产（run 收尾） | 句柄关闭无生产调用 | run 收尾接线+注释统一 | S |
| GAP-54 | p42:134 | darwin 交叉编译旧文档未 supersede | 文档与实现口径漂移 | 修订 t027-replay-convergence.md 条目 | S |
| GAP-55 | p42:135 | 原语层测试未补齐（关闭拒读/超界拒摘要/同源对抗） | 原语边界无测试护栏 | 补测试 | S |
| GAP-56 | p43-state-seal.md:142 | TaskAttemptResolver 未接生产（D3 边界） | attempt 解析缺可靠来源 | console 侧装配接线 | M |
| GAP-57 | p43:143 | 工作区命名固定 attempt-1（D4 边界） | attempt>1 场景无端到端覆盖 | supervised attempt>1 端到端 | M |
| GAP-58 | p43:144 | seals/ 目录级耐久无 crash 注入证据 | 目录项与事件日志同序落盘未证 | crash 注入测试 | M |
| GAP-59 | p43:145 | stale/mismatch 判定矩阵未补测 | 判定边界无表格用例 | 矩阵化补测 | S |
| GAP-60 | p43:146 | router/handoff 非 acceptance 路径缺 K 绑定 | 非最终路径 K 校验缺位 | 补齐 K 校验 | S |
| GAP-61 | p43:147 | WorkspaceReopener 接口未落地 | reopen 语义与 WorkspacePort 耦合 | 独立接口承接 | S |
| GAP-62 | p44 R-3(:155/164) | ReadinessInput 未注入 4 处调用点（C7/N4 未端到端） | 端到端闭合不可宣称 | runtime_agent/supervised/runtime/ai 接线 | M |
| GAP-63 | p44 R-2(:157/163) | `ProductionReadiness` 构造签名缺 L2 类型层绑定 | ① 类型层约束未恢复 | 构造签名接收 L2 值 | S |
| GAP-64 | p44:158/167 + P44 review:33/34 | LoadProbePermit 非 NotExist 读失败、Require ctx==nil、Complete 未知态、真跨进程双写未覆盖 | 门禁 fail 分支覆盖缺口 | ⑥ Section D 补测 | S |
| GAP-65 | p44:159 | AT-23 缺不可伪造签收工件 | L3 验收依赖自证 | 生成签收工件 | S |
| GAP-66 | p45:166/173 + P45 review:44 | 非 Windows OS/locale 并发差异未实测 | ProcessAlive/快照语义跨平台证据缺 | 跨 OS 环境验证 | M |
| GAP-67 | p45:168/172 + P45 review:43 | 真实 UI/浏览器路径无端到端点击验证 | 编辑器内运行路径未对齐 | 真实 UI 端到端测试 | L |
| GAP-68 | p45:175 + P45 review:46 | preview 命令执行证据未纳入自动化测试 | 预览路径仅手工补证 | 自动化测试纳入 | S |
| GAP-69 | p45:176(D7) + dr13-fix.md:165 | `registry.len()==1` 与 slot-owner/mirror 一致性无直接断言（跨 launcher 镜像失效窗口） | 镜像失效窗口无检测 | 补一致性断言 | S |
| GAP-70 | P44 06-scan review.md:35/36 | GoalProduction 非空非法 AT23Hash、GoalAcceptance 带值非法 SealBinding 未单独覆盖 | 门禁非法输入分支无测试证据 | 补测（可并入 GAP-64 批次） | S |
| GAP-71 | dr13-fix.md:160/171 + §六 #4 | release 仍 unlock 后 remove，路径消失晚于锁释放的毫秒窗口（D3） | 跨进程并发窗口残留 | 先 remove 再 unlock + Windows remove-while-open 全量回归 | M |
| GAP-72 | dr13-fix.md:161/172 | os.Remove 失败叠加跨进程扫描无告警/补偿路径（W1） | 跨进程 marker 残留不可观测 | 补可观测告警或补偿 | S |
| GAP-73 | dr13-fix.md:162 | StopAgentRunnerProcess 解析/停止失败时无条件 releaseLaunch | 失败分支释放语义不安全 | 按已证状态分流 | S |
| GAP-74 | dr13-fix.md:163 | W2 flock 内核语义未实证（关闭任一 fd 是否释放全部 flock） | unix flock 行为无证据 | unix 宿主实证 | S |
| GAP-75 | dr13-fix.md:164/169 | unix StopProcessIdentity 对活进程行为本机不可证（B1） | unix 分支判定无本机证据 | CI/unix 宿主取证 | S |
| GAP-76 | dr13-fix.md:170 | D2 无专属确定性测试（B2） | 镜像清理仅概率覆盖 | 确定性测试设计 | S |

## 台账汇总
共 76 项；其中 §12.4 非终态 24 / §六 待裁 9 / P39–P45 候选 36 / ⑥ 未覆盖面 1 / DR13-FIX 遗留 6。
ROUND2: DONE