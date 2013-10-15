package dvs

import (
	"fmt"
)

// This command allows the portal to ping the DVS.
// 6. Operation command address header
// 6.1. 
// Description
// This section defines the format of the address headers used by the ACK and NACK in the VOD
// protocol. Its format is defined by the following table:
func NoCommand(transactionId int, sourceId int, destId int, mopPpid int) []byte {
	header := RootHeader(transactionId, CmdTypeOther, sourceId, destId, mopPpid)
	return DeviceIO(fmt.Sprint(header, NoCommandBody()))
}
