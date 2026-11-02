# P49 RACE-LEG-EXPAND ⑤ 预审结论（逐字归档）

> 载体：⑤ 预审 `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`。以下为代理输出原文，未改写。

---

# P49 RACE-LEG-EXPAND ⑤ 预审结论

## 审查范围与核对动作

只读审查，未修改任何文件。实际核对动作：① 读取 `.github/workflows/ci.yml` 全文与 `git diff`（ci.yml、workflow_test.go 均为 1 行改动）；② node 独立复算新 run 行字节 SHA-256（169 字节）；③ 全仓 `git grep` 检索 `ci\.yml|\.github/workflows` 与新旧摘要 `ff4826f2|825bacc3`；④ 复跑 `go test -count=1 ./internal/release/`（23.776s ok）与 `go test -count=1 ./...`（15/15 ok）；⑤ 独立运行 `yamlcheck` 与两改动文件编码检查；⑥ 重算证据包 `SHA256SUMS.txt` 全部 14 条目并抽查编码；⑦ 读取证据包全部载荷（00/01/03/04/09、freedom/input-list、MANIFEST、VM 日志）。

## 逐条结论

**1. ci.yml 修改安全 —— PASS。** `git diff` 确证唯一改动为 `ci.yml:84` run 行（1 insertion/1 deletion）。新行六包清单顺序 `adapters → chain → console → gates → guard → tickets` 与 R-03 裁定一致；每包均为 `./internal/<pkg>/...` 形态且各出现恰 1 次；`CGO_ENABLED=1 ` 前缀、`-race`、`-count=1` 均保留；`if: runner.os == 'Linux'`（`ci.yml:83`）未动。YAML 风险：run 值为纯标量，含 `=`、`-`、`/` 但无 `冒号+空格`、无 `#`，无需引号；独立复跑 yamlcheck 通过（`PASS … 16432 bytes`）。Windows 腿：Race 步骤带 `if` 条件，Windows 矩阵腿跳过执行，且本片未触碰该腿任何步骤与顶层 `env`（`CGO_ENABLED: "0"` 保持 job 级语义）。

**2. 冻结契约同步正确性 —— PASS。** 独立复算：剥离 `run:` 前缀后的脚本体（169 字节）SHA-256 = `825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9`，与 `internal/release/workflow_test.go:245` 新摘要逐字符一致。`git diff` 确证该文件仅此 1 行改动，契约其它条目（Build/Vet/Test/Contract fixtures/Gates/candidate-*/selfhost-*）未被误改。复跑 `go test ./internal/release/` 全绿，全仓 `go test ./...` 15/15 ok。编码合规（无 BOM、纯 LF）。

**3. ① 遗漏定性与补救 —— FINDING（已闭环，定性为发现级失误）。** ① 声称「无活体一致性校验引用该行（只查了 tools/gates）」为**检索不完整**：真正的活体校验是 `internal/release` 的 `validateWorkflowScripts` 冻结契约（`workflow_test.go:65/245/387/453`、`workflow_mutation_test.go:13` 均读取并锁定 ci.yml 各步骤 `step.Run` 摘要）。遗漏导致 r1 后该契约族 7 个测试同因红（实际摘要 `825bacc3…` ≠ 旧期望 `ff4826f2…`），r2 同步后转绿。补救充分性：**完全闭环**。全仓检索逐一核对：`internal/` 仅 `release` 包 2 文件引用 ci.yml（同一摘要映射，已绿）；`tools/` 仅 `gates/testdata/` 两处 fixture 文本提及 `.github/workflows/**`（触发词语法示例，非活体校验）；`cmd/`、`harnesses/`、`schemas/`、`.github/hooks/` 零命中；旧摘要 `ff4826f2` 全仓零残余引用。**无第三处活体校验。**

