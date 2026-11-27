# P97 执行自由度清单（§1.7.3）

切片：`HYDROGEN-CONTROL-PILOT`（P97）。主控：`DeepSeek V4.1 Flash (deepseek)` High。落盘时点：步骤 0（⑥ 之后、⑦ 之前若有增补则同轮补记）。

> 本清单位于 `docs/validation/evidence/**` ⇒ 无机械判据保护，仅由 ⑦ 人工复核是否越界（§1.7.3 覆盖边界）。

| 编号 | 自由度选择 | 依据 | 越界判定 |
|---|---|---|---|
| F-1 | 对照对象取 P87 自身输入包（`P87-DIRECTIVE-REVIEW-MODEL-POLICY/`），不启用 P96 备选 | 起跑消息「六 实现要点」指定首选；P87 证据包实测可得 | in-bounds |
| F-2 | 对照面固定为四面（结论 / 发现数 / 严重度 / 误判），且只记录不重裁 | 起跑消息「六」明文；不改 P87 已闭环结论 | in-bounds |
| F-3 | 试水输入包只含 P87 交付物与基准结论，剔除本片任何结论 | §7.3 输入包消毒；避免自证 | in-bounds |
| F-4 | 台账 P97 段按起跑消息 §〇.5 先以「⬜ 未立项」落盘，随即翻「执行中」 | 起跑消息「四 步骤 0 必做」第 2 项 | in-bounds |
| F-5 | 证据包目录形态沿用 P96（00-briefing / input-list / freedom-list / 01-arch / 03-verification / 04-reviews / MANIFEST / SHA256SUMS） | §13.1.10、§13.1.12 与既有切片惯例 | in-bounds |
| F-6 | 报告对文件名取用户给定的 `hydrogen-control-pilot{,_EN}.md` | 起跑消息「七」第 4 项 | in-bounds |
| F-7 | 批量辅助派单仅限双语镜像核对与数字核对，字符上限随单写明 | 起跑消息裁定 d（`glm-5.3-flash`，禁 `glm-5.3`） | in-bounds |
| F-8 | 编码归一以断言脚本落盘（先断言长度再写回单 BOM + LF） | `docs/CODING_CONVENTIONS.md`；P87 步骤 0 曾因误写清零五件 | in-bounds |

## 越界自查

- 未改任何规则语义、判据本体、模型白名单与 G5-b 闭集。
- 未新增 ADR；未改 `internal/**` / `tools/gates/**` / `ci.yml` / CONTRACTS / schema / fixtures / workflows。
- 未做 gitee；未启动 Ubuntu VM 与代理。
