package cutter

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
	"time"
)

func Test_LengthBasedCutter(t *testing.T) {
	//TODO::full test example
	in := bytes.NewBuffer([]byte{})
	lengthBuf := []byte{0x00, 0x00, 0x00, 0x00}
	binary.BigEndian.PutUint32(lengthBuf, uint32(5))
	in.Write(lengthBuf)
	payload := []byte{0x00, 0x01, 0x10, 0x11, 0x19}
	in.Write(payload)
	in.Write(lengthBuf)
	in.Write(payload)
	in.Write(lengthBuf)

	out := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	t1 := time.Now().UnixNano()
	n, err := LengthBasedCutter(in, out)
	t.Log(time.Now().UnixNano(), t1)
	if err != nil {
		if err == io.EOF {
			goto Here1
		}
		t.Fatal("test failed.", err)
	}

Here1:
	out = out[0:n]
	if payload == nil || out == nil {
		t.Fatal("test failed.")
	}

	t.Log(payload, out)

	if len(payload) != len(out) {
		t.Fatal("test failed.")
	}

	for i := range payload {
		if payload[i] != out[i] {
			t.Fatal("test failed.")
		}
	}

	//TODO::full test example

	_, err = LengthBasedCutter(in, out)
	if err != nil {
		if err == io.EOF {
			goto Here2
		}
		t.Fatal("test failed.", err)
	}

Here2:
	if payload == nil || out == nil {
		t.Fatal("test failed.")
	}

	if len(payload) != len(out) {
		t.Fatal("test failed.")
	}

	for i := range payload {
		if payload[i] != out[i] {
			t.Fatal("test failed.")
		}
	}
}
