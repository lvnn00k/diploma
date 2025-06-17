package search

import (
	f "dirs/internal/storage"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type File struct {
	File f.File `json:"file"`
	Link string `json:"path"`
}

type Folder struct {
	Folder f.Folder `json:"folder"`
	Link   string   `json:"path"`
}

type Response struct {
	Files   []File   `json:"files"`
	Folders []Folder `json:"folders"`
}

type DataSearch interface {
	SearchFiles(name string) ([]f.File, error)
	SearchFolders(name string) ([]f.Folder, error)
	SelectPath(id int64) ([]f.Path, error)
}

func New(log *slog.Logger, dataSearch DataSearch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response := func(code int, message interface{}) {
			render.Status(r, code)
			render.JSON(w, r, message)
		}

		name := chi.URLParam(r, "name")

		files, err := dataSearch.SearchFiles(name)
		if err != nil {
			log.Error("files not found")
			response(500, "files not found")

			return

		}

		var respFile []File

		for i := range files { // убрать вложенный цикл
			path, err := dataSearch.SelectPath(int64(files[i].ParentId))
			if err != nil {
				log.Error("file path select error")
				response(500, "file path select error")

				return
			}

			var folderName []string

			for j := range path {
				folderName = append(folderName, path[j].Name)
			}

			link := strings.Join(folderName, "/")

			respFile = append(respFile, File{files[i], link})
		}

		folders, err := dataSearch.SearchFolders(name)
		if err != nil {
			log.Error("folders not found")
			response(500, "folders not found")

			return

		}

		var respFolder []Folder

		for i := range folders { // убрать вложенный цикл
			path, err := dataSearch.SelectPath(int64(folders[i].Id))
			if err != nil {
				log.Error("folder path select error")
				response(500, "folder path select error")

				return
			}

			var folderName []string

			for j := range path {
				folderName = append(folderName, path[j].Name)
			}

			link := strings.Join(folderName, "/")

			respFolder = append(respFolder, Folder{folders[i], link})
		}

		response(200, Response{
			Files:   respFile,
			Folders: respFolder,
		})

	}
}
