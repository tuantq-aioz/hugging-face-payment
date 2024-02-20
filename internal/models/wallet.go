package models

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *Wallet) error

	GetActiveWallets(ctx context.Context) ([]*Wallet, error)
}

type Wallet struct {
	Address     string `gorm:"primary_key;index:idx_address;priority:1" json:"address,omitempty"`
	PrivateKey  []byte
	PublicKey   []byte
	Balance     decimal.Decimal `json:"balance" gorm:"type:numeric"`
	Debt        decimal.Decimal `json:"debt" gorm:"type:numeric"`
	FreeBalance decimal.Decimal `json:"free_balance" gorm:"type:numeric"`
	UserID      uuid.UUID       `json:"user_id"`
	CreatedAt   time.Time       `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt   time.Time       `gorm:"not null" json:"updated_at,omitempty"`
}

func NewWallet(passphrase string) (*Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	encryptedPrivateKey, err := encrypt(privateKeyBytes, passphrase)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	return &Wallet{
		Address:     crypto.PubkeyToAddress(*publicKeyECDSA).Hex(),
		PrivateKey:  encryptedPrivateKey,
		PublicKey:   publicKeyBytes,
		Balance:     decimal.Zero,
		Debt:        decimal.Zero,
		FreeBalance: decimal.Zero,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	hash := sha256.Sum256([]byte(passphrase))
	blockCipher, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func Decrypt(data []byte, passphrase string) ([]byte, error) {
	hash := sha256.Sum256([]byte(passphrase))
	blockCipher, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
