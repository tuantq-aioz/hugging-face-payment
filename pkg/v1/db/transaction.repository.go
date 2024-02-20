package db

import (
	"context"

	"gorm.io/gorm"

	"github.com/vangxitrum/payment-host/internal/models"
)

type TransactionRepository struct {
	db *gorm.DB
}

func MustNewTransactionRepository(db *gorm.DB, init bool) models.TransactionRepository {
	if init {
		if err := db.AutoMigrate(&models.Transaction{}); err != nil {
			panic(err)
		}
	}

	return TransactionRepository{
		db: db,
	}
}

func (r TransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	if err := r.db.WithContext(ctx).
		Create(transaction).Error; err != nil {
		return err
	}

	return nil
}

func (r TransactionRepository) GetTransactionByHashIndexAndReceiverAddr(
	ctx context.Context,
	hash string,
	index int,
	recvAddr string,
) (*models.Transaction, error) {
	var tx models.Transaction
	if err := r.db.WithContext(ctx).
		Model(models.Transaction{}).
		Where("cosmos_hash = ? and index = ? and \"to\" = ?", hash, index, recvAddr).
		First(&tx).Error; err != nil {
		return nil, err
	}

	return &tx, nil
}
