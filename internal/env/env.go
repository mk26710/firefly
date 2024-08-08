package env

import (
	"fmt"
	"os"
)

func MustGetEnv(key string) string {
	s, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("Environment variable %s must be set!\n", key))
	}
	return s
}

func BotToken() string {
	return MustGetEnv("BOT_TOKEN")
}

func SaucenaoToken() string {
	return MustGetEnv("SAUCENAO_TOKEN")
}
