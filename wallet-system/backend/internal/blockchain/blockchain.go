package blockchain

// Blockchain represents the entire blockchain
type Blockchain struct {
	Chain  []*Block
	Blocks map[string]*Block
	UTXOs  map[string][]UTXO
}

// NewBlockchain creates a new blockchain with genesis block
func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Chain:  make([]*Block, 0),
		Blocks: make(map[string]*Block),
		UTXOs:  make(map[string][]UTXO),
	}

	// Create genesis block
	genesisBlock := NewBlock(0, []Transaction{}, "0", 4)
	pow := NewProofOfWork(genesisBlock)
	pow.Mine()

	bc.AddBlock(genesisBlock)
	return bc
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(block *Block) error {
	if len(bc.Chain) > 0 {
		lastBlock := bc.Chain[len(bc.Chain)-1]
		if block.PreviousHash != lastBlock.Hash {
			return ErrInvalidPreviousHash
		}
	}

	bc.Chain = append(bc.Chain, block)
	bc.Blocks[block.Hash] = block
	return nil
}

// GetLatestBlock returns the last block in the chain
func (bc *Blockchain) GetLatestBlock() *Block {
	if len(bc.Chain) == 0 {
		return nil
	}
	return bc.Chain[len(bc.Chain)-1]
}

// AddUTXO adds a UTXO to the blockchain
func (bc *Blockchain) AddUTXO(walletAddress string, utxo UTXO) {
	bc.UTXOs[walletAddress] = append(bc.UTXOs[walletAddress], utxo)
}

// GetUTXOs returns all UTXOs for a wallet
func (bc *Blockchain) GetUTXOs(walletAddress string) []UTXO {
	return bc.UTXOs[walletAddress]
}

// GetUnspentUTXOs returns all unspent UTXOs for a wallet
func (bc *Blockchain) GetUnspentUTXOs(walletAddress string) []UTXO {
	var unspent []UTXO
	for _, utxo := range bc.UTXOs[walletAddress] {
		if !utxo.IsSpent {
			unspent = append(unspent, utxo)
		}
	}
	return unspent
}

// MarkUTXOAsSpent marks a UTXO as spent
func (bc *Blockchain) MarkUTXOAsSpent(walletAddress, txHash string, outputIndex int) {
	for i, utxo := range bc.UTXOs[walletAddress] {
		if utxo.TransactionHash == txHash && utxo.OutputIndex == outputIndex {
			bc.UTXOs[walletAddress][i].IsSpent = true
			bc.UTXOs[walletAddress][i].SpentInTx = txHash
		}
	}
}

// ValidateChain validates the entire blockchain
func (bc *Blockchain) ValidateChain() bool {
	for i := 1; i < len(bc.Chain); i++ {
		block := bc.Chain[i]
		previousBlock := bc.Chain[i-1]

		if block.PreviousHash != previousBlock.Hash {
			return false
		}

		pow := NewProofOfWork(block)
		if !pow.ValidateProof() {
			return false
		}
	}
	return true
}

// GetBalance calculates the balance of a wallet from UTXOs
func (bc *Blockchain) GetBalance(walletAddress string) float64 {
	var balance float64
	unspentUTXOs := bc.GetUnspentUTXOs(walletAddress)
	for _, utxo := range unspentUTXOs {
		balance += utxo.Amount
	}
	return balance
}
