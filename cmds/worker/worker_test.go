package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/remogatto/prettytest"
	"github.com/ryanbillingsley/blackmagic"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Setup

type WorkerTestSuite struct {
	prettytest.Suite
}

func TestRunner(t *testing.T) {
	prettytest.RunWithFormatter(t, new(prettytest.TDDFormatter), new(WorkerTestSuite))
}

// EndSetup

func (s *WorkerTestSuite) TestShouldGetForecastFromApi() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dat, _ := ioutil.ReadFile("./_fixtures/test.json")
		io.WriteString(w, string(dat))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	worker := &Worker{}

	res, err := worker.apiForecast(server.URL)
	if err != nil {
		s.Error("Resonse returned an error", err)
	}

	s.Equal(res.Forecast.SimpleForecast.ForecastDays[0].High.Fahrenheit, "36")
}

type mockCollection struct{}

func (c *mockCollection) All(map[string]interface{}, *[]interface{}) error { return nil }
func (c *mockCollection) First(map[string]interface{}, *interface{}) error { return nil }
func (c *mockCollection) FindId(bson.ObjectId, *interface{}) error         { return nil }
func (c *mockCollection) UpdateId(day bson.ObjectId, update map[string]interface{}) error {
	return nil
}
func (c *mockCollection) UpsertId(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return nil, nil
}

type mockDatabase struct{}

func (db *mockDatabase) Connect() error { return nil }
func (db *mockDatabase) Collection(collectionName string) (blackmagic.Collection, error) {
	if collectionName == "forecasts" {
		return &mockCollection{}, nil
	}
	return nil, nil
}

func (s *WorkerTestSuite) TestShouldCreateAForecastForEachDay() {
	dat, _ := ioutil.ReadFile("./_fixtures/test.json")

	var apires blackmagic.WUndergroundResponse
	_ = json.Unmarshal(dat, &apires)

	worker := &Worker{
		Database: &mockDatabase{},
	}
	days, err := worker.parseForecast(apires)

	s.Equal(err, nil)
	s.Equal(len(days), 10)
}

func (s *WorkerTestSuite) TestShouldCreateAForecastFromData() {
	dat, _ := ioutil.ReadFile("./_fixtures/test.json")

	var apires blackmagic.WUndergroundResponse
	_ = json.Unmarshal(dat, &apires)

	worker := &Worker{
		Database: &mockDatabase{},
	}

	forecast, err := worker.createForecast(apires.Forecast.SimpleForecast.ForecastDays[2])
	s.Equal(err, nil)

	s.Equal(forecast.High, 29)
}
