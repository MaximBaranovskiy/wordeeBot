package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"strings"
	"wordeeBot/internal/keyboard"
	"wordeeBot/internal/messages"
	"wordeeBot/internal/model/db"
	"wordeeBot/internal/myErrors"
	"wordeeBot/internal/parsing"
	"wordeeBot/internal/pdf"
)

func handleCallbacks(b *TgBotModel, update tgbotapi.Update) {
	id, err := b.userStorage.CheckUser(update.CallbackQuery.From.ID)
	if err != nil {
		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, err)
		return
	}

	if strings.Contains(update.CallbackQuery.Data, "myDictionaries") {

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleShowDictionaries(b, update, id))

	} else if strings.Contains(update.CallbackQuery.Data, "createDictionary") {

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleCreateDictionary(b, update))

	} else if strings.Contains(update.CallbackQuery.Data, "editDictionary") {

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleEditDicitonary(b, update, id))

	} else if strings.Contains(update.CallbackQuery.Data, "mainMenu") {

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleMainMenu(b, update))

	} else if b.userLastCommand[update.CallbackQuery.From.ID] == "createDictionary_columns" {

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleCreateDictionaryColumns(b, update, id))

	} else if b.userLastCommand[update.CallbackQuery.From.ID] == "editDictionary" {

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleChoosingDictionaryForEditing(b, update))

	} else if b.userLastCommand[update.CallbackQuery.From.ID] == "editDictionary_name" {

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleDictionaryEditing(b, update, id))

	} else if strings.Contains(update.CallbackQuery.Data, "addingWord") {

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleAddingWord(b, update))

	} else if strings.Contains(update.CallbackQuery.Data, "cancelWord") {

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleCancelledAddingWord(b, update))

	} else if b.userLastCommand[update.CallbackQuery.From.ID] == "myDictionaries" {

		myErrors.HandleError(b.bot, update.CallbackQuery.Message.Chat.ID, handleSendParticularDictionary(b, update, id))

	}
}

func handleMessages(b *TgBotModel, update tgbotapi.Update) {
	id, err := b.userStorage.CheckUser(update.Message.From.ID)
	if err != nil {
		myErrors.HandleError(b.bot, update.Message.Chat.ID, err)
		return
	}

	if strings.Contains(update.Message.Text, "start") {

		myErrors.HandleError(b.bot, update.Message.Chat.ID, handleStart(b, update))

	} else if b.userLastCommand[update.Message.From.ID] == "createDictionary_name" {

		myErrors.HandleError(b.bot, update.Message.Chat.ID, handleCreateDictionaryName(b, update, id))

	} else if strings.Contains(b.userLastCommand[update.Message.From.ID], "addingWord") {

		myErrors.HandleError(b.bot, update.Message.Chat.ID, handlePreparationAddingWord(b, update, id))

	}
}

func handleShowDictionaries(b *TgBotModel, update tgbotapi.Update, id int) error {
	b.userLastCommand[update.CallbackQuery.From.ID] = "myDictionaries"

	dictionaries, err := b.dictionaryStorage.GetNamesOfUserDictionaries(id)
	if err != nil {
		return err
	}

	msgEdit := messages.GetEditMyDictionaryMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID, dictionaries)
	b.bot.Send(msgEdit)

	return nil
}

func handleCreateDictionary(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.CallbackQuery.From.ID] = "createDictionary_name"

	msgEdit := messages.GetCreationDictionaryMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)
	b.bot.Send(msgEdit)

	return nil
}

func handleEditDicitonary(b *TgBotModel, update tgbotapi.Update, id int) error {
	b.userLastCommand[update.CallbackQuery.From.ID] = "editDictionary"

	dictionaries, err := b.dictionaryStorage.GetNamesOfUserDictionaries(id)
	if err != nil {
		return err
	}

	msgEdit := messages.GetEditMyDictionaryMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID, dictionaries)
	m, _ := b.bot.Send(msgEdit)

	b.userlastMessageID[update.CallbackQuery.From.ID] = m.MessageID

	return nil
}

