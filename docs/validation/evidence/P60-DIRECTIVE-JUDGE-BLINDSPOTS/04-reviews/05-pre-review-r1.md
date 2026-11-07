PRE-REVIEW: PASS WITH FIXES

# P60 `DIRECTIVE-JUDGE-BLINDSPOTS` ⑤ 预审回执 r1（deep-reasoner · DeepSeek V4 Pro）

审查对象：阶段 B 的 ②（J2 判据 + J3 处置列，r1+r2 落盘态）与 ③（r1+r2 夹具与测试）。结论：J2/J3 实现与测试整体正确；逐项复核 8 项中 6 项 PASS、2 项 Medium（见 §4）。Medium 均非判据逻辑缺陷、可机械修复，不堵 ②③ 成果。

## 1. 审查范围与输入（file:line）

| # | 输入 | 位置 |
|---|---|---|
| 1 | 裁定原文（J2+J3 组合、J1 不选、升版授权、⑦ 重点、⑧ 装配要求） | `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/ruling-2026-11-06-p60.md`（全文） |
| 2 | ② 回执 r1 | `02-impl/impl-r1.md`（全文）；r2 落盘态经磁盘 `git diff --cached` 核对 |
| 3 | ③ 回执 r1/r2（F-1/F-2、O-1/O-2/O-3） | `03-test/test-r1.md`、`03-test/test-r2.md`（全文） |
| 4 | G1(e2) 实现 | `tools/gates/lib/criteria/g1.js:103-140`（inblock 收集）、`:214-232`（判定）、`:234-244`（记录） |
| 5 | G5(e2) 实现 | `tools/gates/lib/criteria/g5.js:98-147`（`evidencePathsFromInBlockRow`）、`:184-200`（rows2/refs2 收集）、`:277-290`（判定与记录） |
| 6 | 夹具 T119–T124 | `tools/gates/selftest.js:2899-3022`（T119/T121）、`:3271-3348`（T120/T122）、`:3354-3433`（T123/T124） |
| 7 | 准则 CN/EN | `docs/DELIVERY_DIRECTIVE{,_EN}.md:3`（版本行）、`:111`（§0.3 v1.64 行）、`:608/612`（§6.3 G1/G5 行）、`:629`（脚本化状态）、`:1171`（§12.4 OB-138 行） |
| 8 | 报告骨架 | `docs/validation/p60-directive-judge-blindspots{,_EN}.md`（§1/§2/artifacts 围栏） |

## 2. 逐项结论

### 2.1 J2 只加子项 — PASS

- `g1.js`（numstat 39+/1−）：唯一删除行是 `parseLedger` 的旧 `return {...}` 被含 `inblock` 的新 return 替换；G1(c)/(d)/(e)/(f) 收集与判定零改动——`rec('c'/'d'/'e'/'f', …)` 行在 HEAD 与工作区各恰 1 处且 diff hunk 未触及；`module.exports` 未变。
- `g5.js`（numstat 91+/1−）：唯一删除行是 `_internals` 单行导出改为多行（新增 `evidencePathsFromInBlockRow` 导出，供夹具引用，非判定改动）；四个 hunk 落点全部在 G5(d)/(e)/(f) 代码区之外，G5(d)/(e)/(f) 判定 +0/−0。
- **`evidencePathsFromRow` 字节级同一**：HEAD 与工作区各自抽取函数体，sha256 双侧均为 `482c3a19e451f1eafaf54ea35da8ebb5555611db971293fc9861f83ebfc2dd18`，长度同为 1146 字节。
- §6.3 的 G1/G5 行（CN/EN 四行）与脚本化状态行（CN/EN 两行）：程序化核验均为**纯插入**——首个差异位即 (e2) 段起点，删除插入段后与旧行全等；脚本化状态行为两处加粗短语内追加「（含 G1(e2)，v1.64 起）」「（含 G5(e2)，v1.64 起）」。

### 2.2 J3 只翻一行 — PASS

`docs/DELIVERY_DIRECTIVE.md` 的 OB 表 hunk 仅 `@@ -1167,7 +1168,7 @@`，只改 OB-138 一行（观察中→**已处置** + 落地注记「J2 覆盖未来切片；P40–P58 历史块内完成行由人工承担（J1 未选，§11.1 不争议）」）；EN 同位（Observing→Resolved）。§12.4 其余 OB 行零改动。

### 2.3 v1.64 授权标注 — PASS

- CN:111 `| v1.64 |` 行末句「本片升版由用户 2026-11-06 的 P60 裁定（J2）授权，非切片内自行决定」；EN:111 同位「The version bump of this slice is authorized by the user's P60 ruling of 2026-11-06 (J2), not decided inside the slice.」。
- 头部版本行 CN/EN:3 均 v1.64；G1-b 机械读数独立复核：`header=v1.64 last=v1.64 rows=65 OK`（双侧）。

### 2.4 向前规则 — PASS

- G1(e2)：`added = ctx.addedLineNumbers(p.rel)`；`!added || !added.has(c.n)` ⇒ 计入 historical、绝不判定（`g1.js:224-226`）。
- G5(e2)：仅遍历 `ctx.addedLines(rel)` 的新增行（`g5.js:190-197`）。
- 独立复跑读数：`G1(e2) :: 0 judged, 0 skipped, 40 historical (v1.64)`、`G5(e2) :: no added in-block completion row`——P40–P58 的 40 条历史块内完成行保持不判，符合 OB-138 目标与裁定【四】1。

### 2.5 ③ 观察项定性（⑤ 意见，供 ⑧/⑩ 登记）

