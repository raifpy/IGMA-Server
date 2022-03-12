package tvmaclari

import "net/http"

func (t *TvMaclari) request() (*http.Response, error) {
	req, err := http.NewRequest("GET", t._url, nil)
	if err != nil {
		return nil, err
	}
	return t.Client.Do(req)
}
