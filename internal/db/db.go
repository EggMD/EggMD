package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/conf"
)

// Init connects to the database.
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

	err = db.AutoMigrate(&User{}, &Document{})
	if err != nil {
		log.Fatal("Failed to auto migrate tables: %v", err)
	}

	// Initialize stores, sorted in alphabetical order.
	Documents = &documents{DB: db}
	Users = &users{DB: db}

	return db, nil
}
