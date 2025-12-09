package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail validates an email format
func ValidateEmail(email string) bool {
	// Simple email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword validates password strength
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return false, "Password must contain uppercase, lowercase, digit, and special character"
	}

	return true, ""
}

// ValidateWalletAddress validates a wallet address format
func ValidateWalletAddress(address string) bool {
	// Wallet address should be 64 characters (SHA256 hash)
	return len(address) == 64 && isHexString(address)
}

// ValidateCNIC validates CNIC format
func ValidateCNIC(cnic string) bool {
	// Pakistani CNIC format: 5 digits - 7 digits - 1 digit
	cnicRegex := regexp.MustCompile(`^\d{5}-\d{7}-\d{1}$`)
	return cnicRegex.MatchString(cnic)
}

// ValidateAmount validates transaction amount
func ValidateAmount(amount float64) (bool, string) {
	if amount <= 0 {
		return false, "Amount must be greater than 0"
	}

	if amount > 1000000000 {
		return false, "Amount exceeds maximum limit"
	}

	return true, ""
}

// ValidateFee validates transaction fee
func ValidateFee(fee float64) (bool, string) {
	if fee < 0 {
		return false, "Fee cannot be negative"
	}

	if fee > 1000000 {
		return false, "Fee exceeds maximum limit"
	}

	return true, ""
}

// SanitizeInput sanitizes user input
func SanitizeInput(input string) string {
	// Remove SQL injection characters
	replacer := strings.NewReplacer(
		"'", "",
		"\"", "",
		";", "",
		"--", "",
		"/*", "",
		"*/", "",
	)
	return replacer.Replace(input)
}

// isHexString checks if a string contains only hexadecimal characters
func isHexString(s string) bool {
	hexRegex := regexp.MustCompile(`^[0-9a-fA-F]+$`)
	return hexRegex.MatchString(s)
}

// ValidateNote validates transaction note
func ValidateNote(note string) (bool, string) {
	if len(note) > 500 {
		return false, "Note cannot exceed 500 characters"
	}
	return true, ""
}

// ValidateOTP validates OTP format
func ValidateOTP(otp string) bool {
	otpRegex := regexp.MustCompile(`^\d{6}$`)
	return otpRegex.MatchString(otp)
}
