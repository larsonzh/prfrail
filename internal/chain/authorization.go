package chain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	authorizationRecordDomain    = "proofrail:authorization-record:1\n"
	authorizationTimestampLayout = "2006-01-02T15:04:05.000Z"
	authorizationSchemaVersion   = "1.0.0"
)

var (
	ErrInvalidAuthorization = errors.New("invalid authorization record")

	authorizationNetworkHostPattern = regexp.MustCompile(`^[A-Za-z0-9.-]+$`)
	authorizationCurrencyPattern    = regexp.MustCompile(`^[A-Z]{3}$`)

	validAuthorizationEffectClasses = map[string]bool{
		"read-only":         true,
		"local-discardable": true,
		"external-write":    true,
	}
	validStopDispositions = map[string]bool{
		"not-required": true,
		"requested":    true,
		"completed":    true,
		"uncertain":    true,
	}
	validNetworkProtocols = map[string]bool{
		"http": true, "https": true, "ssh": true, "custom": true,
	}
)

// AuthorizationNetworkScope mirrors networkScope in authorization-record.schema.json.
type AuthorizationNetworkScope struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Purpose  string `json:"purpose"`
}

// AuthorizationBudgetLimit mirrors budgetLimit in authorization-record.schema.json.
type AuthorizationBudgetLimit struct {
	ModelCalls   int    `json:"modelCalls"`
	Tokens       int    `json:"tokens"`
	WallClockMs  int    `json:"wallClockMs"`
	Attempts     int    `json:"attempts"`
	Currency     string `json:"currency"`
	AmountMicros *int64 `json:"amountMicros"`
}

// AuthorizationScope mirrors scope in authorization-record.schema.json.
type AuthorizationScope struct {
	TaskIDs       []string                    `json:"taskIds"`
	StepIDs       []string                    `json:"stepIds"`
	ToolIDs       []string                    `json:"toolIds"`
	TargetIDs     []string                    `json:"targetIds"`
	Network       []AuthorizationNetworkScope `json:"network"`
	EffectClasses []string                    `json:"effectClasses"`
	Budget        AuthorizationBudgetLimit    `json:"budget"`
}

// AuthorizationGrant is the grant body of an authorization record.
type AuthorizationGrant struct {
	RecordID        string             `json:"recordId"`
	Kind            string             `json:"kind"`
	IssuedAt        string             `json:"issuedAt"`
	IssuedBy        evidence.Actor     `json:"issuedBy"`
	AuthorizationID string             `json:"authorizationId"`
	RunID           string             `json:"runId"`
	RunManifestHash string             `json:"runManifestHash"`
	PolicyHash      string             `json:"policyHash"`
	Scope           AuthorizationScope `json:"scope"`
	ExpiresAt       string             `json:"expiresAt"`
	Evidence        []string           `json:"evidence"`
}

// AuthorizationRevocation is the revocation body of an authorization record.
type AuthorizationRevocation struct {
	RecordID             string         `json:"recordId"`
	Kind                 string         `json:"kind"`
	IssuedAt             string         `json:"issuedAt"`
	IssuedBy             evidence.Actor `json:"issuedBy"`
	AuthorizationID      string         `json:"authorizationId"`
	AuthorizationHash    string         `json:"authorizationHash"`
	ReasonCode           string         `json:"reasonCode"`
	StopDisposition      string         `json:"stopDisposition"`
	StopEvidence         []string       `json:"stopEvidence"`
	ResidualRiskEvidence []string       `json:"residualRiskEvidence"`
}

// AuthorizationRecord wraps either a grant or a revocation with a schema
// version and the canonical recordHash that binds the body.
type AuthorizationRecord struct {
	SchemaVersion string                   `json:"schemaVersion"`
	Grant         *AuthorizationGrant      `json:"-"`
	Revocation    *AuthorizationRevocation `json:"-"`
	RecordHash    string                   `json:"recordHash"`
}

