# P78 验证报告：v0.1.0 发行绑定与安装生命周期复验

日期：2026-11-17。状态：完成。切片标识：P78 `RELEASE-ZIP-BINDING`（`[SLICE]`，M-L 实现类）。

## 1. 切片标识与授权

- 产品所有者批准（同轮）：v0.1.0 Release 创建属 ADR-046 Q10=(C) 第一步的发行部分；该批准不得被解读为 S1 Exit 通过。
- 外部动作授权面：`gh workflow run ci.yml`（workflow_dispatch）与 `gh release create v0.1.0` 及附件上传；发布对象仅 `github.com/larsonzh/prfrail`（origin），无 gitee 动作。
- 版本号来源：ADR-010 冻结首版候选 `v0.1.0`；安装模型由 ADR-009 冻结；发布信任由 ADR-006 冻结（不做代码签名）。
- 探针额度 0；⑦ 额度 1 终审 + 2 复审（§7.5，M-L 视同 M 档）；v61 队列项的「1+3」系引用误差，本片按 §7.5 执行并在台账登记。

## 2. 步骤 0 锚点实测

| 锚点 | 期望 | 实测 |
|---|---|---|
| `HEAD` / `origin/main` | `c935713` | `c935713` / `c935713` |
| 准则版本 | v1.69 | v1.69 |
| `SELFTEST` | 737/737 | 737/737 |
| 三作用域 | `TOTAL_FAIL=0` | index / tree / ci 均为 0 |
| `tmp` 顶层 | `gates-cache` + `.gitkeep`；排除缓存后 1 个文件 | 同左 |
| `git status --porcelain -uall` | 0 行 | 0 行 |
| 滚动窗口 v59 / v60 / v61 | 180 / 242 / 214 | 180 / 242 / 214（CN 与 EN 对称） |

权限预检：`gh auth status` 登录账户 `larsonzh`，token scopes 含 `repo`；`gh api repos/larsonzh/prfrail --jq .permissions` 返回 `admin` / `maintain` / `push` 全 true。

## 3. ① 方案书结论

方案书全文见证据包 `01-arch/design-r1.md`（`DESIGN: DONE`）。要点：
- Q1 构建链路：`workflow_dispatch` 实为 8 个 input（起跑消息写 9 个系引用误差）；因 ADR-006 要求 GitHub CI 作为发布证据、且 `INSTALLATION` §2 已冻结核验断言 `0.1.0+<candidate_commit>`，选定 CI `candidate-build` 与 `selfhost` 路径；本机构建仅作等价对照，不作发布证据。
- Q2 绑定：四元组 = 版本 / candidate commit / CI run / 文件名；三处落点 = `s1-exit` 非门表区、`S1_SUPPORT_MATRIX` 加法式新节、`S1_RELEASE_NOTES`；发布证据由 CI 的 `prfrail-release generate` 与 `verify` 产出并作为 Release 资产。
- Q3 Release：tag `v0.1.0` 指向 candidate commit（发行源码与二进制同源）；非 draft、非 prerelease；EOL 公告义务为无（首版 minor，无前一 minor）。
- Q4 复验：安装 / 升级 / 回滚 / 卸载保留数据四项，全部对最终 ZIP；升级按 `INSTALLATION` §5 的切换机制口径执行（同版本第二独立目录 + 显式路径切换 + 旧目录留作回滚点）。
- Q5 文档归属：第一层六项必改；第二层 `ADR_REGISTER:19`、`INSTALLATION_PLAN` 三处、`OPERATIONS` 一处、`DEV_PLAN` 一处加法式注记由实读升级为「触」；`PRODUCT_REQUIREMENTS` 与 `ADR_REGISTER:62` 不触。
- Q6 组合检索：`RELEASE-ZIP` / `GitHub Release` / `v0.1.0` 三个独立 grep 的命中全部归入「历史不改写」「现役假句已升级」「第一层台账」三类，无新增升级项。

## 4. ② 执行结果

