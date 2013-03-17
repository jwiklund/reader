package reader

import (
	"testing"
)

func TestPutGetFeed(t *testing.T) {
	s := NewStore()
	go ServeStore(":memory:", s)
	err := s.Put(&Feed{Url: "http://localhost/", Id: "feed", Title: "Title"})
	if err != nil {
		t.Fatal(err.Error())
	}
	feed, err := s.Get("feed")
	s.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
	if feed.Title != "Title" {
		t.Fatal("Wrong title " + feed.Title)
	}
	if feed.Url != "http://localhost/" {
		t.Fatal("Wrong url " + feed.Url)
	}
}

func TestUpdateFeed(t *testing.T) {
	s := NewStore()
	go ServeStore(":memory:", s)
	feed := Feed{Url: "http://localhost/", Id: "feed", Title: "Title"}
	err := s.Put(&feed)
	if err != nil {
		t.Fatalf(err.Error())
	}
	feed.Title = "Title2"
	feed.Url = "http://localhost2/"
	err = s.Put(&feed)
	if err != nil {
		t.Fatalf(err.Error())
	}
	feed, err = s.Get("feed")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if feed.Title != "Title2" || feed.Url != "http://localhost2/" {
		t.Fatalf("Wrong feed " + feed.String())
	}
}
