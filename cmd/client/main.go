package main

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/config"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	"github.com/therenotomorrow/gotes/pkg/client"
	"github.com/therenotomorrow/gotes/pkg/slogx"
)

func main() {
	var (
		ctx = context.Background()
		cfg = config.MustNew()
		log = slogx.New(slogx.TEXT, cfg.Debug)
		cli = client.MustNew(client.Config{
			Address: cfg.Server.Address,
			Secure:  false,
		})
	)

	defer cli.Close()

	note, err := cli.CreateNote(ctx, &pb.CreateNoteRequest{Title: "header", Content: "some note"})
	log.Info("CreateNote response", "note", note, "error", err)

	ident := note.GetNote().GetId()

	get, err := cli.RetrieveNote(ctx, &pb.RetrieveNoteRequest{Id: ident})
	log.Info("RetrieveNote response", "get", get, "error", err)

	notes, err := cli.ListNotes(ctx, &pb.ListNotesRequest{})
	log.Info("ListNotes response", "notes", notes, "error", err)

	del, err := cli.DeleteNote(ctx, &pb.DeleteNoteRequest{Id: ident})
	log.Info("DeleteNote response", "del", del, "error", err)
}
