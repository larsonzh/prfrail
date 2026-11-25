# base 四件套覆盖等价论证（P90）

- 日期：2026-11-25
- 依据：`tools/gates/gate.js` 的 `runBaseGates`（磁盘真值）、`.github/workflows/ci.yml`（改动后）、`internal/release/workflow_test.go`（改动后）
- 结论：`--skip-base` 后 CI 对 section 6.1 基础门禁的覆盖**等价或更强**；判据面（G1–G9 + G4b）逐字不变。

## 1. base 四件套的精确命令（gate.js 原文）

| 序 | 命令 | 判定条件 | 位置 |
|---|---|---|---|
| 1 | `gofmt -l .` | `status === 0 && stdout 去空白后为空` | `gate.js:103-106` |
| 2 | `go build ./...` | `status === 0` | `gate.js:108-117` |
| 3 | `go vet ./...` | `status === 0` | `gate.js:108-117` |
| 4 | `go test ./...` | `status === 0 && FAIL 行数 === 0` | `gate.js:119-129` |

## 2. CI 对等覆盖映射（改动后）

| base 件 | CI 步 | 命令 | 等价性 |
|---|---|---|---|
| `gofmt -l .` | `Format`（新增，第 8 位） | `gofmt -l .`（`shell: bash`） | **等价**（见 §3 分支表） |
| `go build ./...` | `Build`（第 3 位） | `go build ./...` | **逐字相同** |
| `go vet ./...` | `Vet`（第 4 位） | `go vet ./...` | **逐字相同** |
| `go test ./...` | `Test`（第 5 位） | `go test -count=1 ./...` | **更强**：`-count=1` 绕过 Go 测试结果缓存，每次真跑；包集同为 `./...` |

## 3. gofmt 处置（三选一裁定）

裁定 = **(a) 新增独立格式步**。

- (b) 仅靠「提交前本机门禁绿」纪律 + ADR：使 push / PR 路径失去 gofmt 自动防线，相对现状是**倒退**，不取。
- (c) 并入 `Build` 步或换用 `go fmt` / `gofumpt`：改变命令字面，失去与 `gate.js` 的逐字对应，等价论证负担反而加重，不取。
- (a) 落盘形态：

```yaml
      - name: Format
        shell: bash
        run: |
          set -e
          files="$(gofmt -l .)"
          if [ -n "$files" ]; then
            printf '%s\n' "$files"
            exit 1
          fi
```

四分支等价表（`gate.js` 的 `gf.status === 0 && gfOut === ''` vs Format 步）：

| gofmt 行为 | gate.js | Format 步 | 一致 |
|---|---|---|---|
| 退出 0，无输出 | PASS | 不进入 `if`，退出 0 | ✅ |
| 退出 0，有输出 | FAIL | `exit 1` 且打印列表 | ✅ |
| 非零退出，无输出 | FAIL（`status !== 0`） | `set -e` 在命令替换处取非零，直接红 | ✅ |
| 非零退出，有输出 | FAIL | 同上（先红，不一定打印） | ✅（均为红） |
| gofmt 不可执行 | UsageError（exit 2） | 127（命令未找到） | ✅（均为红） |

双侧 runner 可行性：`shell: bash` 在 `ubuntu-latest` 为原生 bash；在 `windows-latest` 路由至镜像预装的 Git Bash（GitHub 文档化支持），`gofmt` 经 `actions/setup-go` 已在 PATH。本仓库此前无 `shell: bash` 先例 ⇒ 列为推送后双腿验证点；若 Windows 腿证伪，则为 fail-closed 红腿（不会逃绿），处置为回滚另议，不降级为近似形态。

## 4. 去重不掩盖失败

- 改动后 step 序（`test` job）：`Checkout source` → `Setup Go` → `Build` → `Vet` → `Test` → `Race`(Linux) → `Contract fixtures`(Linux) → `Format` → `Gates`。
- 全 job 无 `continue-on-error`（默认 fail-closed）⇒ 任一步非零即 job 红、`Gates` 不执行。
- 第 6–8 步不写 `.go` 源（`Race` 只跑测试、`Contract fixtures` 只跑 npm、`Format` 只读）⇒ 不存在「早步通过、base 才红」的确定性窗口。
- 旧实现中 `Gates` 的 base 重跑与前述步骤**同树、同命令、同工具链**，属确定性冗余。
- **已知边界（须长期保留）**：`Test` 与 `Gates` 之间不得插入会修改 `.go` 源的步骤；若将来插入，本去重将使该改动失去 base 复核。该不变量为 ① §I R1 登记项。

## 5. 判据面与 `TOTAL_FAIL` 口径

- `--skip-base` 只减少 4 条 `base` 记录行；`GATE REPORT` 头部、G1–G9 + G4b 全部记录行、`TOTAL_FAIL=` 行逐字不变（③ A/B 实测：剔 base 后 sha256 两侧同为 `ee00e354…`）。
- `TOTAL_FAIL` 的实现口径（`lib/report.js`：计 `verdict === 'FAIL'` 的记录数）不变；但**构成集合收缩**为「纯判据 FAIL 数」——不再可能包含 base FAIL 记录。CI 上该收缩不可观测（base 型失败由更早步骤红掉，`Gates` 不产出报告）。此项已如实披露，防止被误读为「判据宽松化」。

## 6. ③ 一致性实测（本机）

- 命令对：`node tools/gates/gate.js --all --scope=ci` 与 `node tools/gates/gate.js --all --skip-base --scope=ci`，同一 HEAD、同一作用域。
- 结果：`strippedBaseLines=4`、`baseLinesInB=0`、两侧 `TOTAL_FAIL=0`、`changed=6` 相等、`identical=true`、`aPrimeSha256 === bSha256 === ee00e354…`。
- 比对器可证伪性：翻转 1 字节后 `identical=false`（`comparator-sanity.json`）。
- **变更集口径披露**：`--scope=ci` 在本地解析为 `HEAD^1..HEAD`，即提交 `6e097fd` 的 6 文件变更集，而非 P90 未提交改动；被测性质是「旗标不改变判据记录」，与变更集内容无关。P90 自身变更集由 ④ `--scope=tree` 与推送后 CI 双腿覆盖。
- 本机耗时（同一机器、同一次会话）：含 base 108.3 s，`--skip-base` 11.2 s（差额主要来自 `go test ./...`）；该读数仅作量级参考，权威对照取 CI 同 runner 读数。
