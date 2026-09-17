# appcontainer-b2（B2 机器纪律脚本）

本目录是 B2“外部强制”证据所依赖的机器纪律脚本的版本化归属，后续 B4/B5 复用同一容器或同一端口时也直接
使用这里的脚本。它们是**工具**，不是产品运行时代码，并且刻意成对存在：任何写入机器状态的动作，都必须
有配套的回退脚本加基线比对。

临时产物（运行目录、转录、快照、canary 文件）一律写入仓库内被 git 忽略的 `tmp/b2/`，**不得**与这些脚本
放在一起。本目录的“运行当时快照”已带哈希归档到 `docs/validation/evidence/b2-2026-09-17/machine-ops/`。

清理策略：脚本本身已版本化，删除 `tmp/` 除丢掉中间运行目录外没有代价。证据包冻结后可以整体删除
`tmp/b2/`，唯一影响是 `assemble-evidence.ps1` 会把历史运行目录报为 `missingArtifact`（属预期；已冻结的
证据包及其 `SHA256SUMS.txt` 不受影响）。

一个已实测的坑：`DeriveAppContainerSidFromAppContainerName` **只按名字推导** SID，即使该名字还没有对应
配置文件也会成功（2026-09-17 实测），因此“配置文件是否存在”必须改读每用户注册表键
`HKCU\Software\Classes\Local Settings\...\AppContainer\Storage\<name>`。

