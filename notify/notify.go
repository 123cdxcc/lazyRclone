package notify

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"qbCopyProject/db"
	"qbCopyProject/entity"
	"strings"
)

const PUSH = "push"
const EMAIL = "email"

func SendMail(email, title, text string) bool {
	client := &http.Client{}
	url := "https://msg.rsss.xyz/send/mail"
	data := strings.NewReader("to=" + email + "&title=" + title + "&text=" + text)
	req, _ := http.NewRequest("POST", url, data)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("给%s发送邮件通知失败，请检查:%v", email, err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	var result entity.EmailResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("解析数据错误->" + err.Error())
		return false
	}
	if result.Status && result.Code == 2000 {
		return true
	}
	return false
}

func SendPush(token, title, text string) bool {
	url := "https://www.pushplus.plus/send"
	client := &http.Client{}
	data := strings.NewReader("token=" + token + "&title=" + title + "&content=" + text)
	req, _ := http.NewRequest("POST", url, data)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("给%s发送push通知失败，请检查:%v", token, err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	var result entity.PushResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("解析数据错误->" + err.Error())
		return false
	}
	if result.Code == 200 {
		log.Printf("消息流水号->%v\n", strings.Replace(result.Data, " ", "", -1))
		return true
	}
	return false
}

func SendNotify(title, text string) {
	var sendNotifyResult bool
	if db.Email != "" {
		sendNotifyResult = SendMail(db.Email, title, text)
	}
	if db.Push != "" {
		sendNotifyResult = SendPush(db.Push, title, text)
	}
	if sendNotifyResult {
		log.Printf("\"%v\"通知消息发送成功", text)
	} else {
		log.Printf("\"%v\"通知消息发送失败", text)
	}
}

func Message() {
	//当前任务显示
	indexTaskName := "空"
	if db.UploadingTasks.Front() != nil {
		indexTaskName = db.UploadingTasks.Front().Value.(entity.RcloneCopyConfig).UploadingFilePath
	}
	log.Printf("当前任务: %v\r\n", indexTaskName)
	//下个任务显示
	nextTaskName := "空"
	if db.UploadingTasks.Front() != nil && db.UploadingTasks.Front().Next() != nil {
		nextTaskName = db.UploadingTasks.Front().Next().Value.(entity.RcloneCopyConfig).UploadingFilePath
	}
	log.Printf("下个任务: %v\r\n", nextTaskName)
	//任务列表显示
	tasksStr := strings.Builder{}
	tasksStr.WriteString("任务列表: [")
	for i := (db.UploadingTasks).Front(); i != nil; i = i.Next() {
		tasksStr.WriteString(fmt.Sprintf("%s,", i.Value.(entity.RcloneCopyConfig).UploadingFilePath))
	}
	tasksStr.WriteString("]")
	log.Println(tasksStr.String())
	//错误列表显示，有错误才显示
	if db.ErrorTasks.Len() > 0 {
		errorStr := strings.Builder{}
		errorStr.WriteString("错误列表: [")
		for i := db.ErrorTasks.Front(); i != nil; i = i.Next() {
			errorStr.WriteString(fmt.Sprintf("%s,", i.Value.(string)))
		}
		errorStr.WriteString("]")
		log.Println(errorStr.String())
	}
	fmt.Println()
}
