package main

/*
// set ldflags
#cgo LDFLAGS: -L/usr/local/modsecurity/lib -lmodsecurity -linjection
#cgo CFLAGS: -I/usr/local/include -I/usr/local/modsecurity/include
#include "modsecurity/modsecurity.h"
#include "modsecurity/transaction.h"
#include "modsecurity/rules_set.h"

void cb(void *log, const void *data)
{
    // swallow it
    return;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type wafModsec struct {
	waf     *C.modsecurity
	ruleset *C.RulesSet
}

func (w *wafModsec) LoadDirectives(path string) error {
	p := C.CString(path)
	var err *C.char
	defer C.free(unsafe.Pointer(p))
	C.msc_rules_add_file(w.ruleset, p, &err)
	if err != nil {
		return fmt.Errorf(C.GoString(err))
	}
	return nil
}

func (w *wafModsec) NewTransaction() transactionIface {
	return &txModsec{
		transaction: C.msc_new_transaction(w.waf, w.ruleset, nil),
	}
}

func (w *wafModsec) Init() {
	w.waf = C.msc_init()
	w.ruleset = C.msc_create_rules_set()
	C.msc_set_log_cb(w.waf, (*[0]byte)(C.cb))
}

type txModsec struct {
	transaction *C.Transaction
}

func (tx *txModsec) ProcessConnection(clientAddr string, clientPort int, serverAddr string, serverPort int) {
	caddr := C.CString(clientAddr)
	saddr := C.CString(serverAddr)
	C.msc_process_connection(tx.transaction, caddr, C.int(clientPort), saddr, C.int(serverPort))
	C.free(unsafe.Pointer(caddr))
	C.free(unsafe.Pointer(saddr))
}

func (tx *txModsec) ProcessURI(method string, uri string, httpVersion string) {
	cm := C.CString(method)
	cu := C.CString(uri)
	cv := C.CString(httpVersion)
	C.msc_process_uri(tx.transaction, cm, cu, cv)
	C.free(unsafe.Pointer(cm))
	C.free(unsafe.Pointer(cu))
	C.free(unsafe.Pointer(cv))
}

func (tx *txModsec) AddRequestHeader(name string, value string) {
	if name == "" || value == "" {
		return
	}
	cn := []byte(name)
	cv := []byte(value)
	C.msc_add_request_header(tx.transaction, (*C.uchar)(&cn[0]), (*C.uchar)(&cv[0]))
}

func (tx *txModsec) ProcessRequestHeaders() {
	C.msc_process_request_headers(tx.transaction)
}

func (tx *txModsec) AppendToRequestBody(data []byte) {
	if len(data) == 0 {
		return
	}
	C.msc_append_request_body(tx.transaction, (*C.uchar)(&data[0]), C.size_t(len(data)))
}

func (tx *txModsec) ProcessRequestBody() {
	C.msc_process_request_body(tx.transaction)
}

func (tx *txModsec) AddResponseHeader(name string, value string) {
	cn := []byte(name)
	cv := []byte(value)
	C.msc_add_response_header(tx.transaction, (*C.uchar)(&cn[0]), (*C.uchar)(&cv[0]))
}

func (tx *txModsec) ProcessResponseHeaders(statusCode int, status string) {
	st := C.CString(status)
	C.msc_process_response_headers(tx.transaction, C.int(statusCode), st)
	C.free(unsafe.Pointer(st))
}

func (tx *txModsec) AppendToResponseBody(data []byte) {
	C.msc_append_response_body(tx.transaction, (*C.uchar)(&data[0]), C.size_t(len(data)))
}

func (tx *txModsec) ProcessResponseBody() {
	C.msc_process_response_body(tx.transaction)
}

func (tx *txModsec) ProcessLogging() {
	C.msc_process_logging(tx.transaction)
}

func (tx *txModsec) Clean() {
	C.msc_transaction_cleanup(tx.transaction)
}

var _ wafIface = &wafModsec{}
