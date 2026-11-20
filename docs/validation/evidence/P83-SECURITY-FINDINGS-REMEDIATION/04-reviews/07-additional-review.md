# P83 ⑦ 追加复审（第 6 轮）逐字转录 · 2026-11-21

- 载体：`GPT-6.2 Sol (copilot)` Max。设定时刻 = 2026-11-21；档位值 = Max；未暂停。
- 授权：用户 2026-11-20 裁定「追加 1 轮 ⑦ 复审（收窄输入）」；**超出 §7.5 与本次切片额度裁定，系用户直接裁定、不得作为先例**。
- 输入收窄 = 四项 Medium 整改 diff + `process_linux.go` 自引入回归修正 + T2/T3/T8 顺延登记的台账/报告措辞。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。
- 结论：`ADDITIONAL RE-REVIEW: FINDINGS`（范围内 2 项 Medium 阻断：T4、T5#4）。

---

本轮缩小范围内仍有 **2 项 Medium 阻断**：T4 的 Windows 界外归属检查，以及 T5#4 的同进程批量 proof 回写。其余受审修复未见 Medium 或以上阻断。

| 项 | 状态 | 防护位置 | 退化时会变红的测试 |
|---|---|---|---|
| T4 枚举 | **Medium，未闭合。** Linux 在受限打开的描述符上调用 `Readdirnames`。Windows 后代末级 junction 现在会进入 reparse 拒绝分支；普通界外路径也会被卷号或最终路径检查拒绝。但在同卷大小写敏感目录中，不同目录的路径可仅大小写不同；中间 reparse 仍可能被跟随，`EqualFold` 无法证明目标属于授权树。根目录用 `rel=""` 时仍枚举最初绑定的 junction 目标，并核对卷号和文件 ID；显式传 `"."` 则在清洗阶段被拒，与注释不符。 | bounddir_unix.go:287、bounddir_windows.go:170、:196、:297、bounddir_clean.go:29 | Linux 描述符枚举（p83_fixround8_linux_test.go:21）、Windows 界内 junction 拒绝（p83_fixround8_windows_test.go:181）、根对象身份（:52）。缺专项：大小写折叠相等但对象不同的界外目标必须拒绝；根 junction 用 `""` 仍应成功。 |
| T5#3 保留代读取 | **通过，有残余。** Windows 在同一打开句柄上拒绝末级 reparse 和普通目录，后续只读该句柄并校验记录。受信本地租约目录前提下，剩余窗口可登记为非阻断。 | lease_open_windows.go:45、lease.go:362、:459 | reparse 分类（guard/p83_fixround8_windows_test.go:45）、目录分类（:67）、超限拒绝（guard/p83_fixround7_tester_test.go:115）。 |
| T5#4 代际到 proof | **Medium，未闭合。** 直接恢复已要求同一持有进程、同一 claim、恰好 `N-1`、匹配的 `PreviousHash` 和可验证的保留记录。真实前驱可推进，且当前任务先写新 proof 才构造完成结果。但批量回写只比较 `ResourceID`，可把同资源、不同 claim 或哈希不符的其他任务 proof 改成当前有效 proof，绕过直接恢复的前驱校验。 | runtime_agent.go:2104、:2156、runtime_supervised.go:1194、:1893 | 拒绝无关前驱（console/p83_fixround8_tester_test.go:83）、接受真实前驱（:56）。缺专项：同资源但无关的其他任务 proof 不得被批量改成可发布 proof；按现实现预期会红。 |
| T8b 最终拒绝证据 | **通过，限 `start-failed`。** 最终摘要只包含两个闭集分类；不读取 executor 的 `Evidence` 或 `Measurements`，自由文本 `FailureCode` 映射为 `unclassified`。不要把此保证扩写到 `unknown` 或退出后的不可用结果。 | ai_probe_copilot.go:121、:241 | 最终拒绝证据闭集测试（adapters/p83_fixround8_tester_test.go:24）。 |
| Linux PATH 回归 | **通过。** 无批准摘要时保留 `exec.Command` 从父进程 PATH 解析的路径。批准分支才改写 `cmd.Path`，通过 `ExtraFiles` 将已验证对象交给目标子进程的 fd 3，启动失败不回退到按原名打开。 | process_linux.go:95 | 批准描述符仅由目标子进程继承（guard/p83_fixround8_linux_test.go:27）。缺专项：无摘要的裸命令仅在 PATH 中可找到，仍须成功启动。 |

## 报告残余（复审原文，可原样录入）

> **T4：** Windows 枚举根时，若绑定入口是 junction，枚举对象是绑定时打开的目标目录。当前调用使用空相对路径表示根；显式 `.` 尚未获接受。Windows 的重开及最终路径复核不是 Linux `openat2` 式的原子目录相对约束。
>
> **T5#3：** Windows 已消除保留代记录末级名称的 `Lstat` 到打开之间的类型竞态。64 KiB 是字节上限，不是时间上限；同步打开和读取无硬截止，也未绑定父目录句柄。该残余仅在运行根及租约父目录由受信本地权限保护时接受。
>
> **T5#4：** 直接恢复仅证明当前持有进程的一代前驱；不承诺跨进程或跨多代自动恢复。租约代际提交与 proof 写入不是同一事务。批量回写缺少前驱验证是本轮阻断，不是已接受残余。
>
> **顺延：** F3（T2 STOP）缺少 AgentRunner 可持久复验的可信冻结策略制品；F6（T3）仍有 Windows 名称到启动映像的绑定窗口，以及 Linux 同 inode 不等于字节不可变、shebang 解释器重读的边界；F8（T8a）未证明 CLI 默认落盘日志已关闭或脱敏，默认也未拒发 BYOK。T9 收窄和 T8b 修复均不关闭其对应的顺延项。

## 部分闭合措辞判定（复审原文）

`PARTIALLY CLOSED` 准确，但拟列的已闭合清单不准确。**F1 因 T4、F4 因 T5#4 暂不能记为闭合**。F2/F5/F7 沿用既有结论，本轮不重签；T9 收窄（ai.go:387）可单列为闭合；F3/F6/F8 按上述边界顺延。不能据此登记「全片安全闭合」。

本轮仅做静态只读复审，未修改、未提交，未运行测试或 lint。
