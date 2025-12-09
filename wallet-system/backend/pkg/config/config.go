package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port                    string
	NodeEnv                 string
	DatabaseURL             string
	JWTSecret               string
	ZakatPoolWallet         string
	ZakatPercentage         float64
	CORSAllowedOrigins      []string
	LogLevel                string
	LogFormat               string
	SupabaseURL             string
	SupabaseAnonKey         string
	SupabaseServiceRoleKey  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:                   getEnv("PORT", "8080"),
		NodeEnv:               getEnv("NODE_ENV", "development"),
		DatabaseURL:           getEnv("DATABASE_URL", ""),
		JWTSecret:             getEnv("JWT_SECRET", "your-jwt-secret-key"),
		ZakatPoolWallet:       getEnv("ZAKAT_POOL_WALLET", "zakat_pool"),
		ZakatPercentage:       2.5,
		CORSAllowedOrigins:    []string{"http://localhost:5173", "http://localhost:3000"},
		LogLevel:              getEnv("LOG_LEVEL", "info"),
		LogFormat:             getEnv("LOG_FORMAT", "json"),
		SupabaseURL:           getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:       getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
