package services

import (
	"context"
	"crypto-wallet-backend/internal/blockchain"
	"crypto-wallet-backend/internal/database"
	"fmt"
)

// MiningService handles mining operations
type MiningService struct {
	db *database.Database
	bc *blockchain.Blockchain
}

// NewMiningService creates a new mining service
func NewMiningService(db *database.Database, bc *blockchain.Blockchain) *MiningService {
	return &MiningService{db: db, bc: bc}
}

// MineBlock mines a new block
func (ms *MiningService) MineBlock(ctx context.Context, minerAddress string) (*blockchain.Block, error) {
	// Create block with pending transactions
	lastBlock := ms.bc.GetLatestBlock()
	if lastBlock == nil {
		return nil, fmt.Errorf("no genesis block found")
	}

	// Adjust difficulty
	newDifficulty := blockchain.AdjustDifficulty(
		lastBlock.Difficulty,
		lastBlock,
		ms.bc.GetLatestBlock().Timestamp,
	)

	// Create new block
	newBlock := blockchain.NewBlock(
		lastBlock.Index+1,
		[]blockchain.Transaction{}, // Empty for now, would be populated with pending transactions
		lastBlock.Hash,
		newDifficulty,
	)
	newBlock.MinedBy = minerAddress

	// Mine the block
	pow := blockchain.NewProofOfWork(newBlock)
	pow.Mine()

	// Add to blockchain
	if err := ms.bc.AddBlock(newBlock); err != nil {
		return nil, err
	}

	// Save to database
	dbBlock := &database.Block{
		BlockIndex:   newBlock.Index,
		Timestamp:    newBlock.Timestamp,
		PreviousHash: newBlock.PreviousHash,
		Hash:         newBlock.Hash,
		Nonce:        newBlock.Nonce,
		MerkleRoot:   newBlock.MerkleRoot,
		Difficulty:   newBlock.Difficulty,
		MinedBy:      newBlock.MinedBy,
	}

	if err := ms.db.CreateBlock(ctx, dbBlock); err != nil {
		return nil, err
	}

	return newBlock, nil
}

// ValidateBlock validates a block
func (ms *MiningService) ValidateBlock(block *blockchain.Block) bool {
	pow := blockchain.NewProofOfWork(block)
	return pow.ValidateProof()
}

// ValidateChain validates the entire blockchain
func (ms *MiningService) ValidateChain() bool {
	return ms.bc.ValidateChain()
}
