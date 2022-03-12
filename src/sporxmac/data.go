package sporxmac

func Get() (Matchs, error) {

	t, err := NewSporxMax()
	if err != nil {
		return nil, err
	}
	req, err := t.request()
	if err != nil {
		return Matchs{}, err
	}
	defer req.Body.Close()
	return t.parse(req.Body)
}
