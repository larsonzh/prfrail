# P83 ⑦ 复审 r1 / r2 / r3 逐字转录 · 2026-11-20

- 载体：`GPT-6.2 Sol (copilot)` Max。设定时刻 = 2026-11-20 本切片 ⑦ 环节；档位值 = Max；未暂停。
- 额度：终审 1/1 用尽；复审 3/3 用尽（另加用户 2026-11-20 裁定的追加 1 轮，另记、不作先例）。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名；文中 `api.deepseek.com` 系测试链配置中的公开提供方域名。
- 口径：⑤（2/2）与 ⑥（2/3）的 PASS 均绑定整改前快照；下表三轮回审是覆盖最终快照的唯一评审。

---

## RE-REVIEW 1

**RE-REVIEW: FINDINGS**

| T-id | status | 已核守卫 | 移除守卫会变红的测试（静态判断） |
|---|---|---|---|
| T1 | partially | agent_runner_replay_dispatcher.go:140 先写 R，再读派生账本，读完才取时钟；派发复核失败不启动。但装配初筛 runtime_agent.go:1545 仍使用内存账本。 | 派发重读撤销（p83_security_findings_test.go:288）断言 R 留存且无启动记录；T1b（p83_fixround5_tester_test.go:120）只钉住 review 的读后取时，不是派发。 |
| T2 | partially | runtime_agent.go:846 只接受同 ID 最新 grant；runtime_supervised.go:1002 的策略摘要与风险确认检查（:1039）分离。Agent 独立策略比对仍是 STOP。 | T2a 旧 grant（:173）、T2b 风险确认（:225）。 |
| T3 | partially | trusted_ledger.go:51 拒绝异址账本；ai.go:181 比对含 `secretRef` 的 profile hash 和完整绑定。默认检查器没有沿用批准摘要，启动前也未证明同一文件。 | T3a 异址（:277）、T3b 换引用（:335）。T3c（adapters/p83_fixround5_tester_test.go:40）手工传入已批准摘要，绕过 CLI 断链；身份测试（:87）只测了长度不同。 |
| T4 | partially | runner.go:73 有 globs 才预绑定，句柄复用；runner.go:263 无句柄时拒绝收集。重试的执行目录和产物根仍可能分离。 | T4b 根替换（gates/p83_fixround5_tester_test.go:76）、根不可打开（p83_f1_f5_tester_test.go:138）；T4a（:40）证明无产物可通过，未直接观察是否调用 open。 |
| T5 | partially | runtime_agent.go:1947 拒绝活跃前任；lease.go:242 按 hash 保留旧代。接管未转移持有人，历史留档也未全程失败关闭。 | T5a 活跃前任（console/p83_fixround5_tester_test.go:432）、T5b 旧代读回（guard/p83_fixround5_tester_test.go:16）；未覆盖两个后继进程或真实旧回执的复核。 |
| T6 | closed | runner.go:58 在 Prepare/Execute 前核对真实定义摘要；runner.go:435 从 Hook 取恢复保证和外部系统，拒绝矛盾。仓内正式调用方都对传入的同一 Hook 算摘要。 | T6a 错摘要（:143）、T6b 改定义（:166）、scope 矛盾（runner_effect_policy_test.go:55）。 |
| T7 | partially | process_windows.go:527 追加终止符前安全扩容并清零旧块；secret_input_windows.go:67 清零终端旧缓冲。两处均有 observer，终止符分支缺反证。 | T7a 普通扩容（:21）、T7b 终端扩容（:21）。删掉终止符守卫不会让 T7a 变红。 |
| T8b | partially | ai_probe_copilot_executor.go:173 只散列稳定分类，但外层证据（ai_probe_copilot.go:219）仍散列完整 argv。 | T8b 稳定证据（:118）只比较执行器输出，不比较最终可用性证据。 |
| T9 | partially | review.go:151 恢复仅 `policy` 回执强制绑定；ai.go:135 的 host-auth 豁免未贯穿 CLI。 | T9a 人工回执（chain/p83_t9a_tester_test.go:18）；T9b（:402）直接调 helper，没经过 CLI 读账本。 |

