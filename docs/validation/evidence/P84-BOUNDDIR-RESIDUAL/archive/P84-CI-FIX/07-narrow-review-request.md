# P84 · R-1/R-2 收窄复审请求记录（§7.5 / §7.8 留痕）

- 日期：2026-11-21。裁定人：用户（直接裁定）。
- 载体：GPT-6.2 Sol Max（独立载体，原 ⑦ 载体；独立-reviewer 角色）。
- **额度说明**：本片 ⑦ 既定额度（盲审 1 + 终审 1 + 复审 3 + 窄域例外 1）此前已全部用尽；
  本次 1 次收窄复审**超出既有额度**，系用户本次直接裁定，**不得作为先例**。
- **范围（唯一输入）**：
  1. 最终 R-1/R-2 暂存 diff：`docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1r2-final-diff.txt`（git diff --cached 全文）；
  2. 产品代码：`internal/evidence/bounddir_windows.go`（重点：
     `openRelativeComponent` 文件叶 FILE_SYNCHRONOUS_IO_NONALERT（R-2）；
     `reparseTargetRelativeSuffix` fallback 接入；`boundRelativeSuffixShortAlias`
     候选根受限走链与三重证明；`directoryEntryShortName` /
     `scanFileIDBothDirInfoForShortName` 目录项解析加固）；
  3. 测试：`internal/evidence/p84_bounddir_residual_test.go`（W12 正例 +
     三组拒绝 + 目录项扫描器 13 子例）；
  4. 证据：`docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/`（01 契约判定、02 验证矩阵、03 R-1 设计含 §8/§9/§10、
     04 复跑记录、双工具链全量输出、探针输出、锚点与归一化清单）。
- **不重审**：P84 全片、N13-UNC、P85、T5#4、租约链——仅限上述差异。
- **判据**：拒绝边界是否放宽（N10/W1 大小写/UNC/末级与枚举 reparse）；目录项解析
  是否可被畸形缓冲误导（NextEntryOffset/UTF-16 长度/嵌入 NUL）；短名放行的第二条
  证据是否与已打开对象身份绑定；观察时点与「最终界内读 ≠ 原目标已证明」表述。
- 输出要求：按 Critical / Major / Minor / Info 分级，每条附 file:line 与理由；
  结论：SIGN-OFF 或 FINDINGS。复审记录归档 `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/08-narrow-re-review.md`。
