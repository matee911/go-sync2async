package dvs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func EncodeNum(i int, size int) string {
	s := strconv.Itoa(i)
	if len(s) >= size {
		return s
	}
	return strings.Repeat("0", size-len(s)) + s
}

func DecodeNum(s string) (int, error) {
	return strconv.Atoi(s)
}

// Hex

func EncodeHex(i int, size int) []byte {
	b := [4]byte{}
	binary.BigEndian.PutUint16(b[:], uint16(i))
	return b[2-size:]
}

func DecodeHex(b []byte) (uint16, error) {
	var i uint16
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, &i)
	return i, err
}

// now YYYYMMDD formated in UTC
func CreationDate() string {
	t := time.Now().UTC()
	return fmt.Sprintf("%02d%02d%02d", t.Year(), t.Month(), t.Day())
}

func Hexlen(s string, size int) []byte {
	return EncodeHex(len(s), size)
}

func DeviceIO(body string) []byte {
	return append(Hexlen(body, 2), []byte(body)...)
}

func RootHeader(transactionId int, cmdType CmdType, sourceId int, destId int, mopPpid int) string {
	return fmt.Sprint(EncodeNum(transactionId, 9), EncodeNum(int(cmdType), 2), EncodeNum(sourceId, 4), EncodeNum(destId, 4), strconv.Itoa(mopPpid), CreationDate())
}
