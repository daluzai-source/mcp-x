package media

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Role1776/mcp-x/app/internal/domain/x"
)

const mediaID = "1146654567674912769"

func TestNewMedia(t *testing.T) {
	cases := []struct {
		name      string
		props     MediaProps
		wantErrIs error
	}{
		{name: "a finalized image without processing info", props: MediaProps{ID: mediaID, ProcessingState: ""}},
		{name: "a pending video", props: MediaProps{ID: mediaID, ProcessingState: StatePending}},
		{name: "a video in progress", props: MediaProps{ID: mediaID, ProcessingState: StateInProgress}},
		{name: "a succeeded video", props: MediaProps{ID: mediaID, ProcessingState: StateSucceeded}},
		{name: "a failed video", props: MediaProps{ID: mediaID, ProcessingState: StateFailed}},
		{name: "an empty id is rejected", props: MediaProps{ID: ""}, wantErrIs: x.ErrInvalidMediaID},
		{name: "a malformed id is rejected", props: MediaProps{ID: "nope"}, wantErrIs: x.ErrInvalidMediaID},
		{
			name:      "an unknown processing state is rejected",
			props:     MediaProps{ID: mediaID, ProcessingState: "transcoding"},
			wantErrIs: x.ErrInvalidRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			media, err := NewMedia(c.props)

			if c.wantErrIs != nil {
				require.ErrorIs(t, err, c.wantErrIs)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.props.ID, media.ID().String())
		})
	}
}

func TestMediaStatePredicates(t *testing.T) {
	cases := []struct {
		name          string
		state         string
		wantIsPending bool
		wantIsFailed  bool
		wantIsReady   bool
	}{
		{name: "an image carries no processing info and is ready", state: "", wantIsReady: true},
		{name: "pending", state: StatePending, wantIsPending: true},
		{name: "in progress", state: StateInProgress, wantIsPending: true},
		{name: "succeeded", state: StateSucceeded, wantIsReady: true},
		{name: "failed", state: StateFailed, wantIsFailed: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			media, err := NewMedia(MediaProps{ID: mediaID, ProcessingState: c.state})
			require.NoError(t, err)

			assert.Equal(t, c.wantIsPending, media.IsPending())
			assert.Equal(t, c.wantIsFailed, media.IsFailed())
			assert.Equal(t, c.wantIsReady, media.IsReady())
		})
	}
}

func TestMediaCheckAfter(t *testing.T) {
	cases := []struct {
		name string
		secs int
		want time.Duration
	}{
		{name: "no hint from the api", secs: 0, want: 0},
		{name: "a five second hint", secs: 5, want: 5 * time.Second},
		{name: "a one minute hint", secs: 60, want: time.Minute},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			media, err := NewMedia(MediaProps{
				ID: mediaID, ProcessingState: StatePending, CheckAfterSecs: c.secs,
			})
			require.NoError(t, err)

			assert.Equal(t, c.want, media.CheckAfter())
		})
	}
}

func TestMediaCheckAfterIsPassedThroughUnclamped(t *testing.T) {
	media, err := NewMedia(MediaProps{ID: mediaID, ProcessingState: StatePending, CheckAfterSecs: -5})
	require.NoError(t, err)

	assert.Equal(t, -5*time.Second, media.CheckAfter())
}

func TestMediaView(t *testing.T) {
	t.Run("a view mirrors the props", func(t *testing.T) {
		expires := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

		media, err := NewMedia(MediaProps{
			ID:              mediaID,
			ExpiresAt:       expires,
			ProcessingState: StateInProgress,
			ProgressPercent: 42,
		})
		require.NoError(t, err)

		view := media.View()

		assert.Equal(t, mediaID, view.ID)
		assert.Equal(t, expires, view.ExpiresAt)
		assert.Equal(t, StateInProgress, view.ProcessingState)
		assert.Equal(t, 42, view.ProgressPercent)
	})

	t.Run("an absent expiry stays the zero time", func(t *testing.T) {
		media, err := NewMedia(MediaProps{ID: mediaID})
		require.NoError(t, err)

		assert.True(t, media.View().ExpiresAt.IsZero())
	})
}
