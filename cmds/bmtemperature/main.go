package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
	temp, err := parseTemp(s)

	fmt.Printf("Temp in F: %.2f", temp)
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
