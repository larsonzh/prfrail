# P84 bounddir 残余修复验证报告（BOUNDDIR-RESIDUAL）

日期：2026-11-21。状态：`PARTIALLY CLOSED`（部分闭合）。切片标识：P84 BOUNDDIR-RESIDUAL

**结论：`PARTIALLY CLOSED`**（部分闭合）——T5#4（F4 残余）已全闭合；T4（F1 残余）大部分闭合，残余 = N13-UNC（UNC 共享根的两种规范拼写等值被误拒；确定性界内误拒、方向 fail-closed、属可用性缺陷）；另有证据归因 Low 一项（W10 变异标注与 r3 原文矛盾）。**不得读作「闭环」**。

---

## 1. 切片标识与授权

| 项 | 内容 |
|---|---|
| 标识 | `BOUNDDIR-RESIDUAL`（承接 P83 报告 §8 的 P84-d 与 P84-e；闭合 F1 / F4） |
| 授权 | 用户 2026-11-21 点名起跑；同日排片裁定三项；收窄复审 1 轮另计额度例外 |
| 依据 | P83 报告 §8 顺延清单；P83 证据包的 T4 / T5#4 原文；v67 快照 §四 队列项 |
| 边界 | 只修 T4 / T5#4；不碰 F3 / F6 / F8（P85 范围）；不改 S1 Exit 八门判据；不改威胁模型基线；不升准则版本；不做 gitee |

## 2. 步骤 0 锚点实测

| # | 锚点 | 读数 | 结论 |
|---|---|---|---|
| 1 | HEAD | `71d03c3` | 相符 |
| 2 | origin/main | `71d03c3` | 相符 |
| 3 | 准则版本 | v1.71 | 相符 |
| 4 | SELFTEST | 740/740 | 相符 |
| 5 | 三作用域 | index / tree / ci | 全为 `TOTAL_FAIL=0` |
| 6 | 滚动窗口 v65（LF 计数） | 220 | 相符 |
| 7 | 滚动窗口 v66（LF 计数） | 216 | 相符 |
| 8 | 滚动窗口 v67（LF 计数） | 216 | 相符 |
| 9 | 机械读数 | ADR 47 / OB 139 / DR 15 / 试点顺延 25 | 相符 |

## 3. 目标与范围

按 P83 报告 §8 逐字承接两项残余，不扩不缩：T4 使产物枚举与后代打开在 Windows 上仍可证明归属；T5#4 使租约代次到任务 proof 的更新成为可恢复事务。每项一条可证伪闭合断言、至少一条反例回归测试与变异检验；① 先判定契约影响再动代码。

## 4. 修复实现与闭合状态

