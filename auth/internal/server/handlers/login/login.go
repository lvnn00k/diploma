package login

import (
	"auth/internal/cookie"
	"auth/internal/storage/mysql"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserSelecter interface {
	SelectUser(login string) (mysql.User, error)
}

type TokenCreator interface {
	CreateToken(user_id int64, role int8) ([]string, error)
}

func New(log *slog.Logger, userSelector UserSelecter, tokenCreator TokenCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		fail := func(code int, message string) {
			render.Status(r, code)
			render.JSON(w, r, message)
		}

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request")
			fail(400, "failed to decode request")

			return
		}

		user, err := userSelector.SelectUser(req.Login)
		if err != nil {
			log.Error("Incorrect login")
			fail(401, "Incorrect login or password")

			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(req.Password))
		if err != nil {
			log.Error("Incorrect password")
			fail(401, "Incorrect login or password")

			return
		}

		token, err := tokenCreator.CreateToken(user.Id, user.Role)
		if err != nil {

			log.Error("failed to create token")
			fail(500, "failed to create token")

			return
		}

		cookie.SetToken(token, &w)

		render.Status(r, 200)
		render.JSON(w, r, nil)

	}
}
