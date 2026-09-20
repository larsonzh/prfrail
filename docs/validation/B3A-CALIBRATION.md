# B3a 断电注入装置与双向对照标定 — 验证报告

- 切片：**B3a**（power-loss injection rig + two-way calibration），依赖 A7
- 日期：2026-09-20
- 结论：**标定通过（`CALIBRATED`）** —— 装置能双向分辨"丢失"与"存活"；**B3b 可以开始**
- 证据包：`docs/validation/evidence/b3a-2026-09-20/`（128 文件，`SHA256SUMS.txt` 126 条，**0 失配**）
- 活体数据：`D:\VirtualBox VMs\Win11B3\b3a-runs\2026-09-20\`（仓库外）

## 0. 判定摘要

| 项 | 结果 |
|---|---|
| 负对照（故意不 flush） | 5 个有效轮中 **4 次丢失**（`lost-zero-filled` ×2、`lost-absent` ×2），1 次存活 → 满足"≥1 次丢失" |
| 正对照（显式 flush + 目录项提交） | 5 个有效轮 **全部 `survived`**（`positiveLosses = 0`） |
| journal 全链 | `journalChainOk = true`，36 条记录（`plan`/`cut`/`reconcile` 各 12），重启后完整且可逐轮对账 |
| 盘点可重复性 | 每轮重启后连续两次盘点 **stdout 逐字节一致**（`inventoryIdentical = true` ×12） |
| 注入真实性 | 12 轮**全部** `auditVerdict = HARD-POWER-LOSS`（无一轮被干净关机"伪造存活"） |
| 卷状态 | 12 轮 `fsutil dirty query R:` 均 `Volume - R: is NOT Dirty` |
| 机器判定 | `docs/validation/evidence/b3a-2026-09-20/calibration-verdict.json`，`gate = CALIBRATED` |

## 1. 切片验收证据映射（逐条）

| 切片要求（`REMAINING_SLICES.md` B3a） | 本次证据 |
|---|---|
| 负对照必须在 N 轮内至少出现一次丢失 | 4/5 丢失（r04/r06 全零、r05/r07 缺失），见 §4 表 |
| 正对照必须 0 丢失 | 5/5 `survived`（r08–r12） |
| journal 重启后完整且可逐轮对账 | 全链校验通过；每轮 `plan`→`cut`→`reconcile` 三记录且末条携带 `outcome`；`journalRecordCount = 36` |
| 盘点结果可重复（≥2 次同结果） | 每轮两次盘点逐字节一致（r01–r12 全真） |
| 固化装置脚本与流程，形成证据包 | 证据包 + 本报告；装置脚本冻结副本在 `evidence/b3a-2026-09-20/rig/`（11 个，与仓库逐字节一致） |

> 边界遵守：**本片不构成任何候选（C2/C3）结论**，也不解除 `unproven`；装置判定的只是"能否分辨丢失"。

## 2. 冻结前提（结论仅覆盖此前提）

| 项 | 值 |
|---|---|
| hypervisor | VirtualBox **7.2.18 r175117** |
| VM | `Win11B3`，UUID `02205fb0-1933-48ba-bbdb-a2c06b4eb9b2`（仓库外 `D:\VirtualBox VMs\Win11B3`） |
| 客机 | **零售 Windows 11 专业版 25H2 build 26200.8037（zh-CN，未激活）**；EFI + vTPM 2.0；2 vCPU / 4096 MB |
| 存储 | SATA/AHCI，`useHostIOCache=false`；`os.vdi` 64 GB 动态（port 0）+ `data.vdi` 8 GB 固定（port 1）；**两盘均 `nonrotational=on`** |
| 注入卷 | `R:` = INJECT，NTFS 7.98 GB（位于 `data.vdi`） |
| 显卡 / 网络 | VBoxSVGA（GA WDDM 已绑定）；`nic1=none`（客机内无真实网卡） |
| 基线 | 快照 **`pristine`**（按名称引用）= `065a25e1-d95c-4f33-8081-4f06216ecce8` |
| 安装介质 | `Win11_25H2_Chinese_Simplified_x64_v2.iso`，8 543 608 832 B，sha256 `7408581e67bc455ebaafb9230e531abf45b1c8864a22114a1b03893f897102e4`（Image #4 = 专业版） |
| 前提差异（vs A7） | A7 主机为 24H2 build 26100；本 guest 为 **25H2 build 26200**（同 GA 主线不同 build） |
| 与设计 §1.6 的差异 | 原评估版 `TIMEBASED_EVAL` 每小时自关机缺陷已由零售版重建消除（70 分钟空转实测无自关机）；窗口/断言/审计降为**廉价安全网**而非许可规避手段 |

**前提修正（本次实测）**：重建后两块盘丢失了 `--nonrotational on`（SSD 标志），S0 门禁拦下（`2-vm-facts` FAIL）；
已重设两盘标志并**重取基线快照**（旧 `7e0287b6…` → 新 `065a25e1…`）。快照会一并保存机器配置，故配置变更后必须重取。

## 3. 装置结构与协议要点

- **脚本（11 个）**：`b3a-rig-lib.ps1`（共享助手）、`b3a-journal-lib.ps1`（journal 模块）、`Get-B3aReconcile.ps1`
  （对账**模块**，无顶层参数）、`Invoke-B3aReconcile.ps1`（其 CLI 包装）、`New-B3aSession.ps1`（S0）、
  `Invoke-B3aRound.ps1`（单轮 S1–S9）、`Invoke-B3aCalibration.ps1`（标定编排）、`New-B3aEvidence.ps1`（证据包）、
  `Test-B3aRig.ps1`（密闭自检）、`guest/12-round-inventory.ps1`、`guest/13-round-write.ps1`（另复用客机 `11-shutdown-audit.ps1`）。
- **每轮步骤（实测均执行）**：S1 恢复 `pristine` → S2 启动 → S3 就绪等待 → **S3.5 工具部署（每轮重做）** →
  S4 注入前审计 → **S4.5 轮目录布置** → **S5 `plan` 记录 + ACK（先记后写）** → S6 客机控制写入 + 探针与计划双向核对 →
  S7 硬断电（`controlvm poweroff`）+ `cut` 记录 → S8 重启 + 断电审计 + 盘点 ×2 → S9 对账 + `reconcile` 记录。
- **journal 哈希约定（冻结）**：`recordHash = SHA256(UTF8(prevHashHex | "null") + UTF8(存储行原文))`，
  存储行 = 规范 JSON 去掉收尾 `}` 后追加 `,"recordHash":"<hex>"}`；校验**要求存储行与规范重建式逐字节相等**。
- **"先记后写"**：`plan` 记录的摘要/大小由**宿主**按种子独立算出（预承诺，不抄探针值），ACK 成功后才允许客机写入；
  探针回传的摘要与大小必须与计划一致，否则作废且不注入。

## 4. 标定结果（12 轮，全部有效）

轮 1 为标定前的**单轮冒烟**（用于验证修好的 S1–S9 链路），随后同会话续跑 11 轮正式标定（journal 续轮号机制生效：
`journal chain ok (3 records); next absolute round = 2`）。

| 轮 | 逻辑位 | 候选 | 阶段 | outcome | 判读 |
|---|---|---|---|---|---|
| r01 | — | parity | parity | `lost-zero-filled` | 冒烟轮（复现既有形态） |
| r02 | R0 | parity | parity | `lost-zero-filled` | 复现轮（设计 §4.4 预期形态） |
| r03 | R1 | negative-control | temp-written | `survived` | 未丢失（负对照允许） |
| r04 | R2 | negative-control | temp-written | `lost-zero-filled` | 丢失 |
| r05 | R3 | negative-control | temp-written | `lost-absent` | 丢失 |
| r06 | R4 | negative-control | temp-written | `lost-zero-filled` | 丢失 |
| r07 | R5 | negative-control | temp-written | `lost-absent` | 丢失 |
| r08 | R6 | positive-control | parent-synced | `survived` | 存活 |
| r09 | R7 | positive-control | parent-synced | `survived` | 存活 |
| r10 | R8 | positive-control | parent-synced | `survived` | 存活 |
| r11 | R9 | positive-control | parent-synced | `survived` | 存活 |
| r12 | R10 | positive-control | parent-synced | `survived` | 存活 |

- 载荷 256 KiB（262144 B），注入延时固定 **2.0 s**，每轮 `restore pristine`，硬断电 = `VBoxManage controlvm poweroff`。
- **分辨力**：负对照丢失率 **4/5 = 0.8**；正对照丢失率 **0/5 = 0.0**。装置在该前提与载荷下能稳定制造可观测丢失，
  且显式 flush 链路无丢失 —— 这正是本片要证明的"双向可分辨"。
- **耗时**：正式标定 04:34:38 → 05:13:20 = **38 分 42 秒 / 11 轮**（≈3 分 31 秒每轮）；冒烟轮 3 分 44 秒。

## 5. 逐轮证据与一致性

每轮目录 `rounds/rNN/` 均含 **9/9 工件**（无缺件）：

`plan.json`、`cut.json`、`write-probe.txt`、`audit-pre.txt`、`audit.txt`、`inventory-1.json`、`inventory-2.json`、
`dirty.txt`、`reconcile.json`

- `auditVerdict = HARD-POWER-LOSS` ×12（**没有任何一轮**出现 `GRACEFUL-SHUTDOWN`/`LICENCE-SHUTDOWN`，
  即不存在"被干净关机冲刷而伪装存活"的情形）；
- `inventoryIdentical = true` ×12（盘点可重复，且对账**只比内容与大小**，从不比时间）；
- `dirty = "Volume - R: is NOT Dirty"` ×12（无一轮出现脏卷，无需 `chkdsk` 恢复）；
- 探针证明写入 API：负对照 `apisCalled = ["File.WriteAllBytes"]`（无任何 flush API）；
  正对照为 `FileStream + Flush(true) + CreateFile(\\.\R:) + FlushFileBuffers + CloseHandle`（`commit = positive-commit-ok`）。

## 6. journal 链自证

- 36 条记录 = `plan` 12 + `cut` 12 + `reconcile` 12；每条写入路径为
  `FileStream(Append, FileShare.Read, WriteThrough)` → `Flush($true)` → 关闭 → **读回重算哈希 + 逐字节行形状校验** → ACK。
- 会话开始时**全链校验**通过（`journalChainOk = true`）；标定器只统计**本次调用**产生的轮次，且额外要求
  每轮 `inventoryIdentical` 与 `auditVerdict` 合格，才可能给出 `CALIBRATED`。
- 物理域隔离：journal 在**宿主** NVMe（`D:\VirtualBox VMs\...`），位于客机电源域之外；注入只终止 VM 进程。
  诚实边界：宿主自身掉电不在本片注入域内（与设计 §2.1 一致）。

## 7. ④ 独立终审记录（Codex，只读）

| 轮次 | 判定 | 要点 |
|---|---|---|
| 第 1 轮 | **FAIL** | 3 High（journal 顺序违反"先记后写"；`Wait-B3aReady` 缺 `-TimeoutS`；`-N` 无下限）+ 3 Medium + 2 Low |
| 第 2 轮 | **FAIL** | 8 项修复全部确认 `fixed`；但新增 1 High：哈希只覆盖 `recordHash` 前缀，"行尾追加字段"仍能保持链绿 |
| 第 3 轮 | **PASS** | 该 High 确认 `closed`；无剩余 Medium+ |
| 第 4 轮（专项） | **FAIL** | 三处晚期缺陷修复本身成立，但 3 条 Medium 均为**守卫强度**（守卫可能空跑通过、未断言 `files` 为数组、重试耗尽分支漏记子输出）→ 已全部加强 |

> 调用次数超出设计上限（1 + 1 复审）。理由：B3a 是硬门切片、且中期引入了新的 High 与新的修复批次，
> 按"④ 直至无 Medium+（不得自行判定通过）"的闭环纪律执行；成本已如实记录。

## 8. ⑤ 原生验证发现的缺陷与修复（自检与评审都未抓到）

| # | 缺陷 | 性质 | 后果 | 修复 + 回归测试 |
|---|---|---|---|---|
| D1 | `Get-B3aReconcile.ps1` 带顶层 `param()` 却被点源 | **Critical（静默）** | 点源按参数名重绑调用方同名变量 → 每轮 `-Round`/`-Candidate` 被改写成 `0`/`negative-control`，标定整体失效（首跑卡在 `r00`） | 模块**不得有顶层参数**；CLI 拆到 `Invoke-B3aReconcile.ps1`；自检 `module-load-preserves-vars`（从轮次脚本解析点源清单并重放，且要求清单非空且含已知模块） |
| D2 | 父子进程并发追加 `session.log`（`StreamWriter` 默认 `FileShare.Read`） | High | 共享冲突 + fail-closed → 父子双双 exit 1，并丢掉解释性日志 | `FileShare.ReadWrite` + 最多 5 次退避重试；子进程非零退出时记录其 stdout/stderr（且置于预算判定**之前**） |
| D3 | `[pscustomobject]@{ files = @($List[object]) }` | High | PS 5.1 binder 抛 `ArgumentException 参数类型不匹配` → `Invoke-B3aToolDeploy` **永不返回**，S0 检查 6 与每轮 S3.5 必失败 | 改 `.ToArray()`（两处）；自检 `deploy-return-shape`（伪 `VBoxManage` 路径跑通返回路径，并要求 `files` 是数组且 3 项） |
| D4 | 轮目录从未创建 | High | 被测写入以 `DirectoryNotFoundException` 失败；客机脚本只声明 `exit 0/1`，VBoxManage 映射成 33，宿主只见无解释的 33 | 新增 **S4.5**（客机 `mkdir R:\calib\rNN`，属装置侧布置、不计入被测写入）；S6 失败时归档探针 stdout/stderr（`write-probe.txt`/`write-stderr.txt`） |
| D5 | `New-B3aEvidence.ps1` 的 `$BundleDir` 未规范化（`Join-Path '..\..'`） | Medium | 字面路径长于真实文件路径 → `Substring` 越界，首次真实打包即失败 | `[System.IO.Path]::GetFullPath()` 规范化 + 越界前显式阻断 |

同类排查（该轮复审）：其余脚本未发现"点源带顶层 `param` 的脚本"形态；装置内统一改用 `.ToArray()` 取 `List[object]`。

## 9. 反例与变异验证（证明测试不是恒真）

| 测试 | 反例/变异 | 结果 |
|---|---|---|
| `journal-timestamp-roundtrip` | 把校验器还原成"解析后重新序列化" | **确定性变红**（时间戳第 7 位小数尾零被丢：`…8416210+08:00` → `…841621+08:00`；旧实现约 1/10 概率误判，11 轮会话约 69% 会在中途 `journal chain broken` 阻断） |
| `journal-tail-append-red` | 去掉行形状（逐字节）校验 | **确定性变红**（行尾追加字段后哈希链仍绿） |
| `verdict-requiredN-floor` | `RequiredN = 1` | 必须抛错（拒绝用单轮样本产出 `CALIBRATED`） |
| `module-load-preserves-vars` | 点源带顶层 `param()` 的模块 | 变量被改写即变红（D1 的类级守卫） |
| `deploy-return-shape` | `@(List[object])` binder 缺陷 | 抛异常或 `files` 非数组即变红（D3 的类级守卫） |

密闭自检规模：**22 项，0 失败**（连续 10 次独立进程 + 证据包归档那次全部 `checks=22 failures=0`）。
静态分析：`Invoke-ScriptAnalyzer -Severity Warning,Error` = **0**（含本片清零的 14 条风格告警）——
原始输出已归档为证据包内 `verification/analyzer.txt`（`findings=0`）。
编码：11 个装置脚本 + 证据包内 `.ps1`/`.md` 全部 **UTF-8 BOM + LF**；证据包内 `.json`/`.txt` 无 BOM + LF。

## 10. 已知限制与诚实边界

1. **单设备标定**：结论只覆盖本前提（VirtualBox + VDI 后端 + 宿主 NVMe + NTFS + 上述 build）。换 hypervisor、
   换盘后端、或改注入延时/载荷大小 ⇒ **必须重做标定**。
2. **丢失率不是候选结论**：负对照 4/5 只说明"装置能制造可观测丢失"；5 阶段 × N × 候选的耐久定案属 **B3b**。
   本片不做单轮/单样本外推。
3. **journal 的威胁模型**：无密钥哈希链提供的是**篡改可见性**，不是对"拥有完整改写能力并重算后续链条"的对抗保证；
   该边界由证据包冻结 + 整包 SHA-256 复核补强。宿主掉电不在注入域内。
4. **证据包字节口径**：包内 `rig/` 副本与仓库当前字节**逐字节一致**，且 10 个测量路径脚本的最后修改时间均**早于**
   标定开始（04:34:38），即"冻结字节 == 实际运行字节"；仅打包脚本 `New-B3aEvidence.ps1`（05:14:28）在运行后
   修正了 D5，不影响测量链路。上述两条声明的可独立复核工件：`verification/rig-vs-repo.txt`（逐脚本 `identical`
   + `repoMtime` + 与标定开始的先后判定；实测 10 个测量路径脚本均为 `before-calibration-start`，仅打包脚本为
   `AFTER-CALIBRATION-START`）。
5. **会话目录含冒烟轮**：`r01` 是标定前的单轮冒烟（同为 parity/负向写入），已在 §4 明确标注；标定器判定只统计
   本次调用产生的 11 轮（验证见 `calibration-verdict.json` 的 `negativeValidRounds`/`positiveValidRounds`）。
6. **凭据文件保留（待用户决定）**：`D:\VirtualBox VMs\Win11B3\host-only\cred.txt` 为 B3b 复用同一环境所必需；
   删除会使该环境不可操作。建议**保留至 B3b 结束**后再删，或在 B3b 前重新设定一次性口令。文件始终在仓库外、从未打印。

## 11. B3b 移交条件

B3b **必须**沿用本次冻结的：注入延时 **2.0 s**、每轮 `restore pristine`、S1–S9（含 S3.5/S4.5）轮次协议、
journal 约定与对账/盘点工具、前提声明。改动任一 ⇒ 需**重新标定**。B3b 的 5 阶段 × N × 候选矩阵与 N 由 B3b 决定
（本片已给出负对照逐轮分布供其选 N 参考）。

## 12. 复现命令（宿主 PowerShell 5.1）

```powershell
# S0 会话前置（全绿才开始任何轮）
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\New-B3aSession.ps1