| 项 | 面 | 主要改动位置 | 闭合状态 |
|---|---|---|---|
| T4-a | F1 | `internal/evidence/bounddir_windows.go` 的 `openDescendantRelative`（以绑定句柄为根逐组件走链） | **已闭合**：W2 界外拒绝与界内放行两侧反例 |
| T4-b | F1 | 同文件的 `isUnderDirectory` 与 `boundDirVerify`（Ordinal 取代 EqualFold） | **已闭合**：W1 仅大小写不同的同名目录被拒，端到端同样拒绝 |
| T4-c | F1 | 同文件的 `reparseStoredTarget` 与 `boundRelativeSuffix`（存储目标为唯一权威，界内由构造保证） | **已闭合**：W5b / C12 / W8 / W9 四条反例 |
| T4-d | F1 | `internal/evidence/bounddir_clean.go` 的 `boundRootName` 与两侧调用点 | **已闭合**：W3 两种根拼写等值且同错误拒绝 |
| T4-e | F1 | `internal/gates/runner.go` 的现状短路（本项零代码改动） | **已闭合（钉住）**：既有 no-glob 反例保持绿 |
| T4-f | F1 | 同文件的 `readDirNames` 与根枚举直读绑定句柄 | **已闭合**：W6 根枚举不经名字解析 |
| T4-g | F1 | 同文件的 `openRelativeComponent` 与组件长度守卫 | **已闭合**：W7 超长组件在打开前拒绝 |
| T4-h | F1 | 同文件的 `reparseStoredTarget` 内嵌 NUL 守卫 | **已闭合**：W11 双解析反例（收窄复审判 N14 为 CLOSED） |
| T4-i | F1 | `boundRelativeSuffix` 的等值形态（盘符根保留尾分隔符） | **残余 N13-UNC**：盘符根等值已修，UNC 共享根两种规范拼写仍被误拒 |
| T5#4-a | F4 | `internal/console/runtime_supervised.go` 的 `persistRenewedPublishLease` 与 `runtime_agent.go` 的 `verifyPublishLeaseHistory` | **已闭合**：C1 异主不改写并报错、C2 多代历史追平 |
| T5#4-b | F4 | `internal/guard/lease.go` 的 `RunLocked` 与 `runtime_agent.go` 的三路恢复分支 | **已闭合**：C3 跨进程重入恢复、C6 并发推进整批拒绝、C7 链深上限 fail-closed |
| T5#4-c | F4 | `runtime_agent.go` 的 `verifyPublishLeasePredecessor`（现状钉住） | **已闭合**：C4 / C10 合形但不匹配仍拒绝 |
| T5#4-d | F4 | `runtime_agent.go` 的 `publishLeaseBelongsToRun` 与 `assertPublishLeaseCurrent` | **已闭合**：C15 异 consumer 拒绝、C16 锁内写、C18 入口持锁 |
| 锁基座 | F4 | `internal/guard/lease_lock_windows.go`、`lease_lock_unix.go`、`lease_lock_unsupported.go` | **已闭合**：C9 进程死亡由内核回收、C11 活锁互斥、C17 旧锁目录独立报错 |
| 收窄复审整改 | F1 | `boundRelativeSuffix` 等值形态与 `reparseStoredTarget` 内嵌 NUL 守卫（⑦ 收窄复审后的整改） | N13 = **部分**（盘符根等值已修；UNC 共享根的等值误拒顺延，定性 = 确定性界内误拒、fail-closed 方向、可用性缺陷）；N14 = **已闭合** |

## 5. 契约影响判定

| 项 | 判定 | 理由 |
|---|---|---|
| T4（F1 残余） | **DOES-NOT-TOUCH / IN-BOUNDARY** | 被引用句只定要求不定机制；句柄相对走链、Ordinal 比较、存储目标解析都是 `internal/evidence` 对该句的实现强化；`ntdll!NtCreateFile` 与既有 `kernel32` LazyDLL 调用同型，属 OS API 而非第三方包，`go.mod` 不动 |
| T5#4（F4 残余） | **DOES-NOT-TOUCH / IN-BOUNDARY** | 被引用句已由 store + proof 架构落实；proof 自述不入 recordHash 集，租约文件属 `internal/guard` 的 store-local 布局；只加内部链式校验与 OS 锁包装，无新 schema、无新 CLI、无新依赖 |

两项均未触契约，故未改 `docs/CONTRACTS{,_EN}.md`、未改 schema 与夹具、未改 `go.mod`。守则遵守：实现未引入新 schema 记录、新 CLI 命令、新依赖或第三方包。

## 6. 分域升级与用户裁定

| 事项 | 裁定 | 留痕 |
|---|---|---|
| 主控载体 | 切换为 `DeepSeek V4 Pro (deepseek)` High，替代 Flash | 用户 2026-11-21 裁定；本报告 §12 登记 |
| ⑦ 额度 | 1 终审 + 3 复审 | 超出 §7.5「M-L 视同 M」字面；系用户直接裁定、不得作为先例 |
| 盲审 | 必抽 1 次，另计 +1，终审前执行 | T4 触身份语义、T5#4 触所有权语义（§7.7） |
| ⑦ 档位 | `GPT-6.2 Sol (copilot)` Max | 同轮取档授权，§8.4 记录设定时刻与暂停状态 |
| 收窄复审 1 轮 | 允许，输入仅 N13 / N14 两处整改 | 额度例外，超出 §7.5 与本片额度裁定，不得作为先例 |
| R 口径 | 复审环的整改重跑不计 R；R = 3 | §9.2 口径 |
| 起跑消息「R 达 2 停下」 | 与 §9.2 冲突，以准则为准 | 按 §13.2.11 纪律 F 处置，已登记、不记为行为漂移 |

