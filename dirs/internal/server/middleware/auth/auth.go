package auth

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

var key = []byte("GYgds6fe64372") //

func JWTAuthMiddleware(next http.Handler) http.Handler {
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

		reqToken := cookie.Value

		token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if err != nil {
			fail(401, "invalid token")
			return
		}

		if !token.Valid {
			fail(401, "invalid token")
			return
		}

		next.ServeHTTP(w, r)

	})
}
