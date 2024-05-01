package quickstub

import (
	"io"
	"log"
	"net/http"
	"sync"

	"gopkg.in/yaml.v3"
)

// reconfReq is a requiest sent to the app
// to reconfigure
type reconfReq struct {
	conf *Config
	resp chan *reconfResp
}

// reconfResp is a response from the app
// whether reconfiguration is possible
type reconfResp struct {
	code int
	body []byte
}

type magicHandler struct {
	conf  *Config
	magic chan *reconfReq
	mu    sync.Mutex
}

func newMagicHandler(ch chan *reconfReq, c *Config) *magicHandler {
	return &magicHandler{magic: ch, conf: c}
}

func (m *magicHandler) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.handleShowConfig(w, r)
	case http.MethodPost:
		m.handleReconfigure(w, r)
	default:
		http.Error(w, "use GET to receive current configuration and POST to push a new configuration to a running server", http.StatusMethodNotAllowed)
	}
}

func (m *magicHandler) handleReconfigure(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conf, err := parseConfig(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ok := m.mu.TryLock(); !ok {
		http.Error(w, "server is already being reconfigured", http.StatusConflict)
		return
	}
	defer m.mu.Unlock()

	log.Printf("handling reconfigure request from %s", r.RemoteAddr)
	log.Printf("new config is as follows:\n%s", string(b))

	respChan := make(chan *reconfResp)
	m.magic <- &reconfReq{
		conf: conf,
		resp: respChan,
	}

	resp := <-respChan
	w.WriteHeader(resp.code)
	if _, err := w.Write(resp.body); err != nil {
		log.Println(err)
	}

	if resp.code == http.StatusAccepted {
		// this tells the swapper thread that we're done
		// sending response and the server may be shut down
		close(m.magic)
	}
}

func (m *magicHandler) handleShowConfig(w http.ResponseWriter, _ *http.Request) {
	if m.conf == nil {
		http.Error(w, "configuration is not available", http.StatusInternalServerError)
	}
	b, err := yaml.Marshal(m.conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}
