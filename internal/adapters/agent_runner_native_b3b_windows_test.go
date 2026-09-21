//go:build b3bnative && windows

package adapters

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// This file is the guest writer the B3b rig drives: it publishes one real product
// record through the production no-replace publication path and stops existing at
// a chosen publish stage so the host can cut power inside that stage.
//
// The seam story it depends on: a host-side timing kill cannot hit a chosen stage
// (A7 measured CP4 and CP5 as indistinguishable), so the writer installs the
// publication seams, writes a marker line when a stage is reached, and then waits
// at the target stage. Only the host can cut power, and it cuts only after the
// marker for the target stage is visible.
//
// Environment contract (all read once, all required unless noted):
//
//	B3B_GUEST_WRITER  must be set; otherwise the test skips so an ordinary
//	                  `go test -tags b3bnative` run stays harmless
//	B3B_ROUND         round identifier, also part of the published request id
//	B3B_CANDIDATE     c1 | c2 | c3 | c2c3
//	B3B_STAGE         the target publish stage to stop in
//	B3B_ROOT          replay store root the record is published into
//	B3B_MARKER        marker file the host polls
//	B3B_PLAN_ONLY     when set, emit the pre-commitment and exit without writing
//	B3B_PAIRED        when set, publish the request record first and then stop the
//	                  completion record inside the target stage (the ordering question)
//	B3B_BLOCK_SECONDS bound on the wait at the target stage (default 300)

const (
	b3bGuestWriterEnv = "B3B_GUEST_WRITER"
	b3bRoundEnv       = "B3B_ROUND"
	b3bCandidateEnv   = "B3B_CANDIDATE"
	b3bStageEnv       = "B3B_STAGE"
	b3bRootEnv        = "B3B_ROOT"
	b3bMarkerEnv      = "B3B_MARKER"
	b3bPlanOnlyEnv    = "B3B_PLAN_ONLY"
	b3bBlockEnv       = "B3B_BLOCK_SECONDS"
	b3bPairedEnv      = "B3B_PAIRED"
)

// b3bExitStageTimeout is the dedicated exit code for "the host never cut power
// inside the target stage". The rig has to treat that exit as a void round: the
// writer neither reports a result nor continues past the stage.
const b3bExitStageTimeout = 90

// b3bGuestWriterPlan is the pre-commitment the host journals before the measured
// write: where the record will land and the digest of the exact bytes that will be
// published. It is computed from the same encoder the publication path uses, so
// the host can verify the surviving record byte for byte after the reboot.
type b3bGuestWriterPlan struct {
	Round     string `json:"round"`
	Candidate string `json:"candidate"`
	Stage     string `json:"stage"`
	Path      string `json:"path"`
	Digest    string `json:"digest"`
	Size      int    `json:"size"`
}

// b3bGuestWriterStage is one marker line: the stage the writer reached, the syscall
// trace observed so far, and whether that stage is the one the host must cut in.
type b3bGuestWriterStage struct {
	Round     string   `json:"round"`
	Candidate string   `json:"candidate"`
	Stage     string   `json:"stage"`
	Path      string   `json:"path"`
	Target    bool     `json:"target"`
	Trace     []string `json:"trace"`
	Timestamp string   `json:"timestamp"`
}

func b3bEnvValue(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}

func b3bKnownStages() []string {
	return []string{
		replayStorePublishStageTempWritten,
		replayStorePublishStageTempSynced,
		replayStorePublishStageTempClosed,
		replayStorePublishStageRecordLinked,
		replayStorePublishStageParentSynced,
	}
}

// b3bStageIndex is the position of a stage in the publication order, or -1.
func b3bStageIndex(stage string) int {
	for index, candidateStage := range b3bKnownStages() {
		if candidateStage == stage {
			return index
		}
	}
	return -1
}

func b3bKnownCandidate(candidate string) bool {
	switch candidate {
	case "c1", "c2", "c3", "c2c3", "positive-control":
		return true
	default:
		return false
	}
}

func b3bBlockDuration() time.Duration {
	raw := b3bEnvValue(b3bBlockEnv)
	if raw == "" {
		return 300 * time.Second
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return 300 * time.Second
	}
	return time.Duration(seconds) * time.Second
}

