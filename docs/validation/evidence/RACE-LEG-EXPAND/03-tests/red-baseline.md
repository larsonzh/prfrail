# RACE-LEG-EXPAND · 红基线报告（03-tests）

- 切片：P49 `RACE-LEG-EXPAND`
- 角色：③ 测试工程师（prfrail-tester）
- 日期：2026-11-02
- 目标环境：`.github/workflows/ci.yml` 的 `Race` 腿（Ubuntu Linux 专属）

## 0. 一句话结论

当前 `Race` 步骤的 `run:` 行仅包含 `./internal/adapters/... ./internal/chain/...`，**未覆盖 `console` / `gates` / `guard` / `tickets` 四包**；断言脚本以退出码 **1** 输出 `RACE-COVERAGE: RED`。

## 1. 断言脚本设计

- 路径：`tmp/RACE-LEG-EXPAND/assert-race-leg-coverage.js`
- 形态：Node.js 单文件，**只读** `ci.yml`，不写任何文件。
- 编码：UTF-8 无 BOM + LF；注释与输出全英文。
- 为什么不是 `.go`：`tmp/` 下的 `.go` 会被 `go build ./...` 纳入编译，污染门禁。

解析流程：

1. 按行读取 `ci.yml`（CRLF 容忍，归一为 LF）。
2. 用正则 `^(\s*)-\s+name:\s*Race\s*$` 定位 `Race` 步骤并记录缩进。
3. 从下一行扫描，遇到缩进不大于该步骤缩进的 `- ...` 列表项即视为步骤结束（当前 span = 第 82–85 行）。
4. 在 span 内提取 `if:` 条件与 `run:` 命令；`run:` 兼容单行值与块标量（`|` / `>`）两种形态。
5. 对 `run` 命令执行下列六项断言。

| # | 断言 | 失败含义 |
|---|---|---|
| 1 | `run` 以 `CGO_ENABLED=1 ` 起头 | 前缀丢失，race 在 CI 上编译不过 |
| 2 | 含 `go test -race` | race 腿被替换 |
| 3 | 含 `-count=1` | 测试缓存导致假绿 |
| 4 | 步骤保留 `if: runner.os == 'Linux'` | Windows/macOS 腿误跑 race |
| 5 | 六个包模式各出现 **恰好 1 次** | 缺包或重复包 |
| 6 | 六个包出现位置严格递增 | 顺序契约漂移 |

退出码契约：`0` = GREEN（六包齐备）；`1` = RED（任一断言失败）；`2` = ERROR（文件或步骤缺失）。

## 2. 红基线运行结论（原始输出已存档）

运行命令：

```
cmd /c "node tmp\RACE-LEG-EXPAND\assert-race-leg-coverage.js > tmp\RACE-LEG-EXPAND\03-tdd-red-run.txt 2>&1"
```

结果：`EXITCODE=1`。关键行摘录：

```
Race step found at ci.yml line 82 (step span to line 85)
Race step if   : runner.os == 'Linux'
Race step run  : CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/...

[CHECK] CGO_ENABLED=1 prefix: PASS - "CGO_ENABLED=1 " present at command head
[CHECK] go test -race: PASS
[CHECK] -count=1: PASS
[CHECK] Linux-only guard preserved: PASS
[CHECK] package covered: ./internal/adapters/...: PASS - occurrences=1
[CHECK] package covered: ./internal/chain/...: PASS - occurrences=1
[CHECK] package covered: ./internal/console/...: FAIL - occurrences=0
[CHECK] package covered: ./internal/gates/...: FAIL - occurrences=0
[CHECK] package covered: ./internal/guard/...: FAIL - occurrences=0
[CHECK] package covered: ./internal/tickets/...: FAIL - occurrences=0
[CHECK] package order matches documented order: FAIL - order differs or a package is absent

uncovered package patterns : ./internal/console/... ./internal/gates/... ./internal/guard/... ./internal/tickets/...

RACE-COVERAGE: RED - 5 check(s) failed; uncovered by Race leg: console/gates/guard/tickets
```

解读：前 6 项（前缀 / race / count / Linux 守卫 / 两个旧包）**全绿**，说明红点精确落在「四包缺覆盖 + 顺序不成立」，而非脚本解析失误造成的伪红（falsifiable red）。

## 3. 独立复核（不依赖断言脚本）

