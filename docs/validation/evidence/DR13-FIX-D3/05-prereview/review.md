PRE-REVIEW: PASS

# DR13-FIX-D3 · ⑤ 预审（只读）· review

审查对象：`impl.diff`（唯一生产改动）、改动后全文副本、③ 两个新测试全文、① design.md、③ 红基线与 tester-report、④ implementer-report 与回归日志（`go-test-adapters.txt`/`go-test-stress.txt` 尾部已亲眼核对）、仓库原文（`pinned_cli.go` 的 D1/D2/clear/release 链、windows/unix 锁实现、ci.yml）。本环未运行任何测试与写操作，结论基于上述证据。

## 1. 发现清单

| 级别 | 位置 | 描述 | 最小修法 |
|---|---|---|---|
| Low | `agent_runner_starting_marker_order_test.go` 的 AST 顺序门 | 理论假绿形态：若将来有人在函数体最前放一个永不执行的 `os.Remove(claim.path)`（如 `if false` 分支），门取最小位置会通过。当前实现为线性五语句，无此形态。 | 不改代码；⑦ 终审确认语句列表形态；将来重构 `release` 时重跑本门并人工核对。 |
| Low | unix 运行时回归 | design §2D(2) 要求另跑 ubuntu 腿，④ 证据仅有 `GOOS=linux go vet`（编译级）；`RefusesLiveHostUnix`、`AdoptsOrphanUnix` 与 flock 语义未在本机验证。 | ⑧c 核对推送后 CI ubuntu 腿（`go test -count=1 ./...` + race）全绿；CI 不可用时在 Linux 环境补跑 adapters 包。 |
| Info | Windows 压力回归 | 3s 窗口、无 `-race`（本机无 gcc，与 R3 登记口径一致）；5× 复跑零违例是概率证据，结构保证由顺序门加 D1/D2 承担。 | 无需修改。 |
| Info | 全仓库测试 | ④ 只跑了 adapters 包全量；`go test ./...` 待门禁/CI 覆盖。改动仅一文件，波及面小。 | 无需修改。 |
| Info | ShareFactProbe | 只记录不断言：`os.OpenFile` 是否阻止删除的两面行为均落入 design 声明的闭合路径，断言任一面都会把环境事实变成契约；本工具链事实（errno 32）已实测并归档进红基线与测试报告。 | 无需修改；保持探针形态。 |

无 Medium 及以上发现，故无需整改重跑。

## 2. 五个重点逐项结论

**① fail-closed 组合安全 — 成立。** 新序为 `removed := os.Remove(...) == nil` → `unlockStartingMarker` → `Close` → `claim.file = nil` → `return removed`，remove 失败时 unlock 与 Close 无条件执行，无 defer/无重试/无 errno 映射；句柄在任何路径都关闭，`claim.file = nil` 保证幂等。remove 失败遗留的「锁空闲 + 路径仍在」是 DR13-FIX 已登记的跨进程残留口径（R5，适用边界同 runRoot、跨 launcher），本片收窄不消除——D1 同进程门拦本 launcher，跨进程由下次扫描的 D2 复核与清理收敛。`removed :=` 先求值无返回值漂移：os.Remove 的结果只依赖路径状态与第三方句柄共享，与 claim 自身句柄是否已关、锁是否已放无关（claim 句柄含 delete 共享，unix 上 unlink 恒允许）。「为保锁而不 Close」在锁随句柄/进程死亡释放的语义下不可实现且必泄漏句柄，实现正确地没有走这条路。

**② 新中间态 — 无新误拒绝路径。** remove 成功后、Close 前（Windows delete-pending 窗口），扫描者 open 得 `ErrNotExist` 走既有 no-op 腿；扫描者 open 先于 remove 时，其句柄（无 delete 共享）会令 remove 失败落入 R5 收敛，或 remove 成功后 Stat 复核（D2）判 `ErrNotExist` 终止；锁空闲可被第三方观察到的前提是 remove 已执行，故「锁空闲 + 路径仍在」与扫描者成功路径互斥。claim 句柄共享 read|write|delete，扫描者 `os.OpenFile` 共享 read|write，两者兼容，不会引入共享冲突型 open 硬错误。release 与 `clearAgentRunnerStartingMarker` 的 remove 互斥（clear 需锁空闲，claim 释放前已 remove），无竞争死锁。

**③ 测试可证伪性 — 成立。** 顺序门要求 `os.Remove(claim.path)`（含别名）与 `unlockStartingMarker` 直接调用均存在、非 defer/go、且 remove 源码位置先于 unlock；删除、移位、defer 化、提取 helper 任一改写都令门红（fail-safe 方向），非「被等价改写假绿」。fail-closed 反测用 raw `CreateFile` 无 delete 共享占名，真实制造删除被拒；再以 `openStartingMarkerForClear` + tryLock 证明锁已释放——若实现「remove 失败即保锁」，此测必红。压力回归防假红口径成立：仅「观察时 releasing + 最终删名成功」计违例，删名被拒残留与孤儿 inode（`os.SameFile` 区分）不计；新序下扫描者 tryLock 成功蕴含 remove 已执行，违例在结构上归零。ShareFactProbe 只记录可接受（理由见发现清单）。

**④ 边界合规 — 未越界。** `git diff --stat` 仅 `agent_runner_starting_marker.go` 一文件 21+/4-；未动 discard 闭包（操作 temp 名，扫描不可见）、`clearAgentRunnerStartingMarker`、windows/unix 锁实现、`closeStartingMarker`/`releaseLaunch`/`discardStartingLaunch` 调用结构、任何测试文件与断言。DR-13 定性结论与 `winners==1` 断言未触碰，且 ④ 压力复跑 5× 含单赢家测试全绿。签名与返回值语义零变化（调用方本就忽略返回值）。

**⑤ 回归完整性 — 基本充分。** 本机 adapters 全量 48.6s 绿、focused 6/6 绿、压力与单赢家 5× 20.8s 绿、build/vet 双平台编译级绿。缺口仅 unix 运行时腿（见 Low 发现），由 CI 双平台矩阵结构性覆盖。

## 3. 无法验证的部分与推翻 PASS 的条件

- 本环受只读约束未重跑测试，结论依赖 ③/④ 证据与已核对的日志尾部；磁盘真实状态已用 `git diff --stat` 与 `git status --porcelain` 独立核对（仅一生产文件改动），排除脏缓冲不一致。
- unix 运行时行为（flock 随 inode、孤儿收养测试）未经本机验证：若 CI ubuntu 腿或 race 腿出现红，PASS 失效。
- 顺序门为源码顺序门而非运行时证明：若生产代码在 ⑦ 终审前被再次改写，需按门重跑并人工核对语句列表。
- 概率型压力测试零违例不构成严格证明；其结构保证来自顺序门与 D1/D2 的组合论证，若 D1/D2 在本片后被改动，本结论需重估。
