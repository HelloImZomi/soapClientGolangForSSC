package gosoap

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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
	c.EnvelopeTitle = m
	c.Method = e
	c.Params = p
	c.payload, err = xml.Marshal(c)
	c.Method = m
	if err != nil {
		return err
	}

	b, err := c.doRequest()
	if err != nil {
		return err
	}

	var soap SoapEnvelope

	if err := xml.Unmarshal(b, &soap); err != nil {
		return fmt.Errorf("an error occurred decoding the body: %s", err)
	}

	c.Body = soap.Body.Contents

	return err
}

// Unmarshal get the body and unmarshal into the interface
func (c *Client) GetResponse() (string, error) {
	if len(c.Body) == 0 {
		return "", errors.New("body is empty")
	}

	sss := string(c.Body)
	SPAuth := &SecurityProviderAuthenticate{}
	if err := xml.Unmarshal([]byte(sss), &SPAuth); err != nil {
		return "", fmt.Errorf("an error occurred decoding the body: %s", err)
	}
	// Extract the response
	sss = SPAuth.Response
	if c.Method == "Execute" {
		if strings.Contains(sss, `status="fail"`) {
			return "", errors.New("the operation could not be performed")
		} else {
			sss = strings.TrimLeft(strings.TrimRight(sss, "&lt;/JournalNumber&gt;"), "&lt;JournalNumber&gt;")
		}
	}

	return sss, nil
}

func (c *Client) Unmarshal(v interface{}) error {
	if len(c.Body) == 0 {
		return fmt.Errorf("body is empty")
	}

	var f Fault

	err := xml.Unmarshal(c.Body, &f)
	if err != nil {
		return fmt.Errorf("an error occurred decoding the body: %s", err)
	}

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

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

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

// SecurityProviderAuthenticate struct
type SecurityProviderAuthenticate struct {
	Response string `xml:"response"`
}
