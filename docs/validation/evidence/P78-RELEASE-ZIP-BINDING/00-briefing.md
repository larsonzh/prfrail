# P78 起跑消息原文（§13.1.10 起跑消息即工件）

> 转录纪律：本节为 2026-11-17 用户起跑消息的逐字转录，不作改写。

【SLICE】P78 RELEASE-ZIP-BINDING 起跑（M-L 实现类 · 一口气跑完 ⑩）

=====================================================================
【〇·零】记忆唤醒（开场必读；本段为起草惯例红线）
=====================================================================

你是 VS Code 主控（DeepSeek V4.1 Flash）。本片是 P78，M-L 档实现类 `[SLICE]`，片名 = `RELEASE-ZIP-BINDING`。目标 = v0.1.0 正式 Release 创建 + 最终 ZIP 候选工件构建 + 发行记录绑定（版本/commit/run/文件名，三处落点）+ 安装生命周期复验，S-G3 从「未通过」转为「通过」（P76 Q5 四步的第②③④步；第①步 Release 批准 = 本起跑消息同轮授权，见【九】）。

**关键锚点（步骤 0 须实测核对）**：
- HEAD / origin = `c935713`（v61 快照提交）；前锚点 CI run = `44004444624`（双腿 success）
- 准则 = v1.69；SELFTEST = 737/737；三作用域 = `TOTAL_FAIL=0`
- 滚动窗口 = v59/v60/v61 = 180/242/214（CN/EN 对称，按 LF 计数；**v58 已挤出**）
- 当前快照 = **v61**（v60 冻结，被 v61 替代）
- v61 §四 已把 P78 列为「`RELEASE-ZIP-BINDING`（S-G3 收口 · `下一项` · P78）」，依赖已满足（S-G4 已通过；版本号 v0.1.0 已由 ADR-010 冻结 ⇒ P76 Q5「绑定与复验必须等版本号冻结」前置已满足）

**P78 依赖与版本号来源**：ADR-010 已批准冻结「首版候选 v0.1.0」；ADR-009 已冻结安装模型（手工解压 + 独立目录 + 显式路径运行 + 不改 PATH/注册表/系统目录）；ADR-006 已冻结发布信任（固定 commit SHA + GitHub CI + SHA256SUMS + SBOM + 许可证清单为必需证据；**不做代码签名**）。

**v61 快照误差登记（本片随台账如实登记，不追改 v61）**：v61 §四 P78/P79 队列项写「⑦ 额度 1 终审 + 3 复审」系引用误差；准则 §7.5 权威口径 = **M-L 视同 M 档 ⇒ 1 终审 + 2 复审**。本片按 §7.5 执行；修正随 v62 或后续 [ROUTINE] 补正。

**⑦ 选档判定登记（用户 2026-11-17 同轮裁定）**：本片位于发布前终审位置，但审查对象为文档与机械证据链、产品代码零改动 ⇒ **采用常规档 `GPT-5.3-Codex`**（不选深度档 `GPT-6.1 Sol`）；若片内出现任何产品代码/安全面变更 ⇒ 停下升级重议档位。

**临时件纪律（D=2 防再犯句）**：「执行角色（`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`）**与主控**在 `tmp/` 下操作时，一律使用 `tmp/<切片ID>/`；**CI 观察捕获件同理（不得落 `tmp/` 根）**；片内 D 累计达 3 次 ⇒ 强制停片」。本片临时件一律落 `tmp/P78-RELEASE-ZIP-BINDING/`。

**OB-139 观察中**（第 5 次已发生；P77 主片 D = 0 未触发第 6 次）：本片为 [SLICE] 实现类，若出现同型违规 ⇒ **OB-139 第 6 次 ⇒ 转待议 + 立项**。

**动态额度预置（`R-ACCUMULATION-MECHANISM`，v61 §五 切片级建议 1）**：R 累积到 ≥2 ⇒ 按 §9.2 标记高风险；**本片 ⑦ 额度已在启动前按 M 档定死（1+2），不得中途升档**；触发高风险 ⇒ 停下升级用户裁决。

=====================================================================
【〇】级别判定与环境
=====================================================================

