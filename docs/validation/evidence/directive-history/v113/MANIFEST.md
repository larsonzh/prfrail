# 归档 MANIFEST · v113

- **源路径**：`tmp/v113/`
- **归档片**：`DIRECTIVE-TMP-LIFECYCLE`（v1.23，2026-10-04）
- **归档口径**：§1.5.1 的 ② 类「早期证据评估归档」——保留证据条目（被受版本控制文档引用其路径者 ∪ 切片级文档），排除运行时残留与一次性脚本
- **归档条目**：76 条 / 2,602,143 字节（归档载荷根 = `archive/`；逐条 sha256 见 `SHA256SUMS.txt`）
- **排除条目**：70 条（逐条见下「排除清单」）
- **未入库条目**：0 条（逐条见下「未入库文件清单」）
- **字节保真**：归档副本的**字符内容不改**；因 §6.3 的作用域（G4a 扫全部变更文件）与仓库编码约定，副本按**归档后类型**归一编码——`.ps1` 带 BOM，其余不带 BOM，一律 LF；逐条变更种类与**原始字节 sha256** 见下「后缀与编码映射」节
- **后缀口径**：归档副本的 `.md` 一律改存为 `.txt`（本单元 4 条），以避开只扫 `.md` 的文档判据；映射逐条见下节
- **编码归一计数**：后缀改写 4 条 / BOM 剔除 0 条 / BOM 补入 0 条 / LF 归一 17 条

## 归档条目清单

