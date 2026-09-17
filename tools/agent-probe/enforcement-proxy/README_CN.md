# enforcement-proxy（B2 外部强制边界工具）

> English: [README.md](README.md)

本工具是 B2 证据阶段使用的**最小可审计正向代理**。它属于**工具**（`tools/agent-probe/`），不是生产运行时代码。
仅使用标准库：不做 TLS 拦截（HTTPS 以 `CONNECT` 隧道转发），无第三方依赖，无需安装或卸载。

## 它所属的边界架构

1. 受管候选 CLI 运行在**零网络能力 AppContainer** 中，操作系统在内核层拒绝其全部公网与 LAN 出网（含 DNS）。
2. 容器 SID 持有**回环豁免**，因此唯一可达的目标是本机回环服务。
3. 本代理监听**一个固定的回环端点**，是容器唯一的出口路径：
   `容器 -> 127.0.0.1:<端口>（本代理）-> 宿主侧出网（默认直连，或用 -upstream 链路上游代理）`。
4. **宿主侧出网这一段属于边界之外**：容器既不能直连公网或 LAN，也无法绕过本代理。因此 `networkControl`
   的准确表述是"容器无直连出口；所有出网都经代理并受白名单约束"，而**不是**"宿主本身没有网络"。

## 端口策略（固定，绝不静默更换）

- 默认监听端点：`127.0.0.1:39877`；可用 `-listen` 配置。
- 监听地址**必须是回环地址**（`127.0.0.1`、`::1`、`localhost`）。其它地址一律在启动时拒绝：把白名单代理
  （观察模式下乃至整个上游）暴露在可路由地址上，等于让它成为面向网络的开放中继。
- 启动时绑定**指定端点**；若端口被占用则**以非零码退出**并记录 `event=start decision=error`，
  绝不回退到其它端口。
- B2 期间端口冻结为 39877，所有证据都指向该端口。事后更改端口是允许的，但必须重跑可达性探针（P1）与
  smoke 探针（P4），并把新端口与结果写入证据包。B4 复用必须与 B2 证据使用同一端口；否则 B4 必须重新验证
  边界并在报告中写明差异。
- 系统规则（AppContainer 回环豁免）**不引用端口**，因此改端口不增加恢复工作量。

## 白名单语义

- `-allow` 接收逗号分隔列表；一条规则匹配该主机**及其所有子域**（点边界），大小写不敏感，忽略端口。
- 规则必须是**裸主机后缀**。带协议、路径、端口、通配符或非法标签的规则会在启动时被拒绝，而不是被静默
  接受为一条永远匹配不上的规则。
- 默认值：`api.githubcopilot.com,api.github.com,github.com,githubusercontent.com,copilot-proxy.githubusercontent.com`。
- 其它目标一律返回 `403`，并记录 `decision=deny reason=not-on-allowlist`。
- `-observe` 模式会转发未知主机**并记录**，仅用于探测发现。观察模式的运行结果**绝不能**作为强制证据；
  证据必须冻结白名单后以 enforce 模式重跑。

## 上游代理

- `-upstream http://host:port` 让代理在宿主侧这一段改为经上游 HTTP 代理出网。容器看不到上游；更改它不触碰
  任何系统规则。
- 默认直连出网。

## 审计日志（即证据）

JSONL，每行一个对象：

| 字段 | 含义 |
|---|---|
| `ts` | UTC 时间戳（RFC 3339，纳秒） |
| `event` | `start`、`stop`、`connect`、`http` |
| `phase` | `authorized`（放行前预记录）；`connect` 另有 `established` 与 `closed`；`http` 另有 `allow` |
| `decision` | `allow`、`deny` 或 `error` |
| `reason` | `allowlisted`、`observe`、`not-on-allowlist`、`malformed-request: …`、`dial-failed: …`、`roundtrip-failed: …` |
| `host`、`port` | 目标（明文 HTTP 事件记录 URL 端口，缺省按 80/443 补全） |
| `upstream` | 正在使用的上游代理（若有） |
| `client` | 调用方地址（容器） |
| `bytesToTarget`、`bytesToClient` | 隧道字节数，**仅**出现在 `connect closed` 事件上 |
| `durationMs` | `http` 事件的处理耗时 |

审计完整性：日志是 `networkControl` 的证据，因此写失败一律 **fail-closed**。放行路径会先写一条
`authorized` 预记录再转发任何流量；该写入失败时请求直接返回 `403`，隧道不会建立（已建立的隧道立即关闭），
且首次故障会输出到 stderr。

## 构建、运行、自检

```
go build -o enforcement-proxy.exe ./tools/agent-probe/enforcement-proxy
enforcement-proxy.exe -listen 127.0.0.1:39877 -log evidence\proxy.jsonl
enforcement-proxy.exe -listen 127.0.0.1:39877 -allow "api.githubcopilot.com,github.com" -upstream http://10.0.0.246:8080
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\enforcement-proxy\Test-EnforcementProxy.ps1
go test ./tools/agent-probe/enforcement-proxy/
```