**发现**
1. **High，阻断是，T3**：ai.go:469 按路径重算 `ExpectedExecutableHash`，丢失了账本核准的值。即使保留原摘要，执行器（ai_probe_copilot_executor.go:53）也只比较大小、mtime 和 mode，随后仍按路径启动；这些元数据不是文件身份，打开句柄也没有锁住启动对象。
2. **High，阻断是，T5**：runtime_agent.go:1950 接管只返回旧 token；续租（lease.go:101）保留旧 PID。前任停止后，两个后继进程可通过同一份租约校验；supervised 发布（runtime_supervised.go:1861）从磁盘 proof 直接续租和验证，不走前任检查。留档（lease.go:268）也把「路径已存在」当作成功，坏档不会阻止覆盖当前代。
3. **Medium，阻断是，T4**：runner.go:132 按路径启动每次尝试，Glob（:268）也按路径枚举，读取却使用首轮绑定的句柄。重试期间根发生更换时，执行可落在新目录而回执采信旧目录同名产物。现有替换测试只跑一次。
4. **Medium，阻断是，T8b**：ai_probe_copilot.go:99 在处理 `start-failed` 前加上包含完整 argv 的摘要。合法配置里的 model 值也会进入 argv；因此同类拒绝的最终证据仍可能随值变化。此项不涉及已延期的 BYOK 落盘拒发。
5. **Medium，阻断是，T1/CLI**：runtime_agent.go:1547 在派发前用内存 grant 和预算做准入；不协调的路径又在装配入口（:229）直接拒绝，尚无 R。CLI 旧的任意账本路径现在会被拒绝，而 cli.go:457 仍要求两个 flag，导致新加的省略参数分支不可达；`approvals` 和 `cost` 的默认写入位置仍是工作目录，不是 run 目录。
6. **Medium，阻断是，T9**：ai.go:336 无条件加载账本。合法的无 `baseUrl` host-auth 配置没有账本时，在到达豁免分支前就失败。
7. **Medium，阻断是（验证缺口），T7**：process_windows.go:497 的 builder scratch 清零调用与失败部分块没有直接反证；现有 T7a 也未命中恰满后追加终止符的分支。代码有清零，测试尚不能防这几处回归。
8. **Low，阻断否，T5 留档增长**：lease.go:254 每次替换多留一份，按代数线性增长，无上限或生命周期。记录保留策略时须保住仍被回执引用的代。

**T2 STOP**
可以登记为后续片的开放项，**不能算本片关闭或自动豁免**。STOP 属实：run-manifest schema（schemas/run-manifest.schema.json:142）只声明字段，没有 AgentRunner 写出可供独立对照的冻结策略；发布复核（runtime_agent.go:994）比较的仍是同一账本来源的 grant 摘要。应明确登记：当前仍可能接受一个合形但未经独立对照的摘要；补齐需要新增受执行者不可写保护的、可持久复核的策略证据，绑定 run/task/step，并在 review 和 publish 都验证 `grant.PolicyHash` 等于该摘要。验收反证须覆盖缺失、不可读和合形但不匹配，均不得完成发布。若要带此缺口关闭 P83，需另有明确用户裁定，并与已批准延期的 T8a 分开登记。

只读核对暂存区；`git diff --cached --check` 通过。未重跑 Go、Windows 或变异测试，未修改文件、未提交。⑤/⑥ 的 PASS 仅适用于整改前快照。

---

## RE-REVIEW 2

**RE-REVIEW: FINDINGS**

