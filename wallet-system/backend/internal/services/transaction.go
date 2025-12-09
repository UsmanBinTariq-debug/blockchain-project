package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"crypto-wallet-backend/internal/database"
)

// TransactionService handles transaction operations
type TransactionService struct {
	db *database.Database
}

// NewTransactionService creates a new transaction service
func NewTransactionService(db *database.Database) *TransactionService {
	return &TransactionService{db: db}
}

// CreateTransaction creates a new transaction
func (ts *TransactionService) CreateTransaction(ctx context.Context, tx database.Transaction) (string, error) {
	// Validate transaction
	if tx.SenderWallet == "" || tx.ReceiverWallet == "" {
		return "", fmt.Errorf("invalid transaction: missing sender or receiver wallet")
	}
	if tx.Amount <= 0 {
		return "", fmt.Errorf("invalid transaction: amount must be positive")
	}

	// Generate transaction hash (SHA256, exactly 64 chars)
	h := sha256.Sum256([]byte(fmt.Sprintf("%s%s%d%.8f", tx.SenderWallet, tx.ReceiverWallet, time.Now().UnixNano(), tx.Amount)))
	txHash := hex.EncodeToString(h[:])

	// Build DB transaction record
	dbTx := &database.Transaction{
		TransactionHash: txHash,
		SenderWallet:    tx.SenderWallet,
		ReceiverWallet:  tx.ReceiverWallet,
		Amount:          tx.Amount,
		Fee:             tx.Fee,
		Note:            tx.Note,
		Signature:       tx.Signature,
		Status:          "confirmed",
		TransactionType: "transfer",
		CreatedAt:       time.Now(),
	}

	// Persist transaction record
	if err := ts.db.CreateTransaction(ctx, dbTx); err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	// Select sender UTXOs to cover amount + fee
	required := tx.Amount + tx.Fee
	utxos, err := ts.db.GetUTXOsByWallet(ctx, tx.SenderWallet)
	if err != nil {
		return "", fmt.Errorf("failed to fetch sender utxos: %w", err)
	}

	total := 0.0
	var used []*database.UTXO
	for _, u := range utxos {
		used = append(used, u)
		total += u.Amount
		if total >= required {
			break
		}
	}

	// If UTXOs are insufficient, try to fallback to the wallet's cached balance
	if total < required {
		// Diagnostic log: no or insufficient UTXOs
		fmt.Printf("[tx] insufficient utxos: found %d utxos, total %.8f required %.8f for wallet %s\n", len(utxos), total, required, tx.SenderWallet)

		wallet, werr := ts.db.GetWalletByAddress(ctx, tx.SenderWallet)
		if werr != nil {
			return "", fmt.Errorf("failed to fetch wallet for fallback: %w", werr)
		}
		if wallet != nil {
			fmt.Printf("[tx] wallet.BalanceCache for %s = %.8f\n", wallet.WalletAddress, wallet.BalanceCache)
		} else {
			fmt.Printf("[tx] wallet not found for address %s\n", tx.SenderWallet)
		}
		if wallet != nil && wallet.BalanceCache >= required {
			// Create a synthetic UTXO representing the cached balance
			// Use a proper SHA256 hash (exactly 64 chars) for the seed transaction
			seedHashBytes := sha256.Sum256([]byte(fmt.Sprintf("seed_%s_%d", wallet.WalletAddress, time.Now().UnixNano())))
			seedHash := hex.EncodeToString(seedHashBytes[:])
			seedUTXO := &database.UTXO{
				TransactionHash: seedHash,
				OutputIndex:     0,
				WalletAddress:   wallet.WalletAddress,
				Amount:          wallet.BalanceCache,
				IsSpent:         false,
				CreatedAt:       time.Now(),
			}
			if err := ts.db.CreateUTXO(ctx, seedUTXO); err != nil {
				return "", fmt.Errorf("failed to create seed utxo for fallback: %w", err)
			}
			used = append(used, seedUTXO)
			total += seedUTXO.Amount
		}
	}

	if total < required {
		return "", fmt.Errorf("insufficient funds: have %.8f required %.8f", total, required)
	}

	// Mark used UTXOs as spent
	for _, u := range used {
		if err := ts.db.MarkUTXOAsSpent(ctx, u.TransactionHash, u.OutputIndex, txHash); err != nil {
			return "", fmt.Errorf("failed to mark utxo spent: %w", err)
		}
	}

	// Create output UTXO to receiver (output index 0)
	out := &database.UTXO{
		TransactionHash: txHash,
		OutputIndex:     0,
		WalletAddress:   tx.ReceiverWallet,
		Amount:          tx.Amount,
		IsSpent:         false,
		CreatedAt:       time.Now(),
	}
	if err := ts.db.CreateUTXO(ctx, out); err != nil {
		return "", fmt.Errorf("failed to create receiver utxo: %w", err)
	}

	// Create change UTXO back to sender if any
	change := total - required
	if change > 0 {
		changeUTXO := &database.UTXO{
			TransactionHash: txHash,
			OutputIndex:     1,
			WalletAddress:   tx.SenderWallet,
			Amount:          change,
			IsSpent:         false,
			CreatedAt:       time.Now(),
		}
		if err := ts.db.CreateUTXO(ctx, changeUTXO); err != nil {
			return "", fmt.Errorf("failed to create change utxo: %w", err)
		}
	}

	// Update wallet cached balances (recalculate from UTXOs)
	// Receiver
	recvUtxos, _ := ts.db.GetUTXOsByWallet(ctx, tx.ReceiverWallet)
	recvBal := 0.0
	for _, u := range recvUtxos {
		if !u.IsSpent {
			recvBal += u.Amount
		}
	}
	fmt.Printf("[tx-balance] Receiver %s: %d utxos, calculated balance = %.8f\n", tx.ReceiverWallet, len(recvUtxos), recvBal)
	_ = ts.db.UpdateWalletBalance(ctx, tx.ReceiverWallet, recvBal)

	// Sender
	senderUtxos, _ := ts.db.GetUTXOsByWallet(ctx, tx.SenderWallet)
	senderBal := 0.0
	for _, u := range senderUtxos {
		if !u.IsSpent {
			senderBal += u.Amount
		}
		fmt.Printf("[tx-balance] Sender UTXO: hash=%s amount=%.8f is_spent=%v\n", u.TransactionHash, u.Amount, u.IsSpent)
	}
	fmt.Printf("[tx-balance] Sender %s: %d utxos, calculated balance = %.8f\n", tx.SenderWallet, len(senderUtxos), senderBal)
	_ = ts.db.UpdateWalletBalance(ctx, tx.SenderWallet, senderBal)

	// Log system event
	_ = ts.db.CreateSystemLog(ctx, &database.SystemLog{
		LogType:       "transaction",
		Message:       fmt.Sprintf("Transaction %s: %s -> %s amount %.8f fee %.8f", txHash, tx.SenderWallet, tx.ReceiverWallet, tx.Amount, tx.Fee),
		WalletAddress: tx.SenderWallet,
		CreatedAt:     time.Now(),
	})

	return txHash, nil
}

// GetTransactionHistory retrieves transaction history for a wallet
func (ts *TransactionService) GetTransactionHistory(ctx context.Context, walletAddress string, limit, offset int) ([]*database.Transaction, error) {
	return ts.db.GetTransactionsByWallet(ctx, walletAddress, limit, offset)
}

// GetTransactionByHash retrieves a transaction by hash
func (ts *TransactionService) GetTransactionByHash(ctx context.Context, hash string) (*database.Transaction, error) {
	return ts.db.GetTransactionByHash(ctx, hash)
}

// UpdateTransactionStatus updates transaction status
func (ts *TransactionService) UpdateTransactionStatus(ctx context.Context, txHash, status string) error {
	// This would require adding an update method to the database service
	return nil
}
