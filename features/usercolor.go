package features

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func (f *Features) handleUserColors(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckCommandRestrictions(m) {
		return errors.New("This command is not allowed in this channel.")
	}

	msgC := "You can choose from the following colors:\n\n"
	for _, v := range f.Config.ColorRoles {
		msgC += "<@&" + v.ID + ">\n"
	}
	msgC += "\n\nUse `" + f.Config.CommandKey + "color <color>` to set.\n"
	msgC += "Example: `" + f.Config.CommandKey + "color red`.\n"

	msg := f.CreateDefinedEmbed("User Colors", msgC, "", m.Author)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleUserColor(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var err error

	if !f.Permissions.CheckCommandRestrictions(m) {
		return errors.New("This command is not allowed in this channel.")
	}

	rUserID := m.Author.ID

	userInput := strings.Split(m.Content, " ")
	if len(userInput) < 2 {
		err = f.handleUserColors(s, m)
		return err
	}
	roleColorName := userInput[1]
	roleColorName = strings.ToLower(roleColorName)

	roleColorID := ""
	for _, role := range f.Config.ColorRoles {
		if role.Name == roleColorName {
			roleColorID = role.ID
			break
		}
	}
	if len(roleColorID) == 0 {
		err = f.handleUserColors(s, m)
		return err
	}

	for _, role := range f.Config.ColorRoles {
		// Attempt to remove all color roles regardless of if they have them or not.
		// Slow because of the REST requests...
		_ = s.GuildMemberRoleRemove(m.GuildID, rUserID, role.ID)
	}

	err = s.GuildMemberRoleAdd(m.GuildID, rUserID, roleColorID)
	if err != nil {
		return err
	} else {
		msg := f.CreateDefinedEmbed("User Color", "<@"+m.Author.ID+">: Your color has been set to <@&"+roleColorID+">!", "success", m.Author)
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
		if err != nil {
			return err
		}
	}

	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}

	return nil
}
