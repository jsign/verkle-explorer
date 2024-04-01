package handlers

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"text/template"

	"github.com/jsign/verkle-explorer/database"
	"github.com/shopspring/decimal"
)

type txContext struct {
	Hash   string
	Exists bool

	BlockNumber uint64
	From        string
	To          string
	Value       string

	TotalGas               uint64
	ExecutionGas           uint64
	ExecutionGasPercentage int
	CodeChunkGas           uint64
	CodeChunkGasPercentage int

	ExecutedInstructions int
	ExecutedBytes        uint64
	ChargedBytes         int
	ExecutionEfficiency  string

	WitnessEvents        []database.WitnessEvent
	WitnessTreeKeyValues []database.WitnessTreeKeyValue
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

			txCtx.BlockNumber = txExec.BlockNumber
			txCtx.From = txExec.From
			txCtx.To = txExec.To
			var biValue big.Int
			biValue.SetString(txExec.Value, 10)
			value := decimal.NewFromBigInt(&biValue, -18)
			txCtx.Value = value.String()

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
				txCtx.ExecutionEfficiency = fmt.Sprintf("%0.02fx", float64(txCtx.ExecutedBytes)/float64(txCtx.ChargedBytes))
			}

			txCtx.WitnessEvents = txExec.Events

			txCtx.WitnessTreeKeyValues = txExec.WitnessTreeKeyValues
		}
		if err := tmpl.Execute(w, txCtx); err != nil {
			log.Printf("failed to execute template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
