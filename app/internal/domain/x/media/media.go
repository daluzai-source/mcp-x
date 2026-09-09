package media

import (
	"time"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

const (
	StatePending    = "pending"
	StateInProgress = "in_progress"
	StateSucceeded  = "succeeded"
	StateFailed     = "failed"
)

type Media struct {
	id              MediaID
	expiresAt       time.Time
	processingState string
	progressPercent int
	checkAfter      time.Duration
}

type MediaProps struct {
	ID              string
	ExpiresAt       time.Time
	ProcessingState string
	ProgressPercent int
	CheckAfterSecs  int
}

func NewMedia(props MediaProps) (Media, error) {
	id, err := NewMediaID(props.ID)
	if err != nil {
		return Media{}, err
	}

	switch props.ProcessingState {
	case "", StatePending, StateInProgress, StateSucceeded, StateFailed:
	default:
		return Media{}, x.ErrInvalidRequest
	}

	return Media{
		id:              id,
		expiresAt:       props.ExpiresAt,
		processingState: props.ProcessingState,
		progressPercent: props.ProgressPercent,
		checkAfter:      time.Duration(props.CheckAfterSecs) * time.Second,
	}, nil
}

type MediaView struct {
	ID              string
	ExpiresAt       time.Time
	ProcessingState string
	ProgressPercent int
}

func (m Media) View() MediaView {
	return MediaView{
		ID:              m.id.String(),
		ExpiresAt:       m.expiresAt,
		ProcessingState: m.processingState,
		ProgressPercent: m.progressPercent,
	}
}

func (m Media) ID() MediaID {
	return m.id
}

func (m Media) CheckAfter() time.Duration {
	return m.checkAfter
}

func (m Media) IsPending() bool {
	return m.processingState == StatePending || m.processingState == StateInProgress
}

func (m Media) IsFailed() bool {
	return m.processingState == StateFailed
}

func (m Media) IsReady() bool {
	return !m.IsPending() && !m.IsFailed()
}
