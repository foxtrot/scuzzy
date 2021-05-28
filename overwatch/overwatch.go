package overwatch

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type UserMessageStat struct {
	UserID           string
	Username         string
	MessagesLastDay  uint64
	MessagesLastHour uint64
	MessagesLastFive uint64
}

type Overwatch struct {
	TotalMessages uint64
	UserMessages  map[string]UserMessageStat
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
	for _, user := range o.UserMessages {
		log.Printf("User: %+v\n", user)
	}

	userID := m.Author.ID
	user, ok := o.UserMessages[userID]
	if !ok {
		log.Println("Couldn't find user, making a new one")
		o.UserMessages[userID] = UserMessageStat{
			UserID:           userID,
			Username:         m.Author.Username,
			MessagesLastDay:  1,
			MessagesLastHour: 2,
			MessagesLastFive: 3,
		}
		user = o.UserMessages[userID]
	}

	user.MessagesLastDay++
	user.MessagesLastHour++
	user.MessagesLastFive++

	return nil
}

func (o *Overwatch) Run() {
	go func() {
		for range time.Tick(10 * time.Second) {
			log.Println("Printing UserMessages (10 Seconds)...")
			for _, user := range o.UserMessages {
				log.Printf("User: %+v\n", user)
			}
		}

		for range time.Tick(5 * time.Minute) {
			log.Printf("[*] Resetting all users 5 minute message counters")
			for _, user := range o.UserMessages {
				user.MessagesLastFive = 0
			}
		}

		for range time.Tick(10 * time.Minute) {
			log.Printf("[*] Resetting all users 10 minute message counters")
			for _, user := range o.UserMessages {
				user.MessagesLastHour = 0
			}
		}
	}()
}
