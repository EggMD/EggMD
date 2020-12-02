package db

import (
	"fmt"

	"github.com/EggMD/EggMD/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"
)

func Init() (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Name,
	)
	db, err := gorm.Open(mysql.Open(dns),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}

	// Initialize stores, sorted in alphabetical order.
	Users = &users{DB: db}

	return db, nil
}
