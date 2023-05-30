package services

import (
	"serviceB/db"
	"serviceB/internal/models"
)

type userBalance struct {
	db db.IDbManager
}
type IUserBalance interface {
	DecreaseBalance(message *models.KafkaMessage) (err error)
}

func NewUserBalance(db db.IDbManager) IUserBalance {
	return &userBalance{
		db: db,
	}
}
