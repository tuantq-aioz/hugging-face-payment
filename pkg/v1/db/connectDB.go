package db

import (
	"fmt"
	"log"
	"os"
	"time"

	appConfig "github.com/vangxitrum/payment-host/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func MustConnectDB(config *appConfig.Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.PostgresHost,
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresDB,
		config.PostgresPort)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             1000 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false, // Ignore ErrRecordNotFound error for logger
			// ParameterizedQueries:      true, // Don't include params in the SQL log
			Colorful: true, // Enable color
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect to the database!")
	}

	fmt.Println("🚀 Connected Successfully to the database.")
}
