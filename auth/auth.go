package auth

import (
	"github.com/bwmarrin/discord.go"
	"github.com/foxtrot/scuzzy/models"
	"strings"
)

type AdminRole struct {
	Name string
	ID   string
}

type Auth struct {
	AdminRoles          []AdminRole
	CommandRestrictions []models.CommandRestriction
	Guild               *discordgo.Guild

	Config *models.Configuration
}

func New(config *models.Configuration, guild *discordgo.Guild) *Auth {
	var ars []AdminRole
	for _, gRole := range guild.Roles {
		for _, aRole := range config.AdminRoles {
			if aRole != gRole.Name {
				continue
			}

			ar := AdminRole{
				Name: gRole.Name,
				ID:   gRole.ID,
			}
			ars = append(ars, ar)
		}
	}

	var crs []models.CommandRestriction
	for _, cRes := range config.CommandRestrictions {
		cr := models.CommandRestriction{
			Command:  cRes.Command,
			Mode:     cRes.Mode,
			Channels: cRes.Channels,
		}
		crs = append(crs, cr)
	}

	return &Auth{
		AdminRoles:          ars,
		CommandRestrictions: crs,
		Config:              config,
		Guild:               guild,
	}
}

func (a *Auth) CheckAdminRole(m *discordgo.Member) bool {
	for _, aR := range a.AdminRoles {
		for _, mID := range m.Roles {
			if aR.ID == mID {
				return true
			}
		}
	}

	return false
}

func (a *Auth) CheckCommandRestrictions(m *discordgo.MessageCreate) bool {
	cName := strings.Split(m.Content, " ")[0]
	cName = strings.Replace(cName, a.Config.CommandKey, "", 1)
	cChanID := m.ChannelID

	for _, cR := range a.CommandRestrictions {
		if cName == cR.Command {
			for _, cID := range cR.Channels {
				if cID == cChanID && cR.Mode == "white" {
					return true
				} else if cID == cChanID && cR.Mode == "black" {
					return false
				}
			}

			if cR.Mode == "white" {
				return false
			} else {
				return true
			}
		}
	}

	return true
}
