package tvmaclari

func Get() (Matchs, error) {
	t := New(Options{})
	req, err := t.request()
	if err != nil {
		return Matchs{}, err
	}
	defer req.Body.Close()
	return t.parse(req.Body)
}
