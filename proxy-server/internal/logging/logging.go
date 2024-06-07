package logging

import (
	"fmt"
	"os"

	as "github.com/palavrapasse/aspirador/pkg"
)

const (
	telegramBotTokenEnvKey = "telegram_token"
	telegramChatIdEnvKey   = "telegram_chat_id"
)

const loggingFilePathEnvKey = "logging_fp"

var (
	telegramBotToken = os.Getenv(telegramBotTokenEnvKey)
	telegramChatId   = os.Getenv(telegramChatIdEnvKey)
)

var loggingFilePath = os.Getenv(loggingFilePathEnvKey)

var Aspirador as.Aspirador

func CreateAspiradorClients(serverAddress string) []as.Client {
	patternLayout := as.PatternLayout(fmt.Sprintf("[%s] [%s] %s %s %s.%s:%s : %s", as.LevelPattern, serverAddress, as.DatePattern, as.TimePattern, as.FileNamePattern, as.MethodPattern, as.LinePattern, as.MessagePattern))

	consoleClient := as.NewConsoleClient()
	consoleClient.SetPatternLayout(patternLayout)

	telegramClient := as.NewTelegramClient(telegramBotToken, telegramChatId, as.WARNING, as.ERROR)
	telegramClient.SetPatternLayout(patternLayout)

	fileClient, err := as.NewFileClient(loggingFilePath, as.INFO, as.WARNING, as.WARNING)
	fileClient.SetPatternLayout(patternLayout)

	if err != nil {
		return []as.Client{&consoleClient, &telegramClient}
	}

	return []as.Client{&consoleClient, &telegramClient, &fileClient}
}
