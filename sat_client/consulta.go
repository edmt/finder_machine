package sat_client

/*

Go code for...

curl -XPOST --data @request.xml \
-H "Content-Type: text/xml; charset=utf-8" \
-H "SOAPAction: http://tempuri.org/IConsultaCFDIService/Consulta" \
https://consultaqr.facturaelectronica.sat.gob.mx/ConsultaCFDIService.svc -i

*/

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	URL = "https://consultaqr.facturaelectronica.sat.gob.mx/ConsultaCFDIService.svc"
)

type ConsultaRequest struct {
	Emisor   string
	Receptor string
	Total    string
	Uuid     string
}

type ConsultaResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    BodyType `xml:"Body"`
}

type BodyType struct {
	XMLName          xml.Name             `xml:"Body"`
	ConsultaResponse ConsultaResponseType `xml:"ConsultaResponse"`
}

type ConsultaResponseType struct {
	XMLName        xml.Name           `xml:"ConsultaResponse"`
	ConsultaResult ConsultaResultType `xml:"ConsultaResult"`
}

type ConsultaResultType struct {
	XMLName       xml.Name `xml:"ConsultaResult"`
	CodigoEstatus string   `xml:"CodigoEstatus"`
	Estado        string   `xml:"Estado"`
}

func encode(input string) string {
	r := strings.NewReplacer("&", "&amp;", "Ñ", "&ntilde;", "ñ", "&ntilde;", "\xd1", "&ntilde;")
	return r.Replace(input)
}

func (r ConsultaRequest) query() string {
	return fmt.Sprintf("&quot;?re=%s&amp;rr=%s&amp;tt=%s&amp;id=%s",
		encode(r.Emisor), encode(r.Receptor), r.Total, r.Uuid)
}

func (r ConsultaRequest) template() string {
	return `
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/">
		   <soapenv:Header/>
		   <soapenv:Body>
		      <tem:Consulta>
		         <!--Optional:-->
		         <tem:expresionImpresa>` + r.query() + `</tem:expresionImpresa>
		      </tem:Consulta>
		   </soapenv:Body>
		</soapenv:Envelope>
	`
}

func (r ConsultaRequest) Consulta() string {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", URL, bytes.NewBufferString(r.template()))

	request.Header.Add("content-type", "text/xml; charset=utf-8")
	request.Header.Add("SOAPAction", "http://tempuri.org/IConsultaCFDIService/Consulta")

	response, _ := client.Do(request)
	defer response.Body.Close()

	satResponse, _ := ioutil.ReadAll(response.Body)

	var res ConsultaResponse
	xml.Unmarshal(satResponse, &res)

	return res.Body.ConsultaResponse.ConsultaResult.CodigoEstatus
}
