package db

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/EggMD/EggMD/internal/conf"
	"github.com/EggMD/EggMD/internal/dbutil"
)

var AllTables = []interface{}{
	&User{},
	&Team{},
	&Document{},
}

// Init initializes the database.
func Init() (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Name,
		conf.Database.SSLMode,
	)
	conf.Database.DSN = dsn

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return dbutil.Now()
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "open connection")
	}

	// Migrate databases.
	if db.AutoMigrate(AllTables...) != nil {
		return nil, errors.Wrap(err, "auto migrate")
	}

	// Create sessions table.
	q := `
CREATE TABLE IF NOT EXISTS sessions (
    key        TEXT PRIMARY KEY,
    data       BYTEA NOT NULL,
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL
);`
	if err := db.Exec(q).Error; err != nil {
		return nil, errors.Wrap(err, "create sessions table")
	}

	SetDatabaseStore(db)

	return nil, nil
}

// SetDatabaseStore sets the database table store.
func SetDatabaseStore(db *gorm.DB) {
	Users = NewUsersStore(db)
	Teams = NewTeamsStore(db)
	Documents = NewDocumentsStore(db)
}
