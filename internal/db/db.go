package db

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/EggMD/EggMD/internal/conf"
)

// Init 连接数据库。
func Init() error {
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
		return errors.Wrap(err, "connect database")
	}

	err = db.AutoMigrate(&User{}, &Document{})
	if err != nil {
		return errors.Wrap(err, "migrate tables")
	}

	// 初始化数据存储，按字母顺序排序。
	Documents = &documents{DB: db}
	Users = &users{DB: db}

	return nil
}
