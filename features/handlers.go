package features

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (f *Features) RegisterCommand(name string, description string, adminonly bool, handler ScuzzyHandler) {
	log.Printf("[*] Registering Command '%s'\n", name)
	c := ScuzzyCommand{
		Name:        name,
		Description: description,
		AdminOnly:   adminonly,
		Handler:     handler,
	}
	f.ScuzzyCommands[name] = c
}

func (f *Features) RegisterHandlers() {
	// Misc Commands
	f.RegisterCommand("help", "Show Help Text", false, f.handleHelp)
	f.RegisterCommand("info", "Show Bot Info", false, f.handleInfo)
	f.RegisterCommand("md", "Show common Discord MarkDown formatting", false, f.handleMarkdownInfo)
	f.RegisterCommand("userinfo", "Display a users information", false, f.handleUserInfo)
	f.RegisterCommand("serverinfo", "Display the current servers information", false, f.handleServerInfo)
	f.RegisterCommand("no", "", false, f.handleCat)

	// User Settings
	f.RegisterCommand("colours", "Display available colour roles", false, f.handleUserColors)
	f.RegisterCommand("colors", "", false, f.handleUserColors)
	f.RegisterCommand("colour", "Set a colour role for yourself", false, f.handleUserColor)
	f.RegisterCommand("color", "", false, f.handleUserColor)
	f.RegisterCommand("listroles", "List user joinable roles", false, f.handleListCustomRoles)
	f.RegisterCommand("joinrole", "Join an available role for yourself", false, f.handleJoinCustomRole)
	f.RegisterCommand("leaverole", "Leave an available role", false, f.handleLeaveCustomRole)

	// Conversion Helpers
	f.RegisterCommand("ctof", "Convert Celsius to Farenheit", false, f.handleCtoF)
	f.RegisterCommand("ftoc", "Convert Farenheit to Celsius", false, f.handleFtoC)
	f.RegisterCommand("metofe", "Convert Meters to Feet", false, f.handleMetersToFeet)
	f.RegisterCommand("fetome", "Convert Feet to Meters", false, f.handleFeetToMeters)
	f.RegisterCommand("cmtoin", "Convert Centimeters to Inches", false, f.handleCentimeterToInch)
	f.RegisterCommand("intocm", "Convert Inches to Centimeters", false, f.handleInchToCentimeter)

	// Admin Commands
	f.RegisterCommand("ping", "Ping Scuzzy", true, f.handlePing)
	f.RegisterCommand("rules", "Display the Server Rules", true, f.handleRules)
	f.RegisterCommand("status", "Set Bot Status", true, f.handleSetStatus)
	f.RegisterCommand("purge", "Purge Channel Messages", true, f.handlePurgeChannel)
	f.RegisterCommand("kick", "Kick a User", true, f.handleKickUser)
	f.RegisterCommand("ban", "Ban a User", true, f.handleBanUser)
	f.RegisterCommand("slow", "Set Channel Slow Mode", true, f.handleSetSlowmode)
	f.RegisterCommand("unslow", "Unset Channel Slow Mode", true, f.handleUnsetSlowmode)
	f.RegisterCommand("ignore", "Add a user to Scuzzy's ignore list", true, f.handleIgnoreUser)
	f.RegisterCommand("unignore", "Remove a user from Scuzzy's ignore list", true, f.handleUnIgnoreUser)
	f.RegisterCommand("setconfig", "Set Configuration", true, f.handleSetConfig)
	f.RegisterCommand("getconfig", "Print Configuration", true, f.handleGetConfig)
	f.RegisterCommand("saveconfig", "Save Configuration to Disk", true, f.handleSaveConfig)
	f.RegisterCommand("reloadconfig", "Reload Configuration", true, f.handleReloadConfig)
	f.RegisterCommand("addrole", "Add a joinable role", true, f.handleAddCustomRole)
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

	// Ignore any users on the ignore list
	if f.Permissions.CheckIgnoredUser(m.Author) {
		log.Printf("[*] Ignoring command from ignored user.")
		return nil
	}

	cName := strings.Split(cCmd, cKey)[1]

	if cmd, ok := f.ScuzzyCommands[cName]; ok {
		if cmd.AdminOnly && !f.Permissions.CheckAdminRole(m.Member) {
			return errors.New("You do not have permissions to use this command.")
		}

		log.Printf("[*] Running command %s (Requested by %s)\n", cName, m.Author.Username)

		err := cmd.Handler(s, m)
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
		err = s.GuildMemberRoleAdd(f.Config.GuildID, m.User.ID, roleID)
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
			log.Println("[!] Error (Guild Member Joined): " + err.Error())
		}
		break
	}
}
