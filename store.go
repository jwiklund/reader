package reader

import (
	"encoding/json"
	"errors"
	sqlite "github.com/gwenn/gosqlite"
)

type FeedOrError struct {
	Feed  Feed
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

type Store struct {
	GetFeedChan chan *GetFeedRequest
	GetTypeChan chan *GetTypeRequest
	PutFeedChan chan *PutFeedRequest
	CloseChan   chan bool
}

func NewStore() *Store {
	store := Store{}
	store.GetFeedChan = make(chan *GetFeedRequest, 100)
	store.GetTypeChan = make(chan *GetTypeRequest, 100)
	store.PutFeedChan = make(chan *PutFeedRequest, 100)
	store.CloseChan = make(chan bool)
	return &store
}

func (s *Store) Put(feed *Feed) error {
	p := PutFeedRequest{*feed, make(chan error)}
	s.PutFeedChan <- &p
	r := <-p.Error
	return r
}

func (s *Store) Get(id string) (Feed, error) {
	g := GetFeedRequest{id, make(chan FeedOrError)}
	s.GetFeedChan <- &g
	r := <-g.Feed
	return r.Feed, r.Error
}

func (s *Store) GetByType(feedType string) ([]string, error) {
	g := GetTypeRequest{feedType, make(chan IdOrError)}
	s.GetTypeChan <- &g
	r := <-g.Feed
	return r.Id, r.Error
}

func (s *Store) Close() {
	s.CloseChan <- true
}

func ServeStore(path string, store *Store) error {
	db, err := sqlite.Open(path)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Exec("create table if not exists feed(id string primary key, type string, json string)")
	if err != nil {
		return err
	}
	for {
		select {
		case c := <-store.PutFeedChan:
			c.Error <- putFeed(db, c.Feed)
		case r := <-store.GetFeedChan:
			r.Feed <- getFeed(db, r.Id)
		case c := <-store.GetTypeChan:
			c.Feed <- getType(db, c.Type)
		case <-store.CloseChan:
			return nil
		}
	}
	return nil
}

func putFeed(conn *sqlite.Conn, feed Feed) error {
	bytes, err := json.Marshal(feed)
	if err != nil {
		return err
	}
	stmt, err := conn.Prepare("insert or replace into feed(id, type, json) values (?, ?, ?)")
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

func getFeed(conn *sqlite.Conn, id string) FeedOrError {
	stmt, err := conn.Prepare("select json from feed where id = ?")
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
	return FeedOrError{feed, nil}
}

func getType(conn *sqlite.Conn, feedType string) IdOrError {
	stmt, err := conn.Prepare("select id from feed where type = ?")
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
