package db

import (
	"context"

	"github.com/vangxitrum/payment-host/internal/models"
	"gorm.io/gorm"
)

type PaymentMarkRepository struct {
	db *gorm.DB
}

func MustNewPaymentMarkRepository(db *gorm.DB, init bool) models.PaymentMarkRepository {
	if init {
		if err := db.AutoMigrate(&models.PaymentMark{}); err != nil {
			panic(err)
		}
	}

	return &PaymentMarkRepository{
		db: db,
	}
}

func (r PaymentMarkRepository) Create(ctx context.Context, paymentMark *models.PaymentMark) error {
	return r.db.WithContext(ctx).Create(paymentMark).Error
}

func (r PaymentMarkRepository) GetPaymentMarkByChainId(ctx context.Context, chainId int64) (*models.PaymentMark, error) {
	var rs models.PaymentMark
	if err := r.db.WithContext(ctx).
		Where("chain_id = ?", chainId).
		First(&rs).Error; err != nil {
		return nil, err
	}

	return &rs, nil
}

func (r PaymentMarkRepository) UpdatePaymentMarkByChainId(ctx context.Context, chainId int64, blockNumber int64) error {
	if err := r.db.WithContext(ctx).
		Model(models.PaymentMark{}).
		Where("chain_id = ?", chainId).
		Update("block_number", blockNumber).Error; err != nil {
		return err
	}

	return nil
}

func (r PaymentMarkRepository) DeletePaymentMarkByChainId(ctx context.Context, chainId int64) error {
	if err := r.db.WithContext(ctx).
		Where("chain_id = ?", chainId).
		Delete(&models.PaymentMark{}).Error; err != nil {
		return err
	}

	return nil
}
