package crypto

import (
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Errorf("Failed to generate key pair: %v", err)
	}

	if kp.PrivateKey == "" {
		t.Errorf("Private key is empty")
	}

	if kp.PublicKey == "" {
		t.Errorf("Public key is empty")
	}

	if kp.WalletID == "" {
		t.Errorf("Wallet ID is empty")
	}

	if len(kp.WalletID) != 64 {
		t.Errorf("Expected wallet ID length 64, got %d", len(kp.WalletID))
	}
}

func TestEncryptDecryptPrivateKey(t *testing.T) {
	kp, _ := GenerateKeyPair()
	password := "testpassword123"

	encrypted, err := EncryptPrivateKey(kp.PrivateKey, password)
	if err != nil {
		t.Errorf("Failed to encrypt private key: %v", err)
	}

	decrypted, err := DecryptPrivateKey(encrypted, password)
	if err != nil {
		t.Errorf("Failed to decrypt private key: %v", err)
	}

	if decrypted != kp.PrivateKey {
		t.Errorf("Decrypted key doesn't match original")
	}

	// Test wrong password
	_, err = DecryptPrivateKey(encrypted, "wrongpassword")
	if err == nil {
		t.Errorf("Should fail with wrong password")
	}
}

func TestSHA256Hash(t *testing.T) {
	data := "test data"
	hash1 := SHA256Hash(data)
	hash2 := SHA256Hash(data)

	if hash1 != hash2 {
		t.Errorf("Same data should produce same hash")
	}

	if len(hash1) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash1))
	}

	hash3 := SHA256Hash("different data")
	if hash1 == hash3 {
		t.Errorf("Different data should produce different hash")
	}
}

func TestSignAndVerifyTransaction(t *testing.T) {
	kp, _ := GenerateKeyPair()
	data := "transaction data"

	signature, err := SignTransaction(data, kp.PrivateKey)
	if err != nil {
		t.Errorf("Failed to sign transaction: %v", err)
	}

	if signature == "" {
		t.Errorf("Signature is empty")
	}

	valid, err := VerifySignature(data, signature, kp.PublicKey)
	if err != nil {
		t.Errorf("Failed to verify signature: %v", err)
	}

	if !valid {
		t.Errorf("Signature verification failed")
	}

	// Test with wrong data
	valid, err = VerifySignature("different data", signature, kp.PublicKey)
	if valid {
		t.Errorf("Should fail verification with different data")
	}
}
