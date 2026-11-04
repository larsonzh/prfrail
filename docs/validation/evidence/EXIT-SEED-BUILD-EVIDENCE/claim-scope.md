# claim-scope — P55 证据的能力边界（本片断言 / 不assert）

口径：本文件闭合 ③ 复核中 F-02 / F-04 / F-05 / F-09 / F-13 一类「证据只支持到此为止」的中等发现——处理方式是**收窄断言**，不是补做实验。凡本文件列为「不 assert」者，报告与台账不得扩大读法。

## 一、本片成立（可复算）

| # | 断言 | 判据 |
|---|---|---|
| S-1 | Linux 侧执行**真两代**验证：seed `fc99f55` 构建 `prfrail`，verifier `18c9340` 构建 `prfrail-release`，candidate `baa144d` 构建 `prfrail`；三角色、三二进制互异 | `02b-21-bindings.json`；`03-test-report.md` 任务 2 |
| S-2 | 双侧 `selfhost` 成功路径退出 0，返回 `SeedHash` / `CandidateHash` / `OracleHash` | `02a-06`、`02b-14` |
| S-3 | 外部 oracle 的四个 case（`version` / `validate-noop` / `run-noop` / `reject-executable`）由候选进程产出 actual，expected 位于候选根之外 | `02a-05`、`02a-12`、`02b-13`、`02b-20` |
| S-4 | 五枚夹具哈希双侧一致，且与 `s1-selfhost.md` 的 CI dispatch 输入表逐字相同 | `02a-bindings`、`02b-21`、`03-test-report.md` 任务 2 |
| S-5 | 四条反例在双侧均 exit≠0 且 stderr 命中预期错误串（根重叠 / oracle 归属 / seed 变更 / oracle 漂移） | `02a-07..10`、`02b-15..18` |
| S-6 | 失败运行（漂移反例）后 seed 根哈希不变（`equal=true`），双侧一致 | `02a-11`、`02b-19` |
| S-7 | Linux 侧构建在 `GOPROXY=off` 下完成（构建期无网络依赖） | `02b-00` note、`02b-08..11` |

## 二、不 assert（禁止扩大读法）

| # | 不成立的读法 | 事实依据 | 对应 ③ 发现 |
|---|---|---|---|
| N-1 | 「Windows 侧证明了核验者独立性」 | Windows 侧 `prfrail-release.exe` 与 `prfrail.exe` 同由 HEAD 工作树构建，`bindings.json` 无 `verifierCommit` ⇒ 该侧只证明「工具在自身产物上可跑通 + oracle 语义 + 反例 fail-close」；**独立性仅在 Linux 侧成立** | F-02 |
| N-2 | 「`SeedHash` 等于 seed 提交的源码树哈希」 | seed 根是**组装目录**：Windows 侧为 `testdata/selfhost/` 五件副本；Linux 侧为四项——`bin/prfrail`（由 `fc99f55` 构建）、`fixtures/`（五枚夹具）、`seed-commit.txt`、`toolchain.txt`。均非 `git archive` 输出 ⇒ `SeedHash` 证明**运行期未被改动**，不证明与提交内容等价 | F-13 |
| N-3 | 「OracleHash 证明 actual 与 expected 字节相同」 | 比较是**结构化语义**（`schemaVersion` + 逐 case `id`/`outcome`）；actual 字节哈希 `adc6e927…` ≠ expected `388f8884…` 仍判成功 | F-09 |
| N-4 | 「四条反例均可按捕获中的 `COMMAND:` 行原地重放」 | 捕获的 `COMMAND:` 行含 `<seed>` / `<hash>` / `<64>` / `<drifted actual.json>` 等**显示占位符**；可重放指令见 `replay.md`，漂移输入件是**确定性重建件**（由 expected 翻转 `run-noop`），非原始二进制 | F-04 |
| N-5 | 「已用文件系统级审计证明任何路径都不写 seed」 | 未挂载只读卷、未开写审计。现有保证来自**代码层顺序**：`internal/release/verify.go:94` 的 overlap 检查先于 `:110` 的 seed 哈希，且 `VerifySelfHost` 全程只读（`HashTree` / `os.ReadFile` 等，无写调用）⇒ 重叠反例**结构上不可能**触及 seed；其余三条反例根本不指向 seed | F-05 |
| N-6 | 「Windows 侧工具链与 CI 同形」 | Windows 实测 `go1.27.0`，CI 固定 `1.22.12`；同形证据仅在 Linux 侧（`02b-00`） | 既有登记（G-6） |
| N-7 | 「oracle 一致 ⇒ 候选无恶意」 | 沿用 `s1-selfhost.md` 既有边界：oracle 比较只测**行为漂移** | 既有声明 |

## 三、未覆盖（本片范围外）

真实 `workflow_dispatch` 触发、付费模型调用、Windows 断电耐久注入、恶意候选证明、SessionBridge 交互面——均不属本片，且本片未触碰 `.github/workflows/**` 与 Go 产品代码。
