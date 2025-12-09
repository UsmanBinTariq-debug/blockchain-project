package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
)

// SignTransaction signs a transaction with a private key
func SignTransaction(data string, privateKeyBase64 string) (string, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, hash[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySignature verifies a transaction signature
func VerifySignature(data string, signatureBase64 string, publicKeyBase64 string) (bool, error) {
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, err
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return false, err
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return false, err
	}

	rsaPublicKey := publicKey.(*rsa.PublicKey)
	hash := sha256.Sum256([]byte(data))

	err = rsa.VerifyPKCS1v15(rsaPublicKey, 0, hash[:], signature)
	if err != nil {
		return false, err
	}

	return true, nil
}
