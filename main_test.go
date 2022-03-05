package main

import (
	"log"
	"testing"
	"time"
)

func run(c <-chan int) {
	for {
		i := <-c
		log.Println(i)
		time.Sleep(4 * time.Second)
	}
}

func Test(t *testing.T) {
	c := make(chan int, 10)
	c <- 1
	go run(c)
	for i := 0; i < 10; i++ {
		c <- i
		time.Sleep(2 * time.Second)
	}
	time.Sleep(20 * time.Second)
}
