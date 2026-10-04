# 归档 MANIFEST · v112

- **源路径**：`tmp/v112/`
- **归档片**：`DIRECTIVE-TMP-LIFECYCLE`（v1.23，2026-10-04）
- **归档口径**：§1.5.1 的 ② 类「早期证据评估归档」——保留证据条目（被受版本控制文档引用其路径者 ∪ 切片级文档），排除运行时残留与一次性脚本
- **归档条目**：22 条 / 450,619 字节（归档载荷根 = `archive/`；逐条 sha256 见 `SHA256SUMS.txt`）
- **排除条目**：18 条（逐条见下「排除清单」）
- **未入库条目**：0 条（逐条见下「未入库文件清单」）
- **字节保真**：归档副本的**字符内容不改**；因 §6.3 的作用域（G4a 扫全部变更文件）与仓库编码约定，副本按**归档后类型**归一编码——`.ps1` 带 BOM，其余不带 BOM，一律 LF；逐条变更种类与**原始字节 sha256** 见下「后缀与编码映射」节
- **后缀口径**：归档副本的 `.md` 一律改存为 `.txt`（本单元 4 条），以避开只扫 `.md` 的文档判据；映射逐条见下节
- **编码归一计数**：后缀改写 4 条 / BOM 剔除 0 条 / BOM 补入 0 条 / LF 归一 17 条

## 归档条目清单

