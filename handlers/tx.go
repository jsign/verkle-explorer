package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/jsign/verkle-explorer/database"
)

type txContext struct {
	Hash   string
	Exists bool

	TotalGas               uint64
	ExecutionGas           uint64
	ExecutionGasPercentage int
	CodeChunkGas           uint64
	CodeChunkGasPercentage int

	ExecutedInstructions int
	ExecutedBytes        uint64
	ChargedBytes         int
	ExecutionEfficiency  string

	WitnessEvents []database.WitnessEvent
}

func HandlerGetTx(tmpl *template.Template, db database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		txHash := r.URL.Query().Get("hash")

		if txHash == "" {
			if err := tmpl.Execute(w, nil); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		txCtx := txContext{Hash: txHash}
		txExec, err := db.GetTxExec(txCtx.Hash)
		if err != nil && err != database.ErrTxNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err == nil {
			txCtx.Exists = true

			txCtx.TotalGas = txExec.TotalGas
			txCtx.ExecutionGas = txExec.TotalGas - txExec.CodeChunkGas
			txCtx.ExecutionGasPercentage = int(txCtx.ExecutionGas * 100 / txCtx.TotalGas)
			txCtx.CodeChunkGas = txExec.CodeChunkGas
			txCtx.CodeChunkGasPercentage = 100 - txCtx.ExecutionGasPercentage

			txCtx.ExecutedInstructions = txExec.ExecutedInstructions
			txCtx.ExecutedBytes = txExec.ExecutedBytes
			txCtx.ChargedBytes = txExec.ChargedBytes
			txCtx.ExecutionEfficiency = "N/A"
			if txCtx.ChargedBytes > 0 {
				txCtx.ExecutionEfficiency = fmt.Sprintf("%0.02f", float64(txCtx.ExecutedBytes)/float64(txCtx.ChargedBytes))
			}

			txCtx.WitnessEvents = txExec.Events
		}
		if err := tmpl.Execute(w, txCtx); err != nil {
			log.Printf("failed to execute template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
