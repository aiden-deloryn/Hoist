package util

import (
	"fmt"
	"net"
)

func GetLocalIPAddress() (string, error) {
	var localAddress string
	conn, err := net.Dial("udp", "1.1.1.1:80")

	if err != nil {
		return localAddress, err
	}

	localAddress = conn.LocalAddr().(*net.UDPAddr).IP.String()

	return localAddress, nil
}

func GenerateProgressBarString(progress int) string {
	progressBarSize := 20
	progressBarValue := progress / (100 / progressBarSize)
	progressBarString := ""

	// Populate the current progress with '='
	for i := 0; i < progressBarValue; i++ {
		progressBarString += "="
	}

	// Pad the remaining progress with ' '
	for i := 0; i < progressBarSize-progressBarValue; i++ {
		progressBarString += " "
	}

	// Convert progress to string
	progressString := fmt.Sprint(progress)
	progressString += "%"

	progressStringSplit := []rune(progressString)
	progressBarSplit := []rune(progressBarString)

	// Add percentage text in the middle of the progress bar
	for i := (len(progressBarSplit) - len(progressString)) / 2; i < (len(progressBarSplit)+len(progressString))/2; i++ {
		j := i - (len(progressBarSplit)-len(progressString))/2
		progressBarSplit[i] = progressStringSplit[j]
	}

	progressBarString = "|" + string(progressBarSplit) + "|"

	return progressBarString
}
