package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ryanbillingsley/blackmagic"
	"gopkg.in/mgo.v2/bson"
)

type Worker struct {
	Database blackmagic.Database
}

func main() {
	mongoUrl := flag.String("mongo", "localhost", "The mongo db address.  It can be as simple as `localhost` or involved as `mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb`")
	databaseName := flag.String("db", "blackmagic", "The name of the database you are connecting to.  Defaults to blackmagic")
	apiKey := flag.String("api", "", "Your WUnderground API Key")
	flag.Parse()

	database := blackmagic.NewDatabase(*mongoUrl, *databaseName)
	err := database.Connect()

	worker := &Worker{Database: database}

	res, err := worker.apiForecast(fmt.Sprintf("http://api.wunderground.com/api/%s/forecast10day/q/IN/Indianapolis.json", *apiKey))

	_, err = worker.parseForecast(res)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func (worker *Worker) apiForecast(apiUrl string) (blackmagic.WUndergroundResponse, error) {
	var apires blackmagic.WUndergroundResponse

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal(err)
		return apires, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return apires, err
	}

	if err := json.Unmarshal(body, &apires); err != nil {
		return apires, err
	}

	return apires, nil
}

func (worker *Worker) parseForecast(res blackmagic.WUndergroundResponse) ([]blackmagic.Day, error) {
	var days []blackmagic.Day

	for _, fday := range res.Forecast.SimpleForecast.ForecastDays {
		form := "Mon Jan 2 15:04 MST 2006"
		t, err := time.Parse(form, fday.Date.StandardFormat())
		if err != nil {
			fmt.Println(err)
			return days, err
		}

		forecast, err := worker.createForecast(fday)
		if err != nil {
			log.Fatal(err)
			return days, err
		}

		cd, err := worker.Database.Collection("days")
		var d blackmagic.Day

		year, month, day := t.Date()
		query := cd.Find(bson.M{"year": year, "month": month, "day": day})
		cnt, err := query.Count()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		if cnt >= 1 {
			log.Println("Found day...")
			query.One(&d)
		} else {
			log.Println("No day, creating a new one...")
			d = blackmagic.Day{
				Id:        bson.NewObjectId(),
				Forecasts: []bson.ObjectId{},
			}
			d.SetDate(t)
		}

		d.Forecasts = append(d.Forecasts, forecast.Id)

		_, err = cd.UpsertId(d.Id, d)

		days = append(days, d)
	}

	return days, nil
}

func (worker *Worker) createForecast(fday blackmagic.ForecastDay) (*blackmagic.Forecast, error) {
	high, err := strconv.Atoi(fday.High.Fahrenheit)
	low, err := strconv.Atoi(fday.Low.Fahrenheit)

	forecast := &blackmagic.Forecast{
		Id:         bson.NewObjectId(),
		High:       high,
		Low:        low,
		Conditions: fday.Conditions,
		QPF:        fday.QpfAllDay.Inches,
		Date:       time.Now(),
	}

	cf, err := worker.Database.Collection("forecasts")
	_, err = cf.UpsertId(forecast.Id, forecast)

	return forecast, err
}
