package features

import (
	"fmt"
	discordgo "github.com/bwmarrin/discord.go"
	"github.com/foxtrot/scuzzy/models"
	"time"
)

func (f *Features) PrintError(component string, error string) {
	fmt.Printf("Error: %s: %s\n", component, error)
}

func (f *Features) CreateDefinedEmbed(title string, desc string, status string, user *discordgo.User) *discordgo.MessageEmbed {
	msgColor := 0x000000

	switch status {
	case "error":
		msgColor = 0xCC0000
		break
	case "success":
		msgColor = 0x00CC00
		break
	default:
		msgColor = 0xFFA500
	}

	ftrText := "Something broken? Tell foxtrot#1337"
	if user != nil {
		ftrText += "\nRequested by " + user.Username + "#" + user.Discriminator
	}

	ftr := discordgo.MessageEmbedFooter{
		Text:         ftrText,
		IconURL:      "https://cdn.discordapp.com/avatars/514163441548656641/a4ede220fea0ad8872b86f3eebc45524.png?size=128",
		ProxyIconURL: "",
	}

	msg := discordgo.MessageEmbed{
		URL:         "",
		Type:        "",
		Title:       title,
		Description: desc,
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       msgColor,
		Footer:      &ftr,
		Image:       nil,
		Thumbnail:   nil,
		Video:       nil,
		Provider:    nil,
		Author:      nil,
		Fields:      nil,
	}

	return &msg
}

func (f *Features) CreateCustomEmbed(embedData *models.CustomEmbed) *discordgo.MessageEmbed {
	var ftr discordgo.MessageEmbedFooter
	var img discordgo.MessageEmbedImage
	var thm discordgo.MessageEmbedThumbnail
	var prv discordgo.MessageEmbedProvider
	var atr discordgo.MessageEmbedAuthor

	ftr.Text = embedData.FooterText
	ftr.IconURL = embedData.FooterImageURL

	img.URL = embedData.ImageURL
	img.Height = embedData.ImageH
	img.Width = embedData.ImageW

	thm.URL = embedData.ThumbnailURL
	thm.Height = embedData.ThumbnailH
	thm.Width = embedData.ThumbnailW

	prv.Name = embedData.ProviderText
	prv.URL = embedData.ProviderURL

	atr.Name = embedData.AuthorText
	atr.URL = embedData.AuthorURL
	atr.IconURL = embedData.AuthorImageURL

	msg := discordgo.MessageEmbed{
		URL:         embedData.URL,
		Type:        embedData.Type,
		Title:       embedData.Title,
		Description: embedData.Desc,
		Timestamp:   "",
		Color:       embedData.Color,
		Footer:      &ftr,
		Image:       &img,
		Thumbnail:   &thm,
		Video:       nil,
		Provider:    &prv,
		Author:      &atr,
		Fields:      nil,
	}

	return &msg
}
