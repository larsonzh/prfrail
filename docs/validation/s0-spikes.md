# S0 T004 技术探针报告

[English](s0-spikes_EN.md)

日期：2026-09-07。结论：`T004 COMPLETE`，但整体 S0 仍为 `NOT_READY`。本报告只证明原语可行性和失败边界，不是 S1 生产实现、发布批准或真实模型能力证明。

## 1. 平台与方法

| 平台 | 实测环境 |
|---|---|
| Windows | Windows 11 Home 10.0.26100 amd64；D: 为 NTFS；Go 1.27.0 windows/amd64 |
| Linux | Ubuntu 24.04.4 LTS；Linux 6.8.0-85-generic x86_64；根文件系统 ext4；固定 Go 1.27.1 linux/amd64 |
| SessionBridge | 姊妹仓库提交 `d2298f3`；Python 3 合同测试与 Node history engine；未调用模型 |

探针均使用真实本地文件和真实子进程。Go 文件/队列/崩溃探针在两平台运行；Linux 代码经 SSH 传入 `/tmp`，运行后删除。Windows 原型位于仓库忽略的 `tmp/`，报告完成后删除。没有读取 whois 资产、调用收费服务、修改 SessionBridge 或提交 Git。

## 2. 十二个切片结果

| # | 切片 | 结果与约束 |
|---|---|---|
| 1 | Go 单二进制 | Windows 与 Ubuntu 均以 `CGO_ENABLED=0` 构建并运行 `prfrail version`，输出 `prfrail 0.1.0-dev`。Go 1.22 模块基线不变。 |
| 2 | TUI 依赖 | Bubble Tea v2.0.9 要求 Go 1.25，v1.3.10 要求 Go 1.24，均不兼容基线；v1.3.4 要求 Go 1.18、MIT，可作为 S1 候选，但尚未成为根模块依赖。 |
| 3 | TUI 行为 | v1.3.4 隔离原型通过 `q` 键退出、12 列截断和 `NO_COLOR` 无 ANSI 测试；CGO=0 的 Linux 二进制在 Ubuntu 实机产生相同输出。交互终端完整可访问性仍属 T016。 |
| 4 | Windows 文件原语 | NTFS 上文件 `Sync` 与无 reader 时的目标覆盖成功；目录句柄 `Sync` 返回 access denied。目标被普通 reader 打开时 `os.Rename` 返回 access denied，不能假定 POSIX 式开放句柄替换。 |
| 5 | 发布崩溃点 | 两平台结果一致：临时文件刷盘后、rename 前崩溃时 reader 见旧状态且无 receipt；rename 后、receipt 前崩溃时 reader 已见新状态但仍无 receipt；只有 receipt 存在才可判 completed。后两种非 completed 状态必须 fail closed。 |
| 6 | Windows 进程树 | 强杀直接 child 后 grandchild 仍存活；对已知 child 使用 `taskkill /T /F` 可停止 child 与 grandchild。S1 应使用可审计的 Job Object/等价受控树原语；PID 不明或归属无法证明时拒绝强杀。 |
| 7 | Windows 路径 | NTFS 大小写不敏感；NFC/NFD 名称可区分；硬链接可创建且 `os.SameFile` 可识别；普通用户 symlink 因缺少权限失败。Go 可创建超过 260 字符路径和名为 `CON` 的路径，所以保留名/别名必须由 canonical checker 显式拒绝。 |
| 8 | Linux 文件与路径 | ext4 上替换、文件 `Sync`、目录 `Sync`、symlink、hardlink 均成功；大小写敏感，NFC/NFD 名称可区分。开放 reader 继续读旧 inode，而路径读取新内容。 |
| 9 | Linux 进程组 | 杀 session leader 会留下 grandchild；对独立 session 的负 PGID 发信号可停止整组。未知归属同样 fail closed。 |
| 10 | 文件队列 | 32 worker 竞争时，Windows 的“把同一 ready 文件 rename 到不同 claim 名”两次出现 5/3 个成功返回，不能作为选主；ext4 为 1。固定 claim 文件的 `O_CREATE|O_EXCL` 在两平台均恰好 1 个赢家；generation 7 对当前 8 被拒，重复 requestId 只应用一次，无临时残留。 |
| 11 | SessionBridge silent | `python tests/test_contract.py`：28/28；`node tests/test_history_engine.js`：12/12；扩展语法和编码门禁通过。覆盖 busy、timeout、错误 requestId/陈旧结果防护、同 ID 重试、历史 reset/截断与去重。 |
| 12 | 双通道故障矩阵 | ProofRail 文件队列必须自管持久 requestId、generation、claim/result 与 dispatch/takeover receipt；SessionBridge 只作 silent transport。宿主重启、cache 丢失或 history 截断后不得由内存缓存宣称 exactly-once。真实 VS Code LM API 未调用，能力记为 `unavailable/not exercised`，不得回退 visible/GUI。 |

## 3. 可重复命令与判定

```powershell
$env:CGO_ENABLED='0'; go build -o tmp\prfrail-t004-windows.exe .\cmd\prfrail
go run .\tmp\t004_fs_probe.go
go run .\tmp\t004_crash_probe.go
go run .\tmp\t004_queue_probe.go
python .\tmp\t004_process_windows.py
go -C .\tmp\t004-tui test ./...
```

```sh
CGO_ENABLED=0 go build ./cmd/prfrail
go run /tmp/proofrail-t004-fs-probe.go
go run /tmp/proofrail-t004-crash.go
go run /tmp/proofrail-t004-queue-probe.go
/tmp/proofrail-t004-process.sh
NO_COLOR=1 /tmp/proofrail-t004-tui --render 12
```

持久回归命令：

```powershell
go build ./...
go vet ./...
go test ./...
npm --prefix tools/contracts test
python ..\sessbridge\tests\test_contract.py
node ..\sessbridge\tests\test_history_engine.js
node --check ..\sessbridge\extension\extension.js
python ..\sessbridge\tools\enforce_encoding.py
```

一次性探针在本任务结束时删除；上表给出其输入、竞争规模、崩溃点和判定条件。S1 必须把相同反例转成受版本控制的平台测试，不得依赖本报告替代测试。

## 4. ADR 约束与剩余风险

- ADR-003：单文件 rename 不等于事务或持久完成；采用“临时写入和文件刷盘 -> 平台替换 -> 平台可用的目录/卷持久化 -> completed receipt”顺序。Windows 必须处理打开 reader 和目录 `Sync` 不可用，不能把 Linux 结论外推；断电持久性仍需 T008/T017 原生故障注入。
- ADR-004：队列选主使用固定 claim 的排他创建，不使用 rename-only；所有 result 必须同时匹配 requestId 与 generation，持久去重先于副作用。SessionBridge cache/history 只提供 transport 上下文，不是核心 exactly-once 证据。
- TUI：S1 若采用 Bubble Tea，应从 v1.3.4 候选重新执行依赖摘要、许可证和漏洞审计。此次为解决网络限制曾在隔离目录以 `GOSUMDB=off` 和镜像取包；该设置不得进入生产构建。
- unavailable：未验证真实 VS Code LM API、Windows 断电/卷 flush、网络文件系统、管理员 symlink、Job Object 生产封装或完整交互式 TUI。相应能力必须标为 unknown/unavailable 并 fail closed，不能降级为成功。