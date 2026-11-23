package adapters

// P85 T8a (finding F8) - the BYOK confinement proof gate: counterexamples and the
// mutation hooks that keep them falsifiable.
//
// The gate has two halves that share one proof. The checker half runs before any
// credential is resolved; the executor half runs again immediately before a
// credential-bearing child would be created. Both are fail-closed: a missing
// capability, an unreadable verdict and a spawn configuration that breaks one of
// the five invariants are all refusals, never admissions.
//
// Mutation checks (each one turns a named test red):
//   - delete the checker's BYOK gate block in ai_preflight.go =>
//     TestP85T8AMissingProofRefusesBeforeResolveOrSpawn and
//     TestP85T8AUnreadableSandboxVerdictIsNotAdmission go red (the credential is
//     resolved and the probe is called);
//   - delete the `if credentialPresent { assert... }` block in Execute =>
//     TestP85T8AExecuteRefusesBrokenBYOKSpawnWithoutSpawning goes red (the marker
//     child runs);
//   - delete one element assertion of assertCopilotCLIProbeBYOKConfinement =>
//     the matching subtest of TestP85T8AConfinementProofRefusesEachBrokenElement
//     goes red (the error falls through to the sandbox verdict).

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

// p85T8ACountingResolver records how many times a credential was read. The gate's
// whole point is that this stays zero on every refusal.
type p85T8ACountingResolver struct {
	calls  int
	secret []byte
}

func (resolver *p85T8ACountingResolver) Resolve(context.Context, string) ([]byte, error) {
	resolver.calls++
	return slices.Clone(resolver.secret), nil
}

// p85T8AProbeWithoutProof is a probe that does NOT expose the BYOK confinement
// capability. T8a makes that a refusal, never a downgrade to "assume confined".
type p85T8AProbeWithoutProof struct {
	calls int
}

func (probe *p85T8AProbeWithoutProof) Probe(_ context.Context, _ AIProviderProfile, _ []byte, _ string, _ int) (AIProbeResult, error) {
	probe.calls++
	return AIProbeResult{Status: "available", RequestsUsed: 1, Evidence: []string{queueHashTwo}}, nil
}

// p85T8AProbeWithFailingProof exposes the capability but the proof fails, which is
// the "the capability answered, the answer was no" case.
type p85T8AProbeWithFailingProof struct {
	probe *p85T8AProbeWithoutProof
	err   error
}

func (probe *p85T8AProbeWithFailingProof) Probe(ctx context.Context, profile AIProviderProfile, credential []byte, channel string, maximumRequests int) (AIProbeResult, error) {
	return probe.probe.Probe(ctx, profile, credential, channel, maximumRequests)
}

func (probe *p85T8AProbeWithFailingProof) VerifyBYOKConfinement() error { return probe.err }

// p85T8ASandboxProofSpecDigest is the specification digest the checker-level
// doubles are asked about. It only has to be stable and non-empty: guard.CheckSandbox
// compares the verdict's ECHO with this value, so a double that answers about
// another specification is refused exactly like the platform layer would be.
const p85T8ASandboxProofSpecDigest = "sha256:0000000000000000000000000000000000000000000000000000000000000000"

// p85T8ASandboxDouble is the interface stand-in of the sandbox check layer: it
// answers one scripted verdict and records the spec it was handed, so a test can
// drive the verdict instead of the host environment.
//
// P85 F13 (design rev 3, section 5): the double ECHOES the specification digest it
// was asked about, which is what a conforming check layer does. A double that
// stayed silent would be refused for carrying no echo instead of for the scripted
// verdict, and the verdict cases below could no longer be told apart.
type p85T8ASandboxDouble struct {
	verdict guard.Verdict
	err     error
	specs   []guard.CheckSpec
}

func (double *p85T8ASandboxDouble) Check(_ context.Context, spec guard.CheckSpec) (guard.Verdict, error) {
	double.specs = append(double.specs, spec)
	verdict := double.verdict
	verdict.SpecDigest = spec.SpecDigest
	return verdict, double.err
}

// p85T8AValueEchoProbe is a probe that DOES prove its confinement and delegates
// the run to the hostile executor below, so the refusal-evidence path is exercised
// without pinning a real executable's install source.
type p85T8AValueEchoProbe struct {
	inner CopilotCLIProbe
}

func (probe *p85T8AValueEchoProbe) Probe(ctx context.Context, profile AIProviderProfile, credential []byte, channel string, maximumRequests int) (AIProbeResult, error) {
	return probe.inner.Probe(ctx, profile, credential, channel, maximumRequests)
}

func (probe *p85T8AValueEchoProbe) VerifyBYOKConfinement() error { return nil }

// p85T8ASandboxProofProbe composes the production shape: the probe's proof IS the
// sandbox gate (guard.CheckSandbox) over a supplied check layer. It lets the
// checker-level counterexamples judge an Unknown verdict without depending on the
// host's platform sandbox.
type p85T8ASandboxProofProbe struct {
	sandbox guard.Sandbox
	probe   *p85T8AProbeWithoutProof
}

func (probe *p85T8ASandboxProofProbe) Probe(ctx context.Context, profile AIProviderProfile, credential []byte, channel string, maximumRequests int) (AIProbeResult, error) {
	return probe.probe.Probe(ctx, profile, credential, channel, maximumRequests)
}

func (probe *p85T8ASandboxProofProbe) VerifyBYOKConfinement() error {
	// P85 F13: a BYOK check that names no specification is refused before the
	// layer is even consulted, so the probe asks about the fixture's stable digest
	// and the double echoes it back.
	return guard.CheckSandbox(context.Background(), probe.sandbox, guard.CheckSpec{Reason: guard.SandboxReasonBYOKProbe, SpecDigest: p85T8ASandboxProofSpecDigest})
}

// p85T8AValueEchoExecutor is a hostile executor: it answers "start-failed" under
// the BYOK class AND tries to smuggle a refused value (the model) into its
// Evidence and Measurements. T8b discipline says neither may reach the record.
type p85T8AValueEchoExecutor struct {
	model string
}

func (executor p85T8AValueEchoExecutor) VerifyBYOKConfinement() error { return nil }

func (executor p85T8AValueEchoExecutor) Execute(_ context.Context, invocation CopilotCLIProbeInvocation) CopilotCLIProbeExecution {
	leak := evidence.Digest("proofrail:test-echo:1\n", []byte(executor.model+"\n"+strings.Join(invocation.Args, "\n")))
	return CopilotCLIProbeExecution{
		Outcome:      "start-failed",
		FailureCode:  copilotCLIProbeBYOKFailure,
		Evidence:     []string{leak},
		Measurements: []string{leak},
	}
}

func p85T8AChecker(probe AIChannelProbe, resolver AISecretResolver) *AIAvailabilityChecker {
	return availabilityChecker(probe, resolver)
}

// TestP85T8AMissingProofRefusesBeforeResolveOrSpawn is the "proof missing"
// counterexample: no capability and a failing capability are both refusals, and in
// both cases the secret is never read, the probe (and therefore the child) is
// never called, and the record stays a valid, complete availability record.
func TestP85T8AMissingProofRefusesBeforeResolveOrSpawn(t *testing.T) {
	silent := &p85T8AProbeWithoutProof{}
	failing := &p85T8AProbeWithFailingProof{
		probe: &p85T8AProbeWithoutProof{},
		err:   errors.New("p85-t8a: the proof value must not reach the record"),
	}
	cases := map[string]AIChannelProbe{
		"capability-missing": silent,
		"capability-refused": failing,
	}
	for name, probe := range cases {
		t.Run(name, func(t *testing.T) {
			resolver := &p85T8ACountingResolver{secret: []byte("p85-t8a-secret")}
			checker := p85T8AChecker(probe, resolver)
			record, err := checker.Check(context.Background(), deepSeekProfile(), "agent-runner-cli", 1)
			if err != nil {
				t.Fatal(err)
			}
			availability := record.Availability
			if availability.Status != "unavailable" || availability.Reason == nil || *availability.Reason != copilotCLIProbeBYOKFailure {
				t.Fatalf("the refusal must be the closed-set BYOK class, got %+v", availability)
			}
			if availability.RequestsUsed != 0 {
				t.Fatalf("a refused probe must spend no request, got %d", availability.RequestsUsed)
			}
			// "Not published": the record is still a complete, valid record.
			if err := ValidateAIAvailabilityRecord(record); err != nil {
				t.Fatalf("the refusal record must stay complete and valid: %v", err)
			}
			// "Not advanced": no credential was resolved and no probe ran.
			if resolver.calls != 0 {
				t.Fatalf("the credential must not be resolved before the proof, got %d resolves", resolver.calls)
			}
			if silent.calls != 0 || failing.probe.calls != 0 {
				t.Fatalf("no probe may run before the proof, got %d/%d calls", silent.calls, failing.probe.calls)
			}
			joined := strings.Join(availability.Evidence, "\n")
			if strings.Contains(joined, "must not reach the record") {
				t.Fatalf("the refusal must not echo the proof error, got %q", joined)
			}
		})
	}
}