**4. 超时与 cgo —— PASS。** VM 实测四新包 `-race` 合计 34.3s（console 20.410s / gates 3.384s / guard 3.043s / tickets 7.420s，rc=0，见 `03-tests/logs/vm-race-four-pkgs-green.txt`，Ubuntu 24.04 + go1.22.12 + gcc 13.3，与 CI 工具链对齐）；adapters 历史 71.5–91.5s、chain 1.6–2.5s ⇒ Race 腿总时长 ≈ ≤130s + 编译开销，低于 GLM-2 的 230s 保守上界（约 4 分钟）。包级 `go test` 默认 10m 按包计时：最大包 adapters ≈91s 裕度 >6×，新包 ≤21s 裕度 >28×，无风险；`timeout-minutes` 在 prfrail `ci.yml` 零命中（GitHub 默认 job 360m）。cgo：行内 `CGO_ENABLED=1` 前缀保留，仅作用于该步骤；ubuntu-latest 自带 gcc（历史 t028×2、haiku-bump 三次 Race 全绿实证）。flake 风险：GLM-2 B1/B2/B3 清单已列（Top3：console e2e 20s CheckTimeout×约 10 次 re-exec、gates 2s ready 轮询、guard/tickets 毫秒窗口），VM 首跑零 flake，且外层均有 5–10s 兜底；六包 `testing.Short` 零命中（无降载通道，属已知代价）。残留风险可接受，条件见 F4。

**5. 红绿证据链 —— PASS。** 链条成立且可证伪：红基线在未改 ci.yml 上静态断言（`03-tdd-red-run.txt`，exit 1 RED，前 6 项形状断言 PASS、4 项包覆盖断言 FAIL，排除解析误报）→ ② 改后翻转（`03-tdd-green-run.txt`，exit 0 GREEN，六包 occurrences=1 + 顺序 PASS）→ VM 真实 `-race` 四包全绿（rc=0，HEAD 81dd49f clean，工具链与 CI 一致；命令为四新包直跑，不依赖 ci.yml 内容，证据独立有效；adapters/chain 已有历史 CI race 实证）。链条足以支撑「扩包可行且不引入已知竞态」；最终权威仍是 CI 腿首跑（⑧ 观察项）。证据包合规：SHA256SUMS 14/14 条目重算一致；`.md` 均 BOM+LF、`.txt` 均无 BOM+LF。

**6. 边界合规 —— PASS（附台账说明）。** 代码侧仅 `ci.yml`（1 行）+ `workflow_test.go`（1 行）；无任何生产 `.go` 改动（`git status` 佐证）；未删、未重排任何 CI 步骤。另发现 `docs/t027/REMAINING_SLICES{,_EN}.md` 各有 +21 行改动（P49 台账登记行 + 切片表格，系起跑消息 §五.3 规定的步骤 0 动作），属切片授权范围内，但未列入本次改动集摘要（见 F3）。

## 发现清单

| ID | 级别 | 描述 | 建议 |
|---|---|---|---|
| F1 | Medium | ① 声称「无活体一致性校验」实为检索遗漏：`internal/release` 冻结契约锁定 `test/Race` 摘要，r1 后 7 测试红 | **已闭环**（r2 同步，本审复算摘要一致、测试全绿）。剩余动作：⑧ 结片报告显式登记该遗漏与补救经过（含「① 类检索必须覆盖 internal/release」的不重犯声明），由 ⑥/⑦ 复核确认无第三处活体校验 |
| F2 | Low | r2 关键原始日志仅存 `tmp/RACE-LEG-EXPAND/`（`04-red-before-contract-sync.txt`、`04-gotest-all-green.txt`、`04-release-workflow-tests-verbose.txt` 等），证据包无副本；⑧c 清理后证据只剩报告转述 | ⑧c 清理前将上述 3 个关键日志复制入 `03-tests/logs/`，重算并更新 `SHA256SUMS.txt` 与 MANIFEST 条目 |
| F3 | Info | `docs/t027/REMAINING_SLICES{,_EN}.md` 各 +21 行（步骤 0 台账登记）未列入改动集摘要 | ⑧ 最终会计将两文件计入本片改动集并回写「已完成」状态；⑨ 前核对台账行与事实一致 |
| F4 | Info | 四新包首次进 CI `-race` 的 flake / 真实竞态残留风险（GLM-2 B1/B2/B3 已列，VM 首跑零触发） | ⑧ CI 首跑单独观察 Race 腿：采集六条 `ok <pkg>` 行与墙钟时长归档；若暴露真实数据竞争需修生产代码，属越界，须上报升级，不得本片擅修 |