| 相对路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/capture.js` | 2897 | `b5162802decfcbd4a5b05d70bf6be813edd1571214e477683d67f4fd0e20afb6` | 被受版本控制文档引用其路径 |
| `archive/enc-tree-r10.txt` | 421 | `fe9278b9be49065d4207f77c2567d91e415ea7d150fa29d9c171134d865d5d19` | 被受版本控制文档引用其路径 |
| `archive/enc-tree-r7.txt` | 421 | `fe9278b9be49065d4207f77c2567d91e415ea7d150fa29d9c171134d865d5d19` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/enc-tree-r8.txt` | 421 | `fe9278b9be49065d4207f77c2567d91e415ea7d150fa29d9c171134d865d5d19` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/enc-tree-r9.txt` | 421 | `fe9278b9be49065d4207f77c2567d91e415ea7d150fa29d9c171134d865d5d19` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/evidence-selfcheck-r10.txt` | 411 | `73edef2f15974e6c3727f3b4f7636c9b7bfa579ed8f5e2970f218fb2386a869d` | 被受版本控制文档引用其路径 |
| `archive/g6-5-r1-r10.txt` | 240 | `14bd2517e5f6d8273b4b742fad987f790272d03644210fe96ace9500a86ae9ef` | 被受版本控制文档引用其路径 |
| `archive/g6-6-r1-r10.txt` | 222 | `926935b22d4a7ec5014d0d80a89a770450dd4ea87948c017accb1811a4f9c9aa` | 被受版本控制文档引用其路径 |
| `archive/gate-ci-after-tester.txt` | 2451 | `7a55bbe84a594efc3c9bc034b16f4fd390fc52f91e777ce0889f81b02e2bdb48` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-index-after-tester.txt` | 2533 | `f720259d237967868bff84e229b8806fb41f5e9415471b68afba4ddd2d4af551` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-index-r1-r10.txt` | 2902 | `c99867df739ec6cf6a4c46374427863ba1e0d9d0547b93215295774683eaf05e` | 被受版本控制文档引用其路径 |
| `archive/gate-index-r1-r2.txt` | 2533 | `f720259d237967868bff84e229b8806fb41f5e9415471b68afba4ddd2d4af551` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-index-r1-r3.txt` | 2533 | `f720259d237967868bff84e229b8806fb41f5e9415471b68afba4ddd2d4af551` | 被受版本控制文档引用其路径 |
| `archive/gate-index-r1-r7.txt` | 2902 | `1f6ecdbe4d6ba8f300b9a4036971d20dfe3c493052431f946f01820566c0448b` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-index-r1-r8.txt` | 2902 | `bc30253c8bf899aaf15f60038eb113815e6ffbb2638753e7c8fda2f42516b52b` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-index-r1-r9.txt` | 2902 | `a0c796426aae08c676d9fcc8558eba65503d5fccbacee4259fec335e62ed5d9f` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-index-r1.txt` | 2533 | `f720259d237967868bff84e229b8806fb41f5e9415471b68afba4ddd2d4af551` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-index.txt` | 2244 | `ded10bd0a4f04808af5126b16b81f6ca1439bb7a0b01eb36d015d8e59e0b47b0` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-selftest.txt` | 17809 | `466088f10c713379f2dca9dec12c35feeb8a0f798085d7803f81206caf09c590` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-final.txt` | 2549 | `74ab5d641ce4c5c483be704a1d29c50e87b2995d939c9c8b09ef8c53f7660693` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-ob45.txt` | 2549 | `5e382cd64e18de35f9dd88cadf25d37accf315b93aeb81948eeef8d7ec6c9863` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r1-r10.txt` | 2901 | `50690ea44db0eaf08b3310fd1356287861a7eb773ac2128a49008d7ab6ee293a` | 被受版本控制文档引用其路径 |
| `archive/gate-tree-r1-r2.txt` | 2549 | `2e73f94a7b6e9658a31d999abe5ee1f904394f6cff9adb1296df1571273353c8` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r1-r3.txt` | 2549 | `5e382cd64e18de35f9dd88cadf25d37accf315b93aeb81948eeef8d7ec6c9863` | 被受版本控制文档引用其路径 |
| `archive/gate-tree-r1-r7.txt` | 2901 | `0786523d8fad81575d49ac44074669c286fbe96b8e6595a54ccd000f675306f0` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r1-r8.txt` | 2901 | `dec929f565fd19c61cd3df10e27a2d7bc720d928befd4abb994edd1aa63f4fa6` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r1-r9.txt` | 2901 | `bf00ad2494680356c61ec6eda9479b29f2b159f34b88a99c023dd1ea886020d2` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r1.txt` | 2549 | `fb565358515f2116b5daa964977d71ee0c423347f45014fc68660edf4e5af02f` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r2.txt` | 2549 | `2e73f94a7b6e9658a31d999abe5ee1f904394f6cff9adb1296df1571273353c8` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r4.txt` | 2549 | `3f9e38007b6ebdf26875aab45849137b7a19eba077a4fc3b3e417dde1ec438d1` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r5.txt` | 7601 | `23b414130ce6efc4666642703f9711bd010a8fd1ab34dc261196e6423a802453` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree-r6.txt` | 7601 | `a274027fe71d616853c68ac29e5dc00198be9508c21d30b8b10d31c7bf9dd71c` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-tree.txt` | 2547 | `486088aa035f31c8f5bf7e92940480eb6472a075d9d0c3f3b7ad336500e6dd75` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/mirror-r1.txt` | 128 | `23cee53e68a26b6d2ae5e11de6d5b99618d397e22cae40685e7432040dd86c6b` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/mirror-r10.txt` | 128 | `ff9efcb90cf9a894d7c875e641b26e66df6b285070e0d9347d3740dd259ee646` | 被受版本控制文档引用其路径 |
| `archive/mirror-r2.txt` | 128 | `13cd8e213e99b2c6f4b733a3879eb7ebcb0fce87f919803398f4e4b86123e51f` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/mirror-r3.txt` | 128 | `ff9efcb90cf9a894d7c875e641b26e66df6b285070e0d9347d3740dd259ee646` | 被受版本控制文档引用其路径 |
| `archive/mirror-r7.txt` | 128 | `ff9efcb90cf9a894d7c875e641b26e66df6b285070e0d9347d3740dd259ee646` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/mirror-r8.txt` | 128 | `ff9efcb90cf9a894d7c875e641b26e66df6b285070e0d9347d3740dd259ee646` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/mirror-r9.txt` | 128 | `ff9efcb90cf9a894d7c875e641b26e66df6b285070e0d9347d3740dd259ee646` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/mirror.txt` | 128 | `23cee53e68a26b6d2ae5e11de6d5b99618d397e22cae40685e7432040dd86c6b` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/PLAN.txt` | 35931 | `58531bb0b2b1f390c8f0a4e05f09494d5408f3f7b6071b6668c6ddf61c79ec46` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/remediation.diff` | 53987 | `19679e6f20d56eff74569d67ba481caa3a2efdb66944bcf72f61aac3893e7131` | 被受版本控制文档引用其路径 |
| `archive/review-input-r2-r10.diff` | 325749 | `66c77d74978ae630908251297b10fa62f40ef8ac926f2fafea60fcdadc9a44e5` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2-r2.diff` | 221141 | `3b3bee9d7db1924bfe8ba94b482638c5c40e0f7ceb778528faa52d516806cd59` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2-r3.diff` | 223199 | `24ca51c9daff08e6186b0130a8aa85ecd812b952f3a5524e7fcb5439d9256e52` | 被受版本控制文档引用其路径 |
| `archive/review-input-r2-r7.diff` | 302060 | `a220bd656650aeac6e5febdfbffab5c256d52041eee8e57e8be6555c0fff8045` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2-r8.diff` | 322503 | `9eb0df211fc021012bb7c785f8c7fba5aecc4532cd694ceb0cca9b15a3808c92` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2-r9.diff` | 322304 | `10600ef3d4e58097eda1ab6834ce92e32b4fbbf808a5c03f34dbb6421f053df5` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2.diff` | 216016 | `77a3266f2ec8cdc2cc3a0b2b2e9068e0e7f2bb80829eea726b14b26619675d2c` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2.stat-r10.txt` | 762 | `0ef6b6283596765c982c618a4953523c10b8b75380925e02c1b3e4cbb0a2acf7` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2.stat-r2.txt` | 136 | `7e1b2c5a1d78a2eea16feb550c7a90e5355f6835279b6c3fc44864e2ea8bed7a` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2.stat-r3.txt` | 136 | `f5f569113a217f39eb9b8df7d178da562adbc88da6af74bd05e2f5f59e54460f` | 被受版本控制文档引用其路径 |
| `archive/review-input-r2.stat-r7.txt` | 762 | `59e6fbaaf4f8be949be909314c268326374ba715eb5b4d4d5cfde2166bb57e47` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2.stat-r8.txt` | 762 | `7d12c3d1dfd6ea501c07af6057b337826745ec4f4cf43f61b8eee99c9adb3aeb` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2.stat-r9.txt` | 762 | `7d12c3d1dfd6ea501c07af6057b337826745ec4f4cf43f61b8eee99c9adb3aeb` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input-r2.stat.txt` | 136 | `1d82f4226de14c782e67c793938e5dd9f4ee3e6c23b190eeea58719f2ca5cd52` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-input.diff` | 205333 | `1eb961e906ca57eaa994e594c4c3e1646736602c972beec01581cf41edc2b898` | 被受版本控制文档引用其路径 |
| `archive/review-input.stat.txt` | 136 | `934a381a21ade9e254d9dcf06a8e77545d8ee4966e9332470e16e204462cacf9` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/review-r1.txt` | 12360 | `b30e9c75dfb6c873e54a2a428a0806627fbacff11ebccbc61bc567c249898de7` | 被受版本控制文档引用其路径 |
| `archive/review-r2.txt` | 8136 | `94da3b51e199c9d8f3310744afa618bec92967730614d9ef2969fe54201aa7d9` | 被受版本控制文档引用其路径 |
| `archive/review-r3.txt` | 8929 | `2b947bc0065c15c16d24cf8914bc260b28cf0b8bd9465bead69d4f977c51be2b` | 被受版本控制文档引用其路径 |
| `archive/selftest-after-tester.txt` | 21948 | `df86641fa3c9f44670018c4ebec9190764e0876614f064af1e262f472265fbe6` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/selftest-master-r10.txt` | 21644 | `a75e7b7a10a6e0f52a15f6f6b112de7026b391a9a82aab7c890619317aa484ff` | 被受版本控制文档引用其路径 |
| `archive/selftest-master-r2.txt` | 21644 | `a75e7b7a10a6e0f52a15f6f6b112de7026b391a9a82aab7c890619317aa484ff` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/selftest-master-r3.txt` | 21644 | `a75e7b7a10a6e0f52a15f6f6b112de7026b391a9a82aab7c890619317aa484ff` | 被受版本控制文档引用其路径 |
| `archive/selftest-master-r7.txt` | 21644 | `a75e7b7a10a6e0f52a15f6f6b112de7026b391a9a82aab7c890619317aa484ff` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/selftest-master-r8.txt` | 21644 | `a75e7b7a10a6e0f52a15f6f6b112de7026b391a9a82aab7c890619317aa484ff` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/selftest-master-r9.txt` | 21644 | `a75e7b7a10a6e0f52a15f6f6b112de7026b391a9a82aab7c890619317aa484ff` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/selftest-master.txt` | 21644 | `a75e7b7a10a6e0f52a15f6f6b112de7026b391a9a82aab7c890619317aa484ff` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/selftest-ob45.txt` | 21948 | `df86641fa3c9f44670018c4ebec9190764e0876614f064af1e262f472265fbe6` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/selftest-r1.txt` | 21948 | `df86641fa3c9f44670018c4ebec9190764e0876614f064af1e262f472265fbe6` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/selftest-r2.txt` | 21948 | `df86641fa3c9f44670018c4ebec9190764e0876614f064af1e262f472265fbe6` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/verify-v113-r10.txt` | 349 | `359c4b8db6ff8e1213eb446bd3b7f380dcdf255523e973041619d34b7916c8a1` | 被受版本控制文档引用其路径 |
| `archive/verify-v113-r3.txt` | 349 | `359c4b8db6ff8e1213eb446bd3b7f380dcdf255523e973041619d34b7916c8a1` | 被受版本控制文档引用其路径 |
| `archive/verify-v113.js` | 2987 | `138d48dc5def4661cd28d7bab64b48e8b72a2a3440530a7017332433f506bcd6` | 被受版本控制文档引用其路径 |

