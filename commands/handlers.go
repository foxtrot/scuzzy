package commands

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (c *Commands) RegisterCommand(name string, description string, adminonly bool, handler ScuzzyHandler) {
	log.Printf("[*] Registering Command '%s'\n", name)
	co := ScuzzyCommand{
		Index:       len(c.ScuzzyCommands) + 1,
		Name:        name,
		Description: description,
		AdminOnly:   adminonly,
		Handler:     handler,
	}
	c.ScuzzyCommands[name] = co
	c.ScuzzyCommandsByIndex[co.Index] = co
}

func (c *Commands) RegisterHandlers() {
	c.ScuzzyCommands = make(map[string]ScuzzyCommand)
	c.ScuzzyCommandsByIndex = make(map[int]ScuzzyCommand)

	// Misc Commands
	c.RegisterCommand("help", "Show Help Text", false, c.handleHelp)
	c.RegisterCommand("info", "Show Bot Info", false, c.handleInfo)
	c.RegisterCommand("md", "Show common Discord MarkDown formatting", false, c.handleMarkdownInfo)
	c.RegisterCommand("userinfo", "Display a users information", false, c.handleUserInfo)
	c.RegisterCommand("serverinfo", "Display the current servers information", false, c.handleServerInfo)
	c.RegisterCommand("no", "", false, c.handleCat)

	// User Settings
	c.RegisterCommand("colours", "Display available colour roles", false, c.handleUserColors)
	c.RegisterCommand("colors", "", false, c.handleUserColors)
	c.RegisterCommand("colour", "Set a colour role for yourself", false, c.handleUserColor)
	c.RegisterCommand("color", "", false, c.handleUserColor)
	c.RegisterCommand("listroles", "List user joinable roles", false, c.handleListCustomRoles)
	c.RegisterCommand("joinrole", "Join an available role for yourself", false, c.handleJoinCustomRole)
	c.RegisterCommand("leaverole", "Leave an available role", false, c.handleLeaveCustomRole)

	// Conversion Helpers
	c.RegisterCommand("ctof", "Convert Celsius to Farenheit", false, c.handleCtoF)
	c.RegisterCommand("ftoc", "Convert Farenheit to Celsius", false, c.handleFtoC)
	c.RegisterCommand("metofe", "Convert Meters to Feet", false, c.handleMetersToFeet)
	c.RegisterCommand("fetome", "Convert Feet to Meters", false, c.handleFeetToMeters)
	c.RegisterCommand("cmtoin", "Convert Centimeters to Inches", false, c.handleCentimeterToInch)
	c.RegisterCommand("intocm", "Convert Inches to Centimeters", false, c.handleInchToCentimeter)

	// Admin Commands
	c.RegisterCommand("ping", "Ping Scuzzy", true, c.handlePing)
	c.RegisterCommand("rules", "Display the Server Rules", true, c.handleRules)
	c.RegisterCommand("status", "Set Bot Status", true, c.handleSetStatus)
	c.RegisterCommand("purge", "Purge Channel Messages", true, c.handlePurgeChannel)
	c.RegisterCommand("kick", "Kick a User", true, c.handleKickUser)
	c.RegisterCommand("ban", "Ban a User", true, c.handleBanUser)
	c.RegisterCommand("slow", "Set Channel Slow Mode", true, c.handleSetSlowmode)
	c.RegisterCommand("unslow", "Unset Channel Slow Mode", true, c.handleUnsetSlowmode)
	c.RegisterCommand("ignore", "Add a user to Scuzzy's ignore list", true, c.handleIgnoreUser)
	c.RegisterCommand("unignore", "Remove a user from Scuzzy's ignore list", true, c.handleUnIgnoreUser)
	c.RegisterCommand("setconfig", "Set Configuration", true, c.handleSetConfig)
	c.RegisterCommand("getconfig", "Print Configuration", true, c.handleGetConfig)
	c.RegisterCommand("saveconfig", "Save Configuration to Disk", true, c.handleSaveConfig)
	c.RegisterCommand("reloadconfig", "Reload Configuration", true, c.handleReloadConfig)
	c.RegisterCommand("addrole", "Add a joinable role", true, c.handleAddCustomRole)
}

