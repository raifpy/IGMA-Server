package scores

import "sync"

type MatchUpdates struct {
	RMutex     *sync.RWMutex          `json:"-"`
	UpdateID   string                 `json:"update_id"`
	PlainScore map[int]PlainGameScore `json:"scores"`
}
