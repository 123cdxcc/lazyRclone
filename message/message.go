package message

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"qbCopyProject/entity"
	"qbCopyProject/rclone"
	"strconv"
	"strings"
)

var queue = make(chan entity.RcloneCopyConfig)
var addr = ":1564"

func uploading() {
	for {
		config := <-queue
		log.Printf("上传队列收到任务，开始执行，任务文件->%s", config.UploadingFilePath)
		err := rclone.Uploading(
			config.RemoteName,
			strconv.Itoa(config.ThreadCount),
			config.LogFilePath, config.UploadingFilePath,
			config.RemoteFilePath)
		if err != nil {
			log.Printf("任务执行错误->%v\n", err)
		}
		log.Printf("执行完毕，任务文件->%s", config.UploadingFilePath)
	}
}

func processed(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	log.Println("开始提交任务")
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
	log.Printf("%v\n", sb.String())
	var rcc entity.RcloneCopyConfig
	err := json.Unmarshal([]byte(sb.String()), &rcc)
	if err != nil {
		log.Printf("反序列化失败,原因->%v\n", err)
		return
	}
	queue <- rcc
	log.Println("任务提交成功")
}

func Server() {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	go uploading()
	log.Println("开始监听")
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("收到连接")
		processed(conn)
	}
}

func Client(data []byte) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		log.Fatalln(err)
	}
}
