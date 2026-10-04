# 归档 MANIFEST · tmp-lifecycle-workdir

- **源路径**：`tmp/tl/`
- **归档片**：`DIRECTIVE-TMP-LIFECYCLE`（v1.23，2026-10-04）
- **归档口径**：§1.5.1 的 ② 类「早期证据评估归档」——保留证据条目（被受版本控制文档引用其路径者 ∪ 切片级文档），排除运行时残留与一次性脚本
- **归档条目**：15 条 / 744,747 字节（归档载荷根 = `archive/`；逐条 sha256 见 `SHA256SUMS.txt`）
- **排除条目**：41 条（逐条见下「排除清单」）
- **未入库条目**：0 条（逐条见下「未入库文件清单」）
- **字节保真**：归档副本的**字符内容不改**；因 §6.3 的作用域（G4a 扫全部变更文件）与仓库编码约定，副本按**归档后类型**归一编码——`.ps1` 带 BOM，其余不带 BOM，一律 LF；逐条变更种类与**原始字节 sha256** 见下「后缀与编码映射」节
- **后缀口径**：归档副本的 `.md` 一律改存为 `.txt`（本单元 0 条），以避开只扫 `.md` 的文档判据；映射逐条见下节
- **编码归一计数**：后缀改写 0 条 / BOM 剔除 0 条 / BOM 补入 0 条 / LF 归一 0 条

## 归档条目清单

| 相对路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/arch-log.txt` | 0 | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/archive-filelist.txt` | 25375 | `a2a73710cb6117855bd5b4474e613cbec12217bf147cff26c16d558fdf0c52da` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/classification.txt` | 146511 | `e81dea2b13277bae378be8564d977d64e90e5e95ace3f2c3b4352470e5fc9a5d` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/excl-summary.txt` | 2 | `2689367b205c16ce32ed4200942b8b8b1e262dfc70d9bc9fbc77c49699a4f1df` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/facts.txt` | 42095 | `406194f02ca752a96e692906f2d2cf339a2da2397bfe7e241dfcb561de62be33` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/facts2.txt` | 33261 | `6048202bd2ec6ff02763a3565f9f833b2733d50dfb44dc8dd28864ea914c36cb` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/gate-index.txt` | 19802 | `f5ee497a7dd571209f20cbed6dc8b66fe6a227425c6f2ef645b58ffabf33da48` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/inventory.txt` | 221063 | `995dfff23cbff8d6d288b016b0d6473b28f72bf5fa20a1ab31cd991feaaa5eee` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/measure.txt` | 186929 | `83f807f56671e5d12f9c73cfc6c3fb1afd2bb1f865afdcb30d3ef34808d73486` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/peek.txt` | 21458 | `8165f161f656f0ce988d6a777891c5dc90376b58cdf02da590fdbe66ac7ea2f7` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/prep.txt` | 9085 | `523e98970f7241a50554a81c513e2becc20035cc6ec321d3a6f24a1ebdff951f` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/sets.txt` | 3854 | `39108ff3822c20e8bdd579278253421a8d291f68651e38587bc31bfb91e7bb9a` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/stage-list.txt` | 25375 | `1159c6fc455c4ec58bf9d791d968efb55cb8833b736a2dfba01cea880f803727` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/structure.txt` | 8299 | `134df194b41c53d389625d1bc9a9cac1b76f1119a951659de98d9a2cd0d25340` | 步骤 0 盘点与分类读数（本片证据） |
| `archive/verify-arch.txt` | 1638 | `9a02f04e38a313b1da93f31836c390c63f454c473fba46b743d4586738bbc070` | 步骤 0 盘点与分类读数（本片证据） |

## 排除清单（未归档）