- 构建：由 CI workflow_dispatch 触发（run `44043705691`）；candidate-build / candidate-probe / selfhost 在 ubuntu 与 windows 双腿全部 success；构建命令与版本注入逐字登记在 Release notes。
- 产物：artifact `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86` 解压后为 `prfrail/` 目录，恰含 `prfrail.exe`、`SHA256SUMS`、`sbom.spdx.json`、`licenses.json`；最终 ZIP 名为 `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`，解压根目录只有 `prfrail/`。
- Release：GitHub Release `v0.1.0` 已创建（tag 指向 candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`），资产 = ZIP + `SHA256SUMS` + `sbom.spdx.json` + `licenses.json`。
- 复验：四项全部对最终 ZIP 执行并通过；逐项断言见 §6。
- 文档同步：第一层六份按目标文本落盘；第二层五处按实读判定落盘；`ADR_REGISTER:62` 与 `PRODUCT_REQUIREMENTS` 未动。

口径偏离披露：相对方案书 r1 的目标文本，落盘在五处作了更保守的增补——`S1_RELEASE_NOTES:5` 末句改为「不构成 S1 exit 整体通过」（贴合 ADR-046 硬边界）、`OPERATIONS:5` 增补 P78 标识与四元组、`INSTALLATION_PLAN:32` 增补资产文件名、`INSTALLATION_PLAN:7` 与 `S1_RELEASE_NOTES:18` 的「无正式发行包」改为「候选发行」。五处均保持 CN/EN 对称，语义只收窄不扩张。

## 5. 发布绑定四元组与三处落点

| 元 | 取值 |
|---|---|
| 版本 | `v0.1.0` |
| commit | `c93571352f456a275d6e71ce3e084c9f96360e86` |
| run | `44043705691` |
| 文件名 | `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip` |

三处落点（`INSTALLATION:208` 明定，方案书引作 `:206`，§2 插入两行后整体后移）均已逐字写入同一四元组：`s1-exit` 剩余风险第 1 条（非门表区）、`S1_SUPPORT_MATRIX` 文末「已绑定发行记录」节、`S1_RELEASE_NOTES` 状态行与文末切片同步引语。
另有两处执行证据同源承载：本报告与 GitHub Release 正文；四者的 commit / run / 文件名完全一致。

## 6. 复验矩阵结果

| 项 | 关键断言 | 结果 |
|---|---|---|
| 安装 | 解压根目录只有 `prfrail/`；四个发行文件；`SHA256SUMS` 逐行通过；`version --json` 等于 `0.1.0+candidate`；`init` / `validate` / `preview` / noop `run` / `report` 全通过 | 通过 |
| 升级 | 第二独立目录解压同一最终 ZIP；`version` 同串；`validate` / `preview` 通过；旧目录未被覆盖；历史 run 报告逐字一致 | 通过 |
| 回滚 | 路径切回 dirA；`version` 同串；历史 run 报告与升级前逐字一致 | 通过 |
| 卸载保留数据 | 无运行进程；仅删 dirB 二进制目录；run 目录与 workspace 仍在；保留数据仍可 `report` 读取 | 通过 |

四项的原始断言输出落 `tmp/P78-RELEASE-ZIP-BINDING/20-reverify/`（一次性，⑧c 清理）；失败处置为停下升级，本片未触发。

## 7. S-G3 收口

S-G3 的判据是「正式安装包与安装教程」与「安装、升级、回滚、卸载实测」两项：最终 ZIP 已绑定发行记录并在 Windows 本机完成四项复验，两项判据的事实前提均已满足，S-G3 由未通过转为通过。
门表零改动：`p57-closeout-assessment{,_EN}` §4.2 与 `s1-exit{,_EN}` 的门表行、行名与判定一字未动；新事实由本报告与三处文档落点承载。
S-G5 记「发行部分已批准」；硬边界不变：本批准不构成 S1 Exit 整体通过，S1 Exit 整体批准待 S2 之后。

## 8. 审查链

- ③ 一致性：S-G3 收口结论与 p57 §4.2、`s1-exit`、v61 §四、ADR-046 不矛盾；四元组在三处落点与两处证据载体逐字一致；CN / EN 改动行数逐对相等。
- ④ 门禁：三作用域 + `SELFTEST` + `go build` / `go vet` / `go test` 全绿，读数见 §9。
- ③ 一致性：S-G3 收口与 p57 §4.2、`s1-exit`、v61 §四、ADR-046 无冲突；四元组在三处落点与两处证据载体逐字一致；CN/EN 改动行数逐对相等。
- ⑤ 预审：2 轮——r1 `PRE-REVIEW: PASS WITH FIXES`（2 中 3 低）；整改后 r2 `PRE-REVIEW: PASS WITH FIXES`（0 中 2 低）；额度 2/2 用尽，r2 的 2 条低风险已由主控顺手闭合。
- ⑥ 独立扫描：1 次——`INDEPENDENT SCAN: PASS`（0 发现，A–E 契约合规，由 ⑦ 独立确认 `COMPLIANT: YES`）；额度 1/3。
- ⑦ 终审：r1 `FINAL REVIEW: FINDINGS`（1 高 1 中，均指向证据包闭环）；整改后 r2 `RE-REVIEW: PASS`（0 新发现，两项闭合）；额度 1 终审 + 1 复审。

## 9. 门禁读数

以下读数为 ⑨ 暂存并注入 ⑥ 围栏后的实跑值；④ 的中间态在 tree 作用域曾判红两处（G7-b 报告工件未暂存、G7-f 围栏未注入），两处均属时序红因，⑨ 前已闭合。

| 判据 | 读数 |
|---|---|
| `--scope=index` | `TOTAL_FAIL=0` |
| `--scope=tree` | `TOTAL_FAIL=0` |
| `--scope=ci` | `TOTAL_FAIL=0` |
| `--selftest` | `SELFTEST: PASS 737/737` |
| `go build ./...` | exit 0 |
| `go vet ./...` | exit 0 |
| `go test ./...` | 全部 ok，`FAIL=0` |

## 10. 已知限制与残留

- `ADR_REGISTER:62`「当前仍不是已发布 SLA」保持不变：候选发布不等于 SLA 生效；该句待 S1 Exit 通过或用户裁定再翻。
- `S1_SUPPORT_MATRIX:5` 与 `S1_RELEASE_NOTES` 的「非已发布 SLA / S1 exit 待 S2」限定保持，不因本次发布翻转为「已发布」。
- `ADR_REGISTER:104` 的 T001 决策单安装模型行属 2026-09-09 时点授权记录，按 §11.1 不改写。
- `DEV_PLAN` 的 T018 状态段为带日期的历史记录，本片只追加 P78 注记，不改写历史句与勾选状态。
- 真实跨版本升级演练留待首个 patch / minor 发布；本片按 `INSTALLATION` §5 的切换机制口径验证升级路径。
- `s1-exit` 的 2026-09-09 结论段与门表三行证据列仍写「最终 ZIP 尚未绑定和复验」：前者属带日期的历史结论（§11.1 不改写），后者属门表（P76 明令门表零改动）；现役事实以本报告 §7 与 `s1-exit` 剩余风险第 1 条为准，后续同步片统一改口径。
- 报告 artifacts 块声明的 `04-reviews/` 与 `06-gates/` 载荷在审查与门禁轮次产出后归档；④ 时点的缺件属时序，不是遗漏。

## 11. 计量

| 环节 | 次数 |
|---|---|
| ① 架构 | 1 次（`deep-reasoner` + `DeepSeek V4 Pro`） |
| 批量辅助 | 2 次（`GLM-5.3-Flash`） |
| 探针 | 0 |
| R / D | R = 0、D = 1（D-1 = 空提交作 verifier 标识的实施顺序调整，自主披露） |

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/INSTALLATION.md	repo
docs/INSTALLATION_EN.md	repo
docs/INSTALLATION_PLAN.md	repo
docs/INSTALLATION_PLAN_EN.md	repo
docs/OPERATIONS.md	repo
docs/OPERATIONS_EN.md	repo
docs/S1_RELEASE_NOTES.md	repo
docs/S1_RELEASE_NOTES_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
docs/validation/release-zip-binding.md	repo
docs/validation/release-zip-binding_EN.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/00-briefing.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/01-arch/design-r1.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/02-batch/batch-01-release-pipeline.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/02-batch/batch-02-release-docs.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/04-reviews/prereview-r2.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/04-reviews/scan-r1.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/04-reviews/final-r1.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/04-reviews/final-r2.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/06-gates/gate-index.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/06-gates/gate-tree.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/06-gates/gate-ci.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/06-gates/base-go.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/MANIFEST.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/SHA256SUMS.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/freedom-list.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/input-list.md	repo
```

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
| # | SEVERITY | file:line | 被违反的句子或判据 | 最小修复 |
| - | none | - | - | - |

