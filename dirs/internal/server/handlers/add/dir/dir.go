package dir

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
)

type Request struct {
	ParentID int64  `json:"parent_id"`
	Name     string `json:"name"`
}

type Response struct {
	Success bool `json:"success"`
}

type FolderAdder interface {
	AddFolder(tx *sql.Tx, data []interface{}) (string, error)
	Begin() (*sql.Tx, error)
}

type DirCreater interface {
	NewDir(path string) error
}

func New(log *slog.Logger, folderAdder FolderAdder, dirCreater DirCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		response := func(code int, message interface{}) {
			render.Status(r, code)
			render.JSON(w, r, message)
		}

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request")
			response(400, "failed to decode request")

			return
		}

		// возможно убрать
		if err := req.emptyVal(); err != nil {
			log.Error("empty field")
			response(400, "empty field")

			return
		}

		folder := []interface{}{req.ParentID, req.Name, time.Now()}

		tx, err := folderAdder.Begin()
		if err != nil {
			log.Error("failed to create transaction")
			response(500, "failed to create transaction")

			return
		}

		path, err := folderAdder.AddFolder(tx, folder)
		if err != nil {
			tx.Rollback()
			log.Error("failed to add folder")
			response(500, "failed to add folder")

			return

		}

		err = dirCreater.NewDir(path)
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), "already exists") {
				log.Error("folder already exists")
				response(409, "folder already exists")
			} else {
				log.Error("failed to create folder")
				response(500, "failed to create folder")
			}

			return

		}

		err = tx.Commit()
		if err != nil {
			log.Error("failed to commit transaction")
			response(500, "failed to commit transaction")

			return
		}

		response(201, Response{Success: true})

	}
}

func (req *Request) emptyVal() error {
	if req.ParentID == 0 || req.Name == "" {
		return fmt.Errorf("field is empty")
	}

	return nil
}
