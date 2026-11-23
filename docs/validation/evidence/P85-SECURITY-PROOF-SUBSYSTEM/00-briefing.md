# P85 起跑消息原文（§13.1.10 起跑消息即工件）

切片：P85 `SECURITY-PROOF-SUBSYSTEM`（`[SLICE]` · 实现类 L 档 · 协议先行）
主控：DeepSeek V4 Pro Max（用户 2026-11-22 授权切换；与 ① 统一 Max、
  思考档位起跑后固定；不留先例）
① 架构：DeepSeek V4 Pro Max（`deep-reasoner`）
② 实现 / ③ 测试：DeepSeek V4.1 Flash
⑤ 预审：DeepSeek V4 Pro（额度 1+1）
⑥ 独立扫描：MAI Code 2 Flash（额度 ≤3）
⑦ 终审：GPT-6.2 Sol Max（深度档；1 终审 + 3 复审 + 盲审必抽 1 另计；
  用户 2026-11-21 排片裁定深度档 + 2026-11-22 同轮授权 Max）
批量辅助：GLM-5.3-Flash ≥ 2（严格串行）

---

## 【〇·零】记忆唤醒（开场必读；本段为起草惯例红线，未登入准则 / OB / ADR）

你是 VS Code 主控（DeepSeek V4 Pro，Max 档）。本片是 P85，L 档实现类
`[SLICE]`，片名 = `SECURITY-PROOF-SUBSYSTEM`，承接 P83 顺延义务、闭合
F3 / F6 / F8（三项均为「证明载体」类能力）。步骤链首位 = ① 契约影响判定：
触 ⇒ 强制停下走协议先行（先改 CONTRACTS 双语 + schema / 夹具再动代码）；
不触 ⇒ 按方案书推进。盲审必抽 1 次（终审前执行）。额度与档位不得中途升档
（主控与 ① 已统一 Max、起跑后固定）。火力全开配置（主控 / ① 均 V4 Pro Max、
⑦ Sol Max）已获用户同轮授权；**不预设临时授权**——任何额度用尽按 §7.8.3
停下报用户。

**关键锚点（步骤 0 须实测核对）**：
- HEAD / origin = `e35e50f`；准则 = v1.71；SELFTEST = 740/740
- 三作用域 = `TOTAL_FAIL=0`；滚动窗口 = v67 / v68 / v69 = 216 / 218 / 209（LF）
- 当前快照 = **v69**；机械读数 ADR 47 / OB 139 / DR 18 / 试点顺延 25

---

【P85 起跑 · SECURITY-PROOF-SUBSYSTEM · 用户 2026-11-22 点名】

〇 切片标识
- 编号 P85；标识 SECURITY-PROOF-SUBSYSTEM（承接 P83 报告 §8 顺延清单
  P84-a / P84-b / P84-c，闭合 F3 / F6 / F8）；分级 [SLICE]；
  类型 = 实现类（策略制品 / 不可变安装来源 / BYOK 证明载体）；规模 = L；
- 步骤链：0 → ① 契约影响判定 → 协议先行（触则先改契约）→ ② → ③ 一致性 →
  ④ 门禁 → ⑤ 预审 → ⑥ 独立扫描 → ⑦ 盲审 + 终审 + 复审 → ⑧a → ⑨ → ⑧b →
  ⑧c → ⑩（盲审在终审前执行）；
- 主控 = **DeepSeek V4 Pro (deepseek) Max**（用户 2026-11-22 同轮授权切换、
  替代默认 Flash，并与 ① 架构统一 Max 档——思考档位起跑后固定、不得中途
  更改，§8.3；不留先例）；
- 档位限制（§8.3）：②③ = DeepSeek V4.1 Flash High；全片不得中途改档；
- 用户裁定（四笔）：
  a. ⑦ 额度 = 1 终审 + 3 复审（§7.5「L 档 / 治理片 / 高风险」字面；v69 快照
     P85 段载用户 2026-11-21 排片裁定、不得作先例）；
  b. ⑦ 档位 = GPT-6.2 Sol **Max**（§5.2 深度档条件「安全」+ 2026-11-21 排片
     裁定 + 2026-11-22 同轮 Max 授权，§8.4 义务照记）；
  c. 盲审必抽 1 次（另计 +1）：策略制品触所有权语义、安装来源触身份语义；
  d. 主控与 ① 统一 V4 Pro Max（2026-11-22 用户裁定，档位起跑后固定）。

