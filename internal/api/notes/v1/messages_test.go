package v1_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "github.com/therenotomorrow/gotes/internal/api/notes/v1"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMarshalNote(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Second)
	note := &entities.Note{
		CreatedAt: now.Add(-time.Hour),
		UpdatedAt: now.Add(time.Hour),
		Owner:     new(entities.User),
		Title:     "title",
		Content:   "content",
		ID:        id.New(42),
	}

	got := v1.MarshalNote(note)
	want := &pb.Note{
		Id:      &typespb.ID{Value: 42},
		Title:   "title",
		Content: "content",
		CreatedAt: &timestamppb.Timestamp{
			Seconds: now.Add(-time.Hour).Unix(),
			Nanos:   0,
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: now.Add(time.Hour).Unix(),
			Nanos:   0,
		},
	}

	assert.Equal(t, want, got)
}

func TestMarshalNotes(t *testing.T) {
	t.Parallel()

	empty := timestamppb.New(time.Time{})
	notes := []*entities.Note{
		{Title: "title 1", ID: id.New(1)},
		{Title: "title 3", ID: id.New(3)},
		{Title: "title 2", ID: id.New(2)},
	}

	got := v1.MarshalNotes(notes)
	want := []*pb.Note{
		{Title: "title 1", Id: &typespb.ID{Value: 1}, CreatedAt: empty, UpdatedAt: empty},
		{Title: "title 3", Id: &typespb.ID{Value: 3}, CreatedAt: empty, UpdatedAt: empty},
		{Title: "title 2", Id: &typespb.ID{Value: 2}, CreatedAt: empty, UpdatedAt: empty},
	}

	assert.Equal(t, want, got)
}
