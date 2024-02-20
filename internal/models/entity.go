package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EntityRepository interface {
	Create(ctx context.Context, entity *Entity) error

	GetEntityById(ctx context.Context, id uuid.UUID) (*Entity, error)
	GetEntityByName(ctx context.Context, name string) (*Entity, error)
	GetEntityByWalletAddress(ctx context.Context, walletAddress string) (*Entity, error)

	DeleteEntityById(ctx context.Context, id uuid.UUID) error
}

type Entity struct {
	Id            uuid.UUID `json:"id" gorm:"primary_key,type:uuid"`
	Name          string    `json:"name" gorm:"text,not null"`
	WalletAddress string    `json:"wallet_address" gorm:"text,not null"`
	CreatedAt     int64     `json:"created_at" gorm:"int8,not null"`
	Wallet        *Wallet   `json:"-" gorm:"foreignkey:WalletAddress;"`
}

func NewEntity(name string, wallet *Wallet) *Entity {
	return &Entity{
		Id:            uuid.New(),
		Name:          name,
		WalletAddress: wallet.Address,
		CreatedAt:     time.Now().UTC().Unix(),
		Wallet:        wallet,
	}
}
