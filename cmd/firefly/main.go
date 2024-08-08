package main

import (
	"firefly/internal/commands"
	"firefly/internal/env"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type FireflyCommand interface {
	Meta() *discordgo.ApplicationCommand
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

var appCommands = map[string]FireflyCommand{
	"userinfo": &commands.UserInfoHandler{},
	"saucenao": &commands.SaucenaoHandler{},
}

func CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if s == nil || i == nil {
		panic("session or interaction is nil")
	}

	cmd, exists := appCommands[i.ApplicationCommandData().Name]
	if !exists {
		slog.Info("Unknown command was invoked - ", i.ApplicationCommandData().Name, "")
		return

	}

	err := cmd.Handle(s, i)
	if err != nil {
		log.Printf("an error happened while executing %s:\n%v\n", cmd.Meta().Name, err)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: ":warning: Unexpected error has occured during execution of the command.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// cfg := gorm.Config{TranslateError: true}
	// dsn := os.Getenv("POSTGRES_DSN")

	// _, err := gorm.Open(postgres.Open(dsn), &cfg)
	// if err != nil {
	// 	log.Fatal("Error while connecting to database")
	// }

	bot, err := discordgo.New("Bot " + env.BotToken())
	if err != nil {
		log.Fatal("Error connecting to Discord")
	}

	if err := bot.Open(); err != nil {
		log.Fatalf("Cannot open Discord session: %v", err)
	}

	bot.AddHandler(CommandHandler)

	for _, cmd := range appCommands {
		if _, err := bot.ApplicationCommandCreate(bot.State.User.ID, "", cmd.Meta()); err != nil {
			log.Printf("Error:\n%v\n", err)
		}
	}

	defer bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	log.Println("Gracefully shutting down.")
}
