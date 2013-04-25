package main

import (
	"flag"
	"fmt"
	"github.com/jwiklund/reader"
	"github.com/jwiklund/reader/rss"
)

func main() {
	var cmd = flag.String("cmd", "fetch-rss", "command")
	var feedName = flag.String("feed", "xkcd", "the feed")
	flag.Parse()

	store := reader.NewStore("data")
	defer store.Close()
	feed, err := store.Get(*feedName)
	if err != nil {
		fmt.Println("Could not fetch " + feed.Id + " due to " + err.Error())
		return
	}
	if *cmd == "fetch-rss" {
		items, err := rss.Fetch(feed.Url)
		if err != nil {
			fmt.Println("Could not fetch " + feed.Url + " due to " + err.Error())
			return
		}
		for _, v := range items {
			fmt.Println(v.Id)
			fmt.Println(v.Description)
		}
		return
	} else if *cmd == "read-items" {
		for _, v := range feed.Items {
			fmt.Println(v.Id)
			fmt.Println(v.Description)
		}
		return
	}
	rss := reader.NewRss(store)
	defer store.Close()
	if *cmd == "update-rss" {
		rss.Fetch(*feedName)
	}
}
