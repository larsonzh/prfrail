# MANIFEST — P64-AT24-SUPERVISED-BLACK-BOX 证据包

建包：步骤 0（2026-11-11）；随片进展追加条目（① design、③ 验收腿证据、④ 门禁读数随环注册）。编码归一注记：本包载荷均经机械归一后入库（`.md` = BOM+LF；`.txt`/`SHA256SUMS.txt` = 无 BOM+LF），原始字节 = 归档字节，两列 sha256 同值为常态（非缺失）。

## 归档条目清单

| 路径 | 原始字节 sha256 | 归档字节 sha256 | 说明 |
|---|---|---|---|
| `00-briefing.md` | `8935dfb3477c4cba20942a5f23d465882f7e3aed5c7641871ca22560c78190bf` | `8935dfb3477c4cba20942a5f23d465882f7e3aed5c7641871ca22560c78190bf` | briefing verbatim |
| `freedom-list.md` | `3b990efe71d81ea6f92ce15ab4a08765058a3a0124d4940c4597e3d7c7b287c6` | `3b990efe71d81ea6f92ce15ab4a08765058a3a0124d4940c4597e3d7c7b287c6` | freedom list |
| `input-list.md` | `d387a61abf6a9f6e8fa7e8fb92fdc94f790c7f4c0b7f131df60c87a6e374a9d9` | `d387a61abf6a9f6e8fa7e8fb92fdc94f790c7f4c0b7f131df60c87a6e374a9d9` | input list |
| `01-arch/design-r1.md` | `57afcf526d0f965e70905f434a3524f7065eed4e5cdd6bf90e99ac2de6501866` | `57afcf526d0f965e70905f434a3524f7065eed4e5cdd6bf90e99ac2de6501866` | stage 1 design |
| `01-arch/r09-3-user-confirmation.md` | `a2bc3a5e6de402bf2d299848490a664375cb25a53efc1802cb17094bfd2cd3fc` | `a2bc3a5e6de402bf2d299848490a664375cb25a53efc1802cb17094bfd2cd3fc` | R-09.3 user confirmation |
| `01-arch/r09-7-user-confirmation.md` | `900b470f7012e88e35a49257e7eafd2048dab3dc9978c863f6ab1a455a9762f8` | `900b470f7012e88e35a49257e7eafd2048dab3dc9978c863f6ab1a455a9762f8` | R-09.7 user confirmation |
| `02-impl/schema-gap-analysis.md` | `677e74c056ee626011f29a70dc7b6cf5872335311af4e468658e801a36ea1144` | `677e74c056ee626011f29a70dc7b6cf5872335311af4e468658e801a36ea1144` | schema gap analysis (ruling step 1) |
| `03-test/SUMMARY.md` | `b4e753b05c55d7f277f428e62e6880459f6a28ca331d7b14d622bb9a9a3fdb74` | `b4e753b05c55d7f277f428e62e6880459f6a28ca331d7b14d622bb9a9a3fdb74` | 17-leg summary table |
| `03-test/L01-risk-confirmation-binding.txt` | `09a6f6a900dbdc8252c1270ae1bd832d772503b772cc014cab4e7e049473132f` | `09a6f6a900dbdc8252c1270ae1bd832d772503b772cc014cab4e7e049473132f` | leg 1 risk-confirmation binding |
| `03-test/L02-delivery-not-completion.txt` | `a3a94e241df1d7053aeb762a30c72bb477596edb0a2bd0a04a57014e5b165ac8` | `a3a94e241df1d7053aeb762a30c72bb477596edb0a2bd0a04a57014e5b165ac8` | leg 2 delivery != completion |
| `03-test/L03-explicit-return-postscan.txt` | `1257e8f794f4692221894c89e0a7d08b0bdd3c0fe62fb5f7883ce13dfdb301ad` | `1257e8f794f4692221894c89e0a7d08b0bdd3c0fe62fb5f7883ce13dfdb301ad` | leg 3 explicit return + postscan |
| `03-test/L04-independent-gates-review.txt` | `2ca06489316b3c2d67abd3af046a2a09a76f5d3f296646fb25ebb23b8091850a` | `2ca06489316b3c2d67abd3af046a2a09a76f5d3f296646fb25ebb23b8091850a` | leg 4 independent gates/review |
| `03-test/L05-unknown-dimensions-cost.txt` | `f022176439f7cbc2153f4e9024ba4f788dacdea9726000883bb3f0d39967234c` | `f022176439f7cbc2153f4e9024ba4f788dacdea9726000883bb3f0d39967234c` | leg 5 unknown dimensions incl cost |
| `03-test/L06-three-sentences.txt` | `35333725567a197750b18c3f7f91f2b59c8b8f277e572eaad173f6cc5c7acd70` | `35333725567a197750b18c3f7f91f2b59c8b8f277e572eaad173f6cc5c7acd70` | leg 6 three sentences |
| `03-test/L07-no-auto-degrade.txt` | `6106002a01fb0292de3b321c2e4da1efd7aa0cf12ecac9240adf8872cf0d6309` | `6106002a01fb0292de3b321c2e4da1efd7aa0cf12ecac9240adf8872cf0d6309` | leg 7 no auto degrade |
| `03-test/L08-unknown-not-safe.txt` | `754673f1d89c99f82d7924a3b39c3a22bbfb3363bb4f7af9b5238773684b84cf` | `754673f1d89c99f82d7924a3b39c3a22bbfb3363bb4f7af9b5238773684b84cf` | leg 8 unknown not shown as safe |
| `03-test/L09-01-tool-use.txt` | `f6f24a79cee242871436676abdd0edcf2b2f04a6abf02505070891d615bd530e` | `f6f24a79cee242871436676abdd0edcf2b2f04a6abf02505070891d615bd530e` | leg 9.1 tool use |
| `03-test/L09-02-resource-reads.txt` | `243d519333c18cb439148de36fe62574ba1f3dd76b6eefa3c3564df09b41a59d` | `243d519333c18cb439148de36fe62574ba1f3dd76b6eefa3c3564df09b41a59d` | leg 9.2 resource reads |
| `03-test/L09-03-network.txt` | `4e73bec827a60a3a7f4ed497ee0a664be9ba0701a71c1ebc6dd1199dc8204f44` | `4e73bec827a60a3a7f4ed497ee0a664be9ba0701a71c1ebc6dd1199dc8204f44` | leg 9.3 network |
| `03-test/L09-04-cost.txt` | `ba9a47ae4e2edb4aee0b78c966a25ed51e4ae5f8925e16654e918053d91519b7` | `ba9a47ae4e2edb4aee0b78c966a25ed51e4ae5f8925e16654e918053d91519b7` | leg 9.4 cost |
| `03-test/L09-05-background-processes.txt` | `345e59901fd4c54fb985f0261c8572ff84454f334d1e38debceea3502e9ccec8` | `345e59901fd4c54fb985f0261c8572ff84454f334d1e38debceea3502e9ccec8` | leg 9.5 background processes |
| `03-test/L09-06-session-continuity.txt` | `8707115b552df2273f99cbf65017c3a6b5f9a99492e2742c7bd810947456f9c3` | `8707115b552df2273f99cbf65017c3a6b5f9a99492e2742c7bd810947456f9c3` | leg 9.6 session continuity |
| `03-test/L09-07-external-side-effects.txt` | `74ed865d726d2744bab85c484850222f3cde2e8e1c8c64b3e202462558ecf527` | `74ed865d726d2744bab85c484850222f3cde2e8e1c8c64b3e202462558ecf527` | leg 9.7 external side effects |
| `03-test/L09-08-mandatory-fields.txt` | `9cd4df6d76fe74b409b87e8e943e318221dde1fd6325f4e271ce93cd1caa0e14` | `9cd4df6d76fe74b409b87e8e943e318221dde1fd6325f4e271ce93cd1caa0e14` | leg 9.8 mandatory fields |
| `03-test/L09-09-judgment-approach.txt` | `50cbfa7aade9d2b27438fce80b946e9c489514f1e254f7ddd741e865c7fde7e7` | `50cbfa7aade9d2b27438fce80b946e9c489514f1e254f7ddd741e865c7fde7e7` | leg 9.9 judgment approach |
| `03-test/L05-r2-schema-seven-mandatory.txt` | `1a1fe67e39494bf9783daec62bc3e8fa80dd12659fff71380ca8ef2de74df608` | `1a1fe67e39494bf9783daec62bc3e8fa80dd12659fff71380ca8ef2de74df608` | leg 5 rerun R2 (schema seven-mandatory) |
| `03-test/L09-08-r2-schema-seven-mandatory.txt` | `cd337e68a48a37ab04d32d3780ab89794f6e14ed3a6767193c6c21ed628f8a56` | `cd337e68a48a37ab04d32d3780ab89794f6e14ed3a6767193c6c21ed628f8a56` | leg 9.8 rerun R2 (schema seven-mandatory) |
| `03-test/L09-01-junction-probe.txt` | `10f4588aeb4f3b8d22c778c709f273f1b78ce23858142719a65e320b9becfe87` | `10f4588aeb4f3b8d22c778c709f273f1b78ce23858142719a65e320b9becfe87` | leg 9.1 junction equivalence probe |
| `03-test/L09-01-r2.txt` | `8e3ca65f9a7d6fd3cafb072458d3b2be2725cadec1dbb388a788812de4427343` | `8e3ca65f9a7d6fd3cafb072458d3b2be2725cadec1dbb388a788812de4427343` | leg 9.1 rerun R2 (junction equivalent cover) |
| `04-reviews/05-prereview-r1.txt` | `861d2479742da5b5546a8dce02acaedc37e18e4efd9d691043e983e6c377131b` | `861d2479742da5b5546a8dce02acaedc37e18e4efd9d691043e983e6c377131b` | stage 5 pre-review r1 (verbatim) |
| `04-reviews/06-scan-r1.txt` | `494ed306e1cbe6a4305262a56d4633cbe2aba0d34fae4f6df13c1f678b3b95fa` | `494ed306e1cbe6a4305262a56d4633cbe2aba0d34fae4f6df13c1f678b3b95fa` | stage 6 independent scan (verbatim) |
| `04-reviews/07-blind-r1.txt` | `7de8eff6edfee27cbc25bfa549b79b68f7dde829c36f6f083119dd01f8be4583` | `7de8eff6edfee27cbc25bfa549b79b68f7dde829c36f6f083119dd01f8be4583` | stage 7 blind review r1 (verbatim) |
| `04-reviews/07-final-r1.txt` | `280f6455a6f46a2e9fef3c7f0ec8773e73a613736316a247e3943534f49c0375` | `280f6455a6f46a2e9fef3c7f0ec8773e73a613736316a247e3943534f49c0375` | stage 7 final review r1 (verbatim) |
| `04-reviews/07-rereview-r1.txt` | `f0f0bc86a45787a002d2ee61fa169786ceadb7d3cc8a9d8b19b961aad36483c3` | `f0f0bc86a45787a002d2ee61fa169786ceadb7d3cc8a9d8b19b961aad36483c3` | stage 7 re-review r1 (verbatim) |
| `06-gates/gates-r1.txt` | `907b58f01570c7c4347c1cb42a5daa10b0f76d9c92784710ab19785db5a64af0` | `907b58f01570c7c4347c1cb42a5daa10b0f76d9c92784710ab19785db5a64af0` | stage 4 gates readings |
