package main

import (
	"github.com/docopt/docopt-go"
	"github.com/edmt/finder_machine/consumers"
	"github.com/edmt/finder_machine/producers"
	// "github.com/edmt/finder_machine/sat_client"
	l4g "github.com/edmt/log4go"
	"os"
	"time"
)

const LOG_CONFIGURATION_FILE = "config/logging.xml"

func init() {
	l4g.LoadConfiguration(LOG_CONFIGURATION_FILE)
}

func main() {
	// l4g.Debug(sat_client.ConsultaRequest{"ATA980601E90", "SP&040526HD3", "45.45", "5adaa059-391e-4461-8d65-87647de235bc"}.Consulta())
	usage := `
	Usage:
	  finder_machine requeue [--start-date=<start_date>] [--end-date=<end_date>]
	  finder_machine -h | --help
	  finder_machine -v | --version

	Options:
	  -h --help     Show this screen.
	  -v --version  Show version.`

	options, _ := docopt.Parse(usage, nil, true, "0.0.1", false)
	l4g.Debug(options)
	l4g.Info("Process ID: %d", os.Getpid())

	if options["requeue"].(bool) {
		consumers.WritePool(
			consumers.XmlToPool(
				producers.ReadXML(options)))
	}

	l4g.Info("Process stopped")
	time.Sleep(time.Millisecond)
}