## 排除清单（未归档）

| 相对路径 | 字节 | sha256 | 未归档理由 |
|---|---|---|---|
| `/apply-r1-fixes.js` | 9944 | `ab9b5603ed8a0c760d2316509d45660affac597cc40cffcb3a32d7af58c86462` | 落地与修复脚本残留 |
| `/apply-v113.js` | 15522 | `32664ec29b02ef2c25744ecb47e515544c3da5264e59fac3e4d4255700f166b8` | 落地与修复脚本残留 |
| `/blocks.txt` | 28746 | `97d0d73d7c301f9498f000b600ded666cb0a4d6967759e81c99d73534e3c2b6b` | 中间转储（未被引用的过程输出） |
| `/c0h-note-cn.txt` | 448 | `57927632f833db2ad6104ef90357767d55581b5264308c494e27d5d6bd7b5871` | 中间转储（未被引用的过程输出） |
| `/c0h-note-en.txt` | 529 | `b338a6c04b5cb8c5d8cb02979e50df8af7c55f7a092de0ff7ee5369f901ff549` | 中间转储（未被引用的过程输出） |
| `/capture-r10.log` | 465 | `8dfd1c0ff5f73b3f25261509a5cf497e342c289105238889c24a258e9883643d` | 运行时状态与日志 |
| `/capture-r8.log` | 456 | `b83faf359b69ecfd3aedd29921d71bea1927edd23a387563a7e8c6db131e3ec0` | 运行时状态与日志 |
| `/capture-r9.log` | 456 | `de882af2ddf3333181d4d56b53abe442374e6ddf480128c9d98a5a65aad1ce4b` | 运行时状态与日志 |
| `/evidence-selfcheck-r7.txt` | 411 | `73edef2f15974e6c3727f3b4f7636c9b7bfa579ed8f5e2970f218fb2386a869d` | 中间转储（未被引用的过程输出） |
| `/evidence-selfcheck-r8.txt` | 411 | `73edef2f15974e6c3727f3b4f7636c9b7bfa579ed8f5e2970f218fb2386a869d` | 中间转储（未被引用的过程输出） |
| `/evidence-selfcheck-r9.txt` | 411 | `73edef2f15974e6c3727f3b4f7636c9b7bfa579ed8f5e2970f218fb2386a869d` | 中间转储（未被引用的过程输出） |
| `/extract-review-r1.js` | 1935 | `8c94caec07f9020040a9e8ef6378b968476588d153bc4cdac26b4a04d0c2d4a5` | 落地与修复脚本残留 |
| `/extract-review.js` | 2230 | `4e0e4d1ec3dcc9f8da1f5708a0ac0eed8a7efa2224941faba3df2b8031c33857` | 落地与修复脚本残留 |
| `/fence-check.js` | 704 | `2879d91e73301449324e609f2044747b4e84f72b7a05475bedbd6e627205606d` | 落地与修复脚本残留 |
| `/final-check.txt` | 27449 | `5dadf4f8bbf461d5157b312fa39de527d5dd6bd8867c77bdc1e40aff5f5a45d3` | 中间转储（未被引用的过程输出） |
| `/fix-c0i-blank-lines.js` | 1714 | `34db27f86045a9a2a896e6f28ed99befb498557570551bdad44f3e84e8a198b2` | 落地与修复脚本残留 |
| `/fix-en-disposed-mirror.js` | 2354 | `034858f575d1d0a6ca1d6e82ddd620994c5216bcc5772a83fef6b5729b320a2d` | 落地与修复脚本残留 |
| `/fix-en-disposed-v113.js` | 2145 | `7195dfef6d281b7e6d54c2ce9e2a9c6215aa430398aa642e4960396bafe50c03` | 落地与修复脚本残留 |
| `/fix-g63-blank-lines.js` | 2128 | `b3bb13e7b3bf7b9318a99f62d75fd5a8033491b8226c89c279371e6afcf7d7a4` | 落地与修复脚本残留 |
| `/fix-g64-v113.js` | 1650 | `e5d6475f61a327a581e5b016fe0d7bbf084a98df54df9d47df5f009c9558ac67` | 落地与修复脚本残留 |
| `/g6-5-r1-r2.txt` | 240 | `2219f38296f493a21d5cc3b66193bc917d79828cb99d1d967d2b743073d1ed25` | 中间转储（未被引用的过程输出） |
| `/g6-5-r1-r3.txt` | 240 | `2219f38296f493a21d5cc3b66193bc917d79828cb99d1d967d2b743073d1ed25` | 中间转储（未被引用的过程输出） |
| `/g6-5-r1-r7.txt` | 240 | `14bd2517e5f6d8273b4b742fad987f790272d03644210fe96ace9500a86ae9ef` | 中间转储（未被引用的过程输出） |
| `/g6-5-r1-r8.txt` | 240 | `14bd2517e5f6d8273b4b742fad987f790272d03644210fe96ace9500a86ae9ef` | 中间转储（未被引用的过程输出） |
| `/g6-5-r1-r9.txt` | 240 | `14bd2517e5f6d8273b4b742fad987f790272d03644210fe96ace9500a86ae9ef` | 中间转储（未被引用的过程输出） |
| `/g6-5-r1.txt` | 240 | `2219f38296f493a21d5cc3b66193bc917d79828cb99d1d967d2b743073d1ed25` | 中间转储（未被引用的过程输出） |
| `/g6-5-range-9f68a52.txt` | 268 | `dbccdd00ffb5abf55d325792a018479a73be6b825a2434f977a790dfb604a8da` | 中间转储（未被引用的过程输出） |
| `/g6-5-range-b738e6b.txt` | 267 | `1277b9c643666bbe7e41348e350cd660e65fddcde9202d3cad93b7e4c006ee68` | 中间转储（未被引用的过程输出） |
| `/g6-5.txt` | 242 | `7dccaa0e99ff3582841018851c14526d4e37aedec7a91f642061cbf1f0f3f1e3` | 中间转储（未被引用的过程输出） |
| `/g6-6-r1-r2.txt` | 222 | `47f9b8b6cf9dc0f6b6f27e7729b2a6afd8369e1b4655c02fee0148bb947b49a0` | 中间转储（未被引用的过程输出） |
| `/g6-6-r1-r3.txt` | 222 | `caa4233626bd35a1939e0b94d7ad16d1ab5ba3b9138189cb94fb04c4ad7b7004` | 中间转储（未被引用的过程输出） |
| `/g6-6-r1-r7.txt` | 222 | `926935b22d4a7ec5014d0d80a89a770450dd4ea87948c017accb1811a4f9c9aa` | 中间转储（未被引用的过程输出） |
| `/g6-6-r1-r8.txt` | 222 | `926935b22d4a7ec5014d0d80a89a770450dd4ea87948c017accb1811a4f9c9aa` | 中间转储（未被引用的过程输出） |
| `/g6-6-r1-r9.txt` | 222 | `926935b22d4a7ec5014d0d80a89a770450dd4ea87948c017accb1811a4f9c9aa` | 中间转储（未被引用的过程输出） |
| `/g6-6-r1.txt` | 222 | `951caf8092452d6063b99a214c72b56dac14b3948d5d86b339727fff09e41a10` | 中间转储（未被引用的过程输出） |
| `/g6-6.txt` | 224 | `0b786c406c65218c2b39de5dd7da44f67722fb80b749707c23d2bbbb89cb1ce0` | 中间转储（未被引用的过程输出） |
| `/g6-range-b738e6b-after-tester.txt` | 908 | `2a2318c3bab64fc52517b8ddd0719c1bf832a2a62d3a2a358199e8013781a66b` | 中间转储（未被引用的过程输出） |
| `/g6-range-b738e6b.txt` | 908 | `2a2318c3bab64fc52517b8ddd0719c1bf832a2a62d3a2a358199e8013781a66b` | 中间转储（未被引用的过程输出） |
| `/map-status.js` | 3293 | `167d8a64988c96debf75937aebb4d6dcc8e078ae2079518776bb04aaaefbd39a` | 落地与修复脚本残留 |
| `/normalize-lf.js` | 1463 | `310d047d173911ed0423587e661fdc570258404e83d2d7d8e9a9928e447175e4` | 落地与修复脚本残留 |
| `/ob39-40-41-cn.md` | 2491 | `5c38b9045126d10fc63f771e7d30212a045e3ac3fcb975d4378ea68e3e944ab4` | 中间转储（未被引用的过程输出） |
| `/ob39-40-41-en.md` | 2972 | `e543142540b074dfed82d77ee8b51c88a4c174aa2851b7a02ab301b3bc5c2f97` | 中间转储（未被引用的过程输出） |
| `/ob43-44-cn.md` | 1604 | `e8e8bd6a6cdb95323e3b2138fb91f043634a66b997e46e994ea6fcc6d228b966` | 中间转储（未被引用的过程输出） |
| `/ob43-44-en.md` | 2016 | `ad54032682479d3036027a6cd9cd7e751593aca77419692ea92c7fdd4b68b973` | 中间转储（未被引用的过程输出） |
| `/ob45-cn.md` | 912 | `9740127bb96aa889d4a6ae9265f6fbd19585ef91222154b5654803b71c3aca1b` | 中间转储（未被引用的过程输出） |
| `/ob45-en.md` | 1137 | `9289cc4fabd714aead7656b99677ed8b332580113279b23d1820b56b93991e47` | 中间转储（未被引用的过程输出） |
| `/probe-anchors.js` | 1778 | `4283fc70819494d57182a81d79644c4730aac1fe2cd9b5bc6b62bffc574ae8f8` | 落地与修复脚本残留 |
| `/probe-g66.js` | 1459 | `801335ae05f48c0b015fc89ea0d3db14a4d8e0a3e6d4c83d48e7fe7c8c92def3` | 落地与修复脚本残留 |
| `/probe-new-criteria.js` | 5054 | `223e052129a4d554555a33a52764c82e95db11f5e12fe843e684cc23b69758b7` | 落地与修复脚本残留 |
| `/probe-tester-fixtures.js` | 3210 | `22744d47449f5f38db6950f507d93921bbcb93f78d0e4d72555663fabcca0137` | 落地与修复脚本残留 |
| `/probe-tester-mutation.js` | 1756 | `dd9677f5c1333a452e667a4daa98ddfedf62ea18b4b12b4b55795b02bcdce07c` | 落地与修复脚本残留 |
| `/r1-cn-draft.md` | 1293 | `01c244c882f751ecbd8c11a4cc7d40a63f84db4d2b7619c11305318c922b15b6` | 中间转储（未被引用的过程输出） |
| `/range-after-tester.txt` | 1069 | `7a4d7fc095a3877e81d4e841de8b878614c2157717b9efa24109679a2328543e` | 中间转储（未被引用的过程输出） |
| `/register-ob39-41.js` | 4503 | `0c0b8793a27a21f249bab425ff07c99ba5831f1f65423c005b6a9c2297177c31` | 落地与修复脚本残留 |
| `/register-ob43-44.js` | 6994 | `3d55c094d9113472ede6e5d92b4d065bf22d877d9496ab8614e359fe5651342d` | 落地与修复脚本残留 |
| `/register-ob45.js` | 5625 | `100fbdd208beacc84dd5a50ed54a1ffe6cd52d9efa40f61013c9634218edd150` | 落地与修复脚本残留 |
| `/scan-g6-4.js` | 1319 | `ac2844384896934de4ea1bbad404a823d3de0e39222c86d8086dcff06082dccc` | 落地与修复脚本残留 |
| `/sec14-archival-note-cn.md` | 803 | `2ec186c7fea2179b73baabcfaf7119f295fd99433fbcc60e67df232219cd4539` | 中间转储（未被引用的过程输出） |
| `/sec6-third-application-cn.md` | 728 | `7029da56973e88eb0573bcbb1199bf0202cd60bc61ce139795320c888c83f5d5` | 中间转储（未被引用的过程输出） |
| `/sec93-row-cn.md` | 421 | `4c861323a4e0f2286882a5118120c557ef0c60308ee0a53227a3f73d3f9381ca` | 中间转储（未被引用的过程输出） |
| `/sec93-row-en.md` | 506 | `9cf72e2b10f51f0eebe8788945663f98e1508d38db749b9db2c32bd7b977647a` | 中间转储（未被引用的过程输出） |
| `/sentence-cn.md` | 123 | `a0610a7cf394f08b54095e9bfe8d30f1d0cbd9c15d04179bc13b9e30c1f4fb7f` | 中间转储（未被引用的过程输出） |
| `/sentence-en.md` | 167 | `11404b992402630c78fad54da6b21707be0ad08d862967086ec91c073ed5cde9` | 中间转储（未被引用的过程输出） |
| `/status-map.txt` | 6453 | `eadbfa430b81bc501182f9c94808567f88545943035d46f761ff5927b909800e` | 中间转储（未被引用的过程输出） |
| `/summary-r9.js` | 1570 | `9715b496dcb1ee7c4ce9a5ffa486766d3a86d7b30f23255b386a22672680d355` | 落地与修复脚本残留 |
| `/verify-v113-r2.txt` | 349 | `c0094350136e2768a77bda3d4ebae7c4e84fb1fa5a407d4d1cc1d410c107b780` | 中间转储（未被引用的过程输出） |
| `/verify-v113-r7.txt` | 349 | `359c4b8db6ff8e1213eb446bd3b7f380dcdf255523e973041619d34b7916c8a1` | 中间转储（未被引用的过程输出） |
| `/verify-v113-r8.txt` | 349 | `359c4b8db6ff8e1213eb446bd3b7f380dcdf255523e973041619d34b7916c8a1` | 中间转储（未被引用的过程输出） |
| `/verify-v113-r9.txt` | 349 | `359c4b8db6ff8e1213eb446bd3b7f380dcdf255523e973041619d34b7916c8a1` | 中间转储（未被引用的过程输出） |
| `/verify-v113.txt` | 349 | `825c0a526da444161f1319ea336f4cf6d3ed26b527aa01b7a9f37a5a7ff3aadb` | 中间转储（未被引用的过程输出） |

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
| `PLAN.md` | `archive/PLAN.txt` | `58531bb0b2b1f390c8f0a4e05f09494d5408f3f7b6071b6668c6ddf61c79ec46` | `58531bb0b2b1f390c8f0a4e05f09494d5408f3f7b6071b6668c6ddf61c79ec46` | 后缀 `.md` → `.txt` |
| `capture.js` | `archive/capture.js` | `b56dda409dc9cb2c503780ed41c5880ba097f0f8ab30f39196bacbe7c921b324` | `b5162802decfcbd4a5b05d70bf6be813edd1571214e477683d67f4fd0e20afb6` | LF 归一 |
| `gate-ci-after-tester.txt` | `archive/gate-ci-after-tester.txt` | `e8acbf8ccef477c484250972e9bd093c3767249bf6bf277ecf57996eac3118f9` | `7a55bbe84a594efc3c9bc034b16f4fd390fc52f91e777ce0889f81b02e2bdb48` | LF 归一 |
| `gate-index-after-tester.txt` | `archive/gate-index-after-tester.txt` | `62240cdd4546d36506266e856de7387dfe393454020efd78039bc38e1d1b958c` | `f720259d237967868bff84e229b8806fb41f5e9415471b68afba4ddd2d4af551` | LF 归一 |
| `gate-index.txt` | `archive/gate-index.txt` | `6becbaa77de3d1a4ee283de7568c529b099afbdeb0d6807bbaff4d515d7168c3` | `ded10bd0a4f04808af5126b16b81f6ca1439bb7a0b01eb36d015d8e59e0b47b0` | LF 归一 |
| `gate-selftest.txt` | `archive/gate-selftest.txt` | `7d133fa4319b94c0b7d84280151d5ba6c6ed3a9d31985da7170133917165b435` | `466088f10c713379f2dca9dec12c35feeb8a0f798085d7803f81206caf09c590` | LF 归一 |
| `gate-tree-final.txt` | `archive/gate-tree-final.txt` | `b21ca4e216326e7722e73fe28b41afffebfc99cb050350db5eb8efd9fff82f57` | `74ab5d641ce4c5c483be704a1d29c50e87b2995d939c9c8b09ef8c53f7660693` | LF 归一 |
| `gate-tree-ob45.txt` | `archive/gate-tree-ob45.txt` | `331e5353adca442cd33e5e1cf39e50a21c2b11674e44ab84d63e99b0adc44772` | `5e382cd64e18de35f9dd88cadf25d37accf315b93aeb81948eeef8d7ec6c9863` | LF 归一 |
| `gate-tree-r2.txt` | `archive/gate-tree-r2.txt` | `f6766fdf27657ebd11bbed616933a1efc6147d2e9968a98edd5d29b2a26da70b` | `2e73f94a7b6e9658a31d999abe5ee1f904394f6cff9adb1296df1571273353c8` | LF 归一 |
| `gate-tree-r4.txt` | `archive/gate-tree-r4.txt` | `873e14f9aacfd2a04edf6f39343f51f7424d041d6ca4cd7c3a53c58f7d7e265e` | `3f9e38007b6ebdf26875aab45849137b7a19eba077a4fc3b3e417dde1ec438d1` | LF 归一 |
| `gate-tree-r5.txt` | `archive/gate-tree-r5.txt` | `4a1f9b1313e7a16ff5b71065fde64393d55be37a76d5e4367cff8b5c49a91dff` | `23b414130ce6efc4666642703f9711bd010a8fd1ab34dc261196e6423a802453` | LF 归一 |
| `gate-tree-r6.txt` | `archive/gate-tree-r6.txt` | `ccd80fd097786d813399c646ee088ac831554ad4fb04246980867dccba9e3cb4` | `a274027fe71d616853c68ac29e5dc00198be9508c21d30b8b10d31c7bf9dd71c` | LF 归一 |
| `gate-tree.txt` | `archive/gate-tree.txt` | `c6cb4bec42b5bcb83de554cb336caf176aaa5e28b946d18dbac9da43b714428d` | `486088aa035f31c8f5bf7e92940480eb6472a075d9d0c3f3b7ad336500e6dd75` | LF 归一 |
| `mirror-r1.txt` | `archive/mirror-r1.txt` | `9e8eeeb62cf26d8dce9c10710cda89f75b4e01cf1bf1ad1acbbf2d8cfcff7337` | `23cee53e68a26b6d2ae5e11de6d5b99618d397e22cae40685e7432040dd86c6b` | LF 归一 |
| `review-r1.md` | `archive/review-r1.txt` | `b30e9c75dfb6c873e54a2a428a0806627fbacff11ebccbc61bc567c249898de7` | `b30e9c75dfb6c873e54a2a428a0806627fbacff11ebccbc61bc567c249898de7` | 后缀 `.md` → `.txt` |
| `review-r2.md` | `archive/review-r2.txt` | `94da3b51e199c9d8f3310744afa618bec92967730614d9ef2969fe54201aa7d9` | `94da3b51e199c9d8f3310744afa618bec92967730614d9ef2969fe54201aa7d9` | 后缀 `.md` → `.txt` |
| `review-r3.md` | `archive/review-r3.txt` | `2b947bc0065c15c16d24cf8914bc260b28cf0b8bd9465bead69d4f977c51be2b` | `2b947bc0065c15c16d24cf8914bc260b28cf0b8bd9465bead69d4f977c51be2b` | 后缀 `.md` → `.txt` |
| `selftest-ob45.txt` | `archive/selftest-ob45.txt` | `0428f81ae7dfcff5c0a0516a8db0e797c4e3678fe277ff558b8ba20e0537e167` | `df86641fa3c9f44670018c4ebec9190764e0876614f064af1e262f472265fbe6` | LF 归一 |
| `selftest-r1.txt` | `archive/selftest-r1.txt` | `0428f81ae7dfcff5c0a0516a8db0e797c4e3678fe277ff558b8ba20e0537e167` | `df86641fa3c9f44670018c4ebec9190764e0876614f064af1e262f472265fbe6` | LF 归一 |
| `selftest-r2.txt` | `archive/selftest-r2.txt` | `0428f81ae7dfcff5c0a0516a8db0e797c4e3678fe277ff558b8ba20e0537e167` | `df86641fa3c9f44670018c4ebec9190764e0876614f064af1e262f472265fbe6` | LF 归一 |
| `verify-v113.js` | `archive/verify-v113.js` | `164b15735bd264ce831b5cbb9182308c8ce2aebe1bd1a30677ccce617287a544` | `138d48dc5def4661cd28d7bab64b48e8b72a2a3440530a7017332433f506bcd6` | LF 归一 |

