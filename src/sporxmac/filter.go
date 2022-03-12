package sporxmac

import (
	"strings"
)

func (m Matchs) ContainesMatchName(name string) bool {
	for _, a := range m {
		if strings.EqualFold(a.MatchName, name) {
			return true
		}
	}
	return false
}

type StringList []string

func (sl StringList) Contains(key string) bool {
	for _, a := range sl {
		if strings.EqualFold(a, key) {
			return true
		}
	}
	return false
}

func (sl StringList) ContainsList(list StringList) bool {
	for _, a := range sl {
		if list.Contains(a) {
			return true
		}
	}
	return false
}

func (m Matchs) FilterWithSportName(sportname string) (yenim Matchs) {
	for _, a := range m {
		//if a.MatchType

		if strings.EqualFold(a.MatchType, sportname) {
			yenim = append(yenim, a)
		}
	}
	return
}

func (m Matchs) FilterWithChannelsName(whitelist StringList) (yenim Matchs) {
	for _, veri := range m {
		/*if yenim.ContainesMatchName(veri.MatchName) {
			log.Println("duplicate: ", veri.MatchName)
			continue
		}*/

		for _, channel := range veri.Channels {
			if whitelist.Contains(channel) {
				yenim = append(yenim, veri)
			}
		}
	}
	return
}
