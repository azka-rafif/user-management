package configs

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config config struct consist of data that provided from env
type Config struct {
	LogLevel       string `mapstructure:"LOG_LEVEL"`
	Port           string `mapstructure:"PORT"`
	ShutdownPeriod int64  `mapstructure:"SHUTDOWN_PERIOD"`

	// Mysql
	WriteDatabaseHost     string `mapstructure:"WRITE_DB_HOST"`
	WriteDatabaseUsername string `mapstructure:"WRITE_DB_USER"`
	WriteDatabasePassword string `mapstructure:"WRITE_DB_PASSWORD"`
	WriteDatabaseName     string `mapstructure:"WRITE_DB_NAME"`
	WriteDatabaseTimeZone string `mapstructure:"WRITE_DB_TIME_ZONE"`

	ReadDatabaseHost     string `mapstructure:"READ_DB_HOST"`
	ReadDatabaseUsername string `mapstructure:"READ_DB_USER"`
	ReadDatabasePassword string `mapstructure:"READ_DB_PASSWORD"`
	ReadDatabaseName     string `mapstructure:"READ_DB_NAME"`
	ReadDatabaseTimeZone string `mapstructure:"READ_DB_TIME_ZONE"`

	// Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	// Message Broker
	NsqHost               string `mapstructure:"NSQ_HOST"`
	NsqPort               string `mapstructure:"NSQ_PORT"`
	EnableExampleConsumer bool   `mapstructure:"ENABLE_EXAMPLE_CONSUMER"`

	// Consumer
	ConsumerBackoffMaxRetry uint64 `mapstructure:"CONSUMER_BACKOFF_MAX_RETRY"`

	// Producer
	ProducerBackoffMaxRetry int `mapstructure:"PRODUCER_BACKOFF_MAX_RETRY"`
	ProducerRetryBackoff    int `mapstructure:"PRODUCER_RETRY_BACKOFF"`
	ProducerDialTimeout     int `mapstructure:"PRODUCER_DIAL_TIMEOUT"`

	// APP
	AppURL string `mapstructure:"APP_URL"`
	Env    string `mapstructure:"ENV"`
}

var conf Config

// Get are responsible to load env and get data an return the struct
func Get() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed reading config file")
	}

	once := sync.Once{}
	once.Do(func() {
		err = viper.Unmarshal(&conf)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
	})

	return &conf
}
