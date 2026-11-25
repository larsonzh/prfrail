# P84 · 方案 A 窄复审记录（2026-11-22）

- 载体：independent-reviewer + GPT-6.2 Sol（copilot）。如实登记：用户裁定署名
  「GPT-6.2 Sol Max」，当前可用模型列表中无该名，按**与前几轮复审同一载体**
  `GPT-6.2 Sol` 执行，不新增模型；超出本片 ⑦ 额度已在
  `10-final-review-request.md` 追加 OB-37 登记（四要素 + 不得作为先例 + 系用户直接裁定）。
- 基线：冻结 diff `r1r2-final-diff-r5.txt`（51,291 字节），SHA256
  `88e866260f7cd03c82e41da2fe82f05aff24f93ac6fbc77d7c54b0f91c98558c`
  （与 r4 逐字节一致，同内容重新冻结）；HEAD=origin=`d597f1e`；
  两文件暂存 +471/−2、+543/−0。只读复审：未改代码、未跑测试。
- 范围：仅方案 A（不重开 Minor 2 seam、不重开上轮 Major、不重开其它任何项）。

## 结论：RELEASE（放行）

- 总评 0 Critical / 0 Major / 0 Medium；仅 2 条 Info（不阻断）。
- 方案 A 的 48 位 NTFS 根记录判定正确：`high&0xFFFF == 0 && low == 5`
  只检查高字低 16 位（记录号第 32–47 位），忽略第 48–63 位序列号，
  正确排除记录 `0x1_0000_0005`。

## FINDINGS

| 严重度 | 位置 | 依据 | 结论 |
|---|---|---|---|
| Info | `bounddir_windows.go:1045` | 注释写 `An MFT segment reference is 48 bits wide`，把记录号位宽写成整个引用位宽；完整引用为 64 位 = 16 位序列号 + 48 位记录号。后文与实现位划分正确 | 放行（注释措辞问题） |
| Info | `bounddir_windows.go:1051` | 实测示例注释 `the record bits are all zero` 应限定为「高字中的记录号位」；High=0x50000/Low=5 解出序列号 5、记录号 5，并非记录号各位全零 | 放行（注释措辞问题） |

## 分项结论
1. **ntfsRootRecordOf 纯函数**：注释与实现一致；掩码语义正确；参数为两个数值字，
   无字节序问题。放行。
2. **seedIsFixedVolumeRoot 调用点**：确认句柄属 NTFS 且成功读取文件信息后，
   按高字、低字顺序传参。放行。
3. **常量注释 48 位语义**：NTFS 保留 MFT 记录 5 归根目录，常量与实测示例结论正确；
   措辞问题见上表 Info。放行。
4. **TestP84W12NtfsRootRecordOf（6 例）**：三组合法根记录均真、高字含记录位两例与
   低字 6 一例均假；删除高字掩码后 `record-bit-in-high-word` 确实变红——判别力成立。放行。
5. **TestP84W12SeedIsFixedVolumeRoot 两行恢复**：卷根路径构造与 driveFixed 跳过条件
   均已恢复，且在打开句柄之前执行。放行。
6. **两处 Info 注释**：扫描注释准确限定为「匹配长名重复」；1290 行 `M2C”`（U+201D）
   为 gofmt 规范形，符合用户裁定 A。放行。
7. **M9 变异证据**：静态判别力已确认（两个反例低字均为 5，删掩码后变 true 触发断言失败）；
   本次复审未运行测试、未修改文件、未独立重验恢复哈希（如实登记）。放行。

## 处置
- 无 Medium+ ⇒ 放行，不阻断。
- 2 条 Info 均为注释措辞，不阻断；**未获修复授权，未自行改动**，
  待用户裁定：随下轮整改一并处理，或保持现状。
- 仍不授权 commit / push / CI 重跑 / §11.2 回填。
