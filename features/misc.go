package features

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discord.go"
	"github.com/foxtrot/scuzzy/models"
	"strconv"
	"strings"
	"time"
)

func (f *Features) handleCat(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if f.Auth.CheckAdminRole(m.Member) {
		_, _ = s.ChannelMessageSend(m.ChannelID, "https://giphy.com/gifs/cat-cute-no-rCxogJBzaeZuU")
		_ = s.ChannelMessageDelete(m.ChannelID, m.ID)
	}

	return nil
}

func (f *Features) handlePing(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var r *discordgo.Message
	var err error

	if !f.Auth.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	} else {
		msg := f.CreateDefinedEmbed("Ping", "Pong", "success")
		r, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
		if err != nil {
			return err
		}
	}

	time.Sleep(5 * time.Second)

	err = s.ChannelMessageDelete(m.ChannelID, r.ID)
	if err != nil {
		return err
	}
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	desc := "**Source**:   https://github.com/foxtrot/scuzzy\n"
	desc += "**Language**: Go\n"
	desc += "**Commands**: See `" + f.Config.CommandKey + "help`\n\n\n"

	gm, err := s.GuildMember(f.Config.GuildID, s.State.User.ID)
	if err != nil {
		return err
	}

	d := models.CustomEmbed{
		Title:          "Scuzzy Information",
		Desc:           desc,
		ImageURL:       "",
		ImageH:         100,
		ImageW:         100,
		Color:          0xFFA500,
		URL:            "",
		Type:           "",
		Timestamp:      "",
		FooterText:     "Made with  ❤  by Foxtrot",
		FooterImageURL: "https://cdn.discordapp.com/avatars/514163441548656641/a4ede220fea0ad8872b86f3eebc45524.png",
		ThumbnailURL:   gm.User.AvatarURL(""),
		ThumbnailH:     150,
		ThumbnailW:     150,
		ProviderURL:    "",
		ProviderText:   "",
		AuthorText:     "",
		AuthorURL:      "",
		AuthorImageURL: "",
	}

	msg := f.CreateCustomEmbed(&d)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) error {
	desc := "**Available Commands**\n"
	desc += "__Misc__\n"
	desc += "`help` - This help dialog\n"
	desc += "`info` - Display Scuzzy info\n"
	desc += "`md` - Display Discord markdown information\n"

	desc += "\n__User Settings__\n"
	desc += "`colors` - Available color roles\n"
	desc += "`color` - Set an available color role\n"

	desc += "\n__Conversion Helpers__\n"
	desc += "`ctof` - Convert Celsius to Farenheit\n"
	desc += "`ftoc` - Convert Farenheit to Celsius\n"
	desc += "`metofe` - Convert Meters to Feet\n"
	desc += "`fetome` - Convert Feet to Meters\n"
	desc += "`cmtoin` - Convert Centimeters to Inches\n"
	desc += "`intocm` - Convert Inches to Centimeters\n"

	if f.Auth.CheckAdminRole(m.Member) {
		desc += "\n"
		desc += "**Admin Commands**\n"
		desc += "`ping` - Ping the bot\n"
		desc += "`status` - Set the bot status\n"
		desc += "`purge` - Purge channel messages\n"
	}

	desc += "\n\nAll commands are prefixed with `" + f.Config.CommandKey + "`\n"

	msg := f.CreateDefinedEmbed("Help", desc, "")

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleMarkdownInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	desc := "*Italic* text goes between `*single asterisks*`\n"
	desc += "**Bold** text goes between `**double asterisks**`\n"
	desc += "***Bold and Italic*** text goes between `***triple asterisks***`\n"
	desc += "__Underlined__ text goes between `__double underscore__`\n"
	desc += "~~Strikethrough~~ text goes between `~~double tilde~~`\n"
	desc += "||Spoilers|| go between `|| double pipe ||`\n\n"
	desc += "You can combine the above styles.\n\n"
	desc += "Inline Code Blocks start and end with a single ``​`​``\n"
	desc += "Multi line Code Blocks start and end with ``​```​``\n"
	desc += "Multi line Code Blocks can also specify a language with ``​```​language`` at the start\n\n"
	desc += "Single line quotes start with `>`\n"
	desc += "Multi line quotes start with `>>>`\n"

	msg := f.CreateDefinedEmbed("Discord Markdown", desc, "")
	r, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	time.Sleep(15 * time.Second)

	err = s.ChannelMessageDelete(m.ChannelID, r.ID)
	if err != nil {
		return err
	}
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleCtoF(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a temperature")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	cels := (inF * 9.0 / 5.0) + 32.0
	celsF := float64(cels)

	msg := fmt.Sprintf("`%.1f°c` is `%.1f°f`", inF, celsF)

	e := f.CreateDefinedEmbed("Celsius to Farenheit", msg, "")
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleFtoC(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a temperature")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	faren := (inF - 32) * 5 / 9
	farenF := float64(faren)

	msg := fmt.Sprintf("`%.1f°f` is `%.1f°c`", inF, farenF)

	e := f.CreateDefinedEmbed("Farenheit to Celsius", msg, "")
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleMetersToFeet(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a distance")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	meters := inF * 3.28
	metersF := float64(meters)

	msg := fmt.Sprintf("`%.1fm` is `%.1fft`", inF, metersF)

	e := f.CreateDefinedEmbed("Meters to Feet", msg, "")
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleFeetToMeters(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a distance")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	feet := inF / 3.28
	feetF := float64(feet)

	msg := fmt.Sprintf("`%.1fft` is `%.1fm`", inF, feetF)

	e := f.CreateDefinedEmbed("Feet to Meters", msg, "")
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleCentimeterToInch(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a distance")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	inch := inF / 2.54
	inchF := float64(inch)

	msg := fmt.Sprintf("`%.1fcm` is `%.1fin`", inF, inchF)

	e := f.CreateDefinedEmbed("Centimeter To Inch", msg, "")
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleInchToCentimeter(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a distance")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	cm := inF * 2.54
	cmF := float64(cm)

	msg := fmt.Sprintf("`%.1fin` is `%.1fcm`", inF, cmF)

	e := f.CreateDefinedEmbed("Inch to Centimeter", msg, "")
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}
