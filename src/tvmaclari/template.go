package tvmaclari

type Matchs []Match
type Match struct {
	MatchType string
	HomeName  string
	AwayName  string
	MatchName string
	Time      string
	Channels  []string
}
