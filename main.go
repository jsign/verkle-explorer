package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	configureStaticFiles(mux)

	mux.HandleFunc("/tx", handleLulz)

	server := http.Server{
		Addr:    ":8181",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func handleLulz(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("webtemplate/tx.html"))
	type witnessEvent struct {
		Name string
		Gas  int
	}
	type txContext struct {
		ID string

		TotalGas               int
		ExecutionGas           int
		ExecutionGasPercentage int
		CodeChunkGas           int
		CodeChunkGasPercentage int

		TouchedCodeChunks int
		ChargedCodeChunks int
		ChunkEfficiency   string

		WitnessEvents []witnessEvent
	}
	data := txContext{
		ID: "0xc6a03a1cc3678f83d5a62de0bf8ca6f1fc0ee46aea90f0ea7a6c56890e0e0613",

		TotalGas:     4300,
		ExecutionGas: 3150,

		TouchedCodeChunks: 120,
		ChargedCodeChunks: 34,

		WitnessEvents: []witnessEvent{
			{Name: "ContractInit", Gas: 100},
			{Name: "TouchFullAddress", Gas: 200},
			{Name: "TouchAddressOnWrite", Gas: 500},
			{Name: "ContractInitiCompletion", Gas: 800},
		},
	}
	data.CodeChunkGas = data.TotalGas - data.ExecutionGas
	data.ExecutionGasPercentage = data.ExecutionGas * 100 / data.TotalGas
	data.CodeChunkGasPercentage = 100 - data.ExecutionGasPercentage

	data.ChunkEfficiency = fmt.Sprintf("%0.02f", float64(data.TouchedCodeChunks)/float64(data.ChargedCodeChunks))
	if err := tmpl.Execute(w, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
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
