package common

import (
	"database/sql"
	"fmt"
	"go-transaction/config"
	"io"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ConnectionDB              *sql.DB
	GormDB                    *gorm.DB
	SQLMigrationResolutionDir string
	RedisClient               *redis.Client
	Validation                ValidationInterface
	ErrorBundle               *i18n.Bundle
	ConstantaBundle           *i18n.Bundle
	CommonBundle              *i18n.Bundle
)

//type logWriter struct {
//}
//
//func (writer logWriter) Write(bytes []byte) (int, error) {
//	return fmt.Print(time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + string(bytes))
//	//return os.Stdout.Write([]byte(time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + string(bytes)))
//}

func SetServerAttribute() error {
	var err error
	GormDB, err = ConnectDB(config.ApplicationConfiguration.GetPostgresqlConfig().Address, config.ApplicationConfiguration.GetPostgresqlConfig().DefaultSchema)
	if err != nil {
		return err
	}

	ConnectionDB, err = GormDB.DB()
	if err != nil {
		return err
	}
	// set log fiber
	// Output to ./test.log file
	file, _ := os.OpenFile("fiber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	iw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(iw)

	// CONNECT REDIS
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.ApplicationConfiguration.GetRedisConfig().Host, config.ApplicationConfiguration.GetRedisConfig().Port),
		Password: config.ApplicationConfiguration.GetRedisConfig().Password, // no password set
		DB:       config.ApplicationConfiguration.GetRedisConfig().Db,       // use default DB
	})
	return nil
}
