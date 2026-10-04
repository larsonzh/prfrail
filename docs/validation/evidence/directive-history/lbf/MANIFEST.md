# 归档 MANIFEST · lbf

- **源路径**：`tmp/lbf/`
- **归档片**：`DIRECTIVE-TMP-LIFECYCLE`（v1.23，2026-10-04）
- **归档口径**：§1.5.1 的 ② 类「早期证据评估归档」——保留证据条目（被受版本控制文档引用其路径者 ∪ 切片级文档），排除运行时残留与一次性脚本
- **归档条目**：12 条 / 273,553 字节（归档载荷根 = `archive/`；逐条 sha256 见 `SHA256SUMS.txt`）
- **排除条目**：9 条（逐条见下「排除清单」）
- **未入库条目**：0 条（逐条见下「未入库文件清单」）
- **字节保真**：归档副本的**字符内容不改**；因 §6.3 的作用域（G4a 扫全部变更文件）与仓库编码约定，副本按**归档后类型**归一编码——`.ps1` 带 BOM，其余不带 BOM，一律 LF；逐条变更种类与**原始字节 sha256** 见下「后缀与编码映射」节
- **后缀口径**：归档副本的 `.md` 一律改存为 `.txt`（本单元 6 条），以避开只扫 `.md` 的文档判据；映射逐条见下节
- **编码归一计数**：后缀改写 6 条 / BOM 剔除 6 条 / BOM 补入 0 条 / LF 归一 2 条

## 归档条目清单

| 相对路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/dd36b9b-trace.txt` | 71 | `0743f0515b82c5a04d257c474c86f349e7908aa40112cc7a791fb64b275bb57c` | 被受版本控制文档引用其路径 |
| `archive/dispositions.txt` | 4781 | `379c2593bab3bb61d4f87ba03c9964bfe38f53d0e97208c7af50fe0797f9fcbf` | 被受版本控制文档引用其路径 |
| `archive/gate-index.txt` | 88 | `1b28f476963371ebbdbce873324f30500f74468df31024d53c5f8490becf59ca` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree.raw.txt` | 3553 | `baa1fdc6ddc672ebe3563862f7fe098f6c0abe2cb17e6fe5877096ec91d81dd6` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree.txt` | 3804 | `1cf6fd85f9b4fef2fbd91f2e739f8f2a96322e20d095e8cba562985e7fe2ee7b` | 被受版本控制文档引用其路径 |
| `archive/head-adr-cn.txt` | 19902 | `bfd3d719e7335060b43594694bedd1d213517da3cce2a9465962f76eb65a79e4` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/head-adr-en.txt` | 22656 | `c9e7fd4094c3f183424b6fadead2ea466bea5998c844fce8fc45874ef9c91d6c` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/report-cn.txt` | 16348 | `f064567fd715d6a0aa20a83f76882c897cd62b317a7729a45b2d29af264211be` | 被受版本控制文档引用其路径 |
| `archive/report-en.txt` | 19366 | `0e4f4a088fd8b26a72f87f6510d6da081cfaf4c78819d9dc4563738b27b41241` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input.diff` | 155330 | `e775e2205fbc23070d8c1e067375473bc709dd69e0bff513d3b15467020702a9` | 被受版本控制文档引用其路径 |
| `archive/review-round2.txt` | 6010 | `eda2111827428a3242b90f4f8229b297f7515cbd8cf6334fd6391e6d5a00a62f` | 被受版本控制文档引用其路径 |
| `archive/selftest.txt` | 21644 | `a75e7b7a10a6e0f52a15f6f6b112de7026b391a9a82aab7c890619317aa484ff` | 切片级文档 / 评审包 / 门禁读数 |

## 排除清单（未归档）

| 相对路径 | 字节 | sha256 | 未归档理由 |
|---|---|---|---|
| `/apply.js` | 12296 | `0bbcadc51efb107eaef167d58089e4f621901f262a2a6cabbcbbd0bdd38d47e8` | 落地与修复脚本残留 |
| `/apply2.js` | 4341 | `119cc146a4ca24df2a38c3aeb5432f822f76089ffea9c9974420e6bbf3f781f6` | 落地与修复脚本残留 |
| `/ci-watch.txt` | 4467 | `9939402fb16bf96287c1aabb062041220c2cbc0134986990487fe514b24e0b6e` | 中间转储（未被引用的过程输出） |
| `/findchar.js` | 1219 | `3d29652133153f9a13593a648289558221ab6e3d7b35216613bc1e542a4c25e2` | 落地与修复脚本残留 |
| `/fix-adr.js` | 1514 | `43fc613fbbefdedc4094e88ac0c95f811c57812c6bc321e1fc028174115aee07` | 落地与修复脚本残留 |
| `/mkpack.js` | 1355 | `04a39620f5597fdccab3e62d60034625896a2c071221986dfb563f5e50ff1ae4` | 落地与修复脚本残留 |
| `/probe-callshape.js` | 886 | `c99722b2c4e0bfe9b8983ad8780b6643f04218e7a63bb7d9e9e15c5ef33deedb` | 落地与修复脚本残留 |
| `/show.js` | 1411 | `08851aa3724e250b9b510962e227a2f7764c642e6cdd3cde7861b4a50b9a3f82` | 落地与修复脚本残留 |
| `/writeback-8b.js` | 8456 | `e870a81408ad4a97918a1e407045fcfada77d014f3ebbd58847d0192c2209659` | 落地与修复脚本残留 |

## 未入库文件清单（`.gitignore` 排除类型）

| 相对路径 | 字节 | sha256 | 理由 |
|---|---|---|---|

