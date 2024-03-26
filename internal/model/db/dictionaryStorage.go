package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"wordeeBot/internal/myErrors"
)

type DictionaryStorage struct {
	db *sqlx.DB
}

func NewDictionaryStorage() (*DictionaryStorage, error) {
	db, err := sqlx.Open("postgres", "user=USER password=PASSWORD dbname=DBNAME sslmode=disable")
	if err != nil {
		return nil, myErrors.ErrSql
	}

	if err := db.Ping(); err != nil {
		return nil, myErrors.ErrSql
	}

	return &DictionaryStorage{db: db}, nil
}

func (storage *DictionaryStorage) GetNamesOfUserDictionaries(id int) ([]string, error) {
	names := make([]string, 0)

	rows, err := storage.db.Query("SELECT id,user_id,name FROM dictionaries WHERE user_id = $1", id)
	if err != nil {
		return nil, myErrors.ErrSql
	}
	defer rows.Close()

	for rows.Next() {
		var dict Dictionary

		err = rows.Scan(&dict.ID, &dict.UserId, &dict.Name)
		if err != nil {
			return nil, myErrors.ErrSql
		}

		names = append(names, dict.Name)
	}

	return names, nil
}

func (storage *DictionaryStorage) GetNamesOfUserDicitonariesWithDefinition(id int) ([]string, error) {
	names := make([]string, 0)

	rows, err := storage.db.Query("SELECT id,user_id,name FROM dictionaries WHERE user_id = $1 AND is_definition = $2", id, true)
	if err != nil {
		return nil, myErrors.ErrSql
	}
	defer rows.Close()

	for rows.Next() {
		var dict Dictionary

		err = rows.Scan(&dict.ID, &dict.UserId, &dict.Name)
		if err != nil {
			return nil, myErrors.ErrSql
		}

		names = append(names, dict.Name)
	}

	return names, nil
}

func (storage *DictionaryStorage) CheckDicitonary(name string, id int) (bool, error) {
	var count int

	err := storage.db.Get(&count, "SELECT COUNT(*) FROM dictionaries WHERE name = $1 AND user_id = $2", strings.ToLower(name), id)
	if err != nil {
		return false, myErrors.ErrSql
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (storage *DictionaryStorage) AddDictionary(name string, id int, columns []string) error {
	_, err := storage.db.Exec("INSERT INTO dictionaries (user_id,name) VALUES($1,$2)", id, name)
	if err != nil {
		return myErrors.ErrSql
	}
	for _, column := range columns {
		query := fmt.Sprintf("UPDATE dictionaries SET %s = $1 WHERE user_id = $2 AND name = $3", "is_"+column)

		_, err := storage.db.Exec(query, true, id, name)
		if err != nil {
			return myErrors.ErrSql
		}
	}

	return nil
}

func (storage *DictionaryStorage) GetNamesOfDictionaryColumns(id int, name string) ([]string, error) {
	var dict Dictionary

	err := storage.db.QueryRow(
		"SELECT id,user_id,name,is_transcription,is_translation,is_synonyms,"+
			"is_antonyms,is_definition,is_collocations,is_idioms FROM dictionaries"+
			" WHERE user_id = $1 AND name = $2", id, name).Scan(&dict.ID, &dict.UserId, &dict.Name, &dict.IsTranscription,
		&dict.IsTranslation, &dict.IsSynonyms, &dict.IsAntonyms, &dict.IsDefinition, &dict.IsCollocations, &dict.IsIdioms)

	if err != nil {
		return nil, myErrors.ErrSql
	}

	columns := make([]string, 0)
	columns = append(columns, "Слово")
	checkField(dict.IsTranscription, "Транскрипция", &columns)
	checkField(dict.IsTranslation, "Перевод", &columns)
	checkField(dict.IsSynonyms, "Синонимы", &columns)
	checkField(dict.IsAntonyms, "Антонимы", &columns)
	checkField(dict.IsDefinition, "Определение", &columns)
	checkField(dict.IsCollocations, "Коллокации", &columns)
	checkField(dict.IsIdioms, "Идиомы", &columns)
	return columns, nil
}

func checkField(field bool, name string, columns *[]string) {
	if field {
		*columns = append(*columns, name)
	}
}

func (storage *DictionaryStorage) GetDictionaryId(user_id int, name string) (int, error) {
	var id int
	err := storage.db.QueryRow("SELECT id FROM dictionaries WHERE user_id = $1 AND name = $2", user_id, name).Scan(&id)
	if err != nil {
		return -1, myErrors.ErrSql
	}
	return id, nil
}
