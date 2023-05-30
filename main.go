package main

import (
	"os"
	"os/signal"
	"serviceB/cfg"
	"serviceB/db"
	"serviceB/internal/services"
	"serviceB/kafka"
	"syscall"
)

func init() {
	cfg.Setup("config.json")
}

func main() {
	// Подписываемся на получение сигналов от операционной системы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	userBalanceService := services.NewUserBalance(db.NewDBManager())

	kafkaConsumer, errKafkaInit := kafka.NewConsumer(&cfg.AppConfigs.Kafka, userBalanceService)
	if errKafkaInit != nil {
		panic(errKafkaInit)
	}

	defer kafkaConsumer.Close()

	kafkaConsumer.StartReading()

	// Говорим, что не надо больше обрабатывать сообщения новые
	kafkaConsumer.Stop()
	kafkaConsumer.WaitExit()

}
