package common

import (
	"database/sql"
	"errors"
	"go-transaction/config"
	"reflect"
	"strconv"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/gofiber/fiber/v2/log"
	migrate "github.com/rubenv/sql-migrate"
)

func MigrateSchema(db *sql.DB, pathFile string, schemaName string) error {
	//class := "[MigrateSql.go,MigrateSchema]"
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations_"+schemaName, pathFile),
	}
	if db == nil {
		return errors.New("error because db is null")
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	if SQLMigrationResolutionDir == "" {
		box := reflect.Indirect(reflect.ValueOf(migrations)).FieldByName("Box")
		resolution := reflect.Indirect(box).Interface().(*packr.Box)
		splitData := strings.Split(resolution.ResolutionDir, "\\")
		SQLMigrationResolutionDir = strings.Join(splitData[0:len(splitData)-1], "\\")
	}
	logModel := LoggerModel{
		Status: 200,
		//Class:    class,
		Message:  "Applied " + strconv.Itoa(n) + " migrations!",
		Version:  config.ApplicationConfiguration.GetServerConfig().Version,
		Resource: config.ApplicationConfiguration.GetServerConfig().ResourceID,
	}

	log.Info(GenerateLogModel(logModel))

	return nil
}
