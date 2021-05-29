package overwatch

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/foxtrot/scuzzy/commands"
)

type UserMessageStat struct {
	UserID               string
	Username             string
	MessagesLastDay      uint64
	MessagesLastHour     uint64
	MessagesLastFiveMins uint64
	MessagesLastTenSecs  uint64
	Kicks                int
}

type ServerStat struct {
	JoinsLastTenMins       uint64
	SlowmodeFlood          bool
	SlowmodeFloodStartTime time.Time
}

type Overwatch struct {
	TotalMessages uint64
	UserMessages  map[string]*UserMessageStat
	ServerStats   ServerStat
	Commands      *commands.Commands
}

func (o *Overwatch) ProcessMessage(s *discordgo.Session, m interface{}) {
	switch m.(type) {
	case *discordgo.MessageCreate:
		err := o.handleUserStat(s, m.(*discordgo.MessageCreate))
		if err != nil {
			log.Printf("[!] Error handling Overwatch user stat: %s\n", err.Error())
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

func (o *Overwatch) handleUserStat(s *discordgo.Session, m *discordgo.MessageCreate) error {
	userID := m.Author.ID
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
	if o.ServerStats.JoinsLastTenMins > 10 {
		log.Printf("[*] User flood detected, enforcing slow mode on all channels for 30 minutes\n")
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
				// load the threshold from the config file, dipshit
				if user.MessagesLastTenSecs > 10 {
					// Set slow mode, kick user? add kick count?
					if user.Kicks > 2 {
						// ban that sucker
						delete(o.UserMessages, user.UserID)
						log.Printf("[*] User %s (%s) was banned due to previous spam-related kicks", user.Username, user.UserID)
					} else {
						user.Kicks++
						// kick user
						user.MessagesLastTenSecs = 0
						log.Printf("[*] User %s (%s) has triggered the message threshold.", user.Username, user.UserID)
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
				if time.Since(o.ServerStats.SlowmodeFloodStartTime) > time.Minute*30 {
					log.Printf("[*] Removing Slowmode for all channels after flood\n")
				}
			}
		}
	}()

	// Clear Counters
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
