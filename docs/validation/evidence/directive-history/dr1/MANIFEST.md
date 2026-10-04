# 归档 MANIFEST · dr1

- **源路径**：`tmp/dr1/`
- **归档片**：`DIRECTIVE-TMP-LIFECYCLE`（v1.23，2026-10-04）
- **归档口径**：§1.5.1 的 ② 类「早期证据评估归档」——保留证据条目（被受版本控制文档引用其路径者 ∪ 切片级文档），排除运行时残留与一次性脚本
- **归档条目**：14 条 / 137,253 字节（归档载荷根 = `archive/`；逐条 sha256 见 `SHA256SUMS.txt`）
- **排除条目**：17 条（逐条见下「排除清单」）
- **未入库条目**：0 条（逐条见下「未入库文件清单」）
- **字节保真**：归档副本的**字符内容不改**；因 §6.3 的作用域（G4a 扫全部变更文件）与仓库编码约定，副本按**归档后类型**归一编码——`.ps1` 带 BOM，其余不带 BOM，一律 LF；逐条变更种类与**原始字节 sha256** 见下「后缀与编码映射」节
- **后缀口径**：归档副本的 `.md` 一律改存为 `.txt`（本单元 5 条），以避开只扫 `.md` 的文档判据；映射逐条见下节
- **编码归一计数**：后缀改写 5 条 / BOM 剔除 0 条 / BOM 补入 0 条 / LF 归一 9 条

## 归档条目清单

| 相对路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/draft-cn.txt` | 2478 | `56f306af99be2bcc45d64f586e6cf47422432a941ffb9e2d77faefeff08971f9` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/draft-en.txt` | 2780 | `6721902a619ab2bfc2245e90793ea2838a2cf9747a1f2d384590ca8ab0395380` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/draft-v114-cn.txt` | 3277 | `392e2227570b4790c57352bf3d5a85d4752b1cbe1d95d86fba3bfb76012a0ea0` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/draft-v114-en.txt` | 4018 | `166a11add4a43e17aff9105c6815ea1662a461a9a82e9837049fc80c756948e7` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-final.txt` | 3618 | `08037718767dde1830b4f12a8d0328371f8457c6a581c8c915972b09f7f8005a` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-final2.txt` | 3618 | `08037718767dde1830b4f12a8d0328371f8457c6a581c8c915972b09f7f8005a` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r2.txt` | 3048 | `499656da95dd95aa8d691f6b393d77da08c9bdc90adcbbadfcf4f043ca9b22bb` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r3.txt` | 7770 | `05a8aa4cfae78cfda425665d12c8243432359c1f5545e00ca9fccab0eea9f18f` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r4.txt` | 3648 | `c0bba3391c6b483a0cba6f918fa2feb88c2ed2fdfa82b3c2e6878a83d53ab32f` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree.txt` | 2831 | `a7c9e0c70147cb5878b5c6cf96dbeb9697ee7fac28146158096a15d1d7922999` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/mutation-check.txt` | 2350 | `1fc8141f2b693278b7af92d20051caa3dfcbac4775fc0aec90692c4fa69faed2` | 被受版本控制文档引用其路径 |
| `archive/review-input-r2.diff` | 53600 | `0550d56026fe1a49a22b3a6ac08cf06fb9b7bc354093ff5220ff88837ced0b38` | 被受版本控制文档引用其路径 |
| `archive/review-input.diff` | 42393 | `dadf7f5e20237ba191336d871291a8cf505d4b28227daeaf97bba96e92a7a7f0` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/rulings.txt` | 1824 | `e0da9a92493029f2dcd56d3121a63a0e6802dd2d6be6296490574aeeec1e07a5` | 被受版本控制文档引用其路径 |

## 排除清单（未归档）

