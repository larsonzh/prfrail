# 切片验证报告 · `DIRECTIVE-MODEL-LISTS`（v1.18）

> 依据 `docs/DELIVERY_DIRECTIVE.md` §5.1 的串行流水线与 §6.2 的证据落盘义务。**CN 为权威版**，
> `docs/validation/directive-model-lists_EN.md` 为同位镜像。本报告的自合规围栏块见文末「产物清单」。

## 1. 变更概要（含执行自由度清单）

**目标**：黑名单收敛为全称 `kimi-k3`、§2.2 增列 `gpt-6.1-sol`（深度档）与 `gpt-6-luna`（平行备选）、
§5.2 增 ⑦ 三模型分档表、§1.7.4 增临时授权五要素、§6.3 的 G5-b 改为白名单闭集（含 `tools/gates` 实现与自检
夹具）、台账与 ADR 留痕。

**执行自由度清单**：主控作为编排工件落盘于
`docs/validation/evidence/directive-model-lists-freedom-list.md`（F-1 至 F-10、T-1 至 T-3），
本节的表述与之一致：F-1 §6.3 的 G6-5 端点措辞补全；F-2 附录 B.6 载体行改写并删「成本最高」；
F-3 §2.5 旧 G5-b 描述改为指向 §6.3；F-4 §2.2 就绪句去掉过期计数「七个」；F-5 `OB-53` 不翻转为已处置
（第 ④ / ⑤ 项未落）；F-6 深度档触发集与 `gpt-6.1-sol` 的绑定（⑩ 待确认）；F-7 OB-65 至 OB-69 的编号选择；
F-8 Luna 探针的载体与显示名串重试、并追加 Sol 最小调用；F-9 不改 `R2.2`，只登记 OB-67；
F-10 追加三枚 G5-b 边界夹具与 T38 至 T40。

## 2. 交付面与逐条落点

| # | 交付 | 落点 |
|---|---|---|
| ① | §2.2 增两行 + 分档指针 | `docs/DELIVERY_DIRECTIVE.md` §2.2（`modelId` 列增 `gpt-6.1-sol` / `gpt-6-luna`） |
| ② | §2.4 黑名单收敛 | 同一文档 §2.4：仅 `kimi-k3` + 兜底句 + 白名单闭集口径 |
| ③ | §5.2 ⑦ 三模型分档表 + 成本/档位边界注记 | 同一文档 §5.2 |
| ④ | §6.3 G5-b 改白名单闭集 | 同一文档 §6.3；实现 `tools/gates/lib/criteria/g5b.js` |
| ⑤ | §1.7.4 临时授权五要素 | 同一文档 §1.7 |
| ⑥ | 判据实现与自检夹具 | `g5b.js`、`tools/gates/selftest.js`（T29 至 T42）、`tools/gates/testdata/fixtures/g5b-*.txt` 下 13 枚 |
| ⑦ | 版本包与附录台账 | 本片于 ⑧b 落地：文件头版本号、§0.3 的版本条目、附录 C.0j |
| ⑧ | Luna 工具的实测 | 本片 L1 至 L5 探针（读数见 ADR-019 与第 6 节） |
| ⑨ | 留痕 | `docs/ADR_REGISTER{,_EN}.md` 的 ADR-019；`docs/DELIVERY_DIRECTIVE.md` §12.4 的 OB-65 至 OB-69 |

## 3. 异常与兜底记录

- **L1 首轮被工具拒绝**：以 `modelId` 形式 `gpt-6-luna` 调用 `runSubagent` 时工具回传
  「not found」及可用模型清单；改用清单中的显示名串后 L1 至 L5 全绿。该差异登记为 **OB-67**（待议）。
- **⑤ 额度用尽后的未闭合发现**：⑤ 第 2 轮报 R2-1（§2.5 残留旧判据描述），已就地整改；因 ⑤ 额度
  （1 + 1 整改后复审）已用尽，未加轮，故按 §7.8.3「额度用尽」行的留痕要求登记为 **OB-69**
  （其处置列首词为待办，内容已并入本片 ⑦ 的审查输入），并由本片 ⑦ 复审覆盖。
- **自由度清单自身引入的红灯**：⑦ 终审指出该清单触发 G4b（引入了不在语料内的新增汉字）与 G6-4
  （行数字面量形态），已改写措辞，复跑门禁 `TOTAL_FAIL=0`。

## 4. 执行流水线

① 架构（`deep-reasoner` / DeepSeek V4 Pro）⇒ ② 实现（主控 + 落盘脚本）⇒ ③ 文档一致性（机械对账 32 项）⇒
④ 门禁（tree / selftest / index）⇒ ⑤ 预审 2 轮（DeepSeek V4 Pro）⇒ ⑥ 独立扫描（MAI-Code-1.1-Flash）⇒
自由度清单 ⇒ ⑦ 终审（盲态）与复审（GPT-5.3-Codex）⇒ ⑧ 报告与回写 ⇒ ⑨ 原生验证 ⇒ ⑩ 停点。
**额度实得**：① 1/1、⑤ 2/2（1 + 整改后复审）、⑥ 1/3、⑦ 2/2（终审 + 复审）、探针 0（本片未声明探针额度）。

