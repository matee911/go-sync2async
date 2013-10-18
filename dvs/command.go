package dvs

type CmdType int

const (
	CmdTypeOther CmdType = 5
	CmdTypeVod CmdType = 8
)

type Command struct {
	Body string
	Type CmdType
}

func NoCommand() Command {
	return Command{Body: "1002", Type: CmdTypeOther}
}

func PushVodCommand() Command {
	//_push_vod_query(self, transaction_number, address, vod_ent_id, content_id, expiration_dt, viewing_duration, metadata, chipset_type_string):
	//vod_address_part = prepare_vod_addr_header(address)
		//assert len(vod_address_part) == VOD_ADDR_HEADER_LEN
	//assert len(chipset_type_string) == 22
	//vod_cmd_part = prepare_vod_load_entitlement_body(vod_ent_id, content_id, expiration_dt, viewing_duration, metadata, chipset_type_string) 
	//self.push(deviceio(root_header_part + vod_address_part + vod_cmd_part))

	return Command{Body: "", Type: CmdTypeVod}
}
