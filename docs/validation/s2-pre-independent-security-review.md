# P79 验证报告：S2 启动前独立安全评审（补评）

日期：2026-11-19。状态：进行中（② 文档实现阶段）。切片标识：P79 `S2-PRE-INDEPENDENT-SECURITY-REVIEW`（`[SLICE]`，评估类）。

## 1. 切片标识与授权

- 编号与标识：P79 `S2-PRE-INDEPENDENT-SECURITY-REVIEW`；分级 `[SLICE]`；类型 = 评估类（独立安全评审补评）；规模 M-L（v65 快照 §四已冻结登记，本片不改快照）。
- 主控：`DeepSeek V4.1 Flash (deepseek)` High。本片判断负载在 ⑦，主控只承担编排 / 材料化 / 转录 / 门禁 / 回执，不以自身读感补写或软化任何评审判断，亦不代 ⑦ 给出任何面结论。
- 触发依据：ADR-047（`docs/ADR_REGISTER.md:58`）。义务三要素 = S1 内不补评、豁免顺延至 S2、须在 S2 启动前完成一次独立安全评审并登记结论；到期日 = S2 启动日与 2027-03-31 中较早者，到期未补评则不得启动 S2 实测。
- 评审对象：`PRODUCT_REQUIREMENTS.md` §7 的「独立安全评审结论」门禁（`:109` 阻断门禁句）；历史状态为从未执行，佐证 `SECURITY.md:5`。该对象与 `docs/validation/s1-exit.md:59` 门表「独立评审」行（T017 基线）是不同对象，证据不可互换，本片不改 s1-exit 门表。
- 授权四要素：裁决人 = 用户；日期 = 2026-11-19；理由 = ADR-047 补偿义务到期前补评；依据 = 本次起跑消息 + v65 快照 §四。
- ⑦ 额度裁定：按 L 档 = 1 终审 + 3 复审。该裁定**超出 §7.5「M-L 视同 M」字面标准，系用户直接裁定、不得作为先例**，报告计量节须如实标注。
- 探针：0（本片无外部真实调用）。
- 盲审：§7.7 必抽三语义（所有权 / 停机 / 身份）在本片不触，故不抽；不抽理由在步骤 0 登记。
- 边界：不改 S1 Exit 八门判据；不改 `SECURITY{,_EN}.md` 威胁模型基线（仅结论归档 + 注记）；不改准则 / `CONTRACTS` / schema / fixtures / 产品代码 / CI；不新增 OB / ADR（用户裁定除外）；无网络调用、无外部命令、不读 `%TEMP%`、不读仓库外文件。

模型与审查链（启动前一次性设定，不得中途升档）：

| 环节 | 角色 | 模型 | 档位 / 额度 |
|---|---|---|---|
| ① | 架构师（`deep-reasoner`） | `DeepSeek V4 Pro (deepseek)` | Max；产出评审方案书 |
| ⑤ | 预审员 | `DeepSeek V4 Pro (deepseek)` | 额度 1 + 1；报告初稿预审 |
| ⑥ | 独立扫描 | `MAI Code 2 Flash (copilot)` | ≤3；初稿输出契约符合性扫描 |
| ⑦ | 独立终审 | `GPT-6.2 Sol (copilot)` | Max；1 终审 + 3 复审，取代 L 档表中的常规终审档位 |
| 批量辅助 | 辅助执行 | `GLM-5.3-Flash (glm)` | ≥2 次、严格串行；仅限枚举 / 摘录与数字核对 / 双语镜像检查 / MANIFEST 草案四类，禁择优 |

## 2. 步骤 0 锚点实测

| 项 | v65 快照值 | 实测值 | 结论 |
|---|---|---|---|
| `HEAD` | `9fd8cd6` | `f211bca` | 相符（v65 自身落盘提交；快照 §一 已声明「v65 自身落盘后 HEAD 前移，届时以实况为准」） |
| `origin/main` | `9fd8cd6` | `f211bca` | 相符（同上口径） |
| 准则版本 | `v1.71`（CN / EN 头部均为 `v1.71`） | `v1.71`（CN / EN 头部一致） | 相符 |
| `SELFTEST` | `node tools/gates/gate.js --selftest` | — | `PASS 740/740`（exit 0）；读数文件 `gate-selftest-final.txt` |
| 全树编码 | `node tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/enc-scan.js` | — | 违规全部位于 `docs/validation/evidence/**` 的冻结证据包（非冻结主树为 0）；编码硬规则见 §6.1 |
| 三作用域 | `TOTAL_FAIL=0` | `index` / `tree` / `ci` 均 `TOTAL_FAIL=0` | 相符 |
| v63 全文件行数 | 222（CN / EN 对称） | 222（按 LF 切分，CN / EN 对称） | 相符 |
| v64 全文件行数 | 224（CN / EN 对称） | 224（按 LF 切分，CN / EN 对称） | 相符 |
| v65 全文件行数 | 由 v66 复核 | 本版落盘后 HEAD 前移 | 按快照口径不判定 |
| `tmp` 顶层 | `gates-cache` + `.gitkeep` | `.gitkeep` + `gates-cache` + `P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW`（本片新增，属合法） | 相符 |
> **时点口径（§11.3，v1.53）**：本节读数为**步骤 0 基线**（`changed=0`）时点值；变更集门禁的当前读数以「门禁结果」章与 ⑩ 回执为准（本片 ⑥ 转录落盘前 `G7-f` 为已知待满足项，⑨ 前必须转绿）。

## 3. 评审范围与三面映射

三面为 ADR-047「未覆盖范围」逐字承接，不得增删：