// b3bAppendMarker appends one JSON line to the marker file. Durability is
// deliberately not attempted: the host reads the marker through a guest process,
// which is served by the same cache the writer wrote into, and the marker is a
// signal rather than part of the measurement.
func b3bAppendMarker(marker string, value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(marker, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := file.Write(append(encoded, '\n')); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

// b3bPublishedPayload reproduces the bytes the publication path writes: the
// canonical encoding of the record plus the terminating newline. The pre-commitment
// digest is computed from exactly this, so a digest mismatch after the reboot means
// the surviving record is not the record that was published.
func b3bPublishedPayload(record any) ([]byte, error) {
	canonical, err := evidence.EncodeCanonical(record)
	if err != nil {
		return nil, err
	}
	return append(canonical, '\n'), nil
}

func b3bDigest(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

// b3bPrimitiveFor installs the candidate's no-replace primitive and records what
// actually ran, because a candidate that quietly degraded to os.Link would
// otherwise produce a green cell that means nothing.
func b3bPrimitiveFor(candidate string, trace *[]string) func(temp string, path string) error {
	if b3bVariantDegrade {
		// Variant arm (design section 6): every candidate degrades to the baseline link, and
		// the trace keeps saying so - a candidate that quietly degrades must never be able to
		// produce a green cell.
		return func(temp string, path string) error {
			if err := os.Link(temp, path); err != nil {
				*trace = append(*trace, "os.Link:error:"+err.Error())
				return err
			}
			*trace = append(*trace, "os.Link")
			return nil
		}
	}
	if candidate == "c2" || candidate == "c2c3" {
		return func(temp string, path string) error {
			if err := b3bMoveFileExWriteThrough(temp, path); err != nil {
				*trace = append(*trace, "MoveFileExW:error:"+err.Error())
				return err
			}
			*trace = append(*trace, "MoveFileExW(WRITE_THROUGH,no-replace)")
			return nil
		}
	}
	return func(temp string, path string) error {
		// C1 is the current baseline: the hard link the store has always used.
		if err := os.Link(temp, path); err != nil {
			*trace = append(*trace, "os.Link:error:"+err.Error())
			return err
		}
		*trace = append(*trace, "os.Link")
		return nil
	}
}

// b3bParentSyncFor installs the candidate's parent-directory durability step. The
// positive control takes the baseline publication primitive and drains the device
// instead, so the control answers "can this rig observe a clean result at this payload
// and window" rather than "is some candidate primitive durable".
func b3bParentSyncFor(candidate string, trace *[]string, target string) func(directory string) error {
	if b3bVariantDegrade {
		// Variant arm: the durability step is a no-op for every arm, and the trace says so.
		return func(directory string) error {
			err := replayStoreSyncParentDirectory(directory)
			if err != nil {
				*trace = append(*trace, "platform-parent-sync:error:"+err.Error())
				return err
			}
			*trace = append(*trace, "platform-parent-sync")
			return nil
		}
	}
	if candidate == "positive-control" {
		return func(directory string) error {
			if err := b3bFlushDirectoryHandle(directory); err != nil {
				*trace = append(*trace, "FlushFileBuffers(dir):error:"+err.Error())
				return err
			}
			*trace = append(*trace, "FlushFileBuffers(dir)")
			if err := b3bFlushDeviceHandle(target); err == nil {
				*trace = append(*trace, "FlushFileBuffers(device)")
				return nil
			}
			// The device handle needs write access to the volume, which an ordinary guest
			// user does not have. The fallback flushes the record itself and says so: the
			// control must never claim a device flush it did not perform.
			if err := b3bFlushTargetFile(target); err != nil {
				*trace = append(*trace, "FlushFileBuffers(file):error:"+err.Error())
				return err
			}
			*trace = append(*trace, "FlushFileBuffers(file)")
			return nil
		}
	}
	if candidate == "c3" || candidate == "c2c3" {
		return func(directory string) error {
			if err := b3bFlushDirectoryHandle(directory); err != nil {
				*trace = append(*trace, "FlushFileBuffers(dir):error:"+err.Error())
				return err
			}
			*trace = append(*trace, "FlushFileBuffers(dir)")
			return nil
		}
	}
	return func(directory string) error {
		// C1 keeps the platform step: on Windows that is the no-op, which is
		// exactly why the durability claim is unproven today.
		err := replayStoreSyncParentDirectory(directory)
		if err != nil {
			*trace = append(*trace, "platform-parent-sync:error:"+err.Error())
			return err
		}
		*trace = append(*trace, "platform-parent-sync")
		return nil
	}
}

// b3bRunGuestWriter is the writer half. It never returns from the target stage:
// either the host cuts power, or the bounded wait expires and the process exits
// with the timeout code so the round is void rather than a result.
func b3bRunGuestWriter(t *testing.T) {
	t.Helper()
	round := b3bEnvValue(b3bRoundEnv)
	candidate := b3bEnvValue(b3bCandidateEnv)
	stage := b3bEnvValue(b3bStageEnv)
	root := b3bEnvValue(b3bRootEnv)
	marker := b3bEnvValue(b3bMarkerEnv)
	if round == "" || candidate == "" || stage == "" || root == "" || marker == "" {
		t.Fatalf("%s requires B3B_ROUND, B3B_CANDIDATE, B3B_STAGE, B3B_ROOT and B3B_MARKER", b3bGuestWriterEnv)
	}
	if !b3bKnownCandidate(candidate) {
		t.Fatalf("unknown B3B_CANDIDATE %q", candidate)
	}
	known := false
	for _, candidateStage := range b3bKnownStages() {
		if candidateStage == stage {
			known = true
		}
	}
	if !known {
		t.Fatalf("unknown B3B_STAGE %q", stage)
	}

	// The store is constructed before any seam is installed: constructing it
	// publishes the ownership record, which is not part of the measurement.
	store, err := newAgentRunnerReplayStoreAt(root)
	if err != nil {
		t.Fatalf("open replay store at %s: %v", root, err)
	}
	// The writer exercises the write path on a platform whose durability is
	// declared unproven; this is the same in-test override the existing
	// experiments use and it never lifts the declaration itself. The store's
	// write gate is policy, not the publication primitive, so bypassing it here is
	// deliberate and is recorded in the slice report.
	store.publicationDurability = PublishDurabilityProven

	paired := b3bEnvValue(b3bPairedEnv) != ""
	requestID := "b3b-" + strings.ToLower(round)
	request := replayStoreRequestRecord(t, requestID)
	// The measured object is a request record in the normal mode and a completion
	// record in the paired mode, so it is held as any.
	var measured any = request
	path, err := store.requestPath(requestID)
	if err != nil {
		t.Fatalf("derive request path: %v", err)
	}
	requestPath := path
	if paired {
		// Paired mode answers the ordering question: the request record is published
		// first and the completion record is the measured object, stopped inside its own
		// target stage. The plan therefore describes the completion record, and the
		// transcript keeps the two phases apart by the path each entry carries.
		measured = replayStoreCompletionRecord(t, request, "b3b-completion-"+strings.ToLower(round))
		path, err = store.completionPath(requestID)
		if err != nil {
			t.Fatalf("derive completion path: %v", err)
		}
	}
	payload, err := b3bPublishedPayload(measured)
	if err != nil {
		t.Fatalf("encode record: %v", err)
	}
	plan := b3bGuestWriterPlan{
		Round:     round,
		Candidate: candidate,
		Stage:     stage,
		Path:      path,
		Digest:    b3bDigest(payload),
		Size:      len(payload),
	}
	if err := b3bAppendMarker(marker, plan); err != nil {
		t.Fatalf("write pre-commitment: %v", err)
	}
	fmt.Printf("PLAN path=%s digest=%s size=%d\n", path, plan.Digest, plan.Size)
	if b3bEnvValue(b3bPlanOnlyEnv) != "" {
		return
	}

	trace := []string{}
	// blockStage is the stage the writer may not continue past. It stays empty during
	// the request phase of the paired mode, so that phase is published without being
	// stopped, and it is armed right before the measured publication.
	blockStage := ""
	if !paired {
		blockStage = stage
	}
	originalStage := replayStorePublishStageHook
	originalPrimitive := replayStorePublishPrimitiveHook
	originalParent := replayStoreParentSyncHook
	replayStorePublishPrimitiveHook = b3bPrimitiveFor(candidate, &trace)
	replayStoreParentSyncHook = b3bParentSyncFor(candidate, &trace, path)
	replayStorePublishStageHook = func(current string, currentPath string) {
		trace = append(trace, "stage:"+current)
		entry := b3bGuestWriterStage{
			Round:     round,
			Candidate: candidate,
			Stage:     current,
			Path:      currentPath,
			Target:    current == blockStage,
			Trace:     append([]string{}, trace...),
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		}
		if err := b3bAppendMarker(marker, entry); err != nil {
			t.Fatalf("write stage marker for %s: %v", current, err)
		}
		if current == blockStage {
			// The host cuts power here. Sleeping rather than blocking on a channel
			// keeps the runtime's deadlock detector quiet (the A7 lesson), and the
			// bound means a host that never cuts turns the round into a void
			// timeout instead of a hung session. os.Exit skips the deferred seam
			// restore, which is harmless because the process is over either way.
			time.Sleep(b3bBlockDuration())
			fmt.Printf("TIMEOUT %s\n", current)
			os.Exit(b3bExitStageTimeout)
		}
	}
	defer func() {
		replayStorePublishStageHook = originalStage
		replayStorePublishPrimitiveHook = originalPrimitive
		replayStoreParentSyncHook = originalParent
	}()

	if paired {
		// The request phase runs through the same seams with no armed target, so the
		// ordering question is asked about a request that the candidate mechanism really
		// published, not about a baseline request that happens to sit next to it.
		if err := writeReplayRecordNoReplace(requestPath, request); err != nil {
			fmt.Printf("FAIL request phase %v\n", err)
			t.Fatalf("publish request record: %v", err)
		}
		blockStage = stage
		trace = trace[:0]
	}
	if err := writeReplayRecordNoReplace(path, measured); err != nil {
		fmt.Printf("FAIL %v\n", err)
		t.Fatalf("publish record: %v", err)
	}
	// Reaching here means no host cut the power inside the target stage: either
	// the target stage never fired or the host is too slow. Both are void rounds.
	fmt.Printf("SURVIVED %s\n", stage)
}

// TestB3bNativeGuestWriter is the entry point the rig executes. It skips unless
// B3B_GUEST_WRITER is set, so an ordinary tagged run stays harmless.
func TestB3bNativeGuestWriter(t *testing.T) {
	if b3bEnvValue(b3bGuestWriterEnv) == "" {
		t.Skip("the B3b guest writer runs only when the rig drives it")
	}
	b3bRunGuestWriter(t)
}

// b3bWriterCommand builds the re-exec invocation of this test binary as the guest
// writer, which is how the rig runs it inside the guest.
func b3bWriterCommand(t *testing.T, env map[string]string) *exec.Cmd {
	t.Helper()
	command := exec.Command(os.Args[0], "-test.run=TestB3bNativeGuestWriter", "-test.v")
	command.Env = append(os.Environ(), b3bGuestWriterEnv+"=1")
	for name, value := range env {
		command.Env = append(command.Env, name+"="+value)
	}
	return command
}

// TestB3bNativeRecordBytesDeterministic pins that the record bytes the host
// pre-commits are reproducible: the same request id and pinned creation timestamp
// have to produce the same digest in two separate processes, otherwise the host's
// plan could never be checked against the surviving record.
func TestB3bNativeRecordBytesDeterministic(t *testing.T) {
	// Both measured objects have to be reproducible across processes: the host
	// pre-commits the digest it will later look for, and a record whose bytes drifted
	// between processes would surface as a plan-digest mismatch in the middle of a VM
	// run - i.e. after hours of cuts, for a reason that has nothing to do with
	// durability. The request record is covered by the normal mode, the completion
	// record by the paired mode.
	for _, paired := range []bool{false, true} {
		name := "request-record"
		if paired {
			name = "completion-record"
		}
		t.Run(name, func(t *testing.T) {
			digests := make([]string, 0, 2)
			for attempt := 0; attempt < 2; attempt++ {
				root := t.TempDir()
				marker := filepath.Join(t.TempDir(), "marker.jsonl")
				environment := map[string]string{
					b3bRoundEnv:     "r01",
					b3bCandidateEnv: "c1",
					b3bStageEnv:     replayStorePublishStageRecordLinked,
					b3bRootEnv:      root,
					b3bMarkerEnv:    marker,
					b3bPlanOnlyEnv:  "1",
				}
				if paired {
					environment[b3bPairedEnv] = "1"
				}
				command := b3bWriterCommand(t, environment)
				output, err := command.CombinedOutput()
				if err != nil {
					t.Fatalf("plan run failed: %v\n%s", err, output)
				}
				lines, err := b3bReadMarkerLines(marker)
				if err != nil {
					t.Fatal(err)
				}
				if len(lines) != 1 {
					t.Fatalf("plan mode must write exactly the pre-commitment, got %d lines", len(lines))
				}
				var plan b3bGuestWriterPlan
				if err := json.Unmarshal([]byte(lines[0]), &plan); err != nil {
					t.Fatalf("pre-commitment is not the expected shape: %v", err)
				}
				if plan.Digest == "" || plan.Size == 0 {
					t.Fatalf("pre-commitment must carry a digest and a size: %+v", plan)
				}
				if paired && !strings.Contains(plan.Path, "completion.") {
					t.Fatalf("the paired plan must describe the completion record, got %s", plan.Path)
				}
				if !paired && !strings.Contains(plan.Path, "request.") {
					t.Fatalf("the normal plan must describe the request record, got %s", plan.Path)
				}
				digests = append(digests, plan.Digest)
			}
			if digests[0] != digests[1] {
				t.Fatalf("record bytes must be reproducible, got %s and %s", digests[0], digests[1])
			}
		})
	}
}

// TestB3bNativeWriterStopsAtTheTargetStage drives the writer with a short block so
// the whole handshake can be checked on the host without a VM: the marker has to
// contain the five stages in order, exactly one of them flagged as the target, a
// syscall trace naming what actually ran, and the writer has to exit with the
// timeout code instead of continuing past the stage it was told to stop in.
func TestB3bNativeWriterStopsAtTheTargetStage(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(t.TempDir(), "marker.jsonl")
	command := b3bWriterCommand(t, map[string]string{
		b3bRoundEnv:     "r02",
		b3bCandidateEnv: "c3",
		b3bStageEnv:     replayStorePublishStageRecordLinked,
		b3bRootEnv:      root,
		b3bMarkerEnv:    marker,
		b3bBlockEnv:     "1",
	})
	output, err := command.CombinedOutput()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("the writer must exit with the timeout code, got err=%v\n%s", err, output)
	}
	if exitErr.ExitCode() != b3bExitStageTimeout {
		t.Fatalf("writer exit code = %d, want %d\n%s", exitErr.ExitCode(), b3bExitStageTimeout, output)
	}
	lines, err := b3bReadMarkerLines(marker)
	if err != nil {
		t.Fatal(err)
	}
	// Nothing after the target stage may ever be reached: the marker holds the
	// pre-commitment plus the stages up to and including the target.
	wantStages := b3bKnownStages()[:b3bStageIndex(replayStorePublishStageRecordLinked)+1]
	if len(lines) != 1+len(wantStages) {
		t.Fatalf("marker must hold the pre-commitment plus the stages up to the target, got %d lines", len(lines))
	}
	var reached []string
	targets := 0
	for _, line := range lines[1:] {
		var entry b3bGuestWriterStage
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("stage marker is not the expected shape: %v", err)
		}
		reached = append(reached, entry.Stage)
		if entry.Target {
			targets++
		}
		if entry.Round != "r02" || entry.Candidate != "c3" {
			t.Fatalf("stage marker lost its round or candidate: %+v", entry)
		}
	}
	if strings.Join(reached, ",") != strings.Join(wantStages, ",") {
		t.Fatalf("stages = %v, want %v", reached, wantStages)
	}
	if targets != 1 {
		t.Fatalf("exactly one stage must be flagged as the target, got %d", targets)
	}
	// The last stage marker is the target: nothing later may be reached.
	if !strings.Contains(string(output), "TIMEOUT "+replayStorePublishStageRecordLinked) {
		t.Fatalf("the writer must report the timeout at the target stage:\n%s", output)
	}
	// C3 must show that the directory flush really ran (the parent stage is after
	// the target here, so the trace at the target still names the primitive).
	lastLine := lines[len(lines)-1]
	if !strings.Contains(lastLine, "os.Link") {
		t.Fatalf("the trace must name the primitive that ran: %s", lastLine)
	}
}

func b3bReadMarkerLines(marker string) ([]string, error) {
	raw, err := os.ReadFile(marker)
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return nil, nil
	}
	return strings.Split(trimmed, "\n"), nil
}

