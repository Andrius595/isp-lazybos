package server

import (
	"context"
	"crypto/rand"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ramasauskas/ispbet/user"
)

func (s *Server) userRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(s.sessions.Auth)

	r.Post("/logout", s.logout)
	r.Post("/verify-email/{token}", s.withUser(s.confirmEmail))

	return r
}

func (s *Server) sendEmailVerification(ctx context.Context, u user.User) error {
	tok, err := randomTextToken(6)
	if err != nil {
		return err
	}

	ver := user.EmailVerification{
		UserUUID: u.UUID,
		Token:    tok,
	}

	if err := s.db.InsertEmailVerification(ctx, ver); err != nil {
		return err
	}

	return s.email.SendEmail(context.Background(), u.Email, tok)
}

func (s *Server) confirmEmail(w http.ResponseWriter, r *http.Request, u user.User) {
	ctx := r.Context()
	log := s.logger("confirmEmail")

	token := chi.URLParamFromCtx(ctx, "token")
	if token == "" {
		respondErr(w, badRequestErr(errors.New("invalid token")))
		return
	}

	ve, ok, err := s.db.FetchEmailVerification(ctx, token)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch email verification")
		respondErr(w, internalErr())
		return
	}

	if !ok {
		log.Info().Str("token", token).Msg("cannot find")
		respondErr(w, notFoundErr())
		return
	}

	if err := user.VerifyUserEmail(&u, &ve); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	if err := s.db.InsertUserVerification(ctx, u, ve); err != nil {
		log.Error().Err(err).Msg("cannot insert user verification")
		respondErr(w, internalErr())
		return
	}

	respondOK(w)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	s.sessions.RevokeAll(r.Context(), w)
}

type EmailSender interface {
	SendEmail(ctx context.Context, to, msg string) error
}

func randomTextToken(n int) (string, error) {
	s := "QWERTYUIOPASDFGHJKLZXCVBNM"

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for pos := range b {
		b[pos] = s[(int(b[pos]) % len(s))]
	}

	return string(b), nil
}
