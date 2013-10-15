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

func PushVodCommandBody() string {
	//_push_vod_query(self, transaction_number, address, vod_ent_id, content_id, expiration_dt, viewing_duration, metadata, chipset_type_string):
	//vod_address_part = prepare_vod_addr_header(address)
    //assert len(vod_address_part) == VOD_ADDR_HEADER_LEN
	//assert len(chipset_type_string) == 22
	//vod_cmd_part = prepare_vod_load_entitlement_body(vod_ent_id, content_id, expiration_dt, viewing_duration, metadata, chipset_type_string) 
	//self.push(deviceio(root_header_part + vod_address_part + vod_cmd_part))
	return ""
}

func RootHeader(transactionId int, cmdType CmdType, sourceId int, destId int, mopPpid int) string {
	return fmt.Sprint(Enum(transactionId, 9), Enum(int(cmdType), 2), Enum(sourceId, 4), Enum(destId, 4), strconv.Itoa(mopPpid), CreationDate())
}
