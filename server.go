package quickstub

import (
	"net/http"
	"os"
)

// NewServer returns *http.Server with Addr set to
// conf.ListenAddr and Handler set to
// ServerMux initialized with handlers created from conf.
func NewServer(conf *Config) (*http.Server, error) {
	mux, err := NewMux(conf)
	if err != nil {
		return nil, err
	}
	srv := &http.Server{
		Addr:    conf.ListenAddr,
		Handler: mux,
	}
	return srv, nil
}

// NewMux creates handler function returning HTTP code and body
// as specified in the coniguration and registers it with an instance of
// *http.ServeMux then return this instance of ServeMux.
func NewMux(conf *Config) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	for endpoint, c := range conf.Endpoints {
		var contents []byte
		if c.Body != "" {
			contents = []byte(c.Body)
		} else if c.File != "" {
			b, err := os.ReadFile(c.File)
			if err != nil {
				return nil, err
			}
			contents = b
		} else {
			contents = []byte{}
		}

		f := func(w http.ResponseWriter, r *http.Request) {
			for k, v := range c.Headers {
				w.Header().Add(k, v)
			}
			w.WriteHeader(c.Code)
			w.Write(contents)
		}
		mux.HandleFunc(endpoint, f)
	}
	return mux, nil
}
