package adapters

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	TimestampLayout = "2006-01-02T15:04:05.000Z"

	adapterRequestDomain         = "proofrail:adapter-request:1\n"
	adapterClaimDomain           = "proofrail:adapter-claim:1\n"
	adapterResultDomain          = "proofrail:adapter-result:1\n"
	adapterDispatchReceiptDomain = "proofrail:adapter-dispatch-receipt:1\n"
	adapterTakeoverReceiptDomain = "proofrail:adapter-takeover-receipt:1\n"
)

var (
	ErrInvalidEnvelope  = errors.New("invalid adapter envelope")
	ErrInvalidReceipt   = errors.New("invalid adapter receipt")
	ErrRecordConflict   = errors.New("adapter record conflict")
	ErrRequestNotFound  = errors.New("adapter request not found")
	ErrClaimNotFound    = errors.New("adapter claim not found")
	ErrFencingConflict  = errors.New("adapter fencing token mismatch")
	ErrTakeoverRequired = errors.New("adapter takeover receipt required")
)

var transportRequestIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]*$`)

type RequestMessage struct {
	Type        string `json:"type"`
	RequestID   string `json:"requestId"`
	RunID       string `json:"runId"`
	TaskID      string `json:"taskId"`
	Attempt     int    `json:"attempt"`
	CreatedAt   string `json:"createdAt"`
	Adapter     string `json:"adapter"`
	ContextHash string `json:"contextHash"`
}

type ClaimMessage struct {
	Type                string  `json:"type"`
	RequestID           string  `json:"requestId"`
	ClaimID             string  `json:"claimId"`
	ConsumerID          string  `json:"consumerId"`
	ClaimedAt           string  `json:"claimedAt"`
	LeaseExpiresAt      string  `json:"leaseExpiresAt"`
	Generation          int     `json:"generation"`
	RequestHash         string  `json:"requestHash"`
	TakeoverReceiptHash *string `json:"takeoverReceiptHash"`
}

type ResultMessage struct {
	Type           string   `json:"type"`
	RequestID      string   `json:"requestId"`
	ClaimID        string   `json:"claimId"`
	Generation     int      `json:"generation"`
	CompletedAt    string   `json:"completedAt"`
	RequestHash    string   `json:"requestHash"`
	Status         string   `json:"status"`
	OutputEvidence []string `json:"outputEvidence"`
	ErrorEvidence  []string `json:"errorEvidence"`
}

type RequestEnvelope struct {
	SchemaVersion string         `json:"schemaVersion"`
	Message       RequestMessage `json:"message"`
	RecordHash    string         `json:"recordHash"`
}

type ClaimEnvelope struct {
	SchemaVersion string       `json:"schemaVersion"`
	Message       ClaimMessage `json:"message"`
	RecordHash    string       `json:"recordHash"`
}

type ResultEnvelope struct {
	SchemaVersion string        `json:"schemaVersion"`
	Message       ResultMessage `json:"message"`
	RecordHash    string        `json:"recordHash"`
}

type DispatchReceipt struct {
	Type               string         `json:"type"`
	ReceiptID          string         `json:"receiptId"`
	OccurredAt         string         `json:"occurredAt"`
	RecordedBy         evidence.Actor `json:"recordedBy"`
	RequestID          string         `json:"requestId"`
	RequestHash        string         `json:"requestHash"`
	Evidence           []string       `json:"evidence"`
	ErrorEvidence      []string       `json:"errorEvidence"`
	Adapter            string         `json:"adapter"`
	Transport          string         `json:"transport"`
	TransportRequestID *string        `json:"transportRequestId"`
	Outcome            string         `json:"outcome"`
}

type TakeoverReason struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

type TakeoverReceipt struct {
	Type                  string         `json:"type"`
	ReceiptID             string         `json:"receiptId"`
	OccurredAt            string         `json:"occurredAt"`
	RecordedBy            evidence.Actor `json:"recordedBy"`
	RequestID             string         `json:"requestId"`
	RequestHash           string         `json:"requestHash"`
	Evidence              []string       `json:"evidence"`
	ErrorEvidence         []string       `json:"errorEvidence"`
	PreviousClaimID       string         `json:"previousClaimId"`
	PreviousGeneration    int            `json:"previousGeneration"`
	NewClaimID            string         `json:"newClaimId"`
	NewGeneration         int            `json:"newGeneration"`
	OldWriterEvidence     []string       `json:"oldWriterEvidence"`
	AuthorizationEvidence []string       `json:"authorizationEvidence"`
	Reason                TakeoverReason `json:"reason"`
	Outcome               string         `json:"outcome"`
}

type DispatchReceiptRecord struct {
	SchemaVersion string          `json:"schemaVersion"`
	Receipt       DispatchReceipt `json:"receipt"`
	ReceiptHash   string          `json:"receiptHash"`
}

type TakeoverReceiptRecord struct {
	SchemaVersion string          `json:"schemaVersion"`
	Receipt       TakeoverReceipt `json:"receipt"`
	ReceiptHash   string          `json:"receiptHash"`
}

type FileQueue struct {
	root string
}

func NewRequestEnvelope(message RequestMessage) (RequestEnvelope, error) {
	if message.Type == "" {
		message.Type = "request"
	}
	if err := validateRequestMessage(message); err != nil {
		return RequestEnvelope{}, err
	}
	hash, err := digestMessage(adapterRequestDomain, message)
	if err != nil {
		return RequestEnvelope{}, err
	}
	return RequestEnvelope{SchemaVersion: evidence.SchemaVersion, Message: message, RecordHash: hash}, nil
}

func NewClaimEnvelope(message ClaimMessage) (ClaimEnvelope, error) {
	if message.Type == "" {
		message.Type = "claim"
	}
	if err := validateClaimMessage(message); err != nil {
		return ClaimEnvelope{}, err
	}
	hash, err := digestMessage(adapterClaimDomain, message)
	if err != nil {
		return ClaimEnvelope{}, err
	}
	return ClaimEnvelope{SchemaVersion: evidence.SchemaVersion, Message: message, RecordHash: hash}, nil
}

func NewResultEnvelope(message ResultMessage) (ResultEnvelope, error) {
	if message.Type == "" {
		message.Type = "result"
	}
	message.OutputEvidence = nonNilStrings(message.OutputEvidence)
	message.ErrorEvidence = nonNilStrings(message.ErrorEvidence)
	if err := validateResultMessage(message); err != nil {
		return ResultEnvelope{}, err
	}
	hash, err := digestMessage(adapterResultDomain, message)
	if err != nil {
		return ResultEnvelope{}, err
	}
	return ResultEnvelope{SchemaVersion: evidence.SchemaVersion, Message: message, RecordHash: hash}, nil
}

func DecodeRequestEnvelope(input []byte) (RequestEnvelope, error) {
	var record RequestEnvelope
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return RequestEnvelope{}, err
	}
	if record.SchemaVersion != evidence.SchemaVersion {
		return RequestEnvelope{}, fmt.Errorf("%w: unsupported schema version", ErrInvalidEnvelope)
	}
	if err := validateRequestMessage(record.Message); err != nil {
		return RequestEnvelope{}, err
	}
	expected, err := digestMessage(adapterRequestDomain, record.Message)
	if err != nil {
		return RequestEnvelope{}, err
	}
	if record.RecordHash != expected {
		return RequestEnvelope{}, fmt.Errorf("%w: request record hash mismatch", ErrInvalidEnvelope)
	}
	return record, nil
}

func DecodeClaimEnvelope(input []byte) (ClaimEnvelope, error) {
	var record ClaimEnvelope
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return ClaimEnvelope{}, err
	}
	if record.SchemaVersion != evidence.SchemaVersion {
		return ClaimEnvelope{}, fmt.Errorf("%w: unsupported schema version", ErrInvalidEnvelope)
	}
	if err := validateClaimMessage(record.Message); err != nil {
		return ClaimEnvelope{}, err
	}
	expected, err := digestMessage(adapterClaimDomain, record.Message)
	if err != nil {
		return ClaimEnvelope{}, err
	}
	if record.RecordHash != expected {
		return ClaimEnvelope{}, fmt.Errorf("%w: claim record hash mismatch", ErrInvalidEnvelope)
	}
	return record, nil
}

func DecodeResultEnvelope(input []byte) (ResultEnvelope, error) {
	var record ResultEnvelope
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return ResultEnvelope{}, err
	}
	if record.SchemaVersion != evidence.SchemaVersion {
		return ResultEnvelope{}, fmt.Errorf("%w: unsupported schema version", ErrInvalidEnvelope)
	}
	record.Message.OutputEvidence = nonNilStrings(record.Message.OutputEvidence)
	record.Message.ErrorEvidence = nonNilStrings(record.Message.ErrorEvidence)
	if err := validateResultMessage(record.Message); err != nil {
		return ResultEnvelope{}, err
	}
	expected, err := digestMessage(adapterResultDomain, record.Message)
	if err != nil {
		return ResultEnvelope{}, err
	}
	if record.RecordHash != expected {
		return ResultEnvelope{}, fmt.Errorf("%w: result record hash mismatch", ErrInvalidEnvelope)
	}
	return record, nil
}

func NewDispatchReceiptRecord(receipt DispatchReceipt) (DispatchReceiptRecord, error) {
	if receipt.Type == "" {
		receipt.Type = "dispatch"
	}
	receipt.Evidence = nonNilStrings(receipt.Evidence)
	receipt.ErrorEvidence = nonNilStrings(receipt.ErrorEvidence)
	if err := validateDispatchReceipt(receipt); err != nil {
		return DispatchReceiptRecord{}, err
	}
	hash, err := digestMessage(adapterDispatchReceiptDomain, receipt)
	if err != nil {
		return DispatchReceiptRecord{}, err
	}
	return DispatchReceiptRecord{SchemaVersion: evidence.SchemaVersion, Receipt: receipt, ReceiptHash: hash}, nil
}

func NewTakeoverReceiptRecord(receipt TakeoverReceipt) (TakeoverReceiptRecord, error) {
	if receipt.Type == "" {
		receipt.Type = "takeover"
	}
	receipt.Evidence = nonNilStrings(receipt.Evidence)
	receipt.ErrorEvidence = nonNilStrings(receipt.ErrorEvidence)
	receipt.OldWriterEvidence = nonNilStrings(receipt.OldWriterEvidence)
	receipt.AuthorizationEvidence = nonNilStrings(receipt.AuthorizationEvidence)
	if err := validateTakeoverReceipt(receipt); err != nil {
		return TakeoverReceiptRecord{}, err
	}
	hash, err := digestMessage(adapterTakeoverReceiptDomain, receipt)
	if err != nil {
		return TakeoverReceiptRecord{}, err
	}
	return TakeoverReceiptRecord{SchemaVersion: evidence.SchemaVersion, Receipt: receipt, ReceiptHash: hash}, nil
}

func DecodeDispatchReceiptRecord(input []byte) (DispatchReceiptRecord, error) {
	var record DispatchReceiptRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return DispatchReceiptRecord{}, err
	}
	if record.SchemaVersion != evidence.SchemaVersion {
		return DispatchReceiptRecord{}, fmt.Errorf("%w: unsupported schema version", ErrInvalidReceipt)
	}
	record.Receipt.Evidence = nonNilStrings(record.Receipt.Evidence)
	record.Receipt.ErrorEvidence = nonNilStrings(record.Receipt.ErrorEvidence)
	if err := validateDispatchReceipt(record.Receipt); err != nil {
		return DispatchReceiptRecord{}, err
	}
	expected, err := digestMessage(adapterDispatchReceiptDomain, record.Receipt)
	if err != nil {
		return DispatchReceiptRecord{}, err
	}
	if record.ReceiptHash != expected {
		return DispatchReceiptRecord{}, fmt.Errorf("%w: dispatch receipt hash mismatch", ErrInvalidReceipt)
	}
	return record, nil
}

func DecodeTakeoverReceiptRecord(input []byte) (TakeoverReceiptRecord, error) {
	var record TakeoverReceiptRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return TakeoverReceiptRecord{}, err
	}
	if record.SchemaVersion != evidence.SchemaVersion {
		return TakeoverReceiptRecord{}, fmt.Errorf("%w: unsupported schema version", ErrInvalidReceipt)
	}
	record.Receipt.Evidence = nonNilStrings(record.Receipt.Evidence)
	record.Receipt.ErrorEvidence = nonNilStrings(record.Receipt.ErrorEvidence)
	record.Receipt.OldWriterEvidence = nonNilStrings(record.Receipt.OldWriterEvidence)
	record.Receipt.AuthorizationEvidence = nonNilStrings(record.Receipt.AuthorizationEvidence)
	if err := validateTakeoverReceipt(record.Receipt); err != nil {
		return TakeoverReceiptRecord{}, err
	}
	expected, err := digestMessage(adapterTakeoverReceiptDomain, record.Receipt)
	if err != nil {
		return TakeoverReceiptRecord{}, err
	}
	if record.ReceiptHash != expected {
		return TakeoverReceiptRecord{}, fmt.Errorf("%w: takeover receipt hash mismatch", ErrInvalidReceipt)
	}
	return record, nil
}

func NewFileQueue(root string) (*FileQueue, error) {
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("%w: empty queue root", ErrInvalidEnvelope)
	}
	queue := &FileQueue{root: root}
	for _, dir := range []string{queue.requestsDir(), queue.inflightDir(), queue.resultsDir(), queue.archiveDir()} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	return queue, nil
}

func (queue *FileQueue) Dispatch(request RequestMessage) (RequestEnvelope, error) {
	envelope, err := NewRequestEnvelope(request)
	if err != nil {
		return RequestEnvelope{}, err
	}
	if existing, ok, err := queue.loadExistingRequest(envelope.Message.RequestID); err != nil {
		return RequestEnvelope{}, err
	} else if ok {
		if existing.RecordHash == envelope.RecordHash {
			return existing, nil
		}
		return RequestEnvelope{}, fmt.Errorf("%w: request %q has different payload", ErrRecordConflict, envelope.Message.RequestID)
	}
	path := queue.requestPath(envelope.Message.RequestID)
	if err := writeCanonicalLineNoReplace(path, envelope); err != nil {
		if errors.Is(err, fs.ErrExist) {
			existing, loadErr := queue.loadRequest(path)
			if loadErr != nil {
				return RequestEnvelope{}, loadErr
			}
			if existing.RecordHash == envelope.RecordHash {
				return existing, nil
			}
			return RequestEnvelope{}, fmt.Errorf("%w: request %q has different payload", ErrRecordConflict, envelope.Message.RequestID)
		}
		return RequestEnvelope{}, err
	}
	return envelope, nil
}

func (queue *FileQueue) Claim(requestID, claimID, consumerID string, generation int, claimedAt, leaseExpiresAt time.Time, takeoverReceiptHash *string) (ClaimEnvelope, error) {
	if !evidence.ValidID(requestID) || !evidence.ValidID(claimID) || !evidence.ValidID(consumerID) || generation < 1 {
		return ClaimEnvelope{}, fmt.Errorf("%w: invalid claim identity", ErrInvalidEnvelope)
	}
	if !leaseExpiresAt.After(claimedAt) {
		return ClaimEnvelope{}, fmt.Errorf("%w: invalid lease timing", ErrInvalidEnvelope)
	}
	request, err := queue.ensureInflightRequest(requestID, generation)
	if err != nil {
		return ClaimEnvelope{}, err
	}
	latest, ok, err := queue.currentClaim(requestID)
	if err != nil {
		return ClaimEnvelope{}, err
	}
	if generation == 1 {
		if takeoverReceiptHash != nil {
			return ClaimEnvelope{}, fmt.Errorf("%w: generation 1 claim must not reference takeover receipt", ErrInvalidEnvelope)
		}
		if ok {
			if latest.Message.Generation != 1 || latest.Message.ClaimID != claimID {
				return ClaimEnvelope{}, fmt.Errorf("%w: generation 1 already claimed", ErrRecordConflict)
			}
		}
	} else {
		if !ok {
			return ClaimEnvelope{}, fmt.Errorf("%w: request %q has no prior claim", ErrClaimNotFound, requestID)
		}
		if generation <= latest.Message.Generation {
			return ClaimEnvelope{}, fmt.Errorf("%w: stale generation %d", ErrRecordConflict, generation)
		}
		if claimID == latest.Message.ClaimID {
			if takeoverReceiptHash != nil {
				return ClaimEnvelope{}, fmt.Errorf("%w: same-claim renewal must not reference takeover receipt", ErrInvalidEnvelope)
			}
		} else if takeoverReceiptHash == nil {
			return ClaimEnvelope{}, ErrTakeoverRequired
		}
	}
	message := ClaimMessage{
		Type:                "claim",
		RequestID:           requestID,
		ClaimID:             claimID,
		ConsumerID:          consumerID,
		ClaimedAt:           claimedAt.UTC().Format(TimestampLayout),
		LeaseExpiresAt:      leaseExpiresAt.UTC().Format(TimestampLayout),
		Generation:          generation,
		RequestHash:         request.RecordHash,
		TakeoverReceiptHash: takeoverReceiptHash,
	}
	claim, err := NewClaimEnvelope(message)
	if err != nil {
		return ClaimEnvelope{}, err
	}
	path := queue.claimPath(requestID, generation, claimID)
	if err := writeCanonicalLineNoReplace(path, claim); err != nil {
		if errors.Is(err, fs.ErrExist) {
			existing, loadErr := queue.loadClaim(path)
			if loadErr != nil {
				return ClaimEnvelope{}, loadErr
			}
			if existing.RecordHash == claim.RecordHash {
				return existing, nil
			}
			return ClaimEnvelope{}, fmt.Errorf("%w: claim %q/%d conflict", ErrRecordConflict, requestID, generation)
		}
		return ClaimEnvelope{}, err
	}
	return claim, nil
}

func (queue *FileQueue) SubmitResult(message ResultMessage) (ResultEnvelope, error) {
	message.Type = "result"
	message.OutputEvidence = nonNilStrings(message.OutputEvidence)
	message.ErrorEvidence = nonNilStrings(message.ErrorEvidence)
	if err := validateResultMessage(message); err != nil {
		return ResultEnvelope{}, err
	}
	request, err := queue.loadInflightRequest(message.RequestID)
	if err != nil {
		return ResultEnvelope{}, err
	}
	if request.RecordHash != message.RequestHash {
		return ResultEnvelope{}, fmt.Errorf("%w: request hash mismatch", ErrInvalidEnvelope)
	}
	latest, ok, err := queue.currentClaim(message.RequestID)
	if err != nil {
		return ResultEnvelope{}, err
	}
	if !ok {
		return ResultEnvelope{}, fmt.Errorf("%w: request %q", ErrClaimNotFound, message.RequestID)
	}
	if latest.Message.ClaimID != message.ClaimID || latest.Message.Generation != message.Generation {
		return ResultEnvelope{}, fmt.Errorf("%w: expected (%s,%d), got (%s,%d)", ErrFencingConflict, latest.Message.ClaimID, latest.Message.Generation, message.ClaimID, message.Generation)
	}
	result, err := NewResultEnvelope(message)
	if err != nil {
		return ResultEnvelope{}, err
	}
	path := queue.resultPath(message.RequestID)
	if err := writeCanonicalLineNoReplace(path, result); err != nil {
		if errors.Is(err, fs.ErrExist) {
			existing, loadErr := queue.loadResult(path)
			if loadErr != nil {
				return ResultEnvelope{}, loadErr
			}
			if existing.RecordHash == result.RecordHash {
				return existing, nil
			}
			return ResultEnvelope{}, fmt.Errorf("%w: result for request %q differs", ErrRecordConflict, message.RequestID)
		}
		return ResultEnvelope{}, err
	}
	return result, nil
}

func (queue *FileQueue) LoadResult(requestID string) (ResultEnvelope, error) {
	if !evidence.ValidID(requestID) {
		return ResultEnvelope{}, fmt.Errorf("%w: invalid request id", ErrInvalidEnvelope)
	}
	return queue.loadResult(queue.resultPath(requestID))
}

func (queue *FileQueue) Root() string {
	return queue.root
}

func (queue *FileQueue) requestsDir() string { return filepath.Join(queue.root, "requests") }
func (queue *FileQueue) inflightDir() string { return filepath.Join(queue.root, "inflight") }
func (queue *FileQueue) resultsDir() string  { return filepath.Join(queue.root, "results") }
func (queue *FileQueue) archiveDir() string  { return filepath.Join(queue.root, "archive") }

func (queue *FileQueue) requestPath(requestID string) string {
	return filepath.Join(queue.requestsDir(), requestID+".request.jsonl")
}

func (queue *FileQueue) inflightRequestPath(requestID string) string {
	return filepath.Join(queue.inflightDir(), requestID+".request.jsonl")
}

func (queue *FileQueue) claimPath(requestID string, generation int, claimID string) string {
	return filepath.Join(queue.inflightDir(), fmt.Sprintf("%s.%d.%s.claim.jsonl", requestID, generation, claimID))
}

func (queue *FileQueue) resultPath(requestID string) string {
	return filepath.Join(queue.resultsDir(), requestID+".result.jsonl")
}

func (queue *FileQueue) loadExistingRequest(requestID string) (RequestEnvelope, bool, error) {
	paths := []string{queue.requestPath(requestID), queue.inflightRequestPath(requestID)}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			record, loadErr := queue.loadRequest(path)
			return record, true, loadErr
		} else if !errors.Is(err, fs.ErrNotExist) {
			return RequestEnvelope{}, false, err
		}
	}
	return RequestEnvelope{}, false, nil
}

func (queue *FileQueue) ensureInflightRequest(requestID string, generation int) (RequestEnvelope, error) {
	inflightPath := queue.inflightRequestPath(requestID)
	if generation == 1 {
		if _, err := os.Stat(inflightPath); errors.Is(err, fs.ErrNotExist) {
			requestPath := queue.requestPath(requestID)
			if moveErr := moveNoReplace(requestPath, inflightPath); moveErr != nil {
				if _, statErr := os.Stat(inflightPath); statErr == nil {
					return queue.loadRequest(inflightPath)
				}
				if errors.Is(moveErr, fs.ErrNotExist) {
					return RequestEnvelope{}, fmt.Errorf("%w: request %q", ErrRequestNotFound, requestID)
				}
				return RequestEnvelope{}, moveErr
			}
		} else if err != nil {
			return RequestEnvelope{}, err
		}
	}
	return queue.loadInflightRequest(requestID)
}

// moveNoReplace publishes the source under the target name without ever
// overwriting an existing target: the hard link step fails when the target
// exists, and the source removal is rolled back if it fails.
func moveNoReplace(source, target string) error {
	if err := os.Link(source, target); err != nil {
		return err
	}
	if err := os.Remove(source); err != nil {
		_ = os.Remove(target)
		return err
	}
	return nil
}

func (queue *FileQueue) loadInflightRequest(requestID string) (RequestEnvelope, error) {
	path := queue.inflightRequestPath(requestID)
	record, err := queue.loadRequest(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return RequestEnvelope{}, fmt.Errorf("%w: request %q", ErrRequestNotFound, requestID)
		}
		return RequestEnvelope{}, err
	}
	return record, nil
}

func (queue *FileQueue) currentClaim(requestID string) (ClaimEnvelope, bool, error) {
	pattern := filepath.Join(queue.inflightDir(), requestID+".*.*.claim.jsonl")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return ClaimEnvelope{}, false, err
	}
	if len(paths) == 0 {
		return ClaimEnvelope{}, false, nil
	}
	sort.Strings(paths)
	latest := ClaimEnvelope{}
	found := false
	for _, path := range paths {
		record, readErr := queue.loadClaim(path)
		if readErr != nil {
			return ClaimEnvelope{}, false, readErr
		}
		if record.Message.RequestID != requestID {
			continue
		}
		if !found || record.Message.Generation > latest.Message.Generation {
			latest = record
			found = true
		}
	}
	return latest, found, nil
}

func (queue *FileQueue) loadRequest(path string) (RequestEnvelope, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return RequestEnvelope{}, err
	}
	return DecodeRequestEnvelope(content)
}

func (queue *FileQueue) loadClaim(path string) (ClaimEnvelope, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ClaimEnvelope{}, err
	}
	return DecodeClaimEnvelope(content)
}

func (queue *FileQueue) loadResult(path string) (ResultEnvelope, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ResultEnvelope{}, err
	}
	return DecodeResultEnvelope(content)
}

func validateRequestMessage(message RequestMessage) error {
	if message.Type != "request" {
		return fmt.Errorf("%w: message.type must be request", ErrInvalidEnvelope)
	}
	if !evidence.ValidID(message.RequestID) || !evidence.ValidID(message.RunID) || !evidence.ValidID(message.TaskID) || !evidence.ValidID(message.Adapter) || message.Attempt < 1 {
		return fmt.Errorf("%w: invalid request identity", ErrInvalidEnvelope)
	}
	if !evidence.ValidHash(message.ContextHash) {
		return fmt.Errorf("%w: invalid context hash", ErrInvalidEnvelope)
	}
	if _, err := time.Parse(TimestampLayout, message.CreatedAt); err != nil {
		return fmt.Errorf("%w: invalid request timestamp", ErrInvalidEnvelope)
	}
	return nil
}

func validateClaimMessage(message ClaimMessage) error {
	if message.Type != "claim" {
		return fmt.Errorf("%w: message.type must be claim", ErrInvalidEnvelope)
	}
	if !evidence.ValidID(message.RequestID) || !evidence.ValidID(message.ClaimID) || !evidence.ValidID(message.ConsumerID) || message.Generation < 1 {
		return fmt.Errorf("%w: invalid claim identity", ErrInvalidEnvelope)
	}
	if !evidence.ValidHash(message.RequestHash) {
		return fmt.Errorf("%w: invalid request hash", ErrInvalidEnvelope)
	}
	claimedAt, err := time.Parse(TimestampLayout, message.ClaimedAt)
	if err != nil {
		return fmt.Errorf("%w: invalid claimedAt", ErrInvalidEnvelope)
	}
	leaseExpiresAt, err := time.Parse(TimestampLayout, message.LeaseExpiresAt)
	if err != nil {
		return fmt.Errorf("%w: invalid leaseExpiresAt", ErrInvalidEnvelope)
	}
	if !leaseExpiresAt.After(claimedAt) {
		return fmt.Errorf("%w: lease expiration must be after claim time", ErrInvalidEnvelope)
	}
	if message.TakeoverReceiptHash != nil && !evidence.ValidHash(*message.TakeoverReceiptHash) {
		return fmt.Errorf("%w: invalid takeover receipt hash", ErrInvalidEnvelope)
	}
	return nil
}

func validateResultMessage(message ResultMessage) error {
	if message.Type != "result" {
		return fmt.Errorf("%w: message.type must be result", ErrInvalidEnvelope)
	}
	if !evidence.ValidID(message.RequestID) || !evidence.ValidID(message.ClaimID) || message.Generation < 1 {
		return fmt.Errorf("%w: invalid result identity", ErrInvalidEnvelope)
	}
	if !evidence.ValidHash(message.RequestHash) {
		return fmt.Errorf("%w: invalid request hash", ErrInvalidEnvelope)
	}
	if _, err := time.Parse(TimestampLayout, message.CompletedAt); err != nil {
		return fmt.Errorf("%w: invalid completedAt", ErrInvalidEnvelope)
	}
	if !validUniqueHashes(message.OutputEvidence) || !validUniqueHashes(message.ErrorEvidence) {
		return fmt.Errorf("%w: invalid output/error evidence", ErrInvalidEnvelope)
	}
	switch message.Status {
	case "completed":
		if len(message.OutputEvidence) == 0 || len(message.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: completed result requires outputEvidence and empty errorEvidence", ErrInvalidEnvelope)
		}
	case "failed", "uncertain":
		if len(message.ErrorEvidence) == 0 {
			return fmt.Errorf("%w: %s result requires errorEvidence", ErrInvalidEnvelope, message.Status)
		}
	default:
		return fmt.Errorf("%w: unknown result status", ErrInvalidEnvelope)
	}
	return nil
}

func validateDispatchReceipt(receipt DispatchReceipt) error {
	if receipt.Type != "dispatch" {
		return fmt.Errorf("%w: receipt.type must be dispatch", ErrInvalidReceipt)
	}
	if !evidence.ValidID(receipt.ReceiptID) || !evidence.ValidID(receipt.RequestID) || !evidence.ValidID(receipt.Adapter) {
		return fmt.Errorf("%w: invalid dispatch identity", ErrInvalidReceipt)
	}
	if !validReceiptActor(receipt.RecordedBy) {
		return fmt.Errorf("%w: invalid dispatch actor", ErrInvalidReceipt)
	}
	if !evidence.ValidHash(receipt.RequestHash) || !validUniqueHashes(receipt.Evidence) || !validUniqueHashes(receipt.ErrorEvidence) {
		return fmt.Errorf("%w: invalid dispatch hashes", ErrInvalidReceipt)
	}
	if _, err := time.Parse(TimestampLayout, receipt.OccurredAt); err != nil {
		return fmt.Errorf("%w: invalid dispatch timestamp", ErrInvalidReceipt)
	}
	switch receipt.Transport {
	case "file-queue":
		if receipt.TransportRequestID != nil {
			return fmt.Errorf("%w: file-queue transportRequestId must be null", ErrInvalidReceipt)
		}
	case "sessbridge":
		if receipt.TransportRequestID == nil || !transportRequestIDPattern.MatchString(*receipt.TransportRequestID) || len(*receipt.TransportRequestID) > 128 {
			return fmt.Errorf("%w: sessbridge transportRequestId is required", ErrInvalidReceipt)
		}
	case "ipc":
		if receipt.TransportRequestID != nil && (!transportRequestIDPattern.MatchString(*receipt.TransportRequestID) || len(*receipt.TransportRequestID) > 128) {
			return fmt.Errorf("%w: invalid ipc transportRequestId", ErrInvalidReceipt)
		}
	default:
		return fmt.Errorf("%w: unknown dispatch transport", ErrInvalidReceipt)
	}
	switch receipt.Outcome {
	case "accepted":
		if len(receipt.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: accepted dispatch must have empty errorEvidence", ErrInvalidReceipt)
		}
	case "rejected", "uncertain":
		if len(receipt.ErrorEvidence) == 0 {
			return fmt.Errorf("%w: %s dispatch requires errorEvidence", ErrInvalidReceipt, receipt.Outcome)
		}
	default:
		return fmt.Errorf("%w: unknown dispatch outcome", ErrInvalidReceipt)
	}
	return nil
}

func validateTakeoverReceipt(receipt TakeoverReceipt) error {
	if receipt.Type != "takeover" {
		return fmt.Errorf("%w: receipt.type must be takeover", ErrInvalidReceipt)
	}
	if !evidence.ValidID(receipt.ReceiptID) || !evidence.ValidID(receipt.RequestID) || !evidence.ValidID(receipt.PreviousClaimID) || !evidence.ValidID(receipt.NewClaimID) {
		return fmt.Errorf("%w: invalid takeover identity", ErrInvalidReceipt)
	}
	if !validReceiptActor(receipt.RecordedBy) {
		return fmt.Errorf("%w: invalid takeover actor", ErrInvalidReceipt)
	}
	if !evidence.ValidHash(receipt.RequestHash) || !validUniqueHashes(receipt.Evidence) || !validUniqueHashes(receipt.ErrorEvidence) || !validUniqueHashes(receipt.OldWriterEvidence) || !validUniqueHashes(receipt.AuthorizationEvidence) {
		return fmt.Errorf("%w: invalid takeover hashes", ErrInvalidReceipt)
	}
	if len(receipt.OldWriterEvidence) == 0 || len(receipt.AuthorizationEvidence) == 0 {
		return fmt.Errorf("%w: takeover requires old-writer and authorization evidence", ErrInvalidReceipt)
	}
	if receipt.Reason.Code == "" || !evidence.ValidID(receipt.Reason.Code) || len(receipt.Reason.Message) > 1024 {
		return fmt.Errorf("%w: invalid takeover reason", ErrInvalidReceipt)
	}
	if _, err := time.Parse(TimestampLayout, receipt.OccurredAt); err != nil {
		return fmt.Errorf("%w: invalid takeover timestamp", ErrInvalidReceipt)
	}
	if receipt.PreviousGeneration < 1 || receipt.NewGeneration < 2 || receipt.NewGeneration <= receipt.PreviousGeneration || receipt.NewClaimID == receipt.PreviousClaimID {
		return fmt.Errorf("%w: invalid takeover generations", ErrInvalidReceipt)
	}
	switch receipt.Outcome {
	case "granted":
		if len(receipt.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: granted takeover must have empty errorEvidence", ErrInvalidReceipt)
		}
	case "rejected":
		if len(receipt.ErrorEvidence) == 0 {
			return fmt.Errorf("%w: rejected takeover requires errorEvidence", ErrInvalidReceipt)
		}
	default:
		return fmt.Errorf("%w: unknown takeover outcome", ErrInvalidReceipt)
	}
	return nil
}

func validReceiptActor(actor evidence.Actor) bool {
	if !evidence.ValidID(actor.ID) {
		return false
	}
	return actor.Type == "system" || actor.Type == "operator"
}

func digestMessage(domain string, value any) (string, error) {
	canonical, err := evidence.EncodeCanonical(value)
	if err != nil {
		return "", err
	}
	return evidence.Digest(domain, canonical), nil
}

func validUniqueHashes(values []string) bool {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !evidence.ValidHash(value) {
			return false
		}
		if _, exists := seen[value]; exists {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}

func nonNilStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	cloned := make([]string, len(values))
	copy(cloned, values)
	return cloned
}

func writeCanonicalLineNoReplace(path string, record any) error {
	canonical, err := evidence.EncodeCanonical(record)
	if err != nil {
		return err
	}
	payload := append(canonical, '\n')
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	// A random unique temporary name keeps fs.ErrExist a truthful signal: it can
	// then only mean that the destination already exists, never that two writers
	// happened to pick the same timestamped scratch name.
	file, err := os.CreateTemp(directory, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	temp := file.Name()
	writeErr := error(nil)
	if _, err := file.Write(payload); err != nil {
		writeErr = err
	} else if err := file.Sync(); err != nil {
		writeErr = err
	}
	closeErr := file.Close()
	if writeErr == nil {
		writeErr = closeErr
	}
	if writeErr != nil {
		_ = os.Remove(temp)
		return writeErr
	}
	// no-replace publish: hard links fail when the destination already exists,
	// unlike os.Rename which replaces existing files on all supported platforms.
	if err := os.Link(temp, path); err != nil {
		_ = os.Remove(temp)
		if errors.Is(err, fs.ErrExist) {
			return fs.ErrExist
		}
		return err
	}
	// The record is durably published once the link exists; removing the
	// temporary name is best-effort cleanup and must not report failure.
	_ = os.Remove(temp)
	return nil
}
