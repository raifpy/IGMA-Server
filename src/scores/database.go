package scores

import (
	"fmt"
	"sync"
)

type Database interface {
	StoreAllScore(AllScore) error
	GetAllScore() (AllScore, bool)

	StoreGameScore(GameScore) error
	CleanGameScore(string)
	GetGameScore(string) (GameScore, bool)
}

type SyncMapDatabase struct {
	Map *sync.Map
}

func (s SyncMapDatabase) StoreAllScore(a AllScore) error {
	s.Map.Store("allscore", a)
	return nil
}
func (a SyncMapDatabase) GetAllScore() (AllScore, bool) {
	veri, ok := a.Map.Load("allscore")
	if !ok {
		return AllScore{}, false
	}
	scoreveri, ok := veri.(AllScore)
	return scoreveri, ok
}

func (a SyncMapDatabase) StoreGameScore(g GameScore) error {
	a.Map.Store(fmt.Sprint(g.GameGame.ID), g)
	return nil
}
func (a SyncMapDatabase) CleanGameScore(id string) {
	a.Map.Delete(id)
}
func (a SyncMapDatabase) GetGameScore(id string) (GameScore, bool) {
	veri, ok := a.Map.Load(id)
	if !ok {
		return GameScore{}, false
	}
	scoreveri, ok := veri.(GameScore)
	return scoreveri, ok
}

func NewSyncMapDatabase() Database {
	return SyncMapDatabase{
		Map: &sync.Map{},
	}
}
