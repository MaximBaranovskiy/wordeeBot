package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"wordeeBot/internal/model/db"
)

type TgBotModel struct {
	bot                        *tgbotapi.BotAPI
	userLastCommand            map[int64]string
	tempColumnsForDictionaries map[DictionaryIdentificator][]string
	userlastMessageID          map[int64]int
	userStorage                *db.UserStorage
	dictionaryStorage          *db.DictionaryStorage
	wordsStorage               *db.WordsStorage
	tempStorageForAddingWords  map[int64]*db.Word
}

type DictionaryIdentificator struct {
	User_id int64
	Name    string
}

func NewClient(token string) (*TgBotModel, error) {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "Проблемы при инициализации бота")
	}
	userLastCommand := make(map[int64]string)
	tempColumnsForDictionaries := make(map[DictionaryIdentificator][]string)
	userLastMessageId := make(map[int64]int)
	tempStorageForAddingWords := make(map[int64]*db.Word)

	userStorage, err := db.NewUserStorage()
	if err != nil {
		return nil, errors.Wrap(err, "Проблемы при инициализации хранилища пользователей")
	}

	dictionaryStorage, err := db.NewDictionaryStorage()
	if err != nil {
		return nil, errors.Wrap(err, "Проблемы при инициализации хранилища словарей")
	}

	wordsStorage, err := db.NewWordsStorage()
	if err != nil {
		return nil, errors.Wrap(err, "Проблемы при инициализации хранилища слов")
	}

	return &TgBotModel{
		bot:                        client,
		userLastCommand:            userLastCommand,
		tempColumnsForDictionaries: tempColumnsForDictionaries,
		userlastMessageID:          userLastMessageId,
		userStorage:                userStorage,
		dictionaryStorage:          dictionaryStorage,
		wordsStorage:               wordsStorage,
		tempStorageForAddingWords:  tempStorageForAddingWords,
	}, nil
}

func (c *TgBotModel) ListenForUpdates() {
	u := tgbotapi.NewUpdate(0)
	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		ProcessCommands(c, update)
	}
}

func ProcessCommands(c *TgBotModel, update tgbotapi.Update) {

	if update.CallbackQuery != nil {
		handleCallbacks(c, update)
	} else if update.Message != nil {
		handleMessages(c, update)
	}
}
