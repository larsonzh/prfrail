# P84 · 方案 A 整改记录（2026-11-21）

用户裁定（同轮授权，严格范围）：48 位记录号修复（纯函数 `ntfsRootRecordOf`
+ `seedIsFixedVolumeRoot` 调用点）+ 判别用例 + 两处 Info + 双工具链验证 +
冻结 diff r4 + 本记录落盘；其余一律不动。**未 commit / 未 push / 未 CI 重跑 /
未做 §11.2 回填 / 未自行发起复审**（窄复审载体与范围由用户另行裁定）。

## 方案 A 落地

### 产品侧（`internal/evidence/bounddir_windows.go`）
- 新增纯函数 `ntfsRootRecordOf(high, low uint32) bool`：
  `high&0xFFFF == 0 && low == ntfsRootDirectoryFileIndex`。
  依据：MFT segment reference = 高 16 位序列号（High 的高 16 位）
  + 48 位记录号（High 的低 16 位 + Low 32 位）⇒ 只判 Low==5 会放行
  记录 `0x1_0000_0005`；必须同时要求 High 的低 16 位为 0。
- `seedIsFixedVolumeRoot` 末行改为 `return ntfsRootRecordOf(info.FileIndexHigh, info.FileIndexLow)`。
- `ntfsRootDirectoryFileIndex` 常量注释重写为 48 位语义（本机卷根实测
  High=0x50000 / Low=5：高 16 位是序列号位，低 16 位为 0）。
- `ntfsRootRecordOf` 带完整 doc comment（说明只判 Low 的放行后果）。
- Info 1：scan 函数 doc bullet 改为「a duplicate OF THE MATCHED long name」。

### 测试侧（`internal/evidence/p84_bounddir_residual_test.go`）
- 恢复 `TestP84W12SeedIsFixedVolumeRoot` 开头被误删的两行
  （volumeRoot 构造 + `driveTypeOf != driveFixed` 的 skip 门）。
- 新增纯单元 `TestP84W12NtfsRootRecordOf`（6 例，全 PASS）：
  root-record（收）、measured-development-volume-root High=0x50000（收）、
  sequence-bit-only High=0x00010000（收：高 16 位序列号不参与）、
  record-bit-in-high-word High=0x00000001（拒：记录 0x1_0000_0005）、
  sequence-and-record-bits High=0x00010001（拒）、other-low-word（拒）。
  （勘误：原记 5 例，漏计 other-low-word；实际为 6 例，2026-11-21 用户裁定补正。）
- Info 2：`M2C''` 标签改正（原智能引号）。

## M9 判别变异（已执行）
- 变异：`ntfsRootRecordOf` 临时改为 `return low == ntfsRootDirectoryFileIndex`
  （删除 High 掩码）。
- 结果：`record-bit-in-high-word` 与 `sequence-and-record-bits` 两子例变红，
  其余三例仍绿 ⇒ 判别用例能区分「缺 High 掩码」这一单项判据。
- 恢复后 `bounddir_windows.go` SHA256 与变异前逐字节一致
  （`f2ea702d1e37e55794bbc39cff518b0f6a47cda9b8dced329104a499a29fa4c6`）。

## 双工具链验证（全为 r4 状态实跑）
- W12（go1.22.12）：EXIT=0，8 个顶层全 PASS、零 SKIP（证据 `w12-go122-r4.txt`）。
- W12（go1.27.0）：EXIT=0，8 个顶层全 PASS、零 SKIP（证据 `w12-go127-r4.txt`）。
- 全量 `go test ./...`：go1.22.12 EXIT=0（`full-go122-r4.txt`）、
  go1.27.0 EXIT=0（`full-go127-r4.txt`），均无 FAIL 行。
- `go build ./...`、`go vet ./...`、`gofmt -l` 均干净。
- 门禁：index 与 tree 两个 scope 分别 `TOTAL_FAIL=0`、EXIT=0
  （`gates-index-r4.txt` / `gates-tree-r4.txt`）；`selftest.js` PASS 740/740。

## 冻结快照（r4，已被 r5 取代，作废）
- 暂存：`internal/evidence/bounddir_windows.go` +471/−2、
  `internal/evidence/p84_bounddir_residual_test.go` +543/−0
  （合计 2 files, 1014 insertions(+), 2 deletions(-)）。
  （勘误：原记 +473/−2，git numstat 实测为 471 增 2 删；2026-11-21 用户裁定补正。）
- 冻结 diff：`r1r2-final-diff-r4.txt`（51,291 字节），SHA256
  `88e866260f7cd03c82e41da2fe82f05aff24f93ac6fbc77d7c54b0f91c98558c`。
- r3 快照（`669fe921…`）已被 r4 取代，作废；r4 随后被 r5 取代（见下方 r5 段）。
- Info 2 标签说明（2026-11-22 用户裁定 A 定稿）：直引号写法 `M2C''` 与
  doc comment 内的 gofmt 规范形 `M2C”`（U+201D）**指同一标签**；
  规范形是 Go 1.19+ doc comment 排版规则对 `''` 的唯一合法形态，最终保留
  `M2C”`（gofmt 双工具链干净）。

## 如实登记
- 执行中发现 `tmp/r1probe/main.go`（前轮探针遗留，gitignore 目录内）被
  `go test ./...` 扫入，已按仓库「临时文件用完即清」规则删除，随后重跑
  双工具链全量以取得干净证据集。
- 首次门禁调用误将 `--scope=index|tree` 传给 PowerShell（`|` 被解析为管道，
  产出两个无效证据文件），已删除无效文件并用 index/tree 分跑取得有效读数。
- 上表双工具链读数以「清理后的重跑」为准。

## r5 冻结（2026-11-22，用户裁定 A）
- 裁定：撤销直引号字面要求，保留 gofmt 规范形 `M2C”`（U+201D）；
  仅恢复该行并精确暂存；两种写法指同一标签；行数与例数勘误保留。
- 执行：1290 行恢复为 `M2C”`（一次性脚本，锚点唯一=1，U+201D 归位、
  全文件 `''` 归零，脚本已清除）；`gofmt -l` 双工具链（go1.22.12 /
  go1.27.0）均无输出；`git diff --cached --check` EXIT=0；两文件精确暂存、
  工作树无未暂存改动。
- 冻结 diff：`r1r2-final-diff-r5.txt`（51,291 字节），SHA256
  `88e866260f7cd03c82e41da2fe82f05aff24f93ac6fbc77d7c54b0f91c98558c`
  = **与 r4 逐字节一致 ⇒ 如实登记为同内容重新冻结，沿用 r4 测试证据**
  （W12 双工具链全 PASS 零 SKIP、全量双 EXIT=0、门禁 TOTAL_FAIL=0、
  SELFTEST 740/740 均不重跑）。
- 候选 `c6429341a6cb…`（含 `M2C''` 的暂存态）按裁定**不作为基线**。

