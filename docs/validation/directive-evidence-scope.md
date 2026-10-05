# 切片验证报告 · `DIRECTIVE-EVIDENCE-SCOPE`（v1.24）

> 类型：准则治理类（证据树判据作用域排除 + 新判据 G8 四子查；触判据本体）。分级 `[SLICE]`。规模 M+。探针 0。
> **授权边界**：commit 否 / push 否 / gitee 否 / 探针 0。

## 1. 目标与范围

- 主体 **OB-79**：§6.3 的 G1-a / G2 / G4a / G4b / G6 目标文件集排除 `docs/validation/evidence/**`（与 G7 既有的排除**同形、单源**定义于实现层），该树内变更改由新判据 **G8** 四子查承担。
- 同批 **OB-80**（归档载荷编码形态口径）、**OB-76**（§1.5.1 第四类「重建归档」）、**OB-78**（§1.5.1 的 ③ 例示修正）。
- **OB-77 仅评估**：值得另立片、低优先级，判据形态 = 变更行 `tmp/<path>` 引用的软扫描；本片不落判据，§12.4 状态词保持 `待议`。

## 2. 变更概要

- 条文（②a 与补丁轮）：`docs/DELIVERY_DIRECTIVE{,_EN}.md`；落 §6.3 的 G8 行、判据作用域排除段、G8 已知边界段、脚本化状态段，以及 §1.5.1 的 ② 与 ④ 与指针同步。
- 实现（②b）：`tools/gates/lib/ctx.js`（含 DEF-2 的批量忽略与记忆化）、`tools/gates/lib/criteria/g8.js`（新增）、`tools/gates/lib/criteria/g4a.js`、`tools/gates/lib/criteria/g7.js`、`tools/gates/gate.js`（含 DEF-1）、`tools/gates/evidence-forms.txt`（新增，3 条登记）。
- 测试（③）：`tools/gates/selftest.js` 扩至 T43–T60。
- 夹具：`tools/gates/testdata/fixtures/*.json` 新增 28 枚（磁盘名一律不以 `.md` 结尾）。
- 证据补录（①）：`docs/validation/evidence/directive-history/tmp-lifecycle-workdir/` 的 `MANIFEST.md` 与 `SHA256SUMS.txt` 补录；载荷字节零改动、历史不追溯。
- **暂存状态**：⑦ 复审时点 = 条文两文件已暂存、实现与夹具未暂存（⑦ 复审/终审的 Low 流程注记）；⑧a 落盘后为使 §6.3 的 `G7-b` tree 预检忠实成立，全部声明工件已逐文件精确 `git add`（未 commit、未 push、未推 gitee）。
- **提交流程注记**（⑦ 复审/终审的 Low）：提交时按 §6.3 的 G5 逐文件精确 `git add` 并复核 staged ≡ 声明集合。

## 3. 执行流水线与环境

① 判定（主控 + V4 Pro）⇒ ②a 规则正文（主控）⇒ ②b 实现（`prfrail-implementer`）⇒ ③ 测试（`prfrail-tester`）⇒ 主控编排（证据补录）⇒ ④ 门禁 ⇒ ⑤ 预审（V4 Pro）⇒ ⑥ 独立扫描（MAI）⇒ ⑦ 盲审（`GPT-5.3-Codex`）⇒ 整改（终审四条文本精度同步落地）⇒ ⑦ 复审 ⇒ ⑨ 原生验证 ⇒ ⑧a 本报告对 ⇒ ⑧b 定稿 ⇒ ⑧c 收尾清洁 ⇒ ⑩ 停点回执。

环境：Windows 11 / Node 24；仓库 `github.com/larsonzh/prfrail`；基线 HEAD = `0c12abc`（本片未提交）。
档位：主控 High；①⑤ V4 Pro Max；② `prfrail-implementer` / ③ `prfrail-tester` = V4.1 Flash；⑥ MAI-Code-1.1-Flash；⑦ `GPT-5.3-Codex` Extra High。

## 4. 异常与兜底记录

