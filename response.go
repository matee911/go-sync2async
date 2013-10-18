package main

import "encoding/json"

type ErrorResponse struct {
	Resp ErrRespJSON `json:"resp"`
}

func (r ErrorResponse) String() (s string) {
	body, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(body)
	return
}

type ErrRespJSON struct {
	Status  string `json:"status"`
	Ts      int    `json:"ts"`
	ErrCode int    `json:"errcode"`
	ErrDesc string `json:"errdesc"`
	ErrText string `json:"err_text"`
}

type SuccessResponse struct {
	Resp SuccessResponseJSON `json:"resp"`
}

func (r SuccessResponse) String() (s string) {
	body, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(body)
	return
}

type SuccessResponseJSON struct {
	Status  string      `json:"status"`
	Ts      int         `json:"ts"`
	License LicenseJSON `json:"license"`
}

type LicenseJSON struct {
	Object           string `json:"object"`
	ValidToTimestamp int    `json:"valid_to_timestamp"`
	MetaData         string `json:"metadata"`
}


