# P88 耗时对照汇总（同机同仓库，2026-11-25）

> 口径：全部为墙钟毫秒（含 Node / git / go 进程启动开销）。机器 = 本机 Windows；仓库 = `d:\LZProjects\prfrail`，HEAD = `fed0f5c`。
> 读数的完整逐项表见 `baseline.tsv`（优化前，含每条命令的退出码）与 `ab-report.md`（优化前后配对）。

## 一 优化前：耗时构成（`baseline.tsv`）

| 项 | 退出码 | 耗时 ms | 备注 |
|---|---|---|---|
| `gofmt -l .` | 0 | 1114 | base 第 1 件 |
| `go build ./...` | 0 | 24503 | base 第 2 件 |
| `go vet ./...` | 0 | 8858 | base 第 3 件 |
| `go test ./...` | 0 | 143628 | base 第 4 件；单件即占 base 合计的 80.6% |
| base 合计 | — | 178103 | 每次 `--all` 串行执行一次 |
| `--all --skip-base --scope=range:HEAD~3..HEAD` | 1 | 490013 | 语料**冷**扫；G4b 逐文件 `git show` 主导 |
| `--all --scope=range:HEAD~3..HEAD` | 1 | 240824 | base 178103 + 判据 62721（语料已热） |
| `--all --skip-base --scope=tree` | 1 | 5247 | HEAD 语料已缓存 |
| `--all --scope=tree` | 1 | 99070 | base 热 93823 + 判据 5247 |
| `--all --skip-base --scope=index` | 1 | 1946 | — |
| `--selftest` | 0 | 265294 | 744 断言 |
| 逐判据（`--check=<ID>`，语料热，含约 1.4–1.9 s 的进程与变更集解析固定开销） | — | 见 `baseline.tsv` | G6 20865 / G7 12525 / G8 10311 / G4a 10236 为最重的四条 |

## 二 优化前后配对（`ab-report.md` 的实测值）

| 路径 | 旧 ms | 新 ms | 加速比 | 口径 |
|---|---|---|---|---|
| A `--all --scope=tree` | 75000 | 79628 | 0.94x | 含 base 四件套；单样本，受 Go 缓存暖度支配，非本片改动面 |
| B `--all --skip-base --scope=index` | 2430 | 1750 | 1.39x | 判据面 |
| B `--all --skip-base --scope=range:HEAD~3..HEAD` | 59272 | 33083 | 1.79x | 判据面（含 O-3c 读数 memo） |
| C `--all --skip-base --scope=ci`（语料**冷**） | 170277 | 12929 | 13.17x | O-1 的直接收益面 |
| C1w 新实现热跑 ci | — | 10137 | — | 跨实现缓存兼容性旁证 |
| D `--selftest` | 209992 | 128881 | 1.63x | 744 断言；两侧结论行同为 `SELFTEST: PASS 744/744` |

**CI 关键路径**（`--selftest && --all --scope=ci`，冷 runner 等价口径）：**380.3 s → 141.8 s（2.68x，−238.5 s）**。

## 三 优化后 ④ 门禁读数（含 base；`verify4` 输出）

| 项 | 退出码 | 耗时 ms |
|---|---|---|
| `gofmt -l .` | 0 | 774 |
| `go build ./...` | 0 | 7220 |
| `go vet ./...` | 0 | 4606 |
| `go test ./...` | 0 | 98978 |
| `--all --scope=tree`（base + 判据） | 1（仅 G8-a，证据包当时尚未建包） | 76085 |
| `--selftest` | 0 | 161118 |

## 四 读数的可复现性边界（如实登记）

1. base 四件套的耗时随 Go 构建/测试缓存的暖度大幅波动（冷 178103、热 93823 为同机两次实测），单次读数不得跨暖度比较；A 组的 0.94x 即该波动的表现，**不构成回退**（`gate.js` 未改动，见 `object-hashes.txt` 的 `IDENTICAL_TO_HEAD` 行）。
2. 逐判据读数含每次调用的固定启动开销，故其合计大于同一进程内跑完全部判据的实测值。
3. 所有耗时读数都只进本文件与 `ab-report.md`，**不进** `GATE REPORT` 的 stdout（stdout 不含任何耗时数字，这是 A/B 逐字比对可行的前提）。
