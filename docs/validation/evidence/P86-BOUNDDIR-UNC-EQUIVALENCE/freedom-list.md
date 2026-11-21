# P86 BOUNDDIR-UNC-EQUIVALENCE 执行自由度清单（freedom-list，步骤 0 骨架）

- 切片：P86 `BOUNDDIR-UNC-EQUIVALENCE`（`[SLICE]` · 实现类 S 档）
- 终落盘路径（台账指定）：`docs/validation/evidence/P86-BOUNDDIR-UNC-EQUIVALENCE/freedom-list.md`
- 纪律（OB-48 / v1.47 载体校正）：⑥ 后、⑦ 前强制补全并纳入 ⑦ 输入包；
  本文件为步骤 0 骨架，随片推进补条目。
- 每项须标注 F 编号、内容、判定（in-bounds / out-of-bounds + 依据）。

## 条目

| 编号 | 自由度条目 | 判定 | 依据 |
|---|---|---|---|
| F-01 | 主控保持 DeepSeek V4.1 Flash（High），不继承 P84 的 V4 Pro 切换 | in-bounds | 起跑消息 〇（2026-10-31 长期裁定默认；不满足升级触发） |
| F-02 | ⑦ 额度 = 1 终审 + 1 复审（§7.5 S 档字面，**无额度例外**） | in-bounds | 用户 2026-11-22 裁定 a（四要素：裁决人 = 用户、日期 = 2026-11-22、理由 = S 档字面、依据 = 起跑消息） |
| F-03 | ⑦ 档位 = `GPT-5.3-Codex (copilot)` 常规档 | in-bounds | 起跑消息 〇（S 档字面，无档位例外） |
| F-04 | ⑤ 额度 1 + 1（DeepSeek V4 Pro） | in-bounds | 起跑消息 二 |
| F-05 | ⑥ 额度 ≤ 3（MAI Code 2 Flash） | in-bounds | 起跑消息 二 |
| F-06 | ② ③ 用 DeepSeek V4.1 Flash（High），与主控同档 | in-bounds | 起跑消息 〇 / 二（角色分离以独立派单保证） |
| F-07 | 批量辅助 GLM-5.3-Flash ≥ 2 次、严格串行 | in-bounds | 起跑消息 二（仅限枚举 / 摘录核对 / 双语镜像 / MANIFEST 草案，禁择优） |
| F-08 | 先 P86 后 P85 的排片顺序 | in-bounds | 用户 2026-11-22 裁定 b |
| F-09 | 盲审是否必抽由 ① 启动前判定并落盘理由 | in-bounds | 起跑消息 〇 / 二（§7.7） |
| F-10 | Ubuntu VM 不开启；⑨ = Windows 本机 + CI 双腿 | in-bounds | 起跑消息 三 |

（后续环节按实补登；⑧b 定稿落盘。）
