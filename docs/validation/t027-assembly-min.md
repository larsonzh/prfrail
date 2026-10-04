# T027-ASSEMBLY-MIN — 最小可跑装配（A 段：拒绝路径）验证报告

> 切片 `T027-ASSEMBLY-MIN`（产品实现类，A 段） | 分级 `[SLICE]` | 规模 M | 日期 2026-10-03。
> **状态：⑨ 原生验证、⑧b 定稿回写、⑧c 收尾已执行；⑩ 停点报告见第 13 节；提交与推送已按用户同轮授权执行（CI 结果与修复见第 13 节）。**

[English](t027-assembly-min_EN.md)

## 1. 变更概要（含执行自由度清单）

- **本片做了什么**：交付 A 段「最小可跑装配」——把已被冻结的 AgentRunner 构件装配成一条可跑的拒绝路径链（候选记录、能力记录、enforcement 记录、可用性记录、授权台账、成本台账、pinned CLI 可执行文件与版本共 8 项输入 flags 全部在场时才路由到装配入口）。新增 3 个生产文件（`internal/console/runtime_agent.go`、`internal/console/runtime_noop.go`、`internal/adapters/agent_runner_postflight.go`）与 2 个测试文件，并修改 `internal/console/runtime.go`、`internal/console/cli.go` 的既有分派。
- **本片没做什么**：不做真实 enforcement 落地（属 B 段）；不改 `internal/chain/**`、既有 adapters 与 console 文件、`schemas/**`、CONTRACTS、fixtures、workflows、`go.mod`；不新增任何哨兵；不为 Windows 提供正向全链证据。
- **产物**：本报告、自由度清单、上列的 5 个 Go 文件、`tmp/assembly-min/DESIGN.md`（方案书与三轮整改记录，工作区专用）、变异脚本与证据读数（`tmp/assembly-min/`，工作区专用）。
- **执行自由度清单**：`docs/validation/evidence/T027-ASSEMBLY-MIN-freedom-list.md`（F-1 至 F-6）。
- **一句话结论**：装配在拒绝路径上可跑、fail-closed 面在第三轮审查后已收紧（停机复核改为「两份 pid 报告都必须存在」，两侧解析器合一）；一个硬门（F-02 的根因）在本片范围外，按用户预授权以**记录态**收口。

## 2. 交付面与逐条落点

| 交付项 | 落点 | 说明 |
| --- | --- | --- |
| 装配入口（拒绝路径） | `internal/console/runtime_agent.go`（`ExecuteAgentRunnerRun`） | 构造顺序固定为 marker → replay store → 拒绝重叠根 → snapshot / launcher / dispatcher / admission / intent / postflight / publisher；任一前置失败即拒绝 |
| 失败路径单入口 | 同上（闭包 `fail`） | 七个 spawn 后失败返回全部经该闭包做 best-effort 兜底停机，唯一裸失败返回位于闭包内 |
| 停机复核口径 | 同上（`managedProcessesGone`） | 两份 pid 报告都必须存在；缺失、不可解析、存活性答案失败均为 fail-closed |
| 后置校验端口 | `internal/adapters/agent_runner_postflight.go` | 六条腿：manifest、diff、log、usage、process-alive、secret scan；产物摘要必须先与冻结事实绑定 |
| 非 AgentRunner 路径保持原状 | `internal/console/runtime_noop.go` | noop 运行时与其 fail-closed 的 `RunPostflight` 桩整体迁移，语义不变 |
| CLI 表面接入 | `internal/console/cli.go`（`executeRun`） | 8 项 flags 全在场且存在可执行步骤时路由到装配入口；否则维持原有报错 |
| 差异护栏测试 | `internal/console/runtime_agent_e2e_test.go`、`internal/adapters/agent_runner_postflight_test.go` | 装配正反向腿、后置校验七类篡改、解析矩阵、兜底闭包行为与源码级不变量 |

## 3. 异常与兜底记录

