# ⑦ 深度档审查原文归档（P96）

载体 `independent-reviewer`；主读 `GPT-6.2 Sol (copilot)` Max；交叉复核 `Gemini 4 Neon (copilot)`。切片 P96 `GATES-BASE-CACHE`。日期：2026-11-27 / 2026-11-28。

> 下列各轮保留原文的 Section A（结论行）与 Section B（发现表）；Section C–E 的完整推理见当轮回执，本文件不重写。

## r1（终审）

```
## Section A
FINAL REVIEW: FINDINGS

## Section B
HIGH | gate.js:145 | 「gofmt 二进制 sha256」 | 哈希与冷跑使用同一个解析后的可执行文件；无法确认则冷跑。
HIGH | baseCache.js:182 | 「键完全命中（四类键成分全等）即有效」 | ignored 路径用 -z 原样解析；外部 Go 输入无法验证时冷跑。
MEDIUM | baseCache.js:56 | 「verdict / detail / 退出码与无缓存运行完全一致」 | 校验记录摘要以检出静默损坏，并明示本地写入者的信任边界。
MEDIUM | baseCache.js:364 | 「任何缓存异常 ⇒ 一律回全量计算」 | 更新时间戳失败时视为未命中，或在条款中明确豁免。
MEDIUM | baseCache.js:34 | 「必须升级格式版本前缀」 | 键和存储结构已变，同步升级前缀、条款与断言。
MEDIUM | verify-pack.ps1:7 | 「基线必须材料化进审查输入包」 | 补递归集合比对及八件逐项哈希实算原文。
```

## r2（复审 1）

```
## Section A
RE-REVIEW: FINDINGS

## Section B
HIGH | baseCache.js:64 | 哈希的文件可能不是实际执行器，热跑可假绿 | 按原生规则解析并执行同一文件，无法确认则冷跑。
HIGH | baseCache.js:211 | trim() 改写合法路径，内容变更可仍命中 | 只丢弃空段，路径不可读则放弃缓存。
MEDIUM | baseCache.js:380 | mtime 写失败仍回放 | 本遍判未命中，或先修订条款明确豁免。
MEDIUM | DELIVERY_DIRECTIVE.md:647 | §6.1 中英正文仍指定 base-v1 | 双语同步为 base-v2。
```

## r3（复审 2）

```
## Section A
RE-REVIEW: FINDINGS

## Section B
HIGH | baseCache.js:66 | H-1：哈希与冷跑可选中不同文件，导致旧 PASS 误命中 | Windows 无法确认与原生解析一致时禁用缓存；修正键语义后升版。
```

## r4（复审 3）

```
## Section A
RE-REVIEW: FINDINGS

## Section B
- **H-1（高，未闭合）**：候选排序抹掉 PATH 优先级，但执行仍依赖该顺序，可同键换执行器并回放旧 PASS。Windows 还扫描父进程 cwd，而子进程 cwd 是仓库根；从子目录启动可漏扫根目录候选。
- **性能**：每目录 Windows 探测 4 次、POSIX 探测 2 次，可能重复整读同一文件；无候选数或字节预算。
```

## Neon 交叉复核

```
## Section A
NEON REVIEW: FINDINGS

## Section B
- HIGH | F1 POSIX 空 PATH 项：`PATH=":/usr/local/go/bin"`，首跑无 `<repo>/gofmt`，次跑放入可执行假文件 ⇒ 键不变而冷跑执行假文件
- HIGH | F2 POSIX 缺 X_OK 校验：0644 首跑跳过，次跑 chmod 0755 ⇒ 同键而执行器改变
- HIGH | F3 Windows 引号 PATH 漏检：`PATH="C:\s";C:\Go\bin` 的引号目录未被扫
- HIGH | F4 相对 PATH 基准错位：PATH 含相对项时按 `process.cwd()` 查而冷跑按 `ctx.repoRoot`
```

## 收窄复审（额度外，OB-37）

