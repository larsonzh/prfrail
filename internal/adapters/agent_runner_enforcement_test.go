package adapters

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

const platformHash = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func admissionPolicy() AgentRunnerAdmissionPolicy {
	return AgentRunnerAdmissionPolicy{
		OS:           "windows",
		Arch:         "amd64",
		PlatformHash: platformHash,
		EvaluatedAt:  time.Date(2026, 9, 11, 7, 0, 0, 0, time.UTC),
		MaximumAge:   24 * time.Hour,
	}
}

func agentRunnerEnforcement(capability AgentRunnerCapabilityRecord) AgentRunnerEnforcement {
	verified := capabilityFinding("verified")
	return AgentRunnerEnforcement{
		EnforcementID:         "enforcement-one",
		ProviderID:            "windows-sandbox",
		ProbedAt:              "2026-09-11T06:00:00.000Z",
		CapabilityRecordHash:  capability.RecordHash,
		ExecutableHash:        capability.Capability.ExecutableHash,
		ConfigHash:            capability.Capability.ConfigHash,
		PlatformHash:          platformHash,
		EnforcementConfigHash: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		Disposition:           "verified",
		Matrix: AgentRunnerEnforcementMatrix{
			ToolControl:    verified,
			NetworkControl: verified,
		},
		Evidence: []string{queueHashFour},
	}
}