| 面 | 评审对象 | 判据来源 | 结论档 |
|---|---|---|---|
| ① 威胁模型面 | `SECURITY.md` §1（资产与信任边界，正文 :11）、§2（写权限矩阵表，:15–:21）及其引用的滥用场景载体 §3 STRIDE 表（SEC-01…SEC-15，:29–:44）与 §3.1（C6 预建槽，:50–:54）；`SECURITY_EN.md` 同结构。评审动作 = 独立复核：资产清单完整性、§1 信任主张与 §2 矩阵行/列权限的一致性、STRIDE「防护及责任」列是否可由既有组件兑现、滥用场景是否闭合 | `SECURITY.md:5`（REQ-004/005/007/013/019/021-026/028）；`PRODUCT_REQUIREMENTS.md:109`（门禁句）；RFC §12、§16.7–16.8（设计来源）；`TEST_STRATEGY.md` §2 AT-01…AT-24 表（AT 编号解析，`SECURITY.md:46` 表尾句）；`ADR_REGISTER.md:58` 未覆盖范围 ① | 报告「威胁模型面」小节：该面 Section A–E 块 + 面结论行 `FACE 1: PASS\|FINDINGS`（模板见 §6） |
| ② 硬安全门禁面 | 判据侧：`docs/validation/sg6-gatelevel-consolidate.md`、`tools/gates/gate.js`、`tools/gates/lib/criteria/{g1,g5a,g7,g9}.js`。实现侧：`internal/gates/policy.go` `Prepare`(:30)/`validateRequest`(:48)/`validateHookEffectPolicy`(:96)；`internal/gates/runner.go` `Run`(:41)/`assess`(:115)/`scanOutput`(:172)；`internal/gates/readiness.go` 级别闭集(:70–83)/`ProbeReadiness`(:113)/`AcceptanceReadiness`(:118)；`internal/gates/probe_permit.go` marker 状态机(:74–91)/`NewProbePermit`(:126)；`internal/guard/lease.go` `Claim`(:82)/`Takeover`(:104)/`validTakeoverProof`(:120)；`internal/guard/sandbox.go` `VerdictStatus`(:21–34)；`internal/chain/postflight.go` `requirePostflightQualifiedReview`(:308)；`internal/evidence/verify.go` `VerifyManifestObjects`(:37)。评审动作 = 独立复核「fail-closed 前置条件的独立性」：每个拒绝路径的前置条件是否可被同一被防主体自行满足或绕过；「修改策略」与「独立批准」是否分离 | `SECURITY.md:23`（同段两句：「限制必须由 runner/OS 能力落实；不能仅靠提示词」与「执行者不得拥有修改策略和独立批准同一变更的合并权限」）、§2 矩阵行语义；`CONTRACTS.md` §4–5（Prepare/Inspect/Validate/Promote、handoff）、§7（admission 纯 preflight、dispatch 前重新确认授权未撤销且 reservation outstanding、任一不满足 fail-closed）；`TEST_STRATEGY.md` AT-05（拒评/自批拒绝）/AT-06/AT-18（硬门禁不可 waive）/AT-19；`ADR_REGISTER.md:58` 未覆盖范围 ② | 报告「硬安全门禁面」小节：Section A–E 块 + `FACE 2: PASS\|FINDINGS` |
| ③ 秘密与凭据面 | `internal/adapters/ai_secret_windows.go` `NewPlatformAISecretManager`(:48)/`Set`(:52)/`Resolve`(:77)/`Exists`(:89)/`Delete`(:101)/`windowsCredentialTarget`(:139)/`containsASCIIControl`(:150)；`internal/adapters/ai_provider.go` `ValidateAIProviderProfile`(:47)/`validHTTPSBaseURL`(:81)/`validDNSName`(:101)/`validAISecretReference`(:115)；`internal/adapters/ai_probe_copilot.go` `Probe`(:60)/`copilotCLIProbeEnvironment`(:134)/`clearCopilotCLIProbeEnvironment`(:150)；`internal/console/ai.go` `executeAICheck`(:33)/`configuredAIProfile`(:230)；`docs/OPERATIONS.md`（`ai check` 读取绑定 SecretStore 凭据、`secret set` 无回显）。评审动作 = 独立复核：secretRef 不透明性、凭据注入/清零生命周期、凭据是否可能进入配置/摘要/日志/错误/证据、BYOK 通道（Copilot CLI probe）凭据仅经子进程环境且用后清零、非 Windows fail-closed | `SECURITY.md` §4（:56–:58：秘密不经模型、不写配置/日志、secret-direct、外发扫描）；`CONTRACTS.md` §7 secretRef 段（「profile 只保存 provider 类型、base URL、模型、传输选项及不透明 secretRef；API key/token 不得进入配置、摘要、日志、错误或证据，运行时只能由系统 SecretStore 解析并临时注入子进程环境」）、`windows-credential:<target>` 段（`prfrail secret set` 仅已确认本地控制台无回显、`status` 只返回引用与存在性、adapter 只收短暂副本用后清零、非 Windows 平台 fail-closed）、「`config explain` 不得显示 secretRef/凭据存在性/解析值」「静态 validate/explain/preview 不读取 SecretStore」；`TEST_STRATEGY.md` AT-10（secret-direct 仅终端、secret canary 不入证据/模型）、AT-06（秘密扫描失败阻断）；`DEV_PLAN.md:119`（实现已完成句）；`ADR_REGISTER.md:58` 未覆盖范围 ③ | 报告「秘密与凭据面」小节：Section A–E 块 + `FACE 3: PASS\|FINDINGS` |

表内行号为盘上定位线索；⑦ 引用时须按 §7.3 贴原文定位断言（块头行号 + 逐字照录），定位或逐字任一不成立者按「线索」处理，不得采信。

## 4. 评审方法与输入包

- 输入包分 A / B / C / D / E 五组登记；件数、行数与字节数以 `input-list.md` 实测为准（本报告不登记件数与字节合计，逐件读数由该清单断言）。
- 附加件：`00-briefing.md`（起跑消息逐字）、`freedom-list.md`、方案书 §2 / §6 的契约性事实列、`01-arch/package-manifest.md`（消毒登记，见下）。

| 组 | 内容 | 件数 |
|---|---|---|
| A | 安全基线与需求：`SECURITY{,_EN}`、`PRODUCT_REQUIREMENTS{,_EN}`、`ADR_REGISTER{,_EN}`、`DEV_PLAN{,_EN}`、`TEST_STRATEGY{,_EN}` | 10 |
| B | 执行准则与快照：`DELIVERY_DIRECTIVE{,_EN}`、`CHECKLIST_v65{,_EN}`、`REMAINING_SLICES{,_EN}` | 6 |
| C | 判据来源：`sg6-gatelevel-consolidate.md`、`s1-exit{,_EN}.md`、`gate.js`、`g1/g5a/g7/g9.js` | 8 |
| D | 产品侧硬安全门禁实现：`internal/gates/*` 4 件、`internal/guard/*` 3 件、`internal/chain/postflight.go`、`internal/evidence/*` 3 件 | 11 |
| E | 秘密与凭据面：`ai_secret_windows.go`、`ai_provider.go`、`ai_probe_copilot.go`、`internal/console/ai.go`、`docs/OPERATIONS.md` | 5 |

