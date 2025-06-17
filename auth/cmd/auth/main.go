package main

import (
	"auth/internal/config"
	"auth/internal/server/handlers/login"
	"auth/internal/server/handlers/logout"
	"auth/internal/server/handlers/new"
	"auth/internal/server/handlers/refresh"
	"auth/internal/server/handlers/validation"
	"auth/internal/server/middleware/admin"
	"auth/internal/storage/mysql"
	"auth/internal/token"
	"log/slog"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger()
	log.Info("starting auth-server")

	storage, err := mysql.New(cfg.StorageConnect)
	if err != nil {
		log.Error("failed connect to storage")
		os.Exit(1)
	}

	token := token.New(cfg.SecretKey)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://client:80"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	}))

	router.Post("/login", login.New(log, storage, token))
	router.Get("/logout", logout.Logout(log))
	router.Get("/refresh-token", refresh.New(log, token))
	router.Get("/validate-token", validation.Validation(log, token))

	router.Group(func(r chi.Router) {
		r.Use(admin.AdminMiddleware(token))
		r.Post("/create", new.New(log, storage))
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

}

func setupLogger() *slog.Logger {

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
