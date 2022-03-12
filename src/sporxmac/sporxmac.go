package sporxmac

import "net/http"

type SporxMax struct {
	requestURL string
	client     *http.Client
	headers    http.Header
}

func NewSporxMax() (*SporxMax, error) {
	return &SporxMax{
		requestURL: "https://www.sporx.com/tvdebugun/",
		client:     http.DefaultClient,
		headers:    http.Header{},
	}, nil
}

func (s *SporxMax) Get() (Matchs, error) {
	res, err := s.request()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return s.parse(res.Body)
}
