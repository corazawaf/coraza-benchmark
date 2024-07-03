package main

import (
	corazav3 "github.com/corazawaf/coraza/v3"
	corazav3types "github.com/corazawaf/coraza/v3/types"
)

type wafCorazaV3 struct {
	waf corazav3.WAF
}

func (w *wafCorazaV3) Init(rulesPath string) (err error) {
	w.waf, err = corazav3.NewWAF(
		corazav3.NewWAFConfig().
			WithDirectivesFromFile(rulesPath).
			WithRequestBodyAccess().
			WithResponseBodyAccess(),
	)

	return err
}

func (w *wafCorazaV3) NewTransaction() transactionIface {
	return &txCorazaV3{
		w.waf.NewTransaction(),
	}
}

var _ wafIface = (*wafCorazaV3)(nil)

type txCorazaV3 struct {
	transaction corazav3types.Transaction
}

func (tx *txCorazaV3) ProcessConnection(clientAddr string, clientPort int, serverAddr string, serverPort int) {
	tx.transaction.ProcessConnection(clientAddr, clientPort, serverAddr, serverPort)
}

func (tx *txCorazaV3) ProcessURI(method string, uri string, httpVersion string) {
	tx.transaction.ProcessURI(method, uri, httpVersion)
}

func (tx *txCorazaV3) AddRequestHeader(name string, value string) {
	tx.transaction.AddRequestHeader(name, value)
}

func (tx *txCorazaV3) ProcessRequestHeaders() {
	tx.transaction.ProcessRequestHeaders()
}

func (tx *txCorazaV3) AppendToRequestBody(data []byte) {
	if _, _, err := tx.transaction.WriteRequestBody(data); err != nil {
		panic(err)
	}
}

func (tx *txCorazaV3) ProcessRequestBody() {
	if _, err := tx.transaction.ProcessRequestBody(); err != nil {
		panic(err)
	}
}

func (tx *txCorazaV3) AddResponseHeader(name string, value string) {
	tx.transaction.AddResponseHeader(name, value)
}

func (tx *txCorazaV3) ProcessResponseHeaders(statusCode int, status string) {
	tx.transaction.ProcessResponseHeaders(statusCode, status)
}

func (tx *txCorazaV3) AppendToResponseBody(data []byte) {
	if _, _, err := tx.transaction.WriteResponseBody(data); err != nil {
		panic(err)
	}
}

func (tx *txCorazaV3) ProcessResponseBody() {
	if _, err := tx.transaction.ProcessResponseBody(); err != nil {
		panic(err)
	}
}

func (tx *txCorazaV3) ProcessLogging() {
	tx.transaction.ProcessLogging()
}

func (tx *txCorazaV3) Clean() {
	if err := tx.transaction.Close(); err != nil {
		panic(err)
	}
}

var _ transactionIface = (*txCorazaV3)(nil)
