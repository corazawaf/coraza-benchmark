package main

import (
	"github.com/jptosso/coraza-waf/v2"
	"github.com/jptosso/coraza-waf/v2/seclang"

	_ "github.com/jptosso/coraza-libinjection"
	_ "github.com/jptosso/coraza-pcre"
)

type wafCoraza struct {
	waf    *coraza.Waf
	parser *seclang.Parser
}

func (w *wafCoraza) Init() {
	w.waf = coraza.NewWaf()
	w.parser, _ = seclang.NewParser(w.waf)
}

func (w *wafCoraza) LoadDirectives(path string) error {
	return w.parser.FromFile(path)
}

func (w *wafCoraza) NewTransaction() transactionIface {
	return &txCoraza{
		w.waf.NewTransaction(),
	}
}

var _ wafIface = (*wafCoraza)(nil)

type txCoraza struct {
	transaction *coraza.Transaction
}

func (tx *txCoraza) ProcessConnection(clientAddr string, clientPort int, serverAddr string, serverPort int) {
	tx.transaction.ProcessConnection(clientAddr, clientPort, serverAddr, serverPort)
}

func (tx *txCoraza) ProcessURI(method string, uri string, httpVersion string) {
	tx.transaction.ProcessURI(method, uri, httpVersion)
}

func (tx *txCoraza) AddRequestHeader(name string, value string) {
	tx.transaction.AddRequestHeader(name, value)
}

func (tx *txCoraza) ProcessRequestHeaders() {
	tx.transaction.ProcessRequestHeaders()
}

func (tx *txCoraza) AppendToRequestBody(data []byte) {
	tx.transaction.RequestBodyBuffer.Write(data)
}

func (tx *txCoraza) ProcessRequestBody() {
	tx.transaction.ProcessRequestBody()
}

func (tx *txCoraza) AddResponseHeader(name string, value string) {
	tx.transaction.AddResponseHeader(name, value)
}

func (tx *txCoraza) ProcessResponseHeaders(statusCode int, status string) {
	tx.transaction.ProcessResponseHeaders(statusCode, status)
}

func (tx *txCoraza) AppendToResponseBody(data []byte) {
	tx.transaction.ResponseBodyBuffer.Write(data)
}

func (tx *txCoraza) ProcessResponseBody() {
	tx.transaction.ProcessResponseBody()
}

func (tx *txCoraza) ProcessLogging() {
	tx.transaction.ProcessLogging()
}

func (tx *txCoraza) Clean() {
	tx.transaction.Clean()
}

var _ transactionIface = (*txCoraza)(nil)
