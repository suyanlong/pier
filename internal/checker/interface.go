package checker

import "github.com/link33/sidecar/model/pb"

type Checker interface {
	Check(ibtp *pb.IBTP) error
}
