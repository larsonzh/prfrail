# P53 事实核对清单（GLM 批辅 2/2，只读，逐字归档）

## 1. openat2(2) 内核事实
- `openat2` 于 **Linux 5.6** 引入（man-pages 6.19 `openat2.2` HISTORY："Linux 5.6"）。
- **RESOLVE_* 引入版本**：v5.6 的 `include/uapi/linux/openat2.h` 已定义全部五个标志：`RESOLVE_NO_XDEV=0x01`、`RESOLVE_NO_MAGICLINKS=0x02`、`RESOLVE_NO_SYMLINKS=0x04`、`RESOLVE_BENEATH=0x08`、`RESOLVE_IN_ROOT=0x10`（torvalds/linux@v5.6 头文件，经 jsDelivr 镜像逐行核验）⇒ **NO_MAGICLINKS 不晚于 5.6**，与 openat2 同批；`RESOLVE_CACHED` 才是后加（man-pages 标注 since Linux **5.12**）。
- **RESOLVE_BENEATH 的 EAGAIN 语义**（man-pages 6.19 ERRORS）：`how.resolve` 含 `RESOLVE_IN_ROOT|RESOLVE_BENEATH` 且内核**无法确保 `..` 组件未逃逸**（竞态或潜在攻击）→ `EAGAIN`，可重试；**已检测到逃逸** → `EXDEV`。BENEATH 拒绝绝对路径与绝对符号链接（仅允许 dirfd 后代），当前版本还隐式禁 magic-link（man 建议仍显式加 NO_MAGICLINKS）。其余错误面：NO_SYMLINKS/NO_MAGICLINKS 命中 → `ELOOP`；NO_XDEV 跨挂载点 → `EXDEV`。

## 2. 系统调用号与 struct open_how
- `SYS_openat2 = 437`（linux/amd64）：`arch/x86/entry/syscalls/syscall_64.tbl` 行 `437 common openat2 sys_openat2`（前 436 close_range、后 438 pidfd_getfd）。
- `struct open_how { __u64 flags; __u64 mode; __u64 resolve; }`：**3×8 = 24 字节（64 位）**，字段序 flags→mode→resolve；v5.6 头文件与 man-pages 6.19 `open_how(2type)` 一致；扩展字段仅可尾部追加、零值 = no-op，调用方须零填充，不支持的非零扩展 → `E2BIG`。

## 3. 支持矩阵（docs/S1_SUPPORT_MATRIX{,_EN}.md，2026-09-09，CANDIDATE）
| 范围 | 状态 |
|---|---|
| Windows 11 amd64 | 正式支持候选（纯 Go `CGO_ENABLED=0` 便携 ZIP） |
| Linux amd64 | 核心 CI；仅构建/测试/候选探针/自托管，**不宣称 S1 生产支持** |
| macOS / Windows on Arm / Linux arm64 | 不支持 |
| generic/C/Go harness | 配置与 noop 编排；executable step 仍 fail-close |
| CLI Agent/AgentRunner | 规划、未交付（T026/T027/AT-23 阻断 S1） |

- 矩阵本身**未声明 Go 工具链版本**；Go 1.22 声明散见他处：`go.mod`（`go 1.22`）、`docs/CODING_CONVENTIONS.md:36`（Go 1.22+）、`PRODUCT_REQUIREMENTS.md` REQ-015（纯 Go、Go 1.22 基线、CGO_ENABLED=0）、`TEST_STRATEGY.md:9`（模块声明 1.22，最低工具链须另验）、ADR-006（固定工具链/依赖、发布证据）。
- **与本片门禁一致性**：「仅 linux/amd64 走 openat2，其余 fail-closed（ENOSYS，无 walk 降级）」与矩阵**一致**：① 唯一 Linux 范围即 amd64（核心 CI），`linux && amd64` 构建标签精确匹配；② Win11 amd64 生产候选不经 openat2 路径，fail-closed 仅"报错拒开"，不引入矩阵外平台声明；③ arm64/macOS 矩阵本就标注不支持，不冲突。**留白**：矩阵未写 Linux 最低内核版本（openat2 需 ≥5.6；<5.6 内核返回 ENOSYS → 走 fail-closed，行为自洽，但矩阵该处无文字声明）。

## 4. CODING_CONVENTIONS.md 原文（供 ⑧ 文档回写核对）
- 第 14 行：`Markdown（.md）| UTF-8 **with BOM** | **LF** | 含 docs/*.md、copilot-instructions.md`。
- 第 17 行：`其它（.go/.py/.js/.txt/.gitignore/LICENSE…）| UTF-8 **without BOM** | **LF** | 代码与明文文本`。
- 第 30 行：`提交前自查：.md/.ps1 带 BOM+LF；.json、.go 等无 BOM+LF`。

未修改任何仓库文件（仅更新了会话记忆）。
