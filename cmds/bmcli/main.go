package main

import (
	"log"
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/ryanbillingsley/termui"
)

type paragraph struct {
	text   string
	color  termui.Attribute
	offset int
}

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(f)

	err = termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	termui.UseTheme("helloworld")

	lc := termui.NewLineChart()
	lc.Border.Label = "dot-mode Line Chart"
	lc.Mode = "dot"
	lc.Data = []float64{65.4, 65.8, 66, 66.2, 66.25, 68.75, 65.3, 62.4}
	lc.Width = 77
	lc.Height = 16
	lc.X = 0
	lc.Y = 0
	lc.AxesColor = termui.ColorWhite
	lc.LineColor = termui.ColorCyan | termui.AttrBold

	parC := termui.NewPar("")
	parC.Height = 9
	parC.Width = 37
	parC.Y = 0
	parC.X = 77
	parC.Border.Label = "Temperatures"
	parC.Border.FgColor = termui.ColorYellow

	temps := []paragraph{
		{text: "High: 65.54", color: termui.ColorWhite, offset: 1},
		{text: "Predicted High: 68 <---> 66", color: termui.ColorRed, offset: 2},
		{text: "Low: 65.54", color: termui.ColorWhite, offset: 4},
		{text: "Predicted Low: 68 <---> 66", color: termui.ColorBlue, offset: 5},
	}

	widgets := []termui.Bufferer{}

	widgets = append(widgets, lc)
	widgets = append(widgets, parC)

	for _, t := range temps {
		widgets = append(widgets, createPar(t.text, t.offset, t.color))
	}

	draw := func() {
		termui.Render(widgets...)
	}

	evt := make(chan termbox.Event)
	go func() {
		for {
			evt <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case e := <-evt:
			if e.Type == termbox.EventKey && e.Ch == 'q' {
				return
			}
		default:
			draw()
			time.Sleep(time.Second / 2)
		}
	}
}

func createPar(text string, offset int, color termui.Attribute) *termui.Par {
	par := termui.NewPar(text)
	par.Height = 1
	par.Width = 30
	par.X = 79
	par.Y = 1 + offset
	par.HasBorder = false
	par.TextFgColor = color

	return par
}
