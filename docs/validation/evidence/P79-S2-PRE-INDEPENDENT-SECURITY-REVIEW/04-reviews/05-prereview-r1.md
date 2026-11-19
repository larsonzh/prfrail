# ⑤ 预审 r1 · 逐字原文（P79）

> 归档说明（§13.1.13）：本文件为 ⑤ 预审员最终消息的逐字转录，主控未按要点重建、未改写措辞。
> 派单：`agentName = deep-reasoner`，`model = DeepSeek V4 Pro (deepseek)`；轮次 = r1。

---

`PRE-REVIEW: PASS WITH FIXES`

## Section A 结论

方案书 §1–§8 主体与起跑消息逐项一致、三面未增删、Go/文档锚点实测大多命中；但 ⑥ `sixstage` 契约（§11.3/G7-f）在设计内整体缺位且门禁现红（TOTAL_FAIL=2），报告 §2 未披露变更集门禁轨迹，另有 4 条 Low。修复均为小改，不阻塞 ⑥/⑦ 启动。

## Section B 发现清单

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| Medium | `docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/design-r1.md:73` | 「机械验收断言（收口前逐条核验）：1. 报告三面各出现独立结论档块……」——§1–§8 无 `sixstage`；准则 §11.3「**新增报告（v1.46 修订）**必须自带 `sixstage` 围栏或给出带非空理由的 ⑥ 跳过声明，两者都没有 ⇒ 判 FAIL」（⑥ 已启用，不能用跳过声明） | 设计 §4 增补「⑥ 最终消息由 ⑧a 逐字转录进 `sixstage` 围栏（结论行 `INDEPENDENT SCAN: PASS\|FINDINGS`）」，§7 断言增「G7-f 绿」 |
| Medium | `docs/validation/s2-pre-independent-security-review.md:35` | 「\| 三作用域 \| `TOTAL_FAIL=0` \| `index` / `tree` / `ci` 均 `TOTAL_FAIL=0` \| 相符 \|」——变更集门禁当前红（`gate-index-r2.txt` TOTAL_FAIL=2：G5(f) MM + G7-f×2），报告未披露 r1=6→r2=2 轨迹 | §2 表加时点注记（步骤 0 基线口径，指向 §8 门禁结果章）；⑩ 前清 G5(f)（同步暂存 `REMAINING_SLICES`）并落 `sixstage` 使 G7-f 转绿 |
| Low | `design-r1.md:19–20` | 「`SECURITY.md` §1（资产与信任边界，:7–8）、§2（写权限矩阵，:10–18）…§3 STRIDE 表（SEC-01…SEC-15，:20–37）与 §3.1（C6 预建槽，:39–45）」及「`SECURITY.md:17`（「限制必须由 runner/OS 能力落实；不能仅靠提示词」）、`:18`…」——实测 §1=:8–10、§2=:12–20、STRIDE=:26–43、§3.1=:49–53、两判据句同在 :22、表尾句=:45、§4=:55–57 | 按实测行号更新线索（引文逐字均在盘上，仅行号系统性漂移 4–18 行） |
| Low | `design-r1.md:10` | 「`PRODUCT_REQUIREMENTS.md:109` 阻断门禁句」——「独立安全评审结论仍是进入 S1 前的阻断门禁」实测在 :108（:109 为 P76 注记行） | 改引 `:108`，或写「:108 门禁句 + :109 注记」 |
| Low | `docs/validation/s2-pre-independent-security-review.md:110` | 「```artifacts …（10 行）```」——设计 §3 要求消毒结果落 `01-arch/package-manifest.md`、报告 §4:56 已列为附加输入件，该文件尚未落盘且不在围栏 | 该文件落盘时同步加围栏行并重建 `MANIFEST.md`/`SHA256SUMS.txt` |
| Low | `design-r1.md:60` | 「R7.1 计量按单一真值源逐字段登记（⑦ = 1 终审 + 3 复审，…）」——准则 §8.4（`DELIVERY_DIRECTIVE.md:885`）「每次改档必须在报告"环境与档位"表中记录：设定时刻、档位值、是否经过暂停、暂停前后门禁结果」未在设计/报告模型表复述 | 设计 §5 与报告「环境与档位」章按四字段登记 ⑦ Max 设定 |

## Section C 逐条证据与推理

