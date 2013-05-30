package main

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/jwiklund/reader/reader"
	"github.com/jwiklund/reader/types"
	"log"
	"net/http"
)

func main() {
	r := reader.NewReader("data")
	defer r.Close()

	registerStaticFiles()
	registerFeedService(service{r})

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
	reader types.Reader
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
	feeds, err := s.reader.GetStore().GetAllFeedsInfo()
	if err != nil {
		response.WriteEntity(Status{"fail", err.Error()})
	} else {
		response.WriteEntity(Feeds{"ok", feeds})
	}
}

func (s *service) getAllUserItems(request *restful.Request, response *restful.Response) {
	items, err := s.reader.GetStore().GetFeedByUser(request.PathParameter(("user-id")), "")
	if err != nil {
		response.WriteEntity(Status{"fail", err.Error()})
	} else {
		response.WriteEntity(Items{"ok", items})
	}
}

func (s *service) getFeed(request *restful.Request, response *restful.Response) {
	feed, err := s.reader.GetStore().GetFeed(request.PathParameter("feed-id"))
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
	_, err = s.reader.GetStore().GetFeed(feed.Id)
	if err == nil {
		response.WriteEntity(Status{"fail", "Feed already exists"})
		return
	}
	err = s.reader.GetStore().PutFeed(&feed)
	if err != nil {
		response.WriteEntity(Status{"fail", "Could not store Feed: " + err.Error()})
		return
	}
	response.WriteEntity(Status{"ok", "Feed created"})
}

func (s *service) refreshFeed(request *restful.Request, response *restful.Response) {
	err := s.reader.GetRss().Fetch(request.PathParameter("feed-id"))
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
	ws.Route(ws.GET("/").To(s.getFeeds).
		Doc("get all feeds").
		Writes(Feeds{}))
	ws.Route(ws.GET("/user/{user-id}/all").To(s.getAllUserItems).
		Doc("get all unread items for user").
		Param(ws.PathParameter("user-id", "the user id").DataType("string")).
		Writes(Feeds{}))
	ws.Route(ws.GET("/{feed-id}").To(s.getFeed).
		Doc("get a feed").
		Param(ws.PathParameter("feed-id", "identifier of the feed").DataType("string")).
		Writes(Feed{}))
	ws.Route(ws.POST("/").To(s.createFeed).
		Doc("create a feed").
		Param(ws.BodyParameter("feed", "the feed").DataType("types.Feed")).
		Writes(Status{}))
	ws.Route(ws.POST("/refresh/{feed-id}").To(s.refreshFeed).
		Doc("refresh a feed").
		Param(ws.PathParameter("feed-id", "identifier of the feed").DataType("string")).
		Writes(Status{}))
	restful.Add(ws)
}

func registerStaticFiles() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/html/", http.StatusMovedPermanently)
	})
	dir := http.FileServer(http.Dir("resources"))
	http.Handle("/css/", dir)
	http.Handle("/js/", dir)
	http.Handle("/html/", dir)
}
