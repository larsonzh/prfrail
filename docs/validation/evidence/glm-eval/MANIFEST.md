# 冻结证据包 · 切片 `DIRECTIVE-GLM-EVAL`

| 项 | 值 |
|---|---|
| 包根 | `docs/validation/evidence/glm-eval/` |
| 切片 | `DIRECTIVE-GLM-EVAL`（评估片；**只评估、不入白名单**） |
| 归档日期 | 2026-10-05 |
| 载荷范围 | `archive/glm-eval/**`（会话工作目录全量：bundle / 两腿输出 / 评分与污染脚本 / 自测器与夹具 / 计量规程 / 交接文档 / 方案书与勘误）；`archive/glm-eval-out/**`（两腿原始输出 + 自测报告时间点快照）；`archive/glm-eval-batch/**`（批量腿 a 的输入与钥匙）；`archive/tmp-root/**`（承重脚本与 `PROMPT-EXACT.txt`）；`archive/tooling/**`（本包的可复现工具：打包器与 MANIFEST 模板） |
| 条目数 | 77 |
| 哈希清单 | `SHA256SUMS.txt`（`<hash><2 spaces><path>` 形态，POSIX 分隔符） |
| 归档口径 | 逐字节复制，**未做任何编码或行尾转换**：载荷原始即为「非 `.ps1` ⇒ 无 BOM」且一律 LF，已满足 §1.5.1「归档后类型归一」。故每条**原始字节 sha256 = 归档字节 sha256**（见下「后缀与编码映射」节） |
| 自洽核验 | 本包须过 `node tools/gates/gate.js --check=G8-b --scope=tree`（归档条目集合 ≡ `SHA256SUMS.txt` 条目集合 ≡ 包内载荷集合，且逐条重算 sha256 一致） |

## 归档条目清单

