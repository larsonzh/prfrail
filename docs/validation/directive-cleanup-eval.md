# 切片验证报告 · `DIRECTIVE-CLEANUP-EVAL`（清账 + R2.2 传参评估）

> 依据 `docs/DELIVERY_DIRECTIVE.md` v1.18 的 §5.1 串行流水线与 §6.2 证据落盘义务。**CN 为权威版**，
> `_EN` 为同位镜像；围栏块见文末「产物清单」。规模 S、分级 `[SLICE]`、探针 0。

## 1. 变更概要（含执行自由度清单）

**目标**：① 清两处台账过期陈述；② 回填 OB-66 口径；③ 完成 OB-70 评估；④ 落地 OB-63。
**本片行使的执行自由度（§1.7.1，实质选择）**：

- **F-1 一字更正**：用户指令原文的「本机为事实记录」按语义改为「**本节**为事实记录」（写入报告 §16 首句）；
  依据 = 该句主语为「报告的一节」，逐字照抄会造成事实错误；已在回执与本表声明。
- **F-2 落点选择**：OB-63 的要求落在 **§11.2「回写口径」**（用户给出「§5.3 或 §11.3」两选）——
  理由 = OB-63 自身的处置列已写「口径见 §11.2」，落 §11.2 可让该指针解析成立；已声明。
- **F-3 处置列翻转**：OB-63 与 OB-70 的处置列由 `已登记` / `待议` 翻为 `已处置`（前者因要求已入正文，
  后者因评估已在本片完成）；两行均为本片范围。
- **F-4 ⑦ 降级**：本片 ⑦ 用 `deepseek-v4-pro`（用户切片定义明示），按 OB-37 元规则以 **ADR-020** 留痕。

## 2. 启动前清单核对（步骤 0；逐条）

| 项 | 内容 | 状态 |
|---|---|---|
| A-1 | `REMAINING_SLICES{,_EN}` 的 T027-ASSEMBLY-MIN 行改为「已提交 `68a1ef0` + `d2d6370`，CI run `37073985144` 双绿」 | **已处理**（由上一切片步骤 0 完成，本片逐字复核确认） |
| A-2 | `REMAINING_SLICES{,_EN}` 的 DIRECTIVE-MODEL-LISTS 行改为「已提交 `35e8198`，CI run `37092310920` 双绿」 | **已处理**（本片 B 之前落盘） |
| B-3 | `directive-model-lists{,_EN}.md` 的 §16 首句补「本节为事实记录，不改写 §11 / §12 的 ⑩ 停点当时表述」 | **已处理**（本片落盘） |
| C-4 | OB-66 处置列明确「以 §12.4 台账为唯一真值源」+ 本片顺延第 3 次 | **已处理**（本片落盘） |
| D-5 | 以 Codex / V4 Pro / Flash 各试一次 `modelId` 形式 | **已处理**（三试全被拒，见 §3） |
| D-6 | 查 `runSubagent` 官方文档 | **已处理**（两页官方文档，见 §3） |
| D-7 | 评估显示名稳定性 | **已处理**（见 §3） |
| D-8 | 综合判断三案并写入报告（不改 `R2.2` 本体、不改 §2.2 表） | **已处理**（见 §3） |
| D-9 | 评估结论作为下一片立项依据 | **已处理**（推荐案与立项要点见 §3 末） |
| E-10 | OB-63 要求落地为规则正文 | **已处理**（落 §11.2，见 §4） |
| 既有-1 | OB-68（§7.5 / §8.2 成本与档位定义） | **待办**（下一片；OB-68 处置列 `待议`） |
| 既有-2 | OB-53 第 ⑤ 项（链配置侧 `model: "auto"`） | **待办**（下一片；触 `internal/**`，本片禁改） |
| 既有-3 | OB-67（`gpt-6-luna` 首轮被拒） | **评估已完成**（与 OB-70 同源；修法裁定仍待下一片，OB-67 处置列维持 `待议`） |
| 既有-4 | 主控交接审计（= T027 结束时） | **不适用**（本片非该锚点） |
| 既有-5 | 战略回看（= S1 结束时） | **不适用**（本片非该锚点） |

## 3. ① 评估：`runSubagent` 的 `model` 传参形式（本片核心）

**实测（三次，D-5）**：以 `modelId` 形式调用 `runSubagent`，三个模型**全部被拒**，工具每次回传同一份可用清单：

