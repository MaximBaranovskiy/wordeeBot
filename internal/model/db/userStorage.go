package db

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"wordeeBot/internal/myErrors"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage() (*UserStorage, error) {
	db, err := sqlx.Open("postgres", "user=USER password=PASSWORD dbname=DBNAME sslmode=disable")
	if err != nil {
		return nil, myErrors.ErrSql
	}

	if err := db.Ping(); err != nil {
		return nil, myErrors.ErrSql
	}

	return &UserStorage{db: db}, nil
}

func (storage *UserStorage) CheckUser(user_id int64) (int, error) {
	var user User

	err := storage.db.QueryRow("SELECT * FROM users WHERE user_id = $1", user_id).Scan(&user.ID, &user.User_id)

	if err != nil {
		if err == sql.ErrNoRows {
			id, err := storage.addUser(user_id)
			if err != nil {
				return -1, myErrors.ErrSql
			}

			return id, nil
		}

		return -1, myErrors.ErrSql
	}

	return user.ID, nil
}

func (storage *UserStorage) addUser(user_id int64) (int, error) {
	var id int
	err := storage.db.QueryRow("INSERT INTO users (user_id) VALUES($1) RETURNING id", user_id).Scan(&id)
	if err != nil {
		return -1, myErrors.ErrSql
	}
	return id, nil
}