// MarshalJSON serializes the wire shape {schemaVersion, record, recordHash}.
func (record AuthorizationRecord) MarshalJSON() ([]byte, error) {
	var body any
	switch {
	case record.Grant != nil && record.Revocation == nil:
		body = record.Grant
	case record.Revocation != nil && record.Grant == nil:
		body = record.Revocation
	default:
		return nil, fmt.Errorf("%w: exactly one of grant/revocation required", ErrInvalidAuthorization)
	}
	return json.Marshal(struct {
		SchemaVersion string `json:"schemaVersion"`
		Record        any    `json:"record"`
		RecordHash    string `json:"recordHash"`
	}{record.SchemaVersion, body, record.RecordHash})
}

// UnmarshalJSON strictly decodes the wire shape and rejects unknown or
// duplicate fields in the record body.
func (record *AuthorizationRecord) UnmarshalJSON(input []byte) error {
	var wire struct {
		SchemaVersion string          `json:"schemaVersion"`
		Record        json.RawMessage `json:"record"`
		RecordHash    string          `json:"recordHash"`
	}
	if err := evidence.DecodeStrictJSON(input, &wire); err != nil {
		return err
	}
	var kind struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(wire.Record, &kind); err != nil {
		return err
	}
	record.SchemaVersion = wire.SchemaVersion
	record.RecordHash = wire.RecordHash
	switch kind.Kind {
	case "grant":
		var grant AuthorizationGrant
		if err := evidence.DecodeStrictJSON(wire.Record, &grant); err != nil {
			return err
		}
		record.Grant = &grant
	case "revocation":
		var revocation AuthorizationRevocation
		if err := evidence.DecodeStrictJSON(wire.Record, &revocation); err != nil {
			return err
		}
		record.Revocation = &revocation
	default:
		return fmt.Errorf("%w: unknown record kind %q", ErrInvalidAuthorization, kind.Kind)
	}
	return nil
}

// NewGrantAuthorizationRecord builds a grant record after validating its body.
func NewGrantAuthorizationRecord(grant AuthorizationGrant, clock Clock, ids IDSource) (AuthorizationRecord, error) {
	if clock == nil {
		clock = defaultClock
	}
	if ids == nil {
		ids = randomID
	}
	if grant.Kind == "" {
		grant.Kind = "grant"
	}
	if grant.RecordID == "" {
		grant.RecordID = ids()
	}
	if grant.IssuedAt == "" {
		grant.IssuedAt = clock().UTC().Format(authorizationTimestampLayout)
	}
	if err := validateAuthorizationGrant(grant); err != nil {
		return AuthorizationRecord{}, err
	}
	hash, err := digestRecord(authorizationRecordDomain, grant)
	if err != nil {
		return AuthorizationRecord{}, err
	}
	return AuthorizationRecord{SchemaVersion: authorizationSchemaVersion, Grant: &grant, RecordHash: hash}, nil
}

// NewRevocationAuthorizationRecord builds a revocation record after validating its body.
func NewRevocationAuthorizationRecord(revocation AuthorizationRevocation, clock Clock, ids IDSource) (AuthorizationRecord, error) {
	if clock == nil {
		clock = defaultClock
	}
	if ids == nil {
		ids = randomID
	}
	if revocation.Kind == "" {
		revocation.Kind = "revocation"
	}
	if revocation.RecordID == "" {
		revocation.RecordID = ids()
	}
	if revocation.IssuedAt == "" {
		revocation.IssuedAt = clock().UTC().Format(authorizationTimestampLayout)
	}
	if err := validateAuthorizationRevocation(revocation); err != nil {
		return AuthorizationRecord{}, err
	}
	hash, err := digestRecord(authorizationRecordDomain, revocation)
	if err != nil {
		return AuthorizationRecord{}, err
	}
	return AuthorizationRecord{SchemaVersion: authorizationSchemaVersion, Revocation: &revocation, RecordHash: hash}, nil
}

