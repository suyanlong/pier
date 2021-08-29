package model

import "github.com/meshplus/pier/model/pb"

type PluginResponse struct {
	Status  bool
	Message string
	Result  *pb.IBTP
}
