# P84 · 第二轮收窄复审记录（GPT-6.2 Sol Max · 超出额度例外）

- 日期：2026-11-21。载体：GPT-6.2 Sol Max（independent reviewer）。
- 授权：用户直接裁定追加，**超出既有 ⑦ 额度，不得作为先例**（见
  `09-narrow-review-request-r2.md`）。范围仅本轮 FINDINGS 处置与冻结暂存差异
  （diff SHA256 `71ddc7b5…`；实现 +398/−2、测试 +428/−0）。

## 结论：FINDINGS（0 Critical / 0 Major / 2 Minor / 2 Info）

**上轮 Major 已确认修复**（跨记录取证：非末条以 NextEntryOffset 为上界 + 整批校验后
出证；判别用例与变异复核成立）。三类判别用例的单项变异均只使对应子例变红
（cross-record ↔ 记录上界；unaligned-but-in-buffer ↔ 8 字节对齐；
file-id-high-half ↔ 64 位比较）。

| 级别 | 位置 | 原因 | 建议最小修复 |
|---|---|---|---|
| Minor | `bounddir_windows.go:1288` | `DRIVE_FIXED` 门槛与 `FILE_FLAG_OPEN_REPARSE_POINT` 已落实，但前者按路径查询，不能证明随后打开的**句柄**就是实际卷根；后续句柄身份校验仍有效，未发现越界读取 | 对已打开的种子句柄验证「固定本地卷 + 卷根身份」，无法证明即拒绝 |
| Minor | `p84_bounddir_residual_test.go:1333` | 大小写变体现被 `EqualFold` 提前拒绝，仅移除同目录项短名证据不再使该用例变红；注释中的 M3 已不能证明最终版覆盖该拒绝路径 | 将 M3 标为历史结果；用不会被大小写门槛提前拒绝的受控夹具验证「短名证据缺失」 |
| Info | `bounddir_windows.go:1187` | 注释称重复长名拒绝整批，实现只检查**匹配长名**是否重复 | 注释限定为匹配长名重复 |
| Info | `p84_bounddir_residual_test.go:1290` | 变异标签写作 `M2C”`，与记录 `M2C''` 不一致 | 使用记录中的准确标签 |

- 未放宽面确认：只有遍历到链尾并核对完整 64 位 FileId 后才可能出证；零长名称、指向
  缓冲末尾的偏移、重复匹配、空缓冲、任一条目短名范围畸形均拒绝；身份冲突
  `refuse=true` 使批次循环硬停；`EqualFold` 门槛在取证之前且不误拒真实短名。
- 按用户裁定第 4 条：**复审后停下回报，不自行加轮、不执行回退**；本轮不修复上述
  Minor/Info（未获授权），不留阻断项。