# 单轮（默认每轮断电前等 operator 确认；-NoConfirm 为会话级选择）
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aRound.ps1 `
  -Round 1 -Candidate parity -NoConfirm

# 双向对照标定（本报告所用）
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aCalibration.ps1 -NoConfirm

# 密闭自检 / 证据包
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Test-B3aRig.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\New-B3aEvidence.ps1
```

## 13. 交付物与开放项

**交付物**：① 环境与前提声明（§2）；② 装置脚本 11 个（仓库 `tools/b3a-rig/`）；③ 双向对照标定结果（§4）
④ 逐轮 journal + 盘点 + 对账工件（证据包 `rounds/`、`journal/`）；⑤ 证据包（`MANIFEST.md` + `SHA256SUMS.txt`，0 失配）；
⑥ 本报告（中英双语）。

**开放项**：
1. 凭据文件保留/轮换决定（§10 第 6 条）。
2. 早期评估版环境曾以明文回显口令，建议轮换该（已退役）环境口令——仅提示，不影响本片。
3. B3b 启动前请确认沿用本片冻结参数（§11）。

**提交与 CI（2026-09-20）**：`136352b`（`feat: add the B3a power-loss injection rig with a hermetic self-test`，`tools/b3a-rig/`
11 个脚本）→ `243ed82`（`docs: record the B3a calibration, its evidence bundle and the verification report`，含 128 文件证据包）
→ `babc041`（`docs: add the T027 working documents (B3a design, slice list, operating directives)`），已推送 `origin/main`
（`815c03f..babc041`；**未推 gitee**）。GitHub Actions run `35482066033`（head `babc041`，**1 分 48 秒**）：
**Go ubuntu-latest success**（含 Linux 专属 Race 与 Contract fixtures）、**Go windows-latest success**，三个
`workflow_dispatch` 专用任务按设计 **skipped**。推送前已用 `git worktree add --detach`（克隆语义）复核整包：
**128/128 条哈希 0 失配**（入库字节与 `SHA256SUMS.txt` 自洽，满足 B2 的字节保真教训）。
