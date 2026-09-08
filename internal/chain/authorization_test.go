package chain

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func validAuthorizationHash(seed byte) string {
	return "sha256:" + strings.Repeat(string(seed), 64)
}

func validAuthorizationGrant() AuthorizationGrant {
	amount := int64(500)
	return AuthorizationGrant{
		RecordID:        "grant-r1",
		Kind:            "grant",
		IssuedAt:        "2026-09-08T10:00:00.000Z",
		IssuedBy:        evidence.Actor{Type: "operator", ID: "alice"},
		AuthorizationID: "auth-tool-build",
		RunID:           "run-one",
		RunManifestHash: validAuthorizationHash('a'),
		PolicyHash:      validAuthorizationHash('b'),
		Scope: AuthorizationScope{
			TaskIDs:   []string{"task-build"},
			StepIDs:   []string{"step-build"},
			ToolIDs:   []string{"tool-go-build"},
			TargetIDs: []string{"target-src"},
			Network: []AuthorizationNetworkScope{
				{Protocol: "https", Host: "example.com", Port: 443, Purpose: "fetch-module"},
			},
			EffectClasses: []string{"local-discardable"},
			Budget: AuthorizationBudgetLimit{
				ModelCalls: 10, Tokens: 10000, WallClockMs: 60000, Attempts: 3,
				Currency: "USD", AmountMicros: &amount,
			},
		},
		ExpiresAt: "2026-09-09T10:00:00.000Z",
		Evidence:  []string{validAuthorizationHash('c')},
	}
}

func validAuthorizationRevocation(grantRecord AuthorizationRecord) AuthorizationRevocation {
	return AuthorizationRevocation{
		RecordID:             "rev-r1",
		Kind:                 "revocation",
		IssuedAt:             "2026-09-08T12:00:00.000Z",
		IssuedBy:             evidence.Actor{Type: "operator", ID: "bob"},
		AuthorizationID:      "auth-tool-build",
		AuthorizationHash:    grantRecord.RecordHash,
		ReasonCode:           "policy-change",
		StopDisposition:      "completed",
		StopEvidence:         []string{validAuthorizationHash('d')},
		ResidualRiskEvidence: []string{},
	}
}

func fixedAuthorizationClock() Clock {
	return func() time.Time { return time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC) }
}

func fixedAuthorizationIDs() IDSource {
	count := 0
	return func() string {
		count++
		return "record-" + string(rune('a'+count-1))
	}
}

