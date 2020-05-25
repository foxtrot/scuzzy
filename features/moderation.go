package features

import (
	"errors"
	"github.com/bwmarrin/discord.go"
	"strconv"
	"strings"
	"time"
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
	msgS, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)

	time.Sleep(time.Second * 10)

	err = s.ChannelMessageDelete(m.ChannelID, msgS.ID)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleKickUser(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Auth.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use this command.")
	}

	var (
		mHandle    *discordgo.Member
		kickReason string
		err        error
	)

	args := strings.Split(m.Content, " ")
	if len(args) < 2 {
		return errors.New("You must specify a user to kick.")
	}
	if len(args) == 3 {
		kickReason = args[2]
	}

	member := args[1]
	idStr := strings.Replace(member, "<@!", "", 1)
	idStr = strings.Replace(idStr, ">", "", 1)
	mHandle, err = s.GuildMember(f.Config.Guild.ID, idStr)
	if err != nil {
		return err
	}

	err = s.GuildMemberDeleteWithReason(f.Config.Guild.ID, mHandle.User.ID, kickReason)
	if err != nil {
		return err
	}

	msg := "User `" + mHandle.User.Username + "#" + mHandle.User.Discriminator + "` was kicked.\n"
	if len(kickReason) > 0 {
		msg += "Reason: `" + kickReason + "`\n"
	}

	embed := f.CreateDefinedEmbed("Kick User", msg, "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleBanUser(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Auth.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use this command.")
	}

	var (
		mHandle   *discordgo.User
		banReason string
		err       error
	)

	args := strings.Split(m.Content, " ")
	if len(args) < 2 {
		return errors.New("You must specify a user to ban.")
	}
	if len(args) == 3 {
		banReason = args[2]
	}

	member := args[1]
	idStr := strings.Replace(member, "<@!", "", 1)
	idStr = strings.Replace(idStr, ">", "", 1)
	mHandle, err = s.User(idStr)
	if err != nil {
		return err
	}

	err = s.GuildBanCreateWithReason(f.Config.Guild.ID, mHandle.ID, banReason, 0)
	if err != nil {
		return err
	}

	msg := "User `" + mHandle.Username + "#" + mHandle.Discriminator + "` was banned.\n"
	if len(banReason) > 0 {
		msg += "Reason: `" + banReason + "`\n"
	}

	embed := f.CreateDefinedEmbed("Ban User", msg, "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return err
	}

	return nil
}
