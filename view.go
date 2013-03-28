package reader

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type ResourceResponse struct {
	Status  string
	Message string
}

func writeError(err error) []byte {
	json := "{\"Status\": \"Error\", \"Message\": \"" + err.Error() + "\"}"
	return []byte(json)
}

func SetupResources(s Store, rss Rss) {
	usagePage := template.Must(template.ParseFiles("tmpl/usage.html"))
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/usage" {
			usagePage.Execute(w, "")
			return
		}
		respond := func(data interface{}, err error) {
			if err != nil {
				w.Write(writeError(err))
				return
			}
			bytes, err := json.Marshal(data)
			if err != nil {
				w.Write(writeError(err))
			} else {
				w.Write(bytes)
			}
		}
		urlMatch := func(start string) bool {
			if len(r.URL.Path) < len(start) {
				return false
			}
			return r.URL.Path[0:len(start)] == start
		}
		w.Header().Add("Content-Type", "application/json")
		if r.Method == "GET" {
			if r.URL.Path == "/feed" || r.URL.Path == "/feed/" {
				feeds, err := s.GetAllInfo()
				respond(feeds, err)
				return
			} else if urlMatch("/feed/") {
				id := r.URL.Path[6:]
				feed, err := s.Get(id)
				respond(feed, err)
				return
			}
		}
		if r.Method == "PUT" {
			if urlMatch("/feed/") {
				id := r.URL.Path[6:]
				if id == "" {
					respond(nil, errors.New("Can not put to empty feed id"))
					return
				}
				feed := Feed{}
				bytes, err := ioutil.ReadAll(r.Body)
				if err != nil {
					respond(nil, err)
					return
				}
				feed.Id = id
				err = json.Unmarshal(bytes, &feed)
				if err != nil {
					respond(nil, err)
					return
				}
				old, err := s.Get(id)
				if err != nil && err.Error() != "not found" {
					respond(nil, err)
					return
				}
				if err != nil {
					old = &feed
				} else {
					old.Title = feed.Title
					old.Url = feed.Url
					old.Type = feed.Type
				}
				err = s.Put(old)
				respond(ResourceResponse{"Ok", "Updated " + id}, err)
				return
			}
		}
		if r.Method == "POST" {
			if urlMatch("/feed/") {
				id := r.URL.Path[6:]
				if id == "" {
					respond(nil, errors.New("Can not operate on empty feed id"))
					return
				}
				ind := strings.Index(id, "/")
				if ind == -1 {
					respond(nil, errors.New("Can not perform no operation on "+id))
					return
				}
				op := id[ind+1:]
				id = id[0:ind]
				if op == "refresh" {
					err := rss.Fetch(id)
					respond(ResourceResponse{"Ok", "Refreshed " + id}, err)
					return
				}
				respond(nil, errors.New("Can not perform operation "+op+" on "+id))
				return
			}
		}
		respond(nil, errors.New("Method not allowed "+r.Method+" "+r.URL.Path))
	}
	http.HandleFunc("/", handler)
}
