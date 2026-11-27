# P95 批量辅助：数字核对（原文）

载体：`deep-reasoner` + `GLM-5.3-Flash (glm)`。以下为输出逐字原文。

---

## P95 批量辅助 · 数字核对报告（只读，未写文件、未派子代理）

| # | 真值项 | 仓库对照 | 判定 |
|---|---|---|---|
| 1 | HEAD=origin=`22370fe` | 实测 HEAD/main/origin/main/origin/HEAD 均 `22370fe`（=「P94 闭环后加 v78 快照」提交） | PASS |
| 2 | 准则 v1.77 | `docs/DELIVERY_DIRECTIVE.md` 头部存在 `\| v1.77 \|` 行（P94 升 MINOR）；报告/台账/CHECKLIST_v78 一致 | PASS |
| 3 | SELFTEST 770/770 | CHECKLIST_v78 §一锚点 + P94 CI run `45910223000` 双腿 `PASS 770/770`；`git diff 49f1609..HEAD -- tools/gates` = 0 行 ⇒ 断言总数未变。未重跑（只读约束，与 ⑤⑦ 同口径） | PASS |
| 4 | 三作用域 TOTAL_FAIL=0 | P94 CI 双腿记录（REMAINING_SLICES L2532）+ 本轮最近 gate run exit 0 | PASS |
| 5 | v76/v77=216/217(LF) | 实测 `CHECKLIST_v76.md`=216、`CHECKLIST_v77.md`=217 LF，与 CHECKLIST_v78 §一「216/217」一致 | PASS |
| 6 | ADR 48 | `ADR_REGISTER.md` 表行 48、max=ADR-048（P87） | PASS |
| 7 | OB 145 | max 口径：OB-145 存在，`OB-14[6-9]/OB-15x` 命中 0；CHECKLIST_v78「OB max=145」。注：directive 表行计数 147，真值源按「max」口径读取，非差异 | PASS |
| 8 | DR 18 | max 口径：DR-18 见 CHECKLIST_v69–v73（P86 守卫轮红 run）；CHECKLIST_v78「DR max=18」，P94 核对未变 | PASS |
| 9 | 试点顺延 25 | CHECKLIST_v78 三处「试点顺延最大=第 25 次」（P94 未变）；briefing 同值 | PASS |
| 10 | 实测读数 | 主报告 L28–L40：1247/1177、1675、1596、1282–1493/796、610–642/471、1803–1905/1372、5799/3812、净自耗 0.6–0.9 s/≈0.2 s、10–15%/≈5%、<1%、base≈94 s、≤250 ms、≤0.3–4%、worker_threads≈0.9 s；P88 冷基线 20865/12525/10311 与 P88 baseline.tsv L26–28 对齐（⑥ S3/⑦ Q2 已核） | PASS |
| 11 | EN 镜像抽查 | `_EN.md` L4/L28/L34/L36/L40/L66/L124 同值（v1.77、1247、5799/3812、≈94、≤250、0.3-4%、1+1） | PASS |
| 12 | 台账 P95 段 | CN L2539–2553（🔄 执行中；额度口径 ⑦ 1+1/⑥ 1/主控 V4 Pro Max/零改动不升版）；EN 段存在（标识命中 6 处） | PASS |
| 13 | 审查链 | ⑤ r1 PASS（2 Low+3 Info）、⑥ r1 PASS（1 轮）、⑦ r1 FINDINGS（1 Medium+1 Low）、⑦ r2 PASS——四件原件逐字核对一致；额度 ⑦ 1+1、⑥ 1（报告 §6+台账）；D=0、探针 0。R=1 与「⑦ r1 返工 1 轮→r2」事实链一致；报告 §1/§6 仍记起跑基线 R=0 且 §5.3 明示 r2 结论由主控回填——属既定待回填态，非数字漂移 | PASS（附注） |

真值表全部 13 组读数与仓库文件一致；唯一软点为第 13 项 R 值回填态，已如实附注，不计差异。

NUMBERS: PASS
