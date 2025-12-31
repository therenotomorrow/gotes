package v1

import (
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MarshalNote(note *entities.Note) *pb.Note {
	return &pb.Note{
		Id:        &typespb.ID{Value: note.ID.Value()},
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: timestamppb.New(note.UpdatedAt),
	}
}

func MarshalNotes(notes []*entities.Note) []*pb.Note {
	pbNotes := make([]*pb.Note, len(notes))
	for i, note := range notes {
		pbNotes[i] = MarshalNote(note)
	}

	return pbNotes
}
