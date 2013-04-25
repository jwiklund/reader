package types

import (
	"errors"
	"strconv"
)

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

func (f *Feed) ValidateNew() error {
	if f.Id == "" || f.Id == "id" {
		return errors.New("Id is required")
	}
	if f.Url == "" || f.Url == "url" {
		return errors.New("Url is required")
	}
	if f.Title == "" {
		f.Title = f.Id
	}
	if f.Type == "" {
		f.Type = "rss"
	}
	return nil
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
