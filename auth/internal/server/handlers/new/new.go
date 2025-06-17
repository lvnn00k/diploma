package new

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     int8   `json:"role"`
}

type Response struct {
	Success bool `json:"success"`
}

type UserCreater interface {
	NewUser(login string, role int8, hash string) error
}

func New(log *slog.Logger, userCreator UserCreater) http.HandlerFunc {
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

		bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			log.Error("failed to create hash")
			fail(500, "failed to create user")

			return
		}

		hash := string(bytes)

		err = userCreator.NewUser(req.Login, req.Role, hash)
		if err != nil {
			if strings.Contains(err.Error(), "login exists") {
				log.Error("login exists")
				fail(409, "login exists")
			} else {
				log.Error("failed to create user")
				fail(500, "failed to create user")
			}

			return
		}

		render.Status(r, 201)
		render.JSON(w, r, Response{
			Success: true,
		})

	}
}
