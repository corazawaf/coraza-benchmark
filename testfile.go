package main

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type testFile struct {
	Name        string `yaml:"name"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`
	Request     struct {
		Method   string            `yaml:"method"`
		Path     string            `yaml:"path"`
		Protocol string            `yaml:"protocol"`
		Headers  map[string]string `yaml:"headers"`
		Body     string            `yaml:"body"`
	} `yaml:"request"`
	Response struct {
		StatusCode int               `yaml:"status_code"`
		StatusMsg  string            `yaml:"status_msg"`
		Headers    map[string]string `yaml:"headers"`
		Body       string            `yaml:"body"`
	} `yaml:"response"`
}

func (tf *testFile) Run(waf wafIface, count int) (res [5]int64, err error) {
	for i := 0; i < count; i++ {
		tx := waf.NewTransaction()
		timeStart := time.Now().UnixNano()
		tx.ProcessConnection("127.0.0.1", 55555, "127.0.0.1", 80)
		tx.ProcessURI(tf.Request.Method, tf.Request.Path, tf.Request.Protocol)
		for k, v := range tf.Request.Headers {
			tx.AddRequestHeader(k, v)
		}
		tx.ProcessRequestHeaders()
		timeEnd := time.Now().UnixNano()
		res[0] += timeEnd - timeStart // phase 1
		timeStart = time.Now().UnixNano()
		tx.AppendToRequestBody([]byte(tf.Request.Body))
		tx.ProcessRequestBody()
		timeEnd = time.Now().UnixNano()
		res[1] += timeEnd - timeStart // phase 2
		timeStart = time.Now().UnixNano()
		for k, v := range tf.Response.Headers {
			tx.AddResponseHeader(k, v)
		}
		tx.ProcessResponseHeaders(tf.Response.StatusCode, tf.Response.StatusMsg)
		timeEnd = time.Now().UnixNano()
		res[2] += timeEnd - timeStart // phase 3
		timeStart = time.Now().UnixNano()
		tx.AppendToResponseBody([]byte(tf.Response.Body))
		tx.ProcessResponseBody()
		timeEnd = time.Now().UnixNano()
		res[3] += timeEnd - timeStart // phase 4
		timeStart = time.Now().UnixNano()
		tx.ProcessLogging()
		tx.Clean()
		timeEnd = time.Now().UnixNano()
		res[4] += timeEnd - timeStart // phase 5
	}
	return
}

// openTest reads a yaml test file and returns a testFile struct
func openTest(path string) (*testFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tf := &testFile{}
	dec := yaml.NewDecoder(f)
	err = dec.Decode(tf)
	if err != nil {
		return nil, err
	}
	return tf, nil
}
