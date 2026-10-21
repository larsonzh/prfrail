# 执行自由度清单 · 切片 `AGENTRUNNER-PROD-PATCH`

> 归属：`docs/validation/agentrunner-prod-patch{,_EN}.md` 的 §13 引用本清单为权威，本清单不改变报告含义。
> 口径：F = 本片自行作出的选择；M = 已裁定、不得重开的判定；D = 偏离（异常）登记。每行给出 `file:line` 或偏移证据。

## F 选择

| # | 选择 | 证据 |
|---|---|---|
| F1 | 取路径 A：三条 spawn 后失败分支就地回收 | `internal/adapters/agent_runner_replay_dispatcher.go:205` / `:216` / `:220` |
| F2 | 停止能力用私有可选接口承载，不扩 chain 端口 | `internal/adapters/agent_runner_replay_dispatcher.go:70` |
| F3 | 无停止能力时给出显式「未回收」出口，不静默 | `internal/adapters/agent_runner_replay_dispatcher.go:84` |
| F4 | grace 缺省取 1s，停止 ctx 用 `grace+1s` 有界 | `internal/adapters/agent_runner_replay_dispatcher.go:86-89` |
| F5 | 停止只尝试一次，不做循环重试 | `internal/adapters/agent_runner_replay_dispatcher.go:92` |
| F6 | cause 始终为首 sentinel，停止错误与结论只作尾注 | `internal/adapters/agent_runner_replay_dispatcher.go:93-100` |
| F7 | 停止错误用 `%v` 拼接，保持单 sentinel 链 | `internal/adapters/agent_runner_replay_dispatcher.go:95` |
| F8 | 不重试身份、不删启动意图、不覆盖既有身份 | `internal/adapters/agent_runner_replay_dispatcher.go:216-221` |
| F9 | 补一个真实进程测试观察 OS 事实 | `internal/adapters/agent_runner_replay_dispatcher_test.go:718` |
| F10 | `reclaimer == nil` 子句保留，配防御性注释 | `internal/adapters/agent_runner_replay_dispatcher.go:83` |
| F11 | 变异证据放在 `tmp/AGENTRUNNER-PROD-PATCH/mutation/`，用 `git archive` 取干净树而非 `git worktree` | `tmp/AGENTRUNNER-PROD-PATCH/mutation/evidence.txt:1-3` |
| F12 | 报告落盘后 artifacts 围栏置于文末，§1 只给出路径数与去向 | `docs/validation/agentrunner-prod-patch.md` §1 与文末围栏 |

## M 判定（不得重开）

| # | 判定 | 依据 |
|---|---|---|
| M1 | 路径 A 为准；路径 B 否决（需三个新 API + CONTRACTS 变更 ⇒ 越界；且把活进程泄漏转为阻塞重派的持久预留） | ① 设计；报告 §2 |
| M2 | 回收目标 = dispatcher 自证 id，不是 `result.LaunchID` | ⑦ r1 盲审独立裁定（slot 语义）；⑥ High 因此不采纳；`agent_runner_replay_dispatcher.go:205` |
| M3 | 禁止身份重试、禁止删除启动意图、禁止覆盖既有身份 | ① G3；no-replace store；`agent_runner_replay_dispatcher.go:216-221` |
| M4 | 停止错误用 `%v` 而非 `%w`（保持单 sentinel 链） | ⑤ r1 背书；`agent_runner_replay_dispatcher.go:94` |
| M5 | D1 继续本片并计为该放置规则的第 7 次；变异证据须重出 | 用户 2026-10-21 裁定；报告 §3 D1 |
| M6 | `reclaimer == nil` 保留 + 防御注释；任何使其可达的改动必须补测试 | 用户裁定；`agent_runner_replay_dispatcher.go:83` |

## D 偏离

| # | 偏离 | 影响 | 证据 |
|---|---|---|---|
| D1 | ③ 在 `%TEMP%` 下建 detached worktree 并删除（自报） | 临时件落点违规，按合并口径计第 7 次；重出计 ③ 1 次返工 | 报告 §3 D1；`tmp/AGENTRUNNER-PROD-PATCH/mutation/evidence.txt` |
| D2 | 首个 M2 变异 `if true {` 编译失败（`declared and not used: ok`） | 工具迭代，非缺陷 | 报告 §3 D2 |
| D3 | `%TEMP%/proofrail-agent-stub-parity` 来自既有助手 `buildAgentRunnerParityStub` 写 `os.TempDir()` | 候选 OB（无号），非本片缺陷 | 报告 §3 D3、§14 U2 |
| D4 | `-race` 本机不可执行（无 cgo / gcc） | 验收只能走 CI ubuntu 腿 | 报告 §4 |
| D5 | ④ 的 `G5(d)` / `G5(f)` 红 | 报告落地前 / 暂存前的结构状态，非缺陷 | 报告 §3 D5、§8 |
| D6 | 派单简报的 `+54/−7` 与 `git diff --numstat` 实测 `+49/−5` 不一致 | 以实测为准，差异如实登记 | 报告 §14 U6 |

- 本清单不新增任何未在报告中登记的自由度；F12 只描述本报告的围栏位置，不构成规格。
