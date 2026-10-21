#!/usr/bin/env bash
#
# C1 - Linux AgentRunner native validation.
#
# Scope: this driver mirrors the f1..f4 groups and the CLI leg of
# tools/b4-e2e/run-faces.ps1 against the FROZEN chain definitions shipped in
# docs/validation/evidence/directive-history/b4/archive/input.
#
# It is a validation harness only: it never edits a source file, never commits,
# never pushes and never touches the network (GOPROXY=off + GOTOOLCHAIN=local).
#
# Face f5 is intentionally absent. Per the design of slice C1, f5 is the only
# face that declares a Windows-only dependency (the gate hook runs `cmd.exe`);
# it is not part of this native Linux matrix. The defensive branch in run_case()
# still exists so that a face which needs `cmd.exe` would be recorded as
# `skipped` with its reason instead of being silently dropped.
#
# Matrix (one invocation per phase; the case column below is the run-root label,
# where a `-r2` suffix marks the in-pack control re-run of the same case):
#
#   f1  pos, neg, neg-r2, neg2, neg2-r2, neg3        phase: run
#   f2  pos, neg, neg-r2                             phase: run
#   f3  pos, neg, neg-r2                             phase: run
#   f4  pos (park -> recover), neg (park -> forge), neg-r2 (park -> forge)
#   cli validate / preview on b4-chain.json and b4-chain-hook.json, plus emit-record
#
# Inputs (environment variables, every path is derived from these):
#   C1_REPO  repository root            (default: three levels above this script)
#   C1_BIN   build output directory     (default: $C1_REPO/tmp/c1/bin)
#   C1_PACK  frozen pack directory      (default: $C1_REPO/docs/validation/evidence/directive-history/b4/archive/input)
#   C1_OUT   evidence output directory  (default: $C1_REPO/tmp/c1/out)
#
# Output layout:
#   $C1_OUT/toolchain-sha256.txt            built binaries: sha256 + version probe
#   $C1_OUT/summary.tsv                     one row per invocation
#   $C1_OUT/cli/<phase>.{stdout,stderr}.txt CLI leg streams, $C1_OUT/cli/driver.log
#   $C1_OUT/<face>/driver.log               per-face command/exit/verdict log
#   $C1_OUT/<face>/<case>/<phase>.stdout.txt
#   $C1_OUT/<face>/<case>/<phase>.stderr.txt
#   $C1_OUT/<face>/<case>/verdict-<phase>.json   written by the harness, never edited here
#   $C1_OUT/runs/<face>-<case>/             run roots (workspace copied from the seed)
#   $C1_OUT/inputs/b4-inputs.json           declared inputs emitted once by -emit-record
#
# summary.tsv columns (tab separated, no header, one row per invocation):
#   face <TAB> case <TAB> phase <TAB> exit <TAB> verdict_path
# A non-zero exit is recorded and the matrix continues. A row that cannot run is
# written with the literal `skipped` in the exit column and the reason string in
# the verdict_path column, so a skip is never silent and never path-like.
#
# Exit status: 0 when every invocation ran (skips are allowed and recorded),
# 1 when a precondition, a build, or at least one invocation failed.

set -uo pipefail

# ---------------------------------------------------------------- helpers

die() {
  printf 'c1: %s\n' "$1" >&2
  exit 1
}

log_line() {
  # Usage: log_line <log_file> <message...>
  # The directory is created here: a log line is bookkeeping that must never be
  # lost because the caller has not created the face directory yet.
  local log_file=$1
  shift
  mkdir -p -- "$(dirname -- "$log_file")" \
    || die "cannot create the log directory for $log_file"
  printf '[%s] %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$*" >>"$log_file" \
    || die "cannot append to $log_file"
}

sha256_file() {
  local path=$1 digest
  if command -v sha256sum >/dev/null 2>&1; then
    digest=$(sha256sum -- "$path") || return 1
  elif command -v shasum >/dev/null 2>&1; then
    digest=$(shasum -a 256 -- "$path") || return 1
  else
    die 'neither sha256sum nor shasum is available'
  fi
  printf '%s' "${digest%% *}"
}

