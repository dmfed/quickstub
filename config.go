package quickstub

// Config is used to initialize Quickstub.
type Config struct {
	Version    uint                `yaml:"version"`     // Version indicates version of the config (may bee ommitted for current version)
	ListenAddr string              `yaml:"listen_addr"` // ListenAddr is the address for webserver to listen in form "host:port" (like "127.0.0.1:8080 or somehost:80") or simply port :8080 to accept connections on any interface.
	Endpoints  map[string]Response `yaml:"endpoints"`   // Endpoints tell Quickstub how to configure stub handlers.
}

type Response struct {
	Code    int               `yaml:"code"`    // HTTP response code to be returned.
	Headers map[string]string `yaml:"headers"` // response headers
	Body    string            `yaml:"body"`    // HTTP response body.
	File    string            `yaml:"file"`    // File tells quickstub to read contents from path specified in this field. If Body field is not wmpty this field is not functional.
}
