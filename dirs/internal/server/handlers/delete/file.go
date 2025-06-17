package delete

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type FileDeleter interface {
	DeleteFile(tx *sql.Tx, id int64) (string, error)
	Begin() (*sql.Tx, error)
}

type DataRemover interface {
	DataRemove(path string) error
}

func DeleteFile(log *slog.Logger, fileDeleter FileDeleter, dataRemover DataRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response := func(code int, message interface{}) {
			render.Status(r, code)
			render.JSON(w, r, message)
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Error("failed parse id")
			response(400, "failed parse id")

			return
		}

		tx, err := fileDeleter.Begin()
		if err != nil {
			log.Error("failed to create transaction")
			response(500, "failed to create transaction")

			return
		}

		link, err := fileDeleter.DeleteFile(tx, id)
		if err != nil {
			tx.Rollback()
			log.Error("failed to delete file")
			response(500, "failed to delete file")

			return
		}

		err = dataRemover.DataRemove(link)
		if err != nil {
			tx.Rollback()
			log.Error("failed to remove file")
			response(500, "failed to remove file")

			return
		}

		err = tx.Commit()
		if err != nil {
			log.Error("failed to commit transaction")
			response(500, "failed to commit transaction")

			return
		}

		render.NoContent(w, r)

	}

}