1. **① 的传输故障（4 次失败）**：架构子代理连续 4 轮返回传输错误；按用户裁定选「稍后重试」——不开代理、间隔两分钟以上、缩载荷；第 5 次成功。重试额度未用尽，如实记录。
2. **DR-8 候选（Windows 临时目录改名抖动）**：`internal/snapshot` 的 `TestRestoreWithHardlinksAndPermissions` 在首次并行全量测试中红（Windows 临时目录 rename 报 `Access is denied`），单包隔离复跑绿、再次全量绿；该包不在本片改动集内 ⇒ 判为环境抖动而非本片引入，首跑红原样保留为证据。
3. **主控自误（须如实留痕）**：为给 ⑤ 准备「改动前」快照，主控把 9 个文件的副本放进 `tmp/assembly-min/before/`；其中 `.go` 副本被 `go build ./...` 当作新包编译，导致基线门禁转红。处置：目录改名为 `_before`（下划线前缀被 Go 工具链忽略）、脚本同步更新、基线恢复绿。
4. **主控自误（须如实留痕）**：首次运行扩展后的变异审计时，M16 的还原写入在 Windows 文件锁（`EBUSY`）下失败 ⇒ M16 变异残留在 `agent_runner_postflight.go`。处置：依据「该变异是单次出现的确定性字面替换」做逆推复原（复原后无变异标记、`gofmt` 输出为空），给审计脚本加上原子写入、`finally` 强制还原与 `EBUSY` 重试，随后全量重跑得到逐字节还原为真。
5. **主控自误（须如实留痕）**：首轮 ④ 之后曾把「13 ok / 3 FAIL」的计数说错——其中 3 条是 marker 行，实际只有 1 个测试失败；已在报告中更正。
6. **兜底 = 2**（第 3、4 项），其余无需兜底；本片未发生任何数据破坏或未授权写入。
7. **门禁判红即整改（两处）**：新增报告首次跑门禁时 `G4b` 报出未登记的新 CJK 字符、`G7-a` 报出报告缺产物块；处置是**改写措辞**与**补齐产物块**，未加白名单、未改判据。

## 4. 执行流水线

| 环 | 角色 | 状态 | 产物 |
| --- | --- | --- | --- |
| 步骤 0 预置 | 主控 | 已执行 | `DELIVERY_DIRECTIVE{,_EN}` 升 v1.17、台账补登与 OB 条目 |
| ① 架构 | V4 Pro | 已执行（v1 与 v2 两版） | 判定方案与装配设计（A 段范围声明、拒绝路径枚举） |
| ②b 实现 | Flash（主控选型） | 已执行（含 r3、r4 两轮整改） | 5 个 Go 文件、测试与变异锚点 |
| ③ 测试 | V4 Pro | 已执行 | 正反向腿、篡改矩阵、解析矩阵、源码级不变量 |
| ④ 主控集成与门禁 | 主控 | 已执行（多轮复跑） | 见第 9 节 |
| ⑤ 预审 | V4 Pro | 已执行（三轮，第三轮为额度例外） | 见第 6 节 |
| ⑥ 独立扫描 | MAI Flash | 已执行 | 见第 6 节 |
| ⑦ 终审与复审 | Codex（不降级） | 已执行（三轮用满） | 见第 6 节 |
| ⑧a 文档初稿 | 主控 | 已执行 | 本报告与自由度清单 |
| ⑨ 原生验证 | 主控 | 已执行 | 三作用域与编码的只读复跑（见第 9 节） |
| ⑧b 文档定稿 | 主控 | 已执行 | 台账回写（本节下方与台账双语的切片段）与状态行更新 |
| ⑧c 收尾清洁 | 主控 | 已执行 | 临时物清理核对与工作区状态核对 |
| ⑩ 停点 | 主控 | 已完成（提交与推送经同轮授权） | 第 13 节 |

## 5. 环境与档位