`[SLICE]`（M-L 档实现类）。主控 = `DeepSeek V4.1 Flash (deepseek)`。① 架构 = `DeepSeek V4 Pro (deepseek)`（`deep-reasoner`，只读）；⑤ 预审 = `DeepSeek V4 Pro (deepseek)`（首选，§2.2）；⑥ 独立扫描 = **必做**（§5.2 矩阵中等行：P78 = 行为变更、多文件、有测试点，不属「简单」——「简单」定义同 §7.2 跳过三条件）——⑥ 模型 = `MAI-Code-1.1-Flash (copilot)`（`independent-reviewer` 载体，只读；成本上限 ≤ 3 次，§7.2）；⑦ 终审 = **`GPT-5.3-Codex (copilot)`（常规档）**——选档判定见【〇·零】登记句；**盲审不抽**（不触必抽三语义；按比例抽查自 P76 窗口起不触发，启动前判定记录理由）。**不派 Sol 双打**。探针额度 = 0（无付费模型外部调用）。

**批量辅助**：`GLM-5.3-Flash (glm)` × 2 次（`deep-reasoner` 载体，只读；详见【六】）。

=====================================================================
【一】目标与背景
=====================================================================

**P78 目标**（按 P76 Q5 四步的第②③④步执行）：
1. 步二：构建**最终 ZIP 候选工件**——基线 = P76 Q5 步二「CI candidate-build 可用」：`ci.yml` 的 `Candidate build` job 为 `workflow_dispatch` 触发（v61 CI 显示 skipped 即常态）；路径由 ① 判定（走 CI 触发构建 vs 本机同参数构建，须论证）；产出 = artifact + `SHA256SUMS` + SBOM + 许可证清单；
2. 步三：绑定发行记录（**版本 / commit / run / 文件名**），落点 = `INSTALLATION.md:206` 明定的**三处**（S1 exit 报告 / 支持矩阵 / 发行说明）；
3. 步四：用最终 ZIP 重跑安装生命周期（安装 / 升级 / 回滚 / 卸载，Windows 本机）；
4. GitHub Release `v0.1.0` 创建（Release notes + ZIP/SHA256SUMS 附件上传）；
5. 文档同步（见【八】三层清单）——S-G3 判据达成。

**现状实况（主控须 ① 前实读复核）**：`ci.yml` 的 `Candidate build` 与 `Bootstrap and release evidence` job 均为 `workflow_dispatch` 触发（带 candidate/verifier/bootstrap 等 SHA inputs）；`cmd/prfrail-release` = 发行工具入口；v60 快照记载候选 artifact 已做过首次运行/升级/回滚/保留数据卸载。**本片新增的本地辅助脚本一律落 `tmp/P78-RELEASE-ZIP-BINDING/` 一次性使用，不入库**（边界见【七】）；CI 既有 candidate-build 流程文件零改动；构建命令与参数须在报告与 Release notes 中如实登记以保可复现。

**P78 边界（本片不改）**：
- **不改产品代码**（`.go` / `go.mod` 零改动；构建失败 ⇒ 按 §1.7.2 停下报实况，不得自行修代码）
- **不改 CI 工作流**（`ci.yml` 零改动；workflow_dispatch 触发是外部动作、不改文件）
- **不改 ADR-009 ZIP 结构**；**不做代码签名**（ADR-006）
- **不改** S1_SUPPORT_MATRIX 平台清单与支持层级（Windows 仍「正式支持候选」、Linux 仍「核心 CI」；**不做 Linux 正式发布**；绑定落点是加法式、不动清单与层级）
- **不改** S1 Exit 八门判据（p57 门表 / s1-exit 门表零改动——P76 已明令；S-G3 翻「通过」的事实由本片报告承载）
- **不改** ADR-010 EOL 通知方式与 ADR-009/ADR-006 任何已批准要点
- **不做** gitee 任何动作
- **不改** 准则 / CONTRACTS / schema / fixtures / 其它 ADR

=====================================================================
【二】步骤 0（锚点核对 + 建包 + 台账登记 + 权限预检）
=====================================================================

