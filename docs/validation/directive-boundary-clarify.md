# 切片 `DIRECTIVE-BOUNDARY-CLARIFY` 验证报告（准则 v1.41）

> 切片：`DIRECTIVE-BOUNDARY-CLARIFY`；准则版本 v1.41；档位 L；本文件为 **⑧a 初稿**（⑨ / ⑩ 的读数与提交状态由 ⑧b 回填）。
> 基线：`d:\LZProjects\prfrail`；⑧a 启动时暂存集为 6 项路径、索引等于工作区（`git diff --stat` 无输出）、无未跟踪件。
> 双语：本文件为 CN 权威；`docs/validation/directive-boundary-clarify_EN.md` 为逐行同位镜像，两侧行数相等。

## 1 变更概要

⑧a 之前 6 项路径，⑧a 之后 9 项；本片 `tools/**` 零改动（`git diff --cached --numstat` 只含文档）。下表增删数为机械读数；两处新报告的增删数取其自身行数（`+279/-0`，两侧相等）。

| 路径 | + | − | 说明 |
|---|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 36 | 11 | v1.41 正文：7 项处置 + 三处边界条款 + 版本包 + 台账翻转 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 36 | 11 | `_EN` 严格同位 |
| `docs/t027/REMAINING_SLICES.md` | 20 | 0 | P16 状态行 |
| `docs/t027/REMAINING_SLICES_EN.md` | 20 | 0 | `_EN` 同位 |
| `docs/validation/directive-boundary-clarify.md` | 279 | 0 | 本文件 |
| `docs/validation/directive-boundary-clarify_EN.md` | 279 | 0 | 同位镜像 |
| `docs/validation/directive-gates-patch.md` | 1 | 1 | §8 行尾内嵌来源标注 |
| `docs/validation/directive-gates-patch_EN.md` | 1 | 1 | `_EN` 同位 |
| `docs/validation/evidence/DIRECTIVE-BOUNDARY-CLARIFY-freedom-list.md` | 53 | 0 | 执行自由度清单（证据树根级） |

