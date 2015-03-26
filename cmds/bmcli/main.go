package main

import (
	"flag"
	"log"
	//"os"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/ryanbillingsley/blackmagic"
	"gopkg.in/mgo.v2/bson"
)

import ui "github.com/ryanbillingsley/termui"

type paragraph struct {
	text   string
	color  ui.Attribute
	offset int
}

func main() {
	mongoUrl := flag.String("mongo", "localhost", "The mongo db address.  It can be as simple as `localhost` or involved as `mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb`")
	databaseName := flag.String("db", "blackmagic", "The name of the database you are connecting to.  Defaults to blackmagic")
	flag.Parse()

	//f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//panic(err)
	//}
	//defer f.Close()

	//log.SetOutput(f)

	db := blackmagic.NewDatabase(*mongoUrl, *databaseName)
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ui.UseTheme("helloworld")

	lc := ui.NewLineChart()
	lc.Border.Label = "24 Hour Temperature"
	lc.Mode = "dot"
	lc.Data = []float64{}
	lc.Width = 60
	lc.Height = 16
	lc.X = 0
	lc.Y = 0
	lc.AxesColor = ui.ColorWhite
	lc.LineColor = ui.ColorCyan | ui.AttrBold

	//parC := termui.NewPar("")
	//parC.Height = 9
	//parC.Width = 37
	//parC.Y = 0
	//parC.X = 77
	//parC.Border.Label = "Temperatures"
	//parC.Border.FgColor = termui.ColorYellow

	//temps := []paragraph{
	//{text: "High: 65.54", color: termui.ColorWhite, offset: 1},
	//{text: "Predicted High: 68 <---> 66", color: termui.ColorRed, offset: 2},
	//{text: "Low: 65.54", color: termui.ColorWhite, offset: 4},
	//{text: "Predicted Low: 68 <---> 66", color: termui.ColorBlue, offset: 5},
	//}

	//widgets := []termui.Bufferer{}

	//widgets = append(widgets, lc)
	//widgets = append(widgets, parC)

	//for _, t := range temps {
	//widgets = append(widgets, createPar(t.text, t.offset, t.color))
	//}

	draw := func(data []float64) {
		lc.Data = data
		ui.Render(lc)
	}

	evt := make(chan termbox.Event)
	go func() {
		for {
			evt <- termbox.PollEvent()
		}
	}()

	poll := make(chan []blackmagic.Reading)
	go func() {
		for {
			cd, err := db.Collection("days")
			cr, err := db.Collection("readings")
			if err != nil {
				log.Fatal(err)
			}

			var d blackmagic.Day
			t := time.Now().Local()
			year, month, day := t.Date()
			err = cd.Find(bson.M{"year": year, "month": month, "day": day}).One(&d)

			var readings []blackmagic.Reading
			err = cr.Find(bson.M{"day": d.Id}).Sort("createdat").All(&readings)
			if err != nil {
				log.Fatal(err)
			} else {
				poll <- readings
			}

			time.Sleep(time.Minute * 5)
		}
	}()

	for {
		select {
		case e := <-evt:
			if e.Type == termbox.EventKey && e.Ch == 'q' {
				return
			}
		case r := <-poll:
			data := make([]float64, 0)
			upper := 50
			if len(r) < 50 {
				upper = len(r)
			}

			for _, rd := range r[:upper] {
				data = append(data, rd.Temperature)
			}

			draw(data)
		}
	}
}

func createPar(text string, offset int, color ui.Attribute) *ui.Par {
	par := ui.NewPar(text)
	par.Height = 1
	par.Width = 30
	par.X = 79
	par.Y = 1 + offset
	par.HasBorder = false
	par.TextFgColor = color

	return par
}