## 5. 环境与档位

Windows 本机；Go 工具链 `go1.27.0 windows/amd64`；Node v24.17.0；仓库 `github.com/larsonzh/prfrail`，
分支 `main`。本片不改 Go 源码，不触 `internal/**`、`CONTRACTS`、`schemas`、`fixtures`、`workflows`、
`go.mod`。

## 6. 审查结论摘要（⑤ ⑥ ⑦）

- **⑤ 第 1 轮**（DeepSeek V4 Pro）：`PRE-REVIEW: FINDINGS`，1 High（F1 版本包时序）+ 2 Medium（F2 四处
  ⑦ Codex 标签、F3 深度档绑定待确认）+ 5 Low；已逐条处置（F1 修引用、F2 登记 OB-68、F3 转 ⑩ 待确认、
  F4/F5/F6/F7 就地整改、F8 记为自由度清单义务）。
- **⑤ 第 2 轮**：`PRE-REVIEW: FINDINGS`，新增 1 Medium（R2-1） + 3 Low 记录；R2-1 已整改（见第 3 节）。
- **⑥ 独立扫描**（MAI-Code-1.1-Flash）：`INDEPENDENT SCAN: FINDINGS`，但报告**未给出 file:line 级缺陷**
  （不符 §4.1 的输出契约），故未开启 ⑥ 整改环；其三条实质主张分别落在 OB-67（已登记）与已声明保留文本
  （§2.4 的 `Kimi K3` 配置引文、§2.5 的族分解历史段）上。判定与结论原文见自由度清单 T-3。
- **⑦ 终审（盲态，不附清单）**：`FINAL REVIEW: FINDINGS`，2 High（H1 自由度清单红灯、H2 `R2.2` 与
  OB-67 冲突）+ 2 Medium（M3 分档表与 §7.5 失同步、M4 新成员无正例夹具）+ 1 High（H5 上一片报告证据链）。
- **⑦ 复审**：`RE-REVIEW: PASS WITH FIXES`——H1 / M4 / H5 关闭；M3 关闭为「已披露、待后续裁决」；
  **H2 未闭合**：⑦ 明确要求 **Option B**（本片不擅改契约语义，保留 `待议` 并升级用户裁决），
  且须在 ⑩ 停点把三案交用户裁决（改 `R2.2` / 增 `modelId` 到调用串的映射 / 保留并在 §2.2 明示差异），
  **裁决前不得宣称语义已闭环**。
- **审查输入包与盲审隔离**：⑦ 终审以盲态执行（不附 ⑤/⑥ 报告与自由度清单），理由 = 本片触及模型**身份**
  语义（白名单 / 黑名单即模型准入身份）⇒ 按 §7.7 必抽盲审；⑦ 复审为收窄输入（本轮整改 diff + 第 1 轮
  发现与处置清单 + 契约原句 + 只读约束）。

## 7. 可证伪性（机械检查与变异检验）

- **G5-b 判据**：T29 至 T42 的 14 项断言覆盖闭集成员、缺键、黑名单残留、大小写、引号、多键全收、
  空值分工（与 G5-a ② 双向锚定）、frontmatter 外忽略、两端修剪、行尾注释、以及两枚新增成员的正例；
  **变异锚 T37 / T42** 分别验证「闭集被放宽」与「成员被删除」两条弱化路径都会翻红，并断言
  `g5b.js` 在变异后字节级还原。
- **CJK 与过期字面量**：本片自身两次触发 G4b / G6-4（自由度清单），均按「改写规避、不扩白名单」处置。
- **A11 对称抽查**：本轮 ⑥ 未给出「删除守卫 X ⇒ 测试 Y 变红」类声明，故无抽查对象；⑦ 的变异类声明由
  T37 / T42 的断言直接承担。

## 8. 门禁结果

| 作用域 | 结果 |
|---|---|
| `--scope=tree` | exit 0；`changed=24`；`TOTAL_FAIL=0`（含 G1-a、G1-b、G2、G3(a)/G3(b)、G4a、G4b、G5-a、G5-b、G6-1 至 G6-6、G7-a 至 G7-d） |
| `--selftest` | `SELFTEST: PASS 264/264`（较切片起点 248 增 16 项） |
| `--scope=index` | 见 ⑨ 原生验证（暂存后执行） |
| 基础门禁 | `go build ./...` / `go vet ./...` / `go test ./...` 见 ⑨ |

## 9. 已知差异与覆盖缺口

- **`R2.2` 的调用串口径与工具现实不符**（OB-67，待议）：本片不改契约语义，只登记并升级。
- **§5.2 分档表与 §7.5 / §8.2 的「⑦ Codex」行失同步**（OB-68，待议）：本片在表下加了成本 / 档位边界
  注记，四处原句保持不动。
