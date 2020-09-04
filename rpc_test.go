package nzrpc

import (
	"bytes"
	"encoding/json"
	"net"
	"testing"
	"time"
)

func TestRPC(t *testing.T) {
	host := "10.10.168.19"
	port := uint16(10717)

	rpc := NzRPC{}
	err := rpc.Login(host, port, "", "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("login success")
	defer rpc.Logout()

	param := map[string]interface{}{
		"username": "a",
		"password": "a",
	}
	data, err := rpc.JsonCall("GetNxsUserInfo", param)
	if err != nil {
		t.Error(err)
		return
	}
	var files []struct {
		Path string `json:"path"`
		Res  string `json:"res"`
	}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&files)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("jsoncall[GetNxsUserInfo]:", files)

	//1分钟执行1次
	count := 0
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		t.Log("keepalive")
		_, err := rpc.JsonCall("Keepalive", nil)
		if _, ok := err.(*net.OpError); ok {
			if err := rpc.Login(host, port, "", ""); err != nil {
				t.Log(err)
				count++
				if count > 5 {
					break
				}
			}
		}
	}
}