## 对应引用点

| 引用文件 | 行号 |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 120 / 863 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 120 / 863 |
| `docs/validation/directive-review-saturation.md` | 28 / 30 / 68 / 99 / 128 / 129 / 130 / 131 / 156 / 173 / 174 / 176 / 177 / 178 / 179 / 180 / 204 / 206 / 207 / 208 / 277 / 285 / 290 / 291 / 292 / 293 / 294 / 295 / 305 / 306 / 311 / 313 / 316 / 317 / 318 / 321 / 322 / 338 / 339 / 351 / 352 |
| `docs/validation/directive-review-saturation_EN.md` | 28 / 30 / 68 / 99 / 128 / 129 / 130 / 131 / 156 / 173 / 174 / 176 / 177 / 178 / 179 / 180 / 204 / 206 / 207 / 208 / 277 / 285 / 290 / 291 / 292 / 293 / 294 / 295 / 305 / 306 / 311 / 313 / 316 / 317 / 318 / 321 / 322 / 338 / 339 / 351 / 352 |
| `docs/validation/evidence/directive-history/dr1/archive/review-input-r2.diff` | 47 / 79 |
| `docs/validation/evidence/directive-history/dr1/archive/review-input.diff` | 23 / 55 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/classification.txt` | 1478 / 1479 / 1480 / 1481 / 1482 / 1483 / 1484 / 1485 / 1486 / 1487 / 1488 / 1489 / 1490 / 1491 / 1492 / 1493 / 1494 / 1495 / 1496 / 1497 / 1498 / 1499 / 1500 / 1501 / 1502 / 1503 / 1504 / 1505 / 1506 / 1507 / 1508 / 1509 / 1510 / 1511 / 1512 / 1513 / 1514 / 1515 / 1516 / 1517 / 1518 / 1519 / 1520 / 1521 / 1522 / 1523 / 1524 / 1525 / 1526 / 1527 / 1528 / 1529 / 1530 / 1531 / 1532 / 1533 / 1534 / 1535 / 1536 / 1537 / 1538 / 1539 / 1540 / 1541 / 1542 / 1543 / 1544 / 1545 / 1546 / 1547 / 1548 / 1549 / 1550 / 1551 / 1552 / 1553 / 1554 / 1555 / 1556 / 1557 / 1558 / 1559 / 1560 / 1561 / 1562 / 1563 / 1564 / 1565 / 1566 / 1567 / 1568 / 1569 / 1570 / 1571 / 1572 / 1573 / 1574 / 1575 / 1576 / 1577 / 1578 / 1579 / 1580 / 1581 / 1582 / 1583 / 1584 / 1585 / 1586 / 1587 / 1588 / 1589 / 1590 / 1591 / 1592 / 1593 / 1594 / 1595 / 1596 / 1597 / 1598 / 1599 / 1600 / 1601 / 1602 / 1603 / 1604 / 1605 / 1606 / 1607 / 1608 / 1609 / 1610 / 1611 / 1612 / 1613 / 1614 / 1615 / 1616 / 1617 / 1618 / 1619 / 1620 / 1621 / 1622 / 1623 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/facts2.txt` | 11 / 21 / 42 / 52 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/inventory.txt` | 209 / 210 / 211 / 212 / 213 / 214 / 215 / 216 / 217 / 218 / 219 / 220 / 221 / 222 / 223 / 224 / 225 / 226 / 227 / 228 / 229 / 230 / 231 / 232 / 233 / 234 / 235 / 236 / 237 / 238 / 239 / 240 / 241 / 242 / 243 / 244 / 245 / 246 / 247 / 248 / 441 / 442 / 443 / 444 / 445 / 461 / 486 / 809 / 810 / 811 / 812 / 813 / 814 / 815 / 816 / 817 / 818 / 819 / 820 / 821 / 822 / 823 / 824 / 825 / 826 / 827 / 828 / 829 / 830 / 831 / 832 / 833 / 834 / 835 / 836 / 837 / 838 / 839 / 840 / 841 / 842 / 843 / 844 / 845 / 846 / 847 / 848 / 849 / 850 / 851 / 852 / 853 / 854 / 855 / 856 / 857 / 858 / 859 / 860 / 861 / 862 / 863 / 864 / 865 / 866 / 867 / 868 / 869 / 870 / 871 / 872 / 873 / 874 / 875 / 876 / 877 / 878 / 879 / 880 / 881 / 882 / 883 / 884 / 885 / 886 / 887 / 888 / 889 / 890 / 891 / 892 / 893 / 894 / 895 / 896 / 897 / 898 / 899 / 900 / 901 / 902 / 903 / 904 / 905 / 906 / 907 / 908 / 909 / 910 / 911 / 912 / 913 / 914 / 915 / 916 / 917 / 918 / 919 / 920 / 921 / 922 / 923 / 924 / 925 / 926 / 927 / 928 / 929 / 930 / 931 / 932 / 933 / 934 / 935 / 936 / 937 / 938 / 939 / 940 / 941 / 942 / 943 / 944 / 945 / 946 / 947 / 948 / 949 / 950 / 951 / 952 / 953 / 954 / 955 / 956 / 957 / 958 / 959 / 960 / 961 / 962 / 963 / 964 / 965 / 966 / 967 / 968 / 969 / 970 / 971 / 972 / 973 / 974 / 975 / 976 / 1054 / 1055 / 1056 / 1057 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/measure.txt` | 29 / 30 / 31 / 32 / 34 / 35 / 37 / 39 |
| `docs/validation/evidence/directive-history/v113/archive/capture.js` | 11 / 13 / 28 / 29 / 30 / 31 / 32 / 34 / 42 / 43 / 44 / 55 / 69 |
| `docs/validation/evidence/directive-history/v113/archive/remediation.diff` | 2 / 4 / 78 / 80 |
| `docs/validation/evidence/directive-history/v113/archive/review-input-r2-r10.diff` | 218 / 491 / 580 / 582 / 620 / 651 / 680 / 681 / 682 / 683 / 708 / 725 / 726 / 728 / 729 / 730 / 731 / 732 / 756 / 758 / 759 / 760 / 829 / 836 / 841 / 842 / 843 / 844 / 845 / 846 / 856 / 857 / 862 / 864 / 867 / 868 / 869 / 872 / 873 / 889 / 890 / 902 / 903 / 938 / 940 / 978 / 1009 / 1038 / 1039 / 1040 / 1041 / 1066 / 1083 / 1084 / 1086 / 1087 / 1088 / 1089 / 1090 / 1114 / 1116 / 1117 / 1118 / 1187 / 1194 / 1199 / 1200 / 1201 / 1202 / 1203 / 1204 / 1214 / 1215 / 1220 / 1222 / 1225 / 1226 / 1227 / 1230 / 1231 / 1247 / 1248 / 1260 / 1261 / 1287 / 1289 |
| `docs/validation/evidence/directive-history/v113/archive/review-input-r2-r2.diff` | 218 / 490 |
| `docs/validation/evidence/directive-history/v113/archive/review-input-r2-r3.diff` | 218 / 491 |
| `docs/validation/evidence/directive-history/v113/archive/review-input-r2-r7.diff` | 218 / 491 / 580 / 582 / 617 / 663 / 664 / 665 / 689 / 693 / 694 / 695 / 696 / 701 / 704 / 789 / 794 / 799 / 800 / 801 / 802 / 803 / 804 / 814 / 815 / 820 / 822 / 825 / 826 / 827 / 830 / 831 / 847 / 848 / 860 / 861 / 896 / 898 / 933 / 979 / 980 / 981 / 1005 / 1009 / 1010 / 1011 / 1012 / 1017 / 1020 / 1105 / 1110 / 1115 / 1116 / 1117 / 1118 / 1119 / 1120 / 1130 / 1131 / 1136 / 1138 / 1141 / 1142 / 1143 / 1146 / 1147 / 1163 / 1164 / 1176 / 1177 / 1203 / 1205 |
| `docs/validation/evidence/directive-history/v113/archive/review-input-r2-r8.diff` | 218 / 491 / 580 / 582 / 620 / 651 / 676 / 677 / 678 / 679 / 704 / 721 / 722 / 724 / 725 / 726 / 727 / 728 / 752 / 754 / 755 / 756 / 825 / 830 / 835 / 836 / 837 / 838 / 839 / 840 / 850 / 851 / 856 / 858 / 861 / 862 / 863 / 866 / 867 / 883 / 884 / 896 / 897 / 932 / 934 / 972 / 1003 / 1028 / 1029 / 1030 / 1031 / 1056 / 1073 / 1074 / 1076 / 1077 / 1078 / 1079 / 1080 / 1104 / 1106 / 1107 / 1108 / 1177 / 1182 / 1187 / 1188 / 1189 / 1190 / 1191 / 1192 / 1202 / 1203 / 1208 / 1210 / 1213 / 1214 / 1215 / 1218 / 1219 / 1235 / 1236 / 1248 / 1249 / 1275 / 1277 |
| `docs/validation/evidence/directive-history/v113/archive/review-input-r2-r9.diff` | 218 / 491 / 580 / 582 / 620 / 651 / 676 / 677 / 678 / 679 / 704 / 721 / 722 / 724 / 725 / 726 / 727 / 728 / 752 / 754 / 755 / 756 / 825 / 830 / 835 / 836 / 837 / 838 / 839 / 840 / 850 / 851 / 856 / 858 / 861 / 862 / 863 / 866 / 867 / 883 / 884 / 896 / 897 / 932 / 934 / 972 / 1003 / 1028 / 1029 / 1030 / 1031 / 1056 / 1073 / 1074 / 1076 / 1077 / 1078 / 1079 / 1080 / 1104 / 1106 / 1107 / 1108 / 1177 / 1182 / 1187 / 1188 / 1189 / 1190 / 1191 / 1192 / 1202 / 1203 / 1208 / 1210 / 1213 / 1214 / 1215 / 1218 / 1219 / 1235 / 1236 / 1248 / 1249 / 1275 / 1277 |
| `docs/validation/evidence/directive-history/v113/archive/review-r1.txt` | 5 / 8 |
| `docs/validation/evidence/directive-history/v113/archive/review-r2.txt` | 5 |
| `docs/validation/evidence/directive-history/v113/archive/review-r3.txt` | 4 / 5 / 6 / 7 / 8 / 9 / 19 / 20 / 25 / 27 / 30 / 31 / 32 / 35 / 36 / 52 / 53 |
| `docs/validation/evidence/directive-history/v113/archive/verify-v113.js` | 11 |
| `docs/validation/evidence/directive-review-saturation-freedom-list.md` | 19 / 21 |