- 平台：Windows 本机（Ubuntu CI 作为正向全链的唯一平台）；档位 `[SLICE]`；主控档位为标准。
- **⑦ 用满三轮且不降级**：第 1 次盲审、第 2 次终审、第 3 次复审，模型固定 `gpt-5.3-codex`；用户红线：不开第 4 轮。
- **⑤ 第三轮为额度例外**：按 OB-37 元规则以 ADR-018 独立留痕（收窄输入、不得作为先例）。
- 探针 = 0（未发起任何真实外部调用）；试点顺延四要素随 ⑩ 一并记录。

## 6. 审查结论摘要（⑤⑥⑦）

- **⑤ 第 1 轮**（V4 Pro）：`PASS WITH FIXES`——F1 为安全相关 Medium（产物未与冻结事实绑定），另有 Medium 与若干 Low；全部在片内整改，并补上对应变异。
- **⑤ 第 2 轮**（V4 Pro）：`PASS WITH FIXES`——无 Medium 以上；遗留 Low 与两条观察项在本片处置。
- **⑤ 第 3 轮**（V4 Pro，额度例外 ADR-018）：`PASS WITH FIXES`——无 Medium 以上；两条 Low 分别以措辞修正与新增变异条目处置。
- **⑥ 独立扫描**（MAI Flash，盲扫）：`PASS`，零发现；因该轮材料未落盘为独立文件，本报告只作如实声明。
- **⑦ 第 1 轮盲审**（Codex，不附 ⑤⑥ 清单）：1 High、1 Medium（硬门）、1 Low——停机复核过弱、spawn 后身份记录失败不回收、注释与实现不符。
- **⑦ 第 2 轮终审**（Codex，附 ⑤⑥ 清单与本轮整改）：`FINAL REVIEW: FINDINGS`——M-1（两侧解析器漂移，实测为 console 放行方向）与 M-2（兜底只覆盖首个非 awaiting 分支，硬门未闭）。
- **⑦ 第 3 轮复审**（Codex）：`RE-REVIEW: PASS WITH FIXES`——M-1 已关闭；F-02 部分关闭、仍判 Medium（按预授权记入记录态）；L-3 维持现状可接受；新增一条 Low（源码级不变量的鲁棒性）登记为 OB-61 候选。

## 7. 审查输入包摘要（含盲审隔离证明）

- **第 3 轮复审的输入包**：`tmp/assembly-min/review-pack-r3.md`（索引）、`change-set.diff`（全切片）、`r3-remediation.diff`、`r4-remediation.diff`、`mutation.txt`、`pre-review-r1.md`、`pre-review-r2.md`、门禁与测试读数、`v3.txt` + `v4.txt` + `snapshot-flake.txt`（DR-8 三份证据）、`contracts-excerpt.txt`。
- **前两轮的输入包**：第 1 轮为盲审（只给变更集与设计，不附任何 ⑤⑥ 清单）；第 2 轮附 ⑤ 第 1、2 轮清单、⑥ 结论与三轮证据读数。
- **每轮整改 diff 的来源**：r3 diff 由「r3 前状态的逐字节重建」与工作区之差生成；r4 diff 由审计前快照与工作区之差生成，并断言参照侧（adapters 源文件）零改动；两份生成脚本都做结构断言，读数留档。
- **消毒与隔离**：送审材料已过滤用户名、主机名、IP、凭据与令牌；⑦ 第 1 轮不做清单投喂以保证盲审独立性。
- 本片触及状态机、授权契约与证据模型语义 ⇒ 必抽盲审条件成立，已按第 1 轮盲审执行。

## 8. 可证伪性（机械检查与变异检验）

- 机械检查：见第 9 节（门禁与全量测试）。
- **变异检验**：脚本对每一条变异施加唯一替换、跑聚焦测试、与未变异基线比对逐测试判定、再逐字节还原（SHA-256 比对）。终态：**18 条 RED、3 条如实标注为被取代（M9、M11、M15），逐字节还原为真**；基线为 23 个测试、22 通过、1 跳过。
- **被取代的三条不是删除而是如实登记**：M9 与 M11 的旧弱守卫被「两份报告都必须存在」取代；M15 的锚点在两侧解析器合一后不复存在，且严格解析器下删掉错误判空是空操作。
- **不可机械判定的部分（如实）**：源码级不变量测试是文本形态计数与字符串匹配，理论上可被等价改写规避 ⇒ 已登记为 OB-61 候选。

