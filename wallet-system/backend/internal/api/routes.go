package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all API routes
func SetupRoutes(router *gin.Engine, handler *Handler) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Authentication routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", handler.RegisterHandler)
		auth.POST("/login", handler.LoginHandler)
	}

	// Wallet routes
	wallet := router.Group("/api/wallet")
	wallet.Use(AuthMiddleware(handler.jwtSecret))
	{
		wallet.GET("/profile", handler.GetWalletHandler)
		wallet.POST("/balance", handler.GetBalanceHandler)
	}

	// Blockchain routes
	blockchain := router.Group("/api/blockchain")
	{
		blockchain.GET("/blocks", handler.GetBlocksHandler)
		blockchain.GET("/latest", handler.GetLatestBlockHandler)
		blockchain.GET("/blocks/:hash", handler.GetBlockByHashHandler)
		blockchain.POST("/mine", handler.MineBlockHandler)
	}

	// Transaction routes
	transaction := router.Group("/api/transaction")
	transaction.Use(AuthMiddleware(handler.jwtSecret))
	{
		transaction.GET("/pending", func(c *gin.Context) {
			// Get pending transactions
			c.JSON(200, gin.H{"transactions": []interface{}{}})
		})
		transaction.GET("/history", handler.GetTransactionHistoryHandler)
		transaction.POST("/send", handler.SendTransactionHandler)
	}

	// Reports routes
	reports := router.Group("/api/reports")
	reports.Use(AuthMiddleware(handler.jwtSecret))
	{
		reports.GET("/monthly", handler.GetMonthlyReportHandler)
		reports.GET("/zakat", handler.GetZakatReportHandler)
	}

	// Beneficiary routes
	beneficiary := router.Group("/api/beneficiary")
	beneficiary.Use(AuthMiddleware(handler.jwtSecret))
	{
		beneficiary.POST("/add", handler.AddBeneficiaryHandler)
		beneficiary.GET("/list", handler.GetBeneficiariesHandler)
	}

	// System routes
	system := router.Group("/api/system")
	system.Use(AuthMiddleware(handler.jwtSecret))
	{
		system.GET("/logs", handler.GetSystemLogsHandler)
		system.GET("/logs/stats", handler.GetSystemLogStatsHandler)
		system.GET("/health", handler.GetSystemHealthHandler)
	}
}
