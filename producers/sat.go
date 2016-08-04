package producers

import (
	"github.com/edmt/finder_machine/sat_client"
)

type CfdiRecord struct {
	Xml       XmlRecord
	Cfdi      CfdiType
	SatStatus string
}

func SatValidation(in <-chan XmlRecord) <-chan CfdiRecord {
	out := make(chan CfdiRecord)
	go func() {
		for record := range in {
			cfdi := ParseXml([]byte(record.Xml))
			status := sat_client.ConsultaRequest{
				cfdi.Emisor.RFC,
				cfdi.Receptor.RFC,
				cfdi.Total,
				cfdi.Complemento.TimbreFiscalDigital.UUID,
			}.Consulta()
			out <- CfdiRecord{record, cfdi, status}
		}
		close(out)
	}()
	return out
}
