package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	base := "/sys/bus/w1/devices"
	matches, err := filepath.Glob(fmt.Sprintf("%s/28*", base))
	if err != nil {
		panic(err)
	}

	deviceFile := filepath.Join(matches[0], "w1_slave")

	deviceFileInfo, err := os.Stat(deviceFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("Device File Info", deviceFileInfo)
}