| 项目 | 状态 | 已核守卫 | 会变红的反例 |
|---|---|---|---|
| T3 | **未闭合** | ai.go:396 传递 grant 原摘要，ai.go:531 拒绝空摘要。执行器拒绝裸命令名（:208），启动前复核身份和摘要（:284），但仍按路径启动。 | T3a 双重比对（adapters/p83_fixround6_tester_test.go:54）、T3b 空摘要（console/p83_fixround6_tester_test.go:226）；没有覆盖最后校验到实际启动的反例。 |
| T5 | **接管通过，保留和恢复未闭合** | lease.go:107 续租匹配 holder，lease.go:135 锁内检查旧 token 和先任证明，再写新 holder 和代次。supervised 发布（runtime_supervised.go:1869）调用同一证明。已有保留文件确实复读校验（lease.go:327）。 | T5a 旧 token（console/…:359）、T5b 外来身份（guard/…:25）、T5c 损坏保留文件（:72）、T5d 外来发布证明（console/…:419）。 |
| T4 | **部分通过** | runner.go:115 每次尝试前验根并更新工作目录；runner.go:299 从 `RootPath()` 取 glob 基址，但 Linux 返回的是旧路径。 | T4a 根已替换时不重试（gates/…:48）；T4b（round5:76）只排除读到替换目录的字节，不排除从它枚举文件名。 |
| T8b | **实现通过，验证缺口** | ai_probe_copilot.go:102 的 `start-failed` 分支绕过通用 argv 摘要，只使用拒绝分类及执行器分类标记。 | T8b 稳定性测试（:115）只检查 Execute；恢复最终 argv 摘要不会让它变红。 |
| T1/CLI | **通过，有时间点残余** | 准入（agent_runner_admission.go:126）重读两份账本，读错即拒绝；装配路径冲突（runtime_agent.go:235）在写 run 状态前拒绝并明确说明。run（cli.go:586）、approvals（:237）、cost（:106）按 run 目录派生，并拒绝显式冲突。 | T1a 省略参数及冲突（:121）、T1b 按 run 目录读取（:179）。 |
| T9 | **通过** | ai.go:387 只在配置 `baseUrl` 时加载接收端授权账本。 | T9a 无账本的 host-auth（:258）。 |
| T7 | **通过静态核对** | process_windows.go:514 延迟清零转换 scratch；局部失败（:521）和容量恰满（:550）都清零后报告实际缓冲区。 | T7a 终止符扩容（:38）、T7b 局部失败及转换 scratch（:87）。 |

**发现**
1. **High，阻断 · T3：最终校验没有绑定实际启动。** 执行器（:288）仍把路径交给 Linux 的 `cmd.Start`（process_linux.go:63）和 Windows 的 `CreateProcessW`（process_windows.go:278）。Windows 的原句柄还允许共享写入和删除（boundopen_windows.go:16）。可执行对象在最后校验后仍可能变化，未获批的映像可能收到凭据。两次路径校验不能证明启动的是已校验对象。
2. **Medium，阻断 · T4：验根后仍可能执行或枚举到另一棵树。** Linux `RootPath`（bounddir_unix.go:252）返回绑定时记录的路径；hook 后的 `Glob`（runner.go:305）可从替换目录取文件名，随后从原句柄读同名文件并成功。重试前的检查也不绑定后续按路径启动的工作目录。
3. **Medium，阻断 · T5：保留记录仍可假成功或卡住锁。** retainGeneration（lease.go:331）用 `Stat` 和 `ReadFile`（:388）跟随路径。若保留名是当前记录的符号链接，复读可通过，但换代后该名字不再给出旧代原件。Linux 上非普通文件还可能使读取在 lease 锁内阻塞，超时无法打断该读取。
4. **Medium，阻断 · T5：换代后写回失败，重入无法恢复。** supervised 发布（:1869）先提交新代，再写回任务凭证（:1879）。若写回前中断或写回失败，任务仍持旧凭证；重入在 token/hash 比对（runtime_agent.go:2086）就被拒绝，没有对账路径。批量刷新还会跳过读取错误并返回成功（runtime_supervised.go:1186）。
5. **Medium，阻断验收 · T8b：最终证据没有反例。** 现有测试（:115）不经过 Probe（ai_probe_copilot.go:102）。需要断言同类、不同输入的拒绝产生相同的最终证据，并且恢复通用 argv 摘要时该断言会失败。
6. **Medium，非阻断 · T1：前置准入在读账本前取钟。** 取钟（agent_runner_admission.go:103）早于读取两份账本；授权在读取期间到期仍可能通过前置准入。后置重确认仍阻止启动，所以这不是启动绕过。