- **DEF-1**（既有缺陷、非本片引入、本片修复）：`gate.js` 的 `run()` 在 catch 中提前 `return` 跳过 stderr 写入 ⇒ 抛出型用法错误 exit 2 却零输出；已修（2 行）并加断言（T57 组，独立捕获 stderr）。
- **DEF-2**（本片自引入、本片修复）：tree / index 分支原为每路径一次 `git check-ignore`；单次 spawn 实测 417 ms，证据树标为 2671 项 ⇒ 单次枚举外推 1,802,368 ms（约 30 分钟），G8-a 祖先遍历外推约 74 分钟；改为单次 `check-ignore --stdin -z` 并加 `treeCache` 记忆化，修复后 1414–1626 ms（热调用 0.12 ms）。
- **F-1 ②b 子代理最终报告未回传**（停在等待后台命令）：验收以主控机械复核为准，不以实现者自报读数。
- **时序偏离**：⑧a 报告对未在 ⑨ 原生验证之前完成（实际序 = ⑦ 复审 → ⑨ → ⑧a）；依据 §5.1 的固定环节顺序，本项属**时序偏离**、如实登记。⑨ 的读数在 ⑧a 落盘前取得 ⇒ §7 表格已标注为 **⑨ 时点**读数（`changed=39` / `2` / `4`）；⑧c 后复跑复核为 `changed=44`、`TOTAL_FAIL=0`、`SELFTEST: PASS 328/328`。
- **自由度清单落盘时序偏离**：自由度清单未在 ⑥ 后 / ⑦ 前落盘，⑦ 三读均未复核该清单（随 ⑧a 落盘）——见 F-14；依据 §1.7.3 与 OB-48 先例，登记为本片偏差项。

## 5. 审查结论摘要（⑤⑥⑦）

- **⑤ 预审**（V4 Pro）：`PASS WITH FINDINGS`。M-1（G8-c 对零行映射节真空通过）裁定不修、入报告；M-2（B-1 强度代价未在正文命名）裁定走 (b)，已加 §6.3 第三点 (iii)；L-1–L-5、L-7 不修、逐条声明口径；L-6(i) 已修 `g8.js` 陈旧注释（非注释变更 0 行）；L-6(ii) 的 C.0n 状态行由本轮回写。
- **⑥ 独立扫描**（MAI）：`PASS WITH FINDINGS`，无阻断。独立复算补录算术 744,747 + 259,061 = 1,003,808；独立验证 `check-ignore --stdin -z` 的 exit 1（无命中）被正确当作数据（`if (!(r.ok || r.code === 1))`）；浅克隆 `--scope=ci` 守卫 fail-closed（stderr 非空、exit 2）。
- **⑦ 三读**（`GPT-5.3-Codex`，1 盲审 + 1 终审 + 1 复审）：盲审（不附任何清单）`BLOCK`，5 条发现（2 High：§0.3 行尾 `fixtures` 口径过宽、C.0n「不改任何历史行措辞」未限定列；3 Medium/Low：G8「只判变更行」口径偏差、G8-d 注记区未披露、暂存拆分流程注记）⇒ 终审 `PASS WITH FINDINGS`（条件通过），四条全接受并给出 CN / EN 定稿措辞、第 5 条作流程注记 ⇒ 四处已落地（§0.3 行尾 / C.0n 边界字段 / §6.3 的 G8 行尾「只判变更集（触及的包）」/ G8 行 ④ 加「头部注记区」）⇒ 复审 `PASS WITH FINDINGS`（独立核实四处落地、三读全绿、实现-正文一致），未发现不可交付事实。
- **落地时的机械裁定（如实载明）**：终审给出的 (2) 处 EN 措辞缺少与 CN `**现象列**` 配对的加粗 ⇒ 会被 G3(a.2) 的整文件 `**` 计数判红（现 3261 = 3261）；文档角色停在报告、未自行换词，由主控按镜像一致性的机械口径裁定取 A1（EN 加粗配平）且 EN 取 `_EN` 自身 §12.4 表列名 `Phenomenon`；EN 侧 G8 行 ④ 保留既有编号字形 `(iv)` 以保 G3(b) 的 glyph 计数。此为**主控机械裁定**，非终审原文。

## 6. 审查输入包摘要（含盲审隔离证明）

