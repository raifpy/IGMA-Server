package dsystem

import (
	"sync"
	"time"
)

type DsystemValue struct {
	Token string
	Ip    string

	//MaxRequestCounter int

	Path string

	Delete time.Duration
}

type Dsystem struct {
	Map map[string]DsystemValue
	r   *sync.RWMutex
	//routergroup fiber.Router
}

func (d *Dsystem) Set(key string, value DsystemValue) {
	d.r.Lock()
	defer d.r.Unlock()

	d.Map[key] = value

	time.AfterFunc(value.Delete, func() {
		d.Del(key)
	})
}
func (d *Dsystem) Del(key string) {
	d.r.Lock()
	defer d.r.Unlock()

	delete(d.Map, key)
}

func (d *Dsystem) Get(key string) (DsystemValue, bool) {
	d.r.RLock()
	defer d.r.RUnlock()

	v, ok := d.Map[key]
	return v, ok

}

func NewDSystem() *Dsystem {
	return &Dsystem{
		r:   &sync.RWMutex{},
		Map: map[string]DsystemValue{},
		//routergroup: routergroup,
	}
}
