# Task Package / Phase Close-out Report Template (SLICE-CLOSEOUT-TEMPLATE)

> Usage: at a task-package / stage close-out write the closing report from this template (the trigger anchors are listed in section 4, the end of T027 being the **first instance**); the **retrospective is written by the driver**, while the **prospective is a pending-question list drafted by the driver and the business direction is ruled by the user** (per §1.7.3 and §3.2 of `docs/DELIVERY_DIRECTIVE.md`). This template is a structural master copy with no rule body; where it conflicts with the directive, the directive prevails.

## 1. Closing Object

The target this report covers: `<task-package-id>` (task package) and `<阶段>` (stage identifier), plus slice scope, time range and premises (one sentence).

## 2. Retrospective (written by the driver)

Summarise verifiable facts only (taken from the slice reports and the ledger, never from memory), in four subsections.

### 2.1 What was done

The completed slices / milestones item by item, each with a report path or commit id.

### 2.2 Done and not done

Completed items listed as a check list; every unfinished item carries its reason.

### 2.3 Outstanding roll-up

Roll up unfinished items and their root causes across slices, so that cross-slice omissions become visible.

### 2.4 Review and gates

Summarise the review verdicts and gate readings of each slice; cite the reports, do not re-run.

## 3. Prospective (driver-drafted + user ruling)

Two subsections: the pending-question list (drafted by the driver) and the ruling record (ruled by the user). **The business direction must not be decided by the driver alone** (per §1.7.3 and §3.2).

### 3.1 Pending-question list (drafted by the driver)

One item per line: question / candidate options / a suggestion (optional) / who rules on it.

### 3.2 User ruling record

One item per line: question / user ruling / date / landing point (OB / ADR / a slice to be opened).

## 4. Anchor and Write-Back

State the trigger anchor (**end of T027 (first instance)** / end of S1 / end of S2 / ...) and the write-back targets (`REMAINING_SLICES{,_EN}`, `DEV_PLAN{,_EN}`, OB / ADR).