// TestB3bNativeCandidateTraceNamesTheRealPrimitive pins that no candidate cell can
// turn green because the candidate quietly degraded to the baseline: the trace the
// host reads has to name the primitive, and for the directory step the flush, that
// the candidate is supposed to execute.
func TestB3bNativeCandidateTraceNamesTheRealPrimitive(t *testing.T) {
	cases := []struct {
		candidate string
		stage     string
		want      string
	}{
		{"c1", replayStorePublishStageRecordLinked, "os.Link"},
		{"c2", replayStorePublishStageRecordLinked, "MoveFileExW(WRITE_THROUGH,no-replace)"},
		{"c3", replayStorePublishStageParentSynced, "FlushFileBuffers(dir)"},
	}
	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.candidate, func(t *testing.T) {
			root := t.TempDir()
			marker := filepath.Join(t.TempDir(), "marker.jsonl")
			command := b3bWriterCommand(t, map[string]string{
				b3bRoundEnv:     "r0" + strings.TrimPrefix(testCase.candidate, "c"),
				b3bCandidateEnv: testCase.candidate,
				b3bStageEnv:     testCase.stage,
				b3bRootEnv:      root,
				b3bMarkerEnv:    marker,
				b3bBlockEnv:     "1",
			})
			output, err := command.CombinedOutput()
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) || exitErr.ExitCode() != b3bExitStageTimeout {
				t.Fatalf("writer must stop in the target stage, got err=%v\n%s", err, output)
			}
			lines, err := b3bReadMarkerLines(marker)
			if err != nil {
				t.Fatal(err)
			}
			if len(lines) < 2 {
				t.Fatalf("expected a pre-commitment and stage markers, got %d lines", len(lines))
			}
			last := lines[len(lines)-1]
			if !strings.Contains(last, testCase.want) {
				t.Fatalf("the trace must name %q for candidate %s, got %s", testCase.want, testCase.candidate, last)
			}
		})
	}
}