用 `Select-String` 直接读 `ci.yml`：

```
=== INDEPENDENT Select-String VERIFY: Race step run line ===
LINE 84:         run: CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/...
=== VERIFY: new packages absent from ci.yml ===
matches_for_new_packages_in_ci.yml = 0
```

即：第 84 行逐字内容只有两包；`internal/(console|gates|guard|tickets)/...` 在整个 `ci.yml` 中匹配数为 **0**。两条独立路径（脚本解析 + 文本检索）结论一致。

## 4. 取证文件清单

| 文件 | 内容 | 结论 |
|---|---|---|
| `tmp/RACE-LEG-EXPAND/03-tdd-red-run.txt` | 红基线断言完整原始输出 | exit 1，`RACE-COVERAGE: RED` |
| `tmp/RACE-LEG-EXPAND/03-tests-norace-green.txt` | 四包常规（非 race）测试输出 | exit 0，四行 `ok <pkg>` |
| `tmp/RACE-LEG-EXPAND/03-race-local-infeasible.txt` | 本机 `CGO_ENABLED=1 go test -race` 尝试 | exit 1，`cgo: C compiler "gcc" not found` |
| `tmp/RACE-LEG-EXPAND/assert-race-leg-coverage.js` | 断言脚本本体 | 可复跑，支持红 → 绿翻转 |

## 5. 红 → 绿翻转判定规则

② 修改 `ci.yml` 后重跑：

```
cmd /c "node tmp\RACE-LEG-EXPAND\assert-race-leg-coverage.js"
```

翻转成立当且仅当以下五条同时满足：

1. 退出码由 `1` 变为 `0`；
2. 末行由 `RACE-COVERAGE: RED ...` 变为 `RACE-COVERAGE: GREEN - all six packages covered by the Race leg`；
3. 六条 `[CHECK] package covered: ...` 全部为 `PASS - occurrences=1`；
4. `[CHECK] package order matches documented order: PASS`；
5. 前缀 / race / count / Linux 守卫四项**仍为 PASS**（防止「扩包时顺手改坏前缀」这类回归）。

任一条不满足即判定翻转失败；不得用「CI 日志看起来绿了」替代机械判定。

## 6. VM 真实 `-race` 运行（2026-11-02 · 主控落盘补记，非 ③ 静态断言部分）

> 来源：用户同轮环境升级裁定——本地开启 Ubuntu 原生 VM（`10.0.0.199`，用户 `larson`）用于真实 `-race` 验证，取代「仅静态文本断言」的口径。以下小节由主控按裁定落盘，工具链与运行事实全部可复现。

| 项 | 值 |
|---|---|
| VM | LZUbuntuLTS（Linux 6.8.0-85-generic x86_64，Ubuntu 24.04） |
| Go | go1.22.12 linux/amd64（与 CI `GO_VERSION` 一致；官方 tarball 安装，sha256 `4fa4f869…34d43` 与 go.dev 发布清单核对一致） |
| gcc | gcc (Ubuntu 13.3.0-6ubuntu2~24.04.1) 13.3.0 |
| 仓库 | `git clone https://github.com/larsonzh/prfrail.git` + 检出 `81dd49f`（与主控工作区 HEAD 一致），工作区干净 |
| 命令 | `CGO_ENABLED=1 go test -race -count=1 ./internal/console/... ./internal/gates/... ./internal/guard/... ./internal/tickets/...` |
| 运行窗口 | 2026-11-02T04:16:48Z – 04:17:38Z（约 50 秒），rc=0 |

**结果：四包全绿，无数据竞争。**

- `ok github.com/larsonzh/prfrail/internal/console 20.410s`
- `ok github.com/larsonzh/prfrail/internal/gates 3.384s`
- `ok github.com/larsonzh/prfrail/internal/guard 3.043s`
- `ok github.com/larsonzh/prfrail/internal/tickets 7.420s`

原始日志已回拉并归档：`03-tests/logs/vm-race-four-pkgs-green.txt`（含环境头部与逐字 `ok` 行）；本机 tmp 副本同文件。**未发现历史数据竞争** ⇒ 不触发「另立专项切片」停点；红绿反转语义 = 修改前的 `ci.yml` 无法覆盖这四包（静态红基线已证），而 VM 真实 `-race` 完美通过（本小节实证），② 可实施扩包。