消毒规则（逐字过滤，不得以「未发现」代替扫描）：

1. 过滤 IP 字面量（IPv4 / IPv6，含内网段形态）。
2. 过滤用户名、主机名、SSH 目标串（形如 `git@<host>:`）、`%TEMP%` 与仓库外绝对路径。
3. 过滤凭据 / token / 密钥 / `password=` 形态。
4. 过滤后以 `<redacted:类别>` 占位，并保持行数不变（行数口径与 `input-list.md` 一致，便于行号定位不漂移）。
5. 保留仓库内相对路径 `docs/...`、`internal/...`（§7.3 明列）；`tools/...`、`examples/...` 经主控确认非敏感后保留。
6. 消毒结果落 `01-arch/package-manifest.md`：逐件登记消毒前 / 后的 SHA-256、行数、字节数与过滤命中计数。

材料化要求（⑦ 只有只读工具，不得声称执行过命令）：

- 凡要求 ⑦ 复核的读数，必须由主控预先落成输入包内的静态材料，⑦ 只依据静态件作答。
- 带行号摘录：`SECURITY{,_EN}` §1–§2 / §4 与 `CONTRACTS.md` §7 secretRef 段。
- 函数级带行号切片：内部实现文件给出函数头 + 关键拒绝分支，不整文件粘贴。
- 机械读数静态件：`go test ./...` / `go vet ./...` 输出、SELFTEST 计数等，由主控在对应步骤执行后落盘为 `.txt` 并入包。
- 引用口径：⑦ 只能写「据输入包内 `<路径>`」，不得写「执行了某命令」。

## 5. 独立性留痕、输出契约与结论档

- 判断主体唯一：全部评审判断由 ⑦ 深度档 `GPT-6.2 Sol (copilot)` Max 作出；主控只做编排 / 材料化 / 转录 / 门禁 / 回执，不补写、不软化、不改写任何发现或结论行。
- 流水线顺序与锚定隔离：⑤ 预审（1 + 1）→ ⑥ 独立扫描（≤3）→ ⑦ 终审盲态；严格串行，⑤ / ⑥ 清单一律不得进入 ⑦ 首轮输入。
- ⑦ 首轮盲态：首轮输入 = 起跑消息原文 + ADR-047 所在行的原文 + 三面定义 + 消毒后输入包 + 方案书 §2 判据来源列与 §6 结论档模板；不附 ⑤ / ⑥ 清单、不附主控工作笔记、不附方案书其余含倾向性内容的节；复审轮方可附清单并标注「完成独立扫描后再读」。
- 转录契约：⑦ 最终消息由主控逐字转录进报告结论摘要节围栏块，Section A–E 标题、面结论行、合成结论行一律原文照录；不得按要点重建、不得重组空行、不得改写措辞。
- 留痕：报告输入包摘要节含盲态证明（首轮是否附清单、输入是否已消毒、MANIFEST 指针）；计量按单一真值源逐字段登记，⑦ 额度含 1 + 3 用户裁定标注。

⑦ 三面各输出一个独立 Section A–E 块，另加一个不含新发现的合成块：

| 段 | 内容 |
|---|---|
| Section A | 该面结论行 `FACE <i>: PASS\|FINDINGS`（i = 1 / 2 / 3；取值限二值） |
| Section B | 发现清单，每行固定 `SEVERITY \| file:line \| 被违反的句子（逐字引用） \| 最小修复`；`SEVERITY` ∈ {`Critical`,`High`,`Medium`,`Low`}；`file:line` 必须盘上可定位；无发现则为空 |
| Section C | 逐条证据与推理，只引用输入包静态材料与契约原句 |
| Section D | 可证伪性审计：逐项给出「若删掉 / 改回某守卫，哪个测试或断言会红、哪个不会红」；未覆盖的不变量以 `UNCOVERED` 独立成句标注，不计入发现数、不影响结论行 |
| Section E | 「我无法验证的部分」+「什么会推翻我的 PASS」 |

结论档定义与合成口径：

- 每面结论档 = 该面 Section A–E 块 + 面结论行 `FACE <i>: PASS\|FINDINGS`；`PASS` 仅记录「本片补评未发现违反判据的缺陷」，不证明安全；`FINDINGS` = 至少一条发现，每条附 `file:line` + 最小修复。
- 首读二值约束：本片为评估类首读、不设修复轮，面结论行限 `PASS\|FINDINGS`；若 ⑦ 输出表外形态（如 `PARTIAL`），主控不自行改写，原词转录并按升级路径交用户裁决。
- 一致性硬约束：Section B 无发现 ⇒ 该面结论行必须 `PASS`；Section B 有发现 ⇒ 该面结论行必须 `FINDINGS`。
- 合成口径（机械规则）：三面均 `PASS` ⇒ `FINAL REVIEW: PASS`；任一面 `FINDINGS` ⇒ `FINAL REVIEW: FINDINGS`。
- 两条限定须与合成口径同时登记：**任一面的 PASS ≠ 整体 PASS**；**整体 PASS ≠ S2 可启动**（S2 启动仍受其它门约束，且整体 PASS 不等于「安全已证明」，`UNCOVERED` 项照录）。
- 格式偏离处置：未按五段 / 无面结论行 ⇒ 重发 1 次（计入 ⑦ 额度）；再偏离 ⇒ 停止并升级用户裁决，不得自行修复格式。
- ⑥ 转录：⑥ 最终消息由主控逐字转录进「审查结论摘要」节的 `sixstage` 围栏（Section A–E 标题 + 块末结论行），块外只放主控注记与指针。

## 6. 审查结论摘要（⑥ 逐字转录）

