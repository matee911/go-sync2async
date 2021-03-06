package main

import (
	"github.com/matee911/go-sync2async/assert"
	"testing"
)

func TestJSONErrorResponse(t *testing.T) {
	assert.Equal(t, ErrorResponse{
		Resp: ErrRespJSON{
			Status:  "err",
			Ts:      1234567890,
			ErrCode: 2,
			ErrDesc: "err desc",
			ErrText: "err text"},
	}.String(),
		`{"resp":{"status":"err","ts":1234567890,"errcode":2,"errdesc":"err desc","err_text":"err text"}}`,
	)

}

func TestJSONSuccessResponse(t *testing.T) {
	assert.Equal(t, SuccessResponse{
		Resp: SuccessResponseJSON{
			Status: "ok",
			Ts:     123456790,
			License: LicenseJSON{
				Object:           "objectstr",
				ValidToTimestamp: 1234567890,
				MetaData:         "metadatastr",
			}}}.String(),
		`{"resp":{"status":"ok","ts":123456790,"license":{"object":"objectstr","valid_to_timestamp":1234567890,"metadata":"metadatastr"}}}`,
	)
}
