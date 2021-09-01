package checker

import "github.com/link33/sidercar/model/pb"

type Checker interface {
	Check(ibtp *pb.IBTP) error
}
