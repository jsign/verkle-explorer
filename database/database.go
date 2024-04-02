package database

import (
	"errors"
)

var ErrTxNotFound = errors.New("tx doesn't exist")

type DB interface {
	GetTxExec(hash string) (TxExec, error)
	GetHighestGasTxs(count int) ([]TxInfo, error)
	GetInefficientCodeAccessTxs(count int) ([]TxInfo, error)
}

type TxExec struct {
	TxInfo

	Events               []WitnessEvent
	WitnessTreeKeyValues []WitnessTreeKeyValue
}

type TxInfo struct {
	Hash string

	BlockNumber uint64
	From        string
	To          string
	Value       string

	TotalGas          uint64
	CodeChunkGas      uint64
	ChargedCodeChunks int

	ExecutedInstructions int
	ExecutedBytes        uint64
}

type WitnessEvent struct {
	Name   string
	Gas    uint64
	Params string
}

type WitnessTreeKeyValue struct {
	Key       string
	PostValue string
}