| 归档路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/tmp-root/PROMPT-EXACT.txt` | 11578 | `bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6` | 冻结评估物 |
| `archive/tmp-root/append-rulings.js` | 3651 | `413b22d6cf61b436b6490b1826b82050a0947f024bdbee328e02e0e5a9175200` | 冻结评估物 |
| `archive/tmp-root/batch-a-gen.js` | 3393 | `b18f3da55757d22296d70cb4d281ef8fa81b3d55216425b31c7333669aef9b03` | 冻结评估物 |
| `archive/tmp-root/batch-d-gen.js` | 860 | `28710cabbf0e4f2ae23b4cf0e9363b7b508d88c00af277029b723c850a67d49f` | 冻结评估物 |
| `archive/tmp-root/fill-dataset.js` | 8161 | `c308ad2654e1abfaa6bca3821f65bed2e4350c86460dab5c4c29c5815292cd33` | 冻结评估物 |
| `archive/tmp-root/fix-bom-and-append.js` | 6160 | `59c9b862ead14223f599ad8dc85cc87b02b3dd01017f84fc5125a943af9c3372` | 冻结评估物 |
| `archive/tmp-root/freeze.js` | 2244 | `38b078c7a9277ad4820e2c157dcc1e5401aac55c86018dec5c4f7b5a31a3d4f5` | 冻结评估物 |
| `archive/glm-eval-batch/a-input.txt` | 1886 | `ea435e99732ebd7436422e0f84baf2e3bafe31e91964af2440010150af9ce564` | 冻结评估物 |
| `archive/glm-eval-batch/a-key.json` | 869 | `3639a989ab0f21a3bd2696cda1b616cad92ece2911bbf869ff4ac66ec979b5a3` | 冻结评估物 |
| `archive/glm-eval-out/leg-control.txt` | 6591 | `2fa1d9eeedae134c758174fff900311983a52b0eda8470889f4ff5d67f463a7f` | 冻结评估物 |
| `archive/glm-eval-out/leg-glm.txt` | 13234 | `54c895f9df821231d9049ac7819077b44335db83383b266d86f3c28fee137b5b` | 冻结评估物 |
| `archive/glm-eval-out/selftest-report-20261005.txt` | 26531 | `7b06d9d757b2d0e7fd00216245dbcef9405f9273afe3230e33d81fc48be9950d` | 冻结评估物 |
| `archive/glm-eval/BRIEFING.md` | 18971 | `18867c95816abfa977c30eec31a74ef54ef17405b6c90e91caa6f05415df3ae8` | 冻结评估物 |
| `archive/glm-eval/FROZEN-MANIFEST.txt` | 9307 | `bda60cc1707b6c945ff2a7ebc48bb7e987af5c36125d276cc9758ffde613a347` | 冻结评估物 |
| `archive/glm-eval/GLM-USAGE-READING.txt` | 2077 | `023426c8020dfc565c09d5176f6120458c1529d0862892f2ae4c2e2d5bff356f` | 冻结评估物 |
| `archive/glm-eval/PLAN-ERRATA.txt` | 4320 | `a3c083559bfe4ca4c866d2260acab3e483b0623241a0900721397e2565b59dbf` | 冻结评估物 |
| `archive/glm-eval/PLAN-PROVENANCE.txt` | 1500 | `66f1a3f5e81fe4619c16cf340f4add5e41a3eff2a5af0d7a09fa593874e18dfa` | 冻结评估物 |
| `archive/glm-eval/PROBE-NOTES.txt` | 2023 | `7f3dc09fae9292ffbbae5989b375f69fa2285a12c3c83fb0401863553ed817e6` | 冻结评估物 |
| `archive/glm-eval/TEST-REPORT-DIRECTIVE-GLM-EVAL.txt` | 30646 | `de5dac345258d0760780db839e74f06bee5e65aeb1372bc1166e29a4dd911472` | 冻结评估物 |
| `archive/glm-eval/bundle/00-plan.md` | 1184 | `70076661d86615d8dfe180aa5884ec9618338ec9496b3893900a749c092f36b8` | 冻结评估物 |
| `archive/glm-eval/bundle/01-clause-cn.txt` | 2346 | `d8fe4bc82df774535c47c39406c96fee1325c87f3153b5d72ee2a378ff03e525` | 冻结评估物 |
| `archive/glm-eval/bundle/02-diff.js` | 6451 | `c939bc9e756696c463571c981a7d3756c87563de913b9c4c58faebdc4d5ff405` | 冻结评估物 |
| `archive/glm-eval/bundle/03-prompt.txt` | 1597 | `3b7a3580edcecb6a92ccb42f8acf37b469283136cb0c8d8d1b2b24d0b54fd758` | 冻结评估物 |
| `archive/glm-eval/contam.js` | 2550 | `dfcb76e737b6ba5dab0453977fa6c1307bb2391e7890e5d1bd9aca3525099455` | 冻结评估物 |
| `archive/glm-eval/hash.js` | 1378 | `ff608449685e5150e5338a8bcad5e36e0e8910f172ed77182b8e1bf3c4c073b3` | 冻结评估物 |
| `archive/glm-eval/key.json` | 3843 | `6fccbe461286abae63e74e52f64d8130f4fa979bb6d27ba71651fad6c3339536` | 冻结评估物 |
| `archive/glm-eval/metering/README.md` | 2058 | `150c1a9bc4e001ab5cadf6e294162e8789048f4cb1ef83af1241baeaca4543a1` | 冻结评估物 |
| `archive/glm-eval/plan-v1.md` | 12860 | `8be638d27cbf467f0b337de7e4362bad6b090b43e876377d1ad7ec7429181d04` | 冻结评估物 |
| `archive/glm-eval/score.js` | 12339 | `2c4052fa8e10273b09465b668a09b0c1e4402065576ec48ca4357341d27ee051` | 冻结评估物 |
| `archive/glm-eval/selftest/fixture-contam-fail.txt` | 263 | `5d2e0631d60db305d9aae6f698797561e18a51e722a8b9fc7adc2d24f115ec04` | 冻结评估物 |
| `archive/glm-eval/selftest/fixture-score.txt` | 666 | `ec698d06c1bcc6a8a8af274c6050121d3cf81c43ce777430bf23c1e828fb9c11` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-cjk1.txt` | 63 | `cd8405ae534bc76db879b992bf33c20974d0d57c707b3e47588b2e94e5f22e06` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-cjk2.txt` | 72 | `ec3962b4d8469a54a19ae86342ac2dd5a68bcee0da01d75578da63406262b1a4` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-directive.txt` | 75 | `4590eda262ce6f0d6e2da605a6df6e90cf1bd4d549e197763f8013cac0d4b5b8` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-g8a.txt` | 55 | `95abf3393bb78966d5c7c8dec217b7dad0c2c9ad86948729b0bbc90d722df0bb` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-hash025.txt` | 58 | `41225e5f54a0c0c8b4b5755a018c6e9763f07981ccc69711779d700c51d1d48f` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-hash0c1.txt` | 58 | `bb04ad75582f95d630e5444f0cd251601262d1f0354f24a62a64c6aff9ed1982` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-ob79.txt` | 59 | `0e3bd8ab01eb11f47d188610cb6b482500d52bc238580eadaeb11512f57bd608` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-ob790.txt` | 63 | `1cb7f0e8e075332adbad194915545d2c06cbd74a9eeaeccf1fe2b994374e4bfd` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-s10.txt` | 57 | `b49bf889a887b750956ee13959169f5408114a4b583150637a0f9b80ae8fdac1` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-s123.txt` | 58 | `83edc34853693971f9630e89fe974c535b58e50648824355dbd27f61423ed23b` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-s8.txt` | 53 | `d0e43de3e6f42731db40751a9466ea394e4215cb08115f87883350fcf01e3394` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/boundary-s9.txt` | 56 | `420c1322871c27fe0039df471e8d364073a74310a6dc71cea22718a0f4d3752e` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/hash-bom.txt` | 9 | `42c1e65b2c948bb754efb6ac171319d6e97ecb3d9afd4f20bd91b3ded25183c0` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/hash-nobom.txt` | 6 | `5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-abstain-line.txt` | 91 | `3a59038bad20c5129765f4d5a53f325efb0be564a28a8aa1ab2d463fead2a593` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-abstain.txt` | 203 | `03f1a97f4d3a0eb521b24b18cc77b392a4a00e35bf241829757c4cdbff6f26d9` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-acceptlines.txt` | 384 | `6e408e1468ed8d9e20981031410908b31b5dc542c7260c663975dabb355b99aa` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-clean.txt` | 420 | `8b60d780b1637921cd1a1d1d9f9c9206e7fc478637da7df1808c42013a76a1da` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-comma.txt` | 77 | `de25a865a50ebe6840019f8346219c6635377ae325334c6b60a837f1eaa5f2e4` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-contam.txt` | 161 | `3354cb1d27957a0728688be21ef04d564732e8ffa0804790769f382588f8cb44` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-dedup.txt` | 140 | `a46b8c016876135a2de5aefcfcf1e302a08c9dffe0fd11eddd66c89fb86e8d76` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-empty.txt` | 43 | `2b544e496f9c88a8ae0d337f57024c38b97f23b8e3602584619258e12094bae2` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-extra.txt` | 86 | `4b55d162afd1515fdd7152c5ad792e5267092fb3cc97afbe9012650c54e571cd` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-fps.txt` | 177 | `7bf1ff635e21b61307bdcda4905c71f127efef379423197c1b9f8a103651840a` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-fullwidth.txt` | 81 | `5f71171e8255f945698124eb63524603f44d6d3da7c3f2b1dac6683b10764387` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-noconcl.txt` | 44 | `5e38192d6c73ed45728378fd96869c05779a1f485224af6403edb88c6a2ba70c` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-outside.txt` | 79 | `7cedeee512ac9b9a21b819f40b23bacc31090f6c5a144ebdf3742a4fcf1dbad2` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-perfect.txt` | 591 | `36d40ca7061d8b94448670f6947ceeb4c7d9047a55c717d1874e07b5a3858672` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-qcheck.txt` | 126 | `0a58d82387e3e944df592663d3e9bdd23554aef9e1ba0c937160a36e10a34e2f` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-range-miss.txt` | 72 | `dfd77a61368639d9e6d94b436df0266ead8486679be730fc22e2e1926cf252f0` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-range.txt` | 93 | `5f041561b478710b0a00b778fb949c1740b327323defc1487b2a77f86ac18f9c` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-sev-cn.txt` | 107 | `ea59fe812117f0418abee250687f6fc01822a49dbbbfbe4ae0db4eeef0423b9f` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-sev-exact.txt` | 56 | `275b6eb589a231ef2815b47b935d53f9b28b1db2102d28b547a851b927b08fdf` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-sev-off1.txt` | 52 | `cd23868bb537cf0343d9d4d7d969ec0997fc08a56c18857d0164150a6d466c01` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-sev-off2.txt` | 51 | `c5d7324412988b23d991a7bcf74c34335121b78ce5a61bde5cbadb2a35fac140` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-sev-table.txt` | 129 | `e3abb722801742b015337b7f1ca3aa021413a6202d916029cbfecc7859775a69` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-slash.txt` | 130 | `fb46611ed342e0e9fbbd9d0f9b7a4bf7fcfc9551075f3b63549d453c9b40e253` | 冻结评估物 |
| `archive/glm-eval/selftest/fixtures/leg-trailing-ws.txt` | 46 | `0c18703e7c83bafc2d6e46ce01eaac6dca8fabbd8096ae29dacf5856e8e8d133` | 冻结评估物 |
| `archive/glm-eval/selftest/run-selftest.js` | 61825 | `2fc646ad9464133b0f6b9211e02d9815414e9b640155e0d9af352651661e0e49` | 冻结评估物 |
| `archive/glm-eval/selftest/selftest-report.txt` | 26531 | `7b06d9d757b2d0e7fd00216245dbcef9405f9273afe3230e33d81fc48be9950d` | 冻结评估物 |
| `archive/tmp-root/norm-score.js` | 1303 | `f07365c0dbddcff94a8f6371003728aea8f6d9efb48cc4c6a556f72d4877322e` | 冻结评估物 |
| `archive/tmp-root/persist.js` | 6995 | `ddd0956fe8ba5e0119054e0e7ab84a0be8f4c531f6f06aa20c8e77041601030a` | 冻结评估物 |
| `archive/tmp-root/recon-briefing.js` | 2814 | `60c394ed3aa87eff6531cd386eb1bc898b9d82aa8b184aacde510072386764bf` | 冻结评估物 |
| `archive/tmp-root/save-leg.js` | 1477 | `96bb7500b08e4bf51e8aa1caa886789b394f2b28cd0ffe005a6890df2c81465a` | 冻结评估物 |
| `archive/tooling/pack-glm-eval.js` | 3067 | `67638513ca06cc83e38b2738334464aa8272c025a090be37d2f1adf8ec366165` | 冻结评估物 |
| `archive/tooling/manifest-template.md` | 3049 | `74ed1697069f448d405735aaf429afc454892d794cdc848ba265278e8667786f` | 冻结评估物 |

## 后缀与编码映射

| 原始路径 | 归档路径 | 原始字节 sha256 | 归档字节 sha256 | 归一 |
|---|---|---|---|---|
| `tmp/PROMPT-EXACT.txt` | `archive/tmp-root/PROMPT-EXACT.txt` | `bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6` | `bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6` | copy |
| `tmp/append-rulings.js` | `archive/tmp-root/append-rulings.js` | `413b22d6cf61b436b6490b1826b82050a0947f024bdbee328e02e0e5a9175200` | `413b22d6cf61b436b6490b1826b82050a0947f024bdbee328e02e0e5a9175200` | copy |
| `tmp/batch-a-gen.js` | `archive/tmp-root/batch-a-gen.js` | `b18f3da55757d22296d70cb4d281ef8fa81b3d55216425b31c7333669aef9b03` | `b18f3da55757d22296d70cb4d281ef8fa81b3d55216425b31c7333669aef9b03` | copy |
| `tmp/batch-d-gen.js` | `archive/tmp-root/batch-d-gen.js` | `28710cabbf0e4f2ae23b4cf0e9363b7b508d88c00af277029b723c850a67d49f` | `28710cabbf0e4f2ae23b4cf0e9363b7b508d88c00af277029b723c850a67d49f` | copy |
| `tmp/fill-dataset.js` | `archive/tmp-root/fill-dataset.js` | `c308ad2654e1abfaa6bca3821f65bed2e4350c86460dab5c4c29c5815292cd33` | `c308ad2654e1abfaa6bca3821f65bed2e4350c86460dab5c4c29c5815292cd33` | copy |
| `tmp/fix-bom-and-append.js` | `archive/tmp-root/fix-bom-and-append.js` | `59c9b862ead14223f599ad8dc85cc87b02b3dd01017f84fc5125a943af9c3372` | `59c9b862ead14223f599ad8dc85cc87b02b3dd01017f84fc5125a943af9c3372` | copy |
| `tmp/freeze.js` | `archive/tmp-root/freeze.js` | `38b078c7a9277ad4820e2c157dcc1e5401aac55c86018dec5c4f7b5a31a3d4f5` | `38b078c7a9277ad4820e2c157dcc1e5401aac55c86018dec5c4f7b5a31a3d4f5` | copy |
| `tmp/glm-eval-batch/a-input.txt` | `archive/glm-eval-batch/a-input.txt` | `ea435e99732ebd7436422e0f84baf2e3bafe31e91964af2440010150af9ce564` | `ea435e99732ebd7436422e0f84baf2e3bafe31e91964af2440010150af9ce564` | copy |
| `tmp/glm-eval-batch/a-key.json` | `archive/glm-eval-batch/a-key.json` | `3639a989ab0f21a3bd2696cda1b616cad92ece2911bbf869ff4ac66ec979b5a3` | `3639a989ab0f21a3bd2696cda1b616cad92ece2911bbf869ff4ac66ec979b5a3` | copy |
| `tmp/glm-eval-out/leg-control.txt` | `archive/glm-eval-out/leg-control.txt` | `2fa1d9eeedae134c758174fff900311983a52b0eda8470889f4ff5d67f463a7f` | `2fa1d9eeedae134c758174fff900311983a52b0eda8470889f4ff5d67f463a7f` | copy |
| `tmp/glm-eval-out/leg-glm.txt` | `archive/glm-eval-out/leg-glm.txt` | `54c895f9df821231d9049ac7819077b44335db83383b266d86f3c28fee137b5b` | `54c895f9df821231d9049ac7819077b44335db83383b266d86f3c28fee137b5b` | copy |
| `tmp/glm-eval-out/selftest-report-20261005.txt` | `archive/glm-eval-out/selftest-report-20261005.txt` | `7b06d9d757b2d0e7fd00216245dbcef9405f9273afe3230e33d81fc48be9950d` | `7b06d9d757b2d0e7fd00216245dbcef9405f9273afe3230e33d81fc48be9950d` | copy |
| `tmp/glm-eval/BRIEFING.md` | `archive/glm-eval/BRIEFING.md` | `18867c95816abfa977c30eec31a74ef54ef17405b6c90e91caa6f05415df3ae8` | `18867c95816abfa977c30eec31a74ef54ef17405b6c90e91caa6f05415df3ae8` | copy |
| `tmp/glm-eval/FROZEN-MANIFEST.txt` | `archive/glm-eval/FROZEN-MANIFEST.txt` | `bda60cc1707b6c945ff2a7ebc48bb7e987af5c36125d276cc9758ffde613a347` | `bda60cc1707b6c945ff2a7ebc48bb7e987af5c36125d276cc9758ffde613a347` | copy |
| `tmp/glm-eval/GLM-USAGE-READING.txt` | `archive/glm-eval/GLM-USAGE-READING.txt` | `023426c8020dfc565c09d5176f6120458c1529d0862892f2ae4c2e2d5bff356f` | `023426c8020dfc565c09d5176f6120458c1529d0862892f2ae4c2e2d5bff356f` | copy |
| `tmp/glm-eval/PLAN-ERRATA.txt` | `archive/glm-eval/PLAN-ERRATA.txt` | `a3c083559bfe4ca4c866d2260acab3e483b0623241a0900721397e2565b59dbf` | `a3c083559bfe4ca4c866d2260acab3e483b0623241a0900721397e2565b59dbf` | copy |
| `tmp/glm-eval/PLAN-PROVENANCE.txt` | `archive/glm-eval/PLAN-PROVENANCE.txt` | `66f1a3f5e81fe4619c16cf340f4add5e41a3eff2a5af0d7a09fa593874e18dfa` | `66f1a3f5e81fe4619c16cf340f4add5e41a3eff2a5af0d7a09fa593874e18dfa` | copy |
| `tmp/glm-eval/PROBE-NOTES.txt` | `archive/glm-eval/PROBE-NOTES.txt` | `7f3dc09fae9292ffbbae5989b375f69fa2285a12c3c83fb0401863553ed817e6` | `7f3dc09fae9292ffbbae5989b375f69fa2285a12c3c83fb0401863553ed817e6` | copy |
| `tmp/glm-eval/TEST-REPORT-DIRECTIVE-GLM-EVAL.txt` | `archive/glm-eval/TEST-REPORT-DIRECTIVE-GLM-EVAL.txt` | `de5dac345258d0760780db839e74f06bee5e65aeb1372bc1166e29a4dd911472` | `de5dac345258d0760780db839e74f06bee5e65aeb1372bc1166e29a4dd911472` | copy |
| `tmp/glm-eval/bundle/00-plan.md` | `archive/glm-eval/bundle/00-plan.md` | `70076661d86615d8dfe180aa5884ec9618338ec9496b3893900a749c092f36b8` | `70076661d86615d8dfe180aa5884ec9618338ec9496b3893900a749c092f36b8` | copy |
| `tmp/glm-eval/bundle/01-clause-cn.txt` | `archive/glm-eval/bundle/01-clause-cn.txt` | `d8fe4bc82df774535c47c39406c96fee1325c87f3153b5d72ee2a378ff03e525` | `d8fe4bc82df774535c47c39406c96fee1325c87f3153b5d72ee2a378ff03e525` | copy |
| `tmp/glm-eval/bundle/02-diff.js` | `archive/glm-eval/bundle/02-diff.js` | `c939bc9e756696c463571c981a7d3756c87563de913b9c4c58faebdc4d5ff405` | `c939bc9e756696c463571c981a7d3756c87563de913b9c4c58faebdc4d5ff405` | copy |
| `tmp/glm-eval/bundle/03-prompt.txt` | `archive/glm-eval/bundle/03-prompt.txt` | `3b7a3580edcecb6a92ccb42f8acf37b469283136cb0c8d8d1b2b24d0b54fd758` | `3b7a3580edcecb6a92ccb42f8acf37b469283136cb0c8d8d1b2b24d0b54fd758` | copy |
| `tmp/glm-eval/contam.js` | `archive/glm-eval/contam.js` | `dfcb76e737b6ba5dab0453977fa6c1307bb2391e7890e5d1bd9aca3525099455` | `dfcb76e737b6ba5dab0453977fa6c1307bb2391e7890e5d1bd9aca3525099455` | copy |
| `tmp/glm-eval/hash.js` | `archive/glm-eval/hash.js` | `ff608449685e5150e5338a8bcad5e36e0e8910f172ed77182b8e1bf3c4c073b3` | `ff608449685e5150e5338a8bcad5e36e0e8910f172ed77182b8e1bf3c4c073b3` | copy |
| `tmp/glm-eval/key.json` | `archive/glm-eval/key.json` | `6fccbe461286abae63e74e52f64d8130f4fa979bb6d27ba71651fad6c3339536` | `6fccbe461286abae63e74e52f64d8130f4fa979bb6d27ba71651fad6c3339536` | copy |
| `tmp/glm-eval/metering/README.md` | `archive/glm-eval/metering/README.md` | `150c1a9bc4e001ab5cadf6e294162e8789048f4cb1ef83af1241baeaca4543a1` | `150c1a9bc4e001ab5cadf6e294162e8789048f4cb1ef83af1241baeaca4543a1` | copy |
| `tmp/glm-eval/plan-v1.md` | `archive/glm-eval/plan-v1.md` | `8be638d27cbf467f0b337de7e4362bad6b090b43e876377d1ad7ec7429181d04` | `8be638d27cbf467f0b337de7e4362bad6b090b43e876377d1ad7ec7429181d04` | copy |
| `tmp/glm-eval/score.js` | `archive/glm-eval/score.js` | `2c4052fa8e10273b09465b668a09b0c1e4402065576ec48ca4357341d27ee051` | `2c4052fa8e10273b09465b668a09b0c1e4402065576ec48ca4357341d27ee051` | copy |
| `tmp/glm-eval/selftest/fixture-contam-fail.txt` | `archive/glm-eval/selftest/fixture-contam-fail.txt` | `5d2e0631d60db305d9aae6f698797561e18a51e722a8b9fc7adc2d24f115ec04` | `5d2e0631d60db305d9aae6f698797561e18a51e722a8b9fc7adc2d24f115ec04` | copy |
| `tmp/glm-eval/selftest/fixture-score.txt` | `archive/glm-eval/selftest/fixture-score.txt` | `ec698d06c1bcc6a8a8af274c6050121d3cf81c43ce777430bf23c1e828fb9c11` | `ec698d06c1bcc6a8a8af274c6050121d3cf81c43ce777430bf23c1e828fb9c11` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-cjk1.txt` | `archive/glm-eval/selftest/fixtures/boundary-cjk1.txt` | `cd8405ae534bc76db879b992bf33c20974d0d57c707b3e47588b2e94e5f22e06` | `cd8405ae534bc76db879b992bf33c20974d0d57c707b3e47588b2e94e5f22e06` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-cjk2.txt` | `archive/glm-eval/selftest/fixtures/boundary-cjk2.txt` | `ec3962b4d8469a54a19ae86342ac2dd5a68bcee0da01d75578da63406262b1a4` | `ec3962b4d8469a54a19ae86342ac2dd5a68bcee0da01d75578da63406262b1a4` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-directive.txt` | `archive/glm-eval/selftest/fixtures/boundary-directive.txt` | `4590eda262ce6f0d6e2da605a6df6e90cf1bd4d549e197763f8013cac0d4b5b8` | `4590eda262ce6f0d6e2da605a6df6e90cf1bd4d549e197763f8013cac0d4b5b8` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-g8a.txt` | `archive/glm-eval/selftest/fixtures/boundary-g8a.txt` | `95abf3393bb78966d5c7c8dec217b7dad0c2c9ad86948729b0bbc90d722df0bb` | `95abf3393bb78966d5c7c8dec217b7dad0c2c9ad86948729b0bbc90d722df0bb` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-hash025.txt` | `archive/glm-eval/selftest/fixtures/boundary-hash025.txt` | `41225e5f54a0c0c8b4b5755a018c6e9763f07981ccc69711779d700c51d1d48f` | `41225e5f54a0c0c8b4b5755a018c6e9763f07981ccc69711779d700c51d1d48f` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-hash0c1.txt` | `archive/glm-eval/selftest/fixtures/boundary-hash0c1.txt` | `bb04ad75582f95d630e5444f0cd251601262d1f0354f24a62a64c6aff9ed1982` | `bb04ad75582f95d630e5444f0cd251601262d1f0354f24a62a64c6aff9ed1982` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-ob79.txt` | `archive/glm-eval/selftest/fixtures/boundary-ob79.txt` | `0e3bd8ab01eb11f47d188610cb6b482500d52bc238580eadaeb11512f57bd608` | `0e3bd8ab01eb11f47d188610cb6b482500d52bc238580eadaeb11512f57bd608` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-ob790.txt` | `archive/glm-eval/selftest/fixtures/boundary-ob790.txt` | `1cb7f0e8e075332adbad194915545d2c06cbd74a9eeaeccf1fe2b994374e4bfd` | `1cb7f0e8e075332adbad194915545d2c06cbd74a9eeaeccf1fe2b994374e4bfd` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-s10.txt` | `archive/glm-eval/selftest/fixtures/boundary-s10.txt` | `b49bf889a887b750956ee13959169f5408114a4b583150637a0f9b80ae8fdac1` | `b49bf889a887b750956ee13959169f5408114a4b583150637a0f9b80ae8fdac1` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-s123.txt` | `archive/glm-eval/selftest/fixtures/boundary-s123.txt` | `83edc34853693971f9630e89fe974c535b58e50648824355dbd27f61423ed23b` | `83edc34853693971f9630e89fe974c535b58e50648824355dbd27f61423ed23b` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-s8.txt` | `archive/glm-eval/selftest/fixtures/boundary-s8.txt` | `d0e43de3e6f42731db40751a9466ea394e4215cb08115f87883350fcf01e3394` | `d0e43de3e6f42731db40751a9466ea394e4215cb08115f87883350fcf01e3394` | copy |
| `tmp/glm-eval/selftest/fixtures/boundary-s9.txt` | `archive/glm-eval/selftest/fixtures/boundary-s9.txt` | `420c1322871c27fe0039df471e8d364073a74310a6dc71cea22718a0f4d3752e` | `420c1322871c27fe0039df471e8d364073a74310a6dc71cea22718a0f4d3752e` | copy |
| `tmp/glm-eval/selftest/fixtures/hash-bom.txt` | `archive/glm-eval/selftest/fixtures/hash-bom.txt` | `42c1e65b2c948bb754efb6ac171319d6e97ecb3d9afd4f20bd91b3ded25183c0` | `42c1e65b2c948bb754efb6ac171319d6e97ecb3d9afd4f20bd91b3ded25183c0` | copy |
| `tmp/glm-eval/selftest/fixtures/hash-nobom.txt` | `archive/glm-eval/selftest/fixtures/hash-nobom.txt` | `5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03` | `5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-abstain-line.txt` | `archive/glm-eval/selftest/fixtures/leg-abstain-line.txt` | `3a59038bad20c5129765f4d5a53f325efb0be564a28a8aa1ab2d463fead2a593` | `3a59038bad20c5129765f4d5a53f325efb0be564a28a8aa1ab2d463fead2a593` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-abstain.txt` | `archive/glm-eval/selftest/fixtures/leg-abstain.txt` | `03f1a97f4d3a0eb521b24b18cc77b392a4a00e35bf241829757c4cdbff6f26d9` | `03f1a97f4d3a0eb521b24b18cc77b392a4a00e35bf241829757c4cdbff6f26d9` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-acceptlines.txt` | `archive/glm-eval/selftest/fixtures/leg-acceptlines.txt` | `6e408e1468ed8d9e20981031410908b31b5dc542c7260c663975dabb355b99aa` | `6e408e1468ed8d9e20981031410908b31b5dc542c7260c663975dabb355b99aa` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-clean.txt` | `archive/glm-eval/selftest/fixtures/leg-clean.txt` | `8b60d780b1637921cd1a1d1d9f9c9206e7fc478637da7df1808c42013a76a1da` | `8b60d780b1637921cd1a1d1d9f9c9206e7fc478637da7df1808c42013a76a1da` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-comma.txt` | `archive/glm-eval/selftest/fixtures/leg-comma.txt` | `de25a865a50ebe6840019f8346219c6635377ae325334c6b60a837f1eaa5f2e4` | `de25a865a50ebe6840019f8346219c6635377ae325334c6b60a837f1eaa5f2e4` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-contam.txt` | `archive/glm-eval/selftest/fixtures/leg-contam.txt` | `3354cb1d27957a0728688be21ef04d564732e8ffa0804790769f382588f8cb44` | `3354cb1d27957a0728688be21ef04d564732e8ffa0804790769f382588f8cb44` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-dedup.txt` | `archive/glm-eval/selftest/fixtures/leg-dedup.txt` | `a46b8c016876135a2de5aefcfcf1e302a08c9dffe0fd11eddd66c89fb86e8d76` | `a46b8c016876135a2de5aefcfcf1e302a08c9dffe0fd11eddd66c89fb86e8d76` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-empty.txt` | `archive/glm-eval/selftest/fixtures/leg-empty.txt` | `2b544e496f9c88a8ae0d337f57024c38b97f23b8e3602584619258e12094bae2` | `2b544e496f9c88a8ae0d337f57024c38b97f23b8e3602584619258e12094bae2` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-extra.txt` | `archive/glm-eval/selftest/fixtures/leg-extra.txt` | `4b55d162afd1515fdd7152c5ad792e5267092fb3cc97afbe9012650c54e571cd` | `4b55d162afd1515fdd7152c5ad792e5267092fb3cc97afbe9012650c54e571cd` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-fps.txt` | `archive/glm-eval/selftest/fixtures/leg-fps.txt` | `7bf1ff635e21b61307bdcda4905c71f127efef379423197c1b9f8a103651840a` | `7bf1ff635e21b61307bdcda4905c71f127efef379423197c1b9f8a103651840a` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-fullwidth.txt` | `archive/glm-eval/selftest/fixtures/leg-fullwidth.txt` | `5f71171e8255f945698124eb63524603f44d6d3da7c3f2b1dac6683b10764387` | `5f71171e8255f945698124eb63524603f44d6d3da7c3f2b1dac6683b10764387` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-noconcl.txt` | `archive/glm-eval/selftest/fixtures/leg-noconcl.txt` | `5e38192d6c73ed45728378fd96869c05779a1f485224af6403edb88c6a2ba70c` | `5e38192d6c73ed45728378fd96869c05779a1f485224af6403edb88c6a2ba70c` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-outside.txt` | `archive/glm-eval/selftest/fixtures/leg-outside.txt` | `7cedeee512ac9b9a21b819f40b23bacc31090f6c5a144ebdf3742a4fcf1dbad2` | `7cedeee512ac9b9a21b819f40b23bacc31090f6c5a144ebdf3742a4fcf1dbad2` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-perfect.txt` | `archive/glm-eval/selftest/fixtures/leg-perfect.txt` | `36d40ca7061d8b94448670f6947ceeb4c7d9047a55c717d1874e07b5a3858672` | `36d40ca7061d8b94448670f6947ceeb4c7d9047a55c717d1874e07b5a3858672` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-qcheck.txt` | `archive/glm-eval/selftest/fixtures/leg-qcheck.txt` | `0a58d82387e3e944df592663d3e9bdd23554aef9e1ba0c937160a36e10a34e2f` | `0a58d82387e3e944df592663d3e9bdd23554aef9e1ba0c937160a36e10a34e2f` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-range-miss.txt` | `archive/glm-eval/selftest/fixtures/leg-range-miss.txt` | `dfd77a61368639d9e6d94b436df0266ead8486679be730fc22e2e1926cf252f0` | `dfd77a61368639d9e6d94b436df0266ead8486679be730fc22e2e1926cf252f0` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-range.txt` | `archive/glm-eval/selftest/fixtures/leg-range.txt` | `5f041561b478710b0a00b778fb949c1740b327323defc1487b2a77f86ac18f9c` | `5f041561b478710b0a00b778fb949c1740b327323defc1487b2a77f86ac18f9c` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-sev-cn.txt` | `archive/glm-eval/selftest/fixtures/leg-sev-cn.txt` | `ea59fe812117f0418abee250687f6fc01822a49dbbbfbe4ae0db4eeef0423b9f` | `ea59fe812117f0418abee250687f6fc01822a49dbbbfbe4ae0db4eeef0423b9f` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-sev-exact.txt` | `archive/glm-eval/selftest/fixtures/leg-sev-exact.txt` | `275b6eb589a231ef2815b47b935d53f9b28b1db2102d28b547a851b927b08fdf` | `275b6eb589a231ef2815b47b935d53f9b28b1db2102d28b547a851b927b08fdf` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-sev-off1.txt` | `archive/glm-eval/selftest/fixtures/leg-sev-off1.txt` | `cd23868bb537cf0343d9d4d7d969ec0997fc08a56c18857d0164150a6d466c01` | `cd23868bb537cf0343d9d4d7d969ec0997fc08a56c18857d0164150a6d466c01` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-sev-off2.txt` | `archive/glm-eval/selftest/fixtures/leg-sev-off2.txt` | `c5d7324412988b23d991a7bcf74c34335121b78ce5a61bde5cbadb2a35fac140` | `c5d7324412988b23d991a7bcf74c34335121b78ce5a61bde5cbadb2a35fac140` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-sev-table.txt` | `archive/glm-eval/selftest/fixtures/leg-sev-table.txt` | `e3abb722801742b015337b7f1ca3aa021413a6202d916029cbfecc7859775a69` | `e3abb722801742b015337b7f1ca3aa021413a6202d916029cbfecc7859775a69` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-slash.txt` | `archive/glm-eval/selftest/fixtures/leg-slash.txt` | `fb46611ed342e0e9fbbd9d0f9b7a4bf7fcfc9551075f3b63549d453c9b40e253` | `fb46611ed342e0e9fbbd9d0f9b7a4bf7fcfc9551075f3b63549d453c9b40e253` | copy |
| `tmp/glm-eval/selftest/fixtures/leg-trailing-ws.txt` | `archive/glm-eval/selftest/fixtures/leg-trailing-ws.txt` | `0c18703e7c83bafc2d6e46ce01eaac6dca8fabbd8096ae29dacf5856e8e8d133` | `0c18703e7c83bafc2d6e46ce01eaac6dca8fabbd8096ae29dacf5856e8e8d133` | copy |
| `tmp/glm-eval/selftest/run-selftest.js` | `archive/glm-eval/selftest/run-selftest.js` | `2fc646ad9464133b0f6b9211e02d9815414e9b640155e0d9af352651661e0e49` | `2fc646ad9464133b0f6b9211e02d9815414e9b640155e0d9af352651661e0e49` | copy |
| `tmp/glm-eval/selftest/selftest-report.txt` | `archive/glm-eval/selftest/selftest-report.txt` | `7b06d9d757b2d0e7fd00216245dbcef9405f9273afe3230e33d81fc48be9950d` | `7b06d9d757b2d0e7fd00216245dbcef9405f9273afe3230e33d81fc48be9950d` | copy |
| `tmp/norm-score.js` | `archive/tmp-root/norm-score.js` | `f07365c0dbddcff94a8f6371003728aea8f6d9efb48cc4c6a556f72d4877322e` | `f07365c0dbddcff94a8f6371003728aea8f6d9efb48cc4c6a556f72d4877322e` | copy |
| `tmp/persist.js` | `archive/tmp-root/persist.js` | `ddd0956fe8ba5e0119054e0e7ab84a0be8f4c531f6f06aa20c8e77041601030a` | `ddd0956fe8ba5e0119054e0e7ab84a0be8f4c531f6f06aa20c8e77041601030a` | copy |
| `tmp/recon-briefing.js` | `archive/tmp-root/recon-briefing.js` | `60c394ed3aa87eff6531cd386eb1bc898b9d82aa8b184aacde510072386764bf` | `60c394ed3aa87eff6531cd386eb1bc898b9d82aa8b184aacde510072386764bf` | copy |
| `tmp/save-leg.js` | `archive/tmp-root/save-leg.js` | `96bb7500b08e4bf51e8aa1caa886789b394f2b28cd0ffe005a6890df2c81465a` | `96bb7500b08e4bf51e8aa1caa886789b394f2b28cd0ffe005a6890df2c81465a` | copy |
| `tmp/_pack-glm-eval.js` | `archive/tooling/pack-glm-eval.js` | `67638513ca06cc83e38b2738334464aa8272c025a090be37d2f1adf8ec366165` | `67638513ca06cc83e38b2738334464aa8272c025a090be37d2f1adf8ec366165` | copy |
| `tmp/_manifest-template.md` | `archive/tooling/manifest-template.md` | `74ed1697069f448d405735aaf429afc454892d794cdc848ba265278e8667786f` | `74ed1697069f448d405735aaf429afc454892d794cdc848ba265278e8667786f` | copy |

