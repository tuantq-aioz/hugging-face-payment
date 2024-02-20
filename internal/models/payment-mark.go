package models

import "context"

type PaymentMarkRepository interface {
	Create(ctx context.Context, paymentMark *PaymentMark) error

	GetPaymentMarkByChainId(ctx context.Context, chainId int64) (*PaymentMark, error)

	UpdatePaymentMarkByChainId(ctx context.Context, chainId int64, blockNumber int64) error

	DeletePaymentMarkByChainId(ctx context.Context, chainId int64) error
}

type PaymentMark struct {
	ChainId     int64 `json:"chain_id" gorm:"primaryKey,int,not null"`
	BlockNumber int64 `json:"block_number" gorm:"int8,not null"`
}