## 重建说明

- **未入库条目**（见上表）：本单元未入库清单为空（无忽略类型条目）。
- **排除条目**：一次性脚本、运行时状态与中间转储；其行为已由本片归档的切片级文档（切片定义 / 方案 / 评审包 / 门禁读数）与对应切片的提交记录固定。
- **行尾与 BOM 归一**：源 `tmp/` 副本的部分文件为 CRLF 或带 BOM；归档副本已按仓库编码约定归一到「`.ps1` 带 BOM / 其余不带 BOM / 一律 LF」（不影响字符内容）。

## 后缀与编码映射（`.md` → `.txt` 及编码归一）

| 原路径 | 归档路径 | 原始字节 sha256 | 归档字节 sha256 | 变更 |
|---|---|---|---|---|
| `dd36b9b-trace.txt` | `archive/dd36b9b-trace.txt` | `64893e1d488e4e40f944ee55c088212ddf1dfb74034aad2d30c775ba05ed4ddf` | `0743f0515b82c5a04d257c474c86f349e7908aa40112cc7a791fb64b275bb57c` | LF 归一 |
| `dispositions.md` | `archive/dispositions.txt` | `69ce117c7a4d992b14270217d7803b7dd6dc2200a85f29c2742d6371c68b6cb5` | `379c2593bab3bb61d4f87ba03c9964bfe38f53d0e97208c7af50fe0797f9fcbf` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `gate-tree.raw.txt` | `archive/gate-tree.raw.txt` | `8bc2170e063638a91ce7460647f523dee7d19096d2ed43881444705635c81482` | `baa1fdc6ddc672ebe3563862f7fe098f6c0abe2cb17e6fe5877096ec91d81dd6` | LF 归一 |
| `head-adr-cn.md` | `archive/head-adr-cn.txt` | `81cfc822dd85a31bc1c24321fc750e44c905005c69b76eeb6c14fe9144656392` | `bfd3d719e7335060b43594694bedd1d213517da3cce2a9465962f76eb65a79e4` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `head-adr-en.md` | `archive/head-adr-en.txt` | `2c9a0590ab10830833d29adc58ed20aa8a850a1375f4176dad8b9bd18f643137` | `c9e7fd4094c3f183424b6fadead2ea466bea5998c844fce8fc45874ef9c91d6c` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `report-cn.md` | `archive/report-cn.txt` | `1b1737102c2d5c027c73f0beb83d9674ee0066f77e6225eaac0c29370a830a55` | `f064567fd715d6a0aa20a83f76882c897cd62b317a7729a45b2d29af264211be` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `report-en.md` | `archive/report-en.txt` | `d01d45c7a9107cec7e890e84dd842ae566d0976d5b0d64820189f297e07af4d8` | `0e4f4a088fd8b26a72f87f6510d6da081cfaf4c78819d9dc4563738b27b41241` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `review-round2.md` | `archive/review-round2.txt` | `4976a14c41165486c2611a542755d91a597fdc9cacab0de7ded3c76963d5eceb` | `eda2111827428a3242b90f4f8229b297f7515cbd8cf6334fd6391e6d5a00a62f` | 后缀 `.md` → `.txt` / BOM 剔除 |

## 对应引用点

| 引用文件 | 行号 |
|---|---|
| `docs/validation/evidence/directive-history/lbf/archive/dispositions.txt` | 39 |
| `docs/validation/evidence/directive-history/lbf/archive/report-cn.txt` | 32 / 36 / 74 / 75 / 93 / 100 / 101 / 104 / 142 / 143 / 144 / 145 / 146 |
| `docs/validation/evidence/directive-history/lbf/archive/report-en.txt` | 32 / 36 / 74 / 75 / 93 / 100 / 101 / 104 / 142 / 143 / 144 / 145 / 146 |
| `docs/validation/evidence/directive-history/lbf/archive/review-input.diff` | 362 / 366 / 404 / 405 / 423 / 430 / 431 / 434 / 472 / 473 / 474 / 475 / 476 / 515 / 519 / 557 / 558 / 576 / 583 / 584 / 587 / 625 / 626 / 627 / 628 / 629 |
| `docs/validation/evidence/directive-history/lbf/archive/review-round2.txt` | 19 / 22 / 37 / 39 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/classification.txt` | 1372 / 1373 / 1374 / 1375 / 1376 / 1377 / 1378 / 1379 / 1380 / 1381 / 1382 / 1383 / 1384 / 1385 / 1386 / 1387 / 1388 / 1389 / 1390 / 1391 / 1392 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/facts.txt` | 68 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/facts2.txt` | 90 / 91 / 92 / 93 / 94 / 95 / 96 / 97 / 98 / 99 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/inventory.txt` | 174 / 175 / 176 / 177 / 178 / 179 / 180 / 1086 / 1087 / 1088 / 1089 / 1090 / 1091 / 1092 / 1093 / 1094 / 1095 / 1096 / 1097 / 1098 / 1099 / 1100 / 1101 / 1102 / 1103 / 1104 / 1105 / 1106 / 1107 / 1108 / 1109 / 1110 / 1111 / 1112 / 1113 / 1114 / 1115 / 1116 / 1117 / 1118 / 1119 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/peek.txt` | 70 / 72 |
| `docs/validation/ledger-backfill.md` | 32 / 36 / 74 / 75 / 93 / 100 / 101 / 104 / 142 / 143 / 144 / 145 / 146 |
| `docs/validation/ledger-backfill_EN.md` | 32 / 36 / 74 / 75 / 93 / 100 / 101 / 104 / 142 / 143 / 144 / 145 / 146 |
