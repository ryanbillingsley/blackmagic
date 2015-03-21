package main

import (
	"log"
	"os"

	"github.com/gizak/termui"
	"github.com/nsf/termbox-go"
)

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("This is a test log entry")

	err = termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	//termui.UseTheme("helloworld")

	data := []int{40, 42, 44, 46, 48, 50, 50, 50, 48, 46, 42, 40, 35, 32, 30, 30}
	spl0 := termui.NewSparkline()
	spl0.Data = data[3:]
	spl0.Height = 3
	spl0.LineColor = termui.ColorGreen

	// single
	spls0 := termui.NewSparklines(spl0)
	spls0.Height = 5
	spls0.Width = 20
	spls0.HasBorder = false

	termui.Render(spls0)
	termbox.PollEvent()
}