command_line() {
  # Render an argv as a single quoted command line for the driver logs.
  local parts=() arg
  for arg in "$@"; do
    parts+=("$(printf '%q' "$arg")")
  done
  printf '%s' "${parts[*]}"
}

probe_version() {
  # Echo the first line of a working version probe, or `unavailable`.
  # Go's flag package accepts both -flag and --flag, so `--version` covers
  # agent-stub; prfrail answers the `version` subcommand instead.
  local bin=$1 out
  if out=$("$bin" --version 2>&1) && [ -n "$out" ]; then
    printf '%s' "$(printf '%s\n' "$out" | head -n 1)"
    return 0
  fi
  if out=$("$bin" version 2>&1) && [ -n "$out" ]; then
    printf '%s' "$(printf '%s\n' "$out" | head -n 1)"
    return 0
  fi
  printf 'unavailable'
}

extract_parent_hash() {
  awk -F'"' '/"parentSnapshotHash"/ { print $4; exit }' "$1" 2>/dev/null
}

# ---------------------------------------------------------------- configuration

C1_SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd) \
  || die 'cannot resolve the script directory'

C1_REPO=${C1_REPO:-$(cd -- "$C1_SCRIPT_DIR/../../.." && pwd)}
C1_REPO=$(cd -- "$C1_REPO" 2>/dev/null && pwd) || die "C1_REPO is not a directory: ${C1_REPO:-}"
C1_BIN=${C1_BIN:-"$C1_REPO/tmp/c1/bin"}
C1_PACK=${C1_PACK:-"$C1_REPO/docs/validation/evidence/directive-history/b4/archive/input"}
C1_OUT=${C1_OUT:-"$C1_REPO/tmp/c1/out"}

CHAIN="$C1_PACK/b4-chain.json"
CHAIN_HOOK="$C1_PACK/b4-chain-hook.json"
FROZEN_INPUTS="$C1_PACK/b4-inputs.json"

C1_SEED="$C1_OUT/seed"
C1_INPUTS_DIR="$C1_OUT/inputs"
INPUTS="$C1_INPUTS_DIR/b4-inputs.json"
SUMMARY="$C1_OUT/summary.tsv"
TOOLCHAIN="$C1_OUT/toolchain-sha256.txt"

PRFRAIL="$C1_BIN/prfrail"
STUB="$C1_BIN/agent-stub"
HARNESS="$C1_BIN/b4-e2e"

# Fail closed on a missing repository or a missing frozen pack.
[ -f "$C1_REPO/go.mod" ] || die "C1_REPO does not look like the prfrail module (go.mod missing): $C1_REPO"
[ -f "$CHAIN" ] || die "frozen chain definition missing: $CHAIN"
[ -f "$CHAIN_HOOK" ] || die "frozen hook chain definition missing: $CHAIN_HOOK"
command -v go >/dev/null 2>&1 || die 'go toolchain not found in PATH'
command -v sha256sum >/dev/null 2>&1 || command -v shasum >/dev/null 2>&1 \
  || die 'neither sha256sum nor shasum is available'

# No network: module downloads and toolchain switches are hard failures instead.
export GOPROXY=off
export GOTOOLCHAIN=local

mkdir -p -- "$C1_BIN" "$C1_OUT" "$C1_SEED" "$C1_INPUTS_DIR" "$C1_OUT/build" "$C1_OUT/cli" \
  || die "cannot create the output directories under $C1_OUT"

# The seed must be byte-identical to the one run-faces.ps1 uses: the declared
# inputs derive parentSnapshotHash from it, and the frozen pack records that hash.
printf 'b4 seed workspace\n' >"$C1_SEED/b4-seed.txt" || die "cannot write the seed workspace: $C1_SEED/b4-seed.txt"

: >"$SUMMARY" || die "cannot initialise $SUMMARY"
: >"$C1_OUT/build/driver.log" || die "cannot initialise $C1_OUT/build/driver.log"
: >"$C1_OUT/cli/driver.log" || die "cannot initialise $C1_OUT/cli/driver.log"

# ---------------------------------------------------------------- invocation

