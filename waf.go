package main

type wafIface interface {
	Init(rulesPath string) error
	NewTransaction() transactionIface
}

type transactionIface interface {
	ProcessConnection(clientAddr string, clientPort int, serverAddr string, serverPort int)
	ProcessURI(method string, uri string, httpVersion string)
	AddRequestHeader(name string, value string)
	ProcessRequestHeaders()
	AppendToRequestBody(data []byte)
	ProcessRequestBody()
	AddResponseHeader(name string, value string)
	ProcessResponseHeaders(statusCode int, status string)
	AppendToResponseBody(data []byte)
	ProcessResponseBody()
	ProcessLogging()
	Clean()
	// Transaction
}

var wafInterfaces = map[string]wafIface{
	"coraza_v2": &wafCorazaV2{},
	"coraza_v3": &wafCorazaV3{},
	"modsec_v3": &wafModsecV3{},
}
