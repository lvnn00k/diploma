package file

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/render"
)

type Request struct {
	ParentId int64
	Name     string
	File     multipart.File
}

type Response struct {
	Success bool `json:"success"`
}

type FileAdder interface {
	AddFile(tx *sql.Tx, file []interface{}) (string, error)
	Begin() (*sql.Tx, error)
}

type FileCreater interface {
	NewFile(name, place string, file io.Reader) error
}

func New(log *slog.Logger, fileAdder FileAdder, fileCreater FileCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response := func(code int, message interface{}) {
			render.Status(r, code)
			render.JSON(w, r, message)

		}

		var req Request

		header, err := req.parseForm(r)
		if err != nil {
			log.Error("failed parse form")
			response(400, "failed parse form")
			return
		}

		// возможно убрать
		if err := req.emptyVal(); err != nil {
			log.Error("empty field")
			response(400, "empty field")

			return
		}

		file, fileReader, err := req.fileProcessing(header)
		if err != nil {
			log.Error("failed to read file")
			response(500, "failed to read file")

			return
		}

		defer req.File.Close()

		tx, err := fileAdder.Begin()
		if err != nil {
			log.Error("failed to create transaction")
			response(500, "failed to create transaction")

			return
		}

		path, err := fileAdder.AddFile(tx, file)
		if err != nil {
			tx.Rollback()
			log.Error("failed to add file")
			response(500, "failed to add file")

			return

		}

		err = fileCreater.NewFile(req.Name, path, fileReader)
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), "file exist") {
				log.Error("file already exists")
				response(409, "file already exists")
			} else {
				log.Error("failed to create file")
				response(500, "failed to create file")
			}

			return
		}

		err = tx.Commit()
		if err != nil {
			log.Error("failed to commit transaction")
			response(500, "failed to commit transaction")

			return
		}

		response(201, Response{true})

	}
}

func (req *Request) fileProcessing(handler *multipart.FileHeader) ([]interface{}, io.Reader, error) {

	fileBytes, err := io.ReadAll(req.File)
	if err != nil {
		return nil, nil, err
	}

	file_type := http.DetectContentType(fileBytes)
	if !strings.Contains(file_type, "image/") {
		return nil, nil, fmt.Errorf("file is not an image")
	}

	size := len(fileBytes)

	req.Name = fmt.Sprintf("%s%s", req.Name, filepath.Ext(handler.Filename))

	file := []interface{}{req.ParentId, req.Name, time.Now(), size}
	fileReader := io.NopCloser(bytes.NewReader(fileBytes))

	return file, fileReader, nil

}

func (req *Request) parseForm(r *http.Request) (*multipart.FileHeader, error) {
	var header *multipart.FileHeader
	var err error

	r.ParseMultipartForm(10 << 20)

	req.File, header, err = r.FormFile("file")
	if err != nil {
		return nil, err
	}

	req.Name = r.FormValue("name")

	req.ParentId, err = strconv.ParseInt(r.FormValue("parent_id"), 10, 64)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (req *Request) emptyVal() error {
	if req.File == nil || req.Name == "" || req.ParentId == 0 {
		return fmt.Errorf("field is empty")
	}
	return nil
}
