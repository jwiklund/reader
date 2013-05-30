package reader

import (
	"github.com/jwiklund/reader/types"
	rss "github.com/ungerik/go-rss"
	"html"
	"time"
)

func NewRss(store types.Store) types.Rss {
	handler := rsshandler{store: store}
	handler.FetchChan = make(chan FetchRssRequest, 100)
	handler.FetchAllChan = make(chan FetchAllRssRequest, 100)
	handler.CloseChan = make(chan bool)
	handler.reader = rss.Read
	handler.store = store
	go handler.serve()
	return &handler
}

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

	store  types.Store
	reader func(string) (*rss.Channel, error)
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

func (handler rsshandler) Close() {
	handler.CloseChan <- true
}

func (handler *rsshandler) serve() error {
	for {
		select {
		case r := <-handler.FetchChan:
			r.Result <- handler.fetch(r.Id)
		case r := <-handler.FetchAllChan:
			r.Result <- handler.fetchAll(r.Type)
		case <-handler.CloseChan:
			close(handler.FetchChan)
			close(handler.FetchAllChan)
			close(handler.CloseChan)
			return nil
		}
	}
	return nil
}

func (handler *rsshandler) fetchAll(feedType string) []error {
	ids, err := handler.store.GetFeedByType(feedType)
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
	feed, err := handler.store.GetFeed(id)
	if err != nil {
		return err
	}
	feed.LastError = ""
	feed.LastFetched = time.Now().Format(types.DateFormat)
	items, err := FetchRss(feed.Url)
	if err != nil {
		feed.LastError = err.Error()
		handler.store.PutFeed(feed)
		return err
	}
	feed.AddNewItems(items)
	err = handler.store.PutFeed(feed)
	return err
}

func FetchRss(url string) ([]types.Item, error) {
	curr, err := rss.Read(url)
	if err != nil {
		return nil, err
	}
	return convertToItems(curr)
}

func convertToItems(channel *rss.Channel) ([]types.Item, error) {
	items := make([]types.Item, len(channel.Item))
	for i := 0; i < len(channel.Item); i++ {
		from := channel.Item[i]
		to := types.Item{}
		to.Id = from.GUID
		to.Title = from.Title
		to.Description = html.UnescapeString(from.Description)
		to.Content = from.Content
		time, err := from.PubDate.Parse()
		if err != nil {
			to.Published = string(from.PubDate)
		} else {
			to.Published = time.Format(types.DateFormat)
		}
		if from.Enclosure.URL != "" {
			to.Type = from.Enclosure.Type
			to.Url = from.Enclosure.URL
		} else {
			to.Type = "Link"
			to.Url = from.Link
		}
		items[i] = to
	}
	return items, nil
}
