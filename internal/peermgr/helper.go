package peermgr

import (
	"github.com/link33/sidercar/model/pb"
)

const (
	V1 = "1.0"
)

func Message(typ pb.Message_Type, ok bool, data []byte) *pb.Message {
	return &pb.Message{
		Type:    typ,
		Version: V1,
		Payload: &pb.Payload{
			Ok:   ok,
			Data: data,
		},
	}
}
