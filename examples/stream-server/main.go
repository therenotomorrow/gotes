package main

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/config"
	notesv1 "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	usersv1 "github.com/therenotomorrow/gotes/pkg/api/users/v1"
	"github.com/therenotomorrow/gotes/pkg/client"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
)

type User struct {
	Name     string
	Email    string
	Password string
	Token    string
}

func authenticate(ctx context.Context, cli *client.Client) context.Context {
	user := &User{
		Name:     "Kirill Kolesnikov",
		Email:    "memes@gmail.com",
		Password: "kekes/memes",
		Token:    "",
	}

	ruResp, err := cli.RegisterUser(ctx, &usersv1.RegisterUserRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	ex.Panic(err)

	user.Token = ruResp.GetUser().GetToken()

	return cli.Authenticate(ctx, user.Token)
}

func populate(ctx context.Context, cli *client.Client) {
	notes := make([]*notesv1.CreateNoteResponse, 0)

	note1, err := cli.CreateNote(ctx, &notesv1.CreateNoteRequest{
		Title:   "header-1",
		Content: "some content-1",
	})
	ex.Panic(err)

	notes = append(notes, note1)

	note2, err := cli.CreateNote(ctx, &notesv1.CreateNoteRequest{
		Title:   "header-2",
		Content: "some content-2",
	})
	ex.Panic(err)

	notes = append(notes, note2)

	note3, err := cli.CreateNote(ctx, &notesv1.CreateNoteRequest{
		Title:   "header-3",
		Content: "some content-3",
	})
	ex.Panic(err)

	notes = append(notes, note3)

	_, err = cli.DeleteNote(ctx, &notesv1.DeleteNoteRequest{Id: notes[1].GetNote().GetId()})
	ex.Panic(err)
}

func main() {
	var (
		ctx = context.Background()
		cfg = config.MustNew()
		log = trace.Logger(trace.TEXT, cfg.Debug)
		cli = client.MustNew(client.Config{
			Address: cfg.Server.Address,
			Secure:  false,
		})
		wait = sync.WaitGroup{}
	)

	defer cli.Close()

	ctx = authenticate(ctx, cli)

	// for simulate previous unread messages
	populate(ctx, cli)

	wait.Go(func() {
		stream, err := cli.SubscribeToEvents(ctx, &notesv1.SubscribeToEventsRequest{})
		if err != nil {
			log.Error("SubscribeToEvents error", "error", err)
			ex.Panic(err)
		}

		log.Info("SubscribeToEvents streaming")

		for {
			resp, err := stream.Recv()

			switch {
			case errors.Is(err, io.EOF):
				log.Info("SubscribeToEvents closed")

				return
			case err != nil:
				log.Error("SubscribeToEvents error", "error", err)
				ex.Panic(err)
			}

			log.Info("SubscribeToEvents received", "payload", resp.GetPayload())
		}
	})

	// simulate some online events
	wait.Go(func() { populate(ctx, cli) })

	wait.Wait()
}
