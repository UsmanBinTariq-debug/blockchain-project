package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"crypto-wallet-backend/internal/blockchain"
	"crypto-wallet-backend/internal/database"

	"github.com/gin-gonic/gin"
)

// GetBlocksRequest represents a request to get blocks
type GetBlocksRequest struct {
	Limit  int `form:"limit" binding:"min=1,max=100"`
	Offset int `form:"offset" binding:"min=0"`
}

// GetBlocksHandler retrieves all blocks
func (h *Handler) GetBlocksHandler(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blocks, err := h.db.GetBlocks(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get blocks: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: "DB_ERROR"})
		return
	}

	if blocks == nil {
		blocks = []*database.Block{}
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Blocks retrieved",
		Data: gin.H{
			"blocks": blocks,
			"count":  len(blocks),
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetLatestBlockHandler retrieves the latest block
func (h *Handler) GetLatestBlockHandler(c *gin.Context) {
	latestBlock := h.bc.GetLatestBlock()
	if latestBlock == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "No blocks found", Code: "NOT_FOUND"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Latest block retrieved",
		Data: gin.H{
			"block": latestBlock,
		},
	})
}

// GetBlockByHashHandler retrieves a specific block by hash
func (h *Handler) GetBlockByHashHandler(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Block hash is required", Code: "INVALID_REQUEST"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	block, err := h.db.GetBlockByHash(ctx, hash)
	if err != nil {
		h.logger.Error("Failed to get block: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: "DB_ERROR"})
		return
	}

	if block == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Block not found", Code: "NOT_FOUND"})
		return
	}

	// Get transactions for this block
	txns, _ := h.db.GetTransactionsByBlockHash(ctx, hash, 100, 0)
	if txns == nil {
		txns = []*database.Transaction{}
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Block retrieved",
		Data: gin.H{
			"block":        block,
			"transactions": txns,
			"tx_count":     len(txns),
		},
	})
}

// MineBlockRequest represents a mine block request
type MineBlockRequest struct {
	MinerAddress string `json:"miner_address" binding:"required"`
}

// MineBlockHandler mines a new block
func (h *Handler) MineBlockHandler(c *gin.Context) {
	var req MineBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Code: "INVALID_REQUEST"})
		return
	}

	// Validate miner wallet address
	if !ValidateWalletAddress(req.MinerAddress) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid miner address", Code: "INVALID_WALLET"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second) // 2 min timeout for mining
	defer cancel()

	// Get pending transactions
	pendingTxns, err := h.db.GetTransactionsByStatus(ctx, "pending", 10)
	if err != nil {
		h.logger.Error("Failed to get pending transactions: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get transactions", Code: "TX_ERROR"})
		return
	}

	if len(pendingTxns) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No pending transactions to mine", Code: "NO_TRANSACTIONS"})
		return
	}

	// Mine the block
	latestBlock := h.bc.GetLatestBlock()
	if latestBlock == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "No genesis block", Code: "NO_GENESIS"})
		return
	}

	// Create new block
	newBlock := &blockchain.Block{
		Index:        latestBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
		PreviousHash: latestBlock.Hash,
		Difficulty:   latestBlock.Difficulty,
		MinedBy:      req.MinerAddress,
	}

	// Add transactions to block (use ID as hash)
	txHashes := make([]string, 0)
	for _, tx := range pendingTxns {
		newBlock.Transactions = append(newBlock.Transactions, blockchain.Transaction{
			ID:             tx.TransactionHash,
			SenderWallet:   tx.SenderWallet,
			ReceiverWallet: tx.ReceiverWallet,
			Amount:         tx.Amount,
			Fee:            tx.Fee,
		})
		txHashes = append(txHashes, tx.TransactionHash)
	}

	// Calculate merkle root
	merkleRoot := blockchain.CalculateMerkleRoot(txHashes)
	newBlock.MerkleRoot = merkleRoot

	// Perform proof of work - call Mine() directly
	pow := blockchain.NewProofOfWork(newBlock)
	pow.Mine() // This sets block.Hash and block.Nonce

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

	if err := h.db.CreateBlock(ctx, dbBlock); err != nil {
		h.logger.Error("Failed to save block: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save block", Code: "DB_ERROR"})
		return
	}

	// Update transactions to confirmed
	for _, tx := range pendingTxns {
		if err := h.db.UpdateTransactionStatus(ctx, tx.TransactionHash, "confirmed", newBlock.Hash); err != nil {
			h.logger.Error("Failed to update transaction status: %v", err)
		}
	}

	// Add block to blockchain
	h.bc.AddBlock(newBlock)

	// Mine reward: add UTXO to miner wallet
	rewardUTXO := &database.UTXO{
		TransactionHash: fmt.Sprintf("reward-%s-%d", newBlock.Hash, time.Now().Unix()),
		OutputIndex:     0,
		WalletAddress:   req.MinerAddress,
		Amount:          5.0, // Mining reward
		IsSpent:         false,
	}
	h.db.CreateUTXO(ctx, rewardUTXO)

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Block mined successfully",
		Data: gin.H{
			"block": gin.H{
				"index":      newBlock.Index,
				"hash":       newBlock.Hash,
				"nonce":      newBlock.Nonce,
				"difficulty": newBlock.Difficulty,
				"tx_count":   len(newBlock.Transactions),
				"mined_by":   newBlock.MinedBy,
				"reward":     5.0,
				"timestamp":  newBlock.Timestamp,
			},
		},
	})
}

// ValidateWalletAddress validates a wallet address format
func ValidateWalletAddress(address string) bool {
	if len(address) != 64 {
		return false
	}
	for _, c := range address {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
