# P92 ⑦ 收窄确认读 r4（原文，逐字）

- 载体：`independent-reviewer` + `GPT-5.3-Codex (copilot)`；时点：2026-11-26；轮次：4（用户同轮授权的额度例外，仅 1 次，不得作为先例）。
- 输入：收窄为三处上下文（§5 G4b 读数 / 围栏正文与外置注记 / 白名单登记处）。
- 结论：`CONFIRM: FINDINGS`（A 处读数陈旧）。

---

A: NONE 依据：报告行写为 G4b（tree/index）newCjk=435 unknown=[]，而捕获原行为 PASS G4b CJK diff minus whitelist :: newCjk=466 unknown=[]（tree 与 index 均如此）；ci 为 newCjk=494 unknown=[]，故未逐字一致。  

B: CLOSED 依据：Section A→E 与 INDEPENDENT SCAN: PASS 的段序一致；围栏内差异只出现在计数口径（共16→计16、共18→合计18、175行/140非空行→行数175、140条非空行）；围栏外存在“转录归一注记…四处措辞…其余逐字保留”。  

C: CLOSED 依据：白名单文件当前仅注释头（数据行=0），且断言原句为 T23 no entry was preloaded into tools/gates/g6-whitelist.txt (the exemption above is in-memory only)，与“白名单应为空、豁免仅内存态”一致，未见以登记掩盖缺陷。  

CONFIRM: FINDINGS

---

## 主控处置（事实记录）

- A 处发现成立：报告 §5 的 tree/index 读数 `435` 与捕获件 `466` 不一致（该值是变更集较小时的旧读数）。已按捕获件改为 `466`，CN/EN 同步；`ci` 的 `494` 未变。
- B / C 两处 `CLOSED`。
- 按用户裁定「若确认读未 PASS ⇒ 停下、不得提交」，本轮**未提交、未推送**，等用户裁决是否追加一次收窄确认读。
