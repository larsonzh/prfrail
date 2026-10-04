# 归档 MANIFEST · assembly-min

- **源路径**：`tmp/assembly-min/`
- **归档片**：`DIRECTIVE-TMP-LIFECYCLE`（v1.23，2026-10-04）
- **归档口径**：§1.5.1 的 ② 类「早期证据评估归档」——保留证据条目（被受版本控制文档引用其路径者 ∪ 切片级文档），排除运行时残留与一次性脚本
- **归档条目**：12 条 / 659,581 字节（归档载荷根 = `archive/`；逐条 sha256 见 `SHA256SUMS.txt`）
- **排除条目**：72 条（逐条见下「排除清单」）
- **未入库条目**：0 条（逐条见下「未入库文件清单」）
- **字节保真**：归档副本的**字符内容不改**；因 §6.3 的作用域（G4a 扫全部变更文件）与仓库编码约定，副本按**归档后类型**归一编码——`.ps1` 带 BOM，其余不带 BOM，一律 LF；逐条变更种类与**原始字节 sha256** 见下「后缀与编码映射」节
- **后缀口径**：归档副本的 `.md` 一律改存为 `.txt`（本单元 3 条），以避开只扫 `.md` 的文档判据；映射逐条见下节
- **编码归一计数**：后缀改写 3 条 / BOM 剔除 1 条 / BOM 补入 0 条 / LF 归一 0 条

## 归档条目清单

| 相对路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/change-set.diff` | 207416 | `ba53cc13090bb69b235a99573849c4ee8c535ac7d42b573d1ac030eb3b648750` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/change-set.pre-r3.diff` | 175686 | `9bf1847f086660e348e884359afc11a279f56a28c95b16d0670123df2214d92b` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/DESIGN.txt` | 71911 | `38a7f85e18417e8a11e663d89545b1e1e3fbb0397f572f94a395de1b5b054c81` | 被受版本控制文档引用其路径 |
| `archive/gate-index.txt` | 2882 | `6023c5ef720b3669dd37039054512808789bb906345f0a0a379b12e578dc7409` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree.txt` | 2881 | `7cfdf298ac0bcf7eb42e104999d4571d34ff4f71b3103dee88a0619153b36e9f` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input.diff` | 159652 | `bb7cb8f337ec044f7eb7304d008c6307e940169076b9036ae1ea679f36103bd7` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-pack-r2.txt` | 7315 | `daaec3e648190aed1abcd18c6a4f8a5329e98744b8ea452af979967384261916` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-pack-r3.txt` | 7466 | `dc7d9e00f519fbd2511230dbf3225f45efb82926829788cf8d1a8c56c5f316f7` | 被受版本控制文档引用其路径 |
| `archive/selftest.txt` | 21644 | `a75e7b7a10a6e0f52a15f6f6b112de7026b391a9a82aab7c890619317aa484ff` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/snapshot-flake.txt` | 58 | `a983c18c7be5da3ae926a1089bccb4cc333ef85bbf7a0512ed2dd568a6a421cc` | 被受版本控制文档引用其路径 |
| `archive/v3.txt` | 1514 | `83fa96b58ebc330c0f32b2d8c630da5841eaa6be87b9ba573a4b373f395e7f8c` | 被受版本控制文档引用其路径 |
| `archive/v4.txt` | 1156 | `040e73869c5e618fd0a6b717b9e543ea8db90ee706847981acc776618ecb70b0` | 被受版本控制文档引用其路径 |

## 排除清单（未归档）

