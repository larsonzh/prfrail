# ⑤ 预审报告（r1）逐字原文 — `G4B-CORPUS-DECOUPLE`

性质：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)` 的**最终消息逐字转录**（主控不按要点重建）。审查时点 2026-11-05，只读。

---

# ⑤ 预审报告 — `G4B-CORPUS-DECOUPLE`（只读）

所有读数来源均为「读取主控落盘的文件路径」（本报告未执行任何命令）。

## 总体结论

改动与判据文本（DELIVERY_DIRECTIVE §6.3 G4 行，已逐字比对 `docs/DELIVERY_DIRECTIVE.md:602`）方向一致：实现向「HEAD 全库 `.md` 语料」对齐，缓存条款（键含版本前缀、读前重算键、异常 fail-open、枚举语义变 ⇒ 升版）逐项满足。A-1..A-5 全部成立；3 项 Low（证据包完备性），无 Medium+。

## 逐条判定

- **A-1 成立**。`ctx.js:545-546`（读取 `tools/gates/lib/ctx.js`）：tree/index 走 `git(['ls-tree','-r','-z','--name-only', 'HEAD'])`，`listAllMarkdown` 内已无 `ls-files`；T114 断言钉住「ls-files 从未被调用」。
- **A-2 成立**。T115（`selftest.js:6593-6612`）冷 miss `corpus=1` / 热 hit `corpus=1` 且双 PASS；实测 `probe-staged-delete-g4b.txt`（`miss 324f2aa18be2 corpus=1429`、`unknown=[]`）；A/B 探针 `ab-probe-adr.txt`（`corpus_old=1428` 出现 `U+6602` unknown vs `corpus_new=1429` 差集空）——语料不收缩、冷热一致。
- **A-3 成立**。`g4b.js:45` `CACHE_VERSION='v2'`、`:152` `validValue v===2`、`:243` put `{v:2}` 三方同步；键首段 = `CACHE_VERSION`。实证：改后冷跑键 `324f2aa18be2` ≠ 基线键 `05793d94629d`（HEAD 未变仍 miss ⇒ 旧 v1 条目未命中）。
- **A-4 成立**。全库检索：`listAllMarkdown` 生产调用仅 `g4b.js:236`（selftest 三处 384/452/6200 为 mock，6579 为 T114）。
- **A-5 成立（等价或更强）**。见下节结构化说明。

## A-5 等价性结构化说明（T115 之外的三处绑定）

- 覆盖面：① 键首段绑定 `CACHE_VERSION`（T112/T113 shape）；② 重写条目的 `v` 绑定常量数值（T112）；③ put 位 `{v:2}` 被 T116 直接读取验证；④ `validValue` 漂移被 T115 的 hit 预期间接钉住（v 错配 ⇒ 热跑变 miss ⇒ 红）。
- 遗漏面：⑤ `validValue` 内硬编码 `2` 与 `CACHE_VERSION` 无直接机械绑定（仅经 T115 间接）；⑥ `ctx.js` 枚举语义 ↔ `CACHE_VERSION` 的耦合无法机械断言（升版义务为人工义务，先验结构限制）。
- 结论：三项替换在 v2 下数值等价，且随版本漂移自动跟随，**不削弱既有断言**；遗漏面不构成回归。

## 独立复算（§4.2.1 ① ②）

- 行数核对：`ctx.js` 773、`g4b.js` 283、`selftest.js` 6660、报告对 56/56、`input-list.md` 70、`freedom-list.md` 39、`REMAINING_SLICES.md` 1855、`_EN.md` 1842 —— 与 `input-manifest.txt` 全等（逐件读到 EOF 确认）。
- `diff.patch` ctx.js 块：头部 `@@ -537,11 +537,16 @@` 与块体一致（旧 3 上下文+5 删+3 上下文=11；新 3+10 增+3=16）；包内无独立 `--stat` 行，以块头算术自洽 + 磁盘文件与补丁后镜像逐字一致为准。

## 发现（3 项 Low）

- `Low | tmp/G4B-CORPUS-DECOUPLE/after-selftest.txt:700 | 「drops the only "涓? and turns this red」 | 控制台 GBK 捕获乱码；源码 `selftest.js:6518` 为 `'\u4e19'` 转义、无代码缺陷。最小修复：重捕获时 `chcp 65001` 或证据说明标注捕获编码。`
- `Low | tmp/G4B-CORPUS-DECOUPLE/after-selftest.txt | 「改动后 SELFTEST: PASS 701/701（墙钟 219458 ms）」 | 包内无任何 after 耗时行（基线 278653 仅见于 `input-list.md:52-57` 表）。最小修复：③ 报告或捕获文件补登记 after 墙钟，满足台账「改前/改后耗时对比表」验收项。`
- `Low | tools/gates/selftest.js:6577-6613 | 「a staged .md deletion must not move the corpus under the cache」 | 新断言仅覆盖 tree 作用域；index 作用域同构等价与 range/ci 的 `-z` 形态切换无专用断言（T113 已 mock 掉 `listAllMarkdown`）。最小修复：可选，补一条 index 作用域同构断言。`

## 反例分析（可证伪判据）

1. tree 分支改回 `git(['ls-files'])` ⇒ T114 红（calls 含 `ls-files`、列表少 `ENUM_DELETED`）+ T115 两检红（`corpus=0`、verdict FAIL）——钉 A-1/A-2。
2. `CACHE_VERSION` 回退 `'v1'` ⇒ T116 红（stored `v=2` ≠ 1）+ T112 重写检红——钉 A-3。
3. `validValue` 的 `v===2` 改回 1 ⇒ T115 热跑变 miss，`hit` 检红——钉 A-2/A-3 联动。
4. 删 `-z`（改回 `--name-only` + `\n` 切分）⇒ fake 仅回答 `-z` 形态 ⇒ 列表空、T114/T115 红——钉 FL-2。

## 无法验证的部分

- 各文件 bytes 与 `diff.patch` 的 23411 字节（只读工具不提供字节数，未执行命令）。
- 墙钟 278653 / 219458 ms（捕获文件均无时间行，见发现 2）。
- 三作用域改后复跑 `TOTAL_FAIL=0`（包内无 after-all-* 捕获；仅 after-selftest、after-g4b-tree 与 probe）。
- `SHA256SUMS.txt` 三行哈希与载荷逐字节一致性、BOM/LF 编码门禁（未执行计算）。

PRE-REVIEW: PASS WITH FIXES
