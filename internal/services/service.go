package services

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/vangxitrum/payment-host/internal/models"
)

type EntityService interface {
	Register(ctx context.Context, name string) (*models.Entity, error)
	Withdraw(ctx context.Context, entityName string, amount decimal.Decimal, receiverAddress common.Address) (string,error)
	WatchTransaction(ctx context.Context) error
}
