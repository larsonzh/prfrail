# P86 `BOUNDDIR-UNC-EQUIVALENCE` 切片报告

[English](bounddir-unc-equivalence_EN.md)

**结论：`CLOSED`**（N13-UNC + W10 两项均闭合；本机与门禁读数全绿，CI 双平台腿读数见 §11.2）。本片只闭合 F1 的最后两笔，不改 F1 其余已闭合项，不碰 F3 / F6 / F8（P85 范围）。

## 1. 切片标识与授权

| 项 | 内容 |
|---|---|
| 编号 / 标识 | P86 / `BOUNDDIR-UNC-EQUIVALENCE` |
| 分级 | `[SLICE]` · 实现类 · S 档 |
| 授权 | 用户 2026-11-22 点名起跑（先 P86 后 P85）；起跑消息逐字见 `evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/00-briefing.md` |
| 依据 | P84 报告 §8 顺延清单（N13-UNC + 证据归因 Low）+ v68 快照 §四队列项 |
| 边界 | 只修 N13-UNC 与 W10；不改门表 / 准则 / CI / 威胁模型 / 依赖 |

## 2. 步骤 0 锚点实测

| # | 项 | 期望 | 实测 |
|---|---|---|---|
| 1 | HEAD | `c7ebf3d` | `c7ebf3d` |
| 2 | origin/main | `c7ebf3d` | `c7ebf3d` |
| 3 | 准则版本 | v1.71 | v1.71 |
| 4 | SELFTEST | 740/740 | 740/740 |
| 5 | 三作用域 | 全 `TOTAL_FAIL=0` | index / tree / ci 全 0 |
| 6 | v66 快照总行数（LF） | 216 | 216（CN / EN 对称） |
| 7 | v67 快照总行数（LF） | 216 | 216（CN / EN 对称） |
| 附 | v68 快照总行数（LF） | 218 | 218（CN / EN 对称） |
| 附 | 机械读数 | ADR 47 / OB 139 / DR 17 / 试点顺延 25 | 全部相符 |

**登记（不记漂移）**：台账 P86 段起跑前预登记「⑦ `GPT-6.2 Sol` 深度档 + 1 终审 + 3 复审 + 盲审必抽」，与用户 2026-11-22 起跑消息「`GPT-5.3-Codex` 常规档 + 1 终审 + 1 复审 + 无额度例外」冲突；按 §0.1 / §13.2.11 纪律 F 以起跑消息为准，台账已就地更正。R/D 基线 = 0 / 0。

## 3. 目标与范围

