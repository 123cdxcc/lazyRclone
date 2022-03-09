package show

import "fmt"

func Demo() {
	fmt.Println("当前剩余任务：[1,23,4,5]")
	fmt.Println("当前正在进行的任务：/data/1.mp4")
	fmt.Println("下一个任务：/data/2.mp4")
	fmt.Print("任务信息：\r\n当前进度49%\r\n")
}