| 相对路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/append-sec12-cn.txt` | 3149 | `7b0780e287082b8972839540ad71c9211a61b0ef25d977dfebaaf263efde28bc` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/append-sec12-en.txt` | 3574 | `2e63bddda852550ee02ae87450ea4f8771bbc3728ad9df2cc7bfca89269dda34` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/c1-recheck.js` | 3973 | `631f8c6c8cdf3995b0d17f481aba7dab9a4a90c6dade7e9536cb65e944efec50` | 被受版本控制文档引用其路径 |
| `archive/c1-recheck.txt` | 49 | `a79a8faf6ecf57e80e28e6b3fc8b6066ac9d7b8836d22fb65cde7d6ede66b910` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/enc-tree-sec12.txt` | 421 | `627992802ac6b6441b57871a04e12d897c9bf431605300825086fed4b9287528` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/enc-tree.txt` | 421 | `627992802ac6b6441b57871a04e12d897c9bf431605300825086fed4b9287528` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/g6-4-check.txt` | 191 | `4328bc577893787464fb1568dc3fcd63f5c5abc04f5426adf35965dd3bf5d9c1` | 被受版本控制文档引用其路径 |
| `archive/gate-index-sec12.txt` | 2145 | `043e108057e87dc83b92efbcbc3412644cc82674402784dd845a647f8d882019` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-index.txt` | 2620 | `3d9ff677a931ca91e83cfdf2d72500952fd7d42f3307a48d401c2d15268d5c1c` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-selftest.txt` | 17597 | `2952826f2fece4291b5c6ad980b5a7b337cbee8d0133e8961733b8650f528217` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-sec12.txt` | 2144 | `ce4cebe1201d7164256ef9c1478e393ba12e607f1fca043f4d7f1c98a8e6c202` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree.txt` | 2619 | `80e30dffaee69e0cf67699e339ca6df9ff6a46d35459de6f255f95bbf80fdd77` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/m4-g6-4-probe.js` | 3992 | `7ea27619463eae3871e61875a25bfc0df0c65c6572cdd5b8aeddec0a772acf0d` | 被受版本控制文档引用其路径 |
| `archive/mirror-sec12.txt` | 250 | `da266310f71a4f1b1caef2966dc67ca4afbbf568c40d9e20ccf43f9261c64f14` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/mirror.txt` | 250 | `c731bcbcba5ca952b8af28fbb21544e07501d5656b32e850b768724d44e728d5` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/PLAN-INPUT-①.txt` | 9595 | `d3aec9599198f9684ff9313af9cfb26366b42f09f448a3fb028b228084aa3505` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/PLAN.txt` | 38792 | `5dcc965aeee2187a64189a3edba6848509ef35dd47ff45b34be8c88a3596aeed` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/refs-check-sec12.txt` | 773 | `554fff9f12efaf51541b7289bb57c22f48a72d396bd09a4b2b0efc39dfa8bc63` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/refs-check.js` | 3715 | `c59b979d298a4096b526276c897e0c233b16338f49c952d81649e46610444e2a` | 被受版本控制文档引用其路径 |
| `archive/refs-check.txt` | 782 | `b3fb015bd39c752282469a7f0fc00c124e72ad6ecb44d500dbf65b81289f6b3f` | 被受版本控制文档引用其路径 |
| `archive/review-input.diff` | 353349 | `8a4a1a000357c8628c4bb8688e92e5633ced0a6e11c91dbd2ddb0398b81673fd` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input.stat.txt` | 218 | `62c26740c7b6cbb85b27ac15ea165d72a107e8cfa3ed61fcda672d2efef55f25` | 切片级文档 / 评审包 / 门禁读数 |

## 排除清单（未归档）

| 相对路径 | 字节 | sha256 | 未归档理由 |
|---|---|---|---|
| `/a5-change-log-column.js` | 6780 | `c57b0a88a380a6b8450e2651ead76775e1b6df238bf4c0e656abcf532baf0abf` | 落地与修复脚本残留 |
| `/align.js` | 1163 | `c4b0a3f60aad2dfca5ff0755d820a41ccdfa6108abd2b18ee83ad7abb7e1103a` | 落地与修复脚本残留 |
| `/align.txt` | 2174 | `0cdb5a6d38281d092c2aedbd924c7d45c67d95a03f085baad9869617013da678` | 中间转储（未被引用的过程输出） |
| `/append-sec12.js` | 2419 | `042eaf804b3a3b7e2d297e3ed75d01bd5f06c0e12ae4f8952541dd457ad4c236` | 落地与修复脚本残留 |
| `/bolds.txt` | 660 | `4871a7b53a7769f83d26caf109f5d5f4b62bc9b6130428d6c308d5c0aba33f95` | 中间转储（未被引用的过程输出） |
| `/charcheck.js` | 1928 | `e0be982aa6807db7a714fb98a330b33480fabcf6d18bddbde792e5a03803912f` | 落地与修复脚本残留 |
| `/check.txt` | 489 | `5620f5a7221fc402e152f72b2a15588bcd3e0ac717e1149c190d5cd133ba8cd2` | 中间转储（未被引用的过程输出） |
| `/enc-tree.js` | 1777 | `c0bda476280446164ddd3efe2db57463f6bd393e67bf9e5e3ca4e29d10ad78b2` | 落地与修复脚本残留 |
| `/evidence-selfcheck.js` | 2355 | `b8b3f73052559c7845956cbec14ec94db661f2ecd8e3e1259272eae1cf4326a1` | 落地与修复脚本残留 |
| `/evidence-selfcheck.txt` | 416 | `c0f49ee468195529b8640e5796a62d1d81d2b0eee23a47848f8f715410db9f57` | 中间转储（未被引用的过程输出） |
| `/fails.txt` | 250 | `1e579caf08bba7429644be8b4af4bcadda8bf50a00eaccac8c5a0d9a53761239` | 中间转储（未被引用的过程输出） |
| `/fix-section8.js` | 4390 | `98263332bb918a569653dcceba5d470413341edbf1f032d0a070401798aa602d` | 落地与修复脚本残留 |
| `/fix-splice.js` | 3066 | `3c938c5acc307af6a631f8fbf7b7827458b2ed105e5d051e26e95dde7eef749d` | 落地与修复脚本残留 |
| `/g6-4-check-sec12.txt` | 192 | `965964ad8ecc1d887114c599eecaeecbf6de7438f6c01be90e81a7eb0f6c5052` | 中间转储（未被引用的过程输出） |
| `/m4-probe.txt` | 1533 | `60be3132e92588668f3dffd4aaf6444ac606324a816d5beb7151de55c0e91038` | 中间转储（未被引用的过程输出） |
| `/mirror-check.js` | 1361 | `875ceea67a8255f9ee3e02e2ad9101efec3aeef6a599aa8c46e9732b12451fcb` | 落地与修复脚本残留 |
| `/normalize-md.js` | 671 | `80edd42a30926c37c73ef2d23d87487e64052c3add82198454a9151ee72ec552` | 落地与修复脚本残留 |
| `/peek-en.js` | 1364 | `405f10565a8c9ece86d62dc7189ef3d46f960fb7785d861486927e420371f35d` | 落地与修复脚本残留 |

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
| `PLAN-INPUT-①.md` | `archive/PLAN-INPUT-①.txt` | `f650d76f3a4b429575e2a822319f7cf3bf57443aa8ebf115548a113efd23f805` | `d3aec9599198f9684ff9313af9cfb26366b42f09f448a3fb028b228084aa3505` | 后缀 `.md` → `.txt` / LF 归一 |
| `PLAN.md` | `archive/PLAN.txt` | `5dcc965aeee2187a64189a3edba6848509ef35dd47ff45b34be8c88a3596aeed` | `5dcc965aeee2187a64189a3edba6848509ef35dd47ff45b34be8c88a3596aeed` | 后缀 `.md` → `.txt` |
| `append-sec12-cn.md` | `archive/append-sec12-cn.txt` | `7b0780e287082b8972839540ad71c9211a61b0ef25d977dfebaaf263efde28bc` | `7b0780e287082b8972839540ad71c9211a61b0ef25d977dfebaaf263efde28bc` | 后缀 `.md` → `.txt` |
| `append-sec12-en.md` | `archive/append-sec12-en.txt` | `2e63bddda852550ee02ae87450ea4f8771bbc3728ad9df2cc7bfca89269dda34` | `2e63bddda852550ee02ae87450ea4f8771bbc3728ad9df2cc7bfca89269dda34` | 后缀 `.md` → `.txt` |
| `c1-recheck.js` | `archive/c1-recheck.js` | `a6d01ca146a5bb35781e9abec727699905be953e4a6f5eb0b189c1bceb25a8af` | `631f8c6c8cdf3995b0d17f481aba7dab9a4a90c6dade7e9536cb65e944efec50` | LF 归一 |
| `c1-recheck.txt` | `archive/c1-recheck.txt` | `f2485b751bb7bd420fe420aff68b1d6cd734df4cf931ad32c556c9dbb1659198` | `a79a8faf6ecf57e80e28e6b3fc8b6066ac9d7b8836d22fb65cde7d6ede66b910` | LF 归一 |
| `enc-tree-sec12.txt` | `archive/enc-tree-sec12.txt` | `b8908c7b2351e1f11e39925ec7aa6b611212491520dacdd20ea8ed43a4fd57c3` | `627992802ac6b6441b57871a04e12d897c9bf431605300825086fed4b9287528` | LF 归一 |
| `enc-tree.txt` | `archive/enc-tree.txt` | `b8908c7b2351e1f11e39925ec7aa6b611212491520dacdd20ea8ed43a4fd57c3` | `627992802ac6b6441b57871a04e12d897c9bf431605300825086fed4b9287528` | LF 归一 |
| `g6-4-check.txt` | `archive/g6-4-check.txt` | `9daa63b878aa6c9463898b735aae579c045dd899a9ce9922eb27718db751836d` | `4328bc577893787464fb1568dc3fcd63f5c5abc04f5426adf35965dd3bf5d9c1` | LF 归一 |
| `gate-index-sec12.txt` | `archive/gate-index-sec12.txt` | `c5127c8b6ee15c14ddf54f407924da4d37ffe9ed9c00fe17b1298f8e882da7cc` | `043e108057e87dc83b92efbcbc3412644cc82674402784dd845a647f8d882019` | LF 归一 |
| `gate-index.txt` | `archive/gate-index.txt` | `b240d05b9ba70ec9fb1389288f5116e01f2998196a6c64a610a9ea94eb9da1a0` | `3d9ff677a931ca91e83cfdf2d72500952fd7d42f3307a48d401c2d15268d5c1c` | LF 归一 |
| `gate-selftest.txt` | `archive/gate-selftest.txt` | `02498ccfed6a9d7dbd7395caa882c2e4fb51d5d9fbacfd17a2bbe879261517c9` | `2952826f2fece4291b5c6ad980b5a7b337cbee8d0133e8961733b8650f528217` | LF 归一 |
| `gate-tree-sec12.txt` | `archive/gate-tree-sec12.txt` | `edf003ee64460873b256fe1ebed222f236187aebd596980085d6b19898b69246` | `ce4cebe1201d7164256ef9c1478e393ba12e607f1fca043f4d7f1c98a8e6c202` | LF 归一 |
| `gate-tree.txt` | `archive/gate-tree.txt` | `ebfdf0dcd17103260ccc38017dff086ee0d4adde0227d1db178a54051b229369` | `80e30dffaee69e0cf67699e339ca6df9ff6a46d35459de6f255f95bbf80fdd77` | LF 归一 |
| `m4-g6-4-probe.js` | `archive/m4-g6-4-probe.js` | `8292a26bbb96450ebc22f47151e2793757713a7febf1bd3640041392bb44efd0` | `7ea27619463eae3871e61875a25bfc0df0c65c6572cdd5b8aeddec0a772acf0d` | LF 归一 |
| `mirror-sec12.txt` | `archive/mirror-sec12.txt` | `9145d2db36f760f397a62b2acad632f53b3466ac81bcae8f9f14a8c64cf3b68b` | `da266310f71a4f1b1caef2966dc67ca4afbbf568c40d9e20ccf43f9261c64f14` | LF 归一 |
| `mirror.txt` | `archive/mirror.txt` | `b1f81c12d2be7e1efad187fb10431e30793be9982e135f2d39660cb25ad93609` | `c731bcbcba5ca952b8af28fbb21544e07501d5656b32e850b768724d44e728d5` | LF 归一 |
| `refs-check-sec12.txt` | `archive/refs-check-sec12.txt` | `e7e1c64e45f063b86ea0754e4c61361a8a06464fcc6c7a5a802e97fcddc2e78a` | `554fff9f12efaf51541b7289bb57c22f48a72d396bd09a4b2b0efc39dfa8bc63` | LF 归一 |
| `refs-check.js` | `archive/refs-check.js` | `8dc8b5bfcb09ea32b9e450443ae0078aba774ff5dea8ac3de49f2c7165dfc632` | `c59b979d298a4096b526276c897e0c233b16338f49c952d81649e46610444e2a` | LF 归一 |
| `refs-check.txt` | `archive/refs-check.txt` | `2bbb5daa3c1200f9ef4a96e21d457a583fa3bee6a75e9bcbde91055bf2e68146` | `b3fb015bd39c752282469a7f0fc00c124e72ad6ecb44d500dbf65b81289f6b3f` | LF 归一 |

## 对应引用点

| 引用文件 | 行号 |
|---|---|
| `"docs/validation/evidence/directive-history/v112/archive/PLAN-INPUT-\342\221\240.txt"` | 72 |
| `docs/DELIVERY_DIRECTIVE.md` | 120 / 860 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 120 / 860 |
| `docs/validation/directive-v1.12.md` | 231 / 246 / 248 / 284 / 296 / 341 |
| `docs/validation/directive-v1.12_EN.md` | 231 / 246 / 248 / 284 / 296 / 341 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/classification.txt` | 1437 / 1438 / 1439 / 1440 / 1441 / 1442 / 1443 / 1444 / 1445 / 1446 / 1447 / 1448 / 1449 / 1450 / 1451 / 1452 / 1453 / 1454 / 1455 / 1456 / 1457 / 1458 / 1459 / 1460 / 1461 / 1462 / 1463 / 1464 / 1465 / 1466 / 1467 / 1468 / 1469 / 1470 / 1471 / 1472 / 1473 / 1474 / 1475 / 1476 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/facts2.txt` | 11 / 20 / 42 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/inventory.txt` | 195 / 196 / 197 / 198 / 199 / 200 / 201 / 441 / 442 / 443 / 444 / 445 / 460 / 485 / 985 / 986 / 987 / 988 / 989 / 990 / 1002 / 1003 / 1004 / 1005 / 1006 / 1007 / 1285 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/measure.txt` | 28 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/sets.txt` | 38 |
| `docs/validation/evidence/directive-history/v112/archive/PLAN.txt` | 8 / 191 |
| `docs/validation/evidence/directive-history/v112/archive/m4-g6-4-probe.js` | 12 |
| `docs/validation/evidence/directive-history/v112/archive/refs-check.js` | 16 |
| `docs/validation/evidence/directive-history/v112/archive/review-input.diff` | 1471 / 1486 / 1488 / 1524 / 1536 / 1581 / 1823 / 1838 / 1840 / 1876 / 1888 / 1933 |
| `tools/gates/lib/criteria/g6.js` | 450 |