// TestP85T8AUnreadableSandboxVerdictIsNotAdmission is the "proof unreadable"
// counterexample at both layers: the guard gate refuses Unknown and Denied (only
// Allowed passes), and the checker propagates that refusal without reading the
// credential or calling the probe - Unknown is never read as a pass.
func TestP85T8AUnreadableSandboxVerdictIsNotAdmission(t *testing.T) {
	t.Run("guard-gate-is-fail-closed", func(t *testing.T) {
		cases := []struct {
			name    string
			verdict guard.Verdict
			wantErr bool
		}{
			{"unknown", guard.Verdict{Status: guard.VerdictUnknown, Reason: "S1 unknown"}, true},
			{"denied", guard.Verdict{Status: guard.VerdictDenied, Reason: "S1 denied"}, true},
			{"allowed", guard.Verdict{Status: guard.VerdictAllowed}, false},
		}
		for _, testCase := range cases {
			t.Run(testCase.name, func(t *testing.T) {
				double := &p85T8ASandboxDouble{verdict: testCase.verdict}
				err := guard.CheckSandbox(context.Background(), double, guard.CheckSpec{Reason: guard.SandboxReasonBYOKProbe, SpecDigest: p85T8ASandboxProofSpecDigest})
				if (err != nil) != testCase.wantErr {
					t.Fatalf("verdict %s gave err=%v, want error=%v", testCase.verdict.Status, err, testCase.wantErr)
				}
				if len(double.specs) != 1 || double.specs[0].Reason != guard.SandboxReasonBYOKProbe {
					t.Fatalf("the gate must label the check with the BYOK reason, got %+v", double.specs)
				}
			})
		}
	})

	t.Run("checker-refuses-unknown-verdict", func(t *testing.T) {
		unknown := &p85T8ASandboxDouble{verdict: guard.Verdict{Status: guard.VerdictUnknown, Reason: "S1 sandbox enforcement checks not implemented"}}
		probe := &p85T8ASandboxProofProbe{sandbox: unknown, probe: &p85T8AProbeWithoutProof{}}
		resolver := &p85T8ACountingResolver{secret: []byte("p85-t8a-secret")}
		record, err := p85T8AChecker(probe, resolver).Check(context.Background(), deepSeekProfile(), "agent-runner-cli", 1)
		if err != nil {
			t.Fatal(err)
		}
		availability := record.Availability
		if availability.Status != "unavailable" || availability.Reason == nil || *availability.Reason != copilotCLIProbeBYOKFailure {
			t.Fatalf("an unknown verdict must refuse as the BYOK class, got %+v", availability)
		}
		if resolver.calls != 0 || probe.probe.calls != 0 {
			t.Fatalf("an unknown verdict must not resolve the credential or run the probe, got resolves=%d calls=%d", resolver.calls, probe.probe.calls)
		}
		if err := ValidateAIAvailabilityRecord(record); err != nil {
			t.Fatalf("the refusal record must stay complete and valid: %v", err)
		}
	})

	t.Run("checker-admits-allowed-verdict", func(t *testing.T) {
		// The negative control: the same composition admits when the verdict is
		// Allowed, so the refusals above are not vacuous.
		allowed := &p85T8ASandboxDouble{verdict: guard.Verdict{Status: guard.VerdictAllowed}}
		probe := &p85T8ASandboxProofProbe{sandbox: allowed, probe: &p85T8AProbeWithoutProof{}}
		resolver := &p85T8ACountingResolver{secret: []byte("p85-t8a-secret")}
		record, err := p85T8AChecker(probe, resolver).Check(context.Background(), deepSeekProfile(), "agent-runner-cli", 1)
		if err != nil {
			t.Fatal(err)
		}
		if record.Availability.Status != "available" || probe.probe.calls != 1 {
			t.Fatalf("an allowed verdict must admit the probe, got %+v calls=%d", record.Availability, probe.probe.calls)
		}
	})
}

// p85T8ASpawnEnvelope is one well-formed BYOK spawn configuration. The subtests of
// TestP85T8AConfinementProofRefusesEachBrokenElement mutate exactly one element of
// it and require that element's refusal, which is what makes the mutation check
// per-element rather than all-or-nothing.
type p85T8ASpawnEnvelope struct {
	invocation    CopilotCLIProbeInvocation
	spec          guard.ProcessSpec
	temporaryRoot string
	usagePath     string
	outputLimit   int64
}

func p85T8ANewSpawnEnvelope(t *testing.T) p85T8ASpawnEnvelope {
	t.Helper()
	root := t.TempDir()
	usagePath := filepath.Join(root, "usage.json")
	invocation := CopilotCLIProbeInvocation{
		Executable:       filepath.Join(root, "copilot"),
		WorkingDirectory: root,
		Args:             []string{"-p", copilotCLIProbePrompt, "--model", "p85-t8a-model"},
		Environment:      copilotCLIProbeEnvironment(deepSeekProfile(), []byte("p85-t8a-secret")),
	}
	return p85T8ASpawnEnvelope{
		invocation: invocation,
		spec: guard.ProcessSpec{
			Command: invocation.Executable,
			Args:    copilotCLIProbeExecArgs(invocation.Args, usagePath),
			Dir:     root,
			Env: []string{
				"COPILOT_MODEL=p85-t8a-model",
				"COPILOT_PROVIDER_API_KEY=p85-t8a-secret",
			},
			Stdout: &bytes.Buffer{},
			Stderr: &bytes.Buffer{},
		},
		temporaryRoot: root,
		usagePath:     usagePath,
		outputLimit:   defaultCopilotCLIOutputBytes,
	}
}

func (envelope p85T8ASpawnEnvelope) assert(t *testing.T) error {
	t.Helper()
	return assertCopilotCLIProbeBYOKConfinement(context.Background(), envelope.invocation, envelope.spec, envelope.temporaryRoot, envelope.usagePath, envelope.outputLimit)
}

