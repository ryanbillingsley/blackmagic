package blackmagic

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Forecast struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Date       time.Time
	High       int
	Low        int
	QPF        float64
	Conditions string
}
