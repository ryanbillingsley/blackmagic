package blackmagic

import (
	"fmt"
	"sort"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Day struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Year      int
	Month     time.Month
	Day       int
	Forecasts []bson.ObjectId `json:",omitempty" bson:",omitempty"`
	Readings  []bson.ObjectId `json:",omitempty" bson:",omitempty"`
	High      float64
	Low       float64
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

func (day *Day) CurrentHighLow(db Database) (float64, float64, error) {
	c, err := db.Collection("readings")
	if err != nil {
		return 0, 0, err
	}

	var readings []Reading
	for _, rId := range day.Readings {
		var r Reading
		err := c.FindId(rId).One(&r)
		if err != nil {
			fmt.Println("Error finding reading", err, rId)
			return 0, 0, err
		}

		readings = append(readings, r)
	}

	fmt.Println("Readings", readings)

	sort.Sort(ByTemperature(readings))

	high := readings[len(readings)-1]
	low := readings[0]

	return high.Temperature, low.Temperature, nil
}