// TestP85T8AConfinementProofRefusesEachBrokenElement is the "well-formed but
// mismatched" counterexample at the spawn-time half of the proof: argv, the usage
// artifact, the environment, the stream set, the working directory and the output
// bound are each broken in turn.
//
// Each case asserts the refusal names its ELEMENT and not the platform verdict
// ("sandbox"): if the element assertion were deleted, the call would fall through
// to the sandbox verdict and the case goes red.
func TestP85T8AConfinementProofRefusesEachBrokenElement(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(t *testing.T, envelope *p85T8ASpawnEnvelope)
		want   string
	}{
		{
			name: "argv-names-a-logging-destination",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				envelope.spec.Args = append(slices.Clone(envelope.spec.Args), "--log-dir", "p85-t8a-logs")
			},
			want: "argv names a logging destination",
		},
		{
			name: "argv-is-not-the-screened-argv",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				envelope.spec.Args = slices.Clone(envelope.invocation.Args)
			},
			want: "the spawn argv is not the screened probe argv",
		},
		{
			name: "usage-artifact-outside-the-executor-directory",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				outside := filepath.Join(t.TempDir(), "usage.json")
				envelope.usagePath = outside
				envelope.spec.Args = copilotCLIProbeExecArgs(envelope.invocation.Args, outside)
			},
			want: "the usage artifact is outside the executor directory",
		},
		{
			name: "child-environment-carries-a-foreign-key",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				envelope.spec.Env = append(slices.Clone(envelope.spec.Env), "P85_T8A_FOREIGN=value")
			},
			want: "the child environment carries a key outside the allowlist",
		},
		{
			name: "byok-environment-carries-a-foreign-key",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				envelope.invocation.Environment = copilotCLIProbeEnvironment(deepSeekProfile(), []byte("p85-t8a-secret"))
				envelope.invocation.Environment["P85_T8A_FOREIGN"] = []byte("value")
			},
			want: "the BYOK environment carries a key outside the profile set",
		},
		{
			name: "streams-not-bound",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				envelope.spec.Stdout = nil
			},
			want: "the child streams are not bound",
		},
		{
			name: "stdout-is-a-file-handle",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				file, err := os.CreateTemp(t.TempDir(), "p85-t8a-stdout-*")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = file.Close() })
				envelope.spec.Stdout = file
			},
			want: "stdout is bound to a file handle",
		},
		{
			name: "working-directory-is-not-the-executor-directory",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				envelope.spec.Dir = t.TempDir()
			},
			want: "the child working directory is not the executor directory",
		},
		{
			name: "output-capture-is-unbounded",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				envelope.outputLimit = 0
			},
			want: "the in-memory output capture is unbounded",
		},
		{
			// F02 (P85 fix round 1): the screened argv and the spawn argv BOTH name
			// the flag, so element 1 (argv equality) passes and the refusal can only
			// come from the argv working-directory element.
			name: "argv-working-directory-outside-the-executor-directory",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				outside := t.TempDir()
				envelope.invocation.Args = append(slices.Clone(envelope.invocation.Args), "-C", outside)
				envelope.spec.Args = copilotCLIProbeExecArgs(envelope.invocation.Args, envelope.usagePath)
			},
			want: "the argv names a working directory outside the executor directory",
		},
		{
			// F02 (P85 fix round 1): the inline `--working-directory=<dir>` form
			// carries its value inside the same token, so it must be screened too.
			name: "argv-inline-working-directory-outside-the-executor-directory",
			mutate: func(t *testing.T, envelope *p85T8ASpawnEnvelope) {
				outside := t.TempDir()
				envelope.invocation.Args = append(slices.Clone(envelope.invocation.Args), "--working-directory="+outside)
				envelope.spec.Args = copilotCLIProbeExecArgs(envelope.invocation.Args, envelope.usagePath)
			},
			want: "the argv names a working directory outside the executor directory",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			envelope := p85T8ANewSpawnEnvelope(t)
			testCase.mutate(t, &envelope)
			err := envelope.assert(t)
			if err == nil {
				t.Fatalf("a spawn configuration that breaks %q must be refused", testCase.name)
			}
			if !errors.Is(err, ErrAIBYOKConfinement) {
				t.Fatalf("the refusal must carry ErrAIBYOKConfinement, got %v", err)
			}
			if !strings.Contains(err.Error(), testCase.want) {
				t.Fatalf("the refusal must name the broken element %q, got %v", testCase.want, err)
			}
			if strings.Contains(err.Error(), "sandbox") {
				t.Fatalf("the refusal must come from element %q, not from the platform verdict, got %v", testCase.name, err)
			}
		})
	}
}

// TestP85T8AConfinementProofAdmitsWorkingDirectoryNamingTheExecutor pins the P85 F13
// replacement of the F02 positive control. The F02 premise - "a `-C` whose value
// cleans to the executor directory is admitted" - no longer holds for a
// credential-bearing spawn: F13 freezes the secret argv WITHOUT any working-directory
// flag, and the freeze refuses a secret argv that names one. What survives is the F02
// assertion as depth-in-depth defence for the host-auth path and for callers that
// bypass the template.
//
// The two halves that replace the old premise:
//   - freeze half: a secret argv carrying `-C` (naming anything, even the executor
//     directory) is refused by the allowlist, so the frozen template cannot carry one;
//   - positive half: the constructor's secret template is admitted, and the spawned
//     ProcessSpec.Dir is the only working directory the child receives.
func TestP85T8AConfinementProofAdmitsWorkingDirectoryNamingTheExecutor(t *testing.T) {
	// Freeze half: the same argv the constructor emits, plus a `-C` naming the
	// executor directory, must not pass the freeze.
	executorDirectory := t.TempDir()
	named := append(slices.Clone(copilotCLIProbeArgs(deepSeekProfile().Model, executorDirectory, 1, "secret")), "-C", executorDirectory)
	if err := copilotCLIProbeValidateWhitelistedArgs(named, "secret"); !errors.Is(err, ErrAIBYOKConfinement) {
		t.Fatalf("a credential-bearing argv that names a working directory must be refused at freeze time, got %v", err)
	}

	// Spawn half: the same shape, with `-C` naming the executor directory EXACTLY
	// (dressed as "<root>/." so cleaning is what decides), is still refused - and the
	// refusal names the working directory instead of falling through to the platform
	// verdict.
	envelope := p85T8ANewSpawnEnvelope(t)
	dressed := filepath.Clean(envelope.temporaryRoot) + string(os.PathSeparator) + "."
	envelope.invocation.Args = append(slices.Clone(envelope.invocation.Args), "-C", dressed)
	envelope.spec.Args = copilotCLIProbeExecArgs(envelope.invocation.Args, envelope.usagePath)
	err := envelope.assert(t)
	if err == nil || !errors.Is(err, ErrAIBYOKConfinement) || !strings.Contains(err.Error(), "working directory") {
		t.Fatalf("a credential-bearing spawn argv must be refused for naming a working directory, got %v", err)
	}
	if strings.Contains(err.Error(), "sandbox") {
		t.Fatalf("the refusal must come from the argv element, not the platform verdict, got %v", err)
	}

	// Positive half: the constructor's secret template carries no working-directory
	// flag, it is admitted, and the child working directory is the executor directory.
	positive := p85T8ANewSpawnEnvelope(t)
	positive.invocation.Args = copilotCLIProbeArgs(deepSeekProfile().Model, positive.temporaryRoot, 1, "secret")
	positive.spec.Args = copilotCLIProbeExecArgs(positive.invocation.Args, positive.usagePath)
	if _, templateErr := copilotCLIProbeLaunchSpec(positive.invocation.Executable, positive.temporaryRoot, "secret", deepSeekProfile().Model, 1, positive.outputLimit); templateErr != nil {
		t.Fatalf("the secret launch template must be constructible: %v", templateErr)
	}
	for _, flag := range []string{"-C", "--working-directory"} {
		if slices.Contains(positive.spec.Args, flag) {
			t.Fatalf("the secret spawn argv must not carry %q: %#v", flag, positive.spec.Args)
		}
	}
	if err := positive.assert(t); err != nil && !strings.Contains(err.Error(), "sandbox") {
		t.Fatalf("only the platform verdict may refuse the well-formed envelope, got %v", err)
	}
	if filepath.Clean(positive.spec.Dir) != filepath.Clean(positive.temporaryRoot) {
		t.Fatalf("the only working directory a BYOK child may receive is ProcessSpec.Dir, got %q", positive.spec.Dir)
	}
}

// TestP85T8AProcessSpecCarriesNoExtraInheritedHandle is the interface half of
// element 3: "ExtraFiles is empty" cannot be constructed through ProcessSpec
// because the spec has no field that can carry an extra inherited handle, and the
// only writer fields are the two standard streams. If such a field is ever added,
// this goes red and the proof has to grow a matching assertion.
func TestP85T8AProcessSpecCarriesNoExtraInheritedHandle(t *testing.T) {
	specType := reflect.TypeOf(guard.ProcessSpec{})
	writerType := reflect.TypeOf((*io.Writer)(nil)).Elem()
	writers := make([]string, 0, 2)
	for index := 0; index < specType.NumField(); index++ {
		field := specType.Field(index)
		switch field.Name {
		case "ExtraFiles", "Files", "SysProcAttr":
			t.Fatalf("guard.ProcessSpec gained an inherited-handle field %q; the BYOK proof must assert it empty", field.Name)
		}
		if field.Type.Implements(writerType) {
			writers = append(writers, field.Name)
		}
	}
	if !slices.Equal(writers, []string{"Stdout", "Stderr"}) {
		t.Fatalf("the only stream handles a spawn spec may bind are stdout and stderr, got %v", writers)
	}
}

// p85T8ASystemShell names a shell whose installation source this host can prove
// immutable, because Execute re-verifies the source before it reaches the BYOK
// assertion (T3). A host that cannot prove it skips instead of failing the BYOK
// counterexample for an unrelated reason.
func p85T8ASystemShell(t *testing.T) string {
	t.Helper()
	var shell string
	switch runtime.GOOS {
	case "windows":
		root := os.Getenv("SystemRoot")
		if root == "" {
			root = `C:\Windows`
		}
		shell = filepath.Join(root, "System32", "cmd.exe")
	case "linux":
		shell = "/bin/sh"
	default:
		t.Skipf("the install source proof has a platform rule only for windows and linux, not %s", runtime.GOOS)
	}
	if _, err := VerifyAIInstallSource(shell); err != nil {
		t.Skipf("the shell %s is not provably immutable for the executing user on this host (%v)", shell, err)
	}
	return shell
}

