package commands

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strings"
	"time"
)

func (c *Commands) handleSetStatus(s *discordgo.Session, m *discordgo.MessageCreate) error {
	stSplit := strings.SplitN(m.Content, " ", 2)
	if len(stSplit) < 2 {
		return errors.New("You did not specify a status.")
	}

	st := stSplit[1]

	err := s.UpdateGameStatus(0, st)
	if err != nil {
		return err
	}

	msg := c.CreateDefinedEmbed("Set Status", "Operation completed successfully.", "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) handleDisconnect(s *discordgo.Session, m *discordgo.MessageCreate) error {
	msg := c.CreateDefinedEmbed("Disconnect", "Attempting Disconnect...", "", m.Author)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	err = s.Close()
	if err != nil {
		return err
	}

	os.Exit(0)

	return nil
}

func (c *Commands) handleReconnect(s *discordgo.Session, m *discordgo.MessageCreate) error {
	t := time.Now()

	err := s.Close()
	if err != nil {
		return err
	}

	err = s.Open()
	if err != nil {
		log.Fatal(err)
	}

	msg := c.CreateDefinedEmbed("Reconnect", "Reconnected Successfully.\nTime: `"+time.Since(t).String()+"`.", "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}
