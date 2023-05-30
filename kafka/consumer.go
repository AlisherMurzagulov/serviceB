package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	kafgo "github.com/segmentio/kafka-go"
	"serviceB/cfg"
	"serviceB/internal/models"
	"serviceB/internal/services"
	"sync"
	"time"
)

type Consumer struct {
	cfg                *cfg.KafkaConfig
	reader             *kafgo.Reader
	wg                 sync.WaitGroup
	stop               bool
	userBalanceService services.IUserBalance
}

// Stop for loop for kafka reading
func (c *Consumer) Stop() {
	c.stop = true
}

// Close kafka reader
func (c *Consumer) Close() {
	err := c.reader.Close()
	if err != nil {
		return
	}
	fmt.Println("Закрытие сессии Consumer-а")
}

// readMessage
func (c *Consumer) readMessage() (*models.KafkaMessage, error) {

	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ReadTimeout*time.Second)
	defer cancel()
	m, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	kafkaMessage := new(models.KafkaMessage)
	err = json.Unmarshal(m.Value, &kafkaMessage)
	if err != nil {
		fmt.Printf("Ошибка парсинга сообщения err: %v", err.Error())
		return nil, err
	}

	return kafkaMessage, err
}

// startProcessing start processing kafka messages to core
func (c *Consumer) startProcessing(kafkaMessage *models.KafkaMessage) error {
	// Говорим, что завершаем обработку
	defer c.wg.Done()

	// Говорим, что начали обработку
	c.wg.Add(1)
	return c.userBalanceService.DecreaseBalance(kafkaMessage)
}

// StartReading start for loop to read kafka messages
func (c *Consumer) StartReading() {
	go func() {
		for {
			if c.stop {
				break
			}
			kafkaMessage, err := c.readMessage()
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("Таймаут ожидание сообщение с kafka")
				continue
			}

			if err != nil {
				fmt.Printf("Ошибка при чтении из kafka")
				continue
			}
			fmt.Printf("Получили сообщение с  kafka")
			err = c.startProcessing(kafkaMessage)
			if err != nil {
				fmt.Printf("Ошибка обработки сообщение и отправки в core")
			}
		}
	}()
}

// WaitExit wait all process for closing or ends with timer
func (c *Consumer) WaitExit() {

	processingFinished := make(chan bool, 1)
	// Запускаем go-рутину ожидания завершения всех процессов
	go func() {
		c.wg.Wait()
		// Если все обработки завершились, то выкидываем в канал сообщение
		processingFinished <- true
	}()

	// Создаем таймер на ожидание завершения всех обработок
	processingFinishedTimer := time.NewTimer(time.Second * time.Duration(cfg.AppConfigs.ProcessingFinishedTimeoutMS))

	// Ждем либо обработки всех сообщений, либо таймаута ожидания обработки
	select {
	case <-processingFinished:
		{
			fmt.Printf("Все текущие процессы завершены, останавливаем приложение штатно")
		}
	case <-processingFinishedTimer.C:
		{
			fmt.Printf("\"Не все текущие процессы завершены, останавливаем приложение по таймауту [%d]", cfg.AppConfigs.ProcessingFinishedTimeoutMS)
		}
	}
}