## 7. 审查链

| 环节 | 载体 | 轮次 | 结论 |
|---|---|---|---|
| ① 方案书 | `deep-reasoner` + `DeepSeek V4 Pro (deepseek)` Max | Rev 1–6 | 契约判定 DOES-NOT-TOUCH / IN-BOUNDARY |
| ② 实现 | `prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)` | 11 轮 | 逐条落地，含 Rev 3b 锁基座重做 |
| ③ 测试 | `prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)` | 8 轮 | 反例 W1–W11 与 C1–C18；变异逐项实测变红并 SHA256 复原 |
| ⑤ 预审 | `DeepSeek V4 Pro (deepseek)` | 2 / 2（用尽） | r1 PASS WITH FIXES → r2 PASS |
| ⑥ 独立扫描 | `MAI Code 2 Flash (copilot)` | 2 / 3 | r1 FINDINGS → r2 PASS |
| ⑦ 盲审 | `GPT-6.2 Sol (copilot)` Max | 1 | FINDINGS（2 High + 3 Medium + 1 Low） |
| ⑦ 终审 | `GPT-6.2 Sol (copilot)` Max | 1 / 1 | FINDINGS |
| ⑦ 复审 | `GPT-6.2 Sol (copilot)` Max | 3 / 3（用尽） | r1 / r2 / r3 均 FINDINGS |
| ⑦ 收窄复审 | `GPT-6.2 Sol (copilot)` Max | 1（额度例外） | FINDINGS（N13-UNC 残余） |

逐字原文全部归档于 `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/`；方案书 Rev 1–6 归档于同包 `01-arch/`。⑦ 收窄复审为最后一读，其 FINDINGS（N13-UNC 残余）按 P83 报告的兜底收口（b）先例处置：残余如实登记并顺延后续切片，不改写为 PASS。

## 8. 顺延 P85 候选清单

| 编号 | 顺延内容 | 来源 |
|---|---|---|
| 无新增 | 本片不新增 P85 载荷；N13-UNC 属 F1 残余，P85 边界仅含 F3 / F6 / F8，是否并入须用户裁定 | ⑦ 收窄复审 |

N12（证据包 MANIFEST 描述更新）不属顺延项，由 ⑧b 收尾时一并落地。

## 9. 文档同步

| 对象 | 处置 |
|---|---|
| 本报告对 | 新增 `docs/validation/bounddir-residual{,_EN}.md`（CN / EN 严格同位） |
| 台账 | `docs/t027/REMAINING_SLICES{,_EN}.md` 的 P84 段状态与完成记录（⑧b 落地） |
| 安装文档 | `docs/INSTALLATION{,_EN}.md` 各加 2 行：旧 Mkdir 目录锁的升级边界与人工处置顺序 |
| `SECURITY` 注记 | `docs/SECURITY{,_EN}.md` 第 5 行注记按本报告结论如实写（⑧b 落地） |
| `CONTRACTS` / schema | 未改（两项判定 DOES-NOT-TOUCH） |
| ADR | 不新增 |
| 未改 | `CHECKLIST_v67{,_EN}`、`s1-exit`、`p57` 门表、准则本体、CI、`DEV_PLAN` |

## 10. 门禁结果

