package soccer

import (
	"context"
	"encoding/json"
	"soccerapi/src/scores"

	"go.mongodb.org/mongo-driver/bson"
)

type MatchApi struct {
	*Soccer
}

func (s *Soccer) SetMatchApi() {
	s.MatchApi = &MatchApi{
		Soccer: s,
	}
}

func (s *Soccer) GetMatchsJSONFiber() ([]byte, error) {

	s.GameScoreWatcher.Scores.Updates.RMutex.RLock()
	defer s.GameScoreWatcher.Scores.Updates.RMutex.RUnlock()

	return json.Marshal(s.GameScoreWatcher.Scores.Updates)
}

func (s *Soccer) UpdateMatchClips(mc MatchClips) error {

	if varmi, err := s.GetMatchClips(mc.MatchID); err == nil {
		mc.Clips = append(mc.Clips, varmi.Clips...)
	}

	_, err := s.Mongo.Database("matchs").Collection("match").ReplaceOne(context.Background(), bson.M{
		"matchid": mc.MatchID,
	}, mc)
	return err
}

//TODO UPDATE
func (s *Soccer) SetMatchClips(mc MatchClips) error {
	_, err := s.Mongo.Database("matchs").Collection("match").InsertOne(context.Background(), mc)
	return err
}
func (s *Soccer) ReplaceMatchClips(mc MatchClips) error {
	_, err := s.Mongo.Database("matchs").Collection("match").ReplaceOne(context.Background(), bson.M{
		"matchid": mc.MatchID,
	}, mc)
	return err
}

func (s *Soccer) GetMatchClips(matchid int64) (c MatchClips, err error) {
	err = s.Mongo.Database("matchs").Collection("match").FindOne(context.Background(), bson.M{
		"matchid": matchid,
	}).Decode(&c)
	return
}

type ClientMatch struct {
	Updates *scores.MatchUpdates        `json:"updates"`
	Clips   map[int64][]ClientMatchClip `json:"clips ,omitempty"`
	//Clips []ClientMatchClip
}

type ClientMatchClip struct {
	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`
	/*HomeTeam  string `json:"home_team"`
	AwayTeam  string `json:"away_team"`*/

	Update   string `json:"update"`
	MatchID  int64  `json:"match_id"`
	UpdateID int64  `json:"update_id"`
	Status   string `json:"status"` // goal,hot,rendering

	Minute string `json:"minute"`

	ClipID int64 `json:"clip_id,omitempty"`

	CheckAfter int `json:"check_after,omitempty"` //second
}

type MatchClip struct {
	Client ClientMatchClip

	Path     string
	GdriveID string

	WorkerID int64
	JobID    int64
}

type MatchClips struct {
	Clips   []MatchClip
	MatchID int64
}

func (mcs MatchClips) GetFromUpdateID(uid int64) (MatchClip, bool) {
	for _, v := range mcs.Clips {
		if v.Client.UpdateID == uid {
			return v, true
		}
	}
	return MatchClip{}, false
}

func (m MatchClips) ToClient() (c []ClientMatchClip) {
	for _, a := range m.Clips {
		c = append(c, a.Client)
	}
	return
}

func (m MatchClips) ToClientJSON() []byte {
	a, _ := json.Marshal(m.ToClient())
	return a
}
