package main

import (
	"fmt"
	rss "github.com/ungerik/go-rss"
	"strconv"
)

func main() {
	channel, err := rss.Read("http://feeds.feedburner.com/DilbertDailyStrip")
	if err != nil {
		fmt.Println("Could not fetch " + err.Error())
		return
	}
	fmt.Println("Title: " + channel.Title)
	fmt.Println("Nr Items: " + strconv.Itoa(len(channel.Item)))
}