invoke() {
  # Usage: invoke <log_file> <stdout_file> <stderr_file> <cwd> <command> [args...]
  # Writes the streams, appends the exact command line + exit code to the log and
  # returns the command's exit code. Never aborts on a non-zero exit.
  local log_file=$1
  local stdout_file=$2
  local stderr_file=$3
  local workdir=$4
  shift 4
  local prog=$1
  shift

  local exit_code=0
  mkdir -p -- "$(dirname -- "$stdout_file")" "$(dirname -- "$stderr_file")" \
    || die "cannot create the log directory for $stdout_file"

  ( cd -- "$workdir" && "$prog" "$@" ) >"$stdout_file" 2>"$stderr_file" || exit_code=$?

  log_line "$log_file" "exit=$exit_code cmd=$(command_line "$prog" "$@")"
  log_line "$log_file" "  cwd=$workdir stdout=$stdout_file stderr=$stderr_file"
  return "$exit_code"
}

# ---------------------------------------------------------------- build

build_binaries() {
  local definitions=(
    'prfrail|./cmd/prfrail'
    'agent-stub|./tools/agent-stub'
    'b4-e2e|./tools/b4-e2e'
  )
  local definition name pkg output exit_code

  : >"$TOOLCHAIN" || die "cannot initialise $TOOLCHAIN"
  {
    printf '# C1 toolchain (native build)\n'
    printf '# repo: %s\n' "$C1_REPO"
    printf '# go: %s\n' "$(go version 2>/dev/null || printf 'unavailable')"
    printf '# built: %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  } >>"$TOOLCHAIN"

  for definition in "${definitions[@]}"; do
    name=${definition%%|*}
    pkg=${definition#*|}
    output="$C1_BIN/$name"
    exit_code=0
    invoke "$C1_OUT/build/driver.log" \
      "$C1_OUT/build/$name.build.stdout.txt" \
      "$C1_OUT/build/$name.build.stderr.txt" \
      "$C1_REPO" go build -o "$output" "$pkg" || exit_code=$?
    [ "$exit_code" -eq 0 ] \
      || die "go build failed for $pkg (exit $exit_code); see $C1_OUT/build/$name.build.stderr.txt"
    [ -x "$output" ] || die "go build did not produce an executable: $output"
    printf '%s  sha256=%s  version=%s\n' "$name" "$(sha256_file "$output")" "$(probe_version "$output")" \
      >>"$TOOLCHAIN"
  done
}

# ---------------------------------------------------------------- recording

record() {
  # Usage: record <face> <case> <phase> <exit> <verdict_path_or_reason>
  local face=$1 case_label=$2 phase=$3 exit_code=$4 verdict_path=$5
  printf '%s\t%s\t%s\t%s\t%s\n' "$face" "$case_label" "$phase" "$exit_code" "$verdict_path" \
    >>"$SUMMARY" || die "cannot append to $SUMMARY"
}

# ---------------------------------------------------------------- faces

face_needs_cmd_exe() {
  # Only the f5 gate face declares a `cmd.exe` hook runner. f5 is not part of this
  # matrix; the predicate exists so the skip branch below stays honest if it ever is.
  case "$1" in
    f5) return 0 ;;
    *) return 1 ;;
  esac
}

