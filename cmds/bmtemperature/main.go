package main

import (
	"bufio"
	"fmt"
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

	file, err := os.Open(deviceFile)
	handleErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println("Line: ", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
