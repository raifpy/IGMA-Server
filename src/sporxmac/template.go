package sporxmac

import "encoding/xml"

type Matchs []Match
type Match struct {
	MatchType string
	HomeName  string
	AwayName  string
	MatchName string
	Time      string
	Channels  []string
}

type XmlMatch struct {
	XMLName xml.Name `xml:"a"`
	Text    string   `xml:",chardata"`
	Href    string   `xml:"href,attr"`
	Target  string   `xml:"target,attr"`
	Class   string   `xml:"class,attr"`
	Style   string   `xml:"style,attr"`
}
