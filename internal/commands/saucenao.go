package commands

import (
	"firefly/internal/utils/sauce"
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type SaucenaoHandler struct{}

func (h *SaucenaoHandler) Meta() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "saucenao",
		Description: "Performs a lookup on saucenao",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "Image URL to use in saucenao search",
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

func (h *SaucenaoHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options

	var public bool
	var urlRaw string

	for _, option := range options {
		if option.Type == discordgo.ApplicationCommandOptionString && option.Name == "url" {
			urlRaw = option.StringValue()
		}

		if option.Type == discordgo.ApplicationCommandOptionBoolean && option.Name == "public" {
			public = option.BoolValue()
		}
	}

	queryUrl, err := url.Parse(urlRaw)
	if err != nil {
		return err
	}

	var flags discordgo.MessageFlags
	if !public {
		flags = discordgo.MessageFlagsEphemeral
	}

	ch, err := s.Channel(i.ChannelID)
	if err != nil {
		return err
	}

	qo := []sauce.QueryOption{sauce.WithMaxResults(25), sauce.WithoutNSFW()}
	if ch.NSFW {
		qo[1] = sauce.WithNSFW()
	}

	results, err := sauce.Query(queryUrl.String(), qo...)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   flags,
				Content: "There are no results.",
			},
		})
	}

	topResult := results[0]

	qe := discordgo.MessageEmbed{}
	qe.Title = "Requested Image"
	qe.Image = &discordgo.MessageEmbedImage{URL: queryUrl.String()}

	re := discordgo.MessageEmbed{}
	re.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: topResult.Header.ThumbnailURL}
	re.Title = topResult.Header.IndexName

	if len(topResult.Header.Similarity) > 0 {
		re.Fields = append(re.Fields, &discordgo.MessageEmbedField{
			Name:  "Similarity",
			Value: topResult.Header.Similarity,
		})
	}

	if len(topResult.Data.Characters) > 0 {
		re.Fields = append(re.Fields, &discordgo.MessageEmbedField{
			Name:  "Characters",
			Value: topResult.Data.Characters,
		})
	}

	if len(topResult.Data.Material) > 0 {
		re.Fields = append(re.Fields, &discordgo.MessageEmbedField{
			Name:  "Material",
			Value: topResult.Data.Material,
		})
	}

	if len(topResult.Data.Creators) > 0 {
		re.Fields = append(re.Fields, &discordgo.MessageEmbedField{
			Name:  "Creators",
			Value: strings.Join(topResult.Data.Creators, "\n"),
		})
	}

	// this part is absolute shitstorm of copies
	// because apparently you cannot use an array
	// of type which is an implementation of interface
	// when you assign that array to a struct field
	// which wants an array where each element is
	// an implementation of the interfact
	//
	// very convenient and simple language!
	rows := []discordgo.ActionsRow{{}} // this should be possible to be used as []discordgo.MessageComponent value

	if len(topResult.Data.SourceURL) > 0 {
		rows[0].Components = append(rows[0].Components, discordgo.Button{
			Label: "Source",
			URL:   topResult.Data.SourceURL,
			Style: discordgo.LinkButton,
		})
	}

	for i, extURL := range topResult.Data.ExtURLs {
		if len(rows[len(rows)-1].Components) == 5 {
			rows = append(rows, discordgo.ActionsRow{})
		}

		if len(extURL) > 0 {
			rows[len(rows)-1].Components = append(rows[len(rows)-1].Components, discordgo.Button{
				URL:   extURL,
				Label: fmt.Sprintf("Mirror #%v", i+1),
				Style: discordgo.LinkButton,
			})
		}
	}

	components := []discordgo.MessageComponent{}
	for _, row := range rows {
		components = append(components, row)
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      flags,
			Embeds:     []*discordgo.MessageEmbed{&qe, &re},
			Components: components,
		},
	})
}
