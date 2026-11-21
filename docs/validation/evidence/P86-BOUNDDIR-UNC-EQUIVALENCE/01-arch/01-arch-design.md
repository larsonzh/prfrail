# P86 BOUNDDIR-UNC-EQUIVALENCE ① 架构方案书（主控定稿）

- 切片：P86 `BOUNDDIR-UNC-EQUIVALENCE`（`[SLICE]` · 实现类 · S 档）
- ① 载体 / 模型：`deep-reasoner` / `DeepSeek V4 Pro (deepseek)`（S 档裁剪：附录 B.1 第 6/7 项不要求）
- 落盘：2026-11-22；工作副本（⑧a 归档至 `docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/01-arch/`）
- 输入：P84 报告 §8 / §5、⑦ 收窄复审原文、⑦ 复审 r3 原文、`SECURITY.md:5`、准则 §1.4/§5.2/§5.3/§7.5/§7.7/§8.2/§9.2、现状代码

## 1 契约影响判定：DOES-NOT-TOUCH / IN-BOUNDARY

| 项 | 判定 | 理由 |
|---|---|---|
| N13-UNC（F1 残余） | **DOES-NOT-TOUCH / IN-BOUNDARY** | `docs/CONTRACTS{,_EN}.md` 全文无 bounddir / reparse / UNC 专属句；最邻近句 `CONTRACTS.md:238` 属 replay store 的写者根/保护根纪律，不约束本文件。改动不改变「Ordinal 字节比较」同源纪律，只补一处共享根等值形态 |

证据：无新 schema（`schemas/` 不动）、无新 CLI（`cmd/` 不动）、无新依赖（仅 `strings` / `filepath`，已在 import 中）。**不触契约 ⇒ 不走协议先行**，直接推进 ②。

## 2 事实基线（2026-11-22 实测，Windows / go1.27.0）

| 输入 | `filepath.Clean` | `filepath.VolumeName` | `IsAbs` |
|---|---|---|---|
| `\\host\share` | `\\host\share` | `\\host\share` | true |
| `\\host\share\` | `\\host\share\`（**保留尾分隔符**） | `\\host\share` | true |
| `\\host\share\sub` | 同形 | `\\host\share` | true |
| `\\host\share\sub\` | 去掉尾分隔符 | `\\host\share` | true |
| `D:\` | `D:\` | `D:` | true |
| `D:` | `D:.` | `D:` | **false** |

结论：UNC 共享根的**两种规范拼写**（有/无尾分隔符）都是 `Clean` 不动点，且 `VolumeName` 相同；严格后代（含尾分隔符者）被 `Clean` 归一为无尾形态，因此**只有共享根本身**存在两拼写。N13 缺口因此是**等值分支**的局部缺口，不涉及后代分支。

## 3 冲突裁定：`大小写变体`（必须留痕）

- 起跑消息 §六 明文要求：`\\server\share` / `\\server\share\` **及大小写变体**「在句柄相对走链下判等值、经生产入口放行」；§一 亦将大小写变体列为「等值被误拒」。
- P84 ⑦ 收窄复审的**最小修复**句为「仅对同一 UNC 共享根接受有、无尾分隔符两种等值形态，**保持主机名和共享名的 Ordinal 比较**」。二者直接冲突。
- **主控裁定（2026-11-22，四要素）**：裁决人 = 用户（起跑消息 §六 为同轮规范性要求）；日期 = 2026-11-22；理由 = 用户起跑消息为本片最新且更具体的要求，⑦ 最小修复系 P84 遗留建议、早于本片起跑；依据 = 起跑消息 §六 + 准则 §0.1（同轮用户指令优先）。处置 = **按起跑消息执行：尾分隔符差异 + 主机名/共享名的 ASCII 大小写差异均判等值**；⑦ 的「保持 Ordinal」作为**下限**，不构成上限。本裁定入 ⑤⑥⑦ 输入包与报告 §6，不得作为其他切片先例。
- **安全论证（须在报告与测试中可核验）**：
  1. UNC 卷命名空间由 SMB 协议定义，**主机名（DNS/NetBIOS）与共享名均大小写不敏感**，同一主机上不存在仅大小写不同的两个共享 ⇒ 大小写不敏感比较**不扩大**被放行的对象集合，只消除误拒。
  2. 折叠范围**仅限卷名（`\\host\share`）**，且**仅 ASCII**（`A`–`Z` ↔ `a`–`z`）；非 ASCII 字节必须逐字节相等。
  3. 后代组件、盘符形态根（`D:\`）、`isUnderDirectory` 一律保持 Ordinal ⇒ **W1 的 case-variant twin 拒绝不变**（`D:\Bound` vs `D:\BOUND\x` 仍拒）。
  4. 不同主机 / 不同共享名 / 中间 reparse / 非规范拼写仍拒 ⇒ fail-closed 方向不削弱。

## 4 设计（精确到 file:line）

**改动点 A —— `internal/evidence/bounddir_windows.go`**

- 插入点：`boundRelativeSuffix`（:991）内，**:1006 的 `if target == cleanBase` 之后、:1009 `trimmedBase` 之前**（严格后代分支零改动）。
- 新增助手（紧随 `boundRelativeSuffix` 之后，包级）：

```go
// uncShareRootName reports the UNC volume name (\\host\share) when p is exactly that
// name, with or without the trailing separator filepath.Clean keeps for a share root,
// and "" otherwise. Host and share bytes stay an Ordinal comparison here; the case
// fold happens only in uncShareRootEqual.
func uncShareRootName(p string) string {
	vol := filepath.VolumeName(p)
	if !strings.HasPrefix(vol, `\\`) {
		return ""
	}
	if p == vol || p == vol+string(filepath.Separator) {
		return vol
	}
	return ""
}

