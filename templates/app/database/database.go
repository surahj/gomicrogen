package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func DbInstance() *sql.DB {

	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&multiStatements=true", username, password, host, port, dbname, "utf8")

	Db, err := otelsql.Open("mysql", dbURI,
		otelsql.WithAttributes(semconv.DBSystemMySQL),
		otelsql.WithDBName(dbname))


	checkErr(err)

	otelsql.ReportDBStatsMetrics(Db)

	//Db, err := sql.Open("mysql", dbURI)

	//checkErr(err)

	idleConnection := os.Getenv("DATABASE_IDLE_CONNECTION")
	ic, err := strconv.Atoi(idleConnection)

	if err != nil {

		ic = 5
	}

	maxConnection := os.Getenv("DATABASE_MAX_CONNECTION")

	mx, err := strconv.Atoi(maxConnection)

	if err != nil {

		mx = 10
	}

	connectionLifetime := os.Getenv("DATABASE_CONNECTION_LIFETIME")

	cl, err := strconv.Atoi(connectionLifetime)

	if err != nil {

		cl = 60
	}

	Db.SetMaxIdleConns(ic)
	Db.SetConnMaxLifetime(time.Second * time.Duration(cl))
	Db.SetMaxOpenConns(mx)
	Db.SetConnMaxIdleTime(time.Second * time.Duration(cl))

	err = Db.Ping()
	checkErr(err)
	return Db
}

func checkErr(err error) {

	if err != nil {

		fmt.Println("db connection error", err)
		log.Printf("DB ERROR %s ", err.Error())
	}
}