func handleMainMenu(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.CallbackQuery.From.ID] = "mainMenu"

	msgEdit := messages.GetEditStartMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)

	b.bot.Send(msgEdit)

	return nil
}

func handleCreateDictionaryColumns(b *TgBotModel, update tgbotapi.Update, id int) error {
	if strings.Contains(update.CallbackQuery.Data, "confirm") {
		return handleConfirm(b, update, id)
	} else {
		handleNewColumn(b, update)
		return nil
	}
}

func handleConfirm(b *TgBotModel, update tgbotapi.Update, id int) error {
	ind := strings.Index(update.CallbackQuery.Data, "_")
	arr := b.tempColumnsForDictionaries[DictionaryIdentificator{User_id: update.CallbackQuery.From.ID,
		Name: update.CallbackQuery.Data[(ind + 1):]}]

	err := b.dictionaryStorage.AddDictionary(update.CallbackQuery.Data[(ind+1):], id, arr)
	if err != nil {
		return err
	}

	msg := messages.GetMessageDictionaryIsCreated(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)
	m, _ := b.bot.Send(msg)

	b.userlastMessageID[update.CallbackQuery.From.ID] = m.MessageID

	return nil
}

func handleNewColumn(b *TgBotModel, update tgbotapi.Update) {
	ind := strings.Index(update.CallbackQuery.Data, "_")
	arr := b.tempColumnsForDictionaries[DictionaryIdentificator{User_id: update.CallbackQuery.From.ID,
		Name: update.CallbackQuery.Data[(ind + 1):]}]

	if !checkColumn(&arr, update.CallbackQuery.Data[0:ind]) {
		arr = append(arr, update.CallbackQuery.Data[0:ind])
		b.tempColumnsForDictionaries[DictionaryIdentificator{User_id: update.CallbackQuery.From.ID,
			Name: update.CallbackQuery.Data[(ind + 1):]}] = arr
	}
}

func handleStart(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.Message.From.ID] = "start"

	msg := messages.GetStartMessage(update.Message.Chat.ID)
	m, _ := b.bot.Send(msg)

	b.userlastMessageID[update.Message.From.ID] = m.MessageID
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

func handleCreateDictionaryName(b *TgBotModel, update tgbotapi.Update, id int) error {
	b.userLastCommand[update.Message.From.ID] = "createDictionary_columns"

	ok, err := b.dictionaryStorage.CheckDicitonary(strings.ToLower(update.Message.Text), id)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Давайте определимся с информацией, которую вы хотите хранить.Выберите нужные столбцы и нажмите подтвердить.Если вы не выберете ничего. То по умолчанию словарь будет состоять лишь из самого слова")
	if ok {
		msg.Text = "Данный словарь уже существует. Пожалуйста перепроверьте введенное название"
		msg.ReplyMarkup = keyboard.CreateMainKeyboard()
	} else {
		b.tempColumnsForDictionaries[DictionaryIdentificator{User_id: update.Message.From.ID,
			Name: strings.ToLower(update.Message.Text)}] = make([]string, 0)
		keyboardMarkup := keyboard.CreateKeyboardWithColumns(strings.ToLower(update.Message.Text))
		msg.ReplyMarkup = &keyboardMarkup
	}

	editMsg := messages.GetEditStartMessage(update.Message.Chat.ID, b.userlastMessageID[update.Message.From.ID])
	b.bot.Send(editMsg)

	m, _ := b.bot.Send(msg)
	b.userlastMessageID[update.Message.From.ID] = m.MessageID
	return nil
}

func handleChoosingDictionaryForEditing(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.CallbackQuery.From.ID] = "editDictionary_name"

	msg := messages.GetMessageToEditDictionary(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Data)

	msgEdit := messages.GetEditStartMessage(update.CallbackQuery.Message.Chat.ID, b.userlastMessageID[update.CallbackQuery.From.ID])
	b.bot.Send(msgEdit)

	m, _ := b.bot.Send(msg)
	b.userlastMessageID[update.CallbackQuery.From.ID] = m.MessageID

	return nil
}

