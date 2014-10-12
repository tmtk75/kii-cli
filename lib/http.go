package kiicli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	p := Profile()
	ep := fmt.Sprintf("%s%s", p.EndpointUrl(), path)
	logger.Printf("%s %s", method, ep)

	body, _ := ioutil.ReadAll(r)
	req, err := http.NewRequest(method, ep, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
		logger.Printf("%s: %s\n", k, v)
	}

	if p.Curl {
		printCurlString(method, headers, ep, body)
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

func printCurlString(method string, header Headers, endpoint string, body []byte) {
	hs := make([]string, 0)
	for k, v := range header {
		hs = append(hs, fmt.Sprintf("-H'%s: %s'", k, v))
	}
	h := strings.Join(hs, " ")

	// ~/.kii/${app_id}/curl.{something}
	p := Profile()
	dataDir := fmt.Sprintf("%v", metaFilePath(p.AppId, ""))
	tmp, err := ioutil.TempFile(dataDir, "curl-data.")
	if err != nil {
		panic(err)
	}
	tmp.Write(body)
	defer tmp.Close()

	if len(body) > 0 {
		logger.Printf("curl -X%s %s %s -d @%v\n\n", method, h, endpoint, tmp.Name())
	} else {
		logger.Printf("curl -X%s %s %s\n\n", method, h, endpoint)
	}
}
