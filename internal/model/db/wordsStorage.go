package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"wordeeBot/internal/myErrors"
)

type WordsStorage struct {
	db *sqlx.DB
}

func NewWordsStorage() (*WordsStorage, error) {
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

	if err = db.Ping(); err != nil {
		return nil, myErrors.ErrSql
	}

	return &WordsStorage{db: db}, nil
}

func (storage *WordsStorage) AddWord(word *Word) error {
	_, err := storage.db.Exec("INSERT INTO words (dictionary_id,writing,transcription,translation,synonyms,"+
		"antonyms,definition,collocations,idioms) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)", word.DictionaryId,
		word.Writing, word.Transcription, word.Translation, word.Synonyms, word.Antonyms, word.Definition,
		word.Collocations, word.Idioms)

	if err != nil {
		return myErrors.ErrSql
	}
	return nil
}

func (storage *WordsStorage) CheckWord(dictionaryId int, writing string) (bool, error) {
	var count int

	err := storage.db.Get(&count, "SELECT COUNT(*) FROM WORDS WHERE dictionary_id = $1 AND writing = $2", dictionaryId, writing)
	if err != nil {
		return false, myErrors.ErrSql
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (storage *WordsStorage) GetAllDictionaryWords(dictionaryId int) ([]Word, error) {
	words := make([]Word, 0)

	rows, err := storage.db.Query("SELECT * FROM words WHERE dictionary_id=$1", dictionaryId)
	defer rows.Close()
	if err != nil {
		return nil, myErrors.ErrSql
	}

	for rows.Next() {
		var word Word

		err := rows.Scan(&word.ID, &word.DictionaryId, &word.Writing, &word.Transcription, &word.Translation,
			&word.Synonyms, &word.Antonyms, &word.Definition, &word.Collocations, &word.Idioms)
		if err != nil {
			return nil, myErrors.ErrSql
		}

		words = append(words, word)
	}

	return words, nil
}
