package reader

import (
	"fmt"
	"strconv"
	rss "github.com/ungerik/go-rss"
	sqlite "github.com/gwenn/gosqlite"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	channel, err := rss.Read("http://feeds.feedburner.com/DilbertDailyStrip")
	if err != nil {
		fmt.Println("Could not fetch feed " + err.Error())
		return
	}
	fmt.Println("Title: " + channel.Title)
	fmt.Println("Nr Items: " + strconv.Itoa(len(channel.Item)))
	db, err := sqlite.Open("data")
	if err != nil {
		fmt.Println("Could not open database " + err.Error())
		return
	}
	defer db.Close()
	err = db.Exec("create table if not exists item(id string primary key)")
	check(err)
	q, err := db.Prepare("select id from item where id like ?")
	check(err)
	defer q.Finalize()
	ids := map[string] bool {}
	q.Select(func(s *sqlite.Stmt) (err error) {
		var id string
		if err = s.Scan(&id) ; err != nil {
			panic(err)
			return err
		}
		ids[id] = true
		return
		}, "dilbert:%")
	err = db.Begin()
	check(err)
	i, err := db.Prepare("insert into item values(?)")
	check(err)
	defer i.Finalize()
	for _, item := range channel.Item {
		if !ids["dilbert:" + item.GUID] {
			fmt.Println(item.Title + " is new")
			i.Insert("dilbert:" + item.GUID)
		}
	}
	db.Commit()
}
