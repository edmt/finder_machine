package producers

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	l4g "github.com/edmt/log4go"
	"github.com/olebedev/config"
	"log"
	"time"
)

func ReadXML(options map[string]interface{}) <-chan XmlRecord {
	l4g.Info("producers/xml: loading configuration")

	cfg, err := config.ParseYamlFile("config/finder_machine.yaml")
	if err != nil {
		log.Fatal(err)
	}

	host, _ := cfg.String("xml.database.host")
	port, _ := cfg.String("xml.database.port")
	user, _ := cfg.String("xml.database.user")
	password, _ := cfg.String("xml.database.password")
	database, _ := cfg.String("xml.database.database")

	l4g.Debug("host:%s port:%s user:%s password:%s database:%s",
		host, port, user, password, database)

	l4g.Info("producers/xml: connecting to database")

	connection := connectionParameters{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}.connect()

	// defer connection.Close()

	ping(connection)

	return produceChannel(connection, options)
}

type connectionParameters struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (c connectionParameters) makeConnectionString() string {
	return fmt.Sprintf(
		"server=%s;port=%s;user id=%s;password=%s;database=%s;log=3;encrypt=disable",
		c.Host, c.Port, c.User, c.Password, c.Database)
}

func (c connectionParameters) connect() *sql.DB {
	database, err := sql.Open("mssql", c.makeConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	return database
}

func ping(db *sql.DB) {
	err := db.Ping()
	if err == nil {
		l4g.Info("producers/xml: pinging database")
	} else {
		l4g.Error(err)
	}
}

func produceChannel(db *sql.DB, options map[string]interface{}) <-chan XmlRecord {
	l4g.Info("producers/xml: querying database")
	out := make(chan XmlRecord)
	go func() {
		startDate := options["--start-date"]
		if startDate == nil {
			startDate = time.Now().AddDate(0, 0, -1).Local().Format("2006-01-02")
		}

		endDate := options["--end-date"]
		if endDate == nil {
			endDate = time.Now().Local().Format("2006-01-02")
		}

		l4g.Debug("StartDate: %s", startDate)
		l4g.Debug("EndDate: %s", endDate)

		var (
			uuid      string
			xml       string
			timestamp time.Time
		)

		rows, err := db.Query("exec FinderMachine_ReadXML ?, ?", startDate, endDate)

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&uuid, &xml, &timestamp)

			if err != nil {
				log.Fatal(err)
			}

			out <- XmlRecord{uuid, xml, timestamp}
		}
		close(out)

		err = rows.Err()

		if err != nil {
			log.Fatal(err)
		}
	}()
	return out
}

type XmlRecord struct {
	Uuid      string
	Xml       string
	Timestamp time.Time
}
