# 执行自由度清单 — 切片 `CONTRACTS-QUOTA-ALIGN`（准则 v1.41）

依据：`docs/DELIVERY_DIRECTIVE.md` §1.7（执行自由度与升级边界）。本清单由 ⑧a 落盘，逐条登记本片行使的实质性选择、由裁定固定的边界与已知偏离；验证报告 §3 / §13 / §14 与本清单逐条一致。

## 1 F — 行使自由度（实现选择）

| # | 选择 | 为何不是「做什么」 | 证据落点 |
|---|---|---|---|
| F1 | OB-52 取**路径 D**（A / B / C / D 四路径中选 D：分层口径、零代码、零新记录字段） | 四路径都合法，取哪一形态是裁定后的实现选择 | `docs/CONTRACTS.md:192` |
| F2 | 分层句写成三段：请求预算独立执行、CLI 额度为候选平台事实、argv 摘要证明实际传参 | 版式与叙述顺序选择 | `docs/CONTRACTS.md:192` |
| F3 | 费用约束以加粗单句承载，不新增记录字段 | 「加字段」与「加约束句」都合法；后者零 schema 变更 | `docs/CONTRACTS.md:192` |
| F4 | pid 报告文件名以括注声明为当前 pinned stub 的约定、非 wire 契约 | 命名口径选择 | `docs/CONTRACTS.md:232` |
| F5 | 三态语义与互推禁令合并在同一物理句内 | 段落切分选择 | `docs/CONTRACTS.md:232` |
| F6 | spawn 写失败选「杀掉已起子进程 + 留缺失报告」 | 三种收口（留空报告 / 直接返回 / 杀掉并留缺口）都可行，取第三种 | `tools/agent-stub/main.go:187-192` |
| F7 | 保留「（恒定写）」/`(constant write)` 括注 | 措辞选择（⑦ r1 整改后仍保留） | `docs/CONTRACTS.md:232` |
| F8 | 测试以「报告路径为目录」制造写失败 | 构造方式选择：跨平台、不依赖权限位 | `tools/agent-stub/main_test.go:263` |
| F9 | `stubVersion` 不 bump，维持 `agent-stub-0.1.0` | 标签策略选择：pin 全在边界外，重放身份靠 hash | `tools/agent-stub/main.go:28` |
| F10 | spawn 模式不预写空报告；空报告只在非 spawn 运行与 `Start()` 失败时写 | 写序选择（⑤ r2 Low 的整改形态） | `tools/agent-stub/main.go:136-141` / `:179-183` |

## 2 M — 由裁定固定、不得重议

| # | 裁定 | 约束（可机械核对） |
|---|---|---|
| M1 | OB-52 走路径 D | 不得改回 A / B / C；`ai_probe_copilot.go` 必须零改动（`git diff --numstat` 零出现） |
| M2 | 记录不含费用维度但不得读作「费用未发生」 | 约束句不得删（`docs/CONTRACTS.md:192`，检索 `费用未发生` = 1） |
| M3 | 三态语义与互推禁令 | 缺失 = fail-closed、空 = 事实、非空 = 按内容解析；不得互推（`docs/CONTRACTS.md:232`，检索 `互推` = 1） |
| M4 | OB-58 移出本片 | 本片不得改 `internal/chain/postflight.go` 与 `internal/adapters/**` |
| M5 | ⑦ 额度 S 档 1+1，不降级 | 两次均为 `GPT-5.3-Codex (copilot)`（偏移 28,567,553 / 28,745,954） |
| M6 | 不推 gitee | 本片无 commit、无 push |
| M7 | ⑦ r2 的「不得为适配宽松读者收窄契约文本」 | ⑤ r1 High 的 `deferred-by-ruling` 不得转为文本让步 |

## 3 D — 偏离（如实登记）

- D1：② 在 `tmp/` 根创建 `tmp/seam.js` 并立即删除。临时件落点违规（OB-100 形状）；本片第 **1** 次，升级触发为第 3 次。
- D2：③ 自报首次写入 `main_test.go` 时重复 `package main` 行且用 CRLF，已用 `gofmt -w` 修正。闭合，仍登记以保持诚实。
- D3：① 的设计断言「恒定写之后宽松分支不会再触发」为**假**：`child` / `hang` / `fail-fast` / `crash` 在写点前返回。已在报告 §3 更正（⑦ r2 要求）。
- D4：⑦ r1 的 High 属自相矛盾：绝对措辞「恒定写」与本片「spawn 模式不预写空报告」的决策冲突。以「写入口径」改写闭合，括注保留。

## 4 越界自查

- 未改 `docs/DELIVERY_DIRECTIVE{,_EN}.md`、`internal/adapters/ai_probe_copilot.go`、`internal/chain/postflight.go`、`internal/console/**`、`tools/b4-e2e/**`、`schemas/**`、`testdata/**`、`.github/workflows/**`、`tools/gates/lib/criteria/**`、`tools/gates/cjk-newwords.txt`。
- 本片无 commit、无 push；未触 gitee；⑨ / ⑩ 未执行（由 ⑧b 承担）。
- 报告 §1 / §3 / §13 / §14 与本清单逐条一致；本清单不新增结论，也不改变任何门禁判据。

DOCS: DONE
