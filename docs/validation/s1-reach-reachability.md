# S1-REACH — OB-21 可达性刷新与 OB-51/52/53 登记（验证报告）

> 切片 `S1-REACH`（业务判定类，非准则治理） | 分级 `[SLICE]` | 规模 S | 日期 2026-09-29。
> **状态：已停止——未完成**（⑩ 停点；未完成清单见 §13；commit / push 须同轮授权）。

[English](s1-reach-reachability_EN.md)

## 1. 变更概要（含执行自由度清单）

- **本片做了什么**：① 以 DR-1 闭环后的当前事实刷新 **OB-21 / T018 / S1 exit** 的**可达性判定**（逐门 × 依赖 × 解锁条件，见 §2）；② 只读确认 **BYOK / secretStorage 五点**（见 §3）；③ 把 **OB-51 / OB-52 / OB-53** 落 `docs/DELIVERY_DIRECTIVE.md` §12.4；④ 按 OB-37 元规则以 **ADR-015** 留痕 ⑦ 降级。
- **本片没做什么**：不改任何 Go 代码 / schema / CONTRACTS / fixtures / workflows；不做真探针（额度 0）；不改任何判据；**OB-21 这一行的历史句不改写**（指令 §11.1 的历史行口径；刷新结论落在本报告 §2 与 §12.4 的 OB-51 这一行）。
- **产物**：本报告（双语两份）、自由度清单、`DELIVERY_DIRECTIVE{,_EN}`（升 v1.15：§0.3 一行 + §12.4 三行）、`ADR_REGISTER{,_EN}`（ADR-015）、`t027/REMAINING_SLICES{,_EN}`（切片定义回写块；收尾行于 ⑧b 补）、`DEV_PLAN{,_EN}`（T018 待议项指针）。其中 **v1.15 是本片的切片内副产物，不是独立的准则治理片**（ADR-015 同理）。
- **执行自由度清单**：`docs/validation/evidence/S1-REACH-freedom-list.md`（§1.7.3 固定路径；本片 F-1 至 F-8 逐条见该文件，本报告 §1 与其内容一致）。
- **一句话结论**：十门中**八门不可达或未通过**（G-A、G-C、G-D、G-E、G-F、G-H、G-I 与 OB-21 自身）、**一门未知**（G-G，人工门）、**一门已执行但不独立成立**（G-B）；DR-1 只消除了 OB-21 ② 的一个子项（`DR-1 未修`），**未改变 AT-23 的不可达判定**。

## 2. 可达性判定（逐门 × 依赖 × 解锁条件）

> 判定闭集 = 不可达 / 未通过 / 已执行不独立成立 / 未知。`未知` 仅用于人工门。解锁句模板固定为「需要 X 具备 Y，当前证据 Z 显示…」，Z 必须是可定位的文件或证据包。

