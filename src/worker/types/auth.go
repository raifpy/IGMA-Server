package types

import "time"

type DatabaseAuth struct {
	Token string `json:"token" mongo:"token"`
	IP    string `json:"ip"`
	Id    int64  `json:"id"`
	//Baerer     string        `json:"baerer" mongo:"baerer"` //??Belki
	Until time.Time `json:"until"`
}