```
[第五轮 · v1.80 收敛后]
## Section A
RE-REVIEW: FINDINGS
## Section B
HIGH | baseCache.js:48 | 键漏 GO111MODULE，GOWORK 也只记路径；同键可回放假 PASS | 补有效模块模式与工作文件摘要并升版

[第六轮 · v1.81 收敛后]
## Section A
RE-REVIEW: FINDINGS
## Section B
HIGH | baseCache.js:214 | 过滤误删 `CGO_CFLAGS=-DGOGCCFLAGS=1` 整行 | 精确识别变量名
HIGH | baseCache.js:203 | `go.work` 不变时改坏其 use 指向的仓外模块源码，键不变 | 仓外本地依赖入键，无法覆盖则绕过缓存
```

## 处置（用户裁决）

- r1–r4 与 Neon 的 gofmt 身份类发现：v1.80 起 **gofmt 恒真跑、不缓存** ⇒ 整体消解。
- GO111MODULE / go.work：v1.81 改 `go env` 全量摘要 + `GOWORK` 保守禁用 ⇒ 消解。
- 过滤过宽（HIGH-1）：改精确变量名过滤（v1.81 内）。
- 仓外模块（HIGH-2）：方案 (b) 保守禁用（v1.81 内）。

## Section C–E 逐字原文（按轮）

> 以下各段自会话转录原样取出，起点为本轮 `## Section C`，终点为本轮文末；未改写、未摘要。
> 第五、六段（收窄复审）原文只含 C / D 两节，无 E 节，此处如实保留缺口。

### r1（终审 · 总体裁决轮）

