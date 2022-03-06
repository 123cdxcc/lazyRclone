package server

import (
	"container/list"
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	arr := list.New()
	/*arr.PushBack("1")
	arr.PushBack("2")
	arr.PushBack("3")*/
	i := arr.Front()
	if i != nil {
		arr.Remove(i)
	}
	fmt.Println(arr.Len())
	for item := arr.Front(); item != nil; item = item.Next() {
		fmt.Println(item.Value)
	}

}

var c = make(chan string)

func func1() {
	fmt.Println(1)
	go func2()
	c <- "2"
}

func func2() {
	fmt.Println(<-c)
}

func func3() {
	fmt.Println(3)
}

func Test2(t *testing.T) {
	go func1()
	time.Sleep(2 * time.Second)
}
