package tvmaclari

import (
	"errors"
	"net/http"
)

type Options struct {
	Client *http.Client
	Error  error
}

type TvMaclari struct {
	Client *http.Client
	Error  error
	_url   string
}

func New(o Options) *TvMaclari {
	if o.Client == nil {
		o.Client = http.DefaultClient
	}

	if o.Error == nil {
		o.Error = errors.New("cannot finding")
	}

	return &TvMaclari{
		Client: o.Client,
		_url:   "https://ajansspor.com/spor-ekrani",
	}
}
