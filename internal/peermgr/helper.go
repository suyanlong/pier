package peermgr

import (
	"github.com/link33/sidecar/model/pb"
)

const (
	V1 = "1.0"
)

func Message(typ pb.Message_Type, ok bool, data []byte) *pb.Message {
	return &pb.Message{
		Type:    typ,
		Version: V1,
		Payload: &pb.Pack{
			Ok:   ok,
			Data: data,
		},
	}
}
