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

func LengthBasedCutter(r io.Reader, payloadMaxLength uint32, payload []byte) error {
	//TODO::check len(payload) or cap(payload)

	if payloadMaxLength > MAX_PAYLOAD_SIZE {
		return errors.New("out of payloadMaxLength limit")
	}

	//read length
	lenBuf := make([]byte, PAYLOAD_LENGTH_SIZE)
	var lengthBufIdx int
	for {
		n, err := r.Read(lenBuf[lengthBufIdx:])
		if err != nil {
			if err == io.EOF {
				return io.EOF
			}
			return err
		}
		lengthBufIdx += n
		if lengthBufIdx == PAYLOAD_LENGTH_SIZE {
			break
		} else if lengthBufIdx > PAYLOAD_LENGTH_SIZE || lengthBufIdx < 0 {
			return errors.New("read payload length error")
		}
	}

	length := binary.BigEndian.Uint32(lenBuf)
	if length > payloadMaxLength {
		return errors.New("payload length out of limit")
	} else if length == 0 {
		return errors.New("payload length equals 0")
	}

	//read payload
	var payloadBufIdx uint32
	for {
		n, err := r.Read(payload[payloadBufIdx:])
		if err != nil {
			if err == io.EOF {
				return io.EOF
			}
			return err
		}

		//TODO::check len(payload) or cap(payload)

		payloadBufIdx += uint32(n)

		if payloadBufIdx == length {
			return nil
		} else if payloadBufIdx > length || payloadBufIdx < 0 {
			return errors.New("read payload error")
		}
	}
}
