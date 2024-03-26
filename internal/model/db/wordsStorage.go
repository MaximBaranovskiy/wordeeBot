package db

import (
	"github.com/jmoiron/sqlx"
	"wordeeBot/internal/myErrors"
)

type WordsStorage struct {
	db *sqlx.DB
}

func NewWordsStorage() (*WordsStorage, error) {
	db, err := sqlx.Open("postgres", "user=USER password=PASSWORD dbname=DBNAME sslmode=disable")
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
