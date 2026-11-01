# P45 执行自由度清单

| 项 | 内容 |
|---|---|
| 切片 | `P45-T027B-E2E-ZERO-MODEL`（M 档 · 装配验证片 · `[SLICE]`） |
| 可自行决定的自由度 | ① 三任务链的驱动入口选择（console 既有装配入口 vs 新增 `examples/` 驱动）；② 声明式 stub 的形态（配置文件 / 环境变量 / 测试替身）；③ T2 自然失败的具体触发方式（必须在协议语义内，禁止伪造状态）；④ 断言的组织方式（单测试 vs 子测试） |
| 必须升级的事项 | 需要改动 `internal/adapters/**`（本片禁改）；需要新增 CLI 子命令（超出 console 既有装配面）；需要改动 `CONTRACTS` / schema；发现 chain 状态机语义缺陷需先协议后代码 |
| 载体铁律 | ① ⑤ 只读 `deep-reasoner`；② `prfrail-implementer`、③ `prfrail-tester`、⑧ `prfrail-documenter`；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`；每次 `runSubagent` 必须同时传 `agentName` 与 `model` |
| 主控禁止 | 禁止代写 ②③ 的代码与测试；主控只做编排、落盘、门禁与回执 |
| 批量辅助 | `GLM-5.3-Flash (glm)` ≥ 2 次（建议 4 次），实派次数须计入 ⑩ 成本与计量 |
| 并发断言（DR13-FIX 强输入） | 必须含「同一 launcher/runRoot 并发 Start」：恰一赢家、败者全命中 `Conflict`/`Duplicate`、closeout 时 `registry.len()==1`、赢家 `ProcessAlive`、slot owner token 与 mirror 一致；**禁止以 10s freshness 墙钟作排他断言** |
| 零模型证明 | 必须以机械判据证明零模型调用（stub 调用计数 / 无真实 provider 记录），不得仅凭"没看到调用" |
| 临时件 | 一律 `tmp/P45-T027B-E2E-ZERO-MODEL/`，只增不删，⑧c 清理（只删本片目录） |
| 禁止 | 推 gitee；在 ⑩ 之前 commit/push |
