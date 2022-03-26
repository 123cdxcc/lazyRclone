package client

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net"
	"qbCopyProject/entity"
)

func client(data []byte, addr string) {
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

func parseArgs() (*entity.RcloneCopyConfig, *string, error) {
	remoteName := flag.String("n", "", "rclone config对应的名字")
	threadCount := flag.Int("t", 3, "线程数")
	logName := flag.String("l", "/tmp/uploading.log", "上传日志文件路径(绝对路径)")
	uploadingFilePath := flag.String("p", "", "待上传的文件或文件夹路径(绝对路径)")
	remoteFilePath := flag.String("r", "", "远程储存路径(绝对路径)")
	addr := flag.String("a", ":1564", "服务端地址和端口")
	flag.Parse()
	if *remoteName == "" {
		return nil, nil, errors.New("rclone名字不能为空")
	}
	if *uploadingFilePath == "" {
		return nil, nil, errors.New("待上传路径不能为空")
	}
	if *remoteFilePath == "" {
		return nil, nil, errors.New("远程目录不能为空")
	}
	config := entity.RcloneCopyConfig{
		RemoteName:        *remoteName,
		ThreadCount:       *threadCount,
		LogFilePath:       *logName,
		UploadingFilePath: *uploadingFilePath,
		RemoteFilePath:    *remoteFilePath,
	}
	return &config, addr, nil
}

func StartClient() {
	config, addr, err := parseArgs()
	data, err := json.Marshal(config)
	if err != nil {
		log.Fatalf("系列化错误->%v\n", err)
	}
	log.Printf("序列化结果->%v\n", string(data))
	client(data, *addr)
	log.Println("任务提交完毕")
}
