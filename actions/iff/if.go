package iff

import (
	"context"

	"go.uber.org/zap"
)

// If represents the set of conditions to running an action
type If struct {
	EventMatch     *EventMatch     `json:"event-match"`
	AlertMatch     *AlertMatch     `json:"alert-match"`
	RoomIssueMatch *RoomIssueMatch `json:"room-issue-match"`
	StateQuery     *StateQuery     `json:"state-query"`
}

// Check returns whether or not the if check passes
func (i *If) Check(ctx context.Context, log *zap.SugaredLogger) (context.Context, bool) {
	if i.EventMatch != nil && !i.EventMatch.DoesEventMatch(ctx) {
		return ctx, false
	}

	if i.AlertMatch != nil && !i.AlertMatch.DoesAlertMatch(ctx) {
		return ctx, false
	}

	if i.RoomIssueMatch != nil && !i.RoomIssueMatch.DoesRoomIssueMatch(ctx) {
		return ctx, false
	}

	if i.StateQuery != nil {
		var cont bool
		cont, ctx = i.StateQuery.CheckStore(ctx)
		if cont {
			return ctx, true
		}
		return ctx, false
	}

	return ctx, true
}
