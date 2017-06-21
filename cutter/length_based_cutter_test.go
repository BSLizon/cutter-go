package cutter

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

func Test_LengthBasedCutter(t *testing.T) {
	//TODO::full test example
	in := bytes.NewBuffer([]byte{})
	lengthBuf := []byte{0x00, 0x00, 0x00, 0x00}
	binary.BigEndian.PutUint32(lengthBuf, uint32(5))
	in.Write(lengthBuf)
	payload := []byte{0x00, 0x01, 0x10, 0x11, 0x19}
	in.Write(payload)

	out := []byte{0x00, 0x00, 0x00, 0x00, 0x00}
	err := LengthBasedCutter(in, uint32(1<<10), out)
	if err != nil {
		if err == io.EOF {
			goto Here1
		}
		t.Error("test failed.", err)
	}

Here1:
	if payload == nil || out == nil {
		t.Error("test failed.")
	}

	if len(payload) != len(out) {
		t.Error("test failed.")
	}

	for i := range payload {
		if payload[i] != out[i] {
			t.Error("test failed.")
		}
	}

	//TODO::full test example

	err = LengthBasedCutter(in, uint32(1<<10), out)
	if err != nil {
		if err == io.EOF {
			goto Here2
		}
		t.Error("test failed.", err)
	}

Here2:
	if payload == nil || out == nil {
		t.Error("test failed.")
	}

	if len(payload) != len(out) {
		t.Error("test failed.")
	}

	for i := range payload {
		if payload[i] != out[i] {
			t.Error("test failed.")
		}
	}
}
