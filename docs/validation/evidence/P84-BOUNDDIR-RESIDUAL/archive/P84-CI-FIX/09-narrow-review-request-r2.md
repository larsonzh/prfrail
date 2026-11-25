# P84 · 第二轮收窄复审请求记录（§7.5 / §7.8 留痕）

- 日期：2026-11-21。裁定人：用户（直接裁定）。载体：GPT-6.2 Sol Max（independent-reviewer）。
- **额度说明**：本片 ⑦ 既定额度已用尽；本轮回合为第二次收窄复审，同样**超出既有额度**，
  系用户本次直接裁定，**不得作为先例**。
- **范围（唯一输入）**：
  1. 本轮修复后的最终暂存 diff：`docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1r2-final-diff-r2.txt`
     （SHA256 = `71ddc7b5d3af8444dd9398b6e156a97e7fed944463702c669da69f5ae7fa690b`；
     暂存统计：实现 +398/−2、测试 +428/−0）。
  2. 上轮 FINDINGS 原文与本轮处置：`docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/08-narrow-re-review.md`、本轮 `10-…`（本文件）。
  3. 代码：`internal/evidence/bounddir_windows.go` 的
     `scanFileIDBothDirInfoForShortName`（本轮改动：整批校验后出证、非末条以
     NextEntryOffset 为上界、空缓冲先拒、同名 FileId 冲突立即拒绝、剥终止符为空不出证、
     每条的短名范围均校验）、`boundRelativeSuffixShortAlias`（仅大小写不同拼写不进入
     放行、卷根种子不跟随 reparse + GetDriveTypeW == DRIVE_FIXED）。
  4. 测试：`internal/evidence/p84_bounddir_residual_test.go` 的 W12 组（含新增判别用例
     跨记录 / 未对齐但入缓冲 / FileId 仅高 32 位不同，及各自的变异判别记录）。
- **不审**：P84 全片、N13-UNC、P85、权威契约、CI 配置。仅核本轮 FINDINGS 的处置与
  最终暂存差异。
- 输出：Critical/Major/Minor/Info + file:line + 理由；结论 SIGN-OFF 或 FINDINGS。
  记录归档 `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/11-narrow-re-review-r2.md`。
