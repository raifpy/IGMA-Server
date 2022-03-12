package iptv

import (
	"context"
	"net/http"
)

func (i *IpTv) getRequest(method, fullurl string) (*http.Request, error) {
	req, err := http.NewRequest(method, fullurl, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range i.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (i *IpTv) getRequestContext(method, fullurl string, context context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context, method, fullurl, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range i.Headers {
		req.Header.Set(key, value)
	}
	return req, nil

}

func (i *IpTv) Request(fullurl string) (*http.Response, error) {
	request, err := i.getRequest("GET", fullurl)
	if err != nil {
		return nil, err
	}
	return i.Client.Do(request)
}

func (i *IpTv) RequestAndControl(fullrul string) (*http.Response, error) {
	res, err := i.Request(fullrul)
	if err != nil {
		return nil, err
	}

	if err := i.ControlResponse(res); err != nil {
		res.Body.Close()
		return nil, err
	}

	return res, nil
}
