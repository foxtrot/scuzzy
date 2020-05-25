package features

import (
	"errors"
	"github.com/bwmarrin/discord.go"
	"log"
	"strconv"
	"strings"
)

type commandHandler func(session *discordgo.Session, m *discordgo.MessageCreate) error

var commandHandlers = make(map[string]commandHandler)

func (f *Features) RegisterCommand(name string, handlerFunc commandHandler) {
	log.Printf("[*] Registering Command '%s'\n", name)
	commandHandlers[name] = handlerFunc
}

func (f *Features) RegisterHandlers() {
	// Misc Commands
	f.RegisterCommand("help", f.handleHelp)
	f.RegisterCommand("info", f.handleInfo)
	f.RegisterCommand("md", f.handleMarkdownInfo)
	f.RegisterCommand("userinfo", f.handleUserInfo)
	f.RegisterCommand("serverinfo", f.handleServerInfo)
	f.RegisterCommand("no", f.handleCat)

	// User Settings
	f.RegisterCommand("colors", f.handleUserColors)
	f.RegisterCommand("color", f.handleUserColor)

	// Conversion Helpers
	f.RegisterCommand("ctof", f.handleCtoF)
	f.RegisterCommand("ftoc", f.handleFtoC)
	f.RegisterCommand("metofe", f.handleMetersToFeet)
	f.RegisterCommand("fetome", f.handleFeetToMeters)
	f.RegisterCommand("cmtoin", f.handleCentimeterToInch)
	f.RegisterCommand("intocm", f.handleInchToCentimeter)

	// Admin Commands
	f.RegisterCommand("ping", f.handlePing)
	f.RegisterCommand("rules", f.handleRules)
	f.RegisterCommand("status", f.handleSetStatus)
	f.RegisterCommand("purge", f.handlePurgeChannel)
	f.RegisterCommand("kick", f.handleKickUser)
	f.RegisterCommand("ban", f.handleBanUser)
}

func (f *Features) ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	cKey := f.Config.CommandKey
	cCmd := strings.Split(m.Content, " ")[0]

	// Ignore the bot itself
	if m.Author.ID == s.State.User.ID {
		return nil
	}

	// Ignore anything not starting with the command prefix
	if !strings.HasPrefix(cCmd, cKey) {
		return nil
	}

	// Ignore Direct Messages
	if m.Member == nil {
		return nil
	}

	cName := strings.Split(cCmd, cKey)[1]

	if cmdFunc, ok := commandHandlers[cName]; ok {
		log.Printf("[*] Running command %s (Requested by %s)\n", cName, m.Author.Username)

		err := cmdFunc(s, m)
		if err != nil {
			log.Printf("[!] Command %s (Requested by %s) had error: '%s'\n", cName, m.Author.Username, err.Error())

			eMsg := f.CreateDefinedEmbed("Error ("+cName+")", err.Error(), "error", m.Author)
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *Features) ProcessMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) error {
	msgChannelID := m.ChannelID

	if m.BeforeDelete == nil {
		return errors.New("Couldn't get deleted message data.")
	}

	msgContent := m.BeforeDelete.Content
	msgAuthor := m.BeforeDelete.Author

	msg := "`Username` - " + msgAuthor.Username + "#" + msgAuthor.Discriminator + "\n"
	msg += "`User ID` - " + msgAuthor.ID + "\n"
	msg += "`Channel` - <#" + msgChannelID + ">\n"
	msg += "`Message` - " + msgContent + "\n"

	embed := f.CreateDefinedEmbed("Deleted Message", msg, "", msgAuthor)
	_, err := s.ChannelMessageSendEmbed(f.Config.LoggingChannel, embed)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) ProcessMessageDeleteBulk(s *discordgo.Session, m *discordgo.MessageDeleteBulk) error {
	msgChannelID := m.ChannelID

	msg := "`Channel` - <#" + msgChannelID + ">\n"
	msg += "Message IDs: \n"
	for k, v := range m.Messages {
		msg += strconv.Itoa(k) + ": `" + v + "`\n"
	}

	embed := f.CreateDefinedEmbed("Deleted Bulk Messages", msg, "", nil)
	_, err := s.ChannelMessageSendEmbed(f.Config.LoggingChannel, embed)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) ProcessUserJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) error {
	userChannel, err := s.UserChannelCreate(m.User.ID)
	if err != nil {
		log.Print("[!] Error (User Join): " + err.Error())
		return err
	}

	_, err = s.ChannelMessageSend(userChannel.ID, f.Config.WelcomeText)
	if err != nil {
		log.Print("[!] Error (User Join): " + err.Error())
		return err
	}

	for _, roleID := range f.Config.JoinRoleIDs {
		err = s.GuildMemberRoleAdd(f.Auth.Guild.ID, m.User.ID, roleID)
		if err != nil {
			log.Print("[!] Error (User Join)" + err.Error())
			return err
		}
	}

	return nil
}

func (f *Features) ProcessMessage(s *discordgo.Session, m interface{}) {
	switch m.(type) {
	case *discordgo.MessageCreate:
		// Pass Messages to the command processor
		err := f.ProcessCommand(s, m.(*discordgo.MessageCreate))
		if err != nil {
			log.Println("[!] Error: " + err.Error())
		}
		break
	case *discordgo.MessageDelete:
		// Log deleted messages to the logging channel.
		err := f.ProcessMessageDelete(s, m.(*discordgo.MessageDelete))
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Message Deleted)", err.Error(), "error", nil)
			_, err = s.ChannelMessageSendEmbed(f.Config.LoggingChannel, eMsg)
			if err != nil {
				log.Println("[!] Error " + err.Error())
			}
		}
		break
	case *discordgo.MessageDeleteBulk:
		err := f.ProcessMessageDeleteBulk(s, m.(*discordgo.MessageDeleteBulk))
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Message Bulk Deleted)", err.Error(), "error", nil)
			_, err = s.ChannelMessageSendEmbed(f.Config.LoggingChannel, eMsg)
			if err != nil {
				log.Println("[!] Error " + err.Error())
			}
		}
		break
	case *discordgo.GuildMemberAdd:
		// Handle new member (Welcome message, etc)
		err := f.ProcessUserJoin(s, m.(*discordgo.GuildMemberAdd))
		if err != nil {
			log.Println("[!] Error: " + err.Error())
		}
		break
	}

}
