// Command agent-stub is the deterministic offline CLI used by the ProofRail
// AgentRunner experiments. It never touches the network and never writes outside
// the evidence directory it is handed, so a run can be replayed byte for byte.
//
// The launcher tells it where to publish artifacts through
// PROOFRAIL_EVIDENCE_DIR and redirects its stdout and stderr into its own log
// artifacts, so the log completeness marker must be the last line of both streams.
//
// Usage:
//
//	agent-stub -version
//	agent-stub -mode run -sleep 2s
//	agent-stub -mode spawn -child-life 1h -sleep 30s
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// stubVersion is the version the launcher must pin. The precheck compares the
// reported string literally, so it is a contract rather than a label.
const stubVersion = "agent-stub-0.1.0"

// logCompletenessMarker mirrors the adapter constant: the marker must terminate
// each log stream, otherwise the run cannot prove the log was written fully.
const logCompletenessMarker = "proofrail-run-complete"

// Artifact names the adapter collects from the evidence directory.
const (
	usageFileName  = "usage.json"
	eventsFileName = "events.jsonl"
	ownPIDFileName = "stub.pid"
	childPIDName   = "stub-child.pid"
	evidenceDirEnv = "PROOFRAIL_EVIDENCE_DIR"
	// mutatedWorkspaceFileName is what the mutate-workspace mode writes into the
	// workspace it was given, so a manifest captured afterwards must show it.
	mutatedWorkspaceFileName = "agent-stub-mutation.txt"
)

// options is one deterministic run: the same flags always produce the same
// artifacts, so an experiment can compare outcomes instead of guessing.
type options struct {
	dir        string
	mode       string
	sleep      time.Duration
	childLife  time.Duration
	exitCode   int
	calls      int
	tokens     int
	durationMs int
	noMarker   bool
	tornUsage  bool
	noEvents   bool
}

func main() {
	mode := flag.String("mode", "run", "run|spawn|spawn-exit|spawn-wait|crash|no-usage|ask|hang|fail-fast|mutate-workspace|torn-events|child")
	sleep := flag.Duration("sleep", time.Minute, "how long the run stays alive before exiting")
	childLife := flag.Duration("child-life", time.Hour, "how long a spawned child stays alive")
	exitCode := flag.Int("exit-code", 3, "exit status for the crash and fail-fast modes")
	calls := flag.Int("calls", 2, "value written to usage.json")
	tokens := flag.Int("tokens", 123, "value written to usage.json")
	durationMs := flag.Int("duration-ms", 5, "value written to usage.json")
	noMarker := flag.Bool("no-marker", false, "omit the log completeness marker, as an interrupted writer would")
	tornUsage := flag.Bool("torn-usage", false, "write an unparsable usage artifact")
	noEvents := flag.Bool("no-events", false, "omit the event log")
	showVersion := flag.Bool("version", false, "print the pinned version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, stubVersion)
		os.Exit(0)
	}

	os.Exit(options{
		dir:        evidenceDir(),
		mode:       *mode,
		sleep:      *sleep,
		childLife:  *childLife,
		exitCode:   *exitCode,
		calls:      *calls,
		tokens:     *tokens,
		durationMs: *durationMs,
		noMarker:   *noMarker,
		tornUsage:  *tornUsage,
		noEvents:   *noEvents,
	}.execute())
}

// execute performs one run and returns the process exit code.
func (run options) execute() int {
	switch run.mode {
	case "child":
		// A spawned descendant: it only has to outlive its parent.
		run.writeOwnPID()
		time.Sleep(run.childLife)
		return 0
	case "fail-fast":
		// The run dies before it produces anything: the start phase fails.
		fmt.Fprintln(os.Stderr, "the run never started")
		return run.exitCode
	case "crash":
		// A crash mid-write leaves the launcher-owned log without its marker.
		fmt.Fprintln(os.Stdout, "interrupted")
		return run.exitCode
	case "hang":
		// No artifacts beyond its own pid: the run is stopped from the outside.
		run.writeOwnPID()
		time.Sleep(time.Hour)
		return 0
	case "run", "spawn", "spawn-exit", "spawn-wait", "no-usage", "ask", "mutate-workspace", "torn-events":
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q\n", run.mode)
		return 2
	}

	if run.mode == "mutate-workspace" {
		// A run that changes its workspace: the pre and post manifests must differ,
		// so a diff that claims "no change" is a broken artifact. The working
		// directory is the workspace the launcher handed over.
		workspace, err := os.Getwd()
		if err != nil {
			return 13
		}
		if err := os.WriteFile(filepath.Join(workspace, mutatedWorkspaceFileName), []byte("the run wrote this\n"), 0o600); err != nil {
			return 13
		}
	}
	run.writeOwnPID()
	run.emitArtifacts()
	run.emitEvents()
	if run.mode == "ask" {
		run.appendOperatorAsk()
	}
	run.emitMarker()
	if run.spawnsChild() {
		if code := run.spawnChild(); code != 0 {
			return code
		}
		if run.mode == "spawn-exit" {
			// The parent leaves cleanly while its child keeps running: apart from the
			// tree, this run looks like a success.
			return 0
		}
	}
	time.Sleep(run.sleep)
	return 0
}

