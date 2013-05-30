package reader

import (
	"github.com/jwiklund/reader/types"
	"strconv"
	"testing"
)

func TestPutGetFeed(t *testing.T) {
	s := NewStore(":memory:")
	defer s.Close()

	err := s.Put(&types.Feed{Url: "http://localhost/", Id: "feed", Title: "Title"})
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
	s := NewStore(":memory:")
	defer s.Close()

	feed := &types.Feed{Url: "http://localhost/", Id: "feed", Title: "Title"}
	err := s.Put(feed)
	if err != nil {
		t.Fatalf(err.Error())
	}
	feed.Title = "Title2"
	feed.Url = "http://localhost2/"
	err = s.Put(feed)
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
	s := NewStore(":memory:")
	defer s.Close()

	feed := &types.Feed{Id: "feed1", Title: "Feed1", Type: "rss1", Url: "http://localhost/1"}
	t.Log("Original 1 " + feed.String())
	err := s.Put(feed)
	if err != nil {
		t.Fatal(err.Error())
	}
	feed = &types.Feed{Id: "feed2", Title: "Feed2", Type: "rss2", Url: "http://localhost/2"}
	feed.AddItem(&types.Item{Id: "hash", Url: "http://localhost/2/item1"})
	t.Log("Original 2 " + feed.String())
	err = s.Put(feed)
	if err != nil {
		t.Fatalf(err.Error())
	}

	feed, err = s.Get("feed1")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("feed1 " + feed.String())
	if feed.Items != nil {
		t.Fatal("feed got to may items")
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

func TestGetInfo(t *testing.T) {
	s := NewStore(":memory:")
	defer s.Close()

	i1 := []types.Item{types.Item{Id: "Id1.1"}}
	i2 := []types.Item{types.Item{Id: "Id2.1"}}
	s.Put(&types.Feed{Id: "id1", Items: i1})
	s.Put(&types.Feed{Id: "id2", Items: i2})
	info, err := s.GetAllInfo()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(info) != 2 {
		t.Fatal("Wrong number of infos " + strconv.Itoa(len(info)))
	}
	for i := 0; i < len(info); i++ {
		if len(info[i].Items) != 0 {
			t.Fatal("Got items from info " + info[i].ItemIds().String())
		}
	}
}