| 传入值（`modelId` 形式） | 结果 | 清单中的对应条目（限定名） |
|---|---|---|
| `gpt-5.3-codex` | **not found** | `GPT-5.3-Codex (copilot)` |
| `deepseek-v4-pro` | **not found** | `DeepSeek V4 Pro (deepseek)` |
| `deepseek-flash` | **not found** | `DeepSeek V4.1 Flash (deepseek)` |

对照既往前片实测：`gpt-6-luna` 被拒，而 `GPT-6 Luna (copilot)` 通过（L1 至 L5 全绿）、
`GPT-6.1 Sol (copilot)` 通过 ⇒ **拒绝原因是名称形式，而非模型不可用或档位不可用**（同一会话下档位更高的
Luna / Sol 反而可用）。

**证据留存**：三次被拒调用的回执原文留存于 `tmp/dce/probe-id-rejections.txt`（`local-only`），本报告的结论不再只有文字。

**官方文档（D-6，两页）**：

- 自定义代理页的 `handoffs.model` 明确：**Use the qualified model name in the format `Model Name (vendor)`**，
  例如 `GPT-5 (copilot)`、`Claude Sonnet 4.5 (copilot)` ⇒ **限定名（显示名 + vendor）是文档规定的模型写法**；
- 子代理页「Select the model for a subagent」（**Local harness 页签下的内容**——该页写明本地 harness 使用 `runSubagent` 工具）规定本地子代理选模顺序：① 主控给 `runSubagent` 的**显式
  `model` 参数** → ② 自定义代理的 `model` 属性 → ③ Auto → ④ 主会话模型；并要求把 `<model name>`
  替换为**本会话可用的模型**；排障表规定「**A requested model doesn't run ⇒ Use one of the models listed
  in the error**」；
- 同页另注：显式选模会**经主会话模型的成本档位校验**，超档时不执行并回报可用模型。该条**不能单独解释**
  本片实测（Luna / Sol 档位更高却可用），故拒绝原因收敛为**名称形式**。

**显示名稳定性（D-7）**：限定名**自带版本与 vendor**（`DeepSeek V4.1 Flash (deepseek)`；文档示例
`Claude Sonnet 4.5 (copilot)`），清单内另有 `Auto (copilot)`（语义随设置变化）⇒ 版本升级与 vendor
变更都会改变调用串。⇒ `R2.2` 的初始立意（防显示名漂移）**经实测与文档双重支持**。

**三案评估（D-8；不改 `R2.2` 本体、不改 §2.2 表）**：

| 案 | 结论 | 理由 |
|---|---|---|
| ① 改 `R2.2` 为显示名口径 | **不推荐** | 虽与文档一致，但把**版本 / vendor 漂移**引入契约，违背 `R2.2` 的防漂移立意 |
| ② 加映射表（`modelId` → 限定名） | **推荐** | 任务包仍写稳定的 `modelId`，调用层按**唯一映射表**取限定名；`R2.2` 本体只需加一句指向映射表；映射表可加机械判据（成员集合相等、限定名逐字来自工具清单） |
| ③ 保留 `R2.2` + 注记 | 现状（前片已落） | 零改动、零风险，但每次调用需人工查清单，规则与注记并存易被忽略 |

**D-9 交付下一片**：若采纳案 ②，须另立 `[SLICE]`（改 `R2.2` 属契约语义变更）：① 建映射表（§2.2 的
`modelId` 列 ↔ 工具清单限定名）；② 在 `R2.2` 增指向句；③ 增机械判据。

## 4. 交付面与逐条落点

| 项 | 落点 |
|---|---|
| A-1 | 复核确认（逐字命中「已提交 `68a1ef0`」「`d2d6370`」「`37073985144`」） |
| A-2 | `docs/t027/REMAINING_SLICES{,_EN}.md` 的 DIRECTIVE-MODEL-LISTS 行 |
| B-3 | `docs/validation/directive-model-lists{,_EN}.md` §16 首句 |
| C-4 | `docs/DELIVERY_DIRECTIVE.md` §12.4 的 OB-66 处置列（`已裁决`） |
| E-10 | `docs/DELIVERY_DIRECTIVE.md` §11.2 正文（OB-63 要求）；OB-63 处置列翻 `已处置` |
| D | OB-70 处置列翻 `已处置`；评估结论入本报告 §3 |
| 留痕 | `docs/ADR_REGISTER{,_EN}.md` 的 ADR-020（⑦ 降级）；上一片报告 `directive-model-lists` 的既有 §16 保持不动 |

## 5. 执行流水线与环境

