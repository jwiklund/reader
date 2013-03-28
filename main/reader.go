package main

import (
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
	run()
}
