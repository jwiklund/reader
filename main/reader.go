package main

import (
	"encoding/json"
	"fmt"
	"github.com/jwiklund/reader"
	"net/http"
)

func run() {
	s := reader.NewStore("data")
	defer s.Close()
	r := reader.NewRss(s)
	defer r.Close()
	reader.SetupResources(s, r)
	fmt.Println("Start listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	i := []reader.Item{reader.Item{Id: "haha"}}
	b, e := json.Marshal(i)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println(string(b))
	}
}
