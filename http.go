package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpResponse http.Response

func (self *HttpResponse) Bytes() []byte {
	b, _ := ioutil.ReadAll(self.Body)
	return b
}

func HttpPostJson(path string, headers map[string]string, body interface{}) *HttpResponse {
	reqbody, _ := json.Marshal(body)
	return HttpPost(path, headers, bytes.NewReader(reqbody))
}

func HttpPost(path string, headers Headers, r io.Reader) *HttpResponse {
	return httpRequest("POST", path, headers, r)
}

func HttpGet(path string, headers Headers) *HttpResponse {
	empty := []byte{}
	return httpRequest("GET", path, headers, bytes.NewReader(empty))
}

func HttpPut(path string, headers Headers, r io.Reader) *HttpResponse {
	return httpRequest("PUT", path, headers, r)
}

func HttpDelete(path string, headers Headers) *HttpResponse {
	empty := []byte{}
	return httpRequest("DELETE", path, headers, bytes.NewReader(empty))
}

func httpRequest(method string, path string, headers Headers, r io.Reader) *HttpResponse {
	ep := fmt.Sprintf("%s%s", globalConfig.EndpointUrl(), path)
	logger.Printf("%s %s", method, ep)
	req, err := http.NewRequest(method, ep, r)
	if err != nil {
		panic(err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
		logger.Printf("%s: %s\n", k, v)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode/100 != 2 {
		b, _ := ioutil.ReadAll(res.Body)
		log.Fatalf("%s\n", string(b))
	}
	hr := HttpResponse(*res)
	return &hr
}
