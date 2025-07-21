package infrastructures

import (
	"fmt"
	"notification-services/config"

	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *PostgresDatabase
)

func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s timezone=%s",
			conf.Db.Host,
			conf.Db.User,
			conf.Db.Password,
			conf.Db.DbName,
			conf.Db.Port,
			conf.Db.SslMode,
			conf.Db.TimeZone,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to database: %v", err))
		}

		dbInstance = &PostgresDatabase{
			db: db,
		}
		fmt.Println("Notification database connected successfully.")
	})

	return dbInstance
}

func (p *PostgresDatabase) GetInstance() *gorm.DB {
	return dbInstance.db
}
