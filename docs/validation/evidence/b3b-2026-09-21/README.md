# B3b 证据包

本目录由 `tools/b3b-rig/New-B3bEvidence.ps1` 生成，请勿手工编辑；任何结论都可由 `rounds/`、
`journal/` 与 `matrix/` 逐字节复算（`matrix.json` 无独立累加器状态）。

## 会话元信息

- 会话根：`D:\VirtualBox VMs\Win11B3\b3b-runs-main\2026-09-21`
- 生成时间：2026-09-21T12:09:05
- 装置快照：`rig/b3b-rig/`（本次运行所用脚本）×`rig/b3a-rig/`（被复用且冻结的 B3a 文件）
- 设计权威：`docs/t027/B3B_DESIGN.md`

## 三档判定

| 臂 | 判定 | 依据 |
| --- | --- | --- |
| c1 | unresolved |  |
| c2 | proven |  |
| c3 | unresolved |  |
| c2c3 | unresolved |  |

- 结论口径（claimScope，逐字引用）：within premise P (VirtualBox Win11B3 / NTFS on SATA-AHCI, useHostIOCache=false, record-level ~1 KiB payload, 2 s cut window, N=5 per cell): no counterexample observed, candidate mechanism demonstrably executed, device discrimination demonstrated at stage temp-written only; not a long-term guarantee and not a claim about other hypervisors, filesystems, payload sizes or the host losing power
- 前提（premise）：VirtualBox 7.2.18 r175117, guest Windows 11 25H2 build 26200.8037 NTFS, SATA/AHCI useHostIOCache=false, inject volume R:=INJECT (data.vdi 7.98 GB, 4096 B clusters), guest hard poweroff via VBoxManage controlvm poweroff, pristine snapshot 065a25e1-d95c-4f33-8081-4f06216ecce8, nic1=none, C: is the signalling domain only, the measured domain is the replay store root on R:, one measured publication per round

## 完整性

- 必需项（`-Required`）：rounds、journal.jsonl、session.log、pin.json、matrix.json、verdict.json
- 构建时缺失的必需项：无
- 未复制的会话文件/可选分区：calibration、u4、guest
- 冻结文件快照：完整
- 日志链校验：chain ok=True records=337 first bad line=0 first error=
- 会话根布局映射：`journal.jsonl`→`journal/journal.jsonl`、`session.log`→`session/session.log`、`pin.json`→`session/pin.json`、`matrix.json`/`verdict.json`→`matrix/`、`rounds/` 原位复制
- 冻结文件比对基准：会话 S0 记录的 `session/pin.json`（非当前工作副本哈希）
- 装置脚本与仓库工作副本哈希比对：`verification/rig-vs-repo.txt`（mismatches=0）
- 静态分析：`verification/analyzer.txt`

## 可复算性边界

矩阵与三档判定可由 `rounds/`、`journal/journal.jsonl` 重算；
`Invoke-B3bMatrix -Finalize` 会从 `plan.json` + `inventory-1.json` 重新派生outcome，
并逐文件比对 `reconcile.json` 与 journal 记录的哈希；不一致即拒裁（exit 4），不作为轮次计入。
因此本包内的任何数字都不得仅凭 `matrix.json` 引用——必须可回到上述原始工件。

## 复算

```powershell
# 重算矩阵与三档判定（不触碰客机）
pwsh -File tools\b3b-rig\Invoke-B3bMatrix.ps1 -Finalize -RigRoot "D:\VirtualBox VMs\Win11B3\b3b-runs-main\2026-09-21"
# 装置密闭自检
pwsh -File tools\b3b-rig\Test-B3bRig.ps1
```