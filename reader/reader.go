package main

import (
	"fmt"
	"github.com/jwiklund/reader"
	"log"
	"net/http"
)

func run() {
	s := reader.NewStore("data")
	defer s.Close()
	r := reader.NewRss(s)
	defer r.Close()
	reader.SetupResources(s, r)
	fmt.Println("Start listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	run()
}
