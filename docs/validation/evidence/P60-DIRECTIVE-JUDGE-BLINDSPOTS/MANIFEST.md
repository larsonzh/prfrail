# MANIFEST — P60-DIRECTIVE-JUDGE-BLINDSPOTS 证据包

建包：步骤 0（2026-11-06）；随片进展追加条目（design-r1 阶段 A 注册；ruling / impl-r1 阶段 B 注册；test-r1/test-r2 ③ 后注册；05/06/07 各环回执与门禁读数随环注册）。编码归一注记：本包载荷均经机械归一（BOM+LF）后入库，原始字节 = 归档字节，两列 sha256 同值为常态（非缺失）。

## 后缀与编码映射

| 原始路径 | 归档路径 | 原始字节 sha256 | 归档字节 sha256 | 说明 |
|---|---|---|---|---|
| `00-briefing.md` | `00-briefing.md` | `2ba0b8ffa3b18b67ba924fd0dfbfae13ab5340b5758ea3aa8439b66863e905d4` | `2ba0b8ffa3b18b67ba924fd0dfbfae13ab5340b5758ea3aa8439b66863e905d4` | briefing verbatim |
| `freedom-list.md` | `freedom-list.md` | `38274467d21d4ee57a98f66e20f2d1da350a04298b2e8bd242b0672d0745fc56` | `38274467d21d4ee57a98f66e20f2d1da350a04298b2e8bd242b0672d0745fc56` | freedom list |
| `input-list.md` | `input-list.md` | `5f1e42a3979be0959d77de79816d22b14e9f3dd9bd941e576f491d8e2fc64294` | `5f1e42a3979be0959d77de79816d22b14e9f3dd9bd941e576f491d8e2fc64294` | input list |
| `01-arch/design-r1.md` | `01-arch/design-r1.md` | `90376ddf8039675d7762fd3f3b67a126ed0bbc9e0793d10f9e44a4ed85767d60` | `90376ddf8039675d7762fd3f3b67a126ed0bbc9e0793d10f9e44a4ed85767d60` | arch three-way design (① output) |
| `01-arch/ruling-2026-11-06-p60.md` | `01-arch/ruling-2026-11-06-p60.md` | `3a1758666e031c5da23ff1ae9c17f4064399eedf82c9ab85e2b11ad095a00c01` | `3a1758666e031c5da23ff1ae9c17f4064399eedf82c9ab85e2b11ad095a00c01` | user ruling verbatim (phase A) |
| `02-impl/impl-r1.md` | `02-impl/impl-r1.md` | `11fdc59184460fd1baef8d2dd58900f81fed746cfd40b50fbd05f8fac6f7ed98` | `11fdc59184460fd1baef8d2dd58900f81fed746cfd40b50fbd05f8fac6f7ed98` | ② implementer receipt verbatim |
| `03-test/test-r1.md` | `03-test/test-r1.md` | `f258268abdcda60786eb958d1398fe8e23fd8061d0b21f8080a3abdf04e9c745` | `f258268abdcda60786eb958d1398fe8e23fd8061d0b21f8080a3abdf04e9c745` | ③ tester receipt verbatim |
| `03-test/test-r2.md` | `03-test/test-r2.md` | `2bd99c099195fcfacc1e5bcab9ce99f1677ebd5d9608fde2a1f252d0255719eb` | `2bd99c099195fcfacc1e5bcab9ce99f1677ebd5d9608fde2a1f252d0255719eb` | ③ tester receipt r2 verbatim |
| `04-reviews/05-pre-review-r1.md` | `04-reviews/05-pre-review-r1.md` | `216f7d02100a31d5844d96592b14ff574f5f5c6b3f7b1fde99b7975c8588daf0` | `216f7d02100a31d5844d96592b14ff574f5f5c6b3f7b1fde99b7975c8588daf0` | ⑤ pre-review receipt verbatim |
| `04-reviews/05-pre-review-r2.md` | `04-reviews/05-pre-review-r2.md` | `b7687453f402adb0f0ad221f4cc94f928a84d5d24da05e69a0cf849d72011f72` | `b7687453f402adb0f0ad221f4cc94f928a84d5d24da05e69a0cf849d72011f72` | ⑤ pre-review r2 receipt verbatim |
| `04-reviews/06-scan-r1.md` | `04-reviews/06-scan-r1.md` | `44c0e68a865141dd66c0ccafca6439ae78b4b9bfe3ad9e6d71adfcbeb262a865` | `44c0e68a865141dd66c0ccafca6439ae78b4b9bfe3ad9e6d71adfcbeb262a865` | ⑥ independent-scan receipt verbatim |
| `04-reviews/07-final-r1.md` | `04-reviews/07-final-r1.md` | `14b2a54fefa1498da898838f657b2f006f05b37733acd17332da39af54a22761` | `14b2a54fefa1498da898838f657b2f006f05b37733acd17332da39af54a22761` | ⑦ final-review receipt verbatim |
| `04-reviews/07-blind-r1.md` | `04-reviews/07-blind-r1.md` | `ee01b330220d96ab437983fb5323cdb24f8138857c7d1e92f00a6b105f7e0c03` | `ee01b330220d96ab437983fb5323cdb24f8138857c7d1e92f00a6b105f7e0c03` | ⑦ blind-review receipt verbatim |
| `04-reviews/07-blind-r2.md` | `04-reviews/07-blind-r2.md` | `36723fe43740b81bb265a33d67207f6cc46784843f906d4265c4013698c27d26` | `36723fe43740b81bb265a33d67207f6cc46784843f906d4265c4013698c27d26` | ⑦ blind-review r2 receipt verbatim |
| `04-reviews/07-blind-r2b.md` | `04-reviews/07-blind-r2b.md` | `e9a6d402a080a6afb03748f35b922bc36e3e0852b425ef87438620aec36606fe` | `e9a6d402a080a6afb03748f35b922bc36e3e0852b425ef87438620aec36606fe` | ⑦ blind-review r2b receipt verbatim |
| `06-gates/gate-readings.txt` | `06-gates/gate-readings.txt` | `8c6affb26c8cb47cc6ccee9ffb903fcb4dacb7f10bd9c781a62fb04f23ff5714` | `8c6affb26c8cb47cc6ccee9ffb903fcb4dacb7f10bd9c781a62fb04f23ff5714` | ⑨ gate readings (three scopes + selftest) |