func blockedAgentRunnerCapability(t *testing.T) AgentRunnerCapabilityRecord {
	t.Helper()
	capability := agentRunnerCapability()
	capability.Disposition = "blocked"
	capability.Matrix.ToolControl = capabilityFinding("unsupported")
	capability.Matrix.NetworkControl = capabilityFinding("unsupported")
	record, err := NewAgentRunnerCapabilityRecord(capability)
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func TestAgentRunnerEnforcementRecordCanonicalRoundTrip(t *testing.T) {
	capability := blockedAgentRunnerCapability(t)
	record, err := NewAgentRunnerEnforcementRecord(agentRunnerEnforcement(capability))
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeAgentRunnerEnforcementRecord(wire)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.RecordHash != record.RecordHash {
		t.Fatalf("record hash mismatch: %s != %s", decoded.RecordHash, record.RecordHash)
	}
}

func TestAgentRunnerEnforcementRejectsTamperingAndUnknownFields(t *testing.T) {
	capability := blockedAgentRunnerCapability(t)
	record, err := NewAgentRunnerEnforcementRecord(agentRunnerEnforcement(capability))
	if err != nil {
		t.Fatal(err)
	}
	record.Enforcement.ConfigHash = queueHashOne
	if err := ValidateAgentRunnerEnforcementRecord(record); !errors.Is(err, ErrInvalidAgentRunnerEnforcement) {
		t.Fatalf("tampered record accepted: %v", err)
	}
	record, err = NewAgentRunnerEnforcementRecord(agentRunnerEnforcement(capability))
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	wire = append(wire[:len(wire)-1], []byte(`,"unknown":true}`)...)
	if _, err := DecodeAgentRunnerEnforcementRecord(wire); err == nil {
		t.Fatal("unknown field accepted")
	}
}

func TestAgentRunnerEnforcementReplayAndConflict(t *testing.T) {
	capability := blockedAgentRunnerCapability(t)
	first, err := NewAgentRunnerEnforcementRecord(agentRunnerEnforcement(capability))
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerEnforcementIndex([]AgentRunnerEnforcementRecord{first})
	if err != nil {
		t.Fatal(err)
	}
	if replayed, err := index.Record(first); err != nil || !replayed {
		t.Fatalf("expected replay: replayed=%v err=%v", replayed, err)
	}
	changed := agentRunnerEnforcement(capability)
	changed.EnforcementConfigHash = queueHashOne
	conflict, err := NewAgentRunnerEnforcementRecord(changed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := index.Record(conflict); !errors.Is(err, ErrAgentRunnerEnforcementConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	var nilIndex *AgentRunnerEnforcementIndex
	if _, err := nilIndex.Record(first); !errors.Is(err, ErrInvalidAgentRunnerEnforcement) {
		t.Fatalf("nil index accepted: %v", err)
	}
}

func TestAgentRunnerEnforcementConcurrentReplay(t *testing.T) {
	capability := blockedAgentRunnerCapability(t)
	record, err := NewAgentRunnerEnforcementRecord(agentRunnerEnforcement(capability))
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerEnforcementIndex(nil)
	if err != nil {
		t.Fatal(err)
	}

	const workers = 24
	results := make([]struct {
		replayed bool
		err      error
	}, workers)
	start := make(chan struct{})
	var ready sync.WaitGroup
	var done sync.WaitGroup
	ready.Add(workers)
	done.Add(workers)
	for i := 0; i < workers; i++ {
		go func(i int) {
			defer done.Done()
			ready.Done()
			<-start
			results[i].replayed, results[i].err = index.Record(record)
		}(i)
	}
	ready.Wait()
	close(start)
	done.Wait()

	inserted := 0
	replayed := 0
	for _, result := range results {
		if result.err != nil {
			t.Fatalf("concurrent replay failed: %v", result.err)
		}
		if result.replayed {
			replayed++
		} else {
			inserted++
		}
	}
	if inserted != 1 || replayed != workers-1 {
		t.Fatalf("unexpected concurrent replay counts: inserted=%d replayed=%d", inserted, replayed)
	}
}

func TestAgentRunnerEnforcementConcurrentConflict(t *testing.T) {
	capability := blockedAgentRunnerCapability(t)
	first, err := NewAgentRunnerEnforcementRecord(agentRunnerEnforcement(capability))
	if err != nil {
		t.Fatal(err)
	}
	changed := agentRunnerEnforcement(capability)
	changed.EnforcementConfigHash = queueHashOne
	second, err := NewAgentRunnerEnforcementRecord(changed)
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerEnforcementIndex(nil)
	if err != nil {
		t.Fatal(err)
	}

	records := []AgentRunnerEnforcementRecord{first, second}
	results := make([]struct {
		replayed bool
		err      error
	}, len(records))
	start := make(chan struct{})
	var ready sync.WaitGroup
	var done sync.WaitGroup
	ready.Add(len(records))
	done.Add(len(records))
	for i := range records {
		go func(i int) {
			defer done.Done()
			ready.Done()
			<-start
			results[i].replayed, results[i].err = index.Record(records[i])
		}(i)
	}
	ready.Wait()
	close(start)
	done.Wait()

	successes := 0
	conflicts := 0
	winner := -1
	loser := -1
	for i, result := range results {
		switch {
		case result.err == nil:
			successes++
			winner = i
		case errors.Is(result.err, ErrAgentRunnerEnforcementConflict):
			conflicts++
			loser = i
		default:
			t.Fatalf("unexpected concurrent conflict result: %v", result.err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("unexpected concurrent conflict counts: successes=%d conflicts=%d", successes, conflicts)
	}

	if replayed, err := index.Record(records[winner]); err != nil || !replayed {
		t.Fatalf("winner was not a replay: replayed=%v err=%v", replayed, err)
	}
	if replayed, err := index.Record(records[loser]); !errors.Is(err, ErrAgentRunnerEnforcementConflict) || replayed {
		t.Fatalf("loser stopped conflicting: replayed=%v err=%v", replayed, err)
	}
}

func TestAgentRunnerEnforcementRejectsDispositionMatrixMismatch(t *testing.T) {
	capability := blockedAgentRunnerCapability(t)
	for name, mutate := range map[string]func(*AgentRunnerEnforcement){
		"verified-with-gap": func(value *AgentRunnerEnforcement) {
			value.Matrix.NetworkControl = capabilityFinding("unsupported")
		},
		"blocked-while-complete": func(value *AgentRunnerEnforcement) {
			value.Disposition = "blocked"
		},
	} {
		t.Run(name, func(t *testing.T) {
			value := agentRunnerEnforcement(capability)
			mutate(&value)
			if _, err := NewAgentRunnerEnforcementRecord(value); !errors.Is(err, ErrInvalidAgentRunnerEnforcement) {
				t.Fatalf("matrix mismatch accepted: %v", err)
			}
		})
	}
}

func TestAgentRunnerAdmissionRequiresBoundCompleteEnforcement(t *testing.T) {
	capability := blockedAgentRunnerCapability(t)
	policy := admissionPolicy()
	if err := ValidateAgentRunnerAdmission(capability, nil, policy); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
		t.Fatalf("missing enforcement admitted: %v", err)
	}
	enforcementValue := agentRunnerEnforcement(capability)
	enforcement, err := NewAgentRunnerEnforcementRecord(enforcementValue)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAgentRunnerAdmission(capability, &enforcement, policy); err != nil {
		t.Fatalf("complete bound enforcement rejected: %v", err)
	}
	for name, mutate := range map[string]func(*AgentRunnerAdmissionPolicy){
		"os":            func(value *AgentRunnerAdmissionPolicy) { value.OS = "linux" },
		"arch":          func(value *AgentRunnerAdmissionPolicy) { value.Arch = "arm64" },
		"platform-hash": func(value *AgentRunnerAdmissionPolicy) { value.PlatformHash = queueHashOne },
	} {
		t.Run("policy-"+name, func(t *testing.T) {
			wrongPlatform := policy
			mutate(&wrongPlatform)
			if err := ValidateAgentRunnerAdmission(capability, &enforcement, wrongPlatform); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
				t.Fatalf("wrong %s admitted: %v", name, err)
			}
		})
	}
	enforcementValue.CapabilityRecordHash = queueHashOne
	wrongCapability, err := NewAgentRunnerEnforcementRecord(enforcementValue)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAgentRunnerAdmission(capability, &wrongCapability, policy); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
		t.Fatalf("wrong candidate admitted: %v", err)
	}
	for name, mutate := range map[string]func(*AgentRunnerEnforcement){
		"executable": func(value *AgentRunnerEnforcement) { value.ExecutableHash = queueHashTwo },
		"config":     func(value *AgentRunnerEnforcement) { value.ConfigHash = queueHashOne },
		"platform":   func(value *AgentRunnerEnforcement) { value.PlatformHash = queueHashFour },
	} {
		t.Run(name, func(t *testing.T) {
			value := agentRunnerEnforcement(capability)
			mutate(&value)
			record, err := NewAgentRunnerEnforcementRecord(value)
			if err != nil {
				t.Fatal(err)
			}
			if err := ValidateAgentRunnerAdmission(capability, &record, policy); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
				t.Fatalf("wrong %s binding admitted: %v", name, err)
			}
		})
	}
}

func TestAgentRunnerAdmissionRejectsInvalidPolicy(t *testing.T) {
	capability, err := NewAgentRunnerCapabilityRecord(agentRunnerCapability())
	if err != nil {
		t.Fatal(err)
	}
	for name, mutate := range map[string]func(*AgentRunnerAdmissionPolicy){
		"os":            func(policy *AgentRunnerAdmissionPolicy) { policy.OS = "freebsd" },
		"arch":          func(policy *AgentRunnerAdmissionPolicy) { policy.Arch = "386" },
		"platform-hash": func(policy *AgentRunnerAdmissionPolicy) { policy.PlatformHash = "invalid" },
		"time":          func(policy *AgentRunnerAdmissionPolicy) { policy.EvaluatedAt = time.Time{} },
		"age":           func(policy *AgentRunnerAdmissionPolicy) { policy.MaximumAge = 0 },
	} {
		t.Run(name, func(t *testing.T) {
			policy := admissionPolicy()
			mutate(&policy)
			if err := ValidateAgentRunnerAdmission(capability, nil, policy); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
				t.Fatalf("invalid %s policy admitted: %v", name, err)
			}
		})
	}
}

