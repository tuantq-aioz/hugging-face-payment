package models

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	CONTRACT_IN_TYPE  = "in"
	CONTRACT_OUT_TYPE = "out"

	AIOZ_CONTRACT_ADDRESS = "aioz"

	TX_STATUS_NEW     = "new"
	TX_STATUS_HANDLED = "handled"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error

	GetTransactionByHashIndexAndReceiverAddr(ctx context.Context, hash string, index int, recvAddr string) (*Transaction, error)
}

type Transaction struct {
	Id              uuid.UUID       `json:"id" gorm:"primary_key,type:uuid"`
	EntityId        uuid.UUID       `json:"entity_id" gorm:"type:uuid"`
	CosmosHash      string          `json:"cosmos_hash" gorm:"text"`
	EvmHash         string          `json:"evm_hash" gorm:"text"`
	ContractAddress string          `json:"contract_address" gorm:"text"`
	From            string          `json:"from" gorm:"text"`
	To              string          `json:"to" gorm:"text"`
	BlockNumber     uint64          `json:"block_number" gorm:"int8"`
	Type            string          `json:"type" gorm:"text"`
	Index           int             `json:"index" gorm:"int4"`
	Denom           string          `json:"denom" gorm:"varchar(255)"`
	Amount          decimal.Decimal `json:"amount" gorm:"type:numeric"`
	Credit          decimal.Decimal `json:"credit" gorm:"type:numeric"`
	Status          string          `json:"status" gorm:"text"`
	CreatedAt       int64           `json:"created_at" gorm:"int8,not null"`
	UpdatedAt       int64           `json:"updated_at" gorm:"int8,not null"`
}

func ParseCoinAmount(amountValue string) (decimal.Decimal, string, error) {
	re := regexp.MustCompile(`^(\d+)([a-z]+)$`)
	match := re.FindStringSubmatch(amountValue)
	if len(match) != 3 {
		return decimal.Zero, "", fmt.Errorf("invalid amount string")
	}
	amount, err := decimal.NewFromString(match[1])
	if err != nil {
		return decimal.Zero, "", err
	}
	return amount, match[2], nil
}
