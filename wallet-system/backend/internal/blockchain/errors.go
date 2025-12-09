package blockchain

import "errors"

var (
	ErrInvalidPreviousHash = errors.New("invalid previous hash")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidTransaction = errors.New("invalid transaction")
	ErrUTXOAlreadySpent = errors.New("UTXO already spent")
	ErrInvalidSignature = errors.New("invalid digital signature")
	ErrInvalidWallet = errors.New("invalid wallet address")
)
