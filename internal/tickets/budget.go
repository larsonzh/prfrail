package tickets

import (
	"errors"
	"fmt"
	"time"
)

type BudgetState string

const (
	BudgetPendingReview  BudgetState = "pending-review"
	BudgetOverrideWindow BudgetState = "override-window"
	BudgetHardBlock      BudgetState = "hard-block"
)

var (
	ErrBudgetHardBlocked        = errors.New("ticket budget hard blocked")
	ErrBudgetNeedsOverride      = errors.New("ticket budget requires override")
	ErrBudgetOverrideExpired    = errors.New("ticket override expired")
	ErrBudgetOverrideExhausted  = errors.New("ticket override attempts exhausted")
	ErrBudgetAttemptExhausted   = errors.New("ticket attempt already consumed")
	ErrBudgetWallClockExhausted = errors.New("ticket wall-clock budget exhausted")
	ErrBudgetCostExhausted      = errors.New("ticket cost budget exhausted")
)

func (ledger *Ledger) BudgetState(now time.Time) BudgetState {
	failures := ledger.failureCount()
	if failures >= ledger.body.Policy.HardBlockThreshold {
		return BudgetHardBlock
	}
	if failures < ledger.body.Policy.ReviewThreshold {
		return BudgetPendingReview
	}
	index, override, valid := ledger.activeOverride(now)
	if !valid {
		return BudgetPendingReview
	}
	if ledger.overrideUsage[index] >= override.MaxAdditionalAttempts {
		return BudgetPendingReview
	}
	return BudgetOverrideWindow
}

func (ledger *Ledger) ReserveAttempt(attempt int, now time.Time, wallClockExhausted, costExhausted bool) error {
	if wallClockExhausted {
		return ErrBudgetWallClockExhausted
	}
	if costExhausted {
		return ErrBudgetCostExhausted
	}
	if attempt < 1 {
		return fmt.Errorf("%w: attempt must be positive", ErrInvalidLedgerOperation)
	}
	if _, exists := ledger.consumedByAttempt[attempt]; exists {
		return ErrBudgetAttemptExhausted
	}
	state := ledger.BudgetState(now)
	if state == BudgetHardBlock {
		return ErrBudgetHardBlocked
	}
	failures := ledger.failureCount()
	if failures < ledger.body.Policy.ReviewThreshold {
		ledger.consumedByAttempt[attempt] = struct{}{}
		return nil
	}
	index, override, valid := ledger.activeOverride(now)
	if !valid {
		return ErrBudgetNeedsOverride
	}
	expiresAt, err := time.Parse(timestampLayout, override.ExpiresAt)
	if err != nil || !expiresAt.After(now.UTC()) {
		return ErrBudgetOverrideExpired
	}
	if ledger.overrideUsage[index] >= override.MaxAdditionalAttempts {
		return ErrBudgetOverrideExhausted
	}
	ledger.overrideUsage[index]++
	ledger.consumedByAttempt[attempt] = struct{}{}
	return nil
}

// CostBudgetExhausted derives whether cost admission should fail closed under
// a shared monetary cap. Unknown committed usage is treated as exhausted.
func CostBudgetExhausted(capMicros *int64, summary CostSummary) bool {
	if capMicros == nil {
		return false
	}
	committed, known := summary.CommittedAmountMicros()
	if !known {
		return true
	}
	return committed >= *capMicros
}

// CostTokenBudgetExhausted checks token caps against settled plus unknown
// reservation holds. It intentionally fails closed once the cap is met.
func CostTokenBudgetExhausted(capTokens int, summary CostSummary) bool {
	if capTokens < 0 {
		return true
	}
	holdTokens := summary.ReservedTokens - summary.SettledTokens
	if holdTokens < 0 {
		holdTokens = 0
	}
	committed := summary.SettledTokens + holdTokens
	return committed >= capTokens
}

func (ledger *Ledger) failureCount() int {
	count := 0
	for _, entry := range ledger.body.Entries {
		if entry.Kind == "failure" {
			count++
		}
	}
	return count
}

func (ledger *Ledger) activeOverride(now time.Time) (int, Entry, bool) {
	lastFailureTime := ""
	for _, entry := range ledger.body.Entries {
		if entry.Kind == "failure" && entry.RecordedAt > lastFailureTime {
			lastFailureTime = entry.RecordedAt
		}
	}
	for index := len(ledger.body.Entries) - 1; index >= 0; index-- {
		entry := ledger.body.Entries[index]
		if entry.Kind != "override-granted" {
			continue
		}
		if entry.RecordedAt <= lastFailureTime {
			continue
		}
		expiresAt, err := time.Parse(timestampLayout, entry.ExpiresAt)
		if err != nil || !expiresAt.After(now.UTC()) {
			continue
		}
		return index, entry, true
	}
	return -1, Entry{}, false
}
