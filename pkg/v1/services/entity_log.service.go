package services

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/vangxitrum/payment-host/internal/models"
	internal_services "github.com/vangxitrum/payment-host/internal/services"
	"github.com/vangxitrum/payment-host/internal/utils"
)

type LogFunc func(time.Time, string, error)

type EntityLogService struct {
	next    internal_services.EntityService
	logFunc LogFunc
}

func NewEntityLogService(next internal_services.EntityService) internal_services.EntityService {
	return &EntityLogService{
		next:    next,
		logFunc: utils.Log,
	}
}

func (s *EntityLogService) Register(ctx context.Context, name string) (entity *models.Entity, err error) {
	defer func(start time.Time) {
		s.logFunc(start, "Register", err)
	}(time.Now().UTC())

	return s.next.Register(ctx, name)
}

func (s *EntityLogService) Withdraw(ctx context.Context, entityName string, amount decimal.Decimal, receiverAddress common.Address) (txHash string, err error) {
	defer func(start time.Time) {
		s.logFunc(start, "Withdraw", err)
	}(time.Now().UTC())

	return s.next.Withdraw(ctx, entityName, amount, receiverAddress)
}

func (s *EntityLogService) WatchTransaction(ctx context.Context) (err error) {
	defer func(start time.Time) {
		s.logFunc(start, "WatchTransaction", err)
	}(time.Now().UTC())

	return s.next.WatchTransaction(ctx)
}