一 背景与目标
- 输入 = P83 报告 §8 顺延清单逐字承接、不扩不缩：
  T2（F3 High）internal/chain/review.go + internal/console/runtime_agent.go +
    runtime_supervised.go：AgentRunner 无可持久复验的冻结策略制品，
    `grant.PolicyHash` 无独立对照物；
  T3（F6 High）internal/console/ai.go + internal/adapters/ai_probe_copilot_executor.go
    + ai_preflight.go：Windows 名称→启动映像绑定窗口、Linux 同 inode 不等于
    字节不可变；
  T8a（F8 Medium）internal/adapters/ai_probe_copilot_executor.go +
    ai_probe_copilot.go：无法证明子进程禁写未脱敏日志即拒发 BYOK。
- 目标：三项各自提供可持久复验的证明载体；三项「缺失 / 不可读 / 合形但不匹配」
  反例均拒绝；BYOK 不能证明 ⇒ 零启动。
- 协议先行（§1.4）：① 方案书必须先判定三项是否触 CONTRACTS 语义（L 档预计触）；
  触则先改 docs/CONTRACTS{,_EN}.md（及 schema / 夹具如涉）再动代码。
  守则：任何实现若需新 schema 记录、新 CLI 命令、新依赖或第三方包
  ⇒ 停下升级，不自行扩域。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro (deepseek) Max —— 修复方案书
  （契约影响判定 + 三项证明载体设计 + 反例与变异测试点 + 跨包共享契约与
  失败路径枚举，L 档必给附录 B.1 第 6/7 项）；Sol 不参与双打；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash (deepseek)；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash (deepseek) —— 每项至少一条
  反例回归测试 + 变异检验（§13.1.2）；
- ⑤ 预审员：DeepSeek V4 Pro (deepseek)，额度 1 + 1；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，行为变更切片启用，额度 ≤3；
- ⑦ 独立终审：深度档 GPT-6.2 Sol (copilot) Max，额度 1 + 3（见 〇-a/b）；
- 盲审：必抽 1 次（策略制品所有权 / 安装来源身份语义，§7.7），另计 +1、
  终审前执行；
- 批量辅助：GLM-5.3-Flash (glm) ≥2 次、严格串行；仅限四类任务
  （枚举 / 摘录核对 / 双语镜像 / MANIFEST 草案），禁择优；L 档机械工作量大、
  应尽量外包；派单写明字符上限；
- 审查输入包消毒（§7.3）；计量（R7.1）以会话调用回执为单一真值源；
- 探针：默认 0；若 T8a 验证需真实外部调用，按 §7.6 停下申请预授权；
- 风险提醒：R 达 2 ⇒ 按 §9.2 高风险标记并停下升级，不自动继续。

三 环境声明
- Ubuntu VM **开启**（10.0.0.199；go1.22.12 @ ~/go1.22；bash -lc；
  依赖走 GOPROXY=https://goproxy.cn,direct GOSUMDB=off）——T3 含 Linux
  不可写安装目录语义，必须 VM 原生验证；网络代理不开（SSH push + gh 直连
  可用）；不做 gitee。
- ⑨ 原生验证三层：Windows 本机（无 gcc 不跑 -race）+ Ubuntu VM 本地
  （Linux-tagged 用例：不可写安装目录 / 同 inode 字节不可变反例）+ CI 双平台腿
  （含 Race 步骤）；提交前 go build/vet/test ./... 全绿（§13.1.3）。

四 步骤 0 必做（§13.1）
- 锚点核对，出处 = docs/t027/CHECKLIST_v69{,_EN}.md 的 §一（7 项）：
  ① HEAD = e35e50f；② origin/main = e35e50f；③ 准则版本 = v1.71；
  ④ SELFTEST = 740/740；⑤ 三作用域 = 全 TOTAL_FAIL=0；
  ⑥ v67 行数 = 216（LF）；⑦ v68 行数 = 218（LF）；
  另核滚动窗口 v67/v68/v69 中 v69 = 209（LF）与机械读数
  （ADR 47 / OB 139 / DR 18 / 试点顺延 25）。≥2 项不符 ⇒ 以仓库实况为准，
  先报终态读数再继续；
- 台账 P85 段翻「执行中」（docs/t027/REMAINING_SLICES{,_EN}.md）；
  R/D = 0 基线；临时件落点 = tmp/P85-SECURITY-PROOF-SUBSYSTEM/
  （⑧c 前只增不删；片内 D 累计达 3 ⇒ 强制停片，§5.1 v1.71）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10）；每次 runSubagent 传 agentName + 限定名 model；
  ⑩ 派单表按 §13.2.2 逐字机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v69.md 及 _EN（§四 P85 队列项）；