## 9. 门禁结果

- `node tools/gates/gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0`（`changed=16`），含基线 `gofmt` / `go build` / `go vet` / `go test` 全过与编码核验全过。
- `G3(a)` 对称：本报告对为新增文件对且两侧新增行数相等、零删除（以当场门禁读数为准）；`G4b`：新增 CJK 字符已全部登记、未知集合为空；`G7-a` 至 `G7-d`：产物块在场且产物行全部可解析。
- `node tools/gates/gate.js --all --scope=index` ⇒ `TOTAL_FAIL=0`（`changed=0`，暂存集为空，判据按设计报空集而非失败）。
- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 248/248`。
- `go test -count=1 ./...` ⇒ 14 个包通过、0 失败（与基线包集合一致）。
- 编码：本片全部改动文件为 UTF-8 无 BOM（Go 与 JS）或带 BOM（Markdown），行尾 LF。

## 10. 已知差异与覆盖缺口

1. **F-01(3) 的替代口径（须逐字）**：`InspectProcess` 报错 ⇒ **视为进程已消失**。理由：字面口径（报错即 fail-closed）在双平台被反例证伪——已正常退出的 pid 在 Windows 报 `ERROR_INVALID_PARAMETER`（实测行：`AgentRunner evidence invalid: cannot observe pid 49656 from stub.pid: The parameter is incorrect.`）、在 Linux 上报 `/proc/<pid>` 的 `ENOENT`，按字面口径会让每一次已完成运行都无法验证，故与运行侧既有读法保持一致。fail-closed 面明确收缩为：存活性答案失败与条目不可解析。
2. **F-02 记录态（用户预授权，须逐字）**：若复审仍将 F-02 判为 Medium（即缓解、未关闭），用户预先接受该状态作为本片最终记录：本片以「M-1 已修 + F-02 缓解已最大化 + OB-57 跟踪根因」收口，不再为 F-02 开第 4 轮，不再扩大范围；这不是破 §5.3 纪律，而是用户明确接受根因在本片范围外的 Medium 作为记录态。复审的判定为「部分关闭、仍为 Medium」，故记录态生效。
3. **DR-8 首跑红的三份证据（须逐字）**：`tmp/assembly-min/v3.txt`（首跑红：Windows 临时目录 rename 报 `Access is denied`）、`tmp/assembly-min/v4.txt`（单包隔离复跑绿）、`tmp/assembly-min/snapshot-flake.txt`（抖动结论摘录）；该包不在本片改动集内。
4. **主控 EBUSY 自误与复原（须逐字）**：首次运行扩展后的变异审计时，M16 的还原写入在 Windows 文件锁下失败，变异残留在 `agent_runner_postflight.go`；主控依据该变异是单次出现的确定性字面替换做逆推复原，并给脚本加上原子写入、强制还原与重试，随后全量重跑得逐字节还原为真。
5. **CLI 八项 flags 无 CLI 级测试（须逐字）**：`executeRun` 的 8 项 flags 路由只有静态生产路径与装配腿担保，没有 CLI 级测试；真实调用者由静态生产路径与装配正反向腿共同担保。
6. **三个非逐字节相等的 parity 字段（须逐字）**：差异测试中 `OutputManifestHash`、`Facts.ManifestHash` 与 `Facts.ProcessStopEvidenceHash` 不逐字节相等，改用替代断言（字段在场、非空、摘要形状正确）；若将来 `snapshot.CapturedAt` 变成可注入，应把这三项收紧为逐字节相等。
7. **F-01 对真实 CLI 的误判风险（须逐字）**：本片的两份 pid 报告判据依赖候选 CLI 恒定写出父子两份报告；现行 stub 只在派生子进程时写子报告，真实 CLI 接入（B 段）时须评估该判据是否会把已完成的运行误判为不可验证。

## 11. 候选 OB（本片不落 §12.4）

| 候选 | 内容 | 来源 |
| --- | --- | --- |
| OB-56 候选 | `AgentRunnerPinnedCLIRun.Run` 无生产调用者、缺非 spawn 等待入口；本片以非 spawn 等待入口闭合该缺口 | 本片 A6 遗留说明 |
| OB-57 候选 | 兜底未闭的根因：spawn 成功后身份记录失败即返回且不回收；根因在 `internal/adapters/agent_runner_replay_dispatcher.go`，超出 `T027-ASSEMBLY-MIN` 范围，另立切片修复 | F-02（⑦ 三轮一致） |
| OB-58 候选 | 把 pid 报告摘要纳入冻结事实并与后置校验请求绑定（触冻结事实与契约域） | F-01(2)（用户裁定另立切片） |
| OB-59 候选 | CLI 契约缺口：stub 与真实 CLI 只在派生子进程时写子报告；在该缺口闭合前，生产正向链会因缺子报告而不可验证；测试夹具补种只证落点、不证生产侧有产出者 | F-01(1) 的生产后果（用户裁定另立切片） |
| OB-60 候选 | 方案书模板增两节——① 跨包共享契约：列明跨文件共享项（常量 / 解析函数 / 判定顺序 / 错误分类），每项定义权威侧与核验方式，配机械核验；② 失败路径枚举：对每个可能失败的操作，列「失败后由谁负责 / 后果 / 如何证明」，逐条穷尽，不留「其他失败同理」。适用裁剪：S 档不要求；M 档以上或跨包改动要求。 | 用户预保留，来源本片 M-1 与 F-02 的设计教训 |
| OB-61 候选 | 源码级不变量测试为文本形态计数与字符串匹配，理论上可被等价改写规避 ⇒ 建议升级为 AST 级结构断言 | ⑦ 第 3 轮复审新增 Low |

## 12. 未完成清单与下一步

- ⑨ 原生验证：已执行（三作用域与编码的只读复跑，读数见第 9 节）。
- ⑧b 文档定稿：已执行（台账双语的切片段与完成记录行、本报告状态行更新）。
- ⑧c 收尾清洁：已执行（编码与残留核对、临时物保留口径核对）。
- ⑩ 停点：已完成（提交与推送经用户同轮授权；CI 结果与修复见第 13 节）。
- commit / push 已获同轮授权并执行（首次提交 `68a1ef0`，含步骤 0 的 v1.17）；不推 gitee；CI 修复提交见第 13 节。

## 13. ⑩ 停点报告

- **⑩ 停点结论**：① 至 ⑧a、⑨、⑧b、⑧c 均已执行；**停在提交授权之前**——不 commit、不 push、不推 gitee，等待用户同轮授权。
- **产物清单与路径**（逐项）：
  - 实现：`internal/console/runtime_agent.go`（装配入口、单入口兜底闭包、停机复核口径）、`internal/console/runtime_noop.go`（noop 运行时与其 fail-closed 桩）、`internal/adapters/agent_runner_postflight.go`（六腿后置校验端口）、`internal/console/cli.go`（八项 flags 路由）、`internal/console/runtime.go`（既有分派保持原状）
  - 测试：`internal/console/runtime_agent_e2e_test.go`、`internal/adapters/agent_runner_postflight_test.go`（含解析矩阵、闭包行为与源码级不变量）
  - 文档：本报告对与自由度清单；步骤 0 的 `DELIVERY_DIRECTIVE{,_EN}`（v1.17）与 `ADR_REGISTER{,_EN}`（ADR-018）；`tmp/assembly-min/`（方案书、三轮整改记录、变异脚本与证据读数，工作区专用）
- **门禁逐项**：`tree` 为 `TOTAL_FAIL=0`（`G3(a)` 报告对对称、`G4b` 未知 CJK 为空、`G7-a` 至 `G7-d` 全 PASS 且产物行可解析）；`index` 为 `TOTAL_FAIL=0`；`SELFTEST PASS 248/248`；全量测试 14 个包 ok 且 0 失败；编码核验全过（Markdown 带 BOM、Go 与 JS 无 BOM、行尾 LF）
- **⑤⑥⑦ 结论行与额度**：⑤ 三轮均 `PASS WITH FIXES`（第 3 轮为额度 2/2 之外的例外，ADR-018）；⑥ `INDEPENDENT SCAN: PASS`（盲扫，0 发现）；⑦ 盲审 `FINDINGS`、终审 `FINAL REVIEW: FINDINGS`、复审 `RE-REVIEW: PASS WITH FIXES`（额度 3/3 用满，不开第 4 轮）
- **成本与计量**（逐行）：
  - ① 架构 2 版；②b 实现 3 轮（r2、r3、r4）；③ 测试 1 轮；① 的传输失败按用户裁定处理（不开代理、间隔两分钟以上、缩载荷，重试额度未用尽）
  - ⑤ 3 次（V4 Pro，含 1 次额度例外）；⑥ 1 次（MAI，盲扫）；⑦ 3 次（Codex，不降级）；真实外部调用 0 次
  - 变异条目 21 条（18 条红、3 条如实标为被取代），逐字节还原为真；聚焦测试基线 23 项（22 通过、1 跳过）；变异审计重跑 4 次（首次因 Windows 文件锁中断并如实留痕）
  - 授权台账与成本台账由装配入口按既有记录文件只读消费，本片不产生新的计费动作
- **8 条意图登记路径**（仅登记路径、不暂存内容；授权提交时转为真实暂存）：
  - `docs/validation/t027-assembly-min.md`
  - `docs/validation/t027-assembly-min_EN.md`
  - `docs/validation/evidence/T027-ASSEMBLY-MIN-freedom-list.md`
  - `internal/adapters/agent_runner_postflight.go`
  - `internal/adapters/agent_runner_postflight_test.go`
  - `internal/console/runtime_agent.go`
  - `internal/console/runtime_agent_e2e_test.go`
  - `internal/console/runtime_noop.go`
- **试点顺延第 6 次（四要素）**：裁决人 = 用户；日期 = 2026-10-03；理由 = 产品实现类切片继续执行而不进入试点；依据 = 本轮授权原文。**不得作为先例**。
- **步骤 0（v1.17）的处置口径**：`DELIVERY_DIRECTIVE{,_EN}` 的 v1.17 这一行、§12.4 的 OB-54 与 OB-55 两者与主控编排器调用登记块均为步骤 0 产物；按用户裁定**与 ⑩ 一并处理**（随本片一次性提交，不单独提交）。
- **CI 观测（2026-10-03）**：提交 `68a1ef0` 已推送 `origin/main`；CI run `37070276731`（head `68a1ef0`）——Windows 腿全绿（含 `Gates` 步实测执行：`SELFTEST PASS 248/248`、`GATE REPORT scope=ci changed=16`、`TOTAL_FAIL=0`）；**Ubuntu 腿红**：`Test` 步 `--- FAIL: TestAgentRunnerAssemblyFullChain`（`runtime_agent_e2e_test.go:464`：`invalid AgentRunner terminal intent: observed settlement requires a non-negative charged amount`）。根因为**真实缺陷**：本片新增 (A′) 映射未履行「金额由调用方供给」的既有义务（`tools/b4-e2e/runmap.go` 的 divergence (b)）；修复方式 = 声明式离线工作负载补零成本金额 + 显式边界注释 + 断言（属实现缺陷而非机械错误；按用户裁定本片修复、**不重跑 ⑦**）。
- **三个经验（须如实记录）**：① **审查盲区的平台根源**——⑦ 三轮（盲审 / 终审 / 复审）都没发现它，因为 Windows 上该路径不可达；⑦ 只能审查「它看到的代码」，无法审查「另一个平台上才能执行的路径」。② **CI 双平台补上了审查看不到的平台**——Ubuntu 腿首次执行正向全链时暴露缺陷，正好验证本报告第 10 节登记的「平台分工」缺口是**真实存在**而非理论问题。③ **「Windows 全绿」不等于全绿**——本次 Windows 腿也为绿，但 Windows 上正向全链不可达 ⇒ Windows 腿的绿只是「没有失败」，**不是**「正向全链通过」。
- **DR-9 记录**：测试 `TestAgentRunnerAssemblyFullChain`；首次红 run `37070276731`（head `68a1ef0`）Ubuntu 腿 `Test` 步；**性质 = 真实缺陷（非 flake）**；**处置 = 本片修复**；边界依据 = 准则 §6.4 边界①（新测试名 / 新失败类型 ⇒ 必须单独登记）；已登记到 `docs/t027/REMAINING_SLICES{,_EN}.md`（不另立 OB）。
- **未执行事项**：B 段的真实 enforcement 落地；正向全链的 Windows 证据；OB-57 的根因修复；OB-58 至 OB-61 的登记与处置（下一片起跑时统一落台账）；`tmp/assembly-min/` 的过程文件按既有指令保留（工作区专用）。
- **待授权项**：`git commit`（门禁读数为 16 条改动路径，其中 8 条为意图登记的新产物）与 `git push`（仅 origin）；**gitee 不推**。

## 14. 产物清单

- 下列行以制表符分隔：仓内相对路径与处置词（`repo` = 属于仓库且可解析；`local-only` = 工作区专用、不入库）。

## 15. 补充记录（闭环后证据回填，2026-10-03，由切片 `DIRECTIVE-MODEL-LISTS` 步骤 0 执行）

- **闭环提交**：本片的 DR-9 修复提交为 `d2d6370`（`fix: settle the declared offline workload with the zero-cost amount`），与首次提交 `68a1ef0` 一并位于 `origin/main`（`git log --oneline` 实读）。
- **CI 双绿（终态）**：CI run `37073985144`（head `d2d6370`）**两腿均通过**——Ubuntu 腿与 Windows 腿的 `Gates` 步各自实测 `SELFTEST: PASS 248/248`、`GATE REPORT scope=ci changed=8`、`TOTAL_FAIL=0`。取证命令 = `gh run view 37073985144 --log`，原始日志留存于 `tmp/dml/ci-run-37073985144.txt`（`local-only`）。
- **与本文档 CI 观测行的关系**：本文档先前记录的 run `37070276731` 是 DR-9 的**首次红**（Ubuntu 腿）；其后的修复与终态绿由本补充记录承接，本片报告**不以该红 run 作为终态**。
- **停点措辞口径**：原「⑩ 停点结论」描述的是**⑩ 停点当时**的状态（提交被明确保留给用户同轮授权）；提交与推送实际发生在该授权之后，二者为时序关系，不构成矛盾。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/validation/t027-assembly-min.md	repo
docs/validation/t027-assembly-min_EN.md	repo
docs/validation/evidence/T027-ASSEMBLY-MIN-freedom-list.md	repo
internal/adapters/agent_runner_postflight.go	repo
internal/adapters/agent_runner_postflight_test.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_agent_e2e_test.go	repo
internal/console/runtime_noop.go	repo
internal/console/cli.go	repo
internal/console/runtime.go	repo
tmp/assembly-min/**	local-only
tmp/dml/ci-run-37073985144.txt	local-only
```
> **归档注记（v1.23）**：本报告引用的 `tmp/` 条目已按 §1.5.1 归档至 `docs/validation/evidence/directive-history/`（单元：`assembly-min/`；各含 `MANIFEST.md` 与 `SHA256SUMS.txt`）；原路径不改写，对应关系以 MANIFEST 为准。（此前切片已清除的 `tmp/` 路径不在本次归档范围，属 OB-77 所指既有脱钩。）

- 注：树作用域下 `repo` 行要求路径**已登记进暂存区**；本片在提交授权之前采用「仅登记路径、不暂存内容」的意图登记（`local-only` 行不受此限），内容不入库，提交仍须同轮授权。
