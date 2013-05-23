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

	restful.Add(NewFeedService(s, r))

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

func NewFeedService(s types.Store, r types.Rss) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/feed").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Route(ws.GET("/{feed-id}").To(getFeed).
		Doc("get a feed").
		Param(ws.PathParameter("feed-id", "identifier of the feed").DataType("string")).
		Writes(types.Feed{}))
	return ws
}

func getFeed(request *restful.Request, response *restful.Response) {
	response.WriteEntity(types.Feed{Id: "Id"})
}
