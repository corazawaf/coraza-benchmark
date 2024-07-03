package main

import (
	_ "github.com/corazawaf/coraza-benchmark/v2/pcre"
	corazav2 "github.com/corazawaf/coraza/v2"
	seclangv2 "github.com/corazawaf/coraza/v2/seclang"
	_ "github.com/corazawaf/libinjection-go"
)

type wafCorazaV2 struct {
	waf    *corazav2.Waf
	parser *seclangv2.Parser
}

func (w *wafCorazaV2) Init(rulesPath string) (err error) {
	w.waf = corazav2.NewWaf()
	w.parser, err = seclangv2.NewParser(w.waf)
	if err != nil {
		return err
	}

	return w.parser.FromFile(rulesPath)
}

func (w *wafCorazaV2) NewTransaction() transactionIface {
	return &txCorazaV2{
		w.waf.NewTransaction(),
	}
}

var _ wafIface = (*wafCorazaV2)(nil)

type txCorazaV2 struct {
	transaction *corazav2.Transaction
}

func (tx *txCorazaV2) ProcessConnection(clientAddr string, clientPort int, serverAddr string, serverPort int) {
	tx.transaction.ProcessConnection(clientAddr, clientPort, serverAddr, serverPort)
}

func (tx *txCorazaV2) ProcessURI(method string, uri string, httpVersion string) {
	tx.transaction.ProcessURI(method, uri, httpVersion)
}

func (tx *txCorazaV2) AddRequestHeader(name string, value string) {
	tx.transaction.AddRequestHeader(name, value)
}

func (tx *txCorazaV2) ProcessRequestHeaders() {
	tx.transaction.ProcessRequestHeaders()
}

func (tx *txCorazaV2) AppendToRequestBody(data []byte) {
	if _, err := tx.transaction.RequestBodyBuffer.Write(data); err != nil {
		panic(err)
	}
}

func (tx *txCorazaV2) ProcessRequestBody() {
	if _, err := tx.transaction.ProcessRequestBody(); err != nil {
		panic(err)
	}
}

func (tx *txCorazaV2) AddResponseHeader(name string, value string) {
	tx.transaction.AddResponseHeader(name, value)
}

func (tx *txCorazaV2) ProcessResponseHeaders(statusCode int, status string) {
	tx.transaction.ProcessResponseHeaders(statusCode, status)
}

func (tx *txCorazaV2) AppendToResponseBody(data []byte) {
	if _, err := tx.transaction.ResponseBodyBuffer.Write(data); err != nil {
		panic(err)
	}
}

func (tx *txCorazaV2) ProcessResponseBody() {
	if _, err := tx.transaction.ProcessResponseBody(); err != nil {
		panic(err)
	}
}

func (tx *txCorazaV2) ProcessLogging() {
	tx.transaction.ProcessLogging()
}

func (tx *txCorazaV2) Clean() {
	if err := tx.transaction.Clean(); err != nil {
		panic(err)
	}
}

var _ transactionIface = (*txCorazaV2)(nil)
