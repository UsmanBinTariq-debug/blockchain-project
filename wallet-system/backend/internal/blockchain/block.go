package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Block represents a single block in the blockchain
type Block struct {
	Index        int64          `json:"index"`
	Timestamp    int64          `json:"timestamp"`
	Transactions []Transaction  `json:"transactions"`
	PreviousHash string         `json:"previous_hash"`
	Nonce        int64          `json:"nonce"`
	Hash         string         `json:"hash"`
	MerkleRoot   string         `json:"merkle_root"`
	Difficulty   int            `json:"difficulty"`
	MinedBy      string         `json:"mined_by,omitempty"`
}

// Transaction represents a transaction in the blockchain
type Transaction struct {
	ID             string        `json:"id"`
	SenderWallet   string        `json:"sender_wallet"`
	ReceiverWallet string        `json:"receiver_wallet"`
	Amount         float64       `json:"amount"`
	Fee            float64       `json:"fee"`
	Note           string        `json:"note,omitempty"`
	Timestamp      int64         `json:"timestamp"`
	Signature      string        `json:"signature"`
	PublicKey      string        `json:"public_key"`
	UTXOInputs     []UTXO        `json:"utxo_inputs"`
	UTXOOutputs    []UTXO        `json:"utxo_outputs"`
	Status         string        `json:"status"` // pending, confirmed, failed
}

// UTXO represents an Unspent Transaction Output
type UTXO struct {
	TransactionHash string  `json:"transaction_hash"`
	OutputIndex     int     `json:"output_index"`
	WalletAddress   string  `json:"wallet_address"`
	Amount          float64 `json:"amount"`
	IsSpent         bool    `json:"is_spent"`
	SpentInTx       string  `json:"spent_in_transaction,omitempty"`
}

// NewBlock creates a new block
func NewBlock(index int64, transactions []Transaction, previousHash string, difficulty int) *Block {
	return &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PreviousHash: previousHash,
		Difficulty:   difficulty,
		MerkleRoot:   calculateMerkleRoot(transactions),
	}
}

// CalculateHash calculates the SHA256 hash of the block
func (b *Block) CalculateHash() string {
	blockData := fmt.Sprintf("%d%d%s%d%s", b.Index, b.Timestamp, b.PreviousHash, b.Nonce, b.MerkleRoot)
	hash := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hash[:])
}

// NewTransaction creates a new transaction
func NewTransaction(senderWallet, receiverWallet string, amount, fee float64, note string) *Transaction {
	return &Transaction{
		SenderWallet:   senderWallet,
		ReceiverWallet: receiverWallet,
		Amount:         amount,
		Fee:            fee,
		Note:           note,
		Timestamp:      time.Now().Unix(),
		Status:         "pending",
	}
}

// CalculateMerkleRoot calculates the Merkle root of transactions
func calculateMerkleRoot(transactions []Transaction) string {
	if len(transactions) == 0 {
		return hashData("")
	}

	if len(transactions) == 1 {
		return hashTransaction(transactions[0])
	}

	var tree []string
	for _, tx := range transactions {
		tree = append(tree, hashTransaction(tx))
	}

	for len(tree) > 1 {
		if len(tree)%2 != 0 {
			tree = append(tree, tree[len(tree)-1])
		}

		var nextLevel []string
		for i := 0; i < len(tree); i += 2 {
			combined := tree[i] + tree[i+1]
			nextLevel = append(nextLevel, hashData(combined))
		}
		tree = nextLevel
	}

	return tree[0]
}

// hashTransaction hashes a transaction
func hashTransaction(tx Transaction) string {
	txData := tx.ID + tx.SenderWallet + tx.ReceiverWallet + 
		floatToString(tx.Amount) + floatToString(tx.Fee) + tx.Signature
	return hashData(txData)
}

// hashData hashes data using SHA256
func hashData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// floatToString converts float64 to string
func floatToString(f float64) string {
	return string(rune(int(f * 100000000)))
}

// CalculateMerkleRoot calculates merkle root from transaction hashes (public version)
func CalculateMerkleRoot(txHashes []string) string {
	if len(txHashes) == 0 {
		return hashData("")
	}

	if len(txHashes) == 1 {
		return hashData(txHashes[0])
	}

	var tree []string
	for _, hash := range txHashes {
		tree = append(tree, hash)
	}

	for len(tree) > 1 {
		if len(tree)%2 != 0 {
			tree = append(tree, tree[len(tree)-1])
		}

		var nextLevel []string
		for i := 0; i < len(tree); i += 2 {
			combined := tree[i] + tree[i+1]
			nextLevel = append(nextLevel, hashData(combined))
		}
		tree = nextLevel
	}

	return tree[0]
}

