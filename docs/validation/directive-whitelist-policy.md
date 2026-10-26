# DIRECTIVE-WHITELIST-POLICY 切片验证报告（协作纪律 D 升格为准则正文条款）

> 切片 `DIRECTIVE-WHITELIST-POLICY`（准则治理片） | 分级 `[SLICE]` | 规模 S | 日期 2026-10-26。
> 状态：✅ **已完成**（2026-10-26；准则 v1.52；⑦ 1 终审 + 1 复审按 §7.8.3 正常收敛；三作用域 `TOTAL_FAIL=0`、`SELFTEST: PASS 697/697`；commit `e6d14a2`；CI 详见 §8）。
> [English](directive-whitelist-policy_EN.md)

## 1 变更概要（含执行自由度清单）

**本片做什么**：把协作纪律 D（白名单政策优先于授权数量）从「会话期协作纪律」升格为准则正文的正式条款，落点 = 用户裁定的候选 (b) `§1.7.2` 末尾新增段（三行，不新增小节编号）；同步 `§0.3` 增 v1.52 变更行、`§12.4` 登记 OB-123（试点继续顺延第 21 次）、附录 C 新增 C.0am、头部升 v1.52、`_EN` 严格同位；台账新增 P30 三处状态行。

**本片不做什么**：不改 `§1.7.4` / `§2.2` 载体表本体 / `§2.4` / `§5.2` 矩阵 / `§7.5` 数值 / `§6.3` 判据本体；不动 `internal/**`、CONTRACTS / schema / fixtures / workflows、`tools/gates/**`；不改 `DEV_PLAN{,_EN}`；不推 gitee。

**产物清单**：与文末 `artifacts` 围栏逐行一致（四个准则与台账文件的修改、本报告对、自由度清单）。

**执行自由度清单**（`docs/validation/evidence/DIRECTIVE-WHITELIST-POLICY-freedom-list.md`；F-1 至 F-8 为自由度、D-1 至 D-8 为用户裁定）：
- **F-1**：把用户裁定的单段文本拆成三个物理行，断点只取句号边界，不增删字符与标点。
- **F-2**：`_EN` 同位稿的术语取 `authorized quantity` / `dispatch / task-package text` / `re-review`。
- **F-3**：`§0.3` 的 v1.52 变更行仿 v1.49 至 v1.51 的四段式结构，并一并声明 OB-123 与 C.0am。
- **F-4**：OB-123 条目追加在 `§12.4` 表尾，处置列首词取 `已裁决`（闭集词）。
- **F-5**：C.0am 沿用 C.0al 的字段骨架，附注只作指针、不复述纪律 D 正文。
- **F-6**：台账 P30 的状态符号取 `🔄`（不取 `⬜`），两者均在闭集内。
- **F-7**：⑥ 的结论行按 G7-f 允许的两处口径转录，两处取值必须相等。
- **F-8**：观察项 A 落点声明的文件名取 `observation-A-landing-declaration.md`。

**一句话结论**：纪律 D 已按用户裁定的候选 (b) 落入 `§1.7.2` 正文并生效；全部门禁三作用域 `TOTAL_FAIL=0`、`SELFTEST: PASS 697/697`；⑦ 结论 `BLIND REVIEW: PASS WITH FIXES` 经 `RE-REVIEW: PASS` 复核，四条 Low 均为片外建议（判据域或逐字文本锁定），本片未改文本。

## 2 执行流水线

| 环 | 角色 | 载体 | 模型（限定名） | 状态 | 轮次 | 产物 |
|---|---|---|---|---|---|---|
| ① | 架构师 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 已执行 | 1 | 落点四候选五维对比表 + 各案拟写文本 + 拆片判定（结论行 `ARCHITECT: OPTIONS READY`，不表立场） |
| 裁定 | 用户 | — | — | 已裁定 | 1 | 选 (b)；并裁定「授权文本」= 主控派单文本 |
| ② | 实现者 | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 已执行 | 1 | 五处改动 × CN/EN（结论行 `IMPL: LANDED WITH NOTES`，两条 Note） |
| ③ | 测试工程师 | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | 已执行 | 1 | 八项一致性检查（结论行 `TEST: CONSISTENT`，偏离零、观察两条） |
| ④ | 主控 | — | — | 已执行 | 2 | 门禁三作用域 + SELFTEST（⑤ r1 整改后重跑一次） |
| ⑤ | 预审员 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 已执行 | 2（r1 `PRE-REVIEW: FINDINGS` / r2 `RE-REVIEW: PASS`） | 一条 Medium（A-1 ⑦ 额度）+ 三条 Low + 语义风险 D-1 至 D-6 |
| ⑥ | 独立扫描员 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 已执行 | 1（`INDEPENDENT SCAN: PASS`，零发现） | Section A 至 Section E |
| ⑦ | 独立终审员 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 已执行 | 2（r1 `BLIND REVIEW: PASS WITH FIXES` / r2 `RE-REVIEW: PASS`） | 四段 + 四条 Low（全部片外） |
| ⑧a | 文档工程师 | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | 本环 | 1 | 本报告 + 台账回写 |