| 相对路径 | 字节 | sha256 | 未归档理由 |
|---|---|---|---|
| `/adr014-rows.md` | 2132 | `3433e56f861b6396939b8a59ee951ba5ad15969d7983dc25f67cb12b0f227a9d` | 中间转储（未被引用的过程输出） |
| `/apply-adr014.js` | 2752 | `f46053cc42d2cd0d9f6f6a72ab72efcc59dbf68f568ba2242d7f512316d990e1` | 落地与修复脚本残留 |
| `/assemble-pack.js` | 4150 | `0fb9a8706d86deb057fe1a853e5468c1a3e4b63467442700738af6eaf8ab684e` | 落地与修复脚本残留 |
| `/ci-win.log` | 74698 | `e5f31c6ac9111517fe9ae7d31584ef0792e97c90ad16653a778ece7d1f539e84` | 运行时状态与日志 |
| `/config-explain.txt` | 3388 | `95114c1443e2bded13b87c867407a5a9ee4ee9317f17c2b569ed0186d3bdfdc1` | 中间转储（未被引用的过程输出） |
| `/contract-excerpt.md` | 905 | `0281f2a847121fc61d17f90f4835fa1d9573473fa3036d589b22f4e0b7e8370a` | 中间转储（未被引用的过程输出） |
| `/focused-test-r2.txt` | 5129 | `2b3fcca1e9992e2b8e3ad6f4f16c00bc3a90fcc9ddcbe581f17bf3af31e32dfa` | 中间转储（未被引用的过程输出） |
| `/focused-test.txt` | 5129 | `f3acc5fba5d6babf33db43430703af47f915dc61a11c5f3bd86a81a2ed958252` | 中间转储（未被引用的过程输出） |
| `/head/ADR_REGISTER.md` | 17923 | `49708db4a0355de900b4caf9bbd5101bde41d64a2900c43374455251ca92b4f2` | 中间转储（未被引用的过程输出） |
| `/head/ADR_REGISTER_EN.md` | 20353 | `92382ab4452e4f7034ddad83c35bdd53f9cf7ca15d8e8437f1c02e1355e8c14c` | 中间转储（未被引用的过程输出） |
| `/head/REMAINING_SLICES.md` | 108623 | `52756bc3c3529dadb7fad22acfecf1c73d4a1ce0063b2dc90ef68b42d9484b11` | 中间转储（未被引用的过程输出） |
| `/head/REMAINING_SLICES_EN.md` | 118787 | `c39b98addf2010cdda41627014f2015bd968ca9217c1c9a7854e646364ba90e8` | 中间转储（未被引用的过程输出） |
| `/probe-01-verify.txt` | 982 | `1974906696632107b23225ccfc1f18af2c7854b1d99ddb1e909c464dd9c1de08` | 中间转储（未被引用的过程输出） |
| `/probe-01.txt` | 981 | `b2c65c953ca750f4cb35af9bcab759385c22b0ab4b5863a0310ac191959eedfe` | 中间转储（未被引用的过程输出） |
| `/stage7-routing.md` | 2015 | `403ad792ff681876e09892b713e0d9804bba3e368012aa520730510b1ddbe957` | 中间转储（未被引用的过程输出） |
| `/writeback-8b-fix.js` | 5351 | `e1d6ffe3653f80e50f003adfa634658f650699deb6bd44b52c3eb1287113ac4d` | 落地与修复脚本残留 |
| `/writeback-8b.js` | 9089 | `ce24fa681b9f334230d0dda5d2e25f1f7439cca8634198b5ab030b623abf3f14` | 落地与修复脚本残留 |

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
| `draft-cn.md` | `archive/draft-cn.txt` | `4e43aa4c20bf30b5457baaebcbff02f36f6c19848013fc3b28c5198cf26c0d25` | `56f306af99be2bcc45d64f586e6cf47422432a941ffb9e2d77faefeff08971f9` | 后缀 `.md` → `.txt` / LF 归一 |
| `draft-en.md` | `archive/draft-en.txt` | `1c89c7707d8504af37528b4319fdbcebe503cc0d18c265f6956bd58af2ffe58b` | `6721902a619ab2bfc2245e90793ea2838a2cf9747a1f2d384590ca8ab0395380` | 后缀 `.md` → `.txt` / LF 归一 |
| `draft-v114-cn.md` | `archive/draft-v114-cn.txt` | `4973d390c4fb479b6462955af89d3aaf4737b6666eefae6196d7fded297d47cf` | `392e2227570b4790c57352bf3d5a85d4752b1cbe1d95d86fba3bfb76012a0ea0` | 后缀 `.md` → `.txt` / LF 归一 |
| `draft-v114-en.md` | `archive/draft-v114-en.txt` | `1231d52f7b0162496f84879c7d4a9cef880602145c49e35b7b72343b5f189751` | `166a11add4a43e17aff9105c6815ea1662a461a9a82e9837049fc80c756948e7` | 后缀 `.md` → `.txt` / LF 归一 |
| `gate-tree-final.txt` | `archive/gate-tree-final.txt` | `622d705cc161b17cd6dae0454518752ffc717990170ca339a17ff7f884bdfb43` | `08037718767dde1830b4f12a8d0328371f8457c6a581c8c915972b09f7f8005a` | LF 归一 |
| `gate-tree-final2.txt` | `archive/gate-tree-final2.txt` | `622d705cc161b17cd6dae0454518752ffc717990170ca339a17ff7f884bdfb43` | `08037718767dde1830b4f12a8d0328371f8457c6a581c8c915972b09f7f8005a` | LF 归一 |
| `gate-tree-r3.txt` | `archive/gate-tree-r3.txt` | `c5e9937b73e55c9870f426abb64054f92862cfe493b24f1bd35a4c67eb033bde` | `05a8aa4cfae78cfda425665d12c8243432359c1f5545e00ca9fccab0eea9f18f` | LF 归一 |
| `gate-tree-r4.txt` | `archive/gate-tree-r4.txt` | `02f3de447cb78e104034a878d9612091ee3dfbb9c86e2f485e3777f5ef0a5dfd` | `c0bba3391c6b483a0cba6f918fa2feb88c2ed2fdfa82b3c2e6878a83d53ab32f` | LF 归一 |
| `mutation-check.txt` | `archive/mutation-check.txt` | `4e08c97271bd6d7d98eea8a6d411616626ce3ad7826d7522d540ddf7fe98cdb3` | `1fc8141f2b693278b7af92d20051caa3dfcbac4775fc0aec90692c4fa69faed2` | LF 归一 |
| `rulings.md` | `archive/rulings.txt` | `e0da9a92493029f2dcd56d3121a63a0e6802dd2d6be6296490574aeeec1e07a5` | `e0da9a92493029f2dcd56d3121a63a0e6802dd2d6be6296490574aeeec1e07a5` | 后缀 `.md` → `.txt` |

## 对应引用点

| 引用文件 | 行号 |
|---|---|
| `docs/validation/dr-1-availability-probe-floor.md` | 31 / 33 / 61 / 123 / 127 |
| `docs/validation/dr-1-availability-probe-floor_EN.md` | 31 / 33 / 61 / 123 / 127 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/classification.txt` | 496 / 497 / 498 / 499 / 500 / 501 / 502 / 503 / 504 / 505 / 506 / 507 / 508 / 509 / 510 / 511 / 512 / 513 / 514 / 515 / 516 / 517 / 518 / 519 / 520 / 521 / 522 / 523 / 524 / 525 / 526 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/inventory.txt` | 118 / 119 / 120 / 121 / 122 / 1019 / 1020 / 1021 / 1022 / 1023 / 1024 / 1025 / 1026 / 1027 / 1028 / 1058 / 1059 |
| `docs/validation/evidence/dr-1-2026-09-24/README.md` | 5 / 32 |
