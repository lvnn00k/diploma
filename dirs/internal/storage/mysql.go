package storage

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

type Path struct {
	Id    int32  `json:"id"`
	Name  string `json:"name"`
	Depth int16  `json:"depth"`
}

type File struct {
	Id         int32  `json:"id"`
	ParentId   int32  `json:"parent_id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	DateCreate string `json:"date_create"`
}

type Folder struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	DateCreate string `json:"date_create"`
}

func New(StorageConnect string) (*Storage, error) {

	db, err := sql.Open("mysql", StorageConnect)
	if err != nil {
		return nil, err
	}

	return &Storage{db}, nil

}

func (s *Storage) Begin() (*sql.Tx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *Storage) SelectFiles(parentId int64) ([]File, error) {

	var files []File

	stmt, err := s.db.Prepare("SELECT * FROM files WHERE parent_id = ? ORDER BY name")
	if err != nil {
		return files, err
	}

	rows, err := stmt.Query(parentId)
	if err != nil {
		return files, err
	}

	for rows.Next() {
		var file File
		var date []byte

		if err := rows.Scan(&file.Id, &file.ParentId, &file.Name, &file.Size, &date); err != nil {
			return files, err
		}

		file.DateCreate = string(date)

		files = append(files, file)
	}

	return files, nil

}

func (s *Storage) SelectFolders(parentId int64) ([]Folder, error) {

	var folders []Folder

	stmt, err := s.db.Prepare(`SELECT folders.* FROM folders 
		JOIN closureFolder AS t ON folders.id = t.child_id 
		WHERE t.parent_id = ? AND t.level = 1 
		ORDER BY folders.name`)
	if err != nil {
		return folders, err
	}

	rows, err := stmt.Query(parentId)
	if err != nil {
		return folders, err
	}

	for rows.Next() {
		var folder Folder
		var date []byte

		if err := rows.Scan(&folder.Id, &folder.Name, &date); err != nil {
			return folders, err
		}

		folder.DateCreate = string(date)

		folders = append(folders, folder)
	}

	return folders, nil

}

func (s *Storage) SearchFiles(name string) ([]File, error) {
	var files []File

	stmt, err := s.db.Prepare("SELECT * FROM files WHERE name LIKE ? ORDER BY name")
	if err != nil {
		return files, err
	}

	rows, err := stmt.Query(fmt.Sprintf("%%%s%%", name))
	if err != nil {
		return files, err
	}

	for rows.Next() {
		var file File
		var date []byte

		if err := rows.Scan(&file.Id, &file.ParentId, &file.Name, &file.Size, &date); err != nil {
			return files, err
		}

		file.DateCreate = string(date)

		files = append(files, file)
	}

	return files, nil
}

func (s *Storage) SearchFolders(name string) ([]Folder, error) {

	var folders []Folder

	stmt, err := s.db.Prepare("SELECT * FROM folders WHERE name LIKE ? ORDER BY name")
	if err != nil {
		return folders, err
	}

	rows, err := stmt.Query(fmt.Sprintf("%%%s%%", name))
	if err != nil {
		return folders, err
	}

	for rows.Next() {
		var folder Folder
		var date []byte

		if err := rows.Scan(&folder.Id, &folder.Name, &date); err != nil {
			return folders, err
		}

		folder.DateCreate = string(date)

		folders = append(folders, folder)
	}

	return folders, nil
}

func (s *Storage) AddFile(tx *sql.Tx, file []interface{}) (string, error) {

	var link string

	stmt, err := tx.Prepare(`INSERT INTO files(parent_id, name, date_create, size) 
		VALUES (?, ?, CONVERT_TZ(?, '+00:00', '+07:00'), ?)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(file...)
	if err != nil {
		return link, err
	}

	link, err = s.GetLink(tx, file[0].(int64))
	if err != nil {
		return link, err
	}

	return link, nil
}

