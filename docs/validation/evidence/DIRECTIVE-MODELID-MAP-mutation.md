# 变异检验日志 · 切片 `DIRECTIVE-MODELID-MAP`

> **重建归档**（v1.22 收尾；原件乱码不可逆，重建依据见下条注记）：源文件 `tmp/dml/mutation.md`（CRLF、无 BOM）按「移到归档」保留，行尾与编码归一为 UTF-8 with BOM + LF。
> 原捕获中两行含中文摘录的 FAIL 详情行是 PowerShell 重定向产生的乱码，其中若干字节不可逆；本副本按 `docs/DELIVERY_DIRECTIVE.md` 的对应原文与门禁自身的截断口径（G6-6 取处置单元格前 60 字符、G6-5 取行首 120 字符）重建这两行，其余 ASCII 状态行逐字保留。
> 三组变异各自的「目标判据转红 ⇒ 还原 ⇒ 字节级一致」结论不受上述编码问题影响：`sha-original` 与 `VERIFY ... byte-identical=true` 均为 ASCII 行。
MUTATED m1 target-check=G6-6 sha-original=0b8a113ca3180e3e
FAIL  G6-6 ledger status closed set :: docs/DELIVERY_DIRECTIVE.md:867 first status word not in the closed set :: **下一片必做**：机制已由本片（`DIRECTIVE-MODELID-MAP`）定义于 §1.5.1（清理时机 / 归 whitelist-source=worktree
TOTAL_FAIL=1
RESTORED m1
VERIFY m1 byte-identical=true sha=0b8a113ca3180e3e
PASS  G6-6 ledger status closed set :: every added ledger row opens with a closed-set status word (scanned 10) whitelist-source=worktree
TOTAL_FAIL=0
MUTATED m2 target-check=G6-5 sha-original=0b8a113ca3180e3e
FAIL  G6-5 section-reference resolution :: docs/DELIVERY_DIRECTIVE.md:68 §1.5.9 unresolved :: | v1.22 | 2026-10-04 | **`modelId` 到限定名的口径反转 + 完成行结构重建 + `tmp/` 生命周期机制（切片 `DIRECTIVE-MODELID-MAP`；含 OB-49 / OB-72 同批评估）* whitelist-source=worktree
TOTAL_FAIL=1
RESTORED m2
VERIFY m2 byte-identical=true sha=0b8a113ca3180e3e
PASS  G6-5 section-reference resolution :: no unresolved section reference on added ledger/changelog lines (scanned 42, refs 111) whitelist-source=worktree
TOTAL_FAIL=0
MUTATED m3 target-check=G3-a sha-original=0b8a113ca3180e3e
FAIL  G3(a.1) docs/DELIVERY_DIRECTIVE.md symmetry :: docs/DELIVERY_DIRECTIVE.md +78/-48 vs _EN +76/-48 MISMATCH
FAIL  G3(a.2) docs/DELIVERY_DIRECTIVE.md mirror parity :: docs/DELIVERY_DIRECTIVE.md {"lines":false,"head":false,"bold":true,"bars":true,"blank":false}
FAIL  G3(a) symmetry + mirror parity :: docs/DELIVERY_DIRECTIVE.md +78/-48 vs _EN +76/-48 MISMATCH ;; docs/DELIVERY_DIRECTIVE.md {"lines":false,"head":false,"bold":true,"bars":true,"blank":false}
TOTAL_FAIL=3
RESTORED m3
VERIFY m3 byte-identical=true sha=0b8a113ca3180e3e
PASS  G3(a) symmetry + mirror parity :: docs/DELIVERY_DIRECTIVE.md +76/-48 vs _EN +76/-48 OK ;; docs/DELIVERY_DIRECTIVE.md {"lines":true,"head":true,"bold":true,"bars":true,"blank":true}
TOTAL_FAIL=0