// p85T8AMarkerArgs makes the child create marker when it runs, which is the
// observable form of "the child was created".
func p85T8AMarkerArgs(marker string) []string {
	if runtime.GOOS == "windows" {
		return []string{"/c", `echo ran> "` + marker + `"`}
	}
	return []string{"-c", "echo ran > '" + marker + "'"}
}

// TestP85T8AExecuteRefusesBrokenBYOKSpawnWithoutSpawning drives the whole Execute
// path with a credential-bearing invocation whose environment carries one foreign
// key. The pre-spawn proof refuses, the child never runs, and the refusal is the
// closed class - not a value-carrying description.
func TestP85T8AExecuteRefusesBrokenBYOKSpawnWithoutSpawning(t *testing.T) {
	shell := p85T8ASystemShell(t)
	root := t.TempDir()
	marker := filepath.Join(root, "p85-t8a-spawned.marker")
	environment := copilotCLIProbeEnvironment(deepSeekProfile(), []byte("p85-t8a-secret"))
	environment["P85_T8A_FOREIGN"] = []byte("value")

	executor := ManagedCopilotCLIProbeExecutor{TemporaryDirectory: root}
	execution := executor.Execute(context.Background(), CopilotCLIProbeInvocation{
		Executable:       shell,
		WorkingDirectory: root,
		Args:             p85T8AMarkerArgs(marker),
		Environment:      environment,
	})
	if execution.Outcome != "start-failed" || execution.FailureCode != copilotCLIProbeBYOKFailure {
		t.Fatalf("a broken BYOK spawn envelope must refuse under the BYOK class, got outcome=%q failure=%q evidence=%v",
			execution.Outcome, execution.FailureCode, execution.Evidence)
	}
	want := copilotCLIExecutorEvidence("byok-confinement-unproven")
	if len(execution.Evidence) != 1 || execution.Evidence[0] != want {
		t.Fatalf("the refusal must be the stable classification %q, got %v", want, execution.Evidence)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("the credential-bearing child must not be started, marker stat err=%v", err)
	}
}

// TestP85T8AExecuteRefusesUnprovableVerdictWithoutSpawning is the end-to-end form
// of the "unreadable proof" counterexample: on a host whose platform verdict is
// not Allowed (every S1 host), a well-formed BYOK spawn still refuses and no child
// exists. A host that proves the environment skips, because the unlock path
// (a platform-level write restriction) is exactly what would legitimately admit.
func TestP85T8AExecuteRefusesUnprovableVerdictWithoutSpawning(t *testing.T) {
	if err := guard.CheckSandbox(context.Background(), nil, guard.CheckSpec{Reason: guard.SandboxReasonBYOKProbe, SpecDigest: p85T8ASandboxProofSpecDigest}); err == nil {
		t.Skip("this host proves the sandbox verdict, so the S1 refusal no longer applies")
	}
	shell := p85T8ASystemShell(t)
	root := t.TempDir()
	marker := filepath.Join(root, "p85-t8a-spawned.marker")
	executor := ManagedCopilotCLIProbeExecutor{TemporaryDirectory: root}
	execution := executor.Execute(context.Background(), CopilotCLIProbeInvocation{
		Executable:       shell,
		WorkingDirectory: root,
		Args:             p85T8AMarkerArgs(marker),
		Environment:      copilotCLIProbeEnvironment(deepSeekProfile(), []byte("p85-t8a-secret")),
	})
	if execution.Outcome != "start-failed" || execution.FailureCode != copilotCLIProbeBYOKFailure {
		t.Fatalf("an unprovable platform verdict must refuse under the BYOK class, got outcome=%q failure=%q",
			execution.Outcome, execution.FailureCode)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("an unprovable verdict must not start the child, marker stat err=%v", err)
	}
}

// TestP85T8ASameClassRefusalsShareEvidenceAcrossValues pins the T8b compatibility
// of the new class: two refusals that differ only in the refused VALUE (two
// different models) produce byte-identical evidence, so the new classification
// cannot become a channel that echoes a refused value.
func TestP85T8ASameClassRefusalsShareEvidenceAcrossValues(t *testing.T) {
	const (
		alphaModel = "p85-t8a-model-alpha"
		betaModel  = "p85-t8a-model-beta"
	)
	runOnce := func(model string) (string, []string) {
		profile := deepSeekProfile()
		profile.Model = model
		probe := &p85T8AValueEchoProbe{inner: CopilotCLIProbe{
			Executable:       "p85-t8a-copilot",
			WorkingDirectory: t.TempDir(),
			Executor:         p85T8AValueEchoExecutor{model: model},
		}}
		resolver := &p85T8ACountingResolver{secret: []byte("p85-t8a-secret")}
		record, err := p85T8AChecker(probe, resolver).Check(context.Background(), profile, "agent-runner-cli", 1)
		if err != nil {
			t.Fatal(err)
		}
		availability := record.Availability
		if availability.Status != "unavailable" || availability.Reason == nil || *availability.Reason != copilotCLIProbeBYOKFailure {
			t.Fatalf("model %s: the refusal must be the BYOK class, got %+v", model, availability)
		}
		return model, availability.Evidence
	}
	alphaName, alphaEvidence := runOnce(alphaModel)
	betaName, betaEvidence := runOnce(betaModel)

	if !slices.Equal(alphaEvidence, betaEvidence) {
		t.Fatalf("refusals of the same class must share evidence, got %v and %v", alphaEvidence, betaEvidence)
	}
	joined := strings.Join(alphaEvidence, "\n")
	for _, value := range []string{alphaName, betaName} {
		if strings.Contains(joined, value) {
			t.Fatalf("the refusal evidence must not echo the refused value %q, got %q", value, joined)
		}
	}
	// The evidence is the closed classification and nothing the executor observed.
	want := copilotCLIProbeRefusalEvidence(CopilotCLIProbeExecution{Outcome: "start-failed", FailureCode: copilotCLIProbeBYOKFailure})
	if !slices.Equal(alphaEvidence, want) {
		t.Fatalf("the refusal evidence must be the closed classification, got %v want %v", alphaEvidence, want)
	}
}

// TestP85T8AHostAuthenticationAndOtherChannelsBypassTheGate pins the boundary of
// the gate: host authentication hands the child no credential of ours, and the
// other channels never take this path, so neither is refused by it.
func TestP85T8AHostAuthenticationAndOtherChannelsBypassTheGate(t *testing.T) {
	hostProfile := AIProviderProfile{ProfileID: "github", ProviderType: "github-copilot", Model: "auto", AuthMode: "host"}
	hostProbe := &p85T8AProbeWithoutProof{}
	hostRecord, err := p85T8AChecker(hostProbe, nil).Check(context.Background(), hostProfile, "agent-runner-cli", 1)
	if err != nil {
		t.Fatal(err)
	}
	if hostRecord.Availability.Status != "available" || hostProbe.calls != 1 {
		t.Fatalf("a host-auth probe must not be gated by the BYOK proof, got %+v calls=%d", hostRecord.Availability, hostProbe.calls)
	}

	secretProbe := &p85T8AProbeWithoutProof{}
	resolver := &p85T8ACountingResolver{secret: []byte("p85-t8a-secret")}
	otherRecord, err := p85T8AChecker(secretProbe, resolver).Check(context.Background(), deepSeekProfile(), "sessionbridge-silent", 1)
	if err != nil {
		t.Fatal(err)
	}
	if otherRecord.Availability.Status != "available" || secretProbe.calls != 1 || resolver.calls != 1 {
		t.Fatalf("a non-agent-runner-cli channel must not be gated by the BYOK proof, got %+v calls=%d resolves=%d",
			otherRecord.Availability, secretProbe.calls, resolver.calls)
	}

	gatedProbe := &p85T8AProbeWithoutProof{}
	gatedResolver := &p85T8ACountingResolver{secret: []byte("p85-t8a-secret")}
	gatedRecord, err := p85T8AChecker(gatedProbe, gatedResolver).Check(context.Background(), deepSeekProfile(), "agent-runner-cli", 1)
	if err != nil {
		t.Fatal(err)
	}
	if gatedRecord.Availability.Reason == nil || *gatedRecord.Availability.Reason != copilotCLIProbeBYOKFailure {
		t.Fatalf("the secret agent-runner-cli path must stay gated, got %+v", gatedRecord.Availability)
	}
	if gatedProbe.calls != 0 || gatedResolver.calls != 0 {
		t.Fatalf("the gated path must not resolve or probe, got calls=%d resolves=%d", gatedProbe.calls, gatedResolver.calls)
	}
}

