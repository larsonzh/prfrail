# P81 ⑥ 独立扫描 r1 原文（格式偏离轮）

> 载体 / 模型：`independent-reviewer` / `MAI-Code-1.1-Flash (copilot)`（只读）。
> 状态：**格式偏离**——未按 §7.2.2 的 Section A–E 契约输出、无指定结论行 ⇒ 按 §7.2.2 的失败处置重发 1 次（含本轮计入 §7.2 的 3 次有效调用上限）。内容仅作留痕。

## Verdict

✅ PASS — the P81 model replacement is within scope, literal, and internally consistent.

### Why this passes

- The active rule surfaces agree on the direct replacement:
  - [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L117-L121)
  - [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L117-L121)
  - [tools/gates/lib/criteria/g5b.js](tools/gates/lib/criteria/g5b.js#L1-L43)

- The closed set in G5-b is the same live allow-list and contains the renamed entries:
  - mai-code-2-flash
  - gpt-6.2-sol
  while the retired values are excluded from the current rule set.

- The self-tests are meaningful anti-regression checks and would turn red if the retired values were restored:
  - T131a accepts mai-code-2-flash
  - T131b rejects gpt-6.1-sol
  - T131c rejects mai-code-1.1-flash
  See [tools/gates/selftest.js](tools/gates/selftest.js#L4223-L4244).

- The cost and tier numbers were not silently changed beyond the authorized rename:
  - §7.5 keeps the Sol row at 200 / 1000 in [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L780-L800)
  - §8.2 keeps the deep-tier mapping as Extra High with Max allowed in [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L860-L868)
  The English mirror says the same in [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L780-L868).

- There is no evidence of an unapproved drift in the CN/EN mirror: the English row mirrors the same scope and replacement wording, and the directive explicitly says the change is a direct carrier-model swap, with no broader rule expansion.

- The P80 ledger block is used as a direct trace for the authorized exception, not as a stray or mis-numbered OB insertion; this is consistent with the scope described in the v1.70 row at [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L117-L121) and [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L117-L121).

### Scope and audit posture

- No repository files were created, modified, or deleted; this was a read-only verification pass.
- The old identifiers remain only in historical/evidence contexts or negative self-tests, not as live active policy values. That is the correct pattern for a narrow direct replacement.

### Final conclusion

The P81 edit is compliant with the requested rule: direct replacement of the stage-6 and stage-7 carrier model IDs, with no unauthorized expansion, no silent cost/tier mutation, no bilingual drift, and no stray OB-number pollution.
