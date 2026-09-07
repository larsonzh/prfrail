package chain

import (
	"context"
	"errors"
)

func (engine *Engine) Recover(ctx context.Context) error {
	events, err := engine.options.Events.Load(ctx)
	if err != nil {
		return err
	}
	projection, err := Rebuild(events)
	if err != nil {
		return err
	}
	engine.state.projection = projection
	if err := engine.options.Reconciler.Reconcile(ctx, engine.options.RunID); err != nil {
		return engine.pauseUncertain(ctx, err)
	}
	for _, state := range projection.StepStates {
		if state == "RUNNING" || state == "WAITING_FOR_OPERATOR" {
			return engine.pauseUncertain(ctx, ErrRecoveryUncertain)
		}
	}
	return nil
}

func (engine *Engine) pauseUncertain(ctx context.Context, cause error) error {
	if engine.state.projection.ChainState == "RUNNING" {
		if err := engine.transition(ctx, chainEntity(engine.options.RunID), "PAUSED", nil, "recovery-uncertain"); err != nil {
			return errors.Join(cause, err)
		}
	}
	return errors.Join(ErrRecoveryUncertain, cause)
}
