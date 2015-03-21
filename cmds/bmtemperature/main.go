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
	dirs, err := filepath.Glob(fmt.Sprintf("%s/28*", base))
	handleErr(err)

	deviceFile := filepath.Join(dirs[0], "w1_slave")

	_, err = os.Stat(deviceFile)
	handleErr(err)

	buf, err := ioutil.ReadFile(deviceFile)
	handleErr(err)

	s := string(buf)

	r, err := regexp.Compile(".*t=([0-9]{5,6})")
	handleErr(err)

	matches := r.FindStringSubmatch(s)

	fmt.Println(matches[1])
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