func TestAgentRunnerAdmissionRejectsStaleAndFutureEnforcement(t *testing.T) {
	capability := blockedAgentRunnerCapability(t)
	policy := admissionPolicy()
	for name, testCase := range map[string]struct {
		probedAt   string
		maximumAge time.Duration
	}{
		"before-capability": {probedAt: "2026-09-10T01:59:59.999Z", maximumAge: 48 * time.Hour},
		"too-old":           {probedAt: "2026-09-10T06:59:59.999Z", maximumAge: 24 * time.Hour},
		"future":            {probedAt: "2026-09-11T07:00:00.001Z", maximumAge: 24 * time.Hour},
	} {
		t.Run(name, func(t *testing.T) {
			value := agentRunnerEnforcement(capability)
			value.ProbedAt = testCase.probedAt
			record, err := NewAgentRunnerEnforcementRecord(value)
			if err != nil {
				t.Fatal(err)
			}
			policy.MaximumAge = testCase.maximumAge
			if err := ValidateAgentRunnerAdmission(capability, &record, policy); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
				t.Fatalf("%s enforcement admitted: %v", name, err)
			}
		})
	}
}

func TestAgentRunnerAdmissionAcceptsFreshnessBoundaries(t *testing.T) {
	capability := blockedAgentRunnerCapability(t)
	policy := admissionPolicy()
	for name, testCase := range map[string]struct {
		probedAt   string
		maximumAge time.Duration
	}{
		"at-capability":  {probedAt: capability.Capability.ProbedAt, maximumAge: 48 * time.Hour},
		"at-maximum-age": {probedAt: "2026-09-10T07:00:00.000Z", maximumAge: 24 * time.Hour},
		"at-evaluation":  {probedAt: "2026-09-11T07:00:00.000Z", maximumAge: 24 * time.Hour},
	} {
		t.Run(name, func(t *testing.T) {
			value := agentRunnerEnforcement(capability)
			value.ProbedAt = testCase.probedAt
			record, err := NewAgentRunnerEnforcementRecord(value)
			if err != nil {
				t.Fatal(err)
			}
			policy.MaximumAge = testCase.maximumAge
			if err := ValidateAgentRunnerAdmission(capability, &record, policy); err != nil {
				t.Fatalf("freshness boundary %s rejected: %v", name, err)
			}
		})
	}
}

