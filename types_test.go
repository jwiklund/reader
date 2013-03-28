package reader

import (
	"testing"
)

func TestAdd(t *testing.T) {
	feed := Feed{}
	feed.AddItem(&Item{Id: "1"})
	feed.AddItem(&Item{Id: "2"})
	if !feed.ItemIds().Equals([]string{"1", "2"}) {
		t.Error("Expected order 1, 2 but got " + feed.ItemIds().String())
	}
}

func TestAddNew(t *testing.T) {
	feed := Feed{}
	feed.AddItem(&Item{Id: "1", Title: "first"})
	feed.AddItem(&Item{Id: "2", Title: "second"})
	feed.AddNewItems([]Item{Item{Id: "2", Title: "second upd"}, Item{Id: "3"}})
	if !feed.ItemIds().Equals([]string{"1", "2", "3"}) {
		t.Fatal("Expected order 1, 2, 3 but got " + feed.ItemIds().String())
	}
	if feed.Items[1].Title != "second upd" {
		t.Fatal("Expected title to be updated but got " + feed.Items[1].Title)
	}
}
