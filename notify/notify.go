package notify

import (
	"encoding/json"
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
		log.Printf("给%s发送邮件通知失败，请检查", email)
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
		log.Printf("给%s发送push通知失败，请检查", token)
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
	if db.Email != "" {
		SendMail(db.Email, title, text)
	}
	if db.Push != "" {
		SendPush(db.Push, title, text)
	}
}