func TestAgentRunnerEnforcementRejectsOneSidedCompensation(t *testing.T) {
	capability := blockedAgentRunnerCapability(t)
	enforcement := agentRunnerEnforcement(capability)
	enforcement.Disposition = "blocked"
	enforcement.Matrix.NetworkControl = capabilityFinding("unknown")
	record, err := NewAgentRunnerEnforcementRecord(enforcement)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAgentRunnerAdmission(capability, &record, admissionPolicy()); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
		t.Fatalf("one-sided enforcement admitted: %v", err)
	}
}

func TestAgentRunnerAdmissionCombinesControlSources(t *testing.T) {
	capabilityValue := agentRunnerCapability()
	capabilityValue.Disposition = "blocked"
	capabilityValue.Matrix.NetworkControl = capabilityFinding("unsupported")
	capability, err := NewAgentRunnerCapabilityRecord(capabilityValue)
	if err != nil {
		t.Fatal(err)
	}
	enforcementValue := agentRunnerEnforcement(capability)
	enforcementValue.Disposition = "blocked"
	enforcementValue.Matrix.ToolControl = capabilityFinding("unsupported")
	enforcement, err := NewAgentRunnerEnforcementRecord(enforcementValue)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAgentRunnerAdmission(capability, &enforcement, admissionPolicy()); err != nil {
		t.Fatalf("complementary control evidence rejected: %v", err)
	}
}