// DecodeAuthorizationRecord strictly decodes one wire record.
func DecodeAuthorizationRecord(input []byte) (AuthorizationRecord, error) {
	var record AuthorizationRecord
	if err := json.Unmarshal(input, &record); err != nil {
		return AuthorizationRecord{}, err
	}
	return record, nil
}

// ValidateAuthorizationRecord enforces the schema contract for a decoded record.
func ValidateAuthorizationRecord(record AuthorizationRecord) error {
	if record.SchemaVersion != authorizationSchemaVersion {
		return fmt.Errorf("%w: unsupported schemaVersion %q", ErrInvalidAuthorization, record.SchemaVersion)
	}
	if !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: invalid recordHash", ErrInvalidAuthorization)
	}
	var (
		expected string
		err      error
	)
	switch {
	case record.Grant != nil && record.Revocation == nil:
		if err := validateAuthorizationGrant(*record.Grant); err != nil {
			return err
		}
		expected, err = digestRecord(authorizationRecordDomain, *record.Grant)
	case record.Revocation != nil && record.Grant == nil:
		if err := validateAuthorizationRevocation(*record.Revocation); err != nil {
			return err
		}
		expected, err = digestRecord(authorizationRecordDomain, *record.Revocation)
	default:
		return fmt.Errorf("%w: exactly one of grant/revocation required", ErrInvalidAuthorization)
	}
	if err != nil {
		return err
	}
	if record.RecordHash != expected {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAuthorization)
	}
	return nil
}

// GrantStatus reports the outcome of evaluating a grant against the ledger.
type GrantStatus string

const (
	GrantMissing GrantStatus = "missing"
	GrantPending GrantStatus = "pending"
	GrantActive  GrantStatus = "active"
	GrantExpired GrantStatus = "expired"
	GrantRevoked GrantStatus = "revoked"
)

// AuthorizationLedger is an append-only collection of grant and revocation records.
type AuthorizationLedger struct {
	Grants      []AuthorizationRecord
	Revocations []AuthorizationRecord
}

// GrantRecordsFor returns every grant record for an authorizationId in append order.
func (ledger AuthorizationLedger) GrantRecordsFor(authorizationID string) []AuthorizationRecord {
	var records []AuthorizationRecord
	for _, record := range ledger.Grants {
		if record.Grant != nil && record.Grant.AuthorizationID == authorizationID {
			records = append(records, record)
		}
	}
	return records
}

// LatestGrant returns the newest grant record for an authorizationId.
func (ledger AuthorizationLedger) LatestGrant(authorizationID string) (AuthorizationRecord, *AuthorizationGrant, bool) {
	records := ledger.GrantRecordsFor(authorizationID)
	if len(records) == 0 {
		return AuthorizationRecord{}, nil, false
	}
	record := records[len(records)-1]
	return record, record.Grant, true
}

// LatestRevocation returns the newest revocation record for an authorizationId.
func (ledger AuthorizationLedger) LatestRevocation(authorizationID string) (AuthorizationRecord, *AuthorizationRevocation, bool) {
	var foundRecord AuthorizationRecord
	var found *AuthorizationRevocation
	for _, record := range ledger.Revocations {
		if record.Revocation != nil && record.Revocation.AuthorizationID == authorizationID {
			foundRecord = record
			found = record.Revocation
		}
	}
	return foundRecord, found, found != nil
}

