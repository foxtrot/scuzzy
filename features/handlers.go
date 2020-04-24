package features

import (
	"github.com/bwmarrin/discord.go"
	"log"
	"strings"
)

func (f *Features) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	// Ignore the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	cKey := f.Config.CommandKey
	cName := strings.Split(m.Content, " ")[0]

	if !strings.HasPrefix(cName, cKey) {
		return
	}

	switch cName {
	/* Misc Commands */
	case cKey + "help":
		err = f.handleHelp(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Help)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "info":
		err = f.handleInfo(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Info)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "ping":
		err = f.handlePing(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Ping)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "no":
		err = f.handleCat(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (No)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "md":
		err = f.handleMarkdownInfo(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Markdown)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "ctof":
		err = f.handleCtoF(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (CtoF)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "ftoc":
		err = f.handleFtoC(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (FtoC)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "metofe":
		err = f.handleMetersToFeet(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Meters to Feet)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "fetome":
		err = f.handleFeetToMeters(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Feet to Meters)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "cmtoin":
		err = f.handleCentimeterToInch(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Meters to Feet)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "intocm":
		err = f.handleInchToCentimeter(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Feet to Meters)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "userinfo":
		err = f.handleUserInfo(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (User Info)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	/* Moderation */
	case cKey + "purge":
		err = f.handlePurgeChannel(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Purge)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
	/* Bot Control */
	case cKey + "status":
		err = f.handleSetStatus(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Status)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
	/* Role Colors */
	case cKey + "colors":
		err = f.listUserColors(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Colors)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	case cKey + "color":
		err = f.setUserColor(s, m)
		if err != nil {
			eMsg := f.CreateDefinedEmbed("Error (Color)", err.Error(), "error")
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		break
	}
}

func (f *Features) OnUserJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
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
