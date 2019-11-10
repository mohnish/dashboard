package subs

import "time"

type Subscription interface {
	Updates() <-chan interface{}
	Close() error
}

type sub struct {
	fetcher Fetcher
	updates chan interface{}
	closing chan chan error
}

func (s *sub) Updates() <-chan interface{} {
	return s.updates
}

func (s *sub) Close() error {
	errc := make(chan error)
	s.closing <- errc

	return <-errc
}

// `loop` is responsible for 3 things:
// 1. fetching updates
// 2. pushing the updates to the stream
// 3. clean up and close when requested
func (s *sub) loop() {
	var pending []interface{}
	var next time.Time
	var err error
	var first interface{}

	for {
		var fetchDelay time.Duration
		var updates chan interface{}

		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}

		startFetch := time.After(fetchDelay)

		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates
		}

		select {
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return
		case <-startFetch:
			var fetched interface{}
			fetched, next, err = s.fetcher.Fetch()

			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}

			pending = append(pending, fetched)
		case updates <- first:
			pending = pending[1:]
		}
	}
}

func Subscribe(fetcher Fetcher) *sub {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan interface{}),
	}

	go s.loop()

	return s
}

type mergedSub struct {
	subscriptions []Subscription
	updates chan interface{}
	quit chan struct{}
	err chan error
}

func (m *mergedSub) Updates() <-chan interface{} {
	return m.updates
}

func (m *mergedSub) Close() error {
	var err error
	close(m.quit)
	for _ = range m.subscriptions {
		if e := <-m.err; e != nil {
			err = e
		}
	}
	close(m.updates)

	return err
}

func Merge(subscriptions []Subscription) Subscription {
	m := &mergedSub{
		subscriptions:    subscriptions,
		updates: make(chan interface{}),
		quit:    make(chan struct{}),
		err:    make(chan error),
	}

	for _, subscription := range subscriptions {
		go func(s Subscription) {
			for {
				var update interface{}
				select {
				case update = <-s.Updates():
				case <-m.quit:
					m.err <- s.Close()
					return
				}
				select {
				case m.updates <- update:
				case <-m.quit:
					m.err <- s.Close()
					return
				}
			}
		}(subscription)
	}

	return m
}
