package cfdi

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/edmt/finder_machine/producers"
	l4g "github.com/edmt/log4go"
	"github.com/olebedev/config"
	"log"
	"os"
	"strings"
)

func WriteCfdi(in <-chan producers.CfdiRecord) {
	l4g.Info("consumers/cfdi: loading configuration")

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

	l4g.Info("consumers/cfdi: connecting to database")

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

func consumeChannel(db *sql.DB, in <-chan producers.CfdiRecord) {
	l4g.Info("consumers/cfdi: ready to write to database")

	stmt, err := db.Prepare("exec FinderMachine_RecoverDeletedCfd ?")
	if err != nil {
		log.Fatal(err)
	}

	for record := range in {
		// Success
		if strings.Contains(record.SatStatus, "satisfactoriamente") {
			l4g.Info("consumers/cfdi: recovering cfdi from cfd_delete with uuid: %s", record.Xml.Uuid)
			_, err = stmt.Exec(record.Xml.Uuid)

			if err != nil {
				log.Fatal(err)
			}

		}

		// 602: not found
		if strings.Contains(record.SatStatus, "602") {
			l4g.Info("consumers/cfdi: this is supposed to write to cfd...")

			// Append to missing.log file
			missingFile, _ := os.OpenFile("./tmp/missing.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
			missingFile.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\n", record.Cfdi.Emisor.RFC, record.Cfdi.Receptor.RFC, record.Xml.Uuid, record.Cfdi.Fecha))
			defer missingFile.Close()

		}

		// 601: bad request
	}
}