- **返工次数**：⑤ 一次（原 Medium 级整改）；⑦ 零次（无 Medium 及以上，四条 Low 判为片外）。
- **⑦ 收口说明**：按 `§7.8.3` 的「正常收敛」在 1 终审 + 1 复审 处收口，法定上限 1+3 未用满；理由 = r1 无 Medium 及以上、r2 确认四条 Low 全属片外。

## 3 异常与兜底记录

- **主控指令错误 · G6-4 覆盖面的误述**：主控在 ① 的任务包里写「expiry 型字面量会被 G6-4 判红」；经机械核对 `tools/gates/lib/criteria/g6.js` 的 `LITERAL_PATTERNS` 只有四类（行数比、行数、文件数、聚合计数），不含日期型。① 顶住该错误指令、四稿仍回避计数与日期形态 ⇒ 对产物无影响。性质 = 主控指令错误，不得据此追究 ①。
- **主控起草错误 · 案 (c) 案名前提不成立**：主控在起跑消息写「`§2.2` 载体表注记段（与 G4b 新字纪律放在一起）」；但准则全文检索「新字」零命中，`§2.2` 无该注记段，G4b 的 CJK 新字纪律实际在 `§6.3` 的 G4 行 (b) 与「两条白名单」段。由 ① 独立发现；不影响裁定（已选 (b)）。
- **起跑消息错误 · ⑦ 额度误写 1+1**：起跑消息三处写「S 档 1+1」；`§7.5` 的分档判定规则为「同时命中多行取复审次数最多者」，本片同时命中 S 档与治理片 ⇒ 应为 1+3（先例 P23 = S 档治理片 1+3）。由 ⑤ r1 独立发现并报为 Medium，用户 2026-10-26 裁定选 A（改为 1+3）；② 照抄起跑消息不是 ② 的错。
- **P29 台账状态字段的时态**：`REMAINING_SLICES.md` 的 P29 状态字段仍含「提交与 CI 待用户授权」，而 P29 已提交并 CI 全绿。属「⑩ 停点快照口径」（同 `s1-reach-reachability.md` 头部的既有三例），不改，登记观察。
- **正例 · ① 顶住错误指令** 与 **正例 · ⑤ 正确识别预算裁定点**：两条并列记录，结论 = 独立审查层存在的理由就是主控会犯错；两例均不得反向追究审查方。
- **临时件违规 · ⑧b 在 `tmp/<切片ID>/` 内自删（非自愈档，特批继续）**：⑧b `prfrail-documenter` 在 `tmp/DIRECTIVE-WHITELIST-POLICY/` 内自建形状核验脚本、取证后自行删除、并在回执中自报。自删先于自报 ⇒ 按 `§5.1` 属**非自愈档** ⇒ 停下升级 ⇒ 用户 2026-10-26 裁定**特批继续**。**本片临时件 D 计数 = 1**；计入 `§9.3` 的「② 提前清理」**第 4 次**同型偏离。**减轻情节（如实登记、不豁免计数）**：只删自建单件、目录完好、零丢失、未删他人材料，目录集合与 ⑧b 前快照双向比对为空。**不得作先例**。
- **证据缺口（纪律 B）**：删除前的 `tmp/` 读数**不可得**（自删先于自报，违反 `§5.1` 的「自报 → 主控记录读数 → 自删」次序）。现有证据 = 「删除后 = 前快照」的集合双向比对（`Compare-Object` 两向均为空）。**未主张证据缺口已消除**。
- **主控行文缺陷 · ⑧b 任务包内的计数型形态**：主控在 ⑧b 任务包里把 `§0.3` 的 v1.52 变更行写成「版本号 + 空格 + 行」的形态（版本号尾两位数字紧邻「行」字），命中 G6-4 的行数判据而红；⑧b 改为「v1.52 变更行」，语义不变。性质 = 主控行文缺陷，门禁正确拦截，不追究 ⑧b。
- **主控格式偏差 · 成本表列数**：主控给的 ⑧b 行是 6 格，而成本表为 5 列；⑧b 按 5 列投影改写，未产生错列。性质 = 主控格式偏差，⑧b 自行适配。
- **临时件 D 计数汇总**：本片 D = 1（见本节第 1 条）；`§9.3` 的「② 提前清理」模式累计至第 4 次，用户已特批继续并要求登记证据缺口与 D 计数。
- **F-9 · 措辞补全（属 `§1.7.1` 执行自由度，不触转录契约）**：主控给的登记文本在「两向 __ 空」处使用了一个不在 HEAD 语料、也不在 48 条白名单内的汉字；执行角色改写为「两向均为空」（改 1 字、语义不变）。**性质 = 对主控登记文本的措辞补全**，属 `§1.7.1` 的「契约措辞的补全」，**不触 `§11.3` 的转录契约**（该契约约束的是审查方 ⑤⑥⑦ 的最终消息原文，非主控登记文本）。处置：**不改回、不注册白名单**（`tools/gates/**` 亦在本片边界外）。

