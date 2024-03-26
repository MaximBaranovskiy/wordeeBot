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

func GetMessageProblemsWithParsing(chat_id int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, "Проблемы с обработкой входной информации."+
		"Пожалуйста, перепроверьте ее на соответствие образцу")
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

func GetCreateDictionaryMessage(chatId int64, messageId int) tgbotapi.EditMessageTextConfig {
	msgEdit := tgbotapi.NewEditMessageText(chatId, messageId, "Введите название нового словаря")
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

var Column2Text = map[string]string{
	"Слово":        "Введите слово:",
	"Транскрипция": "Введите транскрипцию",
	"Перевод":      "Введите перевод.Если есть несколько вариантов, то разделите ;",
	"Синонимы":     "Введите синонимы.Если есть несколько вариантов, то разделите ;",
	"Антонимы":     "Введите антонимы.Если есть несколько вариантов, то разделите ;",
	"Определение":  "Введите определение.Если есть нексколько вариантов, то разделите ;",
	"Коллокации":   "Введите коллокации.Если есть несколько вариантов, то разделите ;",
	"Идиомы":       "Введите идиомы.Если есть несколько вариантов, то разделите ;",
}
