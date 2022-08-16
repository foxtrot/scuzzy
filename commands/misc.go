package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/foxtrot/scuzzy/models"
)

func (c *Commands) handleSetConfig(s *discordgo.Session, m *discordgo.MessageCreate) error {
	configArgs := strings.Split(m.Content, " ")

	if len(configArgs) != 3 {
		return errors.New("Invalid arguments supplied. Usage: " + c.Config.CommandKey + "setconfig <key> <value>")
	}

	configKey := configArgs[1]
	configVal := configArgs[2]

	rt := reflect.TypeOf(c.Config)
	for i := 0; i < rt.NumField(); i++ {
		x := rt.Field(i)
		tagVal := strings.Split(x.Tag.Get("json"), ",")[0]
		tagName := x.Name

		if tagVal == configKey {
			prop := reflect.ValueOf(&c.Config).Elem().FieldByName(tagName)

			switch prop.Interface().(type) {
			case string:
				prop.SetString(configVal)
				break
			case int:
				intVal, err := strconv.ParseInt(configVal, 10, 64)
				if err != nil {
					return err
				}
				prop.SetInt(intVal)
				break
			case float64:
				floatVal, err := strconv.ParseFloat(configVal, 64)
				if err != nil {
					return err
				}
				prop.SetFloat(floatVal)
				break
			case bool:
				boolVal, err := strconv.ParseBool(configVal)
				if err != nil {
					return err
				}
				prop.SetBool(boolVal)
				break
			default:
				return errors.New("Unsupported key value type")
			}

			msgE := c.CreateDefinedEmbed("Set Configuration", "Successfully set property '"+configKey+"'!", "success", m.Author)
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, msgE)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("Unknown key specified")
}

func (c *Commands) handleGetConfig(s *discordgo.Session, m *discordgo.MessageCreate) error {
	//TODO: Handle printing of slices (check the Type, loop accordingly)

	configArgs := strings.Split(m.Content, " ")
	configKey := "all"
	if len(configArgs) == 2 {
		configKey = configArgs[1]
	}

	msg := ""

	rt := reflect.TypeOf(*c.Config)
	for i := 0; i < rt.NumField(); i++ {
		x := rt.Field(i)
		tagVal := strings.Split(x.Tag.Get("json"), ",")[0]
		tagName := x.Name
		prop := reflect.ValueOf(c.Config).Elem().FieldByName(tagName)

		if configKey == "all" {
			switch prop.Interface().(type) {
			case string:
				if len(prop.String()) > 256 {
					// Truncate large values.
					msg += "`" + tagName + "` - " + "Truncated...\n"
				} else {
					msg += "`" + tagName + "` - `" + prop.String() + "`\n"
				}
				break
			default:
				// Ignore non strings for now...
				msg += "`" + tagName + "` - Skipped Value\n"
				continue
			}
		} else {
			if tagVal == configKey {
				switch prop.Interface().(type) {
				case string:
					msg += "`" + tagName + "` - `" + prop.String() + "`\n"
				default:
					// Ignore non strings for now...
					msg += "`" + tagName + "` - Skipped Value\n"
				}

				eMsg := c.CreateDefinedEmbed("Get Configuration", msg, "success", m.Author)
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	if msg == "" {
		return errors.New("Unknown key specified")
	}

	eMsg := c.CreateDefinedEmbed("Get Configuration", msg, "success", m.Author)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleReloadConfig(s *discordgo.Session, m *discordgo.MessageCreate) error {
	fBuf, err := ioutil.ReadFile(c.Config.ConfigPath)
	if err != nil {
		return err
	}

	conf := &models.Configuration{}

	err = json.Unmarshal(fBuf, &conf)
	if err != nil {
		return err
	}

	c.Config = conf
	c.Permissions.Config = conf

	eMsg := c.CreateDefinedEmbed("Reload Configuration", "Successfully reloaded configuration from disk", "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleSaveConfig(s *discordgo.Session, m *discordgo.MessageCreate) error {
	j, err := json.Marshal(c.Config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.Config.ConfigPath, j, os.ModePerm)
	if err != nil {
		return err
	}

	eMsg := c.CreateDefinedEmbed("Save Configuration", "Saved runtime configuration successfully", "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleCat(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://giphy.com/gifs/cat-cute-no-rCxogJBzaeZuU")
	if err != nil {
		return err
	}

	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handlePing(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var r *discordgo.Message
	var err error

	msg := c.CreateDefinedEmbed("Ping", "Pong", "success", m.Author)
	r, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
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

func (c *Commands) handleInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	desc := "**Source**:   https://github.com/foxtrot/scuzzy\n"
	desc += "**Language**: Go\n"
	desc += "**Commands**: See `" + c.Config.CommandKey + "help`\n\n\n"

	gm, err := s.GuildMember(c.Config.GuildID, s.State.User.ID)
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
		FooterImageURL: "https://cdn.discordapp.com/avatars/514163441548656641/a_ac5e022e77e62e7793711ebde8cdf4a1.gif",
		ThumbnailURL:   gm.User.AvatarURL(""),
		ThumbnailH:     150,
		ThumbnailW:     150,
		ProviderURL:    "",
		ProviderText:   "",
		AuthorText:     "",
		AuthorURL:      "",
		AuthorImageURL: "",
	}

	msg := c.CreateCustomEmbed(&d)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) error {
	keys := make([]int, 0, len(c.ScuzzyCommands))
	for _, cmd := range c.ScuzzyCommands {
		keys = append(keys, cmd.Index)
	}
	sort.Ints(keys)

	for _, k := range keys {
		fmt.Println(k, c.ScuzzyCommandsByIndex[k])
	}

	desc := "**Available Commands**\n"
	for _, k := range keys {
		command := c.ScuzzyCommandsByIndex[k]

		if !command.AdminOnly && command.Description != "" {
			desc += "`" + command.Name + "` - " + command.Description + "\n"
		}
	}

	if c.Permissions.CheckAdminRole(m.Member) {
		desc += "\n"
		desc += "**Admin Commands**\n"
		for _, k := range keys {
			command := c.ScuzzyCommandsByIndex[k]

			if command.AdminOnly {
				desc += "`" + command.Name + "` - " + command.Description + "\n"
			}
		}
	}

	desc += "\n\nAll commands are prefixed with `" + c.Config.CommandKey + "`\n"

	msg := c.CreateDefinedEmbed("Help", desc, "", m.Author)

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

func (c *Commands) handleRules(s *discordgo.Session, m *discordgo.MessageCreate) error {
	msg := c.Config.RulesText
	embedTitle := "Rules (" + c.Config.GuildName + ")"
	embed := c.CreateDefinedEmbed(embedTitle, msg, "success", m.Author)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleMarkdownInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	cleanup := true
	args := strings.Split(m.Content, " ")

	if len(args) == 2 {
		if args[1] == "stay" && c.Permissions.CheckAdminRole(m.Member) {
			cleanup = false
		}
	}

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

	msg := c.CreateDefinedEmbed("Discord Markdown", desc, "", m.Author)
	r, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	if cleanup {
		time.Sleep(15 * time.Second)

		err = s.ChannelMessageDelete(m.ChannelID, r.ID)
		if err != nil {
			return err
		}
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Commands) handleCtoF(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

	e := c.CreateDefinedEmbed("Celsius to Farenheit", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleFtoC(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

	e := c.CreateDefinedEmbed("Farenheit to Celsius", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleMetersToFeet(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

	e := c.CreateDefinedEmbed("Meters to Feet", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleFeetToMeters(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

	e := c.CreateDefinedEmbed("Feet to Meters", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleCentimeterToInch(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

	e := c.CreateDefinedEmbed("Centimeter To Inch", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleInchToCentimeter(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

	e := c.CreateDefinedEmbed("Inch to Centimeter", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleUserInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var (
		mHandle   *discordgo.Member
		requester *discordgo.Member
		err       error
	)

	userSplit := strings.Split(m.Content, " ")

	if len(userSplit) < 2 {
		mHandle, err = s.GuildMember(c.Config.GuildID, m.Author.ID)
		requester = mHandle
		if err != nil {
			return err
		}
	} else {
		idStr := strings.ReplaceAll(userSplit[1], "<@!", "")
		idStr = strings.ReplaceAll(idStr, "<@", "")
		idStr = strings.ReplaceAll(idStr, ">", "")
		mHandle, err = s.GuildMember(c.Config.GuildID, idStr)
		if err != nil {
			return err
		}
		requester, err = s.GuildMember(c.Config.GuildID, m.Author.ID)
		if err != nil {
			return err
		}
	}

	rUserID := mHandle.User.ID
	rUserNick := mHandle.Nick
	rUsername := mHandle.User.Username
	rUserDiscrim := mHandle.User.Discriminator
	rUserAvatar := mHandle.User.AvatarURL("4096")
	rJoinTime := mHandle.JoinedAt
	rRoles := mHandle.Roles

	if len(rUserNick) == 0 {
		rUserNick = "No Nickname"
	}

	rJoinTimeP, err := rJoinTime.Parse()
	if err != nil {
		return err
	}

	rRolesTidy := ""
	if len(rRoles) == 0 {
		rRolesTidy = "No Roles"
	} else {
		for _, role := range rRoles {
			rRolesTidy += "<@&" + role + "> "
		}
	}

	msg := "**User ID**: `" + rUserID + "`\n"
	msg += "**User Name**: `" + rUsername + "`\n"
	msg += "**User Nick**: `" + rUserNick + "`\n"
	msg += "**User Discrim**: `#" + rUserDiscrim + "`\n"
	msg += "**User Join**:  `" + rJoinTimeP.String() + "`\n"
	msg += "**User Roles**: " + rRolesTidy + "\n"

	embedData := models.CustomEmbed{
		URL:            "",
		Title:          "User Info (" + rUsername + ")",
		Desc:           msg,
		Type:           "",
		Timestamp:      time.Now().Format(time.RFC3339),
		Color:          0xFFA500,
		FooterText:     "Requested by " + requester.User.Username + "#" + requester.User.Discriminator,
		FooterImageURL: "",
		ImageURL:       "",
		ImageH:         0,
		ImageW:         0,
		ThumbnailURL:   rUserAvatar,
		ThumbnailH:     512,
		ThumbnailW:     512,
		ProviderURL:    "",
		ProviderText:   "",
		AuthorText:     "",
		AuthorURL:      "",
		AuthorImageURL: "",
	}

	embed := c.CreateCustomEmbed(&embedData)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleServerInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	g, err := s.Guild(c.Config.GuildID)
	if err != nil {
		return err
	}

	sID := c.Config.GuildID
	sName := c.Config.GuildName

	chans, _ := s.GuildChannels(c.Config.GuildID)
	sChannels := strconv.Itoa(len(chans))
	sEmojis := strconv.Itoa(len(g.Emojis))
	sRoles := strconv.Itoa(len(g.Roles))
	sRegion := g.Region

	iID, _ := strconv.Atoi(c.Config.GuildID)
	createdMSecs := ((iID / 4194304) + 1420070400000) / 1000
	sCreatedAt := time.Unix(int64(createdMSecs), 0).Format(time.RFC1123)

	sIconURL := g.IconURL()

	user := m.Author

	desc := "**Server ID**: `" + sID + "`\n"
	desc += "**Server Name**: `" + sName + "`\n"
	desc += "**Server Channels**: `" + sChannels + "`\n"
	desc += "**Server Emojis**: `" + sEmojis + "`\n"
	desc += "**Server Roles**: `" + sRoles + "`\n"
	desc += "**Server Region**: `" + sRegion + "`\n"
	desc += "**Server Creation**: `" + sCreatedAt + "`\n"

	embedData := models.CustomEmbed{
		URL:            "",
		Title:          "Server Info (" + sName + ")",
		Desc:           desc,
		Type:           "",
		Timestamp:      time.Now().Format(time.RFC3339),
		Color:          0xFFA500,
		FooterText:     "Requested by " + user.Username + "#" + user.Discriminator,
		FooterImageURL: "",
		ImageURL:       "",
		ImageH:         0,
		ImageW:         0,
		ThumbnailURL:   sIconURL,
		ThumbnailH:     256,
		ThumbnailW:     256,
		ProviderURL:    "",
		ProviderText:   "",
		AuthorText:     "",
		AuthorURL:      "",
		AuthorImageURL: "",
	}

	msg := c.CreateCustomEmbed(&embedData)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}
func (c *Commands) handleGoogle4U(s *discordgo.Session, m *discordgo.MessageCreate) error {
	args := strings.Split(m.Content, " ")

	if len(args) < 2 {
		return errors.New("You did not specify anything to google")
	}

	input := m.Content[strings.Index(m.Content, " "):len(m.Content)]

	desc := "https://letmegooglethat.com/?q=" + url.QueryEscape(input)

	msg := c.CreateDefinedEmbed("Google", desc, "", m.Author)
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
