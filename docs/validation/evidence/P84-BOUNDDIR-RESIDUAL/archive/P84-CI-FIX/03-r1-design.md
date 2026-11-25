# P84 CI 首跑红处置 · R-1 设计（本轮仅设计+验证，不改产品代码）

## 1. 观察事实：8.3 差异仅限绑定根前缀（CI 日志字节级证据）

run 44877224368 attempt 1/2 windows 腿失败信息（原文）：

- `intermediate reparse point j resolves to C:\Users\RUNNER~1\AppData\Local\Temp\Test…\001\D\sub, which is outside the bound directory C:\Users\runneradmin\AppData\Local\Temp\Test…\001\D`
- （W5b 同形：`…\D\hop` vs `…\D`；W9 同形：`…\001\Bound\sub` vs `…\001\Bound`）

字节级对比：
| 区域 | 存储目标 T | 句柄名 base |
|---|---|---|
| 卷 | `C:` | `C:`（同卷） |
| 差异组件 | `Users\RUNNER~1` | `Users\runneradmin` |
| 共同后缀 | `\AppData\Local\Temp\Test…\001\D\sub` | `\AppData\Local\Temp\Test…\001\D`（逐字节相同） |

⇒ 观测到的 8.3 差异**仅存在于绑定根前缀段**（`runneradmin` 目录的短名别名 `RUNNER~1`），
目标在绑定根之下的后缀逐字节一致。本机对照：`cmd /c dir /x C:\` 显示本机 `Users` 无短名
（CJK 用户目录未生成 8.3 别名），故本机字节比较恒等、CI 恒红——与本机/CI 表现完全吻合。
（`fsutil 8dot3name query` 需提权被拒，未采信；`dir /x` 足证。）

## 2. 探针实测（GetLongPathNameW / GetShortPathNameW，本机 C: 卷）

- `C:\PROGRA~1` → long=`C:\Program Files`；`C:\PROGRA~1\Windows Defender` →
  `C:\Program Files\Windows Defender`（真实 8.3 别名可转换，转换组件取磁盘长名）。
- `c:\progra~1\windows defender` → `c:\Program Files\windows defender`：**非转换组件
  保持输入大小写**——GetLongPathName 只重写短名组件，不做大小写归一。
- 不存在的别名（`C:\NO-SUCH-ALIAS~1`）→ ERROR_FILE_NOT_FOUND（fail 信号可得）。
- `C:\Users` / `C:\Windows\System32` 无短名（`dir /x` 空列）⇒ 本机无法构造 8.3 正例；
  CI runner（`RUNNER~1` 别名存在）可构造 ⇒ 短名正例测试必须含主机边界 skip（W1 先例）。

## 3. 收窄设计：候选根组件切分 + 使用时对象同一性（按 2026-11-21 第二次裁定修订）

**只处理本次 CI 已证实形态**：存储目标 T 的前缀段以 8.3 短名拼写**绑定根自身路径**，
绑定根组件边界以下（TComps[n..]）逐字节未改写。CI 锚点四例（`RUNNER~1` vs
`runneradmin`）与本机自造别名复现（`R1PRO~1` vs `r1probe-*`）同形。

**候选根切分（组件边界）**：N10 通过后，若 T 与 base 字节比较失败：
1. `len(TComps) > n`（n = base 组件数，base = 绑定根句柄名）且 `TComps[0] == baseComps[0]`
   （同卷盘符），否则拒绝。
2. 候选根 = `TComps[0..n-1]`（T 按**绑定根的组件边界**切出，即短名拼写的绑定根路径本身）；
   待走后缀 = `TComps[n..]`（未改写，Ordinal 逐字节使用）。

**候选根本身的受限打开与同一性验证（不以祖先一致替代根一致）**：
1. 从卷根起逐组件 `CreateFileW`（`FILE_FLAG_BACKUP_SEMANTICS |
   FILE_FLAG_OPEN_REPARSE_POINT`，share-all），每层 `GetFileInformationByHandle`
   检查 `FILE_ATTRIBUTE_REPARSE_POINT`：**任一层是 reparse 点 ⇒ 拒绝**
   （未经验证的重定向不得参与候选根解析）。
2. 最终层 `(VolumeSerialNumber, FileIndexHigh, FileIndexLow)` 与**绑定根句柄自身**的
   `GetFileInformationByHandle` 同值比较——验证对象是候选根（绑定根组件边界的
   TComps[0..n-1] 解析结果）**本身**，不是其任一祖先；祖先一致性不作为放行依据。
3. 任何失败（打不开、任一层 reparse、卷号或文件 ID 不符、观察互相矛盾）⇒ 维持现有
   "outside the bound directory" 拒绝。

**验证通过后**：从**原绑定句柄**走未改写后缀 `TComps[n..]`（现有走链 +
`confirmOpenedArtifact` 逐组件重证），候选根路径在行走中不再被使用。

**时点与证明边界（不把最终读到界内等同于原目标已获证明）**：
- 判定锚点 = 候选根最终层 (卷号, 文件 ID) 与绑定根句柄同值比较的**那一观察时点**。
  该时点之后，行走只依赖绑定句柄与未改写后缀，候选根路径的任何变化（含前缀别名
  被重指）不影响行走结果（最坏 fail-closed：界内收窄或拒绝）。
- 原目标「已获证明」的构成 = N10 拼写完整性 + 候选根对象同一性 + 未改写尾段边界；
  行走成功读回界内字节是**执行结果**，不是证明本身。二者在报告与代码注释中
  必须分开陈述。

## 4. 改名与竞态的拒绝条件（沿用既有机器 + 上述受限打开）

| 事件 | 拒绝条件 |
|---|---|
| 候选根任一层是 reparse 点 | 受限打开的逐层属性检查拒绝（新增，唯一放行前置） |
| 候选根最终对象 ≠ 绑定根对象（卷号/文件 ID 任一不符） | 同一性比较拒绝 |
| 绑定根在判定后被改名 | `verifyBoundObjectStillNamed` / `confirmOpenedArtifact` 既有拒绝 |
| 后缀组件在判定后被替换为 reparse | 行走逐组件重新分类拒绝（既有） |
| 前缀别名在判定与行走之间被重指 | 行走不经过候选根解析 ⇒ fail-closed（界内收窄或拒绝） |

## 5. 设计（供裁定；未实现）

`reparseTargetRelativeSuffix` 的 N10 检查之后、`boundRelativeSuffix` 之前插入：
1. N10 字节规范检查（不变）。
2. 存储字节与 base 的 Ordinal 比较（不变快路径）——现行为。
3. 仅当 2 失败且 T 满足 §3 形态：候选根受限打开 + 卷号/文件 ID 同值验证；通过 ⇒
   从绑定句柄走未改写后缀；任一前置失败 ⇒ 现有边界消息拒绝。
4. 大小写纪律不变（尾段逐字节匹配）；W9 负例、excursion、末级/枚举 reparse 拒绝原样。

## 6. 夹具验证证据（本机双工具链红绿对，2026-11-21 实测）

隔离夹具（临时探针，仓库外模块、用完即清；输出归档 `r1-probe-output*.txt`）：
1. `root = MkdirTemp` → `D\sub\f`（内容 `inside-bytes`）；`fsutil file setshortname
   root R1PRO~1`（**按文件设短名，不改系统/卷级 8.3 策略**）；`GetShortPathNameW`
   得短名根；`mklink /J D\j → <短名根>\D\sub`。
2. **红（撤掉放行分支 = 现产品）**：`OpenBoundDirFile(dir, "j\f")` 拒绝，错误与本机/CI
   同形（`resolves to C:\...\R1PRO~1\D\sub, which is outside the bound directory
   C:\...\r1probe-…\D`）。
3. **绿（启用放行分支）**：候选根 `R1PRO~1\D` 逐层受限打开（无 reparse）→
   `candidateID == boundID`（vol=1825482194、indexHigh/ Low 相同）→ 从绑定句柄
   走 `sub\f` → **读回 `inside-bytes`**。
4. 工具链：go1.22.12 与 go1.27.0 各跑一遍，红绿对均成立（输出两件）。
5. **反证纪律**：`setshortname` 失败 ⇒ SKIP 属主机边界，**不作为任何反证/通过**；
   红绿对均在别名真实生效的同一夹具上取得；不得以长名夹具替代短名正例
   （变异测试将钉住）。

## 7. 契约影响判定（按裁定重判：P52 行为要求 + SECURITY 归属约束为依据）

- **P52 行为要求**（ARCHITECTURE.md §5.1 反例矩阵 + P84 01-arch ②(c)）：
  「中间 junction 目标是界内 ⇒ 必须放行；界外 ⇒ 拒绝且错误含 outside the bound
  directory」。当前实现把 8.3 短名拼写的界内目标误拒 ⇒ 在含 8.3 别名的主机上
  违反该行为要求（CI 实证 + 本机复现）。本设计是**修复该违约**的最小口径。
- **SECURITY 归属约束**（`docs/SECURITY{,_EN}.md:23` / `:31`）：绑定读取须有
  OS 级归属证明。本设计以 `(卷号, 文件 ID)` 对象同一性 + 逐层 reparse 拒绝强化
  归属证明，方向一致，无放行窗口（§3 时点分析）。
- **`CONTRACTS.md:238`**（replay store 路径安全条款，含「Windows 上同一性判定覆盖
  大小写与 8.3 短名别名」）**仅作参考**：该句管 replay store 路径比较，不直接管
  bounddir 存储目标比较；本设计不改动该条款。
- **确需改动的权威契约句**（报拟改原句，**本轮不落盘**）：
  1. P84 01-arch N10 句（现文）：「……目标须字节等于其 filepath.Clean 形态……再以
     **原字节**做 Ordinal 等于 / 严格后代判定。」
     拟改：「……再以原字节做 Ordinal 等于 / 严格后代判定；仅当字节判定失败且目标
     前缀段以 8.3 短名拼写绑定根路径时，允许以候选根（按绑定根组件边界切分）的
     受限打开对象同一性（同卷、同文件 ID，且未穿越任何 reparse 点）替代该前缀段的
     字节相等，尾段仍逐字节未改写，并由绑定句柄走未改写后缀；任一证明失败即拒绝。」
  2. `bounddir_windows.go` `boundRelativeSuffix` N13 段注释同步修订。
  3. ARCHITECTURE.md §5.1 不需改文（P52 行为要求本就要求界内放行）。
- 结论：本设计触及 P84 01-arch 契约句修订 ⇒ **协议先行**；是否授权改句由你裁定，
  本轮不自行落盘。

---

## 8. 实施修订（2026-11-21 用户裁定：有条件最小实施，不改权威契约，追加本修订）

按裁定「依据 P52 行为要求修复，追加方案书修订，不追改历史审查口径；不改权威契约」，
§7 的契约句修订**不执行**；本修订作为方案书追加节，生效设计为：

1. **候选根验证改为句柄相对走链**：从**可验证的本地卷根目录句柄**起（`CreateFileW`
   打开卷根一次作为种子），此后**每一步只相对上一已验证的父句柄打开一个组件**
   （复用 `openRelativeComponent`，`FILE_OPEN_REPARSE_POINT`），逐层检查 reparse
   属性；**遇 reparse 或任何一步不可验证即拒绝**。**不采用逐层拼绝对路径再
   CreateFileW 重开的方案**（本文件早前 §6 探针的 openLevel 实现方式作废，探针
   结论仍作为可行性证据，不作生产入口验收）。
2. **新分支只接纳经 OS 证实的 8.3 根前缀别名**：分歧组件（`TComps[k] !=
   baseComps[k]`）须同时满足（i）在同一已验证父句柄下按两种拼写打开得到同一文件
   记录（文件 ID 相同）；（ii）`GetShortPathNameW(长拼写前缀)` 的第 k 组件逐字节
   等于 `TComps[k]`——**同一文件 ID 本身不证明差异只是 8.3**（大小写变体也同记录，
   但短名返回仍是长拼写 ⇒ 拒绝）。最终层 (卷号, 文件 ID) 与绑定根句柄自身同值。
3. **边界保留**：N10 不变；Ordinal 大小写纪律不变；仅盘符卷形态参与（`C:` 两字符），
   **UNC 不扩**；界外、不可证、矛盾观察一律拒绝。P52 行为要求为唯一修复依据。
4. **观察时点**：判定锚点 = 候选根最终层与绑定根句柄同值比较的观察时刻；此后行走
   只用原绑定句柄，**不承诺别名随后不变**（变了也不影响本次行走，fail-closed）。
5. 生产入口验证 = `OpenBoundDirFile` 真实短名读回（Go 1.22.12 / 1.27.0 双跑）；
   撤掉新分支同一夹具红；SKIP 不算通过；身份不符/含 reparse 组件/大小写变体三组
   拒绝单测；既有回归与门禁复跑。

## 9. 实施与验证结果（2026-11-21，已落地代码但未暂存/未提交）

- 实现：`boundRelativeSuffixShortAlias`（卷根种子 + 逐步句柄相对打开 + 逐层
  reparse 拒绝 + 分歧组件「同父同记录 + GetShortPathNameW 短名逐字节」双证 +
  最终层与绑定根句柄同卷同文件 ID）+ `splitDosPathComponents` + `shortPathNameOf`；
  `reparseTargetRelativeSuffix` 在字节比较失败后调用该 fallback。
- 测试：`TestP84W12ShortAliasRootPrefixIsAdmitted`（生产入口，真实短名读回
  `inside`）+ 三组拒绝单测（身份不符 / 含 reparse 组件 / 大小写变体）——
  **go1.22.12 与 go1.27.0 双跑全部 PASS、零 SKIP**。
- 变异红基线：临时禁用 fallback ⇒ 同一夹具红（错误为短名误拒原文），恢复后
  SHA256 与基线逐字节一致（`f94a8e9c…`）。
- 回归与门禁：双工具链全量 `go test ./...` EXIT=0（15/15 ok）；SELFTEST 740/740；
  scope=tree TOTAL_FAIL=0、scope=ci TOTAL_FAIL=0；scope=index 仅 G5(f) 红
  （R-1 未暂存所致，预期读数）；IDE 诊断两文件 0。
- 边界确认：N10/W1 大小写纪律不变（W9/W1 全绿）；UNC 不扩（盘符形态硬门槛）；
  未改权威契约；未追改历史审查口径（本修订为追加节）。

## 10. 审核补正（2026-11-21 用户裁定 R-1 补正，已落地并暂存）

1. **第二条证据改为同一目录项证据**：废弃「绝对路径 GetShortPathNameW 返回值
   作独立放行依据」的初版。现实现 `directoryEntryShortName(parent, longName,
   wantHigh, wantLow)`：枚举**已验证父句柄**（类 11 首启 + 类 10 续读，同
   readDirNames 模式）→ 定位长名匹配的目录项 → 该目录项的 64 位 FileId 必须等于
   **已打开对象**的 file id（证据与已打开对象身份绑定）→ 从**同一条目录项**读取
   ShortName 并与候选拼写逐字节相等。无法可靠取得 ⇒ 拒绝。
2. **实现事故登记（如实）**：初版 `directoryEntryShortName` 每轮调用都用 restart
   类 ⇒ 目录枚举永不终止、W12 挂起（`panic: test timed out`）；已改为
   readDirNames 同款的「首启 restart、续读不 restart」模式并复跑全绿。
3. **测试补正与诚实收窄**：
   - 「身份不符」用例改名 `TestP84W12AliasBranchRefusesDifferentRecord` 并收窄
     声明：同记录比较、同目录项证据与最终身份比较对该形态互相重叠（单一检查的
     变异无法区分），最终比较的拒绝侧静态夹具不可达（构造性证明），故该用例钉住
     **合并拒绝**；判别变异 = 整段证明 + 最终比较全删 ⇒ 红（已执行）。
   - reparse 用例收窄声明：仅换 follow=true 时仍被同目录项证据拒绝（已登记实测）；
     判别变异 = follow=true + 证明全删 ⇒ 红（已执行）。
   - 大小写变体用例：删除同目录项证据 ⇒ 红（已执行，证据即该用例的判定检查）。
   - 四例双工具链（go1.22.12 / go1.27.0）逐例 PASS、零 SKIP。
4. **变异与恢复摘要**：fallback 禁用 ⇒ W12 正例红；M1C（证明块+最终比较全删）⇒
   different-record 红；M2 单独 follow 换向 ⇒ reparse 仍拒（重叠证据，登记）；
   M2C''（follow+证明全删）⇒ reparse 红；M3（删同目录项证据）⇒ case-variant 红。
   恢复后 `bounddir_windows.go` SHA256 = `455acdda…` 与变异前逐字节一致。
5. 暂存与门禁：仅精确暂存 `bounddir_windows.go`（R-2+R-1 合并）与
   `p84_bounddir_residual_test.go`；porcelain 仅两条 `M `；scope=index
   TOTAL_FAIL=0、scope=tree TOTAL_FAIL=0、SELFTEST 740/740；双工具链全量
   `go test ./...` EXIT=0（15/15 ok）。
   （行数勘误 2026-11-21：§10 早前「+288/−2、+216」为补正前读数；收窄复审前
   最终暂存快照 = 实现 **+344/−2**、测试 **+365/−0**，见复审记录。）
