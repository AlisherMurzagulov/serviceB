package repositories

import (
	"github.com/jinzhu/gorm"
	"serviceB/db/entity"
)

type userBalanceRepository struct {
	db *gorm.DB
}

type IUserBalanceRepository interface {
	GetBalance(userId string) (balance *entity.UserBalance, err error)
	UpdateBalance(balance entity.UserBalance) (err error)
}

func NewUserBalanceRepository(db *gorm.DB) IUserBalanceRepository {
	return &userBalanceRepository{
		db: db,
	}
}
