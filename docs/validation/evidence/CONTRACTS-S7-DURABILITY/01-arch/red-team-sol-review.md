# ① 双打第二棒：GPT-6.1 Sol 红队抗压审查原文（P54 CONTRACTS-S7-DURABILITY）

载体：`independent-reviewer` + `GPT-6.1 Sol (copilot)`（§1.7.4 临时授权，白名单外，限本片一次性，不得作为先例；2026-11-04，只读）。以下为红队返回文本逐字归档。

**逐条发现**
1. **契约等级边界：无。** [design.md:13](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L13) 明写：“三判据同时满足仅构成该（平台 × 前提）下 C6 的候选级结论，不自动改变平台发布耐久等级”。与 [CONTRACTS.md:232](docs/CONTRACTS.md#L232) 的部署前提覆盖、解锁边界及 [Windows 分支:5](internal/adapters/agent_runner_replay_store_windows.go#L5) 一致。未发现变相降低平台级判定门槛。

2. **High｜异常路径存在三处未闭合。**
   **空槽与损坏槽混淆：** [design.md:24](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L24)：“payloadLen=0 ⇒ 读者见 absent”；但 [design.md:20](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L20) 定义：“空槽 = 长 16 且 payloadLen=0”。写 payload 后即使 flush #1 失败，文件也可能已经增长。此时按零长度直接判 absent，会绕过异常长度检查；先检查长度则会把设计要求不可见的未提交槽判为 corruption。较长失败写后重试较短载荷，还可能遗留尾部。**最小修复：**冻结未提交中间态、长度校验顺序及失败重试规则；不得仅凭 `payloadLen=0` 判 absent，也不得把“未提交”表述为“仍是完整空槽”。

   **丢槽后重建承诺无法执行：** 同一行原文：“已发布槽目录项丢失 → 不静默重建（不可与「从未创建」区分），由 C 的 R+C 绑定/孤儿检测兜底。”但 [design.md:22](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L22) 的入口首先调用 `createReplaySlotWindows`。R-only 的 R 丢失时，C 也不存在，[现有 State:294](internal/adapters/agent_runner_replay_store.go#L294) 会判 absent，孤儿检测无法兜底；按设计重试会重新建槽。**最小修复：**定义可追溯的预建/分配身份，使“应存在槽缺失”可被拒绝；未解决前删除“绝不静默重建已获保证”的论断。

   **可见性与 flush 阶段混同：** [design.md:55](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L55)：“R 允许集：头提交前恒 absent、后恒 dispatched-unknown”。反例是：数据 flush 成功、头写成功、flush #2 前进程被 kill。进程级 kill 不清除 OS 页缓存，锁最终释放后读者可以看到非零头，即便 `slot-header-committed` 阶段尚未触发。**最小修复：**分开定义头写可见点、第二次 flush 耐久点及解锁点；补该窗口及 flush #2 失败的注入，修正候选允许集，不改三判据本体。

3. **High｜共享接缝的替换范围与读者覆盖不匹配；Windows API 泄漏本身无。**
   [design.md:28](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L28)：“replaces the whole no-replace publication”；[design.md:40](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L40) 却只替换 `replayStoreLoadRequestRecord` / `replayStoreLoadCompletionRecord`。现有发布接缝覆盖 ownership、launch、terminal 等记录；[ownership 发布:425](internal/adapters/agent_runner_replay_root.go#L425) 同样经过该函数，而 [ownership 解码:481](internal/adapters/agent_runner_replay_root.go#L481) 仍直接严格解码 JSON。全局写入二进制槽头后，其重开或碰撞读取会失败。
   **最小修复：**明确接管记录域，未迁移记录保持原发布与读取路径；不能只迁移 R/C 读者却替换整个发布器。槽容器格式也须显式协调 [CONTRACTS.md:230](docs/CONTRACTS.md#L230) 的“一条 canonical JSONL”要求，不能以“路径不变”代替格式兼容。

4. **Medium｜核心声明的充分性论证仍有缺项。**
   [design.md:49](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L49)：“发布时点不产生任何目录变更 ⇒ 发布时耐久事实只剩文件数据”。不创建新目录名，并不等于只剩数据：当前预建文件仅 16B，发布还会改变文件长度及可能的分配元数据；小头写与两次 flush 也不自动证明提交更新防撕裂。此外，[design.md:13](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L13) 的“槽文件……永不……追加”与该增长路径按字面冲突。
   **最小修复：**改为有条件声明：“预建目录项及所需元数据耐久已独立建立后，发布不再创建新名字；内容、提交头及必要元数据的耐久仍须取证。”采用真实定长预分配，或明确把增长元数据纳入取证；同步澄清“不得追加”的适用边界。这不是证伪 C6，而是否定当前论证已经充分。

   **Medium｜区分力测量对象未定义。** [design.md:56](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L56)：“预期 slot-created 格 ≥1 丢失”；但该点正负臂的逻辑 R 都是 absent。**最小修复：**分别定义物理槽存活、未提交载荷存活与逻辑记录可见性的测量，避免用同一个 absent 结果判槽丢失或存活。

5. **Medium｜引用纪律未完整传递到输入包。**
   [input-list.md:14](docs/validation/evidence/CONTRACTS-S7-DURABILITY/input-list.md#L14) 原文含“B3b 结论账本（C2 `proven`，决定格……”；[input-list.md:27](docs/validation/evidence/CONTRACTS-S7-DURABILITY/input-list.md#L27) 又写“B3b 的 C2 = MoveFileEx（已有 `proven` 证据）”。这些是缺陷原文引用，**不是本审查的无前提等级断言**。[判定表:54](docs/validation/contracts-s7-durability.md#L54) 的摘要引用也未紧随完整括注。
   [表下注记:60](docs/validation/contracts-s7-durability.md#L60) 确实存在，但不能替代随引用携带限定。**最小修复：**保留历史原行、逐字归档及候选级结论；本片新摘要、摘录和提示词逐处补“（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）”。当前红队请求中的禁止性引用不算等级断言。

6. **High｜“首次发布前预建”不是已证明等价的 A7 解析。**
   [design.md:17](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L17)：“本设计解析为「首次发布前预建 + store 打开时核验既有槽」”；原 [A7:28](docs/validation/t027-windows-publication-durability.md#L28) 是“bootstrap 预建固定槽”。本设计在每个新 requestId 的发布入口创建目录项；每槽只提交一次，创建耐久仍在首次发布临界路径上，没有建立原设想的预建结构优势。“requestId 未知”也不足以排除预留槽池或计划冻结后预建。
   **最小修复：**先明确 bootstrap 完成点、槽容量、身份绑定及缺失槽识别；按需建槽变体须显式报 Sol/用户裁定，不能直接作为原 A7 的等价措辞写入 A2。

7. **Medium｜LockFileEx 原语可用，但当前仲裁协议不足以准入。**
   [design.md:22](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L22) 仅写“`LockFileEx` 排他锁”。[官方契约](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-lockfileex) 确认：进程退出或关闭句柄由 OS 解锁，但释放可能延迟；普通文件 I/O 的锁定区域是强制性排他，不是仅 advisory；映射视图不受其保护。`FILE_SHARE_READ|FILE_SHARE_WRITE` 不会取消字节锁限制。
   设计未给锁范围、读者共享锁、等待上限和完整锁序。读者直接读锁定头会失败；同步等待还可能持有 `store.mu` 无限阻塞。现有“不得回调 store”是必要约束，但不足以独立排除全部死锁。
   **最小修复：**冻结统一仲裁范围、同句柄读写、读者共享锁、精确解锁及失败关闭；固定 `store.mu → 文件锁`，禁止反向取锁和未定义的多槽嵌套；取锁等待有界，失败按 uncertain 收敛；明确映射 I/O 不属于保护范围。

8. **nil 守卫措辞：无，限可观察行为等价。**
   [design.md:30](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L30)：“while nil, the logic path and observable behaviour”；随后写“stay equivalent to today's”。该表述可以接受：新增 nil 判断不进入槽实现，原编码、发布原语、父目录同步、阶段及返回路径可保持等价。
   这不意味着源码、编译产物、指令或调度时序逐字节相同。[design.md:44](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md#L44) 的“字节级行为”只能解释为同一输入下的落盘字节等可观察结果。非 Windows 激活赋值必须继续由构建标签排除；实现后的验证尚待执行。

**三项报 Sol 事项裁定**
① **需修订。**“首次发布前预建”改变了预建生命周期和耐久依赖位置，不能作为原 A7 的直接等价解析。

② **需修订。**增加通用函数接缝本身可接受，nil 默认也可保持 Unix 行为；但当前“替换整个发布器、只替换 R/C 读者”的契约不完整，须先明确接管范围及容器格式。

③ **需修订。**`LockFileEx` 可作为条件候选；当前缺少读者协议、锁范围、有界等待及完整锁序，尚不能批准现稿准入。

本轮已实际读取指定输入及必要邻近调用路径；无文件修改，未生成 diff，未运行 build/test/lint 或故障注入。上述反例是静态推导，不声称当前生产已发生错误派发。

RED-TEAM: FINDINGS — Windows 保持 `unproven` 的等级边界正确，但 C6 的预建生命周期、异常恢复、提交可见性及读写仲裁尚未形成自洽的可判契约，需交回①修订。