- 输入包 = `tmp/ges/REVIEW-BUNDLE-2.md`（终稿态，取代 ⑤ 时点版本）；⑦ 复核时按该包第 0 节的状态变化表校准。
- **自由度清单缺口**：审查输入包 = `tmp/ges/REVIEW-BUNDLE-2.md`；自由度清单未在 ⑥ 后 / ⑦ 前落盘（随 ⑧a 落盘），故 ⑦ 三读未复核该清单——见 F-14。
- **盲审隔离证明**：盲审首轮**不附任何清单**、只看条文与实现 diff；其抽中依据为 §7.7 比例抽查（窗口内触及面最广者），判定理由在 C.0n 留痕；本片 ⑦ = 3 次未超上限（1 盲审 + 1 终审 + 1 复审）。
- 消毒：审查输入只含仓库内路径与文本，无凭据、无外部连接材料。

## 7. 可证伪性与门禁结果

| 项 | 读数 |
|---|---|
| `--selftest` | `SELFTEST: PASS 328/328`（演进：②b 290 → ③ 326 → ③bis 327 → ③ter 328） |
| `--all --scope=tree`（⑨ 时点） | `changed=39`、`TOTAL_FAIL=0` |
| `--all --scope=index`（⑨ 时点） | `changed=2`、`TOTAL_FAIL=0` |
| `--all --scope=ci`（⑨ 时点） | `changed=4`、`TOTAL_FAIL=0` |
| OB-79 核心回归（`range:4284e94^..4284e94`） | `G1-a` / `G2` / `G4a` / `G4b` / `G6` 各 `TOTAL_FAIL=0`；`G4a` 明示 340 file(s) OK (312 excluded: evidence tree) |
| G8 真实语料（同范围） | `G8-a` = 1 红、`G8-b` = 1 红、`G8-c` = 绿（15 包）、`G8-d` = 绿（312 file(s)）；历史红不追溯（§11.1） |
| 补录后 | `--check=G8-b --scope=tree` = PASS，`1 pack(s) self-consistent` |
| DEF-1 | 复现 ⇒ `exit=2` 且 stderr 非空（含 `gate.js:` 前缀与原因）；已修并断言（T57 组） |
| DEF-2 | 修复前 417 ms × 标为 2671 项 ⇒ 外推 1,802,368 ms；修复后 1414–1626 ms（热 0.12 ms）；等价性 = 排序逐字相等、两向差集空；tree 分支 `check-ignore` 恰好 1 次；T60 证伪演示 `shipped=true / perFileMutant=false / doubleBatchMutant=false` |

上表为 ⑨ 时点读数；⑧a/⑧b 落盘报告对、自由度清单与台账行后，变更集增至 `changed=44`；⑧c 后实测：`--all --scope=tree` 与 `--all --scope=index` 均 `changed=44`、`TOTAL_FAIL=0`；`--selftest` 仍 `SELFTEST: PASS 328/328`；`tmp/` 内的 ③ 类过程件已清除（仅留 `.gitkeep`）。

**变异检验**：不适用于语义改判（本片改判据作用域与新增判据本体）；反证手段 = 28 枚夹具的 `--selftest` 断言、OB-79 核心回归、DEF-1 / DEF-2 的复现与修复对照、以及 T60 的证伪演示。

## 8. 已知边界与残余风险

- 逐字引用 §6.3「G8 的已知边界（如实声明）」段的三句：**（i）G8 是形态级守护**——它只核「证据树内路径的形态与包内自洽」，不做内容语义判定；**（ii）「活代码伪装成冻结包」不由 G8 识别**——冻结包合法承载脚本类工件（先例：`docs/validation/evidence/b2-2026-09-17/machine-ops/` 的 `*.ps1`），该残余风险由 ⑤⑥⑦ 审查链承担，不得声称 G8 覆盖。**（iii）载荷的编码形态不由 G8 核验**——G4a 已排除该树，G8-c 只核 `MANIFEST` 已记录的「后缀与编码映射」行；**无该节的包**（遗留包）其载荷编码形态无判据覆盖，该残余强度缺口由 ⑤⑥⑦ 审查链承担。
- **B-1**：无映射节的包其载荷编码形态无判据覆盖，已由 (iii) 声明。
- **B-2**：DEF-2 由 T60 守护；`index` / `ci` 为同一代码路径的传递性覆盖。
- **B-3**：G8-b 不校验 MANIFEST「字节」列；`SHA256SUMS` 的部分损坏形态未单独断言。
- **B-4**：形态级守护、不识别活代码伪装、不判内容语义。
- **B-5**：历史红不追溯（`4284e94` 的 G8-a / G8-b 两条永久保留）。
- **B-6**：M-1——零行映射节 = 声明的零归一，与「无节」同为无记录义务。
- **B-7**：未触及包不判；嵌套按「最近祖先」；G8-d 只扫注记区首 15 条；重建归档 sha256 仅核形态。
- ⑥ 两条 Low：缓存键依赖 repo 相对规范路径；G8 形态级边界。
- ⑦ 两条 Low：双点排除守卫的维护耦合；大历史范围读数耗时。

