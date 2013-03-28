package reader

import (
	"testing"
)

func TestPutGetFeed(t *testing.T) {
	s := NewStore()
	go ServeStore(":memory:", s)
	defer s.Close()

	err := s.Put(&Feed{Url: "http://localhost/", Id: "feed", Title: "Title"})
	if err != nil {
		t.Fatal(err.Error())
	}
	feed, err := s.Get("feed")
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
	defer s.Close()

	feed := Feed{Url: "http://localhost/", Id: "feed", Title: "Title"}
	err := s.Put(&feed)
	if err != nil {
		t.Fatalf(err.Error())
	}
	feed.Title = "Title2"
	feed.Url = "http://localhost2/"
	err = s.Put(&feed)
	if err != nil {
		t.Fatal(err.Error())
	}
	feed, err = s.Get("feed")
	if err != nil {
		t.Fatal(err.Error())
	}
	if feed.Title != "Title2" || feed.Url != "http://localhost2/" {
		t.Fatal("Wrong feed " + feed.String())
	}
}

func TestItemsInFeed(t *testing.T) {
	s := NewStore()
	go ServeStore(":memory:", s)
	defer s.Close()

	feed := Feed{Id: "feed1", Title: "Feed1", Type: "rss1", Url: "http://localhost/1"}
	t.Log("Original 1 " + feed.String())
	err := s.Put(&feed)
	if err != nil {
		t.Fatal(err.Error())
	}
	feed = Feed{Id: "feed2", Title: "Feed2", Type: "rss2", Url: "http://localhost/2"}
	feed.AddItem(&Item{Id: "hash", Url: "http://localhost/2/item1"})
	t.Log("Original 2 " + feed.String())
	err = s.Put(&feed)
	if err != nil {
		t.Fatalf(err.Error())
	}

	feed, err = s.Get("feed1")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("feed1 " + feed.String())
	if feed.Items != nil {
		t.Fatal("feed 1 got some items")
	}
	feed, err = s.Get("feed2")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("feed2 " + feed.String())
	if len(feed.Items) != 1 {
		t.Fatal("Feed2 lost items")
	}
	if feed.Items[0].Id != "hash" || feed.Items[0].Url != "http://localhost/2/item1" {
		t.Fatalf("Wrong item : %s %s", feed.Items[0].Id, feed.Items[0].Url)
	}

	ids, err := s.GetByType("rss1")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(len(ids))
	if ids[0] != "feed1" {
		t.Fatal("Wrong feeds by type1 " + ids[0])
	}
	ids, err = s.GetByType("rss2")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(len(ids))
	if ids[0] != "feed2" {
		t.Fatal("Wrong feeds by type2 " + ids[0])
	}
}
