# B3a evidence bundle 2026-09-20

Reproduction commands (host PowerShell 5.1; see docs/t027/B3A_RIG_DESIGN.md S10):

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\New-B3aSession.ps1 `
  -VmName Win11B3 -CredFile "D:\VirtualBox VMs\Win11B3\host-only\cred.txt"

powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aCalibration.ps1 `
  -VmName Win11B3 -CredFile "D:\VirtualBox VMs\Win11B3\host-only\cred.txt" -N 5

powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aRound.ps1 `
  -Round 3 -Candidate negative-control

powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Test-B3aRig.ps1

powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\New-B3aEvidence.ps1
```

Premise: the full premise and coverage statement is docs/t027/B3A_RIG_DESIGN.md S1
(hypervisor VirtualBox 7.2.18 r175117, guest Windows 11 Pro 25H2 build 26200.8037
zh-CN retail unactivated, VBoxSVGA graphics adapter, no guest NIC, SATA/AHCI
useHostIOCache=false, VDI backend, host NVMe, NTFS, frozen snapshot 'pristine').
Conclusions cover only this premise (design S1.8).
Guest clocks are secondary evidence only; the journal hostTs is the authority.
Calibration results prove only the rig's discriminating power, never a product
candidate (C2/C3 conclusions belong to B3b).

Encoding: .ps1/.md files carry UTF-8 BOM + LF; .json/.txt are UTF-8 without BOM
+ LF (CODING_CONVENTIONS.md).