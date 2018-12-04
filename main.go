package main

import (
	"sync"
	"time"

	"./events"
	"fmt"
)

var round int

func EventListener() {
	for {
		event := events.GetEvent()
		fmt.Println(event)
		if event == true {
			events.ShowGroup()
			time.Sleep(time.Second * 4)
		} else if event == false {
			time.Sleep(time.Second * 4)
			continue
		}
	}
}

//func EventShow() {
//	fmt.Println("Show round " + strconv.Itoa(round))
//	round++
//	events.ShowGroup()
//}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go EventListener()
	wg.Wait()
}
