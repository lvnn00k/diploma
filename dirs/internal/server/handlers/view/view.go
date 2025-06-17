package view

import (
	f "dirs/internal/storage"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Response struct {
	Path    []f.Path   `json:"path"`
	Files   []f.File   `json:"files"`
	Folders []f.Folder `json:"folders"`
}

type DataSelector interface {
	SelectFiles(parentId int64) ([]f.File, error)
	SelectFolders(parentId int64) ([]f.Folder, error)
	SelectPath(id int64) ([]f.Path, error)
}

func New(log *slog.Logger, dataSelector DataSelector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response := func(code int, message interface{}) {
			render.Status(r, code)
			render.JSON(w, r, message)
		}

		parentId, err := strconv.ParseInt(chi.URLParam(r, "parent_id"), 10, 64)
		if err != nil {
			log.Error("failed parse id")
			response(400, "failed parse id")

			return
		}

		path, err := dataSelector.SelectPath(parentId)
		if err != nil {
			log.Error("path not found")
			response(500, "path not found")

			return
		}

		files, err := dataSelector.SelectFiles(parentId)
		if err != nil {
			log.Error("files not found")
			response(500, "files not found")

			return
		}

		folders, err := dataSelector.SelectFolders(parentId)
		if err != nil {
			log.Error("folders not found")
			response(500, "folders not found")

			return
		}

		response(200, Response{
			Path:    path,
			Files:   files,
			Folders: folders,
		})

	}
}
