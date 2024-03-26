package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
	"wordeeBot/internal/keyboard"
	"wordeeBot/internal/messages"
	"wordeeBot/internal/model/db"
)

func handleStart(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.Message.From.ID] = "start"

	msg := messages.GetStartMessage(update.Message.Chat.ID)
	m, _ := b.bot.Send(msg)

	b.userlastMessageID[update.Message.From.ID] = m.MessageID
	return nil
}

func handleCreateDictionaryName(b *TgBotModel, update tgbotapi.Update) error {
	b.userLastCommand[update.Message.From.ID] = "createDictionary_name"

	ok := validateName(update.Message.Text)
	if !ok {
		b.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введено некорректное имя словаря.Имя может содержать только английсские и русские буквы,цифры и знак нижнего подчеркивания"))
		b.bot.Send(messages.GetStartMessage(update.Message.Chat.ID))
		return nil
	}

	ok, err := b.dictionaryStorage.CheckDicitonary(strings.ToLower(update.Message.Text), update.SentFrom().ID)
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
		keyboardMarkup := keyboard.CreateKeyboardWithColumns(strings.ToLower(update.Message.Text), []string{})
		msg.ReplyMarkup = &keyboardMarkup
	}

	editMsg := messages.GetEditStartMessage(update.Message.Chat.ID, b.userlastMessageID[update.Message.From.ID])
	b.bot.Send(editMsg)

	m, _ := b.bot.Send(msg)
	b.userlastMessageID[update.Message.From.ID] = m.MessageID
	return nil
}

func handlePreparationAddingWord(b *TgBotModel, update tgbotapi.Update) error {
	ind := strings.Index(b.userLastCommand[update.Message.From.ID], "_")
	dictionaryId, err := b.dictionaryStorage.GetDictionaryId(update.SentFrom().ID, b.userLastCommand[update.Message.From.ID][(ind+1):])
	if err != nil {
		return err
	}

	text := update.Message.Text
	temp := b.tempStorageForEditingWords[update.SentFrom().ID]

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, " ")

	if temp.Count >= len(temp.Columns) {
		newValue(temp.Word, temp.Columns[temp.Count-1], text)
		b.bot.Send(messages.GetMessageWordNotInDictionary(update.Message.Chat.ID, temp.Word.ToString()))
		return nil
	} else {
		if temp.Count == 1 {
			temp.Word.DictionaryId = dictionaryId

			ok, err := b.wordsStorage.CheckWord(dictionaryId, text)
			if err != nil {
				return err
			}

			if ok {
				b.bot.Send(messages.GetMessageWordInDictionary(update.Message.Chat.ID))
				b.bot.Send(messages.GetStartMessage(update.Message.Chat.ID))
				return nil
			}
		}
		newValue(temp.Word, temp.Columns[temp.Count-1], text)
		msg.Text = messages.Column2Text[temp.Columns[temp.Count]]
		temp.Count++
		b.bot.Send(msg)

	}

	return nil
}

func newValue(word *db.Word, columnName, value string) {
	switch columnName {
	case "Слово":
		word.Writing = value
	case "Транскрипция":
		word.Transcription = value
	case "Перевод":
		word.Translation = value
	case "Синонимы":
		word.Synonyms = value
	case "Антонимы":
		word.Antonyms = value
	case "Определение":
		word.Definition = value
	case "Коллокации":
		word.Collocations = value
	case "Идиомы":
		word.Idioms = value
	}

}

func validateName(name string) bool {
	validNamePattern := `^[a-zA-Zа-яА-Я0-9_]+$`
	matched, _ := regexp.MatchString(validNamePattern, name)
	return matched
}