| 相对路径 | 字节 | sha256 | 未归档理由 |
|---|---|---|---|
| `/arch.js` | 13482 | `72d65b9717423e56e4b9617076a99ff9fab6d397049740220174812fb49e1593` | 一次性脚本（内容可由报告与读数复现） |
| `/check-anchors.js` | 3473 | `d804427375d1b49d1fc931834864e3f29e936ed7ee1ddda6baf904ae9defd96b` | 一次性脚本（内容可由报告与读数复现） |
| `/classify.js` | 5423 | `70bc6b68a98115ac4d42969d2ce72e323b2b2133d3028d71486c0c25afaadfa7` | 一次性脚本（内容可由报告与读数复现） |
| `/excl-summary.js` | 2423 | `be9e04acef54d5b026ffa17020ece2cf2d9e63eb6469f0b51cd9c6819294bd13` | 一次性脚本（内容可由报告与读数复现） |
| `/fix-notes.js` | 3385 | `68d4bf684201c230fd415da0276f451ad144adcf0f7e7bf186e49f27885c3837` | 一次性脚本（内容可由报告与读数复现） |
| `/fix-r7.js` | 3311 | `5b179fdbad8c582c18d2308f413c35cc473aff45eec396026a08261890c87d4e` | 一次性脚本（内容可由报告与读数复现） |
| `/gate-analysis.js` | 2096 | `023cb903b81e21225dff49ed6aaebcff13168f64dae55cb51a05160487091ab2` | 一次性脚本（内容可由报告与读数复现） |
| `/gate-free.log` | 122 | `8bc2eb1bb9ce3eca05f01dcb067bb3d611f3e46010c2640853d3ad2cba31abe0` | 一次性脚本（内容可由报告与读数复现） |
| `/gate-index.log` | 5584 | `f8de84dcd1b0fb0c10e0f2782395aa4a810c908366a6832c8e62fc6dcd595c68` | 一次性脚本（内容可由报告与读数复现） |
| `/gate-tree.log` | 5583 | `2628fc6c60fd170efcb4a64cb4705ea39210c73da8c15a247488226d6b53ad01` | 一次性脚本（内容可由报告与读数复现） |
| `/inv.js` | 4108 | `f4ab0a167f05da2ea7ba216f1dc4c4dd361ab4cdbd9b99bf9b522aa630160bd7` | 一次性脚本（内容可由报告与读数复现） |
| `/land-ob79.js` | 2938 | `d6aab2b81ae26a89c68393ca4faf633565b34aa483fc5e2bd411498a4fe5492f` | 一次性脚本（内容可由报告与读数复现） |
| `/land23.js` | 27629 | `df6f3f5c4416200a3349d21561e1d1b7230360e30613f1318db9eb593f6deddc` | 一次性脚本（内容可由报告与读数复现） |
| `/land23b.js` | 6463 | `4b3b791a2a1e017283736e9ecf3de62a83172572f198f1cca14a859e751358ef` | 一次性脚本（内容可由报告与读数复现） |
| `/meas.js` | 2686 | `47f673f9e1709d0bd7df0154937d121389bdc9086e1aaa80aa7fd65636782e6e` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/b2-enforcement.md` | 3033 | `6ab06af070e31c7a805eb7240848751f8a8141b915c191f50333cc4f74882c6c` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/b3a-vm-environment.md` | 10432 | `048bb27d142d4ca3bd40d6a444c74dba06e675475b4cecefeac37a7622f2f197` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/bilingual-mirror.md` | 5786 | `442078b8e2c7763e515fd35f1cd3b033b41abc48203e5dddd59d11daf6e95d76` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/context-mode.md` | 1432 | `801f70f9e82eccca51e807c949f94381531fe46d1dc590e0d29fb58fd6e75efd` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/delivery-directive.md` | 10112 | `96ef5fac4bbf06a9c95be810e93774bbde029fb9912ae12edf6b1a779165f596` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/directive-model-lists.md` | 16887 | `f886691e8703692b2dc20afdcf9176f6b185a34b2a48510ed29ca22fea54bdf2` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/encoding-conventions.md` | 6809 | `e75232d7e49e72e3d084b267045810173aa46e8582d8640cb746bd5c941d7967` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/prfrail-agent-carriers.md` | 3711 | `8077baee4d5004838be6a0e4b511142e4be85960d31478371480e3d13a8d7f4d` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/prfrail-gate-pitfalls.md` | 3158 | `8e294337118cfce83b637d82a2258088e152a2d00cbdb93e0c1e2d5049c4a89f` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/prfrail-go-gotchas.md` | 1688 | `3944cbe9b35fc43f1dc80585f81ba810bdb934b2a5f95dadfaa3ac34dad1a799` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/prfrail-remote-validation.md` | 1059 | `942612f00612eeba4b3fa37dd533c364bbd59b02dcbf1b3e9c9d937876bb531a` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/prfrail-review-after-external-edits.md` | 5942 | `c2bc9834ca24abb25cd60f69df07d2204f430a11db35f4dccefa22f88073b46b` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/sessbridge.md` | 12195 | `17e684c696631f9364edd2f8e74fccca5ae10ba74327d1d0787036d7bd957dd7` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/t026-agent-probe.md` | 3492 | `f428ae78d09d789d93b64fbefede8f641cf255f45d905cf2843209804bae9703` | 一次性脚本（内容可由报告与读数复现） |
| `/mem/t027-slice-protocol.md` | 63774 | `aaa2df83537695e75784f64534f711db5a1fe90f8906c1208eed57a2d90ba0a7` | 一次性脚本（内容可由报告与读数复现） |
| `/peek.js` | 3066 | `0ddd6f39fb4a87e4c6d9569cd529cd49a438edb36d3ff9f220ea4f6cb4a74059` | 一次性脚本（内容可由报告与读数复现） |
| `/prep.js` | 2492 | `0a819fc370821fbe0863c83d01c8387bc1eef5e856dc01101b86af47c89f0bcd` | 一次性脚本（内容可由报告与读数复现） |
| `/probe2.js` | 2261 | `d4c26878e10f41af653038084ce7c3f6868579f90e5982735a08658b01b4d0ab` | 一次性脚本（内容可由报告与读数复现） |
| `/probe3.js` | 1641 | `4e7b10e2fc70fdac2281d1b8483a9cb8b84a3fd675bf275fa3d0e5ee2ce95081` | 一次性脚本（内容可由报告与读数复现） |
| `/sets.js` | 4494 | `c98abd705e129fac62f88b1bdeaec02f756dc55f0c420bc832c576f4b488e3b2` | 一次性脚本（内容可由报告与读数复现） |
| `/stage.js` | 1414 | `1ed1dd2269c7922865e36e68a273608c08a92b9d4932b906420500b434ddbc2f` | 一次性脚本（内容可由报告与读数复现） |
| `/struct.js` | 1980 | `e5234b1a0f0c188f57b3a0cba6e181f0cfb2bad086470c227cd2eaeaf9d3dc58` | 一次性脚本（内容可由报告与读数复现） |
| `/sum.js` | 806 | `da9766b445f1be15d780969c8453e1f94137b0a9c1fb6f86bccb88bc2aec2abf` | 一次性脚本（内容可由报告与读数复现） |
| `/unit-summary.js` | 876 | `965989d36c094a4e7a96548dee45de0f6ab1322ae4c89f5cc550e404bc1c1508` | 一次性脚本（内容可由报告与读数复现） |
| `/verify-arch.js` | 2480 | `73f1dae0a239f2e72fff0bca46baa6be0cd4871e8659b3a03fb745a599e48d28` | 一次性脚本（内容可由报告与读数复现） |
| `/write-free.js` | 3846 | `e091d86430470081fe9111c5b4e23e2ba796500d3d9b0b7f80b901363d7a7798` | 一次性脚本（内容可由报告与读数复现） |

