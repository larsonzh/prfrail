# anomalies — P55 `EXIT-SEED-BUILD-EVIDENCE` 异常与偏离台账

口径：本文件是 ③ 复核（`03-test-report.md` F-07）指出的「`A-1…A-7` 仅以无编号段落存在」的**编号化落点**。事实源 = 捕获件与本仓库历史；不得作为后续切片的先例。

| 编号 | 事实 | 证据 / 判据 | 处置 | 复核状态 |
|---|---|---|---|---|
| **D = 1** | 无触发词抢跑：步骤 0 在正式起跑消息之前执行，违反 §1.6（无触发词即 `[ROUTINE]`）与 §5.0（流水线只在触发词出现时启动） | 会话时序；`00-briefing.md` 覆盖声明 | 经**用户强制打断**纠正；就地转正；未 commit、未 push（HEAD 保持 `baa144d`），未造成污染 | 用户裁定 |
| A-1 | 首轮 Windows 捕获把 13 个 `.txt` 误写入证据包根而非 `02-win-captures/` | 读取 `02a-06-selfhost-success.txt` 时路径不存在而暴露 | 修正落盘路径、删除误放件、整轮重跑；误放件**从未入库** | **部分可核验**（`index.json` 16 条为修正后产物；误放件已不存在，无法直接比对） |
| A-2 | 主控在会话级导出 `GOCACHE`/`GOMODCACHE` 指向 `tmp/` 下隔离缓存，残留变量使后续 `go test` 重建该目录，`gofmt -l .` 扫到第三方源码并判红一次 | 门禁读数 `FAIL base gofmt -l . :: tmp\…\gomodcache\gopkg.in\yaml.v3@v3.0.1\apic.go` | 清除会话环境变量 + 删除瞬时工作区；`MANIFEST.md` 增「瞬时工作区」条款；复跑 `TOTAL_FAIL=0` | **部分可核验**（现态：`work/` 不存在、门禁绿；当时的判红输出未归档） |
| A-3 | ② 环节由主控自执行（v44 §七 第 9 条要求 `prfrail-implementer`） | 回执派单表 `runSubagent = 0` | 起跑消息【六】明示「主控应亲自执行或指派 Flash 执行命令并收集证据」；用户 2026-11-05 裁定**合规、无需降级** | 用户裁定 |
| A-4 | 计划中的「用 seed 跑 selfhost」不可行：`fc99f553c3066b24230a093d14bc17af0c6a398e` **不含** `cmd/prfrail-release` 与 `internal/release` | `git ls-tree --name-only fc99f55… cmd/` 仅返回 `cmd/prfrail`；③ 独立复算 `git ls-tree -r fc99f55` 亦无 `internal/release`（见 `03-test-report.md` 任务 2） | 改用 T017 既定**三角色结构**：seed `fc99f55`（构建 `prfrail`）、verifier `18c9340`（构建 `prfrail-release`）、candidate `baa144d`（构建 `prfrail`）；三角色与三枚二进制互异 | 用户 2026-11-05 **批准并加重记录**：该替换比原计划更具说服力，实质落实「独立外部 oracle」精神；体现于 `02b-21-bindings.json` |
| A-5 | VM 无法访问 `proxy.golang.org` / `sum.golang.org`（IPv6 黑洞），`GOPROXY=direct` 亦因 `golang.org/x/text?go-get=1` 虚名解析失败 | `02b-00-env.txt` note 行；首两轮捕获中的 `dial tcp [2607:…]:443: i/o timeout` | 预取用 `GOPROXY=https://goproxy.cn,direct` + `GOSUMDB=off`；**构建期强制 `GOPROXY=off`**，反证构建无网络依赖；`go.sum` 校验仍生效 | 用户 2026-11-05 批准（「先拉后断网」策略） |
| A-6 | 远端前两轮执行各产生 9 项 FAIL（模块取不到 ⇒ 构建、探测、selfhost 连锁失败） | 当时捕获已随工作区清理；仅存主控回执记载，**原始件不可复核** | 每轮先读证据原文定位真因再改脚本，未盲目重试；第三轮 0 FAIL | 已闭环（原始件缺失已如实声明） |
| A-7 | 远端 `rm -rf` 因模块缓存只读（`权限不够`）残留一个工作区 | 会话内 `rm` 报错输出 | `chmod -R u+w` 后删除；复核 `remaining=0` | 已闭环（远端已清，副作用不可复核） |

## 记录保真度勘误（⑤ 预审 F-05-1 / F-05-3）

捕获件为**原始输出归档**，不得改写（否则破坏归档保真）。以下两处为捕获时写入的 `TITLE:` / `COMMAND:` 显示文案与真实执行不一致，以本表为**权威更正**：

| 编号 | 位置 | 不一致事实 | 真实情况（权威） |
|---|---|---|---|
| E-1（⑤ F-05-1，Medium） | `02-win-captures/02a-05-probe-actual.txt:3` 的 `COMMAND:` 行；`02-win-captures/index.json` 中 `02a-05` 的 `command` 字段 | 该行缺 `--noop-chain` 与 `--executable-chain`，照此执行会因 `internal/release/probe.go:47`（noop）与 `:51`（executable）的 `existingRealPath` 报 `empty path` 而 exit 1，与捕获的 `EXIT: 0` 矛盾 ⇒ 该行不可作为真实命令 | 真实执行含两条链：`--noop-chain testdata\selfhost\noop-chain.json --executable-chain testdata\selfhost\executable-chain.json`（与 `replay.md` §一、`02b-13` Linux 同形捕获逐参数一致）。捕获的 `EXIT: 0` / 空 stderr / 产出的 `actual.json` 为真；**仅显示文案缩写** |
| E-2（⑤ F-05-3，Low） | `02-linux-captures/02b-05-module-preload-seed.txt` 的 `TITLE:` 行（`02b-06` / `02b-07` 无此残留） | TITLE 末写 `(GOPROXY=direct)`，与 `02b-00-env.txt:16` 的 `GOPROXY_preload=https://goproxy.cn,direct` 及 A-5 不一致（为前两轮失败尝试的旧文案残留） | 预取实际使用 `GOPROXY=https://goproxy.cn,direct`；构建期 `GOPROXY=off`（`02b-08..11`）。不影响 S-7 |

处置口径：不改写两处载荷（保真），以本表勘误；⑧b 验证报告须引用本表 E-1 / E-2。

## 关于 D = 1 的措辞矛盾（③ 复核 F-06）
`00-briefing.md` 是用户消息的**逐字归档**，其中包含用户当时原话「主控已主动发现并停手」。该措辞在同一轮对话内即由用户更正为「经用户强制打断纠正」，并作废我此前回执中的「主动发现」表述。

- **权威表述** = 本表 D = 1 行与 `freedom-list.md` 第三节：**经用户强制打断纠正**（非主控自行发现）。
- **归档件保持原样**：按「归档即原文」纪律，`00-briefing.md` 不改写；其内的旧措辞属**被作废的历史表述**，不得引用为事实。
- 冲突来源是「逐字归档」与「事实更正」并存，不是事实不确定。
