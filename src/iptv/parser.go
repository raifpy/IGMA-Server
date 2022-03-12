package iptv

import (
	"io"
	"net/url"
	"regexp"
	"strings"
)

var httpregex = regexp.MustCompile(`https?:\/\/`)

type IpTvDetails struct {
	StreamInfo string

	Redirect     bool
	RedirectPath string

	Media      bool
	MediaPaths []string

	UsesOtherURL bool
}

func (i *IpTv) cloneURL(path string) *url.URL {

	if path == "" {
		return i._URL
	}

	var query string

	querysplit := strings.Split(path, "?")
	if len(querysplit) > 1 {
		query = querysplit[1]
		path = querysplit[0]
	}

	//fmt.Printf("path: %v\n", path)

	return &url.URL{
		Scheme:   i._URL.Scheme,
		Opaque:   i._URL.Opaque,
		User:     i._URL.User,
		Host:     i._URL.Host,
		Path:     i._URL.Path + path,
		RawPath:  i._URL.RawPath + path,
		RawQuery: query,
	}

}

func parse(r io.Reader) (IpTvDetails, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return IpTvDetails{}, err
	}
	return Parse(string(body)), nil
}

func Parse(body string) (p IpTvDetails) {

	p.UsesOtherURL = httpregex.MatchString(body)
	p.Media = p.UsesOtherURL

	var split = strings.Split(body, "\n")

	if i := strings.Index(body, "#EXT-X-STREAM-INF"); i != -1 {
		p.Redirect = !p.UsesOtherURL

		for _, _a := range split {
			if strings.HasPrefix(_a, "#EXT-X-STREAM-INF") {
				p.StreamInfo = _a[len("#EXT-X-STREAM-INF:"):]
			}
		}

		//fmt.Printf("p.StreamInfo: %v\n", p.StreamInfo)
	}

	//if strings.Contains(body, "PART") {
	if strings.Contains(body, "MEDIA") {
		p.Media = true
	}

	for _, veri := range split {
		if veri == "" {
			continue
		}
		if veri[0] == '#' {
			continue
		}

		if p.Redirect {
			p.RedirectPath = veri
			break
		}

		p.MediaPaths = append(p.MediaPaths, veri)
	}

	return p

}

type stringList []string

func (s stringList) Contains(key string) bool {
	for _, a := range s {
		if a == key {
			return true
		}
	}
	return false
}

func (is *IpTvStand) FilterByChannel(channel string) (ilist []*IpTv) {
	ilist = []*IpTv{}
	for _, ae := range is.List {
		var ok bool
		for ch := range ae.Channels {
			if ch == channel {
				ilist = append(ilist, ae)
				ok = true
				break
			}
		}
		if ok {
			break
		}
	}

	return
}
