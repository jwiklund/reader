package main

import (
	"flag"
	"fmt"
	"github.com/jwiklund/reader/reader"
)

func main() {
	var cmd = flag.String("cmd", "fetch-rss", "command")
	var feedName = flag.String("feed", "xkcd", "the feed")
	flag.Parse()

	r := reader.NewReader("data")
	defer r.Close()
	feed, err := r.GetStore().GetFeed(*feedName)
	if err != nil {
		fmt.Println("Could not fetch " + *feedName + " due to " + err.Error())
		return
	}
	if *cmd == "fetch-rss" {
		items, err := reader.FetchRss(feed.Url)
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
	if *cmd == "update-rss" {
		r.GetRss().Fetch(*feedName)
	}
}