| 门 | 判据（可判真假） | 证据 | 判定 | 解锁条件 |
|---|---|---|---|---|
| G-A AT-23 真实候选腿 | CLI 可真实启动一个经兼容性准入的 AI 候选并计费 | `internal/console/runtime.go:33`（`ExecuteNoopRun`，noop-only 路径）、`:207` 至 `:211`（`RunPostflight` fail-closed 桩）；`docs/DEV_PLAN.md` 的 T026 段（10 verified / 2 unsupported ⇒ 候选不兼容）与 T027 段（`BLOCKED / NOT IMPLEMENTED`） | 不可达 | 需要 T027 完成产品装配（当前证据 `runtime.go:207` 显示 postflight 仍为 fail-closed 桩）且存在兼容候选（当前证据 T026 段显示 2 项 unsupported 未变） |
| G-B AT-23 机制腿 | 隔离 workspace / timeout / 日志缺失 / 未知恢复 / exit 0 后重扫五面在真实宿主执行 | `docs/t027/REMAINING_SLICES.md` 的 B4 小节；`docs/validation/t027-at23-e2e.md` | 已执行不独立成立 | 无需新解锁（已执行）；该腿只在 G-A 满足后才并入 AT-23 整体判定；不得写成"机制腿通过即 AT-23 通过" |
| G-C AT-23 整体 | G-A 与 G-B 同时满足并含真实计费证据 | `docs/validation/s1-exit.md` 的 AT-23 这一行（`BLOCKED`）；B4 小节的四事分立 | 未通过 | 需要 G-A 解锁后重跑真实 E2E，当前证据 `s1-exit.md` 显示 AT-23 仍为 `BLOCKED` |
| G-D T028 与 AT-24 | `supervised-black-box` 显式模式实现完成 | `docs/DEV_PLAN.md` 的 T028 条目（`PLANNED / NOT IMPLEMENTED`）；`docs/validation/s1-exit.md` 的 AT-24 条目（`BLOCKED`） | 不可达 | 需要 T028 实施完成，当前证据显示其状态为 `PLANNED / NOT IMPLEMENTED`；该门**与候选兼容准入无关**（黑箱模式不要求能力矩阵，见 `docs/ADR_REGISTER.md:22` 的 ADR-012） |
| G-E 最终 ZIP 绑定与复验 | 最终 ZIP 绑定发行记录并复验 | `docs/ADR_REGISTER.md:19`（ADR-009：最终 ZIP 仍须绑定发行记录并复验）；`docs/validation/s1-exit.md` 的「正式安装包教程」行（`PARTIAL`） | 未通过 | 需要 T018 其余门闭合后执行绑定与复验，当前证据显示该项仍缺失 |
| G-F EOL 通知方式 | ADR-010 获产品所有者批准 | `docs/ADR_REGISTER.md:20`（ADR-010 状态 `PROPOSED / NOT APPROVED`） | 未通过 | 需要产品所有者明确批准 ADR-010，当前证据显示其状态为 `PROPOSED / NOT APPROVED` |
| G-G 产品所有者 S1 exit 批准 | 产品所有者作出批准（人工门） | `docs/validation/s1-exit.md` 的「产品所有者最终 S1 批准」行（`BLOCKED`） | 未知 | 需要产品所有者作出 S1 exit 批准，当前证据显示尚未作出；本片不预测其时点 |
| G-H T018 整体 | 依赖字段全部满足且四道剩余门闭合 | `docs/DEV_PLAN.md:91`（依赖字段含 T027 与 T028）、`:93`（状态 `IN_PROGRESS / BLOCKED`） | 不可达 | 需要 T027 与 T028 完成，当前证据显示两者未完成；**注意门级与任务级不同源**（见下方注） |
| G-I S1 exit 总体 | AT-01 至 AT-24 全过且 T018 完成且 RFC §14 全部 exit | `docs/validation/s1-exit.md` 头部结论行（`BLOCKED`）与 RFC §14 表（核心任务链 `BLOCKED`、三项 `PARTIAL`、业主批准 `BLOCKED`） | 不可达 | 需要 G-C 与 G-D 与 G-E 与 G-F 与 G-G 全部满足，当前证据显示其中任一均未满足 |
| OB-21 自身 | ①② 两个子项是否仍成立 | ① 未装配（`runtime.go:33` 与 `:207` 未变）；② 中「`DR-1 未修`」子项已消除（DR-1 闭环，`docs/validation/evidence/dr-1-2026-09-24/`），「T026 不兼容」与「B2 §8.3 与 §8.4」未变（`docs/validation/B2-EXTERNAL-ENFORCEMENT.md` §8 第 3 与第 4 条） | 未通过（处置词保持待议） | 需要候选兼容或装配出现进展才触发其「待议时机」，当前证据显示两者均无变化 |

- **门级与任务级不同源（重要）**：G-D 至 G-G 四道门**不依赖**候选或装配进展（T028 是独立模式），但 `docs/DEV_PLAN.md:91` 的 T018 **任务级依赖字段**含 T027 与 T028 ⇒ 四门全部清空后 T018 **仍不可勾选**。本报告按门级与任务级**分开**表述，不给单一结论。
- **DR-1 证据的引用口径**：`available` 记录只支持一条推理——同一 pin 上探针通道现已能驱动候选并产出最小真实响应；**不支持**任何「候选兼容」「候选可启动」「装配完成」推理（依据 `docs/CONTRACTS.md:191` 的「配置合法、秘密存在、CLI 已安装均不足以证明 AI 可用」）。
- **OB-21 这一行的历史句未改写**（指令 §11.1 的历史行口径；本片对授权歧义按 §1.7.2 从严处理）⇒ 其刷新结论落在本报告 §2 与 §12.4 的 OB-51 这一行，§12.4 该行的现象句仍含已被消除的子项。

## 3. 五点只读确认（OB-51 的事实基础）