func TestAuthorizationGrantRecordRoundTrip(t *testing.T) {
	record, err := NewGrantAuthorizationRecord(validAuthorizationGrant(), fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeAuthorizationRecord(wire)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAuthorizationRecord(decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Grant.AuthorizationID != "auth-tool-build" || decoded.Grant.IssuedBy.ID != "alice" {
		t.Fatalf("decoded grant mismatch: %+v", decoded.Grant)
	}
	if decoded.RecordHash != record.RecordHash {
		t.Fatalf("record hash drift: got %q want %q", decoded.RecordHash, record.RecordHash)
	}
	second, err := NewGrantAuthorizationRecord(validAuthorizationGrant(), fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	if second.RecordHash != record.RecordHash {
		t.Fatalf("same input must produce stable hash: %q vs %q", second.RecordHash, record.RecordHash)
	}
}

func TestAuthorizationRecordWireShape(t *testing.T) {
	record, err := NewGrantAuthorizationRecord(validAuthorizationGrant(), fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	var shape map[string]any
	if err := json.Unmarshal(wire, &shape); err != nil {
		t.Fatal(err)
	}
	if shape["schemaVersion"] != "1.0.0" {
		t.Fatalf("schemaVersion mismatch: %v", shape["schemaVersion"])
	}
	recordBody, ok := shape["record"].(map[string]any)
	if !ok || recordBody["kind"] != "grant" {
		t.Fatalf("record body mismatch: %v", shape["record"])
	}
	if _, ok := shape["recordHash"].(string); !ok {
		t.Fatalf("recordHash missing: %v", shape)
	}
}

func TestAuthorizationGrantRejectsSelfIssuedByAgent(t *testing.T) {
	grant := validAuthorizationGrant()
	grant.IssuedBy = evidence.Actor{Type: "agent", ID: "model-x"}
	if _, err := NewGrantAuthorizationRecord(grant, fixedAuthorizationClock(), fixedAuthorizationIDs()); !errors.Is(err, ErrInvalidAuthorization) {
		t.Fatalf("expected agent-issued grant rejection, got %v", err)
	}
}

func TestAuthorizationGrantRejectsInvalidFields(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*AuthorizationGrant)
	}{
		{"bad network host", func(g *AuthorizationGrant) { g.Scope.Network[0].Host = "bad_host!" }},
		{"bad network port", func(g *AuthorizationGrant) { g.Scope.Network[0].Port = 0 }},
		{"unknown protocol", func(g *AuthorizationGrant) { g.Scope.Network[0].Protocol = "ftp" }},
		{"unknown effect class", func(g *AuthorizationGrant) { g.Scope.EffectClasses = []string{"external-delete"} }},
		{"duplicate effect class", func(g *AuthorizationGrant) { g.Scope.EffectClasses = []string{"read-only", "read-only"} }},
		{"empty evidence", func(g *AuthorizationGrant) { g.Evidence = []string{} }},
		{"bad currency", func(g *AuthorizationGrant) { g.Scope.Budget.Currency = "usd" }},
		{"zero wall clock", func(g *AuthorizationGrant) { g.Scope.Budget.WallClockMs = 0 }},
		{"negative amount", func(g *AuthorizationGrant) {
			negative := int64(-1)
			g.Scope.Budget.AmountMicros = &negative
		}},
		{"duplicate task id", func(g *AuthorizationGrant) { g.Scope.TaskIDs = []string{"task-a", "task-a"} }},
		{"expiry not after issue", func(g *AuthorizationGrant) { g.ExpiresAt = g.IssuedAt }},
		{"bad timestamp", func(g *AuthorizationGrant) { g.ExpiresAt = "2026-09-09 10:00:00Z" }},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			grant := validAuthorizationGrant()
			test.mutate(&grant)
			if _, err := NewGrantAuthorizationRecord(grant, fixedAuthorizationClock(), fixedAuthorizationIDs()); !errors.Is(err, ErrInvalidAuthorization) {
				t.Fatalf("expected rejection, got %v", err)
			}
		})
	}
}

func TestAuthorizationRevocationConditionalEvidence(t *testing.T) {
	grantRecord, err := NewGrantAuthorizationRecord(validAuthorizationGrant(), fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name   string
		mutate func(*AuthorizationRevocation)
		ok     bool
	}{
		{"completed with stop evidence", func(r *AuthorizationRevocation) {}, true},
		{"not-required with stop evidence rejected", func(r *AuthorizationRevocation) {
			r.StopDisposition = "not-required"
			r.StopEvidence = []string{validAuthorizationHash('e')}
		}, false},
		{"requested without stop evidence rejected", func(r *AuthorizationRevocation) {
			r.StopDisposition = "requested"
			r.StopEvidence = []string{}
			r.ResidualRiskEvidence = []string{validAuthorizationHash('f')}
		}, false},
		{"requested without residual risk rejected", func(r *AuthorizationRevocation) {
			r.StopDisposition = "requested"
			r.StopEvidence = []string{validAuthorizationHash('e')}
			r.ResidualRiskEvidence = []string{}
		}, false},
		{"uncertain with both accepted", func(r *AuthorizationRevocation) {
			r.StopDisposition = "uncertain"
			r.StopEvidence = []string{validAuthorizationHash('e')}
			r.ResidualRiskEvidence = []string{validAuthorizationHash('f')}
		}, true},
		{"unknown disposition rejected", func(r *AuthorizationRevocation) {
			r.StopDisposition = "teleported"
		}, false},
		{"bad reason code rejected", func(r *AuthorizationRevocation) { r.ReasonCode = "Bad Code" }, false},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			revocation := validAuthorizationRevocation(grantRecord)
			test.mutate(&revocation)
			_, err := NewRevocationAuthorizationRecord(revocation, fixedAuthorizationClock(), fixedAuthorizationIDs())
			if test.ok && err != nil {
				t.Fatalf("expected acceptance, got %v", err)
			}
			if !test.ok && !errors.Is(err, ErrInvalidAuthorization) {
				t.Fatalf("expected rejection, got %v", err)
			}
		})
	}
}

