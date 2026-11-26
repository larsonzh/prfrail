# P92 分区计时实测数据（`--selftest`）

- 生成：2026-11-26；方法 = 进程内加载 `selftest.js`（内存改写 `check()` 打点 + 包裹 `fs` / `child_process`），不落盘改动。
- 原始件：`tmp/P92-GATES-SELFTEST-SPLIT/profile.txt`（第二次运行，含子进程归因）、`profile-run1-nospawnbreakdown.txt`（第一次运行）、`profile-memo.txt`（记忆化变体）。

## 1 总账（改动前实现）

| 口径 | 第一次运行 | 第二次运行 |
|---|---|---|
| 墙钟 ms | 160323.3 | 150934.1 |
| checks / passed | 744 / 744 | 744 / 744 |
| stdout sha256（无尾 LF） | `807b1203…4388a` | `807b1203…4388a` |
| 子进程次数 / ms | 483 / 157884.6 | 483 / 148653.0 |
| 进程内磁盘读 次数 / ms / 字节 | 1036 / 440.6 / 211 MB | 1036 / 407.0 / 211 MB |
| 进程内 CPU 余量（推算） | ≈ 2.0 s | ≈ 2.1 s |

## 2 子进程构成（第二次运行）

| 命令族 | 次数 | ms |
|---|---|---|
| `git rev-parse --show-toplevel` | 413 | 75797.8 |
| `node tools/gates/gate.js …`（CLI 腿） | 4 | 60274.3 |
| `git ls-files --error-unmatch` | 38 | 7345.4 |
| `git rev-parse --verify` | 17 | 3220.5 |
| `git diff --name-status` | 4 | 763.9 |
| 其余（`cat-file` / `show` / `status` / `rev-parse HEAD`） | 7 | ≈ 1.1 s |

## 3 进程内读构成（Top）

| 文件 | 次数 | ms |
|---|---|---|
| `REMAINING_SLICES.md` | 126 | 78.9 |
| `REMAINING_SLICES_EN.md` | 126 | 61.9 |
| `DELIVERY_DIRECTIVE.md` | 79 | 57.4 |
| `DELIVERY_DIRECTIVE_EN.md` | 79 | 52.8 |
| `g6-4-expiring-literals.txt`（夹具） | 43 | 12.0 |
| `g7-c-finding-row.txt`（夹具） | 17 | 6.6 |
| `g6-1-pending-done.txt`（夹具） | 12 | 5.1 |
| 其余夹具单件 | — | ≤ 5 |

⇒ 夹具读合计 < 60 ms（约占总墙钟 0.04%）。

## 4 记忆化变体（内存改写，未落盘）

| 口径 | 数值 |
|---|---|
| 墙钟 ms | 70106.5 |
| 子进程次数 / ms | 71 / 68299.9 |
| `git rev-parse --show-toplevel` | 1 / 191.7 |
| stdout sha256（无尾 LF） | `807b1203…4388a`（与基线相同） |
| checks / passed | 744 / 744 |

## 5 实机落地后复测（② 之后，规范入口 `gate.js --selftest`）

| 运行 | 退出码 | 墙钟 ms | 捕获件 sha256（含尾 LF） |
|---|---|---|---|
| B1 | 0 | 75375 | `36af71e5…bd9d2f` |
| B2 | 0 | 69327 | `36af71e5…bd9d2f` |
| A（改动前，含尾 LF） | 0 | — | `36af71e5…bd9d2f` |

> 口径说明：`807b1203…` 是内存 `res.stdout`（无尾 LF）的哈希；`36af71e5…` 是落盘件（含尾 LF）的哈希。二者对应同一次输出，差异恰为 1 个尾 LF，不构成内容漂移（③ 已逐字节复核）。
