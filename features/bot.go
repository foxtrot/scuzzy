package features

import (
	"errors"
	discordgo "github.com/bwmarrin/discord.go"
	"strings"
)

func (f *Features) handleSetStatus(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Auth.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	}

	stSplit := strings.SplitN(m.Content, " ", 2)
	if len(stSplit) < 2 {
		return errors.New("You did not specify a status.")
	}

	st := stSplit[1]

	err := s.UpdateStatus(0, st)
	if err != nil {
		return err
	}

	msg := f.CreateDefinedEmbed("Set Status", "Operation completed successfully.", "success")
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}
