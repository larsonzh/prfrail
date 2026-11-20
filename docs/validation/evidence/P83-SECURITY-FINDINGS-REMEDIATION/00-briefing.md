【P83 起跑 · SECURITY-FINDINGS-REMEDIATION · 用户 2026-11-20 点名】

〇 切片标识
- 编号 P83；标识 SECURITY-FINDINGS-REMEDIATION（承接 P79 的 F1–F8 修复）；
  分级 [SLICE]；类型 = 实现类（安全修复）；规模 = M-L（8 处位置、7 个文件 +
  可能触契约文本）；
- 步骤链：0 → ① → ② → ③ 一致性 → ④ 门禁 → ⑤ 预审 → ⑥ 独立扫描 →
  ⑦ 盲审 + 终审 + 复审 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩（盲审在终审前执行）；
- 主控 = DeepSeek V4.1 Flash (deepseek) High（实现负载在 ②③、判断负载在 ⑦，
  主控只编排 / 落盘 / 门禁 / 回执；不满足 §13.1.8 升级触发）；
- 档位限制（§8.3 已知限制）：②③ 与主控同为 deepseek-flash ⇒ 档位不可分离、
  整体 High，不得中途改档；
- 用户裁定（2026-11-20，同 P79 先例，两项）：
  a. ⑦ 额度 = **L 档 1 终审 + 3 复审**（超出 §7.5「M-L 视同 M」字面，系用户
     裁定、不得作为先例；第 2 次复审收窄输入、第 3 次只核上一轮 Medium+
     闭合路径）；超出 ⇒ 停下升级，不得自行加轮；本裁定须在报告「审查与
     计量」节与台账「额度口径披露」行如实标注（系用户裁定、不得作先例）；
  b. ⑦ 档位 = **GPT-6.2 Sol Max**（同轮取档授权，§8.2；§8.4 义务：报告
     「环境与档位」表记录设定时刻 / 档位值 / 是否暂停）。

一 背景与目标
- 输入 = P79 报告 §9 的 F1–F8（3 High + 5 Medium），逐字承接、不扩不缩：
  F1 面1 Medium internal/gates/runner.go:208 —— 产物读取绑定授权工作区目录
     句柄，无法安全打开即拒绝收集；
  F2 面2 High   internal/console/runtime_agent.go:1245 —— launch 前向可信
     账本读取授权与预留，读取失败 / 撤销 / 已结算均保留 R 并拒绝启动；
  F3 面2 High   internal/chain/review.go:116 —— 发布前从执行者不可写的可信
     来源验证策略批准与候选证据根；
  F4 面2 Medium internal/console/runtime_agent.go:731 与
     internal/console/runtime_supervised.go:1642 —— 发布前核对排他租约或
     等价写入隔离证明并绑定本次 attempt；
  F5 面2 Medium internal/gates/runner.go:73 —— 执行前从冻结 hook 策略导出
     有效副作用分类；
  F6 面3 High   internal/console/ai.go:94 —— 读取凭据前校验 profile、凭据
     引用、接收端与子进程身份绑定；
  F7 面3 Medium internal/guard/process_windows.go:457 —— 原生环境块与编码
     中间缓冲保留所有权，成功失败均清零；
  F8 面3 Medium internal/adapters/ai_probe_copilot_executor.go:53 —— 关闭
     未证明脱敏的子进程落盘日志，无法保证则拒绝 BYOK 探针。
- 优先级：F2 / F3（面 2 High）→ F6（面 3 High）→ F1 / F4 / F5 / F7 / F8。
- 协议先行（§1.4 硬规则）：F2/F3/F4/F5/F7/F8 均引用 CONTRACTS 句为被违反
  句 ⇒ ① 方案书必须逐条判定是否触契约语义；触则先改
  docs/CONTRACTS{,_EN}.md（及 schema / 夹具如涉）再动代码。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro (deepseek) Max —— 修复方案书
  （含契约影响判定 + 逐 F 修复设计 + 回归测试点）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash (deepseek)；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash (deepseek) —— 每 F 至少一条
  反例回归测试 + 变异检验（§13.1.2）；
- ⑤ 预审员：DeepSeek V4 Pro (deepseek)，额度 1 + 1；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，行为变更切片启用，额度 ≤3；
- ⑦ 独立终审：深度档 GPT-6.2 Sol (copilot) **Max**，额度 1 + 3（见 〇-a）；
- 盲审：F3/F6 触及身份语义 ⇒ 按 §7.7 **必抽 1 次**（另计 +1、不占终审/复审
  额度），终审前执行，理由与判定步骤 0 登记；
- 批量辅助：GLM-5.3-Flash (glm) ≥2 次、严格串行；仅限四类任务（枚举 /
  摘录核对 / 双语镜像 / MANIFEST 草案），禁择优；派单写明字符上限；
- 审查输入包消毒（§7.3）：凭据 / IP / 主机名 / token 一律过滤（F6/F8 涉
  凭据面，输入包与审查归档材料尤其须消毒）；
- 计量（R7.1）：审查用量以会话调用回执为单一真值源逐字段登记；
- 探针：默认 0；若 ⑨ 验证 F8 拒绝路径需真实 BYOK 调用，按 §7.6 停下申请
  预授权，不得擅自调用；
- 风险提醒：R 达 2 ⇒ 按 §9.2 高风险标记并停下升级，不自动继续。

