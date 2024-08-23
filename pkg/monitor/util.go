package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration.Seconds() < 60:
		return "just now"
	case duration.Minutes() < 60:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration.Hours() < 24:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration.Hours() < 48:
		return "yesterday"
	default:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	}
}

func CreateEvent(ctx context.Context, monitorId, fromStatus, toStatus string, env *types.Env) error {
	eventId, err := gonanoid.Generate(NANOID_ALPHABET_LIST, NANOID_LENGTH)
	if err != nil {
		log.Warnf("Error generating nanoid for event creation: %s\n", err.Error())
		return err
	}
	err = env.DB.Query.CreateEvent(ctx, db.CreateEventParams{
		ID:         eventId,
		MonitorID:  monitorId,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
	})
	return err
}
