# P79 ① 评审方案设计（design-r1）

- 切片：P79 `S2-PRE-INDEPENDENT-SECURITY-REVIEW`；分级 `[SLICE]`；类型 = 评估类（独立安全评审补评）；规模 M-L。
- 交付：① 架构师（deep-reasoner，DeepSeek V4 Pro Max）；日期 2026-11-19；编码 = UTF-8 无 BOM + LF。
- **性质声明**：本方案书只定义「评审怎么做」，不含任何评审判断；某面是否通过、是否存在缺陷，全部由 ⑦ 深度档（`GPT-6.2 Sol (copilot)` Max）作出。本方案书任何措辞不得被读作结论倾向；六项必覆盖对应 §2–§7。

## 1. 评审范围与依据

- 触发：ADR-047（`docs/ADR_REGISTER.md:58`）。义务三要素：S1 内不补评、豁免顺延至 S2；须在 S2 启动前完成一次独立安全评审并登记结论；到期日 = S2 启动日与 2027-03-31 中较早者，到期未补评 ⇒ 不得启动 S2 实测（回退条款）。
- 评审对象 = `PRODUCT_REQUIREMENTS.md` §7 的「独立安全评审结论」门禁（`PRODUCT_REQUIREMENTS.md:109` 阻断门禁句，实测复核该句确在 :109、:111 为 P76 加法式注记；历史状态 = 从未执行，佐证 `SECURITY.md:5`）。
- 评审三面 = ADR-047「未覆盖范围」逐字承接，**不得增删**：① 威胁模型面（`SECURITY` §1–§2 的资产与信任边界 / 滥用场景）；② 硬安全门禁面（fail-closed 前置条件的独立性）；③ 秘密与凭据面（SecretStore / BYOK 通道）。
- 对象区分：`docs/validation/s1-exit.md:59` 门表「独立评审」行 = T017 基线（Codex/DeepSeek 双主审），与本片**不同对象**、证据不可互换；本片不改 s1-exit 门表（含八门判据）。
- 边界（承接起跑消息第八节）：不改 S1 Exit 八门判据；不改 `SECURITY{,_EN}.md` 威胁模型基线（仅允许结论归档 + 注记）；不改 DELIVERY_DIRECTIVE / CONTRACTS / schema / fixtures / 产品代码 / CI；不新增 OB / ADR（用户裁定除外）；不做网络调用 / 外部命令 / 读 `%TEMP%` / 仓库外文件。

## 2. 三面逐面映射表（六项之 1）

