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

func (s *service) getFeeds(request *restful.Request, response *restful.Response) {
	feeds, err := s.store.GetAllInfo()
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		response.WriteEntity(feeds)
	}
}

func (s *service) getFeed(request *restful.Request, response *restful.Response) {
	feed, err := s.store.Get(request.PathParameter("feed-id"))
	if err != nil {
		response.WriteError(http.StatusNotFound, err)
	} else {
		response.WriteEntity(feed)
	}
}

func (s *service) refreshFeed(request *restful.Request, response *restful.Response) {
	err := s.rss.Fetch(request.PathParameter("feed-id"))
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		response.WriteEntity(Status{"ok", "feed refreshed"})
	}
}

func registerFeedService(s service) {
	ws := new(restful.WebService)
	ws.Path("/feeds").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Route(ws.GET("/").To(func(req *restful.Request, res *restful.Response) { s.getFeeds(req, res) }).
		Doc("get all feeds").
		Writes(types.Feed{}))
	restful.Add(ws)
	ws = new(restful.WebService)
	ws.Path("/feed").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Route(ws.GET("/{feed-id}").To(func(req *restful.Request, res *restful.Response) { s.getFeed(req, res) }).
		Doc("get a feed").
		Param(ws.PathParameter("feed-id", "identifier of the feed").DataType("string")).
		Writes(types.Feed{}))
	ws.Route(ws.POST("/refresh/{feed-id}").To(func(req *restful.Request, res *restful.Response) { s.refreshFeed(req, res) }).
		Doc("refresh a feed").
		Param(ws.PathParameter("feed-id", "identifier of the feed").DataType("string")).
		Writes(Status{}))
	restful.Add(ws)
}