## 归档条目清单

| 路径 | 原始字节 sha256 | 归档字节 sha256 | 说明 |
|---|---|---|---|
| `00-briefing.md` | 2ba0b8ffa3b18b67ba924fd0dfbfae13ab5340b5758ea3aa8439b66863e905d4 | 2ba0b8ffa3b18b67ba924fd0dfbfae13ab5340b5758ea3aa8439b66863e905d4 | briefing verbatim |
| `freedom-list.md` | 38274467d21d4ee57a98f66e20f2d1da350a04298b2e8bd242b0672d0745fc56 | 38274467d21d4ee57a98f66e20f2d1da350a04298b2e8bd242b0672d0745fc56 | freedom list |
| `input-list.md` | 5f1e42a3979be0959d77de79816d22b14e9f3dd9bd941e576f491d8e2fc64294 | 5f1e42a3979be0959d77de79816d22b14e9f3dd9bd941e576f491d8e2fc64294 | input list |
| `01-arch/design-r1.md` | 90376ddf8039675d7762fd3f3b67a126ed0bbc9e0793d10f9e44a4ed85767d60 | 90376ddf8039675d7762fd3f3b67a126ed0bbc9e0793d10f9e44a4ed85767d60 | arch three-way design (① output) |
| `01-arch/ruling-2026-11-06-p60.md` | 3a1758666e031c5da23ff1ae9c17f4064399eedf82c9ab85e2b11ad095a00c01 | 3a1758666e031c5da23ff1ae9c17f4064399eedf82c9ab85e2b11ad095a00c01 | user ruling verbatim (phase A) |
| `02-impl/impl-r1.md` | 11fdc59184460fd1baef8d2dd58900f81fed746cfd40b50fbd05f8fac6f7ed98 | 11fdc59184460fd1baef8d2dd58900f81fed746cfd40b50fbd05f8fac6f7ed98 | ② implementer receipt verbatim |
| `03-test/test-r1.md` | f258268abdcda60786eb958d1398fe8e23fd8061d0b21f8080a3abdf04e9c745 | f258268abdcda60786eb958d1398fe8e23fd8061d0b21f8080a3abdf04e9c745 | ③ tester receipt verbatim |
| `03-test/test-r2.md` | 2bd99c099195fcfacc1e5bcab9ce99f1677ebd5d9608fde2a1f252d0255719eb | 2bd99c099195fcfacc1e5bcab9ce99f1677ebd5d9608fde2a1f252d0255719eb | ③ tester receipt r2 verbatim |
| `04-reviews/05-pre-review-r1.md` | 216f7d02100a31d5844d96592b14ff574f5f5c6b3f7b1fde99b7975c8588daf0 | 216f7d02100a31d5844d96592b14ff574f5f5c6b3f7b1fde99b7975c8588daf0 | ⑤ pre-review receipt verbatim |
| `04-reviews/05-pre-review-r2.md` | b7687453f402adb0f0ad221f4cc94f928a84d5d24da05e69a0cf849d72011f72 | b7687453f402adb0f0ad221f4cc94f928a84d5d24da05e69a0cf849d72011f72 | ⑤ pre-review r2 receipt verbatim |
| `04-reviews/06-scan-r1.md` | 44c0e68a865141dd66c0ccafca6439ae78b4b9bfe3ad9e6d71adfcbeb262a865 | 44c0e68a865141dd66c0ccafca6439ae78b4b9bfe3ad9e6d71adfcbeb262a865 | ⑥ independent-scan receipt verbatim |
| `04-reviews/07-final-r1.md` | 14b2a54fefa1498da898838f657b2f006f05b37733acd17332da39af54a22761 | 14b2a54fefa1498da898838f657b2f006f05b37733acd17332da39af54a22761 | ⑦ final-review receipt verbatim |
| `04-reviews/07-blind-r1.md` | ee01b330220d96ab437983fb5323cdb24f8138857c7d1e92f00a6b105f7e0c03 | ee01b330220d96ab437983fb5323cdb24f8138857c7d1e92f00a6b105f7e0c03 | ⑦ blind-review receipt verbatim |
| `04-reviews/07-blind-r2.md` | 36723fe43740b81bb265a33d67207f6cc46784843f906d4265c4013698c27d26 | 36723fe43740b81bb265a33d67207f6cc46784843f906d4265c4013698c27d26 | ⑦ blind-review r2 receipt verbatim |
| `04-reviews/07-blind-r2b.md` | e9a6d402a080a6afb03748f35b922bc36e3e0852b425ef87438620aec36606fe | e9a6d402a080a6afb03748f35b922bc36e3e0852b425ef87438620aec36606fe | ⑦ blind-review r2b receipt verbatim |
| `06-gates/gate-readings.txt` | 8c6affb26c8cb47cc6ccee9ffb903fcb4dacb7f10bd9c781a62fb04f23ff5714 | 8c6affb26c8cb47cc6ccee9ffb903fcb4dacb7f10bd9c781a62fb04f23ff5714 | ⑨ gate readings (three scopes + selftest) |
