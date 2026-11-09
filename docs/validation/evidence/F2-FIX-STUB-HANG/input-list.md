# input-list — F2-FIX-STUB-HANG（步骤 0 落盘）

> 本片输入件清单（截至步骤 0）；后续各环派单输入包在此基础上按环节收窄。

## 契约依据输入

1. `docs/CONTRACTS.md` 与 `_EN.md` §7（A6 段）——「pid 报告写入口径（恒定写）」原文（无子进程写 0 字节空文件；缺失按证据缺口 fail-closed；消费方三态语义）。
2. `docs/CONTRACTS.md` A5 段——postflight pid 报告腿「重读重算对账」原文（`stub.pid` + `stub-child.pid`）。

## 修复对象输入

3. `tools/agent-stub/main.go`——`hang` 分支（L112–116 附近：只 `writeOwnPID()` 后睡眠）、`writeEmptyChildPIDReport()`（L257+）、`writeOwnPID()`（L252+）、公共路径写序（L133–141）。
4. `tools/agent-stub/main_test.go`——既有回归测试（`TestRunWithoutChildWritesEmptyChildPIDReport` / `TestSpawnWritesChildPIDList` / `TestUnwritableChildPIDReportFailsClosed` / `TestSpawnWithUnwritableChildPIDReportKillsChildAndFails`）与测试骨架（`buildStub` / `runStub` / `readReport`）。
5. `tools/b4-e2e/main.go`——`process-tree-taken-away` 断言（L766 附近）与 `six-facts-distinct` 断言（L690/L697 附近）。
6. `tools/b4-e2e/run-faces.ps1`——`-OnlyFace f2` 重跑入口（f2-pos / f2-neg / f2-neg-r2 三腿）。

## 判定事实输入

7. P63 报告：`docs/validation/at23-restricted-candidate.md` 与 `_EN.md` §10（豁免不成立 + 后续动作 + v53 观察项候选）。
8. P61 报告：`docs/validation/at23-restricted-candidate.md` §2.4（层 3 F2 事实）、§8.3（⑦ 盲审与终审）。
9. P61 层 3 证据：`docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt`、`l3-f2-evidence/`（三腿 verdict-run.json）、`l3-f2-pack-verdict.json`。
10. B4 基线：`docs/validation/evidence/directive-history/b4/archive/pack-verdict.json` 与 `evidence/f2-pos/verdict-run.json`。

## 切片锚点与队列输入

11. `docs/t027/CHECKLIST_v52.md`（跨会话快照，只读）——§一 锚点、§四 层 1 第 10 项。
12. `docs/t027/REMAINING_SLICES.md` 与 `_EN.md`——F2 修复片登记处、P61/P62/P63 段。
13. `docs/ADR_REGISTER.md` 与 `_EN.md`——ADR-038 格式参照（ADR-039 登记处）。

## 准则引用输入（只读，不改）

14. `docs/DELIVERY_DIRECTIVE.md` 与 `_EN.md` v1.64——§1.7.2 三条 bullet、§1.7.4 独立留痕、§5.2 修复类流水线、§7.5 ⑦ 额度、§7.7 必抽盲审、§11.1 历史行不改写。

## 环境声明

- 不开启 Ubuntu VM、不开启网络代理、不启用 GPT-6.1 Sol 双打；Windows 本机执行；无外部调用。