// uncShareRootEqual reports whether two UNC volume names name the same share root.
// ASCII letters fold both ways; every other byte must match exactly. SMB resolves
// host and share names case-insensitively, so two distinct objects can never differ
// only by ASCII case - folding here removes a false refusal and admits nothing new.
func uncShareRootEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca == cb {
			continue
		}
		if lowerASCII(ca) != lowerASCII(cb) {
			return false
		}
	}
	return true
}

func lowerASCII(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}
```

- 插入的等值分支（唯一新增放行点）：

```go
	// P86 (N13-UNC): the two canonical spellings of a UNC share root - with and
	// without the trailing separator Clean keeps - name the same root, as do
	// spellings that differ only in ASCII case of the host or share. Volume-name
	// equality is confined to this branch: descendants, the drive form and the
	// case-variant twin of a plain directory stay an Ordinal comparison.
	if baseRoot, targetRoot := uncShareRootName(cleanBase), uncShareRootName(target); baseRoot != "" && uncShareRootEqual(baseRoot, targetRoot) {
		return nil, true
	}
```

**判定顺序（修复后）**：

```
boundRelativeSuffix(base, target):
  1. nonCanonicalReparseTarget(target) != ""            -> refuse   (不变)
  2. !IsAbs(target)                                      -> refuse   (不变)
  3. cleanBase := Clean(base); target == cleanBase       -> admit    (不变)
  4. [新增] 两者皆为 UNC 共享根 且卷名 ASCII 折叠相等    -> admit
  5. trimmedBase 空                                      -> refuse   (不变)
  6. 严格后代前缀（逐字节）+ 组件守卫（"", ".", "..")     -> 不变
