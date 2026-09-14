# T027 A0 replay-root composition 前置验证报告

[English](t027-replay-root-composition_EN.md)

日期：2026-09-14。结论：`PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`。

## 验证范围

本报告只覆盖 T027 离线切片 A0 replay-root composition，不宣称 AgentRunner 生产实现。A0 已交付：

- 唯一生产构造入口 `NewAgentRunnerReplayStoreForRun(runRoot, runID)`：replay root 固定由 durable run root 推导为 `<runRoot>/agent-runner-replay`，不再接受任意 caller root；接受任意绝对路径的构造器改为包内非导出 `newAgentRunnerReplayStoreAt`，仅保留测试与内部组成用途。
- durable run root 标记：run root 必须已存在且含常规文件 `events/state-events.jsonl`；跨重启定位只依据 `(runRoot, runID)`，runID 必须为合法 evidence ID，禁止扫描或猜测。
- runID 绑定：加载 R/C 时 `runId` 不符按 corruption fail-closed；写入 runId 不符的 R/C 按 root conflict 拒绝（request 与 completion 两侧）。
- 重叠与别名拒绝：`RejectWriterRoots` 与写者根双向不相交（含 8.3 短名与文件同一性别名）；`RejectProtectedRoots` 相等或位于 replay root 内即拒，包含 replay root 仅在包含链覆盖 run root 时允许；包含或同一性无法判定一律 fail-closed。
- symlink/reparse：构造期与每次读写入口（`State`、`RecordRequest`、`RecordCompletion`）加锁后复检 root、`requests/`、`completions/`、ownership 路径的组件级 symlink（Unix）或 reparse point（Windows，含 junction 与挂载点）；失败 fail-closed 且不得落盘根外。
- ownership：`ownership.json`（kind、schemaVersion、runId、runRootPath、runRootHash、createdAt；域 `proofrail:agent-runner-replay-ownership:1\n`）按 no-replace 单调发布且不依赖发布耐久 `proven`；缺失且已有 R/C 按 corruption；缺失且根为空时重建；runId 或 hash 不符、目录移动或复制按 root conflict；标记发布 link 后父目录同步失败时 fail-closed 并保留已可见标记；构造失败（含 bootstrap 与低层构造）不遗留 ownership 元数据。
- 边界：未引入真实 dispatch；未开放任意 caller root；未改变发布耐久 `unproven` 语义；T027 仍 `BLOCKED / NOT IMPLEMENTED`。

## 执行结果

### Windows 本机执行（原生证据）

| 命令 | 结果 |
|---|---|
| `gofmt -l internal/adapters` | 干净（无输出） |
| `go build ./...` | 通过 |
| `go vet ./...` | 通过 |
| `go test -count=1 ./...` | 全部通过 |
| `go test -count=1 -run 'TestAgentRunnerReplay' ./internal/adapters` | 通过 |
| `node tools/contracts/contracts.test.js` | 4/4 通过（未改 schema/fixtures，作为回归确认） |

### Windows 平台反例测试（实际执行，非跳过）

| 测试 | 结果 |
|---|---|
| junction 替换 replay root → unsafe path 拒绝 | PASS |
| 大小写别名拼写重开同一根 | PASS |
| 8.3 短名别名重叠拒绝与短名重开 | PASS |
| 保护根短名别名 → conflict | PASS |
| `\\?\` 扩展前缀重开同一根 | PASS |
| 盘符相对 `C:relative` → conflict | PASS |
| 并发不同拼写首开全部成功 | PASS |
| 运行期 `requests` 被替换为 junction → `RecordRequest` unsafe-path 且外部目录保持为空 | PASS |

### 审查结论

- V4 Pro 实现预审：PASS（有条件）→ 已修复 Medium（caller root 链级 reparse 检查）与全部 Low（C 侧写前分类、契约措辞、marker 常规文件校验、createdAt 校验、marker IO 错误挂 sentinel、读取路径零值 root 门禁），并补 12 项反例测试。
- Codex 独立终审：需修复 → Medium（运行期 TOCTOU 复检）已实现并通过反例测试；Low ①（ownership 发布时序）已收紧为"低层构造先于 ownership 发布"；Low ②（反例缺口）已补测试。
- Codex 收尾复审：PASS。

## 已知边界

- Unix 侧测试（symlink 拒绝、sync 顺序、sync 失败 fail-closed、构造失败不遗留 ownership、收敛耗尽、运行期换链）尚未在 Linux 原生执行；待授权推送后由 GitHub Actions Ubuntu 回归确认，在此之前不宣称 Unix 原生结论。
- 用户态组件检查与文件操作之间仍存在无法完全消除的 TOCTOU 残余窗口，已写入 CONTRACTS，属显式接受的边界。
- 全链 reparse/symlink 拒绝意味着 runDir 位于 OneDrive 重定向目录、junction 目录树或挂载盘文件夹等布局会被 fail-closed 拒绝；这是既定安全取舍，如需放宽须单独 ADR。
- 8.3 短名与 `\\?\` 用例在能力不可用的环境会显式 skip；本机实际执行，未跳过。

## 明确未执行事项

- 未引入真实 dispatch，未接入 `AgentRunnerPort` 或 Engine；
- 未调用真实模型 CLI，未访问网络，未读取 SecretStore；
- 未执行 commit、push 或 publish；
- T027 仍为 `BLOCKED / NOT IMPLEMENTED`，AT-23 未通过。
