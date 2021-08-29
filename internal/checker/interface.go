package checker

import "github.com/meshplus/pier/model/pb"

type Checker interface {
	Check(ibtp *pb.IBTP) error
}