- **O-1（Low）**：属实。独立复核：共享提取器对 `REMAINING_SLICES.md:1600` 提取出 `bounddir_unix.go`、`:1618` 提取出 `red-team-sol-review-r2.md`（散文 basename 误判为引用）。定性：**登记为观察项/后续片候选，本片维持现状**——缺陷来自 `evidencePathsFromRow`（本片零 diff），修复共享提取器会改变 G5(e) 判定语义、超出 J2 授权；向前规则下对历史行不判、当前无触发面；未来新增此类行时以 FAIL 显性暴露（非静默）。
- **O-2（Low）**：不可展开行「可见跳过 → PASS」**可接受**。理由：F-2 四路 detail 使该状态显性可见（`N added in-block row(s) but no parsable evidence path (v1.64)`），计数留痕、绝不静默；与 G5(e) 既有「nothing to resolve」口径一致。建议（不堵本片）登记为后续硬化候选：可选把「有新增行而证据列完全不可解析」升级为 FAIL/WARN。
- **O-3（Info）**：已消除——独立复跑 `G5(f) index hygiene :: index equals worktree, no untracked path` PASS（selftest.js 已由主控暂存）。

### 2.6 报告合规 — FAIL（Medium，⑧ 补写）

报告 §1 异常登记区仍为「暂无（占位）」，状态行仍为「阶段 A 完成…未进入 ②」，未写 J1 不选事实。按裁定【四】3，⑧ 装配必须补写，现成文案：

> J1 需触碰 P40–P58 历史行，§11.1 争议 + 用户明示授权成本，本片不选；J2 的向前规则 + J3 的历史记录态，已使 OB-138 闭环。

⑧ 同时需同步的正常装配项（非新发现）：状态行（🔄 进行中 → 阶段 B/完成态）、准则版本 v1.63 → v1.64、SELFTEST 数字 703 → 722、G7-f sixstage 围栏（裁定七，⑧ 转录 ⑥ 输出后闭合）。

### 2.7 编码与对称 — FAIL（Medium，test-r2.md）/ 其余 PASS

- **`03-test/test-r2.md` 违反硬规则**：UTF-8 **无 BOM** + **CRLF×71**（首字节 `35 32 226`；prfrail 规则：`.md` 必须 BOM+LF）。其当前字节 sha256 `4515e911…` 与 SHA256SUMS/MANIFEST 记录一致（G8-b 自洽 PASS、G8-c 原始=归档 PASS），故机械门禁全绿——G4a 将证据树排除在作用域外（`21 file(s) OK (10 excluded)`），G8-c 只查哈希一致性、不查规范符合性，此偏差只能由人工层捕获。修法（机械归一，不改任何语义）：BOM+LF 归一；MANIFEST「后缀与编码映射」与「归档条目清单」的原始/归档 sha256 分别更新为归一前/后两值；SHA256SUMS 更新为归档字节哈希。
- 其余全部通过：其余 `.md`（含 P60 证据与报告骨架）BOM+LF ✓；`g1.js`/`g5.js`/`selftest.js` 无 BOM+LF ✓；`SHA256SUMS.txt` 无 BOM+LF ✓。
- 双语对称：DELIVERY_DIRECTIVE CN/EN numstat 均 6+/5−（G3(a) PASS）；六形状 G3(a) 机械判据独立复核 PASS（`lines/head/bold/bars/blank/backticks` 六项全 true）。

### 2.8 白名单 — PASS

`tools/gates/g6-whitelist.txt`、`tools/gates/cjk-newwords.txt` 暂存 diff 均 0 字节；G4b 独立复核 `newCjk=467 unknown=[]`。

## 3. 独立复跑读数（⑤ 不依赖 ③ 自证）

- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 722/722`（exit 0，无 not ok）。
- `node tools/gates/gate.js --all --scope=index` ⇒ `TOTAL_FAIL=1`，唯一红项为 G7-f sixstage 围栏（裁定七的预期项，⑧ 装配后闭合）；base 四门禁全绿（gofmt/build/vet/test，`ok=15 FAIL=0`）。
- 夹具增量核对：703 → 722 = T119(4) + T121(5) + T120(3) + T122(2) + T123(3) + T124(2)，与 r1/r2 回执逐条一致。
- 备注：G9-e INFO 2 处（台账引用 `tmp/P60-DIRECTIVE-JUDGE-BLINDSPOTS` 未暂存，目录内仅余主控 `norm.js`），信息性、不阻断。

## 4. 发现清单（Medium+）

| # | Severity | 位置 | 问题 | 修法 |
|---|---|---|---|---|
| M-1 | Medium | `docs/validation/p60-directive-judge-blindspots{,_EN}.md` §1 | 未写「J1 不选」事实与理由（裁定【四】3；报告尚为阶段 A 骨架、⑧ 未装配） | ⑧ 补写 §2.6 现成文案，并同步状态行/版本/SELFTEST/G7-f |
| M-2 | Medium | `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/03-test/test-r2.md` | 无 BOM + CRLF，违反 `.md` BOM+LF 硬规则（G4a 不覆盖证据树、G8-c 只查一致性，机械门禁全绿） | BOM+LF 机械归一 + MANIFEST/SHA256SUMS 双哈希更新（§2.7） |

无 Critical；无判据逻辑缺陷。

## 5. 声明

- 本次仅审查：未改任何被审文件；未 `git add` / commit / push；未派单；临时验证脚本 `tmp/p60-pre-review-verify.js` 与门禁输出 `tmp/p60-st.txt`、`tmp/p60-gates.txt` 已清除。
- 本回执（本文件）为 ⑤ 的唯一新落盘物；按证据包惯例待主控/⑧ 注册进 MANIFEST 与 SHA256SUMS。