## 9. 明确未执行事项

- **M-1 不修**（理由：识别「显式声明」属内容语义，超出形态级守护）。
- **L-1 不修**（G8-c 只认小写 64 位 hex = 未声明的加严）、**L-2 不修**（未入库条目须 `/` 绝对路径 = 条文外加严）、**L-3 不修**（G8-d 不核日期、按「出现」触发）、**L-4 不修**（登记表「仅根级」靠人工纪律）、**L-5 不修**（tree / range 载荷枚举对 symlink、tracked-but-ignored 不对称）、**L-7 不修**（MANIFEST 重复条目被去重）。
- **OB-77 不落判据**；**历史不追溯**；**legacy 包不修**；**提交授权：本片未获授权**（未 commit、未 push、未推 gitee）；**探针 0**。

## 10. 偏差与自由度

- 与自由度清单同源，逐条见 `docs/validation/evidence/DIRECTIVE-EVIDENCE-SCOPE-freedom-list.md`（F-1 至 F-13）。
- F-1 ②b 子代理最终报告未回传；F-2 ① 补录动了两处元数据；F-3 DEF-1；F-4 DEF-2（本片自引入）；F-5 登记表按端点读取；F-6 遗留单位实测更正；F-7 §6.3 措辞无需改写 + 四项同步复核均判不需落；F-8 ②a 时点 G8 实现尚不存在；F-9 成本向量不入 §0.3 与 C.0n；F-10 ⑤ 发现清单及处置；F-11 B-2 断言与最后一段盲区；F-12 ⑦ 三读引出的四处文本精度同步；F-13 文档角色拦下 G3 冲突 + 主控机械裁定。

## 11. ⑧c 收尾清洁记录

- 本片 ⑧c 由主控在**本报告定稿之后**执行：`tmp/ges/` 内的审查过程件（`REVIEW-BUNDLE-2.md`、各读数与脚本）属 §1.5.1 的 **③ 类**一次性过程产物，用完即时清除、不入库。
- 冻结证据包不参与逐文件编码改写（§5.1 既有豁免）；证据树的编码可核验性由 `MANIFEST` 两枚 sha256 记录承担。

## 12. 成本与计量

- 实得向量：①×1 / ②×3 / ③×3 / ⑤×1 / ⑥×1 / ⑦×3（盲审 1 + 终审 1 + 复审 1）/ ⑧×10 + 主控若干轮。
- 口径：①×1 仅指架构师方案轮（`deep-reasoner`，只读产方案书）；证据补录轮由 ⑧ 文档角色执行、计入 ⑧ 计数，不计入 ①。
- 口径：⑧ = 文档/报告角色的全部实际调用轮次，含未落盘的 BLOCKED 轮与因 `Agent error` 未产生改动的重试轮。
- 未超切片定义上限（①≤1、⑤≤1+1、⑥≤3、⑦ = 1 终审 + 1 复审 + 盲审 1）；探针 0；无 Haiku 调用；无额度例外、无降级 ADR。

## 13. OB-77 评估结论

- **OB-77**：值得另立片、低优先级；判据形态 = 变更行 `tmp/<path>` 引用的**软扫描**；本片不落该判据，§12.4 状态词保持 `待议`。

## 14. 后续建议（含档位建议）

- 后续建议（仅建议、不登记）：B-1 或需新判据（新 OB 候选）；B-2 盲区；**OB-82** 的 §7.7 窗口计数锚点；legacy 包补规范形态（另片）。
- 档位建议：后续判据本体变更片，⑦ 仍用 `GPT-5.3-Codex`、不降级。

