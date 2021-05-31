package commands

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func (c *Commands) handleSetOverwatchEnforcement(s *discordgo.Session, m *discordgo.MessageCreate) error {
	args := strings.Split(m.Content, " ")

	if len(args) != 2 {
		return errors.New("Invalid arguments supplied. Usage: " + c.Config.CommandKey + "enforcement [on/off]")
	}

	mode := strings.ToLower(args[1])

	if mode == "on" {
		c.Config.EnforceMode = true
	} else {
		c.Config.EnforceMode = false
	}

	err := c.handleSaveConfig(s, m)

	return err
}

func (c *Commands) handleSetWordFilter(s *discordgo.Session, m *discordgo.MessageCreate) error {
	args := strings.Split(m.Content, " ")

	if len(args) != 2 {
		return errors.New("Invalid arguments supplied. Usage: " + c.Config.CommandKey + "wordfilter [on/off]")
	}

	mode := strings.ToLower(args[1])

	if mode == "on" {
		c.Config.FilterLanguage = true
	} else {
		c.Config.FilterLanguage = false
	}

	err := c.handleSaveConfig(s, m)

	return err
}
