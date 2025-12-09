package database

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                 string    `json:"id"`
	Email              string    `json:"email"`
	FullName           string    `json:"full_name"`
	CNIC               string    `json:"cnic"`
	WalletID           string    `json:"wallet_id"`
	PublicKey          string    `json:"public_key"`
	EncryptedPrivateKey string   `json:"encrypted_private_key"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	IsVerified         bool      `json:"is_verified"`
}

// Wallet represents a user's wallet
type Wallet struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	WalletAddress    string    `json:"wallet_address"`
	BalanceCache     float64   `json:"balance_cache"`
	LastUpdated      time.Time `json:"last_updated"`
	ZakatDeducted    bool      `json:"zakat_deducted_this_month"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	ID              string     `json:"id"`
	TransactionHash string     `json:"transaction_hash"`
	BlockHash       *string    `json:"block_hash,omitempty"`
	SenderWallet    string     `json:"sender_wallet"`
	ReceiverWallet  string     `json:"receiver_wallet"`
	Amount          float64    `json:"amount"`
	Fee             float64    `json:"fee"`
	Note            string     `json:"note,omitempty"`
	Signature       string     `json:"signature"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	TransactionType string     `json:"transaction_type"`
}

// Block represents a blockchain block
type Block struct {
	ID            string    `json:"id"`
	BlockIndex    int64     `json:"block_index"`
	Timestamp     int64     `json:"timestamp"`
	PreviousHash  string    `json:"previous_hash"`
	Hash          string    `json:"hash"`
	Nonce         int64     `json:"nonce"`
	MerkleRoot    string    `json:"merkle_root"`
	Difficulty    int       `json:"difficulty"`
	MinedBy       string    `json:"mined_by,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// UTXO represents an unspent transaction output
type UTXO struct {
	ID                 string    `json:"id"`
	TransactionHash    string    `json:"transaction_hash"`
	OutputIndex        int       `json:"output_index"`
	WalletAddress      string    `json:"wallet_address"`
	Amount             float64   `json:"amount"`
	IsSpent            bool      `json:"is_spent"`
	SpentInTransaction *string   `json:"spent_in_transaction,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

// ZakatTransaction represents a zakat deduction
type ZakatTransaction struct {
	ID              string    `json:"id"`
	WalletAddress   string    `json:"wallet_address"`
	Amount          float64   `json:"amount"`
	ZakatPercentage float64   `json:"zakat_percentage"`
	TransactionHash string    `json:"transaction_hash"`
	MonthYear       string    `json:"month_year"`
	CreatedAt       time.Time `json:"created_at"`
}

// SystemLog represents a system log entry
type SystemLog struct {
	ID            string    `json:"id"`
	LogType       string    `json:"log_type"`
	Message       string    `json:"message"`
	WalletAddress string    `json:"wallet_address,omitempty"`
	IPAddress     string    `json:"ip_address,omitempty"`
	UserAgent     string    `json:"user_agent,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// Beneficiary represents a saved beneficiary wallet
type Beneficiary struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	BeneficiaryWalletID string   `json:"beneficiary_wallet_id"`
	Nickname           string    `json:"nickname"`
	CreatedAt          time.Time `json:"created_at"`
}
