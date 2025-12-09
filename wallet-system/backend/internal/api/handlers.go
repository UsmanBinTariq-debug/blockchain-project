package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"crypto-wallet-backend/internal/blockchain"
	"crypto-wallet-backend/internal/crypto"
	"crypto-wallet-backend/internal/database"
	"crypto-wallet-backend/internal/services"
	"crypto-wallet-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Handler holds handler dependencies
type Handler struct {
	db                *database.Database
	bc                *blockchain.Blockchain
	walletService     *services.WalletService
	zakatService      *services.ZakatService
	miningService     *services.MiningService
	transactionService *services.TransactionService
	logger            *utils.Logger
	jwtSecret         string
}

// NewHandler creates a new handler
func NewHandler(
	db *database.Database,
	bc *blockchain.Blockchain,
	jwtSecret string,
) *Handler {
	return &Handler{
		db:                 db,
		bc:                 bc,
		walletService:      services.NewWalletService(db, bc),
		zakatService:       services.NewZakatService(db),
		miningService:      services.NewMiningService(db, bc),
		transactionService: services.NewTransactionService(db),
		logger:             utils.NewLogger("info"),
		jwtSecret:          jwtSecret,
	}
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	CNIC     string `json:"cnic" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterHandler handles user registration
func (h *Handler) RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Code: "INVALID_REQUEST"})
		return
	}

	// Validate input
	if !utils.ValidateEmail(req.Email) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid email format", Code: "INVALID_EMAIL"})
		return
	}

	if valid, msg := utils.ValidatePassword(req.Password); !valid {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: msg, Code: "WEAK_PASSWORD"})
		return
	}

	if !utils.ValidateCNIC(req.CNIC) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid CNIC format", Code: "INVALID_CNIC"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user exists
	existingUser, err := h.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error", Code: "DB_ERROR"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "Email already registered", Code: "EMAIL_EXISTS"})
		return
	}

	// Generate key pair
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate keys", Code: "KEY_GEN_ERROR"})
		return
	}

	// Encrypt private key
	encryptedPrivateKey, err := crypto.EncryptPrivateKey(keyPair.PrivateKey, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to encrypt key", Code: "ENCRYPT_ERROR"})
		return
	}

	// Create user
	user := &database.User{
		Email:                req.Email,
		FullName:             req.FullName,
		CNIC:                 req.CNIC,
		WalletID:             keyPair.WalletID,
		PublicKey:            keyPair.PublicKey,
		EncryptedPrivateKey:  encryptedPrivateKey,
	}

	if err := h.db.CreateUser(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create user", Code: "USER_CREATE_ERROR"})
		return
	}

	// Create wallet
	wallet, _, err := h.walletService.CreateWallet(ctx, user.ID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create wallet", Code: "WALLET_CREATE_ERROR"})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data: gin.H{
			"user_id":       user.ID,
			"email":         user.Email,
			"wallet_id":     keyPair.WalletID,
			"wallet_address": wallet.WalletAddress,
		},
	})
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler handles user login
func (h *Handler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Code: "INVALID_REQUEST"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user
	user, err := h.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error", Code: "DB_ERROR"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials", Code: "INVALID_CREDENTIALS"})
		return
	}

	// Decrypt private key to verify password
	_, err = crypto.DecryptPrivateKey(user.EncryptedPrivateKey, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials", Code: "INVALID_CREDENTIALS"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token", Code: "TOKEN_ERROR"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Login successful",
		Data: gin.H{
			"token":   tokenString,
			"user_id": user.ID,
		},
	})
}

// GetWalletHandler returns wallet information
func (h *Handler) GetWalletHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	
	// If no user_id in context (no auth), get from query or use first registered user (dev mode)
	if userID == "" {
		// For testing: accept user_id from query param
		userID = c.Query("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "user_id required", Code: "MISSING_USER_ID"})
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get wallets for this user
	wallets, err := h.db.GetWalletsByUserID(ctx, userID)
	if err != nil || len(wallets) == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Wallet not found", Code: "WALLET_NOT_FOUND"})
		return
	}

	wallet := wallets[0]
	
	// Calculate balance from UTXOs
	utxos, err := h.db.GetUTXOsByWallet(ctx, wallet.WalletAddress)
	if err != nil {
		utxos = []*database.UTXO{}
	}
	
	balance := 0.0
	for _, utxo := range utxos {
		if !utxo.IsSpent {
			balance += utxo.Amount
		}
	}

	// If no UTXOs found but balance_cache exists, use the cached balance
	if balance == 0 && wallet.BalanceCache > 0 {
		balance = wallet.BalanceCache
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Wallet retrieved",
		Data: gin.H{
			"user_id":        userID,
			"wallet_id":      wallet.ID,
			"wallet_address": wallet.WalletAddress,
			"balance":        balance,
			"last_updated":   wallet.LastUpdated,
		},
	})
}

// GetBalanceRequest represents a balance request
type GetBalanceRequest struct {
	WalletAddress string `json:"wallet_address" binding:"required"`
}

// GetBalanceHandler returns wallet balance
func (h *Handler) GetBalanceHandler(c *gin.Context) {
	var req GetBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Code: "INVALID_REQUEST"})
		return
	}

	if !utils.ValidateWalletAddress(req.WalletAddress) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid wallet address", Code: "INVALID_WALLET"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	balance, err := h.walletService.GetWalletBalance(ctx, req.WalletAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get balance", Code: "BALANCE_ERROR"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Balance retrieved",
		Data: gin.H{
			"wallet_address": req.WalletAddress,
			"balance":        balance,
		},
	})
}

// GetSystemLogsHandler retrieves system logs with optional filtering
func (h *Handler) GetSystemLogsHandler(c *gin.Context) {
	logType := c.DefaultQuery("type", "ALL")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	logs, err := h.db.GetSystemLogs(c.Request.Context(), logType, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get system logs", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve system logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"logs":   logs,
			"limit":  limit,
			"offset": offset,
			"count":  len(logs),
		},
	})
}

// GetSystemLogStatsHandler retrieves system log statistics
func (h *Handler) GetSystemLogStatsHandler(c *gin.Context) {
	stats, err := h.db.GetSystemLogStats(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get system log stats", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve log statistics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}

// GetSystemHealthHandler returns system health status
func (h *Handler) GetSystemHealthHandler(c *gin.Context) {
	ctx := c.Request.Context()
	if err := h.db.Ping(ctx); err != nil {
		h.logger.Error("Database health check failed", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"health": gin.H{
				"database": "down",
				"api":      "up",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"health": gin.H{
			"database":   "up",
			"api":        "up",
			"blockchain": "up",
		},
	})
}