func handleDictionaryEditing(b *TgBotModel, update tgbotapi.Update, id int) error {
	ind := strings.Index(update.CallbackQuery.Data, "_")
	name := update.CallbackQuery.Data[(ind + 1):]
	columns, err := b.dictionaryStorage.GetNamesOfDictionaryColumns(id, name)
	if err != nil {
		return err
	}
	b.userLastCommand[update.CallbackQuery.From.ID] = "addingWord_" + name

	msgEdit := messages.GetEditStartMessage(update.CallbackQuery.Message.Chat.ID, b.userlastMessageID[update.CallbackQuery.From.ID])
	b.bot.Send(msgEdit)

	msg := messages.GetMessageToAdding(update.CallbackQuery.Message.Chat.ID, name, columns)
	b.bot.Send(msg)

	msg = messages.GetMessageWithSomeInformationAboutAdding(update.CallbackQuery.Message.Chat.ID)
	b.bot.Send(msg)
	return nil
}

func handlePreparationAddingWord(b *TgBotModel, update tgbotapi.Update, id int) error {
	text := update.Message.Text
	ind := strings.Index(b.userLastCommand[update.Message.From.ID], "_")
	dictionaryId, err := b.dictionaryStorage.GetDictionaryId(id, b.userLastCommand[update.Message.From.ID][(ind+1):])
	if err != nil {
		return err
	}

	word, err := parsing.ParseInfoForAdding(text, dictionaryId)
	if err != nil {
		return err
	}

	if word != nil {
		ok, err := b.wordsStorage.CheckWord(dictionaryId, word.Writing)
		if err != nil {
			return err
		}

		if ok {
			b.bot.Send(messages.GetMessageWordInDictionary(update.Message.Chat.ID))
			b.bot.Send(messages.GetStartMessage(update.Message.Chat.ID))
		} else {
			b.bot.Send(messages.GetMessageWordNotInDictionary(update.Message.Chat.ID, word.ToString()))
			b.tempStorageForAddingWords[update.Message.From.ID] = word
		}
	}

	return nil
}

func handleAddingWord(b *TgBotModel, update tgbotapi.Update) error {
	word := b.tempStorageForAddingWords[update.CallbackQuery.From.ID]
	if err := b.wordsStorage.AddWord(word); err != nil {
		return err
	}

	b.bot.Send(messages.GetCongratulateWithAdding(update.CallbackQuery.Message.Chat.ID))
	b.bot.Send(messages.GetStartMessage(update.CallbackQuery.Message.Chat.ID))

	return nil
}

func handleCancelledAddingWord(b *TgBotModel, update tgbotapi.Update) error {
	b.bot.Send(messages.GetStartMessage(update.CallbackQuery.Message.Chat.ID))
	return nil
}

func handleSendParticularDictionary(b *TgBotModel, update tgbotapi.Update, id int) error {
	dictionaryName := update.CallbackQuery.Data
	dictionaryId, err := b.dictionaryStorage.GetDictionaryId(id, dictionaryName)
	if err != nil {
		return err
	}

	words := make([]db.Word, 0)
	words, err = b.wordsStorage.GetAllDictionaryWords(dictionaryId)
	if err != nil {
		return err
	}

	fields := make(map[string]bool)
	fields["Перевод"] = false
	fields["Транскрипция"] = false
	fields["Синонимы"] = false
	fields["Aнтонимы"] = false
	fields["Определение"] = false
	fields["Коллокации"] = false
	fields["Идиомы"] = false

	columns, err := b.dictionaryStorage.GetNamesOfDictionaryColumns(id, dictionaryName)
	if err != nil {
		return err
	}

	for _, column := range columns {
		fields[column] = true
	}

	err = pdf.MakeDictionaryPDF(dictionaryName, words, fields)
	if err != nil {
		return err
	}

	fileBytes, err := ioutil.ReadFile("dictionary.pdf")
	if err != nil {
		return err
	}

	file := tgbotapi.FileBytes{Name: "dictionary.pdf", Bytes: fileBytes}

	msg := tgbotapi.NewDocument(update.CallbackQuery.Message.Chat.ID, file)

	b.bot.Send(msg)

	return nil
}
