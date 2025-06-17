package logout

import (
	"auth/internal/cookie"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func Logout(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// _, err := r.Cookie("refreshToken")
		// if err != nil {
		// 	render.Status(r, 400)
		// 	render.JSON(w, r, "cookie not found")

		// 	return
		// }

		cookie.DeleteToken(&w)

		render.Status(r, 200)
		render.JSON(w, r, nil)

	}
}