```sixstage
## Section A 总体结论
INDEPENDENT SCAN: PASS

## Section B 发现清单
无

## Section C 逐条证据与推理
- CN/EN 报告均在 §6 建立 `sixstage`围栏：[CN :114](docs/validation/s2-pre-independent-security-review.md#L114)、[EN :114](docs/validation/s2-pre-independent-security-review_EN.md#L114)。A–E 标题各仅出现一次，Section A 与块末结论均为 `INDEPENDENT SCAN: PASS`。
- 行号问题已闭合：CN/EN 主文均引用 §4 `:56–:58`；原始错误的 `:37–40` 仅残留于 r1 发现复述，不是现行引用。
- 字节口径已闭合：`REMAINING_SLICES.md` 在两份静态件中均为 `381800`，未发现旧值 `381794`：[input-list :30](docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/input-list.md#L30)、[package-manifest :25](docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/package-manifest.md#L25)。
- 消毒口径已闭合：逐件登记含行数列，并明确未改写时消毒前后 SHA 相同：[manifest :9](docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/package-manifest.md#L9)、[说明 :60](docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/package-manifest.md#L60)。
- 静态门禁读数报告 `G5 PASS`，`TOTAL_FAIL=0`：[gate readout :1](tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/gate-g7f2.txt#L1)、[汇总行 :6](tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/gate-g7f2.txt#L6)。

## Section D 可证伪性审计
- 删除或遗漏任一份 `sixstage`围栏：CN/EN 形状复核立即产生发现。
- 在主文重新引用 `SECURITY.md §4 :37–40`：当前行号锚点复核立即红。
- 将输入清单恢复为 `381794` 或使两份静态件漂移：字节一致性复核立即红。
- 删除行数列或消毒前后 SHA口径说明：逐件完整性复核立即红。
- 删除门禁读数或恢复非零 `TOTAL_FAIL`：门禁结论立即红。
- UNCOVERED：未读取 `04-reviews/`，未覆盖运行时、网络调用及真实生产环境。

## Section E 我无法验证的部分 + 什么会推翻我的结论
无法验证文档之外的运行面、真实生产安全结论及未读取的预审意见；静态门禁件未逐项列出 `G7-f/G7-a`。若后续证实围栏缺失、行号错误、字节值漂移、消毒登记缺失或门禁非零，本结论应改为 `INDEPENDENT SCAN: FINDINGS`。

INDEPENDENT SCAN: PASS
```

## 7. ⑦ 独立终审逐字转录（安全评审本体）

本节为 ⑦ 首轮最终消息的逐字转录，原文归档见 `docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/07-security-r1.md`；主控未改写措辞。

