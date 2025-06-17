package mysql

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

type User struct {
	Id    int64
	Login string
	Hash  string
	Role  int8
}

func New(StorageConnect string) (*Storage, error) {

	db, err := sql.Open("mysql", StorageConnect)
	if err != nil {
		return nil, err
	}

	return &Storage{db}, nil

}

func (s *Storage) SelectUser(login string) (User, error) {

	var user User

	stmt, err := s.db.Prepare("SELECT * FROM users WHERE login = ?")
	if err != nil {
		return user, err
	}

	err = stmt.QueryRow(login).Scan(&user.Id, &user.Login, &user.Hash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("login does not exist")
		}

		return user, err
	}

	return user, nil

}

func (s *Storage) NewUser(login string, role int8, hash string) error {

	stmt, err := s.db.Prepare("INSERT INTO users(login, password, role_id) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(login, hash, role)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("login exists")
		}

		return err
	}

	return nil

}
