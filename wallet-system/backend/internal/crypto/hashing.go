package crypto

import (
	"crypto/rand"
	"math/big"
)

// GenerateOTP generates a random 6-digit OTP
func GenerateOTP() string {
	max := big.NewInt(1000000)
	randomNum, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "000000"
	}

	otp := randomNum.Int64()
	if otp < 100000 {
		otp += 100000
	}

	return string(rune(otp))
}
