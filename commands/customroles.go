package commands

import (
	"errors"
	discordgo "github.com/bwmarrin/discordgo"
	"github.com/foxtrot/scuzzy/models"
	"strings"
)

func (c *Commands) handleListCustomRoles(s *discordgo.Session, m *discordgo.MessageCreate) error {
	msgC := "You can choose from the following roles:\n\n"
	for _, v := range c.Config.CustomRoles {
		msgC += "<@&" + v.ID + "> (" + v.ShortName + ")\n"
	}
	msgC += "\n\n Use `" + c.Config.CommandKey + "joinrole <role_name>` to join a role.\n"
	msgC += "Example: `" + c.Config.CommandKey + "joinrole pineapple`.\n"

	msg := c.CreateDefinedEmbed("Joinable Roles", msgC, "", m.Author)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleJoinCustomRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var err error

	rUserID := m.Author.ID

	userInput := strings.Split(m.Content, " ")
	if len(userInput) < 2 {
		err = c.handleListCustomRoles(s, m)
		return err
	}

	desiredRole := userInput[1]
	desiredRole = strings.ToLower(desiredRole)
	desiredRoleID := ""

	for _, role := range c.Config.CustomRoles {
		if role.ShortName == desiredRole {
			desiredRoleID = role.ID
			break
		}
	}

	if len(desiredRoleID) == 0 {
		err = c.handleListCustomRoles(s, m)
		return err
	}

	err = s.GuildMemberRoleAdd(m.GuildID, rUserID, desiredRoleID)
	if err != nil {
		return err
	} else {
		msg := c.CreateDefinedEmbed("Join Role", "<@"+m.Author.ID+">: You have joined <@&"+desiredRoleID+">!", "success", m.Author)
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

func (c *Commands) handleLeaveCustomRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var err error

	rUserID := m.Author.ID

	userInput := strings.Split(m.Content, " ")
	if len(userInput) < 2 {
		err = c.handleListCustomRoles(s, m)
		return err
	}

	desiredRole := userInput[1]
	desiredRole = strings.ToLower(desiredRole)
	desiredRoleID := ""

	for _, role := range c.Config.CustomRoles {
		if role.ShortName == desiredRole {
			desiredRoleID = role.ID
			break
		}
	}

	if len(desiredRoleID) == 0 {
		err = c.handleListCustomRoles(s, m)
		return err
	}

	err = s.GuildMemberRoleRemove(m.GuildID, rUserID, desiredRoleID)
	if err != nil {
		return err
	} else {
		msg := c.CreateDefinedEmbed("Leave Role", "<@"+m.Author.ID+">: You have left <@&"+desiredRoleID+">!", "success", m.Author)
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

func (c *Commands) handleAddCustomRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var err error

	userInput := strings.Split(m.Content, " ")
	if len(userInput) < 3 {
		return errors.New("Expected Arguments: short_name role_id")
	}

	shortName := userInput[1]
	shortName = strings.ToLower(shortName)
	roleID := userInput[2]

	customRole := models.CustomRole{
		Name:      "",
		ShortName: shortName,
		ID:        roleID,
	}

	c.Config.CustomRoles = append(c.Config.CustomRoles, customRole)

	err = c.handleSaveConfig(s, m)
	if err != nil {
		return err
	}

	return nil
}
