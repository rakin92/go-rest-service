package server

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitter"
	"github.com/rakin92/go-rest-service/pkg/cfg"
)

// initializeAuthProviders does just that, with Goth providers
func initializeAuthProviders(sc *cfg.Server) error {
	providers := []goth.Provider{}
	// Initialize Goth providers
	for _, p := range sc.AuthProviders {
		switch p.Provider {
		case "facebook":
			providers = append(providers, facebook.New(p.ClientKey, p.Secret,
				sc.SchemaVersionedEndpoint("/auth/"+p.Provider+"/callback"),
				p.Scopes...))
		case "google":
			providers = append(providers, google.New(p.ClientKey, p.Secret,
				sc.SchemaVersionedEndpoint("/auth/"+p.Provider+"/callback"),
				p.Scopes...))
		case "twitter":
			providers = append(providers, twitter.New(p.ClientKey, p.Secret,
				sc.SchemaVersionedEndpoint("/auth/"+p.Provider+"/callback")))
		}
	}
	goth.UseProviders(providers...)
	return nil
}
