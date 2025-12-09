package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"crypto-wallet-backend/internal/database"
	"crypto-wallet-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// SendTransactionRequest represents a send transaction request
type SendTransactionRequest struct {
	SenderWallet   string  `json:"sender_wallet" binding:"required"`
	ReceiverWallet string  `json:"receiver_wallet" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	Fee            float64 `json:"fee" binding:"required"`
	Note           string  `json:"note"`
	Signature      string  `json:"signature" binding:"required"`
}

// SendTransactionHandler sends a transaction
func (h *Handler) SendTransactionHandler(c *gin.Context) {
	var req SendTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("SendTransaction bind error: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("Invalid request: %v", err),
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Validate amount
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Amount must be greater than 0",
			Code:  "INVALID_AMOUNT",
		})
		return
	}

	// Validate fee
	if req.Fee < 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Fee cannot be negative",
			Code:  "INVALID_FEE",
		})
		return
	}

	// Validate wallet addresses
	if req.SenderWallet == "" || req.ReceiverWallet == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Sender and receiver wallets are required",
			Code:  "INVALID_WALLET",
		})
		return
	}

	if !utils.ValidateWalletAddress(req.SenderWallet) || !utils.ValidateWalletAddress(req.ReceiverWallet) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid wallet address format",
			Code:  "INVALID_WALLET",
		})
		return
	}

	if req.SenderWallet == req.ReceiverWallet {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Cannot send to the same wallet",
			Code:  "SAME_WALLET",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create and execute transaction
	txHash, err := h.transactionService.CreateTransaction(ctx, database.Transaction{
		SenderWallet:   req.SenderWallet,
		ReceiverWallet: req.ReceiverWallet,
		Amount:         req.Amount,
		Fee:            req.Fee,
		Note:           req.Note,
		Signature:      req.Signature,
		Status:         "pending",
		CreatedAt:      time.Now(),
	})
	if err != nil {
		h.logger.Error("Failed to create transaction: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
			Code:  "TRANSACTION_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Transaction created",
		Data: gin.H{
			"transaction_hash": txHash,
			"sender":           req.SenderWallet,
			"receiver":         req.ReceiverWallet,
			"amount":           req.Amount,
			"fee":              req.Fee,
		},
	})
}

// GetTransactionHistoryHandler retrieves transaction history for a wallet
func (h *Handler) GetTransactionHistoryHandler(c *gin.Context) {
	wallet := c.Query("wallet_address")
	limit := 10
	offset := 0

	if wallet == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Wallet address is required",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	txns, err := h.transactionService.GetTransactionHistory(ctx, wallet, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get transaction history: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
			Code:  "TRANSACTION_ERROR",
		})
		return
	}

	if txns == nil {
		txns = []*database.Transaction{}
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Transaction history retrieved",
		Data: gin.H{
			"transactions": txns,
			"count":        len(txns),
		},
	})
}
