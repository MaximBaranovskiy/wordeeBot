package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func CreateMainKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Мои словари", "myDictionaries"),
			tgbotapi.NewInlineKeyboardButtonData("Создать словарь", "createDictionary"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Редактировать словарь", "editDictionary"),
		),
	)
}

func CreateKeyboardToShowDict(names []string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, name := range names {
		row := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(name, name)}
		rows = append(rows, row)
	}

	row := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("На главное меню", "mainMenu")}
	rows = append(rows, row)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)

}

func CreateKeyboardWithColumns(name string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Транскрипция", "transcription"+"_"+name),
			tgbotapi.NewInlineKeyboardButtonData("Перевод", "translation"+"_"+name),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Синонимы", "synonyms"+"_"+name),
			tgbotapi.NewInlineKeyboardButtonData("Антонимы", "antonyms"+"_"+name),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Определение", "definition"+"_"+name),
			tgbotapi.NewInlineKeyboardButtonData("Коллокации", "collocations"+"_"+name),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Идиомы", "idioms"+"_"+name),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "confirm"+"_"+name),
		),
	)
}

func CreateKeyboardWithEditingOptions(name string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить слово", "addWord"+"_"+name),
		),
	)
}

func CreateKeyboardToAddingWord() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ДА", "addingWord"),
			tgbotapi.NewInlineKeyboardButtonData("НЕТ", "cancelWord"),
		),
	)
}

func CreateKeyboardCreationDictionary() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("На главное меню", "mainMenu"),
	),
	)
}
