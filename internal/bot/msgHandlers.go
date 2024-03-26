package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"wordeeBot/internal/keyboard"
	"wordeeBot/internal/messages"
	"wordeeBot/internal/parsing"
)

func handleStart(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.Message.From.ID] = "start"

	msg := messages.GetStartMessage(update.Message.Chat.ID)
	m, _ := b.bot.Send(msg)

	b.userlastMessageID[update.Message.From.ID] = m.MessageID
	return nil
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
		b.tempColumnsForDictionaries[DictionaryIdentificator{UserID: update.Message.From.ID,
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
