# GATES-TIGHTEN 自由度清单（§1.7.3）

> 路径固定：`docs/validation/evidence/GATES-TIGHTEN-freedom-list.md`（§1.7.3）。
> 落盘时机：⑥ 之后、⑦ 之前（本片 OB-48 已把该时序写入 §5.1）。
> 判定主体：⑦ 终审逐一复核 in-bounds / out-of-bounds。

| 编号 | 自由度行使 | 判定 |
|---|---|---|
| F-1 | 形状对齐的方向分配（7 处漂移中哪几处落 CN、哪几处落 EN）由 §11.1 的「只补 / 删反引号」与 G3(a.1) 对称共同约束 | in-bounds |
| F-2 | G6-5 端点自检的报错形态选 `UsageError`（退出码 2）而非 FAIL | in-bounds |
| F-3 | G6-6 语言声明窗口取常量 `LANG_HEADER_WINDOW`，不写入准则正文数值 | in-bounds |
| F-4 | G7-e 的取值识别边界（不可识别前缀不判）与排除标记集沿用 G6 五词 | in-bounds |
| F-5 | G5-c 从文档自身读取 §2.2 限定名列，不硬编码模型名单 | in-bounds |
| F-6 | 文本类并入项的落点（§5.1 / §7.2.2 / §7.3 / §9.3 / §11.2 / §12.5 / §6.3 两处） | in-bounds |
| F-7 | 自测探针的构造方式（合成 ctx 与内存补丁，不新增夹具文件） | in-bounds |
| F-8 | 附录 C.0u 块的字段集合按既有 C.0x 块同形选择 | in-bounds |
| F-9 | `tmp/gates-tighten/` 下的过程探针脚本集合 | in-bounds |

