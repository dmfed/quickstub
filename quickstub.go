package quickstub

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

// QuickStub is the main app interface
// under the hood it initializes http.Server accepting
// requests on endpoints provided in the Config
// and sets up the "magic endpoint" for on the fly
// reconfiguration of the server.
type QuickStub interface {
	ListenAndServe() error
	ListenAndServeTLS(certFile string, keyFile string) error
	Shutdown(ctx context.Context) error
}

func NewQuickStub(conf *Config) (QuickStub, error) {
	return &quickStub{
		config: conf,
	}, nil
}

type quickStub struct {
	config   *Config
	server   *http.Server
	pkey     string
	cert     string
	magic    chan *reconfReq
	restart  chan struct{}
	shutdown chan struct{}
	mu       sync.Mutex
}

func (q *quickStub) ListenAndServe() error {
	return q.listenAndServe("", "")
}

func (q *quickStub) ListenAndServeTLS(certFile, keyFile string) error {
	return q.listenAndServe(certFile, keyFile)
}

func (q *quickStub) listenAndServe(certFile, keyFile string) error {
	srv, ch, err := makeNewServer(q.config)
	if err != nil {
		return err
	}
	q.cert = certFile
	q.pkey = keyFile
	q.server = srv
	q.magic = ch
	q.restart = make(chan struct{}, 1)
	q.shutdown = make(chan struct{})

	go q.runServerSwapper()
	for {
		var exitError error
		if q.cert != "" && q.pkey != "" {
			exitError = q.server.ListenAndServeTLS(q.cert, q.pkey)
		} else {
			exitError = q.server.ListenAndServe()
		}

		if !errors.Is(exitError, http.ErrServerClosed) {
			// if error is other than normal ErrServerClosed
			// there might be a problem
			return exitError
		}
		// waiting for signal what to do next
		select {
		case <-q.restart:
			// new server has been created
			// release from select and continue
			// looping
		case <-q.shutdown:
			// Shutdown has been called
			// returning the final error here
			return exitError
		}
	}
}

func (q *quickStub) runServerSwapper() {
	for {
		// wait for request from magic endpoint
		// or for shutdown signal
		var req *reconfReq
		select {
		case req = <-q.magic:
			// we've got the request
		case <-q.shutdown:
			// nothing else to do
			return
		}

		srv, ch, err := makeNewServer(req.conf)
		if err != nil {
			req.resp <- &reconfResp{
				code: http.StatusBadRequest,
				body: []byte(err.Error()),
			}
			close(req.resp)
			continue
		}

		req.resp <- &reconfResp{
			code: http.StatusAccepted,
			body: []byte("OK"),
		}
		close(req.resp)

		// waiting for channel to close
		// this will indicate that the reconfigure handler
		// returned
		<-q.magic

		// replacing the server
		q.mu.Lock() // locking to replace the server instance
		shutdownCtx, stop := context.WithTimeout(context.Background(), time.Second*20)
		q.server.Shutdown(shutdownCtx)
		stop()

		q.config = req.conf
		q.server = srv
		q.magic = ch
		q.mu.Unlock() // unlock

		// when we're done creating server it may be launched by
		// main thread
		q.restart <- struct{}{}

	}
}

func (q *quickStub) Shutdown(ctx context.Context) error {
	q.mu.Lock() // locking to prevent possible operations with server in swapper
	defer q.mu.Unlock()

	close(q.shutdown) // this will tell main thread and swapper that we're done
	return q.server.Shutdown(ctx)
}
