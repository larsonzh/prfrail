# 归档 MANIFEST · drfix

- **源路径**：`tmp/drfix/`
- **归档片**：`DIRECTIVE-TMP-LIFECYCLE`（v1.23，2026-10-04）
- **归档口径**：§1.5.1 的 ② 类「早期证据评估归档」——保留证据条目（被受版本控制文档引用其路径者 ∪ 切片级文档），排除运行时残留与一次性脚本
- **归档条目**：19 条 / 335,863 字节（归档载荷根 = `archive/`；逐条 sha256 见 `SHA256SUMS.txt`）
- **排除条目**：19 条（逐条见下「排除清单」）
- **未入库条目**：0 条（逐条见下「未入库文件清单」）
- **字节保真**：归档副本的**字符内容不改**；因 §6.3 的作用域（G4a 扫全部变更文件）与仓库编码约定，副本按**归档后类型**归一编码——`.ps1` 带 BOM，其余不带 BOM，一律 LF；逐条变更种类与**原始字节 sha256** 见下「后缀与编码映射」节
- **后缀口径**：归档副本的 `.md` 一律改存为 `.txt`（本单元 3 条），以避开只扫 `.md` 的文档判据；映射逐条见下节
- **编码归一计数**：后缀改写 3 条 / BOM 剔除 0 条 / BOM 补入 0 条 / LF 归一 5 条

## 归档条目清单