- **深度档触发集与 `gpt-6.1-sol` 的绑定**（F-6）：属本片新增语义，⑩ 停点交用户确认。
- **⑥ 报告不符合 §4.1 输出契约**：无 file:line 级缺陷，未开启整改环；结论原文与裁定见自由度清单 T-3。
- **⑧b 版本包未经 ⑦ 复核**：⑤ 第 2 轮的附条件要求，⑩ 停点明示供用户否决。
- **残余历史引用**：`附录 C.0–C.0i` 的历史行（§0.3 的 v1.13 条目、OB-12、OB-44）保持原样（§11.1 不重写历史）。

## 10. 候选 OB（本片已落 §12.4）

OB-65（试点顺延第 3 次，用户直接裁定）、OB-66（顺延次数口径定义，下一片候选）、OB-67（`R2.2` 调用串
口径与工具现实不符）、OB-68（§5.2 分档表与四处 ⑦ Codex 标签失同步）、OB-69（⑤ 第 2 轮未闭合发现，
并入 ⑦ 输入）。

## 11. 未完成清单与下一步

| 项 | 性质 | 去向 |
|---|---|---|
| `R2.2` 口径 | 契约语义 | ⑩ 停点用户裁决（三案择一，⑦ 要求 Option B 前置） |
| 深度档触发集绑定 | 新增语义确认 | ⑩ 停点用户确认 |
| §7.5 / §8.2 的成本与档位定义 | 规则缺口 | OB-68，下一片 |
| `model: "auto"` 的链配置侧 | 契约缺口 | OB-53 第 ⑤ 项，下一片 |
| 顺延次数口径 | 台账口径 | OB-66，下一片 |
| ⑤ 是否补第 3 轮 | 额度 | ⑩ 停点用户裁定 |

## 12. ⑩ 停点报告

- 执行至 ⑦ 复审；⑩ 停点为**待授权提交**：本片**不 commit、不 push、不推 gitee**，等待用户同轮授权。
- 状态：`DIRECTIVE-MODEL-LISTS: STARTED | 待用户裁决（H2 / F-6 / T-1 / ⑤ 额度）`。
- 待裁决四项：① `R2.2` 三案；② 深度档绑定确认；③ 版本包于 ⑧b 落地且未经 ⑦ 复核是否接受；
  ④ ⑤ 是否补第 3 轮。
- **回写核对**（OB-63 口径）：验证报告 = 已回写（本文件）；`DEV_PLAN` = **不适用**（本片为准则治理片，
  其权威台账为准则附录 C，见 §11.2）；台账 = 已回写（§12.4 的 OB-65 至 OB-69 与附录 C.0j）；
  ADR = 已回写（ADR-019）。

## 13. 产物清单

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/t027-assembly-min.md	repo
docs/validation/t027-assembly-min_EN.md	repo
docs/validation/directive-model-lists.md	repo
docs/validation/directive-model-lists_EN.md	repo
docs/validation/evidence/directive-model-lists-freedom-list.md	repo
tools/gates/lib/criteria/g5b.js	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/g5b-body-only.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-luna.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-sol.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal.txt	repo
tools/gates/testdata/fixtures/g5b-empty-value.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-case.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-inline-comment.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-kimi.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-legacy.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-multi.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-quoted.txt	repo
tools/gates/testdata/fixtures/g5b-legal-nokey.txt	repo
tools/gates/testdata/fixtures/g5b-legal-trailing-space.txt	repo
tmp/dml/design-1-v4pro.md	local-only
tmp/dml/change-set.diff	local-only
tmp/dml/contract-excerpt.md	local-only
tmp/dml/consistency-out4.txt	local-only
tmp/dml/gate-tree7.txt	local-only
tmp/dml/selftest7.txt	local-only
tmp/dml/ci-run-37073985144.txt	local-only
tmp/dml/luna-probe/target.txt	local-only
```

## 16. ⑩ 停点后的裁决回执（2026-10-03，事实记录）

- 用户在 ⑩ 停点给出四项裁决：① `R2.2` 口径选 **(C)**（保留 `R2.2` 本体，在 §2.2 加实测差异注记，并登记 OB-70）；② **确认**深度档触发集与 `gpt-6.1-sol` 的绑定；③ **接受**版本包于 ⑧b 落地且未经 ⑦ 复核（属事实记录，符合 OB-11 例外标准，不触发重跑 ⑦）；④ ⑤ **不补**第 3 轮（R2-1 由 ⑦ 覆盖闭环）。
- 顺延次数口径：**接受「以 §12.4 为唯一真值源」**；回执中的口头计数不作数，本次记为第 3 次（OB-65），下一片回填 OB-66。
- 本节的 §11 / §12 两节保留 ⑩ 停点当时的表述（`待裁决` 状态），**不改写结论行**：本次只追加本事实记录，未重跑 ⑦（额度用满，用户明示不开第 4 轮）。
