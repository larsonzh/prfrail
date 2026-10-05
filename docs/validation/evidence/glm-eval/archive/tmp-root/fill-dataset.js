// Fill the GLM-eval dataset with the six realised calls (probe + 5 leg calls) and honest metering notes.
const fs = require('fs');
const path = require('path');
const ROOT = 'd:/LZProjects/prfrail';
const DS = path.join(ROOT, 'docs', 'validation', 'directive-glm-eval-dataset.json');
const ds = JSON.parse(fs.readFileSync(DS, 'utf8'));

const glmOut = 'tmp/glm-eval-out/leg-glm.txt';
const ctlOut = 'tmp/glm-eval-out/leg-control.txt';
const noteCommon =
  'metering: GLM console readings were NOT available during the call window (operator away; no console pre-read taken) ' +
  'and the extension log/storage carry no usage field => per-call token attribution is unavailable; the console history ' +
  'remains the designated second source per the slice caliber.';

const rows = [
  {
    call_no: 1,
    role_position: 'LIMIT-NAME-PROBE-REJECTED',
    model_string: 'GLM-5.3',
    quota_package: 'glm-trial-2000w',
    input_summary: 'model-id probe attempt 1 (rejected by the tool: "Requested model not found"; no model invocation, no tokens)',
    input_hash: null,
    output_hash: null,
    conclusion_line: null,
    measured_tokens: null,
    metering_mode: 'unavailable',
    duration_s: null,
    quality_verdict: 'OK',
    notes: 'rejected before invocation; the error payload lists the available model ids and is the evidence for the naming caliber',
  },
  {
    call_no: 2,
    role_position: 'LIMIT-NAME-PROBE-ACCEPTED',
    model_string: 'GLM-5.3 (glm)',
    quota_package: 'glm-trial-2000w',
    input_summary: 'model-id probe attempt 2 (minimal prompt: "PROBE-OK" only)',
    input_hash: null,
    output_hash: null,
    conclusion_line: 'PROBE-OK',
    measured_tokens: null,
    metering_mode: 'unavailable',
    duration_s: null,
    quality_verdict: 'OK',
    notes: 'authoritative limit name for this slice = GLM-5.3 (glm); counted separately from the ⑤ leg per the slice definition',
  },
  {
    call_no: 3,
    role_position: '5-EXPERIMENT',
    model_string: 'GLM-5.3 (glm)',
    quota_package: 'glm-trial-2000w',
    input_summary: 'pre-registered synthetic bundle (00-plan.md + 01-clause-cn.txt + 02-diff.js + 03-prompt.txt), byte-identical to the frozen object',
    input_hash: 'bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6',
    output_hash: '54c895f9df821231d9049ac7819077b44335db83383b266d86f3c28fee137b5b',
    conclusion_line: 'PRE-REVIEW: FINDINGS',
    measured_tokens: null,
    metering_mode: 'unavailable',
    duration_s: null,
    quality_verdict: 'OK',
    notes:
      'recall_all=1.0 (7/7) | recall_high=1.0 (3/3) | sev_exact=0.7857 | conclusion_line=1 | FP=0 | abstentions=2 | ' +
      'extra_true=10 | s8_ok=true | Q=0.9571. Output saved byte-exactly from the tool raw capture. ' +
      'Key/rubric tension (declared for the review chain): the pre-registered key grades S4/S5/S6 as Medium while the ' +
      'rubric given to the model defines High as "重要分支缺失或判断恒真恒假"; under a rubric-consistent re-grading ' +
      'S4/S5/S6 would be High and sev_exact would read 1.0 (Q=1.0). Both readings are reported; the key is NOT altered. ',
  },
  {
    call_no: 4,
    role_position: '5-CONTROL',
    model_string: 'DeepSeek V4 Pro (deepseek)',
    quota_package: 'deepseek-n-a',
    input_summary: 'the SAME pre-registered bundle text as call_no 3 (byte-identical transcription)',
    input_hash: 'bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6',
    output_hash: '2fa1d9eeedae134c758174fff900311983a52b0eda8470889f4ff5d67f463a7f',
    conclusion_line: 'PRE-REVIEW: FINDINGS',
    measured_tokens: null,
    metering_mode: 'unavailable',
    duration_s: null,
    quality_verdict: 'OK',
    notes:
      'recall_all=1.0 (7/7) | recall_high=1.0 (3/3) | sev_exact=0.7857 | conclusion_line=1 | FP=0 | abstentions=1 | ' +
      'extra_true=8 | s8_ok=true | Q=0.9571. Control output is the driver verbatim transcription (the tool returned it inline ' +
      'instead of writing a raw-capture file); the experiment leg was saved byte-exactly from the tool capture - the ' +
      'asymmetry is declared. V4 Pro token usage is declared unavailable (third-party channel; no token field in the call log). ',
  },
  {
    call_no: 5,
    role_position: 'BATCH-a',
    model_string: 'GLM-5.3-Flash (glm)',
    quota_package: 'glm-flash-realname-newuser',
    input_summary: 'evidence-tree enumeration/statistics over a 40-row synthetic listing (inline)',
    input_hash: null,
    output_hash: null,
    conclusion_line: null,
    measured_tokens: null,
    metering_mode: 'interval-sum',
    duration_s: null,
    quality_verdict: 'OK',
    notes:
      'DELIVERABLE 5/5 items correct except item 4: file_count=40 OK, by-ext OK, total_bytes=209020 OK, top-level list OK, ' +
      'but the top-5-by-bytes selection is WRONG (it omitted 9403 and 9196 and listed 8840/8633/8426) => partial quality ' +
      'failure on the ranking sub-task. Batch legs were fired CONCURRENTLY (protocol deviation: the design said 依次执行) ' +
      'so per-call durations are not separable; declared.',
  },
  {
    call_no: 6,
    role_position: 'BATCH-b',
    model_string: 'GLM-5.3-Flash (glm)',
    quota_package: 'glm-flash-realname-newuser',
    input_summary: 'long-document extraction + number cross-check over the v1.24 report pair',
    input_hash: null,
    output_hash: null,
    conclusion_line: null,
    measured_tokens: null,
    metering_mode: 'interval-sum',
    duration_s: null,
    quality_verdict: 'OK',
    notes:
      'Fully correct: 328/328, scopes 39/2/4 with the time-point caveat, cost vector, commit 025158f, run 37278691747, ' +
      'and the CN/EN number cross-check returned ALL MATCH (verified against the actual files). Usable as-is.',
  },
  {
    call_no: 7,
    role_position: 'BATCH-c',
    model_string: 'GLM-5.3-Flash (glm)',
    quota_package: 'glm-flash-realname-newuser',
    input_summary: 'language-side check of the CN/EN report mirror (report pair)',
    input_hash: null,
    output_hash: null,
    conclusion_line: null,
    measured_tokens: null,
    metering_mode: 'interval-sum',
    duration_s: null,
    quality_verdict: 'OK',
    notes:
      '14 language-side items (漏译/错译/术语漂移/语序搭配), boundary respected: no pass/fail verdict, no structural judging, ' +
      'and it correctly refused to flag authoritative EN terms (driver / To be discussed) as drift. Usable as review input.',
  },
  {
    call_no: 8,
    role_position: 'BATCH-d',
    model_string: 'GLM-5.3-Flash (glm)',
    quota_package: 'glm-flash-realname-newuser',
    input_summary: 'archive MANIFEST + SHA256SUMS draft generation for a 6-item synthetic pack',
    input_hash: null,
    output_hash: null,
    conclusion_line: null,
    measured_tokens: null,
    metering_mode: 'interval-sum',
    duration_s: null,
    quality_verdict: 'OK',
    notes:
      'Usable draft: header count/bytes correct (6 / 4437), all six rows carry the given bytes+sha256, archive/ path flattening ' +
      'and dictionary order correct, SHA256SUMS uses <hash><2 spaces><path>. Minor caliber deviation: the empty 未入库 table ' +
      'labels its last column 归档理由 where the pack caliber uses 理由.',
  },
];

ds.calls = rows;
ds.metering_note =
  'The GLM console pre/post readings could not be taken (operator away during the call window; no exclusivity confirmation). ' +
  'Per-call token attribution is therefore unavailable; the console usage history remains the designated second source. ' +
  'This is registered as a slice deviation, not silently filled in.';
ds.batch_concurrency_note =
  'Protocol deviation: the four batch calls were issued concurrently (design said sequentially), so per-call durations are not ' +
  'separable and concurrent load may have influenced latency/quality. Declared in the slice report.';
fs.writeFileSync(DS, JSON.stringify(ds, null, 2) + '\n', 'utf8');
const b = fs.readFileSync(DS);
console.log('dataset rows=' + ds.calls.length + ' bytes=' + b.length + ' bom=' + (b[0] === 0xef) + ' crlf=' + b.includes(Buffer.from('\r\n')));
console.log(JSON.stringify(ds.pre_registration, null, 1));
