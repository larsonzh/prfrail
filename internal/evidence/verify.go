package evidence

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
)

var ErrObjectNotFound = errors.New("evidence object not found")

type ObjectReader interface {
	ReadObject(context.Context, string) ([]byte, error)
}

type VerificationStatus string

const (
	VerificationPassed     VerificationStatus = "passed"
	VerificationFailed     VerificationStatus = "failed"
	VerificationIncomplete VerificationStatus = "incomplete"
)

type VerificationCheck struct {
	CheckID   string
	Status    VerificationStatus
	Evidence  []string
	ErrorCode string
}

type VerificationResult struct {
	Outcome VerificationStatus
	Checks  []VerificationCheck
}

func VerifyManifestObjects(ctx context.Context, manifest EvidenceManifest, reader ObjectReader) VerificationResult {
	if err := VerifyEvidenceManifest(manifest); err != nil {
		return result(VerificationCheck{CheckID: "reference-closure", Status: VerificationFailed, ErrorCode: "schema.invalid-document"})
	}
	if reader == nil {
		return result(VerificationCheck{CheckID: "reference-closure", Status: VerificationIncomplete, ErrorCode: "integrity.reference-missing"})
	}
	verified := make([]string, 0, len(manifest.Manifest.Items))
	var events []StateEvent
	attemptStatus := VerificationPassed
	attemptCode := ""
	eventStatus := VerificationPassed
	eventCode := ""
	for _, item := range manifest.Manifest.Items {
		data, err := reader.ReadObject(ctx, item.ObjectHash)
		if err != nil {
			return result(VerificationCheck{CheckID: "reference-closure", Status: VerificationIncomplete, Evidence: verified, ErrorCode: "integrity.reference-missing"})
		}
		if Digest("", data) != item.ObjectHash {
			return result(VerificationCheck{CheckID: "reference-closure", Status: VerificationFailed, Evidence: verified, ErrorCode: "integrity.hash-mismatch"})
		}
		if item.MediaType == "application/json" {
			canonical, err := Canonicalize(data)
			if err != nil || !bytes.Equal(data, canonical) {
				return result(VerificationCheck{CheckID: "reference-closure", Status: VerificationFailed, Evidence: verified, ErrorCode: "schema.invalid-document"})
			}
		}
		if item.Kind == "state-event" {
			event, err := DecodeStateEvent(data)
			if err != nil {
				eventStatus = VerificationFailed
				eventCode = "schema.invalid-document"
			} else {
				if event.Event.RunID != manifest.Manifest.RunID || (event.Event.Entity.Kind != "chain" && (event.Event.Entity.TaskID != manifest.Manifest.TaskID || event.Event.Entity.Attempt != manifest.Manifest.Attempt)) {
					attemptStatus = VerificationFailed
					attemptCode = "integrity.reference-missing"
				}
				events = append(events, event)
			}
		}
		verified = append(verified, item.ObjectHash)
	}
	if eventStatus == VerificationPassed && len(events) > 0 {
		sort.Slice(events, func(left, right int) bool {
			return events[left].Event.Sequence < events[right].Event.Sequence
		})
		if err := VerifyEventChain(events); err != nil {
			eventStatus = VerificationFailed
			eventCode = "integrity.sequence-invalid"
		}
	}
	checks := make([]VerificationCheck, 0, 4)
	if len(events) > 0 || eventStatus == VerificationFailed {
		checks = append(checks, VerificationCheck{CheckID: "event-chain", Status: eventStatus, Evidence: append([]string(nil), verified...), ErrorCode: eventCode})
	}
	checks = append(checks, VerificationCheck{CheckID: "reference-closure", Status: VerificationPassed, Evidence: append([]string(nil), verified...)})
	if len(events) > 0 {
		checks = append(checks, VerificationCheck{CheckID: "attempt-binding", Status: attemptStatus, Evidence: append([]string(nil), verified...), ErrorCode: attemptCode})
	}
	return result(checks...)
}

func result(checks ...VerificationCheck) VerificationResult {
	outcome := VerificationPassed
	for _, check := range checks {
		if check.Status == VerificationFailed {
			outcome = VerificationFailed
			break
		}
		if check.Status == VerificationIncomplete {
			outcome = VerificationIncomplete
		}
	}
	return VerificationResult{Outcome: outcome, Checks: checks}
}

func EncodeCanonical(value any) ([]byte, error) {
	canonical, err := canonicalValue(value)
	if err != nil {
		return nil, fmt.Errorf("canonical encode: %w", err)
	}
	return canonical, nil
}
