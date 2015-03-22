package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
	flag.Parse()

	log.Println("Trying to connect to", *mongoUrl)
	database := blackmagic.NewDatabase(*mongoUrl, *databaseName)
	err := database.Connect()
	handleErr(err)

	worker := &Worker{Database: database}

	dat, err := readTemp()
	handleErr(err)

	temp, err := parseTemp(dat)

	fmt.Printf("Temp in F: %.2f\n", temp)

	d, err := worker.findDay()

	r := blackmagic.Reading{
		Id:          bson.NewObjectId(),
		CreatedAt:   time.Now(),
		Temperature: temp,
		Day:         d.Id,
	}

	c, err := worker.Database.Collection("readings")
	handleErr(err)

	_, err = c.UpsertId(r.Id, r)
	handleErr(err)

	d.Readings = append(d.Readings, r.Id)

	high, low, err := d.CurrentHighLow(worker.Database)
	handleErr(err)

	d.High = high
	d.Low = low

	c.UpsertId(d.Id, d)
	fmt.Println("Saved day, done")
}

func (worker *Worker) findDay() (blackmagic.Day, error) {
	var d blackmagic.Day

	t := time.Now()
	year, month, day := t.Date()

	cd, err := worker.Database.Collection("days")
	if err != nil {
		return d, err
	}

	cd.Find(bson.M{"year": year, "month": month, "day": day}).One(&d)

	return d, nil
}

func readTemp() (string, error) {
	base := "/sys/bus/w1/devices"
	dirs, err := filepath.Glob(fmt.Sprintf("%s/28*", base))
	if err != nil {
		return "", err
	}

	deviceFile := filepath.Join(dirs[0], "w1_slave")

	if _, err := os.Stat(deviceFile); err != nil {
		return "", err
	}

	buf, err := ioutil.ReadFile(deviceFile)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func parseTemp(data string) (float64, error) {
	r, err := regexp.Compile(".*t=([0-9]{5,6})")
	tempStr := r.FindStringSubmatch(data)[1]
	temp, err := strconv.ParseFloat(tempStr, 64)

	if err != nil {
		return 0, err
	}

	tempC := temp / 1000.0
	tempF := tempC*1.8000 + 32.00

	return tempF, nil
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
