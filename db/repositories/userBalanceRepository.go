package repositories

import (
	"serviceB/db/entity"
)

func (dbm *userBalanceRepository) GetBalance(userId string) (balance *entity.UserBalance, err error) {
	err = dbm.db.Where(&entity.UserBalance{UserId: userId}).First(&balance).Error
	return
}

func (dbm *userBalanceRepository) UpdateBalance(balance entity.UserBalance) (err error) {
	err = dbm.db.Save(&balance).First(&balance).Error
	return
}
