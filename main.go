package main

import (
	"log"
	"net/http"

	"github.com/jsign/verkle-explorer/database"
	"github.com/jsign/verkle-explorer/database/mock"
	"github.com/jsign/verkle-explorer/handlers"
)

func main() {
	db := mock.NewMockDB([]database.TxExec{
		{
			Hash:                 "0xc6a03a1cc3678f83d5a62de0bf8ca6f1fc0ee46aea90f0ea7a6c56890e0e0613",
			TotalGas:             4300,
			CodeChunkGas:         1200,
			ExecutedInstructions: 120,
			ExecutedBytes:        125,
			ChargedBytes:         62,
			Events: []database.WitnessEvent{
				{Name: "ContractInit", Gas: 100},
				{Name: "TouchFullAddress", Gas: 200},
				{Name: "TouchAddressOnWrite", Gas: 500},
				{Name: "ContractInitiCompletion", Gas: 800},
			},
		},
	})

	mux := http.NewServeMux()
	configureStaticFiles(mux)
	mux.HandleFunc("/tx/{hash}", handlers.HandlerGetTx(db))

	server := http.Server{Addr: ":8181", Handler: mux}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
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
