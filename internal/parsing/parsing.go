package parsing

import (
	"strings"
	"wordeeBot/internal/model/db"
	"wordeeBot/internal/myErrors"
)

func ParseInfoForAdding(text string, dictionary_id int) (*db.Word, error) {
	var (
		word db.Word

		err error

		columns = []string{
			"слово",
			"транскрипция",
			"перевод",
			"синонимы",
			"антонимы",
			"определение",
			"коллокации",
			"идиомы",
		}

		fields = []*string{
			&word.Writing,
			&word.Transcription,
			&word.Translation,
			&word.Synonyms,
			&word.Antonyms,
			&word.Definition,
			&word.Collocations,
			&word.Idioms,
		}
	)

	text = strings.ToLower(text)

	word.DictionaryId = dictionary_id

	for i := 0; i < len(fields); i++ {
		ind := strings.Index(text, columns[i])
		if i == 0 && ind == -1 {
			return nil, myErrors.ErrParseWordInfo
		}

		if ind == -1 {
			continue
		}

		*fields[i], err = getSubstring(":", ".", &text)
		if err != nil {
			return nil, err
		}

		if i == 0 && *fields[i] == "" {
			return nil, myErrors.ErrParseWordInfo
		}
	}

	return &word, err
}

func getSubstring(start, end string, text *string) (string, error) {
	ind1 := strings.Index(*text, start)
	ind2 := strings.Index(*text, end)
	if ind1 == -1 || ind2 == -1 {
		return "", myErrors.ErrParseWordInfo
	}

	str := strings.TrimSpace(strings.ToLower((*text)[ind1+1 : ind2]))
	if strings.Contains(str, "\n") {
		return "", myErrors.ErrParseWordInfo
	}

	*text = (*text)[ind2+1:]
	return str, nil
}
