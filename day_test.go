package blackmagic

import (
	"testing"

	"github.com/remogatto/prettytest"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Setup

type DayTestSuite struct {
	prettytest.Suite
}

func TestRunner(t *testing.T) {
	prettytest.RunWithFormatter(t, new(prettytest.TDDFormatter), new(DayTestSuite))
}

// EndSetup

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
func (db *mockDatabase) Collection(collectionName string) (Collection, error) {
	if collectionName == "days" {
		return &mockCollection{}, nil
	}
	return nil, nil
}

func (s *DayTestSuite) TestShouldAddNewForecast() {
	day := &Day{}
	forecast := &Forecast{Id: bson.NewObjectId()}

	err := day.AddForecast(forecast, &mockDatabase{})
	s.Equal(err, nil)
}
