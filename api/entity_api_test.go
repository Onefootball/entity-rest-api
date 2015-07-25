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
	recorded := erat.RunRequest(t, handler, erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/comment", server.URL), nil))
	recorded.CodeIs(200)
	recorded.BodyIs("[]")
}

func TestGETWithExistentSetShouldReturnJsonArray200(t *testing.T) {
	recorded := erat.RunRequest(t, handler, erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/user", server.URL), nil))
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
	recorded := erat.RunRequest(t, handler, erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/user/999", server.URL), nil))
	recorded.CodeIs(404)
}

func TestGETEntityThatExistsReturn200WithJson(t *testing.T) {

}

func TestPOSTWithInvalidEntityShouldReturn400(t *testing.T) {

}

func TestPOSTWithValidEntityShouldReturn201WithHeader(t *testing.T) {

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