## 4 环境与档位

- **档位**：High（`§8.2` 主控固定档）；探针额度零（未使用任何真实外部调用）。
- **工具链**：Node（`gate.js`）、Go 1.22+（`gofmt` / `build` / `vet` / `test`）；`SELFTEST: PASS 697/697`。
- **操作系统**：Windows（本地）。

## 5 审查结论摘要（⑤⑥⑦）

- **⑤ r1**：`PRE-REVIEW: FINDINGS`——一条 Medium（A-1 ⑦ 额度，已按用户裁定 A 整改）+ 三条 Low（A-2 台账时态 / A-3 标注语 / A-4 v18 来源）+ 语义风险 D-1 至 D-6 全部登记。
- **⑤ r2**：`RE-REVIEW: PASS`——A-1 关闭；整改未越界（差异恰 6 行）；未引入新问题；「不得中途升档」不成立（说明见 r2 第 3 节）。
- **⑥**：`INDEPENDENT SCAN: PASS`（零发现）。
- **⑦ r1**：`BLIND REVIEW: PASS WITH FIXES`（四条 Low，全部片外）；**⑦ r2**：`RE-REVIEW: PASS`（确认分类成立、无片内应修项、可按 `§7.8.3` 收口）。
- **⑦ 的四条 Low（片外建议，本片不修）**：

| 编号 | 级别 | 内容 | 处置 |
|---|---|---|---|
| L-a | Low | 派单文本义务无机械判据 | 记录在案（片外建议） |
| L-b | Low | `§11.3` 转录契约对 `§1.7.2` 新增段的适用范围歧义 | 记录在案（片外建议） |
| L-c | Low | 删掉正文三行无判据会红 | 记录在案（片外建议） |
| L-d | Low | 「不得作为先例」无机械交叉校验 | 记录在案（片外建议） |

- 四条均要求新增或修改 G 判据，或改动用户逐字裁定文本 ⇒ 超出本片边界。

> **防误读标注（⑦ r2 要求）**：本轮 `PASS WITH FIXES` 的建议项全部为片外建议，本片未执行对应文本整改；依据 = 切片边界排除 `§6.3` 判据本体与 `tools/gates/**`，且用户 2026-10-26 已裁定「不在本片为纪律 D 加机械判据（另立片）」。

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
（无）

