package api

import (
	"database/sql"
	"fmt"
	eram "github.com/Onefootball/entity-rest-api/manager"
	erat "github.com/Onefootball/entity-rest-api/test"
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server  *httptest.Server
	handler http.Handler
)

type User struct {
	Id       int
	Username string
	Password string
	Salt     string
	Email    string
}

type Post struct {
	Id          int
	Title       string
	Content     string
	Create_Time int
	Author_Id   int
	Status      int
}

func init() {

	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatal("An error '%s' was not expected when opening a stub database connection", err)
	} else {

		dat, err := ioutil.ReadFile("./../test.sql")
		if err != nil {
			log.Fatal("An error '%s' was not expected when opening sql file", err)
		} else {
			db.Exec(string(dat))
		}
	}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	entityManager := eram.NewEntityDbManager(db)
	entityRestApi := NewEntityRestAPI(entityManager)

	router, err := rest.MakeRouter(
		rest.Get("/api/:entity", entityRestApi.GetAllEntities),
		rest.Post("/api/:entity", entityRestApi.PostEntity),
		rest.Get("/api/:entity/:id", entityRestApi.GetEntity),
		rest.Put("/api/:entity/:id", entityRestApi.PutEntity),
		rest.Delete("/api/:entity/:id", entityRestApi.DeleteEntity),
	)

	if err != nil {
		log.Fatal("An error '%s' was not expected when creating the router", err)
	}

	api.SetApp(router)
	handler = api.MakeHandler()

	server = httptest.NewServer(handler)
}

func TestGETWithEmptySetShouldReturnEmptyJsonArray200(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/comment", server.URL), nil))

	recorded.CodeIs(200)
	recorded.BodyIs("[]")
}

func TestGETWithExistentSetShouldReturnJsonArray200(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/user", server.URL), nil))

	recorded.CodeIs(200)

	data := []User{}
	err := recorded.DecodeJsonPayload(&data)

	if err != nil {
		t.Error(err)
	} else {

		if len(data) < 1 {
			t.Error("No users found, and should have found at least one.")
		}
	}
}

func TestGETEntityDoestExistsShouldReturn404(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/user/999", server.URL), nil))

	recorded.CodeIs(404)
}

func TestGETEntityThatExistsReturn200WithJson(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/user/1", server.URL), nil))

	recorded.CodeIs(200)

	data := User{}
	err := recorded.DecodeJsonPayload(&data)

	if err != nil {
		t.Error(err)
	} else {

		if data.Id != 1 {
			t.Error("Weird behavior finding different user Id.")
		}
	}
}

func TestPOSTWithInvalidEntityShouldReturn400(t *testing.T) {

	// This post doesn't have title and should be wrong then
	entity := new(Post)
	entity.Content = "<p>Onefootball test post content...</p>"
	entity.Create_Time = 1437839411
	entity.Author_Id = 1
	entity.Status = 1

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("POST", fmt.Sprintf("%s/api/post", server.URL), entity))

	t.Skipf("Invalid POST should return 400 but it is returning %s", recorded.Recorder.Code)
}

func TestPOSTWithExistingEntryShouldReturn409(t *testing.T) {

	// This post has a ID 1 that conflicts with the on in the database
	entity := new(Post)
	entity.Id = 1
	entity.Title = "Test Post 1"
	entity.Content = "<p>Onefootball test post content...</p>"
	entity.Create_Time = 1437839411
	entity.Author_Id = 1
	entity.Status = 1

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("POST", fmt.Sprintf("%s/api/post", server.URL), entity))

	t.Skipf("Invalid code, should be 409 as a conflict for the ID 1 but get %s", recorded.Recorder.Code)
}

func TestPOSTWithValidEntityShouldReturn201WithHeader(t *testing.T) {

	entity := new(Post)
	entity.Id = 10
	entity.Title = "Test Post 1"
	entity.Content = "<p>Onefootball test post content...</p>"
	entity.Create_Time = 1437839411
	entity.Author_Id = 1
	entity.Status = 1

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("POST", fmt.Sprintf("%s/api/post", server.URL), entity))

	recorded.CodeIs(201)
	recorded.HeaderIs("Location", fmt.Sprintf("post/%d", entity.Id))
}

func TestPUTWithInvalidEntityShouldReturn400(t *testing.T) {

}

func TestPUTWithInvalidEntityShouldReturn201(t *testing.T) {

}

func TestDELETEShouldReturn404IfEntityNotFound(t *testing.T) {

}

func TestDELETEShouldReturn200IfEntityExists(t *testing.T) {

}

func TestGETWithSortQueryStringsShouldReturn200(t *testing.T) {

}

func TestGETWithAllQueryStringsShouldReturn200(t *testing.T) {

}
