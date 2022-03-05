package message

import (
	"encoding/json"
	"log"
	"qbCopyProject/entity"
	"testing"
)

func TestServer(t *testing.T) {
	Server()
}

func TestClient(t *testing.T) {
	config := entity.RcloneCopyConfig{
		RemoteName:        "remoteName",
		ThreadCount:       3,
		LogFilePath:       "logName",
		UploadingFilePath: "uploadingFilePath",
		RemoteFilePath:    "remoteFilePath",
	}
	data, err := json.Marshal(&config)
	if err != nil {
		log.Printf("系列化错误->%v\n", err)
	}
	log.Printf("序列化结果->%v\n", string(data))
	Client(data)
}
