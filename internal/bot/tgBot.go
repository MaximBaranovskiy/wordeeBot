package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"wordeeBot/internal/model/db"
	"wordeeBot/internal/myErrors"
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
	UserID int64
	Name   string
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

func (b *TgBotModel) ListenForUpdates() {
	u := tgbotapi.NewUpdate(0)
	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		ProcessCommands(b, update)
	}
}

func ProcessCommands(b *TgBotModel, update tgbotapi.Update) {
	id, err := b.userStorage.CheckUser(update.SentFrom().ID)
	if err != nil {
		myErrors.HandleError(b.bot, update.Message.Chat.ID, err)
		return
	}

	if update.CallbackQuery != nil {
		handleCallbacks(b, update, id)
	} else if update.Message != nil {
		handleMessages(b, update, id)
	}
}
