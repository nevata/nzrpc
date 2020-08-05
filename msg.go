package nzrpc

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

//size of 48
type Message struct {
	Magic   uint32
	Command uint32
	Length  int32
	Flags   uint32
	Data    [32]byte
}

func NewMessage(cmd uint32) *Message {
	return &Message{
		Magic:   magic_id, //*(*uint32)(unsafe.Pointer(&([]byte("nzR0"))[0])),
		Command: cmd,
	}
}

func (m *Message) Size() uint32 {
	return uint32(unsafe.Sizeof(*m))
}

func (m *Message) Encode(data []byte) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, m.Magic)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, m.Command)
	if err != nil {
		return nil, err
	}

	if data != nil {
		m.Length = int32(len(data))
	}

	err = binary.Write(buffer, binary.LittleEndian, m.Length)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, m.Flags)
	if err != nil {
		return nil, err
	}

	_, err = buffer.Write(m.Data[:])
	if err != nil {
		return nil, err
	}

	if data != nil {
		_, err = buffer.Write(data)
		if err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func (m *Message) Decode(data []byte) error {
	bufReader := bytes.NewReader(data)

	err := binary.Read(bufReader, binary.LittleEndian, &m.Magic)
	if err != nil {
		return err
	}

	err = binary.Read(bufReader, binary.LittleEndian, &m.Command)
	if err != nil {
		return err
	}

	err = binary.Read(bufReader, binary.LittleEndian, &m.Length)
	if err != nil {
		return err
	}

	err = binary.Read(bufReader, binary.LittleEndian, &m.Flags)
	if err != nil {
		return err
	}

	return nil
}