// p85T8ARawMapCarriesCredential is the PRE-F14 classification: it asks the caller's
// raw override map whether it names the credential key. It is kept only so the tests
// below can show the disagreement the fix removes - the child parses the SERIALIZED
// block, so the block is what decides.
func p85T8ARawMapCarriesCredential(environment map[string][]byte) bool {
	for name := range environment {
		if copilotCLIEnvironmentKeyEqual(name, copilotCLIProbeCredentialEnvironmentKey) {
			return true
		}
	}
	return false
}

// p85T8AClearBuffers is the test-side twin of the executor's deferred wipe: the
// serialized environment block holds the credential bytes, so a test that builds one
// clears it before it returns.
func p85T8AClearBuffers(t *testing.T, buffers [][]byte) {
	t.Helper()
	t.Cleanup(func() {
		for _, buffer := range buffers {
			clear(buffer)
		}
	})
}

// TestP85T8AEnvironmentKeyValidationRefusesAmbiguousKeys is the F14 (final review)
// counterexample set. The finding: the executor classified the credential from the
// caller's RAW MAP while the child parses `NAME=VALUE` at the FIRST `=`, so a key name
// that carries a separator, an empty name, or a platform-normalised duplicate let
// host/compatibility branches skip the end gate. The fix validates the names first and
// takes the classification from the SERIALIZED block that is handed to the child.
//
// Mutation checks (each one turns a named subtest red):
//   - delete the '=' / empty-name / NUL branches of copilotCLIProbeValidateEnvironmentKeys
//     => "key-name-containing-a-separator" and "empty-key-name" go red;
//   - take the credential classification from the raw map again (p85T8ARawMapCarriesCredential)
//     inside copilotCLIProbeCarriesCredential => "serialized-block-decides-the-auth-mode"
//     goes red, because the two classifications disagree.
func TestP85T8AEnvironmentKeyValidationRefusesAmbiguousKeys(t *testing.T) {
	t.Run("key-name-containing-a-separator", func(t *testing.T) {
		// The child takes everything before the FIRST '=', so this name delivers
		// COPILOT_PROVIDER_API_KEY to the child while the raw map still reports
		// "no credential".
		overrides := map[string][]byte{
			copilotCLIProbeCredentialEnvironmentKey + "=p85-t8a-smuggled": []byte("p85-t8a-secret"),
		}
		if p85T8ARawMapCarriesCredential(overrides) {
			t.Fatal("the raw map must NOT look like it carries the credential; that is the hole")
		}
		block := []string{copilotCLIProbeCredentialEnvironmentKey + "=p85-t8a-secret"}
		present, err := copilotCLIProbeCarriesCredential(block)
		if err != nil || !present {
			t.Fatalf("the child would receive the credential from that block, got present=%v err=%v", present, err)
		}
		if err := copilotCLIProbeValidateEnvironmentKeys(overrides); err == nil {
			t.Fatal("a key name carrying '=' must be refused before it reaches the block")
		} else if !errors.Is(err, ErrAIBYOKConfinement) {
			t.Fatalf("the refusal must carry ErrAIBYOKConfinement, got %v", err)
		}
		if _, _, err := copilotCLIManagedEnvironment(overrides); err == nil {
			t.Fatal("the managed environment must refuse the same key name")
		}
	})

	t.Run("empty-key-name", func(t *testing.T) {
		overrides := map[string][]byte{"": []byte("p85-t8a-secret")}
		if err := copilotCLIProbeValidateEnvironmentKeys(overrides); err == nil {
			t.Fatal("an empty key name is not addressable and must be refused")
		} else if !errors.Is(err, ErrAIBYOKConfinement) {
			t.Fatalf("the refusal must carry ErrAIBYOKConfinement, got %v", err)
		}
		// The block would parse as the empty name; the classifier refuses an entry
		// without a separator but must not silently accept "" as the credential key.
		if present, err := copilotCLIProbeCarriesCredential([]string{"=p85-t8a-secret"}); err != nil || present {
			t.Fatalf("an empty name is never the credential key, got present=%v err=%v", present, err)
		}
	})

	t.Run("nul-in-key-name", func(t *testing.T) {
		overrides := map[string][]byte{"P85\x00T8A": []byte("value")}
		if err := copilotCLIProbeValidateEnvironmentKeys(overrides); err == nil {
			t.Fatal("a key name carrying a NUL byte terminates the block entry and must be refused")
		}
	})

	t.Run("platform-normalised-duplicate", func(t *testing.T) {
		overrides := map[string][]byte{
			copilotCLIProbeCredentialEnvironmentKey:                  []byte("p85-t8a-canonical"),
			strings.ToLower(copilotCLIProbeCredentialEnvironmentKey): []byte("p85-t8a-variant"),
		}
		err := copilotCLIProbeValidateEnvironmentKeys(overrides)
		if runtime.GOOS == "windows" {
			// The environment block is case-insensitive, so the two names are two
			// spellings of ONE variable whose value depends on ordering.
			if err == nil {
				t.Fatal("two key names the platform cannot tell apart must be refused on windows")
			}
			if !errors.Is(err, ErrAIBYOKConfinement) {
				t.Fatalf("the refusal must carry ErrAIBYOKConfinement, got %v", err)
			}
			return
		}
		// A byte-exact platform really does have two variables, so the key-name gate
		// must not refuse it. The closed BYOK key set is what refuses the variant,
		// which is why the classification below is taken from the block, not the map.
		if err != nil {
			t.Fatalf("byte-exact key names are distinct variables and must not be refused here: %v", err)
		}
		variant := strings.ToLower(copilotCLIProbeCredentialEnvironmentKey)
		if slices.Contains(copilotCLIProbeBYOKEnvironmentKeys, variant) {
			t.Fatalf("the closed BYOK key set must not admit the case variant %q", variant)
		}
	})

	t.Run("credential-key-case-variant-is-not-a-credential", func(t *testing.T) {
		// The block the child would parse names the credential key only up to case.
		present, err := copilotCLIProbeCarriesCredential([]string{strings.ToLower(copilotCLIProbeCredentialEnvironmentKey) + "=p85-t8a-secret"})
		if runtime.GOOS == "windows" {
			if err == nil {
				t.Fatal("on a case-insensitive platform the child receives a credential the gate cannot name, so it must be refused")
			}
			if present {
				t.Fatal("a refused case variant must not be reported as a present credential")
			}
			return
		}
		if err != nil {
			t.Fatalf("on a byte-exact platform the variant is a different variable, not a refusal: %v", err)
		}
		if present {
			t.Fatal("on a byte-exact platform the variant is not the credential key")
		}
	})

	t.Run("serialized-block-decides-the-auth-mode", func(t *testing.T) {
		// Positive control: the map judgement and the block judgement agree, so the
		// declared-secret profile is classified as secret on both sides.
		overrides := copilotCLIProbeEnvironment(deepSeekProfile(), []byte("p85-t8a-secret"))
		block, buffers, err := copilotCLIManagedEnvironment(overrides)
		if err != nil {
			t.Fatalf("a well-formed BYOK override set must serialize: %v", err)
		}
		p85T8AClearBuffers(t, buffers)
		rawPresent := p85T8ARawMapCarriesCredential(overrides)
		blockPresent, blockErr := copilotCLIProbeCarriesCredential(block)
		if blockErr != nil {
			t.Fatalf("the serialized block must classify cleanly: %v", blockErr)
		}
		if rawPresent != blockPresent || !blockPresent {
			t.Fatalf("the block and the raw map must agree that the credential is present, got raw=%v block=%v", rawPresent, blockPresent)
		}
		if mode := copilotCLIProbeEffectiveAuthMode(CopilotCLIProbeInvocation{}, block); mode != "secret" {
			t.Fatalf("the serialized block implies the secret mode, got %q", mode)
		}

		// Positive control for the host branch: it carries no credential, so the end
		// gate passes it as host - for the same reason, measured on the same block.
		hostBlock, hostBuffers, err := copilotCLIManagedEnvironment(nil)
		if err != nil {
			t.Fatalf("the host environment must serialize: %v", err)
		}
		p85T8AClearBuffers(t, hostBuffers)
		hostPresent, hostErr := copilotCLIProbeCarriesCredential(hostBlock)
		if hostErr != nil || hostPresent {
			t.Fatalf("the host branch must clear the end gate, got present=%v err=%v", hostPresent, hostErr)
		}
		if mode := copilotCLIProbeEffectiveAuthMode(CopilotCLIProbeInvocation{AuthMode: "host"}, hostBlock); mode != "host" {
			t.Fatalf("a declared host invocation must stay host, got %q", mode)
		}
	})
}

