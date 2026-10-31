# P41-T027B-SPAWN-OWNERSHIP ⑦ 终审归档（原样）

切片：`P41-T027B-SPAWN-OWNERSHIP`。环节：⑦ 终审。载体：`independent-reviewer`；模型：`GPT-5.3-Codex (copilot)`。输出原样归档。

---

**终审结论**
本次 ⑦ 终审（只读）未发现 Medium 及以上问题；①→⑤→⑥ 证据链在实现语义上可闭合，发布可接受。

**抽查复核（⑤/⑥ 结论，>=3 处）**
1. 创建期归属与 resume 顺序：实现与 ①/⑤ 结论一致。StartManaged 先 attach/identity，再 OnSpawn，再 resume（process.go:190/:224/:232；process_windows.go:259-267/:496-509）。
2. managedCommand 取代 os/exec 后的句柄与退出码语义：单 reaper 写 exitStatus，waitDone 后读取，语义自洽（process.go:102-112/:316；process_windows.go:302-311；process_linux.go:72-75）。
3. RecoverStarting 误删/误杀修复后路径：重开+重锁+重解码在位；⑥ 的两条 Low 与 ⑤ 的 Low-1/Low-2 同源且一一对应（pinned_cli.go:493-530；06-scan L10-11；05-prereview L106-107）。

**逐条发现**
| 编号 | 严重度 | file:line | 一句话问题 | 建议修法 |
|---|---|---|---|---|
| F-01 | Low | pinned_cli.go:350-359 / :450-457 | RecoverStarting 将「无 marker」和「活宿主拒绝」同态返回，且空 launchId 槽位回收分支按 mirror 判定直接删 staging | 后续片：拒绝态哨兵错误或三态返回；空 launchId 仅在 owner 可证明匹配时删除 |
| F-02 | Low | pinned_cli.go:493-530 | clearAgentRunnerStartingMarker 按 path 删除存在极窄 rename 覆盖窗口 | 后续片：对象身份绑定删除（rename-away + 再校验） |
| F-03 | Info | pinned_cli.go:414-416 + 两个回收测试 | 恢复成功路径的 mirror 删除分支无显式断言 | 后续片补断言 |
| F-04 | Info | process_linux.go:44-48 + unix 测试 | F-2 Unix 测试在 Linux CI 语义成立；实现未用 Pdeathsig，Linux 崩溃回收仍以下次扫描为主 | 维持现状可发布；缩短延迟另立片 |
| F-05 | Info | 05-prereview L110 + MANIFEST/SHA256SUMS | D-1..D-6 偏差台账与证据包元数据仍是骨架态 | ⑧c 前补齐并重算 |

**阻断判定**
- 2 Low + 3 Info 均为非阻断；未发现 Medium+。
- 发布可接受性：通过。

FINAL REVIEW: PASS