**步骤 0a：锚点核对（8 项）**
1. `git rev-parse --short HEAD` 与 `origin/main` —— 期望均 = `c935713`
2. 准则 `docs/DELIVERY_DIRECTIVE.md` 头部 = `版本：v1.69`
3. `node tools/gates/gate.js --selftest` = `SELFTEST: PASS 737/737`
4. 三作用域 `--scope=index`、`--scope=tree`、`--scope=ci` = 均 `TOTAL_FAIL=0`
5. `tmp` 顶层 = `gates-cache` + `.gitkeep`；排除 `gates-cache` 后递归文件数 = 1
6. `git status --porcelain -uall` 行数 = 0
7. 滚动窗口实测（按 LF 计数，CN / EN 对称）：**v59 = 180、v60 = 242、v61 = 214**（v58 已挤出）
8. 若 ≥2 项不符 ⇒ **停下报实况**，不得静默修正

**步骤 0d：Release 权限预检（新增，本片特有）**
- `gh auth status` 确认登录态与 `repo` scope；`gh repo view larsonzh/prfrail --json name` 确认可读；确认 `gh workflow run` / `gh release create` 所需 write 权限
- 若 `gh` 无 write 权限 ⇒ **停下报实况**（Release 创建/CI 触发需你处置，主控不得自裁）

**步骤 0b：证据包建包**（`docs/validation/evidence/P78-RELEASE-ZIP-BINDING/`）：`00-briefing.md` / `input-list.md` / `freedom-list.md` 三件初始载荷。
**步骤 0c：台账登记 P78 段**（`docs/t027/REMAINING_SLICES{,_EN}.md`，格式参照 P77 段；十二字段按 §1.2；登记时附带「v61 队列项 ⑦ 额度引用误差」披露句）。步骤 0 结束条件：8 锚点 + 权限预检 + 证据包三件 + P78 台账段 + `tmp` 零残留。

=====================================================================
【三】① 阶段（架构方案书，核心产出）
=====================================================================

**① 架构师**（`deep-reasoner` + `DeepSeek V4 Pro`）产出**方案书**（落 `tmp/P78-RELEASE-ZIP-BINDING/design-r1.md`），必须回答：

**Q1：构建链路**——ZIP 内容清单（prfrail.exe / LICENSE / README 最小集）；**路径判定**：CI candidate-build 触发（`gh workflow run ci.yml`，须实读 `ci.yml` 全部 `workflow_dispatch` inputs 并判定可满足性）vs 本机同参数构建（`go build -trimpath` 等，须论证等价），选定路径并给出理由；版本注入（`main.version=0.1.0`）；`cmd/prfrail-release` 在本片的角色（用或不用，理由）；构建产物下载/落盘方式（全落 `tmp/P78-RELEASE-ZIP-BINDING/`）。

**Q2：发行记录绑定**——绑定四元组（**版本 / commit / run / 文件名**）的取值与 `INSTALLATION.md:206` **三处落点**（S1 exit 报告 / 支持矩阵 / 发行说明）各自的加法式目标文本；发布证据（SHA256SUMS / SPDX SBOM / 许可证清单，ADR-006）的生产方式、SBOM 生成工具选型（仓库既有工具优先）、绑定落点（Release 页面 + 报告 + Release notes）。

**Q3：Release 创建**——`gh release create v0.1.0` 的精确参数（tag、notes 来源、附件清单）；Release notes 内容骨架（含 ADR-010 冻结的 EOL 公告义务判定：v0.1.0 是首版 minor、无前一 minor ⇒ 无 EOL 公告义务，如实记录）；创建与推送/CI 构建的顺序（须论证：CI candidate-build run 号须在绑定与 Release 创建前拿到 ⇒ 顺序由 ① 定案）。

**Q4：复验矩阵**——四项（安装 / 升级 / 回滚 / 卸载保留数据）各在什么状态路径下怎么验；升级复验的路径设计（v0.1.0 是首版 ⇒ ① 据 `INSTALLATION{,_EN}` / P76 报告 §7 实读给出可行口径，如「同版本重装覆盖升级」或「候选→正式」）；每项的成功断言与失败处置（失败 ⇒ 停下升级，不得自裁）。

