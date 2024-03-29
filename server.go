package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/jsign/verkle-explorer/database"
	"github.com/jsign/verkle-explorer/handlers"
)

type WitnessExplorer struct {
	lock    sync.Mutex
	started bool
	server  http.Server
}

func New(addr string, db database.DB) *WitnessExplorer {
	mux := http.NewServeMux()
	configureStaticFiles(mux)
	mux.HandleFunc("/tx", handlers.HandlerGetTx(db))
	mux.HandleFunc("/", handlers.HandlerGetTx(db))

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
	css := http.FileServer(http.Dir("webtemplate/css/"))
	mux.Handle("/css/", http.StripPrefix("/css/", css))

	img := http.FileServer(http.Dir("webtemplate/img/"))
	mux.Handle("/img/", http.StripPrefix("/img/", img))

	js := http.FileServer(http.Dir("webtemplate/js/"))
	mux.Handle("/js/", http.StripPrefix("/js/", js))

	scss := http.FileServer(http.Dir("webtemplate/scss/"))
	mux.Handle("/scss/", http.StripPrefix("/scss/", scss))

	vendor := http.FileServer(http.Dir("webtemplate/vendor/"))
	mux.Handle("/vendor/", http.StripPrefix("/vendor/", vendor))

}