// EvaluateGrant decides whether the grant for authorizationID still permits
// further side effects and acceptance at the given instant. A revocation whose
// authorizationHash does not match the grant record fails closed; once revoked
// an authorizationId never returns to active, and expiry blocks automatically.
func EvaluateGrant(authorizationID string, ledger AuthorizationLedger, now time.Time) (GrantStatus, error) {
	grantRecords := ledger.GrantRecordsFor(authorizationID)
	if len(grantRecords) == 0 {
		return GrantMissing, nil
	}
	grant := grantRecords[len(grantRecords)-1].Grant
	if _, revocation, revoked := ledger.LatestRevocation(authorizationID); revoked {
		matched := false
		for _, record := range grantRecords {
			if record.RecordHash == revocation.AuthorizationHash {
				matched = true
				break
			}
		}
		if !matched {
			return "", fmt.Errorf("%w: revocation hash mismatch for authorization %q", ErrInvalidState, authorizationID)
		}
		return GrantRevoked, nil
	}
	issuedAt, err := time.Parse(authorizationTimestampLayout, grant.IssuedAt)
	if err != nil {
		return "", fmt.Errorf("%w: grant issuedAt for %q: %v", ErrInvalidState, authorizationID, err)
	}
	if now.Before(issuedAt) {
		return GrantPending, nil
	}
	expiresAt, err := time.Parse(authorizationTimestampLayout, grant.ExpiresAt)
	if err != nil {
		return "", fmt.Errorf("%w: grant expiresAt for %q: %v", ErrInvalidState, authorizationID, err)
	}
	if !now.Before(expiresAt) {
		return GrantExpired, nil
	}
	return GrantActive, nil
}

// StopRequestResult describes the outcome of acting on a revocation's stop disposition.
type StopRequestResult struct {
	Disposition          string   `json:"disposition"`
	StopEvidence         []string `json:"stopEvidence"`
	ResidualRiskEvidence []string `json:"residualRiskEvidence"`
}

// RequestControlledStop wires a revocation's stop request to the Stopper
// (guard shutdown). A failed in-flight stop is never reported as completed:
// it stays uncertain with residual-risk evidence so the run remains paused.
func RequestControlledStop(ctx context.Context, revocation AuthorizationRevocation, stopper Stopper, runID string) (StopRequestResult, error) {
	if stopper == nil {
		return StopRequestResult{}, fmt.Errorf("%w: nil stopper", ErrInvalidState)
	}
	switch revocation.StopDisposition {
	case "not-required":
		return StopRequestResult{Disposition: "not-required"}, nil
	case "completed":
		return StopRequestResult{
			Disposition:          "completed",
			StopEvidence:         append([]string{}, revocation.StopEvidence...),
			ResidualRiskEvidence: append([]string{}, revocation.ResidualRiskEvidence...),
		}, nil
	case "uncertain":
		return StopRequestResult{
			Disposition:          "uncertain",
			StopEvidence:         append([]string{}, revocation.StopEvidence...),
			ResidualRiskEvidence: append([]string{}, revocation.ResidualRiskEvidence...),
		}, nil
	case "requested":
		evidenceHashes, err := stopper.Stop(ctx, runID)
		if err != nil {
			return StopRequestResult{
				Disposition:          "uncertain",
				StopEvidence:         uniqueHashes(append(append([]string{}, revocation.StopEvidence...), evidenceHashes...)),
				ResidualRiskEvidence: append([]string{}, revocation.ResidualRiskEvidence...),
			}, fmt.Errorf("%w: controlled stop uncertain: %v", ErrRecoveryUncertain, err)
		}
		return StopRequestResult{
			Disposition:          "completed",
			StopEvidence:         uniqueHashes(append(append([]string{}, revocation.StopEvidence...), evidenceHashes...)),
			ResidualRiskEvidence: append([]string{}, revocation.ResidualRiskEvidence...),
		}, nil
	default:
		return StopRequestResult{}, fmt.Errorf("%w: unknown stop disposition %q", ErrInvalidState, revocation.StopDisposition)
	}
}

func validAuthorizationID(value string) bool {
	return evidence.ValidID(value) && len(value) <= 64
}

func parseAuthorizationTimestamp(value string) (time.Time, error) {
	return time.Parse(authorizationTimestampLayout, value)
}

