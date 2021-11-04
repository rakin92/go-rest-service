// Package cfg is the configuration package hold all config objects
package cfg

// Server defines the configuration for the server
type Server struct {
	ServiceName    string
	Version        string
	Env            string
	Host           string
	Port           string
	URISchema      string
	ServiceVersion string
	SessionSecret  string
	JWT            JWT
	Cache          Cache
	Database       DB
	AuthProviders  []AuthProvider
}

// JWT defines the options for JWT tokens
type JWT struct {
	Secret    string
	Algorithm string
}

// Cache defines the configuration for the cache
type Cache struct {
	Server      string
	Password    string
	TimeoutHour int
}

// DB defines the configuration for the DB config
type DB struct {
	Dialect     string
	DSN         string
	SeedDB      bool
	LogMode     bool
	AutoMigrate bool
	MaxCon      int
	MaxIdleCon  int
}

// AuthProvider defines the configuration for the Goth config
type AuthProvider struct {
	Provider  string
	ClientKey string
	Secret    string
	Domain    string // If needed, like with auth0
	Scopes    []string
}

func getValidHost(host string) string {
	if host == ":" {
		return "localhost"
	}
	return host
}

// ListenEndpoint builds the endpoint string (host + port)
func (s *Server) ListenEndpoint() string {
	if s.Port == "80" {
		return s.Host
	}
	if s.Host == ":" {
		return s.Host + s.Port

	}
	return s.Host + ":" + s.Port
}

// VersionedEndpoint builds the endpoint `string (host + port + version)
func (s *Server) VersionedEndpoint(path string) string {
	if s.ServiceVersion == "" {
		return "/v1" + path
	}
	return "/" + s.ServiceVersion + path
}

// SchemaVersionedEndpoint builds the schema endpoint string (schema + host + port + version)
func (s *Server) SchemaVersionedEndpoint(path string) string {
	if s.Port == "80" {
		return s.URISchema + getValidHost(s.Host) + "/" + s.ServiceVersion + path
	}
	return s.URISchema + getValidHost(s.Host) + ":" + s.Port + "/" + s.ServiceVersion + path
}
