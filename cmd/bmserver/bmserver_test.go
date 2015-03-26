package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remogatto/prettytest"
	"github.com/ryanbillingsley/blackmagic"
)

// Setup

type ServerTestSuite struct {
	prettytest.Suite
}

func TestRunner(t *testing.T) {
	prettytest.RunWithFormatter(t, new(prettytest.TDDFormatter), new(ServerTestSuite))
}

// EndSetup

type mockDatabase struct{}

func (db *mockDatabase) Connect() error { return nil }
func (db *mockDatabase) Collection(collectionName string) (blackmagic.Collection, error) {
	return nil, nil
}

func (s *ServerTestSuite) TestNewServerShouldCreateServer() {
	db := &mockDatabase{}
	server := NewServer(db)

	s.Not(s.Equal(server, nil))
	s.Equal(server.Database, db)
}

func (s *ServerTestSuite) TestHomeHandlerShouldRoute() {
	server := NewServer(nil)

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	server.homeHandler(res, req)

	s.Equal(res.Code, http.StatusOK)
}

func (s *ServerTestSuite) TestHomeHandlerShould404() {
	server := NewServer(nil)

	req, _ := http.NewRequest("GET", "/foo", nil)
	res := httptest.NewRecorder()

	server.homeHandler(res, req)

	s.Equal(res.Code, http.StatusNotFound)
}

func (s *ServerTestSuite) TestHomeHandlerShould405() {
	server := NewServer(nil)

	req, _ := http.NewRequest("POST", "/", nil)
	res := httptest.NewRecorder()

	server.homeHandler(res, req)

	s.Equal(res.Code, http.StatusMethodNotAllowed)
}