- **§1 授权**：报告/设计与起跑消息逐项一致——⑦=1+3（`00-briefing.md` 〇 ⑦ 额度裁定）、⑤=1+1、⑥≤3、⑦ `GPT-6.2 Sol (copilot)` Max（§5.2:566「安全」深度判据属实）、①/批量辅助一致；「超出 §7.5『M-L 视同 M』字面标准」核实属实（`DELIVERY_DIRECTIVE.md:783`）。四要素与回执要求逐字一致。台账 P79 段 16 字段齐全、R/D=0、额度口径披露含「不得作为先例」。
- **§2 锚点**：v65 §一 8 项 + `tmp` 固定断言 = 报告 §2 表 9 行齐全；HEAD 前移口径与快照注记一致；`SELFTEST: PASS 740/740` 由 `gate-selftest.txt` 尾行证实。但变更集门禁 `gate-index-r1=6`→`gate-index-r2=2`（G5(f) MM、G7-f×2）未在任何报告文本披露 → M-2。
- **§3 三面**：与 ADR-047「未覆盖范围」①②③ 逐字一致，未增删；每面 (a)(b)(c) 齐全。函数级锚点全部实测命中：`policy.go` :30/:48/:96、`runner.go` :41/:115/:172、`readiness.go` :70–83/:113/:118、`probe_permit.go` :74–91/:126、`lease.go` :82/:104/:120、`sandbox.go` :21–34、`postflight.go` :308、`verify.go` :37、`ai_secret_windows.go` :48/:52/:77/:89/:101/:139/:150、`ai_provider.go` :47/:81/:101/:115、`ai_probe_copilot.go` :60/:134/:150、`console/ai.go` :33/:230。`ADR_REGISTER.md:58`=ADR-047 ✓、`s1-exit.md:59`=独立评审行 ✓、`DEV_PLAN.md:119` 秘密实现句 ✓。
- **§4 输入包**：三类强制消毒项（IP/用户名主机名 SSH/%TEMP% 仓库外绝对路径/凭据 token 密钥 password=）齐全；材料化要求与 §7.3 工具面标定（`DELIVERY_DIRECTIVE.md` §7.3 ③）逐字一致。`input-list.md` 40 件=10+6+8+11+5 ✓，字节合计 2704708 与我逐列求和一致 ✓，`SECURITY.md` 76 行抽验一致 ✓。
- **§5/§6**：主控只转录、首轮盲态、Section A–E 与 §7.2.2 同构、合成规则闭合并带两条限定句 ✓。
- **报告对**：CN/EN 同构——G3(a) `+123/-0 vs +123/-0 OK` 且六形状全 true、G3(b) 0/0；围栏块 10 行一致；EN face-③ 行盘上完整（截断系读取工具行为）。
- **三类反例判定**：
  - (a) 合成规则反例：**存在无机械护栏的路径**。`g7.js` `VERDICT_PREFIXES` 11 前缀不含 `FACE i`，且前缀正则 `[A-Z][A-Z _-]*` 遇数字断词 ⇒ `FACE 1: FINDINGS` 不受 G7-e 判定；`FINAL REVIEW: PASS` 在闭集内 ⇒「两面 PASS + 一面 FINDINGS 却合成 PASS」可门禁全绿。规则本身闭合，唯一拦截 = §7 断言 3（人工）。判定：接受性残差，⑨ 须逐面→合成显式推导核验。
  - (b) 独立性反例：**未发现**。报告初稿 §1–§5 通读零评审判断、零面结论、零 SEVERITY 行；方案书有「不含任何评审判断」性质声明；盲态输入构成与逐字转录契约已明写。残差为过程层（盲态构建/转录无机械判据，靠 ⑨ 断言 5 对照派单回执）。
  - (c) 机械可核验反例：**成立**。报告 §2 锚点散文（f211bca/740/740/三作用域 0/行数）无 G1–G9 校验；当前 `gate-index-r2.txt` TOTAL_FAIL=2 与 §2「均 TOTAL_FAIL=0」并存且口径未声明（后者为步骤 0 基线，changed=0 的 `gate-index/tree/ci.txt` 均 0 可佐证）。缓解 = v1.53 时点口径（须写明时点、以 ⑩/CI 为当前值）。
- **§7 可证伪锚**：7 个测试名全部 grep 命中（`policy_test.go:36`、`probe_permit_test.go:73/:147`、`readiness_test.go:89/:167`、`runner_test.go:92`、`lease_test.go:72`）。

## Section D 方案书逐条覆盖表

| 方案书条目 | 判定 | 依据 / 修复 |
|---|---|---|
| §1 授权与依据 | 符合 | 与起跑消息/台账逐项一致；§8.4 四字段未复述（L-6） |
| §2 三面映射 | 符合 | 面未增删、(a)(b)(c) 齐全、函数锚全命中；SECURITY.md 行号漂移（L-3） |
| §3 输入包与消毒 | 符合 | 三类强制项齐全、材料化与 §7.3 一致；package-manifest 落盘（L-5） |
| §4 输出契约 | 符合 | A–E 同构 §7.2.2；缺 sixstage（M-1） |
| §5 独立性/转录 | 符合 | 只转录、盲态、清单隔离均明写 |
| §6 结论档/合成 | 符合 | 二值约束、表外升级、合成闭合+两限定；面行无机械判据（反例 a） |
| §7 可证伪/验收断言 | 偏离 | 7 测试名全在盘；断言 7 条缺 G7-f（M-1）；断言 3 人工 |
| §8 风险与停点 | 符合 | 额度/格式/R≥2/表外/到期日停点齐全 |

## Section E 我无法验证的部分

1. `HEAD`/`origin/main` 实测值 `f211bca`——只读无 git 手段，无法核验；口径解释与 v65 §一注记一致。
2. 末节围栏 vs 磁盘实际变更集——`git status --porcelain` 不可得（标「无法核验」）；间接证据 = `gate-index-r2.txt` G5(d)「10 path(s) declared」、G7-b「20 artifact row(s) OK」、changed=10 与围栏 10 行一致。
3. `input-list.md` 各行数（仅抽验 `SECURITY.md`=76 与字节合计 2704708）；其余 39 件行数未逐一重数。
4. ⑦ 首轮盲态输入构成、逐字转录与 ⑤/⑥ 清单隔离的**实执行**（尚未发生，属过程承诺）。
5. EN 全文字面镜像：G3(a) 六形状全等 + 我通读结构一致，但未逐句对译全文。
6. 批量辅助 GLM ≥2 次、R7.1 计量的真值源（尚未发生）。