**报告残余**
- **T2 STOP**：AgentRunner（runtime_agent.go:817）只认最新 grant，但没有持久化的冻结策略对象可与 `grant.PolicyHash` 对证；不能报 F3 全闭合。**T8a** 按裁定顺延 P84：没有证明子进程禁写未脱敏日志，也没有默认禁用 BYOK。
- **租约**：接管成功后，持旧 token 不能在同代夺回；锁内先比较当前 token。未见生产回调反向获取 store 锁；阻塞风险来自上述锁内读取。每次换代保留一份文件，目前无上限或清理（lease.go:304）。
- **兼容性**：`FileID` 的 Windows 句柄查询和 Linux `fstat` 与两条 CI 腿的构建标签匹配。无 glob 的重试不受验根影响；有 glob 时普通文件增删不影响根身份，根在重试时已删除或重建则拒绝。仅传 `--ledger` 的 approvals/cost 保持原义；已有的 `approvals revoke --run-dir` 未显式传 `--ledger` 时会改写目标位置，从 CWD 改为 run 目录，应在报告注明。
- **证据口径**：只读审查了暂存区 75 个变更文件，`git diff --cached --check` 通过。未重跑 Go 测试、lint 或双平台 CI；表中反例是静态核对，不是本轮运行日志。⑤/⑥ 的 PASS 仅绑定整改前快照。本轮未改文件、未提交。

---

## RE-REVIEW 3（额度 3/3，最后一轮）

**RE-REVIEW: FINDINGS**

| 项目 | 状态 | 守卫位置 | 会红的反例 |
|---|---|---|---|
| T3 启动绑定 | **部分关闭**。共用守卫从同一打开句柄计算摘要；有摘要时拒绝裸命令，Darwin 拒绝启动。Windows 未固定启动路径，Linux 未固定已验字节。 | executable_binding.go:63、process_windows.go:306、process_linux.go:68、process_unsupported.go:12 | T3a 直接写入／改名（:44）、T3b 裸命令（:90）、T3c 摘要不符（:132）、T3d Linux 同名替换（:45）。未覆盖 Windows 路径重定向或 Linux 同 inode 内容变更；T3d 允许所有启动均被拒绝，缺少成功对照。 |
| T4 产物枚举 | **部分关闭**。已移除 `workspaceArtifactRel`；Linux 从同一根 FD 解析目录和产物，越界模式被拒。枚举仍有下述平台缺口。 | runner.go:365、bounddir_unix.go:278、bounddir_windows.go:174 | T4a 根改名后仍枚举和收据绑定原对象（:60）、T4b 拒绝越界（:140）。未覆盖 `DT_UNKNOWN` 或 Windows 后代目录解析竞态。 |
| T5#3 保留代读取 | **部分关闭**。Unix 拒绝跟随末级链接并避免 FIFO 阻塞；各平台都检查打开对象的类型并限读 64 KiB。Windows 拒链不原子，限字节也不等于限时。 | lease_open_unix.go:18、lease_open_windows.go:17、lease.go:433 | T5c1 静态符号链接（:61）、T5c2 超限（:106）、T5c3 非普通文件（:139）。没有打开竞态或锁内超时反例。 |
| T5#4 陈旧 proof | **部分关闭**。批量读取的解析／I/O 错误会返回；同进程可追平，异进程持有的陈旧 proof 仍被拒。跨进程恢复未闭合。 | runtime_supervised.go:1186、runtime_agent.go:2087 | T5d 同进程追平（:28）、T5d 异进程拒绝（:71）、T5e 坏 proof 使批次失败（:113）。未覆盖代次提交后、proof 落盘前退出并重入。 |
| T8b 最终拒绝证据 | **部分关闭**。已补最终结果反例，且不再回退到调用参数摘要；但最终构造器仍接收未经分类约束的证据。 | ai_probe_copilot.go:116、:227、ai_probe_copilot_executor.go:332 | T8b 两组配置的最终拒绝证据一致（:74）会抓到调用摘要回流；替身只返回空 `Evidence` 和空 `FailureCode`，不约束非空输入。 |
| T1 准入取钟 | **已关闭**。两个已配置的账本来源均先读，之后统一取钟。 | agent_runner_admission.go:118、:132 | T1a 授权账本读取期间到期（:33）。取钟移到授权读取前会红；移到两次读取之间尚无独立反例。 |

