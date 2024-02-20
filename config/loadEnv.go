package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT" validate:"required"`

	PostgresHost     string `mapstructure:"POSTGRES_HOST" required:"true"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT" required:"true"`
	PostgresDB       string `mapstructure:"POSTGRES_DB" required:"true"`
	PostgresUser     string `mapstructure:"POSTGRES_USER" required:"true"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD" required:"true"`

	PassPhrase   string `mapstructure:"PASSPHRASE" required:"true"`
	BusinessAddr string `mapstructure:"BUSINESS_ADDR" required:"true"`
	RpcUrl       string `mapstructure:"RPC_URL" required:"true"`
	EvmUrl       string `mapstructure:"EVM_URL" required:"true"`

	OathTokenBot string `mapstructure:"OATH_TOKEN_BOT" required:"true"`
	ChannelId    string `mapstructure:"CHANNEL_ID" required:"true"`
}

func LoadEnv(path, filename string) *Config {
	var config Config
	viper.SetConfigName(filename)
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}
