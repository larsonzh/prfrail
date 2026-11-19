# 评审输入包消毒与完整性登记（P79）

依据：方案书 §3（消毒六规则 + 材料化要求）。生成方式：脚本按磁盘实测（`tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/gen-package-manifest.js`），不手工记数。

**结论**：仓内输入件不含 SSH 目标串、用户目录绝对路径与真实凭据。IPv4 命中全部为回环地址 `127.0.0.1`（文档化的零网络能力架构事实，非外部主机），凭据赋值形态命中为代码字段赋值与散文中的 `token =` 措辞（非真实凭据）。`%TEMP%` / `$env:TEMP` 的命中全部来自准则与台账的规则文本（临时件纪律条款本身），属规范性表述。三类命中均按「保留」处置，逐件命中次数见下表。

## 逐件登记

| 文件 | sha256（前 16 位） | 字节数 | 行数 | IP | SSH 目标串 | 用户目录路径 | 凭据赋值 | 临时路径文本 |
|---|---|---|---|---|---|---|---|---|
| `docs/SECURITY.md` | da79548d45343def | 10861 | 76 | 0 | 0 | 0 | 0 | 0 |
| `docs/SECURITY_EN.md` | 8f99b6193f28edb4 | 12146 | 76 | 0 | 0 | 0 | 0 | 0 |
| `docs/PRODUCT_REQUIREMENTS.md` | 2d6cf416c22f265a | 16667 | 132 | 0 | 0 | 0 | 0 | 0 |
| `docs/PRODUCT_REQUIREMENTS_EN.md` | 0f4b73a3b7ef3e50 | 18395 | 130 | 0 | 0 | 0 | 0 | 0 |
| `docs/ADR_REGISTER.md` | c406c5e841c22588 | 65499 | 117 | 0 | 0 | 0 | 0 | 0 |
| `docs/ADR_REGISTER_EN.md` | 0cd7f22a3e741810 | 77732 | 118 | 0 | 0 | 0 | 0 | 0 |
| `docs/DEV_PLAN.md` | a252abf9a0aa45b0 | 86644 | 265 | 1 | 0 | 0 | 0 | 0 |
| `docs/DEV_PLAN_EN.md` | ca3207931925b3b3 | 99231 | 266 | 1 | 0 | 0 | 0 | 0 |
| `docs/TEST_STRATEGY.md` | 76b3b4856974f845 | 13560 | 90 | 0 | 0 | 0 | 0 | 0 |
| `docs/TEST_STRATEGY_EN.md` | 749b567b49d63f95 | 14808 | 90 | 0 | 0 | 0 | 0 | 0 |
| `docs/DELIVERY_DIRECTIVE.md` | f35331cf2897c4f5 | 506266 | 2474 | 0 | 0 | 0 | 0 | 29 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 841fd468f9c396a2 | 617890 | 2474 | 0 | 0 | 0 | 0 | 29 |
| `docs/t027/CHECKLIST_v65.md` | 24a9b9786f70c989 | 28992 | 221 | 0 | 0 | 0 | 0 | 9 |
| `docs/t027/CHECKLIST_v65_EN.md` | 3a48ea8519598733 | 33671 | 221 | 0 | 0 | 0 | 0 | 9 |
| `docs/t027/REMAINING_SLICES.md` | 0ad97e4bafc74fe1 | 381800 | 2448 | 2 | 0 | 0 | 0 | 1 |
| `docs/t027/REMAINING_SLICES_EN.md` | 612e26986caafb12 | 428686 | 2435 | 2 | 0 | 0 | 1 | 1 |
| `docs/validation/sg6-gatelevel-consolidate.md` | de98c8910f92b4ed | 21612 | 263 | 0 | 0 | 0 | 0 | 0 |
| `docs/validation/s1-exit.md` | 0786414927be3685 | 10221 | 77 | 0 | 0 | 0 | 0 | 0 |
| `docs/validation/s1-exit_EN.md` | 6cf72dca8dc4eb63 | 11792 | 77 | 0 | 0 | 0 | 0 | 0 |
| `tools/gates/gate.js` | 0d43b5cc8a627511 | 11126 | 295 | 0 | 0 | 0 | 0 | 0 |
| `tools/gates/lib/criteria/g1.js` | 170c2ced90e6a1b1 | 11994 | 259 | 0 | 0 | 0 | 0 | 0 |
| `tools/gates/lib/criteria/g5a.js` | 8fa276638553ec3b | 2240 | 67 | 0 | 0 | 0 | 0 | 0 |
| `tools/gates/lib/criteria/g7.js` | 75e969cbf3fb419a | 46219 | 1060 | 0 | 0 | 0 | 0 | 0 |
| `tools/gates/lib/criteria/g9.js` | 920f4799eb0b0160 | 8737 | 205 | 0 | 0 | 0 | 0 | 2 |
| `internal/gates/policy.go` | adc85a07b2bffb64 | 8782 | 263 | 0 | 0 | 0 | 0 | 0 |
| `internal/gates/runner.go` | 7b688dc4b43d0164 | 14238 | 399 | 0 | 0 | 0 | 0 | 0 |
| `internal/gates/readiness.go` | e73e884bb7513e88 | 9321 | 238 | 0 | 0 | 0 | 0 | 0 |
| `internal/gates/probe_permit.go` | d6ceca879bfed533 | 12926 | 319 | 0 | 0 | 0 | 0 | 0 |
| `internal/guard/lease.go` | 0743c5e670b5745a | 11563 | 331 | 0 | 0 | 0 | 1 | 0 |
| `internal/guard/sandbox.go` | 6a4f2376056f0305 | 4081 | 115 | 0 | 0 | 0 | 0 | 0 |
| `internal/guard/sandbox_unsupported.go` | c5bf8ea9f06c5898 | 574 | 20 | 0 | 0 | 0 | 0 | 0 |
| `internal/chain/postflight.go` | 77afafbe71ac61b0 | 19294 | 464 | 0 | 0 | 0 | 0 | 0 |
| `internal/evidence/verify.go` | 40e9ae930b825ddc | 4188 | 120 | 0 | 0 | 0 | 0 | 0 |
| `internal/evidence/disposition.go` | af38b48f526299e0 | 26896 | 717 | 0 | 0 | 0 | 0 | 0 |
| `internal/evidence/event.go` | d554002983074ee3 | 9993 | 271 | 0 | 0 | 0 | 0 | 0 |
| `internal/adapters/ai_secret_windows.go` | e8e10be5185d0b18 | 4890 | 165 | 0 | 0 | 0 | 0 | 0 |
| `internal/adapters/ai_provider.go` | 7ff770fe27775160 | 3876 | 122 | 0 | 0 | 0 | 0 | 0 |
| `internal/adapters/ai_probe_copilot.go` | b805ea42748d9731 | 8529 | 222 | 0 | 0 | 0 | 0 | 0 |
| `internal/console/ai.go` | 4d437cb1429da7ef | 10466 | 259 | 0 | 0 | 0 | 0 | 0 |
| `docs/OPERATIONS.md` | 8061c41515c0dcbb | 18308 | 148 | 0 | 0 | 0 | 0 | 0 |

## 合计与判定

- IPv4：命中合计 6；
- SSH 目标串：命中合计 0；
- 用户目录绝对路径：命中合计 0；
- 凭据形态：命中合计 2；
- 仓库外临时路径：命中合计 80；

说明：本片**未执行任何改写**（无命中项需过滤），故每件的「消毒前后 SHA-256」为同一值；表内单列该值，行数列用于与 `input-list.md` 逐件对账。

判定：**需过滤项 = 0**——SSH 目标串与用户目录路径两类为 0；IPv4 命中均为回环地址 `127.0.0.1`；凭据赋值形态命中为字段赋值与措辞，非真实凭据；临时路径文本类命中为规则引文。四类皆按「保留」处置，并逐件登记命中次数以便审查方独立复核。
