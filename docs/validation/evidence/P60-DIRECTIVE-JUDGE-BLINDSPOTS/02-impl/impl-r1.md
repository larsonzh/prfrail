# ② 实现回执（prfrail-implementer · DeepSeek V4.1 Flash）—— 逐字转录

**IMPLEMENT: FINDINGS**

## ①（实现改动）修改文件清单（file:line）

**A. `tools/gates/lib/criteria/g1.js`**
- `:103-140` `parseLedger` 新增收集数组 `inblock`（`!inCompletion` 且首列去空白为 `完成记录`/`Completion record`），并跟踪 `lastHeading`（`:110`、`:121`）→ 每条记录 `{n, line, headKey, headSymbol}`（`:132-137`）；返回值新增 `inblock`（`:140`）。
- `:217-232` 新增判定：仅 `ctx.addedLineNumbers` 的变更行（`!added || !added.has(c.n)` ⇒ `e2Historical++`）；无 `headKey`/`headSymbol` ⇒ `e2Skipped++`；`headSymbol !== '✅'` ⇒ `badE2`。
- `:236-244` 输出 `out.rec('G1(e2)', …, 'in-block completion rows imply completed', …)`，detail 为规定的 `N judged, N skipped, N historical (v1.64)`。
- 现有 c/d/e/f 收集与判定**逐字未改**；`module.exports` 签名未改（`_internals` 未加新键，无测试引用）。

**B. `tools/gates/lib/criteria/g5.js`**
- `:135-147` 新增 `refs2`：`ctx.addedLines` 中匹配 `^\|\s*(?:完成记录|Completion record)\s*\|` 的行，复用 `evidencePathsFromRow` 提取路径。
- `:225-241` 新增 `G5(e2)` 判定与记录（`repoArtifactResolvable` 探测；`refs2.length` 计数文案 `(v1.64)`）。
- `G5(d)/(e)/(f)` 与 `GATES_CHANGESET_JSON` 交叉校验未改。

**C. `tools/gates/selftest.js`**
- `:12` 头部覆盖清单增 T119–T120。
- `:856` `g1Probe` 增 `e2: pick('G1(e2)')`；`:888` `g5Probe` 增 `e2: pick('G5(e2)')`。
- `:2871-2919` T119（4 断言：✅ 正例 / 🔄 反例 / added 为空不追溯 / 标题无状态符号跳过）。
- `:3176-3222` T120（3 断言：解析成功 / ghost 路径失败 / 无新增行 nothing to resolve）。

**D. `docs/DELIVERY_DIRECTIVE{,_EN}.md`（六形状实测 EQUAL）**
- `:3` 版本行 `v1.63 → v1.64`（两版同）。
- `:608` §6.3 G1 行插入 (e2) 句（EN 同位）。
- `:612` §6.3 G5 行插入 (e2)（EN 同位）。
- `:629` 脚本化状态段两处加粗短语内追加 `（含 G1(e2)，v1.64 起）` / `（含 G5(e2)，v1.64 起）`（未新增行）。
- `:111` §0.3 新增 v1.64 行（CN 权威措辞；EN 对应）。
- `:1171` §12.4 OB-138 处置列翻 `已处置`（EN 用台账既有对应词 `Resolved`，与 OB-137 一致），说明按裁定措辞替换。

**两处如实登记的实现偏差**
1. 表格行内的 `| 完成记录 | … |` 以 `\|` 转义书写（否则该行单元格数被撑断表结构）；措辞与内容未变，`G1-a` 判据即按「剔除转义 `\|`」口径计数，实测 PASS。
2. 五份文件已 `git add`（精确路径，**无 commit / push**）——否则 index 作用域的 `G5(f)` 必红；此步不在你的显式清单内，已在此声明。

## ②（selftest 断言数）selftest
`SELFTEST: PASS 710/710`（基线 703 → **710，+7**；T119 四条 + T120 三条，0 失败）。

## ③（G1/G5 自检读数）`--check=G1 --check=G5 --all --scope=index`
```
PASS G1(c)/(d)/(e)/(f)      PASS G1-a
PASS G1(e2) :: 0 judged, 0 skipped, 40 historical (v1.64)
PASS G5(e)                  PASS G5(e2) :: no added in-block completion row (nothing to resolve)
PASS G5(f) index hygiene    PASS G2 / G3(a)(b) / G4a / G4b / G5-a①② / G5-b / G5-c / G5-d
                            PASS G6-1..G6-8 / G7-a..G7-e / G8-a..G8-d / G9 / base 四门禁
FAIL G5(d) staged is a subset of declared :: stagedNotDeclared = 本片 5 个实现路径
FAIL G7-f（你已预授权忽略）
TOTAL_FAIL=2
```
`G3(a.2)` 六形状实测：两版 `lines=2195 / headings=139（行号序列全等）/ bold=5694 / pipes=3692 / blanks=525 / backticks=6140` ⇒ `SIX-SHAPE: EQUAL`；BOM+LF 已归一（`G4a` 12 文件 OK）。

**阻塞项（需片外处置，不在我的改动清单内）**：`G5(d)` 红因 `docs/validation/p60-directive-judge-blindspots{,_EN}.md` 的 `artifacts` 围栏仍是设计期 12 行，未包含本片 5 个实现路径（`docs/DELIVERY_DIRECTIVE{,_EN}.md`、`tools/gates/lib/criteria/g1.js`、`g5.js`、`tools/gates/selftest.js`）。补 5 行围栏后该条转绿；报告不在 ② 的可改清单内，故未擅改。

## ④（越界声明）未越界
- 未改 `G1(c)/(d)/(e)/(f)`、`G5(d)/(e)/(f)` 的任何收集或判定逻辑；未新增 mutation 探针；T 编号未与既有 T1–T118 冲突。
- 未碰 `g6-whitelist.txt` / `cjk-newwords.txt` / CONTRACTS / schema / fixtures / workflows / `internal/**`。
- **未 commit、未 push**（仅 `git add` 上述 5 个文件）；未派单；`tmp/` 临时脚本已删除。
