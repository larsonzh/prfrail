# ③ 测试回执（test-r1）

- 切片：F2-FIX-STUB-HANG；角色：prfrail-tester（DeepSeek V4.1 Flash）；落盘 2026-11-09；仓库 HEAD `de37c58`。
- 纪律：未 commit / 未 push；**未修改 harness 与任何生产/工具脚本**（`tools/b4-e2e` 未动；stub 修复由 ② 落地）；
  临时件落在 `tmp/b4`（⑧c 清除）；证据内无密钥。

## 1. 红基线（修复前）

| 步 | 命令 | 退出码 | 证据 |
|---|---|---|---|
| 预播种 | `tmp/b4` 预建目录；`b4-chain.json`/`b4-chain-hook.json` 复制自 `docs/validation/evidence/directive-history/b4/archive/input/`，摘要复核 `BB8C100D…` / `C7AD22C5…`（t027-at23-e2e.md §2 口径） | 0 | `red/cli/` 各腿日志 |
| layer3 F2 红重跑 | `powershell -ExecutionPolicy Bypass -File tools\b4-e2e\run-faces.ps1 -OnlyFace f2` | **1** | `red/pack-verdict.json`、`red/driver.log`、`red/f2-pos|neg|neg-r2/` |

- 红判据（机械）：`FAILED CHECKS: f2-pos/process-tree-taken-away`；`red/f2-pos/verdict-run.json` 中该检查
  `pass=false`，detail=`reported pids gone: false (AgentRunner evidence invalid: missing required pid report(s) stub-child.pid …)`。
- 环境坑（新登记）：`run-faces.ps1` 在 `tmp/b4` 不存在时首写 `driver.log` 即失败（`DirectoryNotFoundException`），
  需预建 `tmp/b4`；另需先播种两个链配置文件（脚本不自带播种）。

## 2. TDD 回归测试（③ 写入，② 未动测试文件）

- 新增 `TestHangWritesEmptyChildPIDReport`（`tools/agent-stub/main_test.go` 末尾，复用既有骨架与常量）：
  后台启动 `-mode hang` → 轮询 `stub-child.pid`（deadline 5s）→ 断言 0 字节常规文件 + `stub.pid` 十进制 pid
  → `Kill()`+`Wait()` 清理（幂等 stop 闭包，无进程泄漏）。
- 红验证（修复前）：`go test ./tools/agent-stub -run TestHangWritesEmptyChildPIDReport -v` exit=1，
  唯一红点 = deadline `t.Fatalf`（`red-regression.txt` 原文归档）。
- 绿验证（修复后）：同命令 PASS（3.0s）；`go test ./tools/agent-stub -v` 5/5 全 PASS。

## 3. 绿重跑（修复后）

| 步 | 命令 | 退出码 | 证据 |
|---|---|---|---|
| layer3 F2 绿重跑 | 同红重跑命令（脚本重编译含修复 stub） | **0** | `green/pack-verdict.json`、`green/driver.log`、`green/f2-pos|neg|neg-r2/` |

- 绿判据：`all recorded checks passed`；`green/pack-verdict.json` `failedChecks=[]`；
  `green/f2-pos/verdict-run.json` 中 `process-tree-taken-away` `pass=true`，detail=`reported pids gone: true (<nil>)`；
  其余 5 条检查与 B4 基线同值（`error-evidence-present` 恢复为 `error evidence digests: 5`）。

## 4. GLM-1 三腿对照（零判断枚举，全量原文见本轮派单回执）

- pass 漂移 **0 处**；字段漂移 **1 处**：`f2-pos artifacts["stub-child.pid"]` 由基线「键不存在」变为本轮
  「存在（`sha256:e3b0c442…b855` = 空文件哈希）」——即本片修复的预期产物，非异常漂移。
- `failedChecks`：本轮 `[]` = B4 归档 `[]`。
- 备注（既有档案口径，非本片异常）：B4 归档 `SHA256SUMS.txt` 为剪裁前全量清单，git 追踪的 verdict 文件为
  剪裁后子集（8 个），与归档 MANIFEST「中间转储」分类一致；`f2-neg`/`f2-neg-r2` 单腿 verdict 属被剪裁的
  中间转储，故对照用 pack `cases` 作基线。