func TestDecodeAuthorizationRecordRejectsUnknownField(t *testing.T) {
	record, err := NewGrantAuthorizationRecord(validAuthorizationGrant(), fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	tampered := strings.Replace(string(wire), "\"recordId\":\"grant-r1\"", "\"recordId\":\"grant-r1\",\"secretBypass\":true", 1)
	if _, err := DecodeAuthorizationRecord([]byte(tampered)); err == nil {
		t.Fatal("expected unknown field rejection")
	}
}

func TestValidateAuthorizationRecordRejectsGrantHashMismatch(t *testing.T) {
	record, err := NewGrantAuthorizationRecord(validAuthorizationGrant(), fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	tampered := strings.Replace(string(wire), "\"expiresAt\":\"2026-09-09T10:00:00.000Z\"", "\"expiresAt\":\"2026-09-10T10:00:00.000Z\"", 1)
	decoded, err := DecodeAuthorizationRecord([]byte(tampered))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAuthorizationRecord(decoded); !errors.Is(err, ErrInvalidAuthorization) {
		t.Fatalf("expected grant hash mismatch rejection, got %v", err)
	}
}

func TestValidateAuthorizationRecordRejectsRevocationHashMismatch(t *testing.T) {
	grantRecord, err := NewGrantAuthorizationRecord(validAuthorizationGrant(), fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	revocationRecord, err := NewRevocationAuthorizationRecord(validAuthorizationRevocation(grantRecord), fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(revocationRecord)
	if err != nil {
		t.Fatal(err)
	}
	tampered := strings.Replace(string(wire), "\"reasonCode\":\"policy-change\"", "\"reasonCode\":\"policy-shift\"", 1)
	decoded, err := DecodeAuthorizationRecord([]byte(tampered))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAuthorizationRecord(decoded); !errors.Is(err, ErrInvalidAuthorization) {
		t.Fatalf("expected revocation hash mismatch rejection, got %v", err)
	}
}

func TestEvaluateGrantLifecycle(t *testing.T) {
	now := time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)
	grantRecord, err := NewGrantAuthorizationRecord(validAuthorizationGrant(), fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	ledger := AuthorizationLedger{Grants: []AuthorizationRecord{grantRecord}}
	if status, err := EvaluateGrant("auth-tool-build", ledger, now); err != nil || status != GrantActive {
		t.Fatalf("expected active, got %q err=%v", status, err)
	}
	if status, err := EvaluateGrant("auth-unknown", ledger, now); err != nil || status != GrantMissing {
		t.Fatalf("expected missing, got %q err=%v", status, err)
	}
	early := time.Date(2026, 9, 8, 9, 0, 0, 0, time.UTC)
	if status, err := EvaluateGrant("auth-tool-build", ledger, early); err != nil || status != GrantPending {
		t.Fatalf("expected pending before issuedAt, got %q err=%v", status, err)
	}
	late := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	if status, err := EvaluateGrant("auth-tool-build", ledger, late); err != nil || status != GrantExpired {
		t.Fatalf("expected expired after expiresAt, got %q err=%v", status, err)
	}
}

func TestEvaluateGrantRevocationPermanentAndHashMismatchFailsClosed(t *testing.T) {
	now := time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)
	grantRecord, err := NewGrantAuthorizationRecord(validAuthorizationGrant(), fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	revocation := validAuthorizationRevocation(grantRecord)
	revRecord, err := NewRevocationAuthorizationRecord(revocation, fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	ledger := AuthorizationLedger{Grants: []AuthorizationRecord{grantRecord}, Revocations: []AuthorizationRecord{revRecord}}
	if status, err := EvaluateGrant("auth-tool-build", ledger, now); err != nil || status != GrantRevoked {
		t.Fatalf("expected revoked, got %q err=%v", status, err)
	}
	// A later grant with the same authorizationId must not resurrect the revoked id.
	laterGrant := validAuthorizationGrant()
	laterGrant.RecordID = "grant-r2"
	laterRecord, err := NewGrantAuthorizationRecord(laterGrant, fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	ledger.Grants = append(ledger.Grants, laterRecord)
	if status, err := EvaluateGrant("auth-tool-build", ledger, now); err != nil || status != GrantRevoked {
		t.Fatalf("revocation must stay permanent, got %q err=%v", status, err)
	}
	// Hash mismatch fails closed.
	badRevocation := validAuthorizationRevocation(grantRecord)
	badRevocation.AuthorizationHash = validAuthorizationHash('9')
	badRecord, err := NewRevocationAuthorizationRecord(badRevocation, fixedAuthorizationClock(), fixedAuthorizationIDs())
	if err != nil {
		t.Fatal(err)
	}
	badLedger := AuthorizationLedger{Grants: []AuthorizationRecord{grantRecord}, Revocations: []AuthorizationRecord{badRecord}}
	if _, err := EvaluateGrant("auth-tool-build", badLedger, now); err == nil {
		t.Fatal("expected hash mismatch to fail closed")
	}
}

type authorizationTestStopper struct {
	calls    int
	evidence []string
	err      error
}

func (stopper *authorizationTestStopper) Stop(context.Context, string) ([]string, error) {
	stopper.calls++
	return stopper.evidence, stopper.err
}

func TestRequestControlledStopWiring(t *testing.T) {
	ctx := context.Background()

	notRequired := validAuthorizationRevocation(AuthorizationRecord{})
	notRequired.StopDisposition = "not-required"
	notRequired.StopEvidence = []string{}
	notRequired.ResidualRiskEvidence = []string{}
	stopper := &authorizationTestStopper{}
	result, err := RequestControlledStop(ctx, notRequired, stopper, "run-one")
	if err != nil || result.Disposition != "not-required" || stopper.calls != 0 {
		t.Fatalf("not-required must not call stopper: result=%+v err=%v calls=%d", result, err, stopper.calls)
	}

	requested := validAuthorizationRevocation(AuthorizationRecord{})
	requested.StopDisposition = "requested"
	requested.StopEvidence = []string{validAuthorizationHash('d')}
	requested.ResidualRiskEvidence = []string{validAuthorizationHash('f')}
	stopper = &authorizationTestStopper{evidence: []string{validAuthorizationHash('e')}}
	result, err = RequestControlledStop(ctx, requested, stopper, "run-one")
	if err != nil || result.Disposition != "completed" || stopper.calls != 1 || len(result.StopEvidence) != 2 {
		t.Fatalf("requested stop should complete: result=%+v err=%v calls=%d", result, err, stopper.calls)
	}

	failing := &authorizationTestStopper{evidence: []string{validAuthorizationHash('e')}, err: errors.New("process refused to die")}
	result, err = RequestControlledStop(ctx, requested, failing, "run-one")
	if err == nil || !errors.Is(err, ErrRecoveryUncertain) || result.Disposition != "uncertain" {
		t.Fatalf("failed stop must stay uncertain: result=%+v err=%v", result, err)
	}
}
