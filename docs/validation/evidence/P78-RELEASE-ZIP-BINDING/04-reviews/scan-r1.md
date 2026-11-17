# ⑥ 独立扫描 r1（MAI-Code-1.1-Flash，只读）

切片：P78 `RELEASE-ZIP-BINDING`。输入包：方案书 + 落库 diff（28 文件）+ tree 门禁原始输出 + 复验证据路径 + 报告对 + Release notes 草稿（不含 ⑤ 清单）。

## Section A
INDEPENDENT SCAN: PASS

## Section B
| # | SEVERITY | file:line | 被违反的句子或判据 | 最小修复 |
| - | none | - | - | - |

## Section C
- 一致性：四元组在 `docs/validation/release-zip-binding.md`（L50-L53）、`docs/validation/release-zip-binding_EN.md`（L50-L53）、`docs/validation/s1-exit.md`（L66）、`docs/S1_SUPPORT_MATRIX.md`、`docs/S1_RELEASE_NOTES.md`（L5）与 `tmp/P78-RELEASE-ZIP-BINDING/10-build/release-notes.md` 中保持一致：版本为 v0.1.0，candidate commit 为 c93571352f456a275d6e71ce3e084c9f96360e86，run 为 44043705691，文件名为 prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip。没有出现四元组漂移或同值不一致。
- 越界：变更集 `tmp/P78-RELEASE-ZIP-BINDING/impl.diff` 与 `impl.numstat` 只覆盖 28 个文档型变更；没有触及产品代码、CI 工作流、p57 门表、s1-exit 门表、准则、CONTRACTS、schema、fixtures、其它 ADR、S1_SUPPORT_MATRIX 平台支持层级、代码签名、Linux 正式发布。对于「S1_SUPPORT_MATRIX 平台清单与支持层级」这条，实际变更是新增「已绑定发行记录」说明，不是改动平台名单或支持层级。
- 事实与证据：报告中的关键事实在 `tmp/P78-RELEASE-ZIP-BINDING/20-reverify.log` 中得到直接支撑：安装、升级、回滚、卸载保留数据四项均为 PASS。门级原始输出 `tmp/P78-RELEASE-ZIP-BINDING/04-gates-tree-pre.txt` 里 G7-b/G7-f 仅是时序红因（未暂存工件 / 围栏注入），并非实质缺陷。没有出现「未验证却写成已验证」的模式。
- 门级表述：有明确的硬边界，且与 ADR-046 的两步审批一致。例如 `docs/validation/release-zip-binding.md` 与 `docs/S1_RELEASE_NOTES.md`（L5）写明「发行部分已批准，但不构成 S1 Exit 整体通过」，以及 `docs/ADR_REGISTER.md`（L62）与 `docs/S1_SUPPORT_MATRIX.md`（L5）中保留「当前仍不是已发布 SLA」。因此没有将「候选发行」写成「S1 Exit 整体通过」或「已发布 SLA」。

## Section D
- 该结论在什么观察下会被推翻：
  1. 若同一四元组在 `docs/validation/release-zip-binding.md`、`_EN.md`、`docs/S1_RELEASE_NOTES.md`、`docs/S1_SUPPORT_MATRIX.md` 与 `tmp/P78-RELEASE-ZIP-BINDING/10-build/release-notes.md` 中出现不一致。
  2. 若 `tmp/P78-RELEASE-ZIP-BINDING/impl.diff` 发现了产品代码、CI、门表、CONTRACTS、schema、fixtures 或 Linux 正式发布相关越界。
  3. 若 `tmp/P78-RELEASE-ZIP-BINDING/20-reverify.log` 不再显示 install / upgrade / rollback / uninstall 全部 PASS。
  4. 若报告或 Release notes 改成「已发布 SLA」或「S1 Exit 整体通过」。
  5. 若实际远端的 GitHub Release 资产与本地证据发生不一致，或相关文件在后续修订中篡改为与当前证据不符。

## Section E
- 我无法验证的部分：没有直接抓取远端 GitHub Release 页面与发布资产；本扫描基于仓库内证据、日志和文档构成的本地一致性证明。就「仓库内是否符合四元组 / 边界 / 事实」，本地证据足够；就「远端公开页面始终保持一致」，此时点无法做外部实时确认。
- 什么会推翻 PASS：只要本地证据变成与真实远端发布页面或资产不一致；或后续出现新增变更扩大到禁止触碰范围；或四元组 / 门级表述变成「已发布 SLA / S1 Exit 整体通过」，PASS 即应撤销。

INDEPENDENT SCAN: PASS
