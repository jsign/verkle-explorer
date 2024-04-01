package database

import (
	"errors"
)

var ErrTxNotFound = errors.New("tx doesn't exist")

type DB interface {
	GetTxExec(hash string) (TxExec, error)
}

type TxExec struct {
	Hash string

	BlockNumber uint64
	From        string
	To          string
	Value       string

	TotalGas     uint64
	CodeChunkGas uint64

	ExecutedInstructions int
	ExecutedBytes        uint64
	ChargedBytes         int

	Events []WitnessEvent
}

type WitnessEvent struct {
	Name   string
	Gas    uint64
	Params string
}
