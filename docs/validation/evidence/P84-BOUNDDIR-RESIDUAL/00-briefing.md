# P84 起跑消息原文（§13.1.10 起跑消息即工件）

切片：P84 `BOUNDDIR-RESIDUAL`（`[SLICE]` · 实现类 M 档）
主控：DeepSeek V4 Pro（High；用户 2026-11-21 裁定切换，替代 Flash）
① 架构：DeepSeek V4 Pro Max（`deep-reasoner`）
② 实现 / ③ 测试：DeepSeek V4.1 Flash
⑤ 预审：DeepSeek V4 Pro（额度 1+1）
⑥ 独立扫描：MAI Code 2 Flash（额度 ≤3）
⑦ 终审：GPT-6.2 Sol Max（额度 1+3 + 盲审必抽 1，用户 2026-11-21 排片裁定）
批量辅助：GLM-5.3-Flash ≥ 2（严格串行）

---

## 【〇·零】记忆唤醒（开场必读；本段为起草惯例红线，未登入准则 / OB / ADR）

你是 VS Code 主控（DeepSeek V4 Pro，High 档）。本片是 P84，M 档实现类
`[SLICE]`，片名 = `BOUNDDIR-RESIDUAL`，承接 P83 顺延义务、闭合 F1 / F4。
① 阶段必须先判定 T4 / T5#4 是否触 CONTRACTS 语义：触 ⇒ 强制停下走协议
先行（先改契约再动代码）；不触 ⇒ 按方案书推进 ②–⑩。盲审必抽 1 次
（终审前执行），额度与档位不得中途升档。

**关键锚点（步骤 0 须实测核对）**：
- HEAD / origin = `71d03c3`；准则 = v1.71；SELFTEST = 740/740
- 三作用域 = `TOTAL_FAIL=0`；滚动窗口 = v65 / v66 / v67 = 220 / 216 / 216（LF）
- 当前快照 = **v67**；机械读数 ADR 47 / OB 139 / DR 15 / 试点顺延 25

---

【P84 起跑 · BOUNDDIR-RESIDUAL · 用户 2026-11-21 点名】

〇 切片标识
- 编号 P84；标识 BOUNDDIR-RESIDUAL（承接 P83 报告 §8 顺延清单 P84-d / P84-e，
  闭合 F1 / F4）；分级 [SLICE]；类型 = 实现类（平台绑定 + 事务语义）；规模 = M；
- 步骤链：0 → ① → ② → ③ 一致性 → ④ 门禁 → ⑤ 预审 → ⑥ 独立扫描 →
  ⑦ 盲审 + 终审 + 复审 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩（盲审在终审前执行）；
- 主控 = **DeepSeek V4 Pro (deepseek) High**（用户 2026-11-21 裁定切换，替代 Flash；
  编排 / 落盘 / 门禁 / 回执职责；不满足 §13.1.8 升级触发）；
- 档位限制（§8.3）：②③ = DeepSeek V4.1 Flash (deepseek) High（与主控不同模型，
  「同为 deepseek-flash 档位不可分离」条款不适用）；全片不得中途改档；
- 用户裁定（2026-11-21 排片裁定，三项）：
  a. ⑦ 额度 = **1 终审 + 3 复审**（超出 §7.5「M-L 视同 M」的 1+2 字面，
     系用户裁定、不得作为先例；报告「审查与计量」节 + 台账「额度口径披露」
     行四要素留痕）；
  b. **盲审必抽 1 次**（另计 +1、终审前执行）：T4 触身份语义（文件归属证明）、
     T5#4 触所有权语义（proof 跨任务顶替）；
  c. ⑦ 档位 = **GPT-6.2 Sol Max**（承 P79 / P83 先例，本消息发出即同轮授权，
     §8.4 义务：报告「环境与档位」表记录设定时刻 / 档位值 / 是否暂停）。

