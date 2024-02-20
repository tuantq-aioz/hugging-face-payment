package db

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/vangxitrum/payment-host/internal/models"
)

type EntityRepository struct {
	db *gorm.DB
}

func MustNewEntityRepository(db *gorm.DB, init bool) models.EntityRepository {
	if init {
		if err := db.AutoMigrate(&models.Entity{}); err != nil {
			panic(err)
		}
	}

	return &EntityRepository{
		db: db,
	}
}

func (r EntityRepository) Create(ctx context.Context, entity *models.Entity) error {
	if err := r.db.WithContext(ctx).
		Create(entity).Error; err != nil {
		return err
	}

	return nil
}

func (r EntityRepository) GetEntityById(ctx context.Context, id uuid.UUID) (*models.Entity, error) {
	var rs models.Entity
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&rs).Error; err != nil {
		return nil, err
	}

	return &rs, nil
}

func (r EntityRepository) GetEntityByName(
	ctx context.Context,
	name string,
) (*models.Entity, error) {
	var rs models.Entity
	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&rs).Error; err != nil {
		return nil, err
	}

	return &rs, nil
}

func (r EntityRepository) GetEntityByWalletAddress(
	ctx context.Context,
	walletAddress string,
) (*models.Entity, error) {
	var rs models.Entity
	if err := r.db.WithContext(ctx).
		Model(models.Entity{}).
		Where("wallet_address = ?", walletAddress).
		First(&rs).Error; err != nil {
		return nil, err
	}

	return &rs, nil
}

func (r EntityRepository) DeleteEntityById(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Entity{}).Error; err != nil {
		return err
	}

	return nil
}
