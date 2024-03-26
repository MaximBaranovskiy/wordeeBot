package messages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"wordeeBot/internal/keyboard"
)

func GetStartMessage(chat_id int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, "Здравствуйте. Выберите действие")
	keyboardMarkup := keyboard.CreateMainKeyboard()
	msg.ReplyMarkup = &keyboardMarkup
	return msg
}

func GetEditStartMessage(chat_id int64, msg_id int) tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(chat_id, msg_id, "Здравствуйте. Выберите действие")
	keyboardMarkup := keyboard.CreateMainKeyboard()
	msg.ReplyMarkup = &keyboardMarkup
	return msg
}

func GetMessageToAdding(chat_id int64, name string, columns []string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, "Давайте добавим слово в словарь "+name+".Он состоит из следующих столбцов:\nСлово\n")
	for _, column := range columns {
		msg.Text += column + "\n"
	}

	return msg
}

func GetMessageWithSomeInformationAboutAdding(chat_id int64) tgbotapi.MessageConfig {
	text := "Введите информацию о слове в следующем виде. Пример представлен в случае ввода всех столбцов. Но если у вас нет каких-либо," +
		" то просто сразу перейдите к следующему, если в некоторой строке вам надо несколько опций,то разделяйте их ;(точкой с запятой).Если вы " +
		"добавите информацию в столбец которого нет в словаре, то данная информация будет проигнорирована. Также придерживайтесь" +
		"порядка полей, указанного в примере:\n" +
		"*Слово:* Clever.\n" +
		"*Транскрипция:* /ˈklɛvər/.\n" +
		"*Перевод:*Умный *;* сообразительный.\n" +
		"*Синонимы:*intelligent *;* smart *;* bright.\n" +
		"*Антонимы:*foolish *;* stupid *;* unintelligent.\n" +
		"*Определение:*Having the ability to learn, understand, and apply knowledge quickly and effectively *;* mentally agile.\n" +
		"*Коллокации:*a clever idea *;* a clever solution *;* cleverly designed.\n" +
		"*Идиомы:*As clever as a fox."

	msg := tgbotapi.NewMessage(chat_id, text)
	msg.ParseMode = "Markdown"

	return msg
}

func GetMessageProblemsWithParsing(chat_id int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, "Проблемы с обработкой входной информации."+
		"Пожалуйста, перепроверьте ее на соответствие образцу")
	return msg
}

func GetCongratulateWithAdding(chat_id int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, "Слово успешно добавлено в словарь")
	return msg
}

func GetSQLErrorMessage(chat_id int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, "Проблемы при получении пользовательской информации."+
		" Пожалуйста, попробуйте еще раз.")
	return msg
}

func GetUnknownErrorMessage(chat_id int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, "Неизвестная ошибка. Пожалуйста попробуйте повторить действия")
	return msg
}

func GetMessageWordInDictionary(chat_id int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, "Данное слово уже есть у вас в словаре")
	return msg
}

func GetMessageWordNotInDictionary(chat_id int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, text)
	keyboardMarkup := keyboard.CreateKeyboardToAddingWord()
	msg.ReplyMarkup = &keyboardMarkup
	return msg
}

func GetEditMyDictionaryMessage(chatId int64, messageId int, dictionaries []string) tgbotapi.EditMessageTextConfig {
	msgEdit := tgbotapi.NewEditMessageText(chatId, messageId, "Список ваших словарей")
	if len(dictionaries) == 0 {
		msgEdit.Text = "Список ваших словарей пуст. Создайте хотя бы один"
	}
	keyboardMarkup := keyboard.CreateKeyboardToShowDict(dictionaries)
	msgEdit.ReplyMarkup = &keyboardMarkup
	return msgEdit
}

func GetCreationDictionaryMessage(chatId int64, messageId int) tgbotapi.EditMessageTextConfig {
	msgEdit := tgbotapi.NewEditMessageText(chatId, messageId, "Введите название нового словаря")
	keyboardMarkup := keyboard.CreateKeyboardCreationDictionary()
	msgEdit.ReplyMarkup = &keyboardMarkup
	return msgEdit
}

func GetMessageDictionaryIsCreated(chatId int64, messageId int) tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(chatId, messageId, "Словарь успешно создан")
	keyboardMarkup := keyboard.CreateMainKeyboard()
	msg.ReplyMarkup = &keyboardMarkup
	return msg
}

func GetMessageToEditDictionary(chatId int64, name string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, "Выберите одну из возможных опций для редактирования словаря "+name)
	keyboardMarkup := keyboard.CreateKeyboardWithEditingOptions(name)
	msg.ReplyMarkup = &keyboardMarkup
	return msg
}
