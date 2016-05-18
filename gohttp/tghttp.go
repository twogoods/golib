package gohttp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

type Httpclient struct {
	client    *http.Client
	userAgent string
}

type HttpConfig struct {
	userAgent        string
	connectTimeout   time.Duration
	readWriteTimeout time.Duration
	enableCookie     bool

	tLSClientConfig *tls.Config
	proxy           func(*http.Request) (*url.URL, error)
	transport       http.RoundTripper
}

func HttpClientBuilder() *HttpConfig {
	//default setting
	return &HttpConfig{
		userAgent:        "TgHttpLib4Go",
		connectTimeout:   60 * time.Second,
		readWriteTimeout: 60 * time.Second,
		enableCookie:     false,
	}
}

func (this *HttpConfig) SetTimeout(connectTimeout time.Duration, readWriteTimeout time.Duration) *HttpConfig {
	this.connectTimeout = connectTimeout
	this.readWriteTimeout = readWriteTimeout
	return this
}
func (this *HttpConfig) ConnectTimeout(time time.Duration) *HttpConfig {
	this.connectTimeout = time
	return this
}
func (this *HttpConfig) ReadWriteTimeout(time time.Duration) *HttpConfig {
	this.readWriteTimeout = time
	return this
}
func (this *HttpConfig) EnableCookie(flag bool) *HttpConfig {
	this.enableCookie = flag
	return this
}
func (this *HttpConfig) UserAgent(userAgent string) *HttpConfig {
	this.userAgent = userAgent
	return this
}

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		err = conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, err
	}
}
func (this *HttpConfig) Build() *Httpclient {
	// create default transport
	trans := &http.Transport{
		TLSClientConfig: this.tLSClientConfig,
		Proxy:           this.proxy,
		Dial:            TimeoutDialer(this.connectTimeout, this.readWriteTimeout),
	}

	var jar http.CookieJar
	if this.enableCookie {
		newjar, _ := cookiejar.New(nil)
		jar = newjar
	}

	client := &http.Client{
		Transport: trans,
		Jar:       jar,
	}

	return &Httpclient{client: client, userAgent: this.userAgent}
}

func (this *Httpclient) Execute(req *http.Request) (*TGHttpResponse, error) {
	req.Header.Set("User-Agent", this.userAgent)
	resp, err := this.client.Do(req)
	if err != nil {
		return nil, err
	} else {
		httpResponse := &TGHttpResponse{
			Response:   resp,
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
		}
		return httpResponse, nil
	}
}

func (this *Httpclient) Submit(req *http.Request) *Future {
	req.Header.Set("User-Agent", this.userAgent)
	errChan := make(chan error)
	respChan := make(chan *TGHttpResponse)
	go func(client *http.Client, req *http.Request, errChan chan error, respChan chan *TGHttpResponse) {
		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			respChan <- nil
		} else {
			httpResponse := &TGHttpResponse{
				Response:   resp,
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
			}
			errChan <- nil
			respChan <- httpResponse
		}
	}(this.client, req, errChan, respChan)
	return &Future{
		done:     false,
		errChan:  errChan,
		respChan: respChan,
	}
}

func (this *Httpclient) Async(req *http.Request, onResponse func(*TGHttpResponse), onFailure func(error)) {
	req.Header.Set("User-Agent", this.userAgent)
	go func() {
		resp, err := this.client.Do(req)
		if err != nil {
			onFailure(err)
		} else {
			httpResponse := &TGHttpResponse{
				Response:   resp,
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
			}
			onResponse(httpResponse)
		}
	}()
}

type Future struct {
	done     bool
	err      error
	resp     *TGHttpResponse
	errChan  chan error
	respChan chan *TGHttpResponse
}

//this process will block until the request done
func (future *Future) Get() (*TGHttpResponse, error) {
	if !future.done {
		future.err = <-future.errChan
		future.resp = <-future.respChan
	}
	return future.resp, future.err
}

type TGHttpResponse struct {
	Response   *http.Response
	Status     string // e.g. "200 OK"
	StatusCode int
}

func (tgHttpResponse *TGHttpResponse) BodyString() (string, error) {
	defer tgHttpResponse.Response.Body.Close()
	if bodyByte, err := ioutil.ReadAll(tgHttpResponse.Response.Body); err == nil {
		return string(bodyByte), err
	} else {
		return "", err
	}
}

func (tgHttpResponse *TGHttpResponse) BodyByte() ([]byte, error) {
	defer tgHttpResponse.Response.Body.Close()
	if bodyByte, err := ioutil.ReadAll(tgHttpResponse.Response.Body); err == nil {
		return bodyByte, err
	} else {
		return nil, err
	}
}

type TGHttpRequest struct {
	url       string
	method    string
	mediaType string
	params    map[string][]string
	files     map[string]string
	json      []byte
	header    map[string][]string
}

//create request
func RequestBuilder() *TGHttpRequest {
	return &TGHttpRequest{
		method: "GET",
		header: make(map[string][]string),
	}
}

func (this *TGHttpRequest) Header(key string, value string) *TGHttpRequest {
	this.header[key] = append(this.header[key], value)
	return this
}
func (this *TGHttpRequest) Cookie(cookie *http.Cookie) *TGHttpRequest {
	this.header["Cookie"] = append(this.header["Cookie"], cookie.String())
	return this
}

func (this *TGHttpRequest) Url(rawurl string) *TGHttpRequest {
	this.url = rawurl
	return this
}

func (this *TGHttpRequest) Put(reqBody *RequestBody) *TGHttpRequest {
	this.method = "PUT"
	this.mediaType = reqBody.mediaType
	this.params = reqBody.params
	this.files = reqBody.files
	return this
}

