package features

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

func (f *Features) handleSetSlowmode(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	}

	slowmodeSplit := strings.Split(m.Content, " ")
	if len(slowmodeSplit) < 2 {
		return errors.New("You must supply at least an amount of time")
	}

	slowmodeTimeStr := slowmodeSplit[1]
	slowModeTime, err := strconv.Atoi(slowmodeTimeStr)
	if err != nil {
		return err
	}

	if len(slowmodeSplit) == 3 {
		if slowmodeSplit[2] == "all" {
			channels, err := s.GuildChannels(f.Config.GuildID)
			if err != nil {
				return err
			}

			for _, channel := range channels {
				s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
					RateLimitPerUser: slowModeTime,
				})
			}
		}
	} else {
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			return err
		}

		_, err = s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{
			RateLimitPerUser: slowModeTime,
		})
		if err != nil {
			return err
		}
	}

	msg := f.CreateDefinedEmbed("Slow Mode", "Successfully set Slow Mode to `"+slowmodeTimeStr+"`.", "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleUnsetSlowmode(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	}

	slowmodeSplit := strings.Split(m.Content, " ")

	if len(slowmodeSplit) == 2 {
		if slowmodeSplit[1] == "all" {
			channels, err := s.GuildChannels(f.Config.GuildID)
			if err != nil {
				return err
			}

			for _, channel := range channels {
				s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
					RateLimitPerUser: 0,
				})
			}
		}
	} else {
		_, err := s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{
			RateLimitPerUser: 0,
		})
		if err != nil {
			return err
		}
	}

	msg := f.CreateDefinedEmbed("Slow Mode", "Successfully unset Slow Mode", "success", m.Author)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handlePurgeChannel(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
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
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 10)

	err = s.ChannelMessageDelete(m.ChannelID, msgS.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleKickUser(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
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
	idStr := strings.ReplaceAll(member, "<@!", "")
	idStr = strings.ReplaceAll(idStr, "<@", "")
	idStr = strings.ReplaceAll(idStr, ">", "")
	mHandle, err = s.GuildMember(f.Config.GuildID, idStr)
	if err != nil {
		return err
	}

	err = s.GuildMemberDeleteWithReason(f.Config.GuildID, mHandle.User.ID, kickReason)
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
	if !f.Permissions.CheckAdminRole(m.Member) {
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
	idStr := strings.ReplaceAll(member, "<@!", "")
	idStr = strings.ReplaceAll(idStr, "<@", "")
	idStr = strings.ReplaceAll(idStr, ">", "")
	mHandle, err = s.User(idStr)
	if err != nil {
		return err
	}

	err = s.GuildBanCreateWithReason(f.Config.GuildID, mHandle.ID, banReason, 0)
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

func (f *Features) handleIgnoreUser(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use this command.")
	}

	ignArgs := strings.Split(m.Content, " ")
	if len(ignArgs) < 2 {
		return errors.New("You did not specify a user.")
	}

	member := ignArgs[1]
	idStr := strings.ReplaceAll(member, "<@!", "")
	idStr = strings.ReplaceAll(idStr, "<@", "")
	idStr = strings.ReplaceAll(idStr, ">", "")

	f.Config.IgnoredUsers = append(f.Config.IgnoredUsers, idStr)

	eMsg := f.CreateDefinedEmbed("Ignore User", "<@!"+idStr+"> is now being ignored.", "success", m.Author)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
	if err != nil {
		return err
	}

	err = f.handleSaveConfig(s, m)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleUnIgnoreUser(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use this command.")
	}

	ignArgs := strings.Split(m.Content, " ")
	if len(ignArgs) < 2 {
		return errors.New("You did not specify a user.")
	}

	member := ignArgs[1]
	idStr := strings.ReplaceAll(member, "<@!", "")
	idStr = strings.ReplaceAll(idStr, "<@", "")
	idStr = strings.ReplaceAll(idStr, ">", "")

	for k, v := range f.Config.IgnoredUsers {
		if v == idStr {
			f.Config.IgnoredUsers[k] = f.Config.IgnoredUsers[len(f.Config.IgnoredUsers)-1]
			f.Config.IgnoredUsers = f.Config.IgnoredUsers[:len(f.Config.IgnoredUsers)-1]
		}
	}

	eMsg := f.CreateDefinedEmbed("Unignore User", "<@!"+idStr+"> is not being ignored.", "success", m.Author)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
	if err != nil {
		return err
	}

	err = f.handleSaveConfig(s, m)
	if err != nil {
		return err
	}

	return nil
}
