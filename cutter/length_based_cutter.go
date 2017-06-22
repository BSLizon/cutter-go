package cutter

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	PAYLOAD_LENGTH_SIZE = 4       /*len(uint32)*/
	MAX_PAYLOAD_SIZE    = 1 << 20 /*1MB*/
)

func LengthBasedCutter(r io.Reader, payload []byte) (uint32, error) {
	maxlength := len(payload)
	if maxlength > MAX_PAYLOAD_SIZE {
		return 0, errors.New("out of MAX_PAYLOAD_SIZE limit")
	}

	//read length
	lenBuf := make([]byte, PAYLOAD_LENGTH_SIZE)
	var lengthBufIdx int
	for {
		n, err := r.Read(lenBuf[lengthBufIdx:])
		if err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}
			return 0, err
		}
		lengthBufIdx += n
		if lengthBufIdx == PAYLOAD_LENGTH_SIZE {
			break
		} else if lengthBufIdx > PAYLOAD_LENGTH_SIZE || lengthBufIdx < 0 {
			return 0, errors.New("read payload length error")
		}
	}

	length := binary.BigEndian.Uint32(lenBuf)
	if length > uint32(maxlength) {
		return 0, errors.New("payload length out of limit")
	} else if length == 0 {
		return 0, errors.New("payload length equals 0")
	}

	//read payload
	var payloadBufIdx uint32
	for {
		n, err := r.Read(payload[payloadBufIdx:])
		if err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}
			return 0, err
		}

		//TODO::check n == 0

		payloadBufIdx += uint32(n)

		if payloadBufIdx == length {
			return length, nil
		} else if payloadBufIdx > length || payloadBufIdx < 0 {
			return 0, errors.New("read payload error")
		}
	}
}
