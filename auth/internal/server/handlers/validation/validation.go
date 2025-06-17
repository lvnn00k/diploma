package validation

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type TokenValidator interface {
	ValidateToken(refToken string) error
}

func Validation(log *slog.Logger, tokenValidator TokenValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fail := func(code int, message string) {
			render.Status(r, code)
			render.JSON(w, r, message)
		}

		cookies, err := r.Cookie("accessToken") // access
		if err != nil {
			fail(401, "cookie not found")
			return
		}

		token := cookies.Value

		err = tokenValidator.ValidateToken(token)
		if err != nil {
			fail(401, "invalid token")
			return
		}

		render.Status(r, 200)
		render.Respond(w, r, nil)

	}
}
