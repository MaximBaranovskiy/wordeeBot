package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateMainKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Мои словари📚", "myDictionaries"),
			tgbotapi.NewInlineKeyboardButtonData("Создать словарь📖", "createDictionary"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Редактировать словарь✏️", "editDictionary"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изучение слов🏫", "studyWords"),
		),
	)
}

func CreateKeyboardToShowDict(names []string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, name := range names {
		row := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(name, name)}
		rows = append(rows, row)
	}

	row := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("На главное меню↩️", "mainMenu")}
	rows = append(rows, row)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func CreateKeyboardWithColumns(name string, arr []string) tgbotapi.InlineKeyboardMarkup {
	columnName2Sign := getColumnName2Sign()

	for _, columnData := range arr {
		columnName2Sign[columnData2Name[columnData]] = "✅"
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Транскрипция"+columnName2Sign["Транскрипция"], "transcription"+"_"+name),
			tgbotapi.NewInlineKeyboardButtonData("Перевод"+columnName2Sign["Перевод"], "translation"+"_"+name),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Синонимы"+columnName2Sign["Синонимы"], "synonyms"+"_"+name),
			tgbotapi.NewInlineKeyboardButtonData("Антонимы"+columnName2Sign["Антонимы"], "antonyms"+"_"+name),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Определение"+columnName2Sign["Определение"], "definition"+"_"+name),
			tgbotapi.NewInlineKeyboardButtonData("Коллокации"+columnName2Sign["Коллокации"], "collocations"+"_"+name),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Идиомы"+columnName2Sign["Идиомы"], "idioms"+"_"+name),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить🆗", "confirm"+"_"+name),
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
			tgbotapi.NewInlineKeyboardButtonData("ДА🟢", "addingWord"),
			tgbotapi.NewInlineKeyboardButtonData("НЕТ🔴", "cancelWord"),
		),
	)
}

var columnData2Name = map[string]string{
	"transcription": "Транскрипция",
	"translation":   "Перевод",
	"synonyms":      "Синонимы",
	"antonyms":      "Антонимы",
	"definition":    "Определение",
	"collocations":  "Коллокации",
	"idioms":        "Идиомы",
}

func getColumnName2Sign() map[string]string {
	return map[string]string{
		"Транскрипция": "❌",
		"Перевод":      "❌",
		"Синонимы":     "❌",
		"Антонимы":     "❌",
		"Определение":  "❌",
		"Коллокации":   "❌",
		"Идиомы":       "❌",
	}
}
