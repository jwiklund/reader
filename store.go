package reader

import (
	"encoding/json"
	"errors"
	sqlite "github.com/gwenn/gosqlite"
)

type Store interface {
	Get(id string) (*Feed, error)
	GetByType(feedType string) ([]string, error)
	GetAllInfo() ([]Feed, error)
	Put(feed *Feed) error
	Close()
}

func NewStore(path string) Store {
	store := store{}
	store.GetFeedChan = make(chan GetFeedRequest, 100)
	store.GetTypeChan = make(chan GetTypeRequest, 100)
	store.GetInfoChan = make(chan GetInfoRequest, 100)
	store.PutFeedChan = make(chan PutFeedRequest, 100)
	store.CloseChan = make(chan bool)
	store.path = path
	go store.serve()
	return &store
}

func (s *store) Put(feed *Feed) error {
	p := PutFeedRequest{*feed, make(chan error)}
	s.PutFeedChan <- p
	r := <-p.Response
	return r
}

func (s *store) Get(id string) (*Feed, error) {
	g := GetFeedRequest{id, make(chan FeedResponse)}
	s.GetFeedChan <- g
	r := <-g.Response
	return r.Feed, r.Error
}

func (s *store) GetByType(feedType string) ([]string, error) {
	g := GetTypeRequest{feedType, make(chan TypeResponse)}
	s.GetTypeChan <- g
	r := <-g.Response
	return r.Id, r.Error
}

func (s *store) GetAllInfo() ([]Feed, error) {
	g := GetInfoRequest{make(chan InfoResponse)}
	s.GetInfoChan <- g
	r := <-g.Response
	return r.Feed, r.Error
}

func (s *store) Close() {
	s.CloseChan <- true
}

type store struct {
	GetFeedChan chan GetFeedRequest
	GetTypeChan chan GetTypeRequest
	GetInfoChan chan GetInfoRequest
	PutFeedChan chan PutFeedRequest
	CloseChan   chan bool

	path string
	conn *sqlite.Conn
}

type FeedResponse struct {
	Feed  *Feed
	Error error
}

type GetFeedRequest struct {
	Id       string
	Response chan FeedResponse
}

type TypeResponse struct {
	Id    []string
	Error error
}

type GetTypeRequest struct {
	Type     string
	Response chan TypeResponse
}

type InfoResponse struct {
	Feed  []Feed
	Error error
}

type GetInfoRequest struct {
	Response chan InfoResponse
}

type PutFeedRequest struct {
	Feed     Feed
	Response chan error
}

func (s *store) serve() error {
	db, err := sqlite.Open(s.path)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Exec("create table if not exists feed(id string primary key, type string, json string)")
	if err != nil {
		return err
	}
	err = db.Exec("create table if not exists item(id string primary key, json string)")
	if err != nil {
		return err
	}
	s.conn = db
	for {
		select {
		case c := <-s.PutFeedChan:
			c.Response <- s.putFeed(c.Feed)
		case r := <-s.GetFeedChan:
			r.Response <- s.getFeed(r.Id)
		case c := <-s.GetTypeChan:
			c.Response <- s.getType(c.Type)
		case c := <-s.GetInfoChan:
			c.Response <- s.getInfo()
		case <-s.CloseChan:
			close(s.PutFeedChan)
			close(s.GetFeedChan)
			close(s.GetTypeChan)
			close(s.CloseChan)
			return nil
		}
	}
	return nil
}

func (s *store) putFeed(feed Feed) error {
	items := feed.Items
	feed.Items = []Item{}
	info_bytes, err := json.Marshal(feed)
	if err != nil {
		return err
	}
	item_bytes, err := json.Marshal(items)
	if err != nil {
		return err
	}
	stmt, err := s.conn.Prepare("insert or replace into feed(id, type, json) values (?, ?, ?)")
	if err != nil {
		return err
	}
	s.conn.Begin()
	err = stmt.Exec(feed.Id, feed.Type, info_bytes)
	stmt.Finalize()
	if err != nil {
		s.conn.Rollback()
		return err
	}
	stmt, err = s.conn.Prepare("insert or replace into item(id, json) values (?, ?)")
	if err != nil {
		s.conn.Rollback()
		return err
	}
	err = stmt.Exec(feed.Id, item_bytes)
	stmt.Finalize()
	if err != nil {
		s.conn.Rollback()
		return err
	}
	s.conn.Commit()
	return nil
}

func (s *store) getFeed(id string) FeedResponse {
	stmt, err := s.conn.Prepare("select json from feed where id = ?")
	if err != nil {
		return FeedResponse{Error: err}
	}
	feed := Feed{Id: ""}
	err = stmt.Select(func(s *sqlite.Stmt) error {
		bytes := []byte{}
		s.Scan(&bytes)
		return json.Unmarshal(bytes, &feed)
	}, id)
	stmt.Finalize()
	if err != nil {
		return FeedResponse{Error: err}
	}
	if feed.Id == "" {
		return FeedResponse{Error: errors.New("not found")}
	}
	stmt, err = s.conn.Prepare("select json from item where id = ?")
	if err != nil {
		return FeedResponse{Error: err}
	}
	err = stmt.Select(func(s *sqlite.Stmt) error {
		bytes := []byte{}
		s.Scan(&bytes)
		return json.Unmarshal(bytes, &feed.Items)
	}, id)
	stmt.Finalize()
	if err != nil {
		return FeedResponse{Error: err}
	}
	return FeedResponse{&feed, nil}
}

func (s *store) getType(feedType string) TypeResponse {
	stmt, err := s.conn.Prepare("select id from feed where type = ?")
	if err != nil {
		return TypeResponse{Error: err}
	}
	defer stmt.Finalize()
	ids := []string{}
	err = stmt.Select(func(s *sqlite.Stmt) error {
		var id string
		s.Scan(&id)
		ids = append(ids, id)
		return nil
	}, feedType)
	if err != nil {
		return TypeResponse{Error: err}
	}
	return TypeResponse{ids, nil}
}

func (s *store) getInfo() InfoResponse {
	stmt, err := s.conn.Prepare("select json from feed")
	if err != nil {
		return InfoResponse{Error: err}
	}
	f := Feed{}
	feed := []Feed{}
	bytes := []byte{}
	err = stmt.Select(func(s *sqlite.Stmt) error {
		s.Scan(&bytes)
		err := json.Unmarshal(bytes, &f)
		if err != nil {
			return err
		}
		feed = append(feed, f)
		return err
	})
	stmt.Finalize()
	if err != nil {
		return InfoResponse{Error: err}
	}
	return InfoResponse{feed, nil}
}
