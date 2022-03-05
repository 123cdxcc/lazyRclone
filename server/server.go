package server

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"net"
	"qbCopyProject/entity"
	"qbCopyProject/notify"
	"qbCopyProject/rclone"
	"strconv"
	"strings"
)

var queue = make(chan entity.RcloneCopyConfig)

func parseArgs() (string, string, string) {
	addr := flag.String("addr", ":1564", "服务端启动地址和端口")
	email := flag.String("email", "", "任务完成通知邮箱")
	token := flag.String("token", "", "任务完成通知pushPlus token")
	flag.Parse()
	return *addr, *email, *token
}

func uploading(email, token string) {
	for {
		config := <-queue
		log.Printf("上传队列收到任务，任务执行中->%s", config.UploadingFilePath)
		err := rclone.Uploading(
			config.RemoteName,
			strconv.Itoa(config.ThreadCount),
			config.LogFilePath, config.UploadingFilePath,
			config.RemoteFilePath)
		if err != nil {
			log.Printf("任务执行错误->%v\n", err)
			if email != "" {
				notify.SendMail(email, "lazyQB任务失败通知", "任务"+config.UploadingFilePath+"传输失败")
			}
			if token != "" {
				notify.SendPush(token, "lazyQB任务失败通知", "任务"+config.UploadingFilePath+"传输失败")
			}
			continue
		}
		log.Printf("执行完毕，任务文件->%s", config.UploadingFilePath)
		if email != "" {
			notify.SendMail(email, "lazyQB任务完成通知", "任务"+config.UploadingFilePath+"传输完成")
		}
		if token != "" {
			notify.SendPush(token, "lazyQB任务完成通知", "任务"+config.UploadingFilePath+"传输完成")
		}
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
	//log.Printf("%v\n", sb.String())
	var rcc entity.RcloneCopyConfig
	err := json.Unmarshal([]byte(sb.String()), &rcc)
	if err != nil {
		log.Printf("反序列化失败,原因->%v\n", err)
		return
	}
	queue <- rcc
	log.Println("任务提交成功")
}

func server(addr string) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("服务端启动完成")
	log.Println("开始等待任务")
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

func StartServer() {
	addr, email, token := parseArgs()
	log.Println("启动任务队列")
	go uploading(email, token)
	log.Println("启动服务端")
	server(addr)
}
