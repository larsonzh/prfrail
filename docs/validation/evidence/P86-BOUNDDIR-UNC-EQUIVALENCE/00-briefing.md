# P86 起跑消息原文（§13.1.10 起跑消息即工件）

切片：P86 `BOUNDDIR-UNC-EQUIVALENCE`（`[SLICE]` · 实现类 S 档）
主控：DeepSeek V4.1 Flash（High；2026-10-31 长期裁定默认，本片不切换 Pro）
① 架构：DeepSeek V4 Pro（`deep-reasoner`；S 档裁剪，不要求第 6/7 项）
② 实现 / ③ 测试：DeepSeek V4.1 Flash
⑤ 预审：DeepSeek V4 Pro（额度 1+1）
⑥ 独立扫描：MAI Code 2 Flash（额度 ≤3）
⑦ 终审：GPT-5.3-Codex 常规档（额度 1+1；§7.5 S 档字面，无额度例外）
批量辅助：GLM-5.3-Flash ≥ 2（严格串行）

---

## 【〇·零】记忆唤醒（开场必读；本段为起草惯例红线，未登入准则 / OB / ADR）

你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片是 P86，S 档实现类
`[SLICE]`，片名 = `BOUNDDIR-UNC-EQUIVALENCE`，承接 P84 残余、闭合 F1 的
最后两笔（N13-UNC + W10 标注）。① 阶段必须先判定 N13-UNC 等值修复是否触
CONTRACTS 语义：触 ⇒ 强制停下走协议先行（先改契约再动代码）；不触 ⇒ 按
方案书推进 ②–⑩。盲审是否必抽由 ① 在启动前判定并记录理由（bounddir 属
身份/所有权语义域）。额度与档位不得中途升档。

**关键锚点（步骤 0 须实测核对）**：
- HEAD / origin = `c7ebf3d`；准则 = v1.71；SELFTEST = 740/740
- 三作用域 = `TOTAL_FAIL=0`；滚动窗口 = v66 / v67 / v68 = 216 / 216 / 218（LF）
- 当前快照 = **v68**；机械读数 ADR 47 / OB 139 / DR 17 / 试点顺延 25

---

【P86 起跑 · BOUNDDIR-UNC-EQUIVALENCE · 用户 2026-11-22 点名（先 P86 后 P85）】

〇 切片标识
- 编号 P86；标识 BOUNDDIR-UNC-EQUIVALENCE（承接 P84 报告 §8 顺延清单的
  N13-UNC 残余 + W10 证据归因；闭合 F1 最后两笔）；分级 [SLICE]；
  类型 = 实现类（Windows 目录绑定等值判定 + 文档标注）；规模 = S；
- 步骤链：0 → ① → ② → ③ 一致性 → ④ 门禁 → ⑤ 预审 → ⑥ 独立扫描 →
  ⑦ 盲审判定 + 终审 + 复审 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩；
- 主控 = **DeepSeek V4.1 Flash (deepseek) High**（2026-10-31 长期裁定默认；
  P84 的 V4 Pro 系 2026-11-21 用户裁定切换，本片不继承、不满足升级触发）；
- 档位限制（§8.3）：②③ = DeepSeek V4.1 Flash High（②③ 同为 Flash 档、
  不可分离；主控与②③ 同档属默认常态，角色分离以独立派单保证）；全片不得
  中途改档；
- 用户裁定（2026-11-22，两项）：
  a. ⑦ 额度 = **1 终审 + 1 复审**（§7.5 字面 S 档；**无额度例外**，与
     P84/P85 的用户裁定例外不同、不得套用）；
  b. **先 P86 后 P85**（两片无依赖，可并行；本片先行、P85 待本片闭环后再点名）。