func validateAuthorizationActor(actor evidence.Actor) error {
	if actor.Type != "operator" && actor.Type != "policy" {
		return fmt.Errorf("%w: actor must be operator or policy, got %q", ErrInvalidAuthorization, actor.Type)
	}
	if !validAuthorizationID(actor.ID) {
		return fmt.Errorf("%w: invalid actor id", ErrInvalidAuthorization)
	}
	return nil
}

func validateAuthorizationGrant(grant AuthorizationGrant) error {
	if grant.Kind != "grant" {
		return fmt.Errorf("%w: kind must be grant", ErrInvalidAuthorization)
	}
	if !validAuthorizationID(grant.RecordID) {
		return fmt.Errorf("%w: invalid recordId", ErrInvalidAuthorization)
	}
	if _, err := parseAuthorizationTimestamp(grant.IssuedAt); err != nil {
		return fmt.Errorf("%w: invalid issuedAt", ErrInvalidAuthorization)
	}
	if err := validateAuthorizationActor(grant.IssuedBy); err != nil {
		return err
	}
	if !validAuthorizationID(grant.AuthorizationID) {
		return fmt.Errorf("%w: invalid authorizationId", ErrInvalidAuthorization)
	}
	if !evidence.ValidID(grant.RunID) {
		return fmt.Errorf("%w: invalid runId", ErrInvalidAuthorization)
	}
	if !evidence.ValidHash(grant.RunManifestHash) {
		return fmt.Errorf("%w: invalid runManifestHash", ErrInvalidAuthorization)
	}
	if !evidence.ValidHash(grant.PolicyHash) {
		return fmt.Errorf("%w: invalid policyHash", ErrInvalidAuthorization)
	}
	if err := validateAuthorizationScope(grant.Scope); err != nil {
		return err
	}
	issuedAt, _ := parseAuthorizationTimestamp(grant.IssuedAt)
	expiresAt, err := parseAuthorizationTimestamp(grant.ExpiresAt)
	if err != nil {
		return fmt.Errorf("%w: invalid expiresAt", ErrInvalidAuthorization)
	}
	if !expiresAt.After(issuedAt) {
		return fmt.Errorf("%w: expiresAt must be after issuedAt", ErrInvalidAuthorization)
	}
	if len(grant.Evidence) < 1 {
		return fmt.Errorf("%w: grant evidence must not be empty", ErrInvalidAuthorization)
	}
	if !allHashes(grant.Evidence) {
		return fmt.Errorf("%w: invalid grant evidence", ErrInvalidAuthorization)
	}
	return nil
}

func validateAuthorizationScope(scope AuthorizationScope) error {
	if err := validateAuthorizationIDSet(scope.TaskIDs, "taskIds"); err != nil {
		return err
	}
	if err := validateAuthorizationIDSet(scope.StepIDs, "stepIds"); err != nil {
		return err
	}
	if err := validateAuthorizationIDSet(scope.ToolIDs, "toolIds"); err != nil {
		return err
	}
	if err := validateAuthorizationIDSet(scope.TargetIDs, "targetIds"); err != nil {
		return err
	}
	for _, network := range scope.Network {
		if !validNetworkProtocols[network.Protocol] {
			return fmt.Errorf("%w: unknown network protocol %q", ErrInvalidAuthorization, network.Protocol)
		}
		if len(network.Host) < 1 || len(network.Host) > 253 || !authorizationNetworkHostPattern.MatchString(network.Host) {
			return fmt.Errorf("%w: invalid network host %q", ErrInvalidAuthorization, network.Host)
		}
		if network.Port < 1 || network.Port > 65535 {
			return fmt.Errorf("%w: invalid network port %d", ErrInvalidAuthorization, network.Port)
		}
		if !validAuthorizationID(network.Purpose) {
			return fmt.Errorf("%w: invalid network purpose", ErrInvalidAuthorization)
		}
	}
	seenEffects := map[string]bool{}
	for _, class := range scope.EffectClasses {
		if !validAuthorizationEffectClasses[class] {
			return fmt.Errorf("%w: unknown effect class %q", ErrInvalidAuthorization, class)
		}
		if seenEffects[class] {
			return fmt.Errorf("%w: duplicate effect class %q", ErrInvalidAuthorization, class)
		}
		seenEffects[class] = true
	}
	budget := scope.Budget
	if budget.ModelCalls < 0 || budget.Tokens < 0 {
		return fmt.Errorf("%w: budget modelCalls/tokens must be >= 0", ErrInvalidAuthorization)
	}
	if budget.WallClockMs < 1 || budget.Attempts < 1 {
		return fmt.Errorf("%w: budget wallClockMs/attempts must be >= 1", ErrInvalidAuthorization)
	}
	if !authorizationCurrencyPattern.MatchString(budget.Currency) {
		return fmt.Errorf("%w: invalid budget currency %q", ErrInvalidAuthorization, budget.Currency)
	}
	if budget.AmountMicros != nil && *budget.AmountMicros < 0 {
		return fmt.Errorf("%w: budget amountMicros must be >= 0", ErrInvalidAuthorization)
	}
	return nil
}

