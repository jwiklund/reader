package reader

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
)

func writeError(err error) []byte {
	return []byte("{\"Message\": \"" + err.Error() + "\"}")
}

func SetupResources(s Store, r Rss) {
	page := template.Must(template.ParseFiles("tmpl/page.html"))
	usage := func(w http.ResponseWriter, r *http.Request) {
		page.Execute(w, "")
	}
	feed := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/feed" && r.Method == "GET" {
			feeds, err := s.GetAllInfo()
			if err != nil {
				w.Write(writeError(err))
			}
			bytes, err := json.Marshal(feeds)
			if err != nil {
				w.Write(writeError(err))
			} else {
				w.Write(bytes)
			}
		} else {
			w.Write(writeError(errors.New("Method not allowed " + r.Method + " " + r.URL.Path)))
		}
	}
	http.HandleFunc("/", usage)
	http.HandleFunc("/usage", usage)
	http.HandleFunc("/feed", feed)
}
