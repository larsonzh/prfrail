package adapters

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeAISecretResolver struct {
	secret []byte
	err    error
}

func (resolver fakeAISecretResolver) Resolve(context.Context, string) ([]byte, error) {
	return resolver.secret, resolver.err
}

type fakeAIChannelProbe struct {
	called     bool
	credential []byte
	channel    string
	result     AIProbeResult
	err        error
}

func (probe *fakeAIChannelProbe) Probe(_ context.Context, _ AIProviderProfile, credential []byte, channel string, _ int) (AIProbeResult, error) {
	probe.called = true
	probe.credential = credential
	probe.channel = channel
	return probe.result, probe.err
}

func deepSeekProfile() AIProviderProfile {
	return AIProviderProfile{
		ProfileID:    "deepseek-anthropic",
		ProviderType: "anthropic",
		BaseURL:      "https://api.deepseek.com/anthropic",
		Model:        "deepseek-flash",
		AuthMode:     "secret",
		SecretRef:    "windows-credential:proofrail/deepseek",
	}
}

func availabilityChecker(probe AIChannelProbe, secrets AISecretResolver) *AIAvailabilityChecker {
	return &AIAvailabilityChecker{
		Clock:   func() time.Time { return time.Date(2026, 9, 11, 8, 0, 0, 0, time.UTC) },
		IDs:     func() string { return "probe-ai-check" },
		Secrets: secrets,
		Probe:   probe,
	}
}

func TestAIAvailabilityCheckBlocksMissingSecretWithoutProbe(t *testing.T) {
	probe := &fakeAIChannelProbe{}
	checker := availabilityChecker(probe, fakeAISecretResolver{err: ErrAISecretNotFound})
	record, err := checker.Check(context.Background(), deepSeekProfile(), "agent-runner-cli", 1)
	if err != nil {
		t.Fatal(err)
	}
	if probe.called {
		t.Fatal("probe called without a secret")
	}
	if record.Availability.Status != "unavailable" || record.Availability.Reason == nil || *record.Availability.Reason != "secret_missing" {
		t.Fatalf("unexpected missing-secret result: %+v", record.Availability)
	}
}

func TestAIAvailabilityCheckUsesOnlyRequestedChannelAndClearsSecret(t *testing.T) {
	secret := []byte("not-a-real-secret")
	probe := &fakeAIChannelProbe{result: AIProbeResult{
		Status:       "available",
		RequestsUsed: 1,
		Evidence:     []string{queueHashTwo},
	}}
	checker := availabilityChecker(probe, fakeAISecretResolver{secret: secret})
	record, err := checker.Check(context.Background(), deepSeekProfile(), "sessionbridge-silent", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !probe.called || probe.channel != "sessionbridge-silent" || record.Availability.Channel != "sessionbridge-silent" {
		t.Fatalf("probe crossed channels: called=%v probe=%s record=%s", probe.called, probe.channel, record.Availability.Channel)
	}
	for _, value := range probe.credential {
		if value != 0 {
			t.Fatal("resolved secret was not cleared after probing")
		}
	}
}

func TestAIAvailabilityCheckSanitizesInvalidAndUnclassifiedResults(t *testing.T) {
	for name, probe := range map[string]*fakeAIChannelProbe{
		"invalid-response": {result: AIProbeResult{Status: "available", RequestsUsed: 0}},
		"probe-error":      {result: AIProbeResult{RequestsUsed: 1}, err: errors.New("provider detail must not escape")},
	} {
		t.Run(name, func(t *testing.T) {
			checker := availabilityChecker(probe, fakeAISecretResolver{secret: []byte("not-a-real-secret")})
			record, err := checker.Check(context.Background(), deepSeekProfile(), "agent-runner-cli", 1)
			if err != nil {
				t.Fatal(err)
			}
			if record.Availability.Status == "available" || record.Availability.Reason == nil {
				t.Fatalf("invalid probe result did not fail closed: %+v", record.Availability)
			}
		})
	}
}

func TestAIAvailabilityCheckHostAuthenticationDoesNotRequireSecretStore(t *testing.T) {
	profile := AIProviderProfile{ProfileID: "github", ProviderType: "github-copilot", Model: "auto", AuthMode: "host"}
	probe := &fakeAIChannelProbe{result: AIProbeResult{Status: "unavailable", Reason: "account_unavailable", Evidence: []string{queueHashTwo}}}
	checker := availabilityChecker(probe, nil)
	record, err := checker.Check(context.Background(), profile, "agent-runner-cli", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !probe.called || record.Availability.Reason == nil || *record.Availability.Reason != "account_unavailable" {
		t.Fatalf("unexpected host-auth result: %+v", record.Availability)
	}
}
