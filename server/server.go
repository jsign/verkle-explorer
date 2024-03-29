package server

import (
	"fmt"
	"net/http"
	"sync"
	"text/template"

	"github.com/jsign/verkle-explorer/database"
	"github.com/jsign/verkle-explorer/handlers"
)

const prefix = "/home/ignacio/code/verkle-explorer/"

type WitnessExplorer struct {
	lock    sync.Mutex
	started bool
	server  http.Server
}

func New(addr string, db database.DB) *WitnessExplorer {
	mux := http.NewServeMux()
	configureStaticFiles(mux)

	var tmpl = template.Must(template.ParseFiles(prefix + "webtemplate/tx.html"))
	mux.HandleFunc("/tx", handlers.HandlerGetTx(tmpl, db))
	mux.HandleFunc("/", handlers.HandlerGetTx(tmpl, db))

	return &WitnessExplorer{server: http.Server{Addr: addr, Handler: mux}}
}

func (s *WitnessExplorer) Run() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.started {
		return fmt.Errorf("server already started")
	}
	s.started = true
	s.lock.Unlock()

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *WitnessExplorer) Close() error {
	return s.server.Close()
}

func configureStaticFiles(mux *http.ServeMux) {
	paths := []string{"css", "img", "js", "scss", "vendor"}

	for _, path := range paths {
		incl := http.FileServer(http.Dir(prefix + "webtemplate/" + path))
		mux.Handle("/"+path+"/", http.StripPrefix("/"+path+"/", incl))

	}
}
