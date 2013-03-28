package reader

import (
	"encoding/json"
	"errors"
	sqlite "github.com/gwenn/gosqlite"
)

type Store interface {
	Get(id string) (*Feed, error)
	GetByType(feedType string) ([]string, error)
	Put(feed *Feed) error
	Close()
}

func NewStore(path string) Store {
	store := store{}
	store.GetFeedChan = make(chan *GetFeedRequest, 100)
	store.GetTypeChan = make(chan *GetTypeRequest, 100)
	store.PutFeedChan = make(chan *PutFeedRequest, 100)
	store.CloseChan = make(chan bool)
	store.path = path
	go store.serve()
	return &store
}

func (s *store) Put(feed *Feed) error {
	p := PutFeedRequest{*feed, make(chan error)}
	s.PutFeedChan <- &p
	r := <-p.Error
	return r
}

func (s *store) Get(id string) (*Feed, error) {
	g := GetFeedRequest{id, make(chan FeedOrError)}
	s.GetFeedChan <- &g
	r := <-g.Feed
	return r.Feed, r.Error
}

func (s *store) GetByType(feedType string) ([]string, error) {
	g := GetTypeRequest{feedType, make(chan IdOrError)}
	s.GetTypeChan <- &g
	r := <-g.Feed
	return r.Id, r.Error
}

func (s *store) Close() {
	s.CloseChan <- true
}

type store struct {
	GetFeedChan chan *GetFeedRequest
	GetTypeChan chan *GetTypeRequest
	PutFeedChan chan *PutFeedRequest
	CloseChan   chan bool

	path string
	conn *sqlite.Conn
}

type FeedOrError struct {
	Feed  *Feed
	Error error
}

type GetFeedRequest struct {
	Id   string
	Feed chan FeedOrError
}

type IdOrError struct {
	Id    []string
	Error error
}

type GetTypeRequest struct {
	Type string
	Feed chan IdOrError
}

type PutFeedRequest struct {
	Feed  Feed
	Error chan error
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
	s.conn = db
	for {
		select {
		case c := <-s.PutFeedChan:
			c.Error <- s.putFeed(c.Feed)
		case r := <-s.GetFeedChan:
			r.Feed <- s.getFeed(r.Id)
		case c := <-s.GetTypeChan:
			c.Feed <- s.getType(c.Type)
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
	bytes, err := json.Marshal(feed)
	if err != nil {
		return err
	}
	stmt, err := s.conn.Prepare("insert or replace into feed(id, type, json) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Finalize()
	err = stmt.Exec(feed.Id, feed.Type, bytes)
	if err != nil {
		stmt.Finalize()
		return err
	}
	return nil
}

func (s *store) getFeed(id string) FeedOrError {
	stmt, err := s.conn.Prepare("select json from feed where id = ?")
	if err != nil {
		return FeedOrError{Error: err}
	}
	defer stmt.Finalize()
	feed := Feed{Id: ""}
	err = stmt.Select(func(s *sqlite.Stmt) error {
		bytes := []byte{}
		s.Scan(&bytes)
		return json.Unmarshal(bytes, &feed)
	}, id)
	if err != nil {
		return FeedOrError{Error: err}
	}
	if feed.Id == "" {
		return FeedOrError{Error: errors.New("not found")}
	}
	return FeedOrError{&feed, nil}
}

func (s *store) getType(feedType string) IdOrError {
	stmt, err := s.conn.Prepare("select id from feed where type = ?")
	if err != nil {
		return IdOrError{Error: err}
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
		return IdOrError{Error: err}
	}
	return IdOrError{Id: ids}
}
