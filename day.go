package blackmagic

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Day struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Year      int
	Month     time.Month
	Day       int
	Forecasts []bson.ObjectId `json:",omitempty" bson:",omitempty"`
}

func (d *Day) SetDate(t time.Time) {
	year, month, day := t.Date()
	d.Day = day
	d.Month = month
	d.Year = year
}

func (day *Day) AddForecast(f *Forecast, db Database) error {
	c, err := db.Collection("days")
	err = c.UpdateId(day.Id, bson.M{"$push": bson.M{"forecasts": f.Id}})

	return err
}