1. **`--model` flag 与 `configuredAIProfile` 的默认值 / fallback**：`internal/console/ai.go` 的 `ai check` flag 集**不含** `--model`；模型唯一来源是链配置 `ai.profiles[].model`，经 `configuredAIProfile`（`ai.go:230`）按 channel 绑定到 profileId 到 profile 三级查找，**无默认值、无 fallback**（任一缺失即返回未配置并 fail-closed）；空 `model` 被 `internal/adapters/ai_provider.go:51` 拒绝。
2. **schema 里 `model` 是否 required**：`schemas/` 目录**无任何一份覆盖链配置 `ai.profiles`** ⇒ 该问句**不适用**；必填性由 Go 校验承担（`ai_provider.go:51`）。`schemas/chain.schema.json` 的 `model` 出现在 `:116`、`:156`、`:185`，均为 task 或 codeStep 级字段，且对应 `required` 数组（`:108`、`:148`、`:177`）**不含** `model` ⇒ 在那里是可选，与 `ai.profiles` 无关。
3. **`model: "auto"` 是否被 schema 或代码禁止**：**不被禁止**——代码侧无任何针对 `auto` 的等值比较，schema 侧无覆盖；该值作为普通字符串经 `internal/adapters/ai_probe_copilot.go:114` 原样透传给 CLI，仓库测试把它当合法值（`internal/adapters/ai_provider_test.go:13`、`internal/adapters/ai_preflight_test.go:110`）。CLI 侧 `auto` 的实际路由语义属仓库外事实，**本片无探针额度 ⇒ 不可验证**，如实标注。
4. **BYOK 密钥条目的存在性**（只报存在性与条目标识名，不读内容）：① **环境变量通道当前不在位**——本机仅存在 `COPILOT_AGENT`、`COPILOT_DEBUG_NONCE`、`GH_TOKEN`，而 T026 的 BYOK 实跑正是从进程环境读取 `COPILOT_PROVIDER_API_KEY` 或 `DEEPSEEK_API_KEY`（`docs/validation/t026-agent-runner-contract.md:117`）；② **Windows 凭据管理器**存在 Copilot CLI 的 host 条目（`LegacyGeneric:target=https://github.com` 前缀的 CLI 条目，完整标识名按指令 §7.3 不落盘），且**不存在任何 `ProofRail` 前缀条目**；③ **Copilot CLI 自有存储** `%USERPROFILE%\.copilot\` 存在 `config.json`、`session-store.db` 与 `session-store.db-shm`、`session-store.db-wal`、`logs`、`session-state`、`installed-plugins`、`ide`，**未读任何内容**，**文件名级**未见 BYOK 密钥文件；④ **VS Code secretStorage 不可枚举**（外部无枚举 API）⇒ 如实标注为不可枚举，**不作任何结论**。
5. **prfrail 自有 secretStorage 机制**：**有**。后端 = 当前用户的 Windows Credential Manager（generic 型、local-machine 持久性），前缀 `windows-credential:`（`internal/adapters/ai_secret_windows.go:15`），写入时用户名域为 `ProofRail`，目标名 = 引用串去掉前缀（`windowsCredentialTarget`，`internal/adapters/ai_secret_windows.go:139`）；CLI 为 `prfrail secret set`、`status`、`delete`（无回显、禁 argv）；契约见 `docs/CONTRACTS.md:193`。`docs/OPERATIONS.md:56` 的 `secretRef` 示例与 `:65` 的 `secret set` 命令**是同一机制的文档示例**（示例 target `ProofRail/deepseek`），非另有一套机制；**该示例条目当前不在位**（凭据枚举无 `ProofRail` 前缀条目）。⇒ 与 DR-1 的「BYOK 归属未证」一致，并加强为：本次枚举到的**全部**凭据条目中无 BYOK 形态目标名，环境变量通道亦无 BYOK 密钥 ⇒ 就两处可枚举通道而言，**BYOK 凭据当前不在位**。
- **停止判定**：五点结论**均不指向** CONTRACTS / schema / 授权契约的语义改动 ⇒ **未触发**切片定义中的"停下升级"。

## 4. 执行流水线

| 环 | 角色 | 状态 | 产物 |
|---|---|---|---|
| ① 架构 | V4 Pro（`deep-reasoner`） | 已执行 | 方案书（十门判定框架、五点方法、验收判据、文件级改动点、停止条件） |
| ②a 契约 | — | 未触发（无契约改动） | — |
| ②b 实现 | 主控（文档型退化） | 已执行 | §12.4 三行、ADR-015、台账定义块、`DEV_PLAN` 指针 |
| ③ 测试 | 主控（文档一致性检查） | 已执行 | 门禁与镜像读数（见 §10） |
| ④ 主控集成与门禁 | 主控 | 已执行 | 见 §10 |
| ⑤ 预审 | V4 Pro | 见 §7 | 见 §7 |
| ⑥ 独立扫描 | 拟跳过（§7.2.1 三条件） | 见 §7 | 由 ⑦ 确认 |
| ⑦ 终审 | V4 Pro（降级，ADR-015） | 见 §7 | 见 §7 |
| ⑧a 文档初稿 | 主控 | 见 §7 | 本报告 |
| ⑨ 原生验证 | 主控 | 见 §10 与 §11 | 门禁复跑与镜像核验 |
| ⑧b 文档定稿 | 主控 | 见 §14 | 收尾行与完成记录 |
| ⑧c 收尾清洁 | 主控 | 见 §11 | 见 §11 |
| ⑩ 停点 | 主控 | 已停 | 本报告 |

## 5. 异常与兜底记录

1. **结构偏离（F-1）**：本报告在 §11.3 的固定小节序列中插入 §2「可达性判定」与 §3「五点只读确认」两节——该两节是本片交付物主体，置于流水线之前便于阅读；**不改任何判据、不改任何结论行语义**。
2. **主控兜底 = 1**：⑤ 首轮报 1 High 与若干 Medium / Low ⇒ 整改后按 §5.3 重跑 ⑤（未发生 ⑨ 失败）。
3. **行号口径**：本报告引用的 `file:line` 取自 ① 阶段与 ② 阶段的实读；仓库存量文档的行号会随其它切片变动而漂移，引用时以**文件名与节号**为主、行号为辅。

## 6. 环境与档位

- 平台：Windows 本机（只读枚举与门禁）；无 Linux 原生步骤（本片不涉平台验证）。
- 档位：`[SLICE]`（文档型，按 §5.2 退化）；主控档位 = 标准；①、⑤ 与 ⑦ = `deepseek-v4-pro`（V4 Pro）。
- **⑦ 降级（ADR-015）**：本片 ⑦ 由 `deepseek-v4-pro` 承担，**不是** §2.2 白名单的 `gpt-5.3-codex` ⇒ **本报告如实标注：⑦ 独立性降低**，结论交 ⑧c 复核；该例外**不推广**、**不得作为先例**。
- **试点顺延第 4 次（四要素）**：裁决人 = 用户；日期 = 2026-09-29；理由 = 业务判定类切片、非准则治理；依据 = 本轮裁定原文。**用户直接裁定、不得作为先例**。本片**不做试点指标声明**（§12.1 口径），不计入试点评估。
- **探针额度 = 0**：本片未发起任何真实外部调用。

## 7. 审查结论摘要（⑤⑥⑦）

- ⑤ 预审（`deepseek-v4-pro`）**第 1 轮**：`PRE-REVIEW: FINDINGS`——1 High（授权类：OB-21 那一行的「更正注记」）与 5 Medium 与 7 Low。整改：**回退 OB-21 那一行的更正注记**（按 §1.7.2 从严处理授权歧义，改记自由度清单 F-2）；报告 §1 补 v1.15 副产物句（M1）；§1 结论计数改正（M3）；§3 行号改正与「文件名级」限定（M5、L1）；OB-53 出处标注（L7）；用词统一（L4、L5）等。
- ⑤ 预审**第 2 轮（复审）**：`PRE-REVIEW: PASS WITH FIXES`——H1 主体已消除、M1 至 M5 全部消除、L1 至 L7 消除或保留并登记；新增 2 条 Low（N1：报告写 F-1 至 F-7 而清单为 F-1 至 F-8；N2：OB-21 那一行行尾多一空格）。两条 Low 已按字面修正：清单条目数改为 F-1 至 F-8；OB-21 那一行恢复为与 HEAD 逐字节一致（该行已退出变更集）。**⑤ 额度 2/2 用尽** ⇒ 不再重跑 ⑤，两条字面修正的复核交 ⑦ 复审确认（留痕于本节末条）。
- ⑥ 独立扫描：按 §7.2.1 三条件判定——① 改动只落文档与台账类文件；② 不含 `.go` 与 `.json` 与 schema 与 fixtures 与脚本逻辑的语义变更；③ 不涉及状态机、门禁判据、证据模型、角色、授权契约语义 ⇒ **三条同时满足、可跳过**；按 §7.2.1「跳过必须由 ⑦ 或用户确认」**已交 ⑦ 确认**（⑦ 判定「可跳过」）。
- ⑦ 终审（`deepseek-v4-pro`，降级，ADR-015）**第 1 轮**：`FINAL REVIEW: PASS WITH FIXES`——四项均「基本成立或成立」，三条 Medium（M-1 门禁落盘证据过期、M-2 成本计量自相矛盾、M-3 §7 悬空引用）与四条 Low；整改与复审结论见本节末条；该层**独立性降低**（非 §2.2 白名单模型），结论交 ⑧c 复核。
- ⑦ 终审**第 2 轮（复审）**：`FINAL REVIEW (re-review): PASS`——M-1 至 M-3 与 L-1 至 L-4 全部消除；代 ⑤ 确认的两条字面修正（N1/N2）已消除；无 Medium+。新增 1 条 Low（凭据枚举原始输出含用户名与令牌形态标识，与其消毒承诺不符）⇒ ⑧c 已删除该文件并从 artifacts 块移除其行。该轮为 ⑦ 的最后一轮（额度 2/2）。

## 8. 审查输入包摘要（含盲审隔离证明）

- 输入包（主控落盘，路径交给审查方）：`tmp/s1-reach/plan-01.md`（① 方案书）、`tmp/s1-reach/review-input.diff`（`DELIVERY_DIRECTIVE{,_EN}`、`ADR_REGISTER{,_EN}`、`REMAINING_SLICES{,_EN}`、`DEV_PLAN{,_EN}` 的原始 diff）、`tmp/s1-reach/s1-reach-reachability.md` 与 `_EN.md`（本报告全文，审查对象；⑧a 落盘时与 `docs/validation/` 同名文件逐字节一致）、`tmp/s1-reach/gate-tree.txt`（门禁原始输出）、凭据枚举的**消毒结论**（原始输出含用户名，按指令 §7.3 不交审查方）。
- **消毒**：交审查方的材料已过滤用户名、主机名、IP、凭据与令牌（指令 §7.3）；凭据枚举只提供**与主题相关的两条**（CLI host 条目存在、`ProofRail` 前缀不存在），其余无关条目**不给**审查方。
- **盲审隔离**：⑦ 首轮**不附** ⑤ 与 ⑥ 的清单（⑦ 只拿切片定义、diff、本报告全文、只读约束）；⑥ 拟跳过 ⇒ 无 ⑥ 清单可附。
- 本片**不触**「所有权 / 停机 / 身份」语义 ⇒ **不触发 §7.7 必抽盲审**；按「每 4 片抽 1 片」的口径，本片不在必抽位（抽样记录由后续切片按 §7.7 判定）。

## 9. 可证伪性（机械检查 + 变异检验）

- 机械检查：见 §10（G1 至 G7 与 §6.1 四条基础门禁）。
- **文档变异检验**（每条先施加变异、确认判据变红，再逐字节还原并复核哈希一致）：① 删除 CN 侧 §12.4 的 OB-52 这一行 ⇒ G3(a) 对称判据应变红；② 把 OB-51 的首状态词改成闭集外词 ⇒ G6-6 应变红；③ 删除本报告的 `artifacts` 块 ⇒ G7-a 应变红；④ 在 §0.3 的 v1.15 这一行塞入行数与文件数类字面量 ⇒ G6-4 应变红；⑤ 把头部版本改回 v1.14 ⇒ G1-b 应变红。实际读数见 §10。

## 10. 门禁结果

- `node tools/gates/gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0`（`gate-exit=0`，`changed=9`；报告尚未入库 ⇒ `G7-a` 至 `G7-d` 为 `INFO`，落盘后复跑）。
- 逐项：`G1-a` PASS；`G1-b` PASS（`header=v1.15 last=v1.15 2026-09-29/2026-09-29`）；`G1` SKIP（未脚本化）；`G2` PASS；`G3(a)` PASS（四对对称：指令 `+5/-1`、ADR `+1/-0`、`DEV_PLAN` `+1/-1`、台账 `+19/-0`）；`G3(b)` PASS；`G4a` PASS；`G4b` PASS（`unknown=[]`）；`G5-a` 与 `G5-b` PASS；`G5` SKIP（未脚本化）；`R2.5` PASS；`G6-1` 至 `G6-4` PASS；`G6-5` PASS（`scanned 2, refs 6`）；`G6-6` PASS（`scanned 8`）。
- 基础门禁：`gofmt -l` 空、`go build ./...` exit 0、`go vet ./...` exit 0、`go test ./...` `ok=14 FAIL=0`。
- 镜像：报告对与指令对逐行失配均为 0（`_EN` 与 CN 同位）；`ADR_REGISTER{,_EN}` 与 `REMAINING_SLICES{,_EN}` **不做整文件位置镜像**（工作区既有位置性差异），对其适用 `G3(a)` 的改动行数对称。
- ⑧c 复核读数（含报告落盘后的 `G7` 复跑）见 §11。

## 11. ⑧c 收尾清洁记录

- 全仓逐文件编码核验（`enc-tree.js`）：违规 **3** 处，**全部**为 OB-27 的冻结证据包豁免项（`b2-2026-09-17` 的 `.raw.txt` CRLF 与 `b3b-2026-09-21` 的两处 BOM）；**本片新增文件 0 违规**。
- 门禁自检：`SELFTEST: PASS 248/248`；`gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0`（报告落盘后 `G7-a` 至 `G7-d` 均为 PASS；`G7-b` 覆盖 28 条工件行，`G4a` 由该判据自身输出所核文件数）。
- 临时残留清单：`tmp/s1-reach/` 保留 13 个过程文件（方案书、草稿、断言脚本、门禁原始输出、审查包副本），按 OB-49 与用户既有指令**不清理**，处置 = `local-only`；仓库内未发现未跟踪的非临时残留。
- ⑧c 发现并处置 1 项：⑦ 复审的 Low（凭据枚举原始输出含用户名与令牌形态标识）⇒ 已删除该文件并同步移除报告 artifacts 块中的该行。本节与 ⑧b 的写入均为**事实记录**（工具读数、审查结论、台账回写），按 §5.3 的语义边界**不触发重跑 ⑦**；若用户认为应触发，本片 ⑦ 额度（2/2）已用尽 ⇒ 必须停下升级用户裁决。

## 12. 成本与计量

| 环 | 角色 | 模型 | 次数 | reason / retryOf |
| --- | --- | --- | --- | --- |
| ① | 架构师 | `deepseek-v4-pro` | 1 | 切片定义要求 ①（判定框架、五点方法、改动点） |
| ⑤ | 预审员 | `deepseek-v4-pro` | 2 | 第 2 次为整改后复审（retryOf = 第 1 次） |
| ⑥ | 独立扫描 | 无 | 0 | 按 §7.2.1 可跳过，已由 ⑦ 确认 |
| ⑦ | 独立终审 | `deepseek-v4-pro` | 2 | 终审（降级，ADR-015）与整改后复审 |
- 探针：**0 次**（额度 0）。
- 计量真值源：会话调用回执（§7.5 硬规则 R7.1），不口述估算。

## 13. 明确未执行事项

- commit 与 push（须同轮授权）；真探针（本片额度为 0）；BYOK 实测（若判定需要 ⇒ 停下升级，本片不做）；⑥ 独立扫描（拟跳过，由 ⑦ 确认）；试点指标声明（本片不计入试点评估）。

## 14. 下一步建议

- **本片结论不解除任何东西**：AT-23 真实候选腿仍不可达；T018 与 S1 exit 仍阻断。
- 解阻顺序建议（需用户裁决后另立切片）：先 **T027 装配**（`runtime.go` 的 noop-only 路径与 fail-closed postflight 桩是硬前置），再评估候选兼容性；T028 与候选兼容**无依赖关系**，可与候选线并行评估。
- 档位建议：T027 装配切片涉及状态机与崩溃路径 ⇒ 应保留 `gpt-5.3-codex`；本片这类纯判定切片可继续用 V4 Pro。
- 未闭合项移交：`OB-51` 的「VS Code secretStorage 不可枚举」为本机与外部工具的**共同边界**，后续任何切片都无法枚举，建议长期保持如实标注。

```artifacts
docs/validation/s1-reach-reachability.md	repo
docs/validation/s1-reach-reachability_EN.md	repo
docs/validation/evidence/S1-REACH-freedom-list.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
tmp/s1-reach/plan-01.md	local-only
tmp/s1-reach/review-input.diff	local-only
tmp/s1-reach/gate-tree.txt	local-only
```
> **归档注记（v1.23）**：本报告引用的 `tmp/` 条目已按 §1.5.1 归档至 `docs/validation/evidence/directive-history/`（单元：`s1-reach/`；各含 `MANIFEST.md` 与 `SHA256SUMS.txt`）；原路径不改写，对应关系以 MANIFEST 为准。（此前切片已清除的 `tmp/` 路径不在本次归档范围，属 OB-77 所指既有脱钩。）