## Section C
1. **键不完备。** [baseCache.js:207](tools/gates/lib/baseCache.js#L207) 纳入 HEAD、porcelain（含删除和改名）、三类内容摘要、gofmt、Go 版本和七变量。反例：PATH 只换掉 gofmt，GOROOT 摘要不变，热 PASS 而冷 FAIL 或退出 2；默认 Git 引用带非 ASCII 的 ignored 路径，代码读错路径并固定为 `unreadable`；外部 `GOWORK` 路径不变，其替换模块源码改坏也能同键异判。
2. **失效机制通过，前提是键可信。** [gate.js:121](tools/gates/gate.js#L121) 先重算；[baseCache.js:266](tools/gates/lib/baseCache.js#L266) 拒绝值未改、仅放错文件的条目。无 TTL；已编码成分变化即换键。
3. **fail-open 部分通过。** 查键失败、坏 JSON、读失败回冷跑；[gate.js:174](tools/gates/gate.js#L174) 写失败发生在真实结果生成后。时间戳异常仍热回放；第 1 项可实际改变 verdict 和退出码。
4. **PASS-only 通过，真实性不通过。** [baseCache.js:229](tools/gates/lib/baseCache.js#L229) 读写都拒绝 FAIL，固定标签及顺序；但同键的 `ok=15 FAIL=0` 变成 `ok=16 FAIL=0` 仍合法，`key` 绑定不验证内容来源。
5. **冷热输出有条件一致。** [gate.js:131](tools/gates/gate.js#L131) 原序回放四元组，无命中标记。真实条目在同作用域可逐字一致；第 4 项可改变 detail，终版全报告尚待对照。
6. **判据零改动通过。** 十路径补丁中，入口仅新增导入和 base 包裹；[gate.js:275](tools/gates/gate.js#L275) 仍先执行并记录 G1–G9，`--skip-base` 不进缓存。
7. **条款部分通过。** 200 路径、各映射 64 MiB、原子写及尽力 LRU 50；改名双字段计数偏严。[§5.1:559](docs/DELIVERY_DIRECTIVE.md#L559) 例外两成员；flake 仅复用 PASS，CI 双腿 [真跑:82](.github/workflows/ci.yml#L82)。G1-b 的头/末行均为 v1.79，[C.0be:2584](docs/DELIVERY_DIRECTIVE.md#L2584) 可作 v1.78 首落留痕；格式前缀未升。
8. **证据部分通过。** 三轮 ⑥ 正文在档，清单与校验和表八路径同集；未见递归实算或独立原稿比对。[读数:20](tmp/P96-GATES-BASE-CACHE/07-input/gate-readings.txt#L20) 披露三反例及坏 JSON 翻转来自旧版，可作历史佐证；「只扩张键」不准确，值结构和预算也变了。805/805、v1.79 tree/index 跳过 base 均绿；含 base 的 r3 日志仍标 v1.78，不能当终版冷热证据。

## Section D
删内容摘要，T157 同名改内容断言红；删 gofmt 摘要，T165 红，但只换 PATH gofmt 仍绿。删 `parsed.key` 守卫，[T165:8339](tools/gates/selftest.js#L8339) 仍绿，因为 Map 未返回条目；让 `get(当前键)` 返回合法的旧键值，才能证伪该守卫。UNCOVERED：外部依赖、可解析的数字损坏、终版三作用域逐字对照。

## Section E
本轮未执行命令，不能独立确认八件哈希、⑥ 逐字转录或终版实机结果。修复 Section B 并补齐终版冷热证据后可复审；任一非 flake 的同键异判、异 detail、异退出码，或实算哈希不符，均推翻 PASS。

### r2（复审 1）

## Section C
- **H-1 PARTIAL**：Windows `shell:false` 不按 `PATHEXT` 查找 `.cmd`。PATH 首目录只有 `gofmt.cmd`，次目录有 `gofmt.exe`：键哈希前者，冷跑执行后者；只换后者，键不变。
- **H-2 PARTIAL**：`-z` 保留 UTF-8 相对路径，可按仓库根拼接；但开头含 U+3000 的 ignored 路径会被 `trim()` 改名，原文件改字节仍可能同键。
- **M-1 CLOSED**：四条 PASS、固定 label/detail、同键才回放。§6.1 要求「条目缺失/畸形/不可读」时重算；手改后仍合法同键不属于这些形态。键绑定不是来源认证，§5.1 也只是目录例外。
- **M-2 OPEN**：touch 失败虽只影响淘汰顺序，仍是缓存操作异常。§6.1 原文「任何缓存异常 ⇒ 一律回全量计算」未授权 LRU 豁免。
- **M-3 PARTIAL**：实现和 T157 已升 `base-v2`，旧键失效；§0.3、C.0be 的留痕不替代正文。
- **M-4 CLOSED**：递归实算载荷与校验表各 8 项、集合相等、`HASH_MISMATCH=0`、清单 `8/8`；目录也逐项对上。

## Section D
T165 只换模拟执行器字节、增删 ignored 路径；错键用例读的是空槽。应分别只替换真实 `.exe`、只改同名 ignored 文件内容，对比键、冷热报告和退出码；再让读取返回错键非空条目，并注入 touch 失败。现有 `805/805` 和两份 `TOTAL_FAIL=0` 不覆盖这些反例，且 index 跳过 base。

## Section E
本轮只读，未独立重算哈希、实跑反例或验证终版三作用域冷热等价。

### r3（复审 2）

## Section C
- **H-1 PARTIAL**：已排除 `.cmd`，但 [libuv](https://github.com/libuv/libuv/blob/v1.x/src/win/process.c) 可先查子进程 `cwd`，且同目录先 `.com` 后 `.exe`；当前实现只查 PATH，扩展名顺序相反。若 `cwd` 无候选，首个 PATH 目录 `C:\tools` 同时有可运行的 `gofmt.exe` 和 `gofmt.com`，键哈希 `.exe`，冷跑执行 `.com`。
- **H-2 CLOSED**：[baseCache.js:211](tools/gates/lib/baseCache.js#L211) 保留前导空白；读失败时原路径与 `unreadable` 一同进入键摘要，区别于可读状态。连续不可读仍可同键。
- **M-2 CLOSED**：[baseCache.js:379](tools/gates/lib/baseCache.js#L379) 的 `utimesSync` 抛错即返回 `null`，`read` 判为未命中，调用方进入冷跑。
- **M-3 CLOSED**：§6.1 [CN](docs/DELIVERY_DIRECTIVE.md#L647) 和 [EN](docs/DELIVERY_DIRECTIVE_EN.md#L647) 正文均为 `base-v2`。

## Section D
双候选下先填充 PASS 缓存，再仅替换 PATH 中的 `.com`，使格式检查输出非空。若键不变、热跑 PASS 而禁用缓存后 FAIL，即证实 H-1。修复后须变键或放弃缓存，且结果与禁用缓存一致。

## Section E
本轮只读，未独立重跑。r3 记录自测 805/805、tree 冷跑与 index 跳过 base 均零失败；未覆盖 Windows 双候选实跑或访问戳故障注入，终版三作用域冷热对照待⑨。

### r4（复审 3）

## Section C
1. 固定搜索顺序且 cwd 一致时，同目录 `.com`/`.exe` 和 cwd 候选的增删或字节变化会换键；不能推出执行器不变。
2. r4 报告 `SELFTEST 805/805`、六形状 CN=EN（`HEAD_MISMATCH=0`）；原始日志确认 tree 四条 base PASS，tree/index 均 `TOTAL_FAIL=0`（index 跳过 base）。
3. 排序结果稳定，但丢失有语义的优先级。POSIX 还过滤空 PATH 项（代表 cwd），且执行位变化不入键。

## Section D
固定仓库与 `go`，只互换两个含不同行为 `gofmt` 的 PATH 目录。完整键相同而冷跑结果不同，即证伪；T165 只验证字节变化，不覆盖此例。

## Section E
只读审查；未实跑反例或测量 I/O，自测与六形状仅核对 r4 汇总。

### Neon 交叉复核

## Section C

1. **Windows 顺序**：无引号绝对项下，`cwdRoot` 置首、`.com` 先于 `.exe` 与不排序摘要已覆盖（注：[gate-readings-r4.txt](tmp/P96-GATES-BASE-CACHE/07-input/gate-readings-r4.txt#L6) 仍写“排序后”与代码不符）；但被 F3、F4 绕过。
2. **POSIX 等价**：不等价，存在 F1（空项丢失）、F2（缺 `X_OK`）、F4（相对项错位）。
3. **同集替换**：已入集字节替换由 `sha256` 覆盖；纯权限切换（F2）与漏检路径替换（F1/F3/F4）均未覆盖。

## Section D

按每条构造冷跑写缓存，变更后对比默认热跑与 `GATES_BASE_CACHE=off` 冷跑的 `computeKey` 与裁决；若键变化或冷热裁决一致，则该项被证伪。

### 收窄复审 · v1.80 收敛后

## Section C
1. 固定本地工具链，在 `GOPATH/src` 外、无 `go.work` 时先缓存 PASS，再仅设 `GO111MODULE=off`：键不变，冷跑却无法解析[模块导入](cmd/prfrail/main.go#L7)。另一构造是保持外部 `go.work` 路径不变，仅将 `use` 改指另一有效模块：键不变，冷跑拒绝本仓。
2. [冷热两路](tools/gates/gate.js#L115)均先真跑并记录 gofmt，再经同一 `out.rec` 按 build → vet → test 输出，无命中标记。[实测整报告 SHA 相同](tmp/P96-GATES-BASE-CACHE/07-input/gate-readings-v180.txt#L19)，证明该对逐字一致，不能外推到上述同键反例。

## Section D
在隔离副本分别施加两种变异，对照命中与禁缓存冷跑；若键变化或冷跑仍 PASS，对应反例即被证伪。固定输入下，命中路漏记 gofmt、换序或加标记，应使冷热 SHA 不同。本轮只读，未复跑。

### 收窄复审 · v1.81 收敛后

## Section C
1. T157/T165 验证普通 `GOFLAGS` 和 `go.work` 正文变化会失效，未覆盖上述反例。
2. 冷热同序调用 `out.rec`，无命中标记；[实测](tmp/P96-GATES-BASE-CACHE/07-input/gate-readings-v181.txt#L17) `FULL_IDENTICAL=True`。正常回放逐字一致，错命中则不保证。

## Section D
先缓存 PASS，再分别只改宏值或仓外源码。同作用域下，若键不变、热跑仍 PASS 而禁用缓存后 FAIL，即证实错命中。另验证仅随机目录变化时键不变。本轮只读，未复跑。
