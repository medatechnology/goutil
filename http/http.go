package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	// "github.com/thoas/go-funk"
	"github.com/medatechnology/goutil/object"
)

// http status code
type StatusCode int
type basicAuth struct {
	Username string
	Password string
}
type httpClient struct {
	client        *http.Client
	headers       map[string][]string
	params        map[string]string
	useBasicAuth  bool
	basicAuthData basicAuth
}

func NewHttp() HttpClient {
	return &httpClient{
		client:  &http.Client{},
		headers: make(map[string][]string),
		params:  make(map[string]string),
	}
}

type HttpClient interface {
	Post(urll string, body interface{}, result interface{}, errorResponse interface{}) (StatusCode, error)
	Get(urll string, result interface{}, errorResponse interface{}) (StatusCode, error)
	SetHeader(headers map[string][]string) *httpClient
	SetQueryParams(params map[string]string) *httpClient
	SetBasicAuth(username, password string) *httpClient
	// PostStream sends a POST request and returns raw response for streaming
	PostStream(url string, data any) (*http.Response, error)
	// GetStream sends a GET request and returns raw response for streaming
	GetStream(url string) (*http.Response, error)
}

func (h *httpClient) Post(urll string, body interface{}, result interface{}, errorResponse interface{}) (StatusCode, error) {
	var payload *bytes.Buffer
	if body != nil {
		// jika content typenya merupakan application json
		// maka akan di encode menjadi string menggunakan jsonENcode
		// if funk.ContainsString(h.headers["Content-Type"], "application/json") {
		if object.ArrayAContainsBString(h.headers["Content-Type"], "application/json") {
			_body, err := h.marshalPayload(body)
			if err != nil {
				return StatusCode(0), err
			}
			payload = bytes.NewBufferString(string(_body))
			// } else if funk.ContainsString(h.headers["Content-Type"], "application/x-www-form-urlencoded") {
		} else if object.ArrayAContainsBString(h.headers["Content-Type"], "application/x-www-form-urlencoded") {
			//  Jika body requestnya merupakan url-encoded
			// data akan di set pada url values lalau di encode menajdi string
			_body := object.StructToMap(body)
			data := url.Values{}
			for k, v := range _body {
				if val, ok := v.(string); ok {
					data.Set(k, val)
				}
			}
			payload = bytes.NewBufferString(data.Encode())
		}
	}
	request := &http.Request{}
	if payload != nil {
		req, err := http.NewRequest(http.MethodPost, urll, payload)
		if err != nil {

			return 0, err
		}
		request = req
	} else {
		req, err := http.NewRequest(http.MethodPost, urll, nil)
		if err != nil {
			return 0, err
		}
		request = req
	}
	request.Header = h.headers
	if h.useBasicAuth {
		request.SetBasicAuth(h.basicAuthData.Username, h.basicAuthData.Password)
	}
	response, err := h.client.Do(request)
	if err != nil {
		return StatusCode(0), err
	}
	if result != nil && response.StatusCode < 300 {
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return StatusCode(0), err
		}
		return StatusCode(response.StatusCode), nil
	}
	// jika response code tidak sama dengan 200
	// maka dilakukan pengecekan errornya dan akan di return errorr messagennya
	if response.StatusCode >= 400 {
		//  Jika status errornya bukan berbentuk json
		if response.Header.Get("Content-Type") != "application/json" {
			bodyByte, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return StatusCode(0), err
			}
			return StatusCode(response.StatusCode), errors.New(string(bodyByte))
		} else {
			//  jika status errornya merupakan json
			//  maka akan di decode hasil error codenya
			if errorResponse != nil {
				err = json.NewDecoder(response.Body).Decode(&errorResponse)
				if err != nil {
					return StatusCode(0), err
				}
			}
			return StatusCode(response.StatusCode), errors.New("error response")
		}
	}
	response.Body.Close()
	return StatusCode(response.StatusCode), nil
}

// PostStream sends a POST request and returns the raw *http.Response.
// The caller is responsible for reading and closing the response body.
// This is useful for SSE (Server-Sent Events) and streaming APIs.
func (h *httpClient) PostStream(url string, data any) (*http.Response, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// Apply stored headers
	req.Header = h.headers

	// Apply query params if any
	if len(h.params) > 0 {
		q := req.URL.Query()
		for key, value := range h.params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Apply basic auth if set
	if h.useBasicAuth {
		req.SetBasicAuth(h.basicAuthData.Username, h.basicAuthData.Password)
	}

	return h.client.Do(req)
}

// GetStream sends a GET request and returns the raw *http.Response.
// The caller is responsible for reading and closing the response body.
// This is useful for SSE (Server-Sent Events) and streaming APIs.
func (h *httpClient) GetStream(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Apply stored headers
	req.Header = h.headers

	// Apply query params if any
	if len(h.params) > 0 {
		q := req.URL.Query()
		for key, value := range h.params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Apply basic auth if set
	if h.useBasicAuth {
		req.SetBasicAuth(h.basicAuthData.Username, h.basicAuthData.Password)
	}

	return h.client.Do(req)
}

func (h *httpClient) Get(urll string, result interface{}, errorResponse interface{}) (StatusCode, error) {
	baseUrl := urll
	request, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return StatusCode(0), err
	}
	if h.params != nil {
		// Jika paramsnya tidak nil maka akan di set query paramnnya sesuai datanya
		params := url.Values{}
		for k, v := range h.params {
			params.Set(k, v)
		}
		request.URL.RawQuery = params.Encode()
	}
	// set header
	request.Header = h.headers
	if h.useBasicAuth {
		request.SetBasicAuth(h.basicAuthData.Username, h.basicAuthData.Password)
	}
	response, err := h.client.Do(request)
	if err != nil {
		return StatusCode(0), err
	}
	defer response.Body.Close()
	if result != nil && response.StatusCode < 300 {
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return StatusCode(0), err
		}
	}
	if response.StatusCode >= 400 {
		if response.Header.Get("Content-Type") != "application/json" {
			//  Jika reponse errornya bukan merupakan application/json
			// maka akan menggunakan response bodynya
			bodyByte, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return StatusCode(0), err
			}
			return StatusCode(response.StatusCode), errors.New(string(bodyByte))
		} else {
			//  jika status errornya merupakan json
			//  maka akan di decode hasil error codenya
			if errorResponse != nil {
				err = json.NewDecoder(response.Body).Decode(&errorResponse)
				if err != nil {
					return StatusCode(0), err
				}
			}
			return StatusCode(response.StatusCode), errors.New("error response")
		}
	}

	return StatusCode(response.StatusCode), nil
}

func (h *httpClient) marshalPayload(p interface{}) ([]byte, error) {
	var err error
	var data []byte
	if p != nil {
		data, err = json.Marshal(p)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}
func (h *httpClient) SetHeader(headers map[string][]string) *httpClient {
	//  untuk reset value headernya
	h.headers = map[string][]string{}
	h.headers = headers
	return h
}

func (h *httpClient) SetQueryParams(params map[string]string) *httpClient {
	// untuk reset value paramsnnya
	h.params = map[string]string{}
	h.params = params
	return h
}
func (h *httpClient) SetBasicAuth(username, password string) *httpClient {
	// untuk reset value paramsnnya
	h.basicAuthData.Username = username
	h.basicAuthData.Password = password
	h.useBasicAuth = true
	return h
}
