package features

import (
	"github.com/bwmarrin/discordgo"
	"github.com/foxtrot/scuzzy/models"
	"github.com/foxtrot/scuzzy/permissions"
)

type ScuzzyHandler func(session *discordgo.Session, m *discordgo.MessageCreate) error

type ScuzzyCommand struct {
	Name        string
	Description string
	AdminOnly   bool
	Handler     ScuzzyHandler
}

type Features struct {
	Token          string
	Permissions    *permissions.Permissions
	Config         *models.Configuration
	ScuzzyCommands map[string]ScuzzyCommand
}