func TestAgentRunnerAdmissionRejectsNonControlGap(t *testing.T) {
	capabilityValue := agentRunnerCapability()
	capabilityValue.Disposition = "blocked"
	capabilityValue.Matrix.SessionResume = capabilityFinding("unsupported")
	capability, err := NewAgentRunnerCapabilityRecord(capabilityValue)
	if err != nil {
		t.Fatal(err)
	}
	enforcement, err := NewAgentRunnerEnforcementRecord(agentRunnerEnforcement(capability))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAgentRunnerAdmission(capability, &enforcement, admissionPolicy()); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
		t.Fatalf("non-control gap admitted: %v", err)
	}
}

func TestCompatibleAgentRunnerAdmissionNeedsNoEnforcement(t *testing.T) {
	capability, err := NewAgentRunnerCapabilityRecord(agentRunnerCapability())
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAgentRunnerAdmission(capability, nil, admissionPolicy()); err != nil {
		t.Fatalf("compatible candidate rejected: %v", err)
	}
}

func TestCompatibleAgentRunnerAdmissionRejectsPlatformMismatch(t *testing.T) {
	capability, err := NewAgentRunnerCapabilityRecord(agentRunnerCapability())
	if err != nil {
		t.Fatal(err)
	}
	for name, mutate := range map[string]func(*AgentRunnerAdmissionPolicy){
		"os":            func(policy *AgentRunnerAdmissionPolicy) { policy.OS = "linux" },
		"arch":          func(policy *AgentRunnerAdmissionPolicy) { policy.Arch = "arm64" },
		"platform-hash": func(policy *AgentRunnerAdmissionPolicy) { policy.PlatformHash = queueHashFour },
	} {
		t.Run(name, func(t *testing.T) {
			policy := admissionPolicy()
			mutate(&policy)
			if err := ValidateAgentRunnerAdmission(capability, nil, policy); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
				t.Fatalf("compatible candidate admitted on wrong %s: %v", name, err)
			}
		})
	}
}

func TestPinnedCandidateCannotStartWithoutEnforcement(t *testing.T) {
	wire, err := os.ReadFile(filepath.Join("..", "..", "testdata", "agent-runner", "capability-probes", "github-copilot-cli-windows.json"))
	if err != nil {
		t.Fatal(err)
	}
	capability, err := DecodeAgentRunnerCapabilityRecord(wire)
	if err != nil {
		t.Fatal(err)
	}
	verified := 0
	nonVerified := make(map[string]string)
	for name, finding := range map[string]AgentRunnerCapabilityFinding{
		"noninteractive":            capability.Capability.Matrix.Noninteractive,
		"cwd":                       capability.Capability.Matrix.CWD,
		"eventStreamOrCompleteLogs": capability.Capability.Matrix.EventStreamOrCompleteLogs,
		"sessionCreate":             capability.Capability.Matrix.SessionCreate,
		"sessionResume":             capability.Capability.Matrix.SessionResume,
		"cancellation":              capability.Capability.Matrix.Cancellation,
		"processTreeStop":           capability.Capability.Matrix.ProcessTreeStop,
		"toolControl":               capability.Capability.Matrix.ToolControl,
		"networkControl":            capability.Capability.Matrix.NetworkControl,
		"permissionControl":         capability.Capability.Matrix.PermissionControl,
		"usage":                     capability.Capability.Matrix.Usage,
		"unattendedConfirmations":   capability.Capability.Matrix.UnattendedConfirmations,
	} {
		if finding.Status == "verified" {
			verified++
		} else {
			nonVerified[name] = finding.Status
		}
	}
	if verified != 10 || len(nonVerified) != 2 || nonVerified["toolControl"] != "unsupported" || nonVerified["networkControl"] != "unsupported" {
		t.Fatalf("unexpected pinned capability shape: verified=%d nonVerified=%v", verified, nonVerified)
	}
	if err := ValidateAgentRunnerAdmission(capability, nil, admissionPolicy()); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
		t.Fatalf("pinned incompatible candidate admitted: %v", err)
	}
}
