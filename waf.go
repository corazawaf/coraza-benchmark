package main

type wafIface interface {
	Init()
	LoadDirectives(path string) error
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
	"modsec_v3": &wafModsecV3{},
}
