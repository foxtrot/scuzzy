package models

import "github.com/bwmarrin/discord.go"

type ColorRole struct {
	Name string `json:"color"`
	ID   string `json:"id"`
}

type CommandRestriction struct {
	Command  string   `json:"command"`
	Mode     string   `json:"mode"`
	Channels []string `json:"channels"`
}

type Configuration struct {
	CommandKey string `json:"command_key"`

	GuildID string `json:"guild_id"`

	StatusText  string `json:"status_text"`
	WelcomeText string `json:"welcome_text"`
	RulesText   string `json:"rules_text"`

	AdminRoles  []string `json:"admin_roles"`
	JoinRoleIDs []string `json:"join_role_ids"`

	CommandRestrictions []CommandRestriction `json:"command_restrictions"`

	ColorRoles []ColorRole `json:"color_roles"`

	Guild *discordgo.Guild `json:"reserved_guild"`
}
