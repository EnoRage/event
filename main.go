package main

import (
	"./events"
	"fmt"
	"strconv"
	"time"
)

var round int

func main() {
	go EventListener()
	time.Sleep(time.Second * 10)
	EventShow()
}

func EventListener() {
	for {
		for {
			if !events.GetEvent() {
				return
			}
		}
	}
}

func EventShow() {
	for {
		fmt.Println("Show round " + strconv.Itoa(round))
		round++
		events.ShowGroup()
		time.Sleep(time.Second * 15)
	}
}
