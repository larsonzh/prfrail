# B3a 前置分析：断电注入装置与双向对照标定协议（仅设计，不含实现代码）

- 状态：工作稿（自 `babc041` 起已入库于 `docs/t027/`；切片开始时的输入物；`commit`/`push` 仍须同一轮显式授权）。
- 权威来源：`docs/t027/REMAINING_SLICES.md` 的 A7/B3a/B3b 块、`docs/validation/B2-EXTERNAL-ENFORCEMENT.md`、
  `docs/CODING_CONVENTIONS.md`、`D:\VirtualBox VMs\Win11B3\guest-prep\` 与 `docs/t027/` 下的冻结环境事实。
- 本文件自身编码：UTF-8 **with BOM** + **LF**（硬规则）。
- 角色边界：本文只定装置与标定协议。B3b 的候选矩阵、C2/C3 结论、first-dispatch 解锁一律不在本片。

## 0. 范围与不可做清单（B3a 边界，逐字执行）

**必须交付**（切片清单六项）：① 环境与前提声明；② 每轮崩溃点 journal（宿主侧、先记后写、确认后才注入）；
③ 断电执行与恢复预案（执行人、恢复介质、脏卷处置、轮次预算、重启窗口、50 分钟边界）；④ 重启盘点；
⑤ 双向对照标定（负对照 N 轮内至少一次丢失、正对照零丢失，否则装置不可用且 **B3b 不得开始**）；⑥ 装置脚本、复现命令、标定证据包。

**必须不做**：

1. 不对 C2（`MoveFileEx(MOVEFILE_WRITE_THROUGH)`）与 C3（目录句柄 `FlushFileBuffers`）下任何结论——包括标定过程中
   "顺带"观察到的现象，一律只登记、不裁决（裁决权在 B3b）。
2. 不用软件模拟落盘/软件注入替代设备级注入；唯一注入 = `VBoxManage controlvm Win11B3 poweroff`。
3. 不从单轮/单样本外推；每个标定结论都建立在 N 轮有效轮次上。
4. 不触碰宿主自身磁盘的破坏性操作、不触碰用户的远程 Ubuntu 机器。
5. 不把标定对照结果当作候选结论；不解除 first-dispatch 的 `unproven`；不调用
   `10-fix-timebomb.ps1`（许可决定，operator 已于 2026-09-19 否决）。该脚本随评估版环境一并退役：
   2026-09-20 已按 operator 决定改用零售镜像重建（见 §1.6），新环境不需要任何许可规避手段。
6. 装置若日后升级为物理断电：必须回到**真正的另一物理设备**上的 journal，并**重做全部标定**（本片结果作废）。

## 1. 前提声明（完整填写，结论必须与本前提绑定引用）

### 1.1 宿主与 hypervisor

- 宿主：Windows 11（26100 家族），单块 NVMe，C: 与 D: 同盘。VM 目录 `D:\VirtualBox VMs\Win11B3`。
- VirtualBox **7.2.18 r175117**。装置会话开始前必须断言 `VBoxManage --version` 输出恰为此值，漂移即阻断会话。
- 注入虚拟机：`Win11B3`，machine UUID `02205fb0-1933-48ba-bbdb-a2c06b4eb9b2`（2026-09-20 重建；旧评估版
  `Win11Eval-B3` 已退役但保留登记作后备，不再作为 B3a 环境）。

### 1.2 客机（注入对象）

- 客机系统：**Windows 11 专业版 25H2，build 10.0.26200.8037，zh-CN**（镜像 Image #4，零售/消费者通道、
  未激活）；EFI 固件 + vTPM 2.0；2 vCPU / 4096 MB；chipset piix3；显示 = VBoxSVGA + 128 MB VRAM（GA WDDM 已绑定）；
  输入区域 = zh-CN（装有微软拼音）。
- **与 A7 前提的构建差异（必须写明）**：A7 宿主为 24H2 build 26100；本客机为 25H2 build 26200。
  同一 GA 主线代码代际、不同 build。因此 A7 的任何结论都**不自动迁移**到本客机；
  本片结论只覆盖"本客机 25H2/26200 + NTFS + 本存储栈"这一前提。
- 存储后端：SATA/AHCI 控制器（PortCount=4），**`useHostIOCache=false`**（宿主 I/O 缓存关闭）；
  Port 0 = `Win11B3-os.vdi`（64 GB 动态）；Port 1 = `Win11B3-data.vdi`（8 GB **固定** VDI）。
  当前因存在 `pristine` 快照，运行态落在 `Snapshots\*.vdi` 差分链上（每次"恢复快照"即丢弃差分）。
- 网络：**NIC = none（完全离线、无 IP）**为冻结前提。已实测确认：`showvminfo` 八张网卡全 none/disabled，
  且客机内不存在任何真实网卡设备节点（仅 Teredo/6to4 等 Not Present 伪接口）。`.vbox` 中残留的
  `<Adapter slot="0"><NAT/></Adapter>` 元素不呈现硬件、不构成连通性（§8 的 Q1 据此结案）。
- 无光驱；boot order = disk only。
- 客机文件系统（注入卷）：**`R:` = INJECT，NTFS，7.98 GB**（`data.vdi`），**4096 B 群集 / 512 B 扇区**；
  卷序列号在 R0 轮由 `fsutil fsinfo volumeinfo R:` 现捕（旧环境为 C6DD-1858，本环境待捕）。基线时卷不脏。

### 1.3 自动化与提权

- UAC 已关（`EnableLUA=0`）⇒ 宿主驱动的进程获得 **High-integrity 管理员令牌**（finalize.log 实测 `Mandatory Label\High Mandatory Level`）。
- Guest Additions 7.2.18 工作正常，宿主通过 `VBoxManage guestcontrol <vm> run|copyto|copyfrom` 驱动客机（无 TCP/IP）。
- 客机凭据：宿主侧、仓外的凭据文件（`--passwordfile`，首行为密码）；路径经装置参数传入，**永不入仓**。
- 客机工具集：`C:\prfrail-prep`（宿主副本在 guest-prep 目录）。

### 1.4 降噪（已施加，冻结）

`wuauserv`、`UsoSvc`、`WSearch`、`SysMain` 禁用；Defrag 与 Defender 计划扫描禁用；休眠 + 快速启动关闭；
系统还原关闭；`R:\` 纳入 Defender 排除；`RealTimeIsUniversal=1`；显示器/待机/磁盘超时 = 从不。

### 1.5 基线快照与注入机制

- 基线快照 **`pristine`**（按**名称**引用）：重建后为 `7e0287b6-e00e-41ff-975d-2f690a4e981d`
  （VM 关机态、`boot1=disk`、无光驱、`nic1=none`、`graphicscontroller=vboxsvga`）。
- 注入 = `VBoxManage controlvm "Win11B3" poweroff`：设备级硬断电（A7 允许的等价注入），客机无感知。
- 客机重启后约 **48 s** 可再次 `guestcontrol`（实测 30–48.3 s；预算取 60 s，上限 120 s）。
- 客机侧断电证据：System 日志 **Id 41（Kernel-Power）** + **Id 6008**；**6008 滞后真实断电最多 ~50 s**
  ⇒ 客机时间戳永远只作次级证据，**断电权威时间一律取宿主侧 journal**。

### 1.6 许可与"无自关机"前提（2026-09-20 重建后）

- 客机为**零售/消费者通道的 Windows 11 专业版，未激活**。`SoftwareLicensingProduct` 报 `status=5, grace=0`
  ——这是**未激活的常态**（任何未激活的 Windows 都处于 notification），**与评估版的 `TIMEBASED_EVAL` 有本质区别**：
  零售通道没有固定到期日去驱动强制关机，状态码相同而含义不同。
- **行为取证（决定性）**：重建后空转观察 **70 个采样点 × 60 s**（日志 `D:\VirtualBox VMs\Win11B3\idle-watch.log`）：
  末次采样 `01:16:52 VMState="running" resets=1`（resets 为运行前既有值），`verdict=IDLE-NO-SHUTDOWN`；
  客机 uptime 00:05:53 → 01:18:15 ≈ **72 分钟仍在运行**。对照：已退役的评估版镜像在 **61.1 分钟**（两次实测 61.1/61.1）
  由 `wlms.exe` 干净关机。故"评估版每小时自关机"这一缺陷**已被消除**，不再需要任何许可层面的规避手段。
- 装置约束保留但**不再承重**（降为便宜的安全网）：
  1. 每个测量窗口仍断言客机 uptime **< 50 分钟**——理由从"躲开许可计时器"转为"窗口本身足够新鲜、
     防止把长暂停混入测量"；
  2. **注入前断言**（uptime 合理 + 无意外关机事件），不满足则不注入、该轮作废；
  3. **重启后审计** `C:\prfrail-prep\11-shutdown-audit.ps1`，命中 `wlms` 即作废该轮——即使环境已干净，
     它使"环境无自关机"这一主张在每一轮都可机器核验，而不是靠假设。

### 1.7 已观测的可行性现象（dryrun，非标定）

14 字节明文写入（`Set-Content`，无显式 flush）+ 约 2 s 后硬断电：重启后**文件存在、记录大小 14、全部 14 字节为零**
（目录项已提交、数据未持久）；卷未变脏。这是 B3a 必须能稳定复现的负对照丢失形态，但**单样本不构成结论**。

### 1.8 结论覆盖面与装置专用手段声明

- 本片全部结论只覆盖：**"本前提下的客机侧磁盘语义"**——VirtualBox 7.2.18 + SATA/AHCI 无宿主缓存 + VDI 后端 +
  本宿主 NVMe + Windows 11 25H2/26200 客机 + NTFS。**不**外推到物理硬件、其它存储栈、或产品保证。
- 正对照的目录项提交使用**卷句柄 flush**（`\\.\R:` 的 `FlushFileBuffers`）：这是装置标定专用、依赖客机提权令牌的手段，
  **不是产品候选**（C4 因便携定位仍被政策排除，ADR-009）；因此正对照的成立**不构成任何候选结论**。
  之所以不用 C3（目录句柄 flush）做正对照：C3 本身是 `unresolved` 的待测物，用未定论机制当标定锚会循环论证。
- 标定结论只认证**装置的分辨能力**（能分辨"未持久"与"已持久"），不认证任何发布协议。

## 2. 装置架构

### 2.1 宿主侧 journal（每轮崩溃点日志）

**位置（宿主侧、仓外）**：`D:\VirtualBox VMs\Win11B3\b3a-runs\<YYYY-MM-DD>\journal.jsonl`
（`b3a-runs/` 为装置运行目录，不随仓库走；会话结束后只把副本冻结进证据包）。

**格式**：JSONL，UTF-8 **无 BOM** + **LF**，每行一条记录，键序固定，字段如下：

| 字段 | 含义 | 必填 |
|---|---|---|
| `schema` | `"b3a-journal"` | 是 |
| `v` | `1` | 是 |
| `round` | 轮号（R0=对照复现轮，R1–R5 负对照，R6–R10 正对照；作废轮也占用单调递增号） | 是 |
| `candidate` | `negative-control` / `positive-control` / `parity` | 是 |
| `stage` | 五阶段分类法中的阶段号（§5.1）；负对照=`temp-written`，正对照=`parent-synced`，R0=`parity` | 是 |
| `targetPath` | 客机目标路径，如 `R:\calib\r03\neg.bin` | 是 |
| `contentDigest` | 待写内容的 SHA-256（hex） | 是 |
| `contentSize` | 待写内容字节数 | 是 |
| `planCutDelayS` | 写入返回后、断电前的固定延时（**2.0**） | 是 |
| `op` | `plan`（强制，每轮恰一条） / `cut` / `reconcile`（强推扩展） | 是 |
| `hostTs` | 宿主本地时间 ISO-8601（带时区偏移）——**断电权威时间** | 是 |
| `prevHash` | 上一条记录的 `recordHash`（链式防篡改；首条为 `null`） | 是 |
| `recordHash` | `sha256(prevHash || 本记录去掉 recordHash 字段后的规范 UTF-8 字节)` | 是 |

- `plan` = 本片要求的"每轮一条记录"（字段覆盖切片清单最小集并扩展）；`cut` 追加 `cutExitCode`、`cutHostTs`；
  `reconcile` 追加 `outcome`、`inventoryDigests`、`auditDigest`、`verdict`、`voidReason`。

**写入语义（宿主 PowerShell 5.1，.NET）**：`FileStream(path, FileMode.Append, FileAccess.Write, FileShare.Read, 4096, FileOptions.WriteThrough)`
→ 写行 → `Flush($true)`（= `FlushFileBuffers`）→ `Close` → **读回自证**：以只读方式重开、读最后一条记录、
重算 `recordHash` 并比对，一致才输出 `ACK`（exit 0）。追加语义：只追加、不重写历史行；会话开始必须**全链校验**
（从第一条起逐条重算哈希），链断裂即阻断会话（见 §6 第 10 项）。

**哈希覆盖原文（2026-09-20 修正）**：`recordHash` 覆盖的是**存储行原文**（写入时序列化的那段文本；校验时从
存储行里取回同一段文本），**不得**用“解析后重新序列化”的对象去重算——`ConvertFrom-Json` 会把 ISO 时间戳解析为
`[datetime]`，重新序列化时 7 位小数的尾零会丢失（`...8416210+08:00` → `...841621+08:00`），使校验约 1/10 概率误判，
11 轮会话幾乎必然在中途触发 `journal chain broken` 阻断。回归测试 `journal-timestamp-roundtrip` 固定钉住此形态。
校验还必须**同时要求存储行与规范重建式逐字节相等**（`行 = 原文 + ',"recordHash":"' + 哈希 + '"}'`）：
否则 `recordHash` 之后的字节落在哈希覆盖面之外，“哈希链仍绿但行尾被追加/伪造字段”的篡改将无法检出。
对应回归测试 `journal-tail-append-red`（追加行尾字段后链必须变红）。

**为什么注入不可能丢掉它**：
1. **物理域隔离**：journal 在宿主 NVMe 上，位于客机电源域之外；`controlvm poweroff` 只终止 VM 进程，
   不触及宿主 I/O 路径；
2. **写入即刷**：`WriteThrough` + `FlushFileBuffers` 在确认前强制落盘（数据与文件元数据）；
3. **读回自证**：ACK 以"读回并解析成功 + 哈希相等"为前提；
4. **诚实边界**：以上保证的是"对注入免疫"；宿主自身掉电不在本片注入域内（journal 与任何宿主文件同级，如实声明）。

**确认协议（"journal first; inject only after the record is acknowledged"）**：
写 guest 目标文件**之前**先落 `plan` 记录并取得 ACK（比"断电前"更严格）；`plan.contentDigest` 与客机写探针回传的
实际哈希**双向核对**，不一致 ⇒ 作废（不注入）；无 ACK ⇒ 绝不注入。断电后立刻落 `cut` 记录（同样 ACK）。

### 2.2 注入执行器（宿主）

- 命令（宿主）：`& 'C:\Program Files\Oracle\VirtualBox\VBoxManage.exe' controlvm "Win11B3" poweroff`。
- 执行人：**operator 本人**（B2 纪律：机器侧动作由用户执行）。装置默认**每轮断电前等待 operator 确认**
  （屏幕打印该轮 `plan` 记录 + 询问确认）；`-NoConfirm` 仅限 operator 在会话级显式选择，且仍逐轮落 `cut` 记录。
- 时序：写入探针返回 → 固定 `planCutDelayS = 2.0` 秒 → poweroff → 记录退出码与宿主时间戳。
- 严禁以 `controlvm acpipowerbutton`、`guestcontrol run shutdown.exe` 或任何客机内关机路径替代硬断电（会冲刷存储、伪造存活）。

### 2.3 重启盘点（客机执行，宿主采集）

**客机脚本** `12-round-inventory.ps1`（装置新建，部署到 `C:\prfrail-prep`），参数：
`-RoundRoot R:\calib\rNN -ProbeRoot R:\agent-runner-replay`。输出**单一确定性 JSON 到 stdout**（键序固定、路径排序）：

- `roundRoot`：逐文件列举相对路径、大小、SHA-256、CreationTime/LastWriteTime（**标记 `guest-clock`，只作次级证据**）；
- `probeRoot`：状态 `absent | empty | entries`；若 `entries`，按产品持久化布局枚举：
  `requests/*.jsonl`、`completions/*.jsonl`、`launches/intent.*.jsonl`、`launches/identity.*.jsonl`、
  `terminals/terminal-{intent,closure}.*.jsonl`、`runs/<requestId>/`（含 `events/state-events.jsonl`、
  `manifest.pre/post.json`、`diff.json`、`logs/`、`usage.json`、`managed-process.identity.json`）、
  `ownership.json`、全树 `*.tmp` 残留，并做 **R/C 配对**（request 有/无对应 completion、孤儿 completion）。
  本片标定期内产品区应为 `absent`（B3b 才会产生真实产品工件）——"absent"是标定期的**期望值**，照实记录；
- `dirty`：`fsutil dirty query R:` 输出；
- `boot`：`lastBoot`、uptime 分钟数（客机时钟，次级；实现把 uptime 放 stderr，见 §2.5 偏离 3）。

**执行方式（宿主发起）**：

```powershell
# 宿主
VBoxManage guestcontrol "Win11B3" run `
  --exe "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe" `
  --username prfrail --passwordfile "<host-cred-file>" `
  --wait-stdout --timeout 120000 `
  -- -NoProfile -ExecutionPolicy Bypass -File C:\prfrail-prep\12-round-inventory.ps1 `
       -RoundRoot R:\calib\rNN -ProbeRoot R:\agent-runner-replay
```

- **可重复性证明（验收要求"≥2 次同结果"）**：每轮重启后**连续跑两次**，两次 stdout 必须**逐字节一致**；
  不一致 ⇒ 该轮作废（状态不稳定，不可解释）。盘点本身零客机落盘（只读 + stdout），不污染轮次。
- **耐久证明**：盘点在重启后、对活体客机重导出；两次一致 + 与 journal 摘要对账，证明枚举态稳定而非瞬时读。

### 2.4 对账器（宿主）

`Get-B3aReconcile.ps1`：以 `plan.contentDigest/contentSize` 对盘点结果分类并给出**机器可判**的 `outcome`：

| outcome | 判据（全部机器可判） |
|---|---|
| `survived` | 文件存在，大小相等，哈希 = `contentDigest` |
| `lost-absent` | 目标路径不存在 |
| `lost-zero-filled` | 大小相等，哈希 = 同长全零块的预计算摘要（会话开始算好） |
| `lost-torn-size` | 存在但大小不等 |
| `lost-torn-content` | 大小相等但哈希既非内容也非全零 |

对账输出 `reconcile.json`（判定式 + 引用的盘点/审计摘要），并追加 `reconcile` journal 记录。
**对账只比内容与大小，不比时间**（客机时间一律次级）。

### 2.5 装置脚本清单（本片交付物⑥；实现期按此落地）

| 脚本 | 侧 | 职责 |
|---|---|---|
| `tools/b3a-rig/b3a-rig-lib.ps1` | 宿主 | 共享助手：原始字节进程捕获、`guestcontrol` 封装、哈希、确定性载荷派生、outcome 分类器、标定判定、编码合规写文件 |
| `tools/b3a-rig/b3a-journal-lib.ps1` | 宿主 | journal 追加/ACK/全链校验 |
| `tools/b3a-rig/New-B3aSession.ps1` | 宿主 | 会话前置检查 + 客机工具部署（copyto + 客机侧哈希比对）|
| `tools/b3a-rig/Invoke-B3aRound.ps1` | 宿主 | 单轮编排（§3 S1–S9）|
| `tools/b3a-rig/Invoke-B3aCalibration.ps1` | 宿主 | R0 + R1–R5 + R6–R10 + `calibration-verdict.json` |
| `tools/b3a-rig/Get-B3aReconcile.ps1` | 宿主 | §2.4 分类**模块**（不得声明顶层 `param`，可被点源）|
| `tools/b3a-rig/Invoke-B3aReconcile.ps1` | 宿主 | 该模块的独立 CLI 包装（**禁止被点源**）|
| `tools/b3a-rig/New-B3aEvidence.ps1` | 宿主 | 证据包组装（MANIFEST + SHA256SUMS，B2 格式）|
| `tools/b3a-rig/Test-B3aRig.ps1` | 宿主 | 自检：journal 链篡改必红、ACK 失败必红、时间戳往返必不误判（尾零回归）、盘点确定性（夹具树）、outcome 分类器（夹具）|
| `tools/b3a-rig/guest/12-round-inventory.ps1` | 客机 | §2.3 |
| `tools/b3a-rig/guest/13-round-write.ps1` | 客机 | §4 两种对照的写入/提交链（`-Mode negative|positive`）|
| 复用 `11-shutdown-audit.ps1` | 客机 | 断电审计（拷贝自 guest-prep 主机副本，部署前逐字节核对编码）|

编码：`.ps1` = UTF-8 **BOM + LF**；`.json/.txt` = 无 BOM + LF；脚本注释与输出为英文。

### 2.6 实现期偏离登记（2026-09-20，已回写本档）

实现阶段相对本设计正文的三处偏离，现已全部回写，以便实现与设计基线一致：

1. **多一个共享库**：本设计只列 `b3a-journal-lib.ps1` 一个库文件；实现把与 journal 无关的通用助手放进
   `b3a-rig-lib.ps1`，使 `b3a-journal-lib.ps1` 保持为纯 journal 模块（已在 §2.5 表内登记）。
2. **工具部署改为每轮重做（S3.5）**：§3.1 的“会话级部署一次”在首轮 `restore pristine` 后即失效
   （`C:\prfrail-prep` 被回退），故部署落到每轮 S3.5；代价为每轮 3 次 `copyto` + 3 次客机哈希（约 10 s 量级）。
   S0 第 6 项仅在 VM 运行时执行，否则记为“延后”。
3. **盘点脚本把 uptime 放到 stderr**：§2.3 的 `boot` 段要求 `lastBoot` + uptime 进 JSON，但 uptime 逐秒变化，
   会破坏“两次盘点逐字节一致”的验收要求；故 stdout JSON 只保留 boot-stable 字段，uptime 走 stderr 作次级诊断。
4. **reconcile 模块与 CLI 分离（2026-09-20，缺陷修复）**：`Get-B3aReconcile.ps1` 原带顶层 `param()` 并被
   `Invoke-B3aRound.ps1` 点源；而**点源会按参数名重绑调用方同名变量**（`$Round`/`$Candidate`/…），把轮次脚本
   命令行传入的 `-Round`/`-Candidate` 静默改写成该文件的默认值（`0` / `negative-control`）——实测每轮退化为
   “第 0 轮 / 负对照”，整个标定失效。修复：模块文件**不得声明顶层参数**，CLI 入口拆到
   `Invoke-B3aReconcile.ps1`（不被点源）；并新增自检 `module-load-preserves-vars`（按轮次脚本里实际点源的
   模块清单重放加载，断言变量不被改写）。同一轮还加固了两处：`Add-B3aLogLine` 改用 `FileShare.ReadWrite` + 重试
   （父子进程并发写 `session.log` 的共享冲突曾使父子以 exit 1 一起退出）；标定器在子进程非零退出时记录其
   stdout/stderr，便于 ⑤ 取证。
5. **轮目录布置与失败路径取证（2026-09-20，⑤ 实测修复）**：新增 **S4.5**（客机 `mkdir R:\calib\rNN`）——它属于
   **装置侧布置**，不计入“客机本轮落盘 = 恰一次控制写入”的计数；同时 S6 写入失败时**归档探针 stdout 与 stderr**
   （`write-probe.txt` / `write-stderr.txt`），否则失败轮无法定性（当次 `exit 33` 就是靠此复现才定位到
   `DirectoryNotFoundException`）。注：客机脚本只声明 `exit 0/1`，`VBoxManage guestcontrol` 会把非零映射成自己的
   退出码（实测 33），因此宿主侧日志必须区分“客机脚本退出码”与“VBoxManage 退出码”。
6. **PS 5.1 binder 陷阱（缺陷修复）**：`[pscustomobject]@{ files = @($list) }`（`$list` 为 `List[object]`）在 PS 5.1 下
   抛 `ArgumentException 参数类型不匹配`（`PSToObjectArrayBinder`），使 `Invoke-B3aToolDeploy` **永不返回**
   （S0 检查 6 与每轮 S3.5 必失败）。已改为 `.ToArray()`，并加自检 `deploy-return-shape`。
   结论：装置内**不用** `@(<List[object]>)` 取值，统一 `.ToArray()`。
7. **证据包根路径规范化（缺陷修复）**：`New-B3aEvidence.ps1` 的 `$BundleDir` 原用 `Join-Path $PSScriptRoot '..\..'`
   直接拼接，`Join-Path` 保留 `..` 字面段，导致该路径**长于**真实文件路径，`Substring($BundleDir.Length)` 越界
   （首次真实打包即失败）。已用 `[System.IO.Path]::GetFullPath()` 规范化，并在越界前显式阻断。

## 3. 轮次协议（一轮的完整步骤与门禁）

### 3.1 S0 会话前置（每次装置会话一次，全绿才开始任何轮）

[宿主] 1) `VBoxManage --version` == `7.2.18r175117`；2) VM 存在且 `showvminfo` 与冻结前提比对
（存储控制器 `useHostIOCache=false`、两盘 UUID、**NIC 状态**、快照 `pristine` 存在——按名称查 `snapshot list`）；
3) 宿主 D: 剩余空间 ≥ 5 GB；4) 宿主电源计划只读断言（不许睡眠，见 §6 第 5 项）；5) journal 目录新建 +
全链校验通过；6) 客机工具部署：`copyto` 三个脚本到 `C:\prfrail-prep`，再经客机 `Get-FileHash` 与宿主原文件哈希比对一致
（**仅当 VM 正在运行时**；VM 关机时本项记为“延后”，由每轮 S3.5 承担——见 §2.5 偏离 2）。
任一失败 ⇒ 阻断会话，不动 VM。

### 3.2 S1–S9（每轮，含 S3.5 部署步）

| 步 | 侧 | 动作与命令 | 门禁（不满足即停在本步） |
|---|---|---|---|
| S1 恢复基线 | 宿主 | 断言 VM 已关机（未关机先 `controlvm poweroff` 并记录）；`VBoxManage snapshot "Win11B3" restore pristine` | restore exit 0 |
| S2 启动 | 宿主 | `VBoxManage startvm "Win11B3" --type headless` | exit 0 |
| S3 就绪等待 | 宿主 | 轮询 `guestcontrol ... run --exe cmd.exe -- /c ver` 至 exit 0 | ≤120 s（实测记录每轮实际值） |
| S3.5 工具部署 | 宿主→客机 | `copyto` 三个脚本到 `C:\prfrail-prep`，再经客机 `Get-FileHash` 与宿主原文件哈希比对（**每轮重做**：S1 恢复 `pristine` 会回退 `C:\prfrail-prep`，故 S0 的一次性部署在首轮 restore 后失效） | 三件全部 copied 且客机哈希 = 宿主哈希 |
| S4 注入前门禁 | 客机（宿主发起） | `guestcontrol run powershell -File C:\prfrail-prep\11-shutdown-audit.ps1` | ① uptime < 50 分钟；② `VERDICT` ≠ `LICENCE-SHUTDOWN`。违者：作废 + poweroff + restore + 新轮号重跑 |
| S4.5 轮目录布置 | 客机（宿主发起） | `cmd /c mkdir R:\calib\rNN`（**装置侧布置，不计入被测写入**：每轮 `restore pristine` 会清空 `R:`，不建目录则被测写入以 `DirectoryNotFoundException` 失败，VBoxManage 表现为退出码 33——2026-09-20 实测定位） | exit 0 |
| S5 计划记录 | 宿主 | 按种子/字面量在**宿主**重算 `contentDigest`/`contentSize`，写 `plan`（§2.1）→ flush → 读回 ACK（**先记后写**：无 ACK 则不写、不注入）；记录是“预承诺”，不得抄探针回传值 | ACK=真 |
| S6 控制写入 | 客机（宿主发起） | `guestcontrol run powershell -File C:\prfrail-prep\13-round-write.ps1 -Mode negative|positive -Path R:\calib\rNN\<file> -Size 262144 -Seed b3a-<mode>-<round>`（探针 JSON 只进 stdout，**客机本轮落盘 = 恰一次控制写入**）；探针回传的摘要/大小与 S5 计划**双向核对** | exit 0 且 摘要、大小均与计划一致（否则作废，不注入） |
| S7 注入 | 宿主 | operator 确认 → sleep 2.0 → `controlvm "Win11B3" poweroff` → 立即落 `cut` 记录（ACK） | 确认 + poweroff exit 0 + `cut` ACK |
| S8 重启审计盘点 | 宿主→客机→宿主 | startvm → 就绪 → `11-shutdown-audit.ps1` → `12-round-inventory.ps1` ×2 → stdout 存宿主轮目录 | 审计 `VERDICT` == `HARD-POWER-LOSS`；两次盘点逐字节一致（否则作废） |
| S9 对账 | 宿主 | `Get-B3aReconcile.ps1` → `reconcile` 记录（ACK） | 记录完成 |

### 3.3 轮次有效性、pass/fail/void 精确定义

- **有效轮（valid）**：S1–S9 全部门禁通过。只有有效轮计入标定分母 N。
- **pass / fail（可重试）**：pass = 有效轮产出机器可判 `outcome`；fail = 操作性故障
  （S2/S3 超时、copyto 失败、脚本异常等）→ 同轮号重试，**每轮重试预算 2 次**，耗尽则升级 operator。
- **void（作废，绝不解释、绝不重试进证据）**——满足任一条即作废（恢复 pristine 后用新轮号重跑）：
  1. S4 审计 = `LICENCE-SHUTDOWN`（wlms 干净关机）；
  2. S4 uptime ≥ 50 分钟；
  3. `plan` 未 ACK / 无 ACK 时注入 / 哈希核对不一致；
  4. `controlvm poweroff` 非零退出；
  5. S8 审计 `VERDICT` ≠ `HARD-POWER-LOSS`（含 `GRACEFUL-SHUTDOWN`——干净关机会冲刷缓存，**伪造存活**，必须检出）；
  6. 两次盘点不一致（状态不稳定）；
  7. 本轮 `fsutil dirty query R:` 为脏且恢复动作改动过轮目录（见 §6 第 1 项）；
  8. 前置前提漂移（VM 配置与冻结不符——属会话级阻断）。
- 每控制候选的 void 上限 3 轮；超限 ⇒ 标定失败，交 operator。

## 4. 双向对照标定协议（验收门）

### 4.1 阶段映射

负对照 = 在耐久链**之前**断电（等价于五阶段中的 `temp-written`）；正对照 = 在完整耐久链**之后**断电
（等价于 `parent-synced`）。B3a 标定两个极端阶段的分辨力；B3b 在全部 5 阶段注入候选。

### 4.2 负对照（必须出现丢失）

- 载荷：256 KiB（262144 B = 64 个 4 KiB 群集）确定性伪随机块（种子 `b3a-neg-<round>`，SHA-256 派生，宿主可独立重算）。
- 写入构造（客机，`13-round-write.ps1 -Mode negative`）：`[System.IO.File]::WriteAllBytes($Path, $bytes)` ——
  单次缓冲写入、**不调用任何 flush**（无 `Flush(true)`、无 `FileOptions.WriteThrough`、无目录/卷句柄 flush；
  脚本代码被冻结 + 哈希固化，探针 stdout 列出实际调用的 API 名以资证明）。写后 2.0 s 硬断电。
- 通过规则：**N=5 个有效轮内至少一次丢失**（`lost-*` 任一形态）。5/5 全存活 ⇒ 装置不能制造可观测丢失 ⇒ **装置不可用**，
  **B3b 不得开始**（并据此推断单侧 95% 置信上界 p_loss < 0.45，无法支撑任何"证伪"能力）。
- 记录：每轮 `outcome` + 逐轮分布表（供 B3b 选择其 N 时参考）。

### 4.3 正对照（必须零丢失）

- 载荷：256 KiB，种子 `b3a-pos-<round>`。
- 提交链（客机，`13-round-write.ps1 -Mode positive`，每步成功均为**硬断言**，失败即 exit≠0）：
  1. `FileStream($Path,'Create','Write','None',4096,[FileOptions]::WriteThrough)` → `Write` → `Flush($true)`（`FlushFileBuffers`：文件数据 + 文件元数据）→ `Close`；
  2. **目录项提交**：P/Invoke `CreateFile('\\.\R:', GENERIC_WRITE, FILE_SHARE_READ|FILE_SHARE_WRITE, OPEN_EXISTING)` +
     `FlushFileBuffers` + `CloseHandle`（卷句柄 flush 冲刷该卷全部缓存，含 MFT 与目录索引）。客机 High-integrity 令牌使其可用；
     该步骤是**装置专用标定手段**（§1.8），不触碰 C2/C3 结论；
  3. 探针 stdout 报告 `positive-commit-ok` + 每步 API 结果。
- 写后 2.0 s 硬断电。
- 通过规则：**N=5 个有效轮全部 `survived`**。任何一轮 `lost-*` ⇒ **装置不可用**，**B3b 不得开始**，
  并进入 §6 排查（不得以"放宽门禁"收场）。
- **干净关机伪造的检出**：正对照最怕"被意外干净关机冲刷而伪装存活"。三重防线：
  (a) S8 审计必须 `HARD-POWER-LOSS`（任何 1074/wlms 干净关机 ⇒ 作废）；
  (b) 断电发生在开机 ~3 分钟，而 wlms 触发点在 ~61 分钟，两者在时间上互斥（加上 S4 uptime < 50 断言）；
  (c) 宿主 `cut` 记录 + poweroff 退出码构成权威注入证明。三者任一缺失即本轮不成立。

### 4.4 对照复现轮 R0（不计入 N）

完全复刻 dryrun 配置：14 字节 `pre-cut-marker`、`-Mode negative`、2.0 s 断电。目的：在正式标定前
重证已观测丢失形态（存在 + 全零 + 卷不脏）。若 R0 反而 `survived` ⇒ 记录并**暂停**，operator 决定
（建议先排查延时/环境，再启动正式标定）。

### 4.5 标定总判定（`calibration-verdict.json`，全部机器可判）

| 键 | 值 |
|---|---|
| `journalChainOk` | 全链哈希校验通过 |
| `negativeValidRounds >= 5` | 且 `negativeLosses >= 1` |
| `positiveValidRounds >= 5` | 且 `positiveLosses == 0` |
| `inventoryRepeatabilityOk` | 每有效轮两次盘点逐字节一致 |
| `auditVerdictsOk` | 每有效轮 S8 审计 = HARD-POWER-LOSS |
| `gate` | 全真 ⇒ `CALIBRATED`；任一假 ⇒ `UNUSABLE`（B3b 禁入） |

## 5. 矩阵与预算

### 5.1 五阶段分类法（B3b 注入轴；本片只标定两极端）

`temp-written` → `temp-synced` → `temp-closed` → `record-linked` → `parent-synced`（A7 的五个发布阶段）。
B3a 标定 `temp-written`（负）与 `parent-synced`（正）。

### 5.2 B3a 轮次矩阵（11 次注入 + 轮内 2 次盘点重跑）

| 轮 | 候选 | 阶段 | 载荷 | 期望（门禁） |
|---|---|---|---|---|
| R0 | parity | parity | 14 B `pre-cut-marker` | 复现丢失形态（信息性；异常则暂停） |
| R1–R5 | negative-control | temp-written | 256 KiB 种子派生 | ≥1/5 丢失 |
| R6–R10 | positive-control | parent-synced | 256 KiB 种子派生 | 5/5 存活 |

### 5.3 单轮墙钟预算（实测 48 s 重启为依据）

| 步 | 预算 | 上限 |
|---|---|---|
| S1 restore pristine | 15 s | 60 s |
| S2+S3 启动至就绪 | 48 s（实测 30–48.3） | 120 s |
| S4 审计 | 10 s | 60 s |
| S5 plan 记录+ACK | 2 s | 30 s |
| S6 写入+双向核对 | 10 s | 60 s |
| S7 确认+2.0 s+断电+cut 记录 | 10 s（不含 operator 确认时长） | 60 s |
| S8 重启+审计+盘点×2 | 90 s | 180 s |
| S9 对账 | 5 s | 30 s |

单轮名义 **~3.5 min**，预算 5 min，上限 8 min；全标定 ≈ 11 × 3.5 ≈ **40 min**（加会话前置 ~15 min ≈ 1 h，
不含 operator 确认时长）。

### 5.4 与 50 分钟边界的映射

- 每轮客机 uptime ≈ 2.5 min（S2→S7），距 61 min 触发点余量 > 55 min；每次重启计时归零。
- 边界只约束"客机连续运行"：任何步骤间**禁止让客机空转**；S9 完成即 poweroff + restore（下一轮）。
- operator 若在 S7 前停留 > 45 min：S4 门禁在下一轮重跑时以 uptime ≥ 50 min 作废该轮（不注入、恢复、重跑）。
- 任何一轮审计命中 `wlms` 即作废——与 uptime 断言双保险。

### 5.5 B3b 包络（信息性，本片不定案）

5 阶段 × N 轮 × 2 候选（C2、C3）；N 由 B3b 决定（建议每格 ≥3、理想 5）。以 N=5 计 50 轮 ≈ 3 h。
B3b 必须沿用本片冻结的注入延时（2.0 s）、每轮 restore pristine、轮次协议与盘点/对账工具；改动任一 ⇒ 需重标定。

## 6. 失败模式与恢复预案

| # | 失败模式 | 检测 | 恢复 |
|---|---|---|---|
| 1 | 脏卷（断电撕裂致 NTFS 标记脏） | 每轮 S8 `fsutil dirty query R:` | 首选：**先取证后恢复**——盘点/审计产物先拷回宿主，随后 restore pristine 即回到干净基线；如 operator 想对脏态做尸检：`chkdsk R: /f`（R: 非系统卷，可在线；需 operator 批准），若修复改动轮目录 ⇒ 该轮作废 |
| 2 | 断电后客机启动失败/镜像半应用 | S3 超时、连续启动失败 | 再启动一次；仍失败 ⇒ 记录 + restore pristine + 该轮 fail 重试；绝不手工"修"盘 |
| 3 | 意外干净关机（wlms 或其它 1074） | S4/S8 审计 | 作废该轮（若已发生）；恢复 pristine 重跑；核查 uptime 门禁为何未拦截 |
| 4 | guestcontrol 不可用 | S3 轮询超时 | 核对 VBoxService（客机 `sc query VBoxService`）；超限 ⇒ fail + operator |
| 5 | 宿主睡眠/掉电 | S0 只读断言（`powercfg`）；会话中宿主失电 | 装置**不**改动宿主电源设置（若 operator 决定改，按 B2 纪律 apply/revert 成对）；宿主失电：journal 在宿主盘上可存活，重开会话先全链校验再继续 |
| 6 | VBoxManage 漂移/VM 配置漂移 | S0 版本与前提比对 | 阻断会话，operator 裁决 |
| 7 | 凭据文件缺失/权限不足 | S0 | 阻断会话 |
| 8 | 已观测模式：MFT 已提交/数据全零（负对照丢失态） | R0/R1–R5 | 这是**期望行为**，不作恢复处理；逐轮记录分布 |
| 9 | 已观测模式：6008 滞后 ~50 s | 审计文本 | 客机时间只作次级证据（§1.5），对账只比内容摘要 |
| 10 | journal 链断裂/尾行撕裂（宿主异常） | S0 全链校验 / 每轮 ACK | 阻断会话，operator 裁决（禁止"修链"，只允许人工裁定后新建会话目录） |
| 11 | 盘点重复性失败 | S8 两次盘点不一致 | 该轮作废（状态不稳定不可解释） |
| 12 | 标定门禁失败（负 0/5 丢失 或 正出现丢失） | `calibration-verdict.json` | 装置 `UNUSABLE`，**B3b 不得开始**；排查后重新标定，禁止放宽门禁 |

## 7. 证据包结构（与 B2 包一致）

`docs/validation/evidence/b3a-<YYYY-MM-DD>/`（冻结后入库；活体数据先留在 `D:\VirtualBox VMs\Win11B3\b3a-runs\<date>\`）：

```text
b3a-<date>/
  README.md                      # 复现命令 + 前提引用 + 结论覆盖声明（BOM+LF）
  MANIFEST.md                    # B2 格式：`- path  N bytes  sha256 <hex>`（BOM+LF）
  SHA256SUMS.txt                 # B2 格式：`<sha256 hex>  path\relative`（无 BOM+LF；路径用 `\`）
  calibration-verdict.json       # §4.5 全布尔判定（无 BOM+LF）
  rig/                           # 冻结字节副本：本片全部装置脚本（.ps1 BOM+LF）
  journal/
    journal.jsonl                # 会话 journal 副本（无 BOM+LF）
    chain-verify.txt             # 全链校验输出（exit 0 + 记录数）
  rounds/rNN/
    plan.json                    # plan 记录原文
    cut.json                     # cut 记录原文（含 cutHostTs、cutExitCode）
    write-probe.txt              # 客机 13-round-write.ps1 stdout
    write-stderr.txt             # 客机 stderr（写入失败时才有）
    audit.txt                    # 11-shutdown-audit.ps1 stdout（含 VERDICT）
    inventory-1.json / inventory-2.json   # 两次盘点 + 一致性标记
    dirty.txt                    # fsutil dirty query
    reconcile.json               # outcome + verdict + 引用摘要
  session/session.log            # S0 前置检查输出 + VM 配置比对 + 每轮实际耗时
```

`New-B3aEvidence.ps1` 组装后必须整包 `SHA256SUMS` 校验 0 失配（B2 验收同款）；`Test-B3aRig.ps1` 自检
（`checks=N failures=0`）输出归档。

## 8. 待决问题（需 operator 决策，附建议）

- **Q1 NIC 状态矛盾 — 已结案（2026-09-20）**：实测 `showvminfo` 八张网卡全 none/disabled，且客机内不存在任何真实网卡设备节点（仅 Teredo/6to4 等 Not Present 伪接口）；`.vbox` 中残留的 `<Adapter slot="0"><NAT/></Adapter>` 不呈现硬件、不构成连通性。前提句"NIC=none / 无 IP"成立。
- **Q2 快照 UUID 矛盾 — 已结案（2026-09-20）**：重建后基线为 `pristine` = `7e0287b6-e00e-41ff-975d-2f690a4e981d`，与 `snapshot list` 及 `.vbox` 的 `currentSnapshot` 一致；旧捕获已重捕。装置一律**按名称** `pristine` 引用。
- **Q3 N 的取值**：建议 N=5（二项式论证见 §4.2/§4.3；墙钟 ~40 min 可承受）。请批准。
- **Q4 R0 对照复现轮**：建议纳入（不占 N），R0 存活即暂停待裁。请确认。
- **Q5 宿主睡眠策略**：建议装置只做只读断言、不改宿主电源设置；operator 若需改动请按 B2 成对纪律。请确认。
- **Q6 客机凭据文件路径约定 — 已固定（2026-09-20）**：`D:\VirtualBox VMs\Win11B3\host-only\cred.txt`（仓外），已作为装置默认值实现。
- **Q7 guest-prep 内 `.ps1` 编码审计**：部署前逐字节核对 BOM+LF（provision.log 曾出现控制台乱码），
  违规先规范化再 copyto。建议采纳。
- **Q8 基线 R: 的 2 个条目**：finalize.log `items on R: = 2`（疑似 System Volume Information + 写探针残留）。
  建议：R0 前用 `01-inventory.ps1` 记录基线清单，operator 决定是否清理残留。
- **Q9 `10-fix-timebomb.ps1`**：存在但 operator 已否决使用（许可决定）。本片永不调用；请确认该立场持续有效。
- **Q10 每轮断电前 operator 确认**：默认开启（切片清单"由用户本人执行"）；`-NoConfirm` 为会话级可选项。请确认。
- **Q11 每轮 restore pristine**：默认强制（事件日志与 R: 状态逐轮归零，审计无歧义）；请确认。

## 9. 明确留给 B3b 的裁定

1. C2 / C3 的断电耐久结论（三档 proven / disproven / unresolved），含目录项丢失有序、真孤儿 C、torn 可见等反例判定；
2. 5 阶段 × N 轮 × 2 候选的完整注入矩阵及其 N；
3. ADR-013 候选排序修订、CONTRACTS 耐久口径评估与 first-dispatch 解锁（仅 B3b 出现 `proven` 才评估）；
4. `.tmp` 残留清理策略（A7 E8 遗留）；
5. U1–U4 的定案（本片只提供可复现装置与标定证据）；
6. 盘点工具的"产品布局模式"首次真机使用（本片仅夹具自检 + 标定模式）。

## 10. 复现命令（operator 在宿主 PowerShell 会话逐条执行）

```powershell
# 会话前置（宿主）
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\New-B3aSession.ps1 `
  -VmName Win11B3 -CredFile D:\VirtualBox VMs\Win11B3\host-only\cred.txt

# 完整标定（R0 + 5×负 + 5×正，逐轮确认断电）
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aCalibration.ps1 `
  -VmName Win11B3 -CredFile D:\VirtualBox VMs\Win11B3\host-only\cred.txt -N 5

# 单轮（排障用）
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aRound.ps1 `
  -Round 3 -Candidate negative-control

# 装置自检（任何装置改动后必跑）
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Test-B3aRig.ps1

# 证据包组装 + 校验
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\New-B3aEvidence.ps1
```

## 11. 交付物映射（切片清单 → 本文件）

| 切片要求 | 本文件 |
|---|---|
| ① 环境与前提声明 | §1（含 hypervisor、后端与缓存、NTFS、A7-vs-B3a 差异、许可边界、结论覆盖面） |
| ② 每轮 journal | §2.1 + §3 S5/S7 |
| ③ 断电执行与恢复预案 | §2.2、§3、§5、§6 |
| ④ 重启盘点 | §2.3 + §3 S8/S9 |
| ⑤ 双向对照标定 | §4 |
| ⑥ 脚本、复现命令、证据包 | §2.5、§7、§10 |
