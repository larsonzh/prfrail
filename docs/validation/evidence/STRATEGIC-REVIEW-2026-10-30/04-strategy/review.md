引用核对已完成，正在收束八问的结论、替代路径和验收判据。未能从指定输入证实的时长、资产比例及五轮闭环状态，会在最终报告中明确保留。


读取范围：[docs/RFC-proofrail-unattended-ai-engineering-product.md](docs/RFC-proofrail-unattended-ai-engineering-product.md)（全文，1–2446 行）；[docs/PRODUCT_REQUIREMENTS.md](docs/PRODUCT_REQUIREMENTS.md)；[docs/DEV_PLAN.md](docs/DEV_PLAN.md)；[README.md](README.md)；[docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md)（仅 OB-21 及其更正上下文）。

事实来源仅为读取主控落盘的上述文件；“完成、通过”均指落盘记载，不代表本轮复验。未读取前三轮结论。RFC 自身已说明它是设计来源、不是全域唯一权威，因此本轮采用“有记录、有理由、改进得到验证”评价演进，不以偏离旧设计直接判错。[RFC-proofrail-unattended-ai-engineering-product.md:13](docs/RFC-proofrail-unattended-ai-engineering-product.md#L13)

**S1 出口条件**
结论（判断，High）：真实 AI 候选端到端验收仍应保留为 S1 最终出口；应调整验收拆分和推进方式，不应将它改成可选。否则“AI 无人值守工程”尚未证明，就宣布其 MVP 成立，与现行正式执行要求不符。[PRODUCT_REQUIREMENTS.md:69](docs/PRODUCT_REQUIREMENTS.md#L69)

事实：OB-21 登记的是装配缺口和合法候选缺口，DR-1 修复不等于这两项解除；B4 明确区分“机制腿完成”与“AT-23 未通过”；README 仍记载执行型任务链及正式 AI 执行未交付。[DELIVERY_DIRECTIVE.md:1035](docs/DELIVERY_DIRECTIVE.md#L1035)、[DEV_PLAN.md:143](docs/DEV_PLAN.md#L143)、[README.md:29](README.md#L29)

替代路径（建议）：将 AT-23 拆为机制验证、强制边界及生产装配、真实候选执行与计量三个可独立举证的验收单元，全部满足才通过总门。机制单元不能直接沿用 B4 的非等价替代来声称覆盖完整。另可设离线技术预览里程碑，但只能按实际可用命令验收，不能称为 S1 完成。[DEV_PLAN.md:143](docs/DEV_PLAN.md#L143)、[README.md:25](README.md#L25)

性质与记录：这是**重新定义验收颗粒度，不是降低门槛**；若删除真实候选要求却保留原 S1 名称和承诺，就是降低现有门槛。应由所有者书面决定，记录旧新门槛映射、理由、风险、预算及反证条件，同步需求、任务、验收和 OB-21；不得仅靠改状态词完成。[RFC-proofrail-unattended-ai-engineering-product.md:990](docs/RFC-proofrail-unattended-ai-engineering-product.md#L990)、[DEV_PLAN.md:205](docs/DEV_PLAN.md#L205)

**S2 提前项**
结论（判断，Medium）：扩展能力的部分基础已经提前发生，但不足以认定 S2 已经完成或阶段定义整体失真；真正需要纠正的是把“基础、实验、条件完成、正式支持”混为一类。[RFC-proofrail-unattended-ai-engineering-product.md:979](docs/RFC-proofrail-unattended-ai-engineering-product.md#L979)、[DEV_PLAN.md:209](docs/DEV_PLAN.md#L209)

事实一：**Go harness 明确从 S2 提前至 S1**，理由是支持自托管；T015 记载 generic/C/Go、C+Go 作用域和整体恢复验证完成，T017 另有两代自托管记载。这项具有提前记录、明确理由和工程验收依据，可视为工程层的合理演进。[RFC-proofrail-unattended-ai-engineering-product.md:561](docs/RFC-proofrail-unattended-ai-engineering-product.md#L561)、[DEV_PLAN.md:85](docs/DEV_PLAN.md#L85)、[DEV_PLAN.md:90](docs/DEV_PLAN.md#L90)

事实二：**多组件/多语言及文档协同的底座**已经形成，但它们本来就被划为 S1 基础，并非完整 S2 偷跑；**Linux AgentRunner 相关机制**已有 C1 原生验证，但明确排除部分能力、不改变 Windows 阻断结论；**黑箱监督候选**已有 T028 条件闭合，但真实 CLI 端到端仍依赖 T027-B，不等于 S2 的完整实时监督协作。[PRODUCT_REQUIREMENTS.md:82](docs/PRODUCT_REQUIREMENTS.md#L82)、[DEV_PLAN.md:149](docs/DEV_PLAN.md#L149)、[DEV_PLAN.md:113](docs/DEV_PLAN.md#L113)

判断：Python 等新增 harness、正式三语言产品任务、完整实时监督、Web、智能影响助手及模型策略，仍应按正式交付证据判定，不能从上述底座推导完成。阶段台账宜增加成熟度字段，而不是为了维护阶段“纯度”撤回已有成果。[DEV_PLAN.md:209](docs/DEV_PLAN.md#L209)、[README.md:31](README.md#L31)

**S3 反哺排序**
结论（判断，Low）：按 RFC §16.2 明确列出的封装三层，当前反哺优先级应为：**核心 CLI ＞ VS Code 扩展 ＞ Web 控制台**；这是当前阶段的投入排序，不是永久价值排序。[RFC-proofrail-unattended-ai-engineering-product.md:1076](docs/RFC-proofrail-unattended-ai-engineering-product.md#L1076)、[RFC-proofrail-unattended-ai-engineering-product.md:1077](docs/RFC-proofrail-unattended-ai-engineering-product.md#L1077)、[RFC-proofrail-unattended-ai-engineering-product.md:1079](docs/RFC-proofrail-unattended-ai-engineering-product.md#L1079)

排序理由：CLI 是唯一必需交付物，且真实执行闭环仍未交付，首先应吸收合法启动、恢复、接受和成本解释方面的经验；VS Code 层其次，重点是把实际操作摩擦转成核心可复现的问题，而不是增加独占流程；Web 最后，先复用既有事实和控制 API，不以新界面替代核心能力。[README.md:29](README.md#L29)、[PRODUCT_REQUIREMENTS.md:78](docs/PRODUCT_REQUIREMENTS.md#L78)、[RFC-proofrail-unattended-ai-engineering-product.md:832](docs/RFC-proofrail-unattended-ai-engineering-product.md#L832)

判据：优先采纳能减少每个已接受任务的人工介入、降低总成本或提高恢复成功率的反馈。若后续同任务集证明 Web 或扩展的某项改进收益更高，可以上调；目前不能凭界面丰富度作此判断。这一排序不构成提前立项新界面的授权。[PRODUCT_REQUIREMENTS.md:27](docs/PRODUCT_REQUIREMENTS.md#L27)、[PRODUCT_REQUIREMENTS.md:130](docs/PRODUCT_REQUIREMENTS.md#L130)

**S4 战略回看**
结论（判断，Medium）：五轮独立意见是战略回看的输入，不是其成立条件；指定输入包不足以证明 P35 已形成决策闭环，因此本轮不能签认“战略回看完成”。阶段报告要求能力、未决项和下一步，任务完成也要求实际验收，而非报告数量。[RFC-proofrail-unattended-ai-engineering-product.md:982](docs/RFC-proofrail-unattended-ai-engineering-product.md#L982)、[DEV_PLAN.md:205](docs/DEV_PLAN.md#L205)

建议验收标准一：形成一张 **8/8 问题覆盖矩阵**，每问至少一个可定位事实、一项明确判断及其不确定性；对“继续等待、降为可选、重新定义并拆解”至少三条路线，给出收益、风险、资源和放弃理由的比较，不能只汇总赞成意见。[PRODUCT_REQUIREMENTS.md:111](docs/PRODUCT_REQUIREMENTS.md#L111)、[DEV_PLAN.md:211](docs/DEV_PLAN.md#L211)

建议验收标准二：**每项 High 判断至少一个能推翻它的检查**；每项采纳动作都具有责任主体、费用/时长上限、产物、通过/失败条件和复议触发点。最后形成一份所有者决策记录，逐问注明采纳、拒绝或待证，并把采纳项落实到任务和门槛映射中；可合入现有记录，不要求新增模板。[DEV_PLAN.md:15](docs/DEV_PLAN.md#L15)、[DEV_PLAN.md:211](docs/DEV_PLAN.md#L211)、[RFC-proofrail-unattended-ai-engineering-product.md:990](docs/RFC-proofrail-unattended-ai-engineering-product.md#L990)

本片对照：本轮可以提供八问意见，但其他轮次的独立性、汇总采纳、所有者裁定及后续验证结果，我不知道；本轮也按要求未读取前三轮。故判定是**闭环尚不可证**，不是断言相关记录不存在，更不能用“五份评审齐全”替代上述验收。[DEV_PLAN.md:205](docs/DEV_PLAN.md#L205)

**S5 模板与反哺**
结论（判断，Medium）：已经出现“工程过程成熟度领先用户价值验证”的失衡，但现有证据不能证明是模板本身抑制了反哺；不支持因此撤除必要安全纪律。[README.md:29](README.md#L29)、[PRODUCT_REQUIREMENTS.md:130](docs/PRODUCT_REQUIREMENTS.md#L130)

事实：反哺并非没有发生。T015 把文档影响规则、生成物新鲜度等经验转成产品检查；后续切片把 spawn 后回收、PID 报告缺失的三态判定继续转成生产行为。这些是能力改进，不只是增加清单。[DEV_PLAN.md:85](docs/DEV_PLAN.md#L85)、[DEV_PLAN.md:147](docs/DEV_PLAN.md#L147)、[DEV_PLAN.md:151](docs/DEV_PLAN.md#L151)

反面事实：T027 已积累多轮局部成果，但 enforcement 生产落地仍未立项；正式任务链未交付，需求文件仍明确收益无实测结论。因此当前不足主要是**局部经验没有汇成可用闭环及价值证据**，不能仅归因为“纪律太多”。[DEV_PLAN.md:145](docs/DEV_PLAN.md#L145)、[README.md:29](README.md#L29)、[PRODUCT_REQUIREMENTS.md:130](docs/PRODUCT_REQUIREMENTS.md#L130)

建议：下一项流程改进应至少把一种重复人工处置转成可复用的检查、默认行为或可解释诊断，并对比人工介入或漏检情况；未产生这种结果，就不要把新增文档量算作产品反哺。RFC 已明确警示规则复杂度，并要求以安全完成目标而非条目数评价产品。[RFC-proofrail-unattended-ai-engineering-product.md:236](docs/RFC-proofrail-unattended-ai-engineering-product.md#L236)、[RFC-proofrail-unattended-ai-engineering-product.md:2446](docs/RFC-proofrail-unattended-ai-engineering-product.md#L2446)

**S6 T027 破局**
结论（判断，High）：存在非“等待外部条件”的内部推进路径：将已选定的外部 OS 强制方案接入生产路径；但这不是宣称当前已有合法候选。落盘记载已区分“边界实测”与“产品未接入”，T027-B 仍只是待立项概念。[DEV_PLAN.md:112](docs/DEV_PLAN.md#L112)、[DEV_PLAN.md:135](docs/DEV_PLAN.md#L135)、[DEV_PLAN.md:145](docs/DEV_PLAN.md#L145)

路径第一步：由 Sol/所有者处理已登记的边界冲突并裁定 T027-B 立项，限定 Windows 首发范围、既定方案 B、容器内工具执行路径和生产 adapter 接线，冻结资源上限。本轮不代替该裁定，也不另选架构；该立项和缺口冻结动作本身无需新增模型请求。[DEV_PLAN.md:112](docs/DEV_PLAN.md#L112)、[DEV_PLAN.md:145](docs/DEV_PLAN.md#L145)、[DEV_PLAN.md:15](docs/DEV_PLAN.md#L15)

路径第二步：先满足目标部署前提内的持久发布证明，再用确定性候选贯通**实际生产入口**，验证合法执行及独立接受路径，并覆盖超时、日志缺失、未知恢复、退出零后重扫失败等拒绝路径；验收必须包括源目录不变、失败候选不传播、进程停止可证明。仍为 `unproven` 时只能得到阻断证据，不能拿拒绝路径成功冒充启动路径交付，也不能外推实验环境的耐久结论。[DEV_PLAN.md:141](docs/DEV_PLAN.md#L141)、[DEV_PLAN.md:143](docs/DEV_PLAN.md#L143)、[PRODUCT_REQUIREMENTS.md:103](docs/PRODUCT_REQUIREMENTS.md#L103)

路径第三步：合法启动路径成立后，再申请有明确上限的真实候选实验，固定版本、配置和输入，验证产物、日志、用量、独立 gates/review 与恢复；未通过则归属到具体阻断项，不追加无边界重试。与此同时，T018 的 EOL、包绑定复验等非同源门可以分项准备，但最终 T018 和 S1 签署不能绕过 T027。[PRODUCT_REQUIREMENTS.md:29](docs/PRODUCT_REQUIREMENTS.md#L29)、[DEV_PLAN.md:15](docs/DEV_PLAN.md#L15)、[DEV_PLAN.md:95](docs/DEV_PLAN.md#L95)

**S7 长期未决**
结论（判断，High）：三选一，我选择 **重新定义并拆解**；重定义的是验收单元、阻断归属和复议机制，不是取消真实 AI 验收，与 S1 保留最终出口一致。[DEV_PLAN.md:143](docs/DEV_PLAN.md#L143)、[PRODUCT_REQUIREMENTS.md:29](docs/PRODUCT_REQUIREMENTS.md#L29)

事实边界：“半年”无法由指定输入包证实。可定位的需求评审日期是 2026-08-28，实施计划至少在 2026-09-11 已记载候选不兼容、AT-23 未通过。本轮不把“半年”作为已证事实，但长期阻断需要处置的判断并不依赖这个时长精确成立。[RFC-proofrail-unattended-ai-engineering-product.md:296](docs/RFC-proofrail-unattended-ai-engineering-product.md#L296)、[DEV_PLAN.md:3](docs/DEV_PLAN.md#L3)

选择理由：继续等会忽略尚未完成的内部装配工作；降为可选则撤掉了“确定性 fake 与真实 CLI 通过同一契约及恢复用例”的关键验证。拆解可以保留已获得的机制证据，同时暴露尚未成立的生产和候选假设。[DEV_PLAN.md:145](docs/DEV_PLAN.md#L145)、[PRODUCT_REQUIREMENTS.md:29](docs/PRODUCT_REQUIREMENTS.md#L29)

处置要求：下一次评估应分别裁定“生产接线切片”和“真实候选验收切片”的范围及预算，给 OB-21 增补责任主体、下一项区分性实验、停止条件和明确复议触发点。两次同类失败无新证据就停止该实验路线并提交决策，不因日历推进自动消耗更多探针额度。[DELIVERY_DIRECTIVE.md:1035](docs/DELIVERY_DIRECTIVE.md#L1035)、[DEV_PLAN.md:15](docs/DEV_PLAN.md#L15)

**S8 特殊与通用**
结论（判断，Medium）：当前更准确的描述是**通用机制已经抽取，通用产品成效尚未充分证明，专用实验与处置经验的成熟度领先实际复用闭环**；不能据此给出可信的百分比。[DEV_PLAN.md:85](docs/DEV_PLAN.md#L85)、[README.md:29](README.md#L29)、[PRODUCT_REQUIREMENTS.md:130](docs/PRODUCT_REQUIREMENTS.md#L130)

事实：通用一侧已有证据、快照、守卫、事务、任务链及 harness 等完成记录；专用一侧仍有固定 Copilot CLI、Windows AppContainer、特定断电前提等实验资产。whois 影子比对有结果，但不能直接等价为其他项目的真实 AI 使用收益。比例究竟是几比几，我不知道：本包没有统一分类口径、投入统计和跨项目收益对照。[DEV_PLAN.md:49](docs/DEV_PLAN.md#L49)、[DEV_PLAN.md:75](docs/DEV_PLAN.md#L75)、[DEV_PLAN.md:133](docs/DEV_PLAN.md#L133)、[DEV_PLAN.md:139](docs/DEV_PLAN.md#L139)、[DEV_PLAN.md:93](docs/DEV_PLAN.md#L93)

转化方向（建议）：把“某 CLI 的绕过现象”转成独立于供应商的能力要求和反例；把“某虚机的耐久实验”转成携带平台及前提的判定规则；把操作纪律转成外置 harness/profile、检查器或解释性诊断，而不是让使用者继承整个项目台账。这与核心不感知语言、编辑器和提供方的目标一致。[DEV_PLAN.md:111](docs/DEV_PLAN.md#L111)、[DEV_PLAN.md:141](docs/DEV_PLAN.md#L141)、[RFC-proofrail-unattended-ai-engineering-product.md:119](docs/RFC-proofrail-unattended-ai-engineering-product.md#L119)

可检验标准（建议）：至少用两个非 whois、不同 harness 的项目，在不增加项目专名核心分支的情况下完成受控变更、独立接受和失败恢复，并记录人工介入与每个已接受任务成本。达到前，应称“可复用机制/方法”，不称“已证明普适有效”；方法论价值也不能由模板数量替代。[PRODUCT_REQUIREMENTS.md:41](docs/PRODUCT_REQUIREMENTS.md#L41)、[PRODUCT_REQUIREMENTS.md:27](docs/PRODUCT_REQUIREMENTS.md#L27)、[PRODUCT_REQUIREMENTS.md:130](docs/PRODUCT_REQUIREMENTS.md#L130)

未展开项清单：无。  
本轮自评：事实/判断比例约 4:6；事实限于主控落盘记载，排序、替代路径及新增验收标准均为判断。
