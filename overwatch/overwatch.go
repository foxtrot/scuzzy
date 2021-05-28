package overwatch

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type UserMessageStat struct {
	UserID               string
	Username             string
	MessagesLastDay      uint64
	MessagesLastHour     uint64
	MessagesLastFiveMins uint64
	MessagesLastTenSecs  uint64
}

type Overwatch struct {
	TotalMessages uint64
	UserMessages  map[string]*UserMessageStat
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
		break
	}
}

func (o *Overwatch) handleUserStat(s *discordgo.Session, m *discordgo.MessageCreate) error {
	userID := m.Author.ID
	user, ok := o.UserMessages[userID]
	if !ok {
		o.UserMessages[userID] = &UserMessageStat{
			UserID:               userID,
			Username:             m.Author.Username,
			MessagesLastDay:      0,
			MessagesLastHour:     0,
			MessagesLastFiveMins: 0,
			MessagesLastTenSecs:  0,
		}
		user = o.UserMessages[userID]
	}

	user.MessagesLastDay++
	user.MessagesLastHour++
	user.MessagesLastFiveMins++
	user.MessagesLastTenSecs++

	return nil
}

// this is fucking amazing code
func (o *Overwatch) Run() {
	// Five second loop
	go func() {
		for range time.Tick(10 * time.Second) {
			for _, user := range o.UserMessages {
				// load the threshold from the config file
				if user.MessagesLastTenSecs > 10 {
					// Set slow mode, kick user? add kick count?
					log.Printf("[*] User %s (%s) has triggered the message threshold.", user.Username, user.UserID)
				}
			}
		}
	}()

	// Clear Counters
	go func() {
		for range time.Tick(10 * time.Second) {
			log.Printf("[*] Resetting all users 10 second message counters")
			for _, user := range o.UserMessages {
				user.MessagesLastTenSecs = 0
			}
		}

		for range time.Tick(5 * time.Minute) {
			log.Printf("[*] Resetting all users 5 minute message counters")
			for _, user := range o.UserMessages {
				user.MessagesLastFiveMins = 0
			}
		}

		for range time.Tick(1 * time.Hour) {
			log.Printf("[*] Resetting all users 60 minute message counters")
			for _, user := range o.UserMessages {
				user.MessagesLastHour = 0
			}
		}

		for range time.Tick(24 * time.Hour) {
			log.Println("[*] Resetting all users 1 day message counters")
			for _, user := range o.UserMessages {
				user.MessagesLastDay = 0
			}
		}
	}()
}
