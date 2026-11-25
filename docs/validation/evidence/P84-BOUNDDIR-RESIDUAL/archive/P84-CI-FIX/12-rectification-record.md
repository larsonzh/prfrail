# P84 · 第二轮 Minor 整改记录（2026-11-21）

## Minor 1（卷根身份绑定种子句柄）已闭合
- 实现：`seedIsFixedVolumeRoot(handle)` 只凭**已打开种子句柄**取证：
  ① `GetVolumeInformationByHandleW` 必须对该句柄回答 `NTFS` 文件系统；
  ② 该句柄 `GetFileInformationByHandle` 的低 32 位文件索引必须等于 5
  （NTFS MFT 保留记录：卷根目录独占；高 32 位为卷特定扩展标识位，实测
  本机卷根 = High 327680 / Low 5，不参与判定）。
- 任一条件不满足 ⇒ 仅拒绝本次短名 fallback，字节快路径不触碰；无新依赖、
  无平台范围扩大、无契约改动。
- 判别证据（M8 已执行）：删除根记录比较 ⇒ `TestP84W12SeedIsFixedVolumeRoot`
  的「普通子目录句柄不得证明为卷根」断言变红；恢复后 SHA256 与变异前逐字节一致。
- 实测坑（如实登记）：初版要求 High==0，本机卷根 High=327680 导致正例红、
  种子用例 SKIP；探针实测后改为仅 Low==5（MFT 记录号即低 32 位）。

## Minor 2（短名证据门受控用例）已闭合
- 因 NTFS 语义上「同对象 + 非纯大小写分歧 + 缺失/不符短名证据」无法由内核
  构造（真实别名必然被目录项报告），引入最小生产 seam
  `directoryEntryShortNameFn`（与仓库既有 seam 先例同型），受控用例
  `TestP84W12AliasBranchRefusesMissingEntryEvidence` 注入「证据缺失」：
  真实短名夹具、非大小写分歧、候选根走链同记录前置全部真实成立，
  仅证据缺失 ⇒ 现实现拒绝。
- 判别证据（M7 已执行）：仅删除该证据门 ⇒ 同夹具放行（`suffix=["sub"] ok=true`）
  用例变红；恢复后逐字节一致。
- M3（旧 case-variant 变异）已在用例注释中标记为**历史结果**；证据门单独
  判别改由本用例承担。
- Info 1：重复长名注释限定为「匹配长名重复」；Info 2：`M2C''` 标签改正。

## 冻结快照（整改后）
- `internal/evidence/bounddir_windows.go`：+454/−2；
  `internal/evidence/p84_bounddir_residual_test.go`：+509/−0。
- 冻结 diff：`r1r2-final-diff-r3.txt`，SHA256
  `669fe92141cf8c851c483af716701677b68699cae6a2fba74d8d29c1fa6b7fdf`。
- 双工具链（go1.22.12 / go1.27.0）：W12 六顶层 + 18 子例全 PASS、零 SKIP；
  全量 `go test ./...` 双 EXIT=0；index/tree 门禁 TOTAL_FAIL=0；SELFTEST 740/740。
- 旧摘要（+398/−2、+428、71ddc7b5…）不代表整改后版本，作废。