// p85T8ABYOKProfile is the declared profile the envelope fixtures use.
//
// It differs from the bare deepSeekProfile in one field: it names a wire API, so
// copilotCLIProbeEnvironment emits the FULL closed BYOK key set. The frozen template
// carries that closed set, and the spawn-time binding is CONTAINMENT, not equality:
// every key the child would receive must be one the template froze. A profile that
// omits the optional wire-API key therefore yields a strict subset and still binds,
// while a key the template never named is a spill and is refused (T8a C2, see
// TestP85T8AEnvironmentBindingIsContainment).
func p85T8ABYOKProfile() AIProviderProfile {
	profile := deepSeekProfile()
	profile.WireAPI = "chat"
	return profile
}

// p85T8ARecordedInvocation is one deep copy of what the probe handed the executor.
type p85T8ARecordedInvocation struct {
	args        []string
	environment map[string][]byte
	frozen      *copilotCLIProbeSpawnSpec
}

// p85T8ARecordingExecutor records the invocation it was given instead of starting
// anything, so a test can compare the object that was verified with the object that
// would be spawned (P85 F13).
type p85T8ARecordingExecutor struct {
	bound       int64
	invocations []p85T8ARecordedInvocation
}

func (executor *p85T8ARecordingExecutor) VerifyBYOKConfinement() error { return nil }

func (executor *p85T8ARecordingExecutor) CopilotCLIProbeOutputBound() int64 {
	if executor.bound > 0 {
		return executor.bound
	}
	return defaultCopilotCLIOutputBytes
}

func (executor *p85T8ARecordingExecutor) Execute(_ context.Context, invocation CopilotCLIProbeInvocation) CopilotCLIProbeExecution {
	environment := make(map[string][]byte, len(invocation.Environment))
	for name, value := range invocation.Environment {
		environment[name] = slices.Clone(value)
	}
	executor.invocations = append(executor.invocations, p85T8ARecordedInvocation{
		args:        slices.Clone(invocation.Args),
		environment: environment,
		frozen:      invocation.FrozenSpawn,
	})
	return CopilotCLIProbeExecution{Outcome: "exited", ExitCode: 0}
}

// p85T8ADeclaredProbe builds a probe whose spawn inputs are DECLARED, which is the
// shape that freezes a launch template before the credential is resolved.
func p85T8ADeclaredProbe(t *testing.T, executor *p85T8ARecordingExecutor) (*CopilotCLIProbe, copilotCLIProbeSpawnSpec) {
	t.Helper()
	probe := &CopilotCLIProbe{
		Executable:       filepath.Join(t.TempDir(), "copilot"),
		WorkingDirectory: t.TempDir(),
		Executor:         executor,
		AuthMode:         "secret",
		Model:            p85T8ABYOKProfile().Model,
		MaximumRequests:  1,
	}
	spec, err := probe.frozenLaunchTemplate()
	if err != nil {
		t.Fatalf("the declared inputs must derive a launch template: %v", err)
	}
	return probe, spec
}

// TestP85T8AFrozenSpawnSpecIsTheVerifiedObject is the F13 (final review)
// counterexample: the template PROVEN before the credential was resolved must be the
// object the spawn instantiates, and no field of it may be re-derived at spawn time.
//
// Mutation checks: build the template again inside Probe instead of returning
// probe.frozen => "spawn-uses-the-frozen-object" goes red on the pointer comparison;
// drop the copilotCLIProbeSpawnSpecEqual check in launchSpecFor => every case of
// "a-changed-field-refuses-the-spawn" goes red.
func TestP85T8AFrozenSpawnSpecIsTheVerifiedObject(t *testing.T) {
	t.Run("spawn-uses-the-frozen-object", func(t *testing.T) {
		executor := &p85T8ARecordingExecutor{}
		probe, spec := p85T8ADeclaredProbe(t, executor)
		// This is the shared object the pre-Resolve proof validated: the pointer the
		// checker froze must be the pointer the spawn reads.
		probe.frozen = &spec
		if _, err := probe.Probe(context.Background(), p85T8ABYOKProfile(), []byte("p85-t8a-secret"), "agent-runner-cli", 1); err != nil {
			t.Fatalf("the probe must not error on a frozen template: %v", err)
		}
		if len(executor.invocations) != 1 {
			t.Fatalf("the executor must be called once, got %d", len(executor.invocations))
		}
		recorded := executor.invocations[0]
		if recorded.frozen != probe.frozen {
			t.Fatal("the spawn must instantiate the SAME frozen template that was verified, not a second derivation")
		}
		if !slices.Equal(recorded.args, spec.Args) {
			t.Fatalf("the spawned argv must be the frozen argv, got %v want %v", recorded.args, spec.Args)
		}
		if !slices.Equal(copilotCLIProbeSortedEnvironmentKeys(recorded.environment), spec.EnvironmentKeys) {
			t.Fatalf("the spawned key set %v must be the frozen key set %v",
				copilotCLIProbeSortedEnvironmentKeys(recorded.environment), spec.EnvironmentKeys)
		}
	})

	t.Run("a-changed-field-refuses-the-spawn", func(t *testing.T) {
		mutations := map[string]func(*copilotCLIProbeSpawnSpec){
			"authentication-mode": func(spec *copilotCLIProbeSpawnSpec) { spec.AuthMode = "host" },
			"argv": func(spec *copilotCLIProbeSpawnSpec) {
				spec.Args = append(slices.Clone(spec.Args), "--no-color")
			},
			"environment-keys": func(spec *copilotCLIProbeSpawnSpec) {
				spec.EnvironmentKeys = append(slices.Clone(spec.EnvironmentKeys), "P85_T8A_EXTRA")
			},
			"output-bound": func(spec *copilotCLIProbeSpawnSpec) { spec.OutputBound++ },
			"executable":   func(spec *copilotCLIProbeSpawnSpec) { spec.Executable += "-other" },
			"working-dir":  func(spec *copilotCLIProbeSpawnSpec) { spec.WorkingDirectory += "-other" },
		}
		for name, mutate := range mutations {
			t.Run(name, func(t *testing.T) {
				executor := &p85T8ARecordingExecutor{}
				probe, spec := p85T8ADeclaredProbe(t, executor)
				mutate(&spec)
				probe.frozen = &spec
				result, err := probe.Probe(context.Background(), p85T8ABYOKProfile(), []byte("p85-t8a-secret"), "agent-runner-cli", 1)
				if err != nil {
					t.Fatalf("a refused template must be a classified refusal, not an error: %v", err)
				}
				if result.Status != "unavailable" || result.Reason != copilotCLIProbeBYOKFailure {
					t.Fatalf("a drifted template must refuse under the BYOK class, got %+v", result)
				}
				if len(executor.invocations) != 0 {
					t.Fatalf("a drifted template must never reach the executor, got %d invocations", len(executor.invocations))
				}
			})
		}
	})

	t.Run("an-empty-verdict-can-never-be-frozen", func(t *testing.T) {
		invocation, spec, directory, usagePath, limit := p85T8ACompliantSpawn(t)
		frozen, err := copilotCLIProbeLaunchSpec(invocation.Executable, invocation.WorkingDirectory, "secret", p85T8ABYOKProfile().Model, 1, limit)
		if err != nil {
			t.Fatalf("the compliant envelope must derive a template: %v", err)
		}
		if !slices.Equal(frozen.Args, invocation.Args) ||
			!slices.Equal(frozen.EnvironmentKeys, copilotCLIProbeSortedEnvironmentKeys(invocation.Environment)) {
			t.Fatalf("the frozen template must be the one the envelope spawns, got %+v", frozen)
		}
		invocation.FrozenSpawn = &frozen
		if frozen.PlatformWriteProof != "" {
			t.Fatal("a freshly derived template carries no platform verdict")
		}
		if err := assertCopilotCLIProbeBYOKConfinement(context.Background(), invocation, spec, directory, usagePath, limit); err == nil {
			t.Fatal("a template with no recorded platform verdict must never be spawned")
		} else if !errors.Is(err, ErrAIBYOKConfinement) {
			t.Fatalf("the refusal must carry ErrAIBYOKConfinement, got %v", err)
		}

		verdict, verdictErr := copilotCLIProbePlatformWriteVerdict(context.Background(), nil, copilotCLIProbeSpawnSpecDigest(frozen))
		if verdictErr != nil {
			t.Skipf("S1 host: the platform verdict cannot be observed at all (%v), so only the empty-proof refusal is measurable here; the same-verdict leg needs a host that proves it", verdictErr)
		}
		frozen.PlatformWriteProof = verdict
		if err := assertCopilotCLIProbeBYOKConfinement(context.Background(), invocation, spec, directory, usagePath, limit); err != nil {
			t.Fatalf("a template proven under the current verdict must be spawnable: %v", err)
		}
		frozen.PlatformWriteProof = evidence.Digest("p85-t8a-other-verdict:1\n", []byte("p85-t8a"))
		if err := assertCopilotCLIProbeBYOKConfinement(context.Background(), invocation, spec, directory, usagePath, limit); err == nil {
			t.Fatal("a template proven under a DIFFERENT verdict must be refused")
		}
	})
}

