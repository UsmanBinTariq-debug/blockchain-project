package services

import (
	"context"
	"crypto-wallet-backend/internal/blockchain"
	"crypto-wallet-backend/internal/crypto"
	"crypto-wallet-backend/internal/database"
	"fmt"
	"time"
)

// WalletService handles wallet operations
type WalletService struct {
	db *database.Database
	bc *blockchain.Blockchain
}

// NewWalletService creates a new wallet service
func NewWalletService(db *database.Database, bc *blockchain.Blockchain) *WalletService {
	return &WalletService{db: db, bc: bc}
}

// CreateWallet creates a new wallet for a user
func (ws *WalletService) CreateWallet(ctx context.Context, userID, email string) (*database.Wallet, string, error) {
	// Generate key pair
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, "", err
	}

	// Create wallet in database with default balance cache
	wallet := &database.Wallet{
		UserID:        userID,
		WalletAddress: keyPair.WalletID,
		BalanceCache:  200.0,
	}

	if err := ws.db.CreateWallet(ctx, wallet); err != nil {
		return nil, "", err
	}

	// Create initial UTXO so the balance is spendable
	initTxHash := fmt.Sprintf("init-%s-%d", keyPair.WalletID, time.Now().Unix())
	initUTXO := &database.UTXO{
		TransactionHash: initTxHash,
		OutputIndex:     0,
		WalletAddress:   keyPair.WalletID,
		Amount:          200.0,
		IsSpent:         false,
	}

	if err := ws.db.CreateUTXO(ctx, initUTXO); err == nil {
		// also add to in-memory blockchain UTXO set
		ws.bc.AddUTXO(keyPair.WalletID, blockchain.UTXO{
			TransactionHash: initUTXO.TransactionHash,
			OutputIndex:     initUTXO.OutputIndex,
			WalletAddress:   initUTXO.WalletAddress,
			Amount:          initUTXO.Amount,
			IsSpent:         initUTXO.IsSpent,
			SpentInTx:       "",
		})
	}

	// Store keys (private key should be encrypted)
	user, err := ws.db.GetUserByWalletID(ctx, keyPair.WalletID)
	if err != nil || user == nil {
		// Update user with wallet info (already done during user creation)
	}

	return wallet, keyPair.PublicKey, nil
}

// GetWalletBalance calculates the balance from UTXOs
func (ws *WalletService) GetWalletBalance(ctx context.Context, walletAddress string) (float64, error) {
	utxos, err := ws.db.GetUTXOsByWallet(ctx, walletAddress)
	if err != nil {
		return 0, err
	}

	var balance float64
	for _, utxo := range utxos {
		balance += utxo.Amount
	}

	// Update cached balance
	_ = ws.db.UpdateWalletBalance(ctx, walletAddress, balance)

	return balance, nil
}

// ValidateTransaction validates a transaction before adding it to the blockchain
func (ws *WalletService) ValidateTransaction(ctx context.Context, tx *blockchain.Transaction) error {
	// Check wallet exists
	sender, err := ws.db.GetUserByWalletID(ctx, tx.SenderWallet)
	if err != nil || sender == nil {
		return fmt.Errorf("sender wallet not found")
	}

	receiver, err := ws.db.GetUserByWalletID(ctx, tx.ReceiverWallet)
	if err != nil || receiver == nil {
		return fmt.Errorf("receiver wallet not found")
	}

	// Verify signature
	dataToSign := tx.SenderWallet + tx.ReceiverWallet + fmt.Sprintf("%.8f", tx.Amount)
	isValid, err := crypto.VerifySignature(dataToSign, tx.Signature, tx.PublicKey)
	if err != nil || !isValid {
		return fmt.Errorf("invalid digital signature")
	}

	// Check balance
	balance, err := ws.GetWalletBalance(ctx, tx.SenderWallet)
	if err != nil {
		return err
	}

	if balance < (tx.Amount + tx.Fee) {
		return fmt.Errorf("insufficient balance")
	}

	// Check UTXOs
	utxos, err := ws.db.GetUTXOsByWallet(ctx, tx.SenderWallet)
	if err != nil {
		return err
	}

	if len(utxos) == 0 {
		return fmt.Errorf("no unspent outputs")
	}

	return nil
}

// CreateTransaction creates a new transaction
func (ws *WalletService) CreateTransaction(ctx context.Context, senderWallet, receiverWallet string, amount, fee float64, note, signature, publicKey string) (*blockchain.Transaction, error) {
	tx := blockchain.NewTransaction(senderWallet, receiverWallet, amount, fee, note)
	tx.Signature = signature
	tx.PublicKey = publicKey

	// Validate transaction
	if err := ws.ValidateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	// Save to database
	dbTx := &database.Transaction{
		TransactionHash: tx.ID,
		SenderWallet:    tx.SenderWallet,
		ReceiverWallet:  tx.ReceiverWallet,
		Amount:          tx.Amount,
		Fee:             tx.Fee,
		Note:            tx.Note,
		Signature:       tx.Signature,
		Status:          "pending",
		TransactionType: "transfer",
	}

	if err := ws.db.CreateTransaction(ctx, dbTx); err != nil {
		return nil, err
	}

	return tx, nil
}