一 背景与目标
- 输入 = P83 报告 §8 顺延清单逐字承接、不扩不缩：
  T4（F1 残余）internal/evidence/bounddir_*：
    Windows 枚举后代仅比路径字符串、EqualFold 不能证明归属；
    中间 reparse 可被跟随；根以空相对路径表示而显式 "." 被拒；
  T5#4（F4 残余）internal/guard/lease.go + internal/console/runtime_supervised.go：
    批量 proof 回写只比 ResourceID、可被他任务 proof 顶替；
    代次提交与 proof 写入非同一事务、跨进程中断不可恢复。
- 目标：闭合 F1 与 F4（其余 F3 / F6 / F8 属 P85，本片不碰）。
- 协议先行（§1.4）：① 方案书必须先判定两项是否触 CONTRACTS 语义；
  触则先改 docs/CONTRACTS{,_EN}.md（及 schema / 夹具如涉）再动代码。
  守则：任何实现若需新 schema 记录、新 CLI 命令、新依赖或第三方包
  ⇒ 停下升级，不自行扩域。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro (deepseek) Max —— 修复方案书
  （契约影响判定 + T4/T5#4 修复设计 + 反例与变异测试点）；Sol 不参与双打；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash (deepseek)；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash (deepseek) —— 每项至少一条
  反例回归测试 + 变异检验（§13.1.2）；
- ⑤ 预审员：DeepSeek V4 Pro (deepseek)，额度 1 + 1；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，行为变更切片启用，额度 ≤3；
- ⑦ 独立终审：深度档 GPT-6.2 Sol (copilot) Max，额度 1 + 3（见 〇-a/c）；
- 盲审：必抽 1 次（T4 身份 / T5#4 所有权语义，§7.7），另计 +1、终审前执行；
- 批量辅助：GLM-5.3-Flash (glm) ≥2 次、严格串行；仅限四类任务
  （枚举 / 摘录核对 / 双语镜像 / MANIFEST 草案），禁择优；派单写明字符上限；
- 审查输入包消毒（§7.3）；计量（R7.1）以会话调用回执为单一真值源；
- 探针：默认 0；若 ⑨ 验证需真实调用，按 §7.6 停下申请预授权；
- 风险提醒：R 达 2 ⇒ 按 §9.2 高风险标记并停下升级，不自动继续。

三 环境声明
- Ubuntu VM **开启**（10.0.0.199；go1.22.12 @ ~/go1.22；bash -lc；
  依赖走 GOPROXY=https://goproxy.cn,direct GOSUMDB=off）；
  网络代理不开（SSH push + gh 直连可用）；不做 gitee。
- ⑨ 原生验证三层：Windows 本机（无 gcc 不跑 -race）+
  Ubuntu VM 本地（Linux-tagged 用例：bounddir unix 大小写 / reparse 反例、
  lease 事务恢复；P83 教训 = 本机不可执行项必须 VM/CI 验证）+
  CI ubuntu 腿 Race 步骤；提交前 go build/vet/test ./... 全绿（§13.1.3）。

四 步骤 0 必做（§13.1）
- 锚点核对，出处 = docs/t027/CHECKLIST_v67{,_EN}.md 的 §一（7 项）：
  ① HEAD = 71d03c3；② origin/main = 71d03c3；③ 准则版本 = v1.71；
  ④ SELFTEST = 740/740；⑤ 三作用域 = 全 TOTAL_FAIL=0；
  ⑥ v65 行数 = 220（LF）；⑦ v66 行数 = 216（LF）；
  另核滚动窗口 v65/v66/v67 中 v67 = 216（LF）与机械读数
  （ADR 47 / OB 139 / DR 15 / 试点顺延 25）。≥2 项不符 ⇒ 以仓库实况
  为准，先报终态读数再继续；
- 台账 P84 段翻「执行中」（docs/t027/REMAINING_SLICES{,_EN}.md）；
  R/D = 0 基线；临时件落点 = tmp/P84-BOUNDDIR-RESIDUAL/
  （⑧c 前只增不删；片内 D 累计达 3 ⇒ 强制停片，§5.1 v1.71）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10）；每次 runSubagent 传 agentName + 限定名 model；
  ⑩ 派单表按 §13.2.2 逐字机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v67.md 及 docs/t027/CHECKLIST_v67_EN.md（§四 P84 队列项）；
