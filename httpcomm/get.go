package httpcomm

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	neturl "net/url"

	"github.com/dchaykin/mygolib/auth"
	"github.com/dchaykin/mygolib/log"
)

type HTTPResult struct {
	Answer     []byte
	StatusCode int
	Status     string
	url        string
	method     string
	err        error
}

func (hr *HTTPResult) GetError() error {
	if hr.err != nil {
		return hr.err
	}
	switch hr.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		{
			return nil
		}
	}
	return fmt.Errorf("%s %s: %s (%d)", hr.method, hr.GetURL(), hr.Status, hr.StatusCode)
}

func (hr *HTTPResult) GetURL() string {
	u, err := neturl.Parse(hr.url)
	if err != nil {
		log.WrapError(err)
		return ""
	}
	return fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
}

func Get(url string, identity auth.UserIdentity, parameters map[string]string, headers map[string]string) (httpResult HTTPResult) {
	return get(url, identity, parameters, headers, false)
}

func get(url string, identity auth.UserIdentity, parameters map[string]string, headers map[string]string, insecure bool) (httpResult HTTPResult) {
	log.Debug("/GET %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WrapError(err)
	}

	if identity != nil {
		if err = identity.Set(req); err != nil {
			return HTTPResult{err: err}
		}
	}

	for key := range headers {
		req.Header.Set(key, headers[key])
	}

	q := req.URL.Query()
	for key, value := range parameters {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return HTTPResult{err: err}
	}
	defer response.Body.Close()

	hr := HTTPResult{
		StatusCode: response.StatusCode,
		Status:     response.Status,
		url:        url,
		method:     "GET",
	}

	if hr.GetError() == nil {
		if hr.Answer, err = io.ReadAll(response.Body); err != nil {
			return HTTPResult{err: err}
		}
	}

	return hr
}

func getContentType() string {
	return "application/json"
}
