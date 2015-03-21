package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	base := "/sys/bus/w1/devices"
	matches, err := filepath.Glob(fmt.Sprintf("%s/28*", base))
	handleErr(err)

	deviceFile := filepath.Join(matches[0], "w1_slave")

	_, err = os.Stat(deviceFile)
	handleErr(err)

	buf, err := ioutil.ReadFile(deviceFile)
	handleErr(err)

	s := string(buf)
	fmt.Println("Device File Info", s)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