// spawnsChild reports whether the mode leaves a descendant behind or waits for it.
func (run options) spawnsChild() bool {
	return run.mode == "spawn" || run.mode == "spawn-exit" || run.mode == "spawn-wait"
}

// spawnChild starts a real descendant that outlives its parent and reports its pid
// in the artifact the run reads. Only the mode that models a well-behaved CLI waits
// for it. The child's stdout and stderr are discarded: the parent owns the logs.
func (run options) spawnChild() int {
	self, err := os.Executable()
	if err != nil {
		return 4
	}
	command := exec.Command(self, "-mode", "child", "-child-life", run.childLife.String())
	command.Env = append(os.Environ(), evidenceDirEnv+"="+run.dir)
	command.Dir = run.dir
	if err := command.Start(); err != nil {
		return 5
	}
	if err := os.WriteFile(filepath.Join(run.dir, childPIDName), fmt.Appendf(nil, "%d\n", command.Process.Pid), 0o600); err != nil {
		return 6
	}
	if run.mode == "spawn-wait" {
		_, _ = command.Process.Wait()
	}
	return 0
}

// emitArtifacts writes the usage artifact of a clean run, or the torn one that
// models a writer interrupted half-way. The no-usage mode omits it entirely.
func (run options) emitArtifacts() {
	if run.mode == "no-usage" {
		return
	}
	body := fmt.Sprintf("{\"calls\":%d,\"tokens\":%d,\"durationMs\":%d}", run.calls, run.tokens, run.durationMs)
	if run.tornUsage {
		body = fmt.Sprintf("{\"calls\":%d,\"tokens\":", run.calls)
	}
	if err := os.WriteFile(filepath.Join(run.dir, usageFileName), []byte(body), 0o600); err != nil {
		os.Exit(7)
	}
}

// emitEvents writes the event log, which the run reads line by line.
func (run options) emitEvents() {
	if run.noEvents || run.mode == "no-usage" {
		return
	}
	body := "{\"type\":\"session-created\"}\n{\"type\":\"output\"}\n"
	if run.mode == "torn-events" {
		// A writer interrupted half-way: the last line is not decodable, so the file
		// proves nothing about what came after it.
		body = "{\"type\":\"session-created\"}\n{\"type\":\"usa"
	}
	if err := os.WriteFile(filepath.Join(run.dir, eventsFileName), []byte(body), 0o600); err != nil {
		os.Exit(8)
	}
}

// appendOperatorAsk turns the run into one that needs a human: the outcome mapping
// must never report it as completed.
func (run options) appendOperatorAsk() {
	file, err := os.OpenFile(filepath.Join(run.dir, eventsFileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		os.Exit(9)
	}
	defer func() { _ = file.Close() }()
	if _, err := file.WriteString("{\"type\":\"operator-action-required\"}\n"); err != nil {
		os.Exit(10)
	}
}

// emitMarker terminates both log streams, which is what makes the collected log
// provably complete.
func (run options) emitMarker() {
	if run.noMarker {
		return
	}
	fmt.Fprintln(os.Stdout, logCompletenessMarker)
	fmt.Fprintln(os.Stderr, logCompletenessMarker)
}

// writeOwnPID publishes the run's own pid so an experiment can prove it is gone
// after a stop.
func (run options) writeOwnPID() {
	if err := os.WriteFile(filepath.Join(run.dir, ownPIDFileName), fmt.Appendf(nil, "%d\n", os.Getpid()), 0o600); err != nil {
		os.Exit(11)
	}
}

// evidenceDir is where the launcher told the CLI to publish its artifacts. A run
// without that variable falls back to the working directory, never to a shared
// location.
func evidenceDir() string {
	if dir := strings.TrimSpace(os.Getenv(evidenceDirEnv)); dir != "" {
		return dir
	}
	dir, err := os.Getwd()
	if err != nil {
		os.Exit(12)
	}
	return dir
}
