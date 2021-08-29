package checker

import "github.com/meshplus/pier/model/pb"

type MockChecker struct {
}

func (ck *MockChecker) Check(ibtp *pb.IBTP) error {
	return nil
}
