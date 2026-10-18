# 执行自由度清单 — 切片 `G9-c`（判据 id `G9`）

依据：`docs/DELIVERY_DIRECTIVE.md` §1.7（执行自由度与升级边界）。本清单由主控在 ⑥ 之后、⑦ 之前产出，
⑦ 终审复核本清单判定是否越界；⑧b 定稿时须与验证报告“变更概要”节逐条一致。

## 1 实质性选择（行使自由度）

| # | 选择 | 为什么不属「做什么」 | 落点 |
|---|---|---|---|
| F1 | 判据 id 用 **`G9`**，切片名保留 `G9-c` | id 是命名，不改判据覆盖面 | `gate.js` 注册、`g9.js`、§6.3 行 |
| F2 | 形态闭集取 `.log` / `.tmp` / `.bak` / `.out`，加 `/i` 放宽为大小写不敏感 | 闭集元素取 §5.1「临时件」的典型产物；`/i` 是为 Windows 语义显式放宽，只扩大拦截面 | `g9.js` `TEMP_SHAPE_RE`、§6.3 行 |
| F3 | `tmp/` 之外判形态、`tmp/` 之内判位置，**互斥**（命中 `tmp/` 即 `continue`，同一路径不会产生两条 finding） | 避免重复计数，属实现选择 | `g9.js` G9(2) 分支 |
| F4 | 排除面顺序固定为「证据树 → 已删除 → 白名单 → 判形态/位置」 | 排除先行，减少无意义 finding；语义与 §6.3 一致 | `g9.js` 主循环 |
| F5 | 白名单文件格式 `<match>\t<file>\t<reason>`；畸形态抛 `UsageError`（exit 2）；空表（仅注释）合法 | 格式与 fail-closed 是实现选择；空表合法因为本片 `git ls-files` 实测零合法临时件 | `g9-whitelist.txt`；`ctx.js` 的 `parseWhitelist`；`g9.js` 的 `loadWhitelist` 与 `UsageError` 分支 |
| F6 | 详情最多列 5 条，其余折叠为 `(+N more)` | 报告可读性 | `g9.js` `shown` |
| F7 | 空变更集输出 `INFO`（不 PASS 不 FAIL） | 沿用既有判据先例（空集无行可判） | `g9.js` 早退分支 |
| F8 | `tmp/` 根比较**大小写归一**后判定，finding 仍打印原始拼写；`tmp/.gitkeep` 及其大小写变体放行 | ⑦ r1 M1 整改；`Tmp/` 与 `tmp/` 在 Windows/macOS 是同一目录 | `g9.js` `key` 比较、§6.3 行括注 |
| F9 | 常量具名（`TMP_ROOT_PREFIX` / `TMP_ROOT_KEEP` / `TEMP_SHAPE_RE`）而非字面量散落 | 可读性 + 变异测试需要唯一锚点 | `g9.js` |
| F10 | 自测用「读 `g9.js` 源文本 → 字符串变异 → `Module._compile` 内存编译补丁模块 → 断言源文件字节同一」手法，负例编号续用 `N5`/`N6` | 沿用 T71 既有手法（OB-17 恢复断言），编号续排 | `selftest.js` T80–T83 |
| F11 | 准则侧补一个极短括注（`（**大小写不敏感**）` / `(**case-insensitive**)`）而非新增段落 | 让散文与代码一致，避免新增节号 | §6.3 G9 行 CN/EN |
| F12 | G9(2) 的位置守护**只覆盖 `tmp/` 根**；其它 gitignored 根（`/out/`、`/bin/`、`vendor/`、`/tools/tmp/`、`.vscode/`、`.idea/`、`.github/context-mode/content/`）的强入**不在 G9 覆盖内** | 覆盖面的取舍，且在准则 §6.3 G9 行与 `g9.js` docstring 双处明写，不构成沉默缺口 | §6.3 G9 行、`g9.js` docstring |

## 2 连带机械同步（非自由度，登记备查）

- **M4 变异的 needle** 由 `if (p !== TMP_ROOT_KEEP) ...` 同步为 `if (key !== TMP_ROOT_KEEP) ...`——F8 改了实现，不同步则 needle 失配、`applied=false`，T82 M4 与 F1 stale-mutation 会双红。变异动作与被断言结论不变。
- `selftest.js` T80–T83 注释块的负例描述由「four negatives」改为「six negatives (N1-N6)」。

## 3 越界自查（主控声明）

- 未改 §2.2 / §2.4 本体闭集、§5.2 矩阵、§7.5 数值、G1–G8 判据本体。
- 未触 `internal/**`、CONTRACTS、schema、fixtures、workflows、角色集与载体。
- 未推 gitee；未提交、未推送（等待 ⑩ 同轮授权）。
- 未擅自改变 OB-98 状态词（`已处置` 系用户同轮裁定）。

## 4 已知偏离（如实登记）

1. **本清单产出时点晚于 ⑦ r1/r2**：§1.7 要求「⑥ 之后、⑦ 之前」落盘，实际在 ⑦ r2 通过后才产出（⑦ 输入包内无本清单）。补救：以 ⑦ r3 **收窄复审本清单**闭合（占用 §7.5 的 ⑦ 第 3 次额度）。
2. **② 整改轮误删 `tmp/G9-C/`**：⑧c 才应清理，执行角色提前清理导致 ⑤ r2 输入包与 ③ 原始读数丢失，已按原记录重建等价副本并在报告 §5 登记（`03-tester-report.md` 头注自证）。
3. **⑤ 配额 1+1 用尽后**，② 对 `selftest.js` 注释的 F-A 第 4 处收窄改由主控**确定性自证**（`git grep` 零命中 + `SELFTEST` 全绿 + 三作用域 `TOTAL_FAIL=0`），未再开第 3 次 ⑤。
