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
func creationDate() string {
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

func deviceIO(body string) []byte {
	return append(Hexlen(body, 2), []byte(body)...)
}

// This command allows the portal to ping the DVS.
// 6. Operation command address header
// 6.1. 
// Description
// This section defines the format of the address headers used by the ACK and NACK in the VOD
// protocol. Its format is defined by the following table:
func noCommand() string {
	return "1002"
}

func NoCommand(transactionId int, sourceId int, destId int, mopPpid int) []byte {
	header := RootHeader(transactionId, CmdTypeOther, sourceId, destId, mopPpid)
	return deviceIO(fmt.Sprint(header, noCommand()))
}

func RootHeader(transactionId int, cmdType CmdType, sourceId int, destId int, mopPpid int) string {
	return fmt.Sprint(Enum(transactionId, 9), Enum(int(cmdType), 2), Enum(sourceId, 4), Enum(destId, 4), strconv.Itoa(mopPpid), creationDate())
}