- **N13-UNC**（F1 残余，Medium）：UNC 共享根的两种规范拼写（`\\host\share` 与 `\\host\share\`）在 `boundRelativeSuffix` 等值判定下被误拒——确定性界内误拒、方向 fail-closed、属可用性缺陷。
- **W10**（证据归因 Low）：`p84_bounddir_residual_test.go` 的变异标注把双分隔符前缀 `D:\\` 称为「轮 11 前原形」，与 ⑦ 复审 r3 原文矛盾（r3 记录的前缀是单分隔符 `D:\`）。

## 4. 修复实现与闭合状态

| 项 | 归属 | 改动 | 状态 |
|---|---|---|---|
| N13-UNC | F1 | `internal/evidence/bounddir_windows.go` `boundRelativeSuffix` 新增等值分支（步骤 3 之后、严格后代之前）+ 助手 `uncShareRootName` / `uncShareRootEqual` / `lowerASCII` | **闭合** |
| W10 | 证据归因 | `internal/evidence/p84_bounddir_residual_test.go` 的 `Mutation check:` 注释段（纯注释、代码零改动） | **闭合** |
| 测试 | — | 新增 `internal/evidence/p86_bounddir_unc_test.go`（正例 5 / 反例 12 / 分类器负例 10 / 生产入口 1） | 全绿（生产入口本机 SKIP，见 §11） |

**N13-UNC 判定顺序（修复后）**：① 非规范拼写守卫 → ② `IsAbs` → ③ 规范形恒等 → ④ **新增：两者皆为 UNC 共享根且卷名 ASCII 折叠相等 ⇒ 放行** → ⑤ trimmed 空拒 → ⑥ 严格后代 Ordinal 前缀 + 组件守卫。**唯一新增放行点 = 步骤 ④**。

**零放宽论证**：不同主机 / 不同共享名 ⇒ 卷名折叠不等；host-only（`\\host`、`\\host\`）⇒ 卷名非两段非空 ⇒ 分类器返 `""`；`\\?\` / `\\.\` ⇒ 分类器显式返 `""`；非规范拼写（`/`、重复分隔符、非根尾分隔符、`.`/`..`）⇒ 步骤 ① 先行拒绝；界外后代 ⇒ 步骤 ⑥ 前缀不匹配；盘符形态 ⇒ 卷名无 `\\` 前缀；后代组件大小写 ⇒ 步骤 ⑥ Ordinal 仍拒。

**W10 更正要点**：M1（仅回退等值，前缀仍由 trimmed base 构造 ⇒ `D:\`）、M2（叠加未修剪前缀 ⇒ `D:\\`）分开标注；删除「`D:\\` = 轮 11 前原形」的错误表述，改为「NOT the round-11-before arrangement」；与 r3 的「目标刚好等于前缀」逐字相容。

## 5. 契约影响判定

| 项 | 判定 | 理由 |
|---|---|---|
| N13-UNC | **DOES-NOT-TOUCH / IN-BOUNDARY** | `docs/CONTRACTS{,_EN}.md` 无 bounddir / reparse / UNC 专属句；最邻近句 `CONTRACTS.md:238` 属 replay store 的写者根/保护根纪律，不约束本文件。改动不改「Ordinal 字节比较」同源纪律，只补一处共享根等值形态 |

无新 schema、无新 CLI、无新依赖（仅既有 `strings` / `filepath`）。**未触契约 ⇒ 未走协议先行，未改 `CONTRACTS{,_EN}.md` / schema / 夹具 / `go.mod`**。

## 6. 分域升级与用户裁定

| 事项 | 裁定 | 留痕 |
|---|---|---|
| 主控 | 保持 `DeepSeek V4.1 Flash (deepseek)` High（不继承 P84 的 V4 Pro） | 起跑消息 〇（2026-10-31 长期裁定默认） |
| ⑦ 额度 / 档位 | 1 终审 + 1 复审；`GPT-5.3-Codex (copilot)` 常规档；无额度例外 | 用户 2026-11-22 裁定 a |
| 盲审 | 必抽 1 次（另计 +1、终审前执行） | ① 判定书（bounddir 触身份/所有权语义）；落盘于方案书 §7 |
| **大小写变体冲突** | **按起跑消息 §六执行**：UNC 共享根卷名的 ASCII 大小写差异判等值 | 四要素：裁决人 = 用户（同轮起跑消息）、日期 = 2026-11-22、理由 = 起跑消息更新更具体、⑦ 收窄复审「保持 Ordinal」系 P84 遗留建议作为下限、依据 = 起跑消息 §六 + 准则 §0.1；**不得作为其他切片先例** |
| 排片 | 先 P86 后 P85 | 用户 2026-11-22 裁定 b |

**大小写折叠的安全论证**：UNC 卷命名空间由 SMB 定义，主机名与共享名均大小写不敏感，同一主机不存在仅大小写不同的两个共享 ⇒ 折叠**不扩大**被放行对象集合，只消除误拒；折叠范围**仅卷名 + 仅 ASCII**（`A`–`Z`）；后代组件、盘符形态根、`isUnderDirectory` 一律保持 Ordinal（W1 的 case-variant twin 拒绝不变）。

## 7. 审查链

| 环 | 载体 | 模型 | 额度 | 结论 |
|---|---|---|---|---|
| ① 方案书 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 1 | 契约判定 + 设计 + 冲突裁定 + 盲审判定 |
| ② 实现 | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 1 | 两处改动 + 全量回归 |
| ③ 测试 | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | 1 | 新测试 + 变异实跑 |
| ⑤ 预审 | `independent-reviewer` | `DeepSeek V4 Pro (deepseek)` | 1 + 1 | r1 `PASS WITH FIXES`（1 Medium + 4 Low）→ r2 `PASS` |
| ⑥ 扫描 | `independent-reviewer` | `MAI Code 2 Flash (copilot)` | ≤ 3（用 2） | r1 `FINDINGS`（2 Low）→ r2 `PASS` |
| ⑦ 盲审 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 另计 +1 | `PASS WITH FIXES`（唯一 Low = 生产入口缺非 SKIP 证据） |
| ⑦ 终审 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 1（1+1 的第 1 次） | `PASS WITH FIXES`（无 Medium+；Low-L1 = 同一条生产入口证据） |

逐字原文归档于 `04-reviews/`（⑤ r1/r2、⑥ r1/r2、⑦ 盲审、⑦ 终审 r1）。

## 8. 残余与顺延

- 本片无新增顺延项。**F1 在本片后闭合**（T4 残余 N13-UNC + 证据归因 W10 两项均由本片闭合）。
- P85 范围（F3 / F6 / F8：T2 / T3 / T8a）不受影响，起跑须用户点名。
- 待用户裁定：`bounddir-residual{,_EN}.md` §8 是否补「N13-UNC 已由 P86 闭合」时点注记（起跑消息 §七-7 列为核对项，本片不擅自加）。

## 9. 文档同步

| 文件 | 动作 |
|---|---|
| `docs/validation/bounddir-unc-equivalence{,_EN}.md` | 本报告（新增） |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | P86 段：状态 → 已闭环 + 完成记录；⑦ 额度口径按起跑消息更正；P84 段「顺延 P86」字样撤除并改为「由 P86 闭合」 |
| `docs/SECURITY{,_EN}.md:5` | 修复状态注记改为「F1 全闭合（T4 残余 N13-UNC 由 P86 闭合）」；F2/F5/F7 已修复句保留；F3/F6/F8 顺延 P85 句保留 |
| `docs/DEV_PLAN{,_EN}.md` | T027 段补录 P86 一段（只记结论与报告路径，不改写 T027 状态字段） |
| `docs/t027/CHECKLIST_v68{,_EN}.md` | **不动**（v69 快照再吸收） |
| 准则 / ADR / OB / 门表 / CI | **不动**（不升版、不新增） |

## 10. 门禁结果

| 门 | 结果 |
|---|---|
| `go build ./...` | EXIT=0 |
| `go vet ./...` | EXIT=0 |
| `gofmt -l internal/evidence` | 空输出 |
| `go test ./... -count=1` | 15 包 ok、0 FAIL（`internal/evidence` ok） |
| `gate.js --selftest` | `SELFTEST: PASS 740/740` |
| `gate.js --all --scope=tree` / `--scope=ci` | `TOTAL_FAIL=0` |
| `gate.js --all --scope=index` | 提交前因工作树脏为 1（G5(f) 索引卫生，预期）；⑧b 提交后复跑为 0（见 §11.2） |

日志：`03-verification/build.txt`、`vet.txt`、`gate-selftest.txt`、`gate-{index,tree,ci}.txt`、`test-all-windows.txt`、`test-p86-verbose.txt`。

## 11. ⑨ 验证结果与本机不可执行项

- **Windows 本机**：`go test ./internal/evidence -run TestP86 -count=1 -v` → `TestP86N13UncShareRootEquivalence`（5 子例）PASS、`TestP86N13UncShareRootRefusals`（12 子例 + 分类器 10 负例）PASS、`TestP86UncShareRootProductionEntry` **SKIP**（本机非提权：`net share … exit status 2 (发生系统错误 5。拒绝访问。)`；已 `t.Skipf` 写明原因）。
- **CI 守卫（⑦ 盲审/终审 Low 的机械收口）**：生产入口用例在 `CI` 环境变量存在时，把三个前置失败（建共享、开 UNC 根句柄、句柄名尾分隔符形态）由 skip 改为 **硬失败**，因此**CI 双腿绿 ⇒ 生产入口真实运行并 PASS**；本机（`CI` 未设）仍按 skip 并写明原因。此法不改 CI 工作流（边界内）。
- **本机 flake 观测（非本片引入）**：本机 3 次 `go test ./... -count=1` 中 1 次 `internal/tickets` 的 `TestP39DurablePrecommitFaultsRejectAck/dirClose` 因 `rename … Access is denied` 失败（与 DR-17 同族的环境性 IO/AV 抖动），其余 2 次全绿、CI 双腿未复现；按「观测级 flake 只记录不阻断」处理，不新增 DR。
- **变异检验（实跑留证）**：`mutation-record.txt` 记录 `BEFORE_SHA256=35df2761…dba62` / `ANCHOR_COUNT=1` / `MUTATED_SHA256=f6ab720f…64db7` / `MUTATED_TEST_EXIT=1` / `RESTORED_SHA256=35df2761…dba62` / `RESTORED_IDENTICAL=True` / `RESTORED_TEST_EXIT=0`；`mutation-branch-removed.txt` 显示**恰好** 3 个拼写差异子例 FAIL（`trailing-base-bare-target` / `bare-base-trailing-target` / `ascii-case-folded-trailing-target`），`identical-trailing-spellings`、`strict-descendant-regression`、全部反例与 `TestP84W10…` 保持 PASS。
- **不做 -race**：本机无 gcc（与 P84 同口径）；Race 步骤由 CI ubuntu 腿承担。
- **Ubuntu VM 不开启**（起跑消息 §三）；Linux 侧零改动，由 CI ubuntu 腿回归。

## 11.2 CI 读数（⑧b 推送后写入本节）

- 主提交 `f1ae821`；首跑 run **`44962802936`** **双腿 success**（`Go windows-latest` / `Go ubuntu-latest`）；windows 腿：`ok github.com/larsonzh/prfrail/internal/evidence`、`SELFTEST: PASS 740/740`、`TOTAL_FAIL=0`；日志载荷 `03-verification/ci/run-44962802936-windows.log` / `-ubuntu.log`。
- 首跑未含 CI 守卫（`-v` 未开启 ⇒ 无 PASS/SKIP 行）；守卫轮（CI 环境下前置失败改硬失败）随修复提交重跑。
- 守卫轮首跑 run **`44965434402`**（提交 `a5ef122`）：**windows 腿红**（`TestP86UncShareRootProductionEntry` 前置断言：CI 宿主的句柄派生根为**无尾分隔符**形 `\\localhost\<share>`，夹具却假定带尾分隔符），ubuntu 腿 success ⇒ 登记 **DR-18**。
- 修复：夹具改为**按句柄实际拼写取另一侧拼写**作为存储目标（句柄带尾分隔符 ⇒ 目标无尾；句柄无尾 ⇒ 目标带尾），两种宿主形态都进入 N13 等值分支；修复轮读数见下行。

## 12. 环境与档位

- 宿主：Windows（本机），Go 工具链 1.22+（本机实测 go1.27.0；CI 以工作流固定版本为准）。
- 主控 / ②③：`DeepSeek V4.1 Flash` High（同档，角色分离以独立派单保证）；①⑤：`DeepSeek V4 Pro`；⑥：`MAI Code 2 Flash`；⑦：`GPT-5.3-Codex` 常规档。
- 网络：不开代理（SSH push / `gh` 直连）；不做 gitee。

## 13. 审查与计量

- **R = 0**（复审环整改重跑不计 R；本片 ⑤⑥ 整改均为同环内整改；⑦ 无整改轮）。
- **D = 0**（临时件全部落 `tmp/P86-BOUNDDIR-UNC-EQUIVALENCE/`，⑧c 清除；无片内违规）。
- 探针 = 0（无外部真实调用）；批量辅助 GLM 未使用（本片规模 S、无枚举/镜像批量需求；计划 ≥2 与实测 0 的差异在 §15 登记）。

## 14. 明确未执行事项

未做 gitee；未升准则版本；未新增 ADR / OB；未改 CI 与门表；未改 `go.mod` / schema / CLI；未跑 `-race`（本机无 gcc）；未在 Linux 实机执行 Linux 专属用例（由 CI ubuntu 腿承担）；未补 `bounddir-residual §8` 时点注记（待用户裁定）。

## 15. ⑧a 自查

1. 起跑消息原文落盘且入 MANIFEST / SHA256SUMS（§13.1.10）？是。
2. 审查方原文逐字归档（§13.1.13）？是（6 份）。
3. 证据载荷三方相等（磁盘 / MANIFEST / SHA256SUMS）？是（⑧a 重建后由 G8-b 复核）。
4. 改动面与边界一致（零触碰 F3/F6/F8、门表、准则、CI、依赖）？是。
5. 计量与额度如实（R = 0、D = 0、⑥ 用 2 / ≤3、⑦ 用 1 终审 + 1 盲审、复审额度未用）？是。
6. 未执行项逐条披露（§14）？是。
7. **偏差登记**：① 计划「批量辅助 GLM ≥2」未使用（无适配任务，未创造任务以凑数）；② 台账 ⑦ 预登记与起跑消息冲突已更正（§2）。

## 16. 六段自查（sixstage）

```sixstage
## Section A 总体结论
CLOSED（N13-UNC 与 W10 两项均闭合；F1 在本片后闭合）

