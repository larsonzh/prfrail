# P86 BOUNDDIR-UNC-EQUIVALENCE 输入清单（input-list，步骤 0 落盘）

切片：P86 `BOUNDDIR-UNC-EQUIVALENCE`（N13-UNC / W10 残余 · 闭合 F1 最后两笔 · 实现类 S 档）
落盘：2026-11-22（步骤 0）；本文件随片推进追加，⑧c 前只增不删。

## 一 权威输入（逐字承接，不扩不缩）

| 来源 | 内容 |
|---|---|
| P84 报告 §8（顺延清单） | `docs/validation/bounddir-residual.md` §8：N13-UNC + W10 另立 P86（2026-11-22 用户裁定、不并入 P85） |
| P84 报告 §5（证据归因行） | 同上 Section B：N13-UNC = UNC 共享根两种规范拼写的等值误拒（确定性界内误拒、方向 fail-closed、属可用性缺陷）；W10 = 变异标注「轮 11 前原形」与 r3 原文矛盾 |
| ⑦ 收窄复审原文（N13 / W10 逐字） | `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-narrow-re-review.md` |
| ⑦ 复审 r3 原文（N13 最小修复与 M1/M2 描述） | `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-re-review-r3.md` |
| F1 修复状态注记 | `docs/SECURITY.md:5`（F1/F4 部分闭合；T4 残余 = N13-UNC 顺延 P86） |
| 契约被引用句 | `docs/CONTRACTS.md` / `docs/CONTRACTS_EN.md`（bounddir 相关句；以 ① 判定书为准） |
| 准则相关节 | `docs/DELIVERY_DIRECTIVE.md` §1.4 / §5.2 / §5.3 / §7.5 / §7.7 / §8.2 / §9.2 |
| 队列锚点 | `docs/t027/CHECKLIST_v68{,_EN}.md` §一 / §四（P86 = 下一项；先 P86 后 P85） |

## 二 现状代码（步骤 0 读数）

| 文件 | 行数 | 说明 |
|---|---|---|
| `internal/evidence/bounddir_windows.go` | 1735 | `boundRelativeSuffix`（:991）、`nonCanonicalReparseTarget`（:939）、`normalizeReparseTarget`（:910）、`reparseTargetRelativeSuffix`（:734）、`boundRelativeSuffixShortAlias`（:1315，仅盘符形态） |
| `internal/evidence/p84_bounddir_residual_test.go` | 1686 | `TestP84W10VolumeRootKeepsItsSeparatorForEquality`（:911 起；变异标注注释 :913-925） |
| `internal/evidence/bounddir_unix.go` / `_darwin.go` | — | 零改动（Linux 侧由 CI ubuntu 腿回归） |

## 三 N13-UNC 残余（F1，逐字要点）

1. 基准为 `\\host\share\`、目标为 `\\host\share` 时，两者都是 Windows
   `filepath.Clean` 保留的绝对规范形，指向同一共享根。
2. 旧等值分支接受；现等值分支不匹配（`target != cleanBase`），严格后代分支
   又因目标短于前缀而拒绝 ⇒ **确定性界内误拒**（非 Rev 5 有意拒绝的非规范拼写）。
3. 最小修复（⑦ 收窄复审原文）：**仅对同一 UNC 共享根接受有、无尾分隔符两种
   等值形态，保持主机名和共享名的 Ordinal 比较**，并补两个方向的等值测试。
4. 方向纪律：非等值形态（不同主机 / 不同共享名 / 中间 reparse / 非规范拼写）
   仍拒，fail-closed 不削弱；盘符根等值（P84 已修）回归不红。

## 四 W10 证据归因（Low，逐字要点）

1. `p84_bounddir_residual_test.go` 注释把双分隔符前缀（`D:\\`）称为「轮 11 前
   原形」，与 r3 原文矛盾：r3 的 N13 描述为「严格后代分支又因目标刚好等于
   前缀而拒绝」⇒ 当时前缀 = `D:\`（trimmed + separator 构造）。
2. 最小修复（⑦ 收窄复审原文）：**将旧等值回退与额外破坏前缀的复合变异分开
   标注**（纯文档、可字面核验）。

## 五 契约影响判定输入（① 必读）

- 守则：任何实现若需新 schema 记录、新 CLI 命令、新依赖或第三方包
  ⇒ OUT-OF-BOUNDARY，停下升级，不自行扩域。
- 若判定触契约 ⇒ 协议先行：先改 `docs/CONTRACTS{,_EN}.md`（及 schema /
  夹具如涉）再动代码。
- P84 同类项先例：T4 判 `DOES-NOT-TOUCH / IN-BOUNDARY`（被引用句只定要求
  不定机制）；本片须对 N13-UNC 残余**重新判定**。

## 六 必读清单（对应起跑消息五）

1. `docs/t027/CHECKLIST_v68{,_EN}.md`（§一 锚点 / §四 P86 队列项）
2. `docs/validation/bounddir-residual{,_EN}.md`（§8 顺延、§4 T4 残余、§7 收窄复审）
3. `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-narrow-re-review.md` + `07-re-review-r3.md`
4. `docs/SECURITY{,_EN}.md:5`；`docs/CONTRACTS{,_EN}.md`（以 ① 判定书为准）
5. `docs/DELIVERY_DIRECTIVE.md` §1.4 / §5.2 / §5.3 / §7.5 / §7.7 / §8.2 / §9.2
6. `internal/evidence/bounddir_windows.go`、`internal/evidence/p84_bounddir_residual_test.go`

## 七 不扩范围（边界）

只修 N13-UNC + W10 标注；不碰 F3 / F6 / F8（P85 范围）、不碰 F1 已闭合部分；
不改 S1 Exit 八门判据；不改 SECURITY.md 威胁模型基线；不升准则版本；
不做 gitee；不新增依赖、不扩平台。