func (s *Storage) AddFolder(tx *sql.Tx, folder []interface{}) (string, error) {

	var link string
	var id int64

	stmt, err := tx.Prepare(`INSERT INTO folders(name, date_create) 
		VALUES (?, CONVERT_TZ(?, '+00:00', '+07:00'))`)
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(folder[1], folder[2])
	if err != nil {
		return link, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return link, err
	}

	stmt, err = tx.Prepare(`INSERT INTO closureFolder(parent_id, child_id, level) VALUES(?, ?, 0);`)
	if err != nil {
		return link, err
	}

	_, err = stmt.Exec(id, id)
	if err != nil {
		return link, err
	}

	stmt, err = tx.Prepare(`INSERT INTO closureFolder(parent_id, child_id, level) 
		SELECT p.parent_id, c.child_id, p.level + c.level+1 
		FROM closureFolder p, closureFolder c 
		WHERE p.child_id = ? and c.parent_id = ?`)
	if err != nil {
		return link, err
	}

	_, err = stmt.Exec(folder[0], id)
	if err != nil {
		return link, err
	}

	link, err = s.GetLink(tx, id)
	if err != nil {
		return link, err
	}

	return link, nil
}

func (s *Storage) SelectPath(id int64) ([]Path, error) {

	var paths []Path

	stmt, err := s.db.Prepare(`SELECT folders.id, folders.name, t.level
		FROM folders
		JOIN closureFolder AS t ON folders.id = t.parent_id
		WHERE t.child_id = ? ORDER BY level desc`)
	if err != nil {
		return paths, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return paths, err
	}

	for rows.Next() {
		var path Path

		if err := rows.Scan(&path.Id, &path.Name, &path.Depth); err != nil {
			return paths, err
		}

		paths = append(paths, path)
	}

	return paths, nil

}

func (s *Storage) GetLink(tx *sql.Tx, id int64) (string, error) {

	var link string
	var names []string

	stmt, err := tx.Prepare(`SELECT folders.name
		FROM folders
		JOIN closureFolder AS t ON folders.id = t.parent_id
		WHERE t.child_id = ? ORDER BY level desc`)
	if err != nil {
		return link, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return link, err
	}

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return link, err
		}

		names = append(names, name)
	}

	link = strings.Join(names, "/")

	return link, nil

}

func (s *Storage) DeleteFolder(tx *sql.Tx, id int64) (string, error) {

	var link string
	var folderIds []string

	link, err := s.GetLink(tx, id)
	if err != nil {
		return link, err
	}

	stmt, err := tx.Prepare(`SELECT child_id FROM closureFolder WHERE parent_id = ?`)
	if err != nil {
		return link, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return link, err
	}

	for rows.Next() {
		var folderId string

		err = rows.Scan(&folderId)
		if err != nil {
			return link, err
		}

		folderIds = append(folderIds, folderId)
	}

	param := strings.Join(folderIds, ", ")

	stmt, err = tx.Prepare(fmt.Sprintf(`DELETE FROM files WHERE parent_id in (%s)`, param))
	if err != nil {
		return link, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return link, err
	}

	stmt, err = tx.Prepare(fmt.Sprintf(`DELETE FROM closureFolder WHERE child_id in (%s)`, param))
	if err != nil {
		return link, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return link, err
	}

	stmt, err = tx.Prepare(fmt.Sprintf(`DELETE FROM folders WHERE id in (%s)`, param))
	if err != nil {
		return link, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return link, err
	}

	return link, nil
}

func (s *Storage) DeleteFile(tx *sql.Tx, id int64) (string, error) {

	var parent_id int64
	var name, link string

	stmt, err := tx.Prepare("SELECT parent_id, name FROM files WHERE id = ?")
	if err != nil {
		return link, err
	}

	err = stmt.QueryRow(id).Scan(&parent_id, &name)
	if err != nil {
		return link, err
	}

	stmt, err = tx.Prepare("DELETE FROM files WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return link, err
	}

	link, err = s.GetLink(tx, parent_id)
	if err != nil {
		return link, err
	}
	link = fmt.Sprintf("%s/%s", link, name)

	return link, nil
}