## Section C
- 一致性：四元组在 `docs/validation/release-zip-binding.md`（L50-L53）、`docs/validation/release-zip-binding_EN.md`（L50-L53）、`docs/validation/s1-exit.md`（L66）、`docs/S1_SUPPORT_MATRIX.md`、`docs/S1_RELEASE_NOTES.md`（L5）与 `tmp/P78-RELEASE-ZIP-BINDING/10-build/release-notes.md` 中保持一致：版本为 v0.1.0，candidate commit 为 c93571352f456a275d6e71ce3e084c9f96360e86，run 为 44043705691，文件名为 prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip。没有出现四元组漂移或同值不一致。
- 越界：变更集 `tmp/P78-RELEASE-ZIP-BINDING/impl.diff` 与 `impl.numstat` 只覆盖 28 个文档型变更；没有触及产品代码、CI 工作流、p57 门表、s1-exit 门表、准则、CONTRACTS、schema、fixtures、其它 ADR、S1_SUPPORT_MATRIX 平台支持层级、代码签名、Linux 正式发布。对于「S1_SUPPORT_MATRIX 平台清单与支持层级」这条，实际变更是新增「已绑定发行记录」说明，不是改动平台名单或支持层级。
- 事实与证据：报告中的关键事实在 `tmp/P78-RELEASE-ZIP-BINDING/20-reverify.log` 中得到直接支撑：安装、升级、回滚、卸载保留数据四项均为 PASS。门级原始输出 `tmp/P78-RELEASE-ZIP-BINDING/04-gates-tree-pre.txt` 里 G7-b/G7-f 仅是时序红因（未暂存工件 / 围栏注入），并非实质缺陷。没有出现「未验证却写成已验证」的模式。
- 门级表述：有明确的硬边界，且与 ADR-046 的两步审批一致。例如 `docs/validation/release-zip-binding.md` 与 `docs/S1_RELEASE_NOTES.md`（L5）写明「发行部分已批准，但不构成 S1 Exit 整体通过」，以及 `docs/ADR_REGISTER.md`（L62）与 `docs/S1_SUPPORT_MATRIX.md`（L5）中保留「当前仍不是已发布 SLA」。因此没有将「候选发行」写成「S1 Exit 整体通过」或「已发布 SLA」。

