package cutter

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/rand"
	"testing"
	"time"
)

func Test_LengthBasedCutter(t *testing.T) {
	//随机长度，随机包数测试
	b := new(bytes.Buffer)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	times := 1000 //1000轮测试
	for ts := 0; ts < times; ts++ {
		//随机几个包
		pkCnt := r.Intn(300) + 1

		for i := 0; i < pkCnt; i++ {
			length := uint32(r.Intn(10240) + 1)
			lengthBuf := []byte{0x00, 0x00, 0x00, 0x00}
			binary.BigEndian.PutUint32(lengthBuf, length)
			b.Write(lengthBuf)

			payload := make([]byte, length)

			for i, _ := range payload {
				payload[i] = 0x00
			}
			b.Write(payload)
		}

		for {
			pack := make([]byte, 102400)
			n, err := LengthBasedCutter(b, pack)
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Fatal(err)
			}

			pack = pack[0:n]
			for _, v := range pack {
				if v != 0x00 {
					t.Fatal("payload error")
				}
			}
		}
	}
}
