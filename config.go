package quickstub

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Version uint

const (
	Version1 Version = iota + 1
	Version2
)

// Pattern must take a form of either request method
// followed by path as in "POST /post/path" or just the path
// is in "/some/endpoint". See https://pkg.go.dev/net/http#ServeMux for
// description of acceptable patterns.
type Pattern string

// Config is used to initialize Quickstub, http.Server of http.ServeMux
// with this package with specific request handlers.
type Config struct {
	Version       Version   `yaml:"version" validate:"required,gte=1"`             // Version indicates version of the config (may bee ommitted for current version)
	ListenAddr    string    `yaml:"listen_addr" validate:"required,hostname_port"` // ListenAddr is the address for webserver to listen in form "host:port" (like "127.0.0.1:8080 or somehost:80") or simply port :8080 to accept connections on any interface.
	MagicEndpoint string    `yaml:"magic_endpoint"`                                // path to use for magic endpoint
	Endpoints     Endpoints `yaml:"endpoints" validate:"required"`                 // Endpoints tell Quickstub how to configure stub handlers.
}

// Endpoints keys must be populated with Patterns.
// Corresponding Response values tell server what to respond
// when a Pattern matches.
type Endpoints map[Pattern]Response

type Response struct {
	Code    int               `yaml:"code"`    // HTTP response code to be returned.
	Headers map[string]string `yaml:"headers"` // HTTP response headers.
	Body    string            `yaml:"body"`    // HTTP response body. If string in this field is prefixed with "@" symbol it is considered a path to file to read contents from.
}

func ParseConfig(b []byte) (*Config, error) {
	return parseConfig(b)
}

func parseConfig(b []byte) (*Config, error) {
	c := new(Config)
	if err := yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}

	if err := validateConfig(c); err != nil {
		return nil, err
	}
	return c, nil
}

func validateConfig(c *Config) error {
	err := validator.New().Struct(c)
	return err
}