## Section C
- 证据源：只根据 tmp/DIRECTIVE-WHITELIST-POLICY/impl-r2.diff 与下列正文位置并执行只读门禁，不使用任何其他审查环结论。
- 读取位置：docs/DELIVERY_DIRECTIVE.md:221-230 / docs/DELIVERY_DIRECTIVE_EN.md:221-230 / docs/DELIVERY_DIRECTIVE.md:1116-1125 / docs/DELIVERY_DIRECTIVE_EN.md:1116-1125 / docs/t027/REMAINING_SLICES.md:57-63 / docs/t027/REMAINING_SLICES_EN.md:57-63
- 命令执行：node tools/gates/gate.js --all --skip-base
- 原始读数（摘录）：
  - PASS G1-b version metadata :: docs/DELIVERY_DIRECTIVE.md header=v1.52 last=v1.52 2026-10-26/2026-10-26 rows=53 OK | docs/DELIVERY_DIRECTIVE_EN.md header=v1.52 last=v1.52 2026-10-26/2026-10-26 rows=53 OK
  - PASS G6-5 section-reference resolution :: no unresolved section reference on added ledger/changelog lines (scanned 28, refs 52) whitelist-source=worktree
  - PASS G6-6 ledger status closed set :: every added ledger row opens with a closed-set status word (scanned 2) whitelist-source=worktree
  - PASS G3(a) symmetry + mirror parity :: docs/DELIVERY_DIRECTIVE.md +19/-1 vs _EN +19/-1 OK || docs/t027/REMAINING_SLICES.md +21/-0 vs _EN +21/-0 OK
  - TOTAL_FAIL=0
- 结论：新增段位于 §1.7.2 末尾、无新增小节编号；版本头、§0.3、§12.4、附录 C 和 EN 对位均按门禁要求自洽；同一片内的 P30 状态栏和 OB-123 也与当前已知事实一致。

## Section D
- 若删掉 G1-b 守卫：G1-b version metadata 会直接红；这不是「人审」问题，而是机械版本/日期一致性检查会失效。
- 若改回旧逻辑、去掉 G6-5：G6-5 section-reference resolution 会红；这会放出未解析的节号引用。
- 若删掉 G6-6：G6-6 ledger status closed set 会红；新增 ledger 行若不以闭集状态词开头，不能通过。
- 若删掉 G3(a) 的双语张力守卫：CN/EN 变更行的对称性检查会红；本片的 EN 对位追加会被无声放过。
- 若删掉对「授权文本 = 主控派单文本」的条文本身，而只是保留「政策优先」抽象句：当前门禁不会因语义变动而红，因为没有机械判据能判定「授权文本必须是派单文本」与「不得以数量超出判越权」的语义完整性；这部分是「未覆盖的不变量」。
- 我认为未覆盖的不变量：
  1. 实际派单模板是否真的把「数量可能超出预估」作为固定措辞落地，且是否在真实执行中被遵守；门禁只能检查文档，不检查操作流程。
  2. 是否存在真实业务中「派单文本已写，但主控不遵守」的例外；门禁没有对运行时执行强制约束。
  3. P29 首例在真实工作流中的事实性留痕是否与报表一致；门禁不验证业务事实，只验证文档结构与门禁自洽。

## Section E
- 无法验证的部分：
  - 我没有读取到本切片的验证报告对或自由度清单的实际文件内容；它们不在这次 4 文件 diff 里，也不在本次只读范围内，所以我不能对「报告对/自由度清单是否与正文完全一致」作最终结论。
  - 我没有对真实外部工时日志或真实调度记录进行审计，因此对「主控在真实派单中是否写了前置说明」不做运行时确认。
- 什么会推翻我的 PASS：
  1. 若未来证据显示该段不在 §1.7.2 的末尾，而被放到其他节、或被加成新小节编号；
  2. 若 CN/EN 实际语义不等价，尤其是「授权文本」定义、政策优先级或「不作为先例」措辞发生偏离；
  3. 若 P30 状态行、OB-123 与事实状态不一致，或报告对中出现与当前门禁输出冲突的结论；
  4. 若后续门禁新增判据证明「本文真正的语义并不满足政策优先于授权数量」的工作机制要求。
- 综上，本次独立扫描依赖的证据是本次 diff、同位双语正文和项目门禁复测输出；在这些证据范围内，我认定本次扫描不构成发现，不触发阻断级风险。

