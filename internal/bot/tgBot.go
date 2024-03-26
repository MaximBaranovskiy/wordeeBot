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
	tempStorageForEditingWords map[int64]*StructForAddingWord
}

type StructForAddingWord struct {
	Word    *db.Word
	Columns []string
	Count   int
}

type DictionaryIdentificator struct {
	UserID int64
	Name   string
}

func NewBot(token string) (*TgBotModel, error) {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "Проблемы при инициализации бота")
	}
	userLastCommand := make(map[int64]string)
	tempColumnsForDictionaries := make(map[DictionaryIdentificator][]string)
	userLastMessageId := make(map[int64]int)
	tempStorageForAddingWords := make(map[int64]*db.Word)
	tempStorageForEditingWords := make(map[int64]*StructForAddingWord)

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
		tempStorageForEditingWords: tempStorageForEditingWords,
	}, nil
}

func (b *TgBotModel) ListenForUpdates() {
	u := tgbotapi.NewUpdate(0)
	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		go b.processCommands(update)
	}
}

func (b *TgBotModel) processCommands(update tgbotapi.Update) {
	err := b.userStorage.CheckUser(update.SentFrom().ID, update.SentFrom().UserName)
	if err != nil {
		myErrors.HandleError(b.bot, update.Message.Chat.ID, err)
		return
	}

	if update.CallbackQuery != nil {
		handleCallbacks(b, update)
	} else if update.Message != nil {
		handleMessages(b, update)
	}
}