## 15. ⑩ 停点报告（含回写核对行）

- 执行至 **⑨**；⑧b 定稿后由主控执行 **⑧c** 并出具 **⑩ 停点收据**；**停在提交授权之前** —— 提交授权：本片未获授权（未 commit、未 push、未推 gitee）。
- **回写核对**（§11.2 的 OB-63 口径）：验证报告 = **已回写**（本文件对）；`DEV_PLAN` = **不适用**（准则治理片）；台账 = **已回写**（`REMAINING_SLICES{,_EN}` 的完成行，只记报告路径）；ADR = **不适用**（本片无降级、无裁决 ADR）。
- ⑨ 读数：tree / index / ci 三作用域 `TOTAL_FAIL=0`；`SELFTEST: PASS 328/328`；`--check=G8-b --scope=tree` = PASS。

## 16. 产物清单

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-evidence-scope.md	repo
docs/validation/directive-evidence-scope_EN.md	repo
docs/validation/evidence/DIRECTIVE-EVIDENCE-SCOPE-freedom-list.md	repo
docs/validation/evidence/directive-history/tmp-lifecycle-workdir/MANIFEST.md	repo
docs/validation/evidence/directive-history/tmp-lifecycle-workdir/SHA256SUMS.txt	repo
tools/gates/gate.js	repo
tools/gates/lib/ctx.js	repo
tools/gates/lib/criteria/g4a.js	repo
tools/gates/lib/criteria/g7.js	repo
tools/gates/lib/criteria/g8.js	repo
tools/gates/evidence-forms.txt	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/exclusion-negative.json	repo
tools/gates/testdata/fixtures/exclusion-positive.json	repo
tools/gates/testdata/fixtures/g8-a-freedom-list.json	repo
tools/gates/testdata/fixtures/g8-a-member.json	repo
tools/gates/testdata/fixtures/g8-a-real-registered.json	repo
tools/gates/testdata/fixtures/g8-a-registered.json	repo
tools/gates/testdata/fixtures/g8-a-stray-root.json	repo
tools/gates/testdata/fixtures/g8-b-backslash-sums.json	repo
tools/gates/testdata/fixtures/g8-b-hash-mismatch.json	repo
tools/gates/testdata/fixtures/g8-b-malformed-sums.json	repo
tools/gates/testdata/fixtures/g8-b-manifest-truncated.json	repo
tools/gates/testdata/fixtures/g8-b-missing-list.json	repo
tools/gates/testdata/fixtures/g8-b-unlisted-payload.json	repo
tools/gates/testdata/fixtures/g8-b-untracked-conflict.json	repo
tools/gates/testdata/fixtures/g8-bom-in-tree.json	repo
tools/gates/testdata/fixtures/g8-bom-outside.json	repo
tools/gates/testdata/fixtures/g8-c-arch-not-in-payload.json	repo
tools/gates/testdata/fixtures/g8-c-arch-not-in-sums.json	repo
tools/gates/testdata/fixtures/g8-c-bad-hash.json	repo
tools/gates/testdata/fixtures/g8-c-no-section.json	repo
tools/gates/testdata/fixtures/g8-c-ok.json	repo
tools/gates/testdata/fixtures/g8-d-bad-sha.json	repo
tools/gates/testdata/fixtures/g8-d-body-only.json	repo
tools/gates/testdata/fixtures/g8-d-missing-basis.json	repo
tools/gates/testdata/fixtures/g8-d-ok.json	repo
tools/gates/testdata/fixtures/g8-ok-pack.json	repo
tools/gates/testdata/fixtures/g8-t50-in-tree.json	repo
tools/gates/testdata/fixtures/g8-t50-outside.json	repo
tmp/ges/REVIEW-BUNDLE-2.md	local-only
tmp/ges/REVIEW-BUNDLE.md	local-only
tmp/ges/	local-only
```

说明：证据树的冻结包以单元 `MANIFEST` 行声明（载荷逐条见各单元 `SHA256SUMS.txt`）；`local-only` 行为 `tmp/ges/` 的审查过程件，随 ⑧c 清除。
