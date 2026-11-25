# P89 ① 架构方案书：影响判定 + 归档方案（`DeepSeek V4 Pro`，逐字原文）

## 0. 补充核验（步骤 0 事实之上，新增确认）

- `p85-worktree.patch` 与 `05-input/diff.patch` **字节同一**（SHA256 `E75C5783…`）；且与**已入库**的 `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/P87-four-read-data/diff.patch` **三者同字节**。⇒ 冻结原件已在库，tmp 副本可安全归一。
- 110 件合计 2,785,709 字节；`attempt1-unzip/` 内含两个**带空格子目录**（`Go ubuntu-latest`、`Go windows-latest`），顶层亦有带空格文件名（`0_Go windows-latest.txt` 等）。
- 门禁实现层面（`tools/gates/lib/criteria/g8.js`）：`SUM_RE=/^([0-9a-f]{64})[ \t]+(.+)$/` 与 `BACKTICK_ROW_RE=/^\|\s*`([^`]+)`/` 均**兼容路径含空格**；`G8-c` 标题按前缀匹配（`## 后缀与编码映射`），无该节即豁免；`G8-b` 载荷集 = 包内除根级 `MANIFEST.md`+`SHA256SUMS.txt` 外全部文件；`G8-d` 只对**头部注记区**含「重建归档」的变更文件生效。
- `g3a.js`：G3(a.1) 对每个被改 CN/`_EN` 对只要求 **numstat 增/删行数相等**；整文件形状全等（a.2）**仅对** `docs/DELIVERY_DIRECTIVE.md` 强制。⇒ 本片三对文件的对称 = 「同位置同量替换」，不要求全文行数相等。
- 起跑消息原文（`P89…/00-briefing.md:95-138`）验收句精确措辞：「零命中（**除 ① 判定的更正注记处，注记措辞不得含可解析完整路径或按 G9 白名单处置**）」；§七必改 #3 明确「台账 P85 段落点声明字段（:2314/:2301）」为**必改**；§六把「报告注记」归入**活跃引用改写**类。

## D1 影响判定：不触 CONTRACTS / 契约 / CI，不需准则 MINOR

| 判定 | 结论 | 依据 |
|---|---|---|
| CONTRACTS / schema / 契约夹具 / workflows / CI workflow | **不触** | 变更面仅 `.md`（证据包 / 台账 / 报告对）＋ `internal/console/p85_t2_frozen_policy_tester_test.go:6` 一处**注释行**；`grep` 证实两 tmp 路径在 `docs/CONTRACTS{,_EN}.md`、`schemas/`、`.github/workflows/` 零命中 |
| `internal/**` 生产逻辑 | **不触** | 仅测试文件注释；P89 边界字段（`REMAINING_SLICES.md:2402`）已授权该注释行改动 |
| 准则 MINOR 版本行 | **不需要** | 两路径在 `DELIVERY_DIRECTIVE{,_EN}.md` 零命中 ⇒ 准则正文一字不改 ⇒ §0.3 变更日志义务不触发；P89 边界「不碰准则正文（除非 ① 判定需 MINOR）」的例外条件不成立 |
| 证据模型 / 门禁判据语义（§7.2.1 ③） | **不触** | 本片操作**证据包内容**（数据），不动 `internal/evidence` 模型、不动任何 G 判据实现或判定逻辑 |
| ⑥ 跳过 | **不主张** | P89 步骤表含 ⑥ 且改动含 `.go`（注释），不满足 §7.2 ②；按步骤表执行，新报告配 `sixstage` 围栏（G7-f） |

## D2 归档布局裁定

**唯一推荐**：

1. **落点**：`docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/`（43 件）＋ `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/`（67 件）。保留 `archive/` 前缀 + 原目录名——先例 `directive-history/master-subagent-architecture/archive/…`（§1.5.1 ① 归档实作），语义「归档载荷根」；G8-b 把 `archive/**` 当载荷 ⇒ 全部入 MANIFEST + SHA256SUMS（G8-b 强制，无豁免）。
2. **带空格目录名**：**逐字保留**。理由：SUM_RE / BACKTICK_ROW_RE / `normalizeRel(toPosix)` 均兼容空格；改名将制造多余「改名归一」映射行、削弱字节保真。执行约束：复制/哈希脚本一律 `-LiteralPath` 或 node 直写，禁 PowerShell 通配解析。
3. **zip 二进制**：**逐字保留**（41445 B 载荷；无归一、无映射行）。
4. **`p85-worktree.patch` ≡ `05-input/diff.patch`**：**逐字两份均归档**。理由：① §1.5.1 ①「归档与引用改写同步」要求源→归档 1:1 路径对应，删名即破坏改写表可核验性；② 两文件名被不同历史件按名点名；③ 二者与已入库 `P87-four-read-data/diff.patch` 逐字节同源，MANIFEST 说明列登记该同源关系即可，不另立「只留一份」新口径。备选 B（只留一份＋登记等价）否决：引入「源路径未逐字归档」新口径且省 776 KB 收益不值。

## D3 归档副本内容处置（11 文件 39 处自引用）：**路径字面归一化**，不逐字保留

**唯一口径**：把 11 件归档副本内的 `tmp/P84-CI-FIX/` → `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/`、`tmp/P85-SECURITY-PROOF-SUBSYSTEM/` → `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/`（保留 `#L38-L77` 等锚点），并用 **G8-c「后缀与编码映射」节登记原始 vs 归档两枚 sha256**（说明列写「路径字面归一，N 处」；**不得**在说明列复现旧完整路径）。

理由：① 逐字保留会让 `git grep` 命中 39→117，白名单失控，直接违背验收；② 这些副本是**本次落盘的新载荷**，落盘形态属本片裁量（同 `.md`→`.txt` 后缀归一先例，`master-subagent-architecture/MANIFEST.md`），不属 §11.1「历史原文」；③ 字节保真由 G8-c 双 sha 记录承担，证据完整性可复算；④ 逐字原件仍在库（`P87-four-read-data/diff.patch` 等已入库冻结件），信息零损失；⑤ 归一后引用指向**真实存在**的归档路径，可解析性反而增强。

## D4 活跃引用逐处改写（目标写法；`_EN` 由 ③ 等量镜像）

| 位置 | 处置 | 目标写法（要点） |
|---|---|---|
| `internal/console/p85_t2_frozen_policy_tester_test.go:6` | 改写 | `// The design (docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/01-design.md §2.4-§2.6) names` —— 仅换路径、其余逐字不动；不触任何测试语义 |
| `REMAINING_SLICES.md:2314` / `_EN:2301`（P85 段落点声明） | 改写（起跑 §七 #3 必改） | 标注「已归档至 `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/`（归档执行 = P89，2026-11-25）」—— 不含旧完整路径 |
| `REMAINING_SLICES.md:2395` / `_EN:2382`（P89 目标字段） | 改写 + §ref 校正 | ② 改为「按 **§1.5.1** 判定历史文件处置；回写按 **§11.2**」（即起跑 §七 #3 的「§11.1 改指 §1.5.1 + §11.2」） |
| `bounddir-residual.md:156` / `_EN:156` | 改写（起跑 §六「报告注记=活跃引用」） | cell 首句改为「已归档至 `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/`（归档执行 = P89，2026-11-25）」；工具链 / `%TEMP%` / 缓冲补正句逐字保留 |
| `bounddir-residual.md:160` / `_EN:160` | 改写 | 括注改「（OB-37 追加登记原件已归档至 `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/10-final-review-request.md`）」 |

纪律：每对 CN/EN 同位置同量替换，保证 G3(a.1) numstat 对称；落盘后跑 G6-4 预检（禁 `NN 行`/`NN 个文件`/`共 NN`；用「合计 N 件」）。

## D5 冻结面处置 + 「零命中」操作化定义

**逐类裁定**：

| 位置 | 裁定 | 理由 |
|---|---|---|
| P85 包 `00-briefing.md:109` | **原文保留 + 文末追加归档注记**（注记不含旧完整路径；随之更新该文件 SHA256SUMS 行） | §13.1.10 逐字原文义务 + 起跑 §六「倾向加更正注记（不改历史原文）」；注记只追加、不改原行 |
| `07-carbon-r7-targeted-read.md:10`、`07-codex-baseline-read.md:11,14,21,26` | **保留原样、不注记** | §13.1.13 审查方产出逐字原文归档、主控不得重建；文件内加注即破坏逐字与包哈希 |
| `P87-four-read-data/diff.patch:237,11818` | **保留原样、不注记** | 补丁载荷追加会破坏 patch 语法；且同字节原件即其自身——对应关系由 MANIFEST + P89 报告承担（与起跑「倾向注记」的偏差 ⇒ 需 ⑦ 背书） |
| `CHECKLIST_v72{,_EN}.md:137` | **保留原样** | v72 快照冻结，起跑 §七 #7 明示「v73 再吸收」 |

**「零命中」机械判据（② 脚本可复跑）**：

1. `git grep -n -F "tmp/P84-CI-FIX" -- ':!tmp'` 与 `…"tmp/P85-SECURITY-PROOF-SUBSYSTEM" -- ':!tmp'` 的实际命中集，排序后与**白名单 W** 逐字节相等 ⇒ PASS；任何多出 ⇒ FAIL。
2. 白名单 W（随报告落盘，逐条 `path:line:理由`）：W-A 冻结历史原文；W-B 快照冻结；W-C 归档来源登记；W-D 本片报告对 + P89 证据包自身登记行。
3. 与起跑字面的差异：W-A/W-B/W-C/W-D 不全是「更正注记处」⇒ **需 ⑦ 背书该白名单化口径**（起跑已预授权「最终口径 ① 判、⑦ 背书」）。
4. `find-quotes.js` 不扫描（P89 边界排除；命令 `-- ':!tmp'` 后其位于 `tmp/` 内本就不在命中集，但报告须明示该排除）。

## D6 元数据机制

- **归档条目清单**：两包 MANIFEST 各增 43 / 67 行（`| \`archive/<原名>/<rel>\` | 说明 |`），说明列可写「字节同源 / 编码归一 / 路径字面归一」概要。
- **对应关系登记**：两包 MANIFEST 文首要点区各加 1 行「**归档来源登记（P89，2026-11-25）**」（§1.5.1 ①「只由归档 MANIFEST 记录对应关系」的字面履行）。
- **`## 后缀与编码映射`（G8-c）**：两包 MANIFEST **启用**该节。行 = 编码归一变更文件 + 11 个路径字面归一文件；5 格 `原路径 | 归档路径 | 原始 sha256 | 归档 sha256 | 说明`；说明列**不含旧完整路径**（防新增命中）。归档值必须 == SHA256SUMS 记录（G8-c 硬校验）。
- **不新增自定义 `##` 节**：G8 只解析三个已知标题，新节虽被 `sectionBody` 截断不会误判，但无必要。
- **SHA256SUMS 生成口径**：POSIX 相对路径（相对包根）、`sha256<空白>path`、按路径排序、原子重建（先写临时文件再 rename）；两枚根元数据不列入载荷。

## D7 门禁与检查面清单（归档+改写后全部执行）

1. `node tools/gates/gate.js --selftest` = **744/744**。
2. 归档自洽三方相等脚本（载荷集 == MANIFEST 条目集 == SHA256SUMS 条目集；逐条 sha256 复算）——两目标包 + P89 包。
3. 精确 `git add`（禁 `-A`）后：`--all --scope=tree` 与 `--scope=index` **TOTAL_FAIL=0**；本地 `ci` 读数为上一提交（观察口径，权威读数在 CI）。
4. G8-a：全部新增证据树路径 ∈ {P84 包, P85 包, P89 包}；G8-b 三方相等；G8-c 映射行 5 格、64-hex、归档 sha == SHA256SUMS；G8-d 无触发（注记避开「重建归档」首词）。
5. G4a（证据树外变更件）：报告对 / 台账对 BOM+LF；`.go` 无 BOM+LF；G4b 预检新 CJK 字（必要时登记 `tools/gates/cjk-newwords.txt`）。
6. **G3(a) 特别说明**：G3(a.1) 只要求 CN/EN **numstat 增删行数相等**，整文件行数相等仅强制于 `DELIVERY_DIRECTIVE`（a.2）。
7. G9：两 tmp 目录为**未跟踪删除**，不进变更集；P89 报告改写表对旧路径的引用在非证据树 ⇒ G9-e 仅 INFO（不判红）。
8. G6：G6-4 预检；G6-5 新增 §X.Y 引用可解析；G7-f：新报告带 `sixstage` 围栏（⑥ 实跑，结论行闭集）。
9. `gofmt` / `go build ./...` / `go vet ./...` / `go test ./...` 全绿（测试注释改动后）。
10. 提交后 CI 双腿：`SELFTEST` / `scope=ci changed=N` / `TOTAL_FAIL=0`；包级 `ok` 行取证。

## D8 ② 可执行清单（有序、可脚本化）＋ 回滚

0. 锚点复读。
1. **备份**：整树复制两 tmp 目录 → `tmp/P89-TMP-ARCHIVE-P84-P85/src-backup/`；记录 110 件 sha256 清单。
2. **复制**：备份校验通过后，复制到两包 `archive/<原名>/…`；逐件 sha256 == 备份 ⇒ 复制保真。
3. **路径字面归一**：11 文件按预期次数（3/3/4/1/1/2/1/1/1/20/2 = 39）断言式替换两前缀 → 归档路径。
4. **编码归一**（§1.5.1 ②）：全部文本载荷归一「无 BOM、一律 LF」（zip 除外）；记录逐件「原字节 sha / 归档字节 sha」。
5. **元数据重建**：两包 + P89 包 MANIFEST（条目行 + 来源登记行 + G8-c 节）、SHA256SUMS 原子写。
6. **三方相等校验**：脚本比对三集合 + 逐条复算 sha ⇒ 全等才继续。
7. **删除两 tmp 目录**（仅当 2–6 全绿）；`Test-Path` 前后读数取证。
8. **引用改写**：D4 六处（含 00-briefing 尾注）＋ D5 白名单对照；台账 P89 段翻 ✅ + 完成记录。
9. **P89 报告对**：`docs/validation/tmp-archive-p84-p85{,_EN}.md`（逐行镜像 + `artifacts` 围栏 TAB 格式 + sixstage 围栏 + 白名单 W 表）；P89 包补改写对照与归档自洽证据载荷。
10. **暂存 + 门禁**：精确 `git add` 全部变更文件 → D7 第 3–9 项；`git grep` 两路径命中集 == W ⇒ 判据 PASS。
11. **停 ⑩**：等用户授权提交；提交后 CI 双腿取证（§11.2 回填口径）。

**回滚**：未提交前 = `git restore --staged` + `git checkout --` 变更文件 + 从 `src-backup/` 还原两 tmp 目录（备份保留至 ⑧c 才清）；已提交 = `git revert <commit>`（不动历史）；全程不 push gitee。

## 需 ⑦ 背书点

- **B1**「零命中」白名单化口径——与验收字面差异。
- **B2** `diff.patch:237,11818` 不加注记（保护补丁语法）。
- **B3** G8-c 节复用为「路径字面归一」登记（机制原为编码归一）。
- **B4** `07-carbon` / `07-codex` 审查原文保留原样不注记（§13.1.13 逐字归档优先）。

## 不得触碰清单

- 准则正文（无 MINOR）、`docs/CONTRACTS{,_EN}.md`、`schemas/`、契约夹具、`.github/workflows/**`、CI workflow。
- `internal/**` 生产逻辑（仅 `p85_t2_frozen_policy_tester_test.go:6` 注释行可改）。
- `find-quotes.js`、`tmp/` 全清、`CHECKLIST_v72{,_EN}.md`（v73 吸收）、`DEV_PLAN{,_EN}`、`DOCUMENTATION_PLAN{,_EN}`、G5-b 闭集与白名单、ADR、gitee。
- 冻结件（`00-briefing.md` 原文行、四件审查原文、`P87-four-read-data/diff.patch`、P84 包既有载荷）不得改写其原文内容；`SHA256SUMS` 重建不得改变既有条目与既有哈希值。

---

> **实施偏差（P89 记录）**：D5/D6 的 W-C / W-D 两类在实施时被**消除**（MANIFEST 来源登记与报告对照表改用非连续写法 `tmp` 根下 `X`），白名单收缩为「冻结原文 9 行 + 冻结快照 2 行 + 本片起跑消息 1 行」三类 11 行；D8 步骤 4 的「无 BOM」按扩展名惯例更正为「`.md` / `.ps1` 带 BOM，其余不带」（见 `freedom-list.md` F-13）。两项偏差均由 ⑤ 报出并在报告 §2.5 披露，交 ⑦ 背书。
