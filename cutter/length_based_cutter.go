package cutter

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	PAYLOAD_LENGTH_SIZE  = 4       /*len(uint32)*/
	MAX_PAYLOAD_SIZE     = 1 << 20 /*1MB*/
	MAX_READ_RETRY_COUNT = 3
)

func LengthBasedCutter(r io.Reader, payload []byte) (uint32, error) {
	maxlength := len(payload)
	if maxlength > MAX_PAYLOAD_SIZE {
		return 0, errors.New("out of MAX_PAYLOAD_SIZE limit")
	}

	readRetryCount := 0

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

		if n == 0 {
			if readRetryCount >= MAX_READ_RETRY_COUNT-1 {
				return 0, errors.New("max read retry count")
			} else {
				readRetryCount++
			}
		} else {
			readRetryCount = 0

			lengthBufIdx += n
			if lengthBufIdx == PAYLOAD_LENGTH_SIZE {
				break
			} else if lengthBufIdx > PAYLOAD_LENGTH_SIZE || lengthBufIdx < 0 {
				return 0, errors.New("read payload length error")
			}
		}
	}

	length := binary.BigEndian.Uint32(lenBuf)
	if length > uint32(maxlength) {
		return 0, errors.New("receive length out of len([]byte)")
	} else if length == 0 {
		return 0, errors.New("receive length equals 0")
	}

	//read payload
	var payloadBufIdx uint32
	for {
		n, err := r.Read(payload[payloadBufIdx:length])
		if err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}
			return 0, err
		}

		if n == 0 {
			if readRetryCount >= MAX_READ_RETRY_COUNT-1 {
				return 0, errors.New("max read retry count")
			} else {
				readRetryCount++
			}
		} else {
			readRetryCount = 0

			payloadBufIdx += uint32(n)
			if payloadBufIdx == length {
				return length, nil
			} else if payloadBufIdx > length || payloadBufIdx < 0 {
				return 0, errors.New("read payload error")
			}
		}

	}
}
