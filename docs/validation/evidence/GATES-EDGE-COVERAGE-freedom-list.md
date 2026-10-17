# GATES-EDGE-COVERAGE 自由度清单（§1.7.3）

> 路径固定：`docs/validation/evidence/GATES-EDGE-COVERAGE-freedom-list.md`（§1.7.3）。
> 落盘时机：⑥ 之后、⑦ 之前（本片 OB-48 已把该时序写入 §5.1）。
> 判定主体：⑦ 终审逐一复核 in-bounds / out-of-bounds。

| 编号 | 自由度行使 | 判定 |
|---|---|---|
| F-1 | OB-90 的判型选择：单侧变更以既有的 `G3(a.1)` FAIL 承担（不新增判据 ID、不改 `gate.js` 注册） | in-bounds |
| F-2 | OB-90 的配对面保持 `changedPaths`（G3 本就不在证据树排除名单内，面不变） | in-bounds |
| F-3 | OB-90 的「不可读侧」选 FAIL 而非 `UsageError`（删除是合法数据态，非判据前提损坏） | in-bounds |
| F-4 | OB-91 的续延字符集取 `[A-Za-z0-9._-]`（终结符 = 其补集，属封闭列举） | in-bounds |
| F-5 | 两处修法均复用合成探针、不新增磁盘夹具（OB-17 先例） | in-bounds |
| F-6 | 变异的载荷证明用 `loadMutatedModule` 的内存补丁（不落盘改写、可断言 byte-identical） | in-bounds |

