-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    user_id BIGINT PRIMARY KEY,
    username TEXT
);

CREATE TABLE IF NOT EXISTS dictionaries(
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users (user_id),
    name VARCHAR(255) NOT NULL UNIQUE,
    is_transcription boolean NOT NULL,
    is_translation boolean NOT NULL,
    is_synonyms boolean NOT NULL,
    is_antonyms boolean NOT NULL,
    is_definition boolean NOT NULL,
    is_collocations boolean NOT NULL,
    is_idioms boolean
);

CREATE TABLE IF NOT EXISTS words(
    id SERIAL PRIMARY KEY,
    dictionary_id BIGINT REFERENCES dictionaries (id),
    writing TEXT NOT NULL UNIQUE,
    transcription TEXT,
    translation TEXT,
    synonyms TEXT,
    antonyms TEXT,
    definition TEXT,
    collocations TEXT,
    idioms TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS words;
DROP TABLE IF EXISTS dictionaries;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd