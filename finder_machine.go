package main

import (
	"github.com/docopt/docopt-go"
	cfdi "github.com/edmt/finder_machine/consumers/cfdi"
	pool "github.com/edmt/finder_machine/consumers/pool"
	"github.com/edmt/finder_machine/producers"
	l4g "github.com/edmt/log4go"
	"os"
	"time"
)

const LOG_CONFIGURATION_FILE = "config/logging.xml"

func init() {
	l4g.LoadConfiguration(LOG_CONFIGURATION_FILE)
}

func main() {
	usage := `
	Usage:
	  finder_machine requeue [--start-date=<start_date>] [--end-date=<end_date>]
	  finder_machine reprocess [--start-date=<start_date>] [--end-date=<end_date>]
	  finder_machine -h | --help
	  finder_machine -v | --version

	Options:
	  -h --help     Show this screen.
	  -v --version  Show version.`

	options, _ := docopt.Parse(usage, nil, true, "0.0.1", false)
	l4g.Debug(options)
	l4g.Info("Process ID: %d", os.Getpid())

	if options["requeue"].(bool) {
		pool.WritePool(
			pool.XmlToPool(
				producers.ReadXML(options)))
	}
	if options["reprocess"].(bool) {
		cfdi.WriteCfdi(
			producers.SatValidation(
				producers.ReadXMLMissingCfd(options)))
	}

	l4g.Info("Process stopped")
	time.Sleep(time.Millisecond)
}
