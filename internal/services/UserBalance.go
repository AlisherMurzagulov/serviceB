package services

import (
	"fmt"
	"math/rand"
	"serviceB/db/entity"
	"serviceB/internal/models"
	"time"
)

func (u userBalance) DecreaseBalance(message *models.KafkaMessage) (err error) {
	balance, err := u.db.IUserRepository().GetBalance(message.ID)
	if err != nil {
		err = fmt.Errorf("пользователь не найден")
		return
	}

	if balance.Balance >= message.Amount {
		err = u.db.IUserRepository().UpdateBalance(entity.UserBalance{
			UserId:  balance.UserId,
			Balance: balance.Balance - message.Amount,
		})
		if err != nil {
			return
		}
	}

	time.Sleep(30 * time.Second)

	if rand.Intn(2) == 0 {
		return nil
	} else {
		err = u.db.IUserRepository().UpdateBalance(entity.UserBalance{
			UserId:  balance.UserId,
			Balance: balance.Balance + message.Amount,
		})
		if err != nil {
			return
		}
		err = fmt.Errorf("операция не успешна")
	}

	return
}
