# P87 DIRECTIVE-REVIEW-MODEL-POLICY 输入清单（input-list，步骤 0 落盘）

切片：P87 `DIRECTIVE-REVIEW-MODEL-POLICY`（⑦ 评审模型选型矩阵 + 判级门槛共识 + ② 加强轮触发 · 治理片 M 档）
落盘：2026-11-24（步骤 0）；本文件随片推进追加，⑧c 前只增不删。

## 一 权威输入（逐字承接，不扩不缩）

| 来源 | 内容 |
|---|---|
| 起跑消息 〇.5（§1.2 全字段 + §5.4 条款） | `docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/00-briefing.md`（逐字落盘） |
| v70 快照 §四 队列项 + §五 四读数据指针 | `docs/t027/CHECKLIST_v70{,_EN}.md` |
| 四读同输入对照表 | `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/P87-four-read-data/07-three-read-comparison.md` |
| 四读四份报告 + r3 原文 | 同目录 `07-carbon-main-read.md` / `07-neon-cross-read.md` / `07-sol-baseline-read.md` / `07-codex-baseline-read.md` / `07-re-review-r3.md` |
| P85 报告对（四读数据出处） | `docs/validation/security-proof-subsystem{,_EN}.md` |
| 准则相关节 | `docs/DELIVERY_DIRECTIVE.md` §2.2 / §2.4 / §5.2 / §5.3 / §6.3（G5-b）/ §7.5 / §7.7 / §8.2 / §12.4 / §13.1 |
| ADR 先例形态 | `docs/ADR_REGISTER{,_EN}.md`（ADR-019 形态、ADR-047 max） |
| 台账登记位置 | `docs/t027/REMAINING_SLICES{,_EN}.md`（P86 段之后、`## 已知缺陷` 之前） |
| 判据实现 | `tools/gates/lib/criteria/g5b.js`（`CLOSED_SET`）+ `tools/gates/selftest.js`（T29–T42 / T131a–c 夹具族） |
| 夹具目录 | `tools/gates/testdata/fixtures/g5b-*.txt` |

## 二 步骤 0 读数（2026-11-24 实测）

| 锚点 | 期望 | 实测 | 判定 |
|---|---|---|---|
| HEAD | 010dd52 | `010dd52` | 相符 |
| origin/main | 010dd52 | `010dd52` | 相符 |
| 准则版本 | v1.71 | v1.71 | 相符 |
| SELFTEST | 740/740 | `SELFTEST: PASS 740/740` | 相符 |
| 三作用域 | 全 TOTAL_FAIL=0 | index 0 / tree 0 / ci 0（判据面；见 `03-verification/`） | 相符 |
| v68 行数（LF） | 218 | CN 218 / EN 218 | 相符 |
| v69 行数（LF） | 210 | CN 209 / EN 209 | **单项不符**（<2 项，如实登记、不改冻结件） |
| ADR max / OB max / DR max / 试点顺延 | 47 / 144 / 18 / 25 | 47 / 144 / 18 / 25 | 相符 |

## 三 生效面（本片将改的条文与判据）

| 对象 | 现状 | 本片动作 |
|---|---|---|
| §2.2 角色表 | 无 Gemini 系成员 | 增三行（Neon / Carbon / Hydrogen）+ Sol 常规档性价比备选句 |
| §5.2 | ⑦ 三模型分档表；「同环不混用」无组合例外 | 增四档选型矩阵 + 判级门槛共识句 + 组合触发条款 + 适用域排除句 |
| §5.3 之后 | 无 §5.4 | 新增 §5.4「② 加强轮触发与回退」 |
| §6.3 G5-b | 闭集 8 枚 | 扩为 11 枚（加 `gemini-4-neon` / `gemini-4-carbon` / `gemini-4-hydrogen`） |
| §7.5 | 单价表 3 行 | 增 3 行 + 上下文档位计价注记 + 四读实测消耗注记 + 备选额度口径 |
| §8.2 | ⑦ 档位映射 3 行 | 增 3 行（Carbon = High 固定、Neon / Hydrogen = Extra High）+ 取档授权句 |
| `tools/gates/lib/criteria/g5b.js` | `CLOSED_SET` 8 枚 | 11 枚（与文档逐字一致） |
| `tools/gates/selftest.js` | 740 断言 | 新增夹具断言（三枚合法 + 一枚非法 + 载荷性变异） |

## 四 边界（不做清单）

- 不新增本草案清单以外的模型行；不改 ⑤⑥ 现役首选；不改深度档 Sol Max 首选地位；
  不触 CONTRACTS / schema / fixtures / workflows；不升准则 MAJOR；不动 `internal/**`；
  Grok 4.7 不纳入；P85 在飞条款不追改；不预设临时授权。

## 五 片内异常登记（步骤 0）

- 编码归一误写（详见 `freedom-list.md` 起跑补充）：五个证据件被截为 3 字节后同轮重写；
  处置 = 断言长度写回 + 逐件回读字节；影响 = 无残留、无仓库外位置；不计 D。