| 面 | (a) 评审对象（文件 + 章节 / 函数级线索） | (b) 判据来源（精确出处） | (c) 结论档定义 |
|---|---|---|---|
| ① 威胁模型面 | `SECURITY.md` §1（资产与信任边界，正文 :11）、§2（写权限矩阵表，:15–:21）及其引用的滥用场景载体 §3 STRIDE 表（SEC-01…SEC-15，:29–:44）与 §3.1（C6 预建槽，:50–:54）；`SECURITY_EN.md` 同结构。评审动作 = 独立复核：资产清单完整性、§1 信任主张与 §2 矩阵行/列权限的一致性、STRIDE「防护及责任」列是否可由既有组件兑现、滥用场景是否闭合 | `SECURITY.md:5`（REQ-004/005/007/013/019/021-026/028）；`PRODUCT_REQUIREMENTS.md:109`（门禁句）；RFC §12、§16.7–16.8（设计来源）；`TEST_STRATEGY.md` §2 AT-01…AT-24 表（AT 编号解析，`SECURITY.md:46` 表尾句）；`ADR_REGISTER.md:58` 未覆盖范围 ① | 报告「威胁模型面」小节：该面 Section A–E 块 + 面结论行 `FACE 1: PASS\|FINDINGS`（模板见 §6） |
| ② 硬安全门禁面 | 判据侧：`docs/validation/sg6-gatelevel-consolidate.md`、`tools/gates/gate.js`、`tools/gates/lib/criteria/{g1,g5a,g7,g9}.js`。实现侧：`internal/gates/policy.go` `Prepare`(:30)/`validateRequest`(:48)/`validateHookEffectPolicy`(:96)；`internal/gates/runner.go` `Run`(:41)/`assess`(:115)/`scanOutput`(:172)；`internal/gates/readiness.go` 级别闭集(:70–83)/`ProbeReadiness`(:113)/`AcceptanceReadiness`(:118)；`internal/gates/probe_permit.go` marker 状态机(:74–91)/`NewProbePermit`(:126)；`internal/guard/lease.go` `Claim`(:82)/`Takeover`(:104)/`validTakeoverProof`(:120)；`internal/guard/sandbox.go` `VerdictStatus`(:21–34)；`internal/chain/postflight.go` `requirePostflightQualifiedReview`(:308)；`internal/evidence/verify.go` `VerifyManifestObjects`(:37)。评审动作 = 独立复核「fail-closed 前置条件的独立性」：每个拒绝路径的前置条件是否可被同一被防主体自行满足或绕过；「修改策略」与「独立批准」是否分离 | `SECURITY.md:23`（同段两句：「限制必须由 runner/OS 能力落实；不能仅靠提示词」与「执行者不得拥有修改策略和独立批准同一变更的合并权限」）、§2 矩阵行语义；`CONTRACTS.md` §4–5（Prepare/Inspect/Validate/Promote、handoff）、§7（admission 纯 preflight、dispatch 前重新确认授权未撤销且 reservation outstanding、任一不满足 fail-closed）；`TEST_STRATEGY.md` AT-05（拒评/自批拒绝）/AT-06/AT-18（硬门禁不可 waive）/AT-19；`ADR_REGISTER.md:58` 未覆盖范围 ② | 报告「硬安全门禁面」小节：Section A–E 块 + `FACE 2: PASS\|FINDINGS` |
| ③ 秘密与凭据面 | `internal/adapters/ai_secret_windows.go` `NewPlatformAISecretManager`(:48)/`Set`(:52)/`Resolve`(:77)/`Exists`(:89)/`Delete`(:101)/`windowsCredentialTarget`(:139)/`containsASCIIControl`(:150)；`internal/adapters/ai_provider.go` `ValidateAIProviderProfile`(:47)/`validHTTPSBaseURL`(:81)/`validDNSName`(:101)/`validAISecretReference`(:115)；`internal/adapters/ai_probe_copilot.go` `Probe`(:60)/`copilotCLIProbeEnvironment`(:134)/`clearCopilotCLIProbeEnvironment`(:150)；`internal/console/ai.go` `executeAICheck`(:33)/`configuredAIProfile`(:230)；`docs/OPERATIONS.md`（`ai check` 读取绑定 SecretStore 凭据、`secret set` 无回显）。评审动作 = 独立复核：secretRef 不透明性、凭据注入/清零生命周期、凭据是否可能进入配置/摘要/日志/错误/证据、BYOK 通道（Copilot CLI probe）凭据仅经子进程环境且用后清零、非 Windows fail-closed | `SECURITY.md` §4（:56–:58：秘密不经模型、不写配置/日志、secret-direct、外发扫描）；`CONTRACTS.md` §7 secretRef 段（「profile 只保存 provider 类型、base URL、模型、传输选项及不透明 secretRef；API key/token 不得进入配置、摘要、日志、错误或证据，运行时只能由系统 SecretStore 解析并临时注入子进程环境」）、`windows-credential:<target>` 段（`prfrail secret set` 仅已确认本地控制台无回显、`status` 只返回引用与存在性、adapter 只收短暂副本用后清零、非 Windows 平台 fail-closed）、「`config explain` 不得显示 secretRef/凭据存在性/解析值」「静态 validate/explain/preview 不读取 SecretStore」；`TEST_STRATEGY.md` AT-10（secret-direct 仅终端、secret canary 不入证据/模型）、AT-06（秘密扫描失败阻断）；`DEV_PLAN.md:119`（实现已完成句）；`ADR_REGISTER.md:58` 未覆盖范围 ③ | 报告「秘密与凭据面」小节：Section A–E 块 + `FACE 3: PASS\|FINDINGS` |

表内行号均为盘上定位线索；⑦ 引用时须按 §7.3 v1.53 贴原文定位断言（块头行号 + 逐字照录），定位或逐字任一不成立者按「线索」处理，不得采信。

