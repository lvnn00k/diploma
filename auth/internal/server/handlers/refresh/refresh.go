package refresh

import (
	"auth/internal/cookie"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type TokenWorker interface {
	CreateToken(user_id int64, fole int8) ([]string, error)
	GetClaims(refToken string) (map[string]interface{}, error)
}

func New(log *slog.Logger, tokenWorker TokenWorker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fail := func(code int, message string) {
			render.Status(r, code)
			render.JSON(w, r, message)
		}

		cookies, err := r.Cookie("refreshToken")
		if err != nil {
			log.Error("cookie not found")
			fail(401, "cookie not found")
			return
		}

		refToken := cookies.Value

		claims, err := tokenWorker.GetClaims(refToken)
		if err != nil {
			log.Error("Incorrect token")
			fail(401, "Incorrect token")

			return
		}

		token, err := tokenWorker.CreateToken(int64(claims["userid"].(float64)), int8(claims["userRole"].(float64)))
		if err != nil {

			log.Error("failed to create token")
			fail(401, "failed to create token")

			return
		}

		cookie.SetToken(token, &w)

		render.Status(r, 200)
		render.JSON(w, r, nil)

	}
}
