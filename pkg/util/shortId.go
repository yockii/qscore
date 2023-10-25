package util

import "github.com/rs/xid"

func GenerateXid() string {
	uid := xid.New()
	return uid.String()
}
