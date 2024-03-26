package db

type User struct {
	User_id  int64
	Username string
}

type Dictionary struct {
	ID              int
	UserId          int
	Name            string
	IsTranscription bool
	IsTranslation   bool
	IsSynonyms      bool
	IsAntonyms      bool
	IsDefinition    bool
	IsCollocations  bool
	IsIdioms        bool
}

type Word struct {
	ID            int
	DictionaryId  int
	Writing       string
	Transcription string
	Translation   string
	Synonyms      string
	Antonyms      string
	Definition    string
	Collocations  string
	Idioms        string
}

var (
	DbWordFields = []string{
		"Перевод",
		"Транскрипция",
		"Синонимы",
		"Антонимы",
		"Определение",
		"Коллокации",
		"Идиомы",
	}
)

func (w *Word) ToString() string {
	var (
		str = ""

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

		fields = []string{
			w.Writing,
			w.Transcription,
			w.Translation,
			w.Synonyms,
			w.Antonyms,
			w.Definition,
			w.Collocations,
			w.Idioms,
		}
	)

	for i := 0; i < len(fields); i++ {
		if fields[i] != "" {
			str += columns[i] + ":" + fields[i] + ".\n"
		}
	}
	return str
}
