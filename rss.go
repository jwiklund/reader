package reader

import (
	rss "github.com/ungerik/go-rss"
	"time"
)

type FetchRssRequest struct {
	Id     string
	Result chan error
}

type FetchAllRssRequest struct {
	Type   string
	Result chan []error
}

type rsshandler struct {
	FetchChan    chan FetchRssRequest
	FetchAllChan chan FetchAllRssRequest
	CloseChan    chan bool
	store        *Store
	reader       func(string) (*rss.Channel, error)
}

type Rss interface {
	Fetch(id string) error
	FetchAll() []error
}

func (handler *rsshandler) Fetch(id string) error {
	req := FetchRssRequest{Id: id, Result: make(chan error)}
	handler.FetchChan <- req
	return <-req.Result
}

func (handler rsshandler) FetchAll() []error {
	req := FetchAllRssRequest{Type: "rss", Result: make(chan []error)}
	handler.FetchAllChan <- req
	return <-req.Result
}

func NewRss(store *Store) Rss {
	handler := rsshandler{store: store}
	handler.FetchChan = make(chan FetchRssRequest, 100)
	handler.FetchAllChan = make(chan FetchAllRssRequest, 100)
	handler.reader = rss.Read
	go handler.serve()
	return &handler
}

func (handler *rsshandler) serve() error {
	for {
		select {
		case r := <-handler.FetchChan:
			r.Result <- handler.fetch(r.Id)
		case r := <-handler.FetchAllChan:
			r.Result <- handler.fetchAll(r.Type)
		case <-handler.CloseChan:
			return nil
		}
	}
	return nil
}

func (handler *rsshandler) fetchAll(feedType string) []error {
	ids, err := handler.store.GetByType(feedType)
	if err != nil {
		return []error{err}
	}
	result := make([]error, len(ids))
	for i := 0; i < len(ids); i++ {
		result[i] = handler.fetch(ids[i])
	}
	return result
}

func (handler *rsshandler) fetch(id string) error {
	feed, err := handler.store.Get(id)
	feed.LastError = ""
	feed.LastFetched = time.Now()
	if err != nil {
		return err
	}
	curr, err := rss.Read(feed.Url)
	if err != nil {
		feed.LastError = err.Error()
		handler.store.Put(&feed)
		return err
	}
	items, err := convertToItems(curr)
	if err != nil {
		feed.LastError = err.Error()
		handler.store.Put(&feed)
		return err
	}
	feed.AddNewItems(items)
	err = handler.store.Put(&feed)
	return err
}

func convertToItems(channel *rss.Channel) ([]Item, error) {
	items := make([]Item, len(channel.Item))
	for i := 0; i < len(channel.Item); i++ {
		from := channel.Item[i]
		to := Item{}
		to.Id = from.GUID
		to.Title = from.Title
		to.Description = from.Description
		to.Content = from.Content
		to.Type = from.Enclosure.Type
		to.Url = from.Enclosure.URL
		items[i] = to
	}
	return items, nil
}