# input-list — P55 `EXIT-SEED-BUILD-EVIDENCE` 输入包清单与前置核对记录

口径：OB-93（审查输入包须逐件断言申报输入存在；缺件即停）+ 起跑消息【五、5】前置核对。字节数为 2026-11-05 步骤 0 实测值；包内载荷的哈希以 `SHA256SUMS.txt` 为准。

## 一、规范与验收依据

| 路径 | 用途 |
|---|---|
| `docs/TEST_STRATEGY.md` / `docs/TEST_STRATEGY_EN.md` | AT-13 验收定义（`:43`）；本片 §八 需回写的文档之一 |
| `docs/SECURITY.md` / `docs/SECURITY_EN.md` | SEC-08 供应链威胁与 AT-13/AT-14 的对应关系 |
| `docs/PRODUCT_REQUIREMENTS.md` / `docs/PRODUCT_REQUIREMENTS_EN.md` | REQ-015 / REQ-019 / REQ-024 的需求条目 |
| `docs/CONTRACTS.md` / `docs/CONTRACTS_EN.md` | 发布与核验相关契约条款（本片只读引用，不改） |
| `docs/DELIVERY_DIRECTIVE.md` / `docs/DELIVERY_DIRECTIVE_EN.md` | 执行纪律（v1.60）；含测试策略与证据口径（前置核对对象之一） |
| `docs/t027/CHECKLIST_v44.md` | 队列与纪律来源（§四 层 3 第 7 项、§七 起跑通则、§八 跨片协作纪律、§九 门禁坑） |

## 二、现状实现（只读对象；本片禁止改动产品代码）

| 路径 | 用途 |
|---|---|
| `internal/release/verify.go` | 两代自托管核验逻辑（seed / candidate / oracle / SHA256SUMS / SBOM / 许可证） |
| `internal/release/verify_test.go` | 既有核验测试（核验语义的现状证据） |
| `internal/release/probe.go`、`internal/release/probe_test.go` | 候选探测与失败面 |
| `internal/release/generate.go` | 清单与 SBOM 生成 |
| `internal/release/workflow_test.go`、`internal/release/workflow_mutation_test.go` | CI 工作流冻结契约（改动 CI 结构会被该契约拦下） |
| `cmd/prfrail-release/main.go` | 发布工具 CLI 入口 |
| `.github/workflows/ci.yml` | dispatch 专用作业（candidate-build / candidate-probe / bootstrap-release）与双平台测试作业；**本片只读** |
| `testdata/selfhost/` | 固定 bootstrap commit、四项外部 oracle、noop 与 executable 链、许可证策略的夹具来源 |

## 三、既有证据与先例

| 路径 | 用途 |
|---|---|
| `docs/validation/s1-selfhost.md` / `docs/validation/s1-selfhost_EN.md` | T017 自托管两代演练的既有证据（前置核对对象之一；本片须区分「已证」与「未证」） |
| `docs/validation/evidence/CONTRACTS-S7-DURABILITY/MANIFEST.md` | 证据包格式先例（MANIFEST + SHA256SUMS 三方相等） |
| `docs/validation/contracts-s7-durability.md` / `docs/validation/contracts-s7-durability_EN.md` | 报告结构先例（含 `artifacts` 围栏与 `sixstage` 转录块） |
| `docs/validation/at24-wording.md` 等近期报告 | G7-f 围栏形态与 ⑩ 回执结构的可对照样本 |

## 四、§八 需在 ⑧ 阶段同步的六项权威文档

| 序 | 路径（`{,_EN}` 成对） | 同步内容 |
|---|---|---|
| 1 | `docs/TEST_STRATEGY{,_EN}.md` | AT-13 的状态与证据链接 |
| 2 | `docs/OPERATIONS{,_EN}.md` | 「自托管、升级和发布」章节的实测证据 |
| 3 | `docs/DEV_PLAN{,_EN}.md` | T017/T018 的构建与自托管进度 |
| 4 | `docs/S1_SUPPORT_MATRIX{,_EN}.md` | 跨平台构建最终结论（有结论时） |
| 5 | `docs/S1_RELEASE_NOTES{,_EN}.md` | 发行要求是否达成的状态（达成时） |
| 6 | `docs/t027/REMAINING_SLICES{,_EN}.md` | 台账回写 |

## 五、缺件断言（已实跑，2026-11-05）

对上表一至四各组逐件 `Test-Path -LiteralPath` 实测，结果：**18/18 存在，0 缺件**。实测字节数（取样，非全量，用于证明非空文件）：

| 路径 | 状态 | 字节 |
|---|---|---|
| `docs/TEST_STRATEGY.md` | OK | 11737 |
| `docs/DELIVERY_DIRECTIVE.md` | OK | 449639 |
| `docs/validation/s1-selfhost.md` | OK | 14270 |
| `internal/release/verify.go` | OK | 14392 |
| `internal/release/workflow_test.go` | OK | 20382 |
| `.github/workflows/ci.yml` | OK | 16432 |
| `docs/OPERATIONS.md` | OK | 16577 |
| `docs/DEV_PLAN.md` | OK | 81705 |
| `docs/S1_SUPPORT_MATRIX.md` | OK | 2557 |
| `docs/S1_RELEASE_NOTES.md` | OK | 2120 |

