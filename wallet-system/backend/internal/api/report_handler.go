package api

import (
	"context"
	"net/http"
	"time"

	"crypto-wallet-backend/internal/database"

	"github.com/gin-gonic/gin"
)

// GetMonthlyReportHandler gets monthly transaction report
func (h *Handler) GetMonthlyReportHandler(c *gin.Context) {
	walletAddress := c.Query("wallet_address")
	if walletAddress == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "wallet_address is required", Code: "INVALID_REQUEST"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all transactions for this wallet
	txns, err := h.db.GetTransactionsByWallet(ctx, walletAddress, 1000, 0)
	if err != nil {
		h.logger.Error("Failed to get transactions: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: "DB_ERROR"})
		return
	}

	// Group by month
	type MonthData struct {
		Month     string  `json:"month"`
		Incoming  float64 `json:"incoming"`
		Outgoing  float64 `json:"outgoing"`
		Fee       float64 `json:"fee"`
		Net       float64 `json:"net"`
		Count     int     `json:"count"`
	}
	monthlyData := make(map[string]MonthData)

	for _, tx := range txns {
		date, _ := time.Parse(time.RFC3339, tx.CreatedAt.String())
		monthKey := date.Format("2006-01")

		if _, exists := monthlyData[monthKey]; !exists {
			monthlyData[monthKey] = MonthData{
				Month: monthKey,
			}
		}

		isOutgoing := tx.SenderWallet == walletAddress

		m := monthlyData[monthKey]
		m.Count++

		if isOutgoing {
			m.Outgoing += tx.Amount
			m.Fee += tx.Fee
		} else {
			m.Incoming += tx.Amount
		}

		m.Net = m.Incoming - m.Outgoing
		monthlyData[monthKey] = m
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Monthly report retrieved",
		Data: gin.H{
			"monthly_data": monthlyData,
			"total_months": len(monthlyData),
		},
	})
}

// GetZakatReportHandler gets zakat deduction report
func (h *Handler) GetZakatReportHandler(c *gin.Context) {
	walletAddress := c.Query("wallet_address")
	if walletAddress == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "wallet_address is required", Code: "INVALID_REQUEST"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get wallet info
	wallet, err := h.db.GetWalletByAddress(ctx, walletAddress)
	if err != nil {
		h.logger.Error("Failed to get wallet: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: "DB_ERROR"})
		return
	}

	if wallet == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Wallet not found", Code: "NOT_FOUND"})
		return
	}

	// Calculate current zakat (2.5%)
	currentZakat := wallet.BalanceCache * 0.025

	// Mock zakat history - in production fetch from zakat_transactions table
	zakatHistory := []gin.H{
		{
			"month":         time.Now().Format("2006-01"),
			"zakat_amount":  currentZakat,
			"percentage":    2.5,
			"deducted_date": time.Now().Format("2006-01-02"),
		},
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Zakat report retrieved",
		Data: gin.H{
			"current_balance":      wallet.BalanceCache,
			"monthly_zakat":        currentZakat,
			"zakat_history":        zakatHistory,
			"is_deducted_this_month": wallet.ZakatDeducted,
		},
	})
}

// AddBeneficiaryHandler adds a beneficiary
func (h *Handler) AddBeneficiaryHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized", Code: "UNAUTHORIZED"})
		return
	}

	var req struct {
		BeneficiaryWalletID string `json:"beneficiary_wallet_id" binding:"required"`
		Nickname            string `json:"nickname"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Code: "INVALID_REQUEST"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	beneficiary := &database.Beneficiary{
		UserID:              userID,
		BeneficiaryWalletID: req.BeneficiaryWalletID,
		Nickname:            req.Nickname,
	}

	if err := h.db.CreateBeneficiary(ctx, beneficiary); err != nil {
		h.logger.Error("Failed to create beneficiary: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: "DB_ERROR"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Beneficiary added successfully",
		Data: gin.H{
			"beneficiary": beneficiary,
		},
	})
}

// GetBeneficiariesHandler gets all beneficiaries for a user
func (h *Handler) GetBeneficiariesHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized", Code: "UNAUTHORIZED"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	beneficiaries, err := h.db.GetBeneficiariesByUserID(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get beneficiaries: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: "DB_ERROR"})
		return
	}

	if beneficiaries == nil {
		beneficiaries = []*database.Beneficiary{}
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Beneficiaries retrieved",
		Data: gin.H{
			"beneficiaries": beneficiaries,
			"count":         len(beneficiaries),
		},
	})
}
