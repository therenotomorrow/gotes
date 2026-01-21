package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/usecases"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	domuuid "github.com/therenotomorrow/gotes/internal/domain/types/uuid"
)

const (
	ErrMarshalEvent   ex.Error = "marshal event error"
	ErrUnmarshalEvent ex.Error = "unmarshal event error"
)

type EventsRepository struct {
	rdb redis.UniversalClient
}

type event struct {
	EventTime time.Time          `json:"eventTime"`
	EventType entities.EventType `json:"eventType"`
	NoteID    id.ID              `json:"noteId"`
	ID        uuid.UUID          `json:"id"`
}

func MarshalEvent(e *entities.Event) ([]byte, error) {
	val, err := json.Marshal(event{
		ID:        uuid.MustParse(e.ID.Value()),
		EventType: e.EventType,
		NoteID:    e.Note.ID,
		EventTime: e.EventTime,
	})
	if err != nil {
		return nil, ErrMarshalEvent.Because(err)
	}

	return val, nil
}

func UnmarshalEvent(raw []byte) (*entities.Event, error) {
	var event event

	err := json.Unmarshal(raw, &event)
	if err != nil {
		return nil, ErrUnmarshalEvent.Because(err)
	}

	note := new(entities.Note)
	note.ID = event.NoteID

	return &entities.Event{
		ID:        domuuid.Conv(event.ID.String()),
		EventType: event.EventType,
		Note:      note,
		EventTime: event.EventTime,
	}, nil
}

func NewEventsRepository(rdb redis.UniversalClient) *EventsRepository {
	return &EventsRepository{rdb: rdb}
}

func (e *EventsRepository) SaveEvent(ctx context.Context, event *entities.Event) error {
	data, err := MarshalEvent(event)
	if err != nil {
		return err
	}

	key := eventsKey(event.Note.Owner)
	err = e.rdb.RPush(ctx, key, data).Err()

	return ex.Unexpected(err)
}

func (e *EventsRepository) GetEvent(ctx context.Context, user *entities.User) (*entities.Event, error) {
	key := eventsKey(user)
	raw, err := e.rdb.LPop(ctx, key).Bytes()

	switch {
	case errors.Is(err, redis.Nil):
		return nil, usecases.ErrZeroEvents
	case err != nil:
		return nil, ex.Unexpected(err)
	}

	return UnmarshalEvent(raw)
}

func (e *EventsRepository) CountEvents(ctx context.Context, user *entities.User) (int32, error) {
	key := eventsKey(user)
	cnt, err := e.rdb.LLen(ctx, key).Result()

	return int32(cnt), ex.Unexpected(err) //nolint:gosec // allowed conversation
}

func eventsKey(user *entities.User) string {
	return fmt.Sprintf("user:%d:events", user.ID.Value())
}
