package features

import (
	"github.com/bwmarrin/discord.go"
	"log"
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
	f.RegisterCommand("status", f.handleSetStatus)
	f.RegisterCommand("purge", f.handlePurgeChannel)
}

func (f *Features) ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	cKey := f.Config.CommandKey
	cCmd := strings.Split(m.Content, " ")[0]

	// Ignore the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore anything not starting with the command prefix
	if !strings.HasPrefix(cCmd, cKey) {
		return
	}

	// Ignore Direct Messages
	if m.Member == nil {
		return
	}

	cName := strings.Split(cCmd, cKey)[1]

	if cmdFunc, ok := commandHandlers[cName]; ok {
		log.Printf("[*] Running command %s (Requested by %s)\n", cName, m.Author.Username)

		err := cmdFunc(s, m)
		if err != nil {
			log.Printf("[*] Command %s (Requested by %s) had error: '%s'\n", cName, m.Author.Username, err.Error())

			eMsg := f.CreateDefinedEmbed("Error ("+cName+")", err.Error(), "error", m.Author)
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (f *Features) ProcessUserJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	userChannel, err := s.UserChannelCreate(m.User.ID)
	if err != nil {
		log.Print("Error (User Join): " + err.Error())
		return
	}

	_, err = s.ChannelMessageSend(userChannel.ID, f.Config.WelcomeText)
	if err != nil {
		log.Print("Error (User Join): " + err.Error())
		return
	}

	for _, roleID := range f.Config.JoinRoleIDs {
		err = s.GuildMemberRoleAdd(f.Auth.Guild.ID, m.User.ID, roleID)
		if err != nil {
			log.Print("Error (User Join)" + err.Error())
			return
		}
	}
}
