package nzrpc

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
)

const magic_id = 810711662 //nzR0

type NzRPC struct {
	Addr string
	Port uint16
	User string

	loginChallenge []byte
	params         []byte
	isLogin        bool

	conn net.Conn
}

func (c *NzRPC) doOpen(addr string, port uint16) error {
	var err error
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}
	return nil
}

func (c *NzRPC) doChallenge(challenge []byte) error {
	msg := NewMessage(CmdChallenge)
	buf, err := msg.Encode(nil)
	if err != nil {
		return err
	}

	_, e := c.conn.Write(buf)
	if e != nil {
		return e
	}

	_, e = c.conn.Read(buf[:msg.Size()])
	if e != nil {
		return e
	}

	err = msg.Decode(buf)
	if err != nil {
		return e
	}

	if msg.Length > 0 {
		data := make([]byte, msg.Length)
		_, e = io.ReadFull(c.conn, data)
		if e != nil {
			return e
		}
		r := regexp.MustCompile("rpc.challenge=([a-z,0-9]+)[\n]?")
		list := r.FindStringSubmatch(string(data))
		if len(list) < 1 {
			return errors.New("No challenge")
		}
		c.loginChallenge = []byte(list[1])

	}

	return nil
}

func (c *NzRPC) doCheckPwd(challenge []byte, user, pwd string) error {
	h := md5.New()
	h.Write(c.loginChallenge)
	h.Write([]byte(pwd))
	s := fmt.Sprintf("rpc.username=%s\nrpc.password=%s\n", user, hex.EncodeToString(h.Sum(nil)))

	msg := NewMessage(CmdCheckPw)
	buf, err := msg.Encode([]byte(s))
	if err != nil {
		return err
	}

	_, e := c.conn.Write(buf)
	if e != nil {
		return e
	}

	_, e = c.conn.Read(buf[:msg.Size()])
	if e != nil {
		return e
	}

	err = msg.Decode(buf)
	if err != nil {
		return e
	}

	if msg.Length > 0 {
		data := make([]byte, msg.Length)
		_, e = io.ReadFull(c.conn, data)
		if e != nil {
			return e
		}
		c.params = data

	}

	return nil
}

func (c *NzRPC) doClose() error {
	c.conn.Close()
	c.conn = nil
	c.isLogin = false
	return nil
}

func (c *NzRPC) Login(addr string, port uint16, user, password string) error {
	c.Addr = addr
	c.Port = port

	err := c.doOpen(addr, port)
	if err != nil {
		return err
	}

	err = c.doChallenge(c.loginChallenge)
	if err != nil {
		return err
	}

	err = c.doCheckPwd(c.loginChallenge, user, password)
	if err != nil {
		return err
	}

	c.User = user
	return nil
}

func (c *NzRPC) Logout() error {
	if c.isLogin {
		c.doClose()
	}
	return nil
}

func (c *NzRPC) JsonCall(cmd string, param interface{}) ([]byte, error) {
	jsonInput := map[string]interface{}{
		"nzrpc": "1.0",
		"call":  cmd,
		"param": param,
	}

	data := bytes.NewBuffer(nil)
	err := json.NewEncoder(data).Encode(jsonInput)
	if err != nil {
		return nil, err
	}

	msg := NewMessage(CmdJsonCall)
	buf, err := msg.Encode(data.Bytes())
	if err != nil {
		return nil, err
	}

	_, e := c.conn.Write(buf)
	if e != nil {
		return nil, e
	}

	_, e = c.conn.Read(buf[:msg.Size()])
	if e != nil {
		return nil, e
	}

	err = msg.Decode(buf)
	if err != nil {
		return nil, e
	}

	if msg.Length > 0 {
		data := make([]byte, msg.Length)
		_, e = io.ReadFull(c.conn, data)
		if e != nil {
			return nil, e
		}
		var jsonOut struct {
			NZRPC  string          `json:"nzrpc"`
			Param  json.RawMessage `json:"param"`
			Result []interface{}   `json:"result"`
		}
		e = json.NewDecoder(bytes.NewReader(data)).Decode(&jsonOut)
		if e != nil {
			return nil, e
		}
		if jsonOut.Result[1].(float64) != 0 {
			return nil, NewError(
				int(jsonOut.Result[1].(float64)),
				int(jsonOut.Result[2].(float64)),
			)

		}
		return jsonOut.Param, nil

	}

	return nil, nil
}
