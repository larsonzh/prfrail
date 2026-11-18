# P81 审查输入清单（`DIRECTIVE-MODEL-UPGRADE`）

> 作用：§5.1「时序说明（OB-93）」要求的输入清单——逐件列文件名与说明，且与申报输入一一对应；缺件即报错并停止派发。
> 编码：UTF-8 with BOM + LF。

## 1. 申报输入（⑤ / ⑥ / ⑦ 共用基线）

| 序号 | 工件 | 路径 | 说明 |
|---|---|---|---|
| 1 | 起跑消息原文 | `docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/00-briefing.md` | §13.1.10 工件；授权可回溯性依据 |
| 2 | 架构方案书 | `docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/01-arch/design-r1.md` | ① 产出；Q1–Q7 |
| 3 | 变更集 diff（含未跟踪新文件全文） | `tmp/P81-DIRECTIVE-MODEL-UPGRADE/impl.diff` | 由主控按当前磁盘态重新生成 |
| 4 | 变更集 numstat | `tmp/P81-DIRECTIVE-MODEL-UPGRADE/impl.numstat` | CN/EN 对称读数源 |
| 5 | 门禁原始读数（④，tree 与 selftest） | `tmp/P81-DIRECTIVE-MODEL-UPGRADE/gate-tree-04.txt`、`tmp/P81-DIRECTIVE-MODEL-UPGRADE/gate-selftest-04.txt` | 改后实测；index / ci 两作用域的权威读数按下列第 3 节口径在 ⑨ 取得 |
| 6 | 自由度清单 | `docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/freedom-list.md` | §1.7.3 固定形态；⑤ r1 后定稿 |
| 7 | 落点机械清单（主控与批量辅助两份） | `tmp/P81-DIRECTIVE-MODEL-UPGRADE/inventory.txt`、`docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/02-batch/batch-01-inventory.md` | 四组标识的全仓命中与三层归属 |
| 8 | 探针记账 | `tmp/P81-DIRECTIVE-MODEL-UPGRADE/probe-1.md` / `probe-2.md` / `probe-3.md` | §7.6 探针三次 |
| 9 | 判据源码与夹具 | `tools/gates/lib/criteria/g5b.js`、`tools/gates/selftest.js`、`tools/gates/testdata/fixtures/g5b-*.txt` | 闭集与夹具的机械对照面 |
| 10 | 准则正文对照 | `docs/DELIVERY_DIRECTIVE.md`、`docs/DELIVERY_DIRECTIVE_EN.md` | §2.2 / §5.2 / §6.3 / §7.5 / §8.2 / §12.4 / 附录 B.6 / 附录 C |
| 11 | 一致性核验脚本与输出 | `tmp/P81-DIRECTIVE-MODEL-UPGRADE/consistency-check.js`、`tmp/P81-DIRECTIVE-MODEL-UPGRADE/check-output.txt` | ③ 产出；C-1 至 C-10 |

## 2. 不申报输入（明确不递交审查层）

- 冻结证据包内既有历史载荷（`docs/validation/evidence/**` 的其它包）。
- `docs/t027/CHECKLIST_v*.md` 历史快照。
- 任何含凭据 / 主机名 / 用户名的材料（本片无此类输入）。

## 3. 逐轮实际递交记录

| 环节 | 实递条目 | 生成时点 |
|---|---|---|
| ① 架构 | 1、7、9、10（节选） | 步骤 0 完成后 |
| ② 实现 | 方案书全文（A1–A10 与 B1–B3 的精确文本） | ① 之后 |
| ③ 一致性 | 3、4、5、9、10 | ② 之后 |
| ⑤ 预审 r1 | 1、2、3、4、5、6、7、9、10 | ② 与门禁读数落盘后 |
| ⑤ 整改 | 6、11（定稿版）＋ ⑤ r1 的四条 Low | ⑤ r1 之后 |
| ⑥ 独立扫描 | 2、3、4、5、6、9、10、11 | ⑤ 整改后（本快照） |
| ⑦ 终审 r1 | 1、2、3、4、5、6、9、10、11 | ⑥ 之后 |
| ⑦ 复审 rN | 修复 diff ＋ 前轮发现清单 ＋ 机械读数 | 整改后 |

**index / ci 两作用域的读数口径（如实声明）**：本机 `--scope=index` 在未暂存时读 `changed=0` 而 fail-closed；`--scope=ci` 在无 CI 环境变量时退化为 `HEAD^1..HEAD`（不覆盖未提交改动）⇒ 故 ④ 的权威读数取 `--scope=tree` 与 `--selftest`，index / ci 两作用域在 ⑨（暂存后与提交后）取得并落入 `06-gates/`。

> **更新义务**：任何改动（含状态字段等最小修正）之后，派单给下一环前必须重新生成本清单与 `impl.diff` / `impl.numstat` / 门禁原始读数（§7.5 与「整改后检查清单」）。