| 相对路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/A2-restored-green.txt` | 753 | `34daf5d51b7e3f1deb03e3306cf4a7f6c0cbd9c93ff9a36e9237411d13512d25` | 被受版本控制文档引用其路径 |
| `archive/ARCH-DESIGN.txt` | 19695 | `a5984b9f698f102359620deae45fe163d18dbde01450bf75430a897a529c6d16` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/B5-postcommit-count10.txt` | 76 | `170a667267a8e26647187e270c9bafa2525442b3127d97600304c32a1df7648f` | 被受版本控制文档引用其路径 |
| `archive/B6-postcommit-psselftest.txt` | 739 | `e3726be75d6a6583c6802ba92d55526fa634d8200e48be2282aeb44d3aac02cb` | 被受版本控制文档引用其路径 |
| `archive/CI-push1.txt` | 44734 | `017d561c1f4a2c007e21ffac350abf8266ae69953b19c7f73e00aa655d571abb` | 被受版本控制文档引用其路径 |
| `archive/DOCS.diff` | 21379 | `13971a781c7ad0564167c8281d387a758be0b4412b8f45c106d5eee6d1afd569` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/EVIDENCE-SUMMARY.txt` | 8966 | `d9c8d3a59d31b5f9e45cd7a3c1fc3f3a695ec1b023af0b27e40bb93734d2438d` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-remediation.txt` | 1365 | `224f6b82aac999d8c5c2b0c30830734c51ed5a178aab8f4b9dbd2c1f02609580` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/IMPL.diff` | 11217 | `af4ba4d70a9839b08c5959e6a06e69f0d00d10d8284e0281ab98f968afbc2f96` | 被受版本控制文档引用其路径 |
| `archive/REPORT-REQUIREMENTS.txt` | 3892 | `6f2e7bd8ed9bb0b1d272d1557007ea3fdc3ba749f28d3678fc73894bcac95685` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review/DRFIX-IMPL.diff` | 5008 | `7c186c6769e56e5ab5481df7b81cf94f7403642321d8adef6a7b38b8707d6b28` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review/OB9-T3T4-criteria-52c28af.diff` | 15198 | `5fe81f1d6547b5df300dd1daa34894272a6b6b2c85577387d7bc3941b86383e0` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review/OB9-T4-note-move-e0b08ec.diff` | 18033 | `32ea57deb82a21c4d49f872b88ba320b49e1d673ae6f7249fa8a528df30224a0` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review/OB9-v1.7-c473fd2.diff` | 44903 | `7441cf3877f3abeebd2e945e13748251a17b5e2771620068d22170487cad233e` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review/OB9-v1.8-30cf799.diff` | 31826 | `454230826d3fd25e04cbcc21fe8020ffbed148a8f1c12b653909603aa46d793f` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review/OB9-v1.9-6f186fd..a366fe3.diff` | 93344 | `6d345f7d59c80dd0ac6b74c63599cb0c5b6d3d91f1c270f020a3d5f0aa5911c2` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review/README-REREVIEW.txt` | 7700 | `4af66acee870e095f669013916ca4f5d337969085d45d8df96de2e93ffe5c212` | 被受版本控制文档引用其路径 |
| `archive/review/README.txt` | 1833 | `fbcc4e024f01a7d6aa2eda400d8cafe466748a9eb25a46702f2a4e6c787e6182` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/SLICE-DEFINITION.txt` | 5202 | `47d2e9e62cede7ac59f7c55a0e9290cd542e39043155927f0c606f44a72a3c13` | 切片级文档 / 评审包 / 门禁读数 |

## 排除清单（未归档）

| 相对路径 | 字节 | sha256 | 未归档理由 |
|---|---|---|---|
| `/A1-T1T4-injected.txt` | 1110 | `5752dc97192737d1ae8f69302417359d44f6526ef071087e0d910783be910134` | 中间转储（未被引用的过程输出） |
| `/A3-guard-removed.txt` | 1088 | `480d1b1247a7b96296b3296cf20c5c61d1990084c565c695ea01e7aa90df05d0` | 中间转储（未被引用的过程输出） |
| `/A4-restored-green-full.txt` | 2009 | `56e9ae54db54cb26035b4a7c9aba347b895ea9906a9cd727247c21cf1c9c18f2` | 中间转储（未被引用的过程输出） |
| `/A5-T6-injected.txt` | 414 | `c5b7317983a4e11d580c138aace6728d358e4ef451e7434f754b0bbc619e12b8` | 中间转储（未被引用的过程输出） |
| `/A6-T7-Fix3-reverted.txt` | 875 | `681447d30d89d46f1bf82db3a583fa2783a9ab788a98e0a6cc4b49658f84df58` | 中间转储（未被引用的过程输出） |
| `/A7-T1-mutation-red.txt` | 302 | `8bc75a126eddee53f683d791dc6f18a599336f3df0d48c6c50664dd3f0be748b` | 中间转储（未被引用的过程输出） |
| `/A8-T1-tolerant-too-much-red.txt` | 275 | `788f246679d05aa9e9ab347787bf499b464e5dbbcefb1367b0401b2170e8ec96` | 中间转储（未被引用的过程输出） |
| `/B-count10-verbose.txt` | 19370 | `4245fa82cdb65cbf2797df409339550b1d7acf0f8a4e1098c237481b43046123` | 中间转储（未被引用的过程输出） |
| `/B2-count10-T1.txt` | 77 | `99c640869c84e2cc736bc9565237b83fcee9064dc08060f6311fad59f1be411b` | 中间转储（未被引用的过程输出） |
| `/B3-full-suite-T1.txt` | 1112 | `2ea1d157b446e3c1b528ddafeebeb32119e387d7eeddbb6ec79133eca3e76d7e` | 中间转储（未被引用的过程输出） |
| `/B4-psself-test-T1.txt` | 753 | `9d29954bb6845f3806d6a2b85464754c105f5a441a4b7c890ce7344b6090a516` | 中间转储（未被引用的过程输出） |
| `/CI-push2.txt` | 44798 | `1e33ea3bee72973b5969ce5300d3aef63baae70044410d02ec79262f8340e161` | 中间转储（未被引用的过程输出） |
| `/CI-push3.txt` | 44733 | `f6ac0cb2f0fe6c29f8387fa5e0ac79f61087a064c3be78239ed71516d10a6bd2` | 中间转储（未被引用的过程输出） |
| `/CI-push4.txt` | 44799 | `2636b48349d480a6dd803e1baaa4fecafcec5ca3fe7f456d3e3a8fd947f1c1c4` | 中间转储（未被引用的过程输出） |
| `/CI-push5.txt` | 44733 | `e67a059ee0f77659e5628bbcc4b0e0755d5aac9764a7f32ac6ad60bb68f0d4ef` | 中间转储（未被引用的过程输出） |
| `/T1-isolated-edit-diff.txt` | 5612 | `702cf6d04e0b96bc06de3743b66fd0e362e5a96beb2b9588ab91641ab1bcc53b` | 中间转储（未被引用的过程输出） |
| `/ci-log.js` | 1175 | `293b63cf5b64d00a07a70a68115090f1e2dba853ae88f2405f0f56e2830f34fa` | 落地与修复脚本残留 |
| `/pkg4.txt` | 1056 | `64c4aef0b7bb37a18b9d15e7baf13c168057270fef32d5edc78ea5ece9defb61` | 中间转储（未被引用的过程输出） |
| `/pkg5.txt` | 321 | `e806814cc65681e9b0afc1bb9a3c95f64faa8bf3bfcea9c6a83c3d6271b51bd1` | 中间转储（未被引用的过程输出） |

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
| `A2-restored-green.txt` | `archive/A2-restored-green.txt` | `04553c43867770c3a061f4a19a562e3a1fcc716e352b80b30d4dacdc31691981` | `34daf5d51b7e3f1deb03e3306cf4a7f6c0cbd9c93ff9a36e9237411d13512d25` | LF 归一 |
| `ARCH-DESIGN.md` | `archive/ARCH-DESIGN.txt` | `a5984b9f698f102359620deae45fe163d18dbde01450bf75430a897a529c6d16` | `a5984b9f698f102359620deae45fe163d18dbde01450bf75430a897a529c6d16` | 后缀 `.md` → `.txt` |
| `B5-postcommit-count10.txt` | `archive/B5-postcommit-count10.txt` | `281e8a21c0bedcd03209daa50dbe02ab5f6e6762189af777941ef0d206980721` | `170a667267a8e26647187e270c9bafa2525442b3127d97600304c32a1df7648f` | LF 归一 |
| `B6-postcommit-psselftest.txt` | `archive/B6-postcommit-psselftest.txt` | `9d29954bb6845f3806d6a2b85464754c105f5a441a4b7c890ce7344b6090a516` | `e3726be75d6a6583c6802ba92d55526fa634d8200e48be2282aeb44d3aac02cb` | LF 归一 |
| `REPORT-REQUIREMENTS.md` | `archive/REPORT-REQUIREMENTS.txt` | `8a6c3a8371d63fb2534caee74b3eeb1f63c23df24dfbc031e2d4d561eb516261` | `6f2e7bd8ed9bb0b1d272d1557007ea3fdc3ba749f28d3678fc73894bcac95685` | 后缀 `.md` → `.txt` / LF 归一 |
| `SLICE-DEFINITION.md` | `archive/SLICE-DEFINITION.txt` | `7a8aa2008d23f3b14cd82779ef5b942d602173cbf10f9e73b1bcc49ae7863832` | `47d2e9e62cede7ac59f7c55a0e9290cd542e39043155927f0c606f44a72a3c13` | 后缀 `.md` → `.txt` / LF 归一 |

## 对应引用点

| 引用文件 | 行号 |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 120 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 120 |
| `docs/validation/dr-fix-enforcement-proxy.md` | 80 / 84 / 94 / 144 / 159 / 196 / 201 |
| `docs/validation/dr-fix-enforcement-proxy_EN.md` | 80 / 84 / 94 / 144 / 159 / 196 / 201 |
| `docs/validation/evidence/directive-history/drfix/archive/ARCH-DESIGN.txt` | 3 |
| `docs/validation/evidence/directive-history/drfix/archive/SLICE-DEFINITION.txt` | 46 |
| `docs/validation/evidence/directive-history/drfix/archive/review/README.txt` | 7 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/classification.txt` | 528 / 529 / 530 / 531 / 532 / 533 / 534 / 535 / 536 / 537 / 538 / 539 / 540 / 541 / 542 / 543 / 544 / 545 / 546 / 547 / 548 / 549 / 550 / 551 / 552 / 553 / 554 / 555 / 556 / 557 / 558 / 559 / 560 / 561 / 562 / 563 / 564 / 565 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/facts2.txt` | 11 / 42 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/inventory.txt` | 128 / 129 / 130 / 131 / 132 / 133 / 134 / 135 / 136 / 137 / 138 / 441 / 442 / 443 / 444 / 445 / 466 / 467 / 468 / 469 / 470 / 1029 / 1030 / 1031 / 1032 / 1034 / 1035 / 1036 / 1037 / 1038 / 1039 / 1040 / 1041 / 1043 / 1044 / 1045 / 1046 / 1288 / 1289 / 1290 / 1291 / 1293 / 1294 / 1295 |
| `tools/gates/selftest.js` | 831 / 839 |
| `tools/gates/testdata/fixtures/g7-b-no-disposition.txt` | 6 |
| `tools/gates/testdata/fixtures/g7-ok-report.txt` | 13 |
| `tools/gates/testdata/historical/dr-fix-report-structure.txt` | 25 |
