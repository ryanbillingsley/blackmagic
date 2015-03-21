package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

	match, _ := regexp.MatchString("t=([0-9]{5-6})", s)
	fmt.Println(match)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