①（评估与文档实现由主控承担，无独立 ① 调用）⇒ ② 文档实现 ⇒ ③ 清单核对与一致性 ⇒ ④ 门禁 ⇒
⑤ **跳过**（§5.2 的 S 档为「可选」）⇒ ⑥ **跳过**（纯 `.md`、无代码语义，按 §7.2 的跳过条件，由 ⑦ 确认）
⇒ ⑦ 终审（`deepseek-v4-pro`，**ADR-020 降级**；**独立性降低如实标注**）⇒ ⑧ 报告与回写 ⇒ ⑨ 原生验证 ⇒ ⑩ 停点。
**额度实得**：探针 0；本次实测的 3 次 `modelId` 调用均为**被工具拒绝的无效调用**（不产生模型用量）。
环境：Windows 本机；Node v24.17.0；Go 工具链 `go1.27.0 windows/amd64`；仓库 `github.com/larsonzh/prfrail`。

## 6. 审查结论

- ⑤：**跳过**（S 档可选）。
- ⑥：**跳过**（纯文档切片；按 §7.2 的跳过条件，报告 §6 由 ⑦ 确认该跳过的成立性）。
- ⑦：第 1 轮 `FINAL REVIEW: FINDINGS`（3 Medium + 9 Low，已全部整改）；复审 `RE-REVIEW: PASS WITH FIXES`（残留 3 条 Low 由 ⑧b 闭合，⑦ 明示无需再开轮次）。本片 ⑦ 由 `deepseek-v4-pro` 承担 ⇒ 独立性降低，已在 ADR-020 与本节如实标注。

## 7. 可证伪性与门禁结果

| 检查 | 结果 |
|---|---|
| `gate.js --all --scope=tree` | 见 §8；首跑因 G3(b) 的 to-do 词表 0/2 而红，整改后复绿 |
| `gate.js --all --scope=index` | 见 §8 |
| `gate.js --selftest` | 见 §8 |
| 清单核对（§2） | 15 项逐条标注，其中 A-1/A-2/B-3/C-4/D-5…D-9/E-10 已处理 |

## 8. 门禁与原生验证

见 §9（门禁读数已回填，含首跑红与整改说明）。

## 9. 门禁读数

- `--scope=tree`：`TOTAL_FAIL=0`。
- `--scope=index`：`TOTAL_FAIL=0`。
- `--selftest`：`PASS 264/264`（本片未改判据与夹具，读数应与切片起点一致）。
- 基础门禁：`gofmt` 空、`go build` / `go vet` / `go test` 均 exit 0。

## 10. 已知差异与未完成

- OB-68（§7.5 / §8.2 的成本与档位定义）与 OB-53 第 ⑤ 项（链配置侧 `model: "auto"`）**不在本片范围**，
  仍为 `待议`，由下一片处理。
- 本片 ⑦ 为降级模型（ADR-020）⇒ 审查独立性低于白名单模型，结论应并读该限制。
- ⑥ 的跳过依据为「纯文档、无代码语义」，由 ⑦ 在结论行中确认。

## 11. 候选 OB（本片未新增）

无新增 OB；本片收口既有 OB-63 / OB-66 / OB-70。

## 12. ⑩ 停点报告

- 执行至 ⑦ 终审；**停在提交授权之前**：本片不 commit、不 push、不推 gitee，等待用户同轮授权。
- **回写核对**（OB-63 口径，本片为首次适用）：验证报告 = 已回写（本文件）；`DEV_PLAN` = **不适用**
  （准则治理片，按切片定义与既有先例不写 `DEV_PLAN` 段落）；台账 = 已回写（§12.4 的 OB-63 / OB-66 / OB-70
  处置列与本片完成记录行）；ADR = 已回写（ADR-020）。

- **⑦ 复审结论**：`RE-REVIEW: PASS WITH FIXES`；残留 3 条 Low（留存 index 读数为中间态、v1.19 条目的归属措辞、上一片报告 §16 的指代）已在 ⑧b 闭合。

## 13. 产物清单

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-model-lists.md	repo
docs/validation/directive-model-lists_EN.md	repo
docs/validation/directive-cleanup-eval.md	repo
docs/validation/directive-cleanup-eval_EN.md	repo
tmp/dce/land.js	local-only
tmp/dce/peek.js	local-only
tmp/dce/gate-tree.txt	local-only
tmp/dce/gate-index.txt	local-only
tmp/dce/selftest.txt	local-only
tmp/dce/probe-id-rejections.txt	local-only
tmp/dce/gate-tree-first-fail.txt	local-only
tmp/dce/test.txt	local-only
```