**Q5：文档同步判定（三层归属）**——逐份实读【八】候选文档，给出第一层 / 第二层 / 不触归属 + 精确 `file:line` + 精确目标文本；**特别判定两项敏感点**：① `ADR_REGISTER:62`「当前仍不是已发布 SLA」——**默认保留不翻**（v0.1.0 是候选版，候选发布 ≠ 正式 SLA 生效；该句待 S1 Exit 通过后或用户裁定再翻；若 ① 判应翻 ⇒ 停下升级用户裁决）；② `ADR_REGISTER:20` ADR-009 状态列「最终 ZIP 仍须绑定发行记录并复验」闭括号（倾向翻，① 给出目标文本）。

**Q6：组合 grep 排查**——发布相关波及面：`RELEASE-ZIP` / `GitHub Release` / `v0.1.0`（独立 grep 分别跑，禁止单关键词联查噪声面）；命中除 Q5 候选文档外的文档 ⇒ 逐条记录交 Q5 升降级（升降规则见【八】末尾）。

**① 阶段结束条件**：方案书成文（`DESIGN: DONE`）后直接 ②；出现「扩范围 vs 自由度」歧义 ⇒ §1.7.2 停下升级。

=====================================================================
【四】② 阶段（构建 → 绑定 → Release → 复验 → 文档同步）
=====================================================================

