# input-list — P63 AT23-RECHECK（步骤 0 落盘）

> 本片输入件清单（截至步骤 0）；后续各环派单输入包在此基础上按环节收窄。

## 核心判据与裁定输入

1. `docs/validation/r08-r09-ruling.md`（CN 权威）——R-08.5「继承漂移豁免」五要素原文（§2，P62 修订版）。
2. `docs/validation/r08-r09-ruling_EN.md`（EN 镜像）——与 CN 逐行同位。
3. `docs/validation/evidence/R08-R09-RULING/`——P62 修订后的裁定文档副本与证据。

## 判定对象事实输入

4. `docs/validation/at23-restricted-candidate.md` 与 `_EN.md`（P61 报告）——§1 / §2.4 / §3 / §4 / §5 / §8.3（PARTIAL 结论、F2 漂移事实、⑦ 盲审与终审摘要）。
5. P61 层 3 证据：`docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt`、`l3-f2-pack-verdict.json`、`l3-f2-evidence/f2-pos/verdict-run.json`。
6. B4 基线：`docs/validation/t027-at23-e2e.md`（B4 基线报告）与 `docs/validation/evidence/directive-history/b4/archive/`（含 `pack-verdict.json`、`evidence/f2-pos/verdict-run.json`、`probe/f2-tree-calibration.txt`）。
7. 成因提交：`02d5797`（2026-10-22，A5 pid reports 绑定为第六枚冻结事实）改动摘要（`internal/adapters/agent_runner_postflight.go` Leg 7、`tools/b4-e2e/main.go` f1 侧与 facts 适配器）；`9ba621d`（2026-10-21，stub 子报告契约）作为因果链前置事实。
8. 现状代码：`tools/b4-e2e/main.go`（F2 检查表达式）、`tools/agent-stub/main.go`（`hang` 分支）。

## 切片锚点与队列输入

9. `docs/t027/CHECKLIST_v52.md`（跨会话快照，只读；v53 生成时继承）——§一 锚点、§四 执行队列（层 1 第 11 项）。
10. `docs/t027/REMAINING_SLICES.md` 与 `_EN.md`——P63 段登记处、P61 / P62 段、状态表。
11. `docs/ADR_REGISTER.md` 与 `_EN.md`——ADR-037 格式参照（ADR-038 登记处）。

## 准则引用输入（只读，不改）

12. `docs/DELIVERY_DIRECTIVE.md` 与 `_EN.md` v1.64——§1.7.2 三条 bullet、§1.7.4 独立留痕、§4.1 结论行闭集、§5.2 文档型切片、§7.5 ⑦ 额度、§7.7 必抽盲审、§11.1 历史行不改写。

## 文档同步输入（⑧ 阶段使用）

13. `docs/validation/s1-exit.md` 与 `_EN.md`——AT-23 条目 / S-G1 / S-G7 依赖链。
14. `docs/DEV_PLAN.md` 与 `_EN.md`——T027 段 AT-23 状态。
15. `docs/TEST_STRATEGY.md` 与 `_EN.md`——AT-23 行验收结论。

## 环境声明

- 不开启 Ubuntu VM、不开启网络代理、不启用 GPT-6.1 Sol 双打；无外部调用、无 syscall；纯 .md 判定。