func (c *Commands) ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	cKey := c.Config.CommandKey
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
	if c.Permissions.CheckIgnoredUser(m.Author) {
		log.Printf("[*] Ignoring command from ignored user.")
		return nil
	}

	cName := strings.Split(cCmd, cKey)[1]

	if cmd, ok := c.ScuzzyCommands[cName]; ok {
		if cmd.AdminOnly && !c.Permissions.CheckAdminRole(m.Member) {
			log.Printf("[*] User %s tried to run admin command %s\n", m.Author.Username, cName)
			return nil
		}

		log.Printf("[*] Running command %s (Requested by %s)\n", cName, m.Author.Username)

		err := cmd.Handler(s, m)
		if err != nil {
			log.Printf("[!] Command %s (Requested by %s) had error: '%s'\n", cName, m.Author.Username, err.Error())

			eMsg := c.CreateDefinedEmbed("Error ("+cName+")", err.Error(), "error", m.Author)
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Commands) ProcessMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) error {
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

	embed := c.CreateDefinedEmbed("Deleted Message", msg, "", msgAuthor)
	_, err := s.ChannelMessageSendEmbed(c.Config.LoggingChannel, embed)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) ProcessMessageDeleteBulk(s *discordgo.Session, m *discordgo.MessageDeleteBulk) error {
	msgChannelID := m.ChannelID

	msg := "`Channel` - <#" + msgChannelID + ">\n"
	msg += "Message IDs: \n"
	for k, v := range m.Messages {
		msg += strconv.Itoa(k) + ": `" + v + "`\n"
	}

	embed := c.CreateDefinedEmbed("Deleted Bulk Messages", msg, "", nil)
	_, err := s.ChannelMessageSendEmbed(c.Config.LoggingChannel, embed)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) ProcessUserJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) error {
	userChannel, err := s.UserChannelCreate(m.User.ID)
	if err != nil {
		log.Print("[!] Error (User Join): " + err.Error())
		return err
	}

	_, err = s.ChannelMessageSend(userChannel.ID, c.Config.WelcomeText)
	if err != nil {
		log.Print("[!] Error (User Join): " + err.Error())
		return err
	}

	for _, roleID := range c.Config.JoinRoleIDs {
		err = s.GuildMemberRoleAdd(c.Config.GuildID, m.User.ID, roleID)
		if err != nil {
			log.Print("[!] Error (User Join)" + err.Error())
			return err
		}
	}

	return nil
}

func (c *Commands) ProcessMessage(s *discordgo.Session, m interface{}) {
	switch m.(type) {
	case *discordgo.MessageCreate:
		// Pass Messages to the command processor
		err := c.ProcessCommand(s, m.(*discordgo.MessageCreate))
		if err != nil {
			log.Println("[!] Error: " + err.Error())
		}
		break
	case *discordgo.MessageDelete:
		// Log deleted messages to the logging channel.
		err := c.ProcessMessageDelete(s, m.(*discordgo.MessageDelete))
		if err != nil {
			eMsg := c.CreateDefinedEmbed("Error (Message Deleted)", err.Error(), "error", nil)
			_, err = s.ChannelMessageSendEmbed(c.Config.LoggingChannel, eMsg)
			if err != nil {
				log.Println("[!] Error " + err.Error())
			}
		}
		break
	case *discordgo.MessageDeleteBulk:
		err := c.ProcessMessageDeleteBulk(s, m.(*discordgo.MessageDeleteBulk))
		if err != nil {
			eMsg := c.CreateDefinedEmbed("Error (Message Bulk Deleted)", err.Error(), "error", nil)
			_, err = s.ChannelMessageSendEmbed(c.Config.LoggingChannel, eMsg)
			if err != nil {
				log.Println("[!] Error " + err.Error())
			}
		}
		break
	case *discordgo.GuildMemberAdd:
		// Handle new member (Welcome message, etc)
		err := c.ProcessUserJoin(s, m.(*discordgo.GuildMemberAdd))
		if err != nil {
			log.Println("[!] Error (Guild Member Joined): " + err.Error())
		}
		break
	}
}