2. docs/validation/security-findings-remediation.md 及 _EN（§8 顺延清单
   P84-a / P84-b / P84-c 逐字；§4 的 F3 / F6 / F8 行）；
3. docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/ 的
   终审与复审原文（F3 / F6 / F8 相关逐字）；
4. docs/SECURITY.md 及 _EN 的 :23、:60（F6 / F8 被引用句）；
   docs/CONTRACTS.md 及 _EN 的相关句（以 ① 判定书为准）；
5. docs/DELIVERY_DIRECTIVE.md 的 §1.4 / §5.2 / §5.3 / §7.5 / §7.7 /
   §8.2 / §9.2；
6. 现状文件：internal/chain/review.go、internal/console/runtime_agent.go、
   runtime_supervised.go、ai.go、internal/adapters/ai_probe_copilot_executor.go、
   ai_probe_copilot.go、ai_preflight.go。

六 修复方案要点（① 方案书须覆盖，缺一退回）
- 逐项改动点精确到 file:line；契约影响判定（触/不触 + 理由）；回归测试点
  （反例 + 变异）；每项一条可证伪闭合断言；
- T2 策略制品：持久冻结策略制品 + 与 grant.PolicyHash 的独立对照物，
  绑定 run / task / step；「缺失 / 不可读 / 合形但不匹配」三类反例均拒绝；
- T3 不可变安装来源：Windows 稳定路径 + ACL 关闭写入（名称→启动映像绑定
  窗口）；Linux 不可写安装目录 + 同 inode ≠ 字节不可变的证明载体；
- T8a BYOK 证明载体：证明子进程禁写未脱敏日志；不能证明 ⇒ 零启动拒发 BYOK；
- 变异检验：逐项撤守卫 ⇒ 对应反例红；恢复后逐字节一致。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. 协议先行：若 ① 判定触契约 ⇒ docs/CONTRACTS{,_EN}.md 同步改
   （及 schema / 夹具如涉），先契约后代码；
2. 代码：internal/chain/review.go、internal/console/{runtime_agent.go,
   runtime_supervised.go, ai.go}、internal/adapters/{ai_probe_copilot_executor.go,
   ai_probe_copilot.go, ai_preflight.go}（+测试，+新增证明载体文件）；
3. 新增报告对 docs/validation/security-proof-subsystem{,_EN}.md
   （§11.3 结构 + 三项逐条闭合状态 + artifacts 围栏 + 双语严格镜像）；
4. 台账 docs/t027/REMAINING_SLICES{,_EN}.md：P85 段状态翻「✅ 闭环」+
   完成记录（三项逐条 + R/D）；
5. docs/SECURITY{,_EN}.md:5 注记更新为「F1–F8 全闭合」（F3/F6/F8 由 P85 闭合）；
6. docs/DEV_PLAN{,_EN}.md：T027 段回写 P85 一段（2026-10-22 补录先例；
   只记结论与报告路径，不改写 T027 状态字段）。
核对项（由用户裁定）：
7. docs/validation/security-findings-remediation{,_EN}.md 的 P83 报告 §8
   是否补「P84-a/b/c 已由 P85 闭合」时点注记；
8. ADR 不新增（除非用户裁定）。
不动：
9. CHECKLIST_v69{,_EN}（v70 快照再吸收）；s1-exit / p57 门表；准则本体
   （不升版、不新增 OB/ADR 除非用户裁定）；CI；DOCUMENTATION_PLAN{,_EN}
   （已指 v69，v70 生成时再换代）。

八 边界
- 只闭合 F3 / F6 / F8 顺延部分，不碰 F1 / F4（已全闭合）、F2 / F5 / F7
  （已闭合）；不改 S1 Exit 八门判据；不改 SECURITY.md 威胁模型基线；
  不升准则版本；不做 gitee；不预设临时授权（额度用尽按 §7.8.3 停下报用户）。

九 回执要求
- 首行 P85: STARTED；四要素（裁决人=用户、日期=2026-11-22、理由=P83 顺延
  义务与闭合 F3/F6/F8、依据=本次起跑消息 + v69 快照 §四）；
- 锚点 7 项 + 模型 / 档位（主控 V4 Pro Max / ① V4 Pro Max / ⑦ Sol Max 1+3）/
  盲审必抽 / 批量辅助计划 / VM 开启逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 单次输出 ≤ 10,000 字符（§13.1.16），超限分批、禁止自重启；
- 提交信息英文；完成后停于 ⑩，等同轮授权；仅 origin。
