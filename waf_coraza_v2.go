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

func (w *wafCorazaV2) Init() {
	w.waf = corazav2.NewWaf()
	w.parser, _ = seclangv2.NewParser(w.waf)
}

func (w *wafCorazaV2) LoadDirectives(path string) error {
	return w.parser.FromFile(path)
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
	tx.transaction.RequestBodyBuffer.Write(data)
}

func (tx *txCorazaV2) ProcessRequestBody() {
	tx.transaction.ProcessRequestBody()
}

func (tx *txCorazaV2) AddResponseHeader(name string, value string) {
	tx.transaction.AddResponseHeader(name, value)
}

func (tx *txCorazaV2) ProcessResponseHeaders(statusCode int, status string) {
	tx.transaction.ProcessResponseHeaders(statusCode, status)
}

func (tx *txCorazaV2) AppendToResponseBody(data []byte) {
	tx.transaction.ResponseBodyBuffer.Write(data)
}

func (tx *txCorazaV2) ProcessResponseBody() {
	tx.transaction.ProcessResponseBody()
}

func (tx *txCorazaV2) ProcessLogging() {
	tx.transaction.ProcessLogging()
}

func (tx *txCorazaV2) Clean() {
	tx.transaction.Clean()
}

var _ transactionIface = (*txCorazaV2)(nil)
