# GATES-EDGE-COVERAGE 切片验证报告（判据边界收敛：OB-90 + OB-91）

> 切片：`GATES-EDGE-COVERAGE`（v1.34）。状态：② 实现完成、门禁全绿；主提交 `3318fef` 已推送且 CI 双绿（见 §7）。

## 1. 变更概要

- **OB-90（G3(a) 配对语义）**：配对改为**任一侧出发即生成记录**——`_EN.md` 单侧变更以其 CN 基准路径为键判定（基准未动即对侧 `+0/-0` ⇒ `G3(a.1)` FAIL，不再静默）；基准 `.md` 变更仅当 `_EN` 镜像在仓库中存在时成对（无镜像的单语文档不生成记录）；配对任一侧不可读不再整体跳过，指令对形状不可比时 `G3(a.2)` fail-closed。
- **OB-91（G6-6 语言标识）**：token 规则改为**整 token 严格匹配 + 并存 fail-closed**——自每个 `lang:` 出现点起取尽**续延字符集**（字母 / 数字 / `_` / `.` / `-`，封闭列举；终结符 = 其补集），再与闭集全值比较；闭合 v1.33 的 ⑦ r3 残差（`lang:zh_cn` 截断为合法前缀、裸 `lang:` 不可见）。
- **版本包**：头部 v1.34、§0.3 新增 v1.34 日志行、附录 C 新增 C.0v 块、§12.4 处置列翻转与新增 OB-92。
- **边界**：不改 §2.2 / §2.4 / §5.2 矩阵 / §7.5 数值；不动 CONTRACTS / schema / fixtures / workflows 与 `internal/**`；**不新增 ADR**（批量辅助使用 §2.2 白名单内的 `GLM-5.3-Flash (glm)`，不触发 §1.7.4 ⑤）。

## 2. 单元实施与收口读数

| 单元 | 内容 | 收口读数 |
|---|---|---|
| 1 | OB-90 G3(a) 配对改为任一侧出发 | selftest 385/385 |
| 2 | OB-91 G6-6 整 token 严格匹配 | selftest 414/414 |
| 3 | 准则文本 v1.34 + §12.4 处置与 OB-92 | selftest 417/417 |

- 末次全量运行（**`--scope=index`，即全部已暂存变更集**；本报告对定稿快照）：`SELFTEST PASS 417/417`；`TOTAL_FAIL=0`；`G3(a)` 指令对 `+22/-6` 且六形状全 `true`；`G4a` 受检文件全部 OK；`G6-5` 与 `G6-6` 均 PASS。
- 新增自测：T74（16 条，含两支变异证明）与 T75（16 条，含一支变异证明）⇒ 自测总数 385 → 417。
- **⑨ 定稿读数**：全树编码 `ENC-TREE tracked=3614 text=3525 binary=89 violations=0`（`docs/validation/evidence/` 冻结样例 17 件按例豁免）；变更行 `§X.Y` 引用 `total=114 unresolved=0`；`go build ./...` OK、`go vet ./...` OK、`go test ./...` `ok=14 FAIL=0`。

## 3. 片内自查自纠（不计偏离）

门禁自检当场红 ⇒ 就地修 ⇒ 不计偏离（用户 2026-10-16 裁定），逐条如下：

- G4b 命中一处新汉字（原措辞的首字不在基线语料内）⇒ 改写为「取尽」，不扩白名单。
- 配平迭代：章句扩展后六形状一度 `bold` / `backticks` 不等 ⇒ 逐行定位后在同一侧补足加粗对与反引号，最终六形状全 `true` 且逐行零漂移。
- OB-90 的载荷证明首版变异对象选错（只改一行无法复现旧语义）⇒ 改为整段守卫替换，忠实回退到旧「仅 CN 侧 + 双侧同在变更集」。
- **方法教训（保留登记）**：英文侧的长句子要按**标记数**与中文侧配平（加粗对、反引号数逐行相等），落笔后必须机械复测，不得凭肉眼判断。
- **⑦ r1 输入包漏带两份申报输入**（`slice-definition.md` / `design-r1.md`）：事实 = ⑦ 在不完整输入下给出 `PASS`；根因 = 输入包生成器对缺失源**静默跳过**（非 ⑤ / ⑦ 之错）；影响 = 该 `PASS` 的可信度由代码与机械读数独立承担；处置 = 本片**不重跑、不销解**，生成器改为逐件断言存在（r2 前补齐），并登记为 §12.4 的 OB-93（`待议`）。
- **⑧c 后落痕脚本置于 `$env:TEMP`**：清理 `tmp/gec/` 后仍需落痕，故把脚本放在仓库外；事实 = `$env:TEMP\gec*.js` 与 `$env:TEMP\gec-*.txt`；清理方式 = 通配符 + `-ErrorAction SilentlyContinue`（**匹配失败时静默、无法机械证明已清**）；仓库内 `tmp/` 已清零，**仓库外残留无机械证据**；属形态偏离、不计偏离、不新增 OB。

## 4. 未决项与偏离台账

