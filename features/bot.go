package features

import (
	"errors"
	discordgo "github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strings"
	"time"
)

func (f *Features) handleSetStatus(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	}

	stSplit := strings.SplitN(m.Content, " ", 2)
	if len(stSplit) < 2 {
		return errors.New("You did not specify a status.")
	}

	st := stSplit[1]

	err := s.UpdateGameStatus(0, st)
	if err != nil {
		return err
	}

	msg := f.CreateDefinedEmbed("Set Status", "Operation completed successfully.", "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleDisconnect(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("I'm sorry Dave, I'm afraid I can't do that.")
	}

	msg := f.CreateDefinedEmbed("Disconnect", "Attempting Disconnect...", "", m.Author)
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

func (f *Features) handleReconnect(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	}

	t := time.Now()

	err := s.Close()
	if err != nil {
		return err
	}

	err = s.Open()
	if err != nil {
		log.Fatal(err)
	}

	msg := f.CreateDefinedEmbed("Reconnect", "Reconnected Successfully.\nTime: `"+time.Since(t).String()+"`.", "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}
