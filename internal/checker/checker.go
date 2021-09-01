package checker

import "github.com/link33/sidecar/model/pb"

type MockChecker struct {
}

func (ck *MockChecker) Check(ibtp *pb.IBTP) error {
	return nil
}
