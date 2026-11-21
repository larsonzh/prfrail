# P86 ⑤ 预审 r2（整改复审）逐字转录 · 2026-11-22

- 载体：`independent-reviewer`（只读）；模型 `DeepSeek V4 Pro (deepseek)`；额度 1 + 1 的第 2 次。
- 输入：r1 发现 + Rev 2 整改（方案书 §10 / 代码注释与守卫收紧 / 分类器负例）。
- 消毒声明（§7.3）：输入包不含凭据 / token / 私密主机名 / IP。

---

PRE-REVIEW: PASS

**闭环表**

| r1 项 | 状态 | 依据 |
|---|---|---|
| Medium：生产入口验收判据不分 PASS/SKIP | 已闭合 | §10 Rev 2 明确判据：⑨ 证据须含 `--- PASS: TestP86UncShareRootProductionEntry` 字面行；SKIP 则逐字登记原因。代码侧 3 个 skip 点（`net share` 失败 / `OpenBoundDir` 失败 / base 拼写前置断言）均 `t.Skipf` 写明原因，PASS 与 SKIP 可区分 |
| Low：`boundRelativeSuffix` 注释矛盾 | 已闭合 | 注释已补「P86 (N13-UNC) adds ONE exception…」段，明确后代/盘符/平面目录孪生仍 Ordinal |
| Low：`uncShareRootName` 放行 `\\?\`/`\\.\` | 已闭合 | 注释改为「volume name of the ordinary DOS form (\\host\share)」；`strings.HasPrefix(vol, `\\?\`) \|\| …(`\\.\`)` 早退已落地 |
| Low：文档同步清单缺项 | 已闭合 | 00-briefing §七-5 已列 `SECURITY{,_EN}.md:5` 为必改；bounddir-residual §8 按 §七-7 归为核对项，Rev 2 登记「待用户裁定、⑩ 回执上报」，处置路径明确 |
| Low：变异判据措辞含糊 | 已闭合 | 测试注释与 Rev 2 一致：恰好 3 拼写差异子例变红；`identical-trailing-spellings` 与 `strict-descendant-regression` 两世界同绿 |

**新风险核验（`\\?\`/`\\.\` 早退可达性）**：无误拒。
- base：`currentName()`→`handleFinalPath`，flags=`volumeNameDos|fileNameNormalized`=0（:32-34），API 必回普通 DOS 形态，无 `\\?\`；`stripExtendedPathPrefix` 空转。
- target：`reparseStoredTarget`→`normalizeReparseTarget` 剥 `\??\`/`\GLOBAL??\` 后 `stripExtendedPathPrefix` 已将存储形态 `\\?\UNC\…`→`\\…`、`\\?\…`→普通；`\\.\`/GLOBALROOT 设备形态不产出普通 DOS，被 IsAbs/前缀守卫先行拒绝。
- 故早退仅对合成输入可达；且它堵住一真实退化洞——此前 `uncShareRootName(`\\?\UNC\`)` 会返回 `\\?\UNC\`（`VolumeName` 恰取前缀），现显式返 `""`，属纵深防御。
- 实测读数与代码自洽：5 正例/9 反例+8 分类器负例（含 4 个新形态）/3 个 `uncShareRootEqual` 断言均与实现一一对应；`gofmt -l internal/evidence` 实测空输出。

**未决 Medium+**：无。两条执行期注记（⑨ 须贴 PASS 字面行或 SKIP 原文；bounddir-residual §8 待用户裁定）为 Low 级流程项，不构成代码发现。

END OF PRE-REVIEW r2
