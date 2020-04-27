package features

import (
	"errors"
	"github.com/bwmarrin/discord.go"
	"strconv"
	"strings"
)

func (f *Features) handlePurgeChannel(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Auth.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	}

	purgeSplit := strings.SplitN(m.Content, " ", 2)
	if len(purgeSplit) < 2 {
		return errors.New("No message count supplied")
	}

	msgCount, err := strconv.Atoi(purgeSplit[1])
	if err != nil {
		return nil
	}

	if msgCount > 100 {
		return errors.New("You may only purge upto 100 messages at a time.")
	}

	chanMsgs, err := s.ChannelMessages(m.ChannelID, msgCount, "", "", "")
	if err != nil {
		return err
	}

	msg := f.CreateDefinedEmbed("Purge Channel", "Purging `"+purgeSplit[1]+"` messages.", "", m.Author)
	r, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	var delMsgs []string
	for _, v := range chanMsgs {
		delMsgs = append(delMsgs, v.ID)
	}

	err = s.ChannelMessagesBulkDelete(m.ChannelID, delMsgs)
	if err != nil {
		return err
	}

	err = s.ChannelMessageDelete(m.ChannelID, r.ID)
	msg = f.CreateDefinedEmbed("Purge Channel", "Purged `"+purgeSplit[1]+"` messages!", "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}
