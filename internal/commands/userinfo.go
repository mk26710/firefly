package commands

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type UserInfoHandler struct{}

func NewUserInfoHandler() UserInfoHandler {
	return UserInfoHandler{}
}

func (h *UserInfoHandler) Meta() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "userinfo",
		Description: "Prints info about Discord users",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Target user",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "public",
				Description: "Whether if you'd like to show the response to everyone",
				Required:    false,
			},
		},
	}
}

func (h *UserInfoHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options

	var public bool
	var target *discordgo.User

	for _, option := range options {
		if option.Type == discordgo.ApplicationCommandOptionUser && option.Name == "user" {
			target = option.UserValue(s)
		}

		if option.Type == discordgo.ApplicationCommandOptionBoolean && option.Name == "public" {
			public = option.BoolValue()
		}
	}

	if target == nil {
		return errors.New("target user is nil")
	}

	var embed discordgo.MessageEmbed

	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: target.AvatarURL("1024")}

	if len(target.GlobalName) > 0 {
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:  "Name",
				Value: target.GlobalName,
			},
		)
	}

	embed.Fields = append(embed.Fields,
		&discordgo.MessageEmbedField{
			Name:  "Username",
			Value: target.String(),
		},
		&discordgo.MessageEmbedField{
			Name:  "ID",
			Value: target.ID,
		},
		&discordgo.MessageEmbedField{
			Name:  "Avatar",
			Value: fmt.Sprintf("[Avatar URL](%s)", target.AvatarURL("4096")),
		},
	)

	if len(target.Banner) > 1 {
		embed.Image = &discordgo.MessageEmbedImage{URL: target.BannerURL("4096")}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Banner",
			Value: fmt.Sprintf("[Banner URL](%v)", target.BannerURL("4096")),
		})
	}

	var flags discordgo.MessageFlags

	if !public {
		flags = discordgo.MessageFlagsEphemeral
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  flags,
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
}