三 环境声明
- Ubuntu VM 不开；网络代理不开（SSH push + gh 直连可用）；不做 gitee。
- ⑨ 原生验证：Windows 本机（无 gcc 不跑 -race）+ CI ubuntu 腿 Race 步骤；
  提交前 go build/vet/test ./... 全绿（§13.1.3）。

四 步骤 0 必做（§13.1）
- 锚点核对 8 项，出处 = docs/t027/CHECKLIST_v66{,_EN}.md 的 §一：
  ① HEAD = 7a63962；② origin/main = 7a63962；③ 准则版本 = v1.71；
  ④ SELFTEST = 740/740；⑤ 三作用域（index / tree / ci）= 全 TOTAL_FAIL=0；
  ⑥ v64 快照行数 = 224（LF 计数）；⑦ v65 快照行数 = 220（LF 计数）；
  ⑧ v66 快照行数 = 由 v67 复核；≥2 项不符 ⇒ 以仓库实况为准，先报终态
  读数再继续；
- 台账登记 P83 段（docs/t027/REMAINING_SLICES{,_EN}.md，字段结构同 P80 /
  P79 段）；R/D = 0 基线；
- 临时件落点 = tmp/P83-SECURITY-FINDINGS-REMEDIATION/（⑧c 前只增不删；
  片内 D 累计达 3 ⇒ 强制停片，§5.1 v1.71）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10）；每次 runSubagent 传 agentName + 限定名 model；
  ⑩ 派单表按 §13.2.2 逐字机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/validation/s2-pre-independent-security-review.md 及
   docs/validation/s2-pre-independent-security-review_EN.md（§9 F1–F8 表、
   §15 裁定留痕）；
2. docs/t027/CHECKLIST_v66.md 及 docs/t027/CHECKLIST_v66_EN.md
   （§四 P83 队列项）；
3. docs/CONTRACTS.md 及 docs/CONTRACTS_EN.md 的 F 项被引用句
   （:147 / :212 / :214 / :216 / :391，及 P79 报告 §9 表中 F4 行标注的行号）；
4. docs/SECURITY.md 的 :23 / :60；
5. docs/DELIVERY_DIRECTIVE.md 的 §1.4 / §5.2 / §5.3 / §7.5 / §7.7 /
   §8.2 / §9.2；
6. 七个改动文件现状：
   internal/gates/runner.go（:208、:73）、
   internal/console/runtime_agent.go（:1245、:731）、
   internal/console/runtime_supervised.go（:1642）、
   internal/chain/review.go（:116）、
   internal/console/ai.go（:94）、
   internal/guard/process_windows.go（:457）、
   internal/adapters/ai_probe_copilot_executor.go（:53）。

六 修复方案要点（① 方案书须覆盖，缺一退回）
- 逐 F 的改动点精确到 file:line、契约影响判定（触/不触 CONTRACTS + 理由）、
  回归测试点（反例 + 变异）、F1–F8 闭合判据（每 F 一条可证伪断言）；
- F3 的可信来源与 F6 的身份绑定：明确「执行者不可写」的实现载体；
- F8 的「拒绝 BYOK 探针」路径与日志脱敏判定的机械判据。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. 代码：runner.go / runtime_agent.go / review.go / runtime_supervised.go /
   ai.go / process_windows.go / ai_probe_copilot_executor.go + 对应测试；
2. 协议先行：若 ① 判定触契约语义 ⇒ docs/CONTRACTS{,_EN}.md 同步改
   （及 schema / 夹具如涉），先契约后代码；
3. 新增报告对 docs/validation/security-findings-remediation{,_EN}.md
   （§11.3 结构 + F1–F8 逐条闭合状态 + artifacts 围栏 + 双语严格镜像）；
4. 台账 docs/t027/REMAINING_SLICES{,_EN}.md P83 段 + 完成记录（F1–F8 逐条
   闭合 + R/D）；
5. docs/SECURITY{,_EN}.md:5 注记补「F1–F8 已于 <日期> 修复（P83）」状态句
   （P79 的 FINDINGS 结论句保留、后补修复状态）；
核对项（由用户裁定）：
6. docs/PRODUCT_REQUIREMENTS{,_EN}.md §7 的 P79 注记补修复状态；
7. ADR-047 保持 COMPENSATED；如加修复注记由用户裁定。
不动：
8. CHECKLIST_v66{,_EN}（v67 快照再吸收）；s1-exit / p57 门表；准则本体
   （不新增 OB / ADR 除非用户裁定）；CI；DEV_PLAN。

八 边界
- 只修 F1–F8，不扩大安全加固范围（额外发现登记候选 OB，不并入本片）；
- 不改 S1 Exit 八门判据；不改 SECURITY.md 威胁模型基线；不升准则版本；
  不做 gitee。

九 回执要求
- 首行 P83: STARTED；四要素（裁决人=用户、日期=2026-11-20、理由=P79 的
  F1–F8 修复义务、依据=本次起跑消息 + v66 快照 §四）；
- 锚点 8 项 + 模型 / 档位（Sol Max）/ 额度（1+3 裁定留痕）/ 盲审 / 批量辅助
  计划逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 单次输出 ≤ 10,000 字符（§13.1.16），超限分批、禁止自重启；
- 提交信息英文；完成后停于 ⑩，等同轮授权；仅 origin。
