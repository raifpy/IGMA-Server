package iptv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/fsnotify/fsnotify"
)

type Options struct {
	IpTvConfigPath string
	OnUpdate       func()
	OnError        func(error)
	HttpClient     *http.Client
}

type IpTvStand struct {
	watching        bool
	Path            string
	List            []*IpTv
	Watcher         *fsnotify.Watcher
	OnWatcherUpdate func()
	OnWatcherError  func(error)

	HttpClient *http.Client
}

/*

	Yeterli bir algoritmaya sahip değil!

*/

func NewIpTvStand(o Options) (i *IpTvStand, err error) {
	i = &IpTvStand{
		Path:            o.IpTvConfigPath,
		OnWatcherUpdate: o.OnUpdate,
		OnWatcherError:  o.OnError,
		HttpClient:      o.HttpClient,
	}
	i.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}
	list, err := i.ReadIpTvListFromDisk(i.Path)
	if err != nil {
		return i, err
	}
	i.List = list
	return i, i.Watch()
}

func (i *IpTvStand) ReadIpTvListFromDisk(path string) ([]*IpTv, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var obj []*IpTv
	err = json.NewDecoder(file).Decode(&obj)
	if err != nil {
		return obj, err
	}
	for _, l := range obj {
		l.Client = i.HttpClient
		l.OnError = i.OnWatcherError

		var err error
		l._URL, err = url.Parse(l.URL)
		if err != nil {
			return obj, err
		}

	}
	return obj, nil
}

func (i *IpTvStand) Watch() error {
	if i.watching {
		return fmt.Errorf("already watching %s", i.Path)
	}

	if err := i.Watcher.Add(i.Path); err != nil {
		return err
	}

	i.watching = true
	go func(i *IpTvStand) {
		for event := range i.Watcher.Events {
			if event.Name != i.Path {
				continue
			}

			list, err := i.ReadIpTvListFromDisk(i.Path)
			if err != nil {
				i.OnWatcherError(err)
				continue
			}
			i.List = list
			i.OnWatcherUpdate()
		}
	}(i)

	go func() {
		for err := range i.Watcher.Errors {
			i.OnWatcherError(err)
		}
	}()

	return nil
}

type IpTv struct {
	Name                  string            `json:"name"`
	Stream                bool              `json:"stream"`
	DeadlineTimeout       int               `json:"deadline_timeout"`
	URL                   string            `json:"url"`
	ReRequestOnErrorCount int               `json:"rerequest_on_error_count"` // 0
	ResponseStatusCode    int               `json:"response_status_code"`
	ErrorOnStatusCodes    []int             `json:"error_on_status_codes"`
	Channels              map[string]string `json:"channels"`

	Headers   map[string]string `json:"headers"`
	Loop      int               `json:"loop"`
	LoopSleep int               `json:"loop_sleep"`

	Client  *http.Client `json:"-"`
	OnError func(error)  `json:"-"`
	_URL    *url.URL     `json:"-"`
}
