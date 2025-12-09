package blockchain

import (
	"testing"
)

func TestNewBlock(t *testing.T) {
	block := NewBlock(0, []Transaction{}, "0", 4)

	if block.Index != 0 {
		t.Errorf("Expected index 0, got %d", block.Index)
	}

	if block.PreviousHash != "0" {
		t.Errorf("Expected previous hash 0, got %s", block.PreviousHash)
	}

	if block.Difficulty != 4 {
		t.Errorf("Expected difficulty 4, got %d", block.Difficulty)
	}
}

func TestBlockHash(t *testing.T) {
	block1 := NewBlock(0, []Transaction{}, "0", 4)
	block2 := NewBlock(0, []Transaction{}, "0", 4)

	hash1 := block1.CalculateHash()
	hash2 := block2.CalculateHash()

	if hash1 != hash2 {
		t.Errorf("Same blocks should have same hash")
	}

	block1.Nonce = 1
	hash3 := block1.CalculateHash()

	if hash1 == hash3 {
		t.Errorf("Different nonces should produce different hashes")
	}
}

func TestNewTransaction(t *testing.T) {
	tx := NewTransaction("sender", "receiver", 10.5, 0.1, "payment")

	if tx.SenderWallet != "sender" {
		t.Errorf("Expected sender 'sender', got %s", tx.SenderWallet)
	}

	if tx.Amount != 10.5 {
		t.Errorf("Expected amount 10.5, got %f", tx.Amount)
	}

	if tx.Status != "pending" {
		t.Errorf("Expected status 'pending', got %s", tx.Status)
	}
}

func TestNewBlockchain(t *testing.T) {
	bc := NewBlockchain()

	if len(bc.Chain) != 1 {
		t.Errorf("Expected 1 genesis block, got %d", len(bc.Chain))
	}

	genesisBlock := bc.GetLatestBlock()
	if genesisBlock.Index != 0 {
		t.Errorf("Expected genesis block index 0, got %d", genesisBlock.Index)
	}
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockchain()
	lastBlock := bc.GetLatestBlock()

	newBlock := NewBlock(1, []Transaction{}, lastBlock.Hash, 4)
	pow := NewProofOfWork(newBlock)
	pow.Mine()

	err := bc.AddBlock(newBlock)
	if err != nil {
		t.Errorf("Failed to add block: %v", err)
	}

	if len(bc.Chain) != 2 {
		t.Errorf("Expected 2 blocks, got %d", len(bc.Chain))
	}
}

func TestProofOfWork(t *testing.T) {
	block := NewBlock(0, []Transaction{}, "0", 2)
	pow := NewProofOfWork(block)

	pow.Mine()

	if block.Hash[:pow.Difficulty] != pow.Target {
		t.Errorf("Proof of work failed validation")
	}

	if !pow.ValidateProof() {
		t.Errorf("Proof of work validation failed")
	}
}

func TestGetBalance(t *testing.T) {
	bc := NewBlockchain()

	balance := bc.GetBalance("wallet1")
	if balance != 0 {
		t.Errorf("Expected balance 0, got %f", balance)
	}

	bc.AddUTXO("wallet1", UTXO{
		WalletAddress: "wallet1",
		Amount:        50.5,
		IsSpent:       false,
	})

	balance = bc.GetBalance("wallet1")
	if balance != 50.5 {
		t.Errorf("Expected balance 50.5, got %f", balance)
	}
}

func TestValidateChain(t *testing.T) {
	bc := NewBlockchain()

	if !bc.ValidateChain() {
		t.Errorf("Genesis chain should be valid")
	}

	lastBlock := bc.GetLatestBlock()
	newBlock := NewBlock(1, []Transaction{}, lastBlock.Hash, 4)
	pow := NewProofOfWork(newBlock)
	pow.Mine()

	bc.AddBlock(newBlock)

	if !bc.ValidateChain() {
		t.Errorf("Valid chain should pass validation")
	}
}
