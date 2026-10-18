# GATES-EDGE-COVERAGE-RECHECK 自由度清单（§1.7.3）

> 路径固定：`docs/validation/evidence/GATES-EDGE-COVERAGE-RECHECK-freedom-list.md`（§1.7.3）。
> 落盘时机：⑥ 之后、⑦ 之前；本片为纯验证片，⑦ 已按 1+2 额度收口。
> 判定主体：⑦ 终审逐一复核 in-bounds / out-of-bounds。

| 编号 | 自由度行使 | 判定 |
|---|---|---|
| F-1 | 只造反例、不改判据：`tools/gates/lib/criteria/**` 与判据文本零改动，缺陷只以自测反例暴露 | in-bounds |
| F-2 | canary 用冻结字面量而非读取源码文本：证明测试不反向定义规则输入 | in-bounds |
| F-3 | S2 期望按 §6.3 与 `g3a.js` 头注修正（非指令对不产生 `G3(a.2)` 记录），不改判据去配合期望 | in-bounds |
| F-4 | OB-97 只登记不修：越界形态属后续修复片，本片不扩大变更集 | in-bounds |
| F-5 | ⑦ 额度 1+2 用尽后不开第 3 读：r2 的 Low 以报告留痕收口，不申请额度例外 | in-bounds |
| F-6 | `tmp/` 夹具用 `git add -f` 入索引跑 `--scope=index`，跑完 `git reset` 复原到基线树哈希 | in-bounds |