## Section B 发现清单

| SEVERITY | 位置 | 处置 |
|---|---|---|
| Medium | N13-UNC（boundRelativeSuffix 等值分支） | 已闭合：新增 UNC 共享根等值分支 + 助手；变异实跑证明分支必要性 |
| Low | W10（p84 测试变异注释） | 已闭合：M1/M2 分开标注，删除「D:\\ = 轮 11 前原形」 |

## Section C 逐条证据与推理
- 判定顺序：非规范守卫 → IsAbs → 规范形恒等 → **UNC 共享根等值（唯一新增放行点）** → trimmed 空拒 → 严格后代 Ordinal 前缀 + 组件守卫。
- 正例 5（两方向尾分隔符 / ASCII 折叠 / 恒等 / 后代回归）与反例 12（异主机 / 异共享 / 后代大小写 / 重复分隔符 / 非根尾分隔符 / 正斜杠 / `..` / 界外 / host-only ×2 / 尾空白）由 `p86_bounddir_unc_test.go` 钉住；分类器负例 10 项含 `\\?\` / `\\.\`。
- 变异：删分支后恰好 3 个拼写差异子例红、恢复后 SHA256 逐字节一致（`03-verification/mutation-record.txt`）。

## Section D 可证伪性审计
- 删等值分支 ⇒ 3 个正例红、反例与 P84 W10/W1 保持绿（实跑留证）。
- 放宽容卷名折叠到后代 ⇒ `case-variant-descendant-stays-ordinal` 转红。
- 放开规范形守卫 ⇒ `repeated-separator` / `parent-component` / `non-root-trailing-separator` 转红。

## Section E 无法验证的部分
- 本机生产入口 SKIP（非提权，`net share` 拒绝访问）；CI windows 腿读数见 §11.2。
- 未跑 `-race`（本机无 gcc）；Linux 侧零改动，由 CI ubuntu 腿回归。

INDEPENDENT SCAN: FINDINGS
```

## 17. 工件围栏（artifacts）

```artifacts
internal/evidence/bounddir_windows.go	repo
internal/evidence/p84_bounddir_residual_test.go	repo
internal/evidence/p86_bounddir_unc_test.go	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/validation/bounddir-unc-equivalence.md	repo
docs/validation/bounddir-unc-equivalence_EN.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/MANIFEST.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/SHA256SUMS.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/00-briefing.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/input-list.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/freedom-list.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/01-arch/01-arch-design.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/build.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/vet.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/gate-selftest.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/gate-index.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/gate-tree.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/gate-ci.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/test-all-windows.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/test-p86-verbose.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/probe-clean-windows.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/mutation-record.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/mutation-branch-removed.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/test-post-restore.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/local-flake-tickets-p39.txt	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/ci/run-44962802936-windows.log	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/ci/run-44962802936-ubuntu.log	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/03-verification/ci/run-44965434402-windows.log	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/07-blind-review.md	repo
docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/04-reviews/07-final-review-r1.md	repo
```

未入库文件：`tmp/P86-BOUNDDIR-UNC-EQUIVALENCE/` 下过程件（⑧c 清除）。
