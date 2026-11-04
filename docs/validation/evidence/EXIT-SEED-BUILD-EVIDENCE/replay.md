# replay — P55 证据的逐条重放指令

口径：闭合 ③ 复核 F-04。捕获件里的 `COMMAND:` 行是**捕获时的显示形式**（长路径与哈希被替换为 `<seed>` / `<hash>` / `<64>` / `<drifted actual.json>` 等占位），不可直接粘贴执行。本文件给出可直接执行的形式。

## 一、Windows（本机，`go1.27.0`）

```powershell
$W = "tmp\EXIT-SEED-BUILD-EVIDENCE\work"
$env:GOCACHE = "$W\gocache"; $env:GOMODCACHE = "$W\gomodcache"
$env:CGO_ENABLED = "0"; $env:GOTOOLCHAIN = "local"; $env:GOFLAGS = ""
go mod download
go build -trimpath -o "$W\bin\prfrail.exe" ./cmd/prfrail
go build -trimpath -o "$W\bin\prfrail-release.exe" ./cmd/prfrail-release
# seed 根 = testdata\selfhost 五件副本
Copy-Item testdata\selfhost\*.json "$W\seed\" -Force
$seedHash = & "$W\bin\prfrail-release.exe" hash-tree --root "$W\seed"
& "$W\bin\prfrail-release.exe" probe --binary "$W\bin\prfrail.exe" `
  --noop-chain testdata\selfhost\noop-chain.json `
  --executable-chain testdata\selfhost\executable-chain.json `
  --work-root "$W\probe-run" --output "$W\candidate\actual.json" --expected-version 0.1.0-dev
& "$W\bin\prfrail-release.exe" selfhost --commit (git rev-parse HEAD) `
  --seed "$W\seed" --seed-hash $seedHash --candidate "$W\candidate" `
  --expected testdata\selfhost\expected.json --actual "$W\candidate\actual.json"
# 执行完请删除 $W（隔离模块缓存含第三方 .go，会污染 gofmt -l .）；同时移除 GOCACHE/GOMODCACHE 环境变量
```

反例（先 `New-Item -ItemType Directory -Force` 建好候选根目录；否则 `verify.go:85-93` 会先报 `candidate root: …`，拿不到下表期望串。任取其一，期望 exit=1 且 stderr 命中对应串）：

| 反例 | 关键改动 | 期望 stderr |
|---|---|---|
| 根重叠 | `--candidate "$W\seed\candidate"`（放在 seed 根内） | `seed and candidate roots overlap` |
| oracle 归属 | 把 `expected.json` 副本放进候选根，`--expected` 指向它 | `oracle ownership is not independent` |
| seed 变更 | `--seed-hash` 给 64 个 `0` | `seed changed` |
| oracle 漂移 | 另建候选副本，把 `actual.json` 的 `run-noop` 改为 `failed` | `candidate result differs from external oracle` |

## 二、Linux VM（`go1.22.12`）

驱动脚本原文见会话产物 `tmp/EXIT-SEED-BUILD-EVIDENCE/p55-capture-linux.sh`（⑧c 清理前可取）。要点（变量与目录须先自建，不依赖该脚本）：

```bash
export PATH="$HOME/go1.22/bin:$PATH"
export CGO_ENABLED=0 GOTOOLCHAIN=local GOFLAGS=""
export GOCACHE=$WORK/gocache GOMODCACHE=$WORK/gomodcache
export GOPROXY=https://goproxy.cn,direct GOSUMDB=off   # 预取阶段（VM 无法访问 proxy/sum.golang.org）
mkdir -p "$SEED/bin" "$SEED/fixtures" "$CAND/bin" "$WORK/verifier-bin"
git clone --filter=blob:none --no-checkout https://github.com/larsonzh/prfrail.git "$REPO"
git -C "$REPO" worktree add --detach "$SEED_SRC" fc99f553c3066b24230a093d14bc17af0c6a398e
git -C "$REPO" worktree add --detach "$VERIFIER_SRC" 18c9340c60deb30049554491f8905df51d0a3180
git -C "$REPO" worktree add --detach "$CAND_SRC" baa144df5898817e0ac7cd860928bb4a90cbd009
# 各工作树 go mod download 后，进入离线构建阶段：
export GOPROXY=off
go build -C "$SEED_SRC" -trimpath -o "$SEED/bin/prfrail" ./cmd/prfrail
go build -C "$VERIFIER_SRC" -trimpath -o "$VERIFIER_BIN" ./cmd/prfrail-release
go build -C "$CAND_SRC" -trimpath -o "$CAND/bin/prfrail" ./cmd/prfrail
```

漂移输入件重建（确定性，等价于 `02b-18` 所用输入）：

```bash
python3 - "$SEED/fixtures/expected.json" drift.json <<'PY'
import json, sys
doc = json.load(open(sys.argv[1]))
for c in doc["cases"]:
    if c["id"] == "run-noop":
        c["outcome"] = "failed"
open(sys.argv[2], "w").write(json.dumps(doc, indent=2) + "\n")
PY
```

## 三、判定口径

- 成功路径：exit 0 且 stdout 含 `SeedHash` / `CandidateHash` / `OracleHash`。
- 反例：exit 1 且 stdout 为空、stderr 命中上表字符串。
- 「失败不损 seed」：失败运行前后各跑一次 `hash-tree --root <seed>`，两次输出必须逐字相同（`02a-11` / `02b-19` 即此读数）。
