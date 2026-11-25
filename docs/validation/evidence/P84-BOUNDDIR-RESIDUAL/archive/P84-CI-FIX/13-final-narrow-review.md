# P84 · 最终收窄复审记录（GPT-6.2 Sol Max · 超额度，OB-37 已登记）

- 日期：2026-11-21。载体：independent-reviewer + GPT-6.2 Sol Max。
- 基线：冻结 diff `r1r2-final-diff-r3.txt`（SHA256 `669fe92141cf…`）；HEAD=origin=`d597f1e`；
  两文件暂存 +454/−2、+509/−0。只读复审，未改代码、未跑测试。

## 结论：FINDINGS（阻断）

| 严重度 | 位置 | 原因 | 结论 |
|---|---|---|---|
| **Major（§4.3 Medium）** | `bounddir_windows.go:1090`（seedIsFixedVolumeRoot） | NTFS MFT **segment reference 为 48 位**（Microsoft 定义）：`FileIndexHigh` 的低 16 位仍属记录号，仅最高 16 位是序列号。只验 `FileIndexLow == 5` 无法区分记录 `5` 与 `0x1_0000_0005`；本机根句柄读数（High=0x50000）不能证明记录号高位恒零 ⇒ 非根种子未被充分排除 | **阻断** |
| Info（Low） | `bounddir_windows.go:1169` | 函数总注释的 `refuse` 清单仍写 `a duplicate long name`，比实现（仅 MATCHED 名重复）更宽 | 放行 |
| Info（Low） | `p84_bounddir_residual_test.go:1290` | 冻结暂存版本仍为 `M2C”`（U+201D），与记录 `M2C''` 不一致 | 放行 |

## 分项结论
- **Item 1（原 Minor 1）：达 Medium+ ⇒ 阻断。** 认可面：证据绑定已打开种子句柄；
  卷类型须答 NTFS；种子以 FILE_FLAG_OPEN_REPARSE_POINT 打开并先查重解析属性；失败仅拒
  短名回退、不动字节快路径；非 NTFS 卷拒绝、普通网络映射被 DRIVE_FIXED 预筛。
  阻断面：SUBST 指向子目录等场景仍依赖有缺口的根记录比较；M8 判别只覆盖「普通子目录」，
  未覆盖「记录号高位不同、低位相同」。
- **Item 2（原 Minor 2）：未达 Medium+ ⇒ 放行。** seam 默认绑定真实函数、仅受控测试写入、
  `t.Cleanup` 恢复、无并行写入者；夹具真实别名 + 非大小写分歧 + 真实逐级句柄走链，仅注入
  NTFS 无法自然提供的「证据缺失」；M7 判别成立（记录未附原始输出，仅文字）。
  `M3 HISTORICAL` 标注准确。
- 两处 Info：均不阻断；其中 `M2C”` 标签与函数总注释在本轮暂存版仍未更正（如实登记，
  随下轮整改一并处理）。

## 处置（按裁定）
- 出现 Medium+ ⇒ **停止、报用户；不自行修复、不自行加轮、不自行回退**。
- 恢复主线备选方案与代价见下轮回报（主控拟议：48 位记录号判定）。
