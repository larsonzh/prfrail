# P98 input-list（输入清单）

- 起跑消息原文：本目录 `00-briefing.md`（用户 2026-11-28 逐字原文，§13.1.10）。
- 权威输入：起跑消息 〇.5 的 P98 全字段表（§1.2 全字段为本片唯一落地依据）。
- 必读文档：
  1. `docs/t027/CHECKLIST_v81{,_EN}.md`（§一 锚点 / §五 O-1..O-5 登记原文）；
  2. `docs/DELIVERY_DIRECTIVE{,_EN}.md` §5.2 / §7.5 / §8.2 / §0.3 / §12.4 / §13.1.10 / §13.2.2；
  3. `docs/validation/hydrogen-control-pilot{,_EN}.md`（O-2..O-4 原文）与 `docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/04-reviews/07-pilot-hydrogen.txt`（Hydrogen 试水原文，发现 #1..#5）；
  4. `docs/t027/REMAINING_SLICES{,_EN}.md` P97 段（观察项 O-2..O-4 / O-5 登记行）与 P98 登记位；
  5. `docs/CODING_CONVENTIONS.md`（.md = BOM + LF）。
- 步骤 0 实测锚点（2026-11-28 本机实读）：HEAD = `ada55b2`；origin/main = `ada55b2`；准则 = v1.82；SELFTEST = 806/806；三作用域（`--all --skip-base`）= TOTAL_FAIL=0；v79 / v80 / v81 = 209 / 210 / 213（LF）；ADR 48 / OB 146 / DR 18 / 试点顺延 25；tmp 顶层 = `.gitkeep` + `gates-cache`；工作区 porcelain clean。
- 修订目标三节现文（v1.82）：
  - §5.2 ⑦ 评审模型选型矩阵（深度档行 = Sol 主读 + Neon 交叉复核、思考档位列单写 `Max`）+ 判级门槛共识 / 组合触发条款 / 适用域排除句 / 成本与档位边界；
  - §7.5 ⑦ 成本定义单价表 + 上下文档位计价注记 + 备选额度口径 + 四读实测消耗注记 + 兜底句；
  - §8.2 ⑦ 档位映射表 + 取档授权句（末段）。
- O-2..O-4 / O-5 原文出处：`07-pilot-hydrogen.txt` 的 Section B 发现 #1（High）/ #2（Medium）/ #3（Medium）/ #4（Low）/ #5（Low）。