INDEPENDENT SCAN: PASS
```

## 6 审查输入包摘要（含盲审隔离证明）

- **⑤**：输入 = 切片定义 + 用户落点裁定原文 + `tmp/DIRECTIVE-WHITELIST-POLICY/impl.diff`（r2 用 `impl-r2.diff`）+ 契约原句索引；未获得 ① 的对话历史与实现期推理（`§2.2` 的 R2.6 防锚定）。
- **⑥**：输入 = 切片定义摘要 + `impl-r2.diff` + 契约原句；未获得 ⑤ 的任何清单（保持独立）。
- **⑦**：首轮不附任何清单（`§4.2`）——r1 输入 = 最终 diff + 契约原句 + 按 `§1.7.3` 必须交付的自由度清单；r2 为收窄输入（只给前轮四条 Low + 机械读数 + 待裁定的分类问题）。
- **输入消毒**：所有输入均为仓库内文件的路径与切片定义，无临时件内容夹带；审查方均只读。

## 7 可证伪性（机械检查 + 变异检验）

1. 把 `§0.3` 的 v1.52 变更行改回 v1.51 ⇒ G1-b 必红（头部 v1.52 与日志末行不一致）。
2. 删除 `_EN` 侧任一新行 ⇒ G3(a.1) 对称必红。
3. 把 OB-123 处置列首词改为闭集外词 ⇒ G6-6 必红。
4. 在 C.0am 写入计数型字面量（行数、文件数、聚合计数三类形态）⇒ G6-4 必红（四类 `LITERAL_PATTERNS`）。
5. 不会红（如实登记）：删掉 `§1.7.2` 的纪律 D 三行、保留 `§0.3` / OB-123 / C.0am ⇒ G6-7 与 G7-f 均不红，仅靠 ⑤ 与 ⑦ 人工守护；这是本片新增段的结构性事实。
6. 变异检验：本片为纯文档治理片、无代码变异面 ⇒ 不适用；按 `§7.2.1` (iii) 档标注为「非等价替代」。

## 8 门禁结果

- `--scope=index`：`changed=7`；`PASS G5(f) index hygiene :: index equals worktree, no untracked path`；`TOTAL_FAIL=0`。
- `--scope=tree`：`changed=7`；`TOTAL_FAIL=0`。
- `--scope=ci`：`changed=2`；`TOTAL_FAIL=0`。
- `--selftest`：`SELFTEST: PASS 697/697`；`exit=0`。
- `G4a`（`--check=G4a --scope=tree`）：`PASS G4a BOM+LF :: 7 file(s) OK (1 excluded: evidence tree)`；`TOTAL_FAIL=0`。
- 索引与工作区一致：`git status --porcelain` 中无 `??` 与 ` M` 混合残留。

### 8.1 CI 读数（主提交 `e6d14a2`，run `40283773081`）

- 结论：`status=completed`、`conclusion=success`、`sha=e6d14a2e100da5e7c6755d9f5e2ba2b796f02908`。
- `Go windows-latest`：包级 `ok github.com/larsonzh/prfrail/internal/gates 0.468s`；`Gates` 步报 `SELFTEST: PASS 697/697` 与 `TOTAL_FAIL=0`。
- `Go ubuntu-latest`：包级 `ok github.com/larsonzh/prfrail/internal/gates 0.127s`；`Gates` 步报 `SELFTEST: PASS 697/697` 与 `TOTAL_FAIL=0`。
- 首跑双绿 ⇒ `§6.4` 的失败处置**未触发**；本片不新增提交为该 run 追加证据。

## 9 ⑨ 原生验证

- 本片为纯文档治理片、无代码改动面：`internal/**`、`tools/gates/**`、schema、fixtures、workflows 均未改动。
- `gofmt -l .`：无输出（`gofmt clean`）。
- `go build ./...`：`exit=0`。
- `go vet ./...`：`exit=0`。
- `go test ./...`：`exit=0`；全部包 `ok`（含 `internal/gates`、`internal/console`、`internal/chain` 等），无 `FAIL`。
- 如实说明：基础 Go 门禁（`gofmt` / `build` / `vet` / `test`）全部通过，属回归确认而非本片改动的直接验证。

## 10 ⑧c 收尾清洁记录

> 本节为 ⑧c 前的记录（依 `§5.1` 的时序条款：⑧c 的清理动作须在记录落盘之后执行）；实际清理动作与逐件回读证明由主控在本报告落盘后执行并登记。

- `tmp` 递归文件总数 = 69（直属目录 = `DIRECTIVE-WHITELIST-POLICY`、`gates-cache`）。
- `tmp/DIRECTIVE-WHITELIST-POLICY` 递归文件数 = 32。
- `tmp/gates-cache/g4b/` 属 v1.50 运行时缓存例外集成员，本片不清理、保持原样。
- 全仓编码与行尾逐文件核验（`git ls-files`；可读文件总计 3719）：evidence 树外违规 = 0；evidence 树内违规 = 17（全部为冻结证据包与 `*.raw.txt` 保真例外，按 `§5.1` 的豁免与 v1.40 的「B2 保真 > ⑧c 编码归一」次序处理）。
- 本片新增文件（报告对 + 自由度清单）编码核验：`.md` 恰 1 个 BOM + 纯 LF，违规 0。
- 本片临时件 D 计数 = **1**（⑧b 在 `tmp/DIRECTIVE-WHITELIST-POLICY/` 内自删一件自建脚本，属非自愈档、用户特批继续；详见 §3「异常与兜底记录」）；`tmp/DIRECTIVE-WHITELIST-POLICY/` 内其余件均在切片授权目录内、只增不删，由主控在本报告落盘后逐件精确清理并回读。
- **⑧c 清理与回读（主控执行）**：该目录内的 32 件按逐件 `Remove-Item -LiteralPath` 精确清理（未用通配符），随后移除空目录；回读 —— `Test-Path` 该目录 = `False`，`tmp` 递归件数由 69 降至 37，`tmp` 直属项 = `gates-cache` 与 `.gitkeep`，`Test-Path tmp/gates-cache/g4b` = `True`（运行时缓存例外集保持原样，其下 36 件未动），`git status --porcelain -uall` 无非预期未跟踪项。

## 11 成本与计量

| 环 | 角色 | 模型（限定名） | 次数 | 场景 |
|---|---|---|---|---|
| ① | 架构师 | `DeepSeek V4 Pro (deepseek)` | 1 | 落点四候选对比与拟写文本 |
| ② | 实现者 | `DeepSeek V4.1 Flash (deepseek)` | 1 | 五处改动同位落笔 |
| ③ | 测试工程师 | `DeepSeek V4.1 Flash (deepseek)` | 1 | 八项一致性检查 |
| ④ | 主控 | 不适用 | 0 | 主控自执行，不计模型调用 |
| ⑤ | 预审员 | `DeepSeek V4 Pro (deepseek)` | 2 | r1 预审 + 整改后复审 |
| ⑥ | 独立扫描员 | `MAI-Code-1.1-Flash (copilot)` | 1 | 盲输入独立扫描 |
| ⑦ | 独立终审员 | `GPT-5.3-Codex (copilot)` | 2 | r1 终审 + r2 复审 |
| ⑧a | 文档工程师 | `DeepSeek V4.1 Flash (deepseek)` | 1 | 本报告与台账回写 |
| ⑧b | 文档工程师 | `DeepSeek V4.1 Flash (deepseek)` | 1 | 报告定稿 + 台账翻转 |

- **批量辅助** `GLM-5.3-Flash (glm)`：零次；原因 = 本片 S 档、改动面小，枚举与比对与草稿均由主控与 ① 直接完成，未触发批量辅助场景。
- **⑦ 额度**：实际消耗 2 / 4（终审 1、复审 1；`§7.5` 治理片行上限 1 终审 + 3 复审）。
- **探针额度**：0 / 0。

## 12 明确未执行事项

- commit 与 push（须 ⑩ 停点后同轮授权）。
- 真探针（额度 0）。
- `DEV_PLAN` 回写（准则治理片不适用，`§11.2` 与 OB-12）。
- 批量辅助调用（零次，原因见 `§11`）。

## 13 下一步建议

- 建议另立片处理 ⑦ 的 L-a 与 L-c 与 L-d（为纪律 D 补机械判据）与 L-b（澄清 `§11.3` 转录契约与正文条款的关系）。
- 建议在会话 v18 快照的观察项 D/E/F 之后新增观察项 G：本片新增的纪律 D 无机械判据守护——删掉护栏句无任何判据会红，「派单文本预置说明」无判据可核；与 P29 的三条验收缺口同属「规则存在但无判据」类，下片 ① 评估是否补判据或接受人工守护为长期方案（v18 为会话快照、非仓库文件，故只写在本报告）。
- **来源注明（A-4 处置）**：纪律 D 的原始出处为会话 v18 快照（用户提供），非仓库文件；仓库内检索不到该纪律的任何历史文本。
- **标注语注明（A-3 处置）**：正文第 3 句的「政策逼出的超出授权数量」是本片标注语，不是 P29 报告的逐字引文（P29 原文为「转录契约逼出的逐字保留」，见 `docs/validation/directive-p28-followup.md`）。

**工件清单（`artifacts` 围栏；处置列取 G7-b 合法闭集词）**：

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-whitelist-policy.md	repo
docs/validation/directive-whitelist-policy_EN.md	repo
docs/validation/evidence/DIRECTIVE-WHITELIST-POLICY-freedom-list.md	repo
```
