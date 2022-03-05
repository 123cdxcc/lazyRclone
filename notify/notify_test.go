package notify

import "testing"

func TestNotify(t *testing.T) {
	//SendMail("2409880020@qq.com", "test", "text")
	SendPush("c4d51f57cc184ce886bdaa803c00d1be", "test", "c1ont1ext")
}