## 3. 输入包与消毒（六项之 2）

入包清单（件数 / 行数 / 字节数以 `input-list.md` 实测为准：40 件、约 2.7 MB）：

- A 安全基线与需求 10 件（`SECURITY{,_EN}`、`PRODUCT_REQUIREMENTS{,_EN}`、`ADR_REGISTER{,_EN}`、`DEV_PLAN{,_EN}`、`TEST_STRATEGY{,_EN}`）；
- B 执行准则与快照 6 件（`DELIVERY_DIRECTIVE{,_EN}`、`CHECKLIST_v65{,_EN}`、`REMAINING_SLICES{,_EN}`）；
- C 判据来源 8 件（`sg6-gatelevel-consolidate.md`、`s1-exit{,_EN}.md`、`gate.js`、`g1/g5a/g7/g9.js`）；
- D 产品侧硬安全门禁实现 11 件（`internal/gates/*` 4 件、`internal/guard/*` 3 件、`internal/chain/postflight.go`、`internal/evidence/*` 3 件）；
- E 秘密与凭据面 5 件（`ai_secret_windows.go`、`ai_provider.go`、`ai_probe_copilot.go`、`internal/console/ai.go`、`docs/OPERATIONS.md`）；
- 附加件：`00-briefing.md`（起跑消息逐字）、`freedom-list.md`、本方案书 §2/§6 的契约性事实列（不含倾向性内容）、`01-arch/package-manifest.md`（见下）。

消毒规则（§7.3，逐字过滤，不得以「未发现」代替扫描）：① 过滤 IP 字面量（IPv4/IPv6，含内网段形态）；② 过滤用户名、主机名、SSH 目标串（如 `git@<host>:`）、`%TEMP%` 与仓库外绝对路径；③ 过滤凭据 / token / 密钥 / password= 形态；④ 过滤后以 `<redacted:类别>` 占位且**保持行数不变**（行数口径与 `input-list.md` 一致，便于行号定位不漂移）；⑤ 保留：仓库内相对路径 `docs/...`、`internal/...`（§7.3 明列），`tools/...`、`examples/...` 经主控确认非敏感后保留；⑥ 消毒结果落 `01-arch/package-manifest.md`：逐件登记 消毒前/后 的 SHA-256、行数、字节数与过滤命中计数。

材料化要求（§7.3 工具面标定：⑦ 只有只读工具，不得声称执行过命令）：凡要求 ⑦ 复核的读数必须由主控预先落成输入包内的静态材料——① `SECURITY{,_EN}` §1–§2/§4 与 `CONTRACTS.md` §7 secretRef 段的**带行号摘录**；② 内部实现文件的**函数级带行号切片**（函数头 + 关键拒绝分支，不整文件粘贴）；③ 机械读数静态件（`go test ./...` / `go vet ./...` 输出、SELFTEST 计数等）由主控在对应步骤执行后落盘为 `.txt` 并入包；④ ⑦ 只能写「据输入包内 `<路径>`」，不得写「执行了某命令」。

## 4. 输出契约（Section A–E，六项之 3）

⑦ 终审最终消息按 §7.2.2 五段（与 ⑦ 的契约同构）输出，**三面各一个独立 Section A–E 块** + 一个合成块：

- **Section A**：该面结论行 `FACE <i>: PASS|FINDINGS`（i = 1/2/3；取值限二值，见 §6）。
- **Section B**：发现清单，每行固定 `SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复`；`SEVERITY` ∈ {`Critical`,`High`,`Medium`,`Low`}（§4.3）；`file:line` 必须盘上可定位；无发现 ⇒ 空。
- **Section C**：逐条证据与推理（只引用输入包静态材料与契约原句）。
- **Section D**：可证伪性审计（硬性）——逐项给出「若删掉 / 改回某守卫，哪个测试或断言会红、哪个不会红」；列出**未覆盖的不变量**（`UNCOVERED` 独立成句标注；不计入发现数、不影响结论行）。
- **Section E**：「我无法验证的部分」+「什么会推翻我的 PASS」（防橡皮图章）。
- 一致性硬约束（同 §7.2.2 对 ⑥ 的口径）：Section B 无发现 ⇒ 该面结论行必须 `PASS`；有发现 ⇒ `FINDINGS`。
- 合成块：不含新发现，只按 §6 合成规则给出整体结论行。

