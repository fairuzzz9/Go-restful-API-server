package db

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	onlyOnce         sync.Once
	dbPool           *pgxpool.Pool
	SQLStatementsMap = &sync.Map{}
)

func init() {
	InitSQLStatements(sqlStatements)
}

func DBPool() *pgxpool.Pool {

	if dbPool == nil {
		log.Fatal("unable to get connection pool. Already run InitDatabase() ?")
	}
	return dbPool
}

func SetupDatabase() *pgxpool.Pool {

	var DB_SERVER = "localhost"
	var DB_PORT = 5432
	var DB_USER = "myusername"
	var DB_PASSWORD = "mypassword"
	var DB_NAME = "postgres"

	dbPool := InitDatabase(DB_USER, DB_PASSWORD, DB_SERVER, DB_NAME, uint16(DB_PORT))

	err := dbPool.Ping(context.Background())
	if err != nil {
		log.Fatal("db pool ping error")
	}

	return dbPool
}

// InitDatabase initialize the database connections pool
func InitDatabase(username, password, host, database string, port uint16) *pgxpool.Pool {

	if host == "" {
		log.Panic("database host name cannot be empty string")
	}

	// hardcode for now
	// for production, read from Kubernetes config map
	var maxConns = int32(25)
	var minConns = int32(10) // if minConns set to X number. zerolog will show X lines. One line for each connection.
	var maxConnLifetime = time.Second * 30
	var maxConnIdleTime = time.Second * 1
	var healthCheckPeriod = time.Minute // 1 minute

	onlyOnce.Do(func() {

		escapedPassword := url.QueryEscape(password)

		connString := "postgres://" + username + ":" + escapedPassword + "@" + host + ":" + fmt.Sprint(port) + "/" + database

		u, err := url.Parse(connString)
		if err != nil {
			log.Fatal("unable to parse database connection string " + err.Error())
		}

		u.User = url.UserPassword(username, escapedPassword)
		config, err := pgxpool.ParseConfig(u.String())

		if err != nil {
			log.Fatal("unable to parse database connection string " + err.Error())
		}

		config.ConnConfig.Host = host
		config.ConnConfig.Port = port
		config.ConnConfig.User = username
		config.ConnConfig.Password = password
		config.ConnConfig.Database = database
		config.ConnConfig.LogLevel = pgx.LogLevelError

		config.MaxConns = maxConns
		config.MinConns = minConns
		config.MaxConnLifetime = maxConnLifetime
		config.MaxConnIdleTime = maxConnIdleTime
		config.HealthCheckPeriod = healthCheckPeriod

		pool, err := pgxpool.ConnectConfig(context.Background(), config)

		if err != nil {
			// if cannot establish database connection, then the entire program
			// will not run properly
			log.Fatal("couldn't connect to postgres " + err.Error())

		}

		pool.Config()

		dbPool = pool

	})

	return dbPool

}

var sqlStatements = map[string]string{
	// SQL Name : SQL statement
	// ok to have duplicate SQL statements, but not ok to have duplicate SQL names
	// Using map ensure that no duplicate keys. The compiler will stop if there's any duplicate in map literal

	"GetAllFromCities":     "SELECT city_name ,country_id FROM cities",
	"GetCitiesByCountryID": "SELECT city_name,country_id FROM cities WHERE country_id = $1",
	"CreateCity":           "INSERT INTO cities( city_name, country_id)VALUES($1, $2)",
	"UpdateCity":           "UPDATE cities SET city_name=$1 WHERE city_id=$2",
	"DeleteCity":           "DELETE FROM cities WHERE city_name=$1",
}

// InitSQLStatements initialize and populate the SQL statements map.
func InitSQLStatements(sqlStatements map[string]string) {

	// clear out old elements in SQLStatementsMap if any exist
	SQLStatementsMap.Range(func(key, value interface{}) bool {
		SQLStatementsMap.Delete(key)
		return true
	})

	for sqlName, sqlStatement := range sqlStatements {
		// load into map
		SQLStatementsMap.Store(sqlName, sqlStatement)
	}

}

// GetSQLByName returns the SQL statement for a given sqlName
func GetSQLByName(sqlName string) string {
	if sqlStatement, ok := SQLStatementsMap.Load(sqlName); ok {
		return sqlStatement.(string)
	}
	return ""
}
