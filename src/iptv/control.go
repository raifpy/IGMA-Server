package iptv

import (
	"context"
	"errors"
	"net/http"
)

func (i *IpTv) ControlResponse(response *http.Response) error {
	if i.ResponseStatusCode != 0 && response.StatusCode != i.ResponseStatusCode {
		return errors.New(response.Status)
	}

	for _, eosc := range i.ErrorOnStatusCodes {
		if response.StatusCode == eosc {
			return errors.New(response.Status)
		}
	}

	return nil

}

func (i *IpTv) ControlStreamResponseError(err error) (res error) {
	if errors.Is(err, context.DeadlineExceeded) {
		res = nil
	} else {
		res = err
	}

	return
}
