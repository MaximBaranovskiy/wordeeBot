package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"wordeeBot/internal/myErrors"
)

func handleCallbacks(b *TgBotModel, update tgbotapi.Update) {
	switch {

	case strings.Contains(update.CallbackQuery.Data, "myDictionaries"):

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleShowDictionaries(b, update))

	case strings.Contains(update.CallbackQuery.Data, "createDictionary"):

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleCreateDictionary(b, update))

	case strings.Contains(update.CallbackQuery.Data, "editDictionary"):

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleEditDictionary(b, update))

	case strings.Contains(update.CallbackQuery.Data, "mainMenu"):

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleMainMenu(b, update))

	case b.userLastCommand[update.CallbackQuery.From.ID] == "createDictionary_name":

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleCreateDictionaryColumns(b, update))

	case b.userLastCommand[update.CallbackQuery.From.ID] == "editDictionary":

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleChooseDictionaryForEdit(b, update))

	case b.userLastCommand[update.CallbackQuery.From.ID] == "editDictionary_name":

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleDictionaryEditing(b, update))

	case strings.Contains(update.CallbackQuery.Data, "addingWord"):

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleAddingWord(b, update))

	case strings.Contains(update.CallbackQuery.Data, "cancelWord"):

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleCancelledAddingWord(b, update))

	case b.userLastCommand[update.CallbackQuery.From.ID] == "myDictionaries":

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleSendParticularDictionary(b, update))

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
