package handlers

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"text/template"

	"github.com/jsign/verkle-explorer/database"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

type txContext struct {
	IsDashboard bool

	Hash   string
	Exists bool

	BlockNumber uint64
	From        string
	To          string
	Value       string

	TotalGas                    uint64
	ExecutionGas                uint64
	ExecutionGasPercentage      int
	NonCodeWitnessGas           uint64
	NonCodeWitnessGasPercentage int
	CodeChunkGas                uint64
	CodeChunkGasPercentage      int

	ExecutedInstructions int
	ExecutedBytes        uint64
	ChargedBytes         int
	ExecutionEfficiency  string

	WitnessEvents        []database.WitnessEvent
	WitnessTreeKeyValues []database.WitnessTreeKeyValue
	WitnessCharges       []database.WitnessCharges
}

type dashboardContext struct {
	IsDashboard              bool
	TopHighestGas            []txListItem
	TopInefficientCodeAccess []txListItem
}

type txListItem struct {
	Hash   string
	Detail string
}

func HandlerGetTx(tmpl *template.Template, db database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		txHash := r.URL.Query().Get("hash")

		if txHash == "" {
			outputDashboard(db, tmpl, w)
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

			txCtx.ExecutedInstructions = txExec.ExecutedInstructions
			txCtx.ExecutedBytes = txExec.ExecutedBytes
			txCtx.ChargedBytes = txExec.ChargedCodeChunks * 31
			txCtx.ExecutionEfficiency = "N/A"
			if txCtx.ChargedBytes > 0 {
				txCtx.ExecutionEfficiency = fmt.Sprintf("%0.02fx", float64(txCtx.ExecutedBytes)/float64(txCtx.ChargedBytes))
			}

			txCtx.WitnessEvents = txExec.WitnessEvents
			txCtx.WitnessTreeKeyValues = txExec.WitnessTreeKeyValues
			txCtx.WitnessCharges = txExec.WitnessCharges

			for _, c := range txCtx.WitnessCharges {
				txCtx.NonCodeWitnessGas += c.Gas

			}

			txCtx.TotalGas = txExec.TotalGas
			txCtx.ExecutionGas = txExec.TotalGas - txCtx.NonCodeWitnessGas - txExec.CodeChunkGas
			txCtx.ExecutionGasPercentage = int(txCtx.ExecutionGas * 100 / txCtx.TotalGas)
			txCtx.CodeChunkGas = txExec.CodeChunkGas
			txCtx.CodeChunkGasPercentage = int(txCtx.CodeChunkGas * 100 / txCtx.TotalGas)
			txCtx.NonCodeWitnessGasPercentage = 100 - txCtx.ExecutionGasPercentage - txCtx.CodeChunkGasPercentage

		}
		if err := tmpl.Execute(w, txCtx); err != nil {
			log.Printf("failed to execute template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func outputDashboard(db database.DB, tmpl *template.Template, w http.ResponseWriter) {
	group, _ := errgroup.WithContext(context.Background())

	var highGasTxs []database.TxInfo
	var inefficientTxs []database.TxInfo
	group.Go(func() error {
		var err error
		highGasTxs, err = db.GetHighestGasTxs(10)
		if err != nil {
			return err
		}
		return nil
	})
	group.Go(func() error {
		var err error
		inefficientTxs, err = db.GetInefficientCodeAccessTxs(10)
		if err != nil {
			return err
		}
		return nil
	})
	if err := group.Wait(); err != nil {
		log.Printf("failed to retrieve dashboard info: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dashboardContext := dashboardContext{
		IsDashboard:              true,
		TopHighestGas:            make([]txListItem, len(highGasTxs)),
		TopInefficientCodeAccess: make([]txListItem, len(inefficientTxs)),
	}
	for i, hgt := range highGasTxs {
		dashboardContext.TopHighestGas[i] = txListItem{Hash: hgt.Hash, Detail: fmt.Sprintf("%d", hgt.TotalGas)}
	}
	for i, hgt := range inefficientTxs {
		dashboardContext.TopInefficientCodeAccess[i] = txListItem{Hash: hgt.Hash, Detail: fmt.Sprintf("%0.02fx", float64(hgt.ExecutedBytes)/float64(hgt.ChargedCodeChunks*31))}
	}
	if err := tmpl.Execute(w, dashboardContext); err != nil {
		log.Printf("failed to execute template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