自检完全在回环上运行（只用本机目标，不产生外部流量、不调用模型），覆盖：白名单转发请求、被拒请求、
观察模式记录、端口占用 fail-fast、日志可解析性。中间产物位于 `<repo>/tmp/enforcement-proxy-selftest/<时间戳>`
（仓库 `tmp/` 已被 gitignore），成功后自动清理；加 `-KeepArtifacts` 可保留以便排查。

## 本工具承载的流程要求

1. **付费探针前先做恢复排练。** `tools\agent-probe\appcontainer-b2\apply-b2-loopback.ps1` →
   `verify-b2-loopback.ps1` → `revert-b2-loopback.ps1` → `verify-b2-loopback.ps1` →
   `apply-b2-loopback.ps1`（由同目录的 `Invoke-B2LoopbackRehearsal.ps1` 编排），并把 revert 的差异（豁免表 +
   防火墙规则，必须为零差异）写入证据包。
2. **回环基线失效必须可执行。** 用 `tools\agent-probe\appcontainer-b2\verify-b2-loopback.ps1` 取监听快照并与
   B2 基线比较（其底层助手是同目录 `b2-loopback-lib.ps1` 中的 `New-B2Snapshot`），时间点为 **P1 之前**与
   **任何 B4 复用之前**。任何新增的回环监听都会报告给人工负责人，由其决定是否继续。
3. **出口路径必须在证据中声明。** 容器只能到达 `127.0.0.1:<端口>`；宿主侧出网由本代理完成（直连或经
   `-upstream`）。`networkControl` 的声明覆盖**容器**，不覆盖宿主。

## 提权使用后如何恢复环境（收尾流程）

提权运行最多只留下两处机器状态改动；本节即把它们恢复为 B2 之前状态的流程。脚本位于
`tools\agent-probe\appcontainer-b2/`（见其 `README.md`），临时产物写入被 gitignore 的 `tmp/b2/`。
所有命令都必须在**管理员** PowerShell 中、于仓库根目录执行（脚本在非提权下会直接拒绝运行）。

| 机器状态 | 由谁创建 | 由谁移除 |
|---|---|---|
| 容器 SID 的回环豁免（整表写入，保留既有条目） | `tools\agent-probe\appcontainer-b2\apply-b2-loopback.ps1` | `tools\agent-probe\appcontainer-b2\revert-b2-loopback.ps1` |
| AppContainer profile `prfrail-b2-battery` | `tools\agent-probe\appcontainer-b2\apply-b2-loopback.ps1`（仅首次运行） | 同一次 `revert-b2-loopback.ps1`（除非用 `-KeepProfile`）；保留场景用 `remove-b2-container-profile.ps1` |

流程：

1. `verify-b2-loopback.ps1` —— 只读：打印 SID、profile 存在性、当前豁免表以及与基线的严格差异。此时期望容器 SID 仍在表中。
2. `revert-b2-loopback.ps1` —— 移除容器 SID 豁免、**删除 AppContainer profile**（“恢复”= 回到已记录基线；
   `-KeepProfile` 可保留供 B4 复用），并重跑严格差异（豁免表 + 防火墙 + profile 存在性）。运行必须以
   **零差异** 结束；否则说明机器未回到基线，必须先排查再继续。
3. `verify-b2-loopback.ps1 -Expect absent` —— 只读确认豁免已移除。
4. 可选：`remove-b2-container-profile.ps1` —— 仅用于 `-KeepProfile` 保留后的补删。
5. 把第 1–3 步的输出归档进证据包并重算 `SHA256SUMS.txt`（`assemble-evidence.ps1` +
   `normalize-encoding.ps1`，同目录）。

早于 schemaVersion 2 的基线无法回答 profile 存在性问题；`revert`/`verify` 会以信息行
`profile-existence-unverifiable` 报告而不是当作差异，
`apply-b2-loopback.ps1 -RefreshBaseline`（或 `Invoke-B2LoopbackRehearsal.ps1 -RefreshBaseline`）可重录有效基线——
但仅在机器已处于恢复态（无豁免、无 profile）时才允许。

排练脚本 `Invoke-B2LoopbackRehearsal.ps1` 执行 apply → verify → revert → verify → apply，是“该收尾可重复”
的证明；请在管理员会话中运行它并归档转录。

收尾流程**按设计不覆盖**的内容：

- `tmp/b2/` 下每个 run 目录的一次性 ACL 授权（随目录本身删除而消失）；
- `subst` 盘符映射（按会话生效，由 harness 自己撤销）；
- 代理自身出网造成的宿主侧影响（属于边界之外，见“已知限制”）。

## 已知限制

- 不做 TLS 拦截：HTTPS 对代理不透明（白名单在 CONNECT 目标上强制）。
- 明文 HTTP 的协议升级流（WebSocket）不隧道化；仅转发标准请求。
- `-observe` 模式仅用于探测发现，绝不可作为强制证据引用。
- 本工具不管理 AppContainer 及其豁免；那套成对纪律是**仓库内的伴随工具**
  `tools/agent-probe/appcontainer-b2/`（`apply-b2-loopback.ps1` / `revert-b2-loopback.ps1` /
  `verify-b2-loopback.ps1`），其当时的运行快照已归档在
  `docs/validation/evidence/b2-2026-09-17/machine-ops/`。
