package event

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	gonanoid "github.com/matoous/go-nanoid"
)

const (
	NANOID_ALPHABET_LIST = "abcdefghijklmnopqstuvwxyzABCDEFGHIJKLMNOPQSTUVWXYZ"
	NANOID_LENGTH        = 22
)

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
