# Go Soap [![Build Status](https://travis-ci.org/tiaguinho/gosoap.svg?branch=master)](https://travis-ci.org/tiaguinho/gosoap) [![GoDoc](https://godoc.org/github.com/tiaguinho/gosoap?status.png)](https://godoc.org/github.com/tiaguinho/gosoap) [![Go Report Card](https://goreportcard.com/badge/github.com/tiaguinho/gosoap)](https://goreportcard.com/report/github.com/tiaguinho/gosoap) [![codecov](https://codecov.io/gh/tiaguinho/gosoap/branch/master/graph/badge.svg)](https://codecov.io/gh/tiaguinho/gosoap)
package to help with SOAP integrations (client)

### Credits

This project is an adaptation of the github.com/tiaguinho/gosoap, to meet the specific needs of Infor Sunsystems Connect   SSC


### Install

```bash
go get github.com/UnionMexicanaDelNorte/soapClientGolangForSSC
```

### Example

```go
soap, err := gosoap.SoapClient("http://localhost:8080/connect/wsdl/SecurityProvider?wsdl","http://localhost:8080/connect/soap/SecurityProvider")
	if err != nil {
		fmt.Errorf("error not expected: %s", err)
	}
	params := gosoap.Params{
		"name": "AOK",
		"password" : "",
	}
	err = soap.Call("Authenticate", "SecurityProviderAuthenticateRequest", params)
	if err != nil {
		fmt.Errorf("error in soap call: %s", err)
	}
	vaucher := soap.GetResponse()
	fmt.Println(vaucher)
	soapJournal, err := gosoap.SoapClient("http://localhost:8080/connect/wsdl/ComponentExecutor?wsdl","http://localhost:8080/connect/soap/ComponentExecutor")
	if err != nil {
		fmt.Errorf("error not expected: %s", err)
	}

	params = gosoap.Params{
		"authentication": vaucher,
		"licensing" : "",
		"component" : "Journal",
		"method" : "Import",
		"group" : "",
		"payload" : `<SSC>
  <SunSystemsContext>
    <BusinessUnit>CEA</BusinessUnit>
    <BudgetCode>A</BudgetCode>
  </SunSystemsContext>
  <MethodContext>
    <LedgerPostingParameters>
      <JournalType>JV</JournalType>
      <PostingType>2</PostingType>
      <PostProvisional>N</PostProvisional>
      <PostToHold>N</PostToHold>
      <BalancingOptions>T2</BalancingOptions>
      <SuspenseAccount>338100</SuspenseAccount>
      <TransactionAmountAccount>338100</TransactionAmountAccount>
      <ReportingAccount>338100</ReportingAccount>
      <SupressSubstitutedMessages>N</SupressSubstitutedMessages>
      <ReportErrorsOnly>Y</ReportErrorsOnly>
    </LedgerPostingParameters>
  </MethodContext>
  <Payload>
    <Ledger>
      <Line>
        <TransactionReference>651C</TransactionReference>
        <AccountingPeriod>0052017</AccountingPeriod>
        <TransactionDate>07052017</TransactionDate>
        <AccountCode>ERROJAB01</AccountCode>
        <AnalysisCode2/>
        <AnalysisCode3>10</AnalysisCode3>
        <AnalysisCode4/>
        <AnalysisCode5/>
        <AnalysisCode6/>
        <AnalysisCode7/>
        <AnalysisCode8/>
        <AnalysisCode9/>
        <AnalysisCode10/>
        <Description>GONZALEZ ALCUDIA HUMBERTO</Description>
        <Value4Amount>3500</Value4Amount>
        <DebitCredit>D</DebitCredit>
        <Value4CurrencyCode>MXP1</Value4CurrencyCode>
        <DueDate>07052017</DueDate>
      </Line>
      <Line>
        <TransactionReference>651C</TransactionReference>
        <AccountingPeriod>0052017</AccountingPeriod>
        <TransactionDate>07052017</TransactionDate>
        <AccountCode>ERROJAB01</AccountCode>
        <AnalysisCode2/>
        <AnalysisCode3>10</AnalysisCode3>
        <AnalysisCode4/>
        <AnalysisCode5/>
        <AnalysisCode6/>
        <AnalysisCode7/>
        <AnalysisCode8/>
        <AnalysisCode9/>
        <AnalysisCode10/>
        <Description>GONZALEZ ALCUDIA HUMBERTO</Description>
        <Value4Amount>3500</Value4Amount>
        <DebitCredit>C</DebitCredit>
        <Value4CurrencyCode>MXP1</Value4CurrencyCode>
        <DueDate>07052017</DueDate>
      </Line>
    </Ledger>
  </Payload>
</SSC>
`,
	}
	err = soapJournal.Call("Execute", "ComponentExecutorExecuteRequest", params)
	if err != nil {
		fmt.Errorf("error in soap call: %s", err)
	}
	diario := soapJournal.GetResponse()
	fmt.Println(diario)
```