**阻断项与最小修复**
- **High，T3 Windows：锁住文件不等于锁住文件名。** `FILE_SHARE_READ` 阻止普通写打开和需要 `DELETE` 的直接改名／删除；T3a 实测设计也只覆盖直接操作原文件。它不固定符号链接、junction 或可变祖先，`CreateProcessW` 仍重新解析原路径。Windows 另有允许替换已打开目标的 POSIX rename 语义，不能从 T3a 推出所有目标替换都被拒。最小修复：只允许从执行者不可改写、不可重定向的稳定路径启动，并保持文件句柄到创建完成；无法证明则拒绝。补真实启动反例。
- **High，T3 Linux：FD 固定 inode，不固定字节。** `/proc/self/fd` 能越过原路径的改名或删除，且真实 procfs 不可用时没有路径回退。普通 shebang 可借继承 FD 让解释器重开脚本，但脚本看到的是描述符路径，解释器也未被该摘要绑定。同 inode 在哈希后仍可被改写。最小修复：核准并执行同一份不可变对象，或只接受能证明不可写的安装来源；补原地变更与成功启动对照。
- **High，新增 Linux 回归：无摘要启动也覆盖 PATH 解析结果。** process_linux.go:97 无条件改写 `cmd.Path`；`Start` 不再查 PATH，裸命令可能失败或启动工作目录里的同名文件。此外，裸 `Dup`（:87）创建的可继承 FD 也可能被无关的并发子进程继承。最小修复：仅批准分支设置描述符路径，保留普通 `exec.Command`；用 `ExtraFiles` 只向目标子进程传递 FD。补普通 PATH 启动和并发继承回归。
- **Medium，T4：枚举还不是全程句柄绑定。** Go 1.22.12 的 `ReadDir` 遇到 `DT_UNKNOWN` 会按 `os.NewFile` 的名字做 `lstat`；这里的名字是 `.` 或相对目录，可能按进程工作目录漏掉真实条目。Windows 的根目录校验文件 ID，但后代和产物只比较两次取得的路径字符串，不能排除中途换根。最小修复：Linux 用 `Readdirnames`；Windows 从已绑定句柄相对打开后代，不能证明归属则拒绝。
- **Medium，T5#3：Windows `Lstat` 和 `Open` 之间仍可换对象。** 打开后 `Stat` 只证明目标是普通文件，不证明未经过链接；`LimitReader` 也不限制同步打开或读取的等待时间。最小修复：从受信根相对打开、原子拒绝重解析点，并在同一磁盘句柄上检查类型；对不具备可控 I/O 的存储拒绝，或提供真正可取消的读取。
- **Medium，T5#4：提交和批写之间的跨进程中断仍会卡住。** 新进程在进入旧持有者的停机证明前就被陈旧 proof 拒绝；反过来，同 PID 分支会把任意合形但不匹配的 token/hash 当成陈旧。最小修复：让代次到 proof 的更新可持久恢复，并校验旧 proof 属于同一 claim 的真实历史。普通异主 proof 继续拒绝，不放宽接管条件。
- **Medium，T8b：最终出口只校验摘要形状，不校验证据来源。** `start-failed` 仍可带入执行器提供的摘要和自由文本 `FailureCode`；生产执行器也会在标记启动失败后追加包含观察值的摘要。最小修复：最终拒绝证据只从封闭分类生成，不合并任意执行器摘要；计量值可独立保留。补非空证据的最终出口反例。

**切片报告残余（可直接登记）**
- **T2 STOP：** AgentRunner 未持久化独立获批的冻结策略对象，无法将 grant 的 `PolicyHash` 与该对象复核；不得记为闭合。任何后续接纳均须用户明确裁定该 STOP 的处置。
- **T8a 延期：** P83 未实施「无法证明子进程禁写未脱敏日志时拒绝 BYOK 探针」；按既有裁定顺延 P84，不宣称子进程不会落盘。
- **租约保留：** 每次续租保留旧代记录，未设清理上限；后续保留策略不得破坏已签发回执的独立复核。
- **平台验证：** 本 Windows 主机不运行 Linux-tagged T3d，且无所需权限运行 T5c1 符号链接用例；分别在 Linux 和具备权限的 Windows 环境补证，不将未运行或跳过记为通过。

**接纳：不通过，3/3 用尽，必须升级。** 剩余项为 T3、T4、T5#3、T5#4、T8b 和新增 Linux 启动回归。T2 用户裁定是后续接纳的必要条件，不是本轮放行依据。

**核验范围：** 只读核对了所列文件的工作树源码和测试。当前审查模式未能直接取得本轮 `git diff --cached`，未独立确认索引一致，也未重跑测试、lint 或变异。未改文件、未提交。
