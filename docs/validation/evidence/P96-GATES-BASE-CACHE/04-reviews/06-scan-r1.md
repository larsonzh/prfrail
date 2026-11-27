# ⑥ 独立扫描 r1（逐字原文）

载体 `independent-reviewer` + `MAI Code 2 Flash (copilot)`。切片 P96 `GATES-BASE-CACHE`。日期：2026-11-27。

---

## Section A
`INDEPENDENT SCAN: FINDINGS`

## Section B
- HIGH | [tools/gates/lib/baseCache.js:76](tools/gates/lib/baseCache.js#L76) | 键只基于 porcelain，不包含 ignored Go 文件、外部 GOWORK/模块、gofmt 可执行身份；可产生键不变、结论改变 | 将全部真实输入纳入键，或无法证明完整性时禁用缓存并真跑
- HIGH | [tools/gates/gate.js:129](tools/gates/gate.js#L129) | 热命中直接跳过 gofmt/build/vet/test；gofmt 缺失或被 PATH 替换时，热遍可掩盖冷遍 exit 2 | 命中前校验所有执行入口；至少纳入 gofmt 路径与身份
- MEDIUM | [tools/gates/lib/baseCache.js:201](tools/gates/lib/baseCache.js#L201) | `detail` 只校验类型，畸形/伪造字符串可被当作合法 PASS 回放 | 增加逐条目 detail 精确模式或条目完整性校验
- MEDIUM | [tools/gates/lib/baseCache.js:37](tools/gates/lib/baseCache.js#L37) | 仅限制未跟踪路径数，未限制已修改跟踪路径数；违反 §6.1 的变更路径上限 | 同时限制跟踪与未跟踪变更路径数，或将条款阈值明确为仅未跟踪路径
- MEDIUM | [tools/gates/selftest.js:7967](tools/gates/selftest.js#L7967) | T157–T164 大量依赖注入假值、源码字符串或内部 helper，不是端到端冷/热等价证据；可等价改写绕过 | 增加真实进程/磁盘冷热测试，覆盖 ignored、gofmt、GOTOOLCHAIN、GOWORK 与完整报告 SHA/退出码

## Section C
- `baseCache` 的缓存值只允许 `id=base`、四条固定 label、全部 verdict=PASS，见 [tools/gates/lib/baseCache.js:191](tools/gates/lib/baseCache.js#L191)。结构上不能存 G1–G9 判据结论。
- 热路径逐条回放 `r.id/r.verdict/r.label/r.detail` 后返回，见 [tools/gates/gate.js:130](tools/gates/gate.js#L130)，冷路径按相同顺序输出，正常情况下无 marker、顺序一致。
- 键组装见 [tools/gates/lib/baseCache.js:148](tools/gates/lib/baseCache.js#L148)。`git status -z` 被哈希，但仅解析其中 tracked/?? 路径；Git ignored 文件完全未进入键。
- 具体反例：在稳定 HEAD/status 下修改一个 ignored `*.go` 文件。键不变，冷遍 `gofmt -l .` 会 FAIL，而热遍可回放此前 PASS。修改/删除 `gofmt` 或替换 PATH 中的 gofmt 同理：键只查询 `go`，不校验 gofmt，[tools/gates/lib/baseCache.js:48](tools/gates/lib/baseCache.js#L48) 的工具链变量也未覆盖 gofmt 身份。
- 实测 `GOTOOLCHAIN=auto/local` 时，`go version/go env` 六变量相同，但实现不读取 GOTOOLCHAIN。外部 `GOWORK` 文件内容也未纳入键；两者都构成键不变、结论可变的候选输入。
- PASS-only 写入在 [tools/gates/gate.js:173](tools/gates/gate.js#L173) 与 [tools/gates/lib/baseCache.js:237](tools/gates/lib/baseCache.js#L237)：任一 FAIL 即不落盘，读回也拒绝 FAIL。不存在直接落盘或回放 FAIL 记录。
- 但键不完整时，旧 PASS 记录可继续命中并掩盖当前真实 FAIL；这不是 FAIL 被缓存，而是同键窗口失效。
- 异常读写均被 try/catch 或内部 catch 收敛，基本满足 fail-open；但热命中绕过冷命令及任意 detail 回放并未完全满足“异常不得改变 verdict/detail/exit”。
- §6.1 的冷/热、PASS-only、目录、LRU、CI 兜底和升版规则基本一致；CI 已以 `--skip-base` 全量真跑，见 [.github/workflows/ci.yml:108](.github/workflows/ci.yml#L108)。偏离集中于键完备性、上限口径和条目语义校验。

## Section D
- 修改任一 ignored `.go` 文件且不改变 HEAD/status：可证伪键完备性；冷 FAIL 与热 PASS 分化即可翻转 PASS。
- 替换 PATH 中的 `gofmt` 或删除它：冷 exit 2 与热 PASS 即证伪 fail-open。
- 编辑外部 GOWORK/go.work 或切换 GOTOOLCHAIN：若键、GO六变量与 gofmt 均相同但 go build/vet/test结论变化，即证伪工具链覆盖。
- 把现有条目的任一 detail 改成任意字符串：若仍命中，即证伪畸形条目校验。
- 在实现中补齐上述键段、精确 detail 校验、跟踪路径上限，并以真实冷/热进程复跑三反例后报告全字节一致且 T157–T164 无源码替代断言，可使结论翻为 PASS。

## Section E
- 未实际运行冷/热进程，只能交叉核对源码、diff 与门禁原始读数；无法验证 Windows 下 PATH、gofmt 可执行解析、缓存文件锁及自托管 runner 的持续状态。
- 未验证磁盘条目是否可被非实现进程伪造或损坏后的真实回放。
- 真实端到端证据证明上述反例不存在，并证明报告 SHA、detail 和 exit 逐字一致，可推翻本次 FINDINGS。
