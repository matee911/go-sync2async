package dvs

import (
	"strconv"
	"strings"
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
// header=prepare_root_header(tid, CMD_TYPE_OTHER)
// self.push(deviceio(header + no_command()))

func RootHeader(transactionId int, cmdType CmdType, sourceId int, destId, mopPpid) {
	command_type = Enum(cmdType, 2)
	transaction_number = Enum(transactionId, 9)
	source_id =
 //   source_id = enum(source_id, 4)
 //   assert isinstance(source_id, basestring) and len(source_id)==4
 //   dest_id = enum(dest_id, 4)
 //   assert isinstance(dest_id, basestring) and len(dest_id)==4

 //   return transaction_number + command_type + source_id + dest_id + mop_ppid + get_creation_date()

}