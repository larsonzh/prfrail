# P84 BOUNDDIR-RESIDUAL 输入清单（input-list，步骤 0 落盘）

切片：P84 `BOUNDDIR-RESIDUAL`（T4 / T5#4 残余 · 闭合 F1 / F4 · 实现类 M 档）
落盘：2026-11-21（步骤 0）；本文件随片推进追加，⑧c 前只增不删。

## 一 权威输入（逐字承接，不扩不缩）

| 来源 | 内容 |
|---|---|
| P83 报告 §8 顺延清单 P84-d | `docs/validation/security-findings-remediation.md` §8：bounddir 大小写折叠与中间 reparse 残余、根相对路径语义（T4 残余，F1） |
| P83 报告 §8 顺延清单 P84-e | 同上：proof 恢复持久化（批量回写前驱验证、代次到 proof 的可恢复事务）（T5#4 残余，F4） |
| ⑦ 追加复审（第 6 轮）T4 / T5#4 原文 | `docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/07-additional-review.md` |
| ⑦ 复审 r3 T4 / T5#4 原文 | `docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/07-re-reviews-r1-r3.md` 的 RE-REVIEW 3 节 |
| ⑦ 终审 F1 / F4 原文 | `docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/07-final-review.md` |
| P83 ① 裁决分流（T4 / T5 节为准） | `docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/01-arch/ruling-triage.md` |
| T4 被引用句 | `docs/SECURITY.md:23`（限制必须由 runner/OS 能力落实；不能仅靠提示词） |
| T5#4 被引用句 | `docs/CONTRACTS.md:91`（恢复先证明旧写者失效并取得有效 fencing token）、`docs/CONTRACTS.md:150`（promotion-receipt 绑定单个发布事务） |
| 准则相关节 | `docs/DELIVERY_DIRECTIVE.md` §1.4 / §5.2 / §5.3 / §7.5 / §7.7 / §8.2 / §9.2 |
| 队列锚点 | `docs/t027/CHECKLIST_v67.md` §一 / §四（P84 = 下一项） |

## 二 现状代码（步骤 0 读数）

| 文件 | 行数 | 说明 |
|---|---|---|
| internal/evidence/bounddir.go | 43 | 共享面 + 注释 |
| internal/evidence/bounddir_clean.go | 39 | 名字清洗（显式 `.` 被拒点） |
| internal/evidence/bounddir_unix.go | 402 | openat2 + Readdirnames 已落地 |
| internal/evidence/bounddir_windows.go | 440 | T4 主战场（EqualFold / 路径 join / 中间 reparse） |
| internal/evidence/bounddir_darwin.go | 84 | 只读 |
| internal/guard/lease.go | 509 | 租约 + 保留代 |
| internal/guard/lease_open_unix.go | 20 | |
| internal/guard/lease_open_windows.go | 63 | |
| internal/console/runtime_supervised.go | 2481 | `persistRenewedPublishLease`（:1166-1208 附近，批量回写） |
| internal/console/runtime_agent.go | 2231 | `renewPublishLease` / `verifyPublishLeasePredecessor`（:2040-2175 附近） |

## 三 T4 残余（F1，逐字要点）

1. Windows 同卷大小写敏感目录中，不同目录路径可仅大小写不同；`EqualFold`
   不能证明目标属于授权树（bounddir_windows.go `isUnderDirectory`）。
2. 后代与产物打开是 `filepath.Join(base, clean)` 后按路径 `CreateFile`；
   中间 reparse 组件仍可被跟随（仅末级带 FILE_FLAG_OPEN_REPARSE_POINT）。
3. 根语义：`BoundDirNames(dir, "")` 允许空相对路径表示根；显式 `"."`
   在 `cleanBoundRelativePath` 清洗阶段被拒，与注释不符。
4. 无产物时不得打开根（若现状仍打开须短路）。
5. 最小修复方向（复审原文）：Windows 从已绑定句柄相对打开后代，不能证明
   归属则拒绝；Ordinal 比较替换 EqualFold；中间 reparse 拒绝；根 rel 语义
   与显式 `.` 清洗一致。

## 四 T5#4 残余（F4，逐字要点）

1. `persistRenewedPublishLease`（runtime_supervised.go 约 :1166-1208）批量
   回写只比较 `ResourceID`；同资源、不同 claim 或哈希不符的其他任务 proof
   可被改成当前有效 proof，绕过直接恢复的前驱校验。
2. 代次提交与 proof 写入非同一事务；跨进程中断（提交后、proof 落盘前退出
   并重入）不可恢复。
3. 同进程分支不得把任意合形但不匹配的 token/hash 当「陈旧」追平
   （runtime_agent.go 约 :2040-2100 的 `renewPublishLease` 分支）。
4. 最小修复方向（复审原文）：让代次到 proof 的更新可持久恢复，并校验旧
   proof 属于同一 claim 的真实历史（claim / 代次 N-1 / PreviousHash /
   retained 原件）；普通异主 proof 继续拒绝，不放宽接管条件。

## 五 契约影响判定输入（① 必读）

- P83 ① 已判 T4 / T5 = DOES-NOT-TOUCH（以 SECURITY.md:23 / CONTRACTS.md:91、:150
  为被引用句）；P84 ① 须对**残余项**重新判定。
- 守则（P83 ① §五）：任何实现若需新 schema 记录、新 CLI 命令、新依赖或
  第三方包 ⇒ OUT-OF-BOUNDARY，停下升级，不自行扩域。
- 若判定触契约 ⇒ 协议先行：先改 `docs/CONTRACTS{,_EN}.md`（及 schema /
  夹具如涉）再动代码。

## 六 必读清单（对应起跑消息五）

1. `docs/t027/CHECKLIST_v67{,_EN}.md`（§四 P84 队列项）
2. `docs/validation/security-findings-remediation{,_EN}.md`（§8 P84-d / P84-e）
3. `docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/07-final-review.md` + `07-re-reviews-r1-r3.md` + `07-additional-review.md`
4. `docs/SECURITY{,_EN}.md:23`；`docs/CONTRACTS{,_EN}.md:91`、`:150`（以 P83 ① 判定书 T4 / T5#4 节为准）
5. `docs/DELIVERY_DIRECTIVE.md` §1.4 / §5.2 / §5.3 / §7.5 / §7.7 / §8.2 / §9.2
6. 现状文件：`internal/evidence/bounddir{,_unix,_windows}.go`、`internal/guard/lease.go`、`internal/guard/lease_open_windows.go`、`internal/console/runtime_supervised.go`（:1186 附近）、`internal/console/runtime_agent.go`（:2087 附近）

## 七 不扩范围（边界）

只修 T4 / T5#4；不碰 F3 / F6 / F8（P85）；不改 S1 Exit 八门判据；不改
SECURITY.md 威胁模型基线；不升准则版本；不做 gitee。
