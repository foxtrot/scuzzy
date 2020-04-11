package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/foxtrot/scuzzy/auth"
	"github.com/foxtrot/scuzzy/features"
	"github.com/foxtrot/scuzzy/models"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discord.go"
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
		log.Fatal("Error: No Auth Token supplied.")
	}
	if len(ConfigPath) == 0 {
		log.Fatal("Error: No Config Path supplied.")
	}

	// Get Config
	err := getConfig()
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	// Instantiate Bot
	bot, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	// Open Connection
	err = bot.Open()
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	// Setup Auth
	Config.Guild, err = bot.Guild(Config.GuildID)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
	var a *auth.Auth
	a = auth.New(&Config, Config.Guild)

	// Setup Handlers
	f := features.Features{
		Token:  Token,
		Auth:   a,
		Config: Config,
	}

	// Register Handlers
	bot.AddHandler(f.OnMessageCreate)
	bot.AddHandler(f.OnUserJoin)

	fmt.Println("Bot Running.")

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
			log.Fatal("Error: " + err.Error())
		}

		// For some reason the bot's status will regularly disappear...
		for _ = range time.Tick(10 * time.Minute) {
			err := bot.UpdateStatusComplex(usd)
			if err != nil {
				log.Fatal("Error: " + err.Error())
			}
		}
	}()

	// Catch SIGINT
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGKILL)
	<-sc

	err = bot.Close()
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
}
