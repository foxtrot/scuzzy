package actions

import "github.com/bwmarrin/discordgo"

func KickUser(s *discordgo.Session, guild string, user string, reason string) error {
	err := s.GuildMemberDeleteWithReason(guild, user, reason)
	if err != nil {
		return err
	}

	return nil
}

func BanUser(s *discordgo.Session, guild string, user string, reason string) error {
	err := s.GuildBanCreateWithReason(guild, user, reason, 0)
	if err != nil {
		return err
	}

	return nil
}