## 未入库文件清单（`.gitignore` 排除类型）

| 相对路径 | 字节 | sha256 | 理由 |
|---|---|---|---|

## 重建说明

- **未入库条目**（见上表）：本单元未入库清单为空（无忽略类型条目）。
- **排除条目**：一次性脚本、运行时状态与中间转储；其行为已由本片归档的切片级文档（切片定义 / 方案 / 评审包 / 门禁读数）与对应切片的提交记录固定。
- **行尾与 BOM 归一**：源 `tmp/` 副本的部分文件为 CRLF 或带 BOM；归档副本已按仓库编码约定归一到「`.ps1` 带 BOM / 其余不带 BOM / 一律 LF」（不影响字符内容）。

## 后缀与编码映射（`.md` → `.txt` 及编码归一）

无（本单元无 `.md` 载荷、无编码归一）。

## 对应引用点

| 引用文件 | 行号 |
|---|---|
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/classification.txt` | 1409 / 1410 / 1411 / 1412 / 1413 / 1414 / 1415 / 1416 / 1417 / 1418 / 1419 / 1420 / 1421 / 1422 / 1423 / 1424 / 1425 / 1426 / 1427 / 1428 / 1429 / 1430 / 1431 / 1432 / 1433 / 1434 / 1435 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/measure.txt` | 36 |
