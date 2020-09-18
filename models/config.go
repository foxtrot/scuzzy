package models

import discordgo "github.com/bwmarrin/discord.go"

type ColorRole struct {
	Name string `json:"color"`
	ID   string `json:"id"`
}

type CustomRole struct {
	Name      string `json:"role_name"`
	ShortName string `json:"short_name"`
	ID        string `json:"id"`
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

	ColorRoles  []ColorRole  `json:"color_roles"`
	CustomRoles []CustomRole `json:"custom_roles"`

	IgnoredUsers []string `json:"ignored_users"`

	LoggingChannel string `json:"logging_channel"`

	Guild      *discordgo.Guild
	ConfigPath string
}
