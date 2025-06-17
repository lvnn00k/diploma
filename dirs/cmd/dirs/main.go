package main

import (
	"dirs/internal/config"
	"dirs/internal/files"
	"dirs/internal/server/handlers/add/dir"
	"dirs/internal/server/handlers/add/file"
	"dirs/internal/server/handlers/delete"
	"dirs/internal/server/handlers/search"
	"dirs/internal/server/handlers/view"
	"dirs/internal/server/middleware/auth"
	"dirs/internal/storage"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger()
	log.Info("starting dir-server")

	storage, err := storage.New(cfg.StorageConnect)
	if err != nil {
		log.Error("failed to connect with DB")
		os.Exit(1)
	}

	data := files.New(cfg.PathToFiles)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(auth.JWTAuthMiddleware)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://client:80"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
	}))

	router.Get("/d/{parent_id}", view.New(log, storage))
	router.Get("/search/{name}", search.New(log, storage))
	router.Post("/add/file", file.New(log, storage, data))
	router.Post("/add/folder", dir.New(log, storage, data))
	router.Delete("/delete/folder/{id}", delete.DeleteFolder(log, storage, data))
	router.Delete("/delete/file/{id}", delete.DeleteFile(log, storage, data))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(1)
	}

}

func setupLogger() *slog.Logger {

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