```

**零放宽论证（逐条）**：不同主机 ⇒ 卷名不等；不同共享名 ⇒ 卷名不等；中间 reparse ⇒ 走链机制（`reparseWalkLimit` / `reparseTargetRelativeSuffix`）不变；非规范拼写（`/`、重复分隔符、非根尾分隔符、`.`/`..`）⇒ 步骤 1 先行拒绝；界外后代（`\\server\share2\sub`）⇒ 卷名不同且步骤 6 前缀不匹配；盘符根 ⇒ `VolumeName` 无 `\\` 前缀，步骤 4 不命中。**步骤 4 的真值集合恰为「同一共享根的两类拼写差」**。

**不动的相邻面（diff 必须为 0）**：`nonCanonicalReparseTarget`（:939）、`normalizeReparseTarget`（:910）、`boundRelativeSuffixShortAlias`（:1315，UNC 明确拒绝保持不变）、`isUnderDirectory`（:1583）、`reparseTargetRelativeSuffix`（:734）。

**跨平台**：`bounddir_unix.go` / `bounddir_darwin.go` 零改动（无等值比较逻辑；CI ubuntu 腿回归）。

**改动点 B —— `internal/evidence/p84_bounddir_residual_test.go`（W10，纯注释）**

替换 :918–:924 的 Mutation check 段（代码零改动）：

```
// Mutation check: two mutations are kept apart, and each is blamed for the positive
// it kills.
//
//   - M1, the round-11-before shape alone: equality against
//     strings.TrimRight(filepath.Clean(base), `\/`) with the descendant prefix still
//     built from that trimmed value. This IS the arrangement r3 recorded: the trimmed
//     base "D:" produced the single-separator prefix "D:\", so the root itself - the
//     target equal to that prefix - was refused. Reverting equality alone turns the
//     root-itself positive red; the root-descendant positives stay green because
//     their prefix is still "D:\".
//
//   - M2, a second mutation on top of M1: also build the prefix from the UNTRIMMED
//     canonical form, which appends one more separator to the root ("D:\" becomes
//     the double-separator prefix "D:\\"). That shape is NOT the round-11-before
//     arrangement - the prefix r3 records is "D:\" - so three reds are the compound
//     of two mutations, not a faithful reversion of one. M2 with equality restored
//     turns the root-descendant positives red while root-itself stays green.
//
// The case-variant refusal of the drive form is unchanged by either mutation, which
// is why the positives are what the mutations are measured on.
```

## 5 测试点（③ 执行）

新文件 `internal/evidence/p86_bounddir_unc_test.go`（`//go:build windows`，package `evidence`）。

1. **正例（纯函数，永不 skip）** `TestP86N13UncShareRootEquivalence`：
   - `(\\server\share\, \\server\share)` → inside，suffix 空；
   - `(\\server\share, \\server\share\)` → inside，suffix 空（两个方向）；
   - `(\\SERVER\SHARE, \\server\share\)` → inside（ASCII 折叠 + 尾分隔符）；
   - 恒等 `(\\server\share\, \\server\share\)` → inside；
   - 严格后代回归 `(\\server\share\, \\server\share\sub\f)` → suffix `[sub f]`。
   - 前置断言钉住 stdlib 事实：两者 `Clean` 为不动点、`VolumeName` 相等。
