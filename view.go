package reader

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
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
		w.Write(bytes)
	}
}

func SetupResources(s Store, rss Rss) {
	usagePage := template.Must(template.ParseFiles("html/usage.html"))
	usageHandler := func(w http.ResponseWriter, r *http.Request) {
		usagePage.Execute(w, "")
	}
	infoHandler := func(w http.ResponseWriter, r *http.Request) {
		feeds, err := s.GetAllInfo()
		respond(w, feeds, err)
	}
	feedHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if r.Method == "GET" {
			id := r.URL.Path[6:]
			feed, err := s.Get(id)
			respond(w, feed, err)
			return
		}
		if r.Method == "PUT" {
			id := r.URL.Path[6:]
			if id == "" {
				respond(w, nil, errors.New("Can not put to empty feed id"))
				return
			}
			feed := Feed{}
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
			respond(w, ResourceResponse{"Ok", "Updated " + id}, err)
			return
		}
		respond(w, nil, errors.New("Method not allowed "+r.Method+" "+r.URL.Path))
	}
	refreshHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id := r.URL.Path[9:]
			err := rss.Fetch(id)
			respond(w, ResourceResponse{"Ok", "Refreshed " + id}, err)
			return
		}
		respond(w, nil, errors.New("Method not allowed "+r.Method+" "+r.URL.Path))
	}
	dir := http.Dir(".")
	fileHandler := func(w http.ResponseWriter, r *http.Request) {
		file, err := dir.Open("html/index.html")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		buf, err := ioutil.ReadAll(file)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(buf)
	}
	http.HandleFunc("/usage", usageHandler)
	http.HandleFunc("/feed", infoHandler)
	http.HandleFunc("/feed/", feedHandler)
	http.HandleFunc("/refresh/", refreshHandler)
	http.HandleFunc("/", fileHandler)
	http.Handle("/js/", http.FileServer(dir))
}