格式偏离处置（§7.2 同构）：未按五段 / 无面结论行 ⇒ 重发 1 次（计入 ⑦ 额度）；再偏离 ⇒ 停止并升级用户裁决（§7.5 表后句），不得自行修复格式。

- **⑥ 转录契约（§11.3 / `G7-f`）**：⑥ 最终消息由主控逐字转录进报告「审查结论摘要」节的 `sixstage` 围栏（Section A–E 五标题 + 块末结论行 `INDEPENDENT SCAN: PASS|FINDINGS`）；本片 ⑥ 必做，故不使用 ⑥ 跳过声明。

## 5. 独立性留痕与转录契约（六项之 4）

- **档位留痕（§8.4）**：⑦ 的档位在报告「环境与档位」表按四字段登记——设定时刻、档位值（Max）、是否经过暂停、暂停前后门禁结果；主控与 ⑦ 非同一模型，§8.3 的同模型档位限制不适用。
- **判断主体唯一**：全部评审判断由 ⑦ 深度档 `GPT-6.2 Sol (copilot)` Max 作出（§5.2「安全」深度判据触发 + 用户 2026-11-19 同轮授权取 Max，§8.2 取档授权）；主控（DeepSeek V4.1 Flash High）只做编排 / 材料化 / 转录 / 门禁 / 回执，不得以自己的读感补写或软化任何发现。
- **流水线顺序与锚定隔离**：⑤ 预审报告初稿（1+1）→ ⑥ MAI Code 2 Flash 独立扫描报告初稿的输出契约符合性（≤3）→ ⑦ 终审盲态；串行执行，**⑤/⑥ 清单一律不得进入 ⑦ 首轮输入**。
- **首轮盲态**（§4.2 ⑦ 行 + §7.3 锚定硬隔离）：⑦ 首轮输入 = 起跑消息原文 + ADR-047 行原文 + 三面定义 + 消毒后输入包 + 本方案书 §2 判据来源列与 §6 结论档模板（契约性事实）——不附任何 ⑤/⑥ 清单、不附主控工作笔记、不附本方案书其余节内容；复审轮方可附清单并标注「完成独立扫描后再读」。
- **转录契约**（§11.3 同构，v1.54 口径）：⑦ 最终消息由主控**逐字**转录进报告「审查结论摘要」节围栏块（Section A–E 标题、各面结论行、合成结论行、原文逐字），**不得按要点重建、不得重组空行、不得改写措辞**；围栏外只放主控注记与指针。
- **留痕**：报告「审查输入包摘要」节含盲态证明（首轮是否附清单、输入是否已消毒、package-manifest 指针）；R7.1 计量按单一真值源逐字段登记（⑦ = 1 终审 + 3 复审，用户 2026-11-19 裁定 L 档，如实标注「超出 §7.5 M-L 视同 M 字面标准、系用户直接裁定、不得作为先例」）；复审轮按 §7.5 收窄输入（第 2 次只给修复 diff + 前轮发现 + 机械读数，第 3 次只核上一轮 Medium+ 闭合路径）。

## 6. 结论档定义与合成口径（六项之 5）

