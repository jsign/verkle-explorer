package database

import "errors"

var ErrTxNotFound = errors.New("tx doesn't exist")

type DB interface {
	GetTxExec(hash string) (TxExec, error)
}

type TxExec struct {
	Hash string

	TotalGas     int
	CodeChunkGas int

	ExecutedInstructions int
	ExecutedBytes        int
	ChargedBytes         int

	Events []WitnessEvent
}

type WitnessEvent struct {
	Name   string
	Gas    int
	Params string
}
