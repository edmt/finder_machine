package consumers

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/edmt/finder_machine/producers"
	l4g "github.com/edmt/log4go"
	"github.com/olebedev/config"
	"log"
	"time"
)

func XmlToPool(in <-chan producers.XmlRecord) <-chan PoolRecord {
	out := make(chan PoolRecord)
	go func() {
		for record := range in {
			out <- PoolRecord{record.Uuid, time.Now().Local(), 0}
		}
		close(out)
	}()
	return out
}

type PoolRecord struct {
	Uuid          string
	FechaRegistro time.Time
	Status        int
}

func WritePool(in <-chan PoolRecord) {
	l4g.Info("consumers/pool: loading configuration")

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

	l4g.Info("consumers/pool: connecting to database")

	connection := connectionParameters{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}.connect()

	consumeChannel(connection, in)
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

func consumeChannel(db *sql.DB, in <-chan PoolRecord) {
	l4g.Info("consumers/pool: ready to write to database")

	stmt, err := db.Prepare("exec FinderMachine_WritePool ?")
	if err != nil {
		log.Fatal(err)
	}

	for record := range in {
		l4g.Info("consumers/pool: enqueuing cfdi with uuid: %s", record.Uuid)
		_, err = stmt.Exec(record.Uuid)

		if err != nil {
			log.Fatal(err)
		}
	}

}