| 相对路径 | 字节 | sha256 | 未归档理由 |
|---|---|---|---|
| `/_before/docs/DELIVERY_DIRECTIVE.md` | 185729 | `0c1a67f8927294895be04e27758d931127043d61982afe630289aa878faf6d35` | 中间转储（未被引用的过程输出） |
| `/_before/docs/DELIVERY_DIRECTIVE_EN.md` | 226751 | `a3d4cb85f6a7b277efca9e3c20e69c1d457d3d682dcdfb2c02dd17a5f815fbe2` | 中间转储（未被引用的过程输出） |
| `/_before/internal/adapters/agent_runner_postflight.go` | 30716 | `d82d79a6f824ff3a0429042772b4b3522a040f2d75e1a8f0e49fc534292890f1` | 中间转储（未被引用的过程输出） |
| `/_before/internal/adapters/agent_runner_postflight_test.go` | 34389 | `9ab4963cf83bad2419ca0113bdbd06c013f1923e1e9933da0f2848062ff9714d` | 中间转储（未被引用的过程输出） |
| `/_before/internal/console/cli.go` | 23928 | `84f976883a2ff6a3ac19343b54abb56c9535763bcbf75bb6853b62c790dc84be` | 中间转储（未被引用的过程输出） |
| `/_before/internal/console/runtime.go` | 3401 | `c84279b56ac4a68bb992ad42f713e40942125fee0aab2a76334c572e542e2610` | 中间转储（未被引用的过程输出） |
| `/_before/internal/console/runtime_agent.go` | 30576 | `1e7cf34daffe78b5c003cf954382fbbcbbb956e9bd5ba38693fb9acf97421a78` | 中间转储（未被引用的过程输出） |
| `/_before/internal/console/runtime_agent_e2e_test.go` | 23846 | `53212a1a0c3a4f5fee7c0ee0d91dde94b4ea4ad8ff7ff8392cc7676fa51cf213` | 中间转储（未被引用的过程输出） |
| `/_before/internal/console/runtime_noop.go` | 4612 | `f089577a3317c1309f309410dd309beb0553de9fe774457212069d8ddc0bc5f3` | 中间转储（未被引用的过程输出） |
| `/_r3_after/internal/adapters/agent_runner_postflight.go` | 32593 | `9fd49c44b4240fb1b1deb8dd75921d2bc3d22eb87c0d1ae0f86fe875b005f523` | 中间转储（未被引用的过程输出） |
| `/_r3_after/internal/adapters/agent_runner_postflight_test.go` | 37977 | `db13b6cc3b3a4a6aebf0ede1a0c889a1f23b62eb2f129b041f6090dd31479fc2` | 中间转储（未被引用的过程输出） |
| `/_r3_after/internal/console/runtime_agent.go` | 36791 | `1dcc259ab9794bbfff195d2f364b9fb49db14a3a89eb127105cd904919d3d3cb` | 中间转储（未被引用的过程输出） |
| `/_r3_after/internal/console/runtime_agent_e2e_test.go` | 34384 | `50b712af05c089a53e3d728f2750d06c34e9b15f7fa713a9ff394cc553d6ef91` | 中间转储（未被引用的过程输出） |
| `/_r3_before/internal/adapters/agent_runner_postflight.go` | 30716 | `d82d79a6f824ff3a0429042772b4b3522a040f2d75e1a8f0e49fc534292890f1` | 中间转储（未被引用的过程输出） |
| `/_r3_before/internal/adapters/agent_runner_postflight_test.go` | 34907 | `eb27de601008bd5a46a1d82382cba8dc918a50919d357237dacd5e498cc013a6` | 中间转储（未被引用的过程输出） |
| `/_r3_before/internal/console/runtime_agent.go` | 32649 | `5677976ce4658a497ca96fc83869f7c5ee4c8ed17bfe80a4b318f95ef2ddf980` | 中间转储（未被引用的过程输出） |
| `/_r3_before/internal/console/runtime_agent_e2e_test.go` | 27307 | `fd848e256ead2cfd1d80f44a7b78b66f287e885703c96e1160b110c21f1f41cd` | 中间转储（未被引用的过程输出） |
| `/_r4_before/internal/adapters/agent_runner_postflight.go` | 32593 | `9fd49c44b4240fb1b1deb8dd75921d2bc3d22eb87c0d1ae0f86fe875b005f523` | 中间转储（未被引用的过程输出） |
| `/_r4_before/internal/adapters/agent_runner_postflight_test.go` | 37977 | `db13b6cc3b3a4a6aebf0ede1a0c889a1f23b62eb2f129b041f6090dd31479fc2` | 中间转储（未被引用的过程输出） |
| `/_r4_before/internal/console/runtime_agent.go` | 36791 | `1dcc259ab9794bbfff195d2f364b9fb49db14a3a89eb127105cd904919d3d3cb` | 中间转储（未被引用的过程输出） |
| `/_r4_before/internal/console/runtime_agent_e2e_test.go` | 34384 | `50b712af05c089a53e3d728f2750d06c34e9b15f7fa713a9ff394cc553d6ef91` | 中间转储（未被引用的过程输出） |
| `/a-final-test.txt` | 1156 | `fb12fab2891ca74a0b53d5fe8f4449929619fd2cefc53f20c4c0d0c9cc377d52` | 中间转储（未被引用的过程输出） |
| `/baseline-test.txt` | 1156 | `a5f3124d41a0ea4896d58bbb44af2d2091960b67d89b361edba3dea671e87449` | 中间转储（未被引用的过程输出） |
| `/build-after-probe.txt` | 0 | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` | 中间转储（未被引用的过程输出） |
| `/ci-ubuntu-fail.txt` | 3276 | `1fa5786967f778e194829e0db082d9478c58470ed40145991fa6a97d192304ac` | 中间转储（未被引用的过程输出） |
| `/ci-ubuntu2.txt` | 77955 | `bc7414a0c01a16c3dcc55972b831438bf31c099a44c8f8ebc94e5ba1002e8cb5` | 中间转储（未被引用的过程输出） |
| `/ci-watch.txt` | 22373 | `6f9c0d71b7e4d4bc469515b88331c9ce8ebceb4dae21a860ae3567b601e14bc2` | 中间转储（未被引用的过程输出） |
| `/ci-watch2.txt` | 20087 | `0a023e645aa8827d2c2756b819035b2f011b71e4c7b6284403bc204d5dbf5597` | 中间转储（未被引用的过程输出） |
| `/ci-win-log.txt` | 74410 | `115018c18c7f21c1df4e067c3cb722b1603fbd3fe325c55fe5416ac5ce476a66` | 中间转储（未被引用的过程输出） |
| `/ci-win2.txt` | 73141 | `b714c0e66af8b97ea67c374b5a04ed7239269745dfde8af2161e9da4ad507eca` | 中间转储（未被引用的过程输出） |
| `/commit.txt` | 680 | `bd2b05e680435ebeb441855fd5a1b2242c4f823faa8f20507cbec19619694e5d` | 中间转储（未被引用的过程输出） |
| `/commit2.txt` | 164 | `7f81690ae5bbeef6b27d987b68e694d74a8b17aaaa3f1366843e4b186eb58bf8` | 中间转储（未被引用的过程输出） |
| `/contracts-excerpt.txt` | 38256 | `5ef2eb347608fdc60bcc8dac45bb466cfeb3ff408374158a2392e51adeb00904` | 中间转储（未被引用的过程输出） |
| `/diag4b.js` | 1597 | `b336de2fa3b10b62d483cd742c7cc84e8f06ef3510b15ba2664360c1a6bec22e` | 落地与修复脚本残留 |
| `/findg6.js` | 1179 | `58f79acad96be4f413a0af7dca9d3360335d7ccda7525d1a251a68c118608bc4` | 落地与修复脚本残留 |
| `/fixm16.js` | 1628 | `33b2f0647f4b7b34eceb365386dc3cd629bcf4f0ebcf91073cf5483513bf3d2d` | 落地与修复脚本残留 |
| `/fixm16.txt` | 138 | `f08c5a3e0aaf7acca218a0f4dbff2c35af16457e3aba65c68281d28b3c1e29b7` | 中间转储（未被引用的过程输出） |
| `/l1l5-probe.txt` | 51 | `177e6a3209d70ae26bb295a33f821449a8e750d10e51c0f8b6b7a37f84731ed4` | 中间转储（未被引用的过程输出） |
| `/mkfix.js` | 8335 | `be1b6024b3d7e6a6ecb641cd131593fd5d798ce20f8a621f91cee0336b1e30b8` | 落地与修复脚本残留 |
| `/mkfull.js` | 1403 | `a0f666661ac6826b8f2c8fa49bd0d12a4b8a715e42657116dd430c3bde067ae5` | 落地与修复脚本残留 |
| `/mkfull.txt` | 482 | `cb42c7f6e755768c5eed6ccfd2373e02dde13de21f6316a81afc52b6720cbe1a` | 中间转储（未被引用的过程输出） |
| `/mkledger.js` | 11920 | `d7b6ccc1f9d588218594e811cc9394123760c861869ac89c7b4e54d4dc64621a` | 落地与修复脚本残留 |
| `/mkpack.js` | 2147 | `1c0ffc911cf0609f7a4aa3b669ec0597172d53e0ddeca1d34c154e41b40a3cf6` | 落地与修复脚本残留 |
| `/mkr3diff.js` | 5332 | `505bcafdde1b7398282325760e532bf507a69f9afbb24f8cd7631c72b0d5bcb6` | 落地与修复脚本残留 |
| `/mkr4diff.js` | 4201 | `9d4d973edb875f59ccde79cf27b4ad8b724ef9714890d1714a05f013e7985d57` | 落地与修复脚本残留 |
| `/mkremediation.js` | 2071 | `6dd3e44c984cb224e3980aea1810dd8221f0e628c9348240c7c6ef47784ec458` | 落地与修复脚本残留 |
| `/mkreport.js` | 53696 | `a84ab254c4911b937a188c4038f5392d4c7dfff93b0db10fc08d9410cc23a8d7` | 落地与修复脚本残留 |
| `/mkreport.txt` | 288 | `fb1a4a4f04e3cb7f10583043b1ac7758c7626049c160f9a594e9494c75fa5534` | 中间转储（未被引用的过程输出） |
| `/mutate.js` | 17280 | `7a24d54d4d9d4b3f8086f040ab8be1ec1dfc1f441aa70caf50a61fdfa4416390` | 落地与修复脚本残留 |
| `/mutation.txt` | 7015 | `a8a9a95822c285021fe05da3fd3b9fb52ecf3a8f6a9cc2885115cd398775ddfd` | 中间转储（未被引用的过程输出） |
| `/norm.js` | 1246 | `76928a41d09da80ec2ab71890dde267e07e65522eafa36cde2837006ad5e710b` | 落地与修复脚本残留 |
| `/norm.txt` | 239 | `472c704510821db04f6d55bdf34abf9fb43d396a0a93b61539bbf126a5a5d336` | 中间转储（未被引用的过程输出） |
| `/parseprobe.txt` | 1003 | `5c8dca8bb713af93724d00ab7a627decd6d03afdd77bd3908dbf1a4ea441fc30` | 中间转储（未被引用的过程输出） |
| `/pre-review-r1.md` | 13775 | `56c6f1b7de1528ad7eecbec07a0e4351c451403b6ccc758c1f0f479f303abc2f` | 中间转储（未被引用的过程输出） |
| `/pre-review-r2.md` | 5365 | `d803dd13c706b3113ff48c1948533def8e004beead7282905a4ff5bbfdcedc39` | 中间转储（未被引用的过程输出） |
| `/push.txt` | 71 | `c6ee763915e277d2f6f1f810e2d6dd944ec8b5e4e096caf129af90bb6c6ea58c` | 中间转储（未被引用的过程输出） |
| `/push2.txt` | 71 | `57f982bb73292395389a5ea24ce00ed3fa3de2e23988a635ffcc3da7555da8c9` | 中间转储（未被引用的过程输出） |
| `/r3-final-test.txt` | 1156 | `44f6d7f4873625a0781adb1feccd4b03a98e79aab606cfa4a72c83ed0eafec3b` | 中间转储（未被引用的过程输出） |
| `/r3-remediation.diff` | 35493 | `e165b28e98b63d2694a7035395b7f9889f57a37bbbd3d154fe3cb5e2dc30dd53` | 中间转储（未被引用的过程输出） |
| `/r3.txt` | 1156 | `186fdf9440894471d453367f835c0168f27d7ea65518f81c1f38ff6b64e844cc` | 中间转储（未被引用的过程输出） |
| `/r3diff.txt` | 659 | `cf2b5bd5743e58da0050a750269fe14e0c7ca2cc5cba0f064298647dabc2a0d0` | 中间转储（未被引用的过程输出） |
| `/r4-final-test.txt` | 1156 | `6a78fd42c58d59353d9d99272b1571f432e61a7e77c996657f5b16a057d1235e` | 中间转储（未被引用的过程输出） |
| `/r4-remediation.diff` | 20758 | `ca941465c45797480de2ef3bc1118caa04df03dfa0956e529c9b6ca8f361543b` | 中间转储（未被引用的过程输出） |
| `/r4diff.txt` | 332 | `00669e236f3a975c1389165acd1838879f0dc4f41ff22ffe4449370908b6e665` | 中间转储（未被引用的过程输出） |
| `/r4snapshot.txt` | 365 | `0576e2e58f6e1577dae0b46d6c9529821d7341a07a567695269e547107a4dae0` | 中间转储（未被引用的过程输出） |
| `/remediation.diff` | 10883 | `a21f8c4033c09c560e85588e7c8908523cfcbf56b76dd7abea84d2169ea9a566` | 中间转储（未被引用的过程输出） |
| `/v2.txt` | 1156 | `e07a35eed9d7510ad4836944a3b1d4afee2dde3b5f44e5e6aef6b0ab920c1850` | 中间转储（未被引用的过程输出） |
| `/v5.txt` | 1156 | `82fcbd1afc5398f8c7f51d4f82e671c744c8f7fb64c3c9a2fd34ee1ad55dccc2` | 中间转储（未被引用的过程输出） |
| `/v6.txt` | 1155 | `ab3457708ea15af13d96dcd42d521b317339cc702b2f9290ffb1bf08647c98f5` | 中间转储（未被引用的过程输出） |
| `/verify-l1l5.js` | 2203 | `fe36dd75b3fc1708ed47cccce79c62511a5ec1443e2c52a98dc9322e7b37202d` | 落地与修复脚本残留 |
| `/verify.txt` | 1156 | `8a0a0bf133e972e4cb88daa94cbe21f2fdcdc3fc0ede6678130d29d1129d22c3` | 中间转储（未被引用的过程输出） |
| `/wrapup.js` | 2379 | `149ef856a4e209a32325416d12e70ba68f95ee2bbb2f8fe5fdaa02f149d471ea` | 落地与修复脚本残留 |

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
| `DESIGN.md` | `archive/DESIGN.txt` | `d0c3c6f7406ecab457c85ec22d9072bd44b8815ce7c91ada91809213f735e06f` | `38a7f85e18417e8a11e663d89545b1e1e3fbb0397f572f94a395de1b5b054c81` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `review-pack-r2.md` | `archive/review-pack-r2.txt` | `daaec3e648190aed1abcd18c6a4f8a5329e98744b8ea452af979967384261916` | `daaec3e648190aed1abcd18c6a4f8a5329e98744b8ea452af979967384261916` | 后缀 `.md` → `.txt` |
| `review-pack-r3.md` | `archive/review-pack-r3.txt` | `dc7d9e00f519fbd2511230dbf3225f45efb82926829788cf8d1a8c56c5f316f7` | `dc7d9e00f519fbd2511230dbf3225f45efb82926829788cf8d1a8c56c5f316f7` | 后缀 `.md` → `.txt` |

## 对应引用点

| 引用文件 | 行号 |
|---|---|
| `docs/ADR_REGISTER.md` | 29 |
| `docs/ADR_REGISTER_EN.md` | 30 |
| `docs/DELIVERY_DIRECTIVE.md` | 121 / 897 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 121 / 897 |
| `docs/validation/evidence/directive-history/assembly-min/archive/DESIGN.txt` | 6 / 24 / 304 / 347 |
| `docs/validation/evidence/directive-history/assembly-min/archive/review-pack-r2.txt` | 12 / 13 / 14 / 15 / 16 / 17 / 18 |
| `docs/validation/evidence/directive-history/assembly-min/archive/review-pack-r3.txt` | 12 / 13 / 14 / 15 / 16 / 17 / 19 / 20 |
| `docs/validation/evidence/directive-history/dce/archive/change-set-r2.diff` | 8 / 20 |
| `docs/validation/evidence/directive-history/dce/archive/change-set.diff` | 8 / 20 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/classification.txt` | 39 / 40 / 41 / 42 / 43 / 44 / 45 / 46 / 47 / 48 / 49 / 50 / 51 / 52 / 53 / 54 / 55 / 56 / 57 / 58 / 59 / 60 / 61 / 62 / 63 / 64 / 65 / 66 / 67 / 68 / 69 / 70 / 71 / 72 / 73 / 74 / 75 / 76 / 77 / 78 / 79 / 80 / 81 / 82 / 83 / 84 / 85 / 86 / 87 / 88 / 89 / 90 / 91 / 92 / 93 / 94 / 95 / 96 / 97 / 98 / 99 / 100 / 101 / 102 / 103 / 104 / 105 / 106 / 107 / 108 / 109 / 110 / 111 / 112 / 113 / 114 / 115 / 116 / 117 / 118 / 119 / 120 / 121 / 122 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/facts.txt` | 113 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/facts2.txt` | 12 / 43 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/inventory.txt` | 56 / 57 / 58 / 59 / 60 / 61 / 62 / 431 / 433 / 446 / 447 / 448 / 471 / 472 / 473 / 1156 / 1157 / 1158 / 1159 / 1160 / 1161 / 1162 / 1164 / 1166 / 1167 / 1168 / 1169 / 1170 / 1171 / 1172 / 1173 / 1174 / 1176 / 1178 / 1179 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/measure.txt` | 33 / 38 / 40 / 41 / 43 |
| `docs/validation/t027-assembly-min.md` | 12 / 32 / 75 / 101 / 132 / 154 / 183 |
| `docs/validation/t027-assembly-min_EN.md` | 12 / 32 / 75 / 101 / 132 / 154 / 183 |