- 每面结论档 = 该面 Section A–E 块 + 面结论行 `FACE <i>: PASS|FINDINGS`，归档于报告对应小节；`PASS` = 该面未发现违反判据的缺陷（不证明安全，只记录「本片补评未发现」）；`FINDINGS` = 至少一条发现，每条附 `file:line` + 最小修复。
- 首读二值约束：本片为评估类首读、不设修复轮，⑦ 面结论行限 `PASS|FINDINGS`；若 ⑦ 输出 `PASS WITH FIXES|PARTIAL`，主控不自行改写，原词转录并按 §4.1「PARTIAL 视同 FINDINGS 子类」升级用户裁决（修复动作超出本片边界）。
- **三面合成口径**（机械规则，写入报告并作为验收断言）：整体结论行由 ⑦ 给出且必须与三面行机械一致——三面均 `PASS` ⇒ 整体 `FINAL REVIEW: PASS`；任一面 `FINDINGS` ⇒ 整体 `FINAL REVIEW: FINDINGS`。**任一面的 PASS 不得被读作整体 PASS**；整体 PASS 也不等于「S2 可启动」（S2 启动仍受其它门约束）或「安全已证明」（UNCOVERED 项照录）。
- 结论落地（起跑消息 §七 必改 4 项）：新报告 `docs/validation/s2-pre-independent-security-review{,_EN}.md`（§11.3 结构 + artifacts 围栏 + 双语严格镜像）；`SECURITY{,_EN}.md:5` 注记更新「已完成 + 评审日期 + 评审者」；`ADR_REGISTER{,_EN}.md:58` 完成注记（是否翻「已补偿」由用户裁定，未裁定只加注记）；`REMAINING_SLICES{,_EN}.md` P79 段完成记录。核对项：`PRODUCT_REQUIREMENTS{,_EN}.md` §7 阻断门禁句加注记与否由用户裁定。

## 7. 可证伪性与验收断言（六项之 6）

可证伪锚（⑦ Section D 必须针对盘上反例逐项回答；主控 ⑨ 按 §7.7 A11 对称抽查执行）：`internal/gates/policy_test.go` `TestPrepareRejectsUnavailablePolicy`；`internal/gates/probe_permit_test.go` `TestProbePermitFailureDoesNotResetQuota` / `TestProbePermitDamagedMarkerStaysSpent`；`internal/guard/lease_test.go` `TestLeaseTakeoverRequiresFreshStopAndAuthorization`；`internal/gates/readiness_test.go` `TestReadinessProductionRequiresAT23` / `TestReadinessProbeConclusionCannotAdvanceAcceptance`；`internal/gates/runner_test.go` `TestSecretScannerRejectsKnownCredentialMaterial`；`TEST_STRATEGY.md` AT-10 / AT-18 断言句。

机械验收断言（收口前逐条核验）：

1. 报告三面各出现独立结论档块（Section A–E 五标题齐全 + `FACE <i>` 行）。
2. 每条发现满足四段格式且 `file:line` 在盘上可定位、被违反句逐字存在于该行。
3. 整体结论行 ∈ §4.1 ⑦ 闭集且与三面合成规则机械一致。
4. 双语报告形状全等：三面结论行、发现条数、`SEVERITY|file:line` 集合逐条对应，仅语言镜像。
5. ⑦ 转录块存在且逐字（对照派单回执原文）；盲态证明 + 消毒 MANIFEST 存在。
6. 消毒复扫 0 命中（IP/用户名/主机名/SSH 路径/凭据/token 正则）。
7. §七 必改 4 项全部落地；s1-exit 门表、`SECURITY` 威胁模型基线、准则 / CONTRACTS / schema / fixtures / 产品代码 / CI 零改动（git diff 域检查）。
8. `sixstage` 围栏存在且块末结论行 ∈ §4.1 闭集（`G7-f` 绿）；`artifacts` 围栏与磁盘变更集逐一相同（`G5(d)` / `G7-b` 绿）。

## 8. 风险与停点

- **体量风险**：输入包 40 件约 2.7 MB ⇒ 按 §3 材料化切片（基线/需求整件、内部实现函数级切片）；派单提示词要求三面分块、每块独立五段；消息截断按 §7.2 截断重发处置并计入额度。
- **首读即 FINDINGS 的风险**：发现可能指向产品代码 / 准则语义 ⇒ 本片禁改代码与准则，发现只归档；是否另立修复片或修订准则由用户裁定（升级边界，主控不得自行决定）。
- **停点**（逐条触发即停，不得绕行）：⑦ 额度超 1+3 ⇒ §7.5 升级；格式偏离重发后仍偏离 ⇒ 升级；R 达 2 ⇒ §9.2 高风险标记并停下；任一面结论行表外形态 ⇒ 升级；ADR-047 到期日（S2 启动日与 2027-03-31 先到者）逼近仍未收口 ⇒ 升级。
- **边界触碰自检**：本方案书不改任何受保护对象、未新增 OB / ADR、未执行任何命令 / 网络调用、未读 `%TEMP%`、未创建仓库外文件 ⇒ 不触边界。
