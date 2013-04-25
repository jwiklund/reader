package reader

import (
	"encoding/json"
	"errors"
	"github.com/jwiklund/reader/types"
	"io/ioutil"
	"net/http"
	"strings"
)

type ResourceResponse struct {
	Status  string
	Message string
}

func respond(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		w.Write([]byte("{\"Status\": \"Error\", \"Message\": \"" + err.Error() + "\"}"))
		return
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		w.Write([]byte("{\"Status\": \"Error\", \"Message\": \"" + err.Error() + "\"}"))
	} else {
		w.Write([]byte("{\"Status\": \"Ok\", \"Data\": "))
		w.Write(bytes)
		w.Write([]byte("}"))
	}
}

func respondOk(w http.ResponseWriter, r *ResourceResponse) {
	bytes, err := json.Marshal(r)
	if err != nil {
		w.Write([]byte("{\"Status\": \"Error\", \"Message\": \"" + err.Error() + "\"}"))
	} else {
		w.Write(bytes)
	}
}

func createFeed(s types.Store, w http.ResponseWriter, r *http.Request) {
	feed := types.Feed{}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond(w, nil, err)
		return
	}
	err = json.Unmarshal(bytes, &feed)
	if err != nil {
		respond(w, nil, err)
		return
	}
	err = feed.ValidateNew()
	if err != nil {
		respond(w, nil, err)
		return
	}
	_, err = s.Get(feed.Id)
	if err == nil {
		respond(w, nil, errors.New("Feed already exists "+feed.Id))
		return
	}
	err = s.Put(&feed)
	if err != nil {
		respond(w, nil, err)
		return
	}
	respondOk(w, &ResourceResponse{"Ok", "Created " + feed.Id})
}

func refreshFeed(id string, rss types.Rss, w http.ResponseWriter) {
	err := rss.Fetch(id)
	if err == nil {
		respondOk(w, &ResourceResponse{"Ok", "Updated " + id})
	} else {
		respond(w, nil, err)
	}
}

func feedOperation(path string) (string, string) {
	if path == "" {
		return "create", ""
	}
	ind := strings.Index(path, "/")
	if ind == -1 {
		return "Unknown", path
	}
	return path[ind+1:], path[0:ind]
}

func SetupResources(s types.Store, rss types.Rss) {
	infoHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			feeds, err := s.GetAllInfo()
			respond(w, feeds, err)
			return
		}
		respond(w, nil, errors.New("Method not allowed "+r.Method+" "+r.URL.Path))
	}
	feedHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if r.Method == "GET" {
			id := r.URL.Path[6:]
			feed, err := s.Get(id)
			respond(w, feed, err)
			return
		}
		if r.Method == "POST" {
			op, id := feedOperation(r.URL.Path[6:])
			if op == "create" {
				createFeed(s, w, r)
				return
			} else if op == "refresh" {
				refreshFeed(id, rss, w)
				return
			}
		}
		if r.Method == "PUT" {
			id := r.URL.Path[6:]
			if id == "" {
				respond(w, nil, errors.New("Can not put to empty feed id"))
				return
			}
			feed := types.Feed{}
			bytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				respond(w, nil, err)
				return
			}
			feed.Id = id
			err = json.Unmarshal(bytes, &feed)
			if err != nil {
				respond(w, nil, err)
				return
			}
			old, err := s.Get(id)
			if err != nil {
				if len(err.Error()) < 9 || err.Error()[0:9] != "not found" {
					respond(w, nil, err)
					return
				}
			}
			if err != nil {
				old = &feed
			} else {
				old.Title = feed.Title
				old.Url = feed.Url
				old.Type = feed.Type
			}
			err = s.Put(old)
			if err != nil {
				respond(w, nil, err)
				return
			}
			respondOk(w, &ResourceResponse{"Ok", "Updated " + id})
			return
		}
		respond(w, nil, errors.New("Method not allowed "+r.Method+" "+r.URL.Path))
	}
	readHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			items, err := s.GetByUser("jwiklund")
			if err != nil {
				respond(w, nil, err)
				return
			}
			respond(w, items, nil)
			return
		}
		respond(w, nil, errors.New("Method not allowed "+r.Method+" "+r.URL.Path))
	}
	http.HandleFunc("/feed", infoHandler)
	http.HandleFunc("/read", readHandler)
	http.HandleFunc("/feed/", feedHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/html/index.html", http.StatusMovedPermanently)
	})
	dir := http.FileServer(http.Dir("."))
	http.Handle("/css/", dir)
	http.Handle("/js/", dir)
	http.Handle("/html/", dir)
}
