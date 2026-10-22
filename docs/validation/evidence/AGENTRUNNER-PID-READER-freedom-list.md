# 执行自由度清单 — 切片 `AGENTRUNNER-PID-READER`（准则 v1.42）

依据：`docs/DELIVERY_DIRECTIVE.md` §1.7（执行自由度与升级边界）。本清单由 ⑧a 落盘，逐条登记本片行使的实质性选择、由裁定固定的边界与已知偏离；验证报告 §5 / §6 / §9 与本清单逐条一致。

## 1 F — 行使自由度（实现选择）

| # | 选择 | 为何不是「做什么」 | 证据落点 |
|---|---|---|---|
| F1 | 第六枚事实的内层摘要用空域参数 `evidence.Digest("", bytes)` | A6 只固定外层域与「字节级摘要」，内层域留空沿用 `evidenceRootHash` 口径；用户 2026-10-22 裁定属 §1.7.1 契约范围内的实现方式 | `internal/adapters/agent_runner_evidence.go:370-371` |
| F2 | 实施顺序取 OB-110 先、OB-109 后 | 两项均可先行；先修读者三态再绑第六枚，使每个中间态都能单独验证 | `internal/adapters/agent_runner_run.go:314` / `internal/adapters/agent_runner_evidence.go:369` |
| F3 | 读者返回值形态取 `(pids, present, err)` 三值 | 哨兵错误或单独的 `hasFile` 布尔同样合法；三值把「缺失」与「空」在类型层面分开 | `internal/adapters/agent_runner_run.go:314` |
| F4 | 缺失报告立即返回，不烧完等待上界 | 轮询到上界再判同样 fail-closed；提前返回避免空转 | `internal/adapters/agent_runner_run.go:270-274` |
| F5 | 进程观测函数以参数注入 | 直接调用 `guard` 包同样合法；注入使 fail-closed 分支可无竞态测试 | `internal/adapters/agent_runner_run.go:261` |
| F6 | postflight 第 5 腿取 `(alive, missing, err)` 并列返回 | 把 missing 折叠进 error 也可行；并列返回让调用方分别决定 `uncertain` 与 `failed` | `internal/adapters/agent_runner_postflight.go:507` |
| F7 | harness 判据名 `five-facts-distinct` 改为 `six-facts-distinct` | 保留旧名并加注也可行；改名使判据自述与六枚事实一致 | `tools/b4-e2e/main.go:690` |
| F8 | `observation.Facts` 改为直接整体赋值 | 逐字段复制同样可行；整体赋值消除漏字段的可能 | `tools/b4-e2e/runmap.go:355` |
| F9 | `PidReportsComplete` 取独立布尔 | 由「摘要为空」推断缺失也可行；独立布尔不依赖摘要取值 | `internal/adapters/agent_runner_evidence.go:79` |
| F10 | harness `readPIDs` 保留 `fmt.Sscanf` | 与生产同步改 `strconv.Atoi` 同样可行；本片只登记为 Low，不扩大改动面 | `tools/b4-e2e/ports.go:525` |

## 2 M — 由裁定固定、不得重议

| # | 裁定 | 约束（可机械核对） |
|---|---|---|
| M1 | 形态取丙（5 个适配器文件 + chain 改动），不采纳丁 / 乙 / 甲 | 不得改为乙（漏父报告）或甲；`git diff --cached --name-only` 必须覆盖 18 个代码文件 |
| M2 | 三态语义不得互推 | 缺失 = fail-closed、空 = 事实、非空 = 按内容解析（`docs/CONTRACTS.md` A6 段） |
| M3 | 自然失败优先级不得改变 | 自然非零退出仍判 `failed`，`uncertain` 不覆盖已成立的自然失败 |
| M4 | 第 7 腿必须重读重算而非只读冻结值 | 冻结值不得作为自身的证据；重算不一致即 `failed` |
| M5 | 准则只改三处 | §0.3 一行、§12.4 一行、头部版本号；不改任何判据本体与作用域 |
| M6 | 本片无 commit、无 push、不推 gitee | 由主控承担；⑧a 只写三份文档 |

## 3 D — 偏离（如实登记）

- D1：②a 在 `tmp/` 根创建 `tmp/verify-*.js` 与 `tmp/pid-reader-dump.md`，自报 + 自删 + 零残留核验通过。按 §5.1「显式登记，不销解、不追加动作」登记；本片第 **1** 次，§5.1 合并口径第 **11** 次。

## 4 越界自查

- 未改 `internal/console/**`、`tools/agent-stub/**`、`schemas/**`、`testdata/**`、`.github/workflows/**`、`tools/gates/**`。
- 未新增 face；未处置 OB-104 / OB-105 / OB-107。
- 本片无 commit、无 push；未触 gitee。
- 报告 §5 / §6 / §9 与本清单逐条一致；本清单不新增结论，也不改变任何门禁判据。

DOCS: DONE
