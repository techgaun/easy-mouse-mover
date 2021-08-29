package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type MousePosition struct {
	X int
	Y int
}

type MonitorInfo struct {
	X      int
	Y      int
	StartX int
	StartY int
}

const (
	argError = "left or right must be passed for direction"
	left     = "left"
	right    = "right"
)

func getCurrentPosition() MousePosition {
	out, err := exec.Command("xdotool", "getmouselocation").Output()
	if err != nil {
		log.Fatal(err)
	}

	xy := strings.Split(string(out), " ")
	x, _ := strconv.Atoi(strings.Split(xy[0], ":")[1])
	y, _ := strconv.Atoi(strings.Split(xy[1], ":")[1])

	return MousePosition{
		X: x,
		Y: y,
	}
}

func getMonitorInfo() (monitors []MonitorInfo) {
	out, err := exec.Command("xrandr", "--current").Output()
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`.*\sconnected.*\s(\d+)x(\d+)\+(\d+)\+(\d+)`)
	matches := re.FindAllStringSubmatch(string(out), -1)

	for _, match := range matches {
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])
		startX, _ := strconv.Atoi(match[3])
		startY, _ := strconv.Atoi(match[4])
		monitors = append(monitors, MonitorInfo{
			X:      x,
			Y:      y,
			StartX: startX,
			StartY: startY,
		})
	}

	sort.Slice(monitors, func(i, j int) bool {
		return monitors[i].StartX < monitors[j].StartX
	})
	return
}

func getNewPosition(currentPosition MousePosition, monitors []MonitorInfo, direction string) (newPosition MousePosition) {
	var currentMonitorIndex, newMonitorIndex int

	for idx, monitor := range monitors {
		if currentPosition.X > monitor.StartX && currentPosition.X < (monitor.X+monitor.StartX) {
			currentMonitorIndex = idx
		}
	}

	if direction == left {
		if currentMonitorIndex == 0 {
			newMonitorIndex = len(monitors) - 1
		} else {
			newMonitorIndex = currentMonitorIndex - 1
		}
	} else {
		if currentMonitorIndex == (len(monitors) - 1) {
			newMonitorIndex = 0
		} else {
			newMonitorIndex = currentMonitorIndex + 1
		}
	}

	newMonitor := monitors[newMonitorIndex]
	newPosition = MousePosition{
		X: newMonitor.X/2 + newMonitor.StartX,
		Y: newMonitor.Y/2 + newMonitor.StartY/2,
	}
	return
}

func moveToNewPosition(newPosition MousePosition) {
	err := exec.Command("xdotool", "mousemove", strconv.Itoa(newPosition.X), strconv.Itoa(newPosition.Y)).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatal(argError)
	}

	if args[0] != left && args[0] != right {
		log.Fatal(argError)
	}
	currentPosition := getCurrentPosition()
	monitors := getMonitorInfo()
	newPosition := getNewPosition(currentPosition, monitors, args[0])
	moveToNewPosition(newPosition)
}
