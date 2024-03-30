package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"wordeeBot/internal/myErrors"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage() (*UserStorage, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	port := os.Getenv("DB_PORT")

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode))

	if err != nil {
		return nil, myErrors.ErrSql
	}

	if err := db.Ping(); err != nil {
		return nil, myErrors.ErrSql
	}

	return &UserStorage{db: db}, nil
}

func (storage *UserStorage) CheckUser(user_id int64, username string) error {
	var count int

	err := storage.db.QueryRow("SELECT COUNT(*) FROM users WHERE user_id = $1", user_id).Scan(&count)
	if err != nil {
		return myErrors.ErrSql
	}

	if count == 0 {
		err := storage.addUser(user_id, username)
		return err
	}

	return nil
}

func (storage *UserStorage) addUser(user_id int64, username string) error {
	_, err := storage.db.Exec("INSERT INTO users (user_id,username) VALUES($1,$2)", user_id, username)
	if err != nil {
		return myErrors.ErrSql
	}
	return nil
}
