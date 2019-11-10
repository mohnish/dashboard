package subs

import (
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
	"github.com/mohnish/dashboard/component"
)

type Fetcher interface {
	Fetch() (interface{}, time.Time, error)
}

type SimpleFetcher struct {
	// Url is the endpoint to hit during a fetch request
	Url string
	// Interval is the polling interval for every fetch
	// TODO: (MT) Support various formats such as:
	// `10s`, `2m`, `1h`, `1h20m3s` etc
	Interval time.Duration
}

func (sf *SimpleFetcher) Fetch() (interface{}, time.Time, error) {
	// 1. create []interface{} from resp
	// 2. populate err
	// 3. set next fetch time based on interval
	// TODO: (MT) Support other restful verbs form the plugins
	res, err := http.Get(sf.Url)

	// TODO: (MT) extract this error handling out somewhere
	if err != nil {
		// retry 10 seconds later
		return nil, time.Now().Add(10 * time.Second), err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		// retry 10 seconds later
		return nil, time.Now().Add(10 * time.Second), err
	}

	defer res.Body.Close()

	return body, time.Now().Add(sf.Interval), nil
}

type ComponentFetcher struct {
	// Url is the endpoint to hit during a fetch request
	Url string
	// Interval is the polling interval for every fetch
	// TODO: (MT) Support various formats such as:
	// `10s`, `2m`, `1h`, `1h20m3s` etc
	Interval time.Duration
}

func (sf *ComponentFetcher) Fetch() (interface{}, time.Time, error) {
	// 1. create []interface{} from resp
	// 2. populate err
	// 3. set next fetch time based on interval
	// TODO: (MT) Support other restful verbs form the plugins
	res, err := http.Get(sf.Url)

	// TODO: (MT) extract this error handling out somewhere
	if err != nil {
		// retry 10 seconds later
		return nil, time.Now().Add(10 * time.Second), err
	}

	body, err := ioutil.ReadAll(res.Body)

  if err != nil {
		// retry 10 seconds later
		return nil, time.Now().Add(10 * time.Second), err
	}

  var comp component.Component

  err = json.Unmarshal(body, &comp)

	if err != nil {
		// retry 10 seconds later
		return nil, time.Now().Add(10 * time.Second), err
	}

	defer res.Body.Close()

	return comp, time.Now().Add(sf.Interval), nil
}

func Fetch(url string, interval time.Duration) *ComponentFetcher {
  // return &SimpleFetcher{
	// 	Url:      url,
	// 	Interval: interval,
	// }

	return &ComponentFetcher{
		Url:      url,
		Interval: interval,
	}
}
