package templates

import (
	"log"
	"os"
)

var dbCode = `/*
Copyright Interstellar, Inc - All Rights Reserved.
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
Written by Fritz Ekwoge (fritz@interstellar.cm), March 2021.
*/
package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenDb method should only run once to reuse the same connection pool
// https://golang.org/pkg/database/sql/#Open
// The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections.
// Thus, the Open function should be called just once. It is rarely necessary to close a DB.

func OpenDb(db_prefix string) (*gorm.DB, error) {

	dbType := os.Getenv(fmt.Sprintf("%s_DB_TYPE", db_prefix))
	if len(dbType) == 0 {
		dbType = "postgres"
	}
	dbConnectionString := os.Getenv(fmt.Sprintf("%s_DB_CONNECTION_STRING", db_prefix))
	if len(dbConnectionString) == 0 {
		return nil, errors.New("empty connection string. Check env variable: DB_CONNECTION_STRING")
	}

	var err error

	var database *gorm.DB
	var sqlDB *sql.DB

	if dbType == "sqlite" {
		database, err = gorm.Open(sqlite.Open(dbConnectionString), &gorm.Config{
			QueryFields: true,
		})

		if err != nil {
			log.Error().Err(err).Msg("could-not-connect-to-db")
			panic(err)
		}

	}

	if dbType == "postgres" {
		var maxoconn, idleCon string
		if os.Getenv(fmt.Sprintf("%s_DB_MAX_OPEN_CONNECTIONS", db_prefix)) != "" {
			maxoconn = os.Getenv(fmt.Sprintf("%s_DB_MAX_OPEN_CONNECTIONS", db_prefix))
		} else {
			maxoconn = "50"
		}
		if os.Getenv(fmt.Sprintf("%s_DB_MAX_IDLE_CONNECTIONS", db_prefix)) != "" {
			idleCon = os.Getenv(fmt.Sprintf("%s_DB_MAX_IDLE_CONNECTIONS", db_prefix))
		} else {
			idleCon = "50"
		}
		maxOpenConn, _ := strconv.ParseInt(maxoconn, 10, 64)
		maxIdleConns, _ := strconv.ParseInt(idleCon, 10, 64)
		// connsMaxIdleTime, _ := strconv.ParseInt(idletime, 10, 64)
		sqlDB, err = sql.Open("pgx", dbConnectionString)

		if err != nil {
			log.Error().Err(err).Msg("could-not-connect-to-db")
			panic(err)
		}

		sqlDB.SetMaxIdleConns(int(maxIdleConns))
		sqlDB.SetConnMaxIdleTime(5 * time.Second)
		sqlDB.SetMaxOpenConns(int(maxOpenConn))
		sqlDB.SetConnMaxLifetime(24 * time.Hour)

		database, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{
			Logger:      logger.Default.LogMode(logger.Silent),
			QueryFields: true,
		})
		if err != nil {
			log.Error().Err(err).Msg("could-not-connect-to-db")
			return nil, err
		}
	}

	if dbType == "mysql" {
		var maxoconn, idleCon string
		if os.Getenv(fmt.Sprintf("%s_DB_MAX_OPEN_CONNECTIONS", db_prefix)) != "" {
			maxoconn = os.Getenv(fmt.Sprintf("%s_DB_MAX_OPEN_CONNECTIONS", db_prefix))
		} else {
			maxoconn = "50"
		}
		if os.Getenv(fmt.Sprintf("%s_DB_MAX_IDLE_CONNECTIONS", db_prefix)) != "" {
			idleCon = os.Getenv(fmt.Sprintf("%s_DB_MAX_IDLE_CONNECTIONS", db_prefix))
		} else {
			idleCon = "50"
		}
		maxOpenConn, _ := strconv.ParseInt(maxoconn, 10, 64)
		maxIdleConns, _ := strconv.ParseInt(idleCon, 10, 64)
		// connsMaxIdleTime, _ := strconv.ParseInt(idletime, 10, 64)
		sqlDB, err = sql.Open("mysql", dbConnectionString)

		if err != nil {
			log.Error().Err(err).Msg("could-not-connect-to-db")
			panic(err)
		}

		sqlDB.SetMaxIdleConns(int(maxIdleConns))
		sqlDB.SetConnMaxIdleTime(5 * time.Second)
		sqlDB.SetMaxOpenConns(int(maxOpenConn))
		sqlDB.SetConnMaxLifetime(24 * time.Hour)

		database, err = gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}), &gorm.Config{
			Logger:      logger.Default.LogMode(logger.Silent),
			QueryFields: true,
		})

		if err != nil {
			log.Error().Err(err).Msg("could-not-connect-to-db")
			return nil, err
		}
	}

	return database, nil
}

func PrintDBStats(tag string, db *gorm.DB) {
	sqlDB, _ := db.DB()
	dbStats := sqlDB.Stats()
	log.Printf("[%v], open connections: %v/%v, idle connections: %v/%v", tag, dbStats.OpenConnections, dbStats.MaxOpenConnections, dbStats.Idle, dbStats.MaxIdleClosed)

}
`

func DBTemplate() string {
	return dbCode
}

func CreateDB() {
	// Create the db/main.go file
	_, err := os.Create("internal/db/main.go")
	if err != nil {
		log.Fatal(err)
	}

	// Open the db/main.go file for writing
	f, err := os.OpenFile("internal/db/main.go", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

	// Write the code to the file
	_, err = f.WriteString(DBTemplate())
	if err != nil {
		log.Fatal(err)
	}

	// Close the file
	f.Close()
}
