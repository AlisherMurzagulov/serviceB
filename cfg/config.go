package cfg

import (
	"encoding/json"
	"os"
	"time"
)

// Configs - корневая модель конфигурации
type Configs struct {
	Kafka                       KafkaConfig
	DB                          DatabaseSettings
	ProcessingFinishedTimeoutMS int
}

// KafkaConfig - класс описывающий конфигурацию подключения для kafka
type KafkaConfig struct {
	BrokerURLs    []string
	DialerTimeout time.Duration
	MinBytes      int
	MaxBytes      int
	ReadTimeout   time.Duration
	Verbose       bool
	Topic         string
	ConsumerGroup string
	SASL          struct {
		User     string `json:"User"`
		Password string `json:"Password"`
	} `json:"SASL"`
}

// DatabaseSettings Настройки базы данных
type DatabaseSettings struct {
	ConnectionString string
	LogMode          bool
	MigrationMode    bool
}

// AppConfigs - конфиг приложение
var AppConfigs = &Configs{}

// Setup - инициализация конфигов
func Setup(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic("Config file reading error")
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic("Config file close error")
		}
	}(file)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfigs)
	if err != nil {
		panic("Config file unmarshalling error")
	}

}
