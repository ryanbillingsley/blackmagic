package blackmagic

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Reading struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt   time.Time     `json:"createdAt"`
	Day         bson.ObjectId
	Temperature float64
}

type ByTemperature []Reading

func (a ByTemperature) Len() int           { return len(a) }
func (a ByTemperature) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTemperature) Less(i, j int) bool { return a[i].Temperature < a[j].Temperature }
