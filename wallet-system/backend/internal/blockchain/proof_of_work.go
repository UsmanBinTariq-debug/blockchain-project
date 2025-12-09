package blockchain

import (
	"fmt"
)

// ProofOfWork represents a Proof-of-Work mining operation
type ProofOfWork struct {
	Block      *Block
	Difficulty int
	Target     string
}

// NewProofOfWork creates a new ProofOfWork instance
func NewProofOfWork(block *Block) *ProofOfWork {
	target := ""
	for i := 0; i < block.Difficulty; i++ {
		target += "0"
	}

	return &ProofOfWork{
		Block:      block,
		Difficulty: block.Difficulty,
		Target:     target,
	}
}

// Mine performs the Proof-of-Work mining
func (pow *ProofOfWork) Mine() {
	var nonce int64 = 0
	var hash string

	fmt.Printf("Mining block with difficulty %d...\n", pow.Difficulty)

	for {
		pow.Block.Nonce = nonce
		hash = pow.Block.CalculateHash()

		if hash[:pow.Difficulty] == pow.Target {
			fmt.Printf("Block mined! Hash: %s, Nonce: %d\n", hash, nonce)
			pow.Block.Hash = hash
			pow.Block.Nonce = nonce
			break
		}

		nonce++

		if nonce%100000 == 0 {
			fmt.Printf("Nonce: %d, Hash: %s\n", nonce, hash)
		}
	}
}

// ValidateProof validates if the block's hash meets the Proof-of-Work requirements
func (pow *ProofOfWork) ValidateProof() bool {
	hash := pow.Block.CalculateHash()
	pow.Block.Hash = hash
	return hash[:pow.Difficulty] == pow.Target
}

// AdjustDifficulty adjusts the difficulty based on mining time
func AdjustDifficulty(previousDifficulty int, previousBlock *Block, currentTime int64) int {
	miningTime := currentTime - previousBlock.Timestamp
	targetTime := int64(60) // 1 minute target mining time

	// If mining took longer than target, decrease difficulty
	if miningTime > targetTime {
		if previousDifficulty > 1 {
			return previousDifficulty - 1
		}
		return 1
	}

	// If mining took less than target, increase difficulty
	if miningTime < targetTime {
		return previousDifficulty + 1
	}

	return previousDifficulty
}
