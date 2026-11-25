# P84 CI 首跑红处置 · R-2 契约影响判定（① 书面判定）

- 日期：2026-11-21。性质：① 协议先行判定书（临时件，⑧b 回填时并入证据包）。
- 对象：`internal/evidence/bounddir_windows.go` 的 `openRelativeComponent` 对**文件叶**打开追加
  `FILE_SYNCHRONOUS_IO_NONALERT`（CreateOptions 0x20），并更正两处注释。

## 判定：DOES-NOT-TOUCH / IN-BOUNDARY

逐条对照契约面：

| 契约面 | 依据 | 结论 |
|---|---|---|
| 读取原语（一次打开、单句柄、字节/截断/摘要同对象） | `boundread.go` BoundFile 契约；`CONTRACTS.md` §7 证据语义 | 不变：仍一次打开、同一句柄；仅文件对象 I/O 完成模式变化，不可观察于语义 |
| P52 中间 junction 界内放行 / 界外拒绝（含 "outside the bound directory" 文本） | ARCHITECTURE.md §5.1；P84 01-arch ②(c)/④ 口径 | 不变：`reparseTargetRelativeSuffix`、N10/N13、Ordinal、卷号 belt 全部未动 |
| 末级 reparse 一律拒绝、枚举目标 reparse 拒绝 | `p83_fixround8_windows_test.go`；P84 口径 | 不变：分类与拒绝条件逐字节未动 |
| 目录枚举（readDirNames，类 10/11） | P84 T4 fix round 2/6 | 不变：目录组件不加该标志，仍异步对象 + GetFileInformationByHandleEx |
| 句柄走链与替换拒绝 | `confirmOpenedArtifact` / `verifyBoundObjectStillNamed` | 不变 |
| 平台/工具链契约 | `CONTRACTS.md:224`（windows/amd64 交付范围）；CI 固定 go1.22.12 | 本改动**正是**恢复该契约在仓库 CI 工具链上的可读性（修前 87 失败 = 契约违约，修后绿） |

- 无 schema / CLI / 依赖 / wire 记录变更；不改任何拒绝条件、不放宽任何准入。
- 修正的注释原声称「加 FILE_SYNCHRONOUS_IO_* 会使 NtCreateFile 返回 STATUS_INVALID_PARAMETER」，
  被本地探针实测证伪（该旧结论缺 SYNCHRONIZE 访问位所致；本片叶开恒带 SYNCHRONIZE）。
- 结论：不触 CONTRACTS{,_EN}.md / schema / 夹具，可直接修代码。