- 同型偏离 **2/3**（DIRECTIVE-OB-SWEEP 期的 A-3 / A-4），本片**无新增偏离**。
- **⑤ 预审**：r1 `PRE-REVIEW: FINDINGS`（1 Medium —— 已入册的台账状态行缺失 —— 加 3 Low）⇒ 逐条整改并**原位重跑** ⇒ r2 `PRE-REVIEW: PASS`（0 Medium+；整改含新增 3 条断言，自测 414 → 417）。
- **⑥ 独立扫描**：`INDEPENDENT SCAN: PASS`（Section B 无发现；盲输入不含 ⑤ 结论）。
- **⑦ 终审 r1**：`FINAL REVIEW: PASS`（0 发现；首轮按锚定隔离**不附 ⑤ / ⑥ 结论**，仅给变更集 / 设计 / 切片定义 / 机械读数与双态基线）。
- **§12.1 顺延登记**：v1.30 至 v1.34 共五片均在回执发出 §12.1 顺延升级项但**未在 §12.4 另立 OB**；如需正式登记由用户裁定（口径题见 §12.4 的 OB-66）。
- ⑩ 停点：提交与推送须同轮授权；本片主提交 `3318fef` 已按授权推送、CI 双绿（见 §7）。
- **⑧b 定稿**：台账补全行的结论链与收口读数已就位（章标题状态标记按惯例保持 `⬜`，完成态记在该行内）。
- 自由度清单：`docs/validation/evidence/GATES-EDGE-COVERAGE-freedom-list.md`（F-1..F-6 全部 in-bounds）。
- **tmp/ 过程工件**：⑧c 已清除 `tmp/gec/` 下的探针、审查输入包与生成脚本；`tmp/` 现仅存 `.gitkeep`（本报告为清理后定稿，未再重生成）。

## 5. 已知缺口与登记

- **OB-90 已闭环**（`已处置`）：配对语义本体改写 + 自测覆盖重写，见 §6.3 的 G3 子规则与 T74。
- **OB-91 已闭环**（`已处置`）：整 token 严格匹配 + 覆盖不变量，见 §6.3 的 G6-6 子规则与 T75。
- **新增 OB-92**（首状态词 `待议`）：ADR-024 的有效范围措辞与 §1.7.4 ④ 冲突（步骤 0 独立发现）；**本片只登记不修**，作为 `DIRECTIVE-TEMP-AUTH-LIMIT` 的输入之一。
- **本片范围内仍未闭合项**：无（OB-92 / OB-93 属登记类，非本片目标）。
- **已立为 OB-93**（首状态词 `待议`）：审查输入包的完整性无机械判据——生成器对缺失输入静默跳过（见 §3）；本片**只登记不修**，候选修法与去向（独立片或并入 `DIRECTIVE-TEMP-AUTH-LIMIT`）见 §12.4 该行。
- **候选（只在报告登记，不进 §12.4）**：**⑧c 落痕时序**——⑧c 要清 `tmp/`，但 ⑧c 自身的记录又要落盘 ⇒ 时序矛盾；建议后续治理片评估「⑧c 的清理动作须在记录落盘**之后**执行，且清理结果须可机械证明（不得静默）」，作为 `DIRECTIVE-TEMP-AUTH-LIMIT` 或独立片的输入。

## 6. 下一步建议（候选片，待用户裁定）

- **候选片 `DIRECTIVE-TEMP-AUTH-LIMIT`**：修订 §7.5 的额度上限与 §1.7.4 ⑤ 的阈值口径（ADR-025 记第 7 次同类额度例外 ⇒ 已达阈值），并把 OB-92 作为其输入之一；**独立立项**，不与判据修复混片。
- **登记方式**：以上为**建议**，本报告只留痕；**用户未明确回应则默认不登记**。

## 7. 提交与 CI（⑧b 回填，OB-72）

- **主提交**：`3318fef` — `fix: make G3(a) pairing and G6-6 language tags edge-safe (OB-90 + OB-91)`（改动文件 10 个，`+639/-50`）；已推 `origin/main`（SSH，范围 `3393043..3318fef`）；**未推 gitee**。
- **CI 观察（该提交自身的 run）**：run `39019183887`（workflow `CI`，event push，attempt 1）**首跑双腿全绿**——`Go ubuntu-latest` 2m23s、`Go windows-latest` 3m39s；**未触发首跑红处置**。
- **Gates 步硬证据（两腿实测）**：`SELFTEST: PASS 417/417`；`GATE REPORT scope=ci changed=10`；`G3(a)` 指令对六形状全 `true`；`G4b newCjk=506 unknown=[]`；`TOTAL_FAIL=0`。
- **包级证据**：两腿 `Test` 步各覆盖 14 个包（`ok` 行 14/14，`FAIL` 0）；日志中含 `FAIL` 字样的行全部来自通过断言名（如 `must FAIL`），不含失败行。
- **回填口径**：本节由 docs-only 回填写入；回填提交自身的 CI 只作观察，不再为其新增提交。

## 8. 工件

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-gates-edge-coverage.md	repo
docs/validation/directive-gates-edge-coverage_EN.md	repo
docs/validation/evidence/GATES-EDGE-COVERAGE-freedom-list.md	repo
tools/gates/lib/criteria/g3a.js	repo
tools/gates/lib/criteria/g6.js	repo
tools/gates/selftest.js	repo
```