```text
定点核对完成。三面都有能定位到源码和契约的发现；结论会区分已接线的 BYOK 路径和仍被上游门禁挡住的装配路径，并如实注明没有运行测试。

## 面 1 威胁模型
### Section A 结论
FACE 1: FINDINGS
### Section B 发现清单

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| Medium | [internal/gates/runner.go:208](internal/gates/runner.go#L208) | [docs/SECURITY.md:23](docs/SECURITY.md#L23)：「限制必须由 runner/OS 能力落实；不能仅靠提示词。」 | 将产物读取绑定到授权工作区目录句柄，在打开时拒绝越界；无法安全打开即拒绝收集，不入证据库。 |

### Section C 逐条证据与推理
读取 [docs/SECURITY.md:11](docs/SECURITY.md#L11)、[docs/SECURITY_EN.md:11](docs/SECURITY_EN.md#L11) 和 [docs/TEST_STRATEGY.md:29](docs/TEST_STRATEGY.md#L29)：双语基线一致区分了目录分离与进程隔离、内容完整性与身份认证。STRIDE 已列出路径逃逸，缺口在于 SEC-02 承诺的打开时检查未覆盖所有读取入口。读取 [internal/gates/runner.go:194](internal/gates/runner.go#L194)：产物收集按路径展开，只对末端做 `Lstat`，再按路径重新打开。它既未约束中间目录解析，也未把检查与读取绑定到同一对象。现有 `GuardExecutor` 缺必需能力而拒绝执行，这限制了当前可达性，不补足该读取边界。
### Section D 可证伪性审计
- 删除 **BoundDir** 的中间路径限制，`TestBoundDirMiddleSymlinkRefusedWithClosedErrno` 会失去预期拒绝；它不经过 Gate 产物收集。删除证据对象摘要校验，`TestVerifyManifestObjectsFailsClosed` 的 `tampered` 子项会红；删除事件归属检查，`TestVerifyManifestStateEventBinding` 的跨 run 子项会红。
- 删除 C6 的魔数检查或未提交槽判定，`TestC6SlotValidationOrderBadMagicIsCorruption` 或 `TestC6SlotValidationOrderUncommittedPayloadIsAbsent` 会红。基线已明示装配记录丢失后的漂移检测降级，不把它记为完整性已证明。
- UNCOVERED：未见 Gate 产物读取跨目录边界时同时断言拒绝、证据库无越界对象、任务不推进。现有 AT 编号不能代替这些断言。
### Section E 我无法验证的部分 + 什么会推翻我的判断
未验证原生平台上的完整 STRIDE 运行证据。若能指出所有 `Runner` 产物读取调用都由不可绕开的根目录绑定保护，并有针对这条读取路径的反例，才会推翻本项。

## 面 2 硬安全门禁
### Section A 结论
FACE 2: FINDINGS
### Section B 发现清单

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| High | [internal/console/runtime_agent.go:1245](internal/console/runtime_agent.go#L1245) | [docs/CONTRACTS.md:212](docs/CONTRACTS.md#L212)：「生产 admission 是纯 preflight 与不可变绑定：不得消费 request replay 身份、不得发布 R/C 或写 replay store、不得启动进程；首次 dispatch 资格与 replay 身份消费由 dispatch 路径在发布 R 之后依据 replay store 判定；admission 的结论是评估时刻的 preflight verdict，不是锁，可能在 dispatch 前失效，dispatch 在发布 R 之后、允许 launch 之前必须重新确认授权未被撤销且预算 reservation 仍 outstanding，任一不满足即 fail-closed，不得把先前的 admission 通过当作 launch 许可。」 | 在 launch 前向可信账本读取当前授权与预留；读取失败、撤销或已结算均保留 R，拒绝启动。 |
| High | [internal/chain/review.go:116](internal/chain/review.go#L116) | [docs/CONTRACTS.md:147](docs/CONTRACTS.md#L147)：「评审者与候选变更产生者是否同一身份、`evidenceRootHash`/`policyHash` 引用是否真实存在及 waiver 到期后不得复用，均由独立 checker 判定。」 | 发布前从执行者不可写的可信来源验证策略批准与候选证据根；没有有效绑定则拒绝，或停在待人工批准。 |
| Medium | [internal/console/runtime_agent.go:731](internal/console/runtime_agent.go#L731)、[internal/console/runtime_supervised.go:1642](internal/console/runtime_supervised.go#L1642) | [docs/CONTRACTS.md:150](docs/CONTRACTS.md#L150)：「`promotion-receipt.schema.json` 将父/候选 snapshot、evidence root、review receipt、停机与租约证据绑定为单个发布事务。」 | 发布前核对当前工作区的排他租约或等价写入隔离证明，将其真实记录绑定到本次 attempt；缺失或不一致则拒绝。 |
| Medium | [internal/gates/runner.go:73](internal/gates/runner.go#L73) | [docs/CONTRACTS.md:391](docs/CONTRACTS.md#L391)：「只读/本地可丢弃/外部写分离；S1 外部写拒绝，未知副作用暂停对账；诊断不变更锁/日志」 | 执行前从冻结 hook 策略导出有效副作用分类；分类缺失、矛盾或不允许时不能进入 `Execute`。 |

### Section C 逐条证据与推理
- **撤销独立性。** 读取 [internal/console/cli.go:562](internal/console/cli.go#L562) 和 [internal/adapters/agent_runner_replay_dispatcher.go:150](internal/adapters/agent_runner_replay_dispatcher.go#L150)：dispatcher 在正确的位置调用复核，但生产回调返回 CLI 启动时读取的授权副本。复核可发现随时间过期，不能发现初读后另行持久化的撤销。预算也只接到已装载对象；未证明可观察跨进程结算。
- **策略与批准。** 读取 [internal/console/runtime_agent.go:678](internal/console/runtime_agent.go#L678)、[internal/console/runtime_supervised.go:1598](internal/console/runtime_supervised.go#L1598) 和 [internal/chain/engine.go:254](internal/chain/engine.go#L254)：两条装配路径在 postflight 合格后本地计算策略摘要并返回 `approve`。同名拒绝是真守卫，但摘要合形不证明策略事先获独立批准；postflight 与批准不能互作原件。
- **租约与副作用。** 读取 [internal/guard/lease.go:104](internal/guard/lease.go#L104) 和 [internal/console/runtime_supervised.go:1491](internal/console/runtime_supervised.go#L1491)：租约库缺独立 verifier 会拒绝接管，但发布端的 `LeaseEvidence` 只由 run ID 派生，不证明当前独占。归还 hook 接线未填 `Effect`，而 runner 仅在有分类记录时调用准入。现版 CLI 另有 hook 定义和 runner 能力拒绝，不能把这个缺口说成现版已经执行了外部写入。
### Section D 可证伪性审计
- 删除能力检查，`TestPrepareRejectsUnavailablePolicy`、`TestRunnerDoesNotExecuteWhenPolicyIsUnavailable` 会红；删除启用后的沙箱拒绝，`TestAgentRunnerRunSandboxRefusesBeforePermit`、`TestSupervisedRunSandboxRefusesBeforePermit` 会红。契约明订的默认 `off` 不算绕过启用后的检查。
- 删除同路径独占或损坏 marker 拒绝，`TestProbePermitCrossProcessSingleClaim`、`TestProbePermitDamagedMarkerStaysSpent` 分别会红。删除就绪级别阻断，`TestReadinessProbeConclusionCannotAdvanceAcceptance`、`TestReadinessProductionRequiresAT23` 会红；这些是级别和格式断言，不是证据签发者认证。
- 删除同名自批、过期 waiver、postflight 事实绑定或接管 verifier 检查，分别由 `TestRejectSelfReview`、`TestRejectExpiredOrWrongCandidateWaiver`、`TestPostflightGateRefusesAnUnprovenPass`、`TestLeaseTakeoverFailsClosedWithoutVerifier` 反证。
- UNCOVERED：生产接线缺少这四项反例：初读后持久化撤销且 **R 保留、零启动**；策略未获独立批准且 **零发布**；租约证明缺失且 **无完成回执**；副作用分类缺失且 **零执行**。账本读取失败测试和显式外部写分类测试覆盖的是相邻分支。
- UNCOVERED：探针测试证明同一路径只认领一次，不证明独立签发或不同路径不能代表同一授权；`--probe-permit` 的调用方信任边界不能用 marker 的存在代证。
### Section E 我无法验证的部分 + 什么会推翻我的判断
未验证独立控制面是否在这两条装配路径外实施了不可绕过的实时账本、策略批准、排他写入或 hook 分类。能反驳上述判断的是实际接线与对应的零启动、零发布或零执行断言。未交付的路径仍应保持阻断。

## 面 3 秘密与凭据
### Section A 结论
FACE 3: FINDINGS
### Section B 发现清单

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| High | [internal/console/ai.go:94](internal/console/ai.go#L94) | [docs/SECURITY.md:60](docs/SECURITY.md#L60)：「网络默认拒绝，例外授权须列 host/port/purpose、提供方数据策略及期限。」 | 解析凭据前，由可信本机控制路径校验已批准的 profile、凭据引用、接收端与子进程身份绑定；不符即不读库、不启动。 |
| Medium | [internal/guard/process_windows.go:457](internal/guard/process_windows.go#L457) | [docs/CONTRACTS.md:216](docs/CONTRACTS.md#L216)：「adapter 只接收解析后的短暂字节副本，使用后清零；系统错误只返回稳定分类，不得拼接 credential blob。」 | 对 Windows 启动所用的原生环境块及编码中间缓冲保留所有权；成功或失败均及时清零，并核对凭据读取与终端输入的临时副本。 |
| Medium | [internal/adapters/ai_probe_copilot_executor.go:53](internal/adapters/ai_probe_copilot_executor.go#L53) | [docs/CONTRACTS.md:214](docs/CONTRACTS.md#L214)：「profile 只保存 provider 类型、base URL、模型、传输选项及不透明 `secretRef`；API key/token 不得进入配置、摘要、日志、错误或证据，运行时只能由系统 SecretStore 解析并临时注入子进程环境。」 | 只保留经过校验的结构化事实，不对未脱敏流生成持久证据；关闭未证明脱敏的子进程落盘日志，无法保证则拒绝 BYOK 探针。 |

### Section C 逐条证据与推理
- 读取 [internal/adapters/ai_provider.go:115](internal/adapters/ai_provider.go#L115)、[internal/adapters/ai_secret_windows.go:139](internal/adapters/ai_secret_windows.go#L139) 和 [internal/console/ai.go:66](internal/console/ai.go#L66)：profile 验证只检查引用形状和 HTTPS 语法，Windows 后端按当前用户的目标名解析。`ai check` 的配置选择和可执行文件路径没有独立核对凭据与接收方的批准绑定。单请求授权、HTTPS 和配置摘要均不能证明接收方有权持有凭据。
- 读取 [internal/adapters/ai_preflight.go:57](internal/adapters/ai_preflight.go#L57)、[internal/adapters/ai_probe_copilot.go:80](internal/adapters/ai_probe_copilot.go#L80) 和 [internal/guard/process_windows.go:221](internal/guard/process_windows.go#L221)：SecretStore 返回的 Go 切片、探针环境副本及 executor 的字节缓冲有清零；真正传给 Windows 的 UTF-16 环境块及转换副本没有。`CredFree` 调用只能据实记为释放，不能据此认定擦除。
- 读取 [internal/adapters/ai_probe_copilot_executor.go:92](internal/adapters/ai_probe_copilot_executor.go#L92) 和 [internal/adapters/ai_secret_unsupported.go:22](internal/adapters/ai_secret_unsupported.go#L22)：ProofRail 未把解析值放进探针 argv，非 Windows SecretStore 拒绝解析，结构化原因也不转录原始错误。但子进程可写临时日志，原始 stdout/stderr/usage 未经秘密检查即入摘要，目录清理不能证明从未落盘。未观察到实际泄漏。
### Section D 可证伪性审计
- 去掉 SecretStore 返回副本或探针调用环境副本的清零，`TestAIAvailabilityCheckUsesOnlyRequestedChannelAndClearsSecret`、`TestCopilotCLIProbeBuildsIsolatedDeepSeekSmoke` 分别会红。去掉缺秘密时的拒绝，`TestAIAvailabilityCheckBlocksMissingSecretWithoutProbe` 会红。去掉 explain 遮蔽，`TestAIConfigRoundTripAndExplainRedactsSecretReference` 会红。
- UNCOVERED：未见生产入口对未批准的凭据/接收方绑定断言零解析、零启动；也未见覆盖原生环境块的成功及失败清零、真实平台后端拒绝后零启动，以及合成秘密在子进程输出、日志及最终证据中的端到端不泄漏检查。
### Section E 我无法验证的部分 + 什么会推翻我的判断
未验证 Windows 原生调用和子进程自身的日志行为。可推翻相应发现的是生产入口已有不可绕过的批准绑定、原生缓冲的完整清零路径，或能证明日志与摘要只消费脱敏事实的端到端反例。

## 合成
FINAL REVIEW: FINDINGS

读取 [docs/ADR_REGISTER.md:58](docs/ADR_REGISTER.md#L58) 与 [docs/PRODUCT_REQUIREMENTS.md:109](docs/PRODUCT_REQUIREMENTS.md#L109)：三面各有不符合判据之处，不能登记为安全门通过，也不能以此替代 S1 Exit 评审。以上是只读判断，未执行命令、未变更文件。方案书定位回显夹带了许可节之外的片段，未采信；未读取 ⑤/⑥ 清单或 P79 报告结论。
```

