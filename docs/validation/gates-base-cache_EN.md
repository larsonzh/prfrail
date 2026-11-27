# P96 GATES-BASE-CACHE validation report (cross-scope base verdict cache)

Slice: `GATES-BASE-CACHE` (P96). Date: 2026-11-28. Size: M. Type: tooling slice (clause first).

## 1. Conclusion

- Verdict: **done**. The cache covers only the three §6.1 verdicts `build` / `vet` / `test`; `gofmt` always runs for real and is never cached.
- Boundary: zero change to the criteria decision logic; no G1-G9 criterion verdict is ever cached; cold and hot reports are byte-identical.
- Version: the directive moved **v1.78 -> v1.79 -> v1.80 -> v1.81** (four amendments in one slice).

## 2. Benefit (measured)

| Case | Cold | Hot | Note |
|---|---|---|---|
| `--all --scope=tree` | 105705 ms | 11484 ms | one entry hit; whole-report sha256 equal |

- Benefit: a multi-scope run inside one slice saves about **262 s** (conservative); the cost is **+1.5-3 s** per hot pass (`gofmt` plus key computation).
- Verdict: the benefit holds; over-invalidation is fail-safe (extra cold runs only).

## 3. Implementation

- Key: `base-v2` | HEAD | porcelain `-z` digest | untracked content map | modified-tracked content map | ignored `.go` map | `go version` | whole `go env` digest (excluding `GOGCCFLAGS`).
- A non-empty `GOWORK` other than `off` abandons the cache for this pass (out-of-workspace module sources cannot be keyed).
- Only **PASS** is cached; FAIL is never stored; label and detail are exact per position; the stored value carries a `key` binding.
- Any anomaly falls back to the full run; the hit path emits no hit/miss marker.

## 4. Acceptance evidence

| Item | Reading | Verdict |
|---|---|---|
| Cold/hot A/B | `FULL_IDENTICAL=True` | pass |
| Three counterexamples | tracked / untracked / toolchain all cold-run and change the key | pass |
| fail-open | flipping one byte returns to the full run and self-heals | pass |
| Non-empty GOWORK | `computeKey = null` | pass |
| `SELFTEST` | 806/806 | pass |
| tree / index | `TOTAL_FAIL=0` | pass |
| Cold/hot twice x three scopes | tree 125038→26392 ms, index 124263→22862 ms, ci 116902→21920 ms, `TOTAL_FAIL=0` | pass |
| CI both legs run `46053254410` | `SELFTEST: PASS 806/806` x2, `GATE REPORT scope=ci changed=17` x2, `TOTAL_FAIL=0` x2 | pass |
| Main commit | `81a6794` | landed |

## 5. Review chain

- (1) `DeepSeek V4 Pro`: `DESIGN: DONE`; (5) pre-review with the same model `PRE-REVIEW: PASS`.
- (3) `DeepSeek V4.1 Flash`: r1 reported F-1, r2 reported F-2, r3 `TEST: PASS` (the controller ran cold/hot and mutations; (3) reviewed afterwards).
- (6) `MAI Code 2 Flash`: r1 `FINDINGS`, r2 `FINDINGS`, r3 `INDEPENDENT SCAN: PASS`.
- (7) deep tier `GPT-6.2 Sol` Max: r1-r4 `FINDINGS`; the Neon cross-check raised four HIGH items; the quota ran out and the user ruled the convergence.
- Narrow re-review (out of quota, OB-37): the v1.80 round raised 1 HIGH (`GO111MODULE` / `go.work`); the v1.81 round raised 2 HIGH (over-broad filter / out-of-workspace modules), both disposed by user ruling.

## 6. Exception handling

- Self-found High: `emit` in `runBaseGates` wrote detail into label, so the cache was never written and the report detail degraded; fixed in place before (4), not counted as R.
- Deviation: the controller ran cold/hot and mutations while (3) reviewed afterwards (endorsed by the user); independence is covered by (6)/(7) and (9).
- HIGH-1: the `GOGCCFLAGS` filter was too broad and would drop a `CGO_CFLAGS` line; replaced by an exact variable-name filter.
- HIGH-2: out-of-workspace `go.work` module sources cannot be keyed; option (b) disables the cache conservatively.
- Found during implementation: `go env` embeds a random `-ffile-prefix-map`, so `GOGCCFLAGS` is excluded.

## 7. Known limits

- Hand-editing a cache entry from the local machine is outside the threat model (§5.1 runtime cache, same model as `g4b`).
- An irrelevant environment variable change causes extra cold runs (correctness first, stated in the clause).

## 8. Version package

| Version | Content |
|---|---|
| v1.78 | §6.1 first lands the base verdict cache clause |
| v1.79 | (6) remediation: key segments, `GOTOOLCHAIN`, `key` binding |
| v1.80 | `gofmt` always runs for real; the cache narrows to three |
| v1.81 | whole `go env` digest + conservative `GOWORK` disable |

## 9. Artifacts

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
tools/gates/gate.js	repo
tools/gates/selftest.js	repo
tools/gates/lib/baseCache.js	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/gates-base-cache.md	repo
docs/validation/gates-base-cache_EN.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/00-briefing.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/01-design.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/input-list.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/freedom-list.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/MANIFEST.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/SHA256SUMS.txt	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/04-reviews/03-tester-review.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/04-reviews/06-scan-r3.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/04-reviews/07-deep-review-log.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/05-batch/01-mirror-check.md	repo
docs/validation/evidence/P96-GATES-BASE-CACHE/05-batch/02-number-check.md	repo
```

## 10. (6) independent scan verdict (verbatim)

```sixstage
## Section A
`INDEPENDENT SCAN: PASS`

## Section B
`none`

## Section C
- M-1：**CLOSED（条款内）**。条目仅缓存 PASS；残缺、非法 detail、非法 key 均不回放，异常回全量。
- M-2：**CLOSED**。删除路径仍计入 porcelain 总预算，且状态摘要本身进入键。
- M-3：**CLOSED**。`GOTOOLCHAIN` 已在实现变量集与条款键段中。
- LOW：**CLOSED**。条款已同步列出 ignored `.go` 与 gofmt 相关表述；运行时缓存例外固定读数已同步为 2。

## Section D
可证伪点：任一四条形态可伪造记录、FAIL 被写入、201 删除未触发 null、`GOTOOLCHAIN` 变更未失效，即可反推对应项未闭合。

## Section E
未独立重跑；结论依据 r3 门禁读数、源码、差异与条款。
INDEPENDENT SCAN: PASS
```