## 未入库文件清单

无。本包不引用任何未入库文件；全部载荷均已随包入库。

## 排除清单（有意不归档）

- 无。本包覆盖 `tmp/glm-eval/**`、`tmp/glm-eval-out/**`、`tmp/glm-eval-batch/**` 与 `tmp/` 根承重脚本的全部条目（来源清单 = `tmp/glm-eval/FROZEN-MANIFEST.txt` 的 72 条 + 该清单自身）。
- 一次性中间物（`tmp/_g4bchk.js`、`tmp/_g4b.txt`、`tmp/_persist2.log`、`tmp/_append-*.txt`、`tmp/_manifest-*.md`）属 §1.5.1 第 ③ 类，**用完即清**，不入包。

## 重建说明

- **不适用**：本包全部为**原始字节归档**，无「重建归档」条目（§1.5.1 第 ④ 类）。
- 唯一被覆写过的载荷是 `archive/glm-eval/selftest/selftest-report.txt`：其**脚手架期原始字节不可复现**（已作为不可逆取证损失登记于报告与 `TEST-REPORT` §12）；包内存放的是整改后重跑版本，并另存时间点快照 `archive/glm-eval-out/selftest-report-20261005.txt`（与前者同哈希）。

## 对应引用点

- 报告：`docs/validation/directive-glm-eval{,_EN}.md`（§3.2 bundle 表 / §3.6 冻结自证 / §6 可证伪性 / §9 门禁结果）。
- 评估数据集：`docs/validation/directive-glm-eval-dataset.json`。
- 自由度清单：`docs/validation/evidence/DIRECTIVE-GLM-EVAL-freedom-list.md`（证据根级 `*-freedom-list.md` 形态）。
- 台账：`docs/t027/REMAINING_SLICES{,_EN}.md` 的 `DIRECTIVE-GLM-EVAL` 完成行。
