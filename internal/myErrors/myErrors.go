package myErrors

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"wordeeBot/internal/messages"
)

var (
	ErrSql           = errors.New("Problems with db connection")
	ErrParseWordInfo = errors.New("Problems with parsing input info about word")
)

func HandleError(bot *tgbotapi.BotAPI, chat_id int64, err error) {
	if err == nil {
		return
	}

	switch err {
	case ErrSql:
		bot.Send(messages.GetSQLErrorMessage(chat_id))
		bot.Send(messages.GetStartMessage(chat_id))
	case ErrParseWordInfo:
		bot.Send(messages.GetMessageProblemsWithParsing(chat_id))
		bot.Send(messages.GetStartMessage(chat_id))
	default:
		bot.Send(messages.GetUnknownErrorMessage(chat_id))
		bot.Send(messages.GetStartMessage(chat_id))
	}
}