// TestP85T8AEnvironmentBindingIsContainment pins the T8a C2 binding: the frozen
// template's key set is the closed OUTER envelope, and the spawn is bound to it by
// CONTAINMENT - every key the child would receive must be one the template froze -
// not by equality. A legitimate profile may name a strict subset (the wire-API entry
// is optional), so the earlier equality refused a spawn that had nothing spilled;
// containment admits it and still refuses every key outside the frozen envelope.
//
// Mutation checks: restore the equality comparison => "a-strict-subset-binds" goes
// red (the subset is refused as a spill); drop the byte-exact membership test =>
// "a-platform-case-variant-spills-outside-the-frozen-set" goes red on a byte-exact
// platform; delete the per-key ambiguity gate => "a-windows-case-variant-pair" goes
// red on windows.
func TestP85T8AEnvironmentBindingIsContainment(t *testing.T) {
	// frozenSpawn builds a compliant credential-bearing envelope and the template it
	// was proven against - the full closed key set, as a declared probe with a wire
	// API freezes. It records NO platform verdict: each case decides which verdict it
	// needs, so the binding stays measurable on an S1 host that cannot observe one.
	frozenSpawn := func(t *testing.T) (CopilotCLIProbeInvocation, guard.ProcessSpec, string, string, int64, copilotCLIProbeSpawnSpec) {
		t.Helper()
		invocation, spec, directory, usagePath, limit := p85T8ACompliantSpawn(t)
		frozen, err := copilotCLIProbeLaunchSpec(invocation.Executable, invocation.WorkingDirectory, invocation.AuthMode, p85T8ABYOKProfile().Model, 1, limit)
		if err != nil {
			t.Fatalf("the compliant envelope must derive a template: %v", err)
		}
		return invocation, spec, directory, usagePath, limit, frozen
	}

	t.Run("the-membership-test-is-byte-exact-and-subset-tolerant", func(t *testing.T) {
		frozen := copilotCLIProbeSortedEnvironmentKeys(copilotCLIProbeEnvironment(p85T8ABYOKProfile(), []byte("p85-t8a-secret")))
		subset := copilotCLIProbeEnvironment(deepSeekProfile(), []byte("p85-t8a-secret"))
		if len(subset) >= len(frozen) {
			t.Fatalf("the fixture must be a strict subset, got %d of %d keys", len(subset), len(frozen))
		}
		if !copilotCLIProbeEnvironmentKeysWithin(subset, frozen) {
			t.Fatal("a strict subset of the frozen key set must bind")
		}
		spilled := make(map[string][]byte, len(subset)+1)
		for name, value := range subset {
			spilled[name] = value
		}
		spilled["P85_T8A_EXTRA"] = []byte("p85-t8a-extra")
		if copilotCLIProbeEnvironmentKeysWithin(spilled, frozen) {
			t.Fatal("a key the frozen set does not name must not bind")
		}
		variant := make(map[string][]byte, len(subset)+1)
		for name, value := range subset {
			variant[name] = value
		}
		variant[strings.ToLower("COPILOT_PROVIDER_API_KEY")] = []byte("p85-t8a-secret")
		if copilotCLIProbeEnvironmentKeysWithin(variant, frozen) {
			t.Fatal("a case variant is a different name for a byte-exact membership test and must not bind")
		}
	})

	t.Run("a-strict-subset-binds", func(t *testing.T) {
		invocation, spec, directory, usagePath, limit, frozen := frozenSpawn(t)
		// The wire-API entry is optional, so this profile simply does not name it.
		delete(invocation.Environment, "COPILOT_PROVIDER_WIRE_API")
		if !copilotCLIProbeEnvironmentKeysWithin(invocation.Environment, frozen.EnvironmentKeys) {
			t.Fatalf("the fixture must be a strict subset of the frozen key set, got %v in %v",
				copilotCLIProbeSortedEnvironmentKeys(invocation.Environment), frozen.EnvironmentKeys)
		}
		invocation.FrozenSpawn = &frozen

		if verdict, verdictErr := copilotCLIProbePlatformWriteVerdict(context.Background(), nil, copilotCLIProbeSpawnSpecDigest(frozen)); verdictErr == nil {
			frozen.PlatformWriteProof = verdict
			if err := assertCopilotCLIProbeBYOKConfinement(context.Background(), invocation, spec, directory, usagePath, limit); err != nil {
				t.Fatalf("a strict subset of the frozen key set must bind end to end, got %v", err)
			}
			return
		}
		// S1 host: the platform verdict is checked AFTER the key binding, so the
		// binding's own answer is still observable - it must be the platform verdict
		// that refuses, never a spill.
		frozen.PlatformWriteProof = "p85-t8a-recorded-verdict"
		err := assertCopilotCLIProbeBYOKConfinement(context.Background(), invocation, spec, directory, usagePath, limit)
		if err == nil {
			t.Fatal("a host that cannot observe the platform verdict must not admit the spawn")
		}
		if strings.Contains(err.Error(), "spill outside the frozen key set") {
			t.Fatalf("a strict subset must never be refused as a spill: %v", err)
		}
		if !errors.Is(err, ErrAIBYOKConfinement) || !strings.Contains(err.Error(), "sandbox") {
			t.Fatalf("the only remaining refusal for a strict subset must be the platform verdict, got %v", err)
		}
	})

	t.Run("an-extra-key-spills-outside-the-frozen-set", func(t *testing.T) {
		invocation, spec, directory, usagePath, limit, frozen := frozenSpawn(t)
		invocation.Environment["P85_T8A_EXTRA"] = []byte("p85-t8a-extra")
		frozen.PlatformWriteProof = "p85-t8a-recorded-verdict"
		invocation.FrozenSpawn = &frozen
		err := assertCopilotCLIProbeBYOKConfinement(context.Background(), invocation, spec, directory, usagePath, limit)
		if err == nil {
			t.Fatal("a key the frozen template never named must be refused")
		}
		if !errors.Is(err, ErrAIBYOKConfinement) || !strings.Contains(err.Error(), "spill outside the frozen key set") {
			t.Fatalf("the refusal must come from the containment binding, got %v", err)
		}
	})

	t.Run("a-platform-case-variant-spills-outside-the-frozen-set", func(t *testing.T) {
		// F14 is not weakened by loosening equality: the membership test stays
		// byte-exact, so a spelling the frozen set does not name is a spill on every
		// platform - the child would otherwise receive a variable the gate cannot name.
		invocation, spec, directory, usagePath, limit, frozen := frozenSpawn(t)
		variant := strings.ToLower("COPILOT_PROVIDER_WIRE_API")
		delete(invocation.Environment, "COPILOT_PROVIDER_WIRE_API")
		invocation.Environment[variant] = []byte("chat")
		frozen.PlatformWriteProof = "p85-t8a-recorded-verdict"
		invocation.FrozenSpawn = &frozen
		err := assertCopilotCLIProbeBYOKConfinement(context.Background(), invocation, spec, directory, usagePath, limit)
		if err == nil {
			t.Fatal("a platform-normalised spelling the frozen set does not name must be refused")
		}
		if !errors.Is(err, ErrAIBYOKConfinement) || !strings.Contains(err.Error(), "spill outside the frozen key set") {
			t.Fatalf("a case variant must be refused as a spill by the byte-exact binding, got %v", err)
		}
	})

	t.Run("a-windows-case-variant-pair-stays-refused", func(t *testing.T) {
		invocation, spec, directory, usagePath, limit, frozen := frozenSpawn(t)
		invocation.Environment[strings.ToLower("COPILOT_PROVIDER_WIRE_API")] = []byte("chat")
		frozen.PlatformWriteProof = "p85-t8a-recorded-verdict"
		invocation.FrozenSpawn = &frozen
		err := assertCopilotCLIProbeBYOKConfinement(context.Background(), invocation, spec, directory, usagePath, limit)
		if err == nil {
			t.Fatal("two spellings of one variable must be refused")
		}
		if !errors.Is(err, ErrAIBYOKConfinement) {
			t.Fatalf("the refusal must carry ErrAIBYOKConfinement, got %v", err)
		}
		if runtime.GOOS == "windows" {
			if !strings.Contains(err.Error(), "cannot tell apart") {
				t.Fatalf("the per-key ambiguity gate must classify the pair before the spill check, got %v", err)
			}
			return
		}
		if !strings.Contains(err.Error(), "spill outside the frozen key set") {
			t.Fatalf("a byte-exact platform must refuse the variant as a spill, got %v", err)
		}
	})
}

