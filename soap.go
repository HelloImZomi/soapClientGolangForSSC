package gosoap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Params type is used to set the params in soap request
type Params map[string]string

// SoapClient return new *Client to handle the requests with the WSDL
func SoapClient(wsdl string, urlSoap string) (*Client, error) {
	_, err := url.Parse(wsdl)
	if err != nil {
		return nil, err
	}

	d, err := getWsdlDefinitions(wsdl)
	if err != nil {
		return nil, err
	}
	c := &Client{
		WSDL:        wsdl,
		URL:         urlSoap,  //strings.TrimSuffix(d.TargetNamespace, "/"),
		Definitions: d,
	}

	return c, nil
}

// Client struct hold all the informations about WSDL,
// request and response of the server
type Client struct {
	WSDL        	string
	URL         	string
	Method      	string
	EnvelopeTitle   string
	Params      	Params
	Definitions 	*wsdlDefinitions
	Body        	[]byte
	payload 		[]byte
}

// Call call's the method m with Params p
func (c *Client) Call(m string, e string, p Params) (err error) {
	c.Method = e
	c.EnvelopeTitle = m
	c.Params = p
	//c.payload, err = xml.MarshalIndent(c, "", "")
	c.payload, err = xml.Marshal(c)
	c.Method = m
	if err != nil {
		return err
	}

	b, err := c.doRequest()
	if err != nil {
		return err
	}
	//sss := string(b)
	//fmt.Println(sss)
	

	var soap SoapEnvelope
	err = xml.Unmarshal(b, &soap)

	c.Body = soap.Body.Contents

	
	return err
}

// Unmarshal get the body and unmarshal into the interface
func (c *Client) GetResponse() string {
	if len(c.Body) == 0 {
		return fmt.Sprintf("Body is empty")
	}
	sss := string(c.Body)
	indexInicio := strings.Index(sss, "<response>")

	indexFinal := strings.Index(sss, "</response>")
	sss = sss[indexInicio+10:indexFinal]
	if c.Method == "Execute" {
		if strings.Contains(sss, `status="fail"`) {
			return "error"
		} else {
			indexInicio = strings.Index(sss, "&lt;JournalNumber&gt;")
			indexFinal = strings.Index(sss, "&lt;/JournalNumber&gt;")
			sss = sss[indexInicio+21:indexFinal]
		}
	}



	return sss

	
}
func (c *Client) Unmarshal(v interface{}) error {
	if len(c.Body) == 0 {
		return fmt.Errorf("Body is empty")
	}

	var f Fault
	xml.Unmarshal(c.Body, &f)
	if f.Code != "" {
		return fmt.Errorf("[%s]: %s", f.Code, f.Description)
	}

	return xml.Unmarshal(c.Body, v)
}

// doRequest makes new request to the server using the c.Method, c.URL and the body.
// body is enveloped in Call method
func (c *Client) doRequest() ([]byte, error) {
	//c.payload = []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><SOAP-ENV:Envelope xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:ns1=\"http://systemsunion.com/connect/webservices/\"><SOAP-ENV:Body><ns1:SecurityProviderAuthenticateRequest><ns1:name>AOK</ns1:name><ns1:password></ns1:password></ns1:SecurityProviderAuthenticateRequest></SOAP-ENV:Body></SOAP-ENV:Envelope>")

//sss := string(c.payload)
//fmt.Println(sss)
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(c.payload))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req.ContentLength = int64(len(c.payload))

	req.Header.Add("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Add("Accept", "text/xml")
	req.Header.Add("SOAPAction", fmt.Sprintf("%s/%s", c.URL, c.Method))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// SoapEnvelope struct
type SoapEnvelope struct {
	XMLName struct{} `xml:"Envelope"`
	Body    SoapBody
}

// SoapBody struct
type SoapBody struct {
	XMLName  struct{} `xml:"Body"`
	Contents []byte   `xml:",innerxml"`
}
