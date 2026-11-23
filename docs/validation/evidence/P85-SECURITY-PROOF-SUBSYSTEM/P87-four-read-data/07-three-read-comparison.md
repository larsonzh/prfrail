# P85 四读同输入对照表（2026-11-23 · 四模型 · 同输入 = R6 批 · 修复批 R=7）

- 输入一致：r3 三条 + R6 三组修复 diff（diff.patch 776KB）+ 契约槽位句 CONTRACTS:233/_EN:227 + 双平台读数
  （numstat / gate-index / gate-tree / gate-ci / selftest / vm-r6-legs / test-r6 / mutation-r6）。
- 独立性：四腿均只读、未改文件、互不读对方结论（Codex 为盲读，禁读前三腿报告）。
- 模型与归档：Carbon = `Gemini 4 Carbon (copilot)` → 07-carbon-main-read.md；
  Neon = `Gemini 4 Neon (copilot)` → 07-neon-cross-read.md；
  Sol = `GPT-6.2 Sol (copilot)` → 07-sol-baseline-read.md（67.9 credits / 14m50s，三读中成本最低）；
  Codex = `GPT-5.3-Codex (copilot)` Extra High → 07-codex-baseline-read.md（credits 待用户登记）。
- 用途：治理片 DIRECTIVE-REVIEW-MODEL-POLICY 选型矩阵定档数据；分歧行单列，作 P87 数据点。

## 一、三组逐组（同输入 = R6 批）

| 组 | Carbon | Neon | Sol | Codex |
|---|---|---|---|---|
| 组 1 Linux ACL 取证对证 | CLOSED | NOT-CLOSED（登记不实：BoundDir 无 AccessACL()、窗口方向写反；闭名对证/目录槽无测试变异） | NOT-CLOSED（:155 ACL 拒绝分支直接返回，未执行闭名核验） | CLOSED |
| 组 2 摘要槽位对齐契约 | CLOSED | CLOSED | NOT-CLOSED（:243 角色覆盖 + 排序兜底；拒绝摘要可少槽） | CLOSED |
| 组 3 T3 Windows 根链 | CLOSED | CLOSED（限 r3 两处替换点） | CLOSED（限 r3 两条） | CLOSED |

## 二、三项终裁（同输入 = R6 批）

| 终裁 | Carbon | Neon | Sol | Codex |
|---|---|---|---|---|
| T2 | CLOSED | NOT-CLOSED | NOT-CLOSED | CLOSED |
| T3 | CLOSED | NOT-CLOSED | NOT-CLOSED | CLOSED |
| T8a | CLOSED | CLOSED | CLOSED | CLOSED |

## 三、发现融合（修复批 R=7 输入）

| # | 发现 | Carbon | Neon | Sol | Codex | 入批 |
|---|---|---|---|---|---|---|
| 1 | Linux 目录槽 ACL 边界登记不实 + `BoundDir.AccessACL()` 缺失 + 窗口方向写反 | 未报 | **M** | Low（:100 窗口表述） | 未报（组 1 CLOSED） | ✅ |
| 2 | Windows `installSourceDescriptorDigest` 落在闭名对证之后 | Low | **M** | 未报 | 未报 | ✅ |
| 3 | Windows `verifyInstallSourceImmutability` 丢弃 leaf 句柄（ACL/DACL 未与摘要/身份对证） | 未报 | **M** | **M**（:157，同一事实） | 未报 | ✅ |
| 4 | Linux 三明治夹具实为文件槽 + 闭名对证无判别用例 + mutation-r6 缺 2 项 | Low（夹具） | **M** | 未报 | 未报 | ✅ |
| 5 | 显式 Role 与路径后缀矛盾未拒绝 / 后缀匹配区分大小写 | 未报 | Low | Low（:342 角色优先级） | 未报 | ✅（裁定含） |
| 6 | 非生产名义入口按名 Getxattr 后无二次 Lstat | 未报 | Low | 未报 | 未报 | ✅（裁定含） |
| 7 | 拒绝摘要缺槽省略（:60/:51 已占槽观察失败返回空行 → 少槽） | 未报 | 未报 | Low | 未报 | ✅（用户 2026-11-23 裁定并入：拒绝路径保持固定四槽、空槽记 ∅） |
| 8 | ACL 拒绝分支（:155）不执行闭名核验 | 未报 | 未报 | 组判据（Low 级） | 未报 | ✅（用户 2026-11-23 裁定并入：由 #1 `BoundDir.AccessACL()` fd 化自然覆盖） |
| 9 | Windows 发布 `CreateTemp`/`Link` 按路径写入残余 + 「F04 豁免」注释扩大边界 | Low | 延续边界核对 | Low（T3 终裁句） | 未报（组 3 CLOSED 但载体验证按路径写入在终裁句外未提） | ✅ 用户终裁：登记 S1 独立边界（非代码修复） |

## 四、分歧行（P87 数据点 · 单独列示）

- 组 1：Carbon/Codex CLOSED vs Neon/Sol NOT-CLOSED（分歧点 = 边界登记如实性 + 闭名核验完备性）。
- 组 2：Carbon/Neon/Codex CLOSED vs Sol NOT-CLOSED（分歧点 = 拒绝路径槽位完备性 + Role/后缀约束）。
- T2/T3 终裁：Carbon/Codex CLOSED vs Neon/Sol NOT-CLOSED。
- 级别分歧：#2 同事实 Carbon=Low vs Neon=Medium；#1 Sol 同窗口表述=Low vs Neon 整体=Medium。
- 一致点：组 3 Windows 根链四读 CLOSED；T8a 四读 CLOSED；#3 Windows leaf 丢弃 Neon/Sol 双读同为 Medium（Carbon/Codex 未报）。
- 修复批 R=7 融合口径：Neon 4 Medium + 2 Low + Sol 3 Low（#5/#7/#8）+ 用户终裁 S1 独立边界（#9）；Codex 新发现 = 0，R=7 无净新增项。
