package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/tickets"
)

// ErrAgentRunnerLaunchReconfirmation reports a failed post-R launch
// reconfirmation: the authorization or budget facts no longer permit launch.
// The underlying authorization or budget failure is preserved for errors.Is.
var ErrAgentRunnerLaunchReconfirmation = errors.New("AgentRunner launch reconfirmation failed")

// AgentRunnerLaunchReconfirmer performs the narrow post-R reconfirmation
// required between publishing the request record and starting the process.
// An admission verdict is a point-in-time preflight, not a lock: dispatch must
// reconfirm a non-revoked grant and an outstanding budget reservation against
// the current ledgers just before launch, and fail closed otherwise.
//
// The reconfirmer has value semantics, performs no writes, and is safe to call
// repeatedly. Clock is the only evaluation-time authority.
type AgentRunnerLaunchReconfirmer struct {
	RequestRecord       AgentRunnerRequestRecord
	AuthorizationLedger chain.AuthorizationLedger
	CostLedger          *tickets.CostLedger
	Clock               func() time.Time
}

// ReconfirmLaunchAuthorization revalidates the bound authorization grant and
// budget reservation for the request record. Any failure blocks launch.
func (reconfirmer AgentRunnerLaunchReconfirmer) ReconfirmLaunchAuthorization(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %w", ErrAgentRunnerLaunchReconfirmation, err)
	}
	if err := ValidateAgentRunnerRequestRecord(reconfirmer.RequestRecord); err != nil {
		return fmt.Errorf("%w: %w", ErrAgentRunnerLaunchReconfirmation, err)
	}
	if reconfirmer.Clock == nil {
		return fmt.Errorf("%w: evaluation clock missing", ErrAgentRunnerLaunchReconfirmation)
	}
	now := reconfirmer.Clock().UTC()
	if now.IsZero() {
		return fmt.Errorf("%w: evaluation clock returned zero time", ErrAgentRunnerLaunchReconfirmation)
	}
	body := reconfirmer.RequestRecord.Request
	if err := validateAgentRunnerAuthorization(reconfirmer.AuthorizationLedger, body, now); err != nil {
		return fmt.Errorf("%w: %w", ErrAgentRunnerLaunchReconfirmation, err)
	}
	if err := validateAgentRunnerBudget(reconfirmer.CostLedger, body); err != nil {
		return fmt.Errorf("%w: %w", ErrAgentRunnerLaunchReconfirmation, err)
	}
	return nil
}
