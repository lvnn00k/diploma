package cookie

import (
	"net/http"
)

func SetToken(tokens []string, w *http.ResponseWriter) {
	cookies := []http.Cookie{
		{
			Name:     "accessToken",
			Value:    tokens[0],
			Path:     "/",
			MaxAge:   900, //900 15 minutes
			HttpOnly: true,
			Secure:   false,
		},
		{
			Name:     "refreshToken",
			Value:    tokens[1],
			Path:     "/api/auth/", // api/auth
			MaxAge:   86400,        // 1 day
			HttpOnly: true,
			Secure:   false,
		},
	}

	for _, cookie := range cookies {
		http.SetCookie(*w, &cookie)
	}

}

func DeleteToken(w *http.ResponseWriter) {

	cookies := []http.Cookie{
		{
			Name:     "accessToken",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false,
		},
		{
			Name:     "refreshToken",
			Value:    "",
			Path:     "/api/auth/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false,
		},
	}

	for _, cookie := range cookies {
		http.SetCookie(*w, &cookie)
	}

}