2. **反例（inside=false 且 suffix==nil）**：不同主机 `\\other\share\`；不同共享名 `\\server\other`；**大小写变体的后代** `\\SERVER\SHARE\sub`（后代保持 Ordinal）；重复分隔符 `\\server\share\\`；非根尾分隔符 `\\server\share\x\`；含 `/`；`..` 组件；界外 `\\server\share2\sub`。
3. **盘符根回归**：W10/W1 既有用例全绿；`uncShareRootName(D:\)==""`、`uncShareRootName(D:\.)==""`、`(D:\)` 与 `(D:\)` 等值仍走步骤 3。
4. **生产入口（可 skip + 登记；CI windows 腿必须真跑）**：`net share` 建共享夹具 → `OpenBoundDir(\\localhost\<share>)` → 手构 mount-point buffer（复用 P84 夹具助手）以 `\??\UNC\localhost\<share>`（无尾分隔符）为 SubstituteName → `OpenBoundDirFile` 放行并读回字节一致；反例 `\??\UNC\localhost\<share>-OTHER` 拒。skip 时必须 `t.Skipf` 并写明原因；验收口径 = CI windows 腿日志出现 `ok internal/evidence`。
5. **变异检验（两步机械）**：① 单点删去步骤 4 分支（锚定注释 + if 行），`go test ./internal/evidence -run 'TestP86' -count=1` 必须**恰好**正例变红、反例与 W10/W1 保持绿；② 恢复后 `Get-FileHash` 与基线**逐字节一致**，全量 `go test ./... -count=1` 15 包 ok。

## 6 落点与门禁影响

新增独立测试文件（P84 测试文件只作 W10 注释更正的唯一例外）；`tools/gates/**` 零改动 ⇒ **SELFTEST 保持 740/740**；`.go` 不在三作用域文档扫描面，但 ④ 仍须三作用域 `TOTAL_FAIL=0`。P84 证据包 `SHA256SUMS` 不因本片改动而失效（本片不改 P84 证据载荷，仅改仓库测试源码），⑧a 时机械核对并在报告登记。

## 7 盲审判定

**必抽 1 次盲审（另计 +1，⑦ 终审前执行）**。理由：bounddir 等值判定即「目标是否与绑定对象同一」的身份/所有权判定；本片新增一个放行集合（尾分隔符 + ASCII 大小写），落入 §7.7 第一类字面范围。判定理由随本方案书落盘。

## 8 可证伪闭合断言

- **A1**：步骤 4 的放行集合恰为「同一 UNC 共享根的尾分隔符差 + 卷名 ASCII 大小写差」；反例表逐项为 false。
- **A2**：`nonCanonicalReparseTarget` / `boundRelativeSuffixShortAlias` / `isUnderDirectory` / 严格后代分支 diff = 0；总改动 = 步骤 4 分支 + 三个助手 + W10 注释段 + 新测试文件。
- **A3**：撤除步骤 4 ⇒ 正例红、恢复后逐字节一致；全量测试 / 三作用域 / SELFTEST 740/740 恢复基线。
- **A4**：盘符根与 W1 case-variant twin 回归不红（W10 五用例 + W1 全绿）。

## 9 最可能被打回的三点（给 ⑤⑥⑦ 的预案）

1. **大小写变体裁定**：已按 §3 留痕；若 ⑦ 判其超边界，整改方向 = 收窄为「仅 ASCII 折叠 + 仅卷名」并补反例（不接受全量 `EqualFold`）。
2. **生产入口可达性**：`net share` 失败时 skip 可能被读作证据缺口 ⇒ 必须留 `ok internal/evidence` 级日志与 skip 原文。
3. **分支位置**：若实现把判断放到守卫之前或改 `nonCanonical`，将放宽拒绝闭集 ⇒ 严格落在步骤 3 之后、步骤 5 之前。

END OF DESIGN

## 10 Rev 2（⑤ r1 整改，2026-11-22）

- ⑤ r1 = `PASS WITH FIXES`（1 Medium + 4 Low）；逐条处置如下。
- **Medium（生产入口验收判据）**：⑨ 的 windows 证据改为 `go test -v` 输出含 `--- PASS: TestP86UncShareRootProductionEntry`；若 SKIP 则在报告逐字登记 skip 原因（本地非提权已实测 SKIP：`net share ... exit status 2 (发生系统错误 5。拒绝访问。)`）。
- **Low（`boundRelativeSuffix` 注释）**：已补 P86 唯一例外段（UNC 共享根卷名 ASCII 折叠；后代 / 盘符形态 / 平面目录大小写孪生仍 Ordinal）。
- **Low（`uncShareRootName` 扩展与设备命名空间）**：已收紧——卷名以 `\\?\` 或 `\\.\` 起头即返 `""`；分类器负例新增 `\\?\UNC\server\share\`、`\\?\UNC\server\share`、`\\?\C:\`、`\\.\C:\`。
- **Low（文档同步清单）**：补入 `SECURITY{,_EN}.md:5` 为必改（起跑消息 §七-5）；`bounddir-residual{,_EN}.md` §8 的时点注记属起跑消息 §七-7「核对项（由用户裁定）」⇒ 本片不擅自加，⑩ 回执列入待裁定。
- **Low（变异判据措辞）**：对齐测试文件注释——恰好 3 个拼写差异子例（`trailing-base-bare-target` / `bare-base-trailing-target` / `ascii-case-folded-trailing-target`）变红；`identical-trailing-spellings` 与 `strict-descendant-regression` 两世界同绿。