| 检查 | 结果 |
|---|---|
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 740/740` |
| `--all --skip-base --scope=index` | `TOTAL_FAIL=0` |
| `--all --skip-base --scope=tree` | `TOTAL_FAIL=0` |
| `--all --skip-base --scope=ci` | `TOTAL_FAIL=0` |
| `go test ./... -count=1` | EXIT=0，全量 15 包 ok（⑧b 于 Windows 本机实测） |
| `gofmt -l internal cmd` | 无输出 |
| IDE 诊断 | 本片改动件在编辑器内为 0 告警 |
| 跨平台交叉诊断 | 6 条（js / wasm 目标下 guard / console / evidence 各 2 条），不影响 windows / linux / darwin 构建 |
| 编码 | 改动 `.go` 为 UTF-8 无 BOM + LF；新增 `.md` 为 UTF-8 with BOM + LF |

上表为 ④ 门禁时点读数，并以 ⑧b 实测复核（`go test` 15 包 ok、`SELFTEST` 740/740）；⑧a 写盘后与 ⑧b 暂存前各复跑一次树作用域，唯一 FAIL 均为 G7-b 的「围栏工件未入索引」，属 `git add` 之前的固有状态（⑧b 复跑读数见 `tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-postNarrow.txt` 与 `tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-8b.txt`）；⑧b 暂存后由主控复跑 index / tree 两作用域归零——该口径为预期、非绝对断言，以复跑输出为准。

## 11. ⑨ 验证结果与本机不可执行项

| 项 | 结果 |
|---|---|
| Windows 本机 `go build` / `go vet` / `go test ./...` | `go build` / `go vet` ok；`go test ./... -count=1` EXIT=0，全量 15 包 ok |
| `-race` | 本机无 gcc，不执行；由 CI ubuntu 腿的 Race 步骤承担 |
| Ubuntu VM（`10.0.0.199`，go1.22.12，`bash -lc`，`GOPROXY=https://goproxy.cn,direct GOSUMDB=off`） | `go build` / `go vet` ok；`go test ./... -count=1` EXIT=0，全量 15 包 ok；定向 `-race` 4 包（evidence / guard / console / adapters）ok |
| CI ubuntu 腿 Race 步骤 | 修复提交 `32bf234` 推送后 run `44926215282` 双腿绿：ubuntu Race 步骤 6 包 `ok`；windows / ubuntu 两腿 `base go test ./... :: ok=15 FAIL=0`、`SELFTEST 740/740`、`TOTAL_FAIL=0`（详见 §11.2 回填） |

本表读数一律由 ⑧b 回填实测值，⑧a 不写预估数字；若三层任一层出现 FAIL，结论须按实测改口而不是沿用本节结论。本机不可执行项登记：Linux-tagged 用例（bounddir unix 大小写与 reparse、lease 事务恢复）由 VM 层覆盖；W1 端到端大小写敏感目录以 `fsutil` 建立，本机成功执行、未 skip。

## 11.2 CI 首跑红与修复（§11.2 回填，2026-11-22）

主提交 `d597f1e` 推送后 CI 首跑**红**；修复提交 `32bf234` 推送后**双绿**。本节按 §11.2 回填实测读数与登记。

| 项 | 内容 |
|---|---|
| 首跑红 | run `44877224368`（head `d597f1e`）windows 腿 `Test` 步骤红（ubuntu 腿绿）；`gh run rerun --failed` attempt 2 复跑 windows 腿同面红 |
| 失败清单 | attempt 1 = attempt 2 = **26/26 逐条相同**（归一化机械 diff 零差异）⇒ 确定性缺陷、非 flake；含 adapters postflight 全族、adapters 两条 bind、gates 两条、evidence W2/W5b/W8/W9 与 `TestBoundDirJunctionInsideBoundAdmitted` |
| 根因 R-2 | `openRelativeComponent` 对**文件叶**打开缺失 `FILE_SYNCHRONOUS_IO_NONALERT`（CreateOptions 0x20）⇒ 仓库 CI 的 go1.22.12 工具链上不可读（修前 87 失败 = 契约违约） |
| 根因 R-1 | `%TEMP%` 的 8.3 短名（`RUNNER~1`）与句柄长名（`runneradmin`）Ordinal 误拒：短名根前缀被误判为界外 ⇒ 4 条短名误拒锚点 |
| 修复 | 提交 `32bf234`（`fix: recover windows abi compatibility for evidence reads`，2 文件 +1014/−2，仅推 origin）；含方案 A 48 位记录号 + 判别用例 + 两处注释（窄复审 RELEASE，见本节末载体追认与 `15-final-review-A.md`） |
| CI 双绿 | run `44926215282`（head `32bf234`）：windows / ubuntu 两腿 `Test` 步骤各 15 包 `ok`；`SELFTEST 740/740`；`base go test ./... :: ok=15 FAIL=0`；`TOTAL_FAIL=0`；ubuntu Race 步骤 6 包 `ok` |
| C7 单次失败 | go1.22.12 全量首次运行中 `TestP84T5h4C7HistoryDepthIsFailClosed` 因 lease 文件重命名 `Access is denied` 失败一次；隔离复跑 3/3 PASS、全量第二次 15/15 ok ⇒ 环境性 flake、未复现，按 DR-17 登记、不记 R |
| 报告行数勘误 | 报告对磁盘实况 = **246/246**（各 246 个 LF 行、结尾带 LF）；此前回执 247 系 split 尾空元素误差；正文无任何行数表述、无需改文，此处如实登记；**回填后本报告 LF 行数 = 266**（CN 版，EN 版见其对） |
| 临时件落点 | 本片 CI 捕获件与证据/设计件统一落 `tmp/P84-CI-FIX/`；仓库外工具链 `D:\LZProjects\.toolchains\go1.22.12`（不入库 tmp）；`%TEMP%` 探测件已删、编辑器缓冲二次落盘按修订记录补正（§13 D=0 口径） |
| wasm 诊断 | 带 `[js,wasm]` 后缀的 syscall 诊断来自 gopls 额外分析视图，构建标签与调用 P84 前已存在 ⇒ 单列「非验收目标」，不计 P84 新增问题、不阻断验收 |
| DR 登记 | DR-16（CI 首跑红 + R-2/R-1 根因）与 DR-17（C7 单次 flake）登记于 `docs/t027/REMAINING_SLICES{,_EN}.md` 已知缺陷表 |