## 8. 三面结论与合成

| 面 | 结论档 | 发现条数 | 最高严重度 |
|---|---|---|---|
| 面 1 威胁模型 | `FINDINGS` | 1 | `Medium` |
| 面 2 硬安全门禁 | `FINDINGS` | 4 | `High` |
| 面 3 秘密与凭据 | `FINDINGS` | 3 | `High` |
| **合成** | `FINAL REVIEW: FINDINGS` | 8 | `High` |

任一面 `FINDINGS` ⇒ 合成 `FINDINGS`；整体 `FINDINGS` ≠ S2 不可启动，但 ADR-047 的补偿义务以「完成评审并登记结论」为履行标准，见 §9。

**复审链**：复审第 1 轮 `RE-REVIEW: FINDINGS`（1 Medium：`SECURITY_EN` 第 7 行的时态未限定为历史状态）→ 整改（英文第 7 行改为 `Before P79 … had never been executed`，中文第 7 行同步加「在 P79 之前」限定）→ 复审第 2 轮（收窄输入）`RE-REVIEW: PASS`（0 条新发现）；⑦ 额度 1 终审 + 3 复审，实得 1 + 2（第 3 次复审未使用）。

## 9. 发现处置（本片不改产品代码）

- 本片为评估类，边界明确禁止改产品代码与准则。
- 8 条发现（3 High + 5 Medium）全部原样归档，逐条给出「建议承接方式」列。
- 是否立修复片、是否修订准则、是否把 ADR-047 翻「已补偿」——三项均由用户裁定，主控不得自行决定。
- 未修任何代码，故「发现是否被修复」不在本片验收范围；本片验收 = 评审完成 + 结论如实登记。

