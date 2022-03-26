package server

import (
	"bufio"
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"qbCopyProject/db"
	"qbCopyProject/entity"
	"qbCopyProject/notify"
	"qbCopyProject/rclone"
	"strconv"
	"strings"
)

var taskQueue = make(chan entity.RcloneCopyConfig, 14)
var tasks = &db.UploadingTasks

func parseArgs() (string, string, string) {
	addr := flag.String("addr", ":1564", "服务端启动地址和端口")
	email := flag.String("email", "", "任务完成通知邮箱")
	token := flag.String("token", "", "任务完成通知pushPlus token")
	flag.Parse()
	return *addr, *email, *token
}

func getArrConfig(obj entity.RcloneCopyConfig) *list.Element {
	for i := (*tasks).Front(); i != nil; i = i.Next() {
		t := i.Value.(entity.RcloneCopyConfig)
		if obj.RemoteName == t.RemoteName &&
			obj.RemoteFilePath == t.RemoteFilePath &&
			obj.LogFilePath == t.LogFilePath &&
			obj.UploadingFilePath == t.UploadingFilePath &&
			obj.ThreadCount == t.ThreadCount {
			return i
		}
	}
	return nil
}

func uploading() {
	for {
		config := <-taskQueue
		log.Printf("当前任务->%s\r\n", config.UploadingFilePath)
		err := rclone.Uploading(
			config.RemoteName,
			strconv.Itoa(config.ThreadCount),
			config.LogFilePath, config.UploadingFilePath,
			config.RemoteFilePath)
		if err != nil {
			log.Printf("任务错误->%v\n", err)
			notify.SendNotify("lazyQB任务失败通知", "任务"+config.UploadingFilePath+"传输失败")
		} else {
			notify.SendNotify("lazyQB任务完成通知", "任务"+config.UploadingFilePath+"传输完成")
		}
		element := getArrConfig(config)
		if element != nil {
			(*tasks).Remove(element)
		}
		log.Printf("任务结束->%s\r\n", config.UploadingFilePath)
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("队列剩余任务数量:%d,任务列表[", (*tasks).Len()))
		for i := (*tasks).Front(); i != nil; i = i.Next() {
			sb.WriteString(fmt.Sprintf("%s,", i.Value.(entity.RcloneCopyConfig).UploadingFilePath))
		}
		sb.WriteString("]")
		log.Println(sb.String())
		if (*tasks).Len() > 0 {
			log.Printf("下一个任务:%s\r\n", (*tasks).Front().Value.(entity.RcloneCopyConfig).UploadingFilePath)
		} else {
			log.Println("当前没有任务")
		}
		fmt.Println()
	}
}

func processed(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	log.Println("收到任务")
	var sb strings.Builder
	reader := bufio.NewReader(conn)
	var bf [12]byte
	for {
		n, err := reader.Read(bf[:])
		if n <= 0 {
			break
		}
		if err != nil {
			log.Printf("获取任务错误->%v\n", err)
			break
		}
		sb.Write(bf[:n])
	}
	var rcc entity.RcloneCopyConfig
	err := json.Unmarshal([]byte(sb.String()), &rcc)
	if err != nil {
		log.Printf("反序列化失败,原因->%v\n", err)
		return
	}
	(*tasks).PushBack(rcc)
	taskQueue <- rcc
	log.Println("提交成功")
}

func server(addr string) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("服务端启动完成")
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go processed(conn)
	}
}

func StartServer() {
	addr, email, token := parseArgs()
	db.Email = email
	db.Push = token
	log.Println("启动任务队列")
	go uploading()
	log.Println("启动服务端")
	server(addr)
}
