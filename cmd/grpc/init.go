package main

import (
	"github.com/vangxitrum/payment-host/config"
	"github.com/vangxitrum/payment-host/internal/crons"
	"github.com/vangxitrum/payment-host/internal/models"
	"github.com/vangxitrum/payment-host/internal/services"
	"github.com/vangxitrum/payment-host/pkg/v1/db"
	v1 "github.com/vangxitrum/payment-host/pkg/v1/services"
)

var (
	appConfig *config.Config
	cron      *crons.Cron

	entityRepo      models.EntityRepository
	paymentMarkRepo models.PaymentMarkRepository
	walletRepo      models.WalletRepository
	txRepo          models.TransactionRepository

	entityService services.EntityService
)

func init() {
	appConfig = config.LoadEnv(".", "app")
	db.MustConnectDB(appConfig)

	entityRepo = db.MustNewEntityRepository(db.DB, true)
	paymentMarkRepo = db.MustNewPaymentMarkRepository(db.DB, true)
	walletRepo = db.MustNewWalletRepository(db.DB, true)
	txRepo = db.MustNewTransactionRepository(db.DB, true)

	entityService = v1.MustNewEntityService(
		appConfig.RpcUrl,
		appConfig.EvmUrl,
		appConfig.PassPhrase,
		appConfig.BusinessAddr,

		entityRepo,
		paymentMarkRepo,
		walletRepo,
		txRepo,
	)

	entityService = v1.NewEntityLogService(entityService)
	cron = crons.NewCron(entityService)
}
