package server

import (
	"embed"
	"fmt"
	"net/http"
	"sync"
	"text/template"

	"github.com/jsign/verkle-explorer/database"
	"github.com/jsign/verkle-explorer/handlers"
)

//go:embed pages
var pagesFS embed.FS

//go:embed static
var staticFS embed.FS

type WitnessExplorer struct {
	lock    sync.Mutex
	started bool
	server  http.Server
}

func New(addr string, db database.DB) *WitnessExplorer {
	mux := http.NewServeMux()

	var tmpl = template.Must(template.ParseFS(pagesFS, "pages/tx.html"))
	mux.HandleFunc("/tx", handlers.HandlerGetTx(tmpl, db))
	mux.Handle("/static/", http.FileServerFS(staticFS))
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
