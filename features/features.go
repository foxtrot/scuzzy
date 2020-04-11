package features

import (
	"github.com/foxtrot/scuzzy/auth"
	"github.com/foxtrot/scuzzy/models"
)

type Features struct {
	Token  string
	Auth   *auth.Auth
	Config models.Configuration
}
