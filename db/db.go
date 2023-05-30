package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"serviceB/cfg"
	"serviceB/db/entity"
	"serviceB/db/repositories"
)

// DBManager для связывания методов БД
type dbManager struct {
	DB *gorm.DB
}

type IDbManager interface {
	IUserRepository() repositories.IUserBalanceRepository
}

// InitDB инициализация БД
func (dbm *dbManager) initDB() {
	var err error
	dbm.DB, err = gorm.Open("postgres", cfg.AppConfigs.DB.ConnectionString)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	dbm.DB.LogMode(cfg.AppConfigs.DB.LogMode)

	if cfg.AppConfigs.DB.MigrationMode {
		dbm.DB.AutoMigrate(&entity.UserBalance{})
	}
}
func (dbm *dbManager) IUserRepository() repositories.IUserBalanceRepository {
	return repositories.NewUserBalanceRepository(dbm.DB)
}

// NewDBManager конструктор
func NewDBManager() IDbManager {
	var dbm dbManager
	dbm.initDB()
	return &dbm
}