编码纪律（这条真的在提权运行里咬过人）：脚本写出的所有产物都是 UTF-8 **无** BOM，而 PowerShell 5.1 的
`Get-Content` 会把无 BOM 的 UTF-8 当 ANSI 解码，从而破坏非 ASCII 路径与显示名（中文用户配置文件、
本地化的防火墙规则名），并让 `ConvertFrom-Json` 在合法 JSON 上失败。读取一律走 `b2-loopback-lib.ps1`
中的 `Read-B2TextFile` / `Read-B2JsonFile` / `Read-B2Lines`，**不要**用 `Get-Content`。
`Test-B2Tooling.ps1`（密闭：不提权、不碰机器状态、不联网）守护这条及其他不变量——改动本目录脚本后请先跑它：

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\Test-B2Tooling.ps1
```

## 脚本

| 脚本 | 用途 | 写入机器状态 | 需要提权 |
|---|---|---|---|
| `ac-lib.ps1` | 容器启动器（用 AppContainer 令牌调用 `CreateProcess`，自定义环境块）、从 Windows 凭据管理器读取凭据、SID 辅助函数。只推导：从不创建配置文件。 | 否 | 否 |
| `b2-loopback-lib.ps1` | SID 解析（`New-B2ContainerSidIfMissing` 会创建 / `Get-B2ContainerSidByName` 仅按名字推导 / `Get-B2ContainerSidExisting` 推导并要求配置文件存在）、豁免表解析、防火墙快照、回环监听枚举、基线快照（schemaVersion 2）与严格 diff。 | 否（辅助库） | 否 |
| `apply-b2-loopback.ps1` | **唯一被允许创建机器状态的脚本**：先写基线快照，再给容器 SID 授予回环豁免（整表写入并保留既有条目）；幂等；失败即关闭。`-RefreshBaseline` 会丢弃早于 schemaVersion 2 的基线，但仅在已确认回退状态（无豁免、无配置文件）之后才执行。 | 是 | 是 |
| `revert-b2-loopback.ps1` | 移除容器 SID 豁免**并删除 AppContainer 配置文件**（回退=回到记录的基线），随后重跑严格 diff，必须零差异；`-KeepProfile` 保留配置文件供 B4 复用，并把“配置文件存在”那行报为有意保留。 | 是（移除） | 是 |
| `verify-b2-loopback.ps1` | 只读：打印 SID、配置文件是否存在、豁免表以及严格 diff（`-Expect present|absent`）。旧基线下只产出信息行 `profile-existence-unverifiable`，永不判负。 | 否 | 否 |
| `remove-b2-container-profile.ps1` | 配合 `-KeepProfile` 的可选助手：删除配置文件；其豁免仍在时拒绝执行；已删除时幂等。 | 是（移除） | 否 |
| `Invoke-B2LoopbackRehearsal.ps1` | 编排 apply → verify → revert → verify（零差异）→ apply → verify；识别旧基线并自动把 `-RefreshBaseline` 转给第 1 步（apply 仍要求先确认回退状态），某步失败时打印该步输出尾部。 | 是 | 是 |
| `Test-B2Tooling.ps1` | 工具不变量的密闭自检（编码安全读取器、schema 现行性判定、diff 身份闸门、旧基线信息行、配置文件探测）。 | 否 | 否 |
| `boundary-battery.ps1` | 免费电池：代理连通性、LAN/公网/DNS 拒绝、允许清单放行/拒绝、容器内与外写入、canary 完整性。 | 否 | 否 |
| `container-candidate.ps1` | 在容器内运行固定版本候选体（`exec-check` 免费；`smoke`、`tool-effect`、`network-bypass` 计费，需 `-Run`）。 | 否 | 否 |
| `container-nodecall.ps1` | 容器内最小 Node 客户端（`/models` 免费，补全需 `-Run`）；用 `gh` 交换 Copilot 会话令牌。 | 否 | 否 |
| `container-modelcall.ps1` | 同一探测的 .NET 变体（保留用于受限请求头发现与交叉校验）。 | 否 | 否 |
| `path-probe.ps1` | 复现路径拼写证据（`subst` 拼写 vs 主机拼写）。 | 否 | 否 |
| `node-client.js` | 供 `container-nodecall.ps1` 使用的无依赖 Node 客户端。 | 否 | 否 |
| `assemble-evidence.ps1` | 从临时树重建经脱敏、带哈希的证据包。 | 否 | 否 |
| `normalize-encoding.ps1` | 按仓库编码约定规范化文本产物，并重建证据包 `MANIFEST.md` 与 `SHA256SUMS.txt`。 | 否 | 否 |

### `normalize-encoding.ps1` 的两个使用要点

- 默认**跳过运行产物目录**（`candidate-runs`、`nodecall-runs`、`modelcall-runs`、`boundary-battery`、
  `path-probe`、`selftest`，以及通配的 `exec-probe*`、`exec-isolation`、`ac-workspace*`）。这些目录里是
  候选体解包出的外部文件或会被单独规范化的临时输出；不跳过的话，一条 `-Root tmp\b2` 会重写上万个
  载荷文件，看起来就像卡死。需要纳入时可自行传 `-ExcludeDirNames`（支持通配符）。
- 运行时会打印 `normalizing …`、`scanned files / content candidates / skipped`、每 100 个文件的进度、
  `normalization done: candidates/rewritten/unchanged`，最后是 `bom/crlf audit: checked=… violations=…`
  并**只列出违规文件**。因此“命令在跑、只是慢”与“命令真的卡住”可以一眼区分。

## 典型顺序

```powershell
# 提权：先彩排，再 apply（若记录的基线早于 schemaVersion 2，加 -RefreshBaseline）
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\Invoke-B2LoopbackRehearsal.ps1 -RefreshBaseline
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\apply-b2-loopback.ps1

# 不提权：免费证据
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\boundary-battery.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\container-candidate.ps1 -Scenario exec-check
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\path-probe.ps1

# 提权收尾：verify -> revert -> verify absent -> 可选删除配置文件
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\verify-b2-loopback.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\revert-b2-loopback.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\verify-b2-loopback.ps1 -Expect absent
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\remove-b2-container-profile.ps1
```

`revert` 本身就会删除配置文件，因此最后一条命令只在 revert 曾带 `-KeepProfile`（为 B4 复用而保留）且现在
希望把它清掉时才需要。

关于强制代理本体、以及这些脚本所实现的边界，见 `../../agent-probe/enforcement-proxy/README_CN.md`
（英文版 `README.md`）。

> English version: [README.md](README.md)
