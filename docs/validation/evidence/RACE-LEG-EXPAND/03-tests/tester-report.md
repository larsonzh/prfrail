# RACE-LEG-EXPAND · ③ 测试角色报告

- 切片：P49 `RACE-LEG-EXPAND`
- 角色：③ 测试工程师（prfrail-tester）
- 日期：2026-11-02
- 授权边界：不修改 `.github/workflows/ci.yml`（② 的活）、不修改任何 `.go` 源码、产出仅落 `tmp/RACE-LEG-EXPAND/` 与 `docs/validation/evidence/RACE-LEG-EXPAND/03-tests/`

## 1. 交付物

| 类别 | 路径 |
|---|---|
| TDD 断言脚本 | `tmp/RACE-LEG-EXPAND/assert-race-leg-coverage.js` |
| 红基线原始输出 | `tmp/RACE-LEG-EXPAND/03-tdd-red-run.txt` |
| 四包常规测试输出 | `tmp/RACE-LEG-EXPAND/03-tests-norace-green.txt` |
| 本机 race 不可行取证 | `tmp/RACE-LEG-EXPAND/03-race-local-infeasible.txt` |
| 红基线报告 | `docs/validation/evidence/RACE-LEG-EXPAND/03-tests/red-baseline.md` |
| 本报告 | `docs/validation/evidence/RACE-LEG-EXPAND/03-tests/tester-report.md` |

## 2. 覆盖矩阵（编号 / 测试名 / 类型 / 对应测试点）

| 编号 | 测试名 | 类型 | 对应测试点 | 手段 | 结果 |
|---|---|---|---|---|---|
| T1 | `RACE-LEG-EXPAND TDD ASSERTION`（红基线） | 反例（当前缺陷） | 六包必须全部进入 Race 腿 | Node 断言 | **RED**，exit 1 |
| T2 | Race `run:` 行逐字复核 | 正例（独立取证） | 四新包是否真的缺席 | `Select-String` | 仅两包，匹配数 0 |
| T3 | `internal/console` 常规测试 | 正例 | 新包存在可执行测试 | `go test -count=1` | `ok 14.041s` |
| T4 | `internal/gates` 常规测试 | 正例 | 同上 | `go test -count=1` | `ok 3.560s` |
| T5 | `internal/guard` 常规测试 | 正例 | 同上 | `go test -count=1` | `ok 8.850s` |
| T6 | `internal/tickets` 常规测试 | 正例 | 同上 | `go test -count=1` | `ok 8.228s` |
| T7 | 本机 `-race` 可行性 | 环境边界 | race 需 cgo/gcc | `CGO_ENABLED=1 go test -race ./internal/gates/...` | **确定性失败**（build failed） |
| T8 | 六包 `-race` 实跑 | 并发 | 数据竞争检测 | Ubuntu CI 腿 | **未执行（本机不可行）** |

T1 的证伪性：脚本对 6 个包逐条断言 + 4 条形状断言（前缀 / race / count / Linux 守卫）。删掉「包覆盖」守卫 ⇒ 必然 RED；删掉「前缀」守卫 ⇒ 无法拦截前缀回归。前 6 项在基线为 PASS，说明红灯**只**由四包缺席与顺序不成立触发。

四包测试家底（用于支撑 T3–T6 非空论证）：

| 包 | `_test.go` 文件数 | 顶层 `Test/Fuzz/Benchmark` 函数数 |
|---|---|---|
| `internal/console` | 13 | 104 |
| `internal/gates` | 7 | 48 |
| `internal/guard` | 8 | 20 |
| `internal/tickets` | 5 | 35 |

## 3. 本机不可 race 的降级口径

- 环境事实：Windows、Go 1.27.0 windows/amd64、`go env CGO_ENABLED` = `0`、`Get-Command gcc` 为空。
- race detector 依赖 cgo 工具链；实测 `set CGO_ENABLED=1&& go test -race -count=1 ./internal/gates/...` 直接失败：

```
# runtime/cgo
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%
FAIL    github.com/larsonzh/prfrail/internal/gates [build failed]
FAIL
```

- 因此本机只能提供**弱佐证**（四包测试存在且常规模式全绿），**不能**证明六包无数据竞争。
- race 真实性**唯一承担者 = Ubuntu CI 腿**（`if: runner.os == 'Linux'` + `CGO_ENABLED=1`）。
- 不安装 gcc：超出本片授权范围，且会改变本机环境（不留可复现差异）。

## 4. 给 ② 实施者的注意事项

1. **只改 `Race` 步骤的 `run:` 行**。不要动 `- name: Race`、不要动 `if: runner.os == 'Linux'`、不要动相邻的 `Test` / `Contract fixtures` 步骤。
2. **保持 `CGO_ENABLED=1 ` 前缀与 `-count=1`**。去掉前缀会在 Linux 上直接 build fail；去掉 `-count=1` 可能因缓存产生假绿。
3. **六包顺序固定**：`adapters → chain → console → gates → guard → tickets`（断言按位置递增校验）。
4. **每个包模式只能出现一次**，统一定义为 `./internal/<pkg>/...`，不要混入 `./internal/gates` 之类第二种形态。
5. **改完立即复跑** `node tmp\RACE-LEG-EXPAND\assert-race-leg-coverage.js`，必须 exit 0 且末行 GREEN；新输出另存（**不要覆盖** `03-tdd-red-run.txt`，红基线是翻转证据）。

## 5. 本片未执行的测试 / 验证项（显式登记）

| 未执行项 | 原因 | 承担方 |
|---|---|---|
| 六包在 `-race` 下的实际竞争检测 | 本机无 gcc / CGO 不可用 | Ubuntu CI 腿 |
| ② 修改后的绿基线（翻转执行） | ② 尚未实施 | ③ 复跑（规则见 `red-baseline.md` §5） |
| GitHub Actions 端真实运行结果 | 需 push 触发，本片未获提交授权 | ⑧ / 发布流程 |
| 扩包后 Race 腿耗时与超时风险 | 无 CI 实测数据 | ② 后观察 CI 耗时 |

## 6. 结论

红基线成立且可证伪（前 6 项 PASS / 后 5 项 FAIL 的差异化结果排除了脚本解析失败）。四包常规模式全绿，Race 真实性按降级口径交由 Ubuntu CI 腿承担。

TEST: DONE
