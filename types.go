package reader

import (
	"strconv"
	"time"
)

type Item struct {
	Id          string
	Title       string
	Description string
	Content     string
	Type        string
	Url         string
}

type Feed struct {
	Id          string
	Title       string
	Type        string
	Url         string
	LastFetched time.Time
	LastError   string
	Items       []Item
}

type ItemIdList []string

func (f *Feed) String() string {
	return "Feed(" + f.Id + ", title=" + f.Title + ", url=" + f.Url + ", len(items)=" + strconv.Itoa(len(f.Items)) + ")"
}

func (f *Feed) AddItem(item *Item) {
	f.Items = append(f.Items, *item)
}

func (f *Feed) ItemIdMap() map[string]*Item {
	res := make(map[string]*Item)
	for i := 0; i < len(f.Items); i++ {
		res[f.Items[i].Id] = &f.Items[i]
	}
	return res
}

func (f *Feed) ItemIds() ItemIdList {
	res := make([]string, len(f.Items))
	for i := 0; i < len(f.Items); i++ {
		res[i] = f.Items[i].Id
	}
	return ItemIdList(res)
}

func (f ItemIdList) String() string {
	res := ""
	s := []string(f)
	for i := 0; i < len(s); i++ {
		if res == "" {
			res = s[i]
		} else {
			res = res + ", " + s[i]
		}
	}
	return res
}

func (f ItemIdList) Equals(ids []string) bool {
	s := []string(f)
	if len(s) != len(ids) {
		return false
	}
	for i := 0; i < len(ids); i++ {
		if s[i] != ids[i] {
			return false
		}
	}
	return true

}

func (f *Feed) AddNewItems(items []Item) {
	old := f.ItemIdMap()
	for i := 0; i < len(items); i++ {
		prev, ok := old[items[i].Id]
		if ok {
			*prev = items[i]
		} else {
			f.Items = append(f.Items, items[i])
		}
	}
}
