package worker

import (
	"math/rand"
	"reflect"
	"soccerapi/src/worker/types"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

type WorkerConn struct {
	Id   int64
	Conn *websocket.Conn
}

func (c *WorkerConn) CloseConnection() {
	c.Conn.SetReadDeadline(time.Now().Add(time.Second))
	c.Conn.SetWriteDeadline(time.Now().Add(time.Second))
	c.Conn.Conn.WriteMessage(websocket.CloseMessage, nil)
}

type WorkerMap struct {
	lock *sync.RWMutex
	Map  map[int64]*WorkerConn

	chanlock *sync.RWMutex
	ChanMap  map[int64]chan types.WebsocketContact
}

//!! Nilable
func (wm *WorkerMap) GetChan(id int64) chan types.WebsocketContact {
	wm.chanlock.RLock()
	defer wm.chanlock.RUnlock()
	return wm.ChanMap[id]
}

func (wm *WorkerMap) SetChan(id int64, c chan types.WebsocketContact) {
	wm.chanlock.Lock()
	defer wm.chanlock.Unlock()
	wm.ChanMap[id] = c
}

func (wm *WorkerMap) DelChan(id int64) {
	wm.chanlock.Lock()
	defer wm.chanlock.Unlock()
	delete(wm.ChanMap, id)
}

func (wm *WorkerMap) Get(id int64) *WorkerConn {
	wm.lock.RLock()
	defer wm.lock.RUnlock()
	return wm.Map[id]
}

func (wm *WorkerMap) Set(c *WorkerConn) {
	wm.lock.Lock()
	defer wm.lock.Unlock()
	wm.Map[c.Id] = c
}
func (wm *WorkerMap) Del(c *WorkerConn) {
	wm.lock.Lock()
	defer wm.lock.Unlock()
	delete(wm.Map, c.Id)
}

func NewWorkerMap() *WorkerMap {
	return &WorkerMap{
		lock:     &sync.RWMutex{},
		Map:      make(map[int64]*WorkerConn),
		chanlock: &sync.RWMutex{},
		ChanMap:  make(map[int64]chan types.WebsocketContact),
	}
}

func (wm *WorkerMap) Random() *WorkerConn {
	wm.lock.RLock()
	defer wm.lock.RUnlock()
	keys := reflect.ValueOf(wm.Map).MapKeys()
	if len(keys) == 0 {
		return nil
	}
	return wm.Map[keys[rand.Intn(len(keys))].Interface().(int64)]
}
