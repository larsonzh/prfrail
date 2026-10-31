# P44-T027B-GATE-READINESS · 执行自由度清单（步骤 0 基线）

本片为 `[SLICE]` 实现类，规模 M，边界 = `internal/console/**` 与 `internal/gates/**`。

## 可自由裁量
1. 新增文件/类型的命名与拆分（`internal/console/**`、`internal/gates/**` 内），须与既有风格一致。
2. 新增测试文件的组织方式（表驱动 / 单例），须保持 ③ 先行且不得修改既有断言强度。
3. 内部辅助函数形状（不改变导出签名的前提下）。
4. 注释与错误信息措辞（英文，遵循仓库 i18n 约定）。

## 不可动（硬边界）
1. 不得修改 `internal/chain/**`、`internal/adapters/**`、`internal/evidence/**`、`internal/tickets/**`。
2. 契约语义以 `docs/CONTRACTS.md` 为准；若需新增 readiness 三级或 `ProbePermit` 语义，须在 ⑧ 阶段回写契约，且不得与既有 §4.1（封存屏障与 K）冲突。
3. 不得降低既有测试断言强度，不得以 `t.Skip` 转绿。
4. 不得提交/推送（停 ⑩ 等授权）；不得推 gitee。
5. 临时件只落 `tmp/P44-T027B-GATE-READINESS/`，片内只增不删。

## 角色载体（P43 D7 铁律）
- ③ 测试必须由 `prfrail-tester` + DeepSeek V4.1 Flash 派单完成。
- ② 实现必须由 `prfrail-implementer` + DeepSeek V4.1 Flash 派单完成。
- 主控只做编排、落盘、门禁与回执，不得代写产品代码或测试。

## 判定冲突时的优先序
1. `docs/CONTRACTS.md` ＞ 2. 本片 ① 架构方案（含修订注记）＞ 3. 既有测试断言 ＞ 4. 代码现状。