| 编号 | 面 | 严重度 | file:line | 建议承接 |
|---|---|---|---|---|
| F1 | 面 1 | `Medium` | `internal/gates/runner.go:208` | 另立修复片：产物读取绑定授权工作区目录句柄，无法安全打开即拒绝收集 |
| F2 | 面 2 | `High` | `internal/console/runtime_agent.go:1245` | 另立修复片：launch 前向可信账本读取授权与预留，读取失败 / 撤销 / 已结算均保留 R 并拒绝启动 |
| F3 | 面 2 | `High` | `internal/chain/review.go:116` | 另立修复片：发布前从执行者不可写的可信来源验证策略批准与候选证据根 |
| F4 | 面 2 | `Medium` | `internal/console/runtime_agent.go:731`、`internal/console/runtime_supervised.go:1642` | 另立修复片：发布前核对排他租约或等价写入隔离证明并绑定本次 attempt |
| F5 | 面 2 | `Medium` | `internal/gates/runner.go:73` | 另立修复片：执行前从冻结 hook 策略导出有效副作用分类 |
| F6 | 面 3 | `High` | `internal/console/ai.go:94` | 另立修复片：读取凭据前校验 profile、凭据引用、接收端与子进程身份绑定 |
| F7 | 面 3 | `Medium` | `internal/guard/process_windows.go:457` | 另立修复片：对原生环境块与编码中间缓冲保留所有权，成功失败均清零 |
| F8 | 面 3 | `Medium` | `internal/adapters/ai_probe_copilot_executor.go:53` | 另立修复片：关闭未证明脱敏的子进程落盘日志，无法保证则拒绝 BYOK 探针 |

## 10. 文档同步

| 对象 | 处置 | 依据 |
|---|---|---|
| `SECURITY{,_EN}` 注记行 | 已更新为「已完成 + 日期 + 评审者」 | 本片交付物 ②；ADR-047 影响面的加注义务 |
| `ADR_REGISTER{,_EN}` 的 ADR-047 | 已加 P79 完成注记，**未翻「已补偿」** | 用户裁定；状态词变更超出本片授权 |
| `PRODUCT_REQUIREMENTS{,_EN}` §7 | **未加注记** | 起跑消息列为核对项，是否加注由用户裁定 |
| `REMAINING_SLICES{,_EN}` P79 段 | 已登记步骤 0 事实并追加完成记录 | 本片交付物 ④ |

## 11. 审查输入包摘要与盲审隔离证明

- 输入包 = A / B / C / D / E 五组 + 附加件；逐件行数、字节、SHA-256 与过滤命中计数由 `input-list.md` 与 `01-arch/package-manifest.md` 断言。
- 消毒结论指向 `01-arch/package-manifest.md`：逐件登记消毒前 / 后 SHA-256，未改写件前后 SHA 相同。
- ⑦ 首轮盲态 = 未附 ⑤ / ⑥ 清单、未附主控工作笔记；首轮输入 = 起跑消息、ADR-047 原文、三面定义、消毒后输入包、方案书 §2 与 §6。
- ⑦ 只读且未执行命令：全部读数按「据输入包内 `<路径>`」口径引用。
- ⑦ 转录逐字：§7 围栏与归档 `04-reviews/07-security-r1.md` 逐行一致。

## 12. 可证伪性与变异检验

- 本片为文档 / 评估型，无产品代码改动；可证伪对象 = 文档事实与门禁读数。
- 版本元数据：若准则版本、SELFTEST 计数或三作用域读数与实际不符，步骤 0 锚点复核与 ⑩ 回执复跑立即红。
- 双语形状：若 CN / EN 行数、`##` 数量或逐行竖线符与加粗标记计数不一致，G3 镜像对称检查立即红。
- 围栏：若 `sixstage` / `text` / `artifacts` 围栏缺失或未闭合，围栏检查立即红。
- 工件围栏：若围栏路径与盘上文件不符或遗漏本片工件，工件围栏检查立即红。
- 引用可解析性：若章节互指（§7–§9、§14–§15）指向不存在的节，引用检查立即红。
- 编码：若任一文件非 UTF-8 无 BOM + LF，编码门禁立即红。
- 未做产品级变异检验：本片无产品代码改动，无守卫或断言可变异，故不适用。

## 13. 门禁结果

时点：本章读数取 ⑧a 定稿时点；变更集门禁的最终读数以 ⑩ 回执为准。

| 作用域 | 命令 | TOTAL_FAIL | 备注 |
|---|---|---|---|
| `index` | `node tools/gates/gate.js --all --scope=index` | 0 | 红项清零（G5(d) / G5(f) / G7-a / G7-f / G8-b 均 PASS）；读数文件 `gate-index-final.txt` |
| `tree` | `node tools/gates/gate.js --all --scope=tree` | 0 | 同上口径；读数文件 `gate-tree-final.txt` |
| `ci` | `node tools/gates/gate.js --all --scope=ci` | 0 | 与 CI 同口径（本地无 CI 变量时退化为上一提交区间） |
| `SELFTEST` | `node tools/gates/gate.js --selftest` | — | 全绿（exit 0） |

## 14. 环境与档位

| 角色 | 载体 | 模型 | 档位 | 设定时刻 | 是否经过暂停 |
|---|---|---|---|---|---|
| 主控 | `sol-orchestrator` | `DeepSeek V4.1 Flash (deepseek)` | High | 起跑前设定 | 无 |
| ① 架构 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | Max | 起跑前一次性设定 | 无 |
| ⑤ 预审 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | Max | 起跑前一次性设定 | 无 |
| ⑥ 独立扫描 | `independent-reviewer` | `MAI Code 2 Flash (copilot)` | 默认 | 起跑前一次性设定 | 无 |
| ⑦ 独立终审 | `independent-reviewer` | `GPT-6.2 Sol (copilot)` | Max | 启动前一次性设定 | 无 |
| ⑧ 文档 | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | High | 起跑前一次性设定 | 无 |

四字段（§8.4）= 设定时刻、档位值、是否暂停、暂停前后门禁结果；本片无暂停，暂停前后门禁结果不适用。§8.3 同模型限制不适用（⑦ 与主控不同模型）。

## 15. 成本与计量

