package overwatch

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/foxtrot/scuzzy/commands"
	"github.com/foxtrot/scuzzy/models"
)

type UserMessageStat struct {
	UserID               string
	Username             string
	MessagesLastDay      int
	MessagesLastHour     int
	MessagesLastFiveMins int
	MessagesLastTenSecs  int
	Warnings             int
	Kicks                int
}

type ServerStat struct {
	JoinsLastTenMins       int
	SlowmodeFlood          bool
	SlowmodeFloodStartTime time.Time
}

type Overwatch struct {
	TotalMessages uint64
	UserMessages  map[string]*UserMessageStat
	ServerStats   ServerStat
	Commands      *commands.Commands
	Config        *models.Configuration
}

func (o *Overwatch) ProcessMessage(s *discordgo.Session, m interface{}) {
	switch m.(type) {
	case *discordgo.MessageCreate:
		err := o.handleUserStat(s, m.(*discordgo.MessageCreate))
		if err != nil {
			log.Printf("[!] Error handling Overwatch user stat: %s\n", err.Error())
		}

		err = o.filterUserMessage(s, m.(*discordgo.MessageCreate))
		if err != nil {
			log.Printf("[!] Error handling Overwatch user message filter: %s\n", err.Error())
		}
		break
	case *discordgo.GuildMemberAdd:
		err := o.handleServerJoin(s, m.(*discordgo.GuildMemberAdd))
		if err != nil {
			log.Printf("[!] Error handling Overwatch server join: %s\n", err.Error())
		}
		break
	}
}

func (o *Overwatch) filterUserMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if o.Config.FilterLanguage {
		// load regex from config?
	}
	return nil
}

func (o *Overwatch) handleUserStat(s *discordgo.Session, m *discordgo.MessageCreate) error {
	userID := m.Author.ID

	if userID == s.State.User.ID {
		// Ignore bots own actions
		return nil
	}

	user, ok := o.UserMessages[userID]
	if !ok {
		o.UserMessages[userID] = &UserMessageStat{
			UserID:   userID,
			Username: m.Author.Username,
		}
		user = o.UserMessages[userID]
	}

	user.MessagesLastDay++
	user.MessagesLastHour++
	user.MessagesLastFiveMins++
	user.MessagesLastTenSecs++

	return nil
}

func (o *Overwatch) handleServerJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) error {
	o.ServerStats.JoinsLastTenMins++

	// json value
	if o.ServerStats.JoinsLastTenMins > o.Config.JoinFloodThreshold {
		log.Printf("[*] User flood detected, enforcing slow mode on all channels for 30 minutes\n")
		// Set slow mode on all channels
		o.ServerStats.SlowmodeFlood = true
		o.ServerStats.SlowmodeFloodStartTime = time.Now()
		o.ServerStats.JoinsLastTenMins = 0
	}
	return nil
}

// this is fucking amazing code
func (o *Overwatch) Run() {
	// State of the art anti-spam loop
	go func() {
		for range time.Tick(5 * time.Second) {
			for _, user := range o.UserMessages {
				if user.MessagesLastTenSecs > o.Config.UserMessageThreshold {
					// Set slow mode, kick user? add kick count?
					if user.Warnings > o.Config.MaxUserWarnings {
						if user.Kicks > o.Config.MaxUserKicks {
							// ban that sucker
							delete(o.UserMessages, user.UserID)
							log.Printf("[*] User %s (%s) was banned due to previous spam-related kicks\n", user.Username, user.UserID)
						} else {
							user.Kicks++
							// kick user
							user.MessagesLastTenSecs = 0
							log.Printf("[*] User %s (%s) has been kicked for message spam\n", user.Username, user.UserID)
						}
					} else {
						user.Warnings++
						log.Printf("[*] User %s (%s) was warned for spamming\n", user.Username, user.UserID)
					}
				}
			}
		}
	}()

	// State of the art anti-join-flood loop
	go func() {
		for range time.Tick(30 * time.Second) {
			if o.ServerStats.SlowmodeFlood {
				// json value
				if o.ServerStats.JoinsLastTenMins > 10 {
					if time.Since(o.ServerStats.SlowmodeFloodStartTime) > time.Minute*30 {
						log.Printf("[*] Removing Slowmode for all channels after flood\n")
						o.ServerStats.SlowmodeFlood = false
						// remove slow mode
					}
				} else {
					log.Printf("[*] Exending Slowmode due to sustained join flood\n")
					o.ServerStats.SlowmodeFloodStartTime = time.Now()
				}
			}
		}
	}()

	// Clear Server Counters
	go func() {
		for range time.Tick(10 * time.Minute) {
			o.ServerStats.JoinsLastTenMins = 0
		}
	}()

	// Clear User Counters
	go func() {
		for range time.Tick(24 * time.Hour) {
			for _, user := range o.UserMessages {
				if user.MessagesLastDay == 0 {
					delete(o.UserMessages, user.UserID)
				} else {
					user.MessagesLastDay = 0
				}
			}
		}
		for range time.Tick(1 * time.Hour) {
			for _, user := range o.UserMessages {
				user.MessagesLastHour = 0
			}
		}
		for range time.Tick(5 * time.Minute) {
			for _, user := range o.UserMessages {
				user.MessagesLastFiveMins = 0
			}
		}
		for range time.Tick(10 * time.Second) {
			for _, user := range o.UserMessages {
				user.MessagesLastTenSecs = 0
			}
		}
	}()
}