func validateAuthorizationIDSet(values []string, field string) error {
	seen := map[string]bool{}
	for _, value := range values {
		if !validAuthorizationID(value) {
			return fmt.Errorf("%w: invalid %s id %q", ErrInvalidAuthorization, field, value)
		}
		if seen[value] {
			return fmt.Errorf("%w: duplicate %s id %q", ErrInvalidAuthorization, field, value)
		}
		seen[value] = true
	}
	return nil
}

func validateAuthorizationRevocation(revocation AuthorizationRevocation) error {
	if revocation.Kind != "revocation" {
		return fmt.Errorf("%w: kind must be revocation", ErrInvalidAuthorization)
	}
	if !validAuthorizationID(revocation.RecordID) {
		return fmt.Errorf("%w: invalid recordId", ErrInvalidAuthorization)
	}
	if _, err := parseAuthorizationTimestamp(revocation.IssuedAt); err != nil {
		return fmt.Errorf("%w: invalid issuedAt", ErrInvalidAuthorization)
	}
	if err := validateAuthorizationActor(revocation.IssuedBy); err != nil {
		return err
	}
	if !validAuthorizationID(revocation.AuthorizationID) {
		return fmt.Errorf("%w: invalid authorizationId", ErrInvalidAuthorization)
	}
	if !evidence.ValidHash(revocation.AuthorizationHash) {
		return fmt.Errorf("%w: invalid authorizationHash", ErrInvalidAuthorization)
	}
	if !validAuthorizationID(revocation.ReasonCode) {
		return fmt.Errorf("%w: invalid reasonCode", ErrInvalidAuthorization)
	}
	if !validStopDispositions[revocation.StopDisposition] {
		return fmt.Errorf("%w: unknown stopDisposition %q", ErrInvalidAuthorization, revocation.StopDisposition)
	}
	if !allHashes(revocation.StopEvidence) || !allHashes(revocation.ResidualRiskEvidence) {
		return fmt.Errorf("%w: invalid stop/residual-risk evidence", ErrInvalidAuthorization)
	}
	switch revocation.StopDisposition {
	case "not-required":
		if len(revocation.StopEvidence) != 0 {
			return fmt.Errorf("%w: not-required must not carry stop evidence", ErrInvalidAuthorization)
		}
	default:
		if len(revocation.StopEvidence) < 1 {
			return fmt.Errorf("%w: %s requires stop evidence", ErrInvalidAuthorization, revocation.StopDisposition)
		}
	}
	if (revocation.StopDisposition == "requested" || revocation.StopDisposition == "uncertain") && len(revocation.ResidualRiskEvidence) < 1 {
		return fmt.Errorf("%w: %s requires residual-risk evidence", ErrInvalidAuthorization, revocation.StopDisposition)
	}
	return nil
}
