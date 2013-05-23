package main

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/jwiklund/reader"
	"github.com/jwiklund/reader/types"
	"log"
	"net/http"
)

func main() {
	s := reader.NewStore("data")
	defer s.Close()
	r := reader.NewRss(s)
	defer r.Close()

	registerStaticFiles()
	registerFeedService(service{s, r})

	config := swagger.Config{
		WebServicesUrl:  "http://localhost:8080",
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "swagger-ui",
		WebServices:     restful.RegisteredWebServices()}
	swagger.InstallSwaggerService(config)

	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type service struct {
	store types.Store
	rss   types.Rss
}

type Status struct {
	Status  string
	Message string
}

type Feeds struct {
	Status string
	Feeds  []types.Feed
}

type Items struct {
	Status string
	Items  []types.Item
}

type Feed struct {
	Status string
	Feed   types.Feed
}

func (s *service) getFeeds(request *restful.Request, response *restful.Response) {
	feeds, err := s.store.GetAllInfo()
	if err != nil {
		response.WriteEntity(Status{"fail", err.Error()})
	} else {
		response.WriteEntity(Feeds{"ok", feeds})
	}
}

func (s *service) getAllUserItems(request *restful.Request, response *restful.Response) {
	items, err := s.store.GetByUser(request.PathParameter(("user-id")))
	if err != nil {
		response.WriteEntity(Status{"fail", err.Error()})
	} else {
		response.WriteEntity(Items{"ok", items})
	}
}

func (s *service) getFeed(request *restful.Request, response *restful.Response) {
	feed, err := s.store.Get(request.PathParameter("feed-id"))
	if err != nil {
		response.WriteEntity(Status{"fail", err.Error()})
	} else {
		response.WriteEntity(Feed{"ok", *feed})
	}
}

func (s *service) createFeed(request *restful.Request, response *restful.Response) {
	feed := types.Feed{}
	err := request.ReadEntity(&feed)
	if err != nil {
		response.WriteEntity(Status{"fail", "Could not parse Feed: " + err.Error()})
		return
	}
	err = feed.ValidateNew()
	if err != nil {
		response.WriteEntity(Status{"fail", "Could not validate Feed: " + err.Error()})
		return
	}
	_, err = s.store.Get(feed.Id)
	if err == nil {
		response.WriteEntity(Status{"fail", "Feed already exists"})
		return
	}
	err = s.store.Put(&feed)
	if err != nil {
		response.WriteEntity(Status{"fail", "Could not store Feed: " + err.Error()})
		return
	}
	response.WriteEntity(Status{"ok", "Feed created"})
}

func (s *service) refreshFeed(request *restful.Request, response *restful.Response) {
	err := s.rss.Fetch(request.PathParameter("feed-id"))
	if err != nil {
		response.WriteEntity(Status{"fail", err.Error()})
	} else {
		response.WriteEntity(Status{"ok", "feed refreshed"})
	}
}

func registerFeedService(s service) {
	ws := new(restful.WebService)
	ws.Path("/feed").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Route(ws.GET("/").To(func(req *restful.Request, res *restful.Response) { s.getFeeds(req, res) }).
		Doc("get all feeds").
		Writes(Feeds{}))
	ws.Route(ws.GET("/user/{user-id}/all").To(func(req *restful.Request, res *restful.Response) { s.getAllUserItems(req, res) }).
		Doc("get all unread items for user").
		Param(ws.PathParameter("user-id", "the user id").DataType("string")).
		Writes(Feeds{}))
	ws.Route(ws.GET("/{feed-id}").To(func(req *restful.Request, res *restful.Response) { s.getFeed(req, res) }).
		Doc("get a feed").
		Param(ws.PathParameter("feed-id", "identifier of the feed").DataType("string")).
		Writes(Feed{}))
	ws.Route(ws.POST("/").To(func(req *restful.Request, res *restful.Response) { s.createFeed(req, res) }).
		Doc("create a feed").
		Param(ws.BodyParameter("feed", "the feed").DataType("types.Feed")).
		Writes(Status{}))
	ws.Route(ws.POST("/refresh/{feed-id}").To(func(req *restful.Request, res *restful.Response) { s.refreshFeed(req, res) }).
		Doc("refresh a feed").
		Param(ws.PathParameter("feed-id", "identifier of the feed").DataType("string")).
		Writes(Status{}))
	restful.Add(ws)
}

func registerStaticFiles() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/html/", http.StatusMovedPermanently)
	})
	dir := http.FileServer(http.Dir("."))
	http.Handle("/css/", dir)
	http.Handle("/js/", dir)
	http.Handle("/html/", dir)
}