// TestB3bNativePairedModePublishesRequestThenStopsInTheCompletion is the plumbing
// check for the ordering experiment: the request record has to be published for real
// before the completion record stops inside its own target stage, the pre-commitment
// has to describe the completion record, and exactly one stage entry may be flagged
// as the target.
func TestB3bNativePairedModePublishesRequestThenStopsInTheCompletion(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(t.TempDir(), "marker.jsonl")
	command := b3bWriterCommand(t, map[string]string{
		b3bRoundEnv:     "r07",
		b3bCandidateEnv: "c2",
		b3bStageEnv:     replayStorePublishStageRecordLinked,
		b3bRootEnv:      root,
		b3bMarkerEnv:    marker,
		b3bBlockEnv:     "1",
		b3bPairedEnv:    "1",
	})
	output, err := command.CombinedOutput()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != b3bExitStageTimeout {
		t.Fatalf("the paired writer must stop in the completion target stage, got err=%v\n%s", err, output)
	}
	lines, err := b3bReadMarkerLines(marker)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) < 3 {
		t.Fatalf("expected a pre-commitment and both phases, got %d lines", len(lines))
	}
	var plan b3bGuestWriterPlan
	if err := json.Unmarshal([]byte(lines[0]), &plan); err != nil {
		t.Fatalf("pre-commitment is not the expected shape: %v", err)
	}
	if !strings.Contains(plan.Path, "completion.b3b-r07.jsonl") {
		t.Fatalf("the pre-commitment must describe the completion record, got %s", plan.Path)
	}
	targets := 0
	requestPhase := 0
	completionPhase := 0
	for _, line := range lines[1:] {
		var entry b3bGuestWriterStage
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("stage marker is not the expected shape: %v", err)
		}
		if strings.Contains(entry.Path, "completion.b3b-r07.jsonl") {
			completionPhase++
		} else {
			requestPhase++
		}
		if entry.Target {
			targets++
			if !strings.Contains(entry.Path, "completion.b3b-r07.jsonl") {
				t.Fatalf("only the completion phase may be the target, got %s", entry.Path)
			}
		}
	}
	if targets != 1 {
		t.Fatalf("exactly one stage must be flagged as the target, got %d", targets)
	}
	if requestPhase == 0 || completionPhase == 0 {
		t.Fatalf("both phases must appear in the transcript, request=%d completion=%d", requestPhase, completionPhase)
	}
	// Both records really exist: the request phase published for real, and the
	// completion record was linked before its durability step ran.
	for _, name := range []string{"request.b3b-r07.jsonl", "completion.b3b-r07.jsonl"} {
		folder := "requests"
		if strings.HasPrefix(name, "completion") {
			folder = "completions"
		}
		if _, statErr := os.Stat(filepath.Join(root, folder, name)); statErr != nil {
			t.Fatalf("paired mode must publish %s: %v", name, statErr)
		}
	}
}
