# input-list — P56 `G4B-CORPUS-DECOUPLE` 输入包与前置核对

建包时点：2026-11-05（步骤 0）。本文件记录**锚点核对**、**输入清单**、**前置核对结论**与**改动前基线读数**。

## 1. 锚点核对（快照 v45 第一节 5 项）

| 项 | 快照值 | 实测值 | 判定 |
|---|---|---|---|
| `HEAD` | `c76c6ce` | `f021af1` | 已前进 1 个提交（`f021af1` = v45 快照提交），按起跑消息【〇·零 2】以 v45 提交号为锚点 ⇒ 起跑锚点 = `f021af1` |
| `origin/main` | `c76c6ce` | `f021af1` | 与 `HEAD` 同点（无未推送提交） |
| 准则版本 | v1.60 | v1.60 | 一致 |
| `SELFTEST` | 697/697 | `SELFTEST: PASS 697/697` | 一致（改动前实测） |
| 三作用域 | 全 `TOTAL_FAIL=0` | 起跑前工作区干净（`git status --porcelain -uall` 为 0 项） | 一致 |

`tmp` 顶层固定断言：实测 = `.gitkeep` + `gates-cache`；`gates-cache` 下仅 `g4b`。

## 2. 输入清单（本片实际读取的源）

| 路径 | 用途 |
|---|---|
| `docs/t027/CHECKLIST_v45.md` | 跨会话快照（锚点、纪律、门禁坑） |
| `docs/DELIVERY_DIRECTIVE.md` | §6.3 的 G4 行（含 G4b 缓存条款与升版义务）、§5.1 的 `tmp/` 与运行时缓存例外、§11.2 的治理片回写口径 |
| `tools/gates/lib/ctx.js` | `listAllMarkdown()` / `showAt()` / `changedMarkdown()` 的实现与作用域规则 |
| `tools/gates/lib/criteria/g4b.js` | 语料构建、缓存键、缓存可用性与升版常量 |
| `tools/gates/selftest.js` | 既有 G4b 断言（T87 / T108–T113）与 `makeCtx` 的既约夹具 |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | 台账体例与 P54 / P55 条目形态 |
| `docs/validation/exit-seed-build-evidence{,_EN}.md` | 报告体例（`artifacts` 围栏、`sixstage` 围栏、审查链小节） |
| `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md` | 证据包体例 |

## 3. 前置核对：判据文本与实现的当前口径

1. **判据文本**（`docs/DELIVERY_DIRECTIVE.md` §6.3 的 G4 行，CN）：`(b) **异常字扫描**：新增内容的 CJK 字符集合 ⊖ HEAD 全库 \`.md\` 语料 = **差集**` ⇒ 文本口径为 **HEAD 全库语料**。
2. **实现现状**（`tools/gates/lib/ctx.js` 的 `listAllMarkdown()`）：`range` / `ci` 用 `git ls-tree -r --name-only <baseRef>`；**tree / index 用 `git ls-files`（索引枚举）** ⇒ 与第 1 条的文本口径**不一致**，即 OB-120 的根因。
3. **缓存升版义务**（同节 G4b 缓存条款）：`ctx.listAllMarkdown` / `ctx.showAt` / git 封装的**枚举或读取语义变化**时，必须同步升级缓存格式版本前缀，旧条目一律视为陈旧。
4. **结论**：本片是**实现向文本对齐**，不是判据语义变更 ⇒ §6.3 的 G4b 判据文本**本体无需改写**（详见报告 §8 的文档同步声明）；但按第 3 条，`CACHE_VERSION` 必须升版，否则旧语义下的缓存条目仍会命中。

## 4. 影响面（批量辅助第 1 次实派结论 + 主控独立复核）

批量辅助：`deep-reasoner` + `GLM-5.3-Flash (glm)`（只读，末行 `BATCH: DONE`）。其结论与主控 `grep_search` 复核一致：

| 项 | 结论 |
|---|---|
| `listAllMarkdown` 定义 | `tools/gates/lib/ctx.js:536`（导出点 `:693`） |
| 唯一生产调用点 | `tools/gates/lib/criteria/g4b.js:234` |
| 自测 mock | `tools/gates/selftest.js:383` / `:452` / `:6199` |
| 受影响判据 | **仅 G4b**（语料源改变 + 缓存键需升版） |
| 不受影响判据 | G1a / G1b / G2 / G3a / G3b / G4a / G5 / G6 / G7 / G8 / G9（依据：它们经 `changedMarkdown` / `changedPaths` / `readText` / `listTreeFiles` 取数，`listTreeFiles` 早已使用 `git ls-tree -r -z --name-only <headRef>`） |

**边界澄清（语义差异，逐条登记）**：
- **暂存删除 `.md`**：旧实现使语料收缩 ⇒ 冷跑从严多报；新实现以 HEAD 为源 ⇒ 与真值一致。
- **暂存新增 `.md`**：旧实现会枚举到该路径，但 `showAt('HEAD', p)` 返回 `null` 而被跳过；新实现根本枚举不到 ⇒ **判定结果不变**（仍不在基线与语料中）。
- **非 ASCII 路径**：`git ls-tree --name-only` 在 `core.quotePath` 下会加引号转义；`listTreeFiles` 已以 `-z` 解决同类问题。本片在 `listAllMarkdown` 同样采用 `-z`，避免把 `git show` 的路径弄错而静默丢字符（HEAD 当前无此类路径，属防御性一致性）。

## 5. 改动前基线读数（同机、同终端、串行）

| 读数 | 命令 | 值 |
|---|---|---|
| `SELFTEST` | `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 697/697`，墙钟 278653 ms |
| G4b `tree` 冷跑 | 清空 `tmp/gates-cache/g4b` 后 `--check=G4b --scope=tree` | 129151 ms |
| G4b `tree` 热跑 | 紧接上一条再跑一次 | 1554 ms |
| `--all --skip-base --scope=index` | 判据仅跑 | 2790 ms（`TOTAL_FAIL=1`，原因为工作区未暂存，属步骤 0 中间态） |
| `--all --skip-base --scope=tree` | 判据仅跑 | 6150 ms（`TOTAL_FAIL=1`，原因为新证据包尚无 `MANIFEST` / `SHA256SUMS`，属步骤 0 中间态） |
| `--all --skip-base --scope=ci` | 判据仅跑 | 14958 ms（`TOTAL_FAIL=0`） |

> 口径：两个中间态失败均在步骤 0 内闭合（台账与报告围栏落盘后重跑），其原始输出按轮归档于 `tmp/G4B-CORPUS-DECOUPLE/`。

## 6. 缺件断言

输入包生成时逐件断言存在：上表第 2 节列出的**全部**输入均可读（本文件由主控以只读方式读取，未发现缺件）；审查输入包另行按 §5.1 的时序说明逐件断言。
