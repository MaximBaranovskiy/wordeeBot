package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"wordeeBot/internal/keyboard"
	"wordeeBot/internal/messages"
	"wordeeBot/internal/model/db"
	"wordeeBot/internal/pdf"
)

func handleShowDictionaries(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.CallbackQuery.From.ID] = "myDictionaries"

	dictionaries, err := b.dictionaryStorage.GetNamesOfUserDictionaries(update.SentFrom().ID)
	if err != nil {
		return err
	}

	msgEdit := messages.GetEditMyDictionaryMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID, dictionaries)
	b.bot.Send(msgEdit)

	return nil
}

func handleCreateDictionary(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.CallbackQuery.From.ID] = "createDictionary"

	msgEdit := messages.GetCreateDictionaryMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)
	b.bot.Send(msgEdit)

	return nil
}

func handleEditDictionary(b *TgBotModel, update tgbotapi.Update) error {
	userId := update.SentFrom().ID
	chatId := update.CallbackQuery.Message.Chat.ID
	msgId := update.CallbackQuery.Message.MessageID

	b.userLastCommand[userId] = "editDictionary"

	dictionaries, err := b.dictionaryStorage.GetNamesOfUserDictionaries(userId)
	if err != nil {
		return err
	}

	msgEdit := messages.GetEditMyDictionaryMessage(chatId,
		msgId, dictionaries)
	m, _ := b.bot.Send(msgEdit)

	b.userlastMessageID[userId] = m.MessageID

	return nil
}

func handleMainMenu(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.CallbackQuery.From.ID] = "mainMenu"

	msgEdit := messages.GetEditStartMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)
	b.bot.Send(msgEdit)

	return nil
}

func handleCreateDictionaryColumns(b *TgBotModel, update tgbotapi.Update) error {
	if strings.Contains(update.CallbackQuery.Data, "confirm") {
		return handleConfirm(b, update)
	} else {
		handleColumn(b, update)
		return nil
	}
}

func handleConfirm(b *TgBotModel, update tgbotapi.Update) error {
	ind := strings.Index(update.CallbackQuery.Data, "_")
	arr := b.tempColumnsForDictionaries[DictionaryIdentificator{UserID: update.CallbackQuery.From.ID,
		Name: update.CallbackQuery.Data[(ind + 1):]}]

	err := b.dictionaryStorage.AddDictionary(update.CallbackQuery.Data[(ind+1):], arr, update.CallbackQuery.From.ID)
	if err != nil {
		return err
	}

	msg := messages.GetMessageDictionaryIsCreated(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)
	b.bot.Send(msg)

	return nil
}

func handleColumn(b *TgBotModel, update tgbotapi.Update) {
	ind := strings.Index(update.CallbackQuery.Data, "_")
	dictionaryName := update.CallbackQuery.Data[(ind + 1):]
	di := DictionaryIdentificator{UserID: update.CallbackQuery.From.ID,
		Name: dictionaryName}

	arr := b.tempColumnsForDictionaries[di]

	if !checkColumn(&arr, update.CallbackQuery.Data[0:ind]) {
		arr = append(arr, update.CallbackQuery.Data[0:ind])
		b.tempColumnsForDictionaries[di] = arr
	} else {
		for i, item := range arr {
			if item == update.CallbackQuery.Data[0:ind] {
				arr = append(arr[:i], arr[i+1:]...)
				b.tempColumnsForDictionaries[di] = arr
			}
		}
	}

	msgEdit := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID,
		keyboard.CreateKeyboardWithColumns(dictionaryName, arr))

	b.bot.Send(msgEdit)
}

func handleChooseDictionaryForEdit(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.CallbackQuery.From.ID] = "editDictionary_name"

	msg := messages.GetMessageToEditDictionary(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Data)
	b.bot.Send(msg)

	msgEdit := messages.GetEditStartMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	b.bot.Send(msgEdit)

	return nil
}

func handleDictionaryEditing(b *TgBotModel, update tgbotapi.Update) error {
	ind := strings.Index(update.CallbackQuery.Data, "_")
	name := update.CallbackQuery.Data[(ind + 1):]
	b.userLastCommand[update.CallbackQuery.From.ID] = "addingWord_" + name

	columns, err := b.dictionaryStorage.GetNamesOfDictionaryColumns(update.SentFrom().ID, name)
	if err != nil {
		return err
	}

	msgEdit := messages.GetEditStartMessage(update.CallbackQuery.Message.Chat.ID, b.userlastMessageID[update.CallbackQuery.From.ID])
	b.bot.Send(msgEdit)

	b.tempStorageForEditingWords[update.CallbackQuery.From.ID] = &StructForAddingWord{
		Word:    new(db.Word),
		Columns: columns,
		Count:   1,
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите слово:")
	b.bot.Send(msg)
	return nil
}

func handleAddingWord(b *TgBotModel, update tgbotapi.Update) error {
	temp := b.tempStorageForEditingWords[update.CallbackQuery.From.ID]
	if err := b.wordsStorage.AddWord(temp.Word); err != nil {
		return err
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Слово успешно добавлено", keyboard.CreateMainKeyboard())
	b.bot.Send(msg)

	return nil
}

func handleCancelledAddingWord(b *TgBotModel, update tgbotapi.Update) error {
	b.bot.Send(tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Добавление слова отменено", keyboard.CreateMainKeyboard()))
	return nil
}

func handleSendParticularDictionary(b *TgBotModel, update tgbotapi.Update) error {
	fields := getFields()
	dictionaryName := update.CallbackQuery.Data
	dictionaryId, err := b.dictionaryStorage.GetDictionaryId(update.SentFrom().ID, dictionaryName)
	if err != nil {
		return err
	}

	words := make([]db.Word, 0)
	words, err = b.wordsStorage.GetAllDictionaryWords(dictionaryId)
	if err != nil {
		return err
	}

	columns, err := b.dictionaryStorage.GetNamesOfDictionaryColumns(update.SentFrom().ID, dictionaryName)
	if err != nil {
		return err
	}

	for _, column := range columns {
		fields[column] = true
	}

	err, file := pdf.MakeDictionaryPDF(dictionaryName, words, fields)
	if err != nil {
		return err
	}

	b.bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
	b.bot.Send(tgbotapi.NewDocument(update.CallbackQuery.Message.Chat.ID, file))
	b.bot.Send(messages.GetStartMessage(update.CallbackQuery.Message.Chat.ID))

	return nil
}

func checkColumn(arr *[]string, column string) bool {
	for _, item := range *arr {
		if item == column {
			return true
		}
	}
	return false
}

func getFields() map[string]bool {
	fields := make(map[string]bool)
	fields["Перевод"] = false
	fields["Транскрипция"] = false
	fields["Синонимы"] = false
	fields["Aнтонимы"] = false
	fields["Определение"] = false
	fields["Коллокации"] = false
	fields["Идиомы"] = false

	return fields
}
