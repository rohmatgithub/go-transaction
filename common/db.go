package common

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func ConnectDB(address, defaultSchema string) (gormDB *gorm.DB, err error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Adjust the output writer as needed
		logger.Config{
			SlowThreshold: time.Second, // Set the threshold for slow query logging
			LogLevel:      logger.Info, // Set the log level to Info to log queries
			Colorful:      true,        // Enable colored output
		})
	gormDB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: address + fmt.Sprintf(" search_path=%s", defaultSchema),
		//PreferSimpleProtocol: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "oauth.", // schema name
			//SingularTable: false,
		},
		Logger: newLogger,
	})
	if err != nil {
		return
	}

	//ConnectionDB, err = GormDB.DB()
	return
}
