package api

import (
	"github.com/DATA-DOG/go-sqlmock"
	test "github.com/Onefootball/entity-rest-api/test"
	"testing"
)

func TestGETWithEmptySetShouldReturnEmptyJsonArray200(t *testing.T) {

}

func TestGETWithExistentSetShouldReturnJsonArray200() {

}

func TestGETEntityDoestExistsShouldReturn404() {

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