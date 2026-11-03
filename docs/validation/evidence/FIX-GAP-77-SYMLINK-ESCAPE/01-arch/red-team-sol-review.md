只读核对已完成，报告将列出两项 High 和四项 Medium，并逐项覆盖八类载荷。核心区别是：该设计可阻断逐组件软链接解析，却不能据此证明子目录移出、挂载跨越和并发 Close 下都安全；这些新反例本环未实跑，会明确标注为推演。


## Section A 范围

日期：2026-11-03。P53 ①架构双打第二棒，只读安全审查；本次临时授权不延续至其他切片。设计全文已阅读，未修改任何文件，未执行测试、挂载、远程命令或 Git 写操作。

输入文件：
- 待审：[design.md](docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/design.md#L1)，完整 §1–9。
- 实现与接口：[bounddir_unix.go](internal/adapters/bounddir_unix.go#L1)、[bounddir.go](internal/adapters/bounddir.go#L1)、[bounddir_clean.go](internal/adapters/bounddir_clean.go#L1)、[bounddir_windows.go](internal/adapters/bounddir_windows.go#L1)、[boundread.go](internal/evidence/boundread.go#L1)。
- 测试与调用方：[bounddir_escape_unix_test.go](internal/adapters/bounddir_escape_unix_test.go#L1)、[evidence_bind_dir_unix_test.go](internal/adapters/evidence_bind_dir_unix_test.go#L1)、[bounddir_identity_test.go](internal/adapters/bounddir_identity_test.go#L1)、[agent_runner_replay_store.go](internal/adapters/agent_runner_replay_store.go#L220)。
- 证据及规范：[P52 Linux 击穿证据](docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/05-attack/t51-1-linux-escape.md#L1)、[CONTRACTS](docs/CONTRACTS.md#L234)、[SECURITY](docs/SECURITY.md#L1)、[TEST_STRATEGY](docs/TEST_STRATEGY.md#L1)、[CHECKLIST v42](docs/t027/CHECKLIST_v42.md#L65)、[DELIVERY_DIRECTIVE](docs/DELIVERY_DIRECTIVE.md#L495)、[go.mod](go.mod#L1)。

## Section B 发现

1. **High：钉住子目录对象，不等于钉住其祖先关系。** 位置：设计 §1.1 不变式、§1.2 rename 行、§4.2 步骤 3–4。精确构造：同一文件系统内，攻击者可写 `D`、`O`；①游走打开 `D/sub` 得到 fd `n`；②暂停游走，把 `D/sub` 改名为 `O/sub`，或与 `O/spare` 做 `RENAME_EXCHANGE`；③将 `O/secret` 改名为 `O/sub/secret`；④恢复，`openat(n,"secret",...)` 成功。这个叶子从未位于绑定目录 `D` 内，目录 `fstat` 仍通过。A/B 均受影响，无需挂载权限。
   **修改指令：** 打回 V4 Pro，删除“并发 rename 免疫”的无条件结论，明确包含关系的判定时点。若维持 §1.1 强不变式，评估从原始 `dir.fd` 一次执行受限 `openat2`，并定义不支持环境的保守行为；不能把当前逐级 A 当作等价保证。若要求打开后持续保持祖先关系，还须约束目录移动。不得用再次路径检查或字符串前缀比较冒充消除竞态。绑定根整体改名后的既有可读性应继续保留。

2. **High：未定义挂载信任边界，却宣称全部载荷 fail-closed。** 位置：设计 §1.1、§1.2 结论、§2 安全性、§5。精确构造：绑定 `D` 后，在读取者挂载命名空间执行 `mount --bind O D/m`；读取 `m/secret` 将取得 `O/secret`。中间组件仍是目录，`O_NOFOLLOW` 不限制挂载跨越；叶子文件 bind mount 同样可放行。需要挂载能力或既有挂载，不代表普通攻击者能修改另一进程的命名空间。
   **修改指令：** 先由 Sol 冻结“命名空间中可达即界内”还是“禁止外部来源挂载”的策略。后者应评估 `RESOLVE_NO_XDEV` 及同保证降级；`RESOLVE_BENEATH` 单独不禁止 bind mount，比较 `st_dev` 也挡不住同文件系统绑定。若挂载属于可信 OS 前提，须明确排除，不能继续写成已阻断。

3. **Medium：A/B 降级描述与可执行流程不一致。** 位置：设计 §2、§3、§4.2、§5、§9。§4.2 中间组件无条件使用 `O_PATH`，两个叶子分支完全相同，没有实现 A 的目录打开方式；实际也没有任何 `/proc/self/fd/N` 重开。`O_PATH` 目录 fd 可以直接用于 `openat`，不需要 `/proc`。探测成功只证明打开了某个目录，不证明其为真实 procfs，更不证明未来魔法链接可用；探测失败也可能只是 `EACCES/EMFILE`。
   **修改指令：** 统一算法和威胁论证。优先评估完全无 `/proc` 的实现，或明确真正的 A/B 中间组件分支；不得保留伪分支。若坚持魔法链接重开，必须补真实 procfs 验证、持有句柄、命名空间变化及 fd 生命周期论证。将 `O_PATH|O_NOFOLLOW ⇒ ELOOP` 改正为“可取得链接自身 fd，必须由 `fstat` 拒绝”；“无权限检查、无引用计数开销”也应删除。

4. **Medium：最低工具链无法照抄伪代码构建。** 位置：设计 §4.2 的两处 `syscall.O_PATH`。Go 1.22.12 Linux/amd64 的标准库未导出该常量，照抄会产生未定义符号；本环依据对应版本源码核对，未运行编译。
   **修改指令：** 明确经过批准的 Linux/ABI 适配方式，或采用不依赖该常量的 A；不要未经授权增加依赖或扩大平台改动。必须用 `go1.22.12 linux/amd64` 构建核验，不能以本机较新工具链替代。

5. **Medium：共享 fd 表下的生命周期前提缺失。** 位置：设计 §4.2 借用 `dir.fd`、关闭逻辑；现有 `BoundDir.Close` 与调用方取出目录后解锁的路径。精确时序：游走复制 fd `n` → 并发 `Close` 释放它 → 另一线程打开 `O` 恰好复用 `n` → `openat(n,"secret",...)` 读取 `O/secret`。普通线程的 `open` 不能复用仍被持有的 fd；问题前提是并发关闭，不是外部进程能随意关闭本进程 fd。此风险原已存在，设计未覆盖。
   **修改指令：** 明确并证明 Close/读取互斥，或在 Unix 私有实现中同步生命周期；需要复制 fd 时，复制本身也必须在保护下完成。中间 fd 严格单一所有者、只关闭一次。若并发 Close 明确不受支持，应记录调用方约束及其证据，不能宣称任意并发安全。

6. **Medium：反例可以因错误原因假绿。** 位置：设计 §6.1–6.3；T51-1 丢弃返回句柄、只检查 `err != nil`。红判据“目标测试 FAIL 且含专属文案”正确，环境 `EACCES` 不应冒充红；但绿侧及 P52 探针把任意打开错误都算阻断，`EBADF/EMFILE/EACCES` 等故障也可能令反例通过。原测试本身没有读取字节。
   **修改指令：** 固定夹具、确认界外 canary 可直读，同一后端正向读取成功；拒绝错误必须符合冻结的安全拒绝类别，而非任意错误。红证据记录基线、生产代码摘要、测试摘要、工具链、退出码，并机械确认目标测试实际运行且未 SKIP；绿必须运行同份测试。保留字节读取探针，并增加确定性 rename、后端分支及资源故障区分测试。不能仅以包级 PASS 或“未观察到泄漏”结算。

## Section C 证据

以下是按设计流程的推演，不是本环新取得的运行结果。“逃逸”注明其采用的边界口径。

| 对抗载荷 | 判定 | 依据 |
|---|---|---|
| `a→b→c`、`a→O` 后访问 `a/b` | 阻断 | A 拒绝非目录链接；B 必须以 `fstat` 拒绝链接自身，不能依赖 ELOOP |
| `link/../safe`、`sub/../link/secret` | 阻断越界；可能接受净化后的界内名 | 前者折叠为 `safe`，根本不解析 `link`；后者折叠为 `link/secret` 后拒绝 |
| `a/../../b`、重复 `/`、混入 `./` | 阻断越界 | 净化后前导 `..` 被拒，其他形态按净化结果游走 |
| 目录硬链接；叶子硬链接 | 目录阻断；叶子来源隔离设计未覆盖 | Linux 普通目录硬链接返回 EPERM；合法叶子硬链接可引用同时具有界外名字的同一 inode，不能区分“原名” |
| 尚未打开的中间名换入软链接 | 阻断 | 该级原子 open 加类型检查；已打开父目录不因名字替换而重定向 |
| 已打开子目录移到界外，随后放入新叶子 | 逃逸 | R1；fd 稳定性恰使后续读取继续访问已移出的对象 |
| `RENAME_EXCHANGE` 换入普通目录树 | 设计未覆盖身份策略 | 打开前换入的普通树现已位于 D 内，不是路径逃逸；但没有“仍为原目录树”的身份保证 |
| 自环、双环、长链接链 | 阻断 | 第一处实际遇到的链接即拒绝，不展开链接链 |
| 64/65 个合法组件、可折叠超长拼写 | 超限拒绝；合法路径可能误杀 | 帽统计净化后组件且包含叶子；折叠减少组件不是越界绕过，但未限制原始输入资源量 |
| `/proc` 遮蔽、chroot、探测后命名空间变化 | 设计未覆盖真实魔链分支 | 按 §4.2 没有魔链重开，变化不构成该路径的魔链逃逸；若照 §2/§5 实现则缺保证 |
| 把 `/proc/.../fd/N` 魔链作为输入组件 | 阻断 | 中间链接由类型检查拒绝，叶子由 `O_NOFOLLOW` 拒绝 |
| 多线程 fd 复用 | 设计未覆盖；并发关闭时可逃逸 | R5；独占持有的活 fd 不会被普通 open 复用 |
| bind mount 在打开前出现于中间或叶子 | 逃逸，按禁止外部来源口径 | R2；`O_NOFOLLOW` 对挂载点不免疫 |
| 拼写守卫与内核衔接 | 条件自洽 | 只能把 `clean` 交给游走；NUL 等可留给系统调用拒绝。守卫与内核不必返回相同错误；尾斜杠折叠是既有净化语义 |

## Section D 可证伪性审计

- **已有实跑证据：** P52 报告记录界外 23 字节读取；它不是本片当前基线的红测替代物。本环没有新增测试/lint 日志；diff 不适用，因为零修改。
- **可机械复现：** R1 用测试专用屏障停在子目录打开后，再移动目录并安放 canary；R5 用生命周期屏障及确认 fd 复用；挂载用隔离子进程/命名空间。均须确认攻击阶段命中，不能以压力循环未命中判通过。
- **可机械检查：** 最低工具链编译、逐分支覆盖、真实无 `/proc` 环境、63/64/65 组件边界、目标测试 RUN/FAIL/PASS/SKIP 计数及错误类别。
- **推理边界：** fd 稳定、链接自身 fd、挂载跨越及硬链接别名语义依据 [open(2)](https://man7.org/linux/man-pages/man2/open.2.html)、[openat2(2)](https://man7.org/linux/man-pages/man2/openat2.2.html)、[rename(2)](https://man7.org/linux/man-pages/man2/rename.2.html)、[link(2)](https://man7.org/linux/man-pages/man2/link.2.html)；工具链依据 [Go 1.22.12 常量表](https://raw.githubusercontent.com/golang/go/go1.22.12/src/syscall/zerrors_linux_amd64.go)。

## Section E 无法验证部分与推翻条件

**三条红线：** 公共签名不变，设计态满足；红测先行，已设计但尚无本片执行证明，须修正 R6；无 `/proc` 降级，仅文字承诺，§4.2 未落实，当前不能验收。

**§8-1 裁定建议：** SEC-02/AT-02 更新是 CHECKLIST v42 已明确的同片交付物，连同英文镜像由 Sol 纳调用侧锁只保护目录引用的取得，没有覆盖实际打开过程，因此设计还需明确并发 `Close` 的约束。审查材料已齐全；我没有修改文件或运行会创建夹具的测试，下面给出含逐项判定、复现步骤和打回指令的完整报告。


## Section A 范围

P53 `FIX-GAP-77-SYMLINK-ESCAPE`，①架构双打第二棒，2026-11-03 临时授权，仅本片有效。本轮只读，未修改任何文件，未执行会创建夹具的测试、挂载或远程命令。

完整审阅：[design.md](docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/design.md#L1) §1–9。实现输入：[Unix](internal/adapters/bounddir_unix.go#L1)、[冻结公共面](internal/adapters/bounddir.go#L1)、[Windows 对向](internal/adapters/bounddir_windows.go#L1)、[拼写守卫](internal/adapters/bounddir_clean.go#L20)、[BoundFile](internal/evidence/boundread.go#L88)、[调用与生命周期](internal/adapters/agent_runner_replay_store.go#L220)。

测试及规范输入：[T51-1](internal/adapters/bounddir_escape_unix_test.go#L38)、[正向与叶子测试](internal/adapters/evidence_bind_dir_unix_test.go#L36)、[同一性测试](internal/adapters/bounddir_identity_test.go#L20)、[P52 实跑证据](docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/05-attack/t51-1-linux-escape.md#L1)、[契约](docs/CONTRACTS.md#L236)、[安全](docs/SECURITY.md#L9)、[测试策略](docs/TEST_STRATEGY.md#L73)、[P53 交付要求](docs/t027/CHECKLIST_v42.md#L65)、[执行准则](docs/DELIVERY_DIRECTIVE.md#L489)、[最低 Go 版本](go.mod#L3)。

## Section B 发现

共 **2 High、4 Medium**。以下节号与步骤均指待审设计；发现针对设计承诺，不表示本轮已运行修复代码。

**F1 High：fd 钉住对象，不钉住对象对绑定根的从属关系。位置：§1.1 不变式、§1.2 rename 行、§4.2 第 3–4 步。** 精确构造：①同一文件系统建立相邻目录 `D/sub`、`O`，绑定 `D`；②读取线程打开 `D/sub` 得到 `cur` 后暂停；③攻击者执行 `rename(D/sub,O/sub)`，再把原本仅在 `O` 中的常规文件 `secret` 移入 `O/sub/secret`；④恢复 `openat(cur,"secret",...)`，读得界外文件。只需相应目录写权限，不需挂载权限；`RENAME_EXCHANGE` 可原子交换目录位置，再注入文件。**修复指令：**打回 V4 Pro，区分“绑定根整体改名仍可读”与“中间目录移出后继续读取”。明确包含关系的判定时刻。若保留强边界，评估从原始绑定根单次 `openat2(RESOLVE_BENEATH|RESOLVE_NO_SYMLINKS)`；无法提供等价保证的降级不得宣称免疫，须拒绝或依赖可证明的写者冻结。返回 fd 后永久保持当前从属关系，也不能靠一次打开保证。

**F2 High：mount point 不是 symlink，A/B 均会穿过。位置：§1.1、§1.2 总结、§2 安全性、§5。** 精确构造：①建立 `D/m` 与含 canary 的 `O/secret`；②在读取方可见的挂载命名空间中 `mount --bind O D/m`；③读取 `m/secret`。中间对象是目录，`fstat` 通过，`O_NOFOLLOW` 不阻止挂载遍历。已有 bind mount 不要求攻击者临时取得 CAP；动态攻击则需要相应命名空间权限。**修复指令：**由 Sol 裁定边界究竟允许挂载导入，还是要求排除界外来源。后者应评估 `RESOLVE_NO_XDEV`，覆盖目录及叶子文件 bind mount；同文件系统 bind mount 不能仅用 `st_dev` 排除。不得仅凭“低权限不能 mount”删除此威胁。

**F3 Medium：A/B 定义互相矛盾，降级尚不可执行。位置：§2 性能、§3、§4.2、§5。** §2/§5 声称经 `/proc/self/fd/N` 打开，§4.2 却明确不使用魔法链接；中间步骤又无条件使用 `O_PATH`，两个叶子分支完全相同，未写出 A 的目录打开分支。另外，Go 1.22.12 Linux/amd64 的 `syscall` 未定义 `O_PATH`，照抄不能编译。**修复指令：**统一实际算法，明确中间 flags 的唯一分支及常量来源；非 Linux 构建不得无条件引入 Linux 专属能力。`O_PATH` 目录 fd 可直接供 `openat` 使用，本身不依赖 `/proc`。同时纠正两处论据：`O_PATH|O_NOFOLLOW` 返回软链接自身 fd，须由类型检查拒绝，并非必然 ELOOP；它仍持有内核引用，且不免除后续目录搜索权限检查。

**F4 Medium：借用根 fd 的生命周期没有并发闭合。位置：§4.2 第 2 步及“Close 不变”。** 精确时序：读取线程保存 `cur=N`；另一线程 `Close` 关闭 N；同进程打开 `O` 恰好复用 N；读取线程执行 `openat(N,"secret",...)`。普通并发 `open` 不能复用仍存活的 fd，危险来自并发关闭或重复关闭。调用侧锁只保护取得目录引用，没有覆盖实际打开。**修复指令：**先交 Sol 明确并发 `Close` 是否允许；允许则用私有同步覆盖检查、游走和关闭，保持签名不变；禁止则补可核验的调用方生命周期前提。仅先 `dup` 而不保护复制窗口也不成立。

**F5 Medium：红判据较明确，绿判据仍可假绿。位置：§6.1–6.3；T51-1 的打开断言。** 指定 FAIL 文案能排除普通 EACCES、编译失败等“误红”，但测试没有读取并验证 canary；任意打开错误又都会变绿。例如只对多组件路径返回 ENOENT，T51-1 即通过；无差别拒绝多层路径虽可被正向用例抓住，选择性错误拒绝仍未排除。**修复指令：**基线红必须确认真实句柄读得精确 canary，再以唯一标记失败；绿必须检查无返回句柄、拒绝属于预定安全错误，并执行同条件正向读取与摘要对照。归档真实退出码、测试名及非 SKIP 事件。移除逐级防护的变异必须重现红；真实绑定对象上的强制 A 测试须证明实际选择 A，不能仅证明字段为 false。

**F6 Medium：64 是新增合法路径限制，不是软链接环防线。位置：§4.2 第 1 步、§4 实现要点、§8-4。** 64 个短目录加叶子共 65 组件，完全可能低于 Linux 路径长度限制，旧实现可读，新设计直接拒绝。“常见证据深度≤3”不能证明公共入口只允许该深度。**修复指令：**由 Sol 裁定是否授权新增上限；获准后明确按净化后的组件数、包含叶子计数，并测试 63/64/65 边界。否则不要让实施者临时选择数值。

## Section C 证据

判定针对设计所述算法；“阻断”表示不能借该载荷越界，不保证原始拼写一定被拒绝。

| 对抗载荷 | 判定 | 依据 |
|---|---|---|
| `a→b→c`、`a→O; a/secret` | 阻断 | A 拒非目录；B 必须以 `fstat` 拒绝链接自身 fd。 |
| `link/../secret` | 阻断越界 | Clean 后为 `secret`，链接从未被内核解析；可能合法读取界内同名文件。 |
| `sub/../link/secret`、`a/../../b` | 阻断 | 前者剩余链接由游走拒绝；后者净化为上逃路径，由守卫拒绝。 |
| 内部重复 `/`、混入 `./` | 阻断越界 | 净化后逐组件打开；绝对路径仍先拒绝。 |
| 目录硬链接 | 阻断 | 普通 Linux VFS 的 `link` 拒目录，EPERM；并非获得常规 CAP 即可开启。 |
| `link(O/secret,D/file)` | 设计未覆盖来源隔离 | 有权限且同挂载时可成功；界内目录项与界外名字共享 inode。不是路径逃逸，却不能证明内容来源仅在界内。 |
| 打开后替换叶子硬链接 | 阻断句柄重定向 | 已打开 fd 不变；共享 inode 的原地写入不因此被冻结。 |
| 中间名字换成软链接 | 阻断重定向 | 换在打开前则拒绝；换在打开后不改变已取得的目录 fd。 |
| 中间目录移出，再注入叶子；目录 EXCHANGE | 逃逸 | F1；对象身份稳定不等于仍属于绑定根。 |
| 换入另一普通目录树 | 设计未覆盖原身份要求 | 打开前换入会读取新树；若它当时在 D 下，不构成当前路径越界，但不保留旧树身份。 |
| 软链接自环、互环、长链 | 阻断 | 首个链接即拒绝，无需等待 64 层；A 的 errno 不应被强钉为 ELOOP。 |
| 深路径与大量可折叠组件 | 阻断超限／可能误杀 | 65 个净化组件拒绝；折叠后变浅不是解析深度绕过。原始输入处理成本另计。 |
| 输入路径中的 proc 魔法链接，包括叶子 | 阻断 | B 中间类型检查、叶子普通 `O_NOFOLLOW` 均拒链接；无需为叶子增加魔链重开。 |
| proc 遮蔽、chroot、探测后 namespace 变化 | 设计未覆盖 | §4 实际不使用魔链，环境变化不产生该重开攻击；若照 §2/§5实现，缓存 bool 不证明后续路径仍是真 proc。 |
| 共享 fd 表、fd 数字复用 | 有条件逃逸 | F4；必须先发生关闭。不同 OS 线程的 namespace 差异也不能由一次探测排除。 |
| 绑定目录内目录／叶子 bind mount | 逃逸于对象来源边界 | F2；若定义为允许的命名空间内容，则须显式改变承诺，不能仍宣称排除此类来源。 |
| 守卫与内核拼写语义 | 有条件自洽 | 只能统一解释为“读取 Clean 后的相对名”；不能同时承诺原始 `link/..`、尾斜杠等内核语义。 |

**三条红线：**签名不变：满足设计意图，私有字段不改变签名；红测先行：时序满足，但须按 F5 加固证据；无 `/proc` 降级：文字承诺存在，操作设计未闭合，不能验收为满足。

**§8 裁定建议：**
- 第 1 条：安全与测试策略及其英文镜像应为本片伴随规范交付，由 Sol 明确纳入派单；已由 CHECKLIST 要求，不应另片拖欠。建议②先落实规范性验收条件，⑧汇总实证。若确有“仅代码及测试”的授权冲突，先请用户裁定，不能自行扩范围或省略。
- 第 2 条：赞成包内注入，不增加公共接口；使用真实已绑定对象、局部可控选择，避免全局可变钩子与并行测试串扰。还需真实无 proc 环境验证自动选路。
- 第 3 条：保留叶子直接 `openat`；魔链不是必需加固，反而新增可信路径与 fd 生命周期负担。
- 第 4 条：64 不作为已批准数值；先确定兼容性与计数规则。第 5 条归档路径可由③确认，不影响安全结论。

## Section D 可证伪性审计

F1 可用阶段握手确定性复现，不应靠循环碰运气；F2 可在隔离挂载命名空间复现并记录权限前提；F4 可通过关闭后受控 fd 复用复现。均需同时断言界外 canary、返回值和句柄清理。硬链接可机械核对设备号、inode 与内容；深度边界、拼写净化及假绿变异均可表驱动验证。

本轮取得的是源码、设计与历史日志证据，**没有新增测试／lint 日志，没有实现 diff**。P52 的 23 字节输出是既有实跑，不是本轮复跑。A/B 安全结论及新增反例目前依赖 Linux 语义推演。

已核对外部依据：[open(2)](https://man7.org/linux/man-pages/man2/open.2.html)、[openat2(2)](https://man7.org/linux/man-pages/man2/openat2.2.html)、[rename(2)](https://man7.org/linux/man-pages/man2/rename.2.html)、[link(2)](https://man7.org/linux/man-pages/man2/link.2.html)、[Go 1.22.12 Linux/amd64 常量表](https://raw.githubusercontent.com/golang/go/go1.22.12/src/syscall/zerrors_linux_amd64.go)。实现阶段仍须最低工具链编译与 Linux 原生证据。

## Section E 无法验证部分与推翻条件

尚无修订实现，无法验证真实 flags、错误关闭分支、无 proc 自动降级、fd 泄漏及原生并发行为。NFS、overlayfs 与不同 namespace 配置的普遍保证也未被本轮证明。

若能证明写者无法移动已打开中间目录、挂载导入被强制禁止且打开期间根 fd 不可关闭，F1/F2/F4 的利用前提可被排除；必须提供落实机制，不能只补一句假设。若 Sol 缩小边界为“净化路径、允许导入别名、稳定生命周期下拒绝软链接”，应明确这是承诺收窄，而非证明当前强不变式成立。修订后须重审，当前不得宣称 A/B 始终 fail-closed。

RED-TEAM VERDICT: FINDINGS — 逐级 nofollow 可阻断 GAP-77 的静态软链接载荷，但不足以兑现当前设计的并发对象边界与降级保证。