| 角色 | 模型 | 轮次 | 备注 |
|---|---|---|---|
| ① 架构 | `DeepSeek V4 Pro (deepseek)` | 1 | 方案书 `01-arch/design-r1.md` |
| ⑤ 预审 | `DeepSeek V4 Pro (deepseek)` | 2 | 额度 1 + 1 用尽 |
| ⑥ 独立扫描 | `MAI Code 2 Flash (copilot)` | 2 | 额度 3 用 2；r2 `INDEPENDENT SCAN: PASS` |
| ⑦ 独立终审 | `GPT-6.2 Sol (copilot)` | 1 | 额度 1 终审 + 3 复审，本轮用 1 终审 |
| ⑧ 文档 | `DeepSeek V4.1 Flash (deepseek)` | 1 | 本报告对 |
| 批量辅助 | `GLM-5.3-Flash (glm)` | 2 | 枚举 / 摘录与数字核对 / 双语镜像检查 / MANIFEST 草案 |

⑦ 额度 = 1 终审 + 3 复审；该额度超出 §7.5「M-L 视同 M」字面标准、系用户 2026-11-19 直接裁定、**不得作为先例**。盲审不抽：§7.7 必抽三语义（所有权 / 停机 / 身份）本片不触。批量辅助 2 次、探针 0。

## 16. 明确未执行事项

- 未改产品代码。
- 未执行任何真实外部调用。
- 未跑 Ubuntu VM。
- 未做 ATT 级运行时安全实验。
- 未修 §9 所列任一条发现。
- 未给 `PRODUCT_REQUIREMENTS` 加注记（由用户裁定）。
- 未翻 ADR-047「已补偿」（由用户裁定）。

## 17. 先例禁用核对

本片引用了「⑦ 额度 1 + 3 复审」的用户裁定，并如实标注「不得作为先例」；未引用 OB-11 例外，未引用任何其它「不得作为先例」条款。

## 18. ⑧a 自查

- ① 成本表机械审计：有；证据 = §15 表与 `docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/` 各轮派单行。
- ② 「已处置」`file:line`：无；本片未把任何项标为「已处置」。
- ③ 旧口径 grep：有；命令与结果见 `tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/old-caliber-grep.txt`。
- ④ ⑥ / ⑦ 每轮 Info / Low 落入报告：有；对照 `04-reviews/06-scan-r1.md`、`06-scan-r2.md`、`07-security-r1.md`，三轮均无 Info / Low 项。

## 19. A-NN 编号完整性

本报告未使用 `A-NN` 编号；F1–F8 为 §9 的发现编号，实体定义在 §9 表内。

## 20. ⑧c 收尾清洁记录

- 全树逐文件编码核验：违规项全部位于 `docs/validation/evidence/**` 的冻结证据包内，非冻结主树违规为 0（口径与 §13 末行一致；冻结包按 §5.1 做自洽性核验，不做编码改写）。
- `tmp/` 断言：本片临时目录 `tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/` 内过程脚本与中间读数按 §5.1 逐件精确删除，删除后 `tmp/` 顶层仅 `.gitkeep` 与运行时缓存容器 `gates-cache`（§5.1 的例外集）。
- 未跟踪项：`git status --porcelain -uall` 的 `??` 项 = 0；本片变更集全部在索引中。
- 运行时缓存例外：`tmp/gates-cache/g4b/` 属 §5.1 明列的例外集，⑧c 不清除。

## 21. 异常与兜底记录（R / D）

| 编号 | 环 | 根因 | 处置 | 计入 |
|---|---|---|---|---|
| R-1 | ⑤ | 预审 r1 报 2 Medium + 4 Low（`sixstage` 契约缺位、门禁轨迹未披露、行号漂移等） | 整改方案书与报告后原位重跑 ⑤ 至 r2 | R = 1 |
| R-2 | ⑥ | 独立扫描 r1 报 2 High + 3 Medium（围栏缺位、行号漂移、静态件字节不一致、消毒登记不全） | 整改后重扫 ⑥ 至 r2 `INDEPENDENT SCAN: PASS` | R = 2 |
| — | ⑦ | 复审第 1 轮报 1 Medium（`SECURITY_EN` 时态）；复审属额度内义务环，按 §9.2 口径不计 R | 整改后收窄复审至 `PASS` | 不计 |

- **R 实测 = 2**（⑤ 环 1 次 + ⑥ 环 1 次；⑦ 复审环按 §9.2 口径不计入）。达 §9.2 的 `R≥2` 门槛，应标记「高风险切片」；过程中未显式打标记，**如实登记为本片偏离**（§5.2 高风险行要求的 ①②③⑤⑥⑦⑨ 均已实际执行，实质未缺环），交用户在 ⑩ 审核裁定。
- **D 实测 = 1**：⑧c 清理后，主控为落本次补正误建 `tmp/P79-fix/`（非 `tmp/<切片ID>/`）；发现后立即整目录删除，零残留（`tmp` 顶层 = `.gitkeep` + `gates-cache`；排除 `gates-cache/g4b` 的递归文件计数 = 1）。按 §5.1 自愈档登记为 **D = 1**；主控本人违规，**自愈判定待同轮独立审查层或用户确认**，未获确认前不视为自愈成立。

- **用户追认（A，2026-11-20）**：R = 2 未显式打「高风险切片」标记一事，用户追认接受本节的如实登记，**不重做、不另计 R/D**（裁决人 = 用户；日期 = 2026-11-20；理由 = 各环已实际执行且实质未缺环、重做无收益；依据 = 用户 2026-11-20 复核裁定）；**系用户裁定、不得作为先例**。
- **用户确认（B，2026-11-20）**：D = 1 的自愈成立（自报 + 发现即删 + 零残留，`tmp` 顶层已核验为 `gates-cache,.gitkeep`）；**D = 1 保留**并写入台账完成记录。

## 工件围栏

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/PRODUCT_REQUIREMENTS.md	repo
docs/PRODUCT_REQUIREMENTS_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/00-briefing.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/design-r1.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/package-manifest.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/05-prereview-r1.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/05-prereview-r2.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/07-re-review-r1.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/07-re-review-r2.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/07-security-r1.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/06-gates/gate-enc-scan.txt	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/06-gates/gate-index-final.txt	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/06-gates/gate-selftest-final.txt	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/06-gates/gate-tree-final.txt	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/MANIFEST.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/SHA256SUMS.txt	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/freedom-list.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/input-list.md	repo
docs/validation/s2-pre-independent-security-review.md	repo
docs/validation/s2-pre-independent-security-review_EN.md	repo
```
