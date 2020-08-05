package nzrpc

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestRPC(t *testing.T) {
	rpc := NzRPC{}
	err := rpc.Login("10.10.168.19", 10717, "", "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("login success")

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
	rpc.Logout()
}