一 背景与目标
- 输入 = P84 报告 §8 顺延清单逐字承接、不扩不缩：
  N13-UNC（F1 残余）internal/evidence/bounddir_windows.go：
    UNC 共享根两种规范拼写（如 `\\server\share` 与 `\\server\share\`、
    大小写变体）的等值被误拒——确定性界内误拒、方向 fail-closed、
    属可用性缺陷；盘符根等值已由 P84 修、UNC 共享根未覆盖；
  W10（证据归因 Low）internal/evidence/p84_bounddir_residual_test.go：
    变异标注「轮 11 前原形」与 r3 原文矛盾，两项须分开标注（纯文档）。
- 目标：闭合 F1 的残余——N13-UNC 等值修复（不削弱 fail-closed 方向）
  + W10 标注更正。F3 / F6 / F8（P85）与其余 F1 已闭合部分一律不碰。
- 协议先行（§1.4）：① 方案书必须先判定 N13-UNC 是否触 CONTRACTS 语义；
  触则先改 docs/CONTRACTS{,_EN}.md（及 schema / 夹具如涉）再动代码。
  守则：任何实现若需新 schema 记录、新 CLI 命令、新依赖或第三方包
  ⇒ 停下升级，不自行扩域。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro (deepseek) —— 修复方案书
  （契约影响判定 + N13-UNC 等值判定设计 + W10 标注口径 + 反例与变异测试点）；
  按 S 档裁剪（附录 B.1 不要求第 6/7 项）；GPT-6.2 Sol 不参与双打；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash (deepseek)；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash (deepseek) —— 每项至少一条
  反例回归测试 + 变异检验（§13.1.2）；
- ⑤ 预审员：DeepSeek V4 Pro (deepseek)，额度 1 + 1；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，行为变更切片启用，额度 ≤3；
- ⑦ 独立终审：常规档 GPT-5.3-Codex (copilot)，额度 1 + 1（见 〇-a）；
- 盲审：由 ① 启动前判定是否必抽（bounddir 触身份/所有权语义 ⇒ 大概率必抽，
  另计 +1、终审前执行、判定理由落盘）；
- 批量辅助：GLM-5.3-Flash (glm) ≥2 次、严格串行；仅限四类任务
  （枚举 / 摘录核对 / 双语镜像 / MANIFEST 草案），禁择优；派单写明字符上限；
- 审查输入包消毒（§7.3）；计量（R7.1）以会话调用回执为单一真值源；
- 探针：默认 0（无外部真实调用；N13-UNC 为本地目录绑定语义）；
- 风险提醒：R 达 2 ⇒ 按 §9.2 高风险标记并停下升级，不自动继续。

三 环境声明
- Ubuntu VM **不开启**（N13-UNC 为 Windows UNC 语义、落在 bounddir_windows.go；
  Linux 侧零改动，由 CI ubuntu 腿回归）；网络代理不开（SSH push + gh 直连
  可用）；不做 gitee。
- ⑨ 原生验证两层：Windows 本机（跑 N13-UNC 正反例；无 gcc 不跑 -race）
  + CI 双平台腿（windows 腿跑 UNC 用例、ubuntu 腿回归；Race 步骤照常）；
  提交前 go build/vet/test ./... 全绿（§13.1.3）。

四 步骤 0 必做（§13.1）
- 锚点核对，出处 = docs/t027/CHECKLIST_v68{,_EN}.md 的 §一（7 项）：
  ① HEAD = c7ebf3d；② origin/main = c7ebf3d；③ 准则版本 = v1.71；
  ④ SELFTEST = 740/740；⑤ 三作用域 = 全 TOTAL_FAIL=0；
  ⑥ v66 行数 = 216（LF）；⑦ v67 行数 = 216（LF）；
  另核滚动窗口 v66/v67/v68 中 v68 = 218（LF）与机械读数
  （ADR 47 / OB 139 / DR 17 / 试点顺延 25）。≥2 项不符 ⇒ 以仓库实况
  为准，先报终态读数再继续；
- 台账 P86 段翻「执行中」（docs/t027/REMAINING_SLICES{,_EN}.md）；
  R/D = 0 基线；临时件落点 = tmp/P86-BOUNDDIR-UNC-EQUIVALENCE/
  （⑧c 前只增不删；片内 D 累计达 3 ⇒ 强制停片，§5.1 v1.71）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10）；每次 runSubagent 传 agentName + 限定名 model；
  ⑩ 派单表按 §13.2.2 逐字机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v68.md 及 docs/t027/CHECKLIST_v68_EN.md（§四 P86 队列项；
   218 行，先 P86 后 P85）；
2. docs/validation/bounddir-residual.md 及 _EN（§8 顺延清单 N13-UNC / W10 逐字；
   §4 的 T4 残余行；§7 收窄复审 FINDINGS 原文）；
3. docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/ 的收窄复审原文
   （N13-UNC 残余的逐字 FINDINGS）；
4. docs/SECURITY.md 及 _EN 的 :5（T4 残余 = N13-UNC 顺延 P86 注记）；
   docs/CONTRACTS.md 及 _EN 的 bounddir 被引用句（以 ① 判定书为准）；
5. docs/DELIVERY_DIRECTIVE.md 的 §1.4 / §5.2 / §5.3 / §7.5 / §7.7 /
   §8.2 / §9.2；
6. 现状文件：internal/evidence/bounddir_windows.go（boundRelativeSuffix /
   reparseTargetRelativeSuffix / UNC 根处理）、
   internal/evidence/p84_bounddir_residual_test.go（N13 / W10 相关用例）。

六 修复方案要点（① 方案书须覆盖，缺一退回）
- 逐项改动点精确到 file:line；契约影响判定（触/不触 + 理由）；
  回归测试点（反例 + 变异）；每项一条可证伪闭合断言；
- N13-UNC：UNC 共享根等值归一化——`\\server\share` 与 `\\server\share\`
  （尾分隔符）及大小写变体在句柄相对走链下判等值、经生产入口放行；
  非等值形态（不同共享名 / 不同主机 / 中间 reparse）仍拒、fail-closed
  不削弱；盘符根等值（P84 已修）回归不红；
- W10：变异标注与 r3 原文矛盾的两项分开标注（纯文档、可字面核验）；
- 变异检验：撤掉 UNC 等值判定 ⇒ 正例变红；恢复后逐字节一致。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. 代码：internal/evidence/bounddir_windows.go（+测试）、
   internal/evidence/p84_bounddir_residual_test.go（或新增 p86 测试文件）；
2. 协议先行：若 ① 判定触契约 ⇒ docs/CONTRACTS{,_EN}.md 同步改
   （及 schema / 夹具如涉），先契约后代码；
3. 新增报告对 docs/validation/bounddir-unc-equivalence{,_EN}.md
   （§11.3 结构 + N13-UNC / W10 逐条闭合状态 + artifacts 围栏 + 双语严格镜像）；
4. 台账 docs/t027/REMAINING_SLICES{,_EN}.md：P86 段状态翻「✅ 闭环」+
   完成记录（N13-UNC / W10 逐条 + R/D）；P84 段「顺延 P86」字样撤除；
5. docs/SECURITY{,_EN}.md:5 注记更新为「F1 已于 <日期> 全闭合（T4 残余
   N13-UNC 由 P86 闭合）」，F2/F5/F7 已修复句保留；
6. docs/DEV_PLAN{,_EN}.md：T027 段回写 P86 一段（沿用 2026-10-22 补录先例；
   只记结论与报告路径，不改写 T027 状态字段）。
核对项（由用户裁定）：
7. docs/validation/bounddir-residual{,_EN}.md 的 P84 报告 §8 是否补
   「N13-UNC 已由 P86 闭合」时点注记；
8. ADR 不新增（除非用户裁定）。
不动：
9. CHECKLIST_v68{,_EN}（v69 快照再吸收）；s1-exit / p57 门表；准则本体
   （不升版、不新增 OB/ADR 除非用户裁定）；CI；DOCUMENTATION_PLAN{,_EN}
   （已指 v68，v69 生成时再换代）。

八 边界
- 只修 N13-UNC + W10 标注，不碰 F3 / F6 / F8（P85 范围）、不碰 F1 已闭合部分；
- 不改 S1 Exit 八门判据；不改 SECURITY.md 威胁模型基线；不升准则版本；
  不做 gitee；不新增依赖、不扩平台。

九 回执要求
- 首行 P86: STARTED；四要素（裁决人=用户、日期=2026-11-22、理由=P84 残余
  N13-UNC + W10 闭合 F1、依据=本次起跑消息 + v68 快照 §四）；
- 锚点 7 项 + 模型 / 档位（主控 Flash High / ⑦ Codex 常规档 1+1）/ 盲审判定
  与理由 / 批量辅助计划逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 单次输出 ≤ 10,000 字符（§13.1.16），超限分批、禁止自重启；
- 提交信息英文；完成后停于 ⑩，等同轮授权；仅 origin。