func (this *TGHttpRequest) Patch(reqBody *RequestBody) *TGHttpRequest {
	this.method = "PATCH"
	this.mediaType = reqBody.mediaType
	this.params = reqBody.params
	this.files = reqBody.files
	return this
}

func (this *TGHttpRequest) DeleteWithBody(reqBody *RequestBody) *TGHttpRequest {
	this.method = "DELETE"
	this.mediaType = reqBody.mediaType
	this.params = reqBody.params
	return this
}
func (this *TGHttpRequest) Delete() *TGHttpRequest {
	this.method = "DELETE"
	return this
}
func (this *TGHttpRequest) Head() *TGHttpRequest {
	this.method = "HEAD"
	return this
}
func (this *TGHttpRequest) Post(reqBody *RequestBody) *TGHttpRequest {
	this.method = "POST"
	this.mediaType = reqBody.mediaType
	this.params = reqBody.params
	this.files = reqBody.files
	this.json = reqBody.json
	return this
}

func (this *TGHttpRequest) Build() (*http.Request, error) {
	if "" == this.url {
		return nil, errors.New("request url can not be empty")
	}
	u, err := url.Parse(this.url)
	if err != nil {
		return nil, err
	}
	req := http.Request{
		URL:        u,
		Method:     this.method,
		Header:     this.header,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	req.Header.Set("Content-Type", this.mediaType)
	// build POST/PUT/PATCH url and body
	if this.method == "POST" || this.method == "PUT" || this.method == "PATCH" || this.method == "DELETE" {
		if len(this.files) > 0 {
			pr, pw := io.Pipe()
			bodyWriter := multipart.NewWriter(pw)
			go func() {
				for formname, filename := range this.files {
					fileWriter, err := bodyWriter.CreateFormFile(formname, filename)
					if err != nil {
						log.Println("goHttp:", err)
					}
					fh, err := os.Open(filename)
					if err != nil {
						log.Println("goHttp:", err)
					}
					_, err = io.Copy(fileWriter, fh)
					fh.Close()
					if err != nil {
						log.Println("goHttp:", err)
					}
				}
				for k, v := range this.params {
					for _, vv := range v {
						bodyWriter.WriteField(k, vv)
					}
				}
				bodyWriter.Close()
				pw.Close()
			}()
			req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
			req.Body = ioutil.NopCloser(pr)
		} else if len(this.params) > 0 {
			paramBody := buildParamBody(this.params)
			bf := bytes.NewBufferString(paramBody)
			req.Body = ioutil.NopCloser(bf)
			req.ContentLength = int64(len(paramBody))
		} else if len(this.json) > 0 {
			req.Body = ioutil.NopCloser(bytes.NewReader(this.json))
			req.ContentLength = int64(len(this.json))
		}
	}
	return &req, nil
}

func buildParamBody(params map[string][]string) string {
	var paramBody string
	if len(params) > 0 {
		var buf bytes.Buffer
		for k, v := range params {
			for _, vv := range v {
				buf.WriteString(url.QueryEscape(k))
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(vv))
				buf.WriteByte('&')
			}
		}
		paramBody = buf.String()
		paramBody = paramBody[0 : len(paramBody)-1]
	}
	return paramBody
}

func BuildGetUrl(baseUrl string, params map[string][]string) string {
	if len(params) > 0 {
		var buf bytes.Buffer
		for k, v := range params {
			for _, vv := range v {
				buf.WriteString(url.QueryEscape(k))
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(vv))
				buf.WriteByte('&')
			}
		}
		baseUrl = baseUrl + "?" + buf.String()
		baseUrl = baseUrl[0 : len(baseUrl)-1]
	}
	return baseUrl
}

type RequestBody struct {
	mediaType string
	params    map[string][]string
	files     map[string]string
	json      []byte
}

type FormRequestBody struct {
	mediaType string
	params    map[string][]string
}

//form
func FormBodyBuilder() *FormRequestBody {
	return &FormRequestBody{params: make(map[string][]string)}
}
func (this *FormRequestBody) AddParam(key string, value string) *FormRequestBody {
	par := this.params
	par[key] = append(par[key], value)
	return this
}
func (this *FormRequestBody) Build() *RequestBody {
	req := RequestBody{
		mediaType: "application/x-www-form-urlencoded",
		params:    this.params,
	}
	return &req
}

//json
type JsonRequestBody struct {
	mediaType string
	json      []byte
}

func JsonBodyBuilder() *JsonRequestBody {
	return &JsonRequestBody{}
}

func (this *JsonRequestBody) Json(data []byte) *JsonRequestBody {
	this.json = data
	return this
}
func (this *JsonRequestBody) Build() *RequestBody {
	req := &RequestBody{
		mediaType: "application/json; charset=utf-8",
		json:      this.json,
	}
	return req
}

//file
type MuiltRequestBody struct {
	mediaType string
	files     map[string]string
	params    map[string][]string
}

func MuiltBodyBuilder() *MuiltRequestBody {
	return &MuiltRequestBody{files: make(map[string]string), params: make(map[string][]string)}
}

func (this *MuiltRequestBody) AddFile(formname string, filename string) *MuiltRequestBody {
	this.files[formname] = filename
	return this
}
func (this *MuiltRequestBody) AddParam(key string, value string) *MuiltRequestBody {
	par := this.params
	par[key] = append(par[key], value)
	return this
}
func (this *MuiltRequestBody) Build() *RequestBody {
	req := &RequestBody{
		files:  this.files,
		params: this.params,
	}
	return req
}
