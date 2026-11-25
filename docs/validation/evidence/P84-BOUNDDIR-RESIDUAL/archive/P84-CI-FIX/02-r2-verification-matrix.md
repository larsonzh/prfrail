# P84 CI 首跑红处置 · R-2 验证矩阵（Windows 本机）

- 改动：仅 `internal/evidence/bounddir_windows.go`（+23/−2）。文件叶打开追加
  `createOptionSynchronousIONonAlert = 0x00000020`；两处注释更正。目录组件、句柄走链、
  拒绝条件、枚举（类 10/11）均未动。

## 工具链版本（实际输出）
- `go version go1.22.12 windows/amd64`（下载自 mirrors.aliyun.com/golang，置于仓库外
  `D:\LZProjects\.toolchains\go1.22.12`；GOCACHE=`D:\LZProjects\.toolchains\gocache122`）
- `go version go1.27.0 windows/amd64`（本机主工具链）

## 修前红 / 修后绿（同一工具链 go1.22.12、同一测试）
| 测试 | 修前（d597f1e 树） | 修后 |
|---|---|---|
| `TestBoundDirReadSurvivesAncestorReplacement` / `TestStoreReadBoundEvidenceSurvivesAncestorReplacement`（adapters） | RED：`read f.txt: The parameter is incorrect.` | GREEN |
| `TestAgentRunnerPostflightPassesOnRecordedEvidence`（adapters） | RED：uncertain，ErrorEvidence=`sha256:7e4012f9…` | GREEN |
| `TestRunnerPersistsPassedResultAndArtifacts` / `TestP83Fix7T4aGlobNamesComeFromTheBoundObject`（gates） | RED：gate failed | GREEN |
| `TestP84W2` / `W5b` / `W8` / `W9`（evidence，junction/读腿） | RED | GREEN |

- 修后红转绿实测：adapters ok、gates ok、evidence W2/W5b/W8/W9 全部 PASS（见运行记录）。
- ReadAll 证明：修前错误为同步 ReadFile 于异步文件对象返回 ERROR_INVALID_PARAMETER(87)
  （探针 syncOpt=false 实测 87）；修后同步与 OVERLAPPED 两种读均成功（探针 syncOpt=true）。

## 全量验证
| 工具链 | 命令 | 结果 |
|---|---|---|
| go1.22.12 | `go test ./... -count=1 -timeout 1800s` | EXIT=0，15/15 ok（第二次运行；见下「登记」） |
| go1.27.0 | `go test ./... -count=1 -timeout 1800s` | EXIT=0，15/15 ok |

- `gofmt -l` 无输出；`go build ./...`、`go vet ./...` EXIT=0；IDE 诊断（get_errors）0。
- `node tools/gates/selftest.js`：SELFTEST PASS 740/740。
- 门禁：scope=tree TOTAL_FAIL=0；scope=ci TOTAL_FAIL=0；scope=index 在**未暂存**时
  TOTAL_FAIL=1（`G5(f) index hygiene`，要求 index 与工作树一致）；**授权精确暂存后**
  复跑 scope=index TOTAL_FAIL=0（2026-11-21 实测，GATE_EXIT=0）。

## 登记（如实）
1. **工具链首次放置失误**：go1.22.12 一度解压于仓库 `tmp/`，导致 `go test ./...` 命中
   `tmp/go1.22.12/go/src/**` 触发 go1.22 装载器 panic（"LoadImport called with empty
   package path"）。已移至仓库外 `D:\LZProjects\.toolchains` 后全量复跑通过；与仓库代码
   无关、CI 不受影响。
2. **C7 单次失败（按 2026-11-21 裁定口径登记）**：go1.22.12 全量第一次运行中
   `TestP84T5h4C7HistoryDepthIsFailClosed` 因 `rename ... .lease-…: Access is denied`
   失败（lease 文件重命名面，非 bounddir 路径、非 R-2 相关面）。隔离复跑 3/3 PASS；
   全量第二次运行 15/15 ok。**记为单次失败、复验未复现；不定性为环境抖动**。
3. 报告双份行数勘误 **246/246** 仍随 ⑧b 回填登记（未动）。
