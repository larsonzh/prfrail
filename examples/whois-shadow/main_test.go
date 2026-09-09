package main

import (
	"strings"
	"testing"
)

func TestFrozenShadowMatchesBaseline(t *testing.T) {
	if err := compareFiles("baseline.json", "shadow.json"); err != nil {
		t.Fatal(err)
	}
}

func TestShadowEnvelopeFailsClosed(t *testing.T) {
	baseline := fixture{Source: &source{SHA256: "source-digest"}, Cases: []shadowCase{{ID: "one"}}}
	shadow := fixture{BaselineSHA256: "other-digest", Mode: "offline-read-only", Cases: []shadowCase{{ID: "one"}}}

	if err := validateEnvelope(baseline, shadow); err == nil || err.Error() != "baseline-digest-mismatch" {
		t.Fatalf("expected baseline-digest-mismatch, got %v", err)
	}
	shadow.BaselineSHA256 = baseline.Source.SHA256
	shadow.Mode = "execute"
	if err := validateEnvelope(baseline, shadow); err == nil || err.Error() != "invalid-shadow-mode" {
		t.Fatalf("expected invalid-shadow-mode, got %v", err)
	}
}

func TestShadowDriftFailsClosed(t *testing.T) {
	baseline, err := loadFixture("baseline.json")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		want       string
		mutateCase func(*shadowCase)
	}{
		{name: "input", want: "input-mismatch", mutateCase: func(item *shadowCase) { item.Input.StartRIR = "ripe" }},
		{name: "result", want: "result-mismatch", mutateCase: func(item *shadowCase) { item.Result.Host = "unknown" }},
		{name: "failure classification", want: "failure-classification-mismatch", mutateCase: func(item *shadowCase) { item.FailureClassification = "no-authoritative-evidence" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shadow := fixture{SchemaVersion: baseline.SchemaVersion, Cases: append([]shadowCase(nil), baseline.Cases...)}
			test.mutateCase(&shadow.Cases[3])
			err := compareFixtures(baseline, shadow)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected %s, got %v", test.want, err)
			}
		})
	}
}