// p85T8AAllowedSandbox is the "the platform guard is withdrawn" stand-in: it answers
// Allowed for every check, so nothing but that guard can refuse the envelope.
type p85T8AAllowedSandbox struct{}

func (p85T8AAllowedSandbox) Check(_ context.Context, spec guard.CheckSpec) (guard.Verdict, error) {
	// P85 F13: echo the asked specification, which is what a conforming layer does.
	return guard.Verdict{Status: guard.VerdictAllowed, SpecDigest: spec.SpecDigest}, nil
}

// p85T8ACompliantSpawn builds the envelope the executor itself constructs for a
// credential-bearing probe - the screened argv, the serialized environment block that
// really carries the credential, and the bounded in-memory capture. Every element
// assertion passes on it, which is what makes the platform verdict the only guard
// left (P85 F18: the earlier counterexample used a shell argv that the argv allowlist
// would refuse on its own, so it could never isolate the platform guard).
func p85T8ACompliantSpawn(t *testing.T) (CopilotCLIProbeInvocation, guard.ProcessSpec, string, string, int64) {
	t.Helper()
	directory := t.TempDir()
	executable := filepath.Join(directory, "copilot")
	profile := p85T8ABYOKProfile()
	args := copilotCLIProbeArgs(profile.Model, directory, 1, "secret")
	environment := copilotCLIProbeEnvironment(profile, []byte("p85-t8a-secret"))
	if !p85T8ARawMapCarriesCredential(environment) {
		t.Fatal("the compliant envelope must carry the credential")
	}
	block, buffers, err := copilotCLIManagedEnvironment(environment)
	if err != nil {
		t.Fatalf("the compliant environment must serialize: %v", err)
	}
	p85T8AClearBuffers(t, buffers)
	usagePath := filepath.Join(directory, "usage.json")
	limit := int64(defaultCopilotCLIOutputBytes)
	output := newCopilotCLIProbeOutput(limit, func() {})
	spec := guard.ProcessSpec{
		Command: executable,
		Args:    copilotCLIProbeExecArgs(args, usagePath),
		Dir:     directory,
		Env:     block,
		Stdout:  output.stdout,
		Stderr:  output.stderr,
	}
	invocation := CopilotCLIProbeInvocation{
		Executable:       executable,
		WorkingDirectory: directory,
		Args:             args,
		Environment:      environment,
		AuthMode:         "secret",
	}
	return invocation, spec, directory, usagePath, limit
}

// TestP85T8APlatformVerdictAloneControlsTheSpawn is the F18 (final review) mutation
// control: with a COMPLIANT credential-bearing envelope, the platform write verdict
// is the only thing that may decide whether the child exists.
//
// The executor passes a nil sandbox layer to guard.CheckSandbox, so there is no seam
// that lets a test withdraw ONLY that guard and watch a real child start. What is
// measured instead is the envelope-level equivalence:
//   - the compliant envelope is refused, and the refusal text names the platform
//     check (the per-element assertions each name their own element, so a refusal
//     that does not name the sandbox cannot be any of them);
//   - an Allowed double removes exactly that step and nothing else, so "the platform
//     guard withdrawn" is the difference between refused and admitted;
//   - a non-compliant element is refused even under the Allowed double, so the
//     envelope really is the thing that isolates the verdict.
//
// The executor-level zero-start counterexample moved to
// TestP85F18ThePlatformVerdictAloneDecidesTheStart (p85_f18_start_count_test.go)
// once the executor gained an injectable Sandbox field and the
// copilotCLIProbeRunManaged observation point (P85 design rev 3, section 7): the
// "real child starts when only the platform guard is withdrawn" leg is now
// measurable, so the EXEMPT-REGISTERED P85/F18 note that used to live here is
// retired. This test keeps the envelope-level legs: a compliant envelope is refused
// only by the platform check, an Allowed double removes exactly that step, and a
// non-compliant element (M12) is still refused under the Allowed verdict by its own
// element rather than by the sandbox.
func TestP85T8APlatformVerdictAloneControlsTheSpawn(t *testing.T) {
	invocation, spec, directory, usagePath, limit := p85T8ACompliantSpawn(t)
	platformErr := guard.CheckSandbox(context.Background(), nil, guard.CheckSpec{Reason: guard.SandboxReasonBYOKProbe, SpecDigest: p85T8ASandboxProofSpecDigest})
	confinementErr := assertCopilotCLIProbeBYOKConfinement(context.Background(), invocation, spec, directory, usagePath, limit)

	t.Run("compliant-envelope-is-refused-by-the-platform-verdict-alone", func(t *testing.T) {
		if platformErr == nil {
			if confinementErr != nil {
				t.Fatalf("the platform verdict is Allowed, so the compliant envelope must be admitted, got %v", confinementErr)
			}
			return
		}
		if confinementErr == nil {
			t.Fatalf("the platform verdict refused this host, but the compliant envelope was admitted")
		}
		if !errors.Is(confinementErr, ErrAIBYOKConfinement) {
			t.Fatalf("the refusal must carry ErrAIBYOKConfinement, got %v", confinementErr)
		}
		if !strings.Contains(confinementErr.Error(), "sandbox") {
			t.Fatalf("a compliant envelope can only be refused by the platform check, got %q", confinementErr)
		}
	})

	t.Run("withdrawing-the-verdict-is-the-only-difference", func(t *testing.T) {
		if err := guard.CheckSandbox(context.Background(), p85T8AAllowedSandbox{}, guard.CheckSpec{Reason: guard.SandboxReasonBYOKProbe, SpecDigest: p85T8ASandboxProofSpecDigest}); err != nil {
			t.Fatalf("an Allowed verdict must pass the platform gate, got %v", err)
		}
		for name, verdict := range map[string]guard.Verdict{
			"unknown": {Status: guard.VerdictUnknown, Reason: "S1 sandbox enforcement checks not implemented"},
			"denied":  {Status: guard.VerdictDenied, Reason: "the environment provably lacks write restrictions"},
		} {
			if err := guard.CheckSandbox(context.Background(), &p85T8ASandboxDouble{verdict: verdict}, guard.CheckSpec{Reason: guard.SandboxReasonBYOKProbe, SpecDigest: p85T8ASandboxProofSpecDigest}); err == nil {
				t.Fatalf("%s: a non-Allowed verdict must refuse the spawn", name)
			}
		}
	})

	t.Run("a-non-compliant-element-is-refused-even-when-the-verdict-allows", func(t *testing.T) {
		if err := guard.CheckSandbox(context.Background(), p85T8AAllowedSandbox{}, guard.CheckSpec{Reason: guard.SandboxReasonBYOKProbe, SpecDigest: p85T8ASandboxProofSpecDigest}); err != nil {
			t.Fatalf("the control needs an Allowed verdict, got %v", err)
		}
		// The shell argv the earlier counterexample used: refused by the allowlist,
		// not by the platform verdict - which is exactly why it could not isolate it.
		broken := spec
		broken.Args = p85T8AMarkerArgs(filepath.Join(directory, "p85-t8a.marker"))
		err := assertCopilotCLIProbeBYOKConfinement(context.Background(), invocation, broken, directory, usagePath, limit)
		if err == nil {
			t.Fatal("a non-whitelisted argv must be refused")
		}
		if strings.Contains(err.Error(), "sandbox") {
			t.Fatalf("the argv must be diagnosed by the argv assertion, not by the platform check: %q", err)
		}
	})
}