run_case() {
  # Usage: run_case <face> <case_label> <case_id> <phase> <run_root> <chain>
  local face=$1
  local case_label=$2
  local case_id=$3
  local phase=$4
  local run_root=$5
  local chain=$6

  local face_log="$C1_OUT/$face/driver.log"
  local case_dir="$C1_OUT/$face/$case_label"
  local verdict="$case_dir/verdict-$phase.json"
  local stdout_file="$case_dir/$phase.stdout.txt"
  local stderr_file="$case_dir/$phase.stderr.txt"

  # A face whose declared dependency is absent is recorded with its reason, never
  # dropped. For C1's matrix (f1..f4) this branch is defensive only.
  if face_needs_cmd_exe "$face" && ! command -v cmd.exe >/dev/null 2>&1; then
    local reason='cmd.exe is not available on this host'
    record "$face" "$case_label" "$phase" 'skipped' "$reason"
    log_line "$face_log" "skipped $case_label/$phase: $reason"
    return 0
  fi

  # Missing preconditions abort: a partially built matrix must not look complete.
  [ -x "$HARNESS" ] || die "harness binary missing: $HARNESS"
  [ -x "$STUB" ] || die "stub binary missing: $STUB"
  [ -f "$INPUTS" ] || die "declared inputs missing: $INPUTS"
  [ -f "$chain" ] || die "chain definition missing: $chain"
  [ -d "$run_root/workspace" ] || die "run root workspace missing: $run_root/workspace"

  mkdir -p -- "$case_dir" || die "cannot create the case directory: $case_dir"

  local exit_code=0
  invoke "$face_log" "$stdout_file" "$stderr_file" "$C1_REPO" "$HARNESS" \
    -face "$face" -case "$case_id" -run-root "$run_root" \
    -inputs "$INPUTS" -chain "$chain" -stub "$STUB" \
    -phase "$phase" -verdict "$verdict" || exit_code=$?

  # 126/127 mean the shell could not start the harness at all: that is a broken
  # toolchain, not a recorded face result.
  if [ "$exit_code" -eq 126 ] || [ "$exit_code" -eq 127 ]; then
    die "harness could not be executed (exit $exit_code): $HARNESS"
  fi

  local verdict_state='missing'
  [ -f "$verdict" ] && verdict_state='present'
  record "$face" "$case_label" "$phase" "$exit_code" "$verdict"
  log_line "$face_log" "recorded $face/$case_label phase=$phase exit=$exit_code verdict=$verdict ($verdict_state)"
  return 0
}

