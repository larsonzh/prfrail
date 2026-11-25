# P88 GATES-PERF-ASSESS 输入清单（input-list，步骤 0 落盘）

切片：P88 `GATES-PERF-ASSESS`（门禁耗时构成评估 + 轻量效率优化 · 工具片 S 档）
落盘：2026-11-25（步骤 0）；本文件随片推进追加，⑧c 前只增不删。

## 一 权威输入（逐字承接，不扩不缩）

| 来源 | 内容 |
|---|---|
| 起跑消息 〇.5（§1.2 全字段） | `docs/validation/evidence/P88-GATES-PERF-ASSESS/00-briefing.md`（逐字落盘） |
| v71 快照 §四 队列项 | `docs/t027/CHECKLIST_v71{,_EN}.md` |
| 台账 P88 段（§1.2 全字段） | `docs/t027/REMAINING_SLICES{,_EN}.md`（P88 段起于 CN 第 2365 行） |
| 准则相关节 | `docs/DELIVERY_DIRECTIVE.md` §5.2（⑦ 分档表）/ §6.3（判据本体）/ §7.5（S 档 1+1）/ §13.1 / §13.2 |
| 判据实现（优化对象） | `tools/gates/gate.js`、`tools/gates/lib/ctx.js`、`tools/gates/lib/criteria/**`、`tools/gates/selftest.js` |
| P87 报告对（SELFTEST 744 与 G5-b 闭集基线） | `docs/validation/directive-review-model-policy{,_EN}.md` |
| 编码与行尾约定 | `docs/CODING_CONVENTIONS.md`（`.js` = 无 BOM + LF） |

## 二 步骤 0 读数（2026-11-25 实测）

| 锚点 | 期望 | 实测 | 判定 |
|---|---|---|---|
| HEAD | `fed0f5c` | `fed0f5c4daac03c78b7ba87ad0770a0813730783` | 相符 |
| origin/main | `fed0f5c` | `fed0f5c4daac03c78b7ba87ad0770a0813730783` | 相符 |
| 准则版本 | v1.72 | `版本：v1.72。日期：2026-11-24` | 相符 |
| SELFTEST | 744/744 | 见 `baseline.tsv` 与 ④ 读数 | 见下 |
| 三作用域 | 全 TOTAL_FAIL=0 | 见 ④ 读数节 | 见下 |
| v69 行数（LF） | 209 | CN 209 / EN 209 | 相符 |
| v70 行数（LF） | 203 | CN 203 / EN 203 | 相符 |
| v71 行数（LF） | 201 | CN 201 / EN 201 | 相符 |
| 工作区 | 干净 | `git status --porcelain` 空 | 相符 |

## 三 生效面（本片将动的文件）

| 对象 | 现状 | 本片动作 |
|---|---|---|
| `tools/gates/gate.js` | base 四件套串行、三作用域各跑一次 | 仅执行调度层优化（方案由 ① 定） |
| `tools/gates/**` 其它 | 见 ① 方案书 | 仅当 ① 点名，且零判定逻辑改动 |
| `docs/validation/gates-perf-assess{,_EN}.md` | 不存在 | 新建报告对（严格镜像） |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | P88 `🔄 执行中` | 翻 `✅ 闭环` + 完成记录 |
| `tools/gates/cjk-newwords.txt` | 现状白名单 | 仅新增 CJK 字时登记 |

## 四 边界（不做清单）

- 不改任何判据的判定逻辑；不触 §6.3 判据本体；不改 `internal/**`、
  CONTRACTS / schema / fixtures / workflows、CI workflow；不升准则 MAJOR；
  不改 G5-b 闭集与白名单；不做 gitee；不预设临时授权；不动
  `CHECKLIST_v71{,_EN}`、`DEV_PLAN{,_EN}`、s1-exit / p57 门表、
  `DOCUMENTATION_PLAN{,_EN}`。

## 五 片内异常登记

- 尚无（后续如实追加）。
