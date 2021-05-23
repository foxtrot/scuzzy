package main

import (
	"encoding/json"
	"flag"
	"github.com/foxtrot/scuzzy/features"
	"github.com/foxtrot/scuzzy/models"
	"github.com/foxtrot/scuzzy/permissions"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Core Bot Properties
var (
	Token      string
	ConfigPath string
	Config     models.Configuration
)

func getConfig() error {
	cf, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(cf, &Config)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Parse and Check Flags
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&ConfigPath, "c", "", "Config Path")
	flag.Parse()

	if len(Token) == 0 {
		log.Fatal("[!] Error: No Auth Token supplied.")
	}
	if len(ConfigPath) == 0 {
		log.Fatal("[!] Error: No Config Path supplied.")
	}

	// Get Config
	err := getConfig()
	if err != nil {
		log.Fatal("[!] Error: " + err.Error())
	}
	Config.ConfigPath = ConfigPath

	// Instantiate Bot
	bot, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("[!] Error: " + err.Error())
	}

	// Enable Message Caching (Last 1024 Events)
	bot.State.MaxMessageCount = 1024
	bot.State.TrackChannels = true
	bot.State.TrackMembers = true

	// Open Connection
	err = bot.Open()
	if err != nil {
		log.Fatal("[!] Error: " + err.Error())
	}

	// Setup Auth
	g, err := bot.Guild(Config.GuildID)
	if err != nil {
		log.Fatal("[!] Error: " + err.Error())
	}
	var p *permissions.Permissions
	p = permissions.New(&Config, g)

	Config.GuildName = g.Name

	// Setup Handlers
	f := features.Features{
		Token:       Token,
		Permissions: p,
		Config:      &Config,
	}

	// Register Handlers
	f.RegisterHandlers()
	bot.AddHandler(f.ProcessMessage)

	log.Printf("[*] Bot Running.\n")

	// Set Bot Status
	go func() {
		usd := discordgo.UpdateStatusData{
			IdleSince: nil,
			Game: &discordgo.Game{
				Name:          Config.StatusText,
				Type:          0,
				URL:           "",
				Details:       "",
				State:         "",
				TimeStamps:    discordgo.TimeStamps{},
				Assets:        discordgo.Assets{},
				ApplicationID: "",
				Instance:      -1,
			},
			AFK:    false,
			Status: "online",
		}
		err = bot.UpdateStatusComplex(usd)
		if err != nil {
			log.Fatal("[!] Error: " + err.Error())
		}

		// For some reason the bot's status will regularly disappear...
		for range time.Tick(10 * time.Minute) {
			err := bot.UpdateStatusComplex(usd)
			if err != nil {
				log.Fatal("[!] Error: " + err.Error())
			}
		}
	}()

	// Catch SIGINT
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGKILL)
	<-sc

	err = bot.Close()
	if err != nil {
		log.Fatal("[!] Error: " + err.Error())
	}
}
