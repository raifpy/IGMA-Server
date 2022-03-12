package sporxmac

import (
	"net/http"
	"strings"
)

func (s *SporxMax) request() (*http.Response, error) {
	req, err := http.NewRequest("GET", s.requestURL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range s.headers {
		req.Header.Set(key, strings.Join(value, ";"))
	}

	/*response, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()*/
	return s.client.Do(req)

}