**②a 构建**：按 ① Q1 选定路径构建最终 ZIP；构建日志、artifact 与 SHA256SUMS 落证据包（`tmp/P78-RELEASE-ZIP-BINDING/`，不入库）。
**②b 发行记录绑定 + Release 创建**：按 ① Q2/Q3（顺序须确保 CI run 号先拿到；Release 创建在绑定四元组齐备后执行）。
**②c 复验**：按 ① Q4 矩阵执行（对**最终 ZIP**，非候选），逐项记录成功断言证据。
**②d 文档同步**：按 ① Q5 目标文本逐处断言式替换（第一层必做；第二层视判定）。
**②e 报告初稿 + 证据包追加**：`docs/validation/release-zip-binding{,_EN}.md`（新建报告对；含 Q1–Q6 + 构建/绑定/Release/复验证据 + S-G3 收口章 + 门禁读数 + `sixstage` 围栏——本片 ⑥ 必做，报告须含 `sixstage` 围栏块）；证据包追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` / `MANIFEST.md` / `SHA256SUMS.txt`。

**不改清单**（② 阶段）：产品代码 / `ci.yml` / p57 门表 / s1-exit 门表 / 准则 / CONTRACTS / schema / fixtures。

=====================================================================
【五】③–⑦ 审查链
=====================================================================

**③ 一致性**：S-G3 收口结论与 p57 §4.2 / s1-exit / v61 §四 / ADR-046 不矛盾（**p57 / s1-exit 门表 / v61 均不改，仅本片报告承载新事实**）；绑定四元组与三处落点互相自洽；发布证据与绑定记录互相自洽；CN/EN 严句对齐；改动行数 CN/EN 对称。
**④ 门禁**：三作用域 + SELFTEST 全绿；`go build ./...`、`go vet ./...`、`go test ./...` 全绿（构建行为不改代码，但须跑基础门禁作证）。
**⑤ 预审**（`DeepSeek V4 Pro`，deep-reasoner，只读；额度 ≤ 1 + 整改复审 1）：r1 报 Medium+ ⇒ 整改后 r2；r2 仍 Medium+ ⇒ 停下升级。
**⑥ 独立扫描**（**必做**；`MAI-Code-1.1-Flash (copilot)`，`independent-reviewer` 只读载体；§7.2 上限 ≤ 3 次；输出 Section A–E 契约）：输入包 = 方案书 + 落库 diff + 门禁原始输出 + 复验证据路径（不含 ⑤ 清单，防锚定）；报 Medium+ ⇒ 修复后重跑。
**⑦ 终审**（`GPT-5.3-Codex (copilot)` 常规档，`independent-reviewer` 只读载体）：
- **额度 = 1 终审 + 2 复审**（§7.5 M-L 视同 M 档；**v61 队列项「1+3」系引用误差，本片按 §7.5 执行并在台账登记**）；第 2 次复审须收窄输入（只给修复 diff + 前轮发现 + 机械读数）
- **盲审不抽**（不触必抽三语义；按比例抽查 P76 窗口不触发；启动前判定并登记理由）
- 首轮不附任何 ⑤/⑥ 清单；r1 报 Medium+ ⇒ r2；r3（第 2 次复审）仍 Medium+ ⇒ **停下升级用户裁决**（额度用尽，无 r4）
- ⑦ 结论须含「⑥ 必做已执行且 Section A–E 契约合规」的确认

=====================================================================
【六】批量辅助配置（GLM-5.3-Flash × 2）
=====================================================================

| # | 任务 | 产出 |
|---|---|---|
| **batch-01** | 发布链路现状枚举：`ci.yml` 的 `Candidate build` / `Bootstrap and release evidence` job 与 `workflow_dispatch` inputs 逐行摘录 + `cmd/prfrail-release` 功能摘录 + `testdata/selfhost/` 与 inputs 的对应关系 + `tools/` 下可复用工具清单——供 ① Q1 参考 | `02-batch/batch-01-release-pipeline.md` |
| **batch-02** | 六份文档发布相关表述逐条枚举（含 `file:line` + 原文摘录）：`S1_RELEASE_NOTES{,_EN}` / `INSTALLATION{,_EN}`（含 `:206` 三处绑定句）/ `INSTALLATION_PLAN{,_EN}` / `ADR_REGISTER{,_EN}`（ADR-009/010 行）/ `DEV_PLAN{,_EN}`（T017/T018）/ `OPERATIONS{,_EN}`——供 ① Q5 判定 | `02-batch/batch-02-release-docs.md` |

纪律：只做事实枚举，不做判定；产出交 ①；不计入 ⑦ 额度。

=====================================================================
【七】门禁与边界
=====================================================================

**门禁**：步骤 0 前后、② 后、commit 前各跑三作用域 + SELFTEST；Release 创建前跑一遍基础门禁（build/vet/test）。

**边界（不做清单）**：不改产品代码（`.go` / `go.mod`）；不改 `ci.yml`；不改 p57 门表 / `s1-exit` 门表；不改 `docs/DELIVERY_DIRECTIVE{,_EN}`（准则零改动）；不改 `CONTRACTS` / `SECURITY` / `ARCHITECTURE` / `TEST_STRATEGY` / `BUSINESS_WORKFLOWS` / `CODING_CONVENTIONS` / `DOCUMENTATION_PLAN` / `RFC-*.md` / `CHECKLIST_v61{,_EN}`（冻结）/ v50–v61 任何快照 / schema / fixtures；不改 ADR-009 ZIP 结构与 ADR-006 发布信任；不做代码签名；不做 Linux 正式发布；不改 S1_SUPPORT_MATRIX 平台清单与支持层级（绑定为加法式）；不新增 ADR / OB（除非用户裁定）；不做 gitee 任何动作。

=====================================================================
【八】文档同步更新清单（三层）
=====================================================================

**第一层：确定会触**：
1. `docs/t027/REMAINING_SLICES{,_EN}.md` — P78 台账段登记（步骤 0c）+ 完成记录（⑧b；只写报告路径 + R/D）
2. `docs/validation/release-zip-binding{,_EN}.md` — 本片报告对（新建；含 S-G3 收口章 + sixstage 围栏）
3. `docs/validation/evidence/P78-RELEASE-ZIP-BINDING/` — 证据包
4. `docs/S1_RELEASE_NOTES{,_EN}.md` — v0.1.0 正式 Release 事实 + **绑定落点之一**（版本/commit/run/文件名）；「发布前剩余门禁」段更新（① 给精确目标文本）
5. `docs/INSTALLATION{,_EN}.md` — 正式 ZIP 文件名 + 下载入口 + 校验命令；「当前尚无正式 GitHub Release」类句翻转（① 给精确目标文本）
6. `docs/S1_SUPPORT_MATRIX{,_EN}.md` — **绑定落点之一**（加法式绑定版本/commit/run/文件名；**不改平台清单与支持层级**，① 给精确目标文本）

**第二层：视 ① Q5 实读触**：
7. `docs/ADR_REGISTER{,_EN}.md` — ADR-009 状态列闭括号（「最终 ZIP 仍须绑定发行记录并复验」→ 已落实句）；ADR-010 行「当前仍不是已发布 SLA」**默认不翻**（翻须停下升级）
8. `docs/INSTALLATION_PLAN{,_EN}.md` — ZIP 文件名 / 下载入口冻结状态更新
9. `docs/PRODUCT_REQUIREMENTS{,_EN}.md` — §7 发布相关待落实句（若涉）
10. `docs/DEV_PLAN{,_EN}.md` — T017/T018 发布里程碑状态（仅更新、不改历史记录语义）
11. `docs/OPERATIONS{,_EN}.md` — 下载 / 校验 / 升级 / 回滚运维说明
12. `docs/validation/s1-exit{,_EN}.md` — **绑定落点之一**（发行记录绑定加法式，落非门表区；① 给出精确位置）+ 剩余风险第 1 条「正式 ZIP 文件名和下载入口仍待发行绑定」→ P78 完成句（**非门表行**，同 P77 的 s1-exit:66 先例：P76「批准后改动点」清单精神；**门表仍零改动**）

**第三层：确定不触**：`DELIVERY_DIRECTIVE` / `CONTRACTS` / `SECURITY` / `ARCHITECTURE` / `TEST_STRATEGY` / `BUSINESS_WORKFLOWS` / `CODING_CONVENTIONS` / `DOCUMENTATION_PLAN` / `RFC-*` / `CHECKLIST_v61{,_EN}` / v50–v61 快照 / `.go` / `.json` / `.yml` / schema / fixtures / `go.mod` / 其它 ADR。

**① Q5 判定优先于本清单初始分类**；升降须在方案书如实披露；第三层 → 第二层升级若与【七】边界冲突 ⇒ §1.7.2 停下升级。

=====================================================================
【九】授权（同轮授权）
=====================================================================

- **产品所有者批准（Q10=(C) 第一步的 Release 部分）**：本起跑消息即用户对 **v0.1.0 Release 创建**的同轮批准；S-G5 记「发行部分已批准」；**硬边界：本批准不得被解读为 S1 Exit 通过**（S1 Exit 整体批准待 S2 后）
- **Release 创建外部动作授权**：`gh release create v0.1.0` + 附件上传（ZIP / SHA256SUMS / SBOM / 许可证清单），仅在步骤 0d 权限预检通过后执行；发布对象 = `github.com/larsonzh/prfrail`（origin），**不做任何 gitee 动作**
- **CI 触发授权（若 ① 判定走 CI 构建路径）**：`gh workflow run ci.yml`（workflow_dispatch，带 inputs）属外部动作，同一授权面内；触发前须把 inputs 取值与依据登记进报告
- **commit / push 授权**（提交信息英文；建议 `docs: land the v0.1.0 release zip binding (P78)`）；push 仅 origin
- **回填提交授权**（§11.2）：主提交推送 CI 双绿后，回填 commit hash + run 号进报告并再推送（回填 run 只作观察记录）
- **CI 观察**：每次 push / 触发后观察双腿；首跑红 ⇒ 按 §6.4 取证、复跑、登记
- 探针额度 = 0

=====================================================================
【十】⑩ 回执要求（按附录 B.7.3 停点回执模板）
=====================================================================

1. **`SLICE: READY_FOR_REVIEW` / `SLICE: BLOCKED`**（首行）
2. **步骤 0 的 8 项锚点实测**（含滚动窗口 180 / 242 / 214）+ **步骤 0d 权限预检结果**
3. **证据包三件落盘确认** + P78 台账段登记确认（含 v61 ⑦ 额度误差披露句）
4. **批量辅助两次产出摘要**（字符数 + 覆盖范围 + 落盘路径）
5. **① 方案书成文**（`DESIGN: DONE`）：Q1–Q6 逐项产出 + Q5 三层归属（含升降级披露）+ Q6 组合 grep 命中清单 + 每处改动精确目标文本
6. **② 逐项结果**：构建路径与 ZIP 内容清单 + SHA256SUMS；绑定四元组与三处落点；Release tag/URL 与附件清单；复验矩阵逐项成功断言；`git diff --numstat` 逐文件 + CN/EN 对称
7. **③ 一致性检查结果**（S-G3 与 p57 / s1-exit / v61 / ADR-046 不冲突）
8. **④ 门禁原始行**（index / tree / ci / selftest / build / vet / test）
9. **⑤ 预审**（r1 / r2 结论 + R 计数）
10. **⑥ 独立扫描**（必做；结论 + Section A–E 围栏转录 + 调用次数）
11. **⑦ 终审**（Codex 常规档，选档判定见【〇·零】；r1 / 复审结论 + R 计数；额度 1+2；盲审不抽理由；若复审仍 Medium+ ⇒ 停下升级实况）
12. **⑧a / ⑨ / ⑧b / ⑧c**（报告定稿 / 门禁复跑 + SHA256SUMS 自校验 / 台账完成记录 / tmp 清理 + 五条读数）
13. **commit hash + push 后 CI run 号 + 双腿时长与 success/failure**（含 Release 创建时点、CI 构建 run 号（若走 CI 路径）与绑定 run 号）
14. **不做清单逐项确认**（产品代码 / ci.yml / p57 / s1-exit 门表 / 准则零改动；无代码签名；无 gitee；S1_SUPPORT_MATRIX 平台清单与支持层级未动）
15. **R / D 计量**
16. **终态**：HEAD / origin/main / porcelain / tmp 顶层 / 排除缓存文件数
17. **给下一片的锚点**：P78 闭环后的 HEAD + CI run 号 + Release tag
18. **S-G3 收口事实**：本片报告明确 S-G3 从「未通过」转「通过」；p57 / s1-exit 门表零改动；S-G5 记「发行部分已批准」；下游 = S-G6 门级侧并陈
19. **G6-4 自查**：报告中无 `\d{2,}\s*行` 类字面
20. **§13.3.15「给下版清单的事实包」**（10 字段）

=====================================================================
【十一】执行节奏与红线
=====================================================================

1. **全程无中间停点**：0 → ① → ② → ③ → ④ → ⑤ → ⑥ → ⑦ → ⑧a → ⑨ → ⑧b → ⑧c → ⑩；触发停升级边界（§1.7.2）⇒ 停下报实况。
2. **构建/复验/Release 任一失败 ⇒ 停下升级**（不修代码、不自裁换方案）。
3. **临时件**一律 `tmp/P78-RELEASE-ZIP-BINDING/`；`tmp/` 根不得落任何文件；D 达 3 次强制停片；**构建产物/ZIP/日志全部落该目录**，⑧c 清理。
4. **不改准则版本**（维持 v1.69）；**不改 v61 快照**（其 ⑦ 额度误差由本片台账登记披露）。
5. **CN / EN 严格对称**：报告对六形状全等；文档改动逐字同位。
6. **p57 / s1-exit 门表零改动**（S-G3 收口事实由本片报告承载）。
7. **⑦ 额度**：1 终审 + 2 复审；R ≥ 2 ⇒ 高风险标记 + 停下升级，**不自行升 1+3**。
8. **⑥ 必做**：上限 ≤ 3 次；格式偏离允许重发 1 次；达上限 ⇒ 按 §7.2 失败处置如实留痕，由 ⑦ 闭环。
9. **push 仅 origin**；**不做任何 gitee 动作**；Release 发布对象仅 GitHub origin。
10. **G6-4 自查**；**E-1 流程纪律**（⑩ 前二次回读 tmp 五条读数）；**D=2 防再犯句**适用；**OB-139 观察中**（本片 [SLICE]，再犯 ⇒ 第 6 次 ⇒ 转待议 + 立项）；**门禁坑 26 / 27 / 28 / 30** 特别注意。

现在开始执行 P78 步骤 0。
