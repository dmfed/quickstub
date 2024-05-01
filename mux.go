package quickstub

import (
	"fmt"
	"net/http"
	"os"
)

// NewMux creates handler functions returning HTTP code and body
// as specified in the coniguration and registers it with an instance of
// *http.ServeMux then return this instance of ServeMux.
func NewMux(endpoints Endpoints) (*http.ServeMux, error) {
	return newEndpointsMux(endpoints)
}

func newEndpointsMux(endpoints Endpoints) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	for pattern, c := range endpoints {
		var contents []byte
		if len(c.Body) > 1 && c.Body[0] == '@' {
			// this must be the name of the file
			// we sould read
			filepath := c.Body[1:]
			b, err := os.ReadFile(filepath)
			if err != nil {
				return nil, fmt.Errorf("error opening file %s provided for endpoint %s: %w", filepath, pattern, err)
			}
			contents = b
		} else if len(c.Body) > 1 && c.Body[0] == '\\' {
			// first character is escaped
			contents = []byte(c.Body[1:])
		} else {
			contents = []byte(c.Body)
		}

		f := func(w http.ResponseWriter, r *http.Request) {
			for k, v := range c.Headers {
				w.Header().Add(k, v)
			}
			w.WriteHeader(c.Code)
			w.Write(contents)
		}
		mux.HandleFunc(string(pattern), f)
	}
	return mux, nil
}
