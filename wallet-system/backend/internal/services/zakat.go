package services

import (
	"context"
	"crypto-wallet-backend/internal/database"
	"fmt"
	"time"
)

// ZakatService handles zakat operations
type ZakatService struct {
	db *database.Database
}

// NewZakatService creates a new zakat service
func NewZakatService(db *database.Database) *ZakatService {
	return &ZakatService{db: db}
}

// DeductZakat deducts zakat from a wallet
func (zs *ZakatService) DeductZakat(ctx context.Context, walletAddress string, balance float64, zakatPercentage float64) error {
	// Calculate zakat amount
	zakatAmount := balance * (zakatPercentage / 100)

	// Create zakat transaction
	zakatTx := &database.ZakatTransaction{
		WalletAddress:   walletAddress,
		Amount:          zakatAmount,
		ZakatPercentage: zakatPercentage,
		MonthYear:       time.Now().Format("2006-01"),
	}

	if err := zs.db.CreateZakatTransaction(ctx, zakatTx); err != nil {
		return err
	}

	// Mark month as processed
	wallet, err := zs.db.GetWalletByAddress(ctx, walletAddress)
	if err != nil {
		return err
	}

	if wallet != nil {
		// Update wallet to mark zakat as deducted
		newBalance := balance - zakatAmount
		if err := zs.db.UpdateWalletBalance(ctx, walletAddress, newBalance); err != nil {
			return err
		}
	}

	return nil
}

// ProcessMonthlyZakat processes zakat for all wallets
func (zs *ZakatService) ProcessMonthlyZakat(ctx context.Context) error {
	// Get all wallets (you would implement a method in database to get all wallets)
	// For now, this is a placeholder
	fmt.Println("Processing monthly zakat deductions...")
	return nil
}

// GetZakatReports returns zakat reports for a wallet
func (zs *ZakatService) GetZakatReports(ctx context.Context, walletAddress string) ([]*database.ZakatTransaction, error) {
	// Query database for zakat transactions
	// This would require adding a method to the database service
	return []*database.ZakatTransaction{}, nil
}
