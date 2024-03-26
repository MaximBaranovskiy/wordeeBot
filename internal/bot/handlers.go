package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"wordeeBot/internal/myErrors"
)

type handlerFunc func(b *TgBotModel, update tgbotapi.Update) error

var callbackData2Handler = map[string]handlerFunc{
	"myDictionaries":   handleShowDictionaries,
	"createDictionary": handleCreateDictionary,
	"editDictionary":   handleEditDictionary,
	"mainMenu":         handleMainMenu,
	"addingWord":       handleAddingWord,
	"cancelWord":       handleCancelledAddingWord,
}

var usrLstCommand2Handler = map[string]handlerFunc{
	"create_dictionary_name": handleCreateDictionaryColumns,
	"edit_dictionary":        handleChooseDictionaryForEdit,
	"editDictionary_name":    handleDictionaryEditing,
	"myDictionaries":         handleSendParticularDictionary,
}

func handleCallbacks(b *TgBotModel, update tgbotapi.Update) {
	if handler, ok := callbackData2Handler[update.CallbackData()]; ok {
		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handler(b, update))
		return
	}

	if handler, ok := usrLstCommand2Handler[b.userLastCommand[update.CallbackQuery.From.ID]]; ok {
		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handler(b, update))
		return
	}
}

func handleMessages(b *TgBotModel, update tgbotapi.Update) {
	switch {

	case strings.Contains(update.Message.Text, "start"):

		myErrors.HandleError(b.bot, update.Message.Chat.ID, handleStart(b, update))

	case b.userLastCommand[update.Message.From.ID] == "createDictionary":

		myErrors.HandleError(b.bot, update.Message.Chat.ID, handleCreateDictionaryName(b, update))

	case strings.Contains(b.userLastCommand[update.Message.From.ID], "addingWord"):

		myErrors.HandleError(b.bot, update.Message.Chat.ID, handlePreparationAddingWord(b, update))

	}
}
