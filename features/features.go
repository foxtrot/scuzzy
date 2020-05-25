package features

import (
	"github.com/foxtrot/scuzzy/models"
	"github.com/foxtrot/scuzzy/permissions"
)

type Features struct {
	Token       string
	Permissions *permissions.Permissions
	Config      models.Configuration
}
