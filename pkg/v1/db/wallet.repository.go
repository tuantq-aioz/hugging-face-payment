package db

import (
	"context"

	"github.com/vangxitrum/payment-host/internal/models"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func MustNewWalletRepository(db *gorm.DB, init bool) models.WalletRepository {
	if init {
		if err := db.AutoMigrate(&models.Wallet{}); err != nil {
			panic(err)
		}
	}

	return &WalletRepository{
		db: db,
	}
}

func (r WalletRepository) Create(ctx context.Context, wallet *models.Wallet) error {
	if err := r.db.WithContext(ctx).Create(wallet).Error; err != nil {
		return nil
	}

	return nil
}

func (r WalletRepository) GetActiveWallets(ctx context.Context) ([]*models.Wallet, error) {
	var rs []*models.Wallet
	if err := r.db.WithContext(ctx).Model(models.Wallet{}).Find(&rs).Error; err != nil {
		return nil, err
	}

	return rs, nil
}