2. docs/validation/security-findings-remediation.md 及 _EN（§8 顺延清单
   P84-d / P84-e 逐字）；
3. docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/ 的
   `07-final-review.md`（T4 / T5#4 原文）与 `07-re-reviews-r1-r3.md`（r3 逐条）；
4. docs/SECURITY.md 及 _EN 的 :23（T4 被引用句）；
   docs/CONTRACTS.md 及 _EN 的 :91、:150（T5#4 被引用句）；
   （以 P83 ① 判定书 02-ruling-triage.md 的 T4 / T5#4 节为准）
5. docs/DELIVERY_DIRECTIVE.md 的 §1.4 / §5.2 / §5.3 / §7.5 / §7.7 /
   §8.2 / §9.2；
6. 现状文件：internal/evidence/bounddir.go / bounddir_unix.go /
   bounddir_windows.go、internal/guard/lease.go、internal/guard/lease_open_windows.go、
   internal/console/runtime_supervised.go（:1186 与 runtime_agent.go:2087 附近）。

六 修复方案要点（① 方案书须覆盖，缺一退回）
- 逐项改动点精确到 file:line；契约影响判定（触/不触 + 理由）；
  回归测试点（反例 + 变异）；每项一条可证伪闭合断言；
- T4：Windows 枚举与后代打开改为句柄相对（Ordinal 比较替换 EqualFold；
  中间 reparse 拒绝；根 rel 语义与显式 "." 清洗一致）；无产物不打开根；
- T5#4：代次→proof 更新做成可恢复事务；批量回写逐条校验前驱
  （claim / 代次 N-1 / PreviousHash / retained 原件），防跨任务顶替；
  同进程分支不得把任意合形 proof 当「陈旧」追平。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. 代码：internal/evidence/bounddir_*（+测试）、internal/guard/lease.go（+测试）、
   internal/console/runtime_supervised.go（+测试）；
2. 协议先行：若 ① 判定触契约 ⇒ docs/CONTRACTS{,_EN}.md 同步改
   （及 schema / 夹具如涉），先契约后代码；
3. 新增报告对 docs/validation/bounddir-residual{,_EN}.md
   （§11.3 结构 + T4/T5#4 逐条闭合状态 + artifacts 围栏 + 双语严格镜像）；
4. 台账 docs/t027/REMAINING_SLICES{,_EN}.md：P84 段状态翻「✅ 闭环」+
   完成记录（T4/T5#4 逐条闭合 + R/D）；
5. docs/SECURITY{,_EN}.md:5 注记更新为「F1 / F4 已于 <日期> 闭合（P84）」，
   F2/F5/F7 已修复句保留；
核对项（由用户裁定）：
6. docs/validation/security-findings-remediation{,_EN}.md 的 P83 报告
   §8 顺延清单是否补 P84 闭合注记；
7. ADR 不新增（除非用户裁定）。
不动：
8. CHECKLIST_v67{,_EN}（v68 快照再吸收）；s1-exit / p57 门表；准则本体
   （不升版、不新增 OB/ADR 除非用户裁定）；CI；DEV_PLAN；
   DOCUMENTATION_PLAN{,_EN}（已指 v67）。

八 边界
- 只修 T4 / T5#4，不碰 F3 / F6 / F8（P85 范围）；
- 不改 S1 Exit 八门判据；不改 SECURITY.md 威胁模型基线；不升准则版本；
  不做 gitee。

九 回执要求
- 首行 P84: STARTED；四要素（裁决人=用户、日期=2026-11-21、理由=P83 顺延
  义务与闭合 F1/F4、依据=本次起跑消息 + v67 快照 §四）；
- 锚点 7 项 + 模型 / 档位（主控 V4 Pro High）/ 额度（1+3 裁定留痕）/
  盲审必抽 / 批量辅助计划逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 单次输出 ≤ 10,000 字符（§13.1.16），超限分批、禁止自重启；
- 提交信息英文；完成后停于 ⑩，等同轮授权；仅 origin。