run_group() {
  # Usage: run_group <face> <case_label> <case_id> <chain> <phase...>
  # A group owns exactly ONE run root: multi-phase groups (f4 park -> recover/forge)
  # must continue the same run root, otherwise a relaunch would measure a fresh run.
  local face=$1
  local case_label=$2
  local case_id=$3
  local chain=$4
  shift 4
  local phases=("$@")

  local run_root="$C1_OUT/runs/$face-$case_label"
  local face_log="$C1_OUT/$face/driver.log"
  local seed_file phase

  [ -f "$chain" ] || die "chain definition missing: $chain"
  rm -rf -- "$run_root" || die "cannot reset the run root: $run_root"
  mkdir -p -- "$run_root/workspace" || die "cannot create the run root workspace: $run_root/workspace"
  for seed_file in "$C1_SEED"/*; do
    [ -e "$seed_file" ] || die "seed workspace is empty: $C1_SEED"
    cp -f -- "$seed_file" "$run_root/workspace/" || die "cannot seed the workspace: $run_root/workspace"
  done

  log_line "$face_log" "group $face-$case_label case=$case_id chain=$chain run_root=$run_root phases=${phases[*]}"
  for phase in "${phases[@]}"; do
    run_case "$face" "$case_label" "$case_id" "$phase" "$run_root" "$chain"
  done
}

# ---------------------------------------------------------------- CLI leg

run_cli_leg() {
  local log="$C1_OUT/cli/driver.log"
  local exit_code

  [ -x "$PRFRAIL" ] || die "CLI binary missing: $PRFRAIL"

  # The four CLI invocations of the frozen pack (both chains, validate and preview).
  exit_code=0
  invoke "$log" "$C1_OUT/cli/validate.stdout.txt" "$C1_OUT/cli/validate.stderr.txt" \
    "$C1_REPO" "$PRFRAIL" validate --chain "$CHAIN" || exit_code=$?
  record cli - validate "$exit_code" -

  exit_code=0
  invoke "$log" "$C1_OUT/cli/preview.stdout.txt" "$C1_OUT/cli/preview.stderr.txt" \
    "$C1_REPO" "$PRFRAIL" preview --chain "$CHAIN" || exit_code=$?
  record cli - preview "$exit_code" -

  exit_code=0
  invoke "$log" "$C1_OUT/cli/validate-hook.stdout.txt" "$C1_OUT/cli/validate-hook.stderr.txt" \
    "$C1_REPO" "$PRFRAIL" validate --chain "$CHAIN_HOOK" || exit_code=$?
  record cli - validate-hook "$exit_code" -

  exit_code=0
  invoke "$log" "$C1_OUT/cli/preview-hook.stdout.txt" "$C1_OUT/cli/preview-hook.stderr.txt" \
    "$C1_REPO" "$PRFRAIL" preview --chain "$CHAIN_HOOK" || exit_code=$?
  record cli - preview-hook "$exit_code" -

  # The declared inputs every face shares are emitted once, exactly as run-faces.ps1
  # does; a failure here invalidates the whole matrix, so it aborts.
  exit_code=0
  invoke "$log" "$C1_OUT/cli/emit-record.stdout.txt" "$C1_OUT/cli/emit-record.stderr.txt" \
    "$C1_REPO" "$HARNESS" -emit-record -chain "$CHAIN" -workspace "$C1_SEED" -out "$INPUTS" \
    || exit_code=$?
  record cli - emit-record "$exit_code" -
  [ "$exit_code" -eq 0 ] || die "emit-record failed (exit $exit_code); see $C1_OUT/cli/emit-record.stderr.txt"
  [ -f "$INPUTS" ] || die "emit-record did not write the declared inputs: $INPUTS"

  log_line "$log" "declared inputs: path=$INPUTS sha256=$(sha256_file "$INPUTS")"
  # Informational only: the emitted parent hash must reproduce the frozen pack so
  # the faces replay against the same declared baseline.
  local emitted_parent='' frozen_parent=''
  emitted_parent=$(extract_parent_hash "$INPUTS")
  if [ -f "$FROZEN_INPUTS" ]; then
    frozen_parent=$(extract_parent_hash "$FROZEN_INPUTS")
  fi
  log_line "$log" "parentSnapshotHash emitted=$emitted_parent frozen=$frozen_parent match=$([ "$emitted_parent" = "$frozen_parent" ] && printf yes || printf no)"
}

# ---------------------------------------------------------------- summary

summary() {
  local face
  for face in f1 f2 f3 f4; do
    awk -F'\t' -v face="$face" '
      $1 == face {
        if ($4 == "skipped") { skipped++ }
        else if ($4 == "0") { ok++ }
        else { fail++ }
      }
      END { printf "%s: ok=%d fail=%d skipped=%d\n", face, ok + 0, fail + 0, skipped + 0 }
    ' "$SUMMARY"
  done
  printf 'out: %s\n' "$C1_OUT"
}

# ---------------------------------------------------------------- main

build_binaries
run_cli_leg

# f1: 4 cases (pos, neg, neg2, neg3); neg and neg2 carry the r2 control re-run.
run_group f1 pos     pos    "$CHAIN" run
run_group f1 neg     neg    "$CHAIN" run
run_group f1 neg-r2  neg    "$CHAIN" run
run_group f1 neg2    neg2   "$CHAIN" run
run_group f1 neg2-r2 neg2   "$CHAIN" run
run_group f1 neg3    neg3   "$CHAIN" run

# f2 / f3: 2 cases each, the negative case re-run as a whole group.
run_group f2 pos     pos    "$CHAIN" run
run_group f2 neg     neg    "$CHAIN" run
run_group f2 neg-r2  neg    "$CHAIN" run

run_group f3 pos     pos    "$CHAIN" run
run_group f3 neg     neg    "$CHAIN" run
run_group f3 neg-r2  neg    "$CHAIN" run

# f4: 2 cases, each a two-phase group on ONE run root; the negative case is
# re-run as a whole group (park then forge) on a second, independent run root.
run_group f4 pos     pos    "$CHAIN" park recover
run_group f4 neg     neg    "$CHAIN" park forge
run_group f4 neg-r2  neg    "$CHAIN" park forge

summary

failures=$(awk -F'\t' '$4 != "0" && $4 != "skipped" { count++ } END { printf "%d", count + 0 }' "$SUMMARY")
if [ "$failures" -ne 0 ]; then
  printf 'c1: %d invocation(s) exited non-zero; see %s\n' "$failures" "$SUMMARY" >&2
  exit 1
fi
exit 0
