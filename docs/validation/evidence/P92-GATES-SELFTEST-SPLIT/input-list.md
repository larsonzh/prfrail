# P92 输入清单（input-list）

- 切片：P92 `GATES-SELFTEST-SPLIT`（工具片 `[SLICE]`，M 档）；步骤：0 / ①。
- 生成时间：2026-11-26。

## 1 权威输入

| 输入 | 路径 | 用途 |
|---|---|---|
| 起跑消息原文 | `docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md` | 全片唯一落地依据 |
| 快照锚点 | `docs/t027/CHECKLIST_v75.md` / `_EN.md` §一 / §四 / §六 | 锚点核对与下一项出处 |
| 执行准则 | `docs/DELIVERY_DIRECTIVE.md` §7.5 / §5.2 / §13.1 / §13.2 / §0.2 | 额度口径、档位、步骤 0、版本包 |
| 台账 | `docs/t027/REMAINING_SLICES{,_EN}.md` P92 段 | 字段定义与状态机 |
| P88 建议出处 | `docs/validation/gates-perf-assess{,_EN}.md` §7（建议 1/3/5） | 候选来源与耗时基线 |
| 编码规范 | `docs/CODING_CONVENTIONS.md` | `.js` = 无 BOM + LF |

## 2 改动候选面（只读通读）

| 文件 | 关注点 |
|---|---|
| `tools/gates/selftest.js` | `fixtureText` / `makeCtx` / `injectedCtx` / `discoverRepoRoot(__dirname)` 调用点 / `gateCli` |
| `tools/gates/lib/criteria/g4b.js` | `CACHE_VERSION` / `sourceDigest` / `keyString` / 语料内环 |
| `tools/gates/lib/ctx.js` | `discoverRepoRoot` 语义、`listAllMarkdown` / `showAt` 枚举口径 |

## 3 实测数据（本片自产，`tmp/P92-GATES-SELFTEST-SPLIT/`）

| 文件 | 内容 |
|---|---|
| `profile.txt` | 分区耗时与子进程归因（第二次运行，含 `SPAWN_TOP`） |
| `profile-run1-nospawnbreakdown.txt` | 第一次运行（无子进程归因） |
| `profile-memo.txt` | 仓库根发现记忆化变体的实测（含子进程归因） |
| `selftest-stdout.txt` | 基线 stdout（sha256 与第二次运行相同） |
| `selftest-stdout-run1.txt` / `selftest-stdout-memo.txt` | 基线与变体 stdout 的对照件 |

## 4 机械读数（步骤 0 实测）

- `HEAD = origin/main = 76d7fc21c22603e2f938d88e0476af511982dbe2`。
- 准则头部 = `v1.74`；ADR max = 48；§12.4 OB max = 145；DR max = 18；试点顺延最大 = 第 25 次。
- 滚动窗口：v73 = 212 / v74 = 216 / v75 = 215（按 LF 换行计的内容行）。
- `tmp/` 顶层 = `.gitkeep` / `gates-cache` / `find-quotes.js`。
- `--selftest`：`744/744`、退出码 0；两次运行的 stdout sha256 相同（`807b1203…4388a`）。