补充探测（不作为本片输入，仅留痕）：`internal/release/generate_test.go` **不存在**（`generate.go` 的覆盖由 `verify_test.go` 与包内其它测试承担）；本片不以此为由扩范围。

## 六、前置核对记录（起跑消息【五、5】）

### 6.1 AT-13 的验收口径（逐字引自 `docs/TEST_STRATEGY.md:43`）

| 项 | 原文 |
|---|---|
| 编号 | `AT-13` |
| 需求映射 | `REQ-015/019/024` |
| 验收内容 | `seed 与 candidate 分离；外部 oracle 比较；clean-room 成功及失败不损 seed` |
| 归属 | `harness/release；T017` |

### 6.2 `docs/validation/s1-selfhost.md` 的历史读数（要点摘录，不改写原文）

| 序 | 既有陈述（出处节） | 对本片的意义 |
|---|---|---|
| 1 | 报告结论为 `COMPLETE`（更新于 2026-09-10），并称「AT-13/14 证据均已归档」（节：报告头部） | **已有结论，但形态是叙述 + CI run 引用**；本片要补的是**可重跑、带原始 stdout/stderr 与退出码的 seed 驱动构建证据** |
| 2 | 「seed 与 candidate 使用不同源码、构建和运行目录；真实路径重叠、oracle 所有权反转、seed 变更均 fail-close」（节：AT-13/14 证据 第 1 条） | 「分离」已有实现层断言；缺的是**本片口径下的实测捕获** |
| 3 | 「外部 oracle 固定 `version`、`validate-noop`、`run-noop`、`reject-executable` 四项结果；candidate 只生成 actual，不能生成 expected」（节：AT-13/14 证据 第 2 条） | oracle 四项与「actual/expected 分离」已有陈述与安全边界说明 |
| 4 | 「本机两代演练从固定 seed commit 构建 seed，从当前 candidate 构建候选，四项 oracle 一致；seed、candidate 和 oracle 哈希均由独立发布工具复算」（节：AT-13/14 证据 第 4 条） | **本机两代演练的原始产物未入证据包**（该节为叙述），是补证重点 |
| 5 | 「固定 seed 只读使用；失败时丢弃 candidate，不能覆盖或回滚 seed」（节：边界与后续 第 3 条） | 「失败不损 seed」已有口径；缺**失败路径的实测捕获** |
| 6 | 手工 dispatch 的完整矩阵（run `34277671704`、`34325758548`、`34358467083`）8/8 jobs 成功，6 个 artifact 摘要入表 | **CI 侧证据较完整**；本片不改 CI、不再触发 dispatch（起跑消息【六、4】禁止改 workflow） |

### 6.3 缺口编号表（G-1 至 G-8）

| 编号 | 缺口 | 本片闭合位置 |
|---|---|---|
| G-1 | seed 与 candidate 分离（实测路径与不交判定 + 三条反例） | `02a-06/07/08/09`、`02b-14/15/16/17` |
| G-2 | 外部 oracle 比较（expected vs actual 成对；候选不能产 expected） | `02a-05`、`02a-12`、`02b-13`、`02b-20` |
| G-3 | clean-room 成功（完整捕获） | `02a-02/03/06`（Windows）、`02b-08/10/11/14`（Linux，含真两代） |
| G-4 | 失败不损 seed | `02a-10/11`、`02b-18/19` |
| G-5 | 跨平台同形（VM 侧 go1.22.12 复现） | `02b-00` 与全部 `02b-*` |
| G-6 | 环境可绑定（工具链读数） | `02a-00`、`02b-00` |
| G-7 | 证据可追溯（HEAD / 五枚夹具 / 准则版本绑定） | `02a-bindings`、`02b-21` |
| G-8 | 拒绝面旁证（AT-14 面，**不扩范围**，仅引用） | `s1-selfhost.md` 既有陈述，本片未实测 |

### 6.4 初步缺口方向的处置结果

原 §6.3 的四条方向**保留为历史记录**，现已全部由 ②a / ②b 闭合；口径与边界见 `claim-scope.md`（S-1..S-7 / N-1..N-7）。以下四条为起跑时的原始表述：

1. **本机两代演练的原始载体缺失**：仅有叙述与哈希结论，无 `stdout/stderr`、退出码、`go version` 等环境读数的落盘件。
2. **失败路径缺实测**：「失败不损 seed」有口径声明，缺一次可复现的失败注入与 seed 前后哈希对照捕获。
3. **Linux 侧缺少与 Windows 同形的本地捕获**：现有 Linux 证据主要来自 GitHub 托管 runner 的 dispatch run，非本片可控的原始日志。
4. **oracle 比对的「actual vs expected」成对原始输出缺落盘**：四项 oracle 的比较结果未见逐项原始文件。

## 七、边界与禁止项（引用起跑消息，不重复定义）

- 不改 `.github/workflows/*.yml`；不改 Go 核心代码；本片产出仅限 `.md` 与证据包（起跑消息【六、4】）。
- 临时脚本与中间产物一律 `tmp/EXIT-SEED-BUILD-EVIDENCE/`，只增不删，⑧c 统一清理（起跑消息【四】）。
