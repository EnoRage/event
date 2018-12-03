package main

import (
	"./events"
	"fmt"
	"strconv"
	"time"
	"os"
)

var round int

func EventListener() {
	for {

		for {

			 event,err := events.GetEvent()

			 if err != nil{
			 	fmt.Println(err)
			 	os.Exit(1)
			 }

			 if !event{
			 	break
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




func main() {
	go EventListener()
	time.Sleep(time.Second * 10)
	EventShow()
}

