package main

import (
	"context"
	"log/slog"
	"strings"

	"github.com/therenotomorrow/gotes/internal/config"
	notesv1 "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	usersv1 "github.com/therenotomorrow/gotes/pkg/api/users/v1"
	"github.com/therenotomorrow/gotes/pkg/client"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
	"google.golang.org/grpc/status"
)

type User struct {
	Name     string
	Email    string
	Password string
	Token    string
}

func HappyPath(cli *client.Client, log *slog.Logger, user *User) *User {
	ctx := context.Background()

	ruResp, err := cli.RegisterUser(ctx, &usersv1.RegisterUserRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	log.Info("RegisterUser response", "resp", ruResp, "error", err)

	rtResp, err := cli.RefreshToken(ctx, &usersv1.RefreshTokenRequest{
		Email:    user.Email,
		Password: user.Password,
	})
	log.Info("RefreshToken response", "resp", rtResp, "error", err)

	user.Token = rtResp.GetUser().GetToken()

	ctx = cli.Authenticate(ctx, user.Token)

	cnResp, err := cli.CreateNote(ctx, &notesv1.CreateNoteRequest{
		Title:   "header",
		Content: "some content",
	})
	log.Info("CreateNote response", "resp", cnResp, "error", err)

	ident := cnResp.GetNote().GetId()

	rnResp, err := cli.RetrieveNote(ctx, &notesv1.RetrieveNoteRequest{
		Id: ident,
	})
	log.Info("RetrieveNote response", "resp", rnResp, "error", err)

	lnResp, err := cli.ListNotes(ctx, &notesv1.ListNotesRequest{})
	log.Info("ListNotes response", "resp", lnResp, "error", err)

	dnResp, err := cli.DeleteNote(ctx, &notesv1.DeleteNoteRequest{
		Id: ident,
	})
	log.Info("DeleteNote response", "resp", dnResp, "error", err)

	return user
}

func FailPath(cli *client.Client, log *slog.Logger, user *User) {
	var (
		err error
		ctx = context.Background()
	)

	_, err = cli.RegisterUser(ctx, &usersv1.RegisterUserRequest{
		Name:     user.Name,
		Email:    "",
		Password: "",
	})
	detailed(log, "RegisterUser", err)

	_, err = cli.RefreshToken(ctx, &usersv1.RefreshTokenRequest{
		Email:    user.Email,
		Password: "invalid password",
	})
	detailed(log, "RefreshToken", err)

	ctx = cli.Authenticate(context.Background(), "68a478d5-8e87-4d50-a136-a7fa1c7efbe2")

	_, err = cli.CreateNote(ctx, &notesv1.CreateNoteRequest{
		Title:   "header",
		Content: "some content",
	})
	detailed(log, "CreateNote 1", err)

	ctx = cli.Authenticate(context.Background(), user.Token)

	_, err = cli.CreateNote(ctx, &notesv1.CreateNoteRequest{
		Title:   "kek",
		Content: "mem",
	})
	detailed(log, "CreateNote 2", err)

	maxLen := 256
	_, err = cli.CreateNote(ctx, &notesv1.CreateNoteRequest{
		Title:   strings.Repeat("5", maxLen),
		Content: "my cool content",
	})
	detailed(log, "CreateNote 3", err)

	resp, err := cli.CreateNote(ctx, &notesv1.CreateNoteRequest{
		Title:   "header",
		Content: "some content",
	})
	detailed(log, "CreateNote 4", err)

	_, err = cli.RetrieveNote(ctx, &notesv1.RetrieveNoteRequest{
		Id: &typespb.ID{Value: -1},
	})
	detailed(log, "RetrieveNote 1", err)

	ident := resp.GetNote().GetId()

	_, err = cli.RetrieveNote(ctx, &notesv1.RetrieveNoteRequest{
		Id: ident,
	})
	detailed(log, "RetrieveNote 2", err)
}

func detailed(log *slog.Logger, msg string, err error) {
	st := status.Convert(err)

	if len(st.Details()) == 0 {
		log.Info(msg, "extra", "no details", "error", err)

		return
	}

	for _, detail := range st.Details() {
		proto, _ := detail.(*typespb.Error)

		log.Info(msg, "code", proto.GetCode(), "reason", proto.GetReason())
	}
}

func main() {
	var (
		cfg = config.MustNew()
		log = trace.Logger(trace.TEXT, cfg.Debug)
		cli = client.MustNew(client.Config{
			Address: cfg.Server.Address,
			Secure:  false,
		})
	)

	defer cli.Close()

	user := &User{
		Name:     "Kirill Kolesnikov",
		Email:    "kkxnes@gmail.com",
		Password: "kekes/memes",
		Token:    "",
	}

	user = HappyPath(cli, log, user)

	FailPath(cli, log, user)
}