**载体追认（不新增模型）**：用户裁定署名「GPT-6.2 Sol Max」，实际可用载体为 `GPT-6.2 Sol (copilot)`（与本片既有 ⑦ 载体一致），据此如实追认；不新增模型、不新增额度（OB-37 追加登记见 `tmp/P84-CI-FIX/10-final-review-request.md`）。

## 12. 环境与档位

| 角色 | 载体与档位 | 设定时刻 | 是否暂停 |
|---|---|---|---|
| 主控 | `DeepSeek V4 Pro (deepseek)` High | 2026-11-21 起跑 | 否 |
| ① 架构 | `deep-reasoner` + `DeepSeek V4 Pro (deepseek)` Max | 2026-11-21 | 否 |
| ② 实现 | `prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)` High | 2026-11-21 | 否 |
| ③ 测试 | `prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)` High | 2026-11-21 | 否 |
| ⑤ 预审 | `DeepSeek V4 Pro (deepseek)` | 2026-11-21 | 否 |
| ⑥ 扫描 | `MAI Code 2 Flash (copilot)` | 2026-11-21 | 否 |
| ⑦ 盲审 / 终审 / 复审 / 收窄复审 | `GPT-6.2 Sol (copilot)` Max | 2026-11-21 本片 ⑦ 环节 | 否 |
| 批量辅助 | `GLM-5.3-Flash (glm)`，严格串行，限四类任务 | 2026-11-21 | 否 |

## 13. 审查与计量

| 项 | 值 | 说明 |
|---|---|---|
| R | 3 | ⑤ 环整改重跑 + ⑥ 环整改重跑 + ⑦ 终审整改重跑；复审环整改不计 R |
| D | 0 | 无临时件纪律违规；`%TEMP%` 探测件「已删后编辑器缓冲二次落盘」按修订记录补正，不计 D |
| 探针 | 0 | 无外部真实调用 |
| 批量辅助 GLM | 2 次 | 摘录核对 + 双语镜像，全 PASS |
| ⑦ 额度 | 盲审 1 + 终审 1 + 复审 3 + 收窄复审 1（额度例外） | 四要素：裁决人 = 用户、日期 = 2026-11-21、理由 = N13 属确定性界内误拒、修复方向明确且小，收窄复审一轮值得、依据 = 用户裁定消息；**不得作为先例** |
| 高风险标记 | 已启用（§9.2） | R 达 2 起按 §5.2 高风险行全环启用 |
| ⑤ / ⑥ 口径 | PASS 绑定整改前快照 | 最终快照仅由 ⑦ 收窄复审覆盖，且仍为 FINDINGS |
| 制度冲突 | 起跑消息「R 达 2 停下」与 §9.2 冲突 | 按 §13.2.11 纪律 F 以准则为准，已登记 |

