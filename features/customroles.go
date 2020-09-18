package features

import (
	discordgo "github.com/bwmarrin/discord.go"
	"strings"
)

func (f *Features) handleListCustomRoles(s *discordgo.Session, m *discordgo.MessageCreate) error {
	msgC := "You can choose from the following roles:\n\n"
	for _, v := range f.Config.CustomRoles {
		msgC += "<@&" + v.ID + "> (" + v.ShortName + ")\n"
	}
	msgC += "\n\n Use `" + f.Config.CommandKey + "joinrole <role_name>` to join a role.\n"
	msgC += "Example: `" + f.Config.CommandKey + "joinrole pineapple`.\n"

	msg := f.CreateDefinedEmbed("Joinable Roles", msgC, "", m.Author)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleJoinCustomRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var err error

	rUserID := m.Author.ID

	userInput := strings.Split(m.Content, " ")
	if len(userInput) < 2 {
		err = f.handleListCustomRoles(s, m)
		return err
	}

	desiredRole := userInput[1]
	desiredRole = strings.ToLower(desiredRole)
	desiredRoleID := ""

	for _, role := range f.Config.CustomRoles {
		if role.ShortName == desiredRole {
			desiredRoleID = role.ID
			break
		}
	}

	if len(desiredRoleID) == 0 {
		err = f.handleListCustomRoles(s, m)
		return err
	}

	err = s.GuildMemberRoleAdd(m.GuildID, rUserID, desiredRoleID)
	if err != nil {
		return err
	} else {
		msg := f.CreateDefinedEmbed("Join Role", "<@"+m.Author.ID+">: You have joined <@&"+desiredRoleID+">!", "success", m.Author)
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

func (f *Features) handleLeaveCustomRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var err error

	rUserID := m.Author.ID

	userInput := strings.Split(m.Content, " ")
	if len(userInput) < 2 {
		err = f.handleListCustomRoles(s, m)
		return err
	}

	desiredRole := userInput[1]
	desiredRole = strings.ToLower(desiredRole)
	desiredRoleID := ""

	for _, role := range f.Config.CustomRoles {
		if role.ShortName == desiredRole {
			desiredRoleID = role.ID
			break
		}
	}

	if len(desiredRoleID) == 0 {
		err = f.handleListCustomRoles(s, m)
		return err
	}

	err = s.GuildMemberRoleRemove(m.GuildID, rUserID, desiredRoleID)
	if err != nil {
		return err
	} else {
		msg := f.CreateDefinedEmbed("Leave Role", "<@"+m.Author.ID+">: You have left <@&"+desiredRoleID+">!", "success", m.Author)
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
