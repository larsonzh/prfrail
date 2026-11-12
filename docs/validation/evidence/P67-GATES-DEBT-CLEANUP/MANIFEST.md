# MANIFEST — P67-GATES-DEBT-CLEANUP 证据包

建包：步骤 0（2026-11-12）；随片进展追加条目（① 五件、② 步骤 A–G 核验、③ 一致性检查、④ 门禁读数、⑤⑥⑦ 回执随环注册）。编码归一注记：本包 `.md` 载荷均经机械归一（BOM+LF）后入库；`.txt` 与 `SHA256SUMS.txt` 为无 BOM+LF。

## 归档条目清单

| 路径 | 原始字节 sha256 | 归档字节 sha256 | 说明 |
|---|---|---|---|
| `00-briefing.md` | `c00b57004b6d391dba554d5733b35ec9fd844c316c035db7cc00394751c8bdea` | `c00b57004b6d391dba554d5733b35ec9fd844c316c035db7cc00394751c8bdea` | briefing verbatim |
| `01-arch/design-r1.md` | `a339d2fe9c7889ea0c308510fe826181985b06d58e78108464dbdee839dbc6f6` | `a339d2fe9c7889ea0c308510fe826181985b06d58e78108464dbdee839dbc6f6` | stage 1 design (five artefacts; conclusion DESIGN: DONE) |
| `02-impl/impl-r1.md` | `aa8da5bc337941c72dcccec6448ea874efe45b15c9961e7fb71e34c55d65a6be` | `aa8da5bc337941c72dcccec6448ea874efe45b15c9961e7fb71e34c55d65a6be` | step A implementation decision list (controller-recorded) |
| `02-batch/glm-batch1.txt` | `ce6ddb78b0b28bd0b30a2ae0f7fe91702d36b7e841633b9e358f36c4aeb784d6` | `ce6ddb78b0b28bd0b30a2ae0f7fe91702d36b7e841633b9e358f36c4aeb784d6` | batch-assist ① g4b cache-key enumeration (1849 chars) |
| `02-batch/glm-batch2.txt` | `2214be59fd0865b616fe21ee79bc69e8aa45c776e09dfe73088d2754fc3fd612` | `2214be59fd0865b616fe21ee79bc69e8aa45c776e09dfe73088d2754fc3fd612` | batch-assist ② seven-pit comparison (1379 chars; stale pit-6 row annotated) |
| `02-batch/glm-batch3.txt` | `da0ea7d975050ecaa4608bb61a8823e6613cbbf36e67517ee7b978b61a611361` | `da0ea7d975050ecaa4608bb61a8823e6613cbbf36e67517ee7b978b61a611361` | batch-assist ③ ⑥ three-round enumeration (1404 chars) |
| `03-test/step-B-check.md` | `ef9ff7eaaa7b7967a2edd4f84d97da8cbc44dbb06e1cd005d0c47d9cd8fbc1e6` | `ef9ff7eaaa7b7967a2edd4f84d97da8cbc44dbb06e1cd005d0c47d9cd8fbc1e6` | stage artefact (step B verification, T125–T129, +14 assertions) |
| `03-test/step-C-check.md` | `773a1a5b11da238497b5ca0b2b920d5e14a4c0a39e7e9c0020d10065f8b1738b` | `773a1a5b11da238497b5ca0b2b920d5e14a4c0a39e7e9c0020d10065f8b1738b` | stage artefact (step C verification, pitfalls 27–30 appended, 1–26 untouched) |
| `03-test/consistency-r1.md` | `ee0f3d4561e6d93e74f5893d5af21a713a6eb3210fb0b7ff9eecfcf4480e7dd3` | `ee0f3d4561e6d93e74f5893d5af21a713a6eb3210fb0b7ff9eecfcf4480e7dd3` | stage 3 consistency check (9 items PASS; one wording fix re-verified) |
| `03-test/step-A-check.md` | `f030fb2e2ee83642f2ce0868283efc1e79e05471953270bd8ac9f9b137f82824` | `f030fb2e2ee83642f2ce0868283efc1e79e05471953270bd8ac9f9b137f82824` | stage artefact (step A verification, back-filled per blind High) |
| `03-test/step-D-check.md` | `e7c477938c6652a7ec2336ade92889421e9eb0e25c26d1b546235fcb9ee5affd` | `e7c477938c6652a7ec2336ade92889421e9eb0e25c26d1b546235fcb9ee5affd` | stage artefact (step D verification, back-filled per blind High) |
| `03-test/step-E-check.md` | `9b06510ee83fb648fc1102825ae36e3cc3b1d336959e758ab6e54acf9a0a4214` | `9b06510ee83fb648fc1102825ae36e3cc3b1d336959e758ab6e54acf9a0a4214` | stage artefact (step E verification, back-filled per blind High) |
| `03-test/step-F-check.md` | `6839358ead946afff6a6e4565de703594f0ec9d329f35806cee4819d1a11bc53` | `6839358ead946afff6a6e4565de703594f0ec9d329f35806cee4819d1a11bc53` | stage artefact (step F verification, back-filled per blind High) |
| `04-reviews/05-prereview-r1.md` | `bf8eeea6e9be72b8d6c8d6dd8fa59fbdadca33532e79d1162581028138d17505` | `bf8eeea6e9be72b8d6c8d6dd8fa59fbdadca33532e79d1162581028138d17505` | stage artefact (⑤ r1 verbatim archive) |
| `04-reviews/05-prereview-r2.md` | `a103d9e3e05d95d3391b66f1cdcf0b5c607ce7c1a7da06324463acc53a27e6b5` | `a103d9e3e05d95d3391b66f1cdcf0b5c607ce7c1a7da06324463acc53a27e6b5` | stage artefact (⑤ r2 verbatim archive) |
| `04-reviews/06-scan-r1.md` | `c6fcdfa37daa4a40b4cb937091740159272b601cc9496b5447100fd728ab8a32` | `c6fcdfa37daa4a40b4cb937091740159272b601cc9496b5447100fd728ab8a32` | stage artefact (⑥ r1 verbatim archive) |
| `04-reviews/07-blind-r1.md` | `d46d140e69c4f6b61470dad75fc7dd7883a8806eabc3df5f52f3d49f8ec40eac` | `d46d140e69c4f6b61470dad75fc7dd7883a8806eabc3df5f52f3d49f8ec40eac` | stage artefact (⑦ blind r1 verbatim archive) |
| `04-reviews/07-final-r1.md` | `d1f8189161743c0cffc8e3d7229e98f482cf060d671b8058ab95e8e12d1cd73a` | `d1f8189161743c0cffc8e3d7229e98f482cf060d671b8058ab95e8e12d1cd73a` | stage artefact (⑦ final r1 verbatim archive) |
| `04-reviews/07-rereview-r1.md` | `afb3b868e4478bff388bceac1e44b376e47aa63a6fdc879be2fcb6a046aed2ab` | `afb3b868e4478bff388bceac1e44b376e47aa63a6fdc879be2fcb6a046aed2ab` | stage artefact (⑦ re-review r1 verbatim archive) |
| `06-gates/step0-selftest-r1.txt` | `a4e25ae22301cec5dc76c1b24aa68c5a9e5e14b10a64bfd7b1c9de3a8bad36c9` | `a4e25ae22301cec5dc76c1b24aa68c5a9e5e14b10a64bfd7b1c9de3a8bad36c9` | SELFTEST 锚点正式原始行（与 tmp 捕获互证） |
| `06-gates/gates-r1.txt` | `37f600314ba86fd6e65d7e7e3a30111a9c8c3d4573f178e71d11db6675052e8c` | `37f600314ba86fd6e65d7e7e3a30111a9c8c3d4573f178e71d11db6675052e8c` | stage 4 gates raw readings (selftest + tree/index/ci, TOTAL_FAIL=0) |
| `06-gates/gates-r2.txt` | `8b81287915d190fedd7ce0985c761ed8b6e54e17ec297619383f2020ec1ab199` | `8b81287915d190fedd7ce0985c761ed8b6e54e17ec297619383f2020ec1ab199` | stage 9 final gates raw readings (SELFTEST 737/737 + three scopes TOTAL_FAIL=0) |
| `freedom-list.md` | `87e45847b5dbd580dcb7809157e1fc08b110b492b31f3559637d8ad42aa1375a` | `87e45847b5dbd580dcb7809157e1fc08b110b492b31f3559637d8ad42aa1375a` | freedom list (F-1..F-14) |
| `input-list.md` | `7c782c17846749c9bc23647afbedbcdeeac131f785390b84a011a4babe0c75a6` | `7c782c17846749c9bc23647afbedbcdeeac131f785390b84a011a4babe0c75a6` | input list (16 items) |
