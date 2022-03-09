package show

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	//fmt.Printf("结果:%d", 1)
	for i := 0; i < 100; i++ {
		//fmt.Println("当前剩余任务：[1,23,4,5]")
		//fmt.Println("当前正在进行的任务：/data/1.mp4")
		//fmt.Println("下一个任务：/data/2.mp4")
		//fmt.Print("任务信息：\r\n当前进度" + strconv.Itoa(i) + "\r\n")
		fmt.Printf("结果:%d\r\n", i)
		time.Sleep(1 * time.Second)
		cmd := exec.Command("cmd", "/c", "cls")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}

}