**artifacts 围栏**（9 项路径，处置均为 `repo`）：

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-boundary-clarify.md	repo
docs/validation/directive-boundary-clarify_EN.md	repo
docs/validation/directive-gates-patch.md	repo
docs/validation/directive-gates-patch_EN.md	repo
docs/validation/evidence/DIRECTIVE-BOUNDARY-CLARIFY-freedom-list.md	repo
```

7 项落点与机械行号（CN 与 `_EN` 行号相同）：

| 项 | 处置 | 正文落点（`docs/DELIVERY_DIRECTIVE.md`） |
|---|---|---|
| OB-23 | 已处置 | §7.2.1 L661 独立段落；§7.8.2 L743 收敛信号；§4.2.1 L418 标题「六条」+ L425 第 6 条 |
| 组 1 / OB-13 | 已处置 | §5.1 L480 时序说明 |
| 组 2 / OB-25 | 已处置 | §3.4 L357 适用范围句 |
| 组 2 / OB-26 | 已处置 | §3.4.2 L369 基线定义；§12.5 L1077 新缺口行 |
| 组 3 / OB-44 | 已处置 | §7.3 L693 第 ③ 项（写盘边界） |
| 组 3 / OB-45 | 已处置 | §7.3 L693 第 ② 项（数据来源表述） |
| 组 4 / OB-50 | 已处置 | §2.5 L316 句末括注 |

台账与版本包（同一文件）：§12.4 七行处置列已翻转（OB-13 L945 / OB-23 L959 / OB-25 L963 / OB-26 L965 / OB-44 L984 / OB-45 L985 / OB-50 L990）；头部 L3 为 `版本：v1.41。日期：2026-10-20。`、§0.3 的 v1.41 变更行（L87）、附录 C.0ac 起于 L1796（交付物行 L1801 已含 `REMAINING_SLICES{,_EN}` 的 P16 状态行）。

## 2 执行流水线

| 环节 | 角色 / 载体 | 模型 | 结论 | 轮次 | 备注 |
|---|---|---|---|---|---|
| 步骤 0 | `prfrail-implementer` | 主控模型 | 完成 | 1 | P16 状态行 +20/-0；8.2 核对通过；v1.41 头部并入 ② 批（`G1-b` 联动） |
| ① 设计 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 完成 | 1 | 不拆片；OB-23 三形态对比后推荐形态 C；四组编号；发现原计划缺 §7.3 |
| 批量辅助 ① | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 完成 | 1 | OB 状态盘点：7 项全 `待议`，CN/EN 行号一一对应 |
| 批量辅助 ② | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 完成 | 1 | 节号与 OB 号引用核对：18 个节号全存在、无悬空节号 |
| ② 实施 | `prfrail-implementer` | 主控模型 | 完成 | 1 | 6 项路径、逐组门禁全绿 |
| ③ 测试 | `prfrail-tester` | 主控模型 | 完成 | 1 | 三作用域、自测、六形状、判据逐条 |
| ④ 门禁 | 主控 | 主控模型 | 完成 | 1 | 三作用域 `TOTAL_FAIL=0` |
| ⑤ 预审 r1 | `prfrail-tester` | `GLM-5.3 (glm)` | PASS WITH FIXES | 1 | 1 Medium（C.0ac 交付物列漏列 P16）+ 5 Low |
| ② 整改 | `prfrail-implementer` | 主控模型 | 部分完成 | 1 | 派单中途网络中断，改动部分落盘；主控机械核验后继续，FIX-1 至 FIX-4 落地 |
| ⑤ 预审 r2 | `prfrail-tester` | `GLM-5.3 (glm)` | PASS | 1（收窄） | Medium CLOSED、3 项 Low CLOSED、FIX-5 未整改理由成立 |
| ⑥ 独立盲审 | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | PASS | 1 | 盲输入，零 Medium+ |
| ⑦ 终审 r1 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | PASS | 1 | 五组逐组通过，零 Critical / High / Medium / Low |
| 批量辅助 ③ | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | 完成 | 1 | 双语语言侧比对 6 处：否定极性、数量词、术语均无差异 |
| 批量辅助 ④ | 未派 | 无 | 不适用 | 0 | MANIFEST 需逐件 sha256，批量辅助无写权限且无法可靠计算哈希 |
| ⑧a 初稿 | `prfrail-documenter` | 主控指定 | 完成 | 1 | 本文件、`_EN` 镜像与自由度清单 |

- 返工口径：② 共 1 次主派单 + 1 次整改派单（中途中断），未重派、未回滚。
- 耗时字段历史缺失：本片未落盘起止时间戳，不补造。
- ④ 的三作用域读数为独立运行，不复用同一进程内的结果。

## 3 异常与兜底记录

| # | 事件 | 影响 | 处置 |
|---|---|---|---|
| D1 | 步骤 0 的 8.3「v1.41 头部」不能单独落（`G1-b` 要求头部版本等于 §0.3 末行版本） | 机械约束，非自由度 | 并入 ② 主变更批 |
| D2 | ① 的落点清单缺 §7.3（起跑消息未列） | 审查层写盘边界与数据来源两件事无落点 | 经用户裁定 B 扩入，范围只限该两件事 |
| D3 | ② 整改派单中途网络中断（客户端 `UND_ERR_SOCKET`） | 已落盘 FIX-1 至 FIX-4 的部分改动 | 主控先机械核验磁盘状态再继续；六形状仍全等、`G3` 仍 PASS |
| D4 | §7.2.1 的口径段由列表项行尾拆出为独立段落 | 文本形态变化、语义不变 | ⑤ r1 Low 的整改，同时落实用户裁定 A「两件事分开写」 |
| D5 | 批量辅助优先级：用户本轮回执更正为优先 `GLM-5.3-FlashX` | 前两次已用 `GLM-5.3-Flash` | 不重跑（产出可用）；第三次起按更正后优先级 |
| D6 | 批量辅助 ④（MANIFEST 草案）未派 | 无 | 原因见 §2 表，非「无适用工作量」 |
| D7 | `G9-e` 对 OB-44 现象列的历史路径 `tmp/v113/review-r2.md` 出 INFO | 该路径是历史事实引用、非证据工件 | 不阻断 |
| D8 | ⑤ r2 指出证据包缺陷：`review-5r2/01-fix-lines.txt` 的 FIX-2 EN 对照显示 `<missing>` | 取证脚本正则匹配不到 `Skippable` | 证据包缺陷，非文档缺陷（⑤ 已直读原文确认） |
| D9 | ⑤ 的 1 条 `[Info]`：EN 的 v1.41 变更行用 ①–⑦ 编号而旧行用 (1)–(6) | 风格差异 | 不处置 |
| D10 | 起跑消息的 OB-25 / OB-26 / C.0ac 起行号为 961 / 963 / 1794，机械读数为 963 / 965 / 1796 | 引用偏移，均偏 2 | 以机械读数为准，已在 §1 载明 |
| D11 | ⑧a 无法在会话内重放 transcript 的 `runSubagent` 解析：`main.jsonl` 仅 264 字节，session store 的 `assistant_response` 截断至 5000 字节 | §10 的成本表为派单事实转录，非机械审计 | 如实登记；机械审计由 ⑧b 回填，本表不作最终成本事实 |
| D12 | FIX-5 家族（`_EN` 单侧 4 处）保留未整改 | EN 与 CN 权威表述有 4 处不严丝合缝 | 单侧修改会触发 `G3(a.1)` FAIL，留待同时触及两侧的切片配对落地 |
| D13 | 自由度清单产出时点晚于 ⑦ r1（§1.7.3 要求 ⑥ 之后、⑦ 之前） | ⑦ r1 的输入包不含清单（`tmp/DIRECTIVE-BOUNDARY-CLARIFY/review-7/` 无清单文件） | 与 `GATES-PATCH` 的 D10 同型，本片为第 2 次；如实登记，不重建该包 |
| D14 | ③ 派单顺序偏离：本片 ③（`prfrail-tester`）在 ⑤ r1/r2、⑥、⑦ r1 **之后**才补派（transcript 偏移 `29098208`）；补派后产出全部通过 | 见下方 D14-1 / D14-2 / D14-3 | 记录在案（**OB-107**） |
| D15 | `--scope=ci` 的本地语义：`ctx.js` 取 `baseRef=HEAD^1` | **提交前它度量上一片**（本片读数 `changed=4`），不是本次暂存集；本片 9 文件的判据由 `index` / `tree` 承担 | `ci` 只在推送后对本次提交执行 |
| D16 | ⑧a 的 `G4b` 阻塞：首稿措辞含白名单外的未登记字，`G4b` 报 `unknown` 非空 | `G4b` 在 `--scope=tree` 下判红，阻断 ⑧a 交稿 | 按本仓口径**改写规避**（换成白名单内的同义词）、**未扩白名单**（`tools/gates/cjk-newwords.txt` 仍 37 条、`tools` 目录零改动），改后 `newCjk=591 unknown=[]` |
| D17 | ⑧a 的成本表缺机械依据：⑧a 无法在会话内重放 transcript 解析，其 §10 当时逐行标注「待 ⑧b」 | §10 的成本读数为派单事实转录，非机械审计 | 本阶段已按 transcript 机械审计回填（见 §10）；此为本仓 2026-10-20 固化的「⑧a 自查三条」中第 ① 条的首个实证缺口 |
| D18 | 起跑消息的行号偏移：⑧a 收到的事实清单里 OB-25 / OB-26 / C.0ac 行号为 961 / 963 / 1794，机械实测为 963 / 965 / 1796（均偏 2，因 §12.4 在登记 OB-106 之前已插入 OB-104 / OB-105 等行所致） | 引用偏移，不影响结论 | 报告一律用机械读数（见 §1） |
| D19 | FIX-5 家族（EN 单侧残留，未整改）：`docs/DELIVERY_DIRECTIVE_EN.md` 的 **L87 / L693 / L985**（`may only read` 应为 `may only be written as`）与 **L743**（§7.8.2 (i) 档 `must carry` 未显式承载完成体「须**已**给出」） | EN 与 CN 权威表述有 4 处不严丝合缝 | 留待后续同时触及两侧的切片一次性配对落地；理由：EN 单侧改动会使 `G3(a.1)` 的 CN / `_EN` 增减行数不对称而判红，而 CN 侧本身正确、不宜为形状门做无语义补偿改动 |
| D20 | 批量辅助优先级更正：用户在本轮回执中把优先级更正为「一律优先 `GLM-5.3-FlashX`，仅在耗尽 / 不宜 / 撤回时回退 Flash」 | 本片前两次（候选 ① / ②）已用 Flash | 产出可用**不重跑**；第三次起（候选 ③）用 FlashX |

- **D14-1 事实**：③ 于 ⑤ r1/r2、⑥、⑦ r1 之后补派（偏移 `29098208`）；补派后产出全部通过，并独立发现 3 项观察（`--check=base` 退出码 2 / `--scope=ci` 的 `HEAD^1..HEAD` 语义 / 静态 `check(` 呼叫点 408 与运行时 588 的口径差）。
- **D14-2 实质影响**：**⑥ 与 ⑦ r1 的审查输入包不含 ③ 的测试结果**——它们审的是「实现」，审查对象里没有「测试」；③ 事后独立发现的 3 项观察本应在 ⑥/⑦ 之前进入审查范围。⇒ **⑥ 的 `PASS` 与 ⑦ r1 的 `PASS` 建立在不完整输入上**，该缺口如实登记，**不得**以「补派后全部通过」掩盖。
- **D14-3 残余**：③ 的产出未进入任何审查环节的输入 ⇒ 本片「实现 + 测试」的双向审查链**只在本片报告内闭合，不在 ⑥/⑦ 的结论依据内**；已登记 **OB-107**（观察中）。

## 4 环境与档位

| 项 | 值 |
|---|---|
| OS | Windows（本机 `d:\LZProjects\prfrail`） |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| 档位 | L |
| ⑥ | 必派（触 §7.2.1 与 §7.8 的审查口径语义） |
| ⑦ 额度 | 1+3 |
| ⑤ 额度 | 1+1 |
| 批量辅助 | ADR-027（限四类任务） |
| 探针 | 0 |

## 5 审查结论摘要（⑤⑥⑦）

| 轮次 | 编号 | 级别 | 结论 / 处置 | 处置 |
|---|---|---|---|---|
| ⑤ r1 | M1 | Medium | C.0ac 交付物列漏列 `docs/t027/REMAINING_SLICES{,_EN}.md` 的 P16 状态行（暂存集含该对 +20/-0，属 under-report） | 已整改 |
| ⑤ r1 | L1–L4 | Low | FIX-1 至 FIX-4：内容侧整改项（行文与落点修正） | 已整改 |
| ⑤ r1 | L5 | Low | FIX-5 家族：`_EN` 单侧 4 处与 CN 权威不严丝合缝 | 记录在案 |
| ⑤ r2 | — | — | Medium CLOSED；3 项 Low CLOSED；FIX-5 未整改理由成立 | 已闭合 |
| ⑥ r1 | — | — | 盲输入、零 Medium+ | 记录在案 |
| ⑦ r1 | — | — | 五组逐组通过、零 Critical / High / Medium / Low | 记录在案 |
| 批量辅助 ③ | — | — | 双语语言侧 6 处比对无实质差异；2 处轻微 EN 差异并入 FIX-5 家族 | 记录在案 |

- ⑤ r1 的 Medium 闭合有机械复核：`git diff --cached --numstat -- docs/t027/` 为两文件各 `20 0`，且 C.0ac 交付物行 L1801 已含该对路径。
- FIX-5 家族的 4 处（`_EN` 单侧）：L87 / L693 / L985 的 `may only read` 应为 `may only be written as`；L743 的 `must carry` 未显式承载完成体「须已给出」。四处均需 EN 单侧改动，`G3(a.1)` 会判 FAIL；CN 侧本身正确，补偿性改动会污染冻结正文 ⇒ 留待配对片。
- ⑦ 的其他观察（无发现）：已闭环报告的单处 diff 为仅追加、未触结论行与读数；新增规则内部一致性（三档与 §7.2.1 ③ 解耦、§7.3 与相邻条款不冲突、§3.4.2 与 §3.4.3 及 §12.5 一致）均通过。

## 6 审查输入包摘要（含盲审隔离证明）

| 包 | 位置 | 内容 | 清单 |
|---|---|---|---|
| ⑤ r1 | `tmp/DIRECTIVE-BOUNDARY-CLARIFY/review-5/` | brief / 设计 / 全量 diff / numstat / 三作用域门禁 / 自测尾部 / 判据文本 / 报告 diff | `MANIFEST.md` |
| ⑤ r2 | `tmp/DIRECTIVE-BOUNDARY-CLARIFY/review-5r2/` | brief / 整改 hunk / 三作用域门禁 | `MANIFEST.md` |
| ⑥ r1 | `tmp/DIRECTIVE-BOUNDARY-CLARIFY/review-6/` | brief / 全量 diff / numstat / 状态 / 三作用域门禁 / 自测 / 读数 | `MANIFEST.md` |
| ⑦ r1 | `tmp/DIRECTIVE-BOUNDARY-CLARIFY/review-7/` | brief / 设计 / 全量 diff / numstat / 状态 / 三作用域门禁 / 自测 / 台账项 / 节号 / 报告单行 diff / 读数 | `MANIFEST.md` |

- ⑥ 与 ⑦ 的 brief 均按盲输入派单：未附 ⑤ 结论、未附任何前轮发现。
- 输入消毒：包内只有切片定义、设计、测试报告、diff、numstat、门禁读数、自测输出与判据源码。
- ⑦ r1 的输入包未含自由度清单（见 §3 的 D13），故 ⑦ 未据清单审查。
- 包内文件计数取自机械枚举：⑤ r1 为 11 项、⑤ r2 为 6 项、⑥ 为 10 项、⑦ 为 14 项（含各自 `MANIFEST.md`）。

## 7 可证伪性

### 7.1 可证伪断言（本片为纯文档片，无运行时探针）

| # | 断言 | 反例（删去后应有哪条读数变化） | 机械检查 |
|---|---|---|---|
| E1 | §7.3 的写盘边界条款存在 | 删去第 ③ 项后，OB-44 无承载句，⑤ 的逐项核对无法通过 | 检索 `不得创建、修改、删除任何文件` 应得 1 处 |
| E2 | §7.2.1 的口径段为独立段落 | 退回列表项行尾后，用户裁定 A 未落实，⑤ r1 Low 未闭合 | 检索段落首句 `跳过面证据充分性的口径` 应得 1 处 |
| E3 | §4.2.1 标题与条目数一致 | 标题退回「五条」时与第 6 条矛盾 | 标题与第 6 条同文件并存 |
| E4 | §12.4 七行处置列已翻转 | 任一行为 `待议` 时与 §0.3 的 v1.41 声明不一致 | 逐行读取七行处置列 |
| E5 | §3.4 的适用范围句只写一处 | 若在 §1.6 复制第二份，两处将可漂移 | 检索 `[ROUTINE]` 档由主控亲自执行 应得 1 处 |

### 7.2 机械读数复算

| 读数 | 值 | 来源 |
|---|---|---|
| `git diff --cached --numstat` | 6 项路径（⑧a 前） | 命令输出 |
| CN 与 `_EN` 行数 | lines=1876 两侧相等 | 机械计数 |
| 六形状 | heads 123 / bold 4554 / bars 3007 / blank 453 / ticks 4636 | 机械计数，两侧全等 |
| `SELFTEST` | PASS 588/588 | `node tools/gates/selftest.js` |
| `G4b` | newCjk=591、unknown=[] | `--check=G4b --scope=tree` |

### 7.3 变异载荷

- 本片 `tools/**` 零改动（`git diff --cached --numstat` 不含 `tools/`），无判据代码与测试代码变更 ⇒ 无变异载荷。
- 本片未新增断言：`SELFTEST` 维持 588/588，与 `GATES-PATCH` 收尾值相同。

### 7.4 覆盖边界与已知 Low

| 编号 | 内容 | 归属 |
|---|---|---|
| L-1 | FIX-5 家族的 4 处 EN 单侧残留 | 记录在案（见 §11） |
| L-2 | D13 的清单时序偏离 | 记录在案（见 §3） |
| L-3 | D10 的引用行号偏移已按机械读数更正 | 已处置 |

- 本报告的结论不覆盖 `_EN` 的语义等价性；该面由 ⑤ / ⑦ 与 `G3(a.2)` 的形状比较共同承担。

## 8 门禁结果

| 命令 | 结果 |
|---|---|
| `node tools/gates/selftest.js` | `SELFTEST: PASS 588/588` |
| `node tools/gates/gate.js --all --scope=index` | `TOTAL_FAIL=0`（exit 0） |
| `node tools/gates/gate.js --all --scope=ci` | `TOTAL_FAIL=0`（exit 0，`changed=4` 为上一片提交；语义见 §9） |
| `node tools/gates/gate.js --all --scope=tree` | `TOTAL_FAIL=0`（exit 0） |
| `--check=G4b --scope=tree` | `newCjk=613 unknown=[]`（`whitelist=37`） |
| `G5(d)`（`--scope=index`） | `9 path(s) declared`；`declared-not-touched=10`（INFO，围栏为落点快照）；`addMode=explicit` |
| 基础 `go test ./...` | 由 `--all` 内嵌执行 |

§8 的数列原为 ⑧a 时点读数，现替换为定稿读数；定稿读数取自主控在 ⑧b 之后、⑧c 之前的实测。

逐判据自检读数（定稿，`--scope=tree`）：

| 判据 | 读数 |
|---|---|
| G7-a | PASS |
| G8-a | PASS |
| G1-a | PASS |
| G2 | PASS |
| G6-4 | PASS |

- `G7-b` 的 3 项新增路径未入索引：该 FAIL 已在主控暂存 9 个路径后消解，定稿 `TOTAL_FAIL=0`；`G4b` 在 `--scope=tree` 下的未登记字已按改写规避清除（见 D8）。

**CN 与 `_EN` 六形状（`docs/DELIVERY_DIRECTIVE.md` 对 `_EN`）**：

| 形状 | CN | EN |
|---|---|---|
| 行数 | 1876 | 1876 |
| 标题 | 123 | 123 |
| `**` | 4554 | 4554 |
| `|` | 3007 | 3007 |
| 空行 | 453 | 453 |
| 反引号 | 4636 | 4636 |

> 上表为 ⑧a 时点读数；其用途是证明 CN 与 `_EN` 的六形状同位，不是当前绝对值的保证。
- 两侧均为 UTF-8 with BOM 与 LF。

## 9 ⑨ 原生验证与提交状态

⑨ 由 ⑧b 执行并回填；本节读数全部为机械读数。

**逐判据退出码（`--check=<ID> --scope=index`，18 项）**：全部 exit 0。

| 判据 | 退出码 | 备注 |
|---|---|---|
| `G1-a` / `G1-b` / `G1` / `G2` / `G3-a` / `G3-b` | 0 | PASS |
| `G4a` / `G4b` / `G5-a` / `G5-b` / `G5-c` | 0 | PASS |
| `G5-d` | 0 | INFO（`no changed .github/agents/*.agent.md in scope`） |
| `G5` / `R2.5` / `G6` / `G7` / `G8` / `G9` | 0 | PASS |
| `base` | 2 | `unknown criterion id: base` ⇒ **OB-105 的第二次实证**（v1.40 已登记，本片不修） |

两点须与上表连读：

- `G5-d` 的 exit 0 携带 INFO 而非 PASS：本片变更集不含 `.github/agents/*.agent.md`，该判据无对象可判。
- `base` 不计入 18 项：`--list` 把它列为第 19 个 id，`--check` 不认它；`--check` 一侧的 18 项全为 exit 0。

**三作用域 `--all`**：`index` / `tree` / `ci` 均 `TOTAL_FAIL=0`、exit 0。**须连读的限定**：`--scope=ci` 的实现基准是 `HEAD^1..HEAD`（`ctx.js` 的 `baseRef=HEAD^1`），因此**本地提交前**它度量的是**上一片**的提交（本片读数 `changed=4`）；本片 9 文件的判据由 `index` / `tree` 承担，`ci` 只在推送后由 CI 对本次提交执行。

**干净 worktree 重放**（`HEAD` = `29309d3`）：`git worktree add --detach tmp/DIRECTIVE-BOUNDARY-CLARIFY/pristine HEAD` → 导出 `git diff --cached --binary` 补丁 → 在 pristine 内 `git apply --index` → `--all --scope=index` 得 `TOTAL_FAIL=0`、`selftest` 得 `588/588` ⇒ 变更集自洽；随后 `git worktree remove --force` 并删除补丁，`git worktree list` 仅剩主工作区。

**全树编码审计**：`scanned=3332` 个 tracked 文本文件，证据树外违规 0；证据树内 16 处历史 BOM / EOL 偏差（历史归档包与 BOM 兼容性 fixture）与本片无关、未计入。

**引用可解析性**：暂存集新增行含 34 个节号引用，悬空 0。

| 项 | 状态 |
|---|---|
| `selftest.js` | `SELFTEST: PASS 588/588` |
| ⑩ 停点回执 | 未执行（本片 ⑩ 停点等授权） |
| 提交与推送 | 均未执行；未触 gitee；CI 未触发 |

## 10 成本与计量（R7.1）

**机械审计依据**：下表为本片 `runSubagent` 调用的 transcript 解析结果，按发生序排列；「#」为本会话内的调用序号，偏移与行号为 transcript JSONL 的字节偏移与行号。

| # | 环节 | `agentName` | `model` | transcript 偏移 | 行号 |
|---|---|---|---|---|---|
| 104 | ② 步骤 0（P16 状态行） | `prfrail-implementer` | 主控同模型（未指定） | 26629059 | 28405 |
| 105 | 批量辅助 ①（7 项 OB 状态盘点） | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 26709083 | 28504 |
| 106 | 批量辅助 ②（节号 / OB 号引用核对） | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 26890115 | 28565 |
| 107 | ① 设计 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 27060493 | 28663 |
| 108 | ② 实施（7 项 + 版本包） | `prfrail-implementer` | 主控同模型 | 27293578 | 28777 |
| 109 | ⑤ r1 | `prfrail-tester` | `GLM-5.3 (glm)` | 27639436 | 29090 |
| 110 | ② 整改（FIX-1~FIX-4；派单中途网络中断） | `prfrail-implementer` | 主控同模型 | 27823412 | 29234 |
| 111 | ⑤ r2（窄化） | `prfrail-tester` | `GLM-5.3 (glm)` | 28001752 | 29437 |
| 112 | ⑥ 盲审 | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | 28227831 | 29536 |
| 113 | ⑦ r1 终审 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 28284792 | 29616 |
| 114 | 批量辅助 ③（五处双语语言侧检查） | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | 28408082 | 29889 |
| 115 | ⑧a | `prfrail-documenter` | 主控同模型 | 28493087 | 29945 |
| 116 | ⑧（G4b 改写规避） | `prfrail-documenter` | 主控同模型 | 28922118 | 30385 |
| 117 | ③ 测试（在 ⑤ 之后补派 —— 见 D14） | `prfrail-tester` | 主控同模型 | 29098208 | 30598 |

**计量结论**：① 1 / ② 3 / ③ 1 / ⑤ 2 / ⑥ 1 / ⑦ 1 / ⑧ 3 / 批量辅助 `GLM-5.3-Flash` 2 次 / 批量辅助 `GLM-5.3-FlashX` 1 次 / 探针 0。

- ② 计 3 次 = #104 + #108 + #110；⑧ 计 3 次 = #115 + #116 + 本次 ⑧b（本次调用晚于上表的审计时点，故不在表中）。
- 两枚 GLM 模型分别计量（`GLM-5.3-Flash` 2 次、`GLM-5.3-FlashX` 1 次），不合并计数。
- 上表取代 ⑧a 的转录版，D11 与 D17 的机械依据缺口由本轮闭合。

**批量辅助四个候选点逐项登记**：

- 候选 ①（七项 OB 逐项状态盘点）= **已派**（`GLM-5.3-Flash`，偏移 26709083）：读数 = 7 项全 `待议` / `To be discussed`，CN 与 EN 行号一一对应（OB-13 L934、OB-23 L948、OB-25 L952、OB-26 L954、OB-44 L973、OB-45 L974、OB-50 L979，均为本片改动前的台账行号）。
- 候选 ②（准则节号 / OB 号引用核对）= **已派**（`GLM-5.3-Flash`，偏移 26890115）：读数 = 18 个节号全部真实存在、CN / EN 行号相同；7 项现象列引用的节号无悬空；`OB-106` 当时未登记（本片 ⑧b 才登记）；G6-5 正文与 `g6.js` 实现已一致。
- 候选 ③（五处双语对位语言侧检查）= **已派**（`GLM-5.3-FlashX`，偏移 28408082）：读数 = 6 处比对，否定极性 0 差异、数量词 / 档位词 0 差异、术语 0 差异；2 处轻微 EN 侧差异（§7.8.2 (i) 档完成体、§7.3 ② 动词）并入 FIX-5 家族。
- 候选 ④（审查包 MANIFEST 清单草案）= **未派**：MANIFEST 的每行是 `bytes + sha256`，必须由能对文件逐件计算哈希并写盘的进程产生；批量辅助角色（`deep-reasoner` 载体）持只读工具、无写权限且无法可靠计算或回读哈希 ⇒ 该任务不适合批量辅助，改由主控用 `node` 脚本机械生成（本片 `review-5` / `review-5r2` / `review-6` / `review-7` 四个包均已生成 `MANIFEST.md`）。

本轮不给 credits 数值，只给可核的模型与轮次事实。

## 11 明确未执行事项

- 未执行 `git commit` / `git push`；未触 gitee；CI 未触发。
- ⑩ 停点回执未执行（本片 ⑩ 停点等授权）。
- 本片未新增断言：`tools/` 目录零改动 ⇒ `SELFTEST` 仍 588/588。
- 未改 §2.2 / §2.4 的本体闭集、§5.2 矩阵、§7.5 数值、§6.3 判据本体。
- 本阶段未整改 FIX-5 家族的 4 处 `_EN` 单侧残留（理由见 §5 与 D19）。
- 批量辅助 ④（MANIFEST 草案）未派（理由见 §10）。

## 12 下一步建议

1. ⑩ 停点：等授权；主控出 ⑩ 回执（含派单登记表、批量辅助四项逐项、OB-23 裁定、四终读与 OB-106）。
2. 用户授权后提交并推送，观察 CI 双绿。
3. FIX-5 家族在下一处同时触及 CN 与 `_EN` 的切片内配对落地。

## 13 执行自由度

逐条引用 `docs/validation/evidence/DIRECTIVE-BOUNDARY-CLARIFY-freedom-list.md`（清单为权威，本表不改变其含义）。

| # | 选择 | 落点 |
|---|---|---|
| F1 | OB-23 取形态 C（等价覆盖 + 显式声明 + 三档等价性，且 (i) 档不免论证） | §7.2.1 L661 |
| F2 | §7.3 扩到 ⑤ / ⑥ / ⑦，且只限「写盘边界」与「数据来源」两件事 | §7.3 L693 |
| F3 | 口径段以独立段落承载，不写成列表项行尾的括注 | §7.2.1 L661 |
| F4 | OB-13 的 N 参数化，不写死具体数值 | §5.1 L480 |
| F5 | `[ROUTINE]` 档由主控亲自执行诊断扫描，不派子代理、成本不增 | §3.4 L357 |
| F6 | OB-26 取保守口径：无快照时全量视为新增 | §3.4.2 L369 |
| F7 | 适用范围只在 §3.4 声明一处，不复制到 §1.6 | §3.4 L357 |
| F8 | OB-50 零新正文，只做承载标注 | §2.5 L316 |
| F9 | C.0ac 交付物列补列 P16 状态行（⑤ r1 Medium 的整改形态） | C.0ac L1801 |

阻塞项已消解 — `G4b` 的 `unknown` 现为空（本轮改写规避，白名单零新增），本稿以 `DOCS: DONE` 交付；`tools/gates/cjk-newwords.txt` 未改动。

DOCS: DONE