## Section D
- 该结论在什么观察下会被推翻：
  1. 若同一四元组在 `docs/validation/release-zip-binding.md`、`_EN.md`、`docs/S1_RELEASE_NOTES.md`、`docs/S1_SUPPORT_MATRIX.md` 与 `tmp/P78-RELEASE-ZIP-BINDING/10-build/release-notes.md` 中出现不一致。
  2. 若 `tmp/P78-RELEASE-ZIP-BINDING/impl.diff` 发现了产品代码、CI、门表、CONTRACTS、schema、fixtures 或 Linux 正式发布相关越界。
  3. 若 `tmp/P78-RELEASE-ZIP-BINDING/20-reverify.log` 不再显示 install / upgrade / rollback / uninstall 全部 PASS。
  4. 若报告或 Release notes 改成「已发布 SLA」或「S1 Exit 整体通过」。
  5. 若实际远端的 GitHub Release 资产与本地证据发生不一致，或相关文件在后续修订中篡改为与当前证据不符。

## Section E
- 我无法验证的部分：没有直接抓取远端 GitHub Release 页面与发布资产；本扫描基于仓库内证据、日志和文档构成的本地一致性证明。就「仓库内是否符合四元组 / 边界 / 事实」，本地证据足够；就「远端公开页面始终保持一致」，此时点无法做外部实时确认。
- 什么会推翻 PASS：只要本地证据变成与真实远端发布页面或资产不一致；或后续出现新增变更扩大到禁止触碰范围；或四元组 / 门级表述变成「已发布 SLA / S1 Exit 整体通过」，PASS 即应撤销。

INDEPENDENT SCAN: PASS
```
