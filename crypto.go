package authora

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

type KeyPair struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

func GenerateKeyPair() (*KeyPair, error) {
	seed := make([]byte, ed25519.SeedSize)
	if _, err := rand.Read(seed); err != nil {
		return nil, fmt.Errorf("authora: failed to generate key: %w", err)
	}
	privKey := ed25519.NewKeyFromSeed(seed)
	pubKey := privKey.Public().(ed25519.PublicKey)
	return &KeyPair{
		PrivateKey: ToBase64URL(seed),
		PublicKey:  ToBase64URL(pubKey),
	}, nil
}

func GetPublicKey(privateKeyB64 string) (string, error) {
	seed, err := FromBase64URL(privateKeyB64)
	if err != nil {
		return "", fmt.Errorf("authora: invalid private key: %w", err)
	}
	privKey := ed25519.NewKeyFromSeed(seed)
	pubKey := privKey.Public().(ed25519.PublicKey)
	return ToBase64URL(pubKey), nil
}

func ToBase64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func FromBase64URL(b64url string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(b64url)
}

func BuildSignaturePayload(method, path, timestamp string, body *string) string {
	bodyHash := ""
	if body != nil && *body != "" {
		h := sha256.Sum256([]byte(*body))
		bodyHash = ToBase64URL(h[:])
	}
	return strings.ToUpper(method) + "\n" + path + "\n" + timestamp + "\n" + bodyHash
}

func Sign(message, privateKeyB64 string) (string, error) {
	seed, err := FromBase64URL(privateKeyB64)
	if err != nil {
		return "", fmt.Errorf("authora: invalid private key: %w", err)
	}
	privKey := ed25519.NewKeyFromSeed(seed)
	sig := ed25519.Sign(privKey, []byte(message))
	return ToBase64URL(sig), nil
}

func Verify(message, signatureB64, publicKeyB64 string) bool {
	sigBytes, err := FromBase64URL(signatureB64)
	if err != nil {
		return false
	}
	pubBytes, err := FromBase64URL(publicKeyB64)
	if err != nil {
		return false
	}
	return ed25519.Verify(ed25519.PublicKey(pubBytes), []byte(message), sigBytes)
}

func SHA256Hash(input string) string {
	h := sha256.Sum256([]byte(input))
	return ToBase64URL(h[:])
}
