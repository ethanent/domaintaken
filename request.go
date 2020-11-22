package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const QueryURL = "https://cloudflare-dns.com/dns-query"

type DNSResponse struct {
	Status   int  `json:"status"`
	TC       bool `json:"TC"`
	RD       bool `json:"RD"`
	RA       bool `json:"RA"`
	AD       bool `json:"AD"`
	CD       bool `json:"CD"`
	Question []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	} `json:"Question"`
	Answer []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		TTL  int    `json:"TTL"`
		Data string `json:"data"`
	}
}

type DNSRequest struct {
	Name string
	Type string
	DO   *bool
	CD   *bool
}

func PerformRequest(r *DNSRequest) (*DNSResponse, error) {
	headers := http.Header{}

	headers.Add("accept", "application/dns-json")

	reqURL, err := url.Parse(QueryURL)

	if err != nil {
		panic(err)
		// This should not happen, clearly, because we are parsing a constant URL.
	}

	q := reqURL.Query()

	q.Set("name", r.Name)
	q.Set("type", r.Type)

	if r.DO != nil {
		q.Set("do", strconv.FormatBool(*r.DO))
	}

	if r.CD != nil {
		q.Set("cd", strconv.FormatBool(*r.CD))
	}

	reqURL.RawQuery = q.Encode()

	req := &http.Request{
		Method:           "GET",
		URL:              reqURL,
		Header:           headers,
	}

	useClient := &http.Client{
		Timeout: time.Second * 10000,
	}

	resp, err := useClient.Do(req)

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("unexpected status code " + strconv.Itoa(resp.StatusCode))
	}

	dec := json.NewDecoder(resp.Body)

	respData := &DNSResponse{}

	if err := dec.Decode(respData); err != nil {
		return nil, err
	}

	return respData, nil
}
