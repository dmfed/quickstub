package quickstub

import (
	"net/http"
)

// NewServer returns *http.Server with Addr set to
// conf.ListenAddr and Handler set to
// ServeMux initialized with handlers as cpecified by conf.
func NewServer(listenAddr string, endpoints Endpoints) (*http.Server, error) {
	mux, err := newEndpointsMux(endpoints)
	if err != nil {
		return nil, err
	}
	return newServerFromMux(listenAddr, mux), nil
}

func newServerFromMux(listenAddr string, m *http.ServeMux) *http.Server {
	srv := &http.Server{
		Addr:    listenAddr,
		Handler: m,
	}
	return srv
}

func makeNewServer(c *Config) (*http.Server, chan *reconfReq, error) {
	mux, err := newEndpointsMux(c.Endpoints)
	if err != nil {
		return nil, nil, err
	}
	magicChan := make(chan *reconfReq)

	if c.MagicEndpoint != "" {
		mux.HandleFunc(c.MagicEndpoint, newMagicHandler(magicChan, c).HandlerFunc)
	}
	return newServerFromMux(c.ListenAddr, mux), magicChan, nil
}