## 14. 明确未执行事项

未做 gitee 动作；未升准则版本；未新增 ADR 或 OB；未改 CI 与 `DEV_PLAN`；未跑 `-race`（本机无 gcc）；未在 Linux 实机执行 Linux 专属用例（由 Ubuntu VM 与 CI ubuntu 腿承担）；未修 N13-UNC（残余，去向须用户裁定）；⑨ 三层读数中 Windows 本机与 VM 两层已回填（§11），CI 读数已按 §11.2 回填（§11.2）。

## 15. ⑧a 自查

① 每处「已闭合」都附 file:line 或测试名证据（§4 表格逐行对应 `internal/evidence/p84_bounddir_residual_test.go` 与 `internal/console/p84_t5h4_residual_test.go` 的用例名）。② 无被本片改动覆盖的旧口径残留：起跑消息所述「根以空相对路径表示而显式 `.` 被拒」已按实际实现改写为共享分类器口径。③ 读数标注时点：§2 与 §10 为 ④ 门禁时点读数，§11 原为占位符、⑧b 已回填实测读数。

## 16. 六段自查（sixstage）

```sixstage
## Section A 总体结论
PARTIALLY CLOSED（T4 残余 N13-UNC 1 项，其余全闭合）

## Section B 发现清单

| SEVERITY | 位置 | 残余 |
|---|---|---|
| Medium | N13-UNC | UNC 共享根两种规范拼写的等值误拒（确定性界内误拒、方向 fail-closed、属可用性缺陷） |
| Low | W10 证据归因 | 变异标注「轮 11 前原形」与 r3 原文矛盾，两项须分开标注 |

## Section C 逐条证据与推理
- T5#4 的四条分支与锁基座由 C1–C18 反例钉住，逐条见 §4 表格。
- T4 的走链、Ordinal、根语义、NUL 守卫由 W1–W11 钉住，逐条见 §4 表格。
- ⑦ 各轮原文归档于证据包 `04-reviews/`，方案书 Rev 1–6 归档于同包 `01-arch/`。

## Section D 可证伪性审计
- 把 `boundRelativeSuffix` 的等值比较还原为 TrimRight 形态：W10 的根等值腿立即变红。
- 删除内嵌 NUL 单元扫描：W11 的两个直接解析子测试由拒绝变为解码成功。
- 把 `persistRenewedPublishLease` 的链验证退回仅比 `ResourceID`：C1 立即变红。
- 删 `renewPublishLease` 的追平分支或 Reclaim 路径：C3 立即变红。
- 把本报告结论改写为「闭环」：§4 的残余行与本节结论互证立即红。
- UNCOVERED：UNC 共享根在生产入口的可达性未实测；变异逐项实跑与 ⑨ 三层读数由 ⑧b 回填。

## Section E 我无法验证的部分
UNC 共享根在生产入口的可达性、变异逐项实跑、⑨ 三层实测读数均不在本机静态验证范围内。

INDEPENDENT SCAN: FINDINGS
```

## 17. 工件围栏

```artifacts
docs/INSTALLATION.md	repo
docs/INSTALLATION_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/bounddir-residual.md	repo
docs/validation/bounddir-residual_EN.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/00-briefing.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/01-arch/01-arch-design.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-blind-review.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-narrow-re-review.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-re-review-r1.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-re-review-r2.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-re-review-r3.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/MANIFEST.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/SHA256SUMS.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/freedom-list.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/input-list.md	repo
internal/console/p83_fixround7_tester_test.go	repo
internal/console/p84_t5h4_residual_test.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_supervised.go	repo
internal/evidence/bounddir_clean.go	repo
internal/evidence/bounddir_unix.go	repo
internal/evidence/bounddir_windows.go	repo
internal/evidence/p84_bounddir_residual_test.go	repo
internal/guard/lease.go	repo
internal/guard/lease_lock_unix.go	repo
internal/guard/lease_lock_unsupported.go	repo
internal/guard/lease_lock_windows.go	repo
tmp/P84-BOUNDDIR-RESIDUAL/	local-only
```
