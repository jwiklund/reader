package reader

import (
	"errors"
	"fmt"
	sqlite "github.com/gwenn/gosqlite"
)

type Feed struct {
	Url   string
	Id    string
	Title string
	Item  []string
}

func (f Feed) String() string {
	return "Feed(" + f.Id + ", title=" + f.Title + ", url=" + f.Url + ")"
}

type FeedOrError struct {
	Feed Feed
	Error error
}

type GetFeedRequest struct {
	Id  string
	Feed chan FeedOrError
}

type PutFeedRequest struct {
	Feed Feed
	Error chan error
}

type Store struct {
	GetFeedChan chan *GetFeedRequest
	PutFeedChan chan *PutFeedRequest
	CloseChan chan bool
}

func NewStore() *Store {
	store := Store{}
	store.GetFeedChan = make(chan *GetFeedRequest, 100)
	store.PutFeedChan = make(chan *PutFeedRequest, 100)
	store.CloseChan = make(chan bool)
	return &store
}

func (s *Store) Put(feed *Feed) error {
	p := PutFeedRequest{*feed, make(chan error)}
	s.PutFeedChan <- &p
	r := <- p.Error
	return r
}

func (s *Store) Get(id string) (Feed, error) {
	g := GetFeedRequest{id, make(chan FeedOrError)}
	s.GetFeedChan <- &g
	r := <- g.Feed
	return r.Feed, r.Error
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
	err = db.Exec("create table if not exists feed(id string primary key, title string, url string)")
	if err != nil {
		return err
	}
	for {
		select {
		case r := <- store.GetFeedChan:
			r.Feed <- getFeed(db, r.Id)
		case c := <- store.PutFeedChan:
			c.Error <- putFeed(db, c.Feed)
		case <- store.CloseChan:
			return nil
		}
	}
	return nil
}

func putFeed(conn *sqlite.Conn, feed Feed) error {
	fe := getFeed(conn, feed.Id)
	exists := true
	if fe.Error != nil {
		if fe.Error.Error() != "not found" {
			fmt.Errorf("unexpected error in putFeed select " + fe.Error.Error())
			return fe.Error
		}
		exists = false
	}
	if exists {
		stmt, err := conn.Prepare("update feed set url = ?, title = ? where id = ?")
		if err != nil {
			fmt.Errorf("unexpected error in putFeed update " + err.Error())
			return err
		}
		defer stmt.Finalize()
		err = stmt.Exec(feed.Url, feed.Title, feed.Id)
		if err != nil {
			fmt.Errorf("unexpected error in putFeed update " + err.Error())
			return err
		}
	} else {
		stmt, err := conn.Prepare("insert into feed(id, url, title) values (?, ?, ?)")
		if err != nil {
			fmt.Errorf("unexpected error in putFeed insert" + err.Error())
			return err
		}
		defer stmt.Finalize()
		err = stmt.Exec(feed.Id, feed.Url, feed.Title)
		if err != nil {
			fmt.Errorf("unexpected error in putFeed insert" + err.Error())
			return err
		}
	}
	return nil
}

func getFeed(conn *sqlite.Conn, id string) FeedOrError {
	stmt, err := conn.Prepare("select id, url, title from feed where id = ?")
	if err != nil {
		return FeedOrError{Error: err}
	}
	defer stmt.Finalize()
	feed := Feed{Id: ""}
	err = stmt.Select(func(s *sqlite.Stmt) (error) {
		s.Scan(&feed.Id, &feed.Url, &feed.Title)
		return nil
		}, id)
	if err != nil {
		return FeedOrError{Error: err}
	}
	if feed.Id == "" {
		return FeedOrError{Error: errors.New("not found")}
	}
	return FeedOrError{feed, nil}
}
