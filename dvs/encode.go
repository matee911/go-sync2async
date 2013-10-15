package dvs

import (
	"strconv"
	"strings"
	"time"
	"fmt"
	"encoding/binary"
)

type CmdType int

const (
	CmdTypeOther CmdType = 5
	CmdTypeVod CmdType = 8
)


func Enum(i int, size int) string {
	s := strconv.Itoa(i)
	if len(s) >= size {
		return s
	}
	return strings.Repeat("0", size-len(s)) + s
}

// now YYYYMMDD formated in UTC
func CreationDate() string {
	t := time.Now().UTC()
	return fmt.Sprintf("%02d%02d%02d", t.Year(), t.Month(), t.Day())
}

func Ehex(i int, size int) []byte {
	b := [4]byte{}
	binary.BigEndian.PutUint32(b[:], uint32(i))
	return b[4-size:]
}

func Hexlen(s string, size int) []byte {
	return Ehex(len(s), size)
}

func DeviceIO(body string) []byte {
	return append(Hexlen(body, 2), []byte(body)...)
}

func NoCommandBody() string {
	return "1002"
}

func RootHeader(transactionId int, cmdType CmdType, sourceId int, destId int, mopPpid int) string {
	return fmt.Sprint(Enum(transactionId, 9), Enum(int(cmdType), 2), Enum(sourceId, 4), Enum(destId, 4), strconv.Itoa(mopPpid), CreationDate())
}
