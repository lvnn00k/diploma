package admin

import (
	"net/http"

	"github.com/go-chi/render"
)

type TokenValidator interface {
	ValidateToken(refToken string) error
	GetClaims(refToken string) (map[string]interface{}, error)
}

func AdminMiddleware(tokenValidator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fail := func(code int, message string) {
				render.Status(r, code)
				render.JSON(w, r, message)
			}

			cookie, err := r.Cookie("accessToken")
			if err != nil {
				fail(401, "cookie not found")
				return
			}

			token := cookie.Value

			err = tokenValidator.ValidateToken(token)
			if err != nil {
				fail(401, "invalid token")
				return
			}

			cliams, err := tokenValidator.GetClaims(token)
			if err != nil {
				fail(401, "Incorrect token")
				return
			}

			if int8(cliams["userRole"].(float64)) != 1 {
				fail(401, "Insufficient access rights")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