## 总体结论

技术改动（ci.yml 扩包 + 冻结契约摘要同步）经独立复核全部正确：摘要与脚本体逐字节一致、全仓测试 15/15 绿、YAML 校验通过、VM 真实 `-race` 四包全绿、无第三处活体校验、无生产代码改动。F1 虽定性为发现级失误，但已在 r2 完全闭环且经本审验证；剩余整改项均为文档/台账级（F2/F3），无需回炉 ⑤，可在 ⑧ 执行并交 ⑥/⑦ 复核。

**PRE-REVIEW: PASS WITH FIXES**（Fixes = F1 登记义务 + F2 证据副本 + F3 台账会计，均不涉及代码；如主控裁定 F2/F3 由 ⑧ 直接执行，本片可继续派 ⑥）。

---

## 主控裁定（对 ⑤ FIXES 的处置）

- F1：已闭环（r2）。结片报告将显式登记 ① 遗漏与补救经过，含不重犯声明（① 类检索必须覆盖 `internal/release` 冻结契约）。
- F2：已执行——`04-contract-diff.txt`、`04-release-green.txt`、`04-gotest-all-green.txt`、`04-red-before-contract-sync.txt`、`04-release-workflow-tests-verbose.txt`、`04-ciyml.diff` 六件已复制入证据包 `04-impl/logs/`，随下一轮哈希重建入 `SHA256SUMS.txt`（落点按角色目录 `04-impl/logs/` 而非 ⑤ 建议的 `03-tests/logs/`，因该证据属 ② 产出）。
- F3：由 ⑧ 台账回写执行。
- F4：登记为 ⑩ 观察项（CI 首跑 Race 腿六包 `ok` 行采集）。

---

## 主控补记 · 环境事实纠正（2026-11-02，用户同轮裁定，⑥⑦ 审查输入）

**预检盲区如实登记**：VM（10.0.0.199）上其实一直装有 Go——版本 **go1.27.1**，位于 `/home/larson/.local/bin/go`（符号链接至 `/home/larson/.local/opt/go1.27.1/bin/go`，2026-09-07 已存在）。主控此前未探查到，原因是：① 非交互式 SSH 会话未加载用户 PATH（`source ~/.bashrc` 因 bashrc 对非交互 shell 的早退保护而未生效）；② 检索 `/usr/local` 等系统目录时遗漏了用户级目录 `~/.local/bin`。主控随后自行安装的 **go1.22.12（隔离于 `~/go1.22`，官方 tarball + sha256 核验）** 与 CI `GO_VERSION=1.22.12` 完全对齐，其跑出的四包 `-race` 全绿结果**作为隔离验证依然有效**。本片不要求用 1.27.1 重跑（用户裁定：现有证据已足够支撑 P49 闭环，只需记录事实）。

**纠正后的事实口径**（后续报告、⑥⑦ 与 v39 一律采用）：
- VM 既有工具链：go1.27.1 @ `/home/larson/.local/bin/go`（用户环境，登录壳可见）。
- 本片隔离验证工具链：go1.22.12 @ `~/go1.22/bin/go`（与 CI 对齐，四包 `-race` 全绿）。
- 教训（并入 v39 §七 环境预检铁律）：远程 SSH 必须 `source ~/.bashrc` 或 `bash -lc` 加载登录环境，且必须探查用户级目录（`~/.local/bin` 等），确认无既有工具后才可最小化安装